package fakes

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

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
	Manifest birelsetmanifest.Manifest
}

type ValidateOutput struct {
	Err error
}

func (v *FakeValidator) Validate(manifest birelsetmanifest.Manifest) error {
	v.ValidateInputs = append(v.ValidateInputs, ValidateInput{
		Manifest: manifest,
	})

	if len(v.validateOutputs) == 0 {
		return bosherr.Errorf("Unexpected FakeValidator.Validate(manifest) called with manifest: %#v", manifest)
	}
	validateOutput := v.validateOutputs[0]
	v.validateOutputs = v.validateOutputs[1:]
	return validateOutput.Err
}

func (v *FakeValidator) SetValidateBehavior(outputs []ValidateOutput) {
	v.validateOutputs = outputs
}
