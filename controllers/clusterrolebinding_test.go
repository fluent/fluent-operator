package controllers

import (
	"context"
	"testing"

	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

// The roleRef of a ClusterRoleBinding can't be changed, so a binding that still refers to the ClusterRole older
// versions shared between instances is deleted and created again for the instance's own role.
func TestReplaceStaleClusterRoleBinding(t *testing.T) {
	binding := func(role string) *rbacv1.ClusterRoleBinding {
		return &rbacv1.ClusterRoleBinding{
			ObjectMeta: metav1.ObjectMeta{Name: "fluent-operator-fluent-bit-events"},
			RoleRef:    rbacv1.RoleRef{APIGroup: rbacv1.GroupName, Kind: "ClusterRole", Name: role},
		}
	}
	expected := binding("fluent-operator-fluent-bit-events")
	key := client.ObjectKey{Name: expected.Name}

	for _, tc := range []struct {
		name        string
		existing    *rbacv1.ClusterRoleBinding
		wantDeleted bool
	}{
		{name: "no binding yet"},
		{name: "binding to the instance's role", existing: binding("fluent-operator-fluent-bit-events")},
		{name: "binding to the shared role", existing: binding("fluent-operator-fluent-bit"), wantDeleted: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			builder := fake.NewClientBuilder().WithScheme(scheme.Scheme)
			if tc.existing != nil {
				builder = builder.WithObjects(tc.existing)
			}
			c := builder.Build()

			if err := replaceStaleClusterRoleBinding(context.Background(), c, expected); err != nil {
				t.Fatalf("replaceStaleClusterRoleBinding returned error: %v", err)
			}

			err := c.Get(context.Background(), key, &rbacv1.ClusterRoleBinding{})
			switch {
			case tc.existing == nil || tc.wantDeleted:
				if !errors.IsNotFound(err) {
					t.Fatalf("binding still exists: %v", err)
				}
			case err != nil:
				t.Fatalf("binding to the expected role was deleted: %v", err)
			}
		})
	}
}
