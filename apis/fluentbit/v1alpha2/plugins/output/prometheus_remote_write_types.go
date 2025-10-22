package output

import (
	"fmt"

	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins/params"
)

// +kubebuilder:object:generate:=true

// An output plugin to submit Prometheus Metrics using the remote write protocol. <br />
// The prometheus remote write plugin allows you to take metrics from Fluent Bit and submit them to a Prometheus server through the remote write mechanism. <br />
// **Important Note: The prometheus exporter only works with metric plugins, such as Node Exporter Metrics** <br />
// **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/outputs/prometheus-remote-write**
type PrometheusRemoteWrite struct {
	// IP address or hostname of the target HTTP Server, default: 127.0.0.1
	Host string `json:"host"`
	// Basic Auth Username
	HTTPUser *plugins.Secret `json:"httpUser,omitempty"`
	// Basic Auth Password.
	// Requires HTTP_user to be se
	HTTPPasswd *plugins.Secret `json:"httpPasswd,omitempty"`
	// TCP port of the target HTTP Serveri, default:80
	// +kubebuilder:validation:Minimum:=1
	// +kubebuilder:validation:Maximum:=65535
	Port *int32 `json:"port,omitempty"`
	// Specify an HTTP Proxy. The expected format of this value is http://HOST:PORT.
	Proxy string `json:"proxy,omitempty"`
	// Specify an optional HTTP URI for the target web server, e.g: /something ,default: /
	URI string `json:"uri,omitempty"`
	// Add a HTTP header key/value pair. Multiple headers can be set.
	Headers map[string]string `json:"headers,omitempty"`
	// Log the response payload within the Fluent Bit log,default: false
	LogResponsePayload *bool `json:"logResponsePayload,omitempty"`
	// This allows you to add custom labels to all metrics exposed through the prometheus exporter. You may have multiple of these fields
	AddLabels map[string]string `json:"addLabels,omitempty"`
	// Enables dedicated thread(s) for this output. Default value is set since version 1.8.13. For previous versions is 0,default : 2
	Workers *int32 `json:"workers,omitempty"`

	*plugins.TLS `json:"tls,omitempty"`
	// Include fluentbit networking options for this output-plugin
	*plugins.Networking `json:"networking,omitempty"`
}

// implement Section() method
func (*PrometheusRemoteWrite) Name() string {
	return "prometheus_remote_write"
}

// implement Section() method
func (p *PrometheusRemoteWrite) Params(sl plugins.SecretLoader) (*params.KVs, error) {
	kvs := params.NewKVs()

	if err := plugins.InsertKVSecret(kvs, "http_user", p.HTTPUser, sl); err != nil {
		return nil, err
	}
	if err := plugins.InsertKVSecret(kvs, "http_passwd", p.HTTPPasswd, sl); err != nil {
		return nil, err
	}

	plugins.InsertKVString(kvs, "host", p.Host)
	plugins.InsertKVField(kvs, "port", p.Port)
	plugins.InsertKVString(kvs, "proxy", p.Proxy)
	plugins.InsertKVString(kvs, "uri", p.URI)

	kvs.InsertStringMap(p.Headers, func(k, v string) (string, string) {
		return header, fmt.Sprintf(" %s    %s", k, v)
	})

	plugins.InsertKVField(kvs, "log_response_payload", p.LogResponsePayload)

	kvs.InsertStringMap(p.AddLabels, func(k, v string) (string, string) {
		return addLabel, fmt.Sprintf(" %s    %s", k, v)
	})

	plugins.InsertKVField(kvs, "workers", p.Workers)

	if p.TLS != nil {
		tls, err := p.TLS.Params(sl)
		if err != nil {
			return nil, err
		}
		kvs.Merge(tls)
	}
	if p.Networking != nil {
		net, err := p.Networking.Params(sl)
		if err != nil {
			return nil, err
		}
		kvs.Merge(net)
	}

	return kvs, nil
}
