package output

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/fluent/fluent-operator/v3/apis/fluentd/v1alpha1/plugins"
	"github.com/fluent/fluent-operator/v3/apis/fluentd/v1alpha1/plugins/common"
	"github.com/fluent/fluent-operator/v3/apis/fluentd/v1alpha1/plugins/custom"
	"github.com/fluent/fluent-operator/v3/apis/fluentd/v1alpha1/plugins/params"
)

// OutputCommon defines the common parameters for output plugin
type OutputCommon struct {
	Id *string `json:"-"`
	// The @log_level parameter specifies the plugin-specific logging level
	LogLevel *string `json:"logLevel,omitempty"`
	// The @label parameter is to route the events to <label> sections
	Label *string `json:"-"`
	// Which tag to be matched.
	Tag *string `json:"tag,omitempty"`
}

// Output defines all available output plugins and their parameters
type Output struct {
	OutputCommon `json:",inline,omitempty"`
	// match setions
	common.BufferSection `json:",inline,omitempty"`
	// out_forward plugin
	Forward *Forward `json:"forward,omitempty"`
	// out_http plugin
	Http *Http `json:"http,omitempty"`
	// out_es plugin
	Elasticsearch *Elasticsearch `json:"elasticsearch,omitempty"`
	// out_es datastreams plugin
	ElasticsearchDataStream *ElasticsearchDataStream `json:"elasticsearchDataStream,omitempty"`
	// out_opensearch plugin
	Opensearch *Opensearch `json:"opensearch,omitempty"`
	// out_kafka plugin
	Kafka *Kafka2 `json:"kafka,omitempty"`
	// out_s3 plugin
	S3 *S3 `json:"s3,omitempty"`
	// out_stdout plugin
	Stdout *Stdout `json:"stdout,omitempty"`
	// out_loki plugin
	Loki *Loki `json:"loki,omitempty"`
	// Custom plugin type
	CustomPlugin *custom.CustomPlugin `json:"customPlugin,omitempty"`
	// out_cloudwatch plugin
	CloudWatch *CloudWatch `json:"cloudWatch,omitempty"`
	// datadog plugin
	Datadog *Datadog `json:"datadog,omitempty"`
	// copy plugin
	Copy *Copy `json:"copy,omitempty"`
	// null plugin
	Null *Null `json:"nullPlugin,omitempty"`
}

// DeepCopyInto implements the DeepCopyInto interface.
func (in *Output) DeepCopyInto(out *Output) {
	bytes, err := json.Marshal(*in)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(bytes, &out)
	if err != nil {
		panic(err)
	}
}

func (o *Output) Name() string {
	return "match"
}

func (o *Output) Params(loader plugins.SecretLoader) (*params.PluginStore, error) {
	ps := params.NewPluginStore(o.Name())
	childs := make([]*params.PluginStore, 0)

	ps.InsertPairs("@id", fmt.Sprint(*o.Id))

	if o.LogLevel != nil {
		ps.InsertPairs("@log_level", fmt.Sprint(*o.LogLevel))
	}

	if o.Label != nil {
		ps.InsertPairs("@label", fmt.Sprint(*o.Label))
	}

	if o.Tag != nil {
		ps.InsertPairs("tag", fmt.Sprint(*o.Tag))
	}

	if o.Buffer != nil {
		child, _ := o.Buffer.Params(loader)
		childs = append(childs, child)
	}
	if o.Inject != nil {
		child, _ := o.Inject.Params(loader)
		childs = append(childs, child)
	}
	if o.Format != nil {
		child, _ := o.Format.Params(loader)
		childs = append(childs, child)
	}

	ps.InsertChilds(childs...)

	if o.Forward != nil {
		ps.InsertType(string(params.ForwardOutputType))
		return o.forwardPlugin(ps, loader), nil
	}

	if o.Http != nil {
		ps.InsertType(string(params.HttpOutputType))
		return o.httpPlugin(ps, loader), nil
	}

	if o.Kafka != nil {
		ps.InsertType(string(params.KafkaOutputType))

		// kafka format section can not be empty
		if o.Format == nil {
			o.Format = &common.Format{
				FormatCommon: common.FormatCommon{
					Type: &params.DefaultFormatType,
				},
			}
			child, _ := o.Format.Params(loader)
			ps.InsertChilds(child)
		}
		return o.kafka2Plugin(ps, loader), nil
	}

	if o.Elasticsearch != nil {
		ps.InsertType(string(params.ElasticsearchOutputType))
		return o.elasticsearchPlugin(ps, loader)
	}

	if o.ElasticsearchDataStream != nil {
		ps.InsertType(string(params.ElasticsearchDataStreamOutputType))
		return o.elasticsearchDataStreamPlugin(ps, loader)
	}

	if o.Opensearch != nil {
		ps.InsertType(string(params.OpensearchOutputType))
		return o.opensearchPlugin(ps, loader)
	}

	if o.S3 != nil {
		ps.InsertType(string(params.S3OutputType))
		return o.s3Plugin(ps, loader), nil
	}

	if o.Loki != nil {
		ps.InsertType(string(params.LokiOutputType))
		return o.lokiPlugin(ps, loader), nil
	}

	if o.Stdout != nil {
		ps.InsertType(string(params.StdOutputType))
		return o.stdoutPlugin(ps, loader), nil
	}

	if o.CloudWatch != nil {
		ps.InsertType(string(params.CloudWatchOutputType))
		return o.cloudWatchPlugin(ps, loader), nil
	}

	if o.Datadog != nil {
		ps.InsertType(string(params.DatadogOutputType))
		return o.datadogPlugin(ps, loader), nil
	}

	if o.Copy != nil {
		ps.InsertType(string(params.CopyOutputType))
		return o.copyPlugin(ps, loader), nil
	}

	if o.Null != nil {
		ps.InsertType(string(params.NullOutputType))
		return o.nullPlugin(ps, loader), nil
	}

	return o.customOutput(ps, loader), nil

}

