package parser

import (
	"fmt"
	"kubesphere.io/fluentbit-operator/api/v1alpha2/plugins"
)

// +kubebuilder:object:generate:=true

// The JSON parser plugin
type JSON struct {
	TimeKey    string `json:"timeKey,omitempty"`
	TimeFormat string `json:"timeFormat,omitempty"`
	TimeKeep   *bool  `json:"timeKeep,omitempty"`
}

func (_ *JSON) Name() string {
	return "json"
}

func (j *JSON) Params(_ plugins.SecretLoader) (*plugins.KVs, error) {
	kvs := plugins.NewKVs()
	if j.TimeKey != "" {
		kvs.Insert("Time_Key", j.TimeKey)
	}
	if j.TimeFormat != "" {
		kvs.Insert("Time_Format", j.TimeFormat)
	}
	if j.TimeKeep != nil {
		kvs.Insert("Time_Keep", fmt.Sprint(*j.TimeKeep))
	}
	return kvs, nil
}
