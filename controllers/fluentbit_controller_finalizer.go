package controllers

import (
	"context"

	logging "kubesphere.io/fluentbit-operator/api/fluentbitoperator/v1alpha2"
)

func (r *FluentBitReconciler) addFinalizer(ctx context.Context, instance *logging.FluentBit) error {
	instance.AddFinalizer(logging.FluentBitFinalizerName)
	return r.Update(ctx, instance)
}

func (r *FluentBitReconciler) handleFinalizer(ctx context.Context, instance *logging.FluentBit) error {
	if !instance.HasFinalizer(logging.FluentBitFinalizerName) {
		return nil
	}
	if err := r.delete(ctx, instance); err != nil {
		return err
	}
	instance.RemoveFinalizer(logging.FluentBitFinalizerName)
	return r.Update(ctx, instance)
}
