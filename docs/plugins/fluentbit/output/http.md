# HTTP

The http output plugin allows to flush your records into a HTTP endpoint. <br /> For now the functionality is pretty basic and it issues a POST request with the data records in MessagePack (or JSON) format. <br /> **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/outputs/http**


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| host | IP address or hostname of the target HTTP Server | string |
| httpUser | Basic Auth Username | *[plugins.Secret](../secret.md) |
| httpPassword | Basic Auth Password. Requires HTTP_User to be set | *[plugins.Secret](../secret.md) |
| port | TCP port of the target HTTP Server | *int32 |
| proxy | Specify an HTTP Proxy. The expected format of this value is http://host:port. Note that https is not supported yet. | string |
| uri | Specify an optional HTTP URI for the target web server, e.g: /something | string |
| compress | Set payload compression mechanism. Option available is 'gzip' | string |
| format | Specify the data format to be used in the HTTP request body, by default it uses msgpack. Other supported formats are json, json_stream and json_lines and gelf. | string |
| allowDuplicatedHeaders | Specify if duplicated headers are allowed. If a duplicated header is found, the latest key/value set is preserved. | *bool |
| headerTag | Specify an optional HTTP header field for the original message tag. | string |
| headers | Add a HTTP header key/value pair. Multiple headers can be set. | map[string]string |
| jsonDateKey | Specify the name of the time key in the output record. To disable the time key just set the value to false. | string |
| jsonDateFormat | Specify the format of the date. Supported formats are double, epoch and iso8601 (eg: 2018-05-30T09:39:52.000681Z) | string |
| gelfTimestampKey | Specify the key to use for timestamp in gelf format | string |
| gelfHostKey | Specify the key to use for the host in gelf format | string |
| gelfShortMessageKey | Specify the key to use as the short message in gelf format | string |
| gelfFullMessageKey | Specify the key to use for the full message in gelf format | string |
| gelfLevelKey | Specify the key to use for the level in gelf format | string |
| tls | HTTP output plugin supports TTL/SSL, for more details about the properties available and general configuration, please refer to the TLS/SSL section. | *[plugins.TLS](../tls.md) |
