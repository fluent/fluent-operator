package v1alpha2

import (
	"fmt"
	"testing"

	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins/multilineparser"
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins/parser"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins/custom"
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins/filter"
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins/input"
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins/output"
	. "github.com/onsi/gomega"
)

const authorization = "foo:bar"
const xLogHeader0 = "testing"
const xLogHeaderAppID = "9780495d9db3"
const xLogHeaderAppName = "app_name"

var expected = fmt.Sprintf(`[Service]
    Daemon    false
    Flush    1
    Grace    30
    Http_Server    true
    Log_Level    info
    Parsers_File    /fluent-bit/etc/parsers.conf
[Input]
    Name    tail
    Alias    input0_alias
    Path    /logs/containers/apps0
    Exclude_Path    /logs/containers/exclude_path
    Refresh_Interval    10
    Ignore_Older    5m
    Skip_Long_Lines    true
    DB    /fluent-bit/tail/pos.db
    Mem_Buf_Limit    5MB
    Tag    logs.foo.bar
    Inotify_Watcher    false
[Filter]
    Name    modify
    Match    logs.foo.bar
    Condition    Key_value_equals    kve0    kvev0
    Condition    Key_value_equals    kve1    kvev1
    Condition    Key_value_equals    kve2    kvev2
    Condition    Key_does_not_exist    kdn0    kdnv0
    Condition    Key_does_not_exist    kdn1    kdnv1
    Condition    Key_does_not_exist    kdn2    kdnv2
    Set    app    foo
    Set    customer    cus1
    Set    sk0    skv0
    Add    add_k0    k0value
    Add    add_k1    k1v
    Add    add_k2    k2val
    Rename    rk0    r0v
    Rename    rk1    r1v
    Rename    rk2    r2v
    Rename    rk3    r3v
[Output]
    Name    es
    Match    *
    Alias    output_elasticsearch_alias
    Host    https://example2.com
    Port    9200
    Index    my_index
    Type    my_type
    Write_Operation    upsert
[Output]
    Name    http
    Match    logs.foo.bar
    Alias    output_http_alias
    host    https://example2.com
    port    433
    uri    /logs
    format    json_lines
    header     Authorization    %s
    header     X-Log-Header-0    %s
    header     X-Log-Header-App-ID    %s
    header     X-Log-Header-App-Name    %s
    json_date_key    timestamp
    json_date_format    iso8601
    tls    On
    tls.verify    true
[Output]
    Name    kafka
    Topics    fluentbit
    Match    *
    Brokers    192.168.100.32:9092
    rdkafka.debug    All
    rdkafka.request.required.acks    1
    rdkafka.log.connection.close    false
    rdkafka.log_level    7
    rdkafka.metadata.broker.list    192.168.100.32:9092
[Output]
    Name    opensearch
    Match    *
    Alias    output_opensearch_alias
    Host    https://example2.com
    Port    9200
    Index    my_index
    Type    my_type
[Output]
    Name    syslog
    Match    logs.foo.bar
    Alias    output_syslog_alias
    Host    example.com
    port    3300
    mode    tls
    syslog_hostname_key    do_app_name
    syslog_appname_key    do_component_name
    syslog_message_key    log
    tls    On
    tls.verify    true
`, authorization, xLogHeader0, xLogHeaderAppID, xLogHeaderAppName)

var expectedYaml = fmt.Sprintf(`service:
  daemon: false
  flush: 1
  grace: 30
  http_server: true
  log_level: info
  parsers_file: /fluent-bit/etc/parsers.conf
pipeline:
  inputs:
    - name: tail
      alias: input0_alias
      path: /logs/containers/apps0
      exclude_path: /logs/containers/exclude_path
      refresh_interval: 10
      ignore_older: 5m
      skip_long_lines: true
      db: /fluent-bit/tail/pos.db
      mem_buf_limit: 5MB
      tag: logs.foo.bar
      inotify_watcher: false
  filters:
    - name: modify
      match: "logs.foo.bar"
      condition:
        - Key_value_equals    kve0    kvev0
        - Key_value_equals    kve1    kvev1
        - Key_value_equals    kve2    kvev2
        - Key_does_not_exist    kdn0    kdnv0
        - Key_does_not_exist    kdn1    kdnv1
        - Key_does_not_exist    kdn2    kdnv2
      set:
        - app    foo
        - customer    cus1
        - sk0    skv0
      add:
        - add_k0    k0value
        - add_k1    k1v
        - add_k2    k2val
      rename:
        - rk0    r0v
        - rk1    r1v
        - rk2    r2v
        - rk3    r3v
  outputs:
    - name: es
      match: "*"
      alias: output_elasticsearch_alias
      host: https://example2.com
      port: 9200
      index: my_index
      type: my_type
      write_operation: upsert
    - name: http
      match: "logs.foo.bar"
      alias: output_http_alias
      host: https://example2.com
      port: 433
      uri: /logs
      format: json_lines
      header:
        -  Authorization    %s
        -  X-Log-Header-0    %s
        -  X-Log-Header-App-ID    %s
        -  X-Log-Header-App-Name    %s
      json_date_key: timestamp
      json_date_format: iso8601
      tls: On
      tls.verify: true
    - name: kafka
      match: kube.*
      brokers: 192.168.100.32:9092
      topics: fluentbit
    - name: opensearch
      match: "*"
      alias: output_opensearch_alias
      host: https://example2.com
      port: 9200
      index: my_index
      type: my_type
    - name: syslog
      match: "logs.foo.bar"
      alias: output_syslog_alias
      host: example.com
      port: 3300
      mode: tls
      syslog_hostname_key: do_app_name
      syslog_appname_key: do_component_name
      syslog_message_key: log
      tls: On
      tls.verify: true
    - name: kafka-namespace
      match: 098f6bcd4621d373cade4e832627b4f6.kube.namespace.*
      brokers: 127.0.1.1:9092
      topics: fluentbit-namespace
`, authorization, xLogHeader0, xLogHeaderAppID, xLogHeaderAppName)

