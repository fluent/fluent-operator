# Kafka

The Kafka input plugin subscribes to one or more Kafka topics and ingests the messages as records. <br /> **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/inputs/kafka**


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| tag | Tag name associated to all records coming from this plugin. | string |
| brokers | Single entry or list of comma-separated Kafka brokers, e.g: 192.168.1.3:9092,192.168.1.4:9092. | string |
| topics | Single entry or list of comma-separated topics that Fluent Bit will subscribe to. | string |
| groupId | Kafka consumer group ID, defaults to fluent-bit. | string |
| clientId | Kafka client ID, defaults to fluent-bit. | string |
| format | Format of the message payload. none stores the payload as a string, json parses it into a map and falls back to a string if it is not valid JSON. Defaults to none. | string |
| pollMs | Time in milliseconds between checks for new messages, defaults to 500. | *int32 |
| pollTimeoutMs | Timeout in milliseconds for each Kafka consumer poll, defaults to 1. Only takes effect when threaded is enabled. A higher value, such as 1.5x to 2x rdkafka.fetch.wait.max.ms, lets the consumer benefit from Kafka's batching. | *int32 |
| bufferMaxSize | Maximum amount of data fetched from Kafka and ingested in one poll, defaults to 4M. Set it to 0 for no limit. | string |
| enableAutoCommit | Rely on Kafka auto-commit and commit messages in batches instead of committing each message, defaults to false. | *bool |
| threaded | Run the input in its own thread instead of the main event loop, see pollTimeoutMs. | *bool |
| rdkafka | Librdkafka properties of the consumer, for example auto.offset.reset or security.protocol. | map[string]string |
| rdkafkaSecrets | Librdkafka properties whose values are read from a Secret, for example sasl.username or sasl.password. A property set in both rdkafka and rdkafkaSecrets takes its value from rdkafkaSecrets. | map[string][plugins.Secret](../secret.md) |
