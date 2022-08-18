package output

import (
	"fmt"

	"github.com/fluent/fluent-operator/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/apis/fluentbit/v1alpha2/plugins/params"
	"github.com/fluent/fluent-operator/pkg/utils"
)

// +kubebuilder:object:generate:=true

// An output plugin to submit Prometheus Metrics using the remote write protocol
type PrometheusRemoteWrite struct {
	// IP address or hostname of the target HTTP Server, default: 127.0.0.1
	Host string `json:"host"`
	// Basic Auth Username
	HTTPUser *plugins.Secret `json:"httpUser,omitempty"`
	// Basic Auth Password.
	// Requires HTTP_user to be se
	HTTPPasswd *plugins.Secret `json:"httpPassword,omitempty"`
	// TCP port of the target HTTP Serveri, default:80
	// +kubebuilder:validation:Minimum:=1
	// +kubebuilder:validation:Maximum:=65535
	Port *int32 `json:"port,omitempty"`
	// Specify an HTTP Proxy. The expected format of this value is http://HOST:PORT.
	Proxy string `json:"proxy,omitempty"`
	//Specify an optional HTTP URI for the target web server, e.g: /something ,default: /
	Uri string `json:"uri,omitempty"`
	//Add a HTTP header key/value pair. Multiple headers can be set.
	Headers map[string]string `json:"headers,omitempty"`
	//Log the response payload within the Fluent Bit log,default: false
	LogResponsePayload *bool `json:"logResponsePayload,omitempty"`
	//This allows you to add custom labels to all metrics exposed through the prometheus exporter. You may have multiple of these fields
	AddLabels []string `json:"addLabels,omitempty"`

	//to do ?
	Match string `json:"match,omitempty"`

	*plugins.TLS `json:"tls,omitempty"`
}

// implement Section() method
func (_ *PrometheusRemoteWrite) Name() string {
	return "Prometheus_Remote_Write"
}

// implement Section() method
func (p *PrometheusRemoteWrite) Params(sl plugins.SecretLoader) (*params.KVs, error) {
	kvs := params.NewKVs()
	if p.Host != "" {
		kvs.Insert("host", p.Host)
	}
	if p.Port != nil {
		kvs.Insert("port", fmt.Sprint(*p.Port))
	}
	if p.HTTPUser != nil {
		u, err := sl.LoadSecret(*p.HTTPUser)
		if err != nil {
			return nil, err
		}
		kvs.Insert("http_user", u)
	}
	if p.HTTPPasswd != nil {
		pwd, err := sl.LoadSecret(*p.HTTPPasswd)
		if err != nil {
			return nil, err
		}
		kvs.Insert("http_passwd", pwd)
	}
	if p.Proxy != "" {
		kvs.Insert("proxy", p.Proxy)
	}
	if p.Uri != "" {
		kvs.Insert("uri", p.Uri)
	}

	kvs.InsertStringMap(p.Headers, func(k, v string) (string, string) {
		return "header", fmt.Sprintf(" %s    %s", k, v)
	})

	if p.LogResponsePayload != nil {
		kvs.Insert("log_response_payload", fmt.Sprint(*p.LogResponsePayload))
	}
	if p.AddLabels != nil && len(p.AddLabels) > 0 {
		kvs.Insert("add_label", utils.ConcatString(p.AddLabels, ","))
	}
	if p.Match != "" {
		kvs.Insert("match", p.Match)
	}
	if p.TLS != nil {
		tls, err := p.TLS.Params(sl)
		if err != nil {
			return nil, err
		}
		kvs.Merge(tls)
	}
	return kvs, nil
}
