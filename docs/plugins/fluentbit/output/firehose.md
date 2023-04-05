# Firehose

The Firehose output plugin, allows to ingest your records into AWS Firehose. <br /> It uses the new high performance kinesis_firehose plugin (written in C) instead <br /> of the older firehose plugin (written in Go). <br /> The fluent-bit container must have the plugin installed. <br /> https://docs.fluentbit.io/manual/pipeline/outputs/firehose <br /> https://github.com/aws/amazon-kinesis-firehose-for-fluent-bit <br />


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| region | The AWS region. | string |
| deliveryStream | The name of the Kinesis Firehose Delivery stream that you want log records sent to. | string |
| timeKey | Add the timestamp to the record under this key. By default, the timestamp from Fluent Bit will not be added to records sent to Kinesis. | *string |
| timeKeyFormat | strftime compliant format string for the timestamp; for example, %Y-%m-%dT%H *string This option is used with time_key. You can also use %L for milliseconds and %f for microseconds. If you are using ECS FireLens, make sure you are running Amazon ECS Container Agent v1.42.0 or later, otherwise the timestamps associated with your container logs will only have second precision. | *string |
| dataKeys | By default, the whole log record will be sent to Kinesis. If you specify a key name(s) with this option, then only those keys and values will be sent to Kinesis. For example, if you are using the Fluentd Docker log driver, you can specify data_keys log and only the log message will be sent to Kinesis. If you specify multiple keys, they should be comma delimited. | *string |
| logKey | By default, the whole log record will be sent to Firehose. If you specify a key name with this option, then only the value of that key will be sent to Firehose. For example, if you are using the Fluentd Docker log driver, you can specify log_key log and only the log message will be sent to Firehose. | *string |
| roleARN | ARN of an IAM role to assume (for cross account access). | *string |
| endpoint | Specify a custom endpoint for the Kinesis Firehose API. | *string |
| stsEndpoint | Specify a custom endpoint for the STS API; used to assume your custom role provided with role_arn. | *string |
| autoRetryRequests | Immediately retry failed requests to AWS services once. This option does not affect the normal Fluent Bit retry mechanism with backoff. Instead, it enables an immediate retry with no delay for networking errors, which may help improve throughput when there are transient/random networking issues. | *bool |
