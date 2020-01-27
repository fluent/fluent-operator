# Kafka




| Field | Description | Scheme |
| ----- | ----------- | ------ |
| format | Specify data format, options available: json, msgpack. | string |
| messageKey | Optional key to store the message | string |
| messageKeyField | If set, the value of Message_Key_Field in the record will indicate the message key. If not set nor found in the record, Message_Key will be used (if set). | string |
| timestampKey | Set the key to store the record timestamp | string |
| timestampFormat | iso8601 or double | string |
| brokers | Single of multiple list of Kafka Brokers, e.g: 192.168.1.3:9092, 192.168.1.4:9092. | string |
| topics | Single entry or list of topics separated by comma (,) that Fluent Bit will use to send messages to Kafka. If only one topic is set, that one will be used for all records. Instead if multiple topics exists, the one set in the record by Topic_Key will be used. | string |
| topicKey | If multiple Topics exists, the value of Topic_Key in the record will indicate the topic to use. E.g: if Topic_Key is router and the record is {\"key1\": 123, \"router\": \"route_2\"}, Fluent Bit will use topic route_2. Note that if the value of Topic_Key is not present in Topics, then by default the first topic in the Topics list will indicate the topic to be used. | string |
| rdkafka | {property} can be any librdkafka properties | map[string]string |
