package parser

import (
	"kubesphere.io/fluentbit-operator/apis/plugins"
	"kubesphere.io/fluentbit-operator/apis/plugins/params"
)

// +kubebuilder:object:generate:=true

// The logfmt parser plugin
type Logfmt struct{}

func (_ *Logfmt) Name() string {
	return "logfmt"
}

func (_ *Logfmt) Params(_ plugins.SecretLoader) (*params.KVs, error) {
	return nil, nil
}
