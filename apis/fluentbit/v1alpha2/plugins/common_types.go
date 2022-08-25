package plugins

import "github.com/fluent/fluent-operator/apis/fluentbit/v1alpha2/plugins/params"

// +kubebuilder:object:generate:=true

type CommonParams struct {

	// Alias for the plugin
	Alias string `json:"alias,omitempty"`
}

func (c *CommonParams) ParseParams(sl SecretLoader) (*params.KVs, error) {
	kvs := params.NewKVs()
	if c.Alias != "" {
		kvs.Insert("Alias", c.Alias)
	}
	return kvs, nil
}
