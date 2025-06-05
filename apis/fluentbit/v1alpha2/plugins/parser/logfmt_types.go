package parser

import (
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins/params"
)

// +kubebuilder:object:generate:=true

// The logfmt parser allows to parse the logfmt format described in https://brandur.org/logfmt . <br />
// **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/parsers/logfmt**
type Logfmt struct{
	// Time_Key
	TimeKey string `json:"timeKey,omitempty"`
	// Time_Format, eg. %Y-%m-%dT%H:%M:%S %z
	TimeFormat string `json:"timeFormat,omitempty"`
	// Time_Keep
	TimeKeep *bool `json:"timeKeep,omitempty"`
}

func (_ *Logfmt) Name() string {
	return "logfmt"
}

func (j *Logfmt) Params(_ plugins.SecretLoader) (*params.KVs, error) {
	kvs := params.NewKVs()
	if l.TimeKey != "" {
		kvs.Insert("Time_Key", l.TimeKey)
	}
	if l.TimeFormat != "" {
		kvs.Insert("Time_Format", l.TimeFormat)
	}
	if l.TimeKeep != nil {
		kvs.Insert("Time_Keep", fmt.Sprint(*l.TimeKeep))
	}
	return kvs, nil
}
