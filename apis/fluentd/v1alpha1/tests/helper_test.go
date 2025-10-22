package cfgrender

import (
	"fmt"
	"strings"
	"testing"

	"github.com/go-logr/logr"
	"github.com/go-openapi/errors"

	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	fluentdv1alpha1 "github.com/fluent/fluent-operator/v3/apis/fluentd/v1alpha1"
	"github.com/fluent/fluent-operator/v3/apis/fluentd/v1alpha1/plugins"
)

const (
	maxRuntimes = 5
)

func Test_Cfg2ES(t *testing.T) {
	sl := plugins.NewSecretLoader(nil, Fluentd.Namespace, logr.Logger{})
	testNamespacedConfig(t, sl, Fluentd, &FluentdConfig1, []fluentdv1alpha1.ClusterOutput{FluentdclusterOutput2ES}, "./expected/fluentd-namespaced-cfg-output-es.cfg")
}

func Test_ClusterCfgInputTail(t *testing.T) {
	sl := plugins.NewSecretLoader(nil, Fluentd.Namespace, logr.Logger{})
	testClusterConfigWithGlobalInputs(t, sl, FluentdInputTail, &FluentdConfig1, []fluentdv1alpha1.ClusterOutput{FluentdClusterOutputTag}, "./expected/fluentd-global-cfg-input-tail.cfg")
}

func Test_ClusterCfgInputSample(t *testing.T) {
	sl := plugins.NewSecretLoader(nil, Fluentd.Namespace, logr.Logger{})
	testClusterConfigWithGlobalInputs(t, sl, FluentdInputSample, &FluentdConfig1, []fluentdv1alpha1.ClusterOutput{FluentdClusterOutputTag}, "./expected/fluentd-global-cfg-input-sample.cfg")
}

func Test_ClusterCfgInputMonitorAgent(t *testing.T) {
	sl := plugins.NewSecretLoader(nil, Fluentd.Namespace, logr.Logger{})
	testClusterConfigWithGlobalInputs(t, sl, FluentdInputMonitorAgent, &FluentdConfig1, []fluentdv1alpha1.ClusterOutput{FluentdClusterOutputTag}, "./expected/fluentd-global-cfg-input-monitorAgent.cfg")
}

func Test_ClusterCfgOutput2ES(t *testing.T) {
	sl := plugins.NewSecretLoader(nil, Fluentd.Namespace, logr.Logger{})
	testClusterConfigWithFiltersAndOutputs(t, sl, Fluentd, &FluentdClusterFluentdConfig1, []fluentdv1alpha1.ClusterFilter{}, []fluentdv1alpha1.ClusterOutput{FluentdclusterOutput2ES}, "./expected/fluentd-cluster-cfg-output-es.cfg", false)
}

func Test_ClusterCfgOutput2ESDataStream(t *testing.T) {
	sl := plugins.NewSecretLoader(nil, Fluentd.Namespace, logr.Logger{})
	testClusterConfigWithFiltersAndOutputs(t, sl, Fluentd, &FluentdClusterFluentdConfig1, []fluentdv1alpha1.ClusterFilter{}, []fluentdv1alpha1.ClusterOutput{FluentdclusterOutput2ESDataStream}, "./expected/fluentd-cluster-cfg-output-es-data-stream.cfg", false)
}

func Test_ClusterCfgOutput2CopyESDataStream(t *testing.T) {
	g := NewGomegaWithT(t)
	sl := NewSecretLoader(logr.Logger{}, esCredentials)

	psr := fluentdv1alpha1.NewGlobalPluginResources("main")
	psr.CombineGlobalInputsPlugins(sl, Fluentd.Spec.GlobalInputs)

	clustercfgRouter, err := psr.BuildCfgRouter(&FluentdClusterFluentdConfig1)
	g.Expect(err).NotTo(HaveOccurred())
	clusterFilters := []fluentdv1alpha1.ClusterFilter{}
	clusterOutputs := []fluentdv1alpha1.ClusterOutput{FluentdclusterOutput2CopyESDataStream}
	clustercfgResources, _ := psr.PatchAndFilterClusterLevelResources(sl, FluentdClusterFluentdConfig1.GetCfgId(), []fluentdv1alpha1.ClusterInput{}, clusterFilters, clusterOutputs)
	err = psr.IdentifyCopyAndPatchOutput(clustercfgResources)
	g.Expect(err).NotTo(HaveOccurred())
	err = psr.WithCfgResources(*clustercfgRouter.Label, clustercfgResources)
	g.Expect(err).NotTo(HaveOccurred())

	for range maxRuntimes {
		config, errs := psr.RenderMainConfig(false)
		g.Expect(errs).NotTo(HaveOccurred())
		g.Expect(string(getExpectedCfg("./expected/fluentd-cluster-cfg-output-copy-es-data-stream.cfg"))).To(Equal(config))
	}
}

