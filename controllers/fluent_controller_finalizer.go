package controllers

import (
	"context"
	rbacv1 "k8s.io/api/rbac/v1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	fluentbitv1alpha2 "github.com/fluent/fluent-operator/v2/apis/fluentbit/v1alpha2"
	fluentdv1alpha1 "github.com/fluent/fluent-operator/v2/apis/fluentd/v1alpha1"
	"github.com/fluent/fluent-operator/v2/pkg/operator"
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

func (r *CollectorReconciler) addFinalizer(ctx context.Context, instance *fluentbitv1alpha2.Collector) error {
	instance.AddFinalizer(fluentbitv1alpha2.CollectorFinalizerName)
	return r.Update(ctx, instance)
}

func (r *CollectorReconciler) handleFinalizer(ctx context.Context, instance *fluentbitv1alpha2.Collector) error {
	if !instance.HasFinalizer(fluentbitv1alpha2.CollectorFinalizerName) {
		return nil
	}
	if err := r.delete(ctx, instance); err != nil {
		return err
	}
	instance.RemoveFinalizer(fluentbitv1alpha2.CollectorFinalizerName)
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
	sa := corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fd.Name,
			Namespace: fd.Namespace,
		},
	}
	if err := r.Delete(ctx, &sa); err != nil && !errors.IsNotFound(err) {
		return err
	}
	// TODO: clusterrole, clusterrolebinding

	sts := appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fd.Name,
			Namespace: fd.Namespace,
		},
	}
	if err := r.Delete(ctx, &sts); err != nil && !errors.IsNotFound(err) {
		return err
	}

	ds := appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fd.Name,
			Namespace: fd.Namespace,
		},
	}
	if err := r.Delete(ctx, &ds); err != nil && !errors.IsNotFound(err) {
		return err
	}

	svc := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fd.Name,
			Namespace: fd.Namespace,
		},
	}
	if err := r.Delete(ctx, &svc); err != nil && !errors.IsNotFound(err) {
		return err
	}

	return nil
}

func (r *FluentdReconciler) mutate(obj client.Object, fd *fluentdv1alpha1.Fluentd) controllerutil.MutateFn {
	switch o := obj.(type) {
	case *rbacv1.ClusterRole:
		expected, _, _ := operator.MakeRBACObjects(fd.Name, fd.Namespace, "fluentd", fd.Spec.RBACRules, fd.Spec.ServiceAccountAnnotations)

		return func() error {
			o.Rules = expected.Rules
			return nil
		}
	case *corev1.ServiceAccount:
		_, expected, _ := operator.MakeRBACObjects(fd.Name, fd.Namespace, "fluentd", fd.Spec.RBACRules, fd.Spec.ServiceAccountAnnotations)

		return func() error {
			o.Labels = expected.Labels
			o.Annotations = expected.Annotations
			if err := ctrl.SetControllerReference(fd, o, r.Scheme); err != nil {
				return err
			}
			return nil
		}
	case *rbacv1.ClusterRoleBinding:
		_, _, expected := operator.MakeRBACObjects(fd.Name, fd.Namespace, "fluentd", fd.Spec.RBACRules, fd.Spec.ServiceAccountAnnotations)

		return func() error {
			o.RoleRef = expected.RoleRef
			o.Subjects = expected.Subjects
			return nil
		}
	case *appsv1.StatefulSet:
		expected := operator.MakeStatefulSet(*fd)

		return func() error {
			o.Labels = expected.Labels
			o.Spec = expected.Spec
			if err := ctrl.SetControllerReference(fd, o, r.Scheme); err != nil {
				return err
			}
			return nil
		}
	case *appsv1.DaemonSet:
		expected := operator.MakeFluentdDaemonSet(*fd)
		return func() error {
			o.Labels = expected.Labels
			o.Spec = expected.Spec
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
			if err := ctrl.SetControllerReference(fd, o, r.Scheme); err != nil {
				return err
			}
			return nil
		}
	default:
	}

	return nil
}
