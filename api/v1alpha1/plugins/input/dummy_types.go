package input

import (
	"fmt"
	"kubesphere.io/fluentbit-operator/api/v1alpha1/plugins"
)

// +kubebuilder:object:generate:=true

// The dummy input plugin, generates dummy events.
// It is useful for testing, debugging, benchmarking and getting started with Fluent Bit.
type Dummy struct {
	// Tag name associated to all records comming from this plugin.
	Tag string `json:"tag,omitempty"`
	// Dummy JSON record.
	Dummy string `json:"dummy,omitempty"`
	// Events number generated per second.
	Rate *int32 `json:"rate,omitempty"`
}

func (_ *Dummy) Name() string {
	return "dummy"
}

// implement Section() method
func (d *Dummy) Params(_ plugins.SecretLoader) (*plugins.KVs, error) {
	kvs := plugins.NewKVs()
	if d.Tag != "" {
		kvs.Insert("Tag", d.Tag)
	}
	if d.Dummy != "" {
		kvs.Insert("Dummy", d.Dummy)
	}
	if d.Rate != nil {
		kvs.Insert("Rate", fmt.Sprint(*d.Rate))
	}
	return kvs, nil
}
