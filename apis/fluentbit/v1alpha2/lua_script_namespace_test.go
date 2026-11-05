package v1alpha2

import (
	"crypto/md5"
	"fmt"
	"regexp"
	"testing"

	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins/filter"
)

func luaClusterFilter(name string, lua *filter.Lua) ClusterFilter {
	return ClusterFilter{
		ObjectMeta: metav1.ObjectMeta{Name: name},
		Spec:       FilterSpec{FilterItems: []FilterItem{{Lua: lua}}},
	}
}

func luaScript(key string) corev1.ConfigMapKeySelector {
	return corev1.ConfigMapKeySelector{LocalObjectReference: corev1.LocalObjectReference{Name: "scripts"}, Key: key}
}

func namespacedScriptName(namespace, key string) string {
	return fmt.Sprintf("%x-%s", md5.Sum([]byte(namespace)), key)
}

func clusterScriptName(namespace, configMap, key string) string {
	return fmt.Sprintf("%x-%s", md5.Sum([]byte("cluster/"+namespace+"/"+configMap)), key)
}

func luaScriptLoader() plugins.ConfigMapLoader {
	return plugins.NewConfigMapLoader(fake.NewClientBuilder().WithObjects(
		&corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{Name: "scripts", Namespace: "fluent"},
			Data:       map[string]string{"default.lua": "default script", "filter.lua": "fluent script"},
		},
		&corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{Name: "scripts", Namespace: "teams"},
			Data:       map[string]string{"other.lua": "other script", "filter.lua": "teams script"},
		},
	).Build(), "fluent")
}

// A ClusterFilter reads its Lua script from the namespace of the
// ClusterFluentBitConfig unless the script names another namespace.
func TestRenderLuaScriptNamespace(t *testing.T) {
	g := NewGomegaWithT(t)
	cl := luaScriptLoader()

	filters := ClusterFilterList{Items: []ClusterFilter{
		luaClusterFilter("default", &filter.Lua{Script: luaScript("default.lua")}),
		luaClusterFilter("other", &filter.Lua{Script: luaScript("other.lua"), ScriptNamespace: "teams"}),
	}}

	scripts, err := ClusterFluentBitConfig{}.RenderLuaScript(cl, filters, "fluent")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(scripts).To(ConsistOf(
		Script{Name: "default.lua", Content: "default script"},
		Script{Name: clusterScriptName("teams", "scripts", "other.lua"), Content: "other script"},
	))

	// The default namespace does not hold other.lua.
	filters.Items[1].Spec.FilterItems[0].Lua.ScriptNamespace = ""
	_, err = ClusterFluentBitConfig{}.RenderLuaScript(cl, filters, "fluent")
	g.Expect(err).To(HaveOccurred())
}

// Two ClusterFilters using the same key from different namespaces get scripts with different
// names, and each filter's rendered script path points at its own script.
func TestRenderLuaScriptSameKeyFromTwoNamespaces(t *testing.T) {
	g := NewGomegaWithT(t)
	cl := luaScriptLoader()

	fluentLua := &filter.Lua{Script: luaScript("filter.lua"), Call: "run"}
	teamsLua := &filter.Lua{Script: luaScript("filter.lua"), ScriptNamespace: "teams", Call: "run"}
	filters := ClusterFilterList{Items: []ClusterFilter{
		luaClusterFilter("fluent", fluentLua),
		luaClusterFilter("teams", teamsLua),
	}}

	scripts, err := ClusterFluentBitConfig{}.RenderLuaScript(cl, filters, "fluent")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(scripts).To(ConsistOf(
		Script{Name: "filter.lua", Content: "fluent script"},
		Script{Name: clusterScriptName("teams", "scripts", "filter.lua"), Content: "teams script"},
	))

	for _, tc := range []struct {
		lua  *filter.Lua
		path string
	}{
		{fluentLua, "/fluent-bit/config/filter.lua"},
		{teamsLua, "/fluent-bit/config/" + clusterScriptName("teams", "scripts", "filter.lua")},
	} {
		kvs, err := tc.lua.Params(plugins.SecretLoader{})
		g.Expect(err).NotTo(HaveOccurred())
		rendered, err := kvs.String()
		g.Expect(err).NotTo(HaveOccurred())
		g.Expect(rendered).To(MatchRegexp(`(?m)^\s*script\s+` + regexp.QuoteMeta(tc.path) + `$`))
	}
}

