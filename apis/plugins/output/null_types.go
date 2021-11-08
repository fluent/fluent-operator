package output

import (
	"kubesphere.io/fluentbit-operator/apis/plugins"
	"kubesphere.io/fluentbit-operator/apis/plugins/params"
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
