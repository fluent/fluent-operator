package output

import (
	"testing"

	"github.com/fluent/fluent-operator/v2/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v2/apis/fluentbit/v1alpha2/plugins/params"
	. "github.com/onsi/gomega"
)

func TestOutput_S3_Params(t *testing.T) {
	g := NewGomegaWithT(t)

	sl := plugins.NewSecretLoader(nil, "test namespace")

	s3 := S3{
		Region:                   "us-east-1",
		Bucket:                   "fluentbit",
		JsonDateKey:              "2018-05-30T09:39:52.000681Z",
		JsonDateFormat:           "iso8601",
		TotalFileSize:            "100M",
		UploadChunkSize:          "50M",
		UploadTimeout:            "10m",
		StoreDir:                 "/tmp/fluent-bit/s3",
		StoreDirLimitSize:        "0",
		S3KeyFormat:              "/fluent-bit-logs/$TAG/%Y/%m/%d/%H/%M/%S",
		S3KeyFormatTagDelimiters: ".",
		StaticFilePath:           ptrAny(false),
		UsePutObject:             ptrAny(false),
		RoleArn:                  "role",
		Endpoint:                 "endpoint",
		StsEndpoint:              "sts_endpoint",
		CannedAcl:                "canned_acl",
		Compression:              "gzip",
		ContentType:              "text/plain",
		SendContentMd5:           ptrAny(false),
		AutoRetryRequests:        ptrAny(true),
		LogKey:                   "log_key",
		PreserveDataOrdering:     ptrAny(true),
		StorageClass:             "storage_class",
		RetryLimit:               ptrAny(int32(1)),
		ExternalId:               "external_id",
		Profile:                  "my-profile",
	}

	expected := params.NewKVs()
	expected.Insert("region", "us-east-1")
	expected.Insert("bucket", "fluentbit")
	expected.Insert("json_date_key", "2018-05-30T09:39:52.000681Z")
	expected.Insert("json_date_format", "iso8601")
	expected.Insert("total_file_size", "100M")
	expected.Insert("upload_chunk_size", "50M")
	expected.Insert("upload_timeout", "10m")
	expected.Insert("store_dir", "/tmp/fluent-bit/s3")
	expected.Insert("store_dir_limit_size", "0")
	expected.Insert("s3_key_format", "/fluent-bit-logs/$TAG/%Y/%m/%d/%H/%M/%S")
	expected.Insert("s3_key_format_tag_delimiters", ".")
	expected.Insert("static_file_path", "false")
	expected.Insert("use_put_object", "false")
	expected.Insert("role_arn", "role")
	expected.Insert("endpoint", "endpoint")
	expected.Insert("sts_endpoint", "sts_endpoint")
	expected.Insert("canned_acl", "canned_acl")
	expected.Insert("compression", "gzip")
	expected.Insert("content_type", "text/plain")
	expected.Insert("send_content_md5", "false")
	expected.Insert("auto_retry_requests", "true")
	expected.Insert("log_key", "log_key")
	expected.Insert("preserve_data_ordering", "true")
	expected.Insert("storage_class", "storage_class")
	expected.Insert("retry_limit", "1")
	expected.Insert("external_id", "external_id")
	expected.Insert("profile", "my-profile")

	kvs, err := s3.Params(sl)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(kvs).To(Equal(expected))

}