func Test_Cfg2OpenSearch(t *testing.T) {
	sl := plugins.NewSecretLoader(nil, Fluentd.Namespace, logr.Logger{})
	testNamespacedConfig(t, sl, Fluentd, &FluentdConfig1, []fluentdv1alpha1.ClusterOutput{FluentdclusterOutput2OpenSearch}, "./expected/fluentd-namespaced-cfg-output-opensearch.cfg")
}

func Test_ClusterCfgOutput2OpenSearch(t *testing.T) {
	g := NewGomegaWithT(t)
	sl := plugins.NewSecretLoader(nil, Fluentd.Namespace, logr.Logger{})

	psr := fluentdv1alpha1.NewGlobalPluginResources("main")
	psr.CombineGlobalInputsPlugins(sl, Fluentd.Spec.GlobalInputs)

	clustercfgRouter, err := psr.BuildCfgRouter(&FluentdClusterFluentdConfig1)
	g.Expect(err).NotTo(HaveOccurred())
	clusterFilters := []fluentdv1alpha1.ClusterFilter{}
	clusterOutputs := []fluentdv1alpha1.ClusterOutput{FluentdclusterOutput2OpenSearch}
	clustercfgResources, _ := psr.PatchAndFilterClusterLevelResources(sl, FluentdClusterFluentdConfig1.GetCfgId(), []fluentdv1alpha1.ClusterInput{}, clusterFilters, clusterOutputs)
	err = psr.WithCfgResources(*clustercfgRouter.Label, clustercfgResources)
	g.Expect(err).NotTo(HaveOccurred())

	for range maxRuntimes {
		config, errs := psr.RenderMainConfig(false)
		// fmt.Println(config)
		g.Expect(errs).NotTo(HaveOccurred())
		g.Expect(string(getExpectedCfg("./expected/fluentd-cluster-cfg-output-opensearch.cfg"))).To(Equal(config))
	}
}

func Test_ClusterCfgOutput2Kafka(t *testing.T) {
	sl := plugins.NewSecretLoader(nil, Fluentd.Namespace, logr.Logger{})
	testClusterConfigWithFiltersAndOutputs(t, sl, Fluentd, &FluentdClusterFluentdConfig1, []fluentdv1alpha1.ClusterFilter{FluentdClusterFilter1}, []fluentdv1alpha1.ClusterOutput{FluentdClusterOutput2kafka}, "./expected/fluentd-cluster-cfg-output-kafka.cfg", false)
}

func Test_ClusterCfgOutput2Loki(t *testing.T) {
	g := NewGomegaWithT(t)
	sl := NewSecretLoader(logr.Logger{}, lokiHttpCredentials, lokiTenantName)

	psr := fluentdv1alpha1.NewGlobalPluginResources("main")
	psr.CombineGlobalInputsPlugins(sl, Fluentd.Spec.GlobalInputs)

	clustercfgRouter, err := psr.BuildCfgRouter(&FluentdClusterFluentdConfig1)
	g.Expect(err).NotTo(HaveOccurred())
	clusterFilters := []fluentdv1alpha1.ClusterFilter{FluentdClusterFilter1}
	clusterOutputs := []fluentdv1alpha1.ClusterOutput{FluentdClusterOutput2Loki}
	clustercfgResources, _ := psr.PatchAndFilterClusterLevelResources(sl, FluentdClusterFluentdConfig1.GetCfgId(), []fluentdv1alpha1.ClusterInput{}, clusterFilters, clusterOutputs)
	err = psr.WithCfgResources(*clustercfgRouter.Label, clustercfgResources)
	g.Expect(err).NotTo(HaveOccurred())

	for i := 0; i < maxRuntimes; i++ {
		config, errs := psr.RenderMainConfig(false)
		// fmt.Println(config)
		g.Expect(errs).NotTo(HaveOccurred())
		g.Expect(string(getExpectedCfg("./expected/fluentd-cluster-cfg-output-loki.cfg"))).To(Equal(config))
	}
}

func Test_MixedCfgs2ES(t *testing.T) {
	sl := plugins.NewSecretLoader(nil, Fluentd.Namespace, logr.Logger{})
	testMixedConfigs(t, sl, Fluentd, &FluentdClusterFluentdConfig1, &FluentdConfig1, []fluentdv1alpha1.ClusterOutput{FluentdclusterOutput2ES}, "./expected/fluentd-mixed-cfgs-output-es.cfg")
}

func Test_ClusterCfgOutput2CloudWatch(t *testing.T) {
	sl := plugins.NewSecretLoader(nil, Fluentd.Namespace, logr.Logger{})
	testClusterConfigWithFiltersAndOutputs(t, sl, Fluentd, &FluentdClusterFluentdConfig1, []fluentdv1alpha1.ClusterFilter{FluentdClusterFilter1}, []fluentdv1alpha1.ClusterOutput{FluentdClusterOutput2CloudWatch}, "./expected/fluentd-cluster-cfg-output-cloudwatch.cfg", true)
}

