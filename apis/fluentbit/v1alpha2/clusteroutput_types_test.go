package v1alpha2

import (
	"testing"

	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins/output"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	outputExpectedYaml = `outputs:
  - name: http
    match: "logs.foo.bar"
    alias: output_http_alias
    host: https://example2.com
    port: 433
    uri: /logs
    format: json_lines
    header:
      -  Authorization    foo:bar
      -  X-Log-Header-0    testing
      -  X-Log-Header-App-ID    9780495d9db3
      -  X-Log-Header-App-Name    app_name
    json_date_key: timestamp
    json_date_format: iso8601
    tls: On
    tls.verify: true
  - name: null
    match: "*"
  - name: opensearch
    match: "*"
    alias: output_opensearch_alias
    host: https://example2.com
    port: 9200
    index: my_index
    type: my_type
  - name: prometheus_remote_write
    match: "logs.foo.bar"
    alias: output_prometheus_remote_write_alias
    host: https://example3.com
    port: 433
    proxy: https://proxy:533
    uri: /prometheus/v1/write?prometheus_server=YOUR_DATA_SOURCE_NAME
    header:
      -  Authorization    foo:bar
      -  X-Log-Header-0    testing
      -  X-Log-Header-App-ID    9780495d9db3
      -  X-Log-Header-App-Name    app_name
    log_response_payload: true
    add_label:
      -  app    fluent-bit
      -  color    blue
    workers: 3
    tls: On
    tls.verify: true
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
`
	outputExpected = `[Output]
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
)

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
				Port: ptr[int32](3300),
				Mode: "tls",
				TLS: &plugins.TLS{
					Verify: ptr(true),
				},
				SyslogMessageKey:  "log",
				SyslogHostnameKey: "do_app_name",
				SyslogAppnameKey:  "do_component_name",
			},
		},
	}

	headers := map[string]string{
		"Authorization":         authorization,
		"X-Log-Header-App-Name": xLogHeaderAppName,
		"X-Log-Header-0":        xLogHeader0,
		"X-Log-Header-App-ID":   xLogHeaderAppID,
	}

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
				Port:           ptr[int32](433),
				Uri:            "/logs",
				Headers:        headers,
				Format:         "json_lines",
				JsonDateKey:    "timestamp",
				JsonDateFormat: "iso8601",
				TLS: &plugins.TLS{
					Verify: ptr(true),
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
				Port:  ptr[int32](9200),
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
				Port:               ptr[int32](433),
				URI:                "/prometheus/v1/write?prometheus_server=YOUR_DATA_SOURCE_NAME",
				Proxy:              "https://proxy:533",
				Headers:            headers,
				LogResponsePayload: ptr(true),
				AddLabels:          addLabels,
				Workers:            ptr[int32](3),
				TLS: &plugins.TLS{
					Verify: ptr(true),
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

func TestClusterOutputList_Load_As_Yaml(t *testing.T) {
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
				Port: ptr[int32](3300),
				Mode: "tls",
				TLS: &plugins.TLS{
					Verify: ptr(true),
				},
				SyslogMessageKey:  "log",
				SyslogHostnameKey: "do_app_name",
				SyslogAppnameKey:  "do_component_name",
			},
		},
	}

	headers := map[string]string{
		"Authorization":         authorization,
		"X-Log-Header-App-Name": xLogHeaderAppName,
		"X-Log-Header-0":        xLogHeader0,
		"X-Log-Header-App-ID":   xLogHeaderAppID,
	}

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
				Port:           ptr[int32](433),
				Uri:            "/logs",
				Headers:        headers,
				Format:         "json_lines",
				JsonDateKey:    "timestamp",
				JsonDateFormat: "iso8601",
				TLS: &plugins.TLS{
					Verify: ptr(true),
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
				Port:  ptr[int32](9200),
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
				Port:               ptr[int32](433),
				URI:                "/prometheus/v1/write?prometheus_server=YOUR_DATA_SOURCE_NAME",
				Proxy:              "https://proxy:533",
				Headers:            headers,
				LogResponsePayload: ptr(true),
				AddLabels:          addLabels,
				Workers:            ptr[int32](3),
				TLS: &plugins.TLS{
					Verify: ptr(true),
				},
			},
		},
	}
	nullOutput := ClusterOutput{TypeMeta: metav1.TypeMeta{
		APIVersion: "fluentbit.fluent.io/v1alpha2",
		Kind:       "ClusterOutput",
	},
		ObjectMeta: metav1.ObjectMeta{
			Name:   "null",
			Labels: labels,
		},
		Spec: OutputSpec{
			Match: "*",
			Null:  &output.Null{},
		},
	}

	outputs := ClusterOutputList{
		Items: []ClusterOutput{prometheusRemoteWriteOutput, syslogOut, httpOutput, openSearchOutput, nullOutput},
	}

	i := 0
	for i < 5 {
		clusterInputs, err := outputs.LoadAsYaml(sl, 0)
		g.Expect(err).NotTo(HaveOccurred())
		g.Expect(clusterInputs).To(Equal(outputExpectedYaml))

		i++
	}
}

