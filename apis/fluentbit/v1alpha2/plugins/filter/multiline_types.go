package filter

import (
	"fmt"

	"github.com/fluent/fluent-operator/v2/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v2/apis/fluentbit/v1alpha2/plugins/params"
)

// +kubebuilder:object:generate:=true

// The Multiline Filter helps to concatenate messages that originally belong to one context but were split across multiple records or log lines. <br />
// **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/filters/multiline-stacktrace**
type Multiline struct {
	plugins.CommonParams `json:",inline"`
	// The Inline struct helps to concatenate messages that originally belong to one context but were split across multiple records or log lines.
	*Multi `json:",inline"`
}

type Multi struct {
	// Specify one or multiple Multiline Parsing definitions to apply to the content.
	//You can specify multiple multiline parsers to detect different formats by separating them with a comma.
	Parser string `json:"parser"`
	//Key name that holds the content to process.
	//Note that a Multiline Parser definition can already specify the key_content to use, but this option allows to overwrite that value for the purpose of the filter.
	KeyContent string `json:"keyContent,omitempty"`
	// +kubebuilder:validation:Enum:=parser;partial_message
	Mode string `json:"mode,omitempty"`
	// +kubebuilder:default:=false
	Buffer bool `json:"buffer,omitempty"`
	// +kubebuilder:default:=2000
	FlushMS int `json:"flushMs,omitempty"`
	// Name for the emitter input instance which re-emits the completed records at the beginning of the pipeline.
	EmitterName string `json:"emitterName,omitempty"`
	// The storage type for the emitter input instance. This option supports the values memory (default) and filesystem.
	// +kubebuilder:validation:Enum:=memory;filesystem
	// +kubebuilder:default:=memory
	EmitterType string `json:"emitterType,omitempty"`
	// Set a limit on the amount of memory in MB the emitter can consume if the outputs provide backpressure. The default for this limit is 10M. The pipeline will pause once the buffer exceeds the value of this setting. For example, if the value is set to 10MB then the pipeline will pause if the buffer exceeds 10M. The pipeline will remain paused until the output drains the buffer below the 10M limit.
	// +kubebuilder:default:=10
	EmitterMemBufLimit int `json:"emitterMemBufLimit,omitempty"`
}

func (_ *Multiline) Name() string {
	return "multiline"
}

func (m *Multiline) Params(_ plugins.SecretLoader) (*params.KVs, error) {
	kvs := params.NewKVs()
	err := m.AddCommonParams(kvs)
	if err != nil {
		return kvs, err
	}
	if m.Multi != nil {
		if m.Multi.Parser != "" {
			kvs.Insert("multiline.parser", m.Multi.Parser)
		}
		if m.Multi.KeyContent != "" {
			kvs.Insert("multiline.key_content", m.Multi.KeyContent)
		}
		if m.Multi.Mode != "" {
			kvs.Insert("mode", m.Multi.Mode)
		}
		if m.Multi.Buffer != false {
			kvs.Insert("buffer", fmt.Sprint(m.Multi.Buffer))
		}
		if m.Multi.FlushMS != 0 {
			kvs.Insert("flush_ms", fmt.Sprint(m.Multi.FlushMS))
		}
		if m.Multi.EmitterName != "" {
			kvs.Insert("emitter_name", m.Multi.EmitterName)
		}
		if m.Multi.EmitterType != "" {
			kvs.Insert("emitter_storage.type", m.Multi.EmitterType)
		}
		if m.Multi.EmitterMemBufLimit != 0 {
			kvs.Insert("emitter_mem_buf_limit", fmt.Sprintf("%dMB", m.Multi.EmitterMemBufLimit))
		}
	}
	return kvs, nil
}
