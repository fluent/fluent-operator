package filter

import (
	"fmt"

	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins/params"
)

// +kubebuilder:object:generate:=true

// The Log To Metrics Filter plugin allows you to generate log-derived metrics. <br />
// **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/filters/log_to_metrics**
type LogToMetrics struct {
	plugins.CommonParams `json:",inline"`
	// Defines the tag for the generated metrics record
	Tag string `json:"tag,omitempty"`
	// Optional filter for records in which the content of KEY matches the regular expression.
	// Value Format: FIELD REGEX
	Regex []string `json:"regex,omitempty"`
	// Optional filter for records in which the content of KEY does not matches the regular expression.
	// Value Format: FIELD REGEX
	Exclude []string `json:"exclude,omitempty"`
	// Defines the mode for the metric. Valid values are [counter, gauge or histogram]
	MetricMode string `json:"metricMode,omitempty"`
	// Sets the name of the metric.
	MetricName string `json:"metricName,omitempty"`
	// Namespace of the metric
	MetricNamespace string `json:"metricNamespace,omitempty"`
	// Sets a sub-system for the metric.
	MetricSubsystem string `json:"metricSubsystem,omitempty"`
	// Sets a help text for the metric.
	MetricDescription string `json:"metricDescription,omitempty"`
	// Defines a bucket for histogram
	Bucket []string `json:"bucket,omitempty"`
	// Add a custom label NAME and set the value to the value of KEY
	AddLabel []string `json:"addLabel,omitempty"`
	// Includes a record field as label dimension in the metric.
	LabelField []string `json:"labelField,omitempty"`
	// Specify the record field that holds a numerical value
	ValueField string `json:"valueField,omitempty"`
	// If enabled, it will automatically put pod_id, pod_name, namespace_name, docker_id and container_name
	// into the metric as labels. This option is intended to be used in combination with the kubernetes filter plugin.
	KubernetesMode *bool `json:"kubernetesMode,omitempty"`
	// Name of the emitter (advanced users)
	EmitterName string `json:"emitterName,omitempty"`
	// set a buffer limit to restrict memory usage of metrics emitter
	EmitterMemBufLimit string `json:"emitterMemBufLimit,omitempty"`
	// Flag that defines if logs should be discarded after processing. This applies
	// for all logs, no matter if they have emitted metrics or not.
	DiscardLogs *bool `json:"discardLogs,omitempty"`
}

func (_ *LogToMetrics) Name() string {
	return "log_to_metrics"
}

func (l *LogToMetrics) Params(_ plugins.SecretLoader) (*params.KVs, error) {
	kvs := params.NewKVs()
	err := l.AddCommonParams(kvs)
	if err != nil {
		return kvs, err
	}
	if l.Tag != "" {
		kvs.Insert("Tag", l.Tag)
	}
	for _, reg := range l.Regex {
		kvs.Insert("Regex", reg)
	}
	for _, ex := range l.Exclude {
		kvs.Insert("Exclude", ex)
	}
	if l.MetricMode != "" {
		kvs.Insert("Metric_mode", l.MetricMode)
	}
	if l.MetricName != "" {
		kvs.Insert("Metric_name", l.MetricName)
	}
	if l.MetricNamespace != "" {
		kvs.Insert("Metric_namespace", l.MetricNamespace)
	}
	if l.MetricSubsystem != "" {
		kvs.Insert("Metric_subsystem", l.MetricSubsystem)
	}
	if l.MetricDescription != "" {
		kvs.Insert("Metric_description", l.MetricDescription)
	}
	for _, b := range l.Bucket {
		kvs.Insert("Bucket", b)
	}
	for _, al := range l.AddLabel {
		kvs.Insert("Add_label", al)
	}
	for _, lf := range l.LabelField {
		kvs.Insert("Label_field", lf)
	}
	if l.ValueField != "" {
		kvs.Insert("Value_field", l.ValueField)
	}
	if l.KubernetesMode != nil {
		kvs.Insert("Kubernetes_mode", fmt.Sprintf("%t", *l.KubernetesMode))
	}
	if l.EmitterName != "" {
		kvs.Insert("Emitter_Name", l.EmitterName)
	}
	if l.EmitterMemBufLimit != "" {
		kvs.Insert("Emitter_Mem_Buf_Limit", l.EmitterMemBufLimit)
	}
	if l.DiscardLogs != nil {
		kvs.Insert("Discard_logs", fmt.Sprintf("%t", *l.DiscardLogs))
	}
	return kvs, nil
}
