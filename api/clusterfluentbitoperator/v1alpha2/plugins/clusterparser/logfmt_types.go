package clusterparser

import (
	"kubesphere.io/fluentbit-operator/api/clusterfluentbitoperator/v1alpha2/plugins"
	"kubesphere.io/fluentbit-operator/api/clusterfluentbitoperator/v1alpha2/plugins/params"
)

// +kubebuilder:object:generate:=true

// The logfmt clusterparser plugin
type Logfmt struct{}

func (_ *Logfmt) Name() string {
	return "logfmt"
}

func (_ *Logfmt) Params(_ plugins.SecretLoader) (*params.KVs, error) {
	return nil, nil
}
