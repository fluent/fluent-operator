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
	"slices"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	fluentbitv1alpha2 "github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2"
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins/custom"
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins/output"
)

func httpOutputWithPassword(namespace, name, secretName string) *fluentbitv1alpha2.Output {
	return &fluentbitv1alpha2.Output{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    map[string]string{"fluentbit.fluent.io/enabled": "true"},
		},
		Spec: fluentbitv1alpha2.OutputSpec{
			Match: "*",
			HTTP: &output.HTTP{
				Host: "logs.example.com",
				HTTPPasswd: &plugins.Secret{ValueFrom: plugins.ValueSource{
					SecretKeyRef: corev1.SecretKeySelector{
						LocalObjectReference: corev1.LocalObjectReference{Name: secretName},
						Key:                  "password",
					},
				}},
			},
		},
	}
}

// TestNamespacedOutputWithMissingSecretIsSkipped reproduces
// https://github.com/fluent/fluent-operator/issues/1728: an Output whose Secret
// is missing in its namespace used to fail the whole config, so no namespace got
// its config updated. Such an Output is now left out, and the others are kept.
func TestNamespacedOutputWithMissingSecretIsSkipped(t *testing.T) {
	const namespace = "team-a"
	selector := metav1.LabelSelector{MatchLabels: map[string]string{"fluentbit.fluent.io/enabled": "true"}}
	cfg := &fluentbitv1alpha2.FluentBitConfig{
		ObjectMeta: metav1.ObjectMeta{Name: "fluentbit-config", Namespace: namespace},
		Spec:       fluentbitv1alpha2.NamespacedFluentBitCfgSpec{OutputSelector: selector},
	}
	customOutput := &fluentbitv1alpha2.Output{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "custom",
			Namespace: namespace,
			Labels:    map[string]string{"fluentbit.fluent.io/enabled": "true"},
		},
		Spec: fluentbitv1alpha2.OutputSpec{
			CustomPlugin: &custom.CustomPlugin{Config: "Name stdout\nMatch *"},
		},
	}
	objs := []client.Object{
		cfg,
		&corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{Name: "present", Namespace: namespace},
			Data:       map[string][]byte{"password": []byte("s3cret")},
		},
		httpOutputWithPassword(namespace, "with-secret", "present"),
		httpOutputWithPassword(namespace, "missing-secret", "absent"),
		httpOutputWithPassword(namespace, "missing-key", "present-without-key"),
		&corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{Name: "present-without-key", Namespace: namespace},
			Data:       map[string][]byte{"user": []byte("admin")},
		},
		customOutput,
	}

	s := newConfigTestScheme(t)
	cli := fake.NewClientBuilder().WithScheme(s).WithObjects(objs...).Build()
	r := &FluentBitConfigReconciler{Client: cli, Log: logf.Log, Scheme: s}

	_, outputs, _, _, _, _, err := r.ListFluentBitConfigResources(context.Background(), *cfg)
	if err != nil {
		t.Fatalf("ListFluentBitConfigResources returned error: %v", err)
	}

	var names []string
	for _, item := range outputs.Items {
		names = append(names, item.Name)
	}
	if len(names) != 2 || !slices.Contains(names, "with-secret") || !slices.Contains(names, "custom") {
		t.Fatalf("expected only the Outputs that can be rendered, got %v", names)
	}

	for _, item := range outputs.Items {
		if item.Name == "custom" && item.Spec.CustomPlugin.Config != customOutput.Spec.CustomPlugin.Config {
			t.Fatalf("the custom plugin config must not be changed by the check, got %q", item.Spec.CustomPlugin.Config)
		}
	}

	rendered, err := outputs.Load(plugins.NewSecretLoader(cli, namespace))
	if err != nil {
		t.Fatalf("rendering the kept Outputs returned error: %v", err)
	}
	if rendered == "" {
		t.Fatal("expected the kept Outputs to be rendered")
	}
}
