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
	"bytes"
	"context"
	"crypto/md5"
	"fmt"
	"github.com/fluent/fluent-operator/apis/fluentbit/v1alpha2/plugins"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/source"
	"strings"

	fluentbitv1alpha2 "github.com/fluent/fluent-operator/apis/fluentbit/v1alpha2"
)

// FluentBitConfigReconciler reconciles a FluentBitConfig object
type FluentBitConfigReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=fluentbit.fluent.io,resources=clusterfluentbitconfigs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=fluentbit.fluent.io,resources=fluentbitconfigs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=fluentbit.fluent.io,resources=clusterinputs;clusterfilters;clusteroutputs;clusterparsers,verbs=list
// +kubebuilder:rbac:groups=fluentbit.fluent.io,resources=filters;outputs;parsers,verbs=list
// +kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch;create;update;patch;delete

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

	var fbs fluentbitv1alpha2.FluentBitList
	if err := r.List(ctx, &fbs); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	for _, fb := range fbs.Items {

		var cfgs fluentbitv1alpha2.ClusterFluentBitConfigList
		if err := r.List(ctx, &cfgs); err != nil {
			if errors.IsNotFound(err) {
				return ctrl.Result{}, nil
			}
			return ctrl.Result{}, err
		}

		for _, cfg := range cfgs.Items {
			// List all inputs matching the label selector.
			if cfg.Name != fb.Spec.FluentBitConfigName {
				continue
			}
			var inputs fluentbitv1alpha2.ClusterInputList
			selector, err := metav1.LabelSelectorAsSelector(&cfg.Spec.InputSelector)
			if err != nil {
				return ctrl.Result{}, err
			}
			if err = r.List(ctx, &inputs, client.MatchingLabelsSelector{Selector: selector}); err != nil {
				return ctrl.Result{}, err
			}

			// List all filters matching the label selector.
			var filters fluentbitv1alpha2.ClusterFilterList
			selector, err = metav1.LabelSelectorAsSelector(&cfg.Spec.FilterSelector)
			if err != nil {
				return ctrl.Result{}, err
			}
			if err = r.List(ctx, &filters, client.MatchingLabelsSelector{Selector: selector}); err != nil {
				return ctrl.Result{}, err
			}

			// List all outputs matching the label selector.
			var outputs fluentbitv1alpha2.ClusterOutputList
			selector, err = metav1.LabelSelectorAsSelector(&cfg.Spec.OutputSelector)
			if err != nil {
				return ctrl.Result{}, err
			}
			if err = r.List(ctx, &outputs, client.MatchingLabelsSelector{Selector: selector}); err != nil {
				return ctrl.Result{}, err
			}

			// List all parsers matching the label selector.
			var parsers fluentbitv1alpha2.ClusterParserList
			selector, err = metav1.LabelSelectorAsSelector(&cfg.Spec.ParserSelector)
			if err != nil {
				return ctrl.Result{}, err
			}
			if err = r.List(ctx, &parsers, client.MatchingLabelsSelector{Selector: selector}); err != nil {
				return ctrl.Result{}, err
			}

			// List all the namespace level resources if they exist and generate configs to mutate tags
			nsFilterLists, nsOutputLists, nsParserLists, nsClusterParserLists, rewriteTagConfigs, err := r.processNamespacedFluentBitCfgs(ctx, fb, inputs)

			if err != nil {
				return ctrl.Result{}, err
			}
			var ns string
			if cfg.Spec.Namespace != nil {
				ns = fmt.Sprintf(*cfg.Spec.Namespace)
			} else {
				ns = os.Getenv("NAMESPACE")
			}
			// Inject config data into Secret
			sl := plugins.NewSecretLoader(r.Client, ns, r.Log)
			mainAppCfg, err := cfg.RenderMainConfig(sl, inputs, filters, outputs, nsFilterLists, nsOutputLists, rewriteTagConfigs)
			if err != nil {
				return ctrl.Result{}, err
			}
			parserCfg, err := cfg.RenderParserConfig(sl, parsers, nsParserLists, nsClusterParserLists)
			if err != nil {
				return ctrl.Result{}, err
			}

			cl := plugins.NewConfigMapLoader(r.Client, ns)
			scripts, err := cfg.RenderLuaScript(cl, filters, ns)
			if err != nil {
				return ctrl.Result{}, err
			}

			// Create or update the corresponding Secret
			sec := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      cfg.Name,
					Namespace: ns,
				},
			}

			if err := ctrl.SetControllerReference(&cfg, sec, r.Scheme); err != nil {
				return ctrl.Result{}, err
			}

			if _, err := controllerutil.CreateOrPatch(ctx, r.Client, sec, func() error {
				sec.Data = map[string][]byte{
					"fluent-bit.conf": []byte(mainAppCfg),
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

			r.Log.Info("Fluent Bit main configuration has updated", "logging-control-plane", ns, "fluentbitconfig", cfg.Name, "secret", sec.Name)
		}
	}

	return ctrl.Result{}, nil
}

