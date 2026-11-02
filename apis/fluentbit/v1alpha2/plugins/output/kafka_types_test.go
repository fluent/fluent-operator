package output

import (
	"testing"

	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins/params"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
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

func TestOutput_Kafka_ParamsRdkafkaSecrets(t *testing.T) {
	g := NewWithT(t)

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{Name: "kafka-credentials", Namespace: "test namespace"},
		Data: map[string][]byte{
			"username": []byte("api-key\n"),
			"password": []byte("api-secret"),
		},
	}
	sl := plugins.NewSecretLoader(fake.NewClientBuilder().WithObjects(secret).Build(), "test namespace")

	secretRef := func(key string) plugins.Secret {
		return plugins.Secret{ValueFrom: plugins.ValueSource{
			SecretKeyRef: corev1.SecretKeySelector{
				LocalObjectReference: corev1.LocalObjectReference{Name: "kafka-credentials"},
				Key:                  key,
			},
		}}
	}

	kafka := Kafka{
		Brokers: "kafka:9092",
		Rdkafka: map[string]string{
			"security.protocol": "SASL_SSL",
			"sasl.mechanism":    "PLAIN",
			"sasl.password":     "overridden",
		},
		RdkafkaSecrets: map[string]plugins.Secret{
			"sasl.username": secretRef("username"),
			"sasl.password": secretRef("password"),
		},
	}

	expected := params.NewKVs()
	expected.Insert("Brokers", "kafka:9092")
	expected.Insert("rdkafka.sasl.mechanism", "PLAIN")
	expected.Insert("rdkafka.sasl.password", "api-secret")
	expected.Insert("rdkafka.sasl.username", "api-key")
	expected.Insert("rdkafka.security.protocol", "SASL_SSL")

	kvs, err := kafka.Params(sl)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(kvs).To(Equal(expected))

	kafka.RdkafkaSecrets["sasl.username"] = secretRef("missing")
	kvs, err = kafka.Params(sl)
	g.Expect(err).To(MatchError(ContainSubstring(`failed to load rdkafka property "sasl.username" from secret`)))
	g.Expect(kvs).To(BeNil())
}
