package cfgrender

import (
	"encoding/json"
	"os"
	"sync"

	corev1 "k8s.io/api/core/v1"

	fluentdv1alpha1 "github.com/fluent/fluent-operator/v2/apis/fluentd/v1alpha1"
	"github.com/fluent/fluent-operator/v2/apis/fluentd/v1alpha1/plugins/common"
	"github.com/fluent/fluent-operator/v2/apis/fluentd/v1alpha1/plugins/filter"
	"github.com/fluent/fluent-operator/v2/apis/fluentd/v1alpha1/plugins/input"
	"github.com/fluent/fluent-operator/v2/apis/fluentd/v1alpha1/plugins/output"
	"sigs.k8s.io/yaml"
)

var (
	Fluentd    fluentdv1alpha1.Fluentd
	FluentdRaw = `
apiVersion: fluentd.fluent.io/v1alpha1
kind: Fluentd
metadata:
  name: fluentd
  namespace: fluent
  labels:
    app.kubernetes.io/name: fluentd
spec:
  globalInputs:
  - forward: 
      bind: 0.0.0.0                
      port: 24224
  replicas: 1
  image: kubesphere/fluentd:v1.14.4
  fluentdCfgSelector: 
    matchLabels:
      config.fluentd.fluent.io/enabled: "true"
`

	FluentdInputSample    fluentdv1alpha1.Fluentd
	FluentdInputSampleRaw = `
apiVersion: fluentd.fluent.io/v1alpha1
kind: Fluentd
metadata:
  name: fluentd
  namespace: fluent
  labels:
    app.kubernetes.io/name: fluentd
spec:
  globalInputs:
  - sample:
      sample: '{"hello": "world"}'
      tag: "foo.bar"
      rate: 10
      size: 10
      autoIncrementKey: "id"
  replicas: 1
  image: kubesphere/fluentd:v1.15.3
  fluentdCfgSelector:
    matchLabels:
      config.fluentd.fluent.io/enabled: "true"
`

	FluentdInputMonitorAgent    fluentdv1alpha1.Fluentd
	FluentdInputMonitorAgentRaw = `
apiVersion: fluentd.fluent.io/v1alpha1
kind: Fluentd
metadata:
name: fluentd
namespace: fluent
labels:
  app.kubernetes.io/name: fluentd
spec:
  globalInputs:
  - monitorAgent:
      bind: 0.0.0.0
      port: 24220
      tag: example
      emitInterval: 5
      includeConfig: true
      includeRetry: true
  replicas: 1
  image: kubesphere/fluentd:v1.15.3
  fluentdCfgSelector:
    matchLabels:
      config.fluentd.fluent.io/enabled: "true"
`

	FluentdInputTail    fluentdv1alpha1.Fluentd
	FluentdInputTailRaw = `
apiVersion: fluentd.fluent.io/v1alpha1
kind: Fluentd
metadata:
  name: fluentd
  namespace: fluent
  labels:
    app.kubernetes.io/name: fluentd
spec:
  globalInputs:
  - tail: 
      tag: "foo.bar"
      path: /var/log/test.log
      emitUnmatchedLines: true
      enableStatWatcher: true
      enableWatchTimer: true
      followInodes: false
      group:
        pattern: '/^\/home\/logs\/(?<file>.+)\.log$/'
        ratePeriod: 30
        rule:
          limit: 2
          match:
            key1: val1
            key2: val2
      ignoreRepeatedPermissionError: false
      limitRecentlyModified: 3
      maxLineSize: 10000
      multilineFlushInterval: 4
      openOnEveryUpdate: false
      parse:
        type: json
      pathKey: tailed_path
      posFile: /fluentd/pos.db
      posFileCompactionInterval: 5
      readBytesLimitPerSecond: 8192
      readFromHead: false
      readLinesLimit: 15
      refreshInterval: 2
      rotateWait: 30
      skipRefreshOnStartup: false
      excludePath:
      - /var/log/foo.log
      - /var/log/bar
  replicas: 1
  image: kubesphere/fluentd:v1.15.3
  fluentdCfgSelector: 
    matchLabels:
      config.fluentd.fluent.io/enabled: "true"
`

	FluentdClusterFluentdConfig1    fluentdv1alpha1.ClusterFluentdConfig
	FluentdClusterFluentdConfig1Raw = `
apiVersion: fluentd.fluent.io/v1alpha1
kind: ClusterFluentdConfig
metadata:
  name: fluentd-config
  labels:
    config.fluentd.fluent.io/enabled: "true"
spec:
  watchedNamespaces: 
  - kube-system
  - default
  clusterFilterSelector:
    matchLabels:
      filter.fluentd.fluent.io/enabled: "true"
  clusterOutputSelector:
    matchLabels:
      output.fluentd.fluent.io/enabled: "true"
`

	FluentdClusterFluentdConfig2    fluentdv1alpha1.ClusterFluentdConfig
	FluentdClusterFluentdConfig2Raw = `
apiVersion: fluentd.fluent.io/v1alpha1
kind: ClusterFluentdConfig
metadata:
  name: fluentd-config-cluster
  labels:
    config.fluentd.fluent.io/enabled: "true"
spec:
  watchedNamespaces: 
  - kube-system
  - default
  clusterOutputSelector:
    matchLabels:
      output.fluentd.fluent.io/enabled: "true"
      output.fluentd.fluent.io/scope: "cluster"
`

	FluentdConfig1    fluentdv1alpha1.FluentdConfig
	FluentdConfig1Raw = `
apiVersion: fluentd.fluent.io/v1alpha1
kind: FluentdConfig
metadata:
  name: fluentd-config
  namespace: fluent
  labels:
    config.fluentd.fluent.io/enabled: "true"
spec:
  clusterOutputSelector:
    matchLabels:
      output.fluentd.fluent.io/enabled: "true"
`

	FluentdConfig2    fluentdv1alpha1.FluentdConfig
	FluentdConfig2Raw = `
apiVersion: fluentd.fluent.io/v1alpha1
kind: FluentdConfig
metadata:
  name: fluentd-config
  namespace: fluent
  labels:
    config.fluentd.fluent.io/enabled: "true"
spec:
  filterSelector:
    matchLabels:
      filter.fluentd.fluent.io/enabled: "true"
  clusterOutputSelector:
    matchLabels:
      output.fluentd.fluent.io/enabled: "loki"
  outputSelector:
    matchLabels:
      output.fluentd.fluent.io/enabled: "es"
`

	FluentdConfigUser1    fluentdv1alpha1.FluentdConfig
	FluentdConfigUser1Raw = `
apiVersion: fluentd.fluent.io/v1alpha1
kind: FluentdConfig
metadata:
  name: fluentd-config-user1
  namespace: fluent
  labels:
    config.fluentd.fluent.io/enabled: "true"
spec:
  outputSelector:
    matchLabels:
      output.fluentd.fluent.io/enabled: "true"
      output.fluentd.fluent.io/user: "user1"
  clusterOutputSelector:
    matchLabels:
      output.fluentd.fluent.io/enabled: "true"
      output.fluentd.fluent.io/role: "log-operator"
`

	FluentdClusterFilter1    fluentdv1alpha1.ClusterFilter
	FluentdClusterFilter1Raw = `
apiVersion: fluentd.fluent.io/v1alpha1
kind: ClusterFilter
metadata:
  name: fluentd-filter
  labels:
    filter.fluentd.fluent.io/enabled: "true"
spec: 
  filters: 
  - recordTransformer:
      enableRuby: true
      records:
      - key: kubernetes_ns
        value: ${record["kubernetes"]["namespace_name"]}
`

	FluentdFilter    fluentdv1alpha1.Filter
	FluentdFilterRaw = `
apiVersion: fluentd.fluent.io/v1alpha1
kind: Filter
metadata:
  name: fluentd-filter
  namespace: fluent
labels:
  filter.fluentd.fluent.io/enabled: "true"
spec: 
  filters: 
  - recordTransformer:
      records:
      - key: loki-tenant
        value: 9705b9de-d60e-4492-b322-eac870
      - key: environment
        value: uat
      - key: application
        value: my-application
`

	FluentdClusterRecordTransformerFilter fluentdv1alpha1.ClusterFilter
	FluentdClusterRecordTransformerRaw    = `
apiVersion: fluentd.fluent.io/v1alpha1
kind: ClusterFilter
metadata:
  name: fluentd-filter
  labels:
    filter.fluentd.fluent.io/enabled: "true"
spec:
  filters:
  - recordTransformer:
      enableRuby: true
      autoTypeCast: true
      renewRecord: true
      records:
      - key: kubernetes_ns
        value: ${record["kubernetes"]["namespace_name"]}
`

	FluentdClusterOutputBuffer    fluentdv1alpha1.ClusterOutput
	FluentdClusterOutputBufferRaw = `
apiVersion: fluentd.fluent.io/v1alpha1
kind: ClusterOutput
metadata:
  name: fluentd-output
  labels:
    output.fluentd.fluent.io/enabled: "true"
spec: 
  outputs: 
  - stdout: {}
    buffer:
      type: file
      path: /buffers/stdout.log
  - elasticsearch:
      host: elasticsearch-logging-data.kubesphere-logging-system.svc
      port: 9200
      logstashFormat: true
      logstashPrefix: ks-logstash-log
    buffer:
      type: file
      path: /buffers/es.log
  `

	FluentdClusterOutputMemoryBuffer    fluentdv1alpha1.ClusterOutput
	FluentdClusterOutputMemoryBufferRaw = `
apiVersion: fluentd.fluent.io/v1alpha1
kind: ClusterOutput
metadata:
  name: fluentd-output
  labels:
    output.fluentd.fluent.io/enabled: "true"
spec: 
  outputs: 
  - stdout: {}
    buffer:
      type: memory
  - elasticsearch:
      host: elasticsearch-logging-data.kubesphere-logging-system.svc
      port: 9200
      logstashFormat: true
      logstashPrefix: ks-logstash-log
    buffer:
      type: file
      path: /buffers/es.log
`

	FluentdClusterOutputTag    fluentdv1alpha1.ClusterOutput
	FluentdClusterOutputTagRaw = `
apiVersion: fluentd.fluent.io/v1alpha1
kind: ClusterOutput
metadata:
  name: fluentd-output-stdout
  labels:
    output.fluentd.fluent.io/enabled: "true"
spec: 
  outputs: 
  - stdout: {}
    tag: foo.*
`

	FluentdclusterOutput2ES    fluentdv1alpha1.ClusterOutput
	FluentdclusterOutput2ESRaw = `
apiVersion: fluentd.fluent.io/v1alpha1
kind: ClusterOutput
metadata:
  name: fluentd-output-es
  labels:
    output.fluentd.fluent.io/enabled: "true"
spec: 
  outputs: 
  - elasticsearch:
      host: elasticsearch-logging-data.kubesphere-logging-system.svc
      port: 9200
      logstashFormat: true
      logstashPrefix: ks-logstash-log
`

	FluentdOutput2ES1    fluentdv1alpha1.Output
	FluentdOutput2ES1Raw = `
apiVersion: fluentd.fluent.io/v1alpha1
kind: Output
metadata:
  name: fluentd-output-es
  namespace: fluent
labels:
  output.fluentd.fluent.io/enabled: "es"
spec: 
  outputs:
  - copy:
      copyMode: no_copy
  - elasticsearch:
      host: elasticsearch-logging-data.kubesphere-logging-system.svc
      indexName: es1-notag-1
      port: 9243
      scheme: https
      sslVerify: false
  - elasticsearch:
      host: elasticsearch-logging-data.kubesphere-logging-system.svc
      indexName: es1-notag-2
      port: 9243
      scheme: https
      sslVerify: false
      user:
        valueFrom:
          secretKeyRef:
            key: username
            name: es-credentials
      password:
        valueFrom:
          secretKeyRef:
            key: password
            name: es-credentials
`
	FluentdOutput2ES2    fluentdv1alpha1.Output
	FluentdOutput2ES2Raw = `
apiVersion: fluentd.fluent.io/v1alpha1
kind: Output
metadata:
  name: fluentd-output-es2
  namespace: fluent
labels:
  output.fluentd.fluent.io/enabled: "es"
spec: 
  outputs:
  - copy:
      copyMode: no_copy
    tag: specific_tag_1
  - elasticsearch:
      host: elasticsearch-logging-data.kubesphere-logging-system.svc
      indexName: es2-specific-tag-1-1
      port: 9243
      scheme: https
      sslVerify: false
    tag: specific_tag_1
  - elasticsearch:
      host: elasticsearch-logging-data.kubesphere-logging-system.svc
      indexName: es2-specific-tag-1-2
      port: 9243
      scheme: https
      sslVerify: false
    tag: specific_tag_1
`
	FluentdOutput2ES3    fluentdv1alpha1.Output
	FluentdOutput2ES3Raw = `
apiVersion: fluentd.fluent.io/v1alpha1
kind: Output
metadata:
  name: fluentd-output-es3
  namespace: fluent
labels:
  output.fluentd.fluent.io/enabled: "es"
spec: 
  outputs:
  - elasticsearch:
      host: elasticsearch-logging-data.kubesphere-logging-system.svc
      indexName: es3-specific-tag-2-1
      port: 9243
      scheme: https
      sslVerify: false
    tag: specific_tag_2
  - elasticsearch:
      host: elasticsearch-logging-data.kubesphere-logging-system.svc
      indexName: es3-specific-tag-2-2
      port: 9243
      scheme: https
      sslVerify: false
    tag: specific_tag_2
`
	FluentdOutput2ES4    fluentdv1alpha1.Output
	FluentdOutput2ES4Raw = `
apiVersion: fluentd.fluent.io/v1alpha1
kind: Output
metadata:
  name: fluentd-output-es4
  namespace: fluent
labels:
  output.fluentd.fluent.io/enabled: "es"
spec: 
  outputs:
  - elasticsearch:
      host: elasticsearch-logging-data.kubesphere-logging-system.svc
      indexName: es4-notag-1
      port: 9243
      scheme: https
      sslVerify: false
`

	FluentdclusterOutput2OpenSearch    fluentdv1alpha1.ClusterOutput
	FluentdclusterOutput2OpenSearchRaw = `
apiVersion: fluentd.fluent.io/v1alpha1
kind: ClusterOutput
metadata:
  name: fluentd-output-opensearch
  labels:
    output.fluentd.fluent.io/enabled: "true"
spec: 
  outputs: 
  - opensearch:
      host: opensearch-logging-data.kubesphere-logging-system.svc
      port: 9200
      logstashFormat: true
      logstashPrefix: ks-logstash-log
`

	FluentdClusterOutput2kafka    fluentdv1alpha1.ClusterOutput
	FluentdClusterOutput2kafkaRaw = `
apiVersion: fluentd.fluent.io/v1alpha1
kind: ClusterOutput
metadata:
  name: fluentd-output-kafka
  labels:
    output.fluentd.fluent.io/enabled: "true"
spec: 
  outputs: 
  - kafka:
      brokers: my-cluster-kafka-bootstrap.default.svc:9091,my-cluster-kafka-bootstrap.default.svc:9092,my-cluster-kafka-bootstrap.default.svc:9093
      useEventTime: true
      topicKey: kubernetes_ns
`
	FluentdClusterOutput2Loki    fluentdv1alpha1.ClusterOutput
	FluentdClusterOutput2LokiRaw = `
apiVersion: fluentd.fluent.io/v1alpha1
kind: ClusterOutput
metadata:
  name: fluentd-output-loki
  labels:
    output.fluentd.fluent.io/enabled: "true"
spec: 
  outputs: 
  - loki:
      url: http://loki-logging-data.kubesphere-logging-system.svc:3100
      extractKubernetesLabels: true
      tenantID:
        valueFrom:
          secretKeyRef:
            key: tenant_key
            name: loki-tenant-name
      httpPassword:
        valueFrom:
          secretKeyRef:
            key: password_key
            name: loki-http-credentials
      httpUser:
        valueFrom:
          secretKeyRef:
            key: user_key
            name: loki-http-credentials
      labels:
        - key11=value11
        - key12 = value12
        - key13
      labelKeys:
        - key21
        - key22
      removeKeys:
        - key31
        - key32
      dropSingleKey: true
      includeThreadLabel: true
      tlsCaCertFile: /path/to/ca.pem
      tlsClientCertFile: /path/to/certificate.pem
      tlsPrivateKeyFile: /path/to/key.key
      insecure: true
`

	lokiHttpCredentials    corev1.Secret
	lokiHttpCredentialsRaw = `
apiVersion: v1
kind: Secret
metadata:
  name: loki-http-credentials
  namespace: fluent
type: Opaque
stringData:
  password_key: s3cr3tP@ssword
  user_key: s3cr3tUsern4me
`

	lokiTenantName    corev1.Secret
	lokiTenantNameRaw = `
apiVersion: v1
kind: Secret
metadata:
  name: loki-tenant-name
  namespace: fluent
type: Opaque
stringData:
  tenant_key: 0c3ba7a4-3148-4605-b62a-afc92dd1c4d7
`

	FluentdClusterOutput2Loki1    fluentdv1alpha1.ClusterOutput
	FluentdClusterOutput2Loki1Raw = `
apiVersion: fluentd.fluent.io/v1alpha1
kind: ClusterOutput
metadata:
  name: fluentd-output-loki
labels:
  output.fluentd.fluent.io/enabled: "loki"
spec: 
  outputs: 
  - loki:
      url: http://loki-logging-data.kubesphere-logging-system.svc:3100
      extractKubernetesLabels: true
      labels:
        - key11=value11
        - key12 = value12
        - key13
      labelKeys:
        - key21
        - key22
      removeKeys:
        - key31
        - key32
      dropSingleKey: true
      includeThreadLabel: true
      insecure: true
`

	FluentdClusterOutputLogOperator    fluentdv1alpha1.ClusterOutput
	FluentdClusterOutputLogOperatorRaw = `
apiVersion: fluentd.fluent.io/v1alpha1
kind: ClusterOutput
metadata:
  name: fluentd-output-log-operator
  labels:
    output.fluentd.fluent.io/enabled: "true"
    output.fluentd.fluent.io/role: "log-operator"
spec: 
  outputs: 
  - elasticsearch:
      host: elasticsearch-logging-data.kubesphere-logging-system.svc
      port: 9200
      logstashFormat: true
      logstashPrefix: ks-logstash-log-operator
`

	FluentdClusterOutputCluster    fluentdv1alpha1.ClusterOutput
	FluentdClusterOutputClusterRaw = `
apiVersion: fluentd.fluent.io/v1alpha1
kind: ClusterOutput
metadata:
  name: fluentd-output-cluster
  labels:
    output.fluentd.fluent.io/enabled: "true"
    output.fluentd.fluent.io/scope: "cluster"
spec: 
  outputs: 
  - elasticsearch:
      host: elasticsearch-logging-data.kubesphere-logging-system.svc
      port: 9200
      logstashFormat: true
      logstashPrefix: ks-logstash-log
`
	FluentdClusterOutputCustom    fluentdv1alpha1.ClusterOutput
	FluentdClusterOutputCustomRaw = `
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
`
	FluentdOutputUser1    fluentdv1alpha1.Output
	FluentdOutputUser1Raw = `
apiVersion: fluentd.fluent.io/v1alpha1
kind: Output
metadata:
  name: fluentd-output-user1
  namespace: fluent
  labels:
    output.fluentd.fluent.io/enabled: "true"
    output.fluentd.fluent.io/user: "user1"
spec:
  outputs:
  - elasticsearch:
      host: elasticsearch-logging-data.kubesphere-logging-system.svc
      port: 9200
      logstashFormat: true
      logstashPrefix: ks-logstash-log-user1
`

	FluentdOutputMixedCopy1    fluentdv1alpha1.Output
	FluentdOutputMixedCopy1Raw = `
apiVersion: fluentd.fluent.io/v1alpha1
kind: Output
metadata:
  name: fluentd-mixed-copy-es-1
  namespace: fluent
  labels:
    output.fluentd.fluent.io/enabled: "es"
spec: 
  outputs:
  - copy:
      copyMode: no_copy
  - elasticsearch:
      host: elasticsearch-logging-data.kubesphere-logging-system.svc
      port: 9243
      indexName: fluentd-mixed-copy-es-1
      scheme: https
  - elasticsearch:
      host: elasticsearch-logging-data.kubesphere-logging-system.svc
      port: 9243
      indexName: fluentd-mixed-copy-es-2
      scheme: https
`

	FluentdOutputMixedCopy2    fluentdv1alpha1.Output
	FluentdOutputMixedCopy2Raw = `
apiVersion: fluentd.fluent.io/v1alpha1
kind: Output
metadata:
  name: fluentd-mixed-copy-es-2
  namespace: fluent
  labels:
    output.fluentd.fluent.io/enabled: "es"
spec: 
  outputs:
  - copy:
      copyMode: no_copy
    tag: mixed2
  - elasticsearch:
      host: elasticsearch-logging-data.kubesphere-logging-system.svc
      port: 9243
      indexName: fluentd-mixed-copy-es-3
      scheme: https
    tag: mixed2
  - elasticsearch:
      host: elasticsearch-logging-data.kubesphere-logging-system.svc
      port: 9243
      indexName: fluentd-mixed-copy-es-4
      scheme: https
    tag: mixed2
`

	FluentdOutputMixedCopy3    fluentdv1alpha1.Output
	FluentdOutputMixedCopy3Raw = `
apiVersion: fluentd.fluent.io/v1alpha1
kind: Output
metadata:
  name: fluentd-mixed-copy-es-3
  namespace: fluent
  labels:
    output.fluentd.fluent.io/enabled: "es"
spec:
  outputs:
  - elasticsearch:
      host: elasticsearch-logging-data.kubesphere-logging-system.svc
      port: 9243
      indexName: fluentd-mixed-copy-es-5
      scheme: https
    tag: mixed3
  - elasticsearch:
      host: elasticsearch-logging-data.kubesphere-logging-system.svc
      port: 9243
      indexName: fluentd-mixed-copy-es-6
      scheme: https
    tag: mixed3
`

	FluentdClusterOutput2CloudWatch    fluentdv1alpha1.ClusterOutput
	FluentdClusterOutput2CloudWatchRaw = `
apiVersion: fluentd.fluent.io/v1alpha1
kind: ClusterOutput
metadata:
  name: fluentd-output-cloudwatch
  labels:
    output.fluentd.fluent.io/enabled: "true"
spec:
  outputs:
  - cloudWatch:
      logStreamName: loggy-mclogface
      roleArn: abc123
      awsStsRoleArn: xyz789
      webIdentityTokenFile: /var/run/secrets/something/token
`

	FluentdClusterOutput2Datadog    fluentdv1alpha1.ClusterOutput
	FluentdClusterOutput2DatadogRaw = `
apiVersion: fluentd.fluent.io/v1alpha1
kind: ClusterOutput
metadata:
  name: fluentd-output-datadog
  labels:
    output.fluentd.fluent.io/enabled: "true"
spec: 
  outputs: 
  - datadog:
      host: http-intake.logs.datadoghq.com
      port: 443
      ddSource: kubernetes
      ddSourcecategory: kubernetes
`
	FluentdClusterOutputCopy2StdoutAndLoki    fluentdv1alpha1.ClusterOutput
	FluentdClusterOutputCopy2StdoutAndLokiRaw = `
apiVersion: fluentd.fluent.io/v1alpha1
kind: ClusterOutput
metadata:
  name: fluentd-output-copy-stdout-and-loki
  labels:
    output.fluentd.fluent.io/enabled: "true"
spec:
  outputs: 
  - copy:
      copyMode: no_copy
  - stdout: {}
  - loki:
      url: http://loki-logging-data.kubesphere-logging-system.svc:3100
      extractKubernetesLabels: true
      labels:
        - key11=value11
        - key12 = value12
        - key13
      labelKeys:
        - key21
        - key22
      removeKeys:
        - key31
        - key32
      dropSingleKey: true
      includeThreadLabel: true
      insecure: true
`

	esCredentials    corev1.Secret
	esCredentialsRaw = `
apiVersion: v1
kind: Secret
metadata:
  name: es-credentials
  namespace: fluent
type: Opaque
stringData:
  password: s3cr3tP@ssword
  username: s3cr3tUsern4me
`

	once sync.Once
)

