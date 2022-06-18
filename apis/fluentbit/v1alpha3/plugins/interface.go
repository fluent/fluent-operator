package plugins

import "github.com/fluent/fluent-operator/apis/fluentbit/v1alpha3/plugins/params"

// +kubebuilder:object:generate=false

// The Plugin interface defines methods for transferring input, filter
// and output plugins to textual section content.
type Plugin interface {
	Name() string
	Params(SecretLoader) (*params.KVs, error)
}
