# OutputCommon

OutputCommon defines the common parameters for output plugin


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| logLevel | The @log_level parameter specifies the plugin-specific logging level | *string |
| tag | Which tag to be matched. | *string |
# Output

Output defines all available output plugins and their parameters


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| forward | out_forward plugin | *[Forward](#forward) |
| http | out_http plugin | *[Http](#http) |
| elasticsearch | out_es plugin | *[Elasticsearch](#elasticsearch) |
| elasticsearchDataStream | out_es datastreams plugin | *[ElasticsearchDataStream](#elasticsearchdatastream) |
| elasticsearchDynamic | out_es dynamic plugin (elasticsearch_dynamic). It takes the same parameters as elasticsearch, but host, port, index and logstash values can use placeholders like ${tag_parts[1]} or ${record['key']}. | *[Elasticsearch](#elasticsearch) |
| opensearch | out_opensearch plugin | *[Opensearch](#opensearch) |
| kafka | out_kafka plugin | *[Kafka2](#kafka2) |
| s3 | out_s3 plugin | *[S3](#s3) |
| stdout | out_stdout plugin | *[Stdout](#stdout) |
| loki | out_loki plugin | *[Loki](#loki) |
| customPlugin | Custom plugin type | *[custom.CustomPlugin](../custom/custom.md) |
| cloudWatch | out_cloudwatch plugin | *[CloudWatch](#cloudwatch) |
| datadog | datadog plugin | *[Datadog](#datadog) |
| copy | copy plugin | *[Copy](#copy) |
| nullPlugin | null plugin | *[Null](#null) |
