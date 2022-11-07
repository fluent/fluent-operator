package plugins

import (
	"github.com/fluent/fluent-operator/apis/fluentbit/v1alpha2/plugins/params"
)

// +kubebuilder:object:generate:=true

type CommonParams struct {

	// Alias for the plugin
	Alias string `json:"alias,omitempty"`
	// Integer value to set the maximum number of retries allowed. N must be >= 1
	// (default: 2) or false for no limit at all.
	RetryLimit string `json:"retryLimit,omitempty"`
}

func (c *CommonParams) AddCommonParams(kvs *params.KVs) error {
	if c.Alias != "" {
		kvs.Insert("Alias", c.Alias)
	}
	if c.RetryLimit != "" {
		kvs.Insert("Retry_Limit", c.RetryLimit)
	}
	return nil
}
