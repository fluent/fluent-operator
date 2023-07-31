package v1alpha2

import (
	"testing"

	"github.com/fluent/fluent-operator/v2/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v2/apis/fluentbit/v1alpha2/plugins/input"
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
				DisableInotifyWatcher:  ptrBool(true),
				Tag:                    "logs.foo.bar",
				Path:                   "/logs/containers/apps0",
				ExcludePath:            "/logs/containers/exclude_path",
				SkipLongLines:          ptrBool(true),
				IgnoreOlder:            "5m",
				MemBufLimit:            "5MB",
				RefreshIntervalSeconds: ptrInt64(10),
				DB:                     "/fluent-bit/tail/pos.db",
				Parser:                 "docker",
				DockerMode:             ptrBool(true),
				DockerModeFlushSeconds: ptrInt64(4),
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
				Rate:    ptrInt32(3),
				Samples: ptrInt32(5),
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
				Port:           ptrInt32(int32(433)),
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
				ScrapeOnStart:  ptrBool(true),
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
				Port:            ptrInt32(int32(433)),
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
