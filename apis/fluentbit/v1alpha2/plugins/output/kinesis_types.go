package output

import (
	"fmt"
	"github.com/fluent/fluent-operator/v2/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v2/apis/fluentbit/v1alpha2/plugins/params"
)

// +kubebuilder:object:generate:=true

// The Kinesis output plugin, allows to ingest your records into AWS Kinesis. <br />
// It uses the new high performance and highly efficient kinesis plugin is called kinesis_streams instead of the older Golang Fluent Bit plugin released in 2019.
// https://docs.fluentbit.io/manual/pipeline/outputs/kinesis <br />
// https://github.com/aws/amazon-kinesis-streams-for-fluent-bit <br />
type Kinesis struct {
	// The AWS region.
	Region string `json:"region"`
	// The name of the Kinesis Streams Delivery stream that you want log records sent to.
	Stream string `json:"stream"`
	// Add the timestamp to the record under this key. By default the timestamp from Fluent Bit will not be added to records sent to Kinesis.
	TimeKey string `json:"timeKey,omitempty"`
	// strftime compliant format string for the timestamp; for example, the default is '%Y-%m-%dT%H:%M:%S'. Supports millisecond precision with '%3N' and supports nanosecond precision with '%9N' and '%L'; for example, adding '%3N' to support millisecond '%Y-%m-%dT%H:%M:%S.%3N'. This option is used with time_key.
	TimeKeyFormat string `json:"timeKeyFormat,omitempty"`
	// By default, the whole log record will be sent to Kinesis. If you specify a key name with this option, then only the value of that key will be sent to Kinesis. For example, if you are using the Fluentd Docker log driver, you can specify log_key log and only the log message will be sent to Kinesis.
	LogKey string `json:"logKey,omitempty"`
	// ARN of an IAM role to assume (for cross account access).
	RoleARN string `json:"roleARN,omitempty"`
	// Specify a custom endpoint for the Kinesis API.
	Endpoint string `json:"endpoint,omitempty"`
	// Custom endpoint for the STS API.
	STSEndpoint string `json:"stsEndpoint,omitempty"`
	// Immediately retry failed requests to AWS services once. This option does not affect the normal Fluent Bit retry mechanism with backoff. Instead, it enables an immediate retry with no delay for networking errors, which may help improve throughput when there are transient/random networking issues. This option defaults to true.
	AutoRetryRequests *bool `json:"autoRetryRequests,omitempty"`
	// Specify an external ID for the STS API, can be used with the role_arn parameter if your role requires an external ID.
	ExternalID string `json:"externalID,omitempty"`
}

// Name implement Section() method
func (*Kinesis) Name() string {
	return "kinesis_streams"
}

// Params implement Section() method
func (k *Kinesis) Params(_ plugins.SecretLoader) (*params.KVs, error) {
	kvs := params.NewKVs()
	if k.Region != "" {
		kvs.Insert("region", k.Region)
	}
	if k.Stream != "" {
		kvs.Insert("stream", k.Stream)
	}
	if k.TimeKey != "" {
		kvs.Insert("time_key", k.TimeKey)
	}
	if k.TimeKeyFormat != "" {
		kvs.Insert("time_key_format", k.TimeKeyFormat)
	}
	if k.LogKey != "" {
		kvs.Insert("log_key", k.LogKey)
	}
	if k.RoleARN != "" {
		kvs.Insert("role_arn", k.RoleARN)
	}
	if k.Endpoint != "" {
		kvs.Insert("endpoint", k.Endpoint)
	}
	if k.STSEndpoint != "" {
		kvs.Insert("sts_endpoint", k.STSEndpoint)
	}
	if k.AutoRetryRequests != nil {
		kvs.Insert("auto_retry_requests", fmt.Sprint(*k.AutoRetryRequests))
	}
	if k.ExternalID != "" {
		kvs.Insert("external_id", k.ExternalID)
	}

	return kvs, nil
}
