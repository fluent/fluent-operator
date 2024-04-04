package v1alpha1

import (
	"bytes"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/fluent/fluent-operator/v2/apis/fluentd/v1alpha1/plugins"
	"github.com/fluent/fluent-operator/v2/apis/fluentd/v1alpha1/plugins/filter"
	"github.com/fluent/fluent-operator/v2/apis/fluentd/v1alpha1/plugins/input"
	"github.com/fluent/fluent-operator/v2/apis/fluentd/v1alpha1/plugins/output"
	"github.com/fluent/fluent-operator/v2/apis/fluentd/v1alpha1/plugins/params"
	fluentdRouter "github.com/fluent/fluent-operator/v2/pkg/fluentd/router"
)

// +kubebuilder:object:generate=false
type Renderer interface {
	GetNamespace() string
	GetName() string
	GetCfgId() string
	GetWatchedLabels() map[string]string
	GetWatchedNamespaces() []string
	GetWatchedContainers() []string
	GetWatchedHosts() []string
}

// +kubebuilder:object:generate=false
// Global pluginstores for the fluentd.
type PluginResources struct {
	InputPlugins         []params.PluginStore
	MainRouterPlugins    params.PluginStore
	LabelPluginResources []params.PluginStore
}

// +kubebuilder:object:generate=false
// All the filter/output selected to this cfg
type CfgResources struct {
	InputPlugins  []params.PluginStore
	FilterPlugins []params.PluginStore
	OutputPlugins []params.PluginStore

	// the hash codes used to depulicate removel
	InputsHashcodes  map[string]bool
	FiltersHashcodes map[string]bool
	OutputsHashcodes map[string]bool
}

// NewGlobalPluginResources represents a combined global fluentd resources
func NewGlobalPluginResources(globalId string) *PluginResources {
	globalMainRouter := fluentdRouter.NewGlobalRouter(globalId)
	return &PluginResources{
		InputPlugins:         make([]params.PluginStore, 0),
		MainRouterPlugins:    *globalMainRouter,
		LabelPluginResources: make([]params.PluginStore, 0),
	}
}

func NewCfgResources() *CfgResources {
	return &CfgResources{
		FilterPlugins: make([]params.PluginStore, 0),
		OutputPlugins: make([]params.PluginStore, 0),

		InputsHashcodes:  make(map[string]bool),
		FiltersHashcodes: make(map[string]bool),
		OutputsHashcodes: make(map[string]bool),
	}
}

func (pgr *PluginResources) CombineGlobalInputsPlugins(sl plugins.SecretLoader, inputs []input.Input) []string {
	errs := make([]string, 0)
	for _, f := range inputs {
		ps, err := f.Params(sl)
		if err != nil {
			errs = append(errs, err.Error())
			continue
		}
		pgr.InputPlugins = append(pgr.InputPlugins, *ps)
	}
	return errs
}

func (pgr *PluginResources) BuildCfgRouter(cfg Renderer) (*fluentdRouter.Route, error) {
	matches := []*fluentdRouter.RouteMatch{
		{
			Labels:         cfg.GetWatchedLabels(),
			Namespaces:     cfg.GetWatchedNamespaces(),
			Hosts:          cfg.GetWatchedHosts(),
			ContainerNames: cfg.GetWatchedContainers(),
		},
	}

	cfgRoute, err := fluentdRouter.NewRoute(cfg.GetCfgId(), cfg.GetNamespace(), cfg.GetName(), matches)
	if err != nil {
		return nil, err
	}

	// Each fluentd config has its own route plugin
	routePluginStore, err := cfgRoute.NewRoutePlugin()
	if err != nil {
		return nil, err
	}

	// Insert the route to the MainRouter
	pgr.MainRouterPlugins.InsertChilds(routePluginStore)

	return cfgRoute, nil
}

// PatchAndFilterClusterLevelResources will combine and patch all the cluster CRs that the fluentdconfig selected,
// convert the related filter/output pluginstores to the global pluginresources.
func (pgr *PluginResources) PatchAndFilterClusterLevelResources(
	sl plugins.SecretLoader,
	cfgId string,
	clusterInputs []ClusterInput,
	clusterfilters []ClusterFilter,
	clusteroutputs []ClusterOutput,
) (*CfgResources, []string) {
	// To store all filters/outputs plugins that this cfg selected
	cfgResources := NewCfgResources()

	errs := make([]string, 0)
	// sort all the CRs by metadata.name
	sort.SliceStable(clusterInputs[:], func(i, j int) bool {
		return clusterInputs[i].Name < clusterInputs[j].Name
	})
	sort.SliceStable(clusterfilters[:], func(i, j int) bool {
		return clusterfilters[i].Name < clusterfilters[j].Name
	})
	sort.SliceStable(clusteroutputs[:], func(i, j int) bool {
		return clusteroutputs[i].Name < clusteroutputs[j].Name
	})
	// List all inputs matching the label selector.
	for _, i := range clusterInputs {
		// patch filterId
		err := cfgResources.filterForInputs(cfgId, "cluster", i.Name, "clusterinput", sl, i.Spec.Inputs)
		if err != nil {
			errs = append(errs, err.Error())
		}
	}

	// List all filters matching the label selector.
	for _, i := range clusterfilters {
		// patch filterId
		err := cfgResources.filterForFilters(cfgId, "cluster", i.Name, "clusterfilter", sl, i.Spec.Filters)
		if err != nil {
			errs = append(errs, err.Error())
		}
	}

	// List all outputs matching the label selector.
	for _, i := range clusteroutputs {
		// patch outputId
		err := cfgResources.filterForOutputs(cfgId, "cluster", i.Name, "clusteroutput", sl, i.Spec.Outputs)
		if err != nil {
			errs = append(errs, err.Error())
		}
	}

	return cfgResources, errs
}

