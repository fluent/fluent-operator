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
	"fmt"
	"strings"
	"time"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/source"

	fluentdv1alpha1 "github.com/fluent/fluent-operator/v2/apis/fluentd/v1alpha1"
	"github.com/fluent/fluent-operator/v2/apis/fluentd/v1alpha1/plugins"
)

const (
	FluentdConfig        = "FluentdConfig"
	ClusterFluentdConfig = "ClusterFluentdConfig"

	FluentdSecretMainKey   = "fluent.conf"
	FluentdSecretSystemKey = "system.conf"
	FluentdSecretAppKey    = "app.conf"
	FluentdSecretLogKey    = "log.conf"

	FlUENT_INCLUDE = `# includes all files
@include /fluentd/etc/system.conf
@include /fluentd/etc/app.conf
@include /fluentd/etc/log.conf
`

	SYSTEM = `# Enable RPC endpoint
<system>
	rpc_endpoint 127.0.0.1:24444
	log_level %s
	workers %d
</system>
`
	FLUENTD_LOG = `# Do not collect fluentd's own logs to avoid infinite loops.
<match **>
	@type null
	@id main-no-output
</match>
<label @FLUENT_LOG>
	<match fluent.*>
		@type null
		@id main-fluentd-log
	</match>
</label>
`

	EMPTY_CFG = ``
)

// FluentdConfigReconciler reconciles a FluentdConfig object
type FluentdConfigReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=fluentd.fluent.io,resources=fluentdconfigs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=fluentd.fluent.io,resources=clusterfluentdconfigs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=fluentd.fluent.io,resources=inputs;filters;outputs,verbs=list
// +kubebuilder:rbac:groups=fluentd.fluent.io,resources=clusterinputs;clusterfilters;clusteroutputs,verbs=list;
// +kubebuilder:rbac:groups=fluentd.fluent.io,resources=fluentds,verbs=list
// +kubebuilder:rbac:groups=fluentd.fluent.io,resources=fluentds/status,verbs=patch
// +kubebuilder:rbac:groups=fluentd.fluent.io,resources=fluentdconfigs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=fluentd.fluent.io,resources=fluentdconfigs/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the FluentdConfig object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.

// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *FluentdConfigReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = r.Log.WithValues("fluentdconfig", req.NamespacedName)

	var fluentdList fluentdv1alpha1.FluentdList
	// List all fluentd instances to bind the generated runtime configuration to each fluentd
	if err := r.List(ctx, &fluentdList); err != nil {
		if errors.IsNotFound(err) {
			r.Log.Info("can not find fluentd CR definition.")
			return ctrl.Result{Requeue: true, RequeueAfter: time.Duration(1)}, nil
		}
		return ctrl.Result{}, err
	}

	// Loop over all Fluentd CRs
	for _, fd := range fluentdList.Items {
		// Get the config selector in this fluentd instance
		fdSelector, err := metav1.LabelSelectorAsSelector(&fd.Spec.FluentdCfgSelector)
		if err != nil {
			// Patch this fluentd instance if the selectors exit with errors
			if err := r.PatchObjects(ctx, &fd, fluentdv1alpha1.InactiveState, err.Error()); err != nil {
				return ctrl.Result{}, err
			}

			continue
		}
		// A secret loader supports LoadSecret method to parse the targeted secret
		sl := plugins.NewSecretLoader(r.Client, fd.Namespace, r.Log)

		// gpr acts as a global resource to store the related plugin resources
		gpr := fluentdv1alpha1.NewGlobalPluginResources("main")

		// Each cluster/namespace scope fluentd configs will generate their own filters/outputs plugins with their own cfgId/cfgLabel,
		// and they will finally be combined into one fluentd config file
		gpr.CombineGlobalInputsPlugins(sl, fd.Spec.GlobalInputs)

		// Default Output and filter
		// list all namespaced CRs
		inputs, filters, outputs, err := r.ListNamespacedLevelResources(ctx, fd.Namespace, fd.Spec.DefaultInputSelector, fd.Spec.DefaultFilterSelector, fd.Spec.DefaultOutputSelector)
		if err != nil {
			r.Log.Info("List namespace level resources failed", "config", "default", "err", err.Error())
			return ctrl.Result{}, err
		}
		if len(inputs) > 0 || len(filters) > 0 || len(outputs) > 0 {
			// Combine the namespaced filter/output pluginstores in this fluentd config
			cfgResouces, errs := gpr.PatchAndFilterNamespacedLevelResources(sl, fmt.Sprintf("%s-%s-%s", fd.Kind, fd.Namespace, fd.Name), inputs, filters, outputs)
			if len(errs) > 0 {
				r.Log.Info("Patch and filter namespace level resources failed", "config", "default", "err", strings.Join(errs, ","))
				return ctrl.Result{}, fmt.Errorf(strings.Join(errs, ","))
			}

			err = gpr.IdentifyCopyAndPatchOutput(cfgResouces)
			if err != nil {
				r.Log.Info("IdentifyCopy and PatchOutput namespace level resources failed", "config", "default", "err", strings.Join(errs, ","))
				return ctrl.Result{}, fmt.Errorf(strings.Join(errs, ","))
			}

			// WithCfgResources will collect all plugins to generate main config
			err = gpr.WithCfgResources("@default", cfgResouces)
			if err != nil {
				r.Log.Info("Combine resources failed", "config", "default", "err", err.Error())
				return ctrl.Result{}, err
			}
			// Add the default route to the main routing plugin
			gpr.MainRouterPlugins.InsertPairs("default_route", "@default")
		}

		// globalCfgLabels stores cfgLabels, the same cfg label is not allowed
		globalCfgLabels := make(map[string]bool)

		// Combine the resources matching the FluentdClusterConfigs selector into gpr
		var clustercfgs fluentdv1alpha1.ClusterFluentdConfigList
		// Use fluentd selector to match the cluster config.
		if err := r.List(ctx, &clustercfgs, client.MatchingLabelsSelector{Selector: fdSelector}); err != nil {
			if !errors.IsNotFound(err) {
				return ctrl.Result{}, err
			}
		}
		if err := r.ClusterCfgsForFluentd(ctx, clustercfgs, sl, gpr, globalCfgLabels); err != nil {
			return ctrl.Result{}, err
		}

		// Combine the resources matching the FluentdConfigs selector into gpr
		var cfgs fluentdv1alpha1.FluentdConfigList
		// Use fluentd selector to match the cluster config
		if err := r.List(ctx, &cfgs, client.MatchingLabelsSelector{Selector: fdSelector}); err != nil {
			if !errors.IsNotFound(err) {
				return ctrl.Result{}, err
			}
		}
		if err := r.CfgsForFluentd(ctx, cfgs, sl, gpr, globalCfgLabels); err != nil {
			return ctrl.Result{}, err
		}

		// Checks and returns the state of the matched cfgs
		state, msg := r.CheckAllState(cfgs, clustercfgs)
		if err := r.PatchObjects(ctx, &fd, state, msg); err != nil {
			return ctrl.Result{}, err
		}

		var mainAppCfg string
		var systemCfg string
		if state == fluentdv1alpha1.InactiveState {
			mainAppCfg = EMPTY_CFG
			systemCfg = EMPTY_CFG
		} else {
			// Get fluentd workers
			var workers int32 = 1
			if fd.Spec.Workers != nil {
				workers = *fd.Spec.Workers
			}

			// Create or update the secret of the fluentd instance in its namespace
			mainAppCfg, err = gpr.RenderMainConfig(bool(workers > 1))
			if err != nil {
				return ctrl.Result{}, err
			}

			systemCfg = fmt.Sprintf(SYSTEM, fd.Spec.LogLevel, workers)
		}

		secName := fmt.Sprintf("%s-config", fd.Name)

		sec := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      secName,
				Namespace: fd.Namespace,
			},
		}

		if _, err := controllerutil.CreateOrPatch(ctx, r.Client, sec, func() error {
			sec.Data = map[string][]byte{
				FluentdSecretMainKey:   []byte(FlUENT_INCLUDE),
				FluentdSecretAppKey:    []byte(mainAppCfg),
				FluentdSecretSystemKey: []byte(systemCfg),
				FluentdSecretLogKey:    []byte(FLUENTD_LOG),
			}
			// The current fd owns the namespaced secret
			sec.SetOwnerReferences(nil)
			if err := ctrl.SetControllerReference(&fd, sec, r.Scheme); err != nil {
				return err
			}
			return nil
		}); err != nil {
			return ctrl.Result{}, err
		}

		r.Log.Info("Fluentd main configuration has updated", "logging-control-plane", fd.Namespace, "fd", fd.Name, "secret", secName)
	}

	return ctrl.Result{}, nil
}

