package output

import (
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins/params"
)

// +kubebuilder:object:generate:=true

// The S3 output plugin, allows to flush your records into a S3 time series database. <br />
// **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/outputs/s3**
type S3 struct {
	// The AWS region of your S3 bucket
	Region string `json:"Region"`
	// S3 Bucket name
	Bucket string `json:"Bucket"`
	// Specify the name of the time key in the output record. To disable the time key just set the value to false.
	JsonDateKey string `json:"JsonDateKey,omitempty"`
	// Specify the format of the date. Supported formats are double, epoch, iso8601 (eg: 2018-05-30T09:39:52.000681Z) and java_sql_timestamp (eg: 2018-05-30 09:39:52.000681)
	JsonDateFormat string `json:"JsonDateFormat,omitempty"`
	// Specifies the size of files in S3. Minimum size is 1M. With use_put_object On the maximum size is 1G. With multipart upload mode, the maximum size is 50G.
	TotalFileSize string `json:"TotalFileSize,omitempty"`
	// The size of each 'part' for multipart uploads. Max: 50M
	UploadChunkSize string `json:"UploadChunkSize,omitempty"`
	// Whenever this amount of time has elapsed, Fluent Bit will complete an upload and create a new file in S3. For example, set this value to 60m and you will get a new file every hour.
	UploadTimeout string `json:"UploadTimeout,omitempty"`
	// Directory to locally buffer data before sending.
	StoreDir string `json:"StoreDir,omitempty"`
	// The size of the limitation for disk usage in S3.
	StoreDirLimitSize string `json:"StoreDirLimitSize,omitempty"`
	// Format string for keys in S3.
	S3KeyFormat string `json:"S3KeyFormat,omitempty"`
	// A series of characters which will be used to split the tag into 'parts' for use with the s3_key_format option.
	S3KeyFormatTagDelimiters string `json:"S3KeyFormatTagDelimiters,omitempty"`
	// Disables behavior where UUID string is automatically appended to end of S3 key name when $UUID is not provided in s3_key_format. $UUID, time formatters, $TAG, and other dynamic key formatters all work as expected while this feature is set to true.
	StaticFilePath *bool `json:"StaticFilePath,omitempty"`
	// Use the S3 PutObject API, instead of the multipart upload API.
	UsePutObject *bool `json:"UsePutObject,omitempty"`
	// ARN of an IAM role to assume
	RoleArn string `json:"RoleArn,omitempty"`
	// Custom endpoint for the S3 API.
	Endpoint string `json:"Endpoint,omitempty"`
	// Custom endpoint for the STS API.
	StsEndpoint string `json:"StsEndpoint,omitempty"`
	// Predefined Canned ACL Policy for S3 objects.
	CannedAcl string `json:"CannedAcl,omitempty"`
	// Compression type for S3 objects.
	Compression string `json:"Compression,omitempty"`
	// A standard MIME type for the S3 object; this will be set as the Content-Type HTTP header.
	ContentType string `json:"ContentType,omitempty"`
	// Send the Content-MD5 header with PutObject and UploadPart requests, as is required when Object Lock is enabled.
	SendContentMd5 *bool `json:"SendContentMd5,omitempty"`
	// Immediately retry failed requests to AWS services once.
	AutoRetryRequests *bool `json:"AutoRetryRequests,omitempty"`
	// By default, the whole log record will be sent to S3. If you specify a key name with this option, then only the value of that key will be sent to S3.
	LogKey string `json:"LogKey,omitempty"`
	// Normally, when an upload request fails, there is a high chance for the last received chunk to be swapped with a later chunk, resulting in data shuffling. This feature prevents this shuffling by using a queue logic for uploads.
	PreserveDataOrdering *bool `json:"PreserveDataOrdering,omitempty"`
	// Specify the storage class for S3 objects. If this option is not specified, objects will be stored with the default 'STANDARD' storage class.
	StorageClass string `json:"StorageClass,omitempty"`
	// Integer value to set the maximum number of retries allowed.
	RetryLimit *int32 `json:"RetryLimit,omitempty"`
	// Specify an external ID for the STS API, can be used with the role_arn parameter if your role requires an external ID.
	ExternalId string `json:"ExternalId,omitempty"`
	// Option to specify an AWS Profile for credentials.
	Profile string `json:"Profile,omitempty"`
	// Specify number of worker threads to use to output to S3
	Workers *int32 `json:"Workers,omitempty"`
	// Set maximum time expressed in seconds to wait for a TCP connection to be established, this include the TLS handshake time.
	ConnectTimeout *int32 `json:"ConnectTimeout,omitempty"`
	// On connection timeout, specify if it should log an error. When disabled, the timeout is logged as a debug message.
	ConnectTimeoutLogError *bool `json:"connectTimeoutLogError,omitempty"`
	// Select the primary DNS connection type (TCP or UDP).
	// +kubebuilder:validation:Enum:="TCP";"UDP"
	DNSMode *string `json:"DNSMode,omitempty"`
	// Prioritize IPv4 DNS results when trying to establish a connection.
	DNSPreferIPv4 *bool `json:"DNSPreferIPv4,omitempty"`
	// Prioritize IPV6 DNS results when trying to establish a connection.
	DNSPreferIPv6 *bool `json:"DNSPreferIPv6,omitempty"`
	// Set maximum time a connection can stay idle while assigned.
	IoTimeout *int32 `json:"IoTimeout,omitempty"`
	// Set maximum number of times a keepalive connection can be used before it is retired.
	KeepaliveMaxRecycle *int32 `json:"keepaliveMaxRecycle,omitempty"`
	// Set maximum number of TCP connections that can be established per worker.
	MaxWorkerConnections *int32 `json:"maxWorkerConnections,omitempty"`
	// Ignore the environment variables HTTP_PROXY, HTTPS_PROXY and NO_PROXY when set.
	ProxyEnvIgnore *bool `json:"proxyEnvIgnore,omitempty"`
	// Specify network address to bind for data traffic.
	SourceAddress *string `json:"sourceAddress,omitempty"`
	// Enable or disable connection keepalive support. Accepts a boolean value: on / off.
	// +kubebuilder:validation:Enum:="on";"off"
	Keepalive *string `json:"keepalive,omitempty"`
	// Interval between TCP keepalive probes when no response is received on a keepidle probe.
	TCPKeepaliveInterval *int32 `json:"tcpKeepaliveInterval,omitempty"`
	// Number of unacknowledged probes to consider a connection dead.
	TCPKeepaliveProbes *int32 `json:"tcpKeepaliveProbes,omitempty"`
	// Interval between the last data packet sent and the first TCP keepalive probe.
	TCPKeepaliveTime *int32 `json:"tcpKeepaliveTime,omitempty"`
	// Set maximum time expressed in seconds for an idle keepalive connection.
	KeepaliveIdleTimeout *int32 `json:"keepaliveIdleTimeout,omitempty"`
	*plugins.TLS         `json:"tls,omitempty"`
}

