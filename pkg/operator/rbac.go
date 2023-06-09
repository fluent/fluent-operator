package operator

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func MakeRBACObjects(name, namespace, component string, additionalRules []rbacv1.PolicyRule, saAnnotations map[string]string) (*rbacv1.ClusterRole, *corev1.ServiceAccount, *rbacv1.ClusterRoleBinding) {
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

func MakeScopedRBACObjects(name, namespace string, saAnnotations map[string]string) (*rbacv1.Role, *corev1.ServiceAccount, *rbacv1.RoleBinding) {
	rName, saName, rbName := MakeScopedRBACNames(name)
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

func MakeScopedRBACNames(name string) (string, string, string) {
	r := "fluent:fluent-operator"
	rb := fmt.Sprintf("fluent-operator-fluent-bit-%s", name)
	return r, name, rb
}
