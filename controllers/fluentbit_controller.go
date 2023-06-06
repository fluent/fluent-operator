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

	rbacv1 "k8s.io/api/rbac/v1"

	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	fluentbitv1alpha2 "github.com/fluent/fluent-operator/v2/apis/fluentbit/v1alpha2"
	"github.com/fluent/fluent-operator/v2/pkg/operator"
)

// FluentBitReconciler reconciles a FluentBit object
type FluentBitReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme

	ContainerLogRealPath string
	Namespaced           bool
}

// +kubebuilder:rbac:groups=fluentbit.fluent.io,resources=fluentbits;fluentbitconfigs;inputs;filters;outputs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=daemonsets,verbs=get;list;watch;create;update;patch;delete
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
func (r *FluentBitReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = r.Log.WithValues("fluent-bit", req.NamespacedName)

	var fb fluentbitv1alpha2.FluentBit
	if err := r.Get(ctx, req.NamespacedName, &fb); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	if fb.IsBeingDeleted() {
		if err := r.handleFinalizer(ctx, &fb); err != nil {
			return ctrl.Result{}, fmt.Errorf("error when handling finalizer: %v", err)
		}
		return ctrl.Result{}, nil
	}

	if !fb.HasFinalizer(fluentbitv1alpha2.FluentBitFinalizerName) {
		if err := r.addFinalizer(ctx, &fb); err != nil {
			return ctrl.Result{}, fmt.Errorf("error when adding finalizer: %v", err)
		}
		return ctrl.Result{}, nil
	}

	// Check if Secret exists and requeue when not found
	var sec corev1.Secret
	if err := r.Get(ctx, client.ObjectKey{Namespace: fb.Namespace, Name: fb.Spec.FluentBitConfigName}, &sec); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{Requeue: true}, nil
		}
		return ctrl.Result{}, err
	}

	// Install RBAC resources for the filter plugin kubernetes
	var role, sa, binding client.Object
	if r.Namespaced {
		role, sa, binding = operator.MakeScopedRBACObjects(fb.Name, fb.Namespace, fb.Spec.ServiceAccountAnnotations)
	} else {
		role, sa, binding = operator.MakeRBACObjects(fb.Name, fb.Namespace, "fluent-bit", fb.Spec.RBACRules, fb.Spec.ServiceAccountAnnotations)
	}
	if _, err := controllerutil.CreateOrPatch(ctx, r.Client, role, r.mutate(role, &fb)); err != nil {
		return ctrl.Result{}, err
	}
	if _, err := controllerutil.CreateOrPatch(ctx, r.Client, sa, r.mutate(sa, &fb)); err != nil {
		return ctrl.Result{}, err
	}
	if _, err := controllerutil.CreateOrPatch(ctx, r.Client, binding, r.mutate(binding, &fb)); err != nil {
		return ctrl.Result{}, err
	}

	// Deploy Fluent Bit DaemonSet
	logPath := r.getContainerLogPath(fb)
	ds := operator.MakeDaemonSet(fb, logPath)
	if _, err := controllerutil.CreateOrPatch(ctx, r.Client, ds, r.mutate(ds, &fb)); err != nil {
		return ctrl.Result{}, err
	}

	// Deploy FluentBit Service
	if !fb.Spec.DisableService {
		svc := operator.MakeFluentbitService(fb)
		if _, err := controllerutil.CreateOrPatch(ctx, r.Client, svc, r.mutate(svc, &fb)); err != nil {
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

func (r *FluentBitReconciler) getContainerLogPath(fb fluentbitv1alpha2.FluentBit) string {
	if fb.Spec.ContainerLogRealPath != "" {
		return fb.Spec.ContainerLogRealPath
	} else if r.ContainerLogRealPath != "" {
		return r.ContainerLogRealPath
	} else {
		return "/var/lib/docker/containers"
	}
}

func (r *FluentBitReconciler) mutate(obj client.Object, fb *fluentbitv1alpha2.FluentBit) controllerutil.MutateFn {
	logPath := r.getContainerLogPath(*fb)

	switch o := obj.(type) {
	case *appsv1.DaemonSet:
		expected := operator.MakeDaemonSet(*fb, logPath)

		return func() error {
			o.Labels = expected.Labels
			o.Annotations = expected.Annotations
			o.Spec = expected.Spec
			if err := ctrl.SetControllerReference(fb, o, r.Scheme); err != nil {
				return err
			}
			return nil
		}

	case *corev1.Service:
		expected := operator.MakeFluentbitService(*fb)

		return func() error {
			o.Labels = expected.Labels
			o.Spec.Selector = expected.Spec.Selector
			o.Spec.Ports = expected.Spec.Ports
			if err := ctrl.SetControllerReference(fb, o, r.Scheme); err != nil {
				return err
			}
			return nil
		}
	case *rbacv1.Role:
		expected, _, _ := operator.MakeScopedRBACObjects(fb.Name, fb.Namespace, fb.Spec.ServiceAccountAnnotations)

		return func() error {
			o.Rules = expected.Rules
			if err := ctrl.SetControllerReference(fb, o, r.Scheme); err != nil {
				return err
			}
			return nil
		}
	case *rbacv1.ClusterRole:
		expected, _, _ := operator.MakeRBACObjects(fb.Name, fb.Namespace, "fluent-bit", fb.Spec.RBACRules, fb.Spec.ServiceAccountAnnotations)
		return func() error {
			o.Rules = expected.Rules
			return nil
		}
	case *corev1.ServiceAccount:
		_, expected, _ := operator.MakeScopedRBACObjects(fb.Name, fb.Namespace, fb.Spec.ServiceAccountAnnotations)

		return func() error {
			o.Annotations = expected.Annotations
			if err := ctrl.SetControllerReference(fb, o, r.Scheme); err != nil {
				return err
			}
			return nil
		}
	case *rbacv1.RoleBinding:
		_, _, expected := operator.MakeScopedRBACObjects(fb.Name, fb.Namespace, fb.Spec.ServiceAccountAnnotations)
		return func() error {
			o.Subjects = expected.Subjects
			o.RoleRef = expected.RoleRef
			if err := ctrl.SetControllerReference(fb, o, r.Scheme); err != nil {
				return err
			}
			return nil
		}
	case *rbacv1.ClusterRoleBinding:
		_, _, expected := operator.MakeRBACObjects(fb.Name, fb.Namespace, "fluent-bit", fb.Spec.RBACRules, fb.Spec.ServiceAccountAnnotations)
		return func() error {
			o.Subjects = expected.Subjects
			o.RoleRef = expected.RoleRef
			return nil
		}
	default:
	}

	return nil
}

func (r *FluentBitReconciler) delete(ctx context.Context, fb *fluentbitv1alpha2.FluentBit) error {
	sa := corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fb.Name,
			Namespace: fb.Namespace,
		},
	}
	if err := r.Delete(ctx, &sa); err != nil && !errors.IsNotFound(err) {
		return err
	}

	if r.Namespaced {
		roleName, _, roleBindingName := operator.MakeScopedRBACNames(fb.Name)
		role := rbacv1.Role{
			ObjectMeta: metav1.ObjectMeta{
				Name:      roleName,
				Namespace: fb.Namespace,
			},
		}
		if err := r.Delete(ctx, &role); err != nil && !errors.IsNotFound(err) {
			return err
		}

		rolebinding := rbacv1.RoleBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name:      roleBindingName,
				Namespace: fb.Namespace,
			},
		}
		if err := r.Delete(ctx, &rolebinding); err != nil && !errors.IsNotFound(err) {
			return err
		}
	}
	// TODO: clusterrole, clusterrolebinding

	ds := appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fb.Name,
			Namespace: fb.Namespace,
		},
	}
	if err := r.Delete(ctx, &ds); err != nil && !errors.IsNotFound(err) {
		return err
	}

	svc := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fb.Name,
			Namespace: fb.Namespace,
		},
	}
	if err := r.Delete(ctx, &svc); err != nil && !errors.IsNotFound(err) {
		return err
	}

	return nil
}

