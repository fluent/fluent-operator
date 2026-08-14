/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"os"
	"sort"

	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins/filter"
	"github.com/fluent/fluent-operator/v3/pkg/utils"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"strings"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"

	fluentbitv1alpha2 "github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2"
)

// FluentBitConfigReconciler reconciles a FluentBitConfig object
type FluentBitConfigReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

var storeNamespaces map[string]bool

func computeConfigHash(
	configFileName, mainAppCfg, parserCfg, multilineParserCfg string,
	scripts []fluentbitv1alpha2.Script,
) string {
	h := sha256.New()
	// Use null bytes as delimiters to prevent hash collisions between fields
	_, _ = h.Write([]byte(configFileName))
	_, _ = h.Write([]byte{0})
	_, _ = h.Write([]byte(mainAppCfg))
	_, _ = h.Write([]byte{0})
	_, _ = h.Write([]byte(parserCfg))
	_, _ = h.Write([]byte{0})
	_, _ = h.Write([]byte(multilineParserCfg))
	_, _ = h.Write([]byte{0})

	sortedScripts := make([]fluentbitv1alpha2.Script, len(scripts))
	copy(sortedScripts, scripts)
	sort.SliceStable(sortedScripts, func(i, j int) bool {
		if sortedScripts[i].Name == sortedScripts[j].Name {
			return sortedScripts[i].Content < sortedScripts[j].Content
		}
		return sortedScripts[i].Name < sortedScripts[j].Name
	})

	for _, s := range sortedScripts {
		_, _ = h.Write([]byte(s.Name))
		_, _ = h.Write([]byte{0})
		_, _ = h.Write([]byte(s.Content))
		_, _ = h.Write([]byte{0})
	}

	return fmt.Sprintf("%x", h.Sum(nil))
}

func (r *FluentBitConfigReconciler) updateSecretIfNeeded(
	ctx context.Context,
	cfgName, ns, configFileName, mainAppCfg, parserCfg, multilineParserCfg string,
	scripts []fluentbitv1alpha2.Script,
	setControllerRef func(*corev1.Secret) error,
) error {
	newConfigHash := computeConfigHash(configFileName, mainAppCfg, parserCfg, multilineParserCfg, scripts)

	sec := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cfgName,
			Namespace: ns,
		},
	}
	existingSecret := &corev1.Secret{}
	err := r.Get(ctx, client.ObjectKey{Name: cfgName, Namespace: ns}, existingSecret)
	needsUpdate := false

	if err != nil {
		if errors.IsNotFound(err) {
			needsUpdate = true
		} else {
			return err
		}
	} else {
		existingHash, ok := existingSecret.Annotations[fluentBitConfigHashAnnotation]
		if !ok || existingHash != newConfigHash {
			needsUpdate = true
		}
	}

	if needsUpdate {
		opResult, err := controllerutil.CreateOrPatch(
			ctx, r.Client, sec, func() error {
				if sec.Annotations == nil {
					sec.Annotations = make(map[string]string)
				}
				sec.Annotations[fluentBitConfigHashAnnotation] = newConfigHash

				sec.Data = map[string][]byte{
					configFileName:           []byte(mainAppCfg),
					"parsers.conf":           []byte(parserCfg),
					"parsers_multiline.conf": []byte(multilineParserCfg),
				}
				for _, s := range scripts {
					sec.Data[s.Name] = []byte(s.Content)
				}
				sec.SetOwnerReferences(nil)
				if err := setControllerRef(sec); err != nil {
					return err
				}
				return nil
			},
		)
		if err != nil {
			return err
		}

		// Only log update if CreateOrPatch actually performed a write
		if opResult != controllerutil.OperationResultNone {
			r.Log.Info(
				"Fluent Bit main configuration has updated", "logging-control-plane", ns, "fluentbitconfig", cfgName,
				"secret", sec.Name, "config-hash", newConfigHash,
			)
		} else {
			r.Log.V(1).Info(
				"Fluent Bit configuration unchanged (concurrent update)", "logging-control-plane", ns, "fluentbitconfig", cfgName,
				"secret", cfgName, "config-hash", newConfigHash,
			)
		}
	} else {
		r.Log.V(1).Info(
			"Fluent Bit configuration unchanged, skipping update", "logging-control-plane", ns, "fluentbitconfig", cfgName,
			"secret", cfgName, "config-hash", newConfigHash,
		)
	}

	return nil
}

