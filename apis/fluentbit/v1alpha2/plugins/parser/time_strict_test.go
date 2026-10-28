package parser

import (
	"testing"

	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins"
	. "github.com/onsi/gomega"
)

func TestParsers_TimeStrict(t *testing.T) {
	strict := false

	tests := []struct {
		name   string
		plugin plugins.Plugin
	}{
		{"json", &JSON{TimeStrict: &strict}},
		{"logfmt", &Logfmt{TimeStrict: &strict}},
		{"ltsv", &LSTV{TimeStrict: &strict}},
		{"regex", &Regex{Regex: "^(?<log>.*)$", TimeStrict: &strict}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGomegaWithT(t)

			kvs, err := tt.plugin.Params(plugins.NewSecretLoader(nil, "test"))
			g.Expect(err).NotTo(HaveOccurred())
			out, err := kvs.String()
			g.Expect(err).NotTo(HaveOccurred())
			g.Expect(out).To(MatchRegexp(`Time_Strict\s+false\n`))
		})
	}
}

func TestParsers_TimeStrictUnset(t *testing.T) {
	g := NewGomegaWithT(t)

	for _, p := range []plugins.Plugin{&JSON{}, &Logfmt{}, &LSTV{}, &Regex{}} {
		kvs, err := p.Params(plugins.NewSecretLoader(nil, "test"))
		g.Expect(err).NotTo(HaveOccurred())
		out, err := kvs.String()
		g.Expect(err).NotTo(HaveOccurred())
		g.Expect(out).NotTo(ContainSubstring("Time_Strict"))
	}
}
