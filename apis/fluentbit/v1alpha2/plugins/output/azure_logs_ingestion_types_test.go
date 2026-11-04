package output

import (
	"testing"

	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins/params"
	. "github.com/onsi/gomega"
)

func TestOutput_AzureLogsIngestion_Params(t *testing.T) {
	g := NewGomegaWithT(t)

	sl := plugins.NewSecretLoader(nil, "test namespace")

	o := AzureLogsIngestion{
		TenantID:      "tenant",
		ClientID:      "client",
		DceURL:        "https://dce.example.ingest.monitor.azure.com",
		DcrID:         "dcr-0123",
		TableName:     "Custom-Logs_CL",
		TimeKey:       "TimeGenerated",
		TimeGenerated: new(true),
		Compress:      new(true),
	}

	expected := params.NewKVs()
	expected.Insert("tenant_id", "tenant")
	expected.Insert("client_id", "client")
	expected.Insert("dce_url", "https://dce.example.ingest.monitor.azure.com")
	expected.Insert("dcr_id", "dcr-0123")
	expected.Insert("table_name", "Custom-Logs_CL")
	expected.Insert("time_key", "TimeGenerated")
	expected.Insert("time_generated", "true")
	expected.Insert("compress", "true")

	kvs, err := o.Params(sl)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(kvs).To(Equal(expected))
	g.Expect(o.Name()).To(Equal("azure_logs_ingestion"))
}
