/*

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
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"kubesphere.io/fluentbit-operator/api/v1alpha1/plugins"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/source"

	logging "kubesphere.io/fluentbit-operator/api/v1alpha1"
)

// FluentBitConfigReconciler reconciles a FluentBitConfig object
type FluentBitConfigReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=logging.kubesphere.io,resources=fluentbitconfigs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=logging.kubesphere.io,resources=inputs;filters;outputs,verbs=list
// +kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;create;update;patch;delete

func (r *FluentBitConfigReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	_ = r.Log.WithValues("fluentbitconfig", req.NamespacedName)

	var cfgs logging.FluentBitConfigList
	if err := r.List(ctx, &cfgs, client.InNamespace(req.Namespace)); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	for _, cfg := range cfgs.Items {
		// List all inputs matching the label selector.
		var inputs logging.InputList
		selector, err := metav1.LabelSelectorAsSelector(&cfg.Spec.InputSelector)
		if err != nil {
			return ctrl.Result{}, err
		}
		if err = r.List(ctx, &inputs, client.InNamespace(req.Namespace), client.MatchingLabelsSelector{Selector: selector}); err != nil {
			return ctrl.Result{}, err
		}

		// List all filters matching the label selector.
		var filters logging.FilterList
		selector, err = metav1.LabelSelectorAsSelector(&cfg.Spec.FilterSelector)
		if err != nil {
			return ctrl.Result{}, err
		}
		if err = r.List(ctx, &filters, client.InNamespace(req.Namespace), client.MatchingLabelsSelector{Selector: selector}); err != nil {
			return ctrl.Result{}, err
		}

		// List all outputs matching the label selector.
		var outputs logging.OutputList
		selector, err = metav1.LabelSelectorAsSelector(&cfg.Spec.OutputSelector)
		if err != nil {
			return ctrl.Result{}, err
		}
		if err = r.List(ctx, &outputs, client.InNamespace(req.Namespace), client.MatchingLabelsSelector{Selector: selector}); err != nil {
			return ctrl.Result{}, err
		}

		sl := plugins.NewSecretLoader(r.Client, cfg.Namespace, r.Log)
		data, err := cfg.Render(sl, inputs, filters, outputs)
		if err != nil {
			return ctrl.Result{}, err
		}

		// Create or update the corresponding ConfigMap
		var cm corev1.ConfigMap
		if err := r.Get(ctx, types.NamespacedName{Namespace: cfg.Namespace, Name: cfg.Name}, &cm); err == nil {
			cm.Data = map[string]string{"fluent-bit.conf": data}
			if err := r.Update(ctx, &cm); err != nil {
				// resourceVersion conflict
				if errors.IsConflict(err) {
					return ctrl.Result{Requeue: true}, nil
				}
				return ctrl.Result{}, err
			}
		} else if errors.IsNotFound(err) {
			cm = corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      req.Name,
					Namespace: req.Namespace,
				},
				Data: map[string]string{"fluent-bit.conf": data},
			}
			// Set ConfigMap's owner to FluentBitConfig
			if err := ctrl.SetControllerReference(&cfg, &cm, r.Scheme); err != nil {
				return ctrl.Result{}, err
			}
			// Create ConfigMap
			return ctrl.Result{}, r.Create(ctx, &cm)
		} else {
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

func (r *FluentBitConfigReconciler) SetupWithManager(mgr ctrl.Manager) error {
	if err := mgr.GetFieldIndexer().IndexField(&corev1.ConfigMap{}, ownerKey, func(rawObj runtime.Object) []string {
		// Grab the job object, extract the owner.
		cm := rawObj.(*corev1.ConfigMap)
		owner := metav1.GetControllerOf(cm)
		if owner == nil {
			return nil
		}
		// Make sure it's a FluentBitConfig. If so, return it.
		if owner.APIVersion != apiGVStr || owner.Kind != "FluentBitConfig" {
			return nil
		}
		return []string{owner.Name}
	}); err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&logging.FluentBitConfig{}).
		Owns(&corev1.ConfigMap{}).
		Watches(&source.Kind{Type: &logging.Input{}}, &handler.EnqueueRequestForObject{}).
		Watches(&source.Kind{Type: &logging.Filter{}}, &handler.EnqueueRequestForObject{}).
		Watches(&source.Kind{Type: &logging.Output{}}, &handler.EnqueueRequestForObject{}).
		Complete(r)
}
