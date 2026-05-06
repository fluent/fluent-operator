package input

import (
	"testing"

	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins/params"
	. "github.com/onsi/gomega"
)

func TestKubernetesEvents_Name(t *testing.T) {
	g := NewGomegaWithT(t)
	ke := KubernetesEvents{}
	g.Expect(ke.Name()).To(Equal("kubernetes_events"))
}

func TestKubernetesEvents_Params_Storage(t *testing.T) {
	g := NewGomegaWithT(t)
	sl := plugins.NewSecretLoader(nil, "test namespace")

	ke := KubernetesEvents{
		StorageType:            "filesystem",
		PauseOnChunksOverlimit: "on",
	}

	expected := params.NewKVs()
	expected.Insert("storage.type", "filesystem")
	expected.Insert("storage.pause_on_chunks_overlimit", "on")

	kvs, err := ke.Params(sl)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(kvs).To(Equal(expected))
}