func Test_ClusterCfgOutput2Datadog(t *testing.T) {
	sl := plugins.NewSecretLoader(nil, Fluentd.Namespace, logr.Logger{})
	testClusterConfigWithFiltersAndOutputs(t, sl, Fluentd, &FluentdClusterFluentdConfig1, []fluentdv1alpha1.ClusterFilter{FluentdClusterFilter1}, []fluentdv1alpha1.ClusterOutput{FluentdClusterOutput2Datadog}, "./expected/fluentd-cluster-cfg-output-datadog.cfg", true)
}

func Test_ClusterCfgOutput2Null(t *testing.T) {
	sl := plugins.NewSecretLoader(nil, Fluentd.Namespace, logr.Logger{})
	testClusterConfigWithFiltersAndOutputs(t, sl, Fluentd, &FluentdClusterFluentdConfig1, []fluentdv1alpha1.ClusterFilter{FluentdClusterFilter1}, []fluentdv1alpha1.ClusterOutput{FluentdClusterOutput2Null}, "./expected/fluentd-cluster-cfg-output-null.cfg", true)
}

func Test_MixedCfgCopy1(t *testing.T) {
	g := NewGomegaWithT(t)
	sl := NewSecretLoader(logr.Logger{}, lokiHttpCredentials, lokiTenantName)

	psr := fluentdv1alpha1.NewGlobalPluginResources("main")
	psr.CombineGlobalInputsPlugins(sl, Fluentd.Spec.GlobalInputs)

	cfgRouter, err := psr.BuildCfgRouter(&FluentdConfig2)
	g.Expect(err).NotTo(HaveOccurred())
	outputsForCluster := []fluentdv1alpha1.Output{FluentdOutputMixedCopy1}
	clusterOutputsForCluster := []fluentdv1alpha1.ClusterOutput{FluentdClusterOutput2Loki}
	clustercfgResources, _ := psr.PatchAndFilterClusterLevelResources(sl, FluentdConfig2.GetCfgId(), []fluentdv1alpha1.ClusterInput{}, []fluentdv1alpha1.ClusterFilter{}, clusterOutputsForCluster)
	cfgResources, _ := psr.PatchAndFilterNamespacedLevelResources(sl, FluentdConfig2.GetCfgId(), []fluentdv1alpha1.Input{}, []fluentdv1alpha1.Filter{}, outputsForCluster)
	cfgResources.InputPlugins = append(cfgResources.InputPlugins, clustercfgResources.InputPlugins...)
	cfgResources.FilterPlugins = append(cfgResources.FilterPlugins, clustercfgResources.FilterPlugins...)
	cfgResources.OutputPlugins = append(cfgResources.OutputPlugins, clustercfgResources.OutputPlugins...)
	err = psr.IdentifyCopyAndPatchOutput(cfgResources)
	g.Expect(err).NotTo(HaveOccurred())

	err = psr.WithCfgResources(*cfgRouter.Label, cfgResources)
	g.Expect(err).NotTo(HaveOccurred())

	for i := 0; i < maxRuntimes; i++ {
		config, errs := psr.RenderMainConfig(false)
		// fmt.Println(config)
		g.Expect(errs).NotTo(HaveOccurred())
		g.Expect(string(getExpectedCfg("./expected/fluentd-mixed-cfgs-output-copy-1.cfg"))).To(Equal(config))
	}
}

func Test_MixedCfgCopy2(t *testing.T) {
	sl := plugins.NewSecretLoader(nil, Fluentd.Namespace, logr.Logger{})
	testMixedConfigWithCopy(t, sl, Fluentd, &FluentdConfig2, []fluentdv1alpha1.Output{FluentdOutputMixedCopy2}, []fluentdv1alpha1.ClusterOutput{}, "./expected/fluentd-mixed-cfgs-output-copy-2.cfg")
}

func Test_MixedCfgCopy3(t *testing.T) {
	sl := plugins.NewSecretLoader(nil, Fluentd.Namespace, logr.Logger{})
	testMixedConfigWithCopy(t, sl, Fluentd, &FluentdConfig2, []fluentdv1alpha1.Output{FluentdOutputMixedCopy3}, []fluentdv1alpha1.ClusterOutput{}, "./expected/fluentd-mixed-cfgs-output-copy-3.cfg")
}