var (
	forwardId       = "forward-001"
	forwardLogLevel = "info"
	forwardLabel    = "forward-test"

	transportTls          = "tls"
	forwardCertPath       = "/etc/td-agent/certs/fluentd.crt"
	forwardPrivateKeyPath = "/etc/td-agent/certs/fluentd.key"

	forwardPort int32 = 24224

	recordKey1   = "avg"
	recordValue1 = `${record["total"] / record["count"]}`
	recordKey2   = "message"
	recordValue2 = `yay, ${record["message"]}`

	regexpKey1      = "message"
	regexpPattern1  = "/cool/"
	regexpKey2      = "hostname"
	regexpPattern2  = `/^web\d+\.example\.com$/`
	excludeKey1     = "message"
	excludePattern1 = "/exclude/"
	excludeKey2     = "hostname"
	excludePattern2 = `/^web\d+\.error\.com$/`

	regexpParser     = "regexp"
	regexpExpression = `/^(?<host>[^ ]*) [^ ]* (?<user>[^ ]*) [(?<time>[^\]]*)\] "(?<method>\S+)(?: +(?<path>[^ ]*) +\S*)?" (?<code>[^ ]*) (?<size>[^ ]*)$/`

	timeFormat = `%d/%b/%Y:%H:%M:%S %z`
	formatType = "json"

	records = []*filter.Record{
		{
			Key:   &recordKey1,
			Value: &recordValue1,
		},
		{
			Key:   &recordKey2,
			Value: &recordValue2,
		},
	}

	bufferId             = "common_buffer"
	bufferType           = "file"
	bufferPath           = "/buffers/fd.log"
	bufferTag            = "buffertag.*"
	bufferTotalLimitSize = "5GB"

	buffer = common.Buffer{
		BufferCommon: common.BufferCommon{
			Id:   &bufferId,
			Type: &bufferType,
		},
		Path:           &bufferPath,
		Tag:            bufferTag,
		TotalLimitSize: &bufferTotalLimitSize,
	}

	serverHost = "host"
	serverPort = "14423"
	sdType     = "file"
	sdPath     = "/sd/path"

	servers = []*common.Server{
		{
			Host: &serverHost,
			Port: &serverPort,
		},
	}

	serversDiscovery = common.ServiceDiscovery{
		SDCommon: common.SDCommon{Type: &sdType},
		FileServiceDiscovery: &common.FileServiceDiscovery{
			Path: &sdPath,
		},
	}

	authMethod = "basic"

	auth = common.Auth{
		Method: &authMethod,
	}

	endpoint           = "http://logserver.com:9000/api"
	opentimeout uint16 = 2

	brokers = "broker1,broker2"
)

