package output

// Firehose defines the parametes for out_firehose output plugin
type Firehose struct {
	// The AWS region.
	Region *string `json:"region,omitempty"`
	// The name of the Kinesis Firehose Delivery stream that you want log records sent to.
	DeliveryStream *string `json:"deliveryStream,omitempty"`
	// Add the timestamp to the record under this key. By default, the timestamp from Fluent Bit will not be added to records sent to Kinesis.
	TimeKey *string `json:"timeKey,omitempty"`
	// strftime compliant format string for the timestamp; for example, %Y-%m-%dT%H *string This option is used with time_key. You can also use %L for milliseconds and %f for microseconds. If you are using ECS FireLens, make sure you are running Amazon ECS Container Agent v1.42.0 or later, otherwise the timestamps associated with your container logs will only have second precision.
	TimeKeyFormat *string `json:"timeKeyFormat,omitempty"`
	// By default, the whole log record will be sent to Kinesis. If you specify a key name(s) with this option, then only those keys and values will be sent to Kinesis. For example, if you are using the Fluentd Docker log driver, you can specify data_keys log and only the log message will be sent to Kinesis. If you specify multiple keys, they should be comma delimited.
	DataKeys *string `json:"dataKeys,omitempty"`
	// By default, the whole log record will be sent to Firehose. If you specify a key name with this option, then only the value of that key will be sent to Firehose. For example, if you are using the Fluentd Docker log driver, you can specify log_key log and only the log message will be sent to Firehose.
	LogKey *string `json:"logKey,omitempty"`
	// ARN of an IAM role to assume (for cross account access).
	RoleARN *string `json:"roleArn,omitempty"`
	// Web identity token file
	WebIdentityTokenFile *string `json:"webIdentityTokenFile,omitempty"`
	// Role Session name
	RoleSessionName *string `json:"roleSessionName,omitempty"`
	// Specify a custom endpoint for the Kinesis Firehose API.
	Endpoint *string `json:"endpoint,omitempty"`
	// Specify a custom endpoint for the STS API; used to assume your custom role provided with role_arn.
	STSEndpoint *string `json:"stsEndpoint,omitempty"`
	// Immediately retry failed requests to AWS services once. This option does not affect the normal Fluent Bit retry mechanism with backoff. Instead, it enables an immediate retry with no delay for networking errors, which may help improve throughput when there are transient/random networking issues.
	AutoRetryRequests *bool `json:"autoRetryRequests,omitempty"`
	// Debug
	Debug *bool `json:"debug,omitempty"`
}
