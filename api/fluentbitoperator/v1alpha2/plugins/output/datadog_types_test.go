package output

import (
	"testing"

	. "github.com/onsi/gomega"
	"kubesphere.io/fluentbit-operator/api/fluentbitoperator/v1alpha2/plugins"
	"kubesphere.io/fluentbit-operator/api/fluentbitoperator/v1alpha2/plugins/params"
)

func TestOutput_DataDog_Params(t *testing.T) {
	g := NewGomegaWithT(t)

	sl := plugins.NewSecretLoader(nil, "test namespace", nil)

	dd := DataDog{
		Host:          "http-intake.logs.datadoghq.com",
		APIKey:        "1234apikey",
		TLS:           ptrBool(true),
		Compress:      "gzip",
		Service:       "service_name",
		Source:        "app_name",
		Tags:          "foo:bar",
		MessageKey:    "message",
		JSONDateKey:   "timestamp",
		IncludeTagKey: ptrBool(true),
		TagKey:        "tagkey",
	}

	expected := params.NewKVs()
	expected.Insert("Host", "http-intake.logs.datadoghq.com")
	expected.Insert("TLS", "true")
	expected.Insert("compress", "gzip")
	expected.Insert("apikey", "1234apikey")
	expected.Insert("json_date_key", "timestamp")
	expected.Insert("include_tag_key", "true")
	expected.Insert("tag_key", "tagkey")
	expected.Insert("dd_service", "service_name")
	expected.Insert("dd_source", "app_name")
	expected.Insert("dd_tags", "foo:bar")
	expected.Insert("dd_message_key", "message")

	kvs, err := dd.Params(sl)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(kvs).To(Equal(expected))

}

func ptrBool(v bool) *bool {
	return &v
}
