# TCP

The tcp output plugin allows to send records to a remote TCP server. <br /> The payload can be formatted in different ways as required. <br /> **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/outputs/tcp-and-tls**


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| host | Target host where Fluent-Bit or Fluentd are listening for Forward messages. | string |
| port | TCP Port of the target service. | *int32 |
| format | Specify the data format to be printed. Supported formats are msgpack json, json_lines and json_stream. | string |
| jsonDateKey | TSpecify the name of the time key in the output record. To disable the time key just set the value to false. | string |
| jsonDateFormat | Specify the format of the date. Supported formats are double, epoch and iso8601 (eg: 2018-05-30T09:39:52.000681Z) | string |
| tls |  | *[plugins.TLS](../tls.md) |
