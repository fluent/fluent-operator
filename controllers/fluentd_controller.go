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
	"time"

	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	fluentdv1alpha1 "github.com/fluent/fluent-operator/v2/apis/fluentd/v1alpha1"
	"github.com/fluent/fluent-operator/v2/pkg/operator"
)

// FluentdReconciler reconciles a Fluentd object
type FluentdReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=fluentd.fluent.io,resources=fluentds,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=fluentd.fluent.io,resources=fluentds/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=fluentd.fluent.io,resources=fluentds/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Fluentd object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *FluentdReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = r.Log.WithValues("fluentd", req.NamespacedName)

	var fd fluentdv1alpha1.Fluentd
	if err := r.Get(ctx, req.NamespacedName, &fd); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	if fd.IsBeingDeleted() {
		if err := r.handleFinalizer(ctx, &fd); err != nil {
			return ctrl.Result{}, fmt.Errorf("error when handling finalizer: %v", err)
		}
		return ctrl.Result{}, nil
	}

	if !fd.HasFinalizer(fluentdv1alpha1.FluentdFinalizerName) {
		if err := r.addFinalizer(ctx, &fd); err != nil {
			return ctrl.Result{}, fmt.Errorf("error when adding finalizer: %v", err)
		}
		return ctrl.Result{}, nil
	}

	// Check if Secret exists and requeue when not found
	var sec corev1.Secret
	secName := fmt.Sprintf("%s-config", fd.Name)
	if err := r.Get(ctx, client.ObjectKey{Namespace: fd.Namespace, Name: secName}, &sec); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{Requeue: true, RequeueAfter: time.Duration(time.Second)}, nil
		}
		return ctrl.Result{}, err
	}

	// Install RBAC resources for the filter plugin kubernetes
	cr, sa, crb := operator.MakeRBACObjects(fd.Name, fd.Namespace, "fluentd", fd.Spec.RBACRules, fd.Spec.ServiceAccountAnnotations)
	// Deploy Fluentd ClusterRole
	if _, err := controllerutil.CreateOrPatch(ctx, r.Client, cr, r.mutate(cr, &fd)); err != nil {
		return ctrl.Result{}, err
	}
	// Deploy Fluentd ClusterRoleBinding
	if _, err := controllerutil.CreateOrPatch(ctx, r.Client, crb, r.mutate(crb, &fd)); err != nil {
		return ctrl.Result{}, err
	}
	// Deploy Fluentd ServiceAccount
	if _, err := controllerutil.CreateOrPatch(ctx, r.Client, sa, r.mutate(sa, &fd)); err != nil {
		return ctrl.Result{}, err
	}

	var err error
	if fd.Spec.Mode == "agent" {
		// Deploy Fluentd DaemonSet
		ds := operator.MakeFluentdDaemonSet(fd)
		_, err = controllerutil.CreateOrPatch(ctx, r.Client, ds, r.mutate(ds, &fd))
		if err != nil {
			return ctrl.Result{}, err
		}
		sts := appsv1.StatefulSet{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fd.Name,
				Namespace: fd.Namespace,
			},
		}
		if err = r.Delete(ctx, &sts); err != nil && !errors.IsNotFound(err) {
			return ctrl.Result{}, err
		}
	} else {
		// Deploy Fluentd StatefulSet
		sts := operator.MakeStatefulSet(fd)
		_, err = controllerutil.CreateOrPatch(ctx, r.Client, sts, r.mutate(sts, &fd))
		ds := appsv1.DaemonSet{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fd.Name,
				Namespace: fd.Namespace,
			},
		}
		if err = r.Delete(ctx, &ds); err != nil && !errors.IsNotFound(err) {
			return ctrl.Result{}, err
		}
	}
	// Deploy Fluentd Service
	if !fd.Spec.DisableService {
		svc := operator.MakeFluentdService(fd)
		if len(svc.Spec.Ports) > 0 {
			if _, err = controllerutil.CreateOrPatch(ctx, r.Client, svc, r.mutate(svc, &fd)); err != nil {
				return ctrl.Result{}, err
			}
		}
	}

	return ctrl.Result{}, nil
}

func (r *FluentdReconciler) SetupWithManager(mgr ctrl.Manager) error {
	if err := mgr.GetFieldIndexer().IndexField(context.Background(), &corev1.Secret{}, fluentdOwnerKey, func(rawObj client.Object) []string {
		// Grab the job object, extract the owner.
		sec := rawObj.(*corev1.Secret)
		owner := metav1.GetControllerOf(sec)
		if owner == nil {
			return nil
		}

		if owner.APIVersion != fluentdApiGVStr || owner.Kind != "Fluentd" {
			return nil
		}
		return []string{owner.Name}
	}); err != nil {
		return err
	}

	if err := mgr.GetFieldIndexer().IndexField(context.Background(), &appsv1.DaemonSet{}, fluentdOwnerKey, func(rawObj client.Object) []string {
		// grab the job object, extract the owner.
		ds := rawObj.(*appsv1.DaemonSet)
		owner := metav1.GetControllerOf(ds)
		if owner == nil {
			return nil
		}

		if owner.APIVersion != fluentdApiGVStr || owner.Kind != "Fluentd" {
			return nil
		}
		return []string{owner.Name}
	}); err != nil {
		return err
	}

	if err := mgr.GetFieldIndexer().IndexField(context.Background(), &appsv1.StatefulSet{}, fluentdOwnerKey, func(rawObj client.Object) []string {
		// grab the job object, extract the owner.
		sts := rawObj.(*appsv1.StatefulSet)
		owner := metav1.GetControllerOf(sts)
		if owner == nil {
			return nil
		}

		if owner.APIVersion != fluentdApiGVStr || owner.Kind != "Fluentd" {
			return nil
		}
		return []string{owner.Name}
	}); err != nil {
		return err
	}

	if err := mgr.GetFieldIndexer().IndexField(context.Background(), &corev1.Service{}, fluentdOwnerKey, func(rawObj client.Object) []string {
		// grab the job object, extract the owner.
		svc := rawObj.(*corev1.Service)
		owner := metav1.GetControllerOf(svc)
		if owner == nil {
			return nil
		}

		if owner.APIVersion != fluentdApiGVStr || owner.Kind != "Fluentd" {
			return nil
		}
		return []string{owner.Name}
	}); err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&fluentdv1alpha1.Fluentd{}).
		Owns(&corev1.ServiceAccount{}).
		Owns(&appsv1.DaemonSet{}).
		Owns(&appsv1.StatefulSet{}).
		Owns(&corev1.Service{}).
		Complete(r)
}