// +kubebuilder:rbac:groups=fluentbit.fluent.io,resources=clusterfluentbitconfigs,verbs=list;watch
// +kubebuilder:rbac:groups=fluentbit.fluent.io,resources=fluentbitconfigs,verbs=list;watch
// +kubebuilder:rbac:groups=fluentbit.fluent.io,resources=collectors,verbs=list;watch
// +kubebuilder:rbac:groups=fluentbit.fluent.io,resources=clusterinputs;clusterfilters;clusteroutputs;clusterparsers;clustermultilineparsers,verbs=list;watch
// +kubebuilder:rbac:groups=fluentbit.fluent.io,resources=filters;outputs;parsers;multilineparsers,verbs=list;watch
// +kubebuilder:rbac:groups=core,resources=configmaps,verbs=get
// +kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch;create;patch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the FluentBitConfig object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *FluentBitConfigReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = r.Log.WithValues("fluentbitconfig", req.NamespacedName)

	// Re-initialize during each reconcile loop to clear namespace names
	// we will repopulate namespaces that have a FluentBitConfig CR
	storeNamespaces = make(map[string]bool)

	var cfgs fluentbitv1alpha2.ClusterFluentBitConfigList
	if err := r.List(ctx, &cfgs); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// Secrets rendered during this loop, keyed by "<namespace>/<name>". A
	// ClusterFluentBitConfig referenced by a FluentBit always wins, so that
	// adding a Collector can never change the configuration a DaemonSet is
	// already using.
	rendered := make(map[string]bool)

	if err := r.reconcileFluentBits(ctx, cfgs, rendered); err != nil {
		return ctrl.Result{}, err
	}

	if err := r.reconcileCollectors(ctx, cfgs, rendered); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// reconcileFluentBits renders the ClusterFluentBitConfigs referenced by the FluentBit
// DaemonSets, including the namespace level (multi-tenant) plugins they select.
func (r *FluentBitConfigReconciler) reconcileFluentBits(
	ctx context.Context, cfgs fluentbitv1alpha2.ClusterFluentBitConfigList, rendered map[string]bool,
) error {
	var fbs fluentbitv1alpha2.FluentBitList
	if err := r.List(ctx, &fbs); err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}

	for _, fb := range fbs.Items {
		for _, cfg := range cfgs.Items {
			// List all inputs matching the label selector.
			if cfg.Name != fb.Spec.FluentBitConfigName {
				continue
			}

			clusterPlugins, err := r.listClusterPlugins(ctx, &cfg)
			if err != nil {
				return err
			}

			// List all the namespace level resources if they exist and generate configs to mutate tags
			nsFilterLists, nsOutputLists, nsParserLists,
				nsClusterParserLists, nsMultilineParserLists, nsClusterMultilineParserLists,
				rewriteTagConfigs, err := r.processNamespacedFluentBitCfgs(
				ctx, fb, clusterPlugins.inputs, cfg.Spec.ConfigFileFormat,
			)

			if err != nil {
				return err
			}

			var ns string
			if cfg.Spec.Namespace != nil {
				ns = fmt.Sprint(*cfg.Spec.Namespace)
			} else {
				ns = os.Getenv("NAMESPACE")
			}

			if err := r.renderAndStoreConfig(ctx, &cfg, ns, clusterPlugins, namespacedPlugins{
				filters:                 nsFilterLists,
				outputs:                 nsOutputLists,
				parsers:                 nsParserLists,
				clusterParsers:          nsClusterParserLists,
				multilineParsers:        nsMultilineParserLists,
				clusterMultilineParsers: nsClusterMultilineParserLists,
				rewriteTagConfigs:       rewriteTagConfigs,
			}); err != nil {
				return err
			}
			rendered[fmt.Sprintf("%s/%s", ns, cfg.Name)] = true
		}
	}

	return nil
}

