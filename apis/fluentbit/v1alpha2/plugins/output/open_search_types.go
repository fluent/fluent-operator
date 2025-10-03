package output

import (
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins/params"
)

// +kubebuilder:object:generate:=true

// OpenSearch is the opensearch output plugin, allows to ingest your records into an OpenSearch database. <br />
// **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/outputs/opensearch**
type OpenSearch struct {
	// IP address or hostname of the target OpenSearch instance, default `127.0.0.1`
	Host string `json:"host,omitempty"`
	// TCP port of the target OpenSearch instance, default `9200`
	// +kubebuilder:validation:Minimum:=1
	// +kubebuilder:validation:Maximum:=65535
	Port *int32 `json:"port,omitempty"`
	// OpenSearch accepts new data on HTTP query path "/_bulk".
	// But it is also possible to serve OpenSearch behind a reverse proxy on a subpath.
	// This option defines such path on the fluent-bit side.
	// It simply adds a path prefix in the indexing HTTP POST URI.
	Path string `json:"path,omitempty"`
	// Specify the buffer size used to read the response from the OpenSearch HTTP service.
	// This option is useful for debugging purposes where is required to read full responses,
	// note that response size grows depending of the number of records inserted.
	// To set an unlimited amount of memory set this value to False,
	// otherwise the value must be according to the Unit Size specification.
	// +kubebuilder:validation:Pattern:="^\\d+(k|K|KB|kb|m|M|MB|mb|g|G|GB|gb)?$"
	BufferSize string `json:"bufferSize,omitempty"`
	// OpenSearch allows to setup filters called pipelines.
	// This option allows to define which pipeline the database should use.
	// For performance reasons is strongly suggested to do parsing
	// and filtering on Fluent Bit side, avoid pipelines.
	Pipeline string `json:"pipeline,omitempty"`
	// Enable AWS Sigv4 Authentication for Amazon OpenSearch Service.
	AWSAuth string `json:"awsAuth,omitempty"`
	// Specify the AWS region for Amazon OpenSearch Service.
	AWSRegion string `json:"awsRegion,omitempty"`
	// Specify the custom sts endpoint to be used with STS API for Amazon OpenSearch Service.
	AWSSTSEndpoint string `json:"awsSTSEndpoint,omitempty"`
	// AWS IAM Role to assume to put records to your Amazon cluster.
	AWSRoleARN string `json:"awsRoleARN,omitempty"`
	// External ID for the AWS IAM Role specified with aws_role_arn.
	AWSExternalID string `json:"awsExternalID,omitempty"`
	// Optional username credential for access
	HTTPUser *plugins.Secret `json:"httpUser,omitempty"`
	// Password for user defined in HTTP_User
	HTTPPasswd *plugins.Secret `json:"httpPassword,omitempty"`
	// Index name
	Index string `json:"index,omitempty"`
	// Type name
	Type string `json:"type,omitempty"`
	// Enable Logstash format compatibility.
	// This option takes a boolean value: True/False, On/Off
	LogstashFormat *bool `json:"logstashFormat,omitempty"`
	// When Logstash_Format is enabled, the Index name is composed using a prefix and the date,
	// e.g: If Logstash_Prefix is equals to 'mydata' your index will become 'mydata-YYYY.MM.DD'.
	// The last string appended belongs to the date when the data is being generated.
	LogstashPrefix string `json:"logstashPrefix,omitempty"`
	// Time format (based on strftime) to generate the second part of the Index name.
	LogstashDateFormat string `json:"logstashDateFormat,omitempty"`
	// When Logstash_Format is enabled, each record will get a new timestamp field.
	// The Time_Key property defines the name of that field.
	TimeKey string `json:"timeKey,omitempty"`
	// When Logstash_Format is enabled, this property defines the format of the timestamp.
	TimeKeyFormat string `json:"timeKeyFormat,omitempty"`
	// When Logstash_Format is enabled, enabling this property sends nanosecond precision timestamps.
	TimeKeyNanos *bool `json:"timeKeyNanos,omitempty"`
	// When enabled, it append the Tag name to the record.
	IncludeTagKey *bool `json:"includeTagKey,omitempty"`
	// When Include_Tag_Key is enabled, this property defines the key name for the tag.
	TagKey string `json:"tagKey,omitempty"`
	// When enabled, generate _id for outgoing records.
	// This prevents duplicate records when retrying OpenSearch.
	GenerateID *bool `json:"generateID,omitempty"`
	// If set, _id will be the value of the key from incoming record and Generate_ID option is ignored.
	IdKey string `json:"idKey,omitempty"`
	// Operation to use to write in bulk requests.
	WriteOperation string `json:"writeOperation,omitempty"`
	// When enabled, replace field name dots with underscore, required by Opensearch 2.0-2.3.
	ReplaceDots *bool `json:"replaceDots,omitempty"`
	// When enabled print the Opensearch API calls to stdout (for diag only)
	TraceOutput *bool `json:"traceOutput,omitempty"`
	// When enabled print the Opensearch API calls to stdout when Opensearch returns an error
	TraceError *bool `json:"traceError,omitempty"`
	// Use current time for index generation instead of message record
	CurrentTimeIndex *bool `json:"currentTimeIndex,omitempty"`
	// Prefix keys with this string
	LogstashPrefixKey string `json:"logstashPrefixKey,omitempty"`
	// When enabled, mapping types is removed and Type option is ignored. Types are deprecated in APIs in v7.0. This options is for v7.0 or later.
	SuppressTypeName *bool `json:"suppressTypeName,omitempty"`
	// Enables dedicated thread(s) for this output. Default value is set since version 1.8.13. For previous versions is 0.
	Workers      *int32 `json:"workers,omitempty"`
	*plugins.TLS `json:"tls,omitempty"`
	// Include fluentbit networking options for this output-plugin
	*plugins.Networking `json:"networking,omitempty"`
	// Limit the maximum number of Chunks in the filesystem for the current output logical destination.
	TotalLimitSize string `json:"totalLimitSize,omitempty"`
	// +kubebuilder:validation:Enum=gzip
	Compress string `json:"compress,omitempty"`
}