func Test_MixedCfgCopy4(t *testing.T) {
	g := NewGomegaWithT(t)
	sl := NewSecretLoader(logr.Logger{}, esCredentials)

	psr := fluentdv1alpha1.NewGlobalPluginResources("main")
	psr.CombineGlobalInputsPlugins(sl, Fluentd.Spec.GlobalInputs)

	cfgRouter, err := psr.BuildCfgRouter(&FluentdConfig2)
	g.Expect(err).NotTo(HaveOccurred())
	filtersForCluster := []fluentdv1alpha1.Filter{FluentdFilter}
	clusterOutputs := []fluentdv1alpha1.ClusterOutput{FluentdClusterOutput2Loki1}
	outputsForCluster := []fluentdv1alpha1.Output{FluentdOutput2ES1, FluentdOutput2ES2, FluentdOutput2ES3, FluentdOutput2ES4}
	clustercfgResources, _ := psr.PatchAndFilterClusterLevelResources(sl, FluentdConfig2.GetCfgId(), []fluentdv1alpha1.ClusterInput{}, []fluentdv1alpha1.ClusterFilter{}, clusterOutputs)
	cfgResources, _ := psr.PatchAndFilterNamespacedLevelResources(sl, FluentdConfig2.GetCfgId(), []fluentdv1alpha1.Input{}, filtersForCluster, outputsForCluster)

	cfgResources.InputPlugins = append(cfgResources.InputPlugins, clustercfgResources.InputPlugins...)
	cfgResources.FilterPlugins = append(cfgResources.FilterPlugins, clustercfgResources.FilterPlugins...)
	cfgResources.OutputPlugins = append(cfgResources.OutputPlugins, clustercfgResources.OutputPlugins...)
	err = psr.IdentifyCopyAndPatchOutput(cfgResources)
	g.Expect(err).NotTo(HaveOccurred())

	err = psr.WithCfgResources(*cfgRouter.Label, cfgResources)
	g.Expect(err).NotTo(HaveOccurred())

	for i := 0; i < maxRuntimes; i++ {
		config, errs := psr.RenderMainConfig(false)
		g.Expect(errs).NotTo(HaveOccurred())
		g.Expect(string(getExpectedCfg("./expected/fluentd-mixed-cfgs-output-copy-4.cfg"))).To(Equal(config))
	}
}

func Test_ClusterCfgOutput2StdoutAndLoki(t *testing.T) {
	g := NewGomegaWithT(t)
	sl := plugins.NewSecretLoader(nil, Fluentd.Namespace, logr.Logger{})

	psr := fluentdv1alpha1.NewGlobalPluginResources("main")
	psr.CombineGlobalInputsPlugins(sl, Fluentd.Spec.GlobalInputs)

	clustercfgRouter, err := psr.BuildCfgRouter(&FluentdClusterFluentdConfig1)
	g.Expect(err).NotTo(HaveOccurred())
	clusterFilters := []fluentdv1alpha1.ClusterFilter{FluentdClusterFilter1}
	clusterOutputs := []fluentdv1alpha1.ClusterOutput{FluentdClusterOutputCopy2StdoutAndLoki}
	clustercfgResources, _ := psr.PatchAndFilterClusterLevelResources(sl, FluentdClusterFluentdConfig1.GetCfgId(), []fluentdv1alpha1.ClusterInput{}, clusterFilters, clusterOutputs)
	err = psr.IdentifyCopyAndPatchOutput(clustercfgResources)
	g.Expect(err).NotTo(HaveOccurred())

	err = psr.WithCfgResources(*clustercfgRouter.Label, clustercfgResources)
	g.Expect(err).NotTo(HaveOccurred())

	for i := 0; i < maxRuntimes; i++ {
		config, errs := psr.RenderMainConfig(false)
		// fmt.Println(config)
		g.Expect(errs).NotTo(HaveOccurred())
		g.Expect(string(getExpectedCfg("./expected/fluentd-cluster-cfg-output-stdout-and-loki.cfg"))).To(Equal(config))
	}
}

func Test_MixedCfgs2OpenSearch(t *testing.T) {
	sl := plugins.NewSecretLoader(nil, Fluentd.Namespace, logr.Logger{})
	testMixedConfigs(t, sl, Fluentd, &FluentdClusterFluentdConfig1, &FluentdConfig1, []fluentdv1alpha1.ClusterOutput{FluentdclusterOutput2OpenSearch}, "./expected/fluentd-mixed-cfgs-output-opensearch.cfg")
}

