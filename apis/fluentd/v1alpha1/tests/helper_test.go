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

	fluentdv1alpha1 "github.com/fluent/fluent-operator/v2/apis/fluentd/v1alpha1"
	"github.com/fluent/fluent-operator/v2/apis/fluentd/v1alpha1/plugins"
)

const (
	maxRuntimes = 5
)

func Test_Cfg2ES(t *testing.T) {
	g := NewGomegaWithT(t)
	sl := plugins.NewSecretLoader(nil, Fluentd.Namespace, logr.Logger{})

	psr := fluentdv1alpha1.NewGlobalPluginResources("main")
	psr.CombineGlobalInputsPlugins(sl, Fluentd.Spec.GlobalInputs)

	clusterOutputsForCluster := []fluentdv1alpha1.ClusterOutput{FluentdclusterOutput2ES}
	cfgRouter, err := psr.BuildCfgRouter(&FluentdConfig1)
	g.Expect(err).NotTo(HaveOccurred())
	cfgResources, _ := psr.PatchAndFilterClusterLevelResources(sl, FluentdConfig1.GetCfgId(), []fluentdv1alpha1.ClusterInput{}, []fluentdv1alpha1.ClusterFilter{}, clusterOutputsForCluster)
	err = psr.WithCfgResources(*cfgRouter.Label, cfgResources)
	g.Expect(err).NotTo(HaveOccurred())

	// we should not see any permutations in serialized config
	i := 0
	for i < maxRuntimes {
		config, errs := psr.RenderMainConfig(false)
		// fmt.Println(config)
		g.Expect(errs).NotTo(HaveOccurred())
		g.Expect(string(getExpectedCfg("./expected/fluentd-namespaced-cfg-output-es.cfg"))).To(Equal(config))

		i++
	}
}

func Test_ClusterCfgInputTail(t *testing.T) {
	g := NewGomegaWithT(t)
	sl := plugins.NewSecretLoader(nil, Fluentd.Namespace, logr.Logger{})

	psr := fluentdv1alpha1.NewGlobalPluginResources("main")
	psr.CombineGlobalInputsPlugins(sl, FluentdInputTail.Spec.GlobalInputs)

	clustercfgRouter, err := psr.BuildCfgRouter(&FluentdConfig1)
	g.Expect(err).NotTo(HaveOccurred())
	clusterOutputs := []fluentdv1alpha1.ClusterOutput{FluentdClusterOutputTag}
	clustercfgResources, _ := psr.PatchAndFilterClusterLevelResources(sl, FluentdConfig1.GetCfgId(), []fluentdv1alpha1.ClusterInput{}, []fluentdv1alpha1.ClusterFilter{}, clusterOutputs)
	err = psr.WithCfgResources(*clustercfgRouter.Label, clustercfgResources)
	g.Expect(err).NotTo(HaveOccurred())

	for i := 0; i < maxRuntimes; i++ {
		config, errs := psr.RenderMainConfig(false)
		// fmt.Println(config)
		g.Expect(errs).NotTo(HaveOccurred())
		g.Expect(string(getExpectedCfg("./expected/fluentd-global-cfg-input-tail.cfg"))).To(Equal(config))
	}

}

func Test_ClusterCfgInputSample(t *testing.T) {
	g := NewGomegaWithT(t)
	sl := plugins.NewSecretLoader(nil, Fluentd.Namespace, logr.Logger{})

	psr := fluentdv1alpha1.NewGlobalPluginResources("main")
	psr.CombineGlobalInputsPlugins(sl, FluentdInputSample.Spec.GlobalInputs)

	clustercfgRouter, err := psr.BuildCfgRouter(&FluentdConfig1)
	g.Expect(err).NotTo(HaveOccurred())
	clusterOutputs := []fluentdv1alpha1.ClusterOutput{FluentdClusterOutputTag}
	clustercfgResources, _ := psr.PatchAndFilterClusterLevelResources(sl, FluentdConfig1.GetCfgId(), []fluentdv1alpha1.ClusterInput{}, []fluentdv1alpha1.ClusterFilter{}, clusterOutputs)
	err = psr.WithCfgResources(*clustercfgRouter.Label, clustercfgResources)
	g.Expect(err).NotTo(HaveOccurred())

	for i := 0; i < maxRuntimes; i++ {
		config, errs := psr.RenderMainConfig(false)
		// fmt.Println(config)
		g.Expect(errs).NotTo(HaveOccurred())
		g.Expect(string(getExpectedCfg("./expected/fluentd-global-cfg-input-sample.cfg"))).To(Equal(config))
	}
}

