# NodeExporterMetrics

The NodeExporterMetrics input plugin, which based on Prometheus Node Exporter to collect system / host level metrics.


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
