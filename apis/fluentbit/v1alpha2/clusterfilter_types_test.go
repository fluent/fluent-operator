package v1alpha2

import (
	"testing"

	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins/filter"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestClusterFilterList_Load(t *testing.T) {
	filtersExpected := `[Filter]
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
[Filter]
    Name    kubernetes
    Match    logs.foo.bar
    Buffer_Size    10m
    Kube_URL    http://127.0.0.1:6443
    Kube_CA_File    root.ca
    Kube_CA_Path    /root/.kube/crt
    Labels    true
    Annotations    true
    DNS_Wait_Time    30
    Use_Kubelet    true
    Kubelet_Port    10000
    Kube_Meta_Cache_TTL    60s
[Filter]
    Name    throttle
    Match    *
    Alias    throttle.application-xy
    Rate    200
    Window    300
    Interval    1s
`
	g := NewGomegaWithT(t)

	sl := plugins.NewSecretLoader(nil, "testnamespace")

	labels := map[string]string{
		"label0": "lv0",
		"label1": "lv1",
		"label3": "lval3",
		"lbl2":   "lval2",
		"lbl1":   "lvl1",
	}

	filterObj1 := &ClusterFilter{
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
	filterObj2 := &ClusterFilter{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "fluentbit.fluent.io/v1alpha2",
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
						BufferSize:       "10m",
						KubeURL:          "http://127.0.0.1:6443",
						KubeCAFile:       "root.ca",
						KubeCAPath:       "/root/.kube/crt",
						Labels:           ptr(true),
						Annotations:      ptr(true),
						DNSWaitTime:      ptr[int32](30),
						UseKubelet:       ptr(true),
						KubeletPort:      ptr[int32](10000),
						KubeMetaCacheTTL: "60s",
					},
				},
			},
		},
	}
	filterObj3 := &ClusterFilter{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "fluentbit.fluent.io/v1alpha2",
			Kind:       "ClusterFilter",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   "filter3",
			Labels: labels,
		},
		Spec: FilterSpec{
			Match: "*",
			FilterItems: []FilterItem{
				{
					Throttle: &filter.Throttle{
						CommonParams: plugins.CommonParams{
							Alias: "throttle.application-xy",
						},
						Rate:     ptr[int64](200),
						Window:   ptr[int64](300),
						Interval: "1s",
					},
				},
			},
		},
	}
	filters := ClusterFilterList{
		Items: []ClusterFilter{*filterObj1, *filterObj2, *filterObj3},
	}

	i := 0
	for i < 5 {
		clusterFilters, err := filters.Load(sl)
		g.Expect(err).NotTo(HaveOccurred())
		g.Expect(clusterFilters).To(Equal(filtersExpected))

		i++
	}
}

func TestClusterFilterList_Load_With_Ordinals(t *testing.T) {
	filtersExpected := `[Filter]
    Name    grep
    Match    *
    Alias    first
    Regex    ^.*$
[Filter]
    Name    grep
    Match    *
    Alias    second
    Regex    ^.*$
[Filter]
    Name    grep
    Match    *
    Alias    third
    Regex    ^.*$
`

	g := NewGomegaWithT(t)
	sl := plugins.NewSecretLoader(nil, "testnamespace")

	filterObj1 := &ClusterFilter{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "fluentbit.fluent.io/v1alpha2",
			Kind:       "ClusterFilter",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "first",
		},
		Spec: FilterSpec{
			Match: "*",
			FilterItems: []FilterItem{
				{
					Grep: &filter.Grep{
						CommonParams: plugins.CommonParams{
							Alias: "second",
						},
						Regex: "^.*$",
					},
				},
			},
		},
	}

	filterObj2 := &ClusterFilter{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "fluentbit.fluent.io/v1alpha2",
			Kind:       "ClusterFilter",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "second",
		},
		Spec: FilterSpec{
			Ordinal: 10,
			Match:   "*",
			FilterItems: []FilterItem{
				{
					Grep: &filter.Grep{
						CommonParams: plugins.CommonParams{
							Alias: "third",
						},
						Regex: "^.*$",
					},
				},
			},
		},
	}

	filterObj3 := &ClusterFilter{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "fluentbit.fluent.io/v1alpha2",
			Kind:       "ClusterFilter",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "third",
		},
		Spec: FilterSpec{
			Ordinal: -10,
			Match:   "*",
			FilterItems: []FilterItem{
				{
					Grep: &filter.Grep{
						CommonParams: plugins.CommonParams{
							Alias: "first",
						},
						Regex: "^.*$",
					},
				},
			},
		},
	}

	filters := ClusterFilterList{
		Items: []ClusterFilter{*filterObj1, *filterObj2, *filterObj3},
	}

	i := 0
	for i < 5 {
		clusterFilters, err := filters.Load(sl)
		g.Expect(err).NotTo(HaveOccurred())

		g.Expect(clusterFilters).To(Equal(filtersExpected))

		i++
	}
}

