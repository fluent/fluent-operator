package output

import (
	"github.com/fluent/fluent-operator/apis/fluentbit/v1alpha3/plugins"
	"github.com/fluent/fluent-operator/apis/fluentbit/v1alpha3/plugins/params"
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
