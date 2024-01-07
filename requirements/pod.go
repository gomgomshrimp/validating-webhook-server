package requirements

import (
	corev1 "k8s.io/api/core/v1"
)

type PodRequirement struct {
	Spec PodSpec `json:"spec" validate:"required"`
}

type PodSpec struct {
	Containers   []Container       `json:"containers" validate:"required,gt=0,dive,required"`
	NodeSelector map[string]string `json:"nodeSelector" validate:"required"`
	//TopologySpreadConstraints []corev1.TopologySpreadConstraint `json:"topologySpreadConstraints" validate:"required"`
}

type Container struct {
	Name      string `json:"name"`
	Image     string `json:"image"`
	Resources struct {
		Limits   corev1.ResourceList `json:"limits" validate:"required"`
		Requests corev1.ResourceList `json:"requests" validate:"required"`
	} `json:"resources" validate:"required"`
	LivenessProbe  *corev1.Probe `json:"livenessProbe" validate:"required"`
	ReadinessProbe *corev1.Probe `json:"readinessProbe" validate:"required"`
}
