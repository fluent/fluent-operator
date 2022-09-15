# Custom Plugin


You can use custom plugin configuration to support plugins that fluent-operator does not currently support. The following are some examples of how to use the custom plugin configuration.

```yaml 
apiVersion: fluentbit.fluent.io/v1alpha2
kind: ClusterInput
metadata:
  namespace: fluent
  name: cpu-input
  labels:
    fluentbit.fluent.io/enabled: "true"
    fluentbit.fluent.io/mode: "fluentbit-only"
spec:
  customPlugin:
    config: |
      Name    cpu
      Tag    my_cpu
```

```yaml
apiVersion: fluentbit.fluent.io/v1alpha2
kind: ClusterOutput
metadata:
  namespace: fluent
  name: kafka-output
  labels:
    fluentbit.fluent.io/enabled: "true"
    fluentbit.fluent.io/mode: "fluentbit-only"
spec:
  customPlugin:
    config: |
      Name    kafka
      Topics     fluentbit
      Match       *
      Brokers     192.168.100.32:9092
      rdkafka.debug All
      rdkafka.request.required.acks 1
      rdkafka.log.connection.close false
      rdkafka.log_level 7
      rdkafka.metadata.broker.list 192.168.100.32:9092
```