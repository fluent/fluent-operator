package filter

import "fluent.io/fluent-operator/apis/fluentd/v1alpha1/plugins/common"

type Stdout struct {
	// The format section
	Format *common.Format `json:"format,omitempty"`
	// The inject section
	*common.Inject `json:"inject,omitempty"`
}
