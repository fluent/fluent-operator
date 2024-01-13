package output

// S3 defines the parameters for out_s3 output plugin
type S3 struct {
	// The AWS access key id.
	AwsKeyId *string `json:"awsKeyId,omitempty"`
	// The AWS secret key.
	AwsSecKey *string `json:"awsSecKey,omitempty"`
	// The Amazon S3 bucket name.
	S3Bucket *string `json:"s3Bucket,omitempty"`
	// The Amazon S3 region name
	S3Region *string `json:"s3Region,omitempty"`
	// The endpoint URL (like "http://localhost:9000/")
	S3Endpoint *string `json:"s3Endpoint,omitempty"`
	// This prevents AWS SDK from breaking endpoint URL
	ForcePathStyle *bool `json:"forcePathStyle,omitempty"`
	// This timestamp is added to each file name
	TimeSliceFormat *string `json:"timeSliceFormat,omitempty"`
	// The path prefix of the files on S3.
	Path *string `json:"path,omitempty"`
	// The actual S3 path. This is interpolated to the actual path.
	S3ObjectKeyFormat *string `json:"s3ObjectKeyFormat,omitempty"`
	// The compression type.
	// +kubebuilder:validation:Enum:= gzip;lzo;json;txt
	StoreAs *string `json:"storeAs,omitempty"`
	// The proxy URL.
	ProxyUri *string `json:"proxyUri,omitempty"`
	// Verify the SSL certificate of the endpoint.
	SslVerifyPeer *bool `json:"sslVerifyPeer,omitempty"`
	// the following parameters are for S3 kms https://docs.aws.amazon.com/AmazonS3/latest/userguide/UsingKMSEncryption.html
	UseServerSideEncryption *string `json:"useServerSideEncryption,omitempty"`
	// The AWS KMS enctyption algorithm.
	SseCustomerAlgorithm *string `json:"sseCustomerAlgorithm,omitempty"`
	// The AWS KMS key ID.
	SsekmsKeyId *string `json:"ssekmsKeyId,omitempty"`
	// The AWS KMS key.
	SseCustomerKey *string `json:"sseCustomerKey,omitempty"`
	// The AWS KMS key MD5.
	SseCustomerKeyMd5 *string `json:"sseCustomerKeyMd5,omitempty"`
}
