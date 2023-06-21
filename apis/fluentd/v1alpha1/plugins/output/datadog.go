package output

import "github.com/fluent/fluent-operator/v2/apis/fluentd/v1alpha1/plugins"

// Datadog defines the parameters for out_datadog plugin
type Datadog struct {
	// This parameter is required in order to authenticate your fluent agent.
	ApiKey *plugins.Secret `json:"apiKey,omitempty"`
	// Event format, if true, the event is sent in json format. Othwerwise, in plain text.
	UseJson *bool `json:"useJson,omitempty"`
	// Automatically include the Fluentd tag in the record.
	IncludeTagKey *bool `json:"includeTagKey,omitempty"`
	// Where to store the Fluentd tag.
	TagKey *string `json:"tagKey,omitempty"`
	// Name of the attribute which will contain timestamp of the log event. If nil, timestamp attribute is not added.
	TimestampKey *string `json:"timestampKey,omitempty"`
	// If true, the agent initializes a secure connection to Datadog. In clear TCP otherwise.
	UseSSL *bool `json:"useSSL,omitempty"`
	// Disable SSL validation (useful for proxy forwarding)
	NoSSLValidation *bool `json:"noSSLValidation,omitempty"`
	// Port used to send logs over a SSL encrypted connection to Datadog. If use_http is disabled, use 10516 for the US region and 443 for the EU region.
	// +kubebuilder:validation:Minimum:=1
	// +kubebuilder:validation:Maximum:=65535
	SSLPort *uint32 `json:"sslPort,omitempty"`
	// The number of retries before the output plugin stops. Set to -1 for unlimited retries
	MaxRetries *uint32 `json:"maxRetries,omitempty"`
	// The maximum time waited between each retry in seconds
	MaxBackoff *uint32 `json:"maxBackoff,omitempty"`
	// Enable HTTP forwarding. If you disable it, make sure to change the port to 10514 or ssl_port to 10516
	UseHTTP *bool `json:"useHTTP,omitempty"`
	// Enable log compression for HTTP
	UseCompression *bool `json:"useCompression,omitempty"`
	// Set the log compression level for HTTP (1 to 9, 9 being the best ratio)
	CompressionLevel *uint32 `json:"compressionLevel,omitempty"`
	// This tells Datadog what integration it is
	DDSource *string `json:"ddSource,omitempty"`
	// Multiple value attribute. Can be used to refine the source attribute
	DDSourcecategory *string `json:"ddSourcecategory,omitempty"`
	// Custom tags with the following format "key1:value1, key2:value2"
	DDTags *string `json:"ddTags,omitempty"`
	// Used by Datadog to identify the host submitting the logs.
	DDHostname *string `json:"ddHostname,omitempty"`
	// Used by Datadog to correlate between logs, traces and metrics.
	Service *string `json:"service,omitempty"`
	// Proxy port when logs are not directly forwarded to Datadog and ssl is not used
	// +kubebuilder:validation:Minimum:=1
	// +kubebuilder:validation:Maximum:=65535
	Port *uint32 `json:"port,omitempty"`
	// Proxy endpoint when logs are not directly forwarded to Datadog
	Host *string `json:"host,omitempty"`
	// HTTP proxy, only takes effect if HTTP forwarding is enabled (use_http). Defaults to HTTP_PROXY/http_proxy env vars.
	HttpProxy *string `json:"httpProxy,omitempty"`
}