func (o *Output) forwardPlugin(parent *params.PluginStore, loader plugins.SecretLoader) *params.PluginStore {
	childs := make([]*params.PluginStore, 0)

	if len(o.Forward.Servers) > 0 {
		for _, s := range o.Forward.Servers {
			child, _ := s.Params(loader)
			childs = append(childs, child)
		}
	}

	if o.Forward.ServiceDiscovery != nil {
		child, _ := o.Forward.ServiceDiscovery.Params(loader)
		childs = append(childs, child)
	}

	if o.Forward.Security != nil {
		child, _ := o.Forward.Security.Params(loader)
		childs = append(childs, child)
	}

	parent.InsertChilds(childs...)

	params.InsertPairs(parent, "require_ack_response", o.Forward.RequireAckResponse)
	params.InsertPairs(parent, "send_timeout", o.Forward.SendTimeout)
	params.InsertPairs(parent, "connect_timeout", o.Forward.ConnectTimeout)
	params.InsertPairs(parent, "recover_wait", o.Forward.RecoverWait)
	params.InsertPairs(parent, "heartbeat_type", o.Forward.HeartbeatType)
	params.InsertPairs(parent, "heartbeat_interval", o.Forward.HeartbeatInterval)
	params.InsertPairs(parent, "phi_failure_detector", o.Forward.PhiFailureDetector)
	params.InsertPairs(parent, "phi_threshold", o.Forward.PhiThreshold)
	params.InsertPairs(parent, "hard_timeout", o.Forward.HardTimeout)
	params.InsertPairs(parent, "expire_dns_cache", o.Forward.ExpireDnsCache)
	params.InsertPairs(parent, "dns_round_robin", o.Forward.DnsRoundRobin)
	params.InsertPairs(parent, "ignore_network_errors_at_startup", o.Forward.IgnoreNetworkErrorsAtStartup)
	params.InsertPairs(parent, "tls_version", o.Forward.TlsVersion)
	params.InsertPairs(parent, "tls_ciphers", o.Forward.TlsCiphers)
	params.InsertPairs(parent, "tls_insecure_mode", o.Forward.TlsInsecureMode)
	params.InsertPairs(parent, "tls_allow_self_signed_cert", o.Forward.TlsAllowSelfSignedCert)
	params.InsertPairs(parent, "tls_verify_hostname", o.Forward.TlsVerifyHostname)
	params.InsertPairs(parent, "tls_cert_path", o.Forward.TlsCertPath)
	params.InsertPairs(parent, "tls_client_cert_path", o.Forward.TlsClientCertPath)
	params.InsertPairs(parent, "tls_client_private_key_path", o.Forward.TlsClientPrivateKeyPath)
	params.InsertPairs(parent, "tls_client_private_key_passphrase", o.Forward.TlsClientPrivateKeyPassphrase)
	params.InsertPairs(parent, "tls_cert_thumbprint", o.Forward.TlsCertThumbprint)
	params.InsertPairs(parent, "tls_cert_logical_storeName", o.Forward.TlsCertLogicalStoreName)
	params.InsertPairs(parent, "tls_cert_use_enterprise_store", o.Forward.TlsCertUseEnterpriseStore)
	params.InsertPairs(parent, "keepalive", o.Forward.Keepalive)
	params.InsertPairs(parent, "keepalive_timeout", o.Forward.KeepaliveTimeout)
	params.InsertPairs(parent, "verify_connection_at_startup", o.Forward.VerifyConnectionAtStartup)

	return parent
}

