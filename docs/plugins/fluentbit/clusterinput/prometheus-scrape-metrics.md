# Prometheus Scrape Metrics

Fluent Bit 1.9 includes additional metrics features to allow you to collect both logs and metrics with the same collector.
The initial release of the Prometheus Scrape metric allows you to collect metrics from a Prometheus-based endpoint at a set interval. These metrics can be routed to metric supported endpoints such as Prometheus Exporter, InfluxDB, or Prometheus Remote Write


| Field           | Description                                                  | Scheme |
| --------------- | ------------------------------------------------------------ | ------ |
| tag             | Tag name associated to all records comming from this plugin. | string |
| host            | The host of the prometheus metric endpoint that you want to scrape. | string |
| port            | The port of the promethes metric endpoint that you want to scrape. | *int32 |
| scrape_interval | The interval to scrape metrics, default is 10s.              | string |
| metrics_path    | The metrics URI endpoint, that must start with a forward slash, default is /metrics. | string |
