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
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	logging "kubesphere.io/fluentbit-operator/api/fluentbitoperator/v1alpha2"
	"kubesphere.io/fluentbit-operator/pkg/operator"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// FluentBitReconciler reconciles a FluentBit object
type FluentBitReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme

	ContainerLogRealPath string
	Namespaced           bool
}

//+kubebuilder:rbac:groups=logging.kubesphere.io,resources=fluentbits;fluentbitconfigs;inputs;filters;outputs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=daemonsets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=serviceaccounts,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=clusterroles,verbs=create
//+kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=clusterrolebindings,verbs=create
//+kubebuilder:rbac:groups=core,resources=pods,verbs=get

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

	var fb logging.FluentBit
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

	if !fb.HasFinalizer(logging.FluentBitFinalizerName) {
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
	var rbacObj, saObj, bindingObj client.Object
	if r.Namespaced {
		rbacObj, saObj, bindingObj = operator.MakeScopedRBACObjects(fb.Name, fb.Namespace)
	} else {
		rbacObj, saObj, bindingObj = operator.MakeRBACObjects(fb.Name, fb.Namespace)
	}
	// Set ServiceAccount's owner to this fluentbit
	if err := ctrl.SetControllerReference(&fb, saObj, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}
	if err := r.Create(ctx, rbacObj); err != nil && !errors.IsAlreadyExists(err) {
		return ctrl.Result{}, err
	}
	if err := r.Create(ctx, saObj); err != nil && !errors.IsAlreadyExists(err) {
		return ctrl.Result{}, err
	}
	if err := r.Create(ctx, bindingObj); err != nil && !errors.IsAlreadyExists(err) {
		return ctrl.Result{}, err
	}

	// Deploy Fluent Bit DaemonSet
	logPath := r.getContainerLogPath(fb)
	ds := operator.MakeDaemonSet(fb, logPath)
	if err := ctrl.SetControllerReference(&fb, &ds, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	if _, err := controllerutil.CreateOrPatch(ctx, r.Client, &ds, r.mutate(&ds, fb)); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *FluentBitReconciler) getContainerLogPath(fb logging.FluentBit) string {
	if fb.Spec.ContainerLogRealPath != "" {
		return fb.Spec.ContainerLogRealPath
	} else if r.ContainerLogRealPath != "" {
		return r.ContainerLogRealPath
	} else {
		return "/var/lib/docker/containers"
	}
}

func (r *FluentBitReconciler) mutate(ds *appsv1.DaemonSet, fb logging.FluentBit) controllerutil.MutateFn {
	logPath := r.getContainerLogPath(fb)
	expected := operator.MakeDaemonSet(fb, logPath)

	return func() error {
		ds.Labels = expected.Labels
		ds.Spec = expected.Spec
		ds.SetOwnerReferences(nil)
		if err := ctrl.SetControllerReference(&fb, ds, r.Scheme); err != nil {
			return err
		}
		return nil
	}
}

func (r *FluentBitReconciler) delete(ctx context.Context, fb *logging.FluentBit) error {
	var sa corev1.ServiceAccount
	err := r.Get(ctx, client.ObjectKey{Namespace: fb.Namespace, Name: fb.Name}, &sa)
	if err == nil {
		if err := r.Delete(ctx, &sa); err != nil && !errors.IsNotFound(err) {
			return err
		}
	} else if !errors.IsNotFound(err) {
		return err
	}

	var ds appsv1.DaemonSet
	err = r.Get(ctx, client.ObjectKey{Namespace: fb.Namespace, Name: fb.Name}, &ds)
	if err == nil {
		if err := r.Delete(ctx, &ds); err != nil && !errors.IsNotFound(err) {
			return err
		}
	} else if !errors.IsNotFound(err) {
		return err
	}

	return nil
}

func (r *FluentBitReconciler) SetupWithManager(mgr ctrl.Manager) error {
	if err := mgr.GetFieldIndexer().IndexField(context.Background(), &corev1.ServiceAccount{}, ownerKey, func(rawObj client.Object) []string {
		// grab the job object, extract the owner.
		sa := rawObj.(*corev1.ServiceAccount)
		owner := metav1.GetControllerOf(sa)
		if owner == nil {
			return nil
		}
		// Make sure it's a FluentBit. If so, return it.
		if owner.APIVersion != apiGVStr || owner.Kind != "FluentBit" {
			return nil
		}
		return []string{owner.Name}
	}); err != nil {
		return err
	}

	if err := mgr.GetFieldIndexer().IndexField(context.Background(), &appsv1.DaemonSet{}, ownerKey, func(rawObj client.Object) []string {
		// grab the job object, extract the owner.
		ds := rawObj.(*appsv1.DaemonSet)
		owner := metav1.GetControllerOf(ds)
		if owner == nil {
			return nil
		}
		// Make sure it's a FluentBit. If so, return it.
		if owner.APIVersion != apiGVStr || owner.Kind != "FluentBit" {
			return nil
		}
		return []string{owner.Name}
	}); err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&logging.FluentBit{}).
		Owns(&corev1.ServiceAccount{}).
		Owns(&appsv1.DaemonSet{}).
		Complete(r)
}