func (r *FluentBitConfigReconciler) processNamespacedFluentBitCfgs(ctx context.Context, fb fluentbitv1alpha2.FluentBit, inputs fluentbitv1alpha2.ClusterInputList) ([]fluentbitv1alpha2.FilterList, []fluentbitv1alpha2.OutputList,
	[]fluentbitv1alpha2.ParserList, []fluentbitv1alpha2.ClusterParserList, []string, error) {
	var nsCfgs fluentbitv1alpha2.FluentBitConfigList
	var filters []fluentbitv1alpha2.FilterList
	var outputs []fluentbitv1alpha2.OutputList
	var parsers []fluentbitv1alpha2.ParserList
	var clusterParsers []fluentbitv1alpha2.ClusterParserList
	var rewriteTagConfigs []string
	// set of rewrite_tag plugin configs to mutate tags for log records coming out of a namespace
	selector, err := metav1.LabelSelectorAsSelector(&fb.Spec.NamespacedFluentBitCfgSelector)
	if err != nil {
		return filters, outputs, parsers, clusterParsers, nil, err
	}

	if err := r.List(ctx, &nsCfgs, client.MatchingLabelsSelector{Selector: selector}); err != nil {
		return filters, outputs, parsers, clusterParsers, nil, err
	}

	// Form a slice of list of resources per namespace
	for _, cfg := range nsCfgs.Items {
		filterList, outputList, parserList, clusterParserList, err := r.ListNamespacedResources(ctx, cfg)
		if err != nil {
			return filters, outputs, parsers, clusterParsers, nil, err
		}
		filters = append(filters, filterList)
		outputs = append(outputs, outputList)
		parsers = append(parsers, parserList)
		clusterParsers = append(clusterParsers, clusterParserList)

		rewriteTagConfig := r.generateRewriteTagConfig(cfg, inputs)
		if rewriteTagConfig != "" {
			rewriteTagConfigs = append(rewriteTagConfigs, rewriteTagConfig)
		}
	}

	return filters, outputs, parsers, clusterParsers, rewriteTagConfigs, nil
}

