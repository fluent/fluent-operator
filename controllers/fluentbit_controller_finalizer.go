package controllers

import (
	"context"

	loggingv1alpha2 "kubesphere.io/fluentbit-operator/apis/kubesphere.io/v1alpha2"
)

func (r *FluentBitReconciler) addFinalizer(ctx context.Context, instance *loggingv1alpha2.FluentBit) error {
	instance.AddFinalizer(loggingv1alpha2.FluentBitFinalizerName)
	return r.Update(ctx, instance)
}

func (r *FluentBitReconciler) handleFinalizer(ctx context.Context, instance *loggingv1alpha2.FluentBit) error {
	if !instance.HasFinalizer(loggingv1alpha2.FluentBitFinalizerName) {
		return nil
	}
	if err := r.delete(ctx, instance); err != nil {
		return err
	}
	instance.RemoveFinalizer(loggingv1alpha2.FluentBitFinalizerName)
	return r.Update(ctx, instance)
}
