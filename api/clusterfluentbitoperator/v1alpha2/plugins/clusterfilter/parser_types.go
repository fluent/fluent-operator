package clusterfilter

import (

	"fmt"
	"kubesphere.io/fluentbit-operator/api/clusterfluentbitoperator/v1alpha2/plugins"
	"kubesphere.io/fluentbit-operator/api/clusterfluentbitoperator/v1alpha2/plugins/params"
	"strings"


)

// +kubebuilder:object:generate:=true

// The Parser Filter plugin allows to parse field in event records.
type Parser struct {
	// Specify field name in record to parse.
	KeyName string `json:"keyName,omitempty"`
	// Specify the clusterparser name to interpret the field.
	// Multiple Parser entries are allowed (split by comma).
	Parser string `json:"clusterparser,omitempty"`
	// Keep original Key_Name field in the parsed result.
	// If false, the field will be removed.
	PreserveKey *bool `json:"preserveKey,omitempty"`
	// Keep all other original fields in the parsed result.
	// If false, all other original fields will be removed.
	ReserveData *bool `json:"reserveData,omitempty"`
	// If the key is a escaped string (e.g: stringify JSON), unescape the string before to apply the clusterparser.
	UnescapeKey *bool `json:"unescapeKey,omitempty"`
}

func (_ *Parser) Name() string {
	return "clusterparser"
}

func (p *Parser) Params(_ plugins.SecretLoader) (*params.KVs, error) {
	kvs := params.NewKVs()
	if p.KeyName != "" {
		kvs.Insert("Key_Name", p.KeyName)
	}
	if p.Parser != "" {
		parsers := strings.Split(p.Parser, ",")
		for _, parser := range parsers {
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