var (
	GlobalInputs = []input.Input{
		{
			InputCommon: input.InputCommon{
				Id:       &forwardId,
				LogLevel: &forwardLogLevel,
				Label:    &forwardLabel,
			},
			Forward: &input.Forward{
				Transport: &common.Transport{
					Protocol:       &transportTls,
					CertPath:       &forwardCertPath,
					PrivateKeyPath: &forwardPrivateKeyPath,
				},
				Port: &forwardPort,
			},
		},
	}
)

func init() {
	once.Do(
		func() {
			ParseIntoObject(FluentdRaw, &Fluentd)
			ParseIntoObject(FluentdInputTailRaw, &FluentdInputTail)
			ParseIntoObject(FluentdInputSampleRaw, &FluentdInputSample)
			ParseIntoObject(FluentdInputMonitorAgentRaw, &FluentdInputMonitorAgent)
			ParseIntoObject(FluentdClusterOutputTagRaw, &FluentdClusterOutputTag)
			ParseIntoObject(FluentdClusterFluentdConfig1Raw, &FluentdClusterFluentdConfig1)
			ParseIntoObject(FluentdClusterFluentdConfig2Raw, &FluentdClusterFluentdConfig2)
			ParseIntoObject(FluentdConfigUser1Raw, &FluentdConfigUser1)
			ParseIntoObject(FluentdConfig1Raw, &FluentdConfig1)
			ParseIntoObject(FluentdConfig2Raw, &FluentdConfig2)
			ParseIntoObject(FluentdClusterFilter1Raw, &FluentdClusterFilter1)
			ParseIntoObject(FluentdFilterRaw, &FluentdFilter)
			ParseIntoObject(FluentdClusterRecordTransformerRaw, &FluentdClusterRecordTransformerFilter)
			ParseIntoObject(FluentdClusterOutputClusterRaw, &FluentdClusterOutputCluster)
			ParseIntoObject(FluentdClusterOutputLogOperatorRaw, &FluentdClusterOutputLogOperator)
			ParseIntoObject(FluentdClusterOutputBufferRaw, &FluentdClusterOutputBuffer)
			ParseIntoObject(FluentdClusterOutputMemoryBufferRaw, &FluentdClusterOutputMemoryBuffer)
			ParseIntoObject(FluentdclusterOutput2ESRaw, &FluentdclusterOutput2ES)
			ParseIntoObject(FluentdOutput2ES1Raw, &FluentdOutput2ES1)
			ParseIntoObject(FluentdOutput2ES2Raw, &FluentdOutput2ES2)
			ParseIntoObject(FluentdOutput2ES3Raw, &FluentdOutput2ES3)
			ParseIntoObject(FluentdOutput2ES4Raw, &FluentdOutput2ES4)
			ParseIntoObject(FluentdclusterOutput2OpenSearchRaw, &FluentdclusterOutput2OpenSearch)
			ParseIntoObject(FluentdClusterOutput2kafkaRaw, &FluentdClusterOutput2kafka)
			ParseIntoObject(FluentdClusterOutput2LokiRaw, &FluentdClusterOutput2Loki)
			ParseIntoObject(FluentdClusterOutput2Loki1Raw, &FluentdClusterOutput2Loki1)
			ParseIntoObject(FluentdOutputUser1Raw, &FluentdOutputUser1)
			ParseIntoObject(FluentdClusterOutputCustomRaw, &FluentdClusterOutputCustom)
			ParseIntoObject(FluentdClusterOutput2CloudWatchRaw, &FluentdClusterOutput2CloudWatch)
			ParseIntoObject(FluentdClusterOutput2DatadogRaw, &FluentdClusterOutput2Datadog)
			ParseIntoObject(FluentdClusterOutputCopy2StdoutAndLokiRaw, &FluentdClusterOutputCopy2StdoutAndLoki)
			ParseIntoObject(FluentdOutputMixedCopy1Raw, &FluentdOutputMixedCopy1)
			ParseIntoObject(FluentdOutputMixedCopy2Raw, &FluentdOutputMixedCopy2)
			ParseIntoObject(FluentdOutputMixedCopy3Raw, &FluentdOutputMixedCopy3)
			ParseIntoObject(esCredentialsRaw, &esCredentials)
			ParseIntoObject(lokiHttpCredentialsRaw, &lokiHttpCredentials)
			ParseIntoObject(lokiTenantNameRaw, &lokiTenantName)
		},
	)
}

