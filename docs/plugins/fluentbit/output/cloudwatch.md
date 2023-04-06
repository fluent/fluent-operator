# CloudWatch

CloudWatch is the AWS CloudWatch output plugin, allows you to ingest your records into AWS CloudWatch. <br /> **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/outputs/cloudwatch**


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| region | AWS Region | string |
| logGroupName | Name of Cloudwatch Log Group to send log records to | string |
| logGroupTemplate | Template for Log Group name, overrides LogGroupName if set. | string |
| logStreamName | The name of the CloudWatch Log Stream to send log records to | string |
| logStreamPrefix | Prefix for the Log Stream name. Not compatible with LogStreamName setting | string |
| logStreamTemplate | Template for Log Stream name. Overrides LogStreamPrefix and LogStreamName if set. | string |
| logKey | If set, only the value of the key will be sent to CloudWatch | string |
| logFormat | Optional parameter to tell CloudWatch the format of the data | string |
| roleArn | Role ARN to use for cross-account access | string |
| autoCreateGroup | Automatically create the log group. Defaults to False. | *bool |
| logRetentionDays | Number of days logs are retained for | *int32 |
| endpoint | Custom endpoint for CloudWatch logs API | string |
| metricNamespace | Optional string to represent the CloudWatch namespace. | string |
| metricDimensions | Optional lists of lists for dimension keys to be added to all metrics. Use comma separated strings for one list of dimensions and semicolon separated strings for list of lists dimensions. | string |
| stsEndpoint | Specify a custom STS endpoint for the AWS STS API | string |
| autoRetryRequests | Automatically retry failed requests to CloudWatch once. Defaults to True. | *bool |
| externalID | Specify an external ID for the STS API. | string |