// Name implement Section() method
func (*OpenSearch) Name() string {
	return "opensearch"
}

// Params implement Section() method
func (o *OpenSearch) Params(sl plugins.SecretLoader) (*params.KVs, error) {
	kvs := params.NewKVs()

	if err := plugins.InsertKVSecret(kvs, "HTTP_User", o.HTTPUser, sl); err != nil {
		return nil, err
	}
	if err := plugins.InsertKVSecret(kvs, "HTTP_Passwd", o.HTTPPasswd, sl); err != nil {
		return nil, err
	}

	plugins.InsertKVString(kvs, "Host", o.Host)
	plugins.InsertKVField(kvs, "Port", o.Port)
	plugins.InsertKVString(kvs, "Index", o.Index)
	plugins.InsertKVString(kvs, "Type", o.Type)
	plugins.InsertKVString(kvs, "Path", o.Path)
	plugins.InsertKVString(kvs, "Buffer_Size", o.BufferSize)
	plugins.InsertKVString(kvs, "Pipeline", o.Pipeline)
	plugins.InsertKVString(kvs, "AWS_Auth", o.AWSAuth)
	plugins.InsertKVString(kvs, "AWS_Region", o.AWSRegion)
	plugins.InsertKVString(kvs, "AWS_STS_Endpoint", o.AWSSTSEndpoint)
	plugins.InsertKVString(kvs, "AWS_Role_ARN", o.AWSRoleARN)
	plugins.InsertKVString(kvs, "AWS_External_ID", o.AWSExternalID)
	plugins.InsertKVString(kvs, "Logstash_Prefix", o.LogstashPrefix)
	plugins.InsertKVString(kvs, "Logstash_DateFormat", o.LogstashDateFormat)
	plugins.InsertKVString(kvs, "Time_Key", o.TimeKey)
	plugins.InsertKVString(kvs, "Time_Key_Format", o.TimeKeyFormat)
	plugins.InsertKVString(kvs, "Tag_Key", o.TagKey)
	plugins.InsertKVString(kvs, "ID_KEY", o.IdKey)
	plugins.InsertKVString(kvs, "Write_Operation", o.WriteOperation)
	plugins.InsertKVString(kvs, "Logstash_Prefix_Key", o.LogstashPrefixKey)
	plugins.InsertKVString(kvs, "storage.total_limit_size", o.TotalLimitSize)
	plugins.InsertKVString(kvs, "Compress", o.Compress)

	plugins.InsertKVField(kvs, "Logstash_Format", o.LogstashFormat)
	plugins.InsertKVField(kvs, "Time_Key_Nanos", o.TimeKeyNanos)
	plugins.InsertKVField(kvs, "Include_Tag_Key", o.IncludeTagKey)
	plugins.InsertKVField(kvs, "Generate_ID", o.GenerateID)
	plugins.InsertKVField(kvs, "Replace_Dots", o.ReplaceDots)
	plugins.InsertKVField(kvs, "Trace_Output", o.TraceOutput)
	plugins.InsertKVField(kvs, "Trace_Error", o.TraceError)
	plugins.InsertKVField(kvs, "Current_Time_Index", o.CurrentTimeIndex)
	plugins.InsertKVField(kvs, "Suppress_Type_Name", o.SuppressTypeName)
	plugins.InsertKVField(kvs, "Workers", o.Workers)

	if o.TLS != nil {
		tls, err := o.TLS.Params(sl)
		if err != nil {
			return nil, err
		}
		kvs.Merge(tls)
	}
	if o.Networking != nil {
		net, err := o.Networking.Params(sl)
		if err != nil {
			return nil, err
		}
		kvs.Merge(net)
	}

	return kvs, nil
}
