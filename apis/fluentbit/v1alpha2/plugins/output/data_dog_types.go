package output

import (
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins/params"
)

// +kubebuilder:object:generate:=true

// DataDog output plugin allows you to ingest your logs into Datadog. <br />
// **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/outputs/datadog**
type DataDog struct {
	// Host is the Datadog server where you are sending your logs.
	Host string `json:"host,omitempty"`
	// TLS controls whether to use end-to-end security communications security protocol.
	// Datadog recommends setting this to on.
	TLS *bool `json:"tls,omitempty"`
	// Compress  the payload in GZIP format.
	// Datadog supports and recommends setting this to gzip.
	Compress string `json:"compress,omitempty"`
	// Your Datadog API key.
	APIKey *plugins.Secret `json:"apikey,omitempty"`
	// Specify an HTTP Proxy.
	Proxy string `json:"proxy,omitempty"`
	// To activate the remapping, specify configuration flag provider.
	Provider string `json:"provider,omitempty"`
	// Date key name for output.
	JSONDateKey string `json:"json_date_key,omitempty"`
	// If enabled, a tag is appended to output. The key name is used tag_key property.
	IncludeTagKey *bool `json:"include_tag_key,omitempty"`
	// The key name of tag. If include_tag_key is false, This property is ignored.
	TagKey string `json:"tag_key,omitempty"`
	// The human readable name for your service generating the logs.
	Service string `json:"dd_service,omitempty"`
	// A human readable name for the underlying technology of your service.
	Source string `json:"dd_source,omitempty"`
	// The tags you want to assign to your logs in Datadog.
	Tags string `json:"dd_tags,omitempty"`
	// By default, the plugin searches for the key 'log' and remap the value to the key 'message'. If the property is set, the plugin will search the property name key.
	MessageKey string `json:"dd_message_key,omitempty"`

	// *plugins.HTTP `json:"tls,omitempty"`
}

func (*DataDog) Name() string {
	return "datadog"
}

// implement Section() method
func (s *DataDog) Params(sl plugins.SecretLoader) (*params.KVs, error) {
	kvs := params.NewKVs()

	plugins.InsertKVString(kvs, "Host", s.Host)
	plugins.InsertKVField(kvs, "TLS", s.TLS)
	plugins.InsertKVString(kvs, "compress", s.Compress)
	plugins.InsertKVString(kvs, "json_date_key", s.JSONDateKey)
	plugins.InsertKVField(kvs, "include_tag_key", s.IncludeTagKey)
	plugins.InsertKVString(kvs, "tag_key", s.TagKey)
	plugins.InsertKVString(kvs, "dd_service", s.Service)
	plugins.InsertKVString(kvs, "dd_source", s.Source)
	plugins.InsertKVString(kvs, "dd_tags", s.Tags)
	plugins.InsertKVString(kvs, "dd_message_key", s.MessageKey)

	if err := plugins.InsertKVSecret(kvs, "apikey", s.APIKey, sl); err != nil {
		return nil, err
	}
	plugins.InsertKVString(kvs, "proxy", s.Proxy)
	plugins.InsertKVString(kvs, "provider", s.Provider)

	return kvs, nil
}
