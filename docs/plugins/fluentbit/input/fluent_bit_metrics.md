# FluentbitMetrics




| Field | Description | Scheme |
| ----- | ----------- | ------ |
| tag |  | string |
| scrapeInterval | The rate at which metrics are collected from the host operating system. default is 2 seconds. | string |
| scrapeOnStart | Scrape metrics upon start, useful to avoid waiting for 'scrape_interval' for the first round of metrics. | *bool |
