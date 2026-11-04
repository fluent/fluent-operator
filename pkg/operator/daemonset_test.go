package operator

import (
	"testing"

	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	fluentbitv1alpha2 "github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2"
)

func TestMakeDaemonSetAnnotations(t *testing.T) {
	t.Run("pod annotations are also set on the daemonset", func(t *testing.T) {
		g := NewGomegaWithT(t)
		fb := fluentbitv1alpha2.FluentBit{
			ObjectMeta: metav1.ObjectMeta{
				Name:        "fluent-bit",
				Namespace:   "logging",
				Annotations: map[string]string{"owner": "resource", "shared": "resource"},
			},
			Spec: fluentbitv1alpha2.FluentBitSpec{
				Annotations: map[string]string{"team": "obs", "shared": "spec"},
			},
		}

		ds := MakeDaemonSet(fb, "/var/log")

		g.Expect(ds.Annotations).To(Equal(map[string]string{"owner": "resource", "team": "obs", "shared": "spec"}))
		g.Expect(ds.Spec.Template.Annotations).To(Equal(map[string]string{"team": "obs", "shared": "spec"}))
		// the maps of the FluentBit resource are not modified
		g.Expect(fb.Annotations).To(Equal(map[string]string{"owner": "resource", "shared": "resource"}))
	})

	t.Run("no annotations", func(t *testing.T) {
		g := NewGomegaWithT(t)
		ds := MakeDaemonSet(fluentbitv1alpha2.FluentBit{ObjectMeta: metav1.ObjectMeta{Name: "fluent-bit"}}, "/var/log")

		g.Expect(ds.Annotations).To(BeNil())
	})
}
