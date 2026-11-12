package filter

import (
	"testing"

	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins/params"
	. "github.com/onsi/gomega"
)

func TestFilter_Sysinfo_Params(t *testing.T) {
	g := NewGomegaWithT(t)

	sl := plugins.NewSecretLoader(nil, "test namespace")

	si := Sysinfo{
		FluentBitVersionKey: "flb_version",
		OSNameKey:           "os_name",
		HostnameKey:         "hostname",
		OSVersionKey:        "os_version",
		KernelVersionKey:    "kernel_version",
	}

	expected := params.NewKVs()
	expected.Insert("Fluentbit_Version_Key", "flb_version")
	expected.Insert("Os_Name_Key", "os_name")
	expected.Insert("Hostname_Key", "hostname")
	expected.Insert("Os_Version_Key", "os_version")
	expected.Insert("Kernel_Version_Key", "kernel_version")

	kvs, err := si.Params(sl)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(kvs).To(Equal(expected))
}

func TestFilter_Sysinfo_ParamsOnlySetKeys(t *testing.T) {
	g := NewGomegaWithT(t)

	sl := plugins.NewSecretLoader(nil, "test namespace")

	kvs, err := (&Sysinfo{HostnameKey: "host"}).Params(sl)
	g.Expect(err).NotTo(HaveOccurred())
	expected := params.NewKVs()
	expected.Insert("Hostname_Key", "host")
	g.Expect(kvs).To(Equal(expected))

	kvs, err = (&Sysinfo{}).Params(sl)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(kvs).To(Equal(params.NewKVs()))
}