func Test_ClusterCfgInputMonitorAgent(t *testing.T) {
	g := NewGomegaWithT(t)
	sl := plugins.NewSecretLoader(nil, Fluentd.Namespace, logr.Logger{})

	psr := fluentdv1alpha1.NewGlobalPluginResources("main")
	psr.CombineGlobalInputsPlugins(sl, FluentdInputMonitorAgent.Spec.GlobalInputs)

	clustercfgRouter, err := psr.BuildCfgRouter(&FluentdConfig1)
	g.Expect(err).NotTo(HaveOccurred())
	clusterOutputs := []fluentdv1alpha1.ClusterOutput{FluentdClusterOutputTag}
	clustercfgResources, _ := psr.PatchAndFilterClusterLevelResources(sl, FluentdConfig1.GetCfgId(), []fluentdv1alpha1.ClusterInput{}, []fluentdv1alpha1.ClusterFilter{}, clusterOutputs)
	err = psr.WithCfgResources(*clustercfgRouter.Label, clustercfgResources)
	g.Expect(err).NotTo(HaveOccurred())

	for i := 0; i < maxRuntimes; i++ {
		config, errs := psr.RenderMainConfig(false)
		g.Expect(errs).NotTo(HaveOccurred())
		g.Expect(string(getExpectedCfg("./expected/fluentd-global-cfg-input-monitorAgent.cfg"))).To(Equal(config))
	}
}

func Test_ClusterCfgOutput2ES(t *testing.T) {
	g := NewGomegaWithT(t)
	sl := plugins.NewSecretLoader(nil, Fluentd.Namespace, logr.Logger{})

	psr := fluentdv1alpha1.NewGlobalPluginResources("main")
	psr.CombineGlobalInputsPlugins(sl, Fluentd.Spec.GlobalInputs)

	clustercfgRouter, err := psr.BuildCfgRouter(&FluentdClusterFluentdConfig1)
	g.Expect(err).NotTo(HaveOccurred())
	clusterFilters := []fluentdv1alpha1.ClusterFilter{}
	clusterOutputs := []fluentdv1alpha1.ClusterOutput{FluentdclusterOutput2ES}
	clustercfgResources, _ := psr.PatchAndFilterClusterLevelResources(sl, FluentdClusterFluentdConfig1.GetCfgId(), []fluentdv1alpha1.ClusterInput{}, clusterFilters, clusterOutputs)
	err = psr.WithCfgResources(*clustercfgRouter.Label, clustercfgResources)
	g.Expect(err).NotTo(HaveOccurred())

	// we should not see any permutations in serialized config
	i := 0
	for i < maxRuntimes {
		config, errs := psr.RenderMainConfig(false)
		fmt.Println(config)
		g.Expect(errs).NotTo(HaveOccurred())
		g.Expect(string(getExpectedCfg("./expected/fluentd-cluster-cfg-output-es.cfg"))).To(Equal(config))

		i++
	}
}

func Test_Cfg2OpenSearch(t *testing.T) {
	g := NewGomegaWithT(t)
	sl := plugins.NewSecretLoader(nil, Fluentd.Namespace, logr.Logger{})

	psr := fluentdv1alpha1.NewGlobalPluginResources("main")
	psr.CombineGlobalInputsPlugins(sl, Fluentd.Spec.GlobalInputs)

	clusterOutputsForCluster := []fluentdv1alpha1.ClusterOutput{FluentdclusterOutput2OpenSearch}
	cfgRouter, err := psr.BuildCfgRouter(&FluentdConfig1)
	g.Expect(err).NotTo(HaveOccurred())
	cfgResources, _ := psr.PatchAndFilterClusterLevelResources(sl, FluentdConfig1.GetCfgId(), []fluentdv1alpha1.ClusterInput{}, []fluentdv1alpha1.ClusterFilter{}, clusterOutputsForCluster)
	err = psr.WithCfgResources(*cfgRouter.Label, cfgResources)
	g.Expect(err).NotTo(HaveOccurred())

	// we should not see any permutations in serialized config
	i := 0
	for i < maxRuntimes {
		config, errs := psr.RenderMainConfig(false)
		g.Expect(errs).NotTo(HaveOccurred())
		g.Expect(string(getExpectedCfg("./expected/fluentd-namespaced-cfg-output-opensearch.cfg"))).To(Equal(config))

		i++
	}
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

	// we should not see any permutations in serialized config
	i := 0
	for i < maxRuntimes {
		config, errs := psr.RenderMainConfig(false)
		//fmt.Println(config)
		g.Expect(errs).NotTo(HaveOccurred())
		g.Expect(string(getExpectedCfg("./expected/fluentd-cluster-cfg-output-opensearch.cfg"))).To(Equal(config))

		i++
	}
}