func Test_MixedCfgs2MultiTenant(t *testing.T) {
	g := NewGomegaWithT(t)
	sl := plugins.NewSecretLoader(nil, Fluentd.Namespace, logr.Logger{})

	psr := fluentdv1alpha1.NewGlobalPluginResources("main")
	psr.CombineGlobalInputsPlugins(sl, Fluentd.Spec.GlobalInputs)

	clustercfgRouter, err := psr.BuildCfgRouter(&FluentdClusterFluentdConfig2)
	g.Expect(err).NotTo(HaveOccurred())
	clusterOutputsForCluster := []fluentdv1alpha1.ClusterOutput{FluentdClusterOutputCluster}
	clustercfgResources, _ := psr.PatchAndFilterClusterLevelResources(sl, FluentdClusterFluentdConfig2.GetCfgId(), []fluentdv1alpha1.ClusterInput{}, []fluentdv1alpha1.ClusterFilter{}, clusterOutputsForCluster)
	err = psr.WithCfgResources(*clustercfgRouter.Label, clustercfgResources)
	g.Expect(err).NotTo(HaveOccurred())

	cfgRouter, err := psr.BuildCfgRouter(&FluentdConfigUser1)
	g.Expect(err).NotTo(HaveOccurred())
	clusterOutputsForUser1 := []fluentdv1alpha1.ClusterOutput{FluentdClusterOutputLogOperator}
	outputsForUser1 := []fluentdv1alpha1.Output{FluentdOutputUser1}
	clustercfgResourcesForUser1, _ := psr.PatchAndFilterClusterLevelResources(sl, FluentdConfigUser1.GetCfgId(), []fluentdv1alpha1.ClusterInput{}, []fluentdv1alpha1.ClusterFilter{}, clusterOutputsForUser1)
	cfgResourcesForUser1, _ := psr.PatchAndFilterNamespacedLevelResources(sl, FluentdConfigUser1.GetCfgId(), []fluentdv1alpha1.Input{}, []fluentdv1alpha1.Filter{}, outputsForUser1)
	cfgResourcesForUser1.FilterPlugins = append(cfgResourcesForUser1.FilterPlugins, clustercfgResourcesForUser1.FilterPlugins...)
	cfgResourcesForUser1.OutputPlugins = append(cfgResourcesForUser1.OutputPlugins, clustercfgResourcesForUser1.OutputPlugins...)
	err = psr.WithCfgResources(*cfgRouter.Label, cfgResourcesForUser1)
	g.Expect(err).NotTo(HaveOccurred())

	for i := 0; i < maxRuntimes; i++ {
		config, errs := psr.RenderMainConfig(false)
		// fmt.Println(config)
		g.Expect(errs).NotTo(HaveOccurred())
		g.Expect(string(getExpectedCfg("./expected/fluentd-mixed-cfgs-multi-tenant-output.cfg"))).To(Equal(config))
	}
}

func Test_OutputWithBuffer(t *testing.T) {
	sl := plugins.NewSecretLoader(nil, Fluentd.Namespace, logr.Logger{})
	testClusterConfigWithFiltersAndOutputs(t, sl, Fluentd, &FluentdClusterFluentdConfig1, []fluentdv1alpha1.ClusterFilter{FluentdClusterFilter1}, []fluentdv1alpha1.ClusterOutput{FluentdClusterOutputBuffer}, "./expected/fluentd-cluster-cfg-output-buffer-example.cfg", false)
}

func Test_OutputWithMemoryBuffer(t *testing.T) {
	sl := plugins.NewSecretLoader(nil, Fluentd.Namespace, logr.Logger{})
	testClusterConfigWithFiltersAndOutputs(t, sl, Fluentd, &FluentdClusterFluentdConfig1, []fluentdv1alpha1.ClusterFilter{FluentdClusterFilter1}, []fluentdv1alpha1.ClusterOutput{FluentdClusterOutputMemoryBuffer}, "./expected/fluentd-cluster-cfg-output-memory-buffer.cfg", false)
}

