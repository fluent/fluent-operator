package output

import (
	"github.com/fluent/fluent-operator/v2/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v2/apis/fluentbit/v1alpha2/plugins/params"
	"github.com/onsi/gomega"
	"testing"
)

func TestOutput_Kinesis_Params(t *testing.T) {
	g := gomega.NewWithT(t)

	sl := plugins.NewSecretLoader(nil, "test namespace")

	ki := Kinesis{
		Region:            "us-east-1",
		Stream:            "test_stream",
		TimeKey:           "test_time_key",
		TimeKeyFormat:     "%Y-%m-%dT%H:%M:%S.%3N",
		LogKey:            "test_time_key",
		RoleARN:           "arn:aws:iam:test",
		Endpoint:          "test_endpoint",
		STSEndpoint:       "test_sts_endpoint",
		AutoRetryRequests: ptrBool(true),
		ExternalID:        "test_external_id",
	}

	expected := params.NewKVs()
	expected.Insert("region", "us-east-1")
	expected.Insert("stream", "test_stream")
	expected.Insert("time_key", "test_time_key")
	expected.Insert("time_key_format", "%Y-%m-%dT%H:%M:%S.%3N")
	expected.Insert("log_key", "test_time_key")
	expected.Insert("role_arn", "arn:aws:iam:test")
	expected.Insert("endpoint", "test_endpoint")
	expected.Insert("sts_endpoint", "test_sts_endpoint")
	expected.Insert("auto_retry_requests", "true")
	expected.Insert("external_id", "test_external_id")

	kvs, err := ki.Params(sl)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(kvs).To(gomega.Equal(expected))
}
