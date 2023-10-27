# MonitorAgent

The in_monitor_agent Input plugin exports Fluentd's internal metrics via REST API.


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| port | The port to listen to. | *int64 |
| bind | The bind address to listen to. | *string |
| tag | If you set this parameter, this plugin emits metrics as records. | *string |
| emitInterval | The interval time between event emits. This will be used when \"tag\" is configured. | *int64 |
| includeConfig | You can set this option to false to remove the config field from the response. | *bool |
| includeRetry | You can set this option to false to remove the retry field from the response. | *bool |
