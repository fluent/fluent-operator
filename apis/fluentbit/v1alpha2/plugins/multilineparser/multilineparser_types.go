package multilineparser

import (
	"fmt"
	"github.com/fluent/fluent-operator/v2/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v2/apis/fluentbit/v1alpha2/plugins/params"
)

// +kubebuilder:object:generate:=true

// **For full documentation, refer to https://docs.fluentbit.io/manual/administration/configuring-fluent-bit/multiline-parsing**
type MultilineParser struct {
	// Set the multiline mode, for now, we support the type regex.
	// +kubebuilder:validation:Enum:=regex
	// +kubebuilder:default:=regex
	Type string `json:"type,omitempty"`
	// Name of a pre-defined parser that must be applied to the incoming content before applying the regex rule. If no parser is defined, it's assumed that's a raw text and not a structured message.
	Parser string `json:"parser,omitempty"`
	// For an incoming structured message, specify the key that contains the data that should be processed by the regular expression and possibly concatenated.
	KeyContent string `json:"keyContent,omitempty"`
	// Timeout in milliseconds to flush a non-terminated multiline buffer. Default is set to 5 seconds.
	// +kubebuilder:default:=5000
	FlushTimeout int `json:"flushTimeout,omitempty"`
	// Configure a rule to match a multiline pattern. The rule has a specific format described below. Multiple rules can be defined.
	Rules []Rule `json:"rules,omitempty"`
}

type Rule struct {
	Start string `json:"start"`
	Regex string `json:"regex"`
	Next  string `json:"next"`
}

func (_ *MultilineParser) Name() string {
	return "multilineparser"
}

func (m *MultilineParser) Params(_ plugins.SecretLoader) (*params.KVs, error) {
	kvs := params.NewKVs()

	if m.Type != "" {
		kvs.Insert("Type", m.Type)
	}
	if m.Parser != "" {
		kvs.Insert("Parser", m.Parser)
	}
	if m.KeyContent != "" {
		kvs.Insert("Key_Content", m.KeyContent)
	}
	if m.FlushTimeout != 0 {
		kvs.Insert("Flush_Timeout", fmt.Sprint(m.FlushTimeout))
	}
	if len(m.Rules) != 0 {
		for _, rule := range m.Rules {
			// add quotes and don't try to escape characters
			kvs.Insert("Rule", fmt.Sprintf("\"%s\" \"%s\" \"%s\"", rule.Start, rule.Regex, rule.Next))
		}
	}
	return kvs, nil
}