// reconcileCollectors renders the ClusterFluentBitConfigs referenced by the Collector
// StatefulSets. The Collector controller blocks until a Secret named after
// spec.fluentBitConfigName shows up in the Collector's own namespace, so without this
// nothing is ever created for a Collector, see
// https://github.com/fluent/fluent-operator/issues/1436.
//
// Collectors have no namespaced FluentBitConfig selector, hence only cluster scoped
// plugins are rendered for them.
func (r *FluentBitConfigReconciler) reconcileCollectors(
	ctx context.Context, cfgs fluentbitv1alpha2.ClusterFluentBitConfigList, rendered map[string]bool,
) error {
	var collectors fluentbitv1alpha2.CollectorList
	if err := r.List(ctx, &collectors); err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}

	for _, co := range collectors.Items {
		for _, cfg := range cfgs.Items {
			if cfg.Name != co.Spec.FluentBitConfigName {
				continue
			}

			// The Collector expects its configuration next to itself. An explicit
			// spec.namespace on the ClusterFluentBitConfig still wins, to stay
			// consistent with the FluentBit code path.
			ns := co.Namespace
			if cfg.Spec.Namespace != nil {
				ns = fmt.Sprint(*cfg.Spec.Namespace)
				if ns != co.Namespace {
					r.Log.Info(
						"ClusterFluentBitConfig pins its Secret to a namespace other than the Collector's, "+
							"the Collector will keep waiting for its configuration",
						"clusterfluentbitconfig", cfg.Name, "logging-control-plane", ns,
						"collector", co.Name, "collector-namespace", co.Namespace,
					)
				}
			}

			if rendered[fmt.Sprintf("%s/%s", ns, cfg.Name)] {
				r.Log.V(1).Info(
					"Fluent Bit configuration already rendered in this reconcile, skipping",
					"logging-control-plane", ns, "fluentbitconfig", cfg.Name, "collector", co.Name,
				)
				continue
			}

			clusterPlugins, err := r.listClusterPlugins(ctx, &cfg)
			if err != nil {
				return err
			}

			if err := r.renderAndStoreConfig(ctx, &cfg, ns, clusterPlugins, namespacedPlugins{}); err != nil {
				return err
			}
			rendered[fmt.Sprintf("%s/%s", ns, cfg.Name)] = true
		}
	}

	return nil
}

// clusterPlugins holds the cluster scoped plugins selected by a ClusterFluentBitConfig.
type clusterPlugins struct {
	inputs           fluentbitv1alpha2.ClusterInputList
	filters          fluentbitv1alpha2.ClusterFilterList
	outputs          fluentbitv1alpha2.ClusterOutputList
	parsers          fluentbitv1alpha2.ClusterParserList
	multilineParsers fluentbitv1alpha2.ClusterMultilineParserList
}

// namespacedPlugins holds the namespace level (multi-tenant) plugins that are merged
// into the rendered configuration together with the cluster scoped ones. It is empty
// for Collectors.
type namespacedPlugins struct {
	filters                 []fluentbitv1alpha2.FilterList
	outputs                 []fluentbitv1alpha2.OutputList
	parsers                 []fluentbitv1alpha2.ParserList
	clusterParsers          []fluentbitv1alpha2.ClusterParserList
	multilineParsers        []fluentbitv1alpha2.MultilineParserList
	clusterMultilineParsers []fluentbitv1alpha2.ClusterMultilineParserList
	rewriteTagConfigs       []string
}