func (o *Output) httpPlugin(parent *params.PluginStore, loader plugins.SecretLoader) *params.PluginStore {
	if o.Http.Auth != nil {
		child, _ := o.Http.Params(loader)
		parent.InsertChilds(child)
	}

	if o.Http.Endpoint != nil {
		parent.InsertPairs("endpoint", fmt.Sprint(*o.Http.Endpoint))
	}

	if o.Http.HttpMethod != nil {
		parent.InsertPairs("http_method", fmt.Sprint(*o.Http.HttpMethod))
	}

	if o.Http.Proxy != nil {
		parent.InsertPairs("proxy", fmt.Sprint(*o.Http.Proxy))
	}

	if o.Http.ContentType != nil {
		parent.InsertPairs("content_type", fmt.Sprint(*o.Http.ContentType))
	}

	if o.Http.JsonArray != nil {
		parent.InsertPairs("json_array", fmt.Sprint(*o.Http.JsonArray))
	}

	if o.Http.Headers != nil {
		parent.InsertPairs("headers", fmt.Sprint(*o.Http.Headers))
	}

	if o.Http.HeadersFromPlaceholders != nil {
		parent.InsertPairs("headers_from_placeholders", fmt.Sprint(*o.Http.HeadersFromPlaceholders))
	}

	if o.Http.OpenTimeout != nil {
		parent.InsertPairs("open_timeout", fmt.Sprint(*o.Http.OpenTimeout))
	}

	if o.Http.ReadTimeout != nil {
		parent.InsertPairs("read_timeout", fmt.Sprint(*o.Http.ReadTimeout))
	}

	if o.Http.SslTimeout != nil {
		parent.InsertPairs("ssl_timeout", fmt.Sprint(*o.Http.SslTimeout))
	}

	if o.Http.TlsCaCertPath != nil {
		parent.InsertPairs("tls_ca_cert_path", fmt.Sprint(*o.Http.TlsCaCertPath))
	}

	if o.Http.TlsClientCertPath != nil {
		parent.InsertPairs("tls_client_cert_path", fmt.Sprint(*o.Http.TlsClientCertPath))
	}

	if o.Http.TlsPrivateKeyPath != nil {
		parent.InsertPairs("tls_private_key_path", fmt.Sprint(*o.Http.TlsPrivateKeyPath))
	}

	if o.Http.TlsPrivateKeyPassphrase != nil {
		parent.InsertPairs("tls_private_key_passphrase", fmt.Sprint(*o.Http.TlsPrivateKeyPassphrase))
	}

	if o.Http.TlsVerifyMode != nil {
		parent.InsertPairs("tls_verify_mode", fmt.Sprint(*o.Http.TlsVerifyMode))
	}

	if o.Http.TlsVersion != nil {
		parent.InsertPairs("tls_version", fmt.Sprint(*o.Http.TlsVersion))
	}

	if o.Http.TlsCiphers != nil {
		parent.InsertPairs("tls_ciphers", fmt.Sprint(*o.Http.TlsCiphers))
	}

	if o.Http.ErrorResponseAsUnrecoverable != nil {
		parent.InsertPairs("error_response_as_unrecoverable", fmt.Sprint(*o.Http.ErrorResponseAsUnrecoverable))
	}

	if o.Http.RetryableResponseCodes != nil {
		parent.InsertPairs("retryable_response_codes", fmt.Sprint(*o.Http.RetryableResponseCodes))
	}

	if o.Http.Compress != nil {
		parent.InsertPairs("compress", fmt.Sprint(*o.Http.Compress))
	}

	return parent
}

func (o *Output) elasticsearchPluginCommon(cmn *ElasticsearchCommon, parent *params.PluginStore, loader plugins.SecretLoader) (*params.PluginStore, error) {
	params.InsertPairs(parent, "host", cmn.Host)
	params.InsertPairs(parent, "port", cmn.Port)
	params.InsertPairs(parent, "hosts", cmn.Hosts)

	if cmn.User != nil {
		user, err := loader.LoadSecret(*cmn.User)
		if err != nil {
			return nil, err
		}
		parent.InsertPairs("user", user)
	}

	if cmn.Password != nil {
		pwd, err := loader.LoadSecret(*cmn.Password)
		if err != nil {
			return nil, err
		}
		parent.InsertPairs("password", pwd)
	}

	params.InsertPairs(parent, "ssl_verify", cmn.SslVerify)
	params.InsertPairs(parent, "ca_file", cmn.CAFile)

	if cmn.CloudAuth != nil {
		cloudauth, err := loader.LoadSecret(*cmn.CloudAuth)
		if err != nil {
			return nil, err
		}
		parent.InsertPairs("cloud_auth", cloudauth)
	}

	if cmn.CloudId != nil {
		cloudid, err := loader.LoadSecret(*cmn.CloudId)
		if err != nil {
			return nil, err
		}
		parent.InsertPairs("cloud_id", cloudid)
	}

	params.InsertPairs(parent, "client_cert", cmn.ClientCert)
	params.InsertPairs(parent, "client_key", cmn.ClientKey)

	if cmn.ClientKeyPassword != nil {
		pwd, err := loader.LoadSecret(*cmn.ClientKeyPassword)
		if err != nil {
			return nil, err
		}
		parent.InsertPairs("client_key_pass", pwd)
	}

	params.InsertPairs(parent, "scheme", cmn.Scheme)
	params.InsertPairs(parent, "path", cmn.Path)
	params.InsertPairs(parent, "template_overwrite", cmn.TemplateOverwrite)
	params.InsertPairs(parent, "max_retry_putting_template", cmn.MaxRetryPuttingTemplate)
	params.InsertPairs(parent, "fail_on_putting_template_retry_exceed", cmn.FailOnPuttingTemplateRetryExceeded)
	params.InsertPairs(parent, "reconnect_on_error", cmn.ReconnectOnError)
	params.InsertPairs(parent, "reload_after", cmn.ReloadAfter)
	params.InsertPairs(parent, "reload_connections", cmn.ReloadConnections)
	params.InsertPairs(parent, "reload_on_failure", cmn.ReloadOnFailure)
	params.InsertPairs(parent, "request_timeout", cmn.RequestTimeout)
	params.InsertPairs(parent, "sniffer_class_name", cmn.SnifferClassName)
	params.InsertPairs(parent, "suppress_type_name", cmn.SuppressTypeName)
	params.InsertPairs(parent, "enable_ilm", cmn.EnableIlm)
	params.InsertPairs(parent, "ilm_policy_id", cmn.IlmPolicyId)
	params.InsertPairs(parent, "ilm_policy", cmn.IlmPolicy)
	params.InsertPairs(parent, "ilm_policy_overwrite", cmn.IlmPolicyOverwrite)
	params.InsertPairs(parent, "log_es_400_reason", cmn.LogEs400Reason)

	return parent, nil
}

