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

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"kubesphere.io/fluentbit-operator/api/fluentbitoperator/v1alpha2/plugins"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/source"

	logging "kubesphere.io/fluentbit-operator/api/fluentbitoperator/v1alpha2"
)

// FluentBitConfigReconciler reconciles a FluentBitConfig object
type FluentBitConfigReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=logging.kubesphere.io,resources=fluentbitconfigs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=logging.kubesphere.io,resources=inputs;filters;outputs;parsers,verbs=list
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
func (r *FluentBitConfigReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
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

		// List all parsers matching the label selector.
		var parsers logging.ParserList
		selector, err = metav1.LabelSelectorAsSelector(&cfg.Spec.ParserSelector)
		if err != nil {
			return ctrl.Result{}, err
		}
		if err = r.List(ctx, &parsers, client.InNamespace(req.Namespace), client.MatchingLabelsSelector{Selector: selector}); err != nil {
			return ctrl.Result{}, err
		}

		// Inject config data into Secret
		sl := plugins.NewSecretLoader(r.Client, cfg.Namespace, r.Log)
		mainCfg, err := cfg.RenderMainConfig(sl, inputs, filters, outputs)
		if err != nil {
			return ctrl.Result{}, err
		}
		parserCfg, err := cfg.RenderParserConfig(sl, parsers)
		if err != nil {
			return ctrl.Result{}, err
		}

		cl := plugins.NewConfigMapLoader(r.Client, cfg.Namespace)
		scripts, err := cfg.RenderLuaScript(cl, filters)
		if err != nil {
			return ctrl.Result{}, err
		}

		// Create or update the corresponding Secret
		sec := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      cfg.Name,
				Namespace: cfg.Namespace,
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

	return ctrl.Result{}, nil
}

func (r *FluentBitConfigReconciler) SetupWithManager(mgr ctrl.Manager) error {
	if err := mgr.GetFieldIndexer().IndexField(context.Background(), &corev1.Secret{}, ownerKey, func(rawObj client.Object) []string {
		// Grab the job object, extract the owner.
		sec := rawObj.(*corev1.Secret)
		owner := metav1.GetControllerOf(sec)
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
		Owns(&corev1.Secret{}).
		Watches(&source.Kind{Type: &logging.Input{}}, &handler.EnqueueRequestForObject{}).
		Watches(&source.Kind{Type: &logging.Filter{}}, &handler.EnqueueRequestForObject{}).
		Watches(&source.Kind{Type: &logging.Output{}}, &handler.EnqueueRequestForObject{}).
		Watches(&source.Kind{Type: &logging.Parser{}}, &handler.EnqueueRequestForObject{}).
		Complete(r)
}
