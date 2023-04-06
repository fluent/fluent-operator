# NodeExporterMetrics

A plugin based on Prometheus Node Exporter to collect system / host level metrics. <br /> **Note: Metrics collected with Node Exporter Metrics flow through a separate pipeline from logs and current filters do not operate on top of metrics.** <br /> This plugin is currently only supported on Linux based operating systems. <br /> **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/inputs/node-exporter-metrics**


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| tag | Tag name associated to all records comming from this plugin. | string |
| scrapeInterval | The rate at which metrics are collected from the host operating system, default is 5 seconds. | string |
| path |  | *Path |
# Path




| Field | Description | Scheme |
| ----- | ----------- | ------ |
| procfs | The mount point used to collect process information and metrics. | string |
| sysfs | The path in the filesystem used to collect system metrics. | string |
