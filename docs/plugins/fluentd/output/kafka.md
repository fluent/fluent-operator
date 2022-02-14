# Kafka2

Kafka2 defines the parameters for out_kafka output plugin


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| brokers | The list of all seed brokers, with their host and port information. Default: localhost:9092 | *string |
| topicKey | The field name for the target topic. If the field value is app, this plugin writes events to the app topic. | *string |
| defaultTopic | The name of the default topic. (default: nil) | *string |
| useEventTime | Set fluentd event time to Kafka's CreateTime. | *bool |
| requiredAcks | The number of acks required per request. | *int16 |
| compressionCodec | The codec the producer uses to compress messages (default: nil). | *string |
