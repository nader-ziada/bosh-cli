package fakes

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	biinstallmanifest "github.com/cloudfoundry/bosh-cli/v7/installation/manifest"
	birelsetmanifest "github.com/cloudfoundry/bosh-cli/v7/release/set/manifest"
)

type FakeValidator struct {
	ValidateInputs  []ValidateInput
	validateOutputs []ValidateOutput
}

func NewFakeValidator() *FakeValidator {
	return &FakeValidator{
		ValidateInputs:  []ValidateInput{},
		validateOutputs: []ValidateOutput{},
	}
}

type ValidateInput struct {
	InstallationManifest biinstallmanifest.Manifest
	ReleaseSetManifest   birelsetmanifest.Manifest
}

type ValidateOutput struct {
	Err error
}

func (v *FakeValidator) Validate(installationManifest biinstallmanifest.Manifest, releaseSetManifest birelsetmanifest.Manifest) error {
	v.ValidateInputs = append(v.ValidateInputs, ValidateInput{
		InstallationManifest: installationManifest,
		ReleaseSetManifest:   releaseSetManifest,
	})

	if len(v.validateOutputs) == 0 {
		return bosherr.Errorf("Unexpected FakeValidator.Validate(manifest) called with installation manifest: %#v, release set manifest: %#v", installationManifest, releaseSetManifest)
	}
	validateOutput := v.validateOutputs[0]
	v.validateOutputs = v.validateOutputs[1:]
	return validateOutput.Err
}

func (v *FakeValidator) SetValidateBehavior(outputs []ValidateOutput) {
	v.validateOutputs = outputs
}