// A namespaced Filter can't use scriptNamespace to read a script from another namespace: it
// always reads from its own namespace, and its rendered script path matches that script.
func TestRenderNamespacedLuaScriptIgnoresScriptNamespace(t *testing.T) {
	g := NewGomegaWithT(t)
	cl := luaScriptLoader()

	lua := &filter.Lua{Script: luaScript("filter.lua"), ScriptNamespace: "teams", Call: "run"}
	filters := []FilterList{{Items: []Filter{{
		ObjectMeta: metav1.ObjectMeta{Name: "tenant", Namespace: "fluent"},
		Spec:       FilterSpec{FilterItems: []FilterItem{{Lua: lua}}},
	}}}}

	scripts, err := ClusterFluentBitConfig{}.RenderNamespacedLuaScript(cl, filters)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(scripts).To(Equal([]Script{
		{Name: namespacedScriptName("fluent", "filter.lua"), Content: "fluent script"},
	}))

	lua.MakeNamespaced("fluent")
	kvs, err := lua.Params(plugins.SecretLoader{})
	g.Expect(err).NotTo(HaveOccurred())
	rendered, err := kvs.String()
	g.Expect(err).NotTo(HaveOccurred())
	path := "/fluent-bit/config/" + namespacedScriptName("fluent", "filter.lua")
	g.Expect(rendered).To(MatchRegexp(`(?m)^\s*script\s+` + regexp.QuoteMeta(path) + `$`))
}

// A ClusterFilter reading a script from a namespace and a namespaced Filter in that namespace,
// both with the key filter.lua, get different script names, so the tenant's script can't replace
// the one the ClusterFilter runs.
func TestClusterAndNamespacedLuaScriptsDontCollide(t *testing.T) {
	g := NewGomegaWithT(t)
	cl := plugins.NewConfigMapLoader(fake.NewClientBuilder().WithObjects(
		&corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{Name: "admin-scripts", Namespace: "teams"},
			Data:       map[string]string{"filter.lua": "admin script"},
		},
		&corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{Name: "tenant-scripts", Namespace: "teams"},
			Data:       map[string]string{"filter.lua": "tenant script"},
		},
	).Build(), "fluent")

	selector := func(name string) corev1.ConfigMapKeySelector {
		return corev1.ConfigMapKeySelector{LocalObjectReference: corev1.LocalObjectReference{Name: name}, Key: "filter.lua"}
	}
	clusterLua := &filter.Lua{Script: selector("admin-scripts"), ScriptNamespace: "teams", Call: "run"}
	tenantLua := &filter.Lua{Script: selector("tenant-scripts"), Call: "run"}

	clusterScripts, err := ClusterFluentBitConfig{}.RenderLuaScript(
		cl, ClusterFilterList{Items: []ClusterFilter{luaClusterFilter("admin", clusterLua)}}, "fluent")
	g.Expect(err).NotTo(HaveOccurred())
	tenantScripts, err := ClusterFluentBitConfig{}.RenderNamespacedLuaScript(cl, []FilterList{{Items: []Filter{{
		ObjectMeta: metav1.ObjectMeta{Name: "tenant", Namespace: "teams"},
		Spec:       FilterSpec{FilterItems: []FilterItem{{Lua: tenantLua}}},
	}}}})
	g.Expect(err).NotTo(HaveOccurred())

	g.Expect(clusterScripts).To(HaveLen(1))
	g.Expect(tenantScripts).To(HaveLen(1))
	g.Expect(clusterScripts[0].Name).NotTo(Equal(tenantScripts[0].Name))
	g.Expect(clusterScripts[0].Content).To(Equal("admin script"))
	g.Expect(tenantScripts[0].Content).To(Equal("tenant script"))

	tenantLua.MakeNamespaced("teams")
	clusterPath := "/fluent-bit/config/" + clusterLua.ScriptFileName()
	tenantPath := "/fluent-bit/config/" + tenantLua.ScriptFileName()
	g.Expect(clusterPath).To(Equal("/fluent-bit/config/" + clusterScripts[0].Name))
	g.Expect(tenantPath).To(Equal("/fluent-bit/config/" + tenantScripts[0].Name))
}