func ParseIntoObject(data string, obj interface{}) error {
	body, err := yaml.YAMLToJSON([]byte(data))
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &obj)
	if err != nil {
		return err
	}

	return nil
}

func getExpectedCfg(path string) []byte {
	data, _ := os.ReadFile(path)
	return data
}

func CreateFluentdFilterSpecs() (fluentdv1alpha1.FilterSpec, fluentdv1alpha1.FilterSpec, fluentdv1alpha1.FilterSpec) {
	filterSpec1 := fluentdv1alpha1.FilterSpec{
		Filters: []filter.Filter{
			{
				RecordTransformer: &filter.RecordTransformer{
					Records: records,
				},
			},
			{
				Grep: &filter.Grep{
					Regexps: []*filter.Regexp{
						{
							Key:     &regexpKey1,
							Pattern: &regexpPattern1,
						},
						{
							Key:     &regexpKey2,
							Pattern: &regexpPattern2,
						},
					},
					Excludes: []*filter.Exclude{
						{
							Key:     &excludeKey1,
							Pattern: &excludePattern1,
						},
						{
							Key:     &excludeKey2,
							Pattern: &excludePattern2,
						},
					},
				},
			},
		},
	}

	filterSpec2 := fluentdv1alpha1.FilterSpec{
		Filters: []filter.Filter{
			{
				RecordTransformer: &filter.RecordTransformer{
					Records: records,
				},
			},
			{
				Stdout: &filter.Stdout{
					Format: &common.Format{
						FormatCommon: common.FormatCommon{
							Type: &formatType,
						},
					},
				},
			},
		},
	}

	filterSpec3 := fluentdv1alpha1.FilterSpec{
		Filters: []filter.Filter{
			{
				Parser: &filter.Parser{
					Parse: &common.Parse{
						ParseCommon: common.ParseCommon{
							Type: &regexpParser,
						},
						Time: common.Time{
							TimeFormat: &timeFormat,
						},
						Expression: &regexpExpression,
					},
				},
			},
			{
				Stdout: &filter.Stdout{
					Format: &common.Format{
						FormatCommon: common.FormatCommon{
							Type: &formatType,
						},
					},
				},
			},
		},
	}

	return filterSpec1, filterSpec2, filterSpec3
}

