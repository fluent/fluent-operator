# Kinesis

The Kinesis output plugin, allows to ingest your records into AWS Kinesis. <br /> It uses the new high performance and highly efficient kinesis plugin is called kinesis_streams instead of the older Golang Fluent Bit plugin released in 2019. https://docs.fluentbit.io/manual/pipeline/outputs/kinesis <br /> https://github.com/aws/amazon-kinesis-streams-for-fluent-bit <br />


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| region | The AWS region. | string |
| stream | The name of the Kinesis Streams Delivery stream that you want log records sent to. | string |
| timeKey | Add the timestamp to the record under this key. By default the timestamp from Fluent Bit will not be added to records sent to Kinesis. | string |
| timeKeyFormat | strftime compliant format string for the timestamp; for example, the default is '%Y-%m-%dT%H:%M:%S'. Supports millisecond precision with '%3N' and supports nanosecond precision with '%9N' and '%L'; for example, adding '%3N' to support millisecond '%Y-%m-%dT%H:%M:%S.%3N'. This option is used with time_key. | string |
| logKey | By default, the whole log record will be sent to Kinesis. If you specify a key name with this option, then only the value of that key will be sent to Kinesis. For example, if you are using the Fluentd Docker log driver, you can specify log_key log and only the log message will be sent to Kinesis. | string |
| roleARN | ARN of an IAM role to assume (for cross account access). | string |
| endpoint | Specify a custom endpoint for the Kinesis API. | string |
| stsEndpoint | Custom endpoint for the STS API. | string |
| autoRetryRequests | Immediately retry failed requests to AWS services once. This option does not affect the normal Fluent Bit retry mechanism with backoff. Instead, it enables an immediate retry with no delay for networking errors, which may help improve throughput when there are transient/random networking issues. This option defaults to true. | *bool |
| externalID | Specify an external ID for the STS API, can be used with the role_arn parameter if your role requires an external ID. | string |
