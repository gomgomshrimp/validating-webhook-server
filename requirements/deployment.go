package requirements

import corev1 "k8s.io/api/core/v1"

type DeploymentRequirement struct {
	Spec DeploymentSpec `json:"spec" validate:"required"`
}

type DeploymentSpec struct {
	Template struct {
		Spec struct {
			PodSpec
			TopologySpreadConstraints []corev1.TopologySpreadConstraint `json:"topologySpreadConstraints" validate:"required,gt=0,dive,required"`
		} `json:"spec" validate:"required"`
	} `json:"template" validate:"required"`
}
