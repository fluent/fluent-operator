package output

import (
	"encoding/json"
	"fmt"

	"github.com/fluent/fluent-operator/apis/fluentd/v1alpha1/plugins"
	"github.com/fluent/fluent-operator/apis/fluentd/v1alpha1/plugins/common"
	"github.com/fluent/fluent-operator/apis/fluentd/v1alpha1/plugins/params"
)

// OutputCommon defines the common parameters for output plugin
type OutputCommon struct {
	Id *string `json:"-"`
	// The @log_level parameter specifies the plugin-specific logging level
	LogLevel *string `json:"logLevel,omitempty"`
	// The @label parameter is to route the events to <label> sections
	Label *string `json:"-"`
	// Which tag to be matched.
	Tag *string `json:"-"`
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
	// out_kafka plugin
	Kafka *Kafka2 `json:"kafka,omitempty"`
	// out_s3 plugin
	S3 *S3 `json:"s3,omitempty"`
	// out_stdout plugin
	Stdout *Stdout `json:"stdout,omitempty"`
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

	if o.BufferSection.Buffer != nil {
		child, _ := o.BufferSection.Buffer.Params(loader)
		childs = append(childs, child)
	}
	if o.BufferSection.Inject != nil {
		child, _ := o.BufferSection.Inject.Params(loader)
		childs = append(childs, child)
	}
	if o.BufferSection.Format != nil {
		child, _ := o.BufferSection.Format.Params(loader)
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
			child, _ := o.BufferSection.Format.Params(loader)
			ps.InsertChilds(child)
		}
		return o.kafka2Plugin(ps, loader), nil
	}

	if o.Elasticsearch != nil {
		ps.InsertType(string(params.ElasticsearchOutputType))
		return o.elasticsearchPlugin(ps, loader)
	}

	if o.S3 != nil {
		ps.InsertType(string(params.S3OutputType))
		return o.s3Plugin(ps, loader), nil
	}

	// if nothing defined, supposed it is a out_stdout plugin
	ps.InsertType(string(params.StdOutputType))
	return o.stdoutPlugin(ps, loader), nil
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

	if o.Forward.RequireAckResponse != nil {
		parent.InsertPairs("require_ack_response", fmt.Sprint(*o.Forward.RequireAckResponse))
	}

	if o.Forward.SendTimeout != nil {
		parent.InsertPairs("send_timeout", fmt.Sprint(*o.Forward.SendTimeout))
	}

	if o.Forward.ConnectTimeout != nil {
		parent.InsertPairs("connect_timeout", fmt.Sprint(*o.Forward.ConnectTimeout))
	}

	if o.Forward.RecoverWait != nil {
		parent.InsertPairs("recover_wait", fmt.Sprint(*o.Forward.RecoverWait))
	}

	if o.Forward.AckResponseTimeout != nil {
		parent.InsertPairs("heartbeat_type", fmt.Sprint(*o.Forward.HeartbeatType))
	}

	if o.Forward.HeartbeatInterval != nil {
		parent.InsertPairs("heartbeat_interval", fmt.Sprint(*o.Forward.HeartbeatInterval))
	}

	if o.Forward.PhiFailureDetector != nil {
		parent.InsertPairs("phi_failure_detector", fmt.Sprint(*o.Forward.PhiFailureDetector))
	}

	if o.Forward.PhiThreshold != nil {
		parent.InsertPairs("phi_threshold", fmt.Sprint(*o.Forward.PhiThreshold))
	}

	if o.Forward.HardTimeout != nil {
		parent.InsertPairs("hard_timeout", fmt.Sprint(*o.Forward.HardTimeout))
	}

	if o.Forward.ExpireDnsCache != nil {
		parent.InsertPairs("expire_dns_cache", fmt.Sprint(*o.Forward.ExpireDnsCache))
	}

	if o.Forward.DnsRoundRobin != nil {
		parent.InsertPairs("dns_round_robin", fmt.Sprint(*o.Forward.DnsRoundRobin))
	}

	if o.Forward.IgnoreNetworkErrorsAtStartup != nil {
		parent.InsertPairs("ignore_network_errors_at_startup", fmt.Sprint(*o.Forward.IgnoreNetworkErrorsAtStartup))
	}

	if o.Forward.TlsVersion != nil {
		parent.InsertPairs("tls_version", fmt.Sprint(*o.Forward.TlsVersion))
	}

	if o.Forward.TlsCiphers != nil {
		parent.InsertPairs("tls_ciphers", fmt.Sprint(*o.Forward.TlsCiphers))
	}

	if o.Forward.TlsInsecureMode != nil {
		parent.InsertPairs("tls_insecure_mode", fmt.Sprint(*o.Forward.TlsInsecureMode))
	}

	if o.Forward.TlsAllowSelfSignedCert != nil {
		parent.InsertPairs("tls_allow_self_signed_cert", fmt.Sprint(*o.Forward.TlsAllowSelfSignedCert))
	}

	if o.Forward.TlsVerifyHostname != nil {
		parent.InsertPairs("tls_verify_hostname", fmt.Sprint(*o.Forward.TlsVerifyHostname))
	}

	if o.Forward.TlsCertPath != nil {
		parent.InsertPairs("tls_cert_path", fmt.Sprint(*o.Forward.TlsCertPath))
	}
	if o.Forward.TlsClientCertPath != nil {
		parent.InsertPairs("tls_client_cert_path", fmt.Sprint(*o.Forward.TlsClientCertPath))
	}
	if o.Forward.TlsClientPrivateKeyPath != nil {
		parent.InsertPairs("tls_client_private_key_path", fmt.Sprint(*o.Forward.TlsClientPrivateKeyPath))
	}
	if o.Forward.TlsClientPrivateKeyPassphrase != nil {
		parent.InsertPairs("tls_client_private_key_passphrase", fmt.Sprint(*o.Forward.TlsClientPrivateKeyPassphrase))
	}
	if o.Forward.TlsCertThumbprint != nil {
		parent.InsertPairs("tls_cert_thumbprint", fmt.Sprint(*o.Forward.TlsCertThumbprint))
	}
	if o.Forward.TlsCertLogicalStoreName != nil {
		parent.InsertPairs("tls_cert_logical_storeName", fmt.Sprint(*o.Forward.TlsCertLogicalStoreName))
	}
	if o.Forward.TlsCertUseEnterpriseStore != nil {
		parent.InsertPairs("tls_cert_use_enterprise_store", fmt.Sprint(*o.Forward.TlsCertUseEnterpriseStore))
	}
	if o.Forward.Keepalive != nil {
		parent.InsertPairs("keepalive", fmt.Sprint(*o.Forward.Keepalive))
	}
	if o.Forward.KeepaliveTimeout != nil {
		parent.InsertPairs("keepalive_timeout", fmt.Sprint(*o.Forward.KeepaliveTimeout))
	}
	if o.Forward.VerifyConnectionAtStartup != nil {
		parent.InsertPairs("verify_connection_at_startup", fmt.Sprint(*o.Forward.VerifyConnectionAtStartup))
	}

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

	return parent
}