// listClusterPlugins lists every cluster scoped plugin matching the label selectors of
// the given ClusterFluentBitConfig.
func (r *FluentBitConfigReconciler) listClusterPlugins(
	ctx context.Context, cfg *fluentbitv1alpha2.ClusterFluentBitConfig,
) (clusterPlugins, error) {
	var cp clusterPlugins

	if err := listClusterResources(ctx, r.Client, &cfg.Spec.InputSelector, &cp.inputs); err != nil {
		return cp, err
	}
	if err := listClusterResources(ctx, r.Client, &cfg.Spec.FilterSelector, &cp.filters); err != nil {
		return cp, err
	}
	if err := listClusterResources(ctx, r.Client, &cfg.Spec.OutputSelector, &cp.outputs); err != nil {
		return cp, err
	}
	if err := listClusterResources(ctx, r.Client, &cfg.Spec.ParserSelector, &cp.parsers); err != nil {
		return cp, err
	}
	if err := listClusterResources(
		ctx, r.Client, &cfg.Spec.MultilineParserSelector, &cp.multilineParsers,
	); err != nil {
		return cp, err
	}

	return cp, nil
}

// renderAndStoreConfig renders the main, parser and multiline parser configurations
// plus the Lua scripts, and stores them in a Secret named after the
// ClusterFluentBitConfig in the given namespace.
func (r *FluentBitConfigReconciler) renderAndStoreConfig(
	ctx context.Context, cfg *fluentbitv1alpha2.ClusterFluentBitConfig, ns string,
	cluster clusterPlugins, namespaced namespacedPlugins,
) error {
	cl := plugins.NewConfigMapLoader(r.Client, ns)
	// load scripts for namespaced filters
	nsScripts, err := cfg.RenderNamespacedLuaScript(cl, namespaced.filters)
	if err != nil {
		return err
	}
	// load scripts for cluster filters
	scripts, err := cfg.RenderLuaScript(cl, cluster.filters, ns)
	if err != nil {
		return err
	}
	scripts = append(scripts, nsScripts...)

	// Inject config data into Secret
	sl := plugins.NewSecretLoader(r.Client, ns)
	mainAppCfg, err := cfg.RenderMainConfigWithTargetFormat(
		sl, cluster.inputs, cluster.filters, cluster.outputs, namespaced.filters, namespaced.outputs,
		namespaced.rewriteTagConfigs, cfg.Spec.ConfigFileFormat,
	)
	if err != nil {
		return err
	}
	parserCfg, err := cfg.RenderParserConfig(sl, cluster.parsers, namespaced.parsers, namespaced.clusterParsers)
	if err != nil {
		return err
	}
	multilineParserCfg, err := cfg.RenderMultilineParserConfig(
		sl, cluster.multilineParsers, namespaced.multilineParsers, namespaced.clusterMultilineParsers,
	)
	if err != nil {
		return err
	}

	configFileName := "fluent-bit.conf"
	if cfg.Spec.ConfigFileFormat != nil && *cfg.Spec.ConfigFileFormat == configFileFormatYaml {
		configFileName = "fluent-bit.yaml"
	}

	return r.updateSecretIfNeeded(
		ctx, cfg.Name, ns, configFileName, mainAppCfg, parserCfg, multilineParserCfg, scripts,
		func(sec *corev1.Secret) error {
			return ctrl.SetControllerReference(cfg, sec, r.Scheme)
		},
	)
}