var expectedK8s = `[Service]
    Daemon    false
    Flush    1
    Grace    30
    Http_Server    true
    Log_Level    info
    Parsers_File    /fluent-bit/etc/parsers.conf
[Input]
    Name    tail
    Path    /var/log/containers/*.log
    Refresh_Interval    10
    Ignore_Older    5m
    Skip_Long_Lines    true
    DB    /fluent-bit/tail/pos.db
    Mem_Buf_Limit    5MB
    Tag    kube.*
[Filter]
    Name    parser
    Match    acbd18db4cc2f85cedef654fccc4a4d8.kube.*
    Key_Name    log
    Parser    bar-acbd18db4cc2f85cedef654fccc4a4d8
    Reserve_Data    true
[Output]
    Name    opensearch
    Match    acbd18db4cc2f85cedef654fccc4a4d8.kube.*
    Host    foo.bar
    Port    9200
    Index    foo-index
[Output]
    Name    es
    Match    acbd18db4cc2f85cedef654fccc4a4d8.kube.*
    Host    foo.bar
    Port    9200
    Index    foo-index
    Write_Operation    update
`

var expectedK8sYaml = `service:
  daemon: false
  flush: 1
  grace: 30
  http_server: true
  log_level: info
  parsers_file: /fluent-bit/etc/parsers.conf
pipeline:
  inputs:
    - name: tail
      path: /var/log/containers/*.log
      refresh_interval: 10
      ignore_older: 5m
      skip_long_lines: true
      db: /fluent-bit/tail/pos.db
      mem_buf_limit: 5MB
      tag: kube.*
  filters:
    - name: parser
      match: "acbd18db4cc2f85cedef654fccc4a4d8.kube.*"
      key_name: log
      parser: bar-acbd18db4cc2f85cedef654fccc4a4d8-acbd18db4cc2f85cedef654fccc4a4d8
      reserve_data: true
  outputs:
    - name: opensearch
      match: "acbd18db4cc2f85cedef654fccc4a4d8.kube.*"
      host: foo.bar
      port: 9200
      index: foo-index
    - name: es
      match: "acbd18db4cc2f85cedef654fccc4a4d8.kube.*"
      host: foo.bar
      port: 9200
      index: foo-index
      write_operation: update
`

var expectedK8sYamlWithClusterFilterOutput = `service:
  daemon: false
  flush: 1
  grace: 30
  http_server: true
  log_level: info
  parsers_file: /fluent-bit/etc/parsers.conf
pipeline:
  inputs:
    - name: tail
      path: /var/log/containers/*.log
      refresh_interval: 10
      ignore_older: 5m
      skip_long_lines: true
      db: /fluent-bit/tail/pos.db
      mem_buf_limit: 5MB
      tag: kube.*
  filters:
    - name: parser
      match: "kubernetes.*"
      key_name: log
      parser: test
      reserve_data: true
    - name: parser
      match: "acbd18db4cc2f85cedef654fccc4a4d8.kube.*"
      key_name: log
      parser: bar-acbd18db4cc2f85cedef654fccc4a4d8
      reserve_data: true
  outputs:
    - name: es
      match: "kube.*"
      host: foo.bar
      port: 9200
      index: foo-index
    - name: opensearch
      match: "acbd18db4cc2f85cedef654fccc4a4d8.kube.*"
      host: foo.bar
      port: 9200
      index: foo-index
    - name: es
      match: "acbd18db4cc2f85cedef654fccc4a4d8.kube.*"
      host: foo.bar
      port: 9200
      index: foo-index
      write_operation: update
`

