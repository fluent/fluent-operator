package plugins

import "kubesphere.io/fluentbit-operator/api/fluentbitoperator/v1alpha2/plugins/params"

// +kubebuilder:object:generate=false

// The Plugin interface defines methods for transferring input, filter
// and output plugins to textual section content.
type Plugin interface {
	Name() string
	Params(SecretLoader) (*params.KVs, error)
}
