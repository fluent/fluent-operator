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

	fluentdv1alpha1 "fluent.io/fluent-operator/apis/fluentd/v1alpha1"
	"fluent.io/fluent-operator/pkg/operator"
)

// FluentdReconciler reconciles a Fluentd object
type FluentdReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=fluentd.fluent.io,resources=fluentds,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=fluentd.fluent.io,resources=fluentds/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=fluentd.fluent.io,resources=fluentds/finalizers,verbs=update

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
	var rbacObj, saObj, bindingObj = operator.MakeRBACObjects(fd.Name, fd.Namespace, "fluentd")

	// Set ServiceAccount's owner to this Fluentd
	if err := ctrl.SetControllerReference(&fd, saObj, r.Scheme); err != nil {
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

	// Deploy pvc
	// if not given a definition of buffer, will return default buffer configuration.
	if fd.Spec.BufferVolume == nil || (!fd.Spec.BufferVolume.DisableBufferVolume && fd.Spec.BufferVolume.PersistentVolumeClaim != nil) {
		bufferpvc := operator.MakeFluentdPVC(fd)
		if err := ctrl.SetControllerReference(&fd, &bufferpvc, r.Scheme); err != nil {
			return ctrl.Result{}, err
		}

		if _, err := controllerutil.CreateOrPatch(ctx, r.Client, &bufferpvc, r.mutate(&bufferpvc, &fd)); err != nil {
			return ctrl.Result{}, err
		}
	}

	// Deploy Fluentd Deployment
	dp := operator.MakeStatefulset(fd)
	if err := ctrl.SetControllerReference(&fd, &dp, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	if _, err := controllerutil.CreateOrPatch(ctx, r.Client, &dp, r.mutate(&dp, &fd)); err != nil {
		return ctrl.Result{}, err
	}

	// Deploy Fluentd Service
	if !fd.Spec.DisableService {
		svc := operator.MakeFluentdService(fd)
		if err := ctrl.SetControllerReference(&fd, &svc, r.Scheme); err != nil {
			return ctrl.Result{}, err
		}

		if _, err := controllerutil.CreateOrPatch(ctx, r.Client, &svc, r.mutate(&svc, &fd)); err != nil {
			return ctrl.Result{}, err
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
		Owns(&appsv1.StatefulSet{}).
		Owns(&corev1.Service{}).
		Owns(&corev1.PersistentVolumeClaim{}).
		Complete(r)
}
