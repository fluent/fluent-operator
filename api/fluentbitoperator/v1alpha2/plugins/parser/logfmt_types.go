package parser

import (
	"kubesphere.io/fluentbit-operator/api/fluentbitoperator/v1alpha2/plugins"
)

// +kubebuilder:object:generate:=true

// The logfmt parser plugin
type Logfmt struct{}

func (_ *Logfmt) Name() string {
	return "logfmt"
}

func (_ *Logfmt) Params(_ plugins.SecretLoader) (*plugins.KVs, error) {
	return nil, nil
}