func (o *Output) elasticsearchPlugin(parent *params.PluginStore, loader plugins.SecretLoader) (*params.PluginStore, error) {

	parent, err := o.elasticsearchPluginCommon(&o.Elasticsearch.ElasticsearchCommon, parent, loader)
	if err != nil {
		return nil, err
	}

	if o.Elasticsearch.IndexName != nil {
		parent.InsertPairs("index_name", fmt.Sprint(*o.Elasticsearch.IndexName))
	}

	if o.Elasticsearch.LogstashFormat != nil {
		parent.InsertPairs("logstash_format", fmt.Sprint(*o.Elasticsearch.LogstashFormat))
	}

	if o.Elasticsearch.LogstashPrefix != nil {
		parent.InsertPairs("logstash_prefix", fmt.Sprint(*o.Elasticsearch.LogstashPrefix))
	}

	return parent, nil
}

func (o *Output) elasticsearchDataStreamPlugin(parent *params.PluginStore, loader plugins.SecretLoader) (*params.PluginStore, error) {

	parent, err := o.elasticsearchPluginCommon(&o.ElasticsearchDataStream.ElasticsearchCommon, parent, loader)
	if err != nil {
		return nil, err
	}

	if o.ElasticsearchDataStream.DataStreamName != nil {
		parent.InsertPairs("data_stream_name", fmt.Sprint(*o.ElasticsearchDataStream.DataStreamName))
	}

	if o.ElasticsearchDataStream.DataStreamTemplateName != nil {
		parent.InsertPairs("data_stream_template_name", fmt.Sprint(*o.ElasticsearchDataStream.DataStreamTemplateName))
	}

	if o.ElasticsearchDataStream.DataStreamTemplateUseIndexPatternsWildcard != nil {
		parent.InsertPairs("data_stream_template_use_index_patterns_wildcard", fmt.Sprint(*o.ElasticsearchDataStream.DataStreamTemplateUseIndexPatternsWildcard))
	}

	if o.ElasticsearchDataStream.DataStreamIlmName != nil {
		parent.InsertPairs("data_stream_ilm_name", fmt.Sprint(*o.ElasticsearchDataStream.DataStreamIlmName))
	}

	if o.ElasticsearchDataStream.DataStreamIlmPolicy != nil {
		parent.InsertPairs("data_stream_ilm_policy", fmt.Sprint(*o.ElasticsearchDataStream.DataStreamIlmPolicy))
	}

	if o.ElasticsearchDataStream.DataStreamIlmPolicyOverwrite != nil {
		parent.InsertPairs("data_stream_ilm_policy_overwrite", fmt.Sprint(*o.ElasticsearchDataStream.DataStreamIlmPolicyOverwrite))
	}

	return parent, nil
}

func (o *Output) opensearchPlugin(parent *params.PluginStore, loader plugins.SecretLoader) (*params.PluginStore, error) {
	if err := o.opensearchBasicConnection(parent, loader); err != nil {
		return nil, err
	}
	o.opensearchIndexConfig(parent)
	if err := o.opensearchSSLConfig(parent, loader); err != nil {
		return nil, err
	}
	o.opensearchConnectionManagement(parent)
	o.opensearchVersionDetection(parent)
	o.opensearchTemplateManagement(parent)
	o.opensearchPerformanceTuning(parent)
	o.opensearchRecordHandling(parent)
	o.opensearchAdvancedOptions(parent)

	return parent, nil
}

func (o *Output) opensearchBasicConnection(parent *params.PluginStore, loader plugins.SecretLoader) error {
	params.InsertPairs(parent, "host", o.Opensearch.Host)
	params.InsertPairs(parent, "port", o.Opensearch.Port)
	params.InsertPairs(parent, "hosts", o.Opensearch.Hosts)
	if o.Opensearch.User != nil {
		user, err := loader.LoadSecret(*o.Opensearch.User)
		if err != nil {
			return err
		}
		parent.InsertPairs("user", user)
	}
	if o.Opensearch.Password != nil {
		pwd, err := loader.LoadSecret(*o.Opensearch.Password)
		if err != nil {
			return err
		}
		parent.InsertPairs("password", pwd)
	}
	params.InsertPairs(parent, "scheme", o.Opensearch.Scheme)
	params.InsertPairs(parent, "path", o.Opensearch.Path)
	return nil
}

func (o *Output) opensearchIndexConfig(parent *params.PluginStore) {
	params.InsertPairs(parent, "index_name", o.Opensearch.IndexName)
	params.InsertPairs(parent, "logstash_format", o.Opensearch.LogstashFormat)
	params.InsertPairs(parent, "logstash_prefix", o.Opensearch.LogstashPrefix)
	params.InsertPairs(parent, "index_date_pattern", o.Opensearch.IndexDatePattern)
	params.InsertPairs(parent, "utc_index", o.Opensearch.UtcIndex)
}

func (o *Output) opensearchSSLConfig(parent *params.PluginStore, loader plugins.SecretLoader) error {
	params.InsertPairs(parent, "ssl_verify", o.Opensearch.SslVerify)
	params.InsertPairs(parent, "ca_file", o.Opensearch.CAFile)
	params.InsertPairs(parent, "client_cert", o.Opensearch.ClientCert)
	params.InsertPairs(parent, "client_key", o.Opensearch.ClientKey)
	if o.Opensearch.ClientKeyPassword != nil {
		pwd, err := loader.LoadSecret(*o.Opensearch.ClientKeyPassword)
		if err != nil {
			return err
		}
		parent.InsertPairs("client_key_pass", pwd)
	}
	params.InsertPairs(parent, "ssl_version", o.Opensearch.SslVersion)
	params.InsertPairs(parent, "ssl_min_version", o.Opensearch.SslMinVersion)
	params.InsertPairs(parent, "ssl_max_version", o.Opensearch.SslMaxVersion)
	return nil
}

