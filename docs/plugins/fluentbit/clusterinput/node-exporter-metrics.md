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

The node exporter metrics input plugin will collect host node's metrics from specific host paths, so you should mount those host paths to fluentbit container paths. For example:

- Host node's `/proc/` should be mounted to container's `/host/proc` to collect process information and metrics.
- Host node's `/sys/` should be mounted to container's `/host/sys` to collect system metrics.

To do this, you'll need to uncomment the following content in helm chart's `values.yaml`:

```yaml
fluentbit:
  volumes:
    - name: hostProc
      hostPath:
        path: /proc/
    - name: hostSys
      hostPath:
        path: /sys/
  volumesMounts:
    - mountPath: /host/sys
      mountPropagation: HostToContainer
      name: hostSys 
      readOnly: true
    - mountPath: /host/proc
      mountPropagation: HostToContainer
      name: hostProc 
      readOnly: true
  input:
    nodeExporterMetrics:
      tag: node_metrics
      scrapeInterval: 15s
      path:
        procfs: /host/proc
        sysfs: /host/sys
```
