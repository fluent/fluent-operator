package operator

import (
	"testing"

	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	fluentdv1alpha1 "github.com/fluent/fluent-operator/v3/apis/fluentd/v1alpha1"
)

func TestFluentdWorkloadAnnotations(t *testing.T) {
	t.Run("pod annotations are also set on the statefulset and daemonset", func(t *testing.T) {
		g := NewGomegaWithT(t)
		fd := fluentdv1alpha1.Fluentd{
			ObjectMeta: metav1.ObjectMeta{
				Name:        "fluentd",
				Namespace:   "logging",
				Annotations: map[string]string{"owner": "resource", "shared": "resource"},
			},
			Spec: fluentdv1alpha1.FluentdSpec{
				Replicas:    new(int32(1)),
				Annotations: map[string]string{"team": "obs", "shared": "spec"},
			},
		}
		want := map[string]string{"owner": "resource", "team": "obs", "shared": "spec"}
		pod := map[string]string{"team": "obs", "shared": "spec"}

		sts := MakeStatefulSet(fd)
		g.Expect(sts.Annotations).To(Equal(want))
		g.Expect(sts.Spec.Template.Annotations).To(Equal(pod))

		ds := MakeFluentdDaemonSet(fd)
		g.Expect(ds.Annotations).To(Equal(want))
		g.Expect(ds.Spec.Template.Annotations).To(Equal(pod))

		// the maps of the Fluentd resource are not modified
		g.Expect(fd.Annotations).To(Equal(map[string]string{"owner": "resource", "shared": "resource"}))
	})

	t.Run("no annotations", func(t *testing.T) {
		g := NewGomegaWithT(t)
		fd := fluentdv1alpha1.Fluentd{
			ObjectMeta: metav1.ObjectMeta{Name: "fluentd", Namespace: "logging"},
			Spec:       fluentdv1alpha1.FluentdSpec{Replicas: new(int32(1))},
		}

		g.Expect(MakeStatefulSet(fd).Annotations).To(BeNil())
		g.Expect(MakeFluentdDaemonSet(fd).Annotations).To(BeNil())
	})
}