func Test_ClusterCfgOutput2Kafka(t *testing.T) {
	g := NewGomegaWithT(t)
	sl := plugins.NewSecretLoader(nil, Fluentd.Namespace, logr.Logger{})

	psr := fluentdv1alpha1.NewGlobalPluginResources("main")
	psr.CombineGlobalInputsPlugins(sl, Fluentd.Spec.GlobalInputs)

	clustercfgRouter, err := psr.BuildCfgRouter(&FluentdClusterFluentdConfig1)
	g.Expect(err).NotTo(HaveOccurred())
	clusterFilters := []fluentdv1alpha1.ClusterFilter{FluentdClusterFilter1}
	clusterOutputs := []fluentdv1alpha1.ClusterOutput{FluentdClusterOutput2kafka}
	clustercfgResources, _ := psr.PatchAndFilterClusterLevelResources(sl, FluentdClusterFluentdConfig1.GetCfgId(), []fluentdv1alpha1.ClusterInput{}, clusterFilters, clusterOutputs)
	err = psr.WithCfgResources(*clustercfgRouter.Label, clustercfgResources)
	g.Expect(err).NotTo(HaveOccurred())

	// we should not see any permutations in serialized config
	i := 0
	for i < maxRuntimes {
		config, errs := psr.RenderMainConfig(false)
		// fmt.Println(config)
		g.Expect(errs).NotTo(HaveOccurred())
		g.Expect(string(getExpectedCfg("./expected/fluentd-cluster-cfg-output-kafka.cfg"))).To(Equal(config))

		i++
	}
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

	// we should not see any permutations in serialized config
	i := 0
	for i < maxRuntimes {
		config, errs := psr.RenderMainConfig(false)
		// fmt.Println(config)
		g.Expect(errs).NotTo(HaveOccurred())
		g.Expect(string(getExpectedCfg("./expected/fluentd-cluster-cfg-output-loki.cfg"))).To(Equal(config))

		i++
	}
}

func Test_MixedCfgs2ES(t *testing.T) {
	g := NewGomegaWithT(t)
	sl := plugins.NewSecretLoader(nil, Fluentd.Namespace, logr.Logger{})

	psr := fluentdv1alpha1.NewGlobalPluginResources("main")
	psr.CombineGlobalInputsPlugins(sl, Fluentd.Spec.GlobalInputs)

	clustercfgRouter, err := psr.BuildCfgRouter(&FluentdClusterFluentdConfig1)
	g.Expect(err).NotTo(HaveOccurred())
	clusterOutputsForCluster := []fluentdv1alpha1.ClusterOutput{FluentdclusterOutput2ES}
	clustercfgResources, _ := psr.PatchAndFilterClusterLevelResources(sl, FluentdClusterFluentdConfig1.GetCfgId(), []fluentdv1alpha1.ClusterInput{}, []fluentdv1alpha1.ClusterFilter{}, clusterOutputsForCluster)
	err = psr.WithCfgResources(*clustercfgRouter.Label, clustercfgResources)
	g.Expect(err).NotTo(HaveOccurred())

	cfgRouter, err := psr.BuildCfgRouter(&FluentdConfig1)
	g.Expect(err).NotTo(HaveOccurred())
	cfgResources, _ := psr.PatchAndFilterClusterLevelResources(sl, FluentdConfig1.GetCfgId(), []fluentdv1alpha1.ClusterInput{}, []fluentdv1alpha1.ClusterFilter{}, clusterOutputsForCluster)
	err = psr.WithCfgResources(*cfgRouter.Label, cfgResources)
	g.Expect(err).NotTo(HaveOccurred())

	// we should not see any permutations in serialized config
	i := 0
	for i < maxRuntimes {
		config, errs := psr.RenderMainConfig(false)
		// fmt.Println(config)
		g.Expect(errs).NotTo(HaveOccurred())
		g.Expect(string(getExpectedCfg("./expected/fluentd-mixed-cfgs-output-es.cfg"))).To(Equal(config))

		i++
	}
}