func (o *Output) opensearchConnectionManagement(parent *params.PluginStore) {
	params.InsertPairs(parent, "log_os_400_reason", o.Opensearch.LogOs400Reason)
	params.InsertPairs(parent, "request_timeout", o.Opensearch.RequestTimeout)
	params.InsertPairs(parent, "reconnect_on_error", o.Opensearch.ReconnectOnError)
	params.InsertPairs(parent, "reload_connections", o.Opensearch.ReloadConnections)
	params.InsertPairs(parent, "reload_after", o.Opensearch.ReloadAfter)
	params.InsertPairs(parent, "reload_on_failure", o.Opensearch.ReloadOnFailure)
}

func (o *Output) opensearchVersionDetection(parent *params.PluginStore) {
	params.InsertPairs(parent, "max_retry_get_os_version", o.Opensearch.MaxRetryGetOsVersion)
	params.InsertPairs(parent, "fail_on_detecting_os_version_retry_exceed", o.Opensearch.FailOnDetectingOsVersionRetryExceed)
	params.InsertPairs(parent, "default_opensearch_version", o.Opensearch.DefaultOpensearchVersion)
	params.InsertPairs(parent, "verify_os_version_at_startup", o.Opensearch.VerifyOsVersionAtStartup)
}

func (o *Output) opensearchTemplateManagement(parent *params.PluginStore) {
	params.InsertPairs(parent, "template_overwrite", o.Opensearch.TemplateOverwrite)
	params.InsertPairs(parent, "max_retry_putting_template", o.Opensearch.MaxRetryPuttingTemplate)
	params.InsertPairs(parent, "fail_on_putting_template_retry_exceed", o.Opensearch.FailOnPuttingTemplateRetryExceed)
	params.InsertPairs(parent, "use_legacy_template", o.Opensearch.UseLegacyTemplate)
}

func (o *Output) opensearchPerformanceTuning(parent *params.PluginStore) {
	params.InsertPairs(parent, "sniffer_class_name", o.Opensearch.SnifferClassName)
	params.InsertPairs(parent, "selector_class_name", o.Opensearch.SelectorClassName)
	params.InsertPairs(parent, "http_backend", o.Opensearch.HttpBackend)
	params.InsertPairs(parent, "http_backend_excon_nonblock", o.Opensearch.HttpBackendExconNonblock)
	params.InsertPairs(parent, "compression_level", o.Opensearch.CompressionLevel)
	params.InsertPairs(parent, "prefer_oj_serializer", o.Opensearch.PreferOjSerializer)
	params.InsertPairs(parent, "bulk_message_request_threshold", o.Opensearch.BulkMessageRequestThreshold)
}

func (o *Output) opensearchRecordHandling(parent *params.PluginStore) {
	params.InsertPairs(parent, "suppress_type_name", o.Opensearch.SuppressTypeName)
	params.InsertPairs(parent, "content_type", o.Opensearch.ContentType)
	params.InsertPairs(parent, "include_tag_key", o.Opensearch.IncludeTagKey)
	params.InsertPairs(parent, "tag_key", o.Opensearch.TagKey)
	params.InsertPairs(parent, "id_key", o.Opensearch.IdKey)
	params.InsertPairs(parent, "remove_keys", o.Opensearch.RemoveKeys)
	params.InsertPairs(parent, "remove_keys_on_update", o.Opensearch.RemoveKeysOnUpdate)
	params.InsertPairs(parent, "write_operation", o.Opensearch.WriteOperation)
	params.InsertPairs(parent, "emit_error_for_missing_id", o.Opensearch.EmitErrorForMissingId)
	params.InsertPairs(parent, "suppress_doc_wrap", o.Opensearch.SuppressDocWrap)
}

func (o *Output) opensearchAdvancedOptions(parent *params.PluginStore) {
	params.InsertPairs(parent, "custom_headers", o.Opensearch.CustomHeaders)
	params.InsertPairs(parent, "pipeline", o.Opensearch.Pipeline)
	params.InsertPairs(parent, "ignore_exceptions", o.Opensearch.IgnoreExceptions)
	params.InsertPairs(parent, "exception_backup", o.Opensearch.ExceptionBackup)
	params.InsertPairs(parent, "application_name", o.Opensearch.ApplicationName)
}

func (o *Output) kafka2Plugin(parent *params.PluginStore, _ plugins.SecretLoader) *params.PluginStore {
	params.InsertPairs(parent, "brokers", o.Kafka.Brokers)
	params.InsertPairs(parent, "topic_key", o.Kafka.TopicKey)
	params.InsertPairs(parent, "default_topic", o.Kafka.DefaultTopic)
	params.InsertPairs(parent, "use_event_time", o.Kafka.UseEventTime)
	params.InsertPairs(parent, "required_acks", o.Kafka.RequiredAcks)
	params.InsertPairs(parent, "compression_codec", o.Kafka.CompressionCodec)

	return parent
}

