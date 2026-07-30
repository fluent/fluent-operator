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
	"strings"
	"testing"

	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins/input"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"

	fluentbitv1alpha2 "github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2"
)

// TestGenerateRewriteTagConfigYaml reproduces
// https://github.com/fluent/fluent-operator/issues/1689: when
// spec.configFileFormat is "yaml", the auto-generated rewrite_tag filter for
// namespaced FluentBitConfig resources must be rendered as YAML, not as a
// classic TOML snippet spliced into the YAML document.
func TestGenerateRewriteTagConfigYaml(t *testing.T) {
	r := &FluentBitConfigReconciler{}
	cfg := fluentbitv1alpha2.FluentBitConfig{
		ObjectMeta: metav1.ObjectMeta{Namespace: "foobar"},
	}
	inputs := fluentbitv1alpha2.ClusterInputList{
		Items: []fluentbitv1alpha2.ClusterInput{
			{
				Spec: fluentbitv1alpha2.InputSpec{
					Tail: &input.Tail{
						Tag:  "kube.*",
						Path: "/var/log/containers/*.log",
					},
				},
			},
		},
	}

	yamlFormat := "yaml"
	out, err := r.generateRewriteTagConfig(cfg, inputs, &yamlFormat)
	if err != nil {
		t.Fatalf("generateRewriteTagConfig returned error: %v", err)
	}

	if strings.Contains(out, "[Filter]") {
		t.Fatalf("expected YAML output, got classic TOML snippet:\n%s", out)
	}

	// The generated snippet must parse as valid YAML on its own, and must be
	// usable as a "pipeline.filters" fragment (i.e. it starts with a
	// top-level "filters:" key).
	var parsed map[string]interface{}
	if err := yaml.Unmarshal([]byte(out), &parsed); err != nil {
		t.Fatalf("generated rewrite_tag config is not valid YAML: %v\n%s", err, out)
	}
	if _, ok := parsed["filters"]; !ok {
		t.Fatalf("expected a top-level \"filters\" key, got:\n%s", out)
	}

	if !strings.Contains(out, "name: rewrite_tag") {
		t.Fatalf("expected a rewrite_tag filter entry, got:\n%s", out)
	}

	// classic (default) format must remain unchanged (TOML).
	classicOut, err := r.generateRewriteTagConfig(cfg, inputs, nil)
	if err != nil {
		t.Fatalf("generateRewriteTagConfig returned error: %v", err)
	}
	if !strings.Contains(classicOut, "[Filter]") || !strings.Contains(classicOut, "Name    rewrite_tag") {
		t.Fatalf("expected classic TOML output, got:\n%s", classicOut)
	}
}
