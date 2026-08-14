package v1alpha2

import (
	"strings"
	"testing"

	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins/custom"
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins/output"
	"github.com/fluent/fluent-operator/v3/pkg/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// newlineInjection is the GHSA-2j8x-46rv-qmpq payload: a value that, if rendered
// unescaped, breaks onto a new line and opens an attacker-controlled [OUTPUT]
// section that exfiltrates via the file plugin.
const newlineInjection = "http://legit.com\n\n[OUTPUT]\n    Name    file\n    Path    /tmp/exfil"

// TestClusterOutput_Load_RejectsInjectionViaPluginValue ensures a newline in a
// plugin parameter that flows through the KVs renderer is rejected.
func TestClusterOutput_Load_RejectsInjectionViaPluginValue(t *testing.T) {
	list := ClusterOutputList{Items: []ClusterOutput{{
		TypeMeta:   metav1.TypeMeta{APIVersion: "fluentbit.fluent.io/v1alpha2", Kind: "ClusterOutput"},
		ObjectMeta: metav1.ObjectMeta{Name: "evil"},
		Spec: OutputSpec{
			Match: "*",
			HTTP:  &output.HTTP{Host: newlineInjection},
		},
	}}}

	assertRejectsInjection(t, func() (string, error) {
		return list.Load(plugins.NewSecretLoader(nil, "ns"))
	})
}

// TestClusterOutput_Load_RejectsInjectionViaAlias ensures a newline in a
// section-header field written directly (bypassing KVs) is rejected.
func TestClusterOutput_Load_RejectsInjectionViaAlias(t *testing.T) {
	list := ClusterOutputList{Items: []ClusterOutput{{
		TypeMeta:   metav1.TypeMeta{APIVersion: "fluentbit.fluent.io/v1alpha2", Kind: "ClusterOutput"},
		ObjectMeta: metav1.ObjectMeta{Name: "evil"},
		Spec: OutputSpec{
			Match:  "*",
			Alias:  newlineInjection,
			Stdout: &output.Stdout{},
		},
	}}}

	assertRejectsInjection(t, func() (string, error) {
		return list.Load(plugins.NewSecretLoader(nil, "ns"))
	})
	assertRejectsInjection(t, func() (string, error) {
		return list.LoadAsYaml(plugins.NewSecretLoader(nil, "ns"), 1)
	})
}

// TestClusterOutput_Load_AllowsBenignConfig confirms a normal config still renders.
func TestClusterOutput_Load_AllowsBenignConfig(t *testing.T) {
	list := ClusterOutputList{Items: []ClusterOutput{{
		TypeMeta:   metav1.TypeMeta{APIVersion: "fluentbit.fluent.io/v1alpha2", Kind: "ClusterOutput"},
		ObjectMeta: metav1.ObjectMeta{Name: "ok"},
		Spec: OutputSpec{
			Match: "logs.foo.bar",
			Alias: "my_alias",
			HTTP:  &output.HTTP{Host: "https://example.com", Port: utils.ToPtr[int32](443)},
		},
	}}}

	got, err := list.Load(plugins.NewSecretLoader(nil, "ns"))
	if err != nil {
		t.Fatalf("Load() unexpected error for benign config: %v", err)
	}
	if !strings.Contains(got, "https://example.com") || !strings.Contains(got, "Alias    my_alias") {
		t.Fatalf("benign config not rendered as expected:\n%s", got)
	}
}

// TestCustomPlugin_Params_RejectsSectionInjection covers the classic customPlugin
// passthrough: a body that opens its own [OUTPUT] section must be rejected.
func TestCustomPlugin_Params_RejectsSectionInjection(t *testing.T) {
	cp := &custom.CustomPlugin{
		Config: "Name    stdout\n[OUTPUT]\n    Name    file\n    Path    /tmp/exfil",
	}
	if _, err := cp.Params(plugins.NewSecretLoader(nil, "ns")); err == nil {
		t.Fatal("expected CustomPlugin.Params to reject a config that opens a new section, got nil")
	}
}

// TestCustomPlugin_Params_RejectsCommandInjection rejects @-commands (e.g. @INCLUDE).
func TestCustomPlugin_Params_RejectsCommandInjection(t *testing.T) {
	cp := &custom.CustomPlugin{
		Config: "Name    stdout\n@INCLUDE /etc/passwd",
	}
	if _, err := cp.Params(plugins.NewSecretLoader(nil, "ns")); err == nil {
		t.Fatal("expected CustomPlugin.Params to reject an @command, got nil")
	}
}

// TestCustomPlugin_Params_AllowsBenignBody confirms a normal multi-line body passes.
func TestCustomPlugin_Params_AllowsBenignBody(t *testing.T) {
	cp := &custom.CustomPlugin{
		Config: "Name    cpu\nTag    my_cpu\n",
	}
	kvs, err := cp.Params(plugins.NewSecretLoader(nil, "ns"))
	if err != nil {
		t.Fatalf("CustomPlugin.Params unexpected error for benign body: %v", err)
	}
	if kvs == nil || kvs.Content == "" {
		t.Fatal("expected benign body to render content")
	}
}

func assertRejectsInjection(t *testing.T, fn func() (string, error)) {
	t.Helper()
	out, err := fn()
	if err != nil {
		return // rejected as expected
	}
	// If not rejected, the payload must at least not have produced a second live
	// [OUTPUT] section (that would be a config-injection breakout).
	if strings.Count(strings.ToLower(out), "[output]") > 1 {
		t.Fatalf("config injection succeeded: more than one [OUTPUT] section rendered:\n%s", out)
	}
	t.Fatalf("expected injection to be rejected with an error, got nil:\n%s", out)
}
