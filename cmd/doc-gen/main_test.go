package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestToLinkUsesGeneratedPluginDocFile(t *testing.T) {
	pluginPath := t.TempDir() + "/"
	writeSource := func(rel, content string) string {
		t.Helper()
		path := filepath.Join(pluginPath, rel)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
		return path
	}

	srcs := []string{
		writeSource("input/statsd.go", "package input\n\ntype StatsD struct {\n\tListen string\n}\n"),
		writeSource("output/types.go", "package output\n\ntype Output struct {\n\tName string\n}\n"),
		writeSource("custom/custom_types.go", "package custom\n\ntype CustomPlugin struct {\n\tConfig string\n}\n"),
	}

	t.Cleanup(func() { delete(pluginDocs, "test") })
	indexPluginDocs(DocumentsLocation{path: pluginPath, name: "test"}, srcs)

	tests := []struct {
		typeName    string
		inPluginDoc bool
		want        string
	}{
		{"input.StatsD", false, "[input.StatsD](plugins/test/input/statsd.md)"},
		{"output.Output", false, "[output.Output](plugins/test/output/types.md)"},
		{"custom.CustomPlugin", true, "[custom.CustomPlugin](../custom/custom.md)"},
		// Types that were not indexed keep the snake_case fallback.
		{"filter.RecordModifier", false, "[filter.RecordModifier](plugins/test/filter/record_modifier.md)"},
	}
	for _, tt := range tests {
		if got := toLink(tt.typeName, "test", tt.inPluginDoc); got != tt.want {
			t.Errorf("toLink(%q, %v) = %q, want %q", tt.typeName, tt.inPluginDoc, got, tt.want)
		}
	}
}
