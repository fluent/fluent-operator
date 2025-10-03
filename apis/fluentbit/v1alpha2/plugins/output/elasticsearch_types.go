package output

import (
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins/params"
)

// +kubebuilder:object:generate:=true

// Elasticsearch is the es output plugin, allows to ingest your records into an Elasticsearch database. <br />
// **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/outputs/elasticsearch**
type Elasticsearch struct {
	// IP address or hostname of the target Elasticsearch instance
	Host string `json:"host,omitempty"`
	// TCP port of the target Elasticsearch instance
	// +kubebuilder:validation:Minimum:=1
	// +kubebuilder:validation:Maximum:=65535
	Port *int32 `json:"port,omitempty"`
	// Elasticsearch accepts new data on HTTP query path "/_bulk".
	// But it is also possible to serve Elasticsearch behind a reverse proxy on a subpath.
	// This option defines such path on the fluent-bit side.
	// It simply adds a path prefix in the indexing HTTP POST URI.
	Path string `json:"path,omitempty"`
	// Set payload compression mechanism. Option available is 'gzip'
	// +kubebuilder:validation:Enum=gzip
	Compress string `json:"compress,omitempty"`
	// Specify the buffer size used to read the response from the Elasticsearch HTTP service.
	// This option is useful for debugging purposes where is required to read full responses,
	// note that response size grows depending of the number of records inserted.
	// To set an unlimited amount of memory set this value to False,
	// otherwise the value must be according to the Unit Size specification.
	// +kubebuilder:validation:Pattern:="^\\d+(k|K|KB|kb|m|M|MB|mb|g|G|GB|gb)?$"
	BufferSize string `json:"bufferSize,omitempty"`
	// Newer versions of Elasticsearch allows setting up filters called pipelines.
	// This option allows defining which pipeline the database should use.
	// For performance reasons is strongly suggested parsing
	// and filtering on Fluent Bit side, avoid pipelines.
	Pipeline string `json:"pipeline,omitempty"`
	// Enable AWS Sigv4 Authentication for Amazon ElasticSearch Service.
	AWSAuth string `json:"awsAuth,omitempty"`
	// AWSAuthSecret Enable AWS Sigv4 Authentication for Amazon ElasticSearch Service.
	AWSAuthSecret *plugins.Secret `json:"awsAuthSecret,omitempty"`
	// Specify the AWS region for Amazon ElasticSearch Service.
	AWSRegion string `json:"awsRegion,omitempty"`
	// Specify the custom sts endpoint to be used with STS API for Amazon ElasticSearch Service.
	AWSSTSEndpoint string `json:"awsSTSEndpoint,omitempty"`
	// AWS IAM Role to assume to put records to your Amazon ES cluster.
	AWSRoleARN string `json:"awsRoleARN,omitempty"`
	// External ID for the AWS IAM Role specified with aws_role_arn.
	AWSExternalID string `json:"awsExternalID,omitempty"`
	// If you are using Elastic's Elasticsearch Service you can specify the cloud_id of the cluster running.
	CloudID string `json:"cloudID,omitempty"`
	// Specify the credentials to use to connect to Elastic's Elasticsearch Service running on Elastic Cloud.
	CloudAuth string `json:"cloudAuth,omitempty"`
	// CloudAuthSecret Specify the credentials to use to connect to Elastic's Elasticsearch Service running on Elastic Cloud.
	CloudAuthSecret *plugins.Secret `json:"cloudAuthSecret,omitempty"`
	// Optional username credential for Elastic X-Pack access
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
	// This prevents duplicate records when retrying ES.
	GenerateID *bool `json:"generateID,omitempty"`
	// If set, _id will be the value of the key from incoming record and Generate_ID option is ignored.
	IdKey string `json:"idKey,omitempty"`
	// Operation to use to write in bulk requests.
	WriteOperation string `json:"writeOperation,omitempty"`
	// When enabled, replace field name dots with underscore, required by Elasticsearch 2.0-2.3.
	ReplaceDots *bool `json:"replaceDots,omitempty"`
	// When enabled print the elasticsearch API calls to stdout (for diag only)
	TraceOutput *bool `json:"traceOutput,omitempty"`
	// When enabled print the elasticsearch API calls to stdout when elasticsearch returns an error
	TraceError *bool `json:"traceError,omitempty"`
	// Use current time for index generation instead of message record
	CurrentTimeIndex *bool `json:"currentTimeIndex,omitempty"`
	// Prefix keys with this string
	LogstashPrefixKey string `json:"logstashPrefixKey,omitempty"`
	// When enabled, mapping types is removed and Type option is ignored. Types are deprecated in APIs in v7.0. This options is for v7.0 or later.
	SuppressTypeName string `json:"suppressTypeName,omitempty"`
	*plugins.TLS     `json:"tls,omitempty"`
	// Include fluentbit networking options for this output-plugin
	*plugins.Networking `json:"networking,omitempty"`
	// Limit the maximum number of Chunks in the filesystem for the current output logical destination.
	TotalLimitSize string `json:"totalLimitSize,omitempty"`
}

