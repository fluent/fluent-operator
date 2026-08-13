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
	"context"
	"strings"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	fluentbitv1alpha2 "github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2"
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins/input"
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins/output"
)

const (
	testCollectorNamespace = "logging"
	testConfigName         = "collector-config"
)

func newConfigTestScheme(t *testing.T) *runtime.Scheme {
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

// newCollectorTestObjects returns a ClusterFluentBitConfig selecting one cluster input
// and one cluster output, plus a Collector referencing that config by name.
func newCollectorTestObjects() []client.Object {
	selector := metav1.LabelSelector{MatchLabels: map[string]string{"fluentbit.fluent.io/enabled": "true"}}

	return []client.Object{
		&fluentbitv1alpha2.ClusterFluentBitConfig{
			ObjectMeta: metav1.ObjectMeta{Name: testConfigName},
			Spec: fluentbitv1alpha2.FluentBitConfigSpec{
				InputSelector:  selector,
				FilterSelector: selector,
				OutputSelector: selector,
			},
		},
		&fluentbitv1alpha2.ClusterInput{
			ObjectMeta: metav1.ObjectMeta{
				Name:   "collector-forward",
				Labels: map[string]string{"fluentbit.fluent.io/enabled": "true"},
			},
			Spec: fluentbitv1alpha2.InputSpec{
				Forward: &input.Forward{Port: ptrInt32(24224)},
			},
		},
		&fluentbitv1alpha2.ClusterOutput{
			ObjectMeta: metav1.ObjectMeta{
				Name:   "collector-stdout",
				Labels: map[string]string{"fluentbit.fluent.io/enabled": "true"},
			},
			Spec: fluentbitv1alpha2.OutputSpec{
				Match:  "*",
				Stdout: &output.Stdout{},
			},
		},
		&fluentbitv1alpha2.Collector{
			ObjectMeta: metav1.ObjectMeta{Name: "collector", Namespace: testCollectorNamespace},
			Spec:       fluentbitv1alpha2.CollectorSpec{FluentBitConfigName: testConfigName},
		},
	}
}

func ptrInt32(i int32) *int32 { return &i }

// TestCollectorClusterFluentBitConfigIsRendered reproduces
// https://github.com/fluent/fluent-operator/issues/1436: a Collector referencing a
// ClusterFluentBitConfig used to be ignored by the FluentBitConfig controller, which
// only ever looked at FluentBit DaemonSets. As a result no Secret was ever rendered
// and the Collector controller requeued forever without creating anything.
func TestCollectorClusterFluentBitConfigIsRendered(t *testing.T) {
	s := newConfigTestScheme(t)
	cli := fake.NewClientBuilder().WithScheme(s).WithObjects(newCollectorTestObjects()...).Build()

	r := &FluentBitConfigReconciler{Client: cli, Log: logf.Log, Scheme: s}

	if _, err := r.Reconcile(context.Background(), ctrl.Request{}); err != nil {
		t.Fatalf("Reconcile returned error: %v", err)
	}

	var sec corev1.Secret
	key := client.ObjectKey{Namespace: testCollectorNamespace, Name: testConfigName}
	if err := cli.Get(context.Background(), key, &sec); err != nil {
		t.Fatalf("expected the rendered configuration Secret %s, got: %v", key, err)
	}

	mainCfg := string(sec.Data["fluent-bit.conf"])
	if !strings.Contains(mainCfg, "forward") {
		t.Fatalf("expected the selected cluster input in the rendered config, got:\n%s", mainCfg)
	}
	if !strings.Contains(mainCfg, "stdout") {
		t.Fatalf("expected the selected cluster output in the rendered config, got:\n%s", mainCfg)
	}
	if _, ok := sec.Data["parsers.conf"]; !ok {
		t.Fatalf("expected parsers.conf in the rendered Secret, got keys %v", secretKeys(sec))
	}

	owners := sec.GetOwnerReferences()
	if len(owners) != 1 || owners[0].Kind != "ClusterFluentBitConfig" || owners[0].Name != testConfigName {
		t.Fatalf("expected the Secret to be owned by the ClusterFluentBitConfig, got %v", owners)
	}
}

// TestCollectorConfigRenderedInYamlFormat verifies that spec.configFileFormat is
// honoured on the Collector code path as well.
func TestCollectorConfigRenderedInYamlFormat(t *testing.T) {
	s := newConfigTestScheme(t)
	objs := newCollectorTestObjects()
	yamlFormat := configFileFormatYaml
	objs[0].(*fluentbitv1alpha2.ClusterFluentBitConfig).Spec.ConfigFileFormat = &yamlFormat

	cli := fake.NewClientBuilder().WithScheme(s).WithObjects(objs...).Build()
	r := &FluentBitConfigReconciler{Client: cli, Log: logf.Log, Scheme: s}

	if _, err := r.Reconcile(context.Background(), ctrl.Request{}); err != nil {
		t.Fatalf("Reconcile returned error: %v", err)
	}

	var sec corev1.Secret
	key := client.ObjectKey{Namespace: testCollectorNamespace, Name: testConfigName}
	if err := cli.Get(context.Background(), key, &sec); err != nil {
		t.Fatalf("expected the rendered configuration Secret %s, got: %v", key, err)
	}
	if _, ok := sec.Data["fluent-bit.yaml"]; !ok {
		t.Fatalf("expected fluent-bit.yaml in the rendered Secret, got keys %v", secretKeys(sec))
	}
}

// TestCollectorConfigHonoursConfigNamespaceOverride checks that an explicit
// spec.namespace on the ClusterFluentBitConfig still wins over the Collector's own
// namespace, consistently with the FluentBit code path.
func TestCollectorConfigHonoursConfigNamespaceOverride(t *testing.T) {
	s := newConfigTestScheme(t)
	objs := newCollectorTestObjects()
	pinned := "fluent"
	objs[0].(*fluentbitv1alpha2.ClusterFluentBitConfig).Spec.Namespace = &pinned

	cli := fake.NewClientBuilder().WithScheme(s).WithObjects(objs...).Build()
	r := &FluentBitConfigReconciler{Client: cli, Log: logf.Log, Scheme: s}

	if _, err := r.Reconcile(context.Background(), ctrl.Request{}); err != nil {
		t.Fatalf("Reconcile returned error: %v", err)
	}

	var sec corev1.Secret
	if err := cli.Get(
		context.Background(), client.ObjectKey{Namespace: pinned, Name: testConfigName}, &sec,
	); err != nil {
		t.Fatalf("expected the Secret in the pinned namespace %q, got: %v", pinned, err)
	}

	err := cli.Get(
		context.Background(),
		client.ObjectKey{Namespace: testCollectorNamespace, Name: testConfigName},
		&corev1.Secret{},
	)
	if err == nil {
		t.Fatalf("did not expect a Secret in the Collector namespace when spec.namespace is pinned")
	}
}

// TestCollectorDoesNotOverrideFluentBitConfig makes sure that a Collector sharing a
// ClusterFluentBitConfig with a FluentBit DaemonSet cannot overwrite the Secret the
// DaemonSet is already using, which would drop its namespaced (multi-tenant) plugins.
func TestCollectorDoesNotOverrideFluentBitConfig(t *testing.T) {
	s := newConfigTestScheme(t)
	objs := newCollectorTestObjects()
	// Pin the config so both code paths resolve to the very same namespace.
	pinned := testCollectorNamespace
	objs[0].(*fluentbitv1alpha2.ClusterFluentBitConfig).Spec.Namespace = &pinned
	objs = append(objs, &fluentbitv1alpha2.FluentBit{
		ObjectMeta: metav1.ObjectMeta{Name: "fluent-bit", Namespace: testCollectorNamespace},
		Spec: fluentbitv1alpha2.FluentBitSpec{
			FluentBitConfigName: testConfigName,
			NamespacedFluentBitCfgSelector: metav1.LabelSelector{
				MatchLabels: map[string]string{"never": "matches"},
			},
		},
	})

	cli := fake.NewClientBuilder().WithScheme(s).WithObjects(objs...).Build()
	r := &FluentBitConfigReconciler{Client: cli, Log: logf.Log, Scheme: s}

	if _, err := r.Reconcile(context.Background(), ctrl.Request{}); err != nil {
		t.Fatalf("Reconcile returned error: %v", err)
	}

	var sec corev1.Secret
	key := client.ObjectKey{Namespace: testCollectorNamespace, Name: testConfigName}
	if err := cli.Get(context.Background(), key, &sec); err != nil {
		t.Fatalf("expected the rendered configuration Secret %s, got: %v", key, err)
	}
	if !strings.Contains(string(sec.Data["fluent-bit.conf"]), "forward") {
		t.Fatalf("expected a rendered config, got:\n%s", sec.Data["fluent-bit.conf"])
	}
}

// TestCollectorsForSecret verifies the Secret -> Collector mapping used to wake the
// Collector controller up as soon as its configuration has been rendered.
func TestCollectorsForSecret(t *testing.T) {
	s := newConfigTestScheme(t)
	cli := fake.NewClientBuilder().WithScheme(s).WithObjects(
		&fluentbitv1alpha2.Collector{
			ObjectMeta: metav1.ObjectMeta{Name: "matching", Namespace: testCollectorNamespace},
			Spec:       fluentbitv1alpha2.CollectorSpec{FluentBitConfigName: testConfigName},
		},
		&fluentbitv1alpha2.Collector{
			ObjectMeta: metav1.ObjectMeta{Name: "other-config", Namespace: testCollectorNamespace},
			Spec:       fluentbitv1alpha2.CollectorSpec{FluentBitConfigName: "something-else"},
		},
		&fluentbitv1alpha2.Collector{
			ObjectMeta: metav1.ObjectMeta{Name: "other-namespace", Namespace: "elsewhere"},
			Spec:       fluentbitv1alpha2.CollectorSpec{FluentBitConfigName: testConfigName},
		},
	).Build()

	r := &CollectorReconciler{Client: cli, Log: logf.Log, Scheme: s}
	reqs := r.collectorsForSecret(context.Background(), &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{Name: testConfigName, Namespace: testCollectorNamespace},
	})

	if len(reqs) != 1 {
		t.Fatalf("expected exactly one reconcile request, got %v", reqs)
	}
	if reqs[0].Name != "matching" || reqs[0].Namespace != testCollectorNamespace {
		t.Fatalf("unexpected reconcile request: %v", reqs[0])
	}
}

