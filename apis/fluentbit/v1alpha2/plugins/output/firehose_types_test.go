package output

import (
	"testing"

	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins/params"
	"github.com/fluent/fluent-operator/v3/pkg/utils"
	"github.com/onsi/gomega"
)

func TestOutput_Firehose_Params(t *testing.T) {
	g := gomega.NewWithT(t)

	sl := plugins.NewSecretLoader(nil, "test namespace")

	fh := Firehose{
		Region:            "us-east-1",
		DeliveryStream:    "test_stream",
		TimeKey:           utils.ToPtr("test_time_key"),
		TimeKeyFormat:     utils.ToPtr("%Y-%m-%dT%H:%M:%S.%3N"),
		DataKeys:          utils.ToPtr("test_data_keys"),
		LogKey:            utils.ToPtr("test_time_key"),
		RoleARN:           utils.ToPtr("arn:aws:iam:test"),
		Endpoint:          utils.ToPtr("test_endpoint"),
		STSEndpoint:       utils.ToPtr("test_sts_endpoint"),
		AutoRetryRequests: utils.ToPtr(true),
		ExternalID:        utils.ToPtr("test_external_id"),
		Compression:       utils.ToPtr("gzip"),
		SimpleAggregation: utils.ToPtr(true),
		Profile:           utils.ToPtr("my-profile"),
		Workers:           utils.ToPtr[int32](1),
	}

	expected := params.NewKVs()
	expected.Insert("region", "us-east-1")
	expected.Insert("delivery_stream", "test_stream")
	expected.Insert("data_keys", "test_data_keys")
	expected.Insert("log_key", "test_time_key")
	expected.Insert("role_arn", "arn:aws:iam:test")
	expected.Insert("endpoint", "test_endpoint")
	expected.Insert("sts_endpoint", "test_sts_endpoint")
	expected.Insert("time_key", "test_time_key")
	expected.Insert("time_key_format", "%Y-%m-%dT%H:%M:%S.%3N")
	expected.Insert("auto_retry_requests", "true")
	expected.Insert("compression", "gzip")
	expected.Insert("external_id", "test_external_id")
	expected.Insert("simple_aggregation", "true")
	expected.Insert("profile", "my-profile")
	expected.Insert("workers", "1")

	kvs, err := fh.Params(sl)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(kvs).To(gomega.Equal(expected))
}