func Test_ClusterCfgOutput2CloudWatch(t *testing.T) {
	g := NewGomegaWithT(t)
	sl := plugins.NewSecretLoader(nil, Fluentd.Namespace, logr.Logger{})

	psr := fluentdv1alpha1.NewGlobalPluginResources("main")
	psr.CombineGlobalInputsPlugins(sl, Fluentd.Spec.GlobalInputs)

	clustercfgRouter, err := psr.BuildCfgRouter(&FluentdClusterFluentdConfig1)
	g.Expect(err).NotTo(HaveOccurred())
	clusterFilters := []fluentdv1alpha1.ClusterFilter{FluentdClusterFilter1}
	clusterOutputs := []fluentdv1alpha1.ClusterOutput{FluentdClusterOutput2CloudWatch}
	clustercfgResources, _ := psr.PatchAndFilterClusterLevelResources(sl, FluentdClusterFluentdConfig1.GetCfgId(), []fluentdv1alpha1.ClusterInput{}, clusterFilters, clusterOutputs)
	err = psr.WithCfgResources(*clustercfgRouter.Label, clustercfgResources)
	g.Expect(err).NotTo(HaveOccurred())

	// we should not see any permutations in serialized config
	i := 0
	for i < maxRuntimes {
		config, errs := psr.RenderMainConfig(false)
		g.Expect(errs).NotTo(HaveOccurred())
		g.Expect(strings.TrimSpace(string(getExpectedCfg("./expected/fluentd-cluster-cfg-output-cloudwatch.cfg")))).To(Equal(config))

		i++
	}
}

func Test_ClusterCfgOutput2Datadog(t *testing.T) {
	g := NewGomegaWithT(t)
	sl := plugins.NewSecretLoader(nil, Fluentd.Namespace, logr.Logger{})

	psr := fluentdv1alpha1.NewGlobalPluginResources("main")
	psr.CombineGlobalInputsPlugins(sl, Fluentd.Spec.GlobalInputs)

	clustercfgRouter, err := psr.BuildCfgRouter(&FluentdClusterFluentdConfig1)
	g.Expect(err).NotTo(HaveOccurred())
	clusterFilters := []fluentdv1alpha1.ClusterFilter{FluentdClusterFilter1}
	clusterOutputs := []fluentdv1alpha1.ClusterOutput{FluentdClusterOutput2Datadog}
	clustercfgResources, _ := psr.PatchAndFilterClusterLevelResources(sl, FluentdClusterFluentdConfig1.GetCfgId(), []fluentdv1alpha1.ClusterInput{}, clusterFilters, clusterOutputs)
	err = psr.WithCfgResources(*clustercfgRouter.Label, clustercfgResources)
	g.Expect(err).NotTo(HaveOccurred())

	// we should not see any permutations in serialized config
	i := 0
	for i < maxRuntimes {
		config, errs := psr.RenderMainConfig(false)
		g.Expect(errs).NotTo(HaveOccurred())
		g.Expect(strings.TrimSpace(string(getExpectedCfg("./expected/fluentd-cluster-cfg-output-datadog.cfg")))).To(Equal(config))

		i++
	}
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

	// we should not see any permutations in serialized config
	i := 0
	for i < maxRuntimes {
		config, errs := psr.RenderMainConfig(false)
		// fmt.Println(config)
		g.Expect(errs).NotTo(HaveOccurred())
		g.Expect(string(getExpectedCfg("./expected/fluentd-mixed-cfgs-output-copy-1.cfg"))).To(Equal(config))

		i++
	}
}

