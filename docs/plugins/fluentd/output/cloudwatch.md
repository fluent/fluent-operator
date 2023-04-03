# CloudWatch

CloudWatch defines the parametes for out_cloudwatch output plugin


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| autoCreateStream |  | *bool |
| awsKeyId |  | *[plugins.Secret](../secret.md) |
| awsSecKey |  | *[plugins.Secret](../secret.md) |
| awsUseSts |  | *bool |
| awsStsRoleArn |  | *string |
| awsStsSessionName |  | *string |
| awsStsExternalId |  | *string |
| awsStsPolicy |  | *string |
| awsStsDurationSeconds |  | *string |
| awsStsEndpointUrl |  | *string |
| awsEcsAuthentication |  | *bool |
| concurrency |  | *int |
| endpoint | Specify an AWS endpoint to send data to. | *string |
| sslVerifyPeer |  | *bool |
| httpProxy |  | *string |
| includeTimeKey |  | *bool |
| jsonHandler |  | *string |
| localtime |  | *bool |
| logGroupAwsTags |  | *string |
| logGroupAwsTagsKey |  | *string |
| logGroupName |  | *string |
| logGroupNameKey |  | *string |
| logRejectedRequest |  | *string |
| logStreamName |  | *string |
| logStreamNameKey |  | *string |
| maxEventsPerBatch |  | *string |
| maxMessageLength |  | *string |
| messageKeys |  | *string |
| putLogEventsDisableRetryLimit |  | *bool |
| putLogEventsRetryLimit |  | *string |
| putLogEventsRetryWait |  | *string |
| region | The AWS region. | *string |
| removeLogGroupAwsTagsKey |  | *bool |
| removeLogGroupNameKey |  | *bool |
| removeLogStreamNameKey |  | *bool |
| removeRetentionInDaysKey |  | *bool |
| retentionInDays |  | *string |
| retentionInDaysKey |  | *string |
| useTagAsGroup |  | *string |
| useTagAsStream |  | *string |
| roleArn | ARN of an IAM role to assume (for cross account access). | *string |
| roleSessionName | Role Session name | *string |
| webIdentityTokenFile | Web identity token file | *string |
| policy |  | *string |
| durationSeconds |  | *string |
