package output

import (
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins/params"
)

// +kubebuilder:object:generate:=true

// The null output plugin just throws away events.
type Null struct{}

func (*Null) Name() string {
	return "null"
}

// implement Section() method
func (*Null) Params(_ plugins.SecretLoader) (*params.KVs, error) {
	return nil, nil
}
