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
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	fluentbitv1alpha2 "github.com/fluent/fluent-operator/v2/apis/fluentbit/v1alpha2"
	"github.com/fluent/fluent-operator/v2/pkg/operator"
)

// CollectorReconciler reconciles a FluentBit object
type CollectorReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=fluentbit.fluent.io,resources=fluentbits;fluentbitconfigs;collectors;inputs;filters;outputs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=serviceaccounts,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=clusterroles,verbs=create
// +kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=clusterrolebindings,verbs=create
// +kubebuilder:rbac:groups=core,resources=pods,verbs=get

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the FluentBit object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *CollectorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = r.Log.WithValues("collector", req.NamespacedName)

	var co fluentbitv1alpha2.Collector
	if err := r.Get(ctx, req.NamespacedName, &co); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	if co.IsBeingDeleted() {
		if err := r.handleFinalizer(ctx, &co); err != nil {
			return ctrl.Result{}, fmt.Errorf("error when handling finalizer: %v", err)
		}
		return ctrl.Result{}, nil
	}

	if !co.HasFinalizer(fluentbitv1alpha2.CollectorFinalizerName) {
		if err := r.addFinalizer(ctx, &co); err != nil {
			return ctrl.Result{}, fmt.Errorf("error when adding finalizer: %v", err)
		}
		return ctrl.Result{}, nil
	}

	// Check if Secret exists and requeue when not found
	var sec corev1.Secret
	if err := r.Get(ctx, client.ObjectKey{Namespace: co.Namespace, Name: co.Spec.FluentBitConfigName}, &sec); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{Requeue: true}, nil
		}
		return ctrl.Result{}, err
	}

	// Install RBAC resources for the filter plugin kubernetes
	cr, sa, crb := operator.MakeRBACObjects(co.Name, co.Namespace, "collector", co.Spec.RBACRules, co.Spec.ServiceAccountAnnotations)
	// Deploy Fluent Bit Collector ClusterRole
	if _, err := controllerutil.CreateOrPatch(ctx, r.Client, cr, r.mutate(cr, &co)); err != nil {
		return ctrl.Result{}, err
	}
	// Deploy Fluent Bit Collector ClusterRoleBinding
	if _, err := controllerutil.CreateOrPatch(ctx, r.Client, crb, r.mutate(crb, &co)); err != nil {
		return ctrl.Result{}, err
	}
	// Deploy Fluent Bit Collector ServiceAccount
	if _, err := controllerutil.CreateOrPatch(ctx, r.Client, sa, r.mutate(sa, &co)); err != nil {
		return ctrl.Result{}, err
	}

	// Deploy Fluent Bit Collector Statefulset
	sts := operator.MakefbStatefulset(co)
	if _, err := controllerutil.CreateOrPatch(ctx, r.Client, sts, r.mutate(sts, &co)); err != nil {
		return ctrl.Result{}, err
	}

	// Deploy Fluent Bit Collector Service
	if !co.Spec.DisableService {
		svc := operator.MakeCollectorService(co)
		if _, err := controllerutil.CreateOrPatch(ctx, r.Client, svc, r.mutate(svc, &co)); err != nil {
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

func (r *CollectorReconciler) mutate(obj client.Object, co *fluentbitv1alpha2.Collector) controllerutil.MutateFn {
	switch o := obj.(type) {
	case *rbacv1.ClusterRole:
		expected, _, _ := operator.MakeRBACObjects(co.Name, co.Namespace, "collector", co.Spec.RBACRules, co.Spec.ServiceAccountAnnotations)

		return func() error {
			o.Rules = expected.Rules
			return nil
		}
	case *corev1.ServiceAccount:
		_, expected, _ := operator.MakeRBACObjects(co.Name, co.Namespace, "collector", co.Spec.RBACRules, co.Spec.ServiceAccountAnnotations)

		return func() error {
			o.Annotations = expected.Annotations
			if err := ctrl.SetControllerReference(co, o, r.Scheme); err != nil {
				return err
			}
			return nil
		}
	case *rbacv1.ClusterRoleBinding:
		_, _, expected := operator.MakeRBACObjects(co.Name, co.Namespace, "collector", co.Spec.RBACRules, co.Spec.ServiceAccountAnnotations)

		return func() error {
			o.RoleRef = expected.RoleRef
			o.Subjects = expected.Subjects
			return nil
		}
	case *appsv1.StatefulSet:
		expected := operator.MakefbStatefulset(*co)

		return func() error {
			o.Labels = expected.Labels
			o.Spec = expected.Spec
			if err := ctrl.SetControllerReference(co, o, r.Scheme); err != nil {
				return err
			}
			return nil
		}
	case *corev1.Service:
		expected := operator.MakeCollectorService(*co)

		return func() error {
			o.Labels = expected.Labels
			o.Spec.Selector = expected.Spec.Selector
			o.Spec.Ports = expected.Spec.Ports
			if err := ctrl.SetControllerReference(co, o, r.Scheme); err != nil {
				return err
			}
			return nil
		}
	default:
	}

	return nil
}

func (r *CollectorReconciler) delete(ctx context.Context, co *fluentbitv1alpha2.Collector) error {
	sa := corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      co.Name,
			Namespace: co.Namespace,
		},
	}
	if err := r.Delete(ctx, &sa); err != nil && !errors.IsNotFound(err) {
		return err
	}
	// TODO: clusterrole, clusterrolebinding

	sts := appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      co.Name,
			Namespace: co.Namespace,
		},
	}
	if err := r.Delete(ctx, &sts); err != nil && !errors.IsNotFound(err) {
		return err
	}

	svcName := co.Name
	if co.Spec.Service.Name != "" {
		svcName = co.Spec.Service.Name
	}
	svc := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      svcName,
			Namespace: co.Namespace,
		},
	}
	if err := r.Delete(ctx, &svc); err != nil && !errors.IsNotFound(err) {
		return err
	}

	return nil
}

func (r *CollectorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	if err := mgr.GetFieldIndexer().IndexField(context.Background(), &corev1.ServiceAccount{}, collectorOwnerKey, func(rawObj client.Object) []string {
		// grab the job object, extract the owner.
		sa := rawObj.(*corev1.ServiceAccount)
		owner := metav1.GetControllerOf(sa)
		if owner == nil {
			return nil
		}
		// Make sure it's a FluentBit. If so, return it.
		if owner.APIVersion != fluentbitApiGVStr || owner.Kind != "Collector" {
			return nil
		}
		return []string{owner.Name}
	}); err != nil {
		return err
	}

	if err := mgr.GetFieldIndexer().IndexField(context.Background(), &appsv1.StatefulSet{}, collectorOwnerKey, func(rawObj client.Object) []string {
		// grab the job object, extract the owner.
		sts := rawObj.(*appsv1.StatefulSet)
		owner := metav1.GetControllerOf(sts)
		if owner == nil {
			return nil
		}
		// Make sure it's a FluentBit. If so, return it.
		if owner.APIVersion != fluentbitApiGVStr || owner.Kind != "Collector" {
			return nil
		}
		return []string{owner.Name}
	}); err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&fluentbitv1alpha2.Collector{}).
		Owns(&corev1.ServiceAccount{}).
		Owns(&appsv1.StatefulSet{}).
		Owns(&corev1.Service{}).
		Complete(r)
}
