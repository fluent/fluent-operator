# InfluxDB

The influxdb output plugin, allows to flush your records into a InfluxDB time series database. <br /> **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/outputs/influxdb**


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| host | IP address or hostname of the target InfluxDB service. | string |
| port | TCP port of the target InfluxDB service. | *int32 |
| database | InfluxDB database name where records will be inserted. | string |
| bucket | InfluxDB bucket name where records will be inserted - if specified, database is ignored and v2 of API is used | string |
| org | InfluxDB organization name where the bucket is (v2 only) | string |
| sequenceTag | The name of the tag whose value is incremented for the consecutive simultaneous events. | string |
| httpUser | Optional username for HTTP Basic Authentication | *[plugins.Secret](../secret.md) |
| httpPassword | Password for user defined in HTTP_User | *[plugins.Secret](../secret.md) |
| httpToken | Authentication token used with InfluxDB v2 - if specified, both HTTPUser and HTTPPasswd are ignored | *[plugins.Secret](../secret.md) |
| tagKeys | List of keys that needs to be tagged | []string |
| autoTags | Automatically tag keys where value is string. | *bool |
| tagsListEnabled | Dynamically tag keys which are in the string array at Tags_List_Key key. | *bool |
| tagListKey | Key of the string array optionally contained within each log record that contains tag keys for that record | string |
| tls |  | *[plugins.TLS](../tls.md) |