func (o *Output) elasticsearchPlugin(parent *params.PluginStore, loader plugins.SecretLoader) (*params.PluginStore, error) {
	if o.Elasticsearch.Host != nil {
		parent.InsertPairs("host", fmt.Sprint(*o.Elasticsearch.Host))
	}

	if o.Elasticsearch.Port != nil {
		parent.InsertPairs("port", fmt.Sprint(*o.Elasticsearch.Port))
	}

	if o.Elasticsearch.Hosts != nil {
		parent.InsertPairs("hosts", fmt.Sprint(*o.Elasticsearch.Hosts))
	}

	if o.Elasticsearch.User != nil {
		user, err := loader.LoadSecret(*o.Elasticsearch.User)
		if err != nil {
			return nil, err
		}
		parent.InsertPairs("user", user)
	}

	if o.Elasticsearch.Password != nil {
		pwd, err := loader.LoadSecret(*o.Elasticsearch.User)
		if err != nil {
			return nil, err
		}
		parent.InsertPairs("password", pwd)
	}

	if o.Elasticsearch.Scheme != nil {
		parent.InsertPairs("scheme", fmt.Sprint(*o.Elasticsearch.Scheme))
	}

	if o.Elasticsearch.Path != nil {
		parent.InsertPairs("path", fmt.Sprint(*o.Elasticsearch.Path))
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

func (o *Output) kafka2Plugin(parent *params.PluginStore, loader plugins.SecretLoader) *params.PluginStore {
	if o.Kafka.Brokers != nil {
		parent.InsertPairs("brokers", fmt.Sprint(*o.Kafka.Brokers))
	}
	if o.Kafka.TopicKey != nil {
		parent.InsertPairs("topic_key", fmt.Sprint(*o.Kafka.TopicKey))
	}
	if o.Kafka.DefaultTopic != nil {
		parent.InsertPairs("default_topic", fmt.Sprint(*o.Kafka.DefaultTopic))
	}
	if o.Kafka.UseEventTime != nil {
		parent.InsertPairs("use_event_time", fmt.Sprint(*o.Kafka.UseEventTime))
	}
	if o.Kafka.RequiredAcks != nil {
		parent.InsertPairs("required_acks", fmt.Sprint(*o.Kafka.RequiredAcks))
	}
	if o.Kafka.CompressionCodec != nil {
		parent.InsertPairs("compression_codec", fmt.Sprint(*o.Kafka.CompressionCodec))
	}

	return parent
}

func (o *Output) s3Plugin(parent *params.PluginStore, loader plugins.SecretLoader) *params.PluginStore {
	if o.S3.AwsKeyId != nil {
		parent.InsertPairs("aws_key_id", fmt.Sprint(*o.S3.AwsKeyId))
	}
	if o.S3.AwsSecKey != nil {
		parent.InsertPairs("aws_sec_key", fmt.Sprint(*o.S3.AwsSecKey))
	}
	if o.S3.S3Bucket != nil {
		parent.InsertPairs("s3_bucket", fmt.Sprint(*o.S3.S3Bucket))
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
	return parent
}

func (o *Output) stdoutPlugin(parent *params.PluginStore, loader plugins.SecretLoader) *params.PluginStore {
	return parent
}

var _ plugins.Plugin = &Output{}
