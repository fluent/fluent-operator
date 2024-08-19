# LogToMetrics

The Log To Metrics Filter plugin allows you to generate log-derived metrics. <br /> **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/filters/log_to_metrics**


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| tag | Defines the tag for the generated metrics record | string |
| regex | Optional filter for records in which the content of KEY matches the regular expression. Value Format: FIELD REGEX | []string |
| exclude | Optional filter for records in which the content of KEY does not matches the regular expression. Value Format: FIELD REGEX | []string |
| metricMode | Defines the mode for the metric. Valid values are [counter, gauge or histogram] | string |
| metricName | Sets the name of the metric. | string |
| metricNamespace | Namespace of the metric | string |
| metricSubsystem | Sets a sub-system for the metric. | string |
| metricDescription | Sets a help text for the metric. | string |
| bucket | Defines a bucket for histogram | []string |
| addLabel | Add a custom label NAME and set the value to the value of KEY | []string |
| labelField | Includes a record field as label dimension in the metric. | []string |
| valueField | Specify the record field that holds a numerical value | string |
| kubernetesMode | If enabled, it will automatically put pod_id, pod_name, namespace_name, docker_id and container_name into the metric as labels. This option is intended to be used in combination with the kubernetes filter plugin. | *bool |
| emitterName | Name of the emitter (advanced users) | string |
| emitterMemBufLimit | set a buffer limit to restrict memory usage of metrics emitter | string |
| discardLogs | Flag that defines if logs should be discarded after processing. This applies for all logs, no matter if they have emitted metrics or not. | *bool |
