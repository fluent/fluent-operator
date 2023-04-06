package v1alpha2

import (
	"testing"

	"github.com/fluent/fluent-operator/v2/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v2/apis/fluentbit/v1alpha2/plugins/output"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var outputExpected = `[Output]
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
    Name    opensearch
    Match    *
    Alias    output_opensearch_alias
    Host    https://example2.com
    Port    9200
    Index    my_index
    Type    my_type
[Output]
    Name    prometheus_remote_write
    Match    logs.foo.bar
    Alias    output_prometheus_remote_write_alias
    host    https://example3.com
    port    433
    proxy    https://proxy:533
    uri    /prometheus/v1/write?prometheus_server=YOUR_DATA_SOURCE_NAME
    header     Authorization    foo:bar
    header     X-Log-Header-0    testing
    header     X-Log-Header-App-ID    9780495d9db3
    header     X-Log-Header-App-Name    app_name
    log_response_payload    true
    add_label     app    fluent-bit
    add_label     color    blue
    workers    3
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

func TestClusterOutputList_Load(t *testing.T) {
	g := NewGomegaWithT(t)

	sl := plugins.NewSecretLoader(nil, "testnamespace")

	labels := map[string]string{
		"label0": "lv0",
		"label1": "lv1",
		"label3": "lval3",
		"lbl2":   "lval2",
		"lbl1":   "lvl1",
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

	addLabels := map[string]string{}

	addLabels["app"] = "fluent-bit"
	addLabels["color"] = "blue"
	prometheusRemoteWriteOutput := ClusterOutput{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "fluentbit.fluent.io/v1alpha2",
			Kind:       "ClusterOutput",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   "prometheus_remote_write_output_0",
			Labels: labels,
		},
		Spec: OutputSpec{
			Alias: "output_prometheus_remote_write_alias",
			Match: "logs.foo.bar",
			PrometheusRemoteWrite: &output.PrometheusRemoteWrite{
				Host:               "https://example3.com",
				Port:               ptrInt32(int32(433)),
				URI:                "/prometheus/v1/write?prometheus_server=YOUR_DATA_SOURCE_NAME",
				Proxy:              "https://proxy:533",
				Headers:            headers,
				LogResponsePayload: ptrBool(true),
				AddLabels:          addLabels,
				Workers:            ptrInt32(int32(3)),
				TLS: &plugins.TLS{
					Verify: ptrBool(true),
				},
			},
		},
	}

	outputs := ClusterOutputList{
		Items: []ClusterOutput{prometheusRemoteWriteOutput, syslogOut, httpOutput, openSearchOutput},
	}

	i := 0
	for i < 5 {
		clusterInputs, err := outputs.Load(sl)
		g.Expect(err).NotTo(HaveOccurred())
		g.Expect(clusterInputs).To(Equal(outputExpected))

		i++
	}
}
