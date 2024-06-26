package client

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/corehandlers"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/cloudfoundry/bosh-s3cli/config"
)

const (
	signatureVersion = "2"                    //nolint:unused
	signatureMethod  = "HmacSHA1"             //nolint:unused
	timeFormat       = "2006-01-02T15:04:05Z" //nolint:unused
)

type signer struct {
	// Values that must be populated from the request
	Request     *http.Request
	Time        time.Time
	Credentials *credentials.Credentials
	Debug       aws.LogLevelType
	Logger      aws.Logger

	Query        url.Values
	stringToSign string
	signature    string
}

var s3ParamsToSign = map[string]bool{
	"acl":                          true,
	"location":                     true,
	"logging":                      true,
	"notification":                 true,
	"partNumber":                   true,
	"policy":                       true,
	"requestPayment":               true,
	"torrent":                      true,
	"uploadId":                     true,
	"uploads":                      true,
	"versionId":                    true,
	"versioning":                   true,
	"versions":                     true,
	"response-content-type":        true,
	"response-content-language":    true,
	"response-expires":             true,
	"response-cache-control":       true,
	"response-content-disposition": true,
	"response-content-encoding":    true,
	"website":                      true,
	"delete":                       true,
}

func setv2Handlers(svc *s3.S3) {
	svc.Handlers.Build.PushBack(func(r *request.Request) {
		parsedURL, err := url.Parse(r.HTTPRequest.URL.String())
		if err != nil {
			log.Fatal("Failed to parse URL", err)
		}
		r.HTTPRequest.URL.Opaque = parsedURL.Path
	})

	svc.Handlers.Sign.Clear()
	svc.Handlers.Sign.PushBack(Sign)
	svc.Handlers.Sign.PushBackNamed(corehandlers.BuildContentLengthHandler)
}

// Sign requests with signature version 2.
//
// Will sign the requests with the service config's Credentials object
// Signing is skipped if the credentials is the credentials.AnonymousCredentials
// object.
func Sign(req *request.Request) {
	// If the request does not need to be signed ignore the signing of the
	// request if the AnonymousCredentials object is used.
	if req.Config.Credentials == credentials.AnonymousCredentials {
		return
	}

	v2 := signer{
		Request:     req.HTTPRequest,
		Time:        req.Time,
		Credentials: req.Config.Credentials,
		Debug:       req.Config.LogLevel.Value(),
		Logger:      req.Config.Logger,
	}

	req.Error = v2.Sign()

	if req.Error != nil {
		return
	}
}

func (v2 *signer) Sign() error {
	credValue, err := v2.Credentials.Get()
	if err != nil {
		return err
	}
	accessKey := credValue.AccessKeyID
	var (
		md5, ctype, date, xamz string
		xamzDate               bool
		sarray                 []string
	)

	headers := v2.Request.Header
	params := v2.Request.URL.Query()
	parsedURL, err := url.Parse(v2.Request.URL.String())
	if err != nil {
		return err
	}
	host, canonicalPath := parsedURL.Host, parsedURL.Path

	// Host can not be parsed successfully results from v2.Request.URL has a non-empty Opaque.
	if len(host) < 1 {
		host = v2.Request.URL.Host
	}

	v2.Request.Header["Host"] = []string{host}
	v2.Request.Header["x-amz-date"] = []string{v2.Time.In(time.UTC).Format(time.RFC1123)}

	// Alibaba Cloud OSS date's formate must be http.TimeFormat
	// Alibaba Cloud OSS uses virtual hosted and the URL Host's format is <bucket-name>.host or <bucket-name>.host:port
	if config.Provider(strings.Join(strings.Split(strings.Split(host, ":")[0], ".")[1:], ".")) == "alicloud" {
		v2.Request.Header["x-amz-date"] = []string{v2.Time.In(time.UTC).Format(http.TimeFormat)}
		canonicalPath = fmt.Sprintf("/%s%s", strings.Split(host, ".")[0], canonicalPath)
	}

	for k, v := range headers {
		k = strings.ToLower(k)
		switch k {
		case "content-md5":
			md5 = v[0]
		case "content-type":
			ctype = v[0]
		case "date":
			if !xamzDate {
				date = v[0]
			}
		default:
			if strings.HasPrefix(k, "x-amz-") {
				vall := strings.Join(v, ",")
				sarray = append(sarray, k+":"+vall)
				if k == "x-amz-date" {
					xamzDate = true
					date = ""
				}
			}
		}
	}
	if len(sarray) > 0 {
		sort.StringSlice(sarray).Sort()
		xamz = strings.Join(sarray, "\n") + "\n"
	}

	expires := false
	if v, ok := params["Expires"]; ok {
		expires = true
		date = v[0]
		params["AWSAccessKeyId"] = []string{accessKey}
	}

	sarray = sarray[0:0]
	for k, v := range params {
		if s3ParamsToSign[k] {
			for _, vi := range v {
				if vi == "" {
					sarray = append(sarray, k)
				} else {
					sarray = append(sarray, k+"="+vi)
				}
			}
		}
	}
	if len(sarray) > 0 {
		sort.StringSlice(sarray).Sort()
		canonicalPath = canonicalPath + "?" + strings.Join(sarray, "&")
	}

	v2.stringToSign = strings.Join([]string{
		v2.Request.Method,
		md5,
		ctype,
		date,
		xamz + canonicalPath,
	}, "\n")
	hash := hmac.New(sha1.New, []byte(credValue.SecretAccessKey))
	hash.Write([]byte(v2.stringToSign))
	v2.signature = base64.StdEncoding.EncodeToString(hash.Sum(nil))

	if expires {
		params["Signature"] = []string{string(v2.signature)}
	} else {
		headers["Authorization"] = []string{"AWS " + accessKey + ":" + string(v2.signature)}
	}

	if v2.Debug.Matches(aws.LogDebugWithSigning) {
		v2.logSigningInfo()
	}
	return nil
}

const logSignInfoMsg = `DEBUG: Request Signature:
---[ STRING TO SIGN ]--------------------------------
%s
---[ SIGNATURE ]-------------------------------------
%s
-----------------------------------------------------`

func (v2 *signer) logSigningInfo() {
	msg := fmt.Sprintf(logSignInfoMsg, v2.stringToSign, v2.signature)
	v2.Logger.Log(msg)
}
