package clusteroutput

import (
	"kubesphere.io/fluentbit-operator/api/clusterfluentbitoperator/v1alpha2/plugins"
	"kubesphere.io/fluentbit-operator/api/clusterfluentbitoperator/v1alpha2/plugins/params"
)

// +kubebuilder:object:generate:=true

// The null clusteroutput plugin just throws away events.
type Null struct{}

func (_ *Null) Name() string {
	return "null"
}

// implement Section() method
func (_ *Null) Params(_ plugins.SecretLoader) (*params.KVs, error) {
	return nil, nil
}
