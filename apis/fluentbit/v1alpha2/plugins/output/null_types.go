package output

import (
	"fluent.io/fluent-operator/apis/fluentbit/v1alpha2/plugins"
	"fluent.io/fluent-operator/apis/fluentbit/v1alpha2/plugins/params"
)

// +kubebuilder:object:generate:=true

// The null output plugin just throws away events.
type Null struct{}

func (_ *Null) Name() string {
	return "null"
}

// implement Section() method
func (_ *Null) Params(_ plugins.SecretLoader) (*params.KVs, error) {
	return nil, nil
}