var expectedParsers = `[PARSER]
    Name    clusterparser0
    Format    json
    Time_Key    time
    Time_Format    %Y-%m-%dT%H:%M:%S %z
[PARSER]
    Name    parser0-4087ca5ebba883e13a4369122e716be7
    Format    regex
    Regex    .*
    Time_Key    time
    Time_Format    %Y-%m-%dT%H:%M:%S %z
[PARSER]
    Name    clusterparser1
    Format    ltsv
    Time_Key    time
    Time_Format    [%d/%b/%Y:%H:%M:%S %z]
    Types    status:integer size:integer
`

var expectedMultilineParsers = `[MULTILINE_PARSER]
    Name    clustermultilineparser0
    Type    regex
    Parser    go
    Key_Content    log
[MULTILINE_PARSER]
    Name    multilineparser0
    Type    regex
    Flush_Timeout    1000
    Rule    "start_state" "/(Dec \d+ \d+\:\d+\:\d+)(.*)/" "cont"
    Rule    "cont" "/^\s+at.*/" "cont"
[MULTILINE_PARSER]
    Name    clustermultilineparser1
    Type    regex
    Flush_Timeout    500
    Rule    "start_state" "/^(\d+ \d+\:\d+\:\d+)(.*)/" "cont"
    Rule    "cont" "/^\s+at.*/" "cont"
`

var labels = map[string]string{
	"label0": "lv0",
	"label1": "lv1",
	"label3": "lval3",
	"lbl2":   "lval2",
	"lbl1":   "lvl1",
}

var cfg = ClusterFluentBitConfig{
	Spec: FluentBitConfigSpec{
		Service: &Service{
			Daemon:       ptrBool(false),
			FlushSeconds: ptrFloat64(1),
			GraceSeconds: ptrInt64(30),
			HttpServer:   ptrBool(true),
			LogLevel:     "info",
			ParsersFile:  "parsers.conf",
		},
	},
}