// ClusterCfgsForFluentd combines all cluster cfgs selected by this fd
func (r *FluentdConfigReconciler) ClusterCfgsForFluentd(
	ctx context.Context,
	clustercfgs fluentdv1alpha1.ClusterFluentdConfigList,
	sl plugins.SecretLoader,
	gpr *fluentdv1alpha1.PluginResources,
	globalCfgLabels map[string]bool,
) error {

	for _, cfg := range clustercfgs.Items {

		// Build the inner router for this cfg and append it to the MainRouter
		// Each cfg is a workflow.
		cfgRouter, err := gpr.BuildCfgRouter(&cfg)
		if err != nil {
			r.Log.Info("Build router failed", "config", cfg.Name, "err", err.Error())
			if err = r.PatchObjects(ctx, &cfg, fluentdv1alpha1.InvalidState, err.Error()); err != nil {
				return err
			}

			continue
		}
		// Set the route label which was calculated by a sub function of BuildCfgRoute
		cfgRouterLabel := fmt.Sprint(*cfgRouter.Label)

		if err := r.registerCfgLabel(cfgRouterLabel, globalCfgLabels); err != nil {
			r.Log.Info("Register fluentd config label failed", "config", cfg.Name, "err", err.Error())
			if err = r.PatchObjects(ctx, &cfg, fluentdv1alpha1.InvalidState, err.Error()); err != nil {
				return err
			}

			continue
		}

		// list all cluster CRs
		clusterInputs, clusterfilters, clusteroutputs, err := r.ListClusterLevelResources(ctx, cfg.Spec.ClusterInputSelector, cfg.Spec.ClusterFilterSelector, cfg.Spec.ClusterOutputSelector)
		if err != nil {
			r.Log.Info("List cluster level resources failed", "config", cfg.Name, "err", err.Error())
			if err = r.PatchObjects(ctx, &cfg, fluentdv1alpha1.InvalidState, err.Error()); err != nil {
				return err
			}

			continue
		}

		// Combine the filter/output pluginstores in this fluentd config
		cfgResouces, errs := gpr.PatchAndFilterClusterLevelResources(sl, cfg.GetCfgId(), clusterInputs, clusterfilters, clusteroutputs)
		if len(errs) > 0 {
			r.Log.Info("Patch and filter cluster level resources failed", "config", cfg.Name, "err", strings.Join(errs, ","))
			if err = r.PatchObjects(ctx, &cfg, fluentdv1alpha1.InvalidState, strings.Join(errs, ", ")); err != nil {
				return err
			}

			continue
		}

		err = gpr.IdentifyCopyAndPatchOutput(cfgResouces)
		if err != nil {
			r.Log.Info("Patch and filter cluster level resources failed", "config", cfg.Name, "err", strings.Join(errs, ","))
			if err = r.PatchObjects(ctx, &cfg, fluentdv1alpha1.InvalidState, strings.Join(errs, ", ")); err != nil {
				return err
			}
			continue
		}

		// WithCfgResources will collect all plugins to generate main config
		var msg string
		err = gpr.WithCfgResources(cfgRouterLabel, cfgResouces)
		if err != nil {
			r.Log.Info("Combine resources failed", "config", cfg.Name, "err", err.Error())
			msg = err.Error()
		} else {
			msg = "Generate fluentd configs successfully"
		}

		if err = r.PatchObjects(ctx, &cfg, fluentdv1alpha1.ValidState, msg); err != nil {
			return err
		}
	}

	return nil
}

