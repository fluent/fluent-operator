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
| forward | out_forward plugin | *Forward |
| http | out_http plugin | *Http |
| elasticsearch | out_es plugin | *Elasticsearch |
| opensearch | out_opensearch plugin | *Opensearch |
| kafka | out_kafka plugin | *Kafka2 |
| s3 | out_s3 plugin | *S3 |
| stdout | out_stdout plugin | *Stdout |
| loki | out_loki plugin | *Loki |
| customPlugin | Custom plugin type | *custom.CustomPlugin |
| cloudWatch | out_cloudwatch plugin | *CloudWatch |
| datadog | datadog plugin | *Datadog |
| copy | copy plugin | *Copy |
