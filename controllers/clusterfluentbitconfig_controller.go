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

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	fluentbitv1alpha2 "kubesphere.io/fluentbit-operator/apis/fluentbit.io/v1alpha2"
	"kubesphere.io/fluentbit-operator/apis/plugins"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// FluentBitConfigReconciler reconciles a FluentBitConfig object
type ClusterFluentBitConfigReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=fluentbit.fluent.io,resources=clusterfluentbitconfigs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=fluentbit.fluent.io,resources=clusterinputs;clusterfilters;clusteroutputs;clusterparsers,verbs=list
//+kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the FluentBitConfig object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *ClusterFluentBitConfigReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = r.Log.WithValues("clusterfluentbitconfig", req.NamespacedName)

	//var cfgs logging.FluentBitConfigList clusterfluentbit.ClusterFluentBitConfigList
	var cfgs fluentbitv1alpha2.ClusterFluentBitConfigList
	if err := r.List(ctx, &cfgs); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	for _, cfg := range cfgs.Items {
		// List all inputs matching the label selector.
		var clusterInputs fluentbitv1alpha2.ClusterInputList
		selector, err := metav1.LabelSelectorAsSelector(&cfg.Spec.ClusterInputSelector)
		if err != nil {
			return ctrl.Result{}, err
		}
		if err = r.List(ctx, &clusterInputs, client.MatchingLabelsSelector{Selector: selector}); err != nil {
			return ctrl.Result{}, err
		}

		// List all filters matching the label selector.
		var clusterFilters fluentbitv1alpha2.ClusterFilterList
		selector, err = metav1.LabelSelectorAsSelector(&cfg.Spec.ClusterFilterSelector)
		if err != nil {
			return ctrl.Result{}, err
		}
		if err = r.List(ctx, &clusterFilters, client.MatchingLabelsSelector{Selector: selector}); err != nil {
			return ctrl.Result{}, err
		}

		// List all outputs matching the label selector.
		var clusterOutputs fluentbitv1alpha2.ClusterOutputList
		selector, err = metav1.LabelSelectorAsSelector(&cfg.Spec.ClusterOutputSelector)
		if err != nil {
			return ctrl.Result{}, err
		}
		if err = r.List(ctx, &clusterOutputs, client.MatchingLabelsSelector{Selector: selector}); err != nil {
			return ctrl.Result{}, err
		}

		// List all parsers matching the label selector.
		var clusterParsers fluentbitv1alpha2.ClusterParserList
		selector, err = metav1.LabelSelectorAsSelector(&cfg.Spec.ClusterParserSelector)
		if err != nil {
			return ctrl.Result{}, err
		}
		if err = r.List(ctx, &clusterParsers, client.MatchingLabelsSelector{Selector: selector}); err != nil {
			return ctrl.Result{}, err
		}

		var nsList corev1.NamespaceList
		if err := r.Client.List(ctx, &nsList); err != nil {
			return ctrl.Result{}, fmt.Errorf("failed to receive a namespace list: %s", err)
		}
		for _, ns := range nsList.Items {
			// Inject config data into Secret
			sl := plugins.NewSecretLoader(r.Client, ns.Name, r.Log)
			mainCfg, err := cfg.RenderMainConfig(sl, clusterInputs, clusterFilters, clusterOutputs)
			if err != nil {
				return ctrl.Result{}, err
			}
			parserCfg, err := cfg.RenderParserConfig(sl, clusterParsers)
			if err != nil {
				return ctrl.Result{}, err
			}

			cl := plugins.NewConfigMapLoader(r.Client, ns.Name)
			scripts, err := cfg.RenderLuaScript(cl, clusterFilters)
			if err != nil {
				return ctrl.Result{}, err
			}

			// Create or update the corresponding Secret
			sec := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      cfg.Name,
					Namespace: ns.Name,
				},
			}

			if err := ctrl.SetControllerReference(&cfg, sec, r.Scheme); err != nil {
				return ctrl.Result{}, err
			}

			if _, err := controllerutil.CreateOrPatch(ctx, r.Client, sec, func() error {
				sec.Data = map[string][]byte{
					"fluent-bit.conf": []byte(mainCfg),
					"parsers.conf":    []byte(parserCfg),
				}
				for _, s := range scripts {
					sec.Data[s.Name] = []byte(s.Content)
				}
				sec.SetOwnerReferences(nil)
				if err := ctrl.SetControllerReference(&cfg, sec, r.Scheme); err != nil {
					return err
				}
				return nil
			}); err != nil {
				return ctrl.Result{}, err
			}
		}
	}
	return ctrl.Result{}, nil
}

func (r *ClusterFluentBitConfigReconciler) SetupWithManager(mgr ctrl.Manager) error {

	return ctrl.NewControllerManagedBy(mgr).
		For(&fluentbitv1alpha2.ClusterFluentBitConfig{}).
		Owns(&corev1.Secret{}).
		Watches(&source.Kind{Type: &fluentbitv1alpha2.ClusterInput{}}, &handler.EnqueueRequestForObject{}).
		Watches(&source.Kind{Type: &fluentbitv1alpha2.ClusterFilter{}}, &handler.EnqueueRequestForObject{}).
		Watches(&source.Kind{Type: &fluentbitv1alpha2.ClusterOutput{}}, &handler.EnqueueRequestForObject{}).
		Watches(&source.Kind{Type: &fluentbitv1alpha2.ClusterParser{}}, &handler.EnqueueRequestForObject{}).
		Complete(r)
}
