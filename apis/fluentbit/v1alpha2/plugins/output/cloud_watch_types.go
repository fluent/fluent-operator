package output

import (
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins/params"
)

// +kubebuilder:object:generate:=true

// CloudWatch is the AWS CloudWatch output plugin, allows you to ingest your records into AWS CloudWatch. <br />
// **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/outputs/cloudwatch**
type CloudWatch struct {
	// AWS Region
	Region string `json:"region"`
	// Name of Cloudwatch Log Group to send log records to
	LogGroupName string `json:"logGroupName,omitempty"`
	// Template for Log Group name, overrides LogGroupName if set.
	LogGroupTemplate string `json:"logGroupTemplate,omitempty"`
	// The name of the CloudWatch Log Stream to send log records to
	LogStreamName string `json:"logStreamName,omitempty"`
	// Prefix for the Log Stream name. Not compatible with LogStreamName setting
	LogStreamPrefix string `json:"logStreamPrefix,omitempty"`
	// Template for Log Stream name. Overrides LogStreamPrefix and LogStreamName if set.
	LogStreamTemplate string `json:"logStreamTemplate,omitempty"`
	// If set, only the value of the key will be sent to CloudWatch
	LogKey string `json:"logKey,omitempty"`
	// Optional parameter to tell CloudWatch the format of the data
	LogFormat string `json:"logFormat,omitempty"`
	// Role ARN to use for cross-account access
	RoleArn string `json:"roleArn,omitempty"`
	// Automatically create the log group. Defaults to False.
	AutoCreateGroup *bool `json:"autoCreateGroup,omitempty"`
	// Number of days logs are retained for
	// +kubebuilder:validation:Enum:=1;3;5;7;14;30;60;90;120;150;180;365;400;545;731;1827;3653
	LogRetentionDays *int32 `json:"logRetentionDays,omitempty"`
	// Custom endpoint for CloudWatch logs API
	Endpoint string `json:"endpoint,omitempty"`
	// Optional string to represent the CloudWatch namespace.
	MetricNamespace string `json:"metricNamespace,omitempty"`
	// Optional lists of lists for dimension keys to be added to all metrics. Use comma separated strings
	// for one list of dimensions and semicolon separated strings for list of lists dimensions.
	MetricDimensions string `json:"metricDimensions,omitempty"`
	// Specify a custom STS endpoint for the AWS STS API
	StsEndpoint string `json:"stsEndpoint,omitempty"`
	// Automatically retry failed requests to CloudWatch once. Defaults to True.
	AutoRetryRequests *bool `json:"autoRetryRequests,omitempty"`
	// Specify an external ID for the STS API.
	ExternalID string `json:"externalID,omitempty"`
}

// Name implement Section() method
func (*CloudWatch) Name() string {
	return "cloudwatch_logs"
}

// Params implement Section() method
func (o *CloudWatch) Params(sl plugins.SecretLoader) (*params.KVs, error) {
	kvs := params.NewKVs()

	plugins.InsertKVString(kvs, "region", o.Region)
	plugins.InsertKVString(kvs, "log_group_name", o.LogGroupName)
	plugins.InsertKVString(kvs, "log_group_template", o.LogGroupTemplate)
	plugins.InsertKVString(kvs, "log_stream_name", o.LogStreamName)
	plugins.InsertKVString(kvs, "log_stream_prefix", o.LogStreamPrefix)
	plugins.InsertKVString(kvs, "log_stream_template", o.LogStreamTemplate)
	plugins.InsertKVString(kvs, "log_key", o.LogKey)
	plugins.InsertKVString(kvs, "log_format", o.LogFormat)
	plugins.InsertKVString(kvs, "role_arn", o.RoleArn)
	plugins.InsertKVString(kvs, "endpoint", o.Endpoint)
	plugins.InsertKVString(kvs, "metric_namespace", o.MetricNamespace)
	plugins.InsertKVString(kvs, "metric_dimensions", o.MetricDimensions)
	plugins.InsertKVString(kvs, "sts_endpoint", o.StsEndpoint)
	plugins.InsertKVString(kvs, "external_id", o.ExternalID)

	plugins.InsertKVField(kvs, "log_retention_days", o.LogRetentionDays)
	plugins.InsertKVField(kvs, "auto_retry_requests", o.AutoRetryRequests)
	plugins.InsertKVField(kvs, "auto_create_group", o.AutoCreateGroup)

	return kvs, nil
}