func secretKeys(sec corev1.Secret) []string {
	keys := make([]string, 0, len(sec.Data))
	for k := range sec.Data {
		keys = append(keys, k)
	}
	return keys
}

// TestFluentBitNamespacedConfigStillRendered guards the refactoring that introduced
// the Collector code path: the FluentBit DaemonSet path must keep rendering cluster
// scoped plugins, namespace level plugins and the generated rewrite_tag filters into
// the operator namespace.
func TestFluentBitNamespacedConfigStillRendered(t *testing.T) {
	t.Setenv("NAMESPACE", "fluent")
	s := newConfigTestScheme(t)
	selector := metav1.LabelSelector{MatchLabels: map[string]string{"fluentbit.fluent.io/enabled": "true"}}

	objs := []client.Object{
		&fluentbitv1alpha2.ClusterFluentBitConfig{
			ObjectMeta: metav1.ObjectMeta{Name: testConfigName},
			Spec: fluentbitv1alpha2.FluentBitConfigSpec{
				InputSelector:  selector,
				FilterSelector: selector,
				OutputSelector: selector,
			},
		},
		&fluentbitv1alpha2.ClusterInput{
			ObjectMeta: metav1.ObjectMeta{
				Name:   "tail",
				Labels: map[string]string{"fluentbit.fluent.io/enabled": "true"},
			},
			Spec: fluentbitv1alpha2.InputSpec{
				Tail: &input.Tail{Tag: "kube.*", Path: "/var/log/containers/*.log"},
			},
		},
		&fluentbitv1alpha2.FluentBit{
			ObjectMeta: metav1.ObjectMeta{Name: "fluent-bit", Namespace: "fluent"},
			Spec: fluentbitv1alpha2.FluentBitSpec{
				FluentBitConfigName:            testConfigName,
				NamespacedFluentBitCfgSelector: selector,
			},
		},
		&fluentbitv1alpha2.FluentBitConfig{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "tenant",
				Namespace: "tenant-a",
				Labels:    map[string]string{"fluentbit.fluent.io/enabled": "true"},
			},
			Spec: fluentbitv1alpha2.NamespacedFluentBitCfgSpec{OutputSelector: selector},
		},
		&fluentbitv1alpha2.Output{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "tenant-stdout",
				Namespace: "tenant-a",
				Labels:    map[string]string{"fluentbit.fluent.io/enabled": "true"},
			},
			Spec: fluentbitv1alpha2.OutputSpec{Match: "*", Stdout: &output.Stdout{}},
		},
	}

	cli := fake.NewClientBuilder().WithScheme(s).WithObjects(objs...).Build()
	r := &FluentBitConfigReconciler{Client: cli, Log: logf.Log, Scheme: s}

	if _, err := r.Reconcile(context.Background(), ctrl.Request{}); err != nil {
		t.Fatalf("Reconcile returned error: %v", err)
	}

	var sec corev1.Secret
	if err := cli.Get(
		context.Background(), client.ObjectKey{Namespace: "fluent", Name: testConfigName}, &sec,
	); err != nil {
		t.Fatalf("expected the rendered configuration Secret fluent/%s, got: %v", testConfigName, err)
	}

	mainCfg := string(sec.Data["fluent-bit.conf"])
	t.Logf("rendered fluent-bit.conf:\n%s", mainCfg)
	for _, want := range []string{"tail", "rewrite_tag", "stdout"} {
		if !strings.Contains(mainCfg, want) {
			t.Fatalf("expected %q in the rendered config, got:\n%s", want, mainCfg)
		}
	}
}

