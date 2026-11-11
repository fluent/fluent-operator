/*
Copyright 2026.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	fluentbitv1alpha2 "github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2"
)

func newFluentBitTestScheme(t *testing.T) *runtime.Scheme {
	t.Helper()
	s := runtime.NewScheme()
	if err := scheme.AddToScheme(s); err != nil {
		t.Fatalf("failed to add client-go scheme: %v", err)
	}
	if err := fluentbitv1alpha2.AddToScheme(s); err != nil {
		t.Fatalf("failed to add fluentbit scheme: %v", err)
	}
	return s
}

func TestDeleteServiceAccountOnlyDeletesItsOwn(t *testing.T) {
	s := newFluentBitTestScheme(t)
	fb := &fluentbitv1alpha2.FluentBit{
		ObjectMeta: metav1.ObjectMeta{Name: "fluent-bit", Namespace: "logging", UID: types.UID("fb-uid")},
	}
	key := client.ObjectKey{Namespace: fb.Namespace, Name: fb.Name}

	for _, tc := range []struct {
		name        string
		sa          *corev1.ServiceAccount
		wantDeleted bool
	}{
		{name: "no ServiceAccount"},
		{
			name:        "owned by the FluentBit",
			sa:          &corev1.ServiceAccount{ObjectMeta: metav1.ObjectMeta{Name: fb.Name, Namespace: fb.Namespace}},
			wantDeleted: true,
		},
		{
			name: "unrelated ServiceAccount with the same name",
			sa:   &corev1.ServiceAccount{ObjectMeta: metav1.ObjectMeta{Name: fb.Name, Namespace: fb.Namespace}},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			builder := fake.NewClientBuilder().WithScheme(s)
			if tc.sa != nil {
				if tc.wantDeleted {
					if err := ctrl.SetControllerReference(fb, tc.sa, s); err != nil {
						t.Fatal(err)
					}
				}
				builder = builder.WithObjects(tc.sa)
			}
			cli := builder.Build()
			r := &FluentBitReconciler{Client: cli, Log: logf.Log, Scheme: s}

			if err := r.deleteServiceAccount(context.Background(), fb); err != nil {
				t.Fatalf("deleteServiceAccount returned error: %v", err)
			}

			err := cli.Get(context.Background(), key, &corev1.ServiceAccount{})
			switch {
			case tc.sa == nil || tc.wantDeleted:
				if !errors.IsNotFound(err) {
					t.Fatalf("ServiceAccount still exists: %v", err)
				}
			default:
				if err != nil {
					t.Fatalf("unrelated ServiceAccount was deleted: %v", err)
				}
			}
		})
	}
}

func TestReconcileWithServiceAccountNameSkipsTheServiceAccount(t *testing.T) {
	s := newFluentBitTestScheme(t)
	fb := &fluentbitv1alpha2.FluentBit{
		ObjectMeta: metav1.ObjectMeta{
			Name:       "fluent-bit",
			Namespace:  "logging",
			Finalizers: []string{fluentbitv1alpha2.FluentBitFinalizerName},
		},
		Spec: fluentbitv1alpha2.FluentBitSpec{
			FluentBitConfigName: "fluent-bit-config",
			ServiceAccountName:  "fluent-bit-irsa",
		},
	}
	secret := &corev1.Secret{ObjectMeta: metav1.ObjectMeta{Name: "fluent-bit-config", Namespace: "logging"}}
	cli := fake.NewClientBuilder().WithScheme(s).WithObjects(fb, secret).Build()
	r := &FluentBitReconciler{Client: cli, Log: logf.Log, Scheme: s}

	key := client.ObjectKey{Namespace: fb.Namespace, Name: fb.Name}
	if _, err := r.Reconcile(context.Background(), ctrl.Request{NamespacedName: key}); err != nil {
		t.Fatalf("Reconcile returned error: %v", err)
	}

	var accounts corev1.ServiceAccountList
	if err := cli.List(context.Background(), &accounts); err != nil {
		t.Fatal(err)
	}
	if len(accounts.Items) != 0 {
		t.Errorf("expected no ServiceAccount to be created, got %d", len(accounts.Items))
	}

	var bindings rbacv1.ClusterRoleBindingList
	if err := cli.List(context.Background(), &bindings); err != nil {
		t.Fatal(err)
	}
	if len(bindings.Items) != 1 {
		t.Fatalf("expected one ClusterRoleBinding, got %d", len(bindings.Items))
	}
	if subjects := bindings.Items[0].Subjects; len(subjects) != 1 ||
		subjects[0].Name != "fluent-bit-irsa" || subjects[0].Namespace != "logging" {
		t.Errorf("ClusterRoleBinding subjects = %+v, want ServiceAccount logging/fluent-bit-irsa", subjects)
	}

	var ds appsv1.DaemonSet
	if err := cli.Get(context.Background(), key, &ds); err != nil {
		t.Fatal(err)
	}
	if got := ds.Spec.Template.Spec.ServiceAccountName; got != "fluent-bit-irsa" {
		t.Errorf("DaemonSet serviceAccountName = %q, want fluent-bit-irsa", got)
	}
}