// CfgsForFluentd combines all namespaced cfgs selected by this fd
func (r *FluentdConfigReconciler) CfgsForFluentd(ctx context.Context, cfgs fluentdv1alpha1.FluentdConfigList, sl plugins.SecretLoader,
	gpr *fluentdv1alpha1.PluginResources, globalCfgLabels map[string]bool) error {

	for _, cfg := range cfgs.Items {
		// Build the inner router for this cfg and append it to the MainRouter
		cfgRouter, err := gpr.BuildCfgRouter(&cfg)
		if err != nil {
			r.Log.Info("Build router failed", "config", cfg.Name, "err", err.Error())
			if err = r.PatchObjects(ctx, &cfg, fluentdv1alpha1.InvalidState, err.Error()); err != nil {
				return err
			}

			continue
		}

		// register routeLabel and fail if it is already present
		cfgRouterLabel := fmt.Sprint(*cfgRouter.Label)
		if err := r.registerCfgLabel(cfgRouterLabel, globalCfgLabels); err != nil {
			r.Log.Info("Register fluentd config label failed", "config", cfg.Name, "err", err.Error())
			if err = r.PatchObjects(ctx, &cfg, fluentdv1alpha1.InvalidState, err.Error()); err != nil {
				return err
			}

			continue
		}

		// list all cluster CRs
		clusterInputs, clusterfilters, clusteroutputs, err := r.ListClusterLevelResources(ctx, cfg.Spec.ClusterInputSelector, cfg.Spec.ClusterFilterSelector, cfg.Spec.ClusterOutputSelector)
		if err != nil {
			r.Log.Info("List cluster level resources failed", "config", cfg.Name, "err", err.Error())
			if err = r.PatchObjects(ctx, &cfg, fluentdv1alpha1.InvalidState, err.Error()); err != nil {
				return err
			}

			continue
		}

		// list all namespaced CRs
		inputs, filters, outputs, err := r.ListNamespacedLevelResources(ctx, cfg.Namespace, cfg.Spec.InputSelector, cfg.Spec.FilterSelector, cfg.Spec.OutputSelector)
		if err != nil {
			r.Log.Info("List namespace level resources failed", "config", cfg.Name, "err", err.Error())
			if err = r.PatchObjects(ctx, &cfg, fluentdv1alpha1.InvalidState, err.Error()); err != nil {
				return err
			}

			continue
		}

		// Combine the cluster input/filter/output pluginstores in this fluentd config
		clustercfgResouces, errs := gpr.PatchAndFilterClusterLevelResources(sl, cfg.GetCfgId(), clusterInputs, clusterfilters, clusteroutputs)
		if len(errs) > 0 {
			r.Log.Info("Patch and filter cluster level resources failed", "config", cfg.Name, "err", strings.Join(errs, ","))
			if err = r.PatchObjects(ctx, &cfg, fluentdv1alpha1.InvalidState, strings.Join(errs, ", ")); err != nil {
				return err
			}

			continue
		}

		// Combine the namespaced filter/output pluginstores in this fluentd config
		cfgResouces, errs := gpr.PatchAndFilterNamespacedLevelResources(sl, cfg.GetCfgId(), inputs, filters, outputs)
		if len(errs) > 0 {
			r.Log.Info("Patch and filter namespace level resources failed", "config", cfg.Name, "err", strings.Join(errs, ","))
			if err = r.PatchObjects(ctx, &cfg, fluentdv1alpha1.InvalidState, strings.Join(errs, ", ")); err != nil {
				return err
			}

			continue
		}

		cfgResouces.InputPlugins = append(cfgResouces.InputPlugins, clustercfgResouces.InputPlugins...)
		cfgResouces.FilterPlugins = append(cfgResouces.FilterPlugins, clustercfgResouces.FilterPlugins...)
		cfgResouces.OutputPlugins = append(cfgResouces.OutputPlugins, clustercfgResouces.OutputPlugins...)

		err = gpr.IdentifyCopyAndPatchOutput(cfgResouces)
		if err != nil {
			r.Log.Info("Patch and filter namespace level resources failed", "config", cfg.Name, "err", strings.Join(errs, ","))
			if err = r.PatchObjects(ctx, &cfg, fluentdv1alpha1.InvalidState, strings.Join(errs, ", ")); err != nil {
				return err
			}
			continue
		}

		// WithCfgResources will collect all plugins to generate main config
		var msg string
		err = gpr.WithCfgResources(cfgRouterLabel, cfgResouces)
		if err != nil {
			r.Log.Info("Combine resources failed", "config", cfg.Name, "err", err.Error())
			msg = err.Error()
		} else {
			msg = "Generate fluentd configs successfully"
		}

		if err = r.PatchObjects(ctx, &cfg, fluentdv1alpha1.ValidState, msg); err != nil {
			return err
		}
	}

	return nil
}

