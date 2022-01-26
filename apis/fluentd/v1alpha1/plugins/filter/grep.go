package filter

import (
	"fmt"

	"fluent.io/fluent-operator/apis/fluentd/v1alpha1/plugins"
	"fluent.io/fluent-operator/apis/fluentd/v1alpha1/plugins/params"
)

type Grep struct {
	Regexps  []*Regexp  `json:"regexp,omitempty"`
	Excludes []*Exclude `json:"exclude,omitempty"`
	Ands     []*And     `json:"and,omitempty"`
	Ors      []*Or      `json:"or,omitempty"`
}

type Regexp struct {
	Key     *string `json:"key,omitempty"`
	Pattern *string `json:"pattern,omitempty"`
}

type Exclude struct {
	Key     *string `json:"key,omitempty"`
	Pattern *string `json:"pattern,omitempty"`
}

type And struct {
	*Regexp  `json:"regexp,omitempty"`
	*Exclude `json:"exclude,omitempty"`
}

type Or struct {
	*Regexp  `json:"regexp,omitempty"`
	*Exclude `json:"exclude,omitempty"`
}

func (r *Regexp) Name() string {
	return "regexp"
}

func (e *Exclude) Name() string {
	return "exclude"
}

func (r *Regexp) Params(_ plugins.SecretLoader) (*params.PluginStore, error) {
	ps := params.NewPluginStore(r.Name())
	ps.InsertPairs("key", fmt.Sprint(*r.Key))
	ps.InsertPairs("pattern", fmt.Sprint(*r.Pattern))
	return ps, nil
}

func (e *Exclude) Params(_ plugins.SecretLoader) (*params.PluginStore, error) {
	ps := params.NewPluginStore(e.Name())
	ps.InsertPairs("key", fmt.Sprint(*e.Key))
	ps.InsertPairs("pattern", fmt.Sprint(*e.Pattern))
	return ps, nil
}
