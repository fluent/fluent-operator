package input

import (
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins/params"
)

// +kubebuilder:object:generate:=true

// NGINX Exporter Metrics input plugin scrapes metrics from the NGINX stub status handler. <br />
// **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/inputs/nginx**
type Nginx struct {
	// Name of the target host or IP address to check, default: localhost
	Host string `json:"host,omitempty"`
	// Port of the target nginx service to connect to, default: 80
	// +kubebuilder:validation:Minimum:=1
	// +kubebuilder:validation:Maximum:=65535
	Port *int32 `json:"port,omitempty"`
	// The URL of the Stub Status Handler,default: /status
	StatusURL string `json:"statusURL,omitempty"`
	// Turn on NGINX plus mode,default: true
	NginxPlus *bool `json:"nginxPlus,omitempty"`
}

func (*Nginx) Name() string {
	return "nginx_metrics"
}

// implement Section() method
func (n *Nginx) Params(_ plugins.SecretLoader) (*params.KVs, error) {
	kvs := params.NewKVs()

	plugins.InsertKVString(kvs, "Host", n.Host)
	plugins.InsertKVString(kvs, "Status_URL", n.StatusURL)

	plugins.InsertKVField(kvs, "Port", n.Port)
	plugins.InsertKVField(kvs, "Nginx_Plus", n.NginxPlus)

	return kvs, nil
}