func (r *FluentBitReconciler) SetupWithManager(mgr ctrl.Manager) error {
	if err := mgr.GetFieldIndexer().IndexField(context.Background(), &corev1.ServiceAccount{}, fluentbitOwnerKey, func(rawObj client.Object) []string {
		// grab the job object, extract the owner.
		sa := rawObj.(*corev1.ServiceAccount)
		owner := metav1.GetControllerOf(sa)
		if owner == nil {
			return nil
		}
		// Make sure it's a FluentBit. If so, return it.
		if owner.APIVersion != fluentbitApiGVStr || owner.Kind != "FluentBit" {
			return nil
		}
		return []string{owner.Name}
	}); err != nil {
		return err
	}

	if err := mgr.GetFieldIndexer().IndexField(context.Background(), &appsv1.DaemonSet{}, fluentbitOwnerKey, func(rawObj client.Object) []string {
		// grab the job object, extract the owner.
		ds := rawObj.(*appsv1.DaemonSet)
		owner := metav1.GetControllerOf(ds)
		if owner == nil {
			return nil
		}
		// Make sure it's a FluentBit. If so, return it.
		if owner.APIVersion != fluentbitApiGVStr || owner.Kind != "FluentBit" {
			return nil
		}
		return []string{owner.Name}
	}); err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&fluentbitv1alpha2.FluentBit{}).
		Owns(&corev1.ServiceAccount{}).
		Owns(&appsv1.DaemonSet{}).
		Complete(r)
}