func (o *Output) s3Plugin(parent *params.PluginStore, loader plugins.SecretLoader) *params.PluginStore {
	if o.S3.AwsKeyIdFromSecret != nil {
		value, err := loader.LoadSecret(*o.S3.AwsKeyIdFromSecret)
		if err != nil {
			return nil
		}
		parent.InsertPairs("aws_key_id", value)
	} else if o.S3.AwsKeyId != nil {
		parent.InsertPairs("aws_key_id", fmt.Sprint(*o.S3.AwsKeyId))
	}
	if o.S3.AwsSecKeyFromSecret != nil {
		value, err := loader.LoadSecret(*o.S3.AwsSecKeyFromSecret)
		if err != nil {
			return nil
		}
		parent.InsertPairs("aws_sec_key", value)
	} else if o.S3.AwsSecKey != nil {
		parent.InsertPairs("aws_sec_key", fmt.Sprint(*o.S3.AwsSecKey))
	}
	if o.S3.S3Bucket != nil {
		parent.InsertPairs("s3_bucket", fmt.Sprint(*o.S3.S3Bucket))
	}
	if o.S3.S3Region != nil {
		parent.InsertPairs("s3_region", fmt.Sprint(*o.S3.S3Region))
	}
	if o.S3.S3Endpoint != nil {
		parent.InsertPairs("s3_endpoint", fmt.Sprint(*o.S3.S3Endpoint))
	}
	if o.S3.ForcePathStyle != nil {
		parent.InsertPairs("force_path_style", fmt.Sprint(*o.S3.ForcePathStyle))
	}
	if o.S3.TimeSliceFormat != nil {
		parent.InsertPairs("time_slice_format", fmt.Sprint(*o.S3.TimeSliceFormat))
	}
	if o.S3.Path != nil {
		parent.InsertPairs("path", fmt.Sprint(*o.S3.Path))
	}
	if o.S3.S3ObjectKeyFormat != nil {
		parent.InsertPairs("s3_object_key_format", fmt.Sprint(*o.S3.S3ObjectKeyFormat))
	}
	if o.S3.StoreAs != nil {
		parent.InsertPairs("store_as", fmt.Sprint(*o.S3.StoreAs))
	}
	if o.S3.ProxyUri != nil {
		parent.InsertPairs("proxy_uri", fmt.Sprint(*o.S3.ProxyUri))
	}
	if o.S3.SslVerifyPeer != nil {
		parent.InsertPairs("ssl_verify_peer", fmt.Sprint(*o.S3.SslVerifyPeer))
	}
	if o.S3.UseServerSideEncryption != nil {
		parent.InsertPairs("use_server_side_encryption", fmt.Sprint(*o.S3.UseServerSideEncryption))
	}
	if o.S3.SseCustomerAlgorithm != nil {
		parent.InsertPairs("sse_customer_algorithm", fmt.Sprint(*o.S3.SseCustomerAlgorithm))
	}
	if o.S3.SsekmsKeyId != nil {
		parent.InsertPairs("ssekms_key_id", fmt.Sprint(*o.S3.SsekmsKeyId))
	}
	if o.S3.SseCustomerKey != nil {
		parent.InsertPairs("sse_customer_key", fmt.Sprint(*o.S3.SseCustomerKey))
	}
	if o.S3.SseCustomerKeyMd5 != nil {
		parent.InsertPairs("sse_customer_key_md5", fmt.Sprint(*o.S3.SseCustomerKeyMd5))
	}
	return parent
}

func (o *Output) lokiPlugin(parent *params.PluginStore, loader plugins.SecretLoader) *params.PluginStore {
	if o.Loki.Url != nil {
		parent.InsertPairs("url", fmt.Sprint(*o.Loki.Url))
	}
	if o.Loki.HTTPUser != nil {
		u, err := loader.LoadSecret(*o.Loki.HTTPUser)
		if err != nil {
			return nil
		}
		parent.InsertPairs("username", u)
	}
	if o.Loki.HTTPPasswd != nil {
		passwd, err := loader.LoadSecret(*o.Loki.HTTPPasswd)
		if err != nil {
			return nil
		}
		parent.InsertPairs("password", passwd)
	}
	if o.Loki.BearerTokenFile != nil {
		parent.InsertPairs("bearer_token_file", fmt.Sprint(*o.Loki.BearerTokenFile))
	}
	if o.Loki.TenantID != nil {
		id, err := loader.LoadSecret(*o.Loki.TenantID)
		if err != nil {
			return nil
		}
		parent.InsertPairs("tenant", id)
	}
	if len(o.Loki.Labels) > 0 {
		labels := make(map[string]string)
		for _, l := range o.Loki.Labels {
			key, value, found := strings.Cut(l, "=")
			if !found {
				continue
			}
			labels[strings.TrimSpace(key)] = strings.TrimSpace(value)
		}
		if len(labels) > 0 {
			jsonStr, err := json.Marshal(labels)
			if err != nil {
				fmt.Printf("Error: %s", err.Error())
			} else {
				parent.InsertPairs("extra_labels", string(jsonStr))
			}
		}
	}
	if len(o.Loki.RemoveKeys) > 0 {
		parent.InsertPairs("remove_keys", strings.Join(o.Loki.RemoveKeys, ","))
	}
	if len(o.Loki.LabelKeys) > 0 {
		ps := params.NewPluginStore("label")
		for _, n := range o.Loki.LabelKeys {
			ps.InsertPairs(n, n)
		}
		parent.InsertChilds(ps)
	}
	if o.Loki.LineFormat != "" {
		parent.InsertPairs("line_format", o.Loki.LineFormat)
	}
	if o.Loki.ExtractKubernetesLabels != nil {
		parent.InsertPairs("extract_kubernetes_labels", fmt.Sprint(*o.Loki.ExtractKubernetesLabels))
	}
	if o.Loki.DropSingleKey != nil {
		parent.InsertPairs("drop_single_key", fmt.Sprint(*o.Loki.DropSingleKey))
	}
	if o.Loki.IncludeThreadLabel != nil {
		parent.InsertPairs("include_thread_label", fmt.Sprint(*o.Loki.IncludeThreadLabel))
	}
	if o.Loki.Insecure != nil {
		parent.InsertPairs("insecure_tls", fmt.Sprint(*o.Loki.Insecure))
	}
	if o.Loki.TlsCaCertFile != nil {
		parent.InsertPairs("ca_cert", fmt.Sprint(*o.Loki.TlsCaCertFile))
	}
	if o.Loki.TlsClientCertFile != nil {
		parent.InsertPairs("cert", fmt.Sprint(*o.Loki.TlsClientCertFile))
	}
	if o.Loki.TlsPrivateKeyFile != nil {
		parent.InsertPairs("key", fmt.Sprint(*o.Loki.TlsPrivateKeyFile))
	}
	return parent
}