func Test_DuplicateRemovalCRSpecs(t *testing.T) {
	g := NewGomegaWithT(t)
	sl := plugins.NewSecretLoader(nil, "testnamespace", logr.Logger{})

	labels := map[string]string{
		"label0": "lv0",
		"label1": "lv1",
		"label2": "lv2",
		"label3": "lv3",
	}

	filterspec1, filterspec2, filterspec3 := CreateFluentdFilterSpecs()
	outputspec1, outputspec2, outputspec3 := CreateFluentdOutputSpecs()

	clusterFilters := []fluentdv1alpha1.ClusterFilter{
		{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "fluentd.fluent.io/v1alpha1",
				Kind:       "ClusterFilter",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: "clusterfilter1",
			},
			Spec: fluentdv1alpha1.ClusterFilterSpec(filterspec1),
		},
		{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "fluentd.fluent.io/v1alpha1",
				Kind:       "ClusterFilter",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: "clusterfilter2",
			},
			Spec: fluentdv1alpha1.ClusterFilterSpec(filterspec2),
		},
		{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "fluentd.fluent.io/v1alpha1",
				Kind:       "ClusterFilter",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: "clusterfilter3",
			},
			Spec: fluentdv1alpha1.ClusterFilterSpec(filterspec3),
		},
	}

	clusterOutputs := []fluentdv1alpha1.ClusterOutput{
		{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "fluentd.fluent.io/v1alpha1",
				Kind:       "ClusterOutput",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: "clusteroutput1",
			},
			Spec: fluentdv1alpha1.ClusterOutputSpec(outputspec1),
		},
		{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "fluentd.fluent.io/v1alpha1",
				Kind:       "ClusterOutput",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: "clusteroutput2",
			},
			Spec: fluentdv1alpha1.ClusterOutputSpec(outputspec2),
		},
		{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "fluentd.fluent.io/v1alpha1",
				Kind:       "ClusterOutput",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: "clusteroutput3",
			},
			Spec: fluentdv1alpha1.ClusterOutputSpec(outputspec3),
		},
	}

	filters := []fluentdv1alpha1.Filter{
		{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "fluentd.fluent.io/v1alpha1",
				Kind:       "Filter",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "filter1",
				Namespace: "testnamespace",
			},
			Spec: filterspec1,
		},
		{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "fluentd.fluent.io/v1alpha1",
				Kind:       "Filter",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "filter2",
				Namespace: "testnamespace",
			},
			Spec: filterspec2,
		},
		{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "fluentd.fluent.io/v1alpha1",
				Kind:       "Filter",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "filter3",
				Namespace: "testnamespace",
			},
			Spec: filterspec3,
		},
	}

	outputs := []fluentdv1alpha1.Output{
		{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "fluentd.fluent.io/v1alpha1",
				Kind:       "Output",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "output1",
				Namespace: "testnamespace",
			},
			Spec: outputspec1,
		},
		{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "fluentd.fluent.io/v1alpha1",
				Kind:       "Output",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "output2",
				Namespace: "testnamespace",
			},
			Spec: outputspec2,
		},
		{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "fluentd.fluent.io/v1alpha1",
				Kind:       "Output",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "output3",
				Namespace: "testnamespace",
			},
			Spec: outputspec3,
		},
	}

	clustercfg := fluentdv1alpha1.ClusterFluentdConfig{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "fluentd.fluent.io/v1alpha1",
			Kind:       "ClusterFluentdConfig",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "clusterfluentdconfig0",
		},
		Spec: fluentdv1alpha1.ClusterFluentdConfigSpec{
			WatchedNamespaces: []string{"ns1", "ns2"},
			WatchedLabels:     labels,
		},
	}

	cfg := fluentdv1alpha1.FluentdConfig{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "fluentd.fluent.io/v1alpha1",
			Kind:       "FluentdConfig",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "fluentdconfig0",
			Namespace: "testnamespace",
		},
		Spec: fluentdv1alpha1.FluentdConfigSpec{
			WatchedLabels: labels,
		},
	}

	psr := fluentdv1alpha1.NewGlobalPluginResources("main")
	psr.CombineGlobalInputsPlugins(sl, GlobalInputs)

	clustercfgRouter, err := psr.BuildCfgRouter(&clustercfg)
	g.Expect(err).NotTo(HaveOccurred())
	clustercfgResources, _ := psr.PatchAndFilterClusterLevelResources(sl, clustercfg.GetCfgId(), []fluentdv1alpha1.ClusterInput{}, clusterFilters, clusterOutputs)
	err = psr.WithCfgResources(*clustercfgRouter.Label, clustercfgResources)
	g.Expect(err).NotTo(HaveOccurred())

	cfgRouter, err := psr.BuildCfgRouter(&cfg)
	g.Expect(err).NotTo(HaveOccurred())
	cfgResources, _ := psr.PatchAndFilterNamespacedLevelResources(sl, cfg.GetCfgId(), []fluentdv1alpha1.Input{}, filters, outputs)
	err = psr.WithCfgResources(*cfgRouter.Label, cfgResources)
	g.Expect(err).NotTo(HaveOccurred())

	for i := 0; i < maxRuntimes; i++ {
		config, errs := psr.RenderMainConfig(false)
		// fmt.Println(config)
		g.Expect(errs).NotTo(HaveOccurred())
		g.Expect(string(getExpectedCfg("./expected/duplicate-removal-cr-specs.cfg"))).To(Equal(config))
	}
}

func Test_RecordTransformer(t *testing.T) {
	sl := plugins.NewSecretLoader(nil, Fluentd.Namespace, logr.Logger{})
	testClusterConfigWithFiltersAndOutputs(t, sl, Fluentd, &FluentdClusterFluentdConfig1, []fluentdv1alpha1.ClusterFilter{FluentdClusterRecordTransformerFilter}, []fluentdv1alpha1.ClusterOutput{FluentdClusterOutputCluster}, "./expected/fluentd-cluster-cfg-filter-recordTransformer.cfg", false)
}

