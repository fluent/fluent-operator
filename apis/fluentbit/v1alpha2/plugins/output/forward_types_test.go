package output

import (
	"testing"

	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins/params"
	. "github.com/onsi/gomega"
)

func TestOutput_Forward_Params(t *testing.T) {
	g := NewGomegaWithT(t)

	sl := plugins.NewSecretLoader(nil, "test namespace")

	fwd := Forward{
		Host:     "fluentd.logging",
		Port:     new(int32(24224)),
		Compress: "gzip",
	}

	expected := params.NewKVs()
	expected.Insert("Host", "fluentd.logging")
	expected.Insert("Compress", "gzip")
	expected.Insert("Port", "24224")

	kvs, err := fwd.Params(sl)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(kvs).To(Equal(expected))
}

func TestOutput_Forward_Params_NoCompress(t *testing.T) {
	g := NewGomegaWithT(t)

	sl := plugins.NewSecretLoader(nil, "test namespace")

	fwd := Forward{Host: "fluentd.logging"}

	expected := params.NewKVs()
	expected.Insert("Host", "fluentd.logging")

	kvs, err := fwd.Params(sl)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(kvs).To(Equal(expected))
}
