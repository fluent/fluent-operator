package controllers

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	fluentbitv1alpha2 "fluent.io/fluent-operator/apis/fluentbit/v1alpha2"
	fluentdv1alpha1 "fluent.io/fluent-operator/apis/fluentd/v1alpha1"
	"fluent.io/fluent-operator/pkg/operator"
)

func (r *FluentBitReconciler) addFinalizer(ctx context.Context, instance *fluentbitv1alpha2.FluentBit) error {
	instance.AddFinalizer(fluentbitv1alpha2.FluentBitFinalizerName)
	return r.Update(ctx, instance)
}

func (r *FluentBitReconciler) handleFinalizer(ctx context.Context, instance *fluentbitv1alpha2.FluentBit) error {
	if !instance.HasFinalizer(fluentbitv1alpha2.FluentBitFinalizerName) {
		return nil
	}
	if err := r.delete(ctx, instance); err != nil {
		return err
	}
	instance.RemoveFinalizer(fluentbitv1alpha2.FluentBitFinalizerName)
	return r.Update(ctx, instance)
}

func (r *FluentdReconciler) addFinalizer(ctx context.Context, instance *fluentdv1alpha1.Fluentd) error {
	instance.AddFinalizer(fluentdv1alpha1.FluentdFinalizerName)
	return r.Update(ctx, instance)
}

func (r *FluentdReconciler) handleFinalizer(ctx context.Context, instance *fluentdv1alpha1.Fluentd) error {
	if !instance.HasFinalizer(fluentdv1alpha1.FluentdFinalizerName) {
		return nil
	}
	if err := r.delete(ctx, instance); err != nil {
		return err
	}
	instance.RemoveFinalizer(fluentdv1alpha1.FluentdFinalizerName)
	return r.Update(ctx, instance)
}

func (r *FluentdReconciler) delete(ctx context.Context, fd *fluentdv1alpha1.Fluentd) error {
	var sa corev1.ServiceAccount
	err := r.Get(ctx, client.ObjectKey{Namespace: fd.Namespace, Name: fd.Name}, &sa)
	if err == nil {
		if err := r.Delete(ctx, &sa); err != nil && !errors.IsNotFound(err) {
			return err
		}
	} else if !errors.IsNotFound(err) {
		return err
	}

	var svc corev1.Service
	err = r.Get(ctx, client.ObjectKey{Namespace: fd.Namespace, Name: fd.Name}, &svc)
	if err == nil {
		if err := r.Delete(ctx, &svc); err != nil && !errors.IsNotFound(err) {
			return err
		}
	} else if !errors.IsNotFound(err) {
		return err
	}

	var dp appsv1.Deployment
	err = r.Get(ctx, client.ObjectKey{Namespace: fd.Namespace, Name: fd.Name}, &dp)
	if err == nil {
		if err := r.Delete(ctx, &dp); err != nil && !errors.IsNotFound(err) {
			return err
		}
	} else if !errors.IsNotFound(err) {
		return err
	}

	var pvc corev1.PersistentVolumeClaim
	err = r.Get(ctx, client.ObjectKey{Namespace: fd.Namespace, Name: fmt.Sprintf("%s-buffer-pvc", fd.Name)}, &pvc)
	if err == nil {
		if err := r.Delete(ctx, &pvc); err != nil && !errors.IsNotFound(err) {
			return err
		}
	} else if !errors.IsNotFound(err) {
		return err
	}

	return nil
}

func (r *FluentdReconciler) mutate(obj client.Object, fd *fluentdv1alpha1.Fluentd) controllerutil.MutateFn {
	switch o := obj.(type) {
	case *corev1.PersistentVolumeClaim:
		expected := operator.MakeFluentdPVC(*fd)

		return func() error {
			o.Labels = expected.Labels
			o.Spec.AccessModes = expected.Spec.AccessModes
			o.Spec.Resources = expected.Spec.Resources
			o.Spec.VolumeMode = expected.Spec.VolumeMode
			o.SetOwnerReferences(nil)
			if err := ctrl.SetControllerReference(fd, o, r.Scheme); err != nil {
				return err
			}
			return nil
		}
	case *appsv1.Deployment:
		expected := operator.MakeDeployment(*fd)

		return func() error {
			o.Labels = expected.Labels
			o.Spec = expected.Spec
			o.SetOwnerReferences(nil)
			if err := ctrl.SetControllerReference(fd, o, r.Scheme); err != nil {
				return err
			}
			return nil
		}
	case *corev1.Service:
		expected := operator.MakeFluentdService(*fd)

		return func() error {
			o.Labels = expected.Labels
			o.Spec.Selector = expected.Spec.Selector
			o.Spec.Ports = expected.Spec.Ports
			o.SetOwnerReferences(nil)
			if err := ctrl.SetControllerReference(fd, o, r.Scheme); err != nil {
				return err
			}
			return nil
		}

	default:
	}

	return nil
}
