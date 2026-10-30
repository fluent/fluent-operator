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

	policyv1 "k8s.io/api/policy/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	fluentdv1alpha1 "github.com/fluent/fluent-operator/v3/apis/fluentd/v1alpha1"
)

func newFluentdTestScheme(t *testing.T) *runtime.Scheme {
	t.Helper()
	s := runtime.NewScheme()
	if err := scheme.AddToScheme(s); err != nil {
		t.Fatalf("failed to add client-go scheme: %v", err)
	}
	if err := fluentdv1alpha1.AddToScheme(s); err != nil {
		t.Fatalf("failed to add fluentd scheme: %v", err)
	}
	return s
}

// The operator only deletes a PodDisruptionBudget it controls: one written by hand
// under the Fluentd resource's name, the workaround before the operator could manage
// one, survives turning the feature off and running in agent mode.
func TestDeletePodDisruptionBudgetOnlyDeletesItsOwn(t *testing.T) {
	s := newFluentdTestScheme(t)
	fd := &fluentdv1alpha1.Fluentd{
		ObjectMeta: metav1.ObjectMeta{Name: "fluentd", Namespace: "logging", UID: "fluentd-uid"},
	}
	key := client.ObjectKey{Namespace: fd.Namespace, Name: fd.Name}

	for _, tc := range []struct {
		name        string
		pdb         *policyv1.PodDisruptionBudget
		wantDeleted bool
	}{
		{name: "no PodDisruptionBudget"},
		{
			name:        "owned by the Fluentd resource",
			pdb:         &policyv1.PodDisruptionBudget{ObjectMeta: metav1.ObjectMeta{Name: fd.Name, Namespace: fd.Namespace}},
			wantDeleted: true,
		},
		{
			name: "written by hand",
			pdb:  &policyv1.PodDisruptionBudget{ObjectMeta: metav1.ObjectMeta{Name: fd.Name, Namespace: fd.Namespace}},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			builder := fake.NewClientBuilder().WithScheme(s)
			if tc.pdb != nil {
				if tc.wantDeleted {
					if err := ctrl.SetControllerReference(fd, tc.pdb, s); err != nil {
						t.Fatal(err)
					}
				}
				builder = builder.WithObjects(tc.pdb)
			}
			cli := builder.Build()
			r := &FluentdReconciler{Client: cli, Log: logf.Log, Scheme: s}

			if err := r.deletePodDisruptionBudget(context.Background(), fd); err != nil {
				t.Fatalf("deletePodDisruptionBudget returned error: %v", err)
			}

			err := cli.Get(context.Background(), key, &policyv1.PodDisruptionBudget{})
			switch {
			case tc.pdb == nil || tc.wantDeleted:
				if !errors.IsNotFound(err) {
					t.Fatalf("PodDisruptionBudget still exists: %v", err)
				}
			case err != nil:
				t.Fatalf("PodDisruptionBudget the operator does not control was deleted: %v", err)
			}
		})
	}
}
