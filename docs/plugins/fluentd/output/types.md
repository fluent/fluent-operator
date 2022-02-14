# OutputCommon

OutputCommon defines the common parameters for output plugin


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| logLevel | The @log_level parameter specifies the plugin-specific logging level | *string |
# Output

Output defines all types for output plugin


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| forward | out_forward plugin | *Forward |
| http | out_http plugin | *Http |
| elasticsearch | out_es plugin | *Elasticsearch |
| kafka | out_kafka plugin | *Kafka2 |
| s3 | out_s3 plugin | *S3 |
| stdout | out_stdout plugin | *Stdout |