// testNamespacedConfig tests a namespaced config with cluster outputs
func testNamespacedConfig(
	t *testing.T,
	sl plugins.SecretLoader,
	fluentd fluentdv1alpha1.Fluentd,
	config *fluentdv1alpha1.FluentdConfig,
	clusterOutputs []fluentdv1alpha1.ClusterOutput,
	expectedCfgPath string,
) {
	g := NewGomegaWithT(t)

	psr := fluentdv1alpha1.NewGlobalPluginResources("main")
	psr.CombineGlobalInputsPlugins(sl, fluentd.Spec.GlobalInputs)

	cfgRouter, err := psr.BuildCfgRouter(config)
	g.Expect(err).NotTo(HaveOccurred())
	cfgResources, _ := psr.PatchAndFilterClusterLevelResources(sl, config.GetCfgId(), []fluentdv1alpha1.ClusterInput{}, []fluentdv1alpha1.ClusterFilter{}, clusterOutputs)
	err = psr.WithCfgResources(*cfgRouter.Label, cfgResources)
	g.Expect(err).NotTo(HaveOccurred())

	for i := 0; i < maxRuntimes; i++ {
		config, errs := psr.RenderMainConfig(false)
		g.Expect(errs).NotTo(HaveOccurred())
		g.Expect(string(getExpectedCfg(expectedCfgPath))).To(Equal(config))
	}
}

// testClusterConfigWithGlobalInputs tests a cluster config with custom global inputs
func testClusterConfigWithGlobalInputs(
	t *testing.T,
	sl plugins.SecretLoader,
	fluentd fluentdv1alpha1.Fluentd,
	config *fluentdv1alpha1.FluentdConfig,
	clusterOutputs []fluentdv1alpha1.ClusterOutput,
	expectedCfgPath string,
) {
	g := NewGomegaWithT(t)

	psr := fluentdv1alpha1.NewGlobalPluginResources("main")
	psr.CombineGlobalInputsPlugins(sl, fluentd.Spec.GlobalInputs)

	clustercfgRouter, err := psr.BuildCfgRouter(config)
	g.Expect(err).NotTo(HaveOccurred())
	clustercfgResources, _ := psr.PatchAndFilterClusterLevelResources(sl, config.GetCfgId(), []fluentdv1alpha1.ClusterInput{}, []fluentdv1alpha1.ClusterFilter{}, clusterOutputs)
	err = psr.WithCfgResources(*clustercfgRouter.Label, clustercfgResources)
	g.Expect(err).NotTo(HaveOccurred())

	for i := 0; i < maxRuntimes; i++ {
		config, errs := psr.RenderMainConfig(false)
		g.Expect(errs).NotTo(HaveOccurred())
		g.Expect(string(getExpectedCfg(expectedCfgPath))).To(Equal(config))
	}
}

// testClusterConfigWithFiltersAndOutputs tests a cluster config with filters and outputs
func testClusterConfigWithFiltersAndOutputs(
	t *testing.T,
	sl plugins.SecretLoader,
	fluentd fluentdv1alpha1.Fluentd,
	clusterConfig *fluentdv1alpha1.ClusterFluentdConfig,
	clusterFilters []fluentdv1alpha1.ClusterFilter,
	clusterOutputs []fluentdv1alpha1.ClusterOutput,
	expectedCfgPath string,
	useTrimSpace bool,
) {
	g := NewGomegaWithT(t)

	psr := fluentdv1alpha1.NewGlobalPluginResources("main")
	psr.CombineGlobalInputsPlugins(sl, fluentd.Spec.GlobalInputs)

	clustercfgRouter, err := psr.BuildCfgRouter(clusterConfig)
	g.Expect(err).NotTo(HaveOccurred())
	clustercfgResources, _ := psr.PatchAndFilterClusterLevelResources(sl, clusterConfig.GetCfgId(), []fluentdv1alpha1.ClusterInput{}, clusterFilters, clusterOutputs)
	err = psr.WithCfgResources(*clustercfgRouter.Label, clustercfgResources)
	g.Expect(err).NotTo(HaveOccurred())

	for i := 0; i < maxRuntimes; i++ {
		config, errs := psr.RenderMainConfig(false)
		g.Expect(errs).NotTo(HaveOccurred())
		expectedCfg := string(getExpectedCfg(expectedCfgPath))
		if useTrimSpace {
			g.Expect(strings.TrimSpace(expectedCfg)).To(Equal(config))
		} else {
			g.Expect(expectedCfg).To(Equal(config))
		}
	}
}

