package v1alpha3

import (
	"github.com/fluent/fluent-operator/apis/fluentbit/v1alpha3/plugins"
	"github.com/fluent/fluent-operator/apis/fluentbit/v1alpha3/plugins/filter"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
	"testing"
)

var filtersExpected = `[Filter]
    Name    modify
    Match    logs.foo.bar
    Condition    Key_value_equals    kve0    kvev0
    Condition    Key_value_equals    kve1    kvev1
    Condition    Key_value_equals    kve2    kvev2
    Condition    Key_does_not_exist    kdn1
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
[Filter]
    Name    kubernetes
    Match    logs.foo.bar
    Buffer_Size    10m
    Kube_URL    http://127.0.0.1:6443
    Kube_CA_File    root.ca
    Kube_CA_Path    /root/.kube/crt
    Labels    true
    Annotations    true
`

func TestClusterFilterList_Load(t *testing.T) {
	g := NewGomegaWithT(t)

	sl := plugins.NewSecretLoader(nil, "testnamespace", nil)

	labels := map[string]string{
		"label0": "lv0",
		"label1": "lv1",
		"label3": "lval3",
		"lbl2":   "lval2",
		"lbl1":   "lvl1",
	}

	filterObj1 := &ClusterFilter{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "fluentbit.fluent.io/v1alpha3",
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
								KeyDoesNotExist: "kdn1",
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
	filterObj2 := &ClusterFilter{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "fluentbit.fluent.io/v1alpha3",
			Kind:       "ClusterFilter",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   "filter2",
			Labels: labels,
		},
		Spec: FilterSpec{
			Match: "logs.foo.bar",
			FilterItems: []FilterItem{
				{
					Kubernetes: &filter.Kubernetes{
						BufferSize:  "10m",
						KubeURL:     "http://127.0.0.1:6443",
						KubeCAFile:  "root.ca",
						KubeCAPath:  "/root/.kube/crt",
						Labels:      pointer.Bool(true),
						Annotations: pointer.Bool(true),
					},
				},
			},
		},
	}

	filters := ClusterFilterList{
		Items: []ClusterFilter{*filterObj1, *filterObj2},
	}

	i := 0
	for i < 5 {
		clusterFilters, err := filters.Load(sl)
		g.Expect(err).NotTo(HaveOccurred())
		g.Expect(clusterFilters).To(Equal(filtersExpected))

		i++
	}
}
