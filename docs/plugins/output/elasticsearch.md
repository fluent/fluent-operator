# Elasticsearch

The es output plugin, allows to ingest your records into a Elasticsearch database.


| Field | Description | Scheme | Default |
| ----- | ----------- | ------ | ----- |
| host | IP address or hostname of the target Elasticsearch instance | string | 127.0.0.1 |
| port | TCP port of the target Elasticsearch instance | *int32 | 9200 |
| path | Elasticsearch accepts new data on HTTP query path \"/_bulk\". But it is also possible to serve Elasticsearch behind a reverse proxy on a subpath. This option defines such path on the fluent-bit side. It simply adds a path prefix in the indexing HTTP POST URI. | string | Empty string |
| bufferSize | Specify the buffer size used to read the response from the Elasticsearch HTTP service. This option is useful for debugging purposes where is required to read full responses, note that response size grows depending of the number of records inserted. To set an unlimited amount of memory set this value to False, otherwise the value must be according to the Unit Size specification. | string | 4KB |
| pipeline | Newer versions of Elasticsearch allows to setup filters called pipelines. This option allows to define which pipeline the database should use. For performance reasons is strongly suggested to do parsing and filtering on Fluent Bit side, avoid pipelines. | string |  |
| httpUser | Optional username credential for Elastic X-Pack access | *[plugins.Secret](../secret.md) |  |
| httpPassword | Password for user defined in HTTP_User | *[plugins.Secret](../secret.md) |  |
| index | Index name | string | fluent-bit |
| type | Type name | string | _doc |
| logstashFormat | Enable Logstash format compatibility. This option takes a boolean value: True/False, On/Off | *bool | False |
| logstashPrefix | When Logstash_Format is enabled, the Index name is composed using a prefix and the date, e.g: If Logstash_Prefix is equals to 'mydata' your index will become 'mydata-YYYY.MM.DD'. The last string appended belongs to the date when the data is being generated. | string | logstash |
| logstashDateFormat | Time format (based on strftime) to generate the second part of the Index name. | string | %Y.%m.%d |
| timeKey | When Logstash_Format is enabled, each record will get a new timestamp field. The Time_Key property defines the name of that field. | string | @timestamp |
| timeKeyFormat | When Logstash_Format is enabled, this property defines the format of the timestamp. | string | %Y-%m-%dT%H:%M:%S |
| includeTagKey | When enabled, it append the Tag name to the record. | *bool | False |
| tagKey | When Include_Tag_Key is enabled, this property defines the key name for the tag. | string | _flb-key |
| generateID | When enabled, generate _id for outgoing records. This prevents duplicate records when retrying ES. | *bool | False |
| replaceDots | When enabled, replace field name dots with underscore, required by Elasticsearch 2.0-2.3. | *bool | False |
| traceOutput | When enabled print the elasticsearch API calls to stdout (for diag only) | *bool | False |
| traceError | When enabled print the elasticsearch API calls to stdout when elasticsearch returns an error | *bool | False |
| currentTimeIndex | Use current time for index generation instead of message record | *bool | False |
| logstashPrefixKey | Prefix keys with this string | string |  |
| tls |  | *[plugins.TLS](../tls.md) |  |