func (r *FluentBitConfigReconciler) processNamespacedFluentBitCfgs(
	ctx context.Context, fb fluentbitv1alpha2.FluentBit, inputs fluentbitv1alpha2.ClusterInputList,
	configFileFormat *string,
) (
	[]fluentbitv1alpha2.FilterList, []fluentbitv1alpha2.OutputList,
	[]fluentbitv1alpha2.ParserList, []fluentbitv1alpha2.ClusterParserList,
	[]fluentbitv1alpha2.MultilineParserList, []fluentbitv1alpha2.ClusterMultilineParserList, []string, error,
) {
	var nsCfgs fluentbitv1alpha2.FluentBitConfigList
	// set of rewrite_tag plugin configs to mutate tags for log records coming out of a namespace
	selector, err := metav1.LabelSelectorAsSelector(&fb.Spec.NamespacedFluentBitCfgSelector)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, err
	}

	if err := r.List(ctx, &nsCfgs, client.MatchingLabelsSelector{Selector: selector}); err != nil {
		return nil, nil, nil, nil, nil, nil, nil, err
	}

	// Sort namespace configs by namespace and name for deterministic rendering
	sort.Slice(nsCfgs.Items, func(i, j int) bool {
		if nsCfgs.Items[i].Namespace == nsCfgs.Items[j].Namespace {
			return nsCfgs.Items[i].Name < nsCfgs.Items[j].Name
		}
		return nsCfgs.Items[i].Namespace < nsCfgs.Items[j].Namespace
	})

	filters := make([]fluentbitv1alpha2.FilterList, 0, len(nsCfgs.Items))
	outputs := make([]fluentbitv1alpha2.OutputList, 0, len(nsCfgs.Items))
	parsers := make([]fluentbitv1alpha2.ParserList, 0, len(nsCfgs.Items))
	clusterParsers := make([]fluentbitv1alpha2.ClusterParserList, 0, len(nsCfgs.Items))
	multilineParsers := make([]fluentbitv1alpha2.MultilineParserList, 0, len(nsCfgs.Items))
	clusterMultilineParsers := make([]fluentbitv1alpha2.ClusterMultilineParserList, 0, len(nsCfgs.Items))
	rewriteTagConfigs := make([]string, 0, len(nsCfgs.Items))

	// Form a slice of list of resources per namespace
	for _, cfg := range nsCfgs.Items {
		filterList, outputList, parserList,
			clusterParserList, multilineParsersList,
			clusterMultilineParsersList, err := r.ListFluentBitConfigResources(
			ctx, cfg,
		)
		if err != nil {
			return filters, outputs, parsers,
				clusterParsers, multilineParsers, clusterMultilineParsers,
				nil, err
		}
		filters = append(filters, filterList)
		outputs = append(outputs, outputList)
		parsers = append(parsers, parserList)
		clusterParsers = append(clusterParsers, clusterParserList)
		multilineParsers = append(multilineParsers, multilineParsersList)
		clusterMultilineParsers = append(clusterMultilineParsers, clusterMultilineParsersList)

		if _, ok := storeNamespaces[cfg.Namespace]; !ok {
			rewriteTagConfig, err := r.generateRewriteTagConfig(cfg, inputs, configFileFormat)
			if err != nil {
				return filters, outputs, parsers,
					clusterParsers, multilineParsers, clusterMultilineParsers,
					nil, err
			}
			if rewriteTagConfig != "" {
				rewriteTagConfigs = append(rewriteTagConfigs, rewriteTagConfig)
				storeNamespaces[cfg.Namespace] = true
			}
		}
	}

	return filters, outputs, parsers, clusterParsers, multilineParsers, clusterMultilineParsers, rewriteTagConfigs, nil
}

func listClusterResources[T client.ObjectList](
	ctx context.Context, cli client.Client, selector *metav1.LabelSelector, list T,
) error {
	sel, err := metav1.LabelSelectorAsSelector(selector)
	if err != nil {
		return err
	}
	if err := cli.List(ctx, list, client.MatchingLabelsSelector{Selector: sel}); err != nil {
		return err
	}
	return nil
}

func listNamespacedResources[T client.ObjectList](
	ctx context.Context, cli client.Client, list T, namespace string, selector *metav1.LabelSelector,
) error {
	sel, err := metav1.LabelSelectorAsSelector(selector)
	if err != nil {
		return err
	}
	matchingLabelSelector := client.MatchingLabelsSelector{Selector: sel}
	if err := cli.List(ctx, list, client.InNamespace(namespace), matchingLabelSelector); err != nil {
		return err
	}
	return nil
}

