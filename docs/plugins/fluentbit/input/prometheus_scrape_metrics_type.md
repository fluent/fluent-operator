# PrometheusScrapeMetrics

Fluent Bit 1.9 includes additional metrics features to allow you to collect both logs and metrics with the same collector. <br /> The initial release of the Prometheus Scrape metric allows you to collect metrics from a Prometheus-based <br /> endpoint at a set interval. These metrics can be routed to metric supported endpoints such as Prometheus Exporter, InfluxDB, or Prometheus Remote Write. <br /> **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/inputs/prometheus-scrape-metrics**


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| tag | Tag name associated to all records comming from this plugin | string |
| host | The host of the prometheus metric endpoint that you want to scrape | string |
| port | The port of the promethes metric endpoint that you want to scrape | *int32 |
| scrapeInterval | The interval to scrape metrics, default: 10s | string |
| metricsPath | The metrics URI endpoint, that must start with a forward slash, deflaut: /metrics | string |
