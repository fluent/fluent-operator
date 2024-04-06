package v1alpha2

import (
	"fmt"
	"github.com/fluent/fluent-operator/v2/apis/fluentbit/v1alpha2/plugins/multilineparser"
	"github.com/fluent/fluent-operator/v2/apis/fluentbit/v1alpha2/plugins/parser"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/fluent/fluent-operator/v2/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v2/apis/fluentbit/v1alpha2/plugins/custom"
	"github.com/fluent/fluent-operator/v2/apis/fluentbit/v1alpha2/plugins/filter"
	"github.com/fluent/fluent-operator/v2/apis/fluentbit/v1alpha2/plugins/input"
	"github.com/fluent/fluent-operator/v2/apis/fluentbit/v1alpha2/plugins/output"
	. "github.com/onsi/gomega"
)

var expected = `[Service]
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
    header     Authorization    foo:bar
    header     X-Log-Header-0    testing
    header     X-Log-Header-App-ID    9780495d9db3
    header     X-Log-Header-App-Name    app_name
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
`
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
			FlushSeconds: ptrInt64(1),
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

	headers["Authorization"] = "foo:bar"
	headers["X-Log-Header-App-Name"] = "app_name"
	headers["X-Log-Header-0"] = "testing"
	headers["X-Log-Header-App-ID"] = "9780495d9db3"

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
}

func TestClusterFluentBitConfig_RenderMainConfig_WithParsersFiles(t *testing.T) {
	g := NewGomegaWithT(t)
	sl := plugins.NewSecretLoader(nil, "testnamespace")

	cfbc := ClusterFluentBitConfig{
		Spec: FluentBitConfigSpec{
			Service: &Service{
				Daemon:       ptrBool(false),
				FlushSeconds: ptrInt64(1),
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
