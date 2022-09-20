package custom

import (
	"errors"
	"fmt"
	"strings"

	"github.com/fluent/fluent-operator/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/apis/fluentbit/v1alpha2/plugins/params"
)

// +kubebuilder:object:generate:=true

// CustomPlugin is used to support filter plugins that are not implemented yet
type CustomPlugin struct {
	Config string `json:"config,omitempty"`
}

func (c *CustomPlugin) Name() string {
	return ""
}

func (a *CustomPlugin) Params(_ plugins.SecretLoader) (*params.KVs, error) {
	kvs := params.NewKVs()
	splits := strings.Split(a.Config, "\n")
	for _, i := range splits {
		if len(i) == 0 {
			continue
		}
		fields := strings.Fields(i)
		if len(fields) < 2 {
			return nil, errors.New(fmt.Sprintf("invalid plugin config: %s", i))
		}
		kvs.Insert(fields[0], strings.Join(fields[1:], " "))
	}
	return kvs, nil
}
