package output

import (
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins/params"
)

// +kubebuilder:object:generate:=true

// Azure Logs Ingestion is the Azure Logs Ingestion output plugin, allows you to ingest your records into Azure Monitor tables through a data collection endpoint and rule. <br />
// **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/outputs/azure_logs_ingestion**
type AzureLogsIngestion struct {
	// Azure Active Directory tenant ID
	TenantID string `json:"tenantID"`
	// Azure Active Directory application (client) ID
	ClientID string `json:"clientID"`
	// Azure Active Directory application client secret
	ClientSecret *plugins.Secret `json:"clientSecret"`
	// Data Collection Endpoint (DCE) URL
	DceURL string `json:"dceURL"`
	// Data Collection Rule (DCR) immutable ID
	DcrID string `json:"dcrID"`
	// Name of the destination table the records are ingested into
	TableName string `json:"tableName"`
	// Specify the name of the key where the timestamp is stored.
	TimeKey string `json:"timeKey,omitempty"`
	// If set, overrides the timeKey value with the `time-generated-field` HTTP header value.
	TimeGenerated *bool `json:"timeGenerated,omitempty"`
	// Enable gzip compression of the request body.
	Compress *bool `json:"compress,omitempty"`
}

// Name implement Section() method
func (*AzureLogsIngestion) Name() string {
	return "azure_logs_ingestion"
}

// Params implement Section() method
func (o *AzureLogsIngestion) Params(sl plugins.SecretLoader) (*params.KVs, error) {
	kvs := params.NewKVs()
	plugins.InsertKVString(kvs, "tenant_id", o.TenantID)
	plugins.InsertKVString(kvs, "client_id", o.ClientID)
	if o.ClientSecret != nil {
		u, err := sl.LoadSecret(*o.ClientSecret)
		if err != nil {
			return nil, err
		}
		kvs.Insert("client_secret", u)
	}
	plugins.InsertKVString(kvs, "dce_url", o.DceURL)
	plugins.InsertKVString(kvs, "dcr_id", o.DcrID)
	plugins.InsertKVString(kvs, "table_name", o.TableName)
	plugins.InsertKVString(kvs, "time_key", o.TimeKey)
	plugins.InsertKVField(kvs, "time_generated", o.TimeGenerated)
	plugins.InsertKVField(kvs, "compress", o.Compress)

	return kvs, nil
}