// ListFluentBitConfigResources lists all resources (both namespaced and cluster-scoped)
// needed by the given FluentBitConfig.
func (r *FluentBitConfigReconciler) ListFluentBitConfigResources(
	ctx context.Context, cfg fluentbitv1alpha2.FluentBitConfig,
) (
	fluentbitv1alpha2.FilterList,
	fluentbitv1alpha2.OutputList, fluentbitv1alpha2.ParserList, fluentbitv1alpha2.ClusterParserList,
	fluentbitv1alpha2.MultilineParserList, fluentbitv1alpha2.ClusterMultilineParserList, error,
) {
	var filters fluentbitv1alpha2.FilterList
	var outputs fluentbitv1alpha2.OutputList
	var parsers fluentbitv1alpha2.ParserList
	var clusterParsers fluentbitv1alpha2.ClusterParserList
	var multipleParsers fluentbitv1alpha2.MultilineParserList
	var clusterMultipleParsers fluentbitv1alpha2.ClusterMultilineParserList

	if err := listNamespacedResources(ctx, r.Client, &filters, cfg.Namespace, &cfg.Spec.FilterSelector); err != nil {
		return filters, outputs, parsers, clusterParsers, multipleParsers, clusterMultipleParsers, err
	}

	if err := listNamespacedResources(ctx, r.Client, &outputs, cfg.Namespace, &cfg.Spec.OutputSelector); err != nil {
		return filters, outputs, parsers, clusterParsers, multipleParsers, clusterMultipleParsers, err
	}
	if err := listNamespacedResources(ctx, r.Client, &parsers, cfg.Namespace, &cfg.Spec.ParserSelector); err != nil {
		return filters, outputs, parsers, clusterParsers, multipleParsers, clusterMultipleParsers, err
	}

	if err := listClusterResources(ctx, r.Client, &cfg.Spec.ClusterParserSelector, &clusterParsers); err != nil {
		return filters, outputs, parsers, clusterParsers, multipleParsers, clusterMultipleParsers, err
	}

	if err := listNamespacedResources(
		ctx,
		r.Client,
		&multipleParsers,
		cfg.Namespace,
		&cfg.Spec.MultilineParserSelector,
	); err != nil {
		return filters, outputs, parsers, clusterParsers, multipleParsers, clusterMultipleParsers, err
	}

	if err := listClusterResources(
		ctx, r.Client, &cfg.Spec.ClusterMultilineParserSelector, &clusterMultipleParsers,
	); err != nil {
		return filters, outputs, parsers, clusterParsers, multipleParsers, clusterMultipleParsers, err
	}

	// Update the name of the local copies of cluster level parsers.
	// The intention is to have each namespace use their own copy
	// of the cluster parser.
	for i := range clusterParsers.Items {
		clusterParsers.Items[i].Name = fmt.Sprintf(
			"%s-%x", clusterParsers.Items[i].Name, md5.Sum([]byte(cfg.Namespace)),
		)
	}

	for i := range clusterMultipleParsers.Items {
		clusterMultipleParsers.Items[i].Name = fmt.Sprintf(
			"%s-%x", clusterMultipleParsers.Items[i].Name, md5.Sum([]byte(cfg.Namespace)),
		)
	}

	return filters, outputs, parsers, clusterParsers, multipleParsers, clusterMultipleParsers, nil
}

