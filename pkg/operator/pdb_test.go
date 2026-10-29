package operator

import (
	"testing"

	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	fluentdv1alpha1 "github.com/fluent/fluent-operator/v3/apis/fluentd/v1alpha1"
)

func newCollectorFluentd(pdb *fluentdv1alpha1.PodDisruptionBudget) fluentdv1alpha1.Fluentd {
	replicas := int32(3)
	return fluentdv1alpha1.Fluentd{
		ObjectMeta: metav1.ObjectMeta{Name: "fluentd", Namespace: "logging", Labels: map[string]string{"team": "obs"}},
		Spec: fluentdv1alpha1.FluentdSpec{
			Replicas:            &replicas,
			PodDisruptionBudget: pdb,
		},
	}
}

func TestMakeFluentdPodDisruptionBudget(t *testing.T) {
	t.Run("no budget requested", func(t *testing.T) {
		g := NewGomegaWithT(t)
		g.Expect(MakeFluentdPodDisruptionBudget(newCollectorFluentd(nil))).To(BeNil())
	})

	t.Run("minAvailable selects the statefulset pods", func(t *testing.T) {
		g := NewGomegaWithT(t)
		minAvailable := intstr.FromInt32(2)
		fd := newCollectorFluentd(&fluentdv1alpha1.PodDisruptionBudget{MinAvailable: &minAvailable})

		pdb := MakeFluentdPodDisruptionBudget(fd)

		g.Expect(pdb).NotTo(BeNil())
		g.Expect(pdb.Name).To(Equal("fluentd"))
		g.Expect(pdb.Namespace).To(Equal("logging"))
		g.Expect(pdb.Spec.MinAvailable).To(Equal(&minAvailable))
		g.Expect(pdb.Spec.MaxUnavailable).To(BeNil())
		g.Expect(pdb.Spec.Selector).To(Equal(MakeStatefulSet(fd).Spec.Selector))
	})

	t.Run("maxUnavailable as a percentage", func(t *testing.T) {
		g := NewGomegaWithT(t)
		maxUnavailable := intstr.FromString("50%")
		fd := newCollectorFluentd(&fluentdv1alpha1.PodDisruptionBudget{MaxUnavailable: &maxUnavailable})

		pdb := MakeFluentdPodDisruptionBudget(fd)

		g.Expect(pdb.Spec.MinAvailable).To(BeNil())
		g.Expect(pdb.Spec.MaxUnavailable).To(Equal(&maxUnavailable))
	})
}