func (r *FluentBitConfigReconciler) ListNamespacedResources(ctx context.Context, cfg fluentbitv1alpha2.FluentBitConfig) (fluentbitv1alpha2.FilterList,
	fluentbitv1alpha2.OutputList, fluentbitv1alpha2.ParserList, fluentbitv1alpha2.ClusterParserList, error) {
	var filters fluentbitv1alpha2.FilterList
	var outputs fluentbitv1alpha2.OutputList
	var parsers fluentbitv1alpha2.ParserList
	var clusterParsers fluentbitv1alpha2.ClusterParserList
	selector, err := metav1.LabelSelectorAsSelector(&cfg.Spec.FilterSelector)
	if err != nil {
		return filters, outputs, parsers, clusterParsers, err
	}
	if err := r.List(ctx, &filters, client.InNamespace(cfg.Namespace), client.MatchingLabelsSelector{Selector: selector}); err != nil {
		return filters, outputs, parsers, clusterParsers, err
	}

	selector, err = metav1.LabelSelectorAsSelector(&cfg.Spec.OutputSelector)
	if err != nil {
		return filters, outputs, parsers, clusterParsers, err
	}
	if err := r.List(ctx, &outputs, client.InNamespace(cfg.Namespace), client.MatchingLabelsSelector{Selector: selector}); err != nil {
		return filters, outputs, parsers, clusterParsers, err
	}

	selector, err = metav1.LabelSelectorAsSelector(&cfg.Spec.ParserSelector)
	if err != nil {
		return filters, outputs, parsers, clusterParsers, err
	}
	if err := r.List(ctx, &parsers, client.InNamespace(cfg.Namespace), client.MatchingLabelsSelector{Selector: selector}); err != nil {
		return filters, outputs, parsers, clusterParsers, err
	}

	selector, err = metav1.LabelSelectorAsSelector(&cfg.Spec.ClusterParserSelector)
	if err != nil {
		return filters, outputs, parsers, clusterParsers, err
	}
	if err := r.List(ctx, &clusterParsers, client.MatchingLabelsSelector{Selector: selector}); err != nil {
		return filters, outputs, parsers, clusterParsers, err
	}

	// Update the name of the local copies of cluster level parsers.
	// The intention is to have each namespace use their own copy
	// of the cluster parser.
	for i := range clusterParsers.Items {
		clusterParsers.Items[i].Name = fmt.Sprintf("%s-%x", clusterParsers.Items[i].Name, md5.Sum([]byte(cfg.Namespace)))
	}

	return filters, outputs, parsers, clusterParsers, nil
}

func (r *FluentBitConfigReconciler) generateRewriteTagConfig(cfg fluentbitv1alpha2.FluentBitConfig, inputs fluentbitv1alpha2.ClusterInputList) string {
	var tag string
	for _, input := range inputs.Items {
		if input.Spec.Tail == nil || !strings.Contains(input.Spec.Tail.Path, "/var/log/containers") {
			continue
		}
		tag = input.Spec.Tail.Tag
	}
	if tag == "" {
		return ""
	}
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("[Filter]\n"))
	buf.WriteString(fmt.Sprintf("    Name    rewrite_tag\n"))
	buf.WriteString(fmt.Sprintf("    Match    %s\n", tag))
	buf.WriteString(fmt.Sprintf("    Rule    $kubernetes['namespace_name'] ^(%s)$ %x.$TAG false\n", cfg.Namespace, md5.Sum([]byte(cfg.Namespace))))
	buf.WriteString(fmt.Sprintf("    Emitter_Name    re_emitted_%x\n", md5.Sum([]byte(cfg.Namespace))))
	return buf.String()
}

func (r *FluentBitConfigReconciler) SetupWithManager(mgr ctrl.Manager) error {
	if err := mgr.GetFieldIndexer().IndexField(context.Background(), &corev1.Secret{}, fluentbitOwnerKey, func(rawObj client.Object) []string {
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
	}); err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&fluentbitv1alpha2.FluentBit{}).
		Owns(&corev1.Secret{}).
		Watches(&source.Kind{Type: &fluentbitv1alpha2.ClusterFluentBitConfig{}}, &handler.EnqueueRequestForObject{}).
		Watches(&source.Kind{Type: &fluentbitv1alpha2.FluentBitConfig{}}, &handler.EnqueueRequestForObject{}).
		Watches(&source.Kind{Type: &fluentbitv1alpha2.ClusterInput{}}, &handler.EnqueueRequestForObject{}).
		Watches(&source.Kind{Type: &fluentbitv1alpha2.ClusterFilter{}}, &handler.EnqueueRequestForObject{}).
		Watches(&source.Kind{Type: &fluentbitv1alpha2.ClusterOutput{}}, &handler.EnqueueRequestForObject{}).
		Watches(&source.Kind{Type: &fluentbitv1alpha2.ClusterParser{}}, &handler.EnqueueRequestForObject{}).
		Watches(&source.Kind{Type: &fluentbitv1alpha2.Filter{}}, &handler.EnqueueRequestForObject{}).
		Watches(&source.Kind{Type: &fluentbitv1alpha2.Output{}}, &handler.EnqueueRequestForObject{}).
		Watches(&source.Kind{Type: &fluentbitv1alpha2.Parser{}}, &handler.EnqueueRequestForObject{}).
		Complete(r)
}
