package validator

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/gomgomshrimp/validating-webhook-server/spec"
	"github.com/sirupsen/logrus"
)

type AdmissionValidator struct {
	Logger    *logrus.Logger
	Validator *validator.Validate
}

func NewAdmissionValidator(logger *logrus.Logger) *AdmissionValidator {
	return &AdmissionValidator{
		Logger:    logger,
		Validator: validator.New(),
	}
}

func (v *AdmissionValidator) ValidateDeployment(rawObject []byte) error {
	deploymentRequirement := &spec.DeploymentRequirement{}
	if err := json.Unmarshal(rawObject, deploymentRequirement); err != nil {
		v.Logger.Fatal(err)
		return err
	}
	if err := v.validateWithRequirements(deploymentRequirement); err != nil {
		v.Logger.Info(err)
		return err
	}

	return nil
}

func (v *AdmissionValidator) ValidatePod(rawObject []byte) error {
	deploymentRequirement := &spec.PodRequirement{}
	if err := json.Unmarshal(rawObject, deploymentRequirement); err != nil {
		v.Logger.Fatal(err)
		return err
	}
	if err := v.validateWithRequirements(deploymentRequirement); err != nil {
		v.Logger.Info(err)
		return err
	}

	return nil
}

func (v *AdmissionValidator) validateWithRequirements(input interface{}) error {
	return v.Validator.Struct(input)
}
