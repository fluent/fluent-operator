package parser

import (
	"fmt"

	"kubesphere.io/fluentbit-operator/api/fluentbitoperator/v1alpha2/plugins"
	"kubesphere.io/fluentbit-operator/api/fluentbitoperator/v1alpha2/plugins/params"
)

// +kubebuilder:object:generate:=true

// The regex parser plugin
type Regex struct {
	Regex string `json:"regex,omitempty"`
	// Time_Key
	TimeKey string `json:"timeKey,omitempty"`
	// Time_Format, eg. %Y-%m-%dT%H:%M:%S %z
	TimeFormat string `json:"timeFormat,omitempty"`
	// Time_Keep
	TimeKeep *bool  `json:"timeKeep,omitempty"`
	Types    string `json:"types,omitempty"`
}

func (_ *Regex) Name() string {
	return "regex"
}

func (re *Regex) Params(_ plugins.SecretLoader) (*params.KVs, error) {
	kvs := params.NewKVs()
	if re.Regex != "" {
		kvs.Insert("Regex", re.Regex)
	}
	if re.TimeKey != "" {
		kvs.Insert("Time_Key", re.TimeKey)
	}
	if re.TimeFormat != "" {
		kvs.Insert("Time_Format", re.TimeFormat)
	}
	if re.TimeKeep != nil {
		kvs.Insert("Time_Keep", fmt.Sprint(*re.TimeKeep))
	}
	if re.Types != "" {
		kvs.Insert("Types", re.Types)
	}
	return kvs, nil
}
