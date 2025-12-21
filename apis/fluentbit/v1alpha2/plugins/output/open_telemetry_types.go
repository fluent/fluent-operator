package output

import (
	"fmt"

	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins/params"
)

// +kubebuilder:object:generate:=true

// The OpenTelemetry plugin allows you to take logs, metrics, and traces from Fluent Bit and submit them to an OpenTelemetry HTTP endpoint. <br />
// **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/outputs/opentelemetry**
type OpenTelemetry struct {
	// IP address or hostname of the target HTTP Server, default `127.0.0.1`
	Host string `json:"host,omitempty"`
	// TCP port of the target OpenSearch instance, default `80`
	// +kubebuilder:validation:Minimum:=1
	// +kubebuilder:validation:Maximum:=65535
	Port *int32 `json:"port,omitempty"`
	// Optional username credential for access
	HTTPUser *plugins.Secret `json:"httpUser,omitempty"`
	// Password for user defined in HTTP_User
	HTTPPasswd *plugins.Secret `json:"httpPassword,omitempty"`
	// Specify an HTTP Proxy. The expected format of this value is http://HOST:PORT. Note that HTTPS is not currently supported.
	// It is recommended not to set this and to configure the HTTP proxy environment variables instead as they support both HTTP and HTTPS.
	Proxy string `json:"proxy,omitempty"`
	// Specify an optional HTTP URI for the target web server listening for metrics, e.g: /v1/metrics
	MetricsUri string `json:"metricsUri,omitempty"`
	// Specify an optional HTTP URI for the target web server listening for logs, e.g: /v1/logs
	LogsUri string `json:"logsUri,omitempty"`
	// Specify an optional HTTP URI for the target web server listening for traces, e.g: /v1/traces
	TracesUri string `json:"tracesUri,omitempty"`
	// Add a HTTP header key/value pair. Multiple headers can be set.
	Header map[string]string `json:"header,omitempty"`
	// Log the response payload within the Fluent Bit log.
	LogResponsePayload *bool `json:"logResponsePayload,omitempty"`
	// This allows you to add custom labels to all metrics exposed through the OpenTelemetry exporter. You may have multiple of these fields.
	AddLabel map[string]string `json:"addLabel,omitempty"`
	// If true, remaining unmatched keys are added as attributes.
	LogsBodyKeyAttributes *bool `json:"logsBodyKeyAttributes,omitempty"`
	// The log body key to look up in the log events body/message. Sets the Body field of the opentelemtry logs data model.
	LogsBodyKey  string `json:"logsBodyKey,omitempty"`
	*plugins.TLS `json:"tls,omitempty"`
	// Include fluentbit networking options for this output-plugin
	*plugins.Networking `json:"networking,omitempty"`
	// Limit the maximum number of Chunks in the filesystem for the current output logical destination.
	TotalLimitSize string `json:"totalLimitSize,omitempty"`
}

// Name implement Section() method
func (*OpenTelemetry) Name() string {
	return "opentelemetry"
}

// Params implement Section() method
func (o *OpenTelemetry) Params(sl plugins.SecretLoader) (*params.KVs, error) {
	kvs := params.NewKVs()

	plugins.InsertKVString(kvs, "host", o.Host)
	plugins.InsertKVField(kvs, "port", o.Port)

	if err := plugins.InsertKVSecret(kvs, "http_user", o.HTTPUser, sl); err != nil {
		return nil, err
	}
	if err := plugins.InsertKVSecret(kvs, "http_passwd", o.HTTPPasswd, sl); err != nil {
		return nil, err
	}

	plugins.InsertKVString(kvs, "proxy", o.Proxy)
	plugins.InsertKVString(kvs, "metrics_uri", o.MetricsUri)
	plugins.InsertKVString(kvs, "logs_uri", o.LogsUri)
	plugins.InsertKVString(kvs, "traces_uri", o.TracesUri)

	kvs.InsertStringMap(o.Header, func(k, v string) (string, string) {
		return header, fmt.Sprintf(" %s    %s", k, v)
	})

	plugins.InsertKVField(kvs, "log_response_payload", o.LogResponsePayload)

	kvs.InsertStringMap(o.AddLabel, func(k, v string) (string, string) {
		return addLabel, fmt.Sprintf(" %s    %s", k, v)
	})

	plugins.InsertKVField(kvs, "logs_body_key_attributes", o.LogsBodyKeyAttributes)
	plugins.InsertKVString(kvs, "logs_body_key", o.LogsBodyKey)

	if o.TLS != nil {
		tls, err := o.TLS.Params(sl)
		if err != nil {
			return nil, err
		}
		kvs.Merge(tls)
	}
	if o.Networking != nil {
		net, err := o.Networking.Params(sl)
		if err != nil {
			return nil, err
		}
		kvs.Merge(net)
	}

	plugins.InsertKVString(kvs, "storage.total_limit_size", o.TotalLimitSize)

	return kvs, nil
}
