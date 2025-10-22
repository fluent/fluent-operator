package input

import (
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins/params"
)

// +kubebuilder:object:generate:=true

// The StatsD input plugin allows you to receive metrics via StatsD protocol.<br />
// **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/inputs/statsd**
type StatsD struct {
	// Listener network interface, default: 0.0.0.0
	Listen string `json:"listen,omitempty"`
	// UDP port where listening for connections, default: 8125
	// +kubebuilder:validation:Minimum:=1
	// +kubebuilder:validation:Maximum:=65535
	Port *int32 `json:"port,omitempty"`
}

func (*StatsD) Name() string {
	return "statsd"
}

// implement Section() method
func (s *StatsD) Params(_ plugins.SecretLoader) (*params.KVs, error) {
	kvs := params.NewKVs()

	plugins.InsertKVString(kvs, "Listen", s.Listen)

	plugins.InsertKVField(kvs, "Port", s.Port)

	return kvs, nil
}