// testMixedConfigs tests a combination of cluster and namespace configs
func testMixedConfigs(
	t *testing.T,
	sl plugins.SecretLoader,
	fluentd fluentdv1alpha1.Fluentd,
	clusterConfig *fluentdv1alpha1.ClusterFluentdConfig,
	namespacedConfig *fluentdv1alpha1.FluentdConfig,
	clusterOutputs []fluentdv1alpha1.ClusterOutput,
	expectedCfgPath string,
) {
	g := NewGomegaWithT(t)

	psr := fluentdv1alpha1.NewGlobalPluginResources("main")
	psr.CombineGlobalInputsPlugins(sl, fluentd.Spec.GlobalInputs)

	clustercfgRouter, err := psr.BuildCfgRouter(clusterConfig)
	g.Expect(err).NotTo(HaveOccurred())
	clustercfgResources, _ := psr.PatchAndFilterClusterLevelResources(sl, clusterConfig.GetCfgId(), []fluentdv1alpha1.ClusterInput{}, []fluentdv1alpha1.ClusterFilter{}, clusterOutputs)
	err = psr.WithCfgResources(*clustercfgRouter.Label, clustercfgResources)
	g.Expect(err).NotTo(HaveOccurred())

	cfgRouter, err := psr.BuildCfgRouter(namespacedConfig)
	g.Expect(err).NotTo(HaveOccurred())
	cfgResources, _ := psr.PatchAndFilterClusterLevelResources(sl, namespacedConfig.GetCfgId(), []fluentdv1alpha1.ClusterInput{}, []fluentdv1alpha1.ClusterFilter{}, clusterOutputs)
	err = psr.WithCfgResources(*cfgRouter.Label, cfgResources)
	g.Expect(err).NotTo(HaveOccurred())

	for i := 0; i < maxRuntimes; i++ {
		config, errs := psr.RenderMainConfig(false)
		g.Expect(errs).NotTo(HaveOccurred())
		g.Expect(string(getExpectedCfg(expectedCfgPath))).To(Equal(config))
	}
}

// testMixedConfigWithCopy tests a mixed config with copy outputs
func testMixedConfigWithCopy(
	t *testing.T,
	sl plugins.SecretLoader,
	fluentd fluentdv1alpha1.Fluentd,
	namespacedConfig *fluentdv1alpha1.FluentdConfig,
	outputs []fluentdv1alpha1.Output,
	clusterOutputs []fluentdv1alpha1.ClusterOutput,
	expectedCfgPath string,
) {
	g := NewGomegaWithT(t)

	psr := fluentdv1alpha1.NewGlobalPluginResources("main")
	psr.CombineGlobalInputsPlugins(sl, fluentd.Spec.GlobalInputs)

	cfgRouter, err := psr.BuildCfgRouter(namespacedConfig)
	g.Expect(err).NotTo(HaveOccurred())
	clustercfgResources, _ := psr.PatchAndFilterClusterLevelResources(sl, namespacedConfig.GetCfgId(), []fluentdv1alpha1.ClusterInput{}, []fluentdv1alpha1.ClusterFilter{}, clusterOutputs)
	cfgResources, _ := psr.PatchAndFilterNamespacedLevelResources(sl, namespacedConfig.GetCfgId(), []fluentdv1alpha1.Input{}, []fluentdv1alpha1.Filter{}, outputs)
	cfgResources.InputPlugins = append(cfgResources.InputPlugins, clustercfgResources.InputPlugins...)
	cfgResources.FilterPlugins = append(cfgResources.FilterPlugins, clustercfgResources.FilterPlugins...)
	cfgResources.OutputPlugins = append(cfgResources.OutputPlugins, clustercfgResources.OutputPlugins...)
	err = psr.IdentifyCopyAndPatchOutput(cfgResources)
	g.Expect(err).NotTo(HaveOccurred())

	err = psr.WithCfgResources(*cfgRouter.Label, cfgResources)
	g.Expect(err).NotTo(HaveOccurred())

	for i := 0; i < maxRuntimes; i++ {
		config, errs := psr.RenderMainConfig(false)
		g.Expect(errs).NotTo(HaveOccurred())
		g.Expect(string(getExpectedCfg(expectedCfgPath))).To(Equal(config))
	}
}

type SecretLoaderStruct struct {
	secrets map[string]corev1.Secret
}

func NewSecretLoader(l logr.Logger, sec ...corev1.Secret) plugins.SecretLoader {
	secrets := make(map[string]corev1.Secret)
	for _, s := range sec {
		secrets[s.Name] = s
	}
	return SecretLoaderStruct{
		secrets: secrets,
	}
}

func (sl SecretLoaderStruct) LoadSecret(s plugins.Secret) (string, error) {
	var secret corev1.Secret
	var ok bool
	if secret, ok = sl.secrets[s.ValueFrom.SecretKeyRef.Name]; !ok {
		return "", errors.NotFound(fmt.Sprintf("The secret %s is not found.", s.ValueFrom.SecretKeyRef.Name))
	}

	if v, ok := secret.StringData[s.ValueFrom.SecretKeyRef.Key]; !ok {
		return "", errors.NotFound(fmt.Sprintf("The key %s is not found.", s.ValueFrom.SecretKeyRef.Key))
	} else {
		return strings.TrimSuffix(v, "\n"), nil
	}
}
