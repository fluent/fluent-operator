package input

import (
	"testing"

	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins/params"
	. "github.com/onsi/gomega"
)

func TestTail_Name(t *testing.T) {
	g := NewGomegaWithT(t)
	tail := Tail{}
	g.Expect(tail.Name()).To(Equal("tail"))
}

func TestTail_Params_StoragePath(t *testing.T) {
	g := NewGomegaWithT(t)
	sl := plugins.NewSecretLoader(nil, "test namespace")

	tail := Tail{
		StorageType: "filesystem",
		StoragePath: "/var/log/flb-storage/",
	}

	expected := params.NewKVs()
	expected.Insert("storage.type", "filesystem")
	expected.Insert("storage.path", "/var/log/flb-storage/")

	kvs, err := tail.Params(sl)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(kvs).To(Equal(expected))
}