// PatchAndFilterNamespacedLevelResources will combine and patch all the cluster CRs that the fluentdconfig selected,
// convert the related filter/output pluginstores to the global pluginresources.
func (pgr *PluginResources) PatchAndFilterNamespacedLevelResources(
	sl plugins.SecretLoader,
	cfgId string,
	inputs []Input,
	filters []Filter,
	outputs []Output,
) (*CfgResources, []string) {
	// To store all filters/outputs plugins that this cfg selected
	cfgResources := NewCfgResources()

	errs := make([]string, 0)

	// List all inputs matching the label selector.
	for _, i := range inputs {
		// patch filterId
		err := cfgResources.filterForInputs(cfgId, i.Namespace, i.Name, "filter", sl, i.Spec.Inputs)
		if err != nil {
			errs = append(errs, err.Error())
		}
	}

	// List all filters matching the label selector.
	for _, i := range filters {
		// patch filterId
		err := cfgResources.filterForFilters(cfgId, i.Namespace, i.Name, "filter", sl, i.Spec.Filters)
		if err != nil {
			errs = append(errs, err.Error())
		}
	}

	// List all outputs matching the label selector.
	for _, i := range outputs {
		// patch outputId
		err := cfgResources.filterForOutputs(cfgId, i.Namespace, i.Name, "output", sl, i.Spec.Outputs)
		if err != nil {
			errs = append(errs, err.Error())
		}
	}

	return cfgResources, errs
}

func (r *CfgResources) filterForInputs(
	cfgId, namespace, name, crdtype string,
	sl plugins.SecretLoader,
	inputs []input.Input,
) error {
	for n, input := range inputs {
		inputId := fmt.Sprintf("%s::%s::%s::%s-%d", cfgId, namespace, crdtype, name, n)
		input.InputCommon.Id = &inputId
		// if input.InputCommon.Tag == nil {
		// 	input.InputCommon.Tag = &params.DefaultTag
		// }

		ps, err := input.Params(sl)
		if err != nil {
			return err
		}

		hashcode := ps.Hash()
		if _, ok := r.InputsHashcodes[hashcode]; ok {
			continue
		}

		r.InputsHashcodes[hashcode] = true
		r.InputPlugins = append(r.InputPlugins, *ps)
	}

	return nil
}

func (r *CfgResources) filterForFilters(
	cfgId, namespace, name, crdtype string,
	sl plugins.SecretLoader,
	filters []filter.Filter,
) error {
	for n, filter := range filters {
		filterId := fmt.Sprintf("%s::%s::%s::%s-%d", cfgId, namespace, crdtype, name, n)
		filter.FilterCommon.Id = &filterId
		if filter.FilterCommon.Tag == nil {
			filter.FilterCommon.Tag = &params.DefaultTag
		}

		ps, err := filter.Params(sl)
		if err != nil {
			return err
		}

		hashcode := ps.Hash()
		if _, ok := r.FiltersHashcodes[hashcode]; ok {
			continue
		}

		r.FiltersHashcodes[hashcode] = true
		r.FilterPlugins = append(r.FilterPlugins, *ps)
	}

	return nil
}

func (r *CfgResources) filterForOutputs(
	cfgId, namespace, name, crdtype string,
	sl plugins.SecretLoader,
	outputs []output.Output,
) error {
	for n, output := range outputs {
		outputId := fmt.Sprintf("%s::%s::%s::%s-%d", cfgId, namespace, crdtype, name, n)
		output.OutputCommon.Id = &outputId
		if output.OutputCommon.Tag == nil {
			output.OutputCommon.Tag = &params.DefaultTag
		}

		ps, err := output.Params(sl)
		if err != nil {
			return err
		}

		hashcode := ps.Hash()
		if _, ok := r.OutputsHashcodes[hashcode]; ok {
			continue
		}

		r.OutputsHashcodes[hashcode] = true
		r.OutputPlugins = append(r.OutputPlugins, *ps)
	}

	return nil
}

