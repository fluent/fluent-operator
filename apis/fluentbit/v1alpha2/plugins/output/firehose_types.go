package output

import (
	"strconv"

	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins/params"
)

// +kubebuilder:object:generate:=true

// The Firehose output plugin, allows to ingest your records into AWS Firehose. <br />
// It uses the new high performance kinesis_firehose plugin (written in C) instead <br />
// of the older firehose plugin (written in Go). <br />
// The fluent-bit container must have the plugin installed. <br />
// https://docs.fluentbit.io/manual/pipeline/outputs/firehose <br />
// https://github.com/aws/amazon-kinesis-firehose-for-fluent-bit <br />
type Firehose struct {
	// The AWS region.
	Region string `json:"region"`
	// The name of the Kinesis Firehose Delivery stream that you want log records sent to.
	DeliveryStream string `json:"deliveryStream"`
	// Add the timestamp to the record under this key. By default, the timestamp from Fluent Bit will not be added to records sent to Kinesis.
	TimeKey *string `json:"timeKey,omitempty"`
	// strftime compliant format string for the timestamp; for example, %Y-%m-%dT%H *string This option is used with time_key. You can also use %L for milliseconds and %f for microseconds. If you are using ECS FireLens, make sure you are running Amazon ECS Container Agent v1.42.0 or later, otherwise the timestamps associated with your container logs will only have second precision.
	TimeKeyFormat *string `json:"timeKeyFormat,omitempty"`
	// By default, the whole log record will be sent to Kinesis. If you specify a key name(s) with this option, then only those keys and values will be sent to Kinesis. For example, if you are using the Fluentd Docker log driver, you can specify data_keys log and only the log message will be sent to Kinesis. If you specify multiple keys, they should be comma delimited.
	DataKeys *string `json:"dataKeys,omitempty"`
	// By default, the whole log record will be sent to Firehose. If you specify a key name with this option, then only the value of that key will be sent to Firehose. For example, if you are using the Fluentd Docker log driver, you can specify log_key log and only the log message will be sent to Firehose.
	LogKey *string `json:"logKey,omitempty"`
	// ARN of an IAM role to assume (for cross account access).
	RoleARN *string `json:"roleARN,omitempty"`
	// Specify a custom endpoint for the Kinesis Firehose API.
	Endpoint *string `json:"endpoint,omitempty"`
	// Specify a custom endpoint for the STS API; used to assume your custom role provided with role_arn.
	STSEndpoint *string `json:"stsEndpoint,omitempty"`
	// Immediately retry failed requests to AWS services once. This option does not affect the normal Fluent Bit retry mechanism with backoff. Instead, it enables an immediate retry with no delay for networking errors, which may help improve throughput when there are transient/random networking issues.
	AutoRetryRequests *bool `json:"autoRetryRequests,omitempty"`
	// Compression type to use when compressing the data. Valid values are: gzip, snappy, lz4, zstd. If you do not specify a compression type, the data will be sent uncompressed.
	// +kubebuilder:validation:Enum=gzip;snappy;lz4;zstd
	Compression *string `json:"compression,omitempty"`
	// Specify an external ID for the STS API, can be used with the role_arn parameter if your role requires an external ID.
	ExternalID *string `json:"externalID,omitempty"`
	// Option to specify an AWS Profile for credentials.
	Profile *string `json:"profile,omitempty"`
	// Option to enable simple aggregation for the Firehose output plugin.
	SimpleAggregation *bool `json:"simpleAggregation,omitempty"`
	// Specify number of worker threads to use to output to Firehose
	Workers *int32 `json:"workers,omitempty"`
}

// implement Section() method
func (*Firehose) Name() string {
	return "kinesis_firehose"
}

// implement Section() method
func (l *Firehose) Params(sl plugins.SecretLoader) (*params.KVs, error) {
	kvs := params.NewKVs()

	plugins.InsertKVString(kvs, "region", l.Region)
	plugins.InsertKVString(kvs, "delivery_stream", l.DeliveryStream)

	if l.DataKeys != nil && *l.DataKeys != "" {
		kvs.Insert("data_keys", *l.DataKeys)
	}
	if l.LogKey != nil && *l.LogKey != "" {
		kvs.Insert("log_key", *l.LogKey)
	}
	if l.RoleARN != nil && *l.RoleARN != "" {
		kvs.Insert("role_arn", *l.RoleARN)
	}
	if l.Endpoint != nil && *l.Endpoint != "" {
		kvs.Insert("endpoint", *l.Endpoint)
	}
	if l.STSEndpoint != nil && *l.STSEndpoint != "" {
		kvs.Insert("sts_endpoint", *l.STSEndpoint)
	}
	if l.TimeKey != nil && *l.TimeKey != "" {
		kvs.Insert("time_key", *l.TimeKey)
	}
	if l.TimeKeyFormat != nil && *l.TimeKeyFormat != "" {
		kvs.Insert("time_key_format", *l.TimeKeyFormat)
	}

	if l.AutoRetryRequests != nil {
		kvs.Insert("auto_retry_requests", strconv.FormatBool(*l.AutoRetryRequests))
	}

	if l.Compression != nil && *l.Compression != "" {
		kvs.Insert("compression", *l.Compression)
	}
	if l.ExternalID != nil && *l.ExternalID != "" {
		kvs.Insert("external_id", *l.ExternalID)
	}
	if l.SimpleAggregation != nil {
		kvs.Insert("simple_aggregation", strconv.FormatBool(*l.SimpleAggregation))
	}
	if l.Profile != nil && *l.Profile != "" {
		kvs.Insert("profile", *l.Profile)
	}
	if l.Workers != nil {
		kvs.Insert("workers", strconv.FormatInt(int64(*l.Workers), 10))
	}

	return kvs, nil
}
