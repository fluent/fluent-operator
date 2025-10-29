package v1alpha2

import (
	"testing"

	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins/input"
	"github.com/fluent/fluent-operator/v3/pkg/utils"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var inputExpected = `[Input]
    Name    tail
    Alias    input0_alias
    Path    /logs/containers/apps0
    Exclude_Path    /logs/containers/exclude_path
    Refresh_Interval    10
    Ignore_Older    5m
    Skip_Long_Lines    true
    DB    /fluent-bit/tail/pos.db
    Mem_Buf_Limit    5MB
    Parser    docker
    Tag    logs.foo.bar
    Docker_Mode    true
    Docker_Mode_Flush    4
    Docker_Mode_Parser    docker-mode-parser
    Inotify_Watcher    false
[Input]
    Name    dummy
    Alias    input2_alias
    Tag    logs.foo.bar
    Rate    3
    Samples    5
[Input]
    Name    prometheus_scrape
    Alias    input3_alias
    tag    logs.foo.bar
    host    https://example3.com
    port    433
    scrape_interval    10s
    metrics_path    /metrics
`

var inputExpectedYaml = `inputs:
  - name: dummy
    alias: input0_alias
    processors:
      logs:
        - add: hostname test
          name: modify
        - call: append_tag
          code: |-
            function append_tag(tag, timestamp, record)
                new_record = record
                new_record["tag"] = tag
                return 1, timestamp, new_record
            end
          name: lua
    tag: logs.foo.bar
    dummy: {"key":"value"}
`

func TestClusterInputList_Load(t *testing.T) {
	g := NewGomegaWithT(t)

	sl := plugins.NewSecretLoader(nil, "testnamespace")

	labels := map[string]string{
		"label0": "lv0",
		"label1": "lv1",
		"label3": "lval3",
		"lbl2":   "lval2",
		"lbl1":   "lvl1",
	}

	inputObj1 := &ClusterInput{
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
				DisableInotifyWatcher:  utils.ToPtr(true),
				Tag:                    "logs.foo.bar",
				Path:                   "/logs/containers/apps0",
				ExcludePath:            "/logs/containers/exclude_path",
				SkipLongLines:          utils.ToPtr(true),
				IgnoreOlder:            "5m",
				MemBufLimit:            "5MB",
				RefreshIntervalSeconds: utils.ToPtr[int64](10),
				DB:                     "/fluent-bit/tail/pos.db",
				Parser:                 "docker",
				DockerMode:             utils.ToPtr(true),
				DockerModeFlushSeconds: utils.ToPtr[int64](4),
				DockerModeParser:       "docker-mode-parser",
			},
		},
	}
	inputObj2 := &ClusterInput{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "fluentbit.fluent.io/v1alpha2",
			Kind:       "ClusterInput",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   "input2",
			Labels: labels,
		},
		Spec: InputSpec{
			Alias: "input2_alias",
			Dummy: &input.Dummy{
				Tag:     "logs.foo.bar",
				Rate:    utils.ToPtr[int32](3),
				Samples: utils.ToPtr[int32](5),
			},
		},
	}
	inputObj3 := &ClusterInput{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "fluentbit.fluent.io/v1alpha2",
			Kind:       "ClusterInput",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   "input3",
			Labels: labels,
		},
		Spec: InputSpec{
			Alias: "input3_alias",
			PrometheusScrapeMetrics: &input.PrometheusScrapeMetrics{
				Tag:            "logs.foo.bar",
				Host:           "https://example3.com",
				Port:           utils.ToPtr[int32](433),
				ScrapeInterval: "10s",
				MetricsPath:    "/metrics",
			},
		},
	}

	inputs := ClusterInputList{
		Items: []ClusterInput{*inputObj1, *inputObj2, *inputObj3},
	}

	i := 0
	for i < 5 {
		clusterInputs, err := inputs.Load(sl)
		g.Expect(err).NotTo(HaveOccurred())
		g.Expect(clusterInputs).To(Equal(inputExpected))

		i++
	}
}

var fluentbitExpected = `[Input]
    Name    fluentbit_metrics
    Alias    input0_alias
    Tag    logs.foo.bar
    scrape_interval    2
    scrape_on_start    true
[Input]
    Name    forward
    Alias    input1_alias
    Port    433
    Listen    0.0.0.0
    Buffer_Chunk_Size    1M
    Buffer_Max_Size    6M
    threaded    on
`

func TestFluentbitMetricClusterInputList_Load(t *testing.T) {
	g := NewGomegaWithT(t)

	sl := plugins.NewSecretLoader(nil, "testnamespace")

	labels := map[string]string{
		"label0": "lv0",
		"label1": "lv1",
		"label3": "lval3",
		"lbl2":   "lval2",
		"lbl1":   "lvl1",
	}

	inputObj1 := &ClusterInput{
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
			FluentBitMetrics: &input.FluentbitMetrics{
				Tag:            "logs.foo.bar",
				ScrapeInterval: "2",
				ScrapeOnStart:  utils.ToPtr(true),
			},
		},
	}
	inputObj2 := &ClusterInput{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "fluentbit.fluent.io/v1alpha2",
			Kind:       "ClusterInput",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   "input1",
			Labels: labels,
		},
		Spec: InputSpec{
			Alias: "input1_alias",
			Forward: &input.Forward{
				Port:            utils.ToPtr[int32](433),
				Listen:          "0.0.0.0",
				BufferChunkSize: "1M",
				BufferMaxSize:   "6M",
				Threaded:        "on",
			},
		},
	}
	inputs := ClusterInputList{
		Items: []ClusterInput{*inputObj1, *inputObj2},
	}

	i := 0
	for i < 5 {
		clusterInputs, err := inputs.Load(sl)
		g.Expect(err).NotTo(HaveOccurred())
		g.Expect(clusterInputs).To(Equal(fluentbitExpected))
		i++
	}
}