func (r *FluentdConfigReconciler) CheckAllState(
	cfgs fluentdv1alpha1.FluentdConfigList,
	clustercfgs fluentdv1alpha1.ClusterFluentdConfigList,
) (fluentdv1alpha1.StatusState, string) {
	invalidCfgIds := make([]string, 0, len(cfgs.Items)+len(cfgs.Items))

	for _, cfg := range cfgs.Items {
		if cfg.Status.State == fluentdv1alpha1.InvalidState {
			invalidCfgIds = append(invalidCfgIds, cfg.GetCfgId())
		}
	}

	for _, cfg := range clustercfgs.Items {
		if cfg.Status.State == fluentdv1alpha1.InvalidState {
			invalidCfgIds = append(invalidCfgIds, cfg.GetCfgId())
		}
	}

	if len(invalidCfgIds) == 0 {
		return fluentdv1alpha1.ActiveState, "all matched cfgs is valid"
	}

	if len(invalidCfgIds) < cap(invalidCfgIds) {
		return fluentdv1alpha1.ActiveState, "part of the cfgs are invalid. Invalid cfgs: " + strings.Join(invalidCfgIds, ", ")
	}

	return fluentdv1alpha1.InactiveState, "all matched cfgs is invalid. Invalid cfgs: " + strings.Join(invalidCfgIds, ", ")
}

// registerCfgLabel registers a cfglabel for this clustercfg/cfg
func (r *FluentdConfigReconciler) registerCfgLabel(cfgLabel string, globalCfgLabels map[string]bool) error {
	// cfgRouterLabel contains the important information for this cfg.
	// check if the calculated cfgLabel is already present
	if ok := globalCfgLabels[cfgLabel]; ok {
		return fmt.Errorf("the current configuration already exists: %s", cfgLabel)
	}

	// register the cfg labels, the same cfg labels is not allowed
	globalCfgLabels[cfgLabel] = true
	return nil
}

func (r *FluentdConfigReconciler) ListClusterLevelResources(
	ctx context.Context,
	inputSelector, filterSelector, outputSelector *metav1.LabelSelector,
) ([]fluentdv1alpha1.ClusterInput, []fluentdv1alpha1.ClusterFilter, []fluentdv1alpha1.ClusterOutput, error) {
	// List all inputs matching the label selector
	var clusterInputs fluentdv1alpha1.ClusterInputList
	if inputSelector != nil {
		selector, err := metav1.LabelSelectorAsSelector(inputSelector)
		if err != nil {
			return nil, nil, nil, err
		}
		if err = r.List(ctx, &clusterInputs, client.MatchingLabelsSelector{Selector: selector}); err != nil {
			return nil, nil, nil, err
		}
	}
	// List all filters matching the label selector
	var clusterfilters fluentdv1alpha1.ClusterFilterList
	if filterSelector != nil {
		selector, err := metav1.LabelSelectorAsSelector(filterSelector)
		if err != nil {
			return nil, nil, nil, err
		}
		if err = r.List(ctx, &clusterfilters, client.MatchingLabelsSelector{Selector: selector}); err != nil {
			return nil, nil, nil, err
		}
	}

	// List all outputs matching the label selector
	var clusteroutputs fluentdv1alpha1.ClusterOutputList
	if outputSelector != nil {
		selector, err := metav1.LabelSelectorAsSelector(outputSelector)
		if err != nil {
			return nil, nil, nil, err
		}
		if err = r.List(ctx, &clusteroutputs, client.MatchingLabelsSelector{Selector: selector}); err != nil {
			return nil, nil, nil, err
		}
	}

	return clusterInputs.Items, clusterfilters.Items, clusteroutputs.Items, nil
}

