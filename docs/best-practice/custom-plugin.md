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

```yaml
apiVersion: fluentd.fluent.io/v1alpha1
kind: ClusterOutput
metadata:
  name: cluster-fluentd-output-os
  labels:
    output.fluentd.fluent.io/scope: "cluster"
    output.fluentd.fluent.io/enabled: "true"
spec:
  outputs:
    - customPlugin:
        config: |
          <match **>
            @type opensearch
            host opensearch-logging-data.kubesphere-logging-system.svc
            port 9200
            logstash_format  true
            logstash_prefix  ks-logstash-log
          </match>
```