func Test_MixedCfgCopy2(t *testing.T) {
	g := NewGomegaWithT(t)
	sl := plugins.NewSecretLoader(nil, Fluentd.Namespace, logr.Logger{})

	psr := fluentdv1alpha1.NewGlobalPluginResources("main")
	psr.CombineGlobalInputsPlugins(sl, Fluentd.Spec.GlobalInputs)

	cfgRouter, err := psr.BuildCfgRouter(&FluentdConfig2)
	g.Expect(err).NotTo(HaveOccurred())
	outputsForCluster := []fluentdv1alpha1.Output{FluentdOutputMixedCopy2}
	clustercfgResources, _ := psr.PatchAndFilterClusterLevelResources(sl, FluentdConfig2.GetCfgId(), []fluentdv1alpha1.ClusterInput{}, []fluentdv1alpha1.ClusterFilter{}, []fluentdv1alpha1.ClusterOutput{})
	cfgResources, _ := psr.PatchAndFilterNamespacedLevelResources(sl, FluentdConfig2.GetCfgId(), []fluentdv1alpha1.Input{}, []fluentdv1alpha1.Filter{}, outputsForCluster)
	cfgResources.InputPlugins = append(cfgResources.InputPlugins, clustercfgResources.InputPlugins...)
	cfgResources.FilterPlugins = append(cfgResources.FilterPlugins, clustercfgResources.FilterPlugins...)
	cfgResources.OutputPlugins = append(cfgResources.OutputPlugins, clustercfgResources.OutputPlugins...)
	err = psr.IdentifyCopyAndPatchOutput(cfgResources)
	g.Expect(err).NotTo(HaveOccurred())

	err = psr.WithCfgResources(*cfgRouter.Label, cfgResources)
	g.Expect(err).NotTo(HaveOccurred())

	// we should not see any permutations in serialized config
	i := 0
	for i < maxRuntimes {
		config, errs := psr.RenderMainConfig(false)
		// fmt.Println(config)
		g.Expect(errs).NotTo(HaveOccurred())
		g.Expect(string(getExpectedCfg("./expected/fluentd-mixed-cfgs-output-copy-2.cfg"))).To(Equal(config))

		i++
	}
}

