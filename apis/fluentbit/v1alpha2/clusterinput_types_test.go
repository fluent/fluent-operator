package v1alpha2

import (
	"github.com/fluent/fluent-operator/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/apis/fluentbit/v1alpha2/plugins/input"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
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
    Tag    logs.foo.bar
    Inotify_Watcher    false
[Input]
    Name    dummy
    Alias    input2_alias
    Tag    logs.foo.bar
    Rate    3
    Samples    5
`

func TestClusterInputList_Load(t *testing.T) {
	g := NewGomegaWithT(t)

	sl := plugins.NewSecretLoader(nil, "testnamespace", nil)

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

	inputs := ClusterInputList{
		Items: []ClusterInput{*inputObj1, *inputObj2},
	}

	i := 0
	for i < 5 {
		clusterInputs, err := inputs.Load(sl)
		g.Expect(err).NotTo(HaveOccurred())
		g.Expect(clusterInputs).To(Equal(inputExpected))

		i++
	}
}
