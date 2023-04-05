# DataDog

DataDog output plugin allows you to ingest your logs into Datadog. <br /> **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/outputs/datadog**


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| host | Host is the Datadog server where you are sending your logs. | string |
| tls | TLS controls whether to use end-to-end security communications security protocol. Datadog recommends setting this to on. | *bool |
| compress | Compress  the payload in GZIP format. Datadog supports and recommends setting this to gzip. | string |
| apikey | Your Datadog API key. | string |
| proxy | Specify an HTTP Proxy. | string |
| provider | To activate the remapping, specify configuration flag provider. | string |
| json_date_key | Date key name for output. | string |
| include_tag_key | If enabled, a tag is appended to output. The key name is used tag_key property. | *bool |
| tag_key | The key name of tag. If include_tag_key is false, This property is ignored. | string |
| dd_service | The human readable name for your service generating the logs. | string |
| dd_source | A human readable name for the underlying technology of your service. | string |
| dd_tags | The tags you want to assign to your logs in Datadog. | string |
| dd_message_key | By default, the plugin searches for the key 'log' and remap the value to the key 'message'. If the property is set, the plugin will search the property name key. | string |
