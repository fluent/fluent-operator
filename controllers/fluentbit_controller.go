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

	logging "kubesphere.io/fluentbit-operator/api/v1alpha2"
)

// FluentBitReconciler reconciles a FluentBit object
type FluentBitReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme

	ContainerLogRealPath string
}

// +kubebuilder:rbac:groups=logging.kubesphere.io,resources=fluentbits;fluentbitconfigs;inputs;filters;outputs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=daemonsets,verbs=get;create;update;patch;delete

func (r *FluentBitReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
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

	// check if Secret exists and requeue when not found
	var sec corev1.Secret
	if err := r.Get(ctx, client.ObjectKey{Namespace: fb.Namespace, Name: fb.Spec.FluentBitConfigName}, &sec); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{Requeue: true}, nil
		}
		return ctrl.Result{}, err
	}

	// install rbac resources for kubernetes filter plugin
	cr := rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: "fluent-bit",
		},
		Rules: []rbacv1.PolicyRule{
			{
				Verbs:     []string{"get"},
				APIGroups: []string{""},
				Resources: []string{"pods"},
			},
		},
	}
	if err := r.Create(ctx, &cr); err != nil && !errors.IsAlreadyExists(err) {
		return ctrl.Result{}, err
	}

	crb := rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: "fluent-bit",
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      rbacv1.ServiceAccountKind,
				Name:      fb.Name,
				Namespace: fb.Namespace,
			},
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: rbacv1.GroupName,
			Kind:     "ClusterRole",
			Name:     "fluent-bit",
		},
	}
	if err := r.Create(ctx, &crb); err != nil && !errors.IsAlreadyExists(err) {
		return ctrl.Result{}, err
	}

	// Deploy Fluent Bit DaemonSet
	var ds appsv1.DaemonSet
	err := r.Get(ctx, req.NamespacedName, &ds)
	if err == nil {
		return ctrl.Result{}, r.update(ctx, ds, fb)
	} else if errors.IsNotFound(err) {
		return ctrl.Result{}, r.create(ctx, fb)
	} else {
		return ctrl.Result{}, err
	}
}

func (r *FluentBitReconciler) update(ctx context.Context, raw appsv1.DaemonSet, fb logging.FluentBit) error {
	ds := r.constructDaemonSet(fb)
	raw.Labels = fb.Labels
	raw.Spec = ds.Spec
	if err := ctrl.SetControllerReference(&fb, &ds, r.Scheme); err != nil {
		return err
	}
	return r.Update(ctx, &ds)
}

func (r *FluentBitReconciler) create(ctx context.Context, fb logging.FluentBit) error {
	sa := corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fb.Name,
			Namespace: fb.Namespace,
		},
	}
	if err := ctrl.SetControllerReference(&fb, &sa, r.Scheme); err != nil {
		return err
	}
	if err := r.Create(ctx, &sa); err != nil && !errors.IsAlreadyExists(err) {
		return err
	}

	ds := r.constructDaemonSet(fb)
	if err := ctrl.SetControllerReference(&fb, &ds, r.Scheme); err != nil {
		return err
	}
	if err := r.Create(ctx, &ds); err != nil && !errors.IsAlreadyExists(err) {
		return err
	}
	return nil
}

func (r *FluentBitReconciler) constructDaemonSet(fb logging.FluentBit) appsv1.DaemonSet {
	logPath := func() string {
		if fb.Spec.ContainerLogRealPath != "" {
			return fb.Spec.ContainerLogRealPath
		} else if r.ContainerLogRealPath != "" {
			return r.ContainerLogRealPath
		} else {
			return "/var/lib/docker/containers"
		}
	}()

	ds := appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fb.Name,
			Namespace: fb.Namespace,
			Labels:    fb.Labels,
		},
		Spec: appsv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: fb.Labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:      fb.Name,
					Namespace: fb.Namespace,
					Labels:    fb.Labels,
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: fb.Name,
					Volumes: []corev1.Volume{
						{
							Name: "varlibcontainers",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: logPath,
								},
							},
						},
						{
							Name: "config",
							VolumeSource: corev1.VolumeSource{
								Secret: &corev1.SecretVolumeSource{
									SecretName: fb.Spec.FluentBitConfigName,
								},
							},
						},
						{
							Name: "varlogs",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "/var/log",
								},
							},
						},
						{
							Name:         "positions",
							VolumeSource: fb.Spec.PositionDB,
						},
					},
					Containers: []corev1.Container{
						{
							Name:            "fluent-bit",
							Image:           fb.Spec.Image,
							ImagePullPolicy: fb.Spec.ImagePullPolicy,
							Ports: []corev1.ContainerPort{
								{
									Name:          "metrics",
									ContainerPort: 2020,
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "varlibcontainers",
									ReadOnly:  true,
									MountPath: logPath,
								},
								{
									Name:      "config",
									ReadOnly:  true,
									MountPath: "/fluent-bit/config",
								},
								{
									Name:      "positions",
									MountPath: "/fluent-bit/tail",
								},
								{
									Name:      "varlogs",
									ReadOnly:  true,
									MountPath: "/var/log/",
								},
							},
						},
					},
					Tolerations: fb.Spec.Tolerations,
				},
			},
		},
	}

	// Mount Secrets
	for _, secret := range fb.Spec.Secrets {
		ds.Spec.Template.Spec.Volumes = append(ds.Spec.Template.Spec.Volumes, corev1.Volume{
			Name: secret,
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: secret,
				},
			},
		})
		ds.Spec.Template.Spec.Containers[0].VolumeMounts = append(ds.Spec.Template.Spec.Containers[0].VolumeMounts, corev1.VolumeMount{
			Name:      secret,
			ReadOnly:  true,
			MountPath: fmt.Sprintf("/fluent-bit/secrets/%s", secret),
		})
	}
	return ds
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
	if err := mgr.GetFieldIndexer().IndexField(&corev1.ServiceAccount{}, ownerKey, func(rawObj runtime.Object) []string {
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

	if err := mgr.GetFieldIndexer().IndexField(&appsv1.DaemonSet{}, ownerKey, func(rawObj runtime.Object) []string {
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
