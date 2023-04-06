# Splunk

Splunk output plugin allows to ingest your records into a Splunk Enterprise service through the HTTP Event Collector (HEC) interface. <br /> **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/outputs/splunk**


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| host | IP address or hostname of the target OpenSearch instance, default `127.0.0.1` | string |
| port | TCP port of the target Splunk instance, default `8088` | *int32 |
| splunkToken | Specify the Authentication Token for the HTTP Event Collector interface. | *[plugins.Secret](../secret.md) |
| httpBufferSize | Buffer size used to receive Splunk HTTP responses: Default `2M` | string |
| compress | Set payload compression mechanism. The only available option is gzip. | string |
| channel | Specify X-Splunk-Request-Channel Header for the HTTP Event Collector interface. | string |
| httpUser | Optional username credential for access | *[plugins.Secret](../secret.md) |
| httpPassword | Password for user defined in HTTP_User | *[plugins.Secret](../secret.md) |
| httpDebugBadRequest | If the HTTP server response code is 400 (bad request) and this flag is enabled, it will print the full HTTP request and response to the stdout interface. This feature is available for debugging purposes. | *bool |
| splunkSendRaw | When enabled, the record keys and values are set in the top level of the map instead of under the event key. Refer to the Sending Raw Events section from the docs more details to make this option work properly. | *bool |
| eventKey | Specify the key name that will be used to send a single value as part of the record. | string |
| eventHost | Specify the key name that contains the host value. This option allows a record accessors pattern. | string |
| eventSource | Set the source value to assign to the event data. | string |
| eventSourcetype | Set the sourcetype value to assign to the event data. | string |
| eventSourcetypeKey | Set a record key that will populate 'sourcetype'. If the key is found, it will have precedence over the value set in event_sourcetype. | string |
| eventIndex | The name of the index by which the event data is to be indexed. | string |
| eventIndexKey | Set a record key that will populate the index field. If the key is found, it will have precedence over the value set in event_index. | string |
| eventFields | Set event fields for the record. This option is an array and the format is \"key_name record_accessor_pattern\". | []string |
| Workers | Enables dedicated thread(s) for this output. Default value `2` is set since version 1.8.13. For previous versions is 0. | *int32 |
| tls |  | *[plugins.TLS](../tls.md) |
