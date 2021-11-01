package clusterinput

import (

	"fmt"
	"kubesphere.io/fluentbit-operator/api/clusterfluentbitoperator/v1alpha2/plugins"
	"kubesphere.io/fluentbit-operator/api/clusterfluentbitoperator/v1alpha2/plugins/params"
)

// +kubebuilder:object:generate:=true

// The dummy clusterinput plugin, generates dummy events.
// It is useful for testing, debugging, benchmarking and getting started with Fluent Bit.
type Dummy struct {
	// Tag name associated to all records comming from this plugin.
	Tag string `json:"tag,omitempty"`
	// Dummy JSON record.
	Dummy string `json:"dummy,omitempty"`
	// Events number generated per second.
	Rate *int32 `json:"rate,omitempty"`
	// Sample events to generate.
	Samples *int32 `json:"samples,omitempty"`
}

func (_ *Dummy) Name() string {
	return "dummy"
}

// implement Section() method
func (d *Dummy) Params(_ plugins.SecretLoader) (*params.KVs, error) {
	kvs := params.NewKVs()
	if d.Tag != "" {
		kvs.Insert("Tag", d.Tag)
	}
	if d.Dummy != "" {
		kvs.Insert("Dummy", d.Dummy)
	}
	if d.Rate != nil {
		kvs.Insert("Rate", fmt.Sprint(*d.Rate))
	}
	if d.Samples != nil {
		kvs.Insert("Samples", fmt.Sprint(*d.Samples))
	}
	return kvs, nil
}
