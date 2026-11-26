package input

import (
	"testing"

	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins/params"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestInput_Kafka_Params(t *testing.T) {
	g := NewGomegaWithT(t)

	sl := plugins.NewSecretLoader(nil, "test namespace")

	kafka := Kafka{
		Tag:              "kafka.logs",
		Brokers:          "kafka-0:9092,kafka-1:9092",
		Topics:           "app-logs,audit-logs",
		GroupID:          "fluent-bit-consumers",
		ClientID:         "fluent-bit-node",
		Format:           "json",
		PollMs:           new(int32(100)),
		PollTimeoutMs:    new(int32(150)),
		BufferMaxSize:    "8M",
		EnableAutoCommit: new(true),
		Threaded:         new(true),
		Rdkafka: map[string]string{
			"fetch.wait.max.ms": "100",
			"auto.offset.reset": "earliest",
		},
	}

	expected := params.NewKVs()
	expected.Insert("tag", "kafka.logs")
	expected.Insert("brokers", "kafka-0:9092,kafka-1:9092")
	expected.Insert("topics", "app-logs,audit-logs")
	expected.Insert("group_id", "fluent-bit-consumers")
	expected.Insert("client_id", "fluent-bit-node")
	expected.Insert("format", "json")
	expected.Insert("poll_ms", "100")
	expected.Insert("poll_timeout_ms", "150")
	expected.Insert("buffer_max_size", "8M")
	expected.Insert("enable_auto_commit", "true")
	expected.Insert("threaded", "true")
	expected.Insert("rdkafka.auto.offset.reset", "earliest")
	expected.Insert("rdkafka.fetch.wait.max.ms", "100")

	kvs, err := kafka.Params(sl)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(kvs).To(Equal(expected))
}

func TestInput_Kafka_ParamsOnlySetFields(t *testing.T) {
	g := NewGomegaWithT(t)

	sl := plugins.NewSecretLoader(nil, "test namespace")

	kafka := Kafka{
		Brokers:          "kafka:9092",
		Topics:           "logs",
		PollTimeoutMs:    new(int32(0)),
		EnableAutoCommit: new(false),
	}

	expected := params.NewKVs()
	expected.Insert("brokers", "kafka:9092")
	expected.Insert("topics", "logs")
	expected.Insert("poll_timeout_ms", "0")
	expected.Insert("enable_auto_commit", "false")

	kvs, err := kafka.Params(sl)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(kvs).To(Equal(expected))
}

func TestInput_Kafka_ParamsRdkafkaSecrets(t *testing.T) {
	g := NewGomegaWithT(t)

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
		Topics:  "logs",
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
	expected.Insert("brokers", "kafka:9092")
	expected.Insert("topics", "logs")
	expected.Insert("rdkafka.sasl.mechanism", "PLAIN")
	expected.Insert("rdkafka.sasl.password", "api-secret")
	expected.Insert("rdkafka.sasl.username", "api-key")
	expected.Insert("rdkafka.security.protocol", "SASL_SSL")

	kvs, err := kafka.Params(sl)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(kvs).To(Equal(expected))
	g.Expect(kafka.Rdkafka).To(HaveKeyWithValue("sasl.password", "overridden"))

	kafka.RdkafkaSecrets["sasl.username"] = secretRef("missing")
	kvs, err = kafka.Params(sl)
	g.Expect(err).To(MatchError(ContainSubstring(`failed to load rdkafka property "sasl.username" from secret`)))
	g.Expect(kvs).To(BeNil())
}
