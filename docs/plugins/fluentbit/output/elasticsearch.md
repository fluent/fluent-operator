# Elasticsearch

Elasticsearch is the es output plugin, allows to ingest your records into an Elasticsearch database. <br /> **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/outputs/elasticsearch**


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| host | IP address or hostname of the target Elasticsearch instance | string |
| port | TCP port of the target Elasticsearch instance | *int32 |
| path | Elasticsearch accepts new data on HTTP query path \"/_bulk\". But it is also possible to serve Elasticsearch behind a reverse proxy on a subpath. This option defines such path on the fluent-bit side. It simply adds a path prefix in the indexing HTTP POST URI. | string |
| compress | Set payload compression mechanism. Option available is 'gzip' | string |
| bufferSize | Specify the buffer size used to read the response from the Elasticsearch HTTP service. This option is useful for debugging purposes where is required to read full responses, note that response size grows depending of the number of records inserted. To set an unlimited amount of memory set this value to False, otherwise the value must be according to the Unit Size specification. | string |
| pipeline | Newer versions of Elasticsearch allows setting up filters called pipelines. This option allows defining which pipeline the database should use. For performance reasons is strongly suggested parsing and filtering on Fluent Bit side, avoid pipelines. | string |
| awsAuth | Enable AWS Sigv4 Authentication for Amazon ElasticSearch Service. | string |
| awsRegion | Specify the AWS region for Amazon ElasticSearch Service. | string |
| awsSTSEndpoint | Specify the custom sts endpoint to be used with STS API for Amazon ElasticSearch Service. | string |
| awsRoleARN | AWS IAM Role to assume to put records to your Amazon ES cluster. | string |
| awsExternalID | External ID for the AWS IAM Role specified with aws_role_arn. | string |
| cloudID | If you are using Elastic's Elasticsearch Service you can specify the cloud_id of the cluster running. | string |
| cloudAuth | Specify the credentials to use to connect to Elastic's Elasticsearch Service running on Elastic Cloud. | string |
| httpUser | Optional username credential for Elastic X-Pack access | *[plugins.Secret](../secret.md) |
| httpPassword | Password for user defined in HTTP_User | *[plugins.Secret](../secret.md) |
| index | Index name | string |
| type | Type name | string |
| logstashFormat | Enable Logstash format compatibility. This option takes a boolean value: True/False, On/Off | *bool |
| logstashPrefix | When Logstash_Format is enabled, the Index name is composed using a prefix and the date, e.g: If Logstash_Prefix is equals to 'mydata' your index will become 'mydata-YYYY.MM.DD'. The last string appended belongs to the date when the data is being generated. | string |
| logstashDateFormat | Time format (based on strftime) to generate the second part of the Index name. | string |
| timeKey | When Logstash_Format is enabled, each record will get a new timestamp field. The Time_Key property defines the name of that field. | string |
| timeKeyFormat | When Logstash_Format is enabled, this property defines the format of the timestamp. | string |
| timeKeyNanos | When Logstash_Format is enabled, enabling this property sends nanosecond precision timestamps. | *bool |
| includeTagKey | When enabled, it append the Tag name to the record. | *bool |
| tagKey | When Include_Tag_Key is enabled, this property defines the key name for the tag. | string |
| generateID | When enabled, generate _id for outgoing records. This prevents duplicate records when retrying ES. | *bool |
| idKey | If set, _id will be the value of the key from incoming record and Generate_ID option is ignored. | string |
| writeOperation | Operation to use to write in bulk requests. | string |
| replaceDots | When enabled, replace field name dots with underscore, required by Elasticsearch 2.0-2.3. | *bool |
| traceOutput | When enabled print the elasticsearch API calls to stdout (for diag only) | *bool |
| traceError | When enabled print the elasticsearch API calls to stdout when elasticsearch returns an error | *bool |
| currentTimeIndex | Use current time for index generation instead of message record | *bool |
| logstashPrefixKey | Prefix keys with this string | string |
| suppressTypeName | When enabled, mapping types is removed and Type option is ignored. Types are deprecated in APIs in v7.0. This options is for v7.0 or later. | string |
| tls |  | *[plugins.TLS](../tls.md) |