// Name implement Section() method
func (*S3) Name() string {
	return "s3"
}

func (o *S3) Params(sl plugins.SecretLoader) (*params.KVs, error) {
	kvs := params.NewKVs()

	plugins.InsertKVString(kvs, "region", o.Region)
	plugins.InsertKVString(kvs, "bucket", o.Bucket)
	plugins.InsertKVString(kvs, "json_date_key", o.JsonDateKey)
	plugins.InsertKVString(kvs, "json_date_format", o.JsonDateFormat)
	plugins.InsertKVString(kvs, "total_file_size", o.TotalFileSize)
	plugins.InsertKVString(kvs, "upload_chunk_size", o.UploadChunkSize)
	plugins.InsertKVString(kvs, "upload_timeout", o.UploadTimeout)
	plugins.InsertKVString(kvs, "store_dir", o.StoreDir)
	plugins.InsertKVString(kvs, "store_dir_limit_size", o.StoreDirLimitSize)
	plugins.InsertKVString(kvs, "s3_key_format", o.S3KeyFormat)
	plugins.InsertKVString(kvs, "s3_key_format_tag_delimiters", o.S3KeyFormatTagDelimiters)
	plugins.InsertKVField(kvs, "static_file_path", o.StaticFilePath)
	plugins.InsertKVField(kvs, "use_put_object", o.UsePutObject)
	plugins.InsertKVString(kvs, "role_arn", o.RoleArn)
	plugins.InsertKVString(kvs, "endpoint", o.Endpoint)
	plugins.InsertKVString(kvs, "sts_endpoint", o.StsEndpoint)
	plugins.InsertKVString(kvs, "canned_acl", o.CannedAcl)
	plugins.InsertKVString(kvs, "compression", o.Compression)
	plugins.InsertKVString(kvs, "content_type", o.ContentType)
	plugins.InsertKVField(kvs, "send_content_md5", o.SendContentMd5)
	plugins.InsertKVField(kvs, "auto_retry_requests", o.AutoRetryRequests)
	plugins.InsertKVString(kvs, "log_key", o.LogKey)
	plugins.InsertKVField(kvs, "preserve_data_ordering", o.PreserveDataOrdering)
	plugins.InsertKVString(kvs, "storage_class", o.StorageClass)
	plugins.InsertKVField(kvs, "retry_limit", o.RetryLimit)
	plugins.InsertKVString(kvs, "external_id", o.ExternalId)
	plugins.InsertKVString(kvs, "profile", o.Profile)
	plugins.InsertKVField(kvs, "workers", o.Workers)
	plugins.InsertKVField(kvs, "net.connect_timeout", o.ConnectTimeout)
	plugins.InsertKVField(kvs, "net.connect_timeout_log_error", o.ConnectTimeoutLogError)
	plugins.InsertKVField(kvs, "net.dns.mode", o.DNSMode)
	plugins.InsertKVField(kvs, "net.dns.prefer_ipv4", o.DNSPreferIPv4)
	plugins.InsertKVField(kvs, "net.dns.prefer_ipv6", o.DNSPreferIPv6)
	plugins.InsertKVField(kvs, "net.io_timeout", o.IoTimeout)
	plugins.InsertKVField(kvs, "net.keepalive_max_recycle", o.KeepaliveMaxRecycle)
	plugins.InsertKVField(kvs, "net.max_worker_connections", o.MaxWorkerConnections)
	plugins.InsertKVField(kvs, "net.proxy_env_ignore", o.ProxyEnvIgnore)
	plugins.InsertKVField(kvs, "net.source_address", o.SourceAddress)
	plugins.InsertKVField(kvs, "net.tcp_keepalive", o.Keepalive)
	plugins.InsertKVField(kvs, "net.tcp_keepalive_interval", o.TCPKeepaliveInterval)
	plugins.InsertKVField(kvs, "net.tcp_keepalive_probes", o.TCPKeepaliveProbes)
	plugins.InsertKVField(kvs, "net.tcp_keepalive_time", o.TCPKeepaliveTime)
	plugins.InsertKVField(kvs, "net.keepalive_idle_timeout", o.KeepaliveIdleTimeout)

	if o.TLS != nil {
		tls, err := o.TLS.Params(sl)
		if err != nil {
			return nil, err
		}
		kvs.Merge(tls)
	}

	return kvs, nil
}
