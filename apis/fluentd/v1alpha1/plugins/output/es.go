package output

import "github.com/fluent/fluent-operator/v3/apis/fluentd/v1alpha1/plugins"

// Elasticsearch defines the parameters for out_es output plugin
type ElasticsearchCommon struct {
	// The hostname of your Elasticsearch node (default: localhost).
	Host *string `json:"host,omitempty"`
	// The port number of your Elasticsearch node (default: 9200).
	// +kubebuilder:validation:Minimum:=1
	// +kubebuilder:validation:Maximum:=65535
	Port *uint32 `json:"port,omitempty"`
	// Hosts defines a list of hosts if you want to connect to more than one Elasticsearch nodes
	Hosts *string `json:"hosts,omitempty"`
	// Specify https if your Elasticsearch endpoint supports SSL (default: http).
	Scheme *string `json:"scheme,omitempty"`
	// Path defines the REST API endpoint of Elasticsearch to post write requests (default: nil).
	Path *string `json:"path,omitempty"`
	// Authenticate towards Elastic Cloud using CloudId. If set, cloudAuth must
	// be set as well and host, port, user and password are ignored.
	CloudId *plugins.Secret `json:"cloudId,omitempty"`
	// Authenticate towards Elastic Cloud using cloudAuth.
	CloudAuth *plugins.Secret `json:"cloudAuth,omitempty"`
	// Optional, The login credentials to connect to Elasticsearch
	User *plugins.Secret `json:"user,omitempty"`
	// Optional, The login credentials to connect to Elasticsearch
	Password *plugins.Secret `json:"password,omitempty"`
	// Optional, Force certificate validation
	SslVerify *bool `json:"sslVerify,omitempty"`
	// Optional, Absolute path to CA certificate file
	CAFile *string `json:"caFile,omitempty"`
	// Optional, Absolute path to client Certificate file
	ClientCert *string `json:"clientCert,omitempty"`
	// Optional, Absolute path to client private Key file
	ClientKey *string `json:"clientKey,omitempty"`
	// Optional, password for ClientKey file
	ClientKeyPassword *plugins.Secret `json:"clientKeyPassword,omitempty"`
	// Optional, Always update the template, even if it already exists (default: false)
	TemplateOverwrite *bool `json:"templateOverwrite,omitempty"`
	// Optional, You can specify times of retry putting template (default: 10)
	MaxRetryPuttingTemplate *uint32 `json:"maxRetryPuttingTemplate,omitempty"`
	// Optional, Indicates whether to fail when max_retry_putting_template is exceeded. If you have multiple output plugin, you could use this property to do not fail on fluentd statup (default: false)
	FailOnPuttingTemplateRetryExceeded *bool `json:"failOnPuttingTemplateRetryExceeded,omitempty"`
	// Optional, Indicates that the plugin should reset connection on any error (reconnect on next send) (default: false)
	ReconnectOnError *bool `json:"reconnectOnError,omitempty"`
	// Optional, Automatically reload connection after 10000 documents (default: true)
	ReloadConnections *bool `json:"reloadConnections,omitempty"`
	// Optional, Indicates that the elasticsearch-transport will try to reload the nodes addresses if there is a failure while making the request, this can be useful to quickly remove a dead node from the list of addresses (default: false)
	ReloadOnFailure *bool `json:"reloadOnFailure,omitempty"`
	// Optional, HTTP Timeout (default: 5)
	// +kubebuilder:validation:Pattern:="^\\d+(s|m|h|d)$"
	RequestTimeout *string `json:"requestTimeout,omitempty"`
	// Optional, Suppress '[types removal]' warnings on elasticsearch 7.x
	SuppressTypeName *bool `json:"suppressTypeName,omitempty"`
	// Optional, Enable Index Lifecycle Management (ILM)
	EnableIlm *bool `json:"enableIlm,omitempty"`
	// Optional, Specify ILM policy id
	IlmPolicyId *string `json:"ilmPolicyId,omitempty"`
	// Optional, Specify ILM policy contents as Hash
	IlmPolicy *string `json:"ilmPolicy,omitempty"`
	// Optional, Specify whether overwriting ilm policy or not
	IlmPolicyOverwrite *bool `json:"ilmPolicyOverride,omitempty"`
	// Optional, Enable logging of 400 reason without enabling debug log level
	LogEs400Reason *bool `json:"logEs400Reason,omitempty"`
}

type Elasticsearch struct {
	ElasticsearchCommon `json:",inline"`

	// IndexName defines the placeholder syntax of Fluentd plugin API. See https://docs.fluentd.org/configuration/buffer-section.
	IndexName *string `json:"indexName,omitempty"`
	// If true, Fluentd uses the conventional index name format logstash-%Y.%m.%d (default: false). This option supersedes the index_name option.
	LogstashFormat *bool `json:"logstashFormat,omitempty"`
	// LogstashPrefix defines the logstash prefix index name to write events when logstash_format is true (default: logstash).
	LogstashPrefix *string `json:"logstashPrefix,omitempty"`
}

type ElasticsearchDataStream struct {
	ElasticsearchCommon `json:",inline"`

	// You can specify Elasticsearch data stream name by this parameter. This parameter is mandatory for elasticsearch_data_stream
	DataStreamName *string `json:"dataStreamName"`
	// Optional, You can specify an existing matching index template for the data stream. If not present, it creates a new matching index template
	DataStreamTemplateName *string `json:"dataStreamTemplateName,omitempty"`
	// Optional, Specify whether index patterns should include a wildcard (*) when creating an index template. This is particularly useful to prevent errors in scenarios where index templates are generated automatically, and multiple services with distinct suffixes are in use
	DataStreamTemplateUseIndexPatternsWildcard *bool `json:"dataStreamTemplateUseIndexPatternsWildcard,omitempty"`
	// Optional, You can specify the name of an existing ILM policy, which will be applied to the data stream. If not present, it creates a new ILM default policy (unless data_stream_template_name is defined, in that case the ILM will be set to the one specified in the matching index template)
	DataStreamIlmName *string `json:"dataStreamIlmName,omitempty"`
	// Optional, You can specify the ILM policy contents as hash. If not present, it will apply the ILM default policy
	DataStreamIlmPolicy *string `json:"dataStreamIlmPolicy,omitempty"`
	// Optional, Specify whether the data stream ILM policy should be overwritten
	DataStreamIlmPolicyOverwrite *bool `json:"dataStreamIlmPolicyOverwrite,omitempty"`
}
