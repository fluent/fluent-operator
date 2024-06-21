package output

import (
	"fmt"

	"github.com/fluent/fluent-operator/v2/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v2/apis/fluentbit/v1alpha2/plugins/params"
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
	AddLabel     map[string]string `json:"addLabel,omitempty"`
	*plugins.TLS `json:"tls,omitempty"`
	// Include fluentbit networking options for this output-plugin
	*plugins.Networking `json:"networking,omitempty"`
}

// Name implement Section() method
func (_ *OpenTelemetry) Name() string {
	return "opentelemetry"
}

// Params implement Section() method
func (o *OpenTelemetry) Params(sl plugins.SecretLoader) (*params.KVs, error) {
	kvs := params.NewKVs()
	if o.Host != "" {
		kvs.Insert("host", o.Host)
	}
	if o.Port != nil {
		kvs.Insert("port", fmt.Sprint(*o.Port))
	}
	if o.HTTPUser != nil {
		u, err := sl.LoadSecret(*o.HTTPUser)
		if err != nil {
			return nil, err
		}
		kvs.Insert("http_user", u)
	}
	if o.HTTPPasswd != nil {
		pwd, err := sl.LoadSecret(*o.HTTPPasswd)
		if err != nil {
			return nil, err
		}
		kvs.Insert("http_passwd", pwd)
	}
	if o.Proxy != "" {
		kvs.Insert("proxy", o.Proxy)
	}
	if o.MetricsUri != "" {
		kvs.Insert("metrics_uri", o.MetricsUri)
	}
	if o.LogsUri != "" {
		kvs.Insert("logs_uri", o.LogsUri)
	}
	if o.TracesUri != "" {
		kvs.Insert("traces_uri", o.TracesUri)
	}
	kvs.InsertStringMap(o.Header, func(k, v string) (string, string) {
		return "header", fmt.Sprintf(" %s    %s", k, v)
	})
	if o.LogResponsePayload != nil {
		kvs.Insert("log_response_payload", fmt.Sprint(*o.LogResponsePayload))
	}
	kvs.InsertStringMap(o.AddLabel, func(k, v string) (string, string) {
		return "add_label", fmt.Sprintf(" %s    %s", k, v)
	})
	if o.TLS != nil {
		tls, err := o.TLS.Params(sl)
		if err != nil {
			return nil, err
		}
		kvs.Merge(tls)
	}
	return kvs, nil
}
