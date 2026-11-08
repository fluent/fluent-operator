package parser

import (
	"testing"

	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins"
	. "github.com/onsi/gomega"
)

func TestParsers_TimeKeep(t *testing.T) {
	keep := true

	tests := []struct {
		name   string
		plugin plugins.Plugin
	}{
		{"json", &JSON{TimeFormat: "%d/%b/%Y:%H:%M:%S %z", TimeKeep: &keep}},
		{"logfmt", &Logfmt{TimeFormat: "%d/%b/%Y:%H:%M:%S %z", TimeKeep: &keep}},
		{"ltsv", &LSTV{TimeFormat: "%d/%b/%Y:%H:%M:%S %z", TimeKeep: &keep}},
		{"regex", &Regex{Regex: "^(?<log>.*)$", TimeFormat: "%d/%b/%Y:%H:%M:%S %z", TimeKeep: &keep}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGomegaWithT(t)

			kvs, err := tt.plugin.Params(plugins.NewSecretLoader(nil, "test"))
			g.Expect(err).NotTo(HaveOccurred())
			out, err := kvs.String()
			g.Expect(err).NotTo(HaveOccurred())
			g.Expect(out).To(MatchRegexp(`Time_Keep\s+true\n`))
			g.Expect(out).To(MatchRegexp(`Time_Format\s+%d/%b/%Y:%H:%M:%S %z\n`))
			g.Expect(out).NotTo(MatchRegexp(`Time_Format\s+true\n`))
		})
	}
}
