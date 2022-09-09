package v1alpha2

import (
	"github.com/fluent/fluent-operator/apis/fluentbit/v1alpha2/plugins"
	. "github.com/onsi/gomega"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

var customs = `
Match    *
Host    127.0.0.1
Port    8088
TLS    On
TLS.Verify    Off
`

var res = `[Output]
    Name    splunk
    Match    *
    Host    127.0.0.1
    Port    8088
    TLS    On
    TLS.Verify    Off
[Output]
    Name    splunk
    Match    *
    Host    127.0.0.1
    Port    8088
    TLS    On
    TLS.Verify    Off
[Output]
    Name    splunk
    Match    *
    Host    127.0.0.1
    Port    8088
    TLS    On
    TLS.Verify    Off
`

func TestClusterClustomPlugin_Load(t *testing.T) {
	g := NewGomegaWithT(t)
	plugin := ClusterCustomPlugin{
		TypeMeta: v1.TypeMeta{
			Kind:       "ClusterCustomPlugin",
			APIVersion: "fluentbit.fluent.io/v1alpha2",
		},
		ObjectMeta: v1.ObjectMeta{
			Name:   "splunk-output",
			Labels: map[string]string{"type": "output", "fluentbit.fluent.io/enabled": "true"},
			Annotations: map[string]string{
				"plugin.config": customs,
			},
		},
		Spec: CustomPluginSpec{
			PluginName: "splunk",
			PluginType: "output",
		},
	}
	plugin1 := ClusterCustomPlugin{
		TypeMeta: v1.TypeMeta{
			Kind:       "ClusterCustomPlugin",
			APIVersion: "fluentbit.fluent.io/v1alpha2",
		},
		ObjectMeta: v1.ObjectMeta{
			Name:   "splunk-output",
			Labels: map[string]string{"type": "output", "fluentbit.fluent.io/enabled": "true"},
			Annotations: map[string]string{
				"plugin.config": customs,
			},
		},
		Spec: CustomPluginSpec{
			PluginName: "splunk",
			PluginType: "output",
		},
	}
	plugin2 := ClusterCustomPlugin{
		TypeMeta: v1.TypeMeta{
			Kind:       "ClusterCustomPlugin",
			APIVersion: "fluentbit.fluent.io/v1alpha2",
		},
		ObjectMeta: v1.ObjectMeta{
			Name:   "splunk-output",
			Labels: map[string]string{"type": "output", "fluentbit.fluent.io/enabled": "true"},
			Annotations: map[string]string{
				"plugin.config": customs,
			},
		},
		Spec: CustomPluginSpec{
			PluginName: "splunk",
			PluginType: "output",
		},
	}

	list := ClusterCustomPluginList{}
	list.Items = append(list.Items, plugin)
	list.Items = append(list.Items, plugin1)
	list.Items = append(list.Items, plugin2)
	sl := plugins.NewSecretLoader(nil, "testnamespace", nil)
	load, err := list.Load(sl)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(load)
	t.Log(res)
	g.Expect(load).To(Equal(res))
}
