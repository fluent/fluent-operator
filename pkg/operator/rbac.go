package operator

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func MakeRBACObjects(
	name, namespace, component string,
	additionalRules []rbacv1.PolicyRule,
	saAnnotations map[string]string,
) (*rbacv1.ClusterRole, *corev1.ServiceAccount, *rbacv1.ClusterRoleBinding) {
	crName, saName, crbName := MakeRBACNames(name, component)
	cr := rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: crName,
		},
		Rules: []rbacv1.PolicyRule{
			{
				Verbs:     []string{"get"},
				APIGroups: []string{""},
				Resources: []string{"pods"},
			},
		},
	}

	if additionalRules != nil {
		cr.Rules = append(cr.Rules, additionalRules...)
	}

	sa := corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:        saName,
			Namespace:   namespace,
			Annotations: saAnnotations,
		},
	}

	crb := rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: crbName,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      rbacv1.ServiceAccountKind,
				Name:      name,
				Namespace: namespace,
			},
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: rbacv1.GroupName,
			Kind:     "ClusterRole",
			Name:     crName,
		},
	}

	return &cr, &sa, &crb
}

func MakeScopedRBACObjects(
	name,
	namespace,
	component string,
	additionalRules []rbacv1.PolicyRule,
	saAnnotations map[string]string,
) (*rbacv1.Role, *corev1.ServiceAccount, *rbacv1.RoleBinding) {
	rName, saName, rbName := MakeScopedRBACNames(name, component)
	r := rbacv1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name:      rName,
			Namespace: namespace,
		},
		Rules: []rbacv1.PolicyRule{
			{
				Verbs:     []string{"get"},
				APIGroups: []string{""},
				Resources: []string{"pods"},
			},
		},
	}

	r.Rules = append(r.Rules, additionalRules...)

	sa := corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:        saName,
			Namespace:   namespace,
			Annotations: saAnnotations,
		},
	}

	rb := rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      rbName,
			Namespace: namespace,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      rbacv1.ServiceAccountKind,
				Name:      name,
				Namespace: namespace,
			},
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: rbacv1.GroupName,
			Kind:     "Role",
			Name:     rName,
		},
	}

	return &r, &sa, &rb
}

func MakeRBACNames(name, component string) (string, string, string) {
	cr := fmt.Sprintf("fluent-operator-%s", component)
	crb := fmt.Sprintf("fluent-operator-%s-%s", component, name)
	return cr, name, crb
}

// MakeScopedRBACNames returns the namespaced Role/ServiceAccount/RoleBinding
// names. They follow the same convention as the cluster-scoped names.
func MakeScopedRBACNames(name, component string) (string, string, string) {
	return MakeRBACNames(name, component)
}

// MakeRBACObjectsForScope builds the RBAC objects an agent needs, returning the
// namespaced (Role/RoleBinding) variants when namespaced is true and the
// cluster-scoped (ClusterRole/ClusterRoleBinding) variants otherwise.
func MakeRBACObjectsForScope(
	namespaced bool,
	name, namespace, component string,
	additionalRules []rbacv1.PolicyRule,
	saAnnotations map[string]string,
) (role, sa, binding client.Object) {
	if namespaced {
		r, s, rb := MakeScopedRBACObjects(name, namespace, component, additionalRules, saAnnotations)
		return r, s, rb
	}
	cr, s, crb := MakeRBACObjects(name, namespace, component, additionalRules, saAnnotations)
	return cr, s, crb
}

// DeletePerInstanceBinding removes the per-instance RoleBinding created for an
// agent when running in namespaced mode. In cluster-scoped mode nothing is
// deleted: the per-instance ClusterRoleBinding references a shared ClusterRole,
// and the operator is intentionally not granted delete on cluster-scoped RBAC
// (keeping its footprint minimal), so the binding is left in place.
func DeletePerInstanceBinding(
	ctx context.Context,
	c client.Client,
	namespaced bool,
	name, namespace, component string,
) error {
	if !namespaced {
		return nil
	}
	_, _, rbName := MakeScopedRBACNames(name, component)
	rb := &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{Name: rbName, Namespace: namespace},
	}
	if err := c.Delete(ctx, rb); err != nil && !apierrors.IsNotFound(err) {
		return err
	}
	return nil
}
