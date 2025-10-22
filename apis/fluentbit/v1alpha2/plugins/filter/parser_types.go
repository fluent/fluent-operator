package filter

import (
	"crypto/md5"
	"fmt"
	"strings"

	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins/params"
)

// +kubebuilder:object:generate:=true

// The Parser Filter plugin allows to parse field in event records. <br />
// **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/filters/parser**
type Parser struct {
	plugins.CommonParams `json:",inline"`
	// Specify field name in record to parse.
	KeyName string `json:"keyName,omitempty"`
	// Specify the parser name to interpret the field.
	// Multiple Parser entries are allowed (split by comma).
	Parser string `json:"parser,omitempty"`
	// Keep original Key_Name field in the parsed result.
	// If false, the field will be removed.
	PreserveKey *bool `json:"preserveKey,omitempty"`
	// Keep all other original fields in the parsed result.
	// If false, all other original fields will be removed.
	ReserveData *bool `json:"reserveData,omitempty"`
	// If the key is a escaped string (e.g: stringify JSON), unescape the string before to apply the parser.
	UnescapeKey *bool `json:"unescapeKey,omitempty"`
}

func (*Parser) Name() string {
	return "parser"
}

func (p *Parser) Params(_ plugins.SecretLoader) (*params.KVs, error) {
	kvs := params.NewKVs()
	err := p.AddCommonParams(kvs)
	if err != nil {
		return kvs, err
	}
	if p.KeyName != "" {
		kvs.Insert("Key_Name", p.KeyName)
	}
	if p.Parser != "" {
		for parser := range strings.SplitSeq(p.Parser, ",") {
			kvs.Insert("Parser", strings.Trim(parser, " "))
		}
	}
	if p.PreserveKey != nil {
		kvs.Insert("Preserve_Key", fmt.Sprint(*p.PreserveKey))
	}
	if p.ReserveData != nil {
		kvs.Insert("Reserve_Data", fmt.Sprint(*p.ReserveData))
	}
	if p.UnescapeKey != nil {
		kvs.Insert("Unescape_Key", fmt.Sprint(*p.UnescapeKey))
	}
	return kvs, nil
}

func (p *Parser) MakeNamespaced(ns string) {
	parsers := strings.Split(p.Parser, ",")
	for i := range parsers {
		parsers[i] = strings.Trim(parsers[i], " ")
		parsers[i] = fmt.Sprintf("%s-%x", parsers[i], md5.Sum([]byte(ns)))
	}
	p.Parser = strings.Join(parsers, ",")
}