func TestLokiOutputWithStructuredMetadata_Load(t *testing.T) {
	g := NewGomegaWithT(t)
	sl := plugins.NewSecretLoader(nil, "testnamespace")

	lokiOutput := ClusterOutput{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "fluentbit.fluent.io/v1alpha2",
			Kind:       "ClusterOutput",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "loki_output_with_metadata",
		},
		Spec: OutputSpec{
			Match: "kube.*",
			Loki: &output.Loki{
				Host: "loki-gateway",
				Port: ptr[int32](3100),
				Labels: []string{
					"job=fluentbit",
					"environment=production",
				},
				StructuredMetadata: map[string]string{
					"pod":       "${record['kubernetes']['pod_name']}",
					"container": "${record['kubernetes']['container_name']}",
					"trace_id":  "${record['trace_id']}",
				},
				StructuredMetadataKeys: []string{
					"level",
					"caller",
				},
			},
		},
	}

	outputs := ClusterOutputList{
		Items: []ClusterOutput{lokiOutput},
	}

	expected := `[Output]
    Name    loki
    Match    kube.*
    host    loki-gateway
    port    3100
    labels    environment=production,job=fluentbit
    structured_metadata    container=${record['kubernetes']['container_name']},pod=${record['kubernetes']['pod_name']},trace_id=${record['trace_id']}
    structured_metadata_keys    level,caller
`

	result, err := outputs.Load(sl)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(result).To(Equal(expected))
}

func TestLokiOutputWithStructuredMetadata_LoadAsYaml(t *testing.T) {
	g := NewGomegaWithT(t)
	sl := plugins.NewSecretLoader(nil, "testnamespace")

	lokiOutput := ClusterOutput{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "fluentbit.fluent.io/v1alpha2",
			Kind:       "ClusterOutput",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "loki_output_with_metadata",
		},
		Spec: OutputSpec{
			Match: "kube.*",
			Loki: &output.Loki{
				Host: "loki-gateway",
				Port: ptr[int32](3100),
				Labels: []string{
					"job=fluentbit",
					"environment=production",
				},
				StructuredMetadata: map[string]string{
					"pod":       "${record['kubernetes']['pod_name']}",
					"container": "${record['kubernetes']['container_name']}",
					"trace_id":  "${record['trace_id']}",
				},
				StructuredMetadataKeys: []string{
					"level",
					"caller",
				},
			},
		},
	}

	outputs := ClusterOutputList{
		Items: []ClusterOutput{lokiOutput},
	}

	expected := `outputs:
  - name: loki
    match: "kube.*"
    host: loki-gateway
    port: 3100
    labels: environment=production,job=fluentbit
    structured_metadata: container=${record['kubernetes']['container_name']},pod=${record['kubernetes']['pod_name']},trace_id=${record['trace_id']}
    structured_metadata_keys: level,caller
`

	result, err := outputs.LoadAsYaml(sl, 0)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(result).To(Equal(expected))
}
