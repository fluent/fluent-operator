package output

import "github.com/fluent/fluent-operator/v3/apis/fluentd/v1alpha1/plugins"

// Opensearch defines the parameters for out_opensearch plugin
type Opensearch struct {
	// The hostname of your Opensearch node (default: localhost).
	Host *string `json:"host,omitempty"`
	// The port number of your Opensearch node (default: 9200).
	// +kubebuilder:validation:Minimum:=1
	// +kubebuilder:validation:Maximum:=65535
	Port *uint32 `json:"port,omitempty"`
	// Hosts defines a list of hosts if you want to connect to more than one Openearch nodes
	Hosts *string `json:"hosts,omitempty"`
	// Specify https if your Opensearch endpoint supports SSL (default: http).
	Scheme *string `json:"scheme,omitempty"`
	// Path defines the REST API endpoint of Opensearch to post write requests (default: nil).
	Path *string `json:"path,omitempty"`
	// IndexName defines the placeholder syntax of Fluentd plugin API. See https://docs.fluentd.org/configuration/buffer-section.
	IndexName *string `json:"indexName,omitempty"`
	// If true, Fluentd uses the conventional index name format logstash-%Y.%m.%d (default: false). This option supersedes the index_name option.
	LogstashFormat *bool `json:"logstashFormat,omitempty"`
	// LogstashPrefix defines the logstash prefix index name to write events when logstash_format is true (default: logstash).
	LogstashPrefix *string `json:"logstashPrefix,omitempty"`
	// Optional, The login credentials to connect to Opensearch
	User *plugins.Secret `json:"user,omitempty"`
	// Optional, The login credentials to connect to Opensearch
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
	// Optional, You can specify SSL/TLS version (default: TLSv1_2)
	SslVersion *string `json:"sslVersion,omitempty"`
	// Optional, Minimum SSL/TLS version
	SslMinVersion *string `json:"sslMinVersion,omitempty"`
	// Optional, Maximum SSL/TLS version
	SslMaxVersion *string `json:"sslMaxVersion,omitempty"`
	// Optional, Enable logging of 400 reason without enabling debug log level (default: false)
	LogOs400Reason *bool `json:"logOs400Reason,omitempty"`
	// Optional, HTTP request timeout in seconds (default: 5s)
	// +kubebuilder:validation:Pattern:="^\\d+(s|m|h|d)$"
	RequestTimeout *string `json:"requestTimeout,omitempty"`
	// Optional, Indicates that the plugin should reset connection on any error (reconnect on next send) (default: false)
	ReconnectOnError *bool `json:"reconnectOnError,omitempty"`
	// Optional, Automatically reload connection after 10000 documents (default: true)
	ReloadConnections *bool `json:"reloadConnections,omitempty"`
	// Optional, When ReloadConnections true, this is the integer number of operations after which the plugin will reload the connections (default: 10000)
	ReloadAfter *uint32 `json:"reloadAfter,omitempty"`
	// Optional, Indicates that the opensearch-transport will try to reload the nodes addresses if there is a failure while making the request (default: false)
	ReloadOnFailure *bool `json:"reloadOnFailure,omitempty"`
	// Optional, You can specify times of retry obtaining OpenSearch version (default: 15)
	MaxRetryGetOsVersion *uint32 `json:"maxRetryGetOsVersion,omitempty"`
	// Optional, Indicates whether to fail when max_retry_get_os_version is exceeded (default: true)
	FailOnDetectingOsVersionRetryExceed *bool `json:"failOnDetectingOsVersionRetryExceed,omitempty"`
	// Optional, Default OpenSearch version (default: 1)
	DefaultOpensearchVersion *uint32 `json:"defaultOpensearchVersion,omitempty"`
	// Optional, Validate OpenSearch version at startup (default: true)
	VerifyOsVersionAtStartup *bool `json:"verifyOsVersionAtStartup,omitempty"`
	// Optional, Always update the template, even if it already exists (default: false)
	TemplateOverwrite *bool `json:"templateOverwrite,omitempty"`
	// Optional, You can specify times of retry putting template (default: 10)
	MaxRetryPuttingTemplate *uint32 `json:"maxRetryPuttingTemplate,omitempty"`
	// Optional, Indicates whether to fail when max_retry_putting_template is exceeded (default: true)
	FailOnPuttingTemplateRetryExceed *bool `json:"failOnPuttingTemplateRetryExceed,omitempty"`
	// Optional, Provide a different sniffer class name
	SnifferClassName *string `json:"snifferClassName,omitempty"`
	// Optional, Provide a selector class name
	SelectorClassName *string `json:"selectorClassName,omitempty"`
	// Optional, You can specify HTTP backend (default: excon). Options: excon, typhoeus
	HttpBackend *string `json:"httpBackend,omitempty"`
	// Optional, With http_backend_excon_nonblock false, plugin uses excon with nonblock=false (default: true)
	HttpBackendExconNonblock *bool `json:"httpBackendExconNonblock,omitempty"`
	// Optional, You can specify the compression level (default: no_compression). Options: no_compression, best_compression, best_speed, default_compression
	CompressionLevel *string `json:"compressionLevel,omitempty"`
	// Optional, With default behavior, plugin uses Yajl as JSON encoder/decoder. Set to true to use Oj (default: false)
	PreferOjSerializer *bool `json:"preferOjSerializer,omitempty"`
	// Optional, Suppress '[types removal]' warnings on OpenSearch 2.x (default: true for OS2+)
	SuppressTypeName *bool `json:"suppressTypeName,omitempty"`
	// Optional, With content_type application/x-ndjson, plugin adds application/x-ndjson as Content-Type (default: application/json)
	ContentType *string `json:"contentType,omitempty"`
	// Optional, Include tag key in record (default: false)
	IncludeTagKey *bool `json:"includeTagKey,omitempty"`
	// Optional, Tag key name when include_tag_key is true (default: tag)
	TagKey *string `json:"tagKey,omitempty"`
	// Optional, Record accessor syntax to specify the field to use as _id in OpenSearch
	IdKey *string `json:"idKey,omitempty"`
	// Optional, Remove specified keys from the event record
	RemoveKeys *string `json:"removeKeys,omitempty"`
	// Optional, Remove keys when record is being updated
	RemoveKeysOnUpdate *string `json:"removeKeysOnUpdate,omitempty"`
	// Optional, The write operation (default: index). Options: index, create, update, upsert
	WriteOperation *string `json:"writeOperation,omitempty"`
	// Optional, When write_operation is not index, setting this true will cause plugin to emit_error_event of records which do not include _id field (default: false)
	EmitErrorForMissingId *bool `json:"emitErrorForMissingId,omitempty"`
	// Optional, Custom headers in Hash format
	CustomHeaders *string `json:"customHeaders,omitempty"`
	// Optional, Pipeline name
	Pipeline *string `json:"pipeline,omitempty"`
	// Optional, UTC index (default: false for local time)
	UtcIndex *bool `json:"utcIndex,omitempty"`
	// Optional, Suppress doc_wrap (default: false)
	SuppressDocWrap *bool `json:"suppressDocWrap,omitempty"`
	// Optional, List of exception classes to ignore
	IgnoreExceptions *string `json:"ignoreExceptions,omitempty"`
	// Optional, Backup chunk when ignore exception occurs (default: true)
	ExceptionBackup *bool `json:"exceptionBackup,omitempty"`
	// Optional, Configure bulk_message request splitting threshold size (default: -1 unlimited)
	BulkMessageRequestThreshold *int32 `json:"bulkMessageRequestThreshold,omitempty"`
	// Optional, Specify the application name for the rollover index to be created (default: default)
	ApplicationName *string `json:"applicationName,omitempty"`
	// Optional, Specify the index date pattern for creating a rollover index (default: now/d)
	IndexDatePattern *string `json:"indexDatePattern,omitempty"`
	// Optional, Use legacy template or not (default: false for composable templates)
	UseLegacyTemplate *bool `json:"useLegacyTemplate,omitempty"`
}