func TestClusterInputList_Load_As_Yaml(t *testing.T) {
	g := NewGomegaWithT(t)

	sl := plugins.NewSecretLoader(nil, "testnamespace")

	labels := map[string]string{
		"label0": "lv0",
		"label1": "lv1",
		"label3": "lval3",
		"lbl2":   "lval2",
		"lbl1":   "lvl1",
	}

	inputObj1 := &ClusterInput{
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
				DisableInotifyWatcher:  utils.ToPtr(true),
				Tag:                    "logs.foo.bar",
				Path:                   "/logs/containers/apps0",
				ExcludePath:            "/logs/containers/exclude_path",
				SkipLongLines:          utils.ToPtr(true),
				IgnoreOlder:            "5m",
				MemBufLimit:            "5MB",
				RefreshIntervalSeconds: utils.ToPtr[int64](10),
				DB:                     "/fluent-bit/tail/pos.db",
				Parser:                 "docker",
				DockerMode:             utils.ToPtr(true),
				DockerModeFlushSeconds: utils.ToPtr[int64](4),
				DockerModeParser:       "docker-mode-parser",
			},
		},
	}
	inputObj2 := &ClusterInput{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "fluentbit.fluent.io/v1alpha2",
			Kind:       "ClusterInput",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   "input2",
			Labels: labels,
		},
		Spec: InputSpec{
			Alias: "input2_alias",
			Dummy: &input.Dummy{
				Tag:     "logs.foo.bar",
				Rate:    utils.ToPtr[int32](3),
				Samples: utils.ToPtr[int32](5),
			},
		},
	}
	inputObj3 := &ClusterInput{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "fluentbit.fluent.io/v1alpha2",
			Kind:       "ClusterInput",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   "input3",
			Labels: labels,
		},
		Spec: InputSpec{
			Alias: "input3_alias",
			PrometheusScrapeMetrics: &input.PrometheusScrapeMetrics{
				Tag:            "logs.foo.bar",
				Host:           "https://example3.com",
				Port:           utils.ToPtr[int32](433),
				ScrapeInterval: "10s",
				MetricsPath:    "/metrics",
			},
		},
	}

	expectedYaml := `inputs:
  - name: tail
    alias: input0_alias
    path: /logs/containers/apps0
    exclude_path: /logs/containers/exclude_path
    refresh_interval: 10
    ignore_older: 5m
    skip_long_lines: true
    db: /fluent-bit/tail/pos.db
    mem_buf_limit: 5MB
    parser: docker
    tag: logs.foo.bar
    docker_mode: true
    docker_mode_flush: 4
    docker_mode_parser: docker-mode-parser
    inotify_watcher: false
  - name: dummy
    alias: input2_alias
    tag: logs.foo.bar
    rate: 3
    samples: 5
  - name: prometheus_scrape
    alias: input3_alias
    tag: logs.foo.bar
    host: https://example3.com
    port: 433
    scrape_interval: 10s
    metrics_path: /metrics
`
	inputs := ClusterInputList{
		Items: []ClusterInput{*inputObj1, *inputObj2, *inputObj3},
	}

	i := 0
	for i < 5 {
		clusterInputs, err := inputs.LoadAsYaml(sl, 0)
		g.Expect(err).NotTo(HaveOccurred())
		g.Expect(clusterInputs).To(Equal(expectedYaml))

		i++
	}
}

func TestClusterInputListProcessors_Load_As_Yaml(t *testing.T) {
	g := NewGomegaWithT(t)
	sl := plugins.NewSecretLoader(nil, "testnamespace")
	inputObj1 := &ClusterInput{
		TypeMeta: metav1.TypeMeta{

			APIVersion: "fluentbit.fluent.io/v1alpha2",
			Kind:       "ClusterInput",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   "random",
			Labels: labels,
		},
		Spec: InputSpec{
			Alias: "input0_alias",
			Dummy: &input.Dummy{
				Tag:   "logs.foo.bar",
				Dummy: "{\"key\":\"value\"}",
			},
			Processors: &plugins.Config{Data: map[string]interface{}{
				"logs": []interface{}{
					map[string]interface{}{"add": "hostname test", "name": "modify"},
					map[string]interface{}{"name": "lua", "call": "append_tag", "code": `function append_tag(tag, timestamp, record)
    new_record = record
    new_record["tag"] = tag
    return 1, timestamp, new_record
end`},
				},
			},
			},
		},
	}
	inputs := ClusterInputList{
		Items: []ClusterInput{*inputObj1},
	}
	i := 0
	for i < 5 {
		clusterInputs, err := inputs.LoadAsYaml(sl, 0)
		g.Expect(err).NotTo(HaveOccurred())
		g.Expect(clusterInputs).To(Equal(inputExpectedYaml))

		i++
	}
}
