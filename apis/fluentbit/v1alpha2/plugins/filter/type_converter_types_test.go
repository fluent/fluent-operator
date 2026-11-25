package filter

import (
	"testing"

	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins/params"
	. "github.com/onsi/gomega"
)

func TestFilter_TypeConverter_Params(t *testing.T) {
	g := NewGomegaWithT(t)

	sl := plugins.NewSecretLoader(nil, "test namespace")

	tc := TypeConverter{
		IntKeys: []TypeConversion{
			{FromKey: "status", ToKey: "status_str", Type: "string"},
		},
		UintKeys: []TypeConversion{
			{FromKey: "Mem.total", ToKey: "Mem.total_str", Type: "string"},
			{FromKey: "Mem.used", ToKey: "Mem.used_str", Type: "string"},
		},
		FloatKeys: []TypeConversion{
			{FromKey: "latency", ToKey: "latency_int", Type: "int"},
		},
		StrKeys: []TypeConversion{
			{FromKey: "bytes", ToKey: "bytes_uint", Type: "uint"},
			{FromKey: "ratio", ToKey: "ratio_float", Type: "float"},
		},
	}

	expected := params.NewKVs()
	expected.Insert("Int_Key", "status status_str string")
	expected.Insert("Uint_Key", "Mem.total Mem.total_str string")
	expected.Insert("Uint_Key", "Mem.used Mem.used_str string")
	expected.Insert("Float_Key", "latency latency_int int")
	expected.Insert("Str_Key", "bytes bytes_uint uint")
	expected.Insert("Str_Key", "ratio ratio_float float")

	kvs, err := tc.Params(sl)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(kvs).To(Equal(expected))
}

func TestFilter_TypeConverter_ParamsEmpty(t *testing.T) {
	g := NewGomegaWithT(t)

	sl := plugins.NewSecretLoader(nil, "test namespace")

	kvs, err := (&TypeConverter{}).Params(sl)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(kvs).To(Equal(params.NewKVs()))
}
