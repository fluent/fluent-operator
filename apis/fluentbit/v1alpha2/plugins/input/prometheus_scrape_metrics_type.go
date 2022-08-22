package input

import (
	"fmt"

	"github.com/fluent/fluent-operator/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/apis/fluentbit/v1alpha2/plugins/params"
	"strings"
)

// +kubebuilder:object:generate:=true

/*  Fluent Bit 1.9 includes additional metrics features to allow you to collect both logs and metrics with the same collector.
 *  The initial release of the Prometheus Scrape metric allows you to collect metrics from a Prometheus-based
 *  endpoint at a set interval. These metrics can be routed to metric supported endpoints such as Prometheus Exporter, InfluxDB, or Prometheus Remote Write
 */
type PrometheusScrapeMetrics struct {
	// Tag name associated to all records comming from this plugin
	Tag string `json:"tag,omitempty"`
	// The host of the prometheus metric endpoint that you want to scrape
	Host string `json:"host,omitempty"`
	// The port of the promethes metric endpoint that you want to scrape
	// +kubebuilder:validation:Minimum:=1
	// +kubebuilder:validation:Maximum:=65535
	Port *int32 `json:"port,omitempty"`
	// The interval to scrape metrics, default: 10s
	ScrapeInterval *int32 `json:"scrapeInterval,omitempty"`
	// The metrics URI endpoint, that must start with a forward slash, deflaut: /metrics
	MetricsPath string `json:"metricsPath,omitempty"`
}

func (_ *PrometheusScrapeMetrics) Name() string {
	return "prometheus_scrape_metrics"
}

// Params implement Section() method
func (p *PrometheusScrapeMetrics) Params(_ plugins.SecretLoader) (*params.KVs, error) {
	kvs := params.NewKVs()
	if p.Tag != "" {
		kvs.Insert("tag", p.Tag)
	}
        host := strings.ToLower(p.Host)
	if host == "" || host == "localhost" || host == "127.0.0.1" || host == "host" {
		kvs.Insert("host", "${HOST_IP}")
	} else if strings.ToLower(p.Host) == "host" || strings.ToLower(p.Host) == "hostnode" {
		kvs.Insert("host", "${HOST_IP}")
	} else {
		kvs.Insert("host", p.Host)
	}
	if p.Port != nil {
		kvs.Insert("port", fmt.Sprint(*p.Port))
	}
	if p.ScrapeInterval != nil {
		kvs.Insert("scrape_interval", fmt.Sprint(*p.ScrapeInterval))
	}
	if p.MetricsPath != "" {
		kvs.Insert("metrics_path", p.MetricsPath)
	}
	return kvs, nil
}