func CreateFluentdOutputSpecs() (fluentdv1alpha1.OutputSpec, fluentdv1alpha1.OutputSpec, fluentdv1alpha1.OutputSpec) {
	outputSpec1 := fluentdv1alpha1.OutputSpec{
		Outputs: []output.Output{
			{
				BufferSection: common.BufferSection{
					Buffer: &buffer,
				},
				Forward: &output.Forward{
					Servers:          servers,
					ServiceDiscovery: &serversDiscovery,
				},
			},
			{
				BufferSection: common.BufferSection{
					Buffer: &buffer,
				},
				Kafka: &output.Kafka2{
					Brokers: &brokers,
				},
			},
		},
	}

	outputSpec2 := fluentdv1alpha1.OutputSpec{
		Outputs: []output.Output{
			{
				BufferSection: common.BufferSection{
					Buffer: &buffer,
				},
				Forward: &output.Forward{
					Servers:          servers,
					ServiceDiscovery: &serversDiscovery,
				},
			},
			{
				BufferSection: common.BufferSection{
					Buffer: &buffer,
				},
				Http: &output.Http{
					Auth:        &auth,
					Endpoint:    &endpoint,
					OpenTimeout: &opentimeout,
				},
			},
		},
	}

	outputSpec3 := fluentdv1alpha1.OutputSpec{
		Outputs: []output.Output{
			{
				BufferSection: common.BufferSection{
					Buffer: &buffer,
				},
				Http: &output.Http{
					Auth:        &auth,
					Endpoint:    &endpoint,
					OpenTimeout: &opentimeout,
				},
			},
			{
				BufferSection: common.BufferSection{
					Buffer: &buffer,
				},
				Kafka: &output.Kafka2{
					Brokers: &brokers,
				},
			},
		},
	}

	return outputSpec1, outputSpec2, outputSpec3
}
