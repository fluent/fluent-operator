package input

import (
	"fmt"

	"github.com/fluent/fluent-operator/v3/apis/fluentd/v1alpha1/plugins/params"
)

// Prometheus is the in_prometheus plugin from fluent-plugin-prometheus. It serves the metrics
// collected by the prometheus plugins over HTTP so that Prometheus can scrape them.
type Prometheus struct {
	// The bind address to listen to. Defaults to 127.0.0.1, so set it to 0.0.0.0 to scrape from outside the pod.
	Bind *string `json:"bind,omitempty"`
	// The port to listen to. Defaults to 24231. Use 2021 to serve the metrics on the "metrics" port
	// that the operator already exposes on the Fluentd pods.
	Port *int64 `json:"port,omitempty"`
	// The HTTP path of the metrics endpoint. Defaults to /metrics.
	MetricsPath *string `json:"metricsPath,omitempty"`
	// The HTTP path of the endpoint that aggregates the metrics of all workers. Defaults to /aggregated_metrics.
	AggregatedMetricsPath *string `json:"aggregatedMetricsPath,omitempty"`
	// The encoding of the exposed metrics.
	// +kubebuilder:validation:Enum:=identity;gzip
	ContentEncoding *string `json:"contentEncoding,omitempty"`
}

// PrometheusMonitor is the in_prometheus_monitor plugin from fluent-plugin-prometheus.
// It collects Fluentd's buffer and retry metrics, such as fluentd_status_buffer_queue_length.
type PrometheusMonitor struct {
	// The interval in seconds between metric updates. Defaults to 5.
	Interval *uint32 `json:"interval,omitempty"`
	// Additional labels added to the metrics.
	Labels map[string]string `json:"labels,omitempty"`
}

// PrometheusOutputMonitor is the in_prometheus_output_monitor plugin from fluent-plugin-prometheus.
// It collects metrics about the output plugins, such as fluentd_output_status_retry_count.
type PrometheusOutputMonitor struct {
	// The interval in seconds between metric updates. Defaults to 5.
	Interval *uint32 `json:"interval,omitempty"`
	// Whether to expose the counters as gauges. Defaults to true.
	GaugeAll *bool `json:"gaugeAll,omitempty"`
	// Additional labels added to the metrics.
	Labels map[string]string `json:"labels,omitempty"`
}

func (i *Input) prometheusPlugin(parent *params.PluginStore) *params.PluginStore {
	model := i.Prometheus
	if model.Bind != nil {
		parent.InsertPairs("bind", *model.Bind)
	}
	if model.Port != nil {
		parent.InsertPairs("port", fmt.Sprint(*model.Port))
	}
	if model.MetricsPath != nil {
		parent.InsertPairs("metrics_path", *model.MetricsPath)
	}
	if model.AggregatedMetricsPath != nil {
		parent.InsertPairs("aggregated_metrics_path", *model.AggregatedMetricsPath)
	}
	if model.ContentEncoding != nil {
		parent.InsertPairs("content_encoding", *model.ContentEncoding)
	}
	return parent
}

func (i *Input) prometheusMonitorPlugin(parent *params.PluginStore) *params.PluginStore {
	model := i.PrometheusMonitor
	if model.Interval != nil {
		parent.InsertPairs("interval", fmt.Sprint(*model.Interval))
	}
	insertPrometheusLabels(parent, model.Labels)
	return parent
}

func (i *Input) prometheusOutputMonitorPlugin(parent *params.PluginStore) *params.PluginStore {
	model := i.PrometheusOutputMonitor
	if model.Interval != nil {
		parent.InsertPairs("interval", fmt.Sprint(*model.Interval))
	}
	if model.GaugeAll != nil {
		parent.InsertPairs("gauge_all", fmt.Sprint(*model.GaugeAll))
	}
	insertPrometheusLabels(parent, model.Labels)
	return parent
}

// insertPrometheusLabels renders the labels as a <labels> section.
func insertPrometheusLabels(parent *params.PluginStore, labels map[string]string) {
	if len(labels) == 0 {
		return
	}
	child := params.NewPluginStore("labels")
	for k, v := range labels {
		child.InsertPairs(k, v)
	}
	parent.InsertChilds(child)
}