// Name implement Section() method
func (*Elasticsearch) Name() string {
	return "es"
}

// Params implement Section() method
func (es *Elasticsearch) Params(sl plugins.SecretLoader) (*params.KVs, error) {
	kvs := params.NewKVs()

	if es.AWSAuthSecret != nil {
		if err := plugins.InsertKVSecret(kvs, "AWS_Auth", es.AWSAuthSecret, sl); err != nil {
			return nil, err
		}
	}
	if es.CloudAuthSecret != nil {
		if err := plugins.InsertKVSecret(kvs, "Cloud_Auth", es.CloudAuthSecret, sl); err != nil {
			return nil, err
		}
	}
	if err := plugins.InsertKVSecret(kvs, "HTTP_User", es.HTTPUser, sl); err != nil {
		return nil, err
	}
	if err := plugins.InsertKVSecret(kvs, "HTTP_Passwd", es.HTTPPasswd, sl); err != nil {
		return nil, err
	}

	plugins.InsertKVString(kvs, "Host", es.Host)
	plugins.InsertKVField(kvs, "Port", es.Port)
	plugins.InsertKVString(kvs, "Index", es.Index)
	plugins.InsertKVString(kvs, "Type", es.Type)
	plugins.InsertKVString(kvs, "Path", es.Path)
	plugins.InsertKVString(kvs, "Compress", es.Compress)
	plugins.InsertKVString(kvs, "Buffer_Size", es.BufferSize)
	plugins.InsertKVString(kvs, "Pipeline", es.Pipeline)
	plugins.InsertKVString(kvs, "AWS_Auth", es.AWSAuth)
	plugins.InsertKVString(kvs, "AWS_Region", es.AWSRegion)
	plugins.InsertKVString(kvs, "AWS_STS_Endpoint", es.AWSSTSEndpoint)
	plugins.InsertKVString(kvs, "AWS_Role_ARN", es.AWSRoleARN)
	plugins.InsertKVString(kvs, "Cloud_ID", es.CloudID)
	plugins.InsertKVString(kvs, "Cloud_Auth", es.CloudAuth)
	plugins.InsertKVString(kvs, "AWS_External_ID", es.AWSExternalID)
	plugins.InsertKVString(kvs, "Logstash_Prefix", es.LogstashPrefix)
	plugins.InsertKVString(kvs, "Logstash_DateFormat", es.LogstashDateFormat)
	plugins.InsertKVString(kvs, "Time_Key", es.TimeKey)
	plugins.InsertKVString(kvs, "Time_Key_Format", es.TimeKeyFormat)
	plugins.InsertKVString(kvs, "Tag_Key", es.TagKey)
	plugins.InsertKVString(kvs, "ID_KEY", es.IdKey)
	plugins.InsertKVString(kvs, "Write_Operation", es.WriteOperation)
	plugins.InsertKVString(kvs, "Logstash_Prefix_Key", es.LogstashPrefixKey)
	plugins.InsertKVString(kvs, "Suppress_Type_Name", es.SuppressTypeName)
	plugins.InsertKVString(kvs, "storage.total_limit_size", es.TotalLimitSize)

	plugins.InsertKVField(kvs, "Logstash_Format", es.LogstashFormat)
	plugins.InsertKVField(kvs, "Time_Key_Nanos", es.TimeKeyNanos)
	plugins.InsertKVField(kvs, "Include_Tag_Key", es.IncludeTagKey)
	plugins.InsertKVField(kvs, "Generate_ID", es.GenerateID)
	plugins.InsertKVField(kvs, "Replace_Dots", es.ReplaceDots)
	plugins.InsertKVField(kvs, "Trace_Output", es.TraceOutput)
	plugins.InsertKVField(kvs, "Trace_Error", es.TraceError)
	plugins.InsertKVField(kvs, "Current_Time_Index", es.CurrentTimeIndex)

	if es.TLS != nil {
		tls, err := es.TLS.Params(sl)
		if err != nil {
			return nil, err
		}
		kvs.Merge(tls)
	}
	if es.Networking != nil {
		net, err := es.Networking.Params(sl)
		if err != nil {
			return nil, err
		}
		kvs.Merge(net)
	}

	return kvs, nil
}