func Test_FluentBitConfig_RenderMainConfig(t *testing.T) {
	g := NewGomegaWithT(t)

	sl := plugins.NewSecretLoader(nil, "testnamespace")

	disableInotifyWatcher := ptrBool(true)

	inputObj := &ClusterInput{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "fluentbit.fluent.io/v1alpha2",
			Kind:       "ClusterInput",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   "input0",
			Labels: labels,
		},
		Spec: InputSpec{
			Alias: "input0_alias",
			Tail: &input.Tail{
				DisableInotifyWatcher:  disableInotifyWatcher,
				Tag:                    "logs.foo.bar",
				Path:                   "/logs/containers/apps0",
				ExcludePath:            "/logs/containers/exclude_path",
				SkipLongLines:          ptrBool(true),
				IgnoreOlder:            "5m",
				MemBufLimit:            "5MB",
				RefreshIntervalSeconds: ptrInt64(10),
				DB:                     "/fluent-bit/tail/pos.db",
			},
		},
	}

	inputs := ClusterInputList{
		Items: []ClusterInput{*inputObj},
	}

	filterObj := &ClusterFilter{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "fluentbit.fluent.io/v1alpha2",
			Kind:       "ClusterFilter",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   "filter0",
			Labels: labels,
		},
		Spec: FilterSpec{
			Match: "logs.foo.bar",
			FilterItems: []FilterItem{
				{
					Modify: &filter.Modify{
						Conditions: []filter.Condition{
							{
								KeyValueEquals: map[string]string{
									"kve1": "kvev1",
									"kve0": "kvev0",
									"kve2": "kvev2",
								},
							},
							{
								KeyDoesNotExist: map[string]string{
									"kdn1": "kdnv1",
									"kdn0": "kdnv0",
									"kdn2": "kdnv2",
								},
							},
						},
						Rules: []filter.Rule{
							{
								Set: map[string]string{
									"sk0":      "skv0",
									"customer": "cus1",
									"app":      "foo",
								},
								Add: map[string]string{
									"add_k1": "k1v",
									"add_k2": "k2val",
									"add_k0": "k0value",
								},
								Rename: map[string]string{
									"rk1": "r1v",
									"rk0": "r0v",
									"rk3": "r3v",
									"rk2": "r2v",
								},
							},
						},
					},
				},
			},
		},
	}

	filters := ClusterFilterList{
		Items: []ClusterFilter{*filterObj},
	}

	syslogOut := ClusterOutput{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "fluentbit.fluent.io/v1alpha2",
			Kind:       "ClusterOutput",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   "syslog_output0",
			Labels: labels,
		},
		Spec: OutputSpec{
			Alias: "output_syslog_alias",
			Match: "logs.foo.bar",
			Syslog: &output.Syslog{
				Host: "example.com",
				Port: ptrInt32(int32(3300)),
				Mode: "tls",
				TLS: &plugins.TLS{
					Verify: ptrBool(true),
				},
				SyslogMessageKey:  "log",
				SyslogHostnameKey: "do_app_name",
				SyslogAppnameKey:  "do_component_name",
			},
		},
	}

	headers := map[string]string{}

	headers["Authorization"] = authorization
	headers["X-Log-Header-App-Name"] = xLogHeaderAppName
	headers["X-Log-Header-0"] = xLogHeader0
	headers["X-Log-Header-App-ID"] = xLogHeaderAppID

	httpOutput := ClusterOutput{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "fluentbit.fluent.io/v1alpha2",
			Kind:       "ClusterOutput",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   "http_output_0",
			Labels: labels,
		},
		Spec: OutputSpec{
			Alias: "output_http_alias",
			Match: "logs.foo.bar",
			HTTP: &output.HTTP{
				Host:           "https://example2.com",
				Port:           ptrInt32(int32(433)),
				Uri:            "/logs",
				Headers:        headers,
				Format:         "json_lines",
				JsonDateKey:    "timestamp",
				JsonDateFormat: "iso8601",
				TLS: &plugins.TLS{
					Verify: ptrBool(true),
				},
			},
		},
	}

	openSearchOutput := ClusterOutput{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "fluentbit.fluent.io/v1alpha2",
			Kind:       "ClusterOutput",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   "opensearch_output_0",
			Labels: labels,
		},
		Spec: OutputSpec{
			Alias: "output_opensearch_alias",
			Match: "*",
			OpenSearch: &output.OpenSearch{
				Host:  "https://example2.com",
				Port:  ptrInt32(int32(9200)),
				Index: "my_index",
				Type:  "my_type",
			},
		},
	}

	elasticSearchOutput := ClusterOutput{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "fluentbit.fluent.io/v1alpha2",
			Kind:       "ClusterOutput",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   "elasticsearch_output_0",
			Labels: labels,
		},
		Spec: OutputSpec{
			Alias: "output_elasticsearch_alias",
			Match: "*",
			Elasticsearch: &output.Elasticsearch{
				Host:           "https://example2.com",
				Port:           ptrInt32(int32(9200)),
				Index:          "my_index",
				Type:           "my_type",
				WriteOperation: "upsert",
			},
		},
	}

	kafkaOutput := ClusterOutput{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "fluentbit.fluent.io/v1alpha2",
			Kind:       "ClusterOutput",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   "kafka_output",
			Labels: labels,
		},
		Spec: OutputSpec{
			CustomPlugin: &custom.CustomPlugin{
				Config: "    Name    kafka\n    Topics    fluentbit\n    Match    *\n    Brokers    192.168.100.32:9092\n    rdkafka.debug    All\n    rdkafka.request.required.acks    1\n    rdkafka.log.connection.close    false\n    rdkafka.log_level    7\n    rdkafka.metadata.broker.list    192.168.100.32:9092",
			},
		},
	}

	outputs := ClusterOutputList{
		Items: []ClusterOutput{syslogOut, httpOutput, openSearchOutput, elasticSearchOutput, kafkaOutput},
	}

	var nsFilterList []FilterList
	var nsOutputList []OutputList
	var rewriteTagCfgs []string
	// we should not see any permutations in serialized config
	i := 0
	for i < 5 {
		config, err := cfg.RenderMainConfig(sl, inputs, filters, outputs, nsFilterList, nsOutputList, rewriteTagCfgs)
		g.Expect(err).NotTo(HaveOccurred())
		g.Expect(config).To(Equal(expected))

		i++
	}
}
func Test_FluentBitConfig_RenderMainConfigYaml(t *testing.T) {
	g := NewGomegaWithT(t)

	sl := plugins.NewSecretLoader(nil, "testnamespace")

	disableInotifyWatcher := ptrBool(true)

	inputObj := &ClusterInput{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "fluentbit.fluent.io/v1alpha2",
			Kind:       "ClusterInput",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   "input0",
			Labels: labels,
		},
		Spec: InputSpec{
			Alias: "input0_alias",
			Tail: &input.Tail{
				DisableInotifyWatcher:  disableInotifyWatcher,
				Tag:                    "logs.foo.bar",
				Path:                   "/logs/containers/apps0",
				ExcludePath:            "/logs/containers/exclude_path",
				SkipLongLines:          ptrBool(true),
				IgnoreOlder:            "5m",
				MemBufLimit:            "5MB",
				RefreshIntervalSeconds: ptrInt64(10),
				DB:                     "/fluent-bit/tail/pos.db",
			},
		},
	}

	inputs := ClusterInputList{
		Items: []ClusterInput{*inputObj},
	}

	filterObj := &ClusterFilter{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "fluentbit.fluent.io/v1alpha2",
			Kind:       "ClusterFilter",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   "filter0",
			Labels: labels,
		},
		Spec: FilterSpec{
			Match: "logs.foo.bar",
			FilterItems: []FilterItem{
				{
					Modify: &filter.Modify{
						Conditions: []filter.Condition{
							{
								KeyValueEquals: map[string]string{
									"kve1": "kvev1",
									"kve0": "kvev0",
									"kve2": "kvev2",
								},
							},
							{
								KeyDoesNotExist: map[string]string{
									"kdn1": "kdnv1",
									"kdn0": "kdnv0",
									"kdn2": "kdnv2",
								},
							},
						},
						Rules: []filter.Rule{
							{
								Set: map[string]string{
									"sk0":      "skv0",
									"customer": "cus1",
									"app":      "foo",
								},
								Add: map[string]string{
									"add_k1": "k1v",
									"add_k2": "k2val",
									"add_k0": "k0value",
								},
								Rename: map[string]string{
									"rk1": "r1v",
									"rk0": "r0v",
									"rk3": "r3v",
									"rk2": "r2v",
								},
							},
						},
					},
				},
			},
		},
	}

	filters := ClusterFilterList{
		Items: []ClusterFilter{*filterObj},
	}

	syslogOut := ClusterOutput{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "fluentbit.fluent.io/v1alpha2",
			Kind:       "ClusterOutput",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   "syslog_output0",
			Labels: labels,
		},
		Spec: OutputSpec{
			Alias: "output_syslog_alias",
			Match: "logs.foo.bar",
			Syslog: &output.Syslog{
				Host: "example.com",
				Port: ptrInt32(int32(3300)),
				Mode: "tls",
				TLS: &plugins.TLS{
					Verify: ptrBool(true),
				},
				SyslogMessageKey:  "log",
				SyslogHostnameKey: "do_app_name",
				SyslogAppnameKey:  "do_component_name",
			},
		},
	}

	headers := map[string]string{}

	headers["Authorization"] = authorization
	headers["X-Log-Header-App-Name"] = xLogHeaderAppName
	headers["X-Log-Header-0"] = xLogHeader0
	headers["X-Log-Header-App-ID"] = xLogHeaderAppID

	httpOutput := ClusterOutput{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "fluentbit.fluent.io/v1alpha2",
			Kind:       "ClusterOutput",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   "http_output_0",
			Labels: labels,
		},
		Spec: OutputSpec{
			Alias: "output_http_alias",
			Match: "logs.foo.bar",
			HTTP: &output.HTTP{
				Host:           "https://example2.com",
				Port:           ptrInt32(int32(433)),
				Uri:            "/logs",
				Headers:        headers,
				Format:         "json_lines",
				JsonDateKey:    "timestamp",
				JsonDateFormat: "iso8601",
				TLS: &plugins.TLS{
					Verify: ptrBool(true),
				},
			},
		},
	}

	openSearchOutput := ClusterOutput{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "fluentbit.fluent.io/v1alpha2",
			Kind:       "ClusterOutput",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   "opensearch_output_0",
			Labels: labels,
		},
		Spec: OutputSpec{
			Alias: "output_opensearch_alias",
			Match: "*",
			OpenSearch: &output.OpenSearch{
				Host:  "https://example2.com",
				Port:  ptrInt32(int32(9200)),
				Index: "my_index",
				Type:  "my_type",
			},
		},
	}

	elasticSearchOutput := ClusterOutput{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "fluentbit.fluent.io/v1alpha2",
			Kind:       "ClusterOutput",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   "elasticsearch_output_0",
			Labels: labels,
		},
		Spec: OutputSpec{
			Alias: "output_elasticsearch_alias",
			Match: "*",
			Elasticsearch: &output.Elasticsearch{
				Host:           "https://example2.com",
				Port:           ptrInt32(int32(9200)),
				Index:          "my_index",
				Type:           "my_type",
				WriteOperation: "upsert",
			},
		},
	}

	kafkaOutput := ClusterOutput{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "fluentbit.fluent.io/v1alpha2",
			Kind:       "ClusterOutput",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   "kafka_output",
			Labels: labels,
		},
		Spec: OutputSpec{
			CustomPlugin: &custom.CustomPlugin{
				YamlConfig: &plugins.Config{Data: map[string]interface{}{
					"name":    "kafka",
					"topics":  "fluentbit",
					"match":   "kube.*",
					"brokers": "192.168.100.32:9092",
				}},
			},
		},
	}

	outputs := ClusterOutputList{
		Items: []ClusterOutput{syslogOut, httpOutput, openSearchOutput, elasticSearchOutput, kafkaOutput},
	}

	var nsFilterList []FilterList
	nsOutputList := []OutputList{
		{
			Items: []Output{
				{
					TypeMeta: metav1.TypeMeta{
						APIVersion: "fluentbit.fluent.io/v1alpha2",
						Kind:       "ClusterOutput",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:      "kafka_output_2",
						Namespace: "test",
						Labels:    labels,
					},
					Spec: OutputSpec{
						CustomPlugin: &custom.CustomPlugin{
							YamlConfig: &plugins.Config{Data: map[string]interface{}{
								"name":    "kafka-namespace",
								"topics":  "fluentbit-namespace",
								"match":   "kube.namespace.*",
								"brokers": "127.0.1.1:9092",
							}},
						},
					},
				}},
		},
	}
	var rewriteTagCfgs []string
	// we should not see any permutations in serialized config
	i := 0
	for i < 5 {
		config, err := cfg.RenderMainConfigInYaml(sl, inputs, filters, outputs, nsFilterList, nsOutputList, rewriteTagCfgs)
		g.Expect(err).NotTo(HaveOccurred())
		g.Expect(config).To(Equal(expectedYaml))

		i++
	}
}
func TestRenderMainConfigK8s(t *testing.T) {
	g := NewGomegaWithT(t)
	sl := plugins.NewSecretLoader(nil, "testnamespace")

	inputObj := &ClusterInput{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "fluentbit.fluent.io/v1alpha2",
			Kind:       "ClusterInput",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   "input0",
			Labels: labels,
		},
		Spec: InputSpec{
			Tail: &input.Tail{
				Tag:                    "kube.*",
				Path:                   "/var/log/containers/*.log",
				SkipLongLines:          ptrBool(true),
				IgnoreOlder:            "5m",
				MemBufLimit:            "5MB",
				RefreshIntervalSeconds: ptrInt64(10),
				DB:                     "/fluent-bit/tail/pos.db",
			},
		},
	}
	inputList := ClusterInputList{
		Items: []ClusterInput{*inputObj},
	}
	filterList := ClusterFilterList{
		Items: []ClusterFilter{},
	}
	outputList := ClusterOutputList{
		Items: []ClusterOutput{},
	}
	filterObj := &Filter{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "fluentbit.fluent.io/v1alpha2",
			Kind:       "Filter",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "filter0",
			Namespace: "foo",
			Labels:    labels,
		},
		Spec: FilterSpec{
			Match: "kube.*",
			FilterItems: []FilterItem{
				{
					Parser: &filter.Parser{
						KeyName:     "log",
						Parser:      "bar",
						ReserveData: ptrBool(true),
					},
				},
			},
		},
	}
	nsFilterList := []FilterList{
		{
			Items: []Filter{*filterObj},
		},
	}
	outputObj := &Output{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "fluentbit.fluent.io/v1alpha2",
			Kind:       "Output",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "output0",
			Namespace: "foo",
			Labels:    labels,
		},
		Spec: OutputSpec{
			Match: "kube.*",
			OpenSearch: &output.OpenSearch{
				Host:  "foo.bar",
				Port:  ptrInt32(9200),
				Index: "foo-index",
			},
		},
	}

	outputObjEs := &Output{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "fluentbit.fluent.io/v1alpha2",
			Kind:       "Output",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "output1",
			Namespace: "foo",
			Labels:    labels,
		},
		Spec: OutputSpec{
			Match: "kube.*",
			Elasticsearch: &output.Elasticsearch{
				Host:           "foo.bar",
				Port:           ptrInt32(9200),
				Index:          "foo-index",
				WriteOperation: "update",
			},
		},
	}

	nsOutputList := []OutputList{
		{
			Items: []Output{*outputObj, *outputObjEs},
		},
	}
	var rewriteTagCfg []string

	config, err := cfg.RenderMainConfig(
		sl, inputList, filterList, outputList, nsFilterList, nsOutputList, rewriteTagCfg,
	)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(config).To(Equal(expectedK8s))

	yamlConfig, err := cfg.RenderMainConfigInYaml(
		sl, inputList, filterList, outputList, nsFilterList, nsOutputList, rewriteTagCfg,
	)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(yamlConfig).To(Equal(expectedK8sYaml))
}

