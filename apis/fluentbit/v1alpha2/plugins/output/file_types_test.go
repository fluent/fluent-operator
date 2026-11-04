package output

import (
	"testing"

	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins/params"
	. "github.com/onsi/gomega"
)

func TestOutput_File_Params(t *testing.T) {
	g := NewWithT(t)

	sl := plugins.NewSecretLoader(nil, "test namespace")
	enabled := true

	file := File{
		Path:  "/var/log/out/app",
		File:  "records.log",
		Mkdir: &enabled,
	}

	expected := params.NewKVs()
	expected.Insert("Path", "/var/log/out/app")
	expected.Insert("File", "records.log")
	expected.Insert("Mkdir", "true")

	kvs, err := file.Params(sl)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(kvs).To(Equal(expected))
}

func TestOutput_File_ParamsMkdirUnsetAndFalse(t *testing.T) {
	g := NewWithT(t)

	sl := plugins.NewSecretLoader(nil, "test namespace")
	disabled := false

	kvs, err := (&File{Path: "/out"}).Params(sl)
	g.Expect(err).NotTo(HaveOccurred())
	expected := params.NewKVs()
	expected.Insert("Path", "/out")
	g.Expect(kvs).To(Equal(expected))

	kvs, err = (&File{Path: "/out", Mkdir: &disabled}).Params(sl)
	g.Expect(err).NotTo(HaveOccurred())
	expected.Insert("Mkdir", "false")
	g.Expect(kvs).To(Equal(expected))
}
