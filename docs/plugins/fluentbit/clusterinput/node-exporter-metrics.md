# Node Exporter Metrics

A plugin based on Prometheus Node Exporter to collect system / host level metrics.

>Note: Metrics collected with Node Exporter Metrics flow through a separate pipeline from logs and current filters do not operate on top of metrics.

This plugin is currently only supported on Linux based operating systems


| Field          | Description                                                                           | Scheme |
|----------------|---------------------------------------------------------------------------------------|--------|
| tag            | Tag name associated to all records comming from this plugin.                          | string |
| scrapeInterval | The rate at which metrics are collected from the host operating system.               | *int32 |
| path.procfs    | The mount point used to collect process information and metrics, default is `/proc/`. | string |
| path.sysfs     | The path in the filesystem used to collect system metrics, default is `/sys/`.        | string |


Important Note: The input plugin of node exporter metrics will collect system / host level metrics from specified path,
so we should mount those hostpath to containers. For example, This plugin will mount `/proc/` to collect process information and metrics 
and `/sys/` to collect system metrics by default if we don't specify path, thus we can add it to `values.yaml` in charts, like:

```yaml
fluentbit:
  volumes:
    - name: node-exporter-metrics-proc
      hostPath:
        path: /proc/
    - name: node-exporter-metrics-sys
      hostPath:
        path: /sys/
  volumesMounts:
    - mountPath: /host/sys
      mountPropagation: HostToContainer
      name: node-exporter-metrics-proc
      readOnly: true
    - mountPath: /host/proc
      mountPropagation: HostToContainer
      name: node-exporter-metrics-sys
      readOnly: true
```