func TestRenderMainConfigK8sInYaml(t *testing.T) {
	g := NewGomegaWithT(t)
	sl := plugins.NewSecretLoader(nil, "testnamespace")

	inputObj := &ClusterInput{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "fluentbit.fluent.io/v1alpha2",
			Kind:       "ClusterInput",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   "input0",
			Labels: labels,
		},
		Spec: InputSpec{
			Tail: &input.Tail{
				Tag:                    "kube.*",
				Path:                   "/var/log/containers/*.log",
				SkipLongLines:          ptrBool(true),
				IgnoreOlder:            "5m",
				MemBufLimit:            "5MB",
				RefreshIntervalSeconds: ptrInt64(10),
				DB:                     "/fluent-bit/tail/pos.db",
			},
		},
	}
	inputList := ClusterInputList{
		Items: []ClusterInput{*inputObj},
	}
	filterList := ClusterFilterList{
		Items: []ClusterFilter{
			{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "fluentbit.fluent.io/v1alpha2",
					Kind:       "Filter",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:   "cluster-filter",
					Labels: labels,
				},
				Spec: FilterSpec{
					Match: "kubernetes.*",
					FilterItems: []FilterItem{
						{Parser: &filter.Parser{
							KeyName:     "log",
							Parser:      "test",
							ReserveData: ptrBool(true),
						}},
					},
				},
			},
		},
	}
	outputList := ClusterOutputList{
		Items: []ClusterOutput{
			{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "fluentbit.fluent.io/v1alpha2",
					Kind:       "Output",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:   "cluster-output0",
					Labels: labels,
				},
				Spec: OutputSpec{
					Match: "kube.*",
					Elasticsearch: &output.Elasticsearch{
						Host:  "foo.bar",
						Port:  ptrInt32(9200),
						Index: "foo-index",
					},
				},
			},
		},
	}
	filterObj := &Filter{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "fluentbit.fluent.io/v1alpha2",
			Kind:       "Filter",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "filter0",
			Namespace: "foo",
			Labels:    labels,
		},
		Spec: FilterSpec{
			Match: "kube.*",
			FilterItems: []FilterItem{
				{
					Parser: &filter.Parser{
						KeyName:     "log",
						Parser:      "bar",
						ReserveData: ptrBool(true),
					},
				},
			},
		},
	}
	nsFilterList := []FilterList{
		{
			Items: []Filter{*filterObj},
		},
	}
	outputObj := &Output{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "fluentbit.fluent.io/v1alpha2",
			Kind:       "Output",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "output0",
			Namespace: "foo",
			Labels:    labels,
		},
		Spec: OutputSpec{
			Match: "kube.*",
			OpenSearch: &output.OpenSearch{
				Host:  "foo.bar",
				Port:  ptrInt32(9200),
				Index: "foo-index",
			},
		},
	}

	outputObjEs := &Output{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "fluentbit.fluent.io/v1alpha2",
			Kind:       "Output",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "output1",
			Namespace: "foo",
			Labels:    labels,
		},
		Spec: OutputSpec{
			Match: "kube.*",
			Elasticsearch: &output.Elasticsearch{
				Host:           "foo.bar",
				Port:           ptrInt32(9200),
				Index:          "foo-index",
				WriteOperation: "update",
			},
		},
	}

	nsOutputList := []OutputList{
		{
			Items: []Output{*outputObj, *outputObjEs},
		},
	}
	var rewriteTagCfg []string

	yamlConfig, err := cfg.RenderMainConfigInYaml(
		sl, inputList, filterList, outputList, nsFilterList, nsOutputList, rewriteTagCfg,
	)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(yamlConfig).To(Equal(expectedK8sYamlWithClusterFilterOutput))
}

