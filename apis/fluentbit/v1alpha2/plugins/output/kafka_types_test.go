package output

import (
	"testing"

	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins/params"
	. "github.com/onsi/gomega"
)

func TestOutput_Kafka_Params(t *testing.T) {
	g := NewWithT(t)

	sl := plugins.NewSecretLoader(nil, "test namespace")

	kafka := Kafka{
		Format:    "raw",
		RawLogKey: "message",
		Brokers:   "kafka:9092",
		Topics:    "logs",
	}

	expected := params.NewKVs()
	expected.Insert("Format", "raw")
	expected.Insert("Raw_Log_Key", "message")
	expected.Insert("Brokers", "kafka:9092")
	expected.Insert("Topics", "logs")

	kvs, err := kafka.Params(sl)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(kvs).To(Equal(expected))
}

func TestOutput_Kafka_ParamsRequiresRawLogKeyForRawFormat(t *testing.T) {
	g := NewWithT(t)

	sl := plugins.NewSecretLoader(nil, "test namespace")

	kafka := Kafka{
		Format: "RAW",
	}

	kvs, err := kafka.Params(sl)
	g.Expect(err).To(MatchError("rawLogKey is required when format is raw"))
	g.Expect(kvs).To(BeNil())
}
