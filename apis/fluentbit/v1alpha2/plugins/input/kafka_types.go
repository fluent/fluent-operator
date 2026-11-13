package input

import (
	"fmt"

	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins/params"
)

// +kubebuilder:object:generate:=true

// The Kafka input plugin subscribes to one or more Kafka topics and ingests the messages as records. <br />
// **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/inputs/kafka**
type Kafka struct {
	// Tag name associated to all records coming from this plugin.
	Tag string `json:"tag,omitempty"`
	// Single entry or list of comma-separated Kafka brokers, e.g: 192.168.1.3:9092,192.168.1.4:9092.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	Brokers string `json:"brokers"`
	// Single entry or list of comma-separated topics that Fluent Bit will subscribe to.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	Topics string `json:"topics"`
	// Kafka consumer group ID, defaults to fluent-bit.
	GroupID string `json:"groupId,omitempty"`
	// Kafka client ID, defaults to fluent-bit.
	ClientID string `json:"clientId,omitempty"`
	// Format of the message payload. none stores the payload as a string, json parses it into a map
	// and falls back to a string if it is not valid JSON. Defaults to none.
	// +kubebuilder:validation:Enum:=none;json
	Format string `json:"format,omitempty"`
	// Time in milliseconds between checks for new messages, defaults to 500.
	// +kubebuilder:validation:Minimum:=1
	PollMs *int32 `json:"pollMs,omitempty"`
	// Timeout in milliseconds for each Kafka consumer poll, defaults to 1.
	// Only takes effect when threaded is enabled. A higher value, such as 1.5x to 2x rdkafka.fetch.wait.max.ms,
	// lets the consumer benefit from Kafka's batching.
	// +kubebuilder:validation:Minimum:=0
	PollTimeoutMs *int32 `json:"pollTimeoutMs,omitempty"`
	// Maximum amount of data fetched from Kafka and ingested in one poll, defaults to 4M.
	// Set it to 0 for no limit.
	// +kubebuilder:validation:Pattern:="^\\d+(k|K|KB|kb|m|M|MB|mb|g|G|GB|gb)?$"
	BufferMaxSize string `json:"bufferMaxSize,omitempty"`
	// Rely on Kafka auto-commit and commit messages in batches instead of committing each message, defaults to false.
	EnableAutoCommit *bool `json:"enableAutoCommit,omitempty"`
	// Run the input in its own thread instead of the main event loop, see pollTimeoutMs.
	Threaded *bool `json:"threaded,omitempty"`
	// Librdkafka properties of the consumer, for example auto.offset.reset or security.protocol.
	Rdkafka map[string]string `json:"rdkafka,omitempty"`
	// Librdkafka properties whose values are read from a Secret, for example sasl.username or sasl.password.
	// A property set in both rdkafka and rdkafkaSecrets takes its value from rdkafkaSecrets.
	RdkafkaSecrets map[string]plugins.Secret `json:"rdkafkaSecrets,omitempty"`
}

func (*Kafka) Name() string {
	return "kafka"
}

// Params implement Section() method
func (k *Kafka) Params(sl plugins.SecretLoader) (*params.KVs, error) {
	kvs := params.NewKVs()

	plugins.InsertKVString(kvs, "tag", k.Tag)
	plugins.InsertKVString(kvs, "brokers", k.Brokers)
	plugins.InsertKVString(kvs, "topics", k.Topics)
	plugins.InsertKVString(kvs, "group_id", k.GroupID)
	plugins.InsertKVString(kvs, "client_id", k.ClientID)
	plugins.InsertKVString(kvs, "format", k.Format)
	plugins.InsertKVField(kvs, "poll_ms", k.PollMs)
	plugins.InsertKVField(kvs, "poll_timeout_ms", k.PollTimeoutMs)
	plugins.InsertKVString(kvs, "buffer_max_size", k.BufferMaxSize)
	plugins.InsertKVField(kvs, "enable_auto_commit", k.EnableAutoCommit)
	plugins.InsertKVField(kvs, "threaded", k.Threaded)

	rdkafka := k.Rdkafka
	if len(k.RdkafkaSecrets) > 0 {
		rdkafka = make(map[string]string, len(k.Rdkafka)+len(k.RdkafkaSecrets))
		for name, value := range k.Rdkafka {
			rdkafka[name] = value
		}
		for name, secret := range k.RdkafkaSecrets {
			value, err := sl.LoadSecret(secret)
			if err != nil {
				return nil, fmt.Errorf("failed to load rdkafka property %q from secret: %w", name, err)
			}
			rdkafka[name] = value
		}
	}
	kvs.InsertStringMap(rdkafka, func(k, v string) (string, string) {
		return fmt.Sprintf("rdkafka.%s", k), v
	})

	return kvs, nil
}