func TestClusterFluentBitConfig_RenderMainConfig_WithParsersFiles(t *testing.T) {
	g := NewGomegaWithT(t)
	sl := plugins.NewSecretLoader(nil, "testnamespace")

	cfbc := ClusterFluentBitConfig{
		Spec: FluentBitConfigSpec{
			Service: &Service{
				Daemon:       ptrBool(false),
				FlushSeconds: ptrFloat64(1),
				GraceSeconds: ptrInt64(30),
				HttpServer:   ptrBool(true),
				LogLevel:     "info",
				ParsersFiles: []string{"parsers.conf", "parsers_multiline.conf"},
			},
		},
	}

	config, err := cfbc.RenderMainConfig(
		sl, ClusterInputList{}, ClusterFilterList{}, ClusterOutputList{}, nil, nil, nil,
	)
	g.Expect(err).NotTo(HaveOccurred())
	fmt.Println(config)
}

func TestClusterFluentBitConfig_RenderParserConfig(t *testing.T) {
	g := NewGomegaWithT(t)
	sl := plugins.NewSecretLoader(nil, "testnamespace")

	clusterParser := &ClusterParser{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "fluentbit.fluent.io/v1alpha2",
			Kind:       "ClusterParser",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   "clusterparser0",
			Labels: labels,
		},
		Spec: ParserSpec{
			JSON: &parser.JSON{
				TimeKey:    "time",
				TimeFormat: "%Y-%m-%dT%H:%M:%S %z",
			},
		},
	}
	clusterParsers := ClusterParserList{
		Items: []ClusterParser{*clusterParser},
	}

	nsParser := &Parser{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "fluentbit.fluent.io/v1alpha2",
			Kind:       "Parser",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "parser0",
			Namespace: "testnamespace",
			Labels:    labels,
		},
		Spec: ParserSpec{
			Regex: &parser.Regex{
				Regex:      ".*",
				TimeKey:    "time",
				TimeFormat: "%Y-%m-%dT%H:%M:%S %z",
			},
		},
	}
	nsParsers := []ParserList{
		{
			Items: []Parser{*nsParser},
		},
	}

	nsClusterParser := &ClusterParser{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "fluentbit.fluent.io/v1alpha2",
			Kind:       "ClusterParser",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   "clusterparser1",
			Labels: labels,
		},
		Spec: ParserSpec{
			LTSV: &parser.LSTV{
				TimeKey:    "time",
				TimeFormat: "[%d/%b/%Y:%H:%M:%S %z]",
				Types:      "status:integer size:integer",
			},
		},
	}
	nsClusterParsers := []ClusterParserList{
		{
			Items: []ClusterParser{*nsClusterParser},
		},
	}

	config, err := cfg.RenderParserConfig(sl, clusterParsers, nsParsers, nsClusterParsers)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(config).To(Equal(expectedParsers))
}

