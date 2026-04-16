package input

import (
	"testing"

	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins/params"
	. "github.com/onsi/gomega"
)

func TestSystemd_Name(t *testing.T) {
	g := NewGomegaWithT(t)
	systemd := Systemd{}
	g.Expect(systemd.Name()).To(Equal("systemd"))
}

func TestSystemd_Params_StoragePath(t *testing.T) {
	g := NewGomegaWithT(t)
	sl := plugins.NewSecretLoader(nil, "test namespace")

	systemd := Systemd{
		StorageType: "filesystem",
		StoragePath: "/var/log/flb-storage/",
	}

	expected := params.NewKVs()
	expected.Insert("storage.type", "filesystem")
	expected.Insert("storage.path", "/var/log/flb-storage/")

	kvs, err := systemd.Params(sl)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(kvs).To(Equal(expected))
}
