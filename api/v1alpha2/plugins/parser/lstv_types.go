package parser

import (
	"fmt"
	"kubesphere.io/fluentbit-operator/api/v1alpha2/plugins"
)

// +kubebuilder:object:generate:=true

// The LSTV parser plugin
type LSTV struct {
	TimeKey    string `json:"timeKey,omitempty"`
	TimeFormat string `json:"timeFormat,omitempty"`
	TimeKeep   *bool  `json:"timeKeep,omitempty"`
	Types      string `json:"types,omitempty"`
}

func (_ *LSTV) Name() string {
	return "ltsv"
}

func (l *LSTV) Params(_ plugins.SecretLoader) (*plugins.KVs, error) {
	kvs := plugins.NewKVs()
	if l.TimeKey != "" {
		kvs.Insert("Time_Key", l.TimeKey)
	}
	if l.TimeFormat != "" {
		kvs.Insert("Time_Format", l.TimeFormat)
	}
	if l.TimeKeep != nil {
		kvs.Insert("Time_Format", fmt.Sprint(*l.TimeKeep))
	}
	if l.Types != "" {
		kvs.Insert("Types", l.Types)
	}
	return kvs, nil
}
