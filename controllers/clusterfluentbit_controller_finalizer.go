package controllers

import (
	"context"
	clusterfluentbit "kubesphere.io/fluentbit-operator/api/clusterfluentbitoperator/v1alpha2"

)

func (r *ClusterFluentBitReconciler) addclusterFinalizer(ctx context.Context, instance *clusterfluentbit.ClusterFluentBit) error {
	instance.AddFinalizer(clusterfluentbit.ClusterFluentBitFinalizerName)
	return r.Update(ctx, instance)
}

func (r *ClusterFluentBitReconciler) handleClusterFinalizer(ctx context.Context, instance *clusterfluentbit.ClusterFluentBit) error {
	if !instance.HasFinalizer(clusterfluentbit.ClusterFluentBitFinalizerName) {
		return nil
	}
	if err := r.delete(ctx, instance); err != nil {
		return err
	}
	instance.RemoveFinalizer(clusterfluentbit.ClusterFluentBitFinalizerName)
	return r.Update(ctx, instance)
}
