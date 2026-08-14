package params

import (
	"strings"
	"testing"
)

func TestEscapeValue(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  string
	}{
		{
			name:  "plain value is left unchanged",
			value: "http://legit.com",
			want:  "http://legit.com",
		},
		{
			name:  "value with hash but no newline is left bare (no interpolation in bare values)",
			value: "pass#{word}",
			want:  "pass#{word}",
		},
		{
			name:  "match pattern is left unchanged",
			value: "kube.**",
			want:  "kube.**",
		},
		{
			name:  "newline injection is quoted and escaped",
			value: "http://legit.com\n</match>\n<match **>\n  @type exec",
			want:  `"http://legit.com\n</match>\n<match **>\n  @type exec"`,
		},
		{
			name:  "carriage return is escaped",
			value: "a\rb",
			want:  `"a\rb"`,
		},
		{
			name:  "backslash, quote, tab and hash are escaped when quoting",
			value: "a\n\t\"b\"\\c#{d}",
			want:  `"a\n\t\"b\"\\c\#{d}"`,
		},
		{
			// Leading '"' makes Fluentd evaluate '#{...}' as Ruby; re-quoting with
			// '#' escaped defuses it (GHSA-2j8x-46rv-qmpq).
			name:  "leading double-quote is re-quoted to disable Ruby interpolation",
			value: `"#{File.read('/etc/passwd')}"`,
			want:  `"\"\#{File.read('/etc/passwd')}\""`,
		},
		{
			name:  "leading whitespace before the quote still triggers quoting",
			value: ` "#{cmd}"`,
			want:  `" \"\#{cmd}\""`,
		},
		{
			name:  "numeric value is left unchanged (not quoted, preserves type)",
			value: "24224",
			want:  "24224",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := escapeValue(tt.value); got != tt.want {
				t.Errorf("escapeValue(%q) = %q, want %q", tt.value, got, tt.want)
			}
		})
	}
}

// TestPluginStoreString_NoInjectionViaBody reproduces the GHSA-2j8x-46rv-qmpq
// proof of concept: a newline in a plugin parameter value must not be able to
// close the enclosing <match> block and open a top-level <match> ... @type exec.
func TestPluginStoreString_NoInjectionViaBody(t *testing.T) {
	ps := NewPluginStore("match")
	ps.InsertPairs("@type", "http")
	ps.InsertPairs("tag", "**")
	ps.InsertPairs("endpoint",
		"http://legit.com\n</match>\n<match **>\n  @type exec\n  command id > /tmp/pwned")

	got := ps.String()

	// Only line-leading tags are directives. The rendered config must contain
	// exactly one closing </match> directive line: the one this plugin emits. An
	// injected </match> would appear on its own line; escaped, it stays inside the
	// quoted value on the parameter line.
	closingLines := 0
	for _, line := range strings.Split(got, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "</match>") {
			closingLines++
		}
		// The injected exec directive must never be rendered as a live directive.
		if trimmed == "@type exec" {
			t.Fatalf("injected `@type exec` directive was rendered live:\n%s", got)
		}
	}
	if closingLines != 1 {
		t.Fatalf("expected exactly one </match> directive line, got %d:\n%s", closingLines, got)
	}
}

// TestPluginStoreString_NoInjectionViaHead ensures a newline in a head-directive
// value (e.g. a match tag) cannot inject sibling directives.
func TestPluginStoreString_NoInjectionViaHead(t *testing.T) {
	ps := NewPluginStore("match")
	ps.InsertPairs("@type", "null")
	ps.InsertPairs("tag", "**>\n<match **")

	got := ps.String()

	// Only line-leading directives change block nesting; the escaped value keeps
	// everything on one physical line, so at most one line may open a <match.
	directiveLines := 0
	for _, line := range strings.Split(got, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "<match") {
			directiveLines++
		}
	}
	if directiveLines != 1 {
		t.Fatalf("expected exactly one <match directive line, got %d:\n%s", directiveLines, got)
	}
}

// TestPluginStoreString_NoInterpolationViaBody covers the newline-free RCE variant
// of GHSA-2j8x-46rv-qmpq: a double-quoted literal value must not render with a live
// '#{...}' interpolation.
func TestPluginStoreString_NoInterpolationViaBody(t *testing.T) {
	ps := NewPluginStore("match")
	ps.InsertPairs("@type", "http")
	ps.InsertPairs("endpoint", `"#{File.read('/etc/passwd')}"`)

	got := ps.String()

	// Live interpolation renders as `"#{`; when neutralized it is escaped to `\#{`.
	if strings.Contains(got, `"#{`) {
		t.Fatalf("value rendered as a live interpolating quoted literal:\n%s", got)
	}
	if !strings.Contains(got, `\#{`) {
		t.Fatalf("expected the interpolation metacharacter to be escaped (\\#{):\n%s", got)
	}
}
