package clusteroutput

import (
	"kubesphere.io/fluentbit-operator/api/clusterfluentbitoperator/v1alpha2/plugins"
	"kubesphere.io/fluentbit-operator/api/clusterfluentbitoperator/v1alpha2/plugins/params"
)

// +kubebuilder:object:generate:=true

// The stdout clusteroutput plugin allows to print to the standard clusteroutput the data received through the clusterinput plugin.
type Stdout struct {
	// Specify the data format to be printed. Supported formats are msgpack json, json_lines and json_stream.
	// +kubebuilder:validation:Enum:=msgpack;json;json_lines;json_stream
	Format string `json:"format,omitempty"`
	// Specify the name of the date field in clusteroutput.
	JsonDateKey string `json:"jsonDateKey,omitempty"`
	// Specify the format of the date. Supported formats are double,  iso8601 (eg: 2018-05-30T09:39:52.000681Z) and epoch.
	// +kubebuilder:validation:Enum:= double;iso8601;epoch
	JsonDateFormat string `json:"jsonDateFormat,omitempty"`
}

func (_ *Stdout) Name() string {
	return "stdout"
}

// implement Section() method
func (s *Stdout) Params(_ plugins.SecretLoader) (*params.KVs, error) {
	kvs := params.NewKVs()
	if s.Format != "" {
		kvs.Insert("Format", s.Format)
	}
	if s.JsonDateKey != "" {
		kvs.Insert("json_date_key", s.JsonDateKey)
	}
	if s.JsonDateFormat != "" {
		kvs.Insert("json_date_format", s.JsonDateFormat)
	}
	return kvs, nil
}
