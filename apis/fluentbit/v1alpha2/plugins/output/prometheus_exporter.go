package output

import (
	"fmt"

	"github.com/fluent/fluent-operator/v2/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v2/apis/fluentbit/v1alpha2/plugins/params"
)

// +kubebuilder:object:generate:=true

// PrometheusExporter An output plugin to expose Prometheus Metrics. <br />
// The prometheus exporter allows you to take metrics from Fluent Bit and expose them such that a Prometheus instance can scrape them. <br />
// **Important Note: The prometheus exporter only works with metric plugins, such as Node Exporter Metrics** <br />
// **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/outputs/prometheus-exporter**
type PrometheusExporter struct {
	// IP address or hostname of the target HTTP Server, default: 0.0.0.0
	Host string `json:"host"`
	// This is the port Fluent Bit will bind to when hosting prometheus metrics.
	// +kubebuilder:validation:Minimum:=1
	// +kubebuilder:validation:Maximum:=65535
	Port *int32 `json:"port,omitempty"`
	//This allows you to add custom labels to all metrics exposed through the prometheus exporter. You may have multiple of these fields
	AddLabels map[string]string `json:"addLabels,omitempty"`
}

// implement Section() method
func (_ *PrometheusExporter) Name() string {
	return "prometheus_exporter"
}

// implement Section() method
func (p *PrometheusExporter) Params(sl plugins.SecretLoader) (*params.KVs, error) {
	kvs := params.NewKVs()
	if p.Host != "" {
		kvs.Insert("host", p.Host)
	}
	if p.Port != nil {
		kvs.Insert("port", fmt.Sprint(*p.Port))
	}
	kvs.InsertStringMap(p.AddLabels, func(k, v string) (string, string) {
		return "add_label", fmt.Sprintf(" %s    %s", k, v)
	})
	return kvs, nil
}
