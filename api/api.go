package api

import (
	"github.com/gofiber/fiber/v2"
	admissionvalidator "github.com/gomgomshrimp/validating-webhook-server/validator"
	"github.com/sirupsen/logrus"
	admissionv1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Api struct {
	Logger             *logrus.Logger
	AdmissionValidator *admissionvalidator.AdmissionValidator
}

func NewApi(logger *logrus.Logger) *Api {
	return &Api{
		Logger:             logger,
		AdmissionValidator: admissionvalidator.NewAdmissionValidator(logger),
	}
}

func (a *Api) Validate(c *fiber.Ctx) error {
	// Create AdmissionReview for response
	admissionResponse := admissionv1.AdmissionReview{
		TypeMeta: metav1.TypeMeta{
			Kind:       "AdmissionReview",
			APIVersion: "admission.k8s.io/v1",
		},
		Response: &admissionv1.AdmissionResponse{
			Allowed: false,
			Result:  &metav1.Status{},
		},
	}
	// Parsing request body
	admissionRequest := &admissionv1.AdmissionReview{}
	if err := c.BodyParser(admissionRequest); err != nil {
		admissionResponse.Response.Result.Message = "AdmissionReview parsing error"
		return c.JSON(admissionResponse)
	}

	switch admissionRequest.Request.RequestKind.Kind {
	case "Deployment":
		if err := a.AdmissionValidator.ValidateDeployment(admissionRequest.Request.Object.Raw); err != nil {
			admissionResponse.Response.Result.Message = err.Error()
		} else {
			admissionResponse.Response.Allowed = true
		}
	case "Pod":
		if err := a.AdmissionValidator.ValidatePod(admissionRequest.Request.Object.Raw); err != nil {
			admissionResponse.Response.Result.Message = err.Error()
		} else {
			admissionResponse.Response.Allowed = true
		}
	default:
		admissionResponse.Response.Allowed = true
	}

	admissionResponse.Response.UID = admissionRequest.Request.UID
	a.Logger.Info(admissionResponse)

	return c.JSON(admissionResponse)
}
