package input

import (
	"fmt"

	"github.com/fluent/fluent-operator/v2/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v2/apis/fluentbit/v1alpha2/plugins/params"
)

// +kubebuilder:object:generate:=true

// The MQTT input plugin, allows to retrieve messages/data from MQTT control packets over a TCP connection. <br />
// The incoming data to receive must be a JSON map. <br />
// **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/inputs/mqtt**
type MQTT struct {
	// Listener network interface, default: 0.0.0.0
	Listen string `json:"listen,omitempty"`
	// TCP port where listening for connections, default: 1883
	// +kubebuilder:validation:Minimum:=1
	// +kubebuilder:validation:Maximum:=65535
	Port *int32 `json:"port,omitempty"`
}

func (_ *MQTT) Name() string {
	return "mqtt"
}

// implement Section() method
func (m *MQTT) Params(_ plugins.SecretLoader) (*params.KVs, error) {
	kvs := params.NewKVs()
	if m.Listen != "" {
		kvs.Insert("Listen", m.Listen)
	}
	if m.Port != nil {
		kvs.Insert("Port", fmt.Sprint(*m.Port))
	}
	return kvs, nil
}