func TestClusterFluentBitConfig_RenderMultilineParserConfig(t *testing.T) {
	g := NewGomegaWithT(t)
	sl := plugins.NewSecretLoader(nil, "testnamespace")

	clusterMultilineParser := &ClusterMultilineParser{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "fluentbit.fluent.io/v1alpha2",
			Kind:       "ClusterMultilineParser",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   "clustermultilineparser0",
			Labels: labels,
		},
		Spec: MultilineParserSpec{
			MultilineParser: &multilineparser.MultilineParser{
				Type:       "regex",
				Parser:     "go",
				KeyContent: "log",
			},
		},
	}
	clusterMultilineParsers := ClusterMultilineParserList{
		Items: []ClusterMultilineParser{*clusterMultilineParser},
	}

	nsMultilineParser := &MultilineParser{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "fluentbit.fluent.io/v1alpha2",
			Kind:       "MultilineParser",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "multilineparser0",
			Namespace: "testnamespace",
			Labels:    labels,
		},
		Spec: MultilineParserSpec{
			MultilineParser: &multilineparser.MultilineParser{
				Type:         "regex",
				FlushTimeout: 1000,
				Rules: []multilineparser.Rule{
					{
						Start: "start_state",
						Regex: `/(Dec \d+ \d+\:\d+\:\d+)(.*)/`,
						Next:  "cont",
					},
					{
						Start: "cont",
						Regex: `/^\s+at.*/`,
						Next:  "cont",
					},
				},
			},
		},
	}
	nsMultilineParsers := []MultilineParserList{
		{
			Items: []MultilineParser{*nsMultilineParser},
		},
	}

	nsClusterMultilineParser := &ClusterMultilineParser{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "fluentbit.fluent.io/v1alpha2",
			Kind:       "ClusterMultilineParser",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   "clustermultilineparser1",
			Labels: labels,
		},
		Spec: MultilineParserSpec{
			MultilineParser: &multilineparser.MultilineParser{
				Type:         "regex",
				FlushTimeout: 500,
				Rules: []multilineparser.Rule{
					{
						Start: "start_state",
						Regex: `/^(\d+ \d+\:\d+\:\d+)(.*)/`,
						Next:  "cont",
					},
					{
						Start: "cont",
						Regex: `/^\s+at.*/`,
						Next:  "cont",
					},
				},
			},
		},
	}
	nsClusterMultilineParsers := []ClusterMultilineParserList{
		{
			Items: []ClusterMultilineParser{*nsClusterMultilineParser},
		},
	}

	config, err := cfg.RenderMultilineParserConfig(
		sl, clusterMultilineParsers, nsMultilineParsers, nsClusterMultilineParsers,
	)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(config).To(Equal(expectedMultilineParsers))
}

func ptrBool(v bool) *bool {
	return &v
}

func ptrInt64(v int64) *int64 {
	return &v
}

func ptrInt32(v int32) *int32 {
	return &v
}

func ptrFloat64(v float64) *float64 {
	return &v
}
