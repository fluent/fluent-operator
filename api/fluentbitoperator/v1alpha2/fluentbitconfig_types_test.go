package v1alpha2

import (
	"testing"

	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kubesphere.io/fluentbit-operator/api/fluentbitoperator/v1alpha2/plugins"
	"kubesphere.io/fluentbit-operator/api/fluentbitoperator/v1alpha2/plugins/filter"
	"kubesphere.io/fluentbit-operator/api/fluentbitoperator/v1alpha2/plugins/input"
	"kubesphere.io/fluentbit-operator/api/fluentbitoperator/v1alpha2/plugins/output"
)

var expected = `[Service]
    Daemon    false
    Flush    1
    Grace    30
    Http_Server    true
    Log_Level    info
    Parsers_File    parsers.conf
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

func Test_FluentBitConfig_RenderMainConfig(t *testing.T) {
	g := NewGomegaWithT(t)

	sl := plugins.NewSecretLoader(nil, "testnamespace", nil)

	disableInotifyWatcher := ptrBool(true)

	labels := map[string]string{
		"label0": "lv0",
		"label1": "lv1",
		"label3": "lval3",
		"lbl2":   "lval2",
		"lbl1":   "lvl1",
	}

	inputObj := &Input{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "logging.kubesphere.io/v1alpha2",
			Kind:       "Input",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "input0",
			Namespace: "testnamespace",
			Labels:    labels,
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
			}},
	}

	inputs := InputList{
		Items: []Input{*inputObj},
	}

	filterObj := &Filter{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "logging.kubesphere.io/v1alpha2",
			Kind:       "Filter",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "filter0",
			Namespace: "testnamespace",
			Labels:    labels,
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

	filters := FilterList{
		Items: []Filter{*filterObj},
	}

	syslogOut := Output{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "logging.kubesphere.io/v1alpha2",
			Kind:       "Output",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "syslog_output0",
			Namespace: "testnamespace",
			Labels:    labels,
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

	httpOutput := Output{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "logging.kubesphere.io/v1alpha2",
			Kind:       "Output",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "http_output_0",
			Namespace: "testnamespace",
			Labels:    labels,
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

	outputs := OutputList{
		Items: []Output{syslogOut, httpOutput},
	}

	cfg := FluentBitConfig{
		Spec: FluentBitConfigSpec{Service: &Service{
			Daemon:       ptrBool(false),
			FlushSeconds: ptrInt64(1),
			GraceSeconds: ptrInt64(30),
			HttpServer:   ptrBool(true),
			LogLevel:     "info",
			ParsersFile:  "parsers.conf",
		}},
	}

	// we should not see any permutations in serialized config
	i := 0
	for i < 5 {
		config, err := cfg.RenderMainConfig(sl, inputs, filters, outputs)
		g.Expect(err).NotTo(HaveOccurred())
		g.Expect(config).To(Equal(expected))

		i++
	}
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