func TestClusterFilter_RecordModifier_Generated(t *testing.T) {
	filtersExpected := `[Filter]
    Name    record_modifier
    Match    logs.foo.bar
    Record    hostname ${HOSTNAME}
    Record    product Awesome_Tool
    Remove_key    Swap.total
    Remove_key    Swap.free
    Remove_key    Swap.used
    Allowlist_key    Mem.total
    Allowlist_key    Mem.free
    Allowlist_key    Mem.used
    Whitelist_key    Disk.total
    Whitelist_key    Disk.free
    Whitelist_key    Disk.used
    Uuid_key    ID1
    Uuid_key    UiD2
`
	g := NewGomegaWithT(t)

	sl := plugins.NewSecretLoader(nil, "test namespace")

	rmFilter := &ClusterFilter{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "fluentbit.fluent.io/v1alpha2",
			Kind:       "ClusterFilter",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "filterRecordModifier",
		},
		Spec: FilterSpec{
			Match: "logs.foo.bar",
			FilterItems: []FilterItem{
				{
					RecordModifier: &filter.RecordModifier{
						CommonParams: plugins.CommonParams{},
						Records: []string{
							"hostname ${HOSTNAME}",
							"product Awesome_Tool",
						},
						RemoveKeys: []string{
							"Swap.total",
							"Swap.free",
							"Swap.used",
						},
						AllowlistKeys: []string{
							"Mem.total",
							"Mem.free",
							"Mem.used",
						},
						WhitelistKeys: []string{
							"Disk.total",
							"Disk.free",
							"Disk.used",
						},
						UUIDKeys: []string{
							"ID1",
							"UiD2",
						},
					},
				},
			},
		},
	}

	filters := ClusterFilterList{
		Items: []ClusterFilter{*rmFilter},
	}

	clusterFilters, err := filters.Load(sl)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(clusterFilters).To(Equal(filtersExpected))
}

func TestClusterFilterList_Load_As_Yaml(t *testing.T) {
	filtersExpected := `filters:
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
  - name: kubernetes
    match: "logs.foo.bar"
    buffer_size: 10m
    kube_url: http://127.0.0.1:6443
    kube_ca_file: root.ca
    kube_ca_path: /root/.kube/crt
    labels: true
    annotations: true
    dns_wait_time: 30
    use_kubelet: true
    kubelet_port: 10000
    kube_meta_cache_ttl: 60s
  - name: throttle
    match: "*"
    alias: throttle.application-xy
    rate: 200
    window: 300
    interval: 1s
`
	g := NewGomegaWithT(t)

	sl := plugins.NewSecretLoader(nil, "testnamespace")

	labels := map[string]string{
		"label0": "lv0",
		"label1": "lv1",
		"label3": "lval3",
		"lbl2":   "lval2",
		"lbl1":   "lvl1",
	}

	filterObj1 := &ClusterFilter{
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
	filterObj2 := &ClusterFilter{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "fluentbit.fluent.io/v1alpha2",
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
						BufferSize:       "10m",
						KubeURL:          "http://127.0.0.1:6443",
						KubeCAFile:       "root.ca",
						KubeCAPath:       "/root/.kube/crt",
						Labels:           ptr(true),
						Annotations:      ptr(true),
						DNSWaitTime:      ptr[int32](30),
						UseKubelet:       ptr(true),
						KubeletPort:      ptr[int32](10000),
						KubeMetaCacheTTL: "60s",
					},
				},
			},
		},
	}
	filterObj3 := &ClusterFilter{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "fluentbit.fluent.io/v1alpha2",
			Kind:       "ClusterFilter",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   "filter3",
			Labels: labels,
		},
		Spec: FilterSpec{
			Match: "*",
			FilterItems: []FilterItem{
				{
					Throttle: &filter.Throttle{
						CommonParams: plugins.CommonParams{
							Alias: "throttle.application-xy",
						},
						Rate:     ptr[int64](200),
						Window:   ptr[int64](300),
						Interval: "1s",
					},
				},
			},
		},
	}
	filters := ClusterFilterList{
		Items: []ClusterFilter{*filterObj1, *filterObj2, *filterObj3},
	}

	i := 0
	for i < 5 {
		clusterFilters, err := filters.LoadAsYaml(sl, 0)
		g.Expect(err).NotTo(HaveOccurred())
		g.Expect(clusterFilters).To(Equal(filtersExpected))

		i++
	}
}

func TestClusterFilter_RecordModifier_Generated_Load_As_Yaml(t *testing.T) {
	filtersExpected := `filters:
  - name: record_modifier
    match: "logs.foo.bar"
    record:
      - hostname ${HOSTNAME}
      - product Awesome_Tool
    remove_key:
      - Swap.total
      - Swap.free
      - Swap.used
    allowlist_key:
      - Mem.total
      - Mem.free
      - Mem.used
    whitelist_key:
      - Disk.total
      - Disk.free
      - Disk.used
    uuid_key:
      - ID1
      - UiD2
`
	g := NewGomegaWithT(t)

	sl := plugins.NewSecretLoader(nil, "test namespace")

	rmFilter := &ClusterFilter{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "fluentbit.fluent.io/v1alpha2",
			Kind:       "ClusterFilter",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "filterRecordModifier",
		},
		Spec: FilterSpec{
			Match: "logs.foo.bar",
			FilterItems: []FilterItem{
				{
					RecordModifier: &filter.RecordModifier{
						CommonParams: plugins.CommonParams{},
						Records: []string{
							"hostname ${HOSTNAME}",
							"product Awesome_Tool",
						},
						RemoveKeys: []string{
							"Swap.total",
							"Swap.free",
							"Swap.used",
						},
						AllowlistKeys: []string{
							"Mem.total",
							"Mem.free",
							"Mem.used",
						},
						WhitelistKeys: []string{
							"Disk.total",
							"Disk.free",
							"Disk.used",
						},
						UUIDKeys: []string{
							"ID1",
							"UiD2",
						},
					},
				},
			},
		},
	}

	filters := ClusterFilterList{
		Items: []ClusterFilter{*rmFilter},
	}

	clusterFilters, err := filters.LoadAsYaml(sl, 0)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(clusterFilters).To(Equal(filtersExpected))
}