func (o *Output) cloudWatchPlugin(parent *params.PluginStore, sl plugins.SecretLoader) *params.PluginStore {
	childs := make([]*params.PluginStore, 0)

	params.InsertPairs(parent, "auto_create_stream", o.CloudWatch.AutoCreateStream)
	if o.CloudWatch.AwsKeyId != nil {
		value, err := sl.LoadSecret(*o.CloudWatch.AwsKeyId)
		if err != nil {
			return nil
		}
		parent.InsertPairs("aws_key_id", value)
	}
	if o.CloudWatch.AwsSecKey != nil {
		value, err := sl.LoadSecret(*o.CloudWatch.AwsSecKey)
		if err != nil {
			return nil
		}
		parent.InsertPairs("aws_sec_key", value)
	}
	params.InsertPairs(parent, "aws_use_sts", o.CloudWatch.AwsUseSts)
	params.InsertPairs(parent, "aws_sts_role_arn", o.CloudWatch.AwsStsRoleARN)
	params.InsertPairs(parent, "aws_sts_session_name", o.CloudWatch.AwsStsSessionName)
	params.InsertPairs(parent, "aws_sts_external_id", o.CloudWatch.AwsStsExternalId)
	params.InsertPairs(parent, "aws_sts_policy", o.CloudWatch.AwsStsPolicy)
	params.InsertPairs(parent, "aws_sts_duration_seconds", o.CloudWatch.AwsStsDurationSeconds)
	params.InsertPairs(parent, "aws_sts_endpoint_url", o.CloudWatch.AwsStsEndpointUrl)
	params.InsertPairs(parent, "aws_ecs_authentication", o.CloudWatch.AwsEcsAuthentication)
	params.InsertPairs(parent, "concurrency", o.CloudWatch.Concurrency)
	params.InsertPairs(parent, "endpoint", o.CloudWatch.Endpoint)
	params.InsertPairs(parent, "ssl_verify_peer", o.CloudWatch.SslVerifyPeer)
	params.InsertPairs(parent, "http_proxy", o.CloudWatch.HttpProxy)
	params.InsertPairs(parent, "include_time_key", o.CloudWatch.IncludeTimeKey)
	params.InsertPairs(parent, "json_handler", o.CloudWatch.JsonHandler)
	params.InsertPairs(parent, "localtime", o.CloudWatch.Localtime)
	params.InsertPairs(parent, "log_group_aws_tags", o.CloudWatch.LogGroupAwsTags)
	params.InsertPairs(parent, "log_group_aws_tags_key", o.CloudWatch.LogGroupAwsTagsKey)
	params.InsertPairs(parent, "log_group_name", o.CloudWatch.LogGroupName)
	params.InsertPairs(parent, "log_group_name_key", o.CloudWatch.LogGroupNameKey)
	params.InsertPairs(parent, "log_rejected_request", o.CloudWatch.LogRejectedRequest)
	params.InsertPairs(parent, "log_stream_name", o.CloudWatch.LogStreamName)
	params.InsertPairs(parent, "log_stream_name_key", o.CloudWatch.LogStreamNameKey)
	params.InsertPairs(parent, "max_events_per_batch", o.CloudWatch.MaxEventsPerBatch)
	params.InsertPairs(parent, "max_message_length", o.CloudWatch.MaxMessageLength)
	params.InsertPairs(parent, "message_keys", o.CloudWatch.MessageKeys)
	params.InsertPairs(parent, "put_log_events_disable_retry_limit", o.CloudWatch.PutLogEventsDisableRetryLimit)
	params.InsertPairs(parent, "put_log_events_retry_limit", o.CloudWatch.PutLogEventsRetryLimit)
	params.InsertPairs(parent, "put_log_events_retry_wait", o.CloudWatch.PutLogEventsRetryWait)
	params.InsertPairs(parent, "region", o.CloudWatch.Region)
	params.InsertPairs(parent, "remove_log_group_aws_tags_key", o.CloudWatch.RemoveLogGroupAwsTagsKey)
	params.InsertPairs(parent, "remove_log_group_name_key", o.CloudWatch.RemoveLogGroupNameKey)
	params.InsertPairs(parent, "remove_log_stream_name_key", o.CloudWatch.RemoveLogStreamNameKey)
	params.InsertPairs(parent, "remove_retention_in_days_key", o.CloudWatch.RemoveRetentionInDaysKey)
	params.InsertPairs(parent, "retention_in_days", o.CloudWatch.RetentionInDays)
	params.InsertPairs(parent, "retention_in_days_key", o.CloudWatch.RetentionInDaysKey)
	params.InsertPairs(parent, "use_tag_as_group", o.CloudWatch.UseTagAsGroup)
	params.InsertPairs(parent, "use_tag_as_stream", o.CloudWatch.UseTagAsStream)
	params.InsertPairs(parent, "policy", o.CloudWatch.Policy)
	params.InsertPairs(parent, "duration_seconds", o.CloudWatch.DurationSeconds)

	// web_identity_credentials is a subsection of its own containing AWS credential settings
	child := params.NewPluginStore("web_identity_credentials")
	params.InsertPairs(child, "role_arn", o.CloudWatch.RoleARN)
	params.InsertPairs(child, "web_identity_token_file", o.CloudWatch.WebIdentityTokenFile)
	params.InsertPairs(child, "role_session_name", o.CloudWatch.RoleSessionName)
	childs = append(childs, child)

	// format is a subsection of its own.  Not implemented yet.
	parent.InsertChilds(childs...)
	return parent
}

