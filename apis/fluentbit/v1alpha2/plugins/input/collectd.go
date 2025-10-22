package input

import (
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins/params"
)

// +kubebuilder:object:generate:=true

// The Collectd input plugin allows you to receive datagrams from collectd service. <br />
// **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/inputs/collectd**
type Collectd struct {
	// Set the address to listen to, default: 0.0.0.0
	Listen string `json:"listen,omitempty"`
	// Set the port to listen to, default: 25826
	// +kubebuilder:validation:Minimum:=1
	// +kubebuilder:validation:Maximum:=65535
	Port *int32 `json:"port,omitempty"`
	// Set the data specification file,default: /usr/share/collectd/types.db
	TypesDB string `json:"typesDB,omitempty"`
}

func (*Collectd) Name() string {
	return "collectd"
}

// implement Section() method
func (c *Collectd) Params(_ plugins.SecretLoader) (*params.KVs, error) {
	kvs := params.NewKVs()

	plugins.InsertKVString(kvs, "Listen", c.Listen)
	plugins.InsertKVString(kvs, "TypesDB", c.TypesDB)

	plugins.InsertKVField(kvs, "Port", c.Port)

	return kvs, nil
}
