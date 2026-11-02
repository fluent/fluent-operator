package controllers

import (
	"context"

	rbacv1 "k8s.io/api/rbac/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// replaceStaleClusterRoleBinding deletes the existing ClusterRoleBinding when it refers to another role than
// expected, because the roleRef of a binding can't be changed. The following CreateOrPatch creates it again. This
// happens once after upgrading from versions that shared one ClusterRole between all instances of a component.
func replaceStaleClusterRoleBinding(ctx context.Context, c client.Client, expected *rbacv1.ClusterRoleBinding) error {
	var existing rbacv1.ClusterRoleBinding
	if err := c.Get(ctx, client.ObjectKey{Name: expected.Name}, &existing); err != nil {
		return client.IgnoreNotFound(err)
	}
	if existing.RoleRef == expected.RoleRef {
		return nil
	}
	return client.IgnoreNotFound(c.Delete(ctx, &existing))
}