func (o *Output) stdoutPlugin(parent *params.PluginStore, loader plugins.SecretLoader) *params.PluginStore {
	return parent
}

func (o *Output) customOutput(parent *params.PluginStore, loader plugins.SecretLoader) *params.PluginStore {
	if o.CustomPlugin == nil {
		return parent
	}
	customPlugin, _ := o.CustomPlugin.Params(loader)
	parent.Content = customPlugin.Content
	return parent
}

func (o *Output) datadogPlugin(parent *params.PluginStore, sl plugins.SecretLoader) *params.PluginStore {
	if o.Datadog.ApiKey != nil {
		apiKey, err := sl.LoadSecret(*o.Datadog.ApiKey)
		if err != nil {
			return nil
		}
		parent.InsertPairs("api_key", apiKey)
	}

	if o.Datadog.UseJson != nil {
		parent.InsertPairs("use_json", fmt.Sprint(*o.Datadog.UseJson))
	}

	if o.Datadog.IncludeTagKey != nil {
		parent.InsertPairs("include_tag_key", fmt.Sprint(*o.Datadog.IncludeTagKey))
	}

	if o.Datadog.TagKey != nil {
		parent.InsertPairs("tag_key", fmt.Sprint(*o.Datadog.TagKey))
	}

	if o.Datadog.TimestampKey != nil {
		parent.InsertPairs("timestamp_key", fmt.Sprint(*o.Datadog.TimestampKey))
	}

	if o.Datadog.UseSSL != nil {
		parent.InsertPairs("use_ssl", fmt.Sprint(*o.Datadog.UseSSL))
	}

	if o.Datadog.NoSSLValidation != nil {
		parent.InsertPairs("no_ssl_validation", fmt.Sprint(*o.Datadog.NoSSLValidation))
	}

	if o.Datadog.SSLPort != nil {
		parent.InsertPairs("ssl_port", fmt.Sprint(*o.Datadog.SSLPort))
	}

	if o.Datadog.MaxRetries != nil {
		parent.InsertPairs("max_retries", fmt.Sprint(*o.Datadog.MaxRetries))
	}

	if o.Datadog.MaxBackoff != nil {
		parent.InsertPairs("max_backoff", fmt.Sprint(*o.Datadog.MaxBackoff))
	}

	if o.Datadog.UseHTTP != nil {
		parent.InsertPairs("use_http", fmt.Sprint(*o.Datadog.UseHTTP))
	}

	if o.Datadog.UseCompression != nil {
		parent.InsertPairs("use_compression", fmt.Sprint(*o.Datadog.UseCompression))
	}

	if o.Datadog.CompressionLevel != nil {
		parent.InsertPairs("compression_level", fmt.Sprint(*o.Datadog.CompressionLevel))
	}

	if o.Datadog.DDSource != nil {
		parent.InsertPairs("dd_source", fmt.Sprint(*o.Datadog.DDSource))
	}

	if o.Datadog.DDSourcecategory != nil {
		parent.InsertPairs("dd_sourcecategory", fmt.Sprint(*o.Datadog.DDSourcecategory))
	}

	if o.Datadog.DDTags != nil {
		parent.InsertPairs("dd_tags", fmt.Sprint(*o.Datadog.DDTags))
	}

	if o.Datadog.DDHostname != nil {
		parent.InsertPairs("dd_hostname", fmt.Sprint(*o.Datadog.DDHostname))
	}

	if o.Datadog.Service != nil {
		parent.InsertPairs("service", fmt.Sprint(*o.Datadog.Service))
	}

	if o.Datadog.Port != nil {
		parent.InsertPairs("port", fmt.Sprint(*o.Datadog.Port))
	}

	if o.Datadog.Host != nil {
		parent.InsertPairs("host", fmt.Sprint(*o.Datadog.Host))
	}

	if o.Datadog.HttpProxy != nil {
		parent.InsertPairs("http_proxy", fmt.Sprint(*o.Datadog.HttpProxy))
	}
	return parent
}

func (o *Output) copyPlugin(parent *params.PluginStore, _ plugins.SecretLoader) *params.PluginStore {
	if o.Copy.CopyMode != nil {
		parent.InsertPairs("copy_mode", fmt.Sprint(*o.Copy.CopyMode))
	}
	return parent
}

func (o *Output) nullPlugin(parent *params.PluginStore, _ plugins.SecretLoader) *params.PluginStore {
	if o.Null.NeverFlush != nil {
		parent.InsertPairs("never_flush", fmt.Sprint(*o.Null.NeverFlush))
	}
	return parent
}

var _ plugins.Plugin = &Output{}