// IdentifyCopyAndPatchOutput patches up the controller with the Manager
func (pgr *PluginResources) IdentifyCopyAndPatchOutput(cfgResources *CfgResources) error {
	// patched structure for OutputPlugins
	patchedOutputPlugins := []params.PluginStore{}
	// copyOutputs stores the id if the output is a `copy`
	copyOutputs := map[string]int{}
	// outputs stores the id if the output is not a `copy`
	outputs := map[string][]int{}

	// Iterate over cfgResources.OutputPlugins to identify Copy output
	for id, ps := range cfgResources.OutputPlugins {
		if ps.Store["@type"] == string(params.CopyOutputType) {
			// We store last output when 2 output with the same tag
			copyOutputs[ps.Store["tag"]] = id
		} else {
			outputs[ps.Store["tag"]] = append(outputs[ps.Store["tag"]], id)
		}
	}

	// Patch the outputs
	for k, output := range outputs {
		// Does it exist a copy output for this tag ?
		if c, ok := copyOutputs[k]; ok {
			// Yes, so we patch
			for _, id := range output {
				o := cfgResources.OutputPlugins[id]
				o.Name = "store"
				cfgResources.OutputPlugins[c].InsertChilds(&o)
			}
			patchedOutputPlugins = append(patchedOutputPlugins, cfgResources.OutputPlugins[c])
		} else {
			// No, we don't patch
			for _, id := range output {
				o := cfgResources.OutputPlugins[id]
				patchedOutputPlugins = append(patchedOutputPlugins, o)
			}
		}
	}
	cfgResources.OutputPlugins = patchedOutputPlugins
	return nil
}

// convert the cfg plugins to a label plugin, appends to the global label plugins
func (pgr *PluginResources) WithCfgResources(cfgRouteLabel string, r *CfgResources) error {
	if len(r.InputPlugins) == 0 && len(r.FilterPlugins) == 0 && len(r.OutputPlugins) == 0 {
		return errors.New("no filter plugins and no output plugins matched")
	}

	// insert input plugins of this fluentd config
	pgr.InputPlugins = append(pgr.InputPlugins, r.InputPlugins...)

	cfgLabelPlugin := params.NewPluginStore("label")
	cfgLabelPlugin.InsertPairs("tag", cfgRouteLabel)

	// insert filter plugins of this fluentd config
	for _, filter := range r.FilterPlugins {
		childFilter := filter
		cfgLabelPlugin.InsertChilds(&childFilter)
	}

	// insert output plugins of this fluentd config
	for _, output := range r.OutputPlugins {
		childOutput := output
		cfgLabelPlugin.InsertChilds(&childOutput)
	}

	pgr.LabelPluginResources = append(pgr.LabelPluginResources, *cfgLabelPlugin)
	return nil
}

func (pgr *PluginResources) RenderMainConfig(enableMultiWorkers bool) (string, error) {
	if len(pgr.InputPlugins) == 0 && len(pgr.LabelPluginResources) == 0 {
		return "", fmt.Errorf("no plugins detect")
	}

	var buf bytes.Buffer

	// sort global inputs
	inputs := ByHashcode(pgr.InputPlugins)
	for _, pluginStore := range inputs {
		if enableMultiWorkers {
			pluginStore.SetIgnorePath()
		}
		buf.WriteString(pluginStore.String())
	}

	// sort main routers
	childRouters := ByRouteLabelsPointers(pgr.MainRouterPlugins.Childs)
	sort.SliceStable(childRouters[:], childRouters.Less)
	pgr.MainRouterPlugins.Childs = childRouters
	if enableMultiWorkers {
		pgr.MainRouterPlugins.SetIgnorePath()
	}
	buf.WriteString(pgr.MainRouterPlugins.String())

	// sort label plugins
	labelPlugins := ByTags(pgr.LabelPluginResources)
	sort.SliceStable(labelPlugins[:], labelPlugins.Less)
	for _, labelPlugin := range labelPlugins {
		if enableMultiWorkers {
			labelPlugin.SetIgnorePath()
		}
		buf.WriteString(labelPlugin.String())
	}

	return strings.TrimRight(buf.String(), "\n"), nil
}

// +kubebuilder:object:generate:=false
type ByHashcode []params.PluginStore

// +kubebuilder:object:generate:=false
type ByRouteLabelsPointers []*params.PluginStore

// +kubebuilder:object:generate:=false
type ByRouteLabels []params.PluginStore

// +kubebuilder:object:generate:=false
type ByTags []params.PluginStore

func (a ByHashcode) Less(i, j int) bool            { return a[i].Hash() < a[j].Hash() }
func (a ByRouteLabelsPointers) Less(i, j int) bool { return a[i].RouteLabel() < a[j].RouteLabel() }
func (a ByRouteLabels) Less(i, j int) bool         { return a[i].RouteLabel() < a[j].RouteLabel() }
func (a ByTags) Less(i, j int) bool                { return a[i].GetTag() < a[j].GetTag() }

var _ Renderer = &FluentdConfig{}
var _ Renderer = &ClusterFluentdConfig{}
