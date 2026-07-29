/*
Copyright 2021.

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
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"

	fluentbitv1alpha2 "github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2"
)

// testAppLabelValue is the "app" label value used throughout these tests.
const testAppLabelValue = "fluentbit"

// TestFluentBitMutateDaemonSetPreservesImmutableSelector reproduces
// https://github.com/fluent/fluent-operator/issues/1944: changing
// spec.labels on a FluentBit CR must not cause the reconciler to try to
// rewrite the DaemonSet's immutable spec.selector.
func TestFluentBitMutateDaemonSetPreservesImmutableSelector(t *testing.T) {
	s := runtime.NewScheme()
	if err := scheme.AddToScheme(s); err != nil {
		t.Fatalf("failed to add client-go scheme: %v", err)
	}
	if err := fluentbitv1alpha2.AddToScheme(s); err != nil {
		t.Fatalf("failed to add fluentbit scheme: %v", err)
	}
	r := &FluentBitReconciler{Scheme: s}

	// The DaemonSet as it exists in the cluster today, created when
	// spec.labels only contained "app".
	existing := &appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{Name: "fluent-bit", Namespace: "default"},
		Spec: appsv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{MatchLabels: map[string]string{"app": testAppLabelValue}},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{"app": testAppLabelValue}},
			},
		},
	}

	// The user then adds a label to the FluentBit CR.
	fb := &fluentbitv1alpha2.FluentBit{
		ObjectMeta: metav1.ObjectMeta{Name: "fluent-bit", Namespace: "default"},
		Spec: fluentbitv1alpha2.FluentBitSpec{
			Labels: map[string]string{"app": testAppLabelValue, "fluentbit.fluent.io/enabled": "true"},
		},
	}

	if err := r.mutate(existing, fb)(); err != nil {
		t.Fatalf("mutate returned an unexpected error: %v", err)
	}

	if got := existing.Spec.Selector.MatchLabels; len(got) != 1 || got["app"] != testAppLabelValue {
		t.Fatalf("expected the immutable selector to stay unchanged, got %v", got)
	}

	// The pod template must still satisfy the preserved selector.
	if existing.Spec.Template.Labels["app"] != testAppLabelValue {
		t.Fatalf("expected pod template to retain the selector label, got %v", existing.Spec.Template.Labels)
	}
	// The new label should still be applied to the pod template.
	if existing.Spec.Template.Labels["fluentbit.fluent.io/enabled"] != "true" {
		t.Fatalf("expected new label to be applied to pod template, got %v", existing.Spec.Template.Labels)
	}
}

// TestFluentBitMutateDaemonSetRemovedLabelDoesNotMutateSpec covers removing a
// label that was part of the original selector: the pod template must still
// carry it (to satisfy the preserved selector) without mutating fb.Spec.Labels
// in place, since MakeDaemonSet reuses that map for the DaemonSet's own labels
// and selector.
func TestFluentBitMutateDaemonSetRemovedLabelDoesNotMutateSpec(t *testing.T) {
	s := runtime.NewScheme()
	if err := scheme.AddToScheme(s); err != nil {
		t.Fatalf("failed to add client-go scheme: %v", err)
	}
	if err := fluentbitv1alpha2.AddToScheme(s); err != nil {
		t.Fatalf("failed to add fluentbit scheme: %v", err)
	}
	r := &FluentBitReconciler{Scheme: s}

	existing := &appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{Name: "fluent-bit", Namespace: "default"},
		Spec: appsv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{MatchLabels: map[string]string{"app": testAppLabelValue, "team": "logging"}},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{"app": testAppLabelValue, "team": "logging"}},
			},
		},
	}

	// The user removes "team" from the FluentBit CR's labels.
	fbLabels := map[string]string{"app": testAppLabelValue}
	fb := &fluentbitv1alpha2.FluentBit{
		ObjectMeta: metav1.ObjectMeta{Name: "fluent-bit", Namespace: "default"},
		Spec:       fluentbitv1alpha2.FluentBitSpec{Labels: fbLabels},
	}

	if err := r.mutate(existing, fb)(); err != nil {
		t.Fatalf("mutate returned an unexpected error: %v", err)
	}

	got := existing.Spec.Selector.MatchLabels
	if len(got) != 2 || got["team"] != "logging" || got["app"] != testAppLabelValue {
		t.Fatalf("expected the immutable selector to stay unchanged, got %v", got)
	}
	if existing.Spec.Template.Labels["team"] != "logging" {
		t.Fatalf("expected pod template to retain the removed selector label, got %v", existing.Spec.Template.Labels)
	}

	// fb.Spec.Labels must not have been mutated as a side effect: MakeDaemonSet
	// reuses it directly as the DaemonSet's labels/selector map.
	if _, ok := fbLabels["team"]; ok {
		t.Fatalf("mutate must not write back into fb.Spec.Labels, got %v", fbLabels)
	}
	if len(fbLabels) != 1 {
		t.Fatalf("expected fb.Spec.Labels to remain untouched, got %v", fbLabels)
	}
}
