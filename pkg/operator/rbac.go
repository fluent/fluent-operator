package operator

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func MakeRBACObjects(name, namespace, component string, additionalRules []rbacv1.PolicyRule, saAnnotations map[string]string) (*rbacv1.ClusterRole, *corev1.ServiceAccount, *rbacv1.ClusterRoleBinding) {
	rbacName := fmt.Sprintf("fluent-operator-%s", component)
	cr := rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: rbacName,
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
			Name:        name,
			Namespace:   namespace,
			Annotations: saAnnotations,
		},
	}

	crb := rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("fluent-operator-%s-%s", component, name),
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
			Name:     rbacName,
		},
	}

	return &cr, &sa, &crb
}

func MakeScopedRBACObjects(fbName, fbNamespace string, saAnnotations map[string]string) (*rbacv1.Role, *corev1.ServiceAccount, *rbacv1.RoleBinding) {
	r := rbacv1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "fluent:fluent-operator",
			Namespace: fbNamespace,
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
			Name:        fbName,
			Namespace:   fbNamespace,
			Annotations: saAnnotations,
		},
	}

	rb := rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("fluent-operator-fluent-bit-%s", fbName),
			Namespace: fbNamespace,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      rbacv1.ServiceAccountKind,
				Name:      fbName,
				Namespace: fbNamespace,
			},
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: rbacv1.GroupName,
			Kind:     "Role",
			Name:     "fluent:fluent-operator",
		},
	}

	return &r, &sa, &rb
}
