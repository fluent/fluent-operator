package input

import (
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins/params"
)

// +kubebuilder:object:generate:=true

// The dummy input plugin, generates dummy events. <br />
// It is useful for testing, debugging, benchmarking and getting started with Fluent Bit. <br />
// **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/inputs/dummy**
type Dummy struct {
	// Tag name associated to all records coming from this plugin.
	Tag string `json:"tag,omitempty"`
	// Dummy JSON record.
	Dummy string `json:"dummy,omitempty"`
	// Events number generated per second.
	Rate *int32 `json:"rate,omitempty"`
	// Sample events to generate.
	Samples *int32 `json:"samples,omitempty"`
}

func (*Dummy) Name() string {
	return "dummy"
}

// implement Section() method
func (d *Dummy) Params(_ plugins.SecretLoader) (*params.KVs, error) {
	kvs := params.NewKVs()

	plugins.InsertKVString(kvs, "Tag", d.Tag)
	plugins.InsertKVString(kvs, "Dummy", d.Dummy)

	plugins.InsertKVField(kvs, "Rate", d.Rate)
	plugins.InsertKVField(kvs, "Samples", d.Samples)

	return kvs, nil
}