func Test_MixedCfgCopy3(t *testing.T) {
	g := NewGomegaWithT(t)
	sl := plugins.NewSecretLoader(nil, Fluentd.Namespace, logr.Logger{})

	psr := fluentdv1alpha1.NewGlobalPluginResources("main")
	psr.CombineGlobalInputsPlugins(sl, Fluentd.Spec.GlobalInputs)

	cfgRouter, err := psr.BuildCfgRouter(&FluentdConfig2)
	g.Expect(err).NotTo(HaveOccurred())
	outputsForCluster := []fluentdv1alpha1.Output{FluentdOutputMixedCopy3}
	clustercfgResources, _ := psr.PatchAndFilterClusterLevelResources(sl, FluentdConfig2.GetCfgId(), []fluentdv1alpha1.ClusterInput{}, []fluentdv1alpha1.ClusterFilter{}, []fluentdv1alpha1.ClusterOutput{})
	cfgResources, _ := psr.PatchAndFilterNamespacedLevelResources(sl, FluentdConfig2.GetCfgId(), []fluentdv1alpha1.Input{}, []fluentdv1alpha1.Filter{}, outputsForCluster)
	cfgResources.InputPlugins = append(cfgResources.InputPlugins, clustercfgResources.InputPlugins...)
	cfgResources.FilterPlugins = append(cfgResources.FilterPlugins, clustercfgResources.FilterPlugins...)
	cfgResources.OutputPlugins = append(cfgResources.OutputPlugins, clustercfgResources.OutputPlugins...)
	err = psr.IdentifyCopyAndPatchOutput(cfgResources)
	g.Expect(err).NotTo(HaveOccurred())

	err = psr.WithCfgResources(*cfgRouter.Label, cfgResources)
	g.Expect(err).NotTo(HaveOccurred())

	// we should not see any permutations in serialized config
	i := 0
	for i < maxRuntimes {
		config, errs := psr.RenderMainConfig(false)
		// fmt.Println(config)
		g.Expect(errs).NotTo(HaveOccurred())
		g.Expect(string(getExpectedCfg("./expected/fluentd-mixed-cfgs-output-copy-3.cfg"))).To(Equal(config))

		i++
	}
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

	// we should not see any permutations in serialized config
	i := 0
	for i < maxRuntimes {
		config, errs := psr.RenderMainConfig(false)
		g.Expect(errs).NotTo(HaveOccurred())
		g.Expect(string(getExpectedCfg("./expected/fluentd-mixed-cfgs-output-copy-4.cfg"))).To(Equal(config))

		i++
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

	// we should not see any permutations in serialized config
	i := 0
	for i < maxRuntimes {
		config, errs := psr.RenderMainConfig(false)
		// fmt.Println(config)
		g.Expect(errs).NotTo(HaveOccurred())
		g.Expect(string(getExpectedCfg("./expected/fluentd-cluster-cfg-output-stdout-and-loki.cfg"))).To(Equal(config))

		i++
	}
}

func Test_MixedCfgs2OpenSearch(t *testing.T) {
	g := NewGomegaWithT(t)
	sl := plugins.NewSecretLoader(nil, Fluentd.Namespace, logr.Logger{})

	psr := fluentdv1alpha1.NewGlobalPluginResources("main")
	psr.CombineGlobalInputsPlugins(sl, Fluentd.Spec.GlobalInputs)

	clustercfgRouter, err := psr.BuildCfgRouter(&FluentdClusterFluentdConfig1)
	g.Expect(err).NotTo(HaveOccurred())
	clusterOutputsForCluster := []fluentdv1alpha1.ClusterOutput{FluentdclusterOutput2OpenSearch}
	clustercfgResources, _ := psr.PatchAndFilterClusterLevelResources(sl, FluentdClusterFluentdConfig1.GetCfgId(), []fluentdv1alpha1.ClusterInput{}, []fluentdv1alpha1.ClusterFilter{}, clusterOutputsForCluster)
	err = psr.WithCfgResources(*clustercfgRouter.Label, clustercfgResources)
	g.Expect(err).NotTo(HaveOccurred())

	cfgRouter, err := psr.BuildCfgRouter(&FluentdConfig1)
	g.Expect(err).NotTo(HaveOccurred())
	cfgResources, _ := psr.PatchAndFilterClusterLevelResources(sl, FluentdConfig1.GetCfgId(), []fluentdv1alpha1.ClusterInput{}, []fluentdv1alpha1.ClusterFilter{}, clusterOutputsForCluster)
	err = psr.WithCfgResources(*cfgRouter.Label, cfgResources)
	g.Expect(err).NotTo(HaveOccurred())

	// we should not see any permutations in serialized config
	i := 0
	for i < maxRuntimes {
		config, errs := psr.RenderMainConfig(false)
		// fmt.Println(config)
		g.Expect(errs).NotTo(HaveOccurred())
		g.Expect(string(getExpectedCfg("./expected/fluentd-mixed-cfgs-output-opensearch.cfg"))).To(Equal(config))

		i++
	}
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

	// we should not see any permutations in serialized config
	i := 0
	for i < maxRuntimes {
		config, errs := psr.RenderMainConfig(false)
		// fmt.Println(config)
		g.Expect(errs).NotTo(HaveOccurred())
		g.Expect(string(getExpectedCfg("./expected/fluentd-mixed-cfgs-multi-tenant-output.cfg"))).To(Equal(config))

		i++
	}
}

func Test_OutputWithBuffer(t *testing.T) {
	g := NewGomegaWithT(t)
	sl := plugins.NewSecretLoader(nil, Fluentd.Namespace, logr.Logger{})

	psr := fluentdv1alpha1.NewGlobalPluginResources("main")
	psr.CombineGlobalInputsPlugins(sl, Fluentd.Spec.GlobalInputs)

	clustercfgRouter, err := psr.BuildCfgRouter(&FluentdClusterFluentdConfig1)
	g.Expect(err).NotTo(HaveOccurred())
	clusterFilters := []fluentdv1alpha1.ClusterFilter{FluentdClusterFilter1}
	clusterOutputs := []fluentdv1alpha1.ClusterOutput{FluentdClusterOutputBuffer}
	clustercfgResources, _ := psr.PatchAndFilterClusterLevelResources(sl, FluentdClusterFluentdConfig1.GetCfgId(), []fluentdv1alpha1.ClusterInput{}, clusterFilters, clusterOutputs)
	err = psr.WithCfgResources(*clustercfgRouter.Label, clustercfgResources)
	g.Expect(err).NotTo(HaveOccurred())

	// we should not see any permutations in serialized config
	i := 0
	for i < maxRuntimes {
		config, errs := psr.RenderMainConfig(false)
		// fmt.Println(config)
		g.Expect(errs).NotTo(HaveOccurred())
		g.Expect(string(getExpectedCfg("./expected/fluentd-cluster-cfg-output-buffer-example.cfg"))).To(Equal(config))

		i++
	}
}

func Test_OutputWithMemoryBuffer(t *testing.T) {
	g := NewGomegaWithT(t)
	sl := plugins.NewSecretLoader(nil, Fluentd.Namespace, logr.Logger{})

	psr := fluentdv1alpha1.NewGlobalPluginResources("main")
	psr.CombineGlobalInputsPlugins(sl, Fluentd.Spec.GlobalInputs)

	clustercfgRouter, err := psr.BuildCfgRouter(&FluentdClusterFluentdConfig1)
	g.Expect(err).NotTo(HaveOccurred())
	clusterFilters := []fluentdv1alpha1.ClusterFilter{FluentdClusterFilter1}
	clusterOutputs := []fluentdv1alpha1.ClusterOutput{FluentdClusterOutputMemoryBuffer}
	clustercfgResources, _ := psr.PatchAndFilterClusterLevelResources(sl, FluentdClusterFluentdConfig1.GetCfgId(), []fluentdv1alpha1.ClusterInput{}, clusterFilters, clusterOutputs)
	err = psr.WithCfgResources(*clustercfgRouter.Label, clustercfgResources)
	g.Expect(err).NotTo(HaveOccurred())

	// we should not see any permutations in serialized config
	i := 0
	for i < maxRuntimes {
		config, errs := psr.RenderMainConfig(false)
		// fmt.Println(config)
		g.Expect(errs).NotTo(HaveOccurred())
		g.Expect(string(getExpectedCfg("./expected/fluentd-cluster-cfg-output-memory-buffer.cfg"))).To(Equal(config))

		i++
	}
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
			Spec: fluentdv1alpha1.FilterSpec(filterspec1),
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
			Spec: fluentdv1alpha1.FilterSpec(filterspec2),
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
			Spec: fluentdv1alpha1.FilterSpec(filterspec3),
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
			Spec: fluentdv1alpha1.OutputSpec(outputspec1),
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
			Spec: fluentdv1alpha1.OutputSpec(outputspec2),
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
			Spec: fluentdv1alpha1.OutputSpec(outputspec3),
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

	// we should not see any permutations in serialized config
	i := 0
	for i < maxRuntimes {
		config, errs := psr.RenderMainConfig(false)
		// fmt.Println(config)
		g.Expect(errs).NotTo(HaveOccurred())
		g.Expect(string(getExpectedCfg("./expected/duplicate-removal-cr-specs.cfg"))).To(Equal(config))

		i++
	}
}

func Test_RecordTransformer(t *testing.T) {
	g := NewGomegaWithT(t)
	sl := plugins.NewSecretLoader(nil, Fluentd.Namespace, logr.Logger{})

	psr := fluentdv1alpha1.NewGlobalPluginResources("main")
	psr.CombineGlobalInputsPlugins(sl, Fluentd.Spec.GlobalInputs)

	clustercfgRouter, err := psr.BuildCfgRouter(&FluentdClusterFluentdConfig1)
	g.Expect(err).NotTo(HaveOccurred())
	clusterFilters := []fluentdv1alpha1.ClusterFilter{FluentdClusterRecordTransformerFilter}
	clusterOutputs := []fluentdv1alpha1.ClusterOutput{FluentdClusterOutputCluster}
	clustercfgResources, _ := psr.PatchAndFilterClusterLevelResources(sl, FluentdClusterFluentdConfig1.GetCfgId(), []fluentdv1alpha1.ClusterInput{}, clusterFilters, clusterOutputs)
	err = psr.WithCfgResources(*clustercfgRouter.Label, clustercfgResources)
	g.Expect(err).NotTo(HaveOccurred())

	// we should not see any permutations in serialized config
	i := 0
	for i < maxRuntimes {
		config, errs := psr.RenderMainConfig(false)
		// fmt.Println(config)
		g.Expect(errs).NotTo(HaveOccurred())
		g.Expect(string(getExpectedCfg("./expected/fluentd-cluster-cfg-filter-recordTransformer.cfg"))).To(Equal(config))

		i++
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
		return strings.TrimSuffix(fmt.Sprintf("%s", v), "\n"), nil
	}
}