// TestIsFluentBitConfigSecret guards the predicate that filters the Collector's
// Secret watch. Unrelated Secrets must be dropped before collectorsForSecret
// runs, while Secrets rendered by FluentBitConfigReconciler must pass.
func TestIsFluentBitConfigSecret(t *testing.T) {
	tests := []struct {
		name        string
		annotations map[string]string
		want        bool
	}{
		{
			name:        "config secret rendered by the operator",
			annotations: map[string]string{fluentBitConfigHashAnnotation: "abc123"},
			want:        true,
		},
		{
			name:        "unrelated secret with no annotations",
			annotations: nil,
			want:        false,
		},
		{
			name:        "unrelated secret with other annotations",
			annotations: map[string]string{"helm.sh/release": "foo"},
			want:        false,
		},
		{
			// A Secret written before the config-hash annotation existed is
			// filtered out, but FluentBitConfigReconciler treats a missing
			// annotation as a change and rewrites it, so the next event passes.
			name:        "legacy config secret without the annotation",
			annotations: nil,
			want:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sec := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:        testConfigName,
					Namespace:   testCollectorNamespace,
					Annotations: tt.annotations,
				},
			}
			if got := isFluentBitConfigSecret(sec); got != tt.want {
				t.Errorf("isFluentBitConfigSecret() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestConfigSecretCarriesHashAnnotation asserts the contract the predicate relies
// on: every Secret rendered for a Collector is stamped with the annotation, so
// the watch is never silently starved.
func TestConfigSecretCarriesHashAnnotation(t *testing.T) {
	s := newConfigTestScheme(t)
	cl := fake.NewClientBuilder().WithScheme(s).WithObjects(newCollectorTestObjects()...).Build()
	r := &FluentBitConfigReconciler{Client: cl, Log: logf.Log, Scheme: s}

	if _, err := r.Reconcile(context.Background(), ctrl.Request{}); err != nil {
		t.Fatalf("reconcile: %v", err)
	}

	var sec corev1.Secret
	if err := cl.Get(context.Background(),
		client.ObjectKey{Name: testConfigName, Namespace: testCollectorNamespace}, &sec); err != nil {
		t.Fatalf("expected rendered secret: %v", err)
	}

	if !isFluentBitConfigSecret(&sec) {
		t.Fatalf("rendered config secret lacks %s; the Collector Secret watch would ignore it",
			fluentBitConfigHashAnnotation)
	}
}