func (r *FluentBitConfigReconciler) generateRewriteTagConfig(
	cfg fluentbitv1alpha2.FluentBitConfig, inputs fluentbitv1alpha2.ClusterInputList, configFileFormat *string,
) (string, error) {
	var tag string
	for _, input := range inputs.Items {
		if input.Spec.Tail == nil || !strings.Contains(input.Spec.Tail.Path, "/var/log/containers") {
			continue
		}
		tag = input.Spec.Tail.Tag
		idx := strings.Index(tag, ".")
		if idx == -1 {
			tag = ""
		} else {
			tag = tag[:idx+1] + "*"
		}
	}
	if tag == "" {
		return "", nil
	}

	rewriteTag := &filter.RewriteTag{
		Rules: []string{
			fmt.Sprintf("$kubernetes['namespace_name'] ^(%s)$ %x.$TAG false", cfg.Namespace, md5.Sum([]byte(cfg.Namespace))),
		},
	}
	if cfg.Spec.Service != nil {
		if cfg.Spec.Service.EmitterName != "" {
			rewriteTag.EmitterName = cfg.Spec.Service.EmitterName
		} else {
			rewriteTag.EmitterName = fmt.Sprintf("re_emitted_%x", md5.Sum([]byte(cfg.Namespace)))
		}
		rewriteTag.EmitterStorageType = cfg.Spec.Service.EmitterStorageType
		rewriteTag.EmitterMemBufLimit = cfg.Spec.Service.EmitterMemBufLimit
	}

	filterList := fluentbitv1alpha2.ClusterFilterList{
		Items: []fluentbitv1alpha2.ClusterFilter{
			{
				Spec: fluentbitv1alpha2.FilterSpec{
					Match:       tag,
					FilterItems: []fluentbitv1alpha2.FilterItem{{RewriteTag: rewriteTag}},
				},
			},
		},
	}

	sl := plugins.NewSecretLoader(nil, "")
	if configFileFormat != nil && *configFileFormat == configFileFormatYaml {
		rendered, err := filterList.LoadAsYaml(sl, 1)
		if err != nil {
			return "", err
		}
		// Strip the "filters:" header so callers can merge this into the
		// single "filters:" section of the main YAML config instead of
		// emitting a second, duplicate key (see RenderMainConfigInYaml).
		header := fmt.Sprintf("%sfilters:\n", utils.YamlIndent(1))
		return strings.TrimPrefix(rendered, header), nil
	}
	return filterList.Load(sl)
}

func (r *FluentBitConfigReconciler) SetupWithManager(mgr ctrl.Manager) error {
	ctx := context.Background()
	if err := mgr.GetFieldIndexer().IndexField(
		ctx, &corev1.Secret{}, fluentbitOwnerKey, func(rawObj client.Object) []string {
			// Grab the job object, extract the owner.
			sec := rawObj.(*corev1.Secret)
			owner := metav1.GetControllerOf(sec)
			if owner == nil {
				return nil
			}
			// Make sure it's a FluentBitConfig. If so, return it.
			if owner.APIVersion != fluentbitApiGVStr || owner.Kind != "FluentBitConfig" {
				return nil
			}
			return []string{owner.Name}
		},
	); err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		Named("FluentBitConfig").
		For(&fluentbitv1alpha2.FluentBit{}).
		Owns(&corev1.Secret{}).
		Watches(&fluentbitv1alpha2.ClusterFluentBitConfig{}, &handler.EnqueueRequestForObject{}).
		Watches(&fluentbitv1alpha2.FluentBitConfig{}, &handler.EnqueueRequestForObject{}).
		Watches(&fluentbitv1alpha2.Collector{}, &handler.EnqueueRequestForObject{}).
		Watches(&fluentbitv1alpha2.ClusterInput{}, &handler.EnqueueRequestForObject{}).
		Watches(&fluentbitv1alpha2.ClusterFilter{}, &handler.EnqueueRequestForObject{}).
		Watches(&fluentbitv1alpha2.ClusterOutput{}, &handler.EnqueueRequestForObject{}).
		Watches(&fluentbitv1alpha2.ClusterParser{}, &handler.EnqueueRequestForObject{}).
		Watches(&fluentbitv1alpha2.ClusterMultilineParser{}, &handler.EnqueueRequestForObject{}).
		Watches(&fluentbitv1alpha2.Filter{}, &handler.EnqueueRequestForObject{}).
		Watches(&fluentbitv1alpha2.Output{}, &handler.EnqueueRequestForObject{}).
		Watches(&fluentbitv1alpha2.Parser{}, &handler.EnqueueRequestForObject{}).
		Watches(&fluentbitv1alpha2.MultilineParser{}, &handler.EnqueueRequestForObject{}).
		Complete(r)
}