func (r *FluentdConfigReconciler) ListNamespacedLevelResources(
	ctx context.Context,
	namespace string,
	inputSelector, filterSelector, outputSelector *metav1.LabelSelector,
) ([]fluentdv1alpha1.Input, []fluentdv1alpha1.Filter, []fluentdv1alpha1.Output, error) {
	// List all inputs matching the label selector
	var inputs fluentdv1alpha1.InputList
	if inputSelector != nil {
		selector, err := metav1.LabelSelectorAsSelector(inputSelector)
		if err != nil {
			return nil, nil, nil, err
		}
		if err = r.List(ctx, &inputs, client.InNamespace(namespace), client.MatchingLabelsSelector{Selector: selector}); err != nil {
			return nil, nil, nil, err
		}
	}

	// List and patch the related cluster CRs
	var filters fluentdv1alpha1.FilterList
	if filterSelector != nil {
		selector, err := metav1.LabelSelectorAsSelector(filterSelector)
		if err != nil {
			return nil, nil, nil, err
		}
		if err = r.List(ctx, &filters, client.InNamespace(namespace), client.MatchingLabelsSelector{Selector: selector}); err != nil {
			return nil, nil, nil, err
		}
	}

	// List all outputs matching the label selector
	var outputs fluentdv1alpha1.OutputList
	if outputSelector != nil {
		selector, err := metav1.LabelSelectorAsSelector(outputSelector)
		if err != nil {
			return nil, nil, nil, err
		}
		if err = r.List(ctx, &outputs, client.InNamespace(namespace), client.MatchingLabelsSelector{Selector: selector}); err != nil {
			return nil, nil, nil, err
		}
	}

	return inputs.Items, filters.Items, outputs.Items, nil
}

// PatchObjects patches the errors to the obj
func (r *FluentdConfigReconciler) PatchObjects(ctx context.Context, obj client.Object, state fluentdv1alpha1.StatusState, msg string) error {
	switch o := obj.(type) {
	case *fluentdv1alpha1.ClusterFluentdConfig:
		o.Status.State = state
		o.Status.Messages = msg
		if err := r.Status().Update(ctx, o); err != nil {
			return err
		}
	case *fluentdv1alpha1.FluentdConfig:
		o.Status.State = state
		o.Status.Messages = msg
		if err := r.Status().Update(ctx, o); err != nil {
			return err
		}
	case *fluentdv1alpha1.Fluentd:
		o.Status.State = state
		o.Status.Messages = msg
		if err := r.Status().Update(ctx, o); err != nil {
			return err
		}
	default:
	}

	return nil
}

// SetupWithManager sets up the controller with the Manager
func (r *FluentdConfigReconciler) SetupWithManager(mgr ctrl.Manager) error {
	if err := mgr.GetFieldIndexer().IndexField(context.Background(), &corev1.ServiceAccount{}, fluentdOwnerKey, func(rawObj client.Object) []string {
		// grab the job object, extract the owner
		sa := rawObj.(*corev1.ServiceAccount)
		owner := metav1.GetControllerOf(sa)
		if owner == nil {
			return nil
		}
		// Make sure it's a Fluentd. If so, return it
		if owner.APIVersion != fluentdApiGVStr || owner.Kind != "Fluentd" {
			return nil
		}
		return []string{owner.Name}
	}); err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&fluentdv1alpha1.Fluentd{}).
		Owns(&corev1.Secret{}).
		Watches(&source.Kind{Type: &fluentdv1alpha1.ClusterFluentdConfig{}}, &handler.EnqueueRequestForObject{}).
		Watches(&source.Kind{Type: &fluentdv1alpha1.FluentdConfig{}}, &handler.EnqueueRequestForObject{}).
		Watches(&source.Kind{Type: &fluentdv1alpha1.Filter{}}, &handler.EnqueueRequestForObject{}).
		Watches(&source.Kind{Type: &fluentdv1alpha1.ClusterFilter{}}, &handler.EnqueueRequestForObject{}).
		Watches(&source.Kind{Type: &fluentdv1alpha1.Output{}}, &handler.EnqueueRequestForObject{}).
		Watches(&source.Kind{Type: &fluentdv1alpha1.ClusterOutput{}}, &handler.EnqueueRequestForObject{}).
		Complete(r)
}
