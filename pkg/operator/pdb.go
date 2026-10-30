package operator

import (
	policyv1 "k8s.io/api/policy/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	fluentdv1alpha1 "github.com/fluent/fluent-operator/v3/apis/fluentd/v1alpha1"
)

// MakeFluentdPodDisruptionBudget returns the PodDisruptionBudget for the Fluentd StatefulSet, selecting its pods with
// the StatefulSet's own selector. It returns nil when the Fluentd spec does not ask for one.
func MakeFluentdPodDisruptionBudget(fd fluentdv1alpha1.Fluentd) *policyv1.PodDisruptionBudget {
	if fd.Spec.PodDisruptionBudget == nil {
		return nil
	}

	sts := MakeStatefulSet(fd)

	return &policyv1.PodDisruptionBudget{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fd.Name,
			Namespace: fd.Namespace,
			Labels:    sts.Labels,
		},
		Spec: policyv1.PodDisruptionBudgetSpec{
			MinAvailable:   fd.Spec.PodDisruptionBudget.MinAvailable,
			MaxUnavailable: fd.Spec.PodDisruptionBudget.MaxUnavailable,
			Selector:       sts.Spec.Selector,
		},
	}
}
