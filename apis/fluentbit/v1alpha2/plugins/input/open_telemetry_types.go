package input

import (
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins/params"
)

// +kubebuilder:object:generate:=true

// The OpenTelemetry plugin allows you to ingest telemetry data as per the OTLP specification, <br />
// from various OpenTelemetry exporters, the OpenTelemetry Collector, or Fluent Bit's OpenTelemetry output plugin. <br />
// **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/inputs/opentelemetry**
type OpenTelemetry struct {
	// The address to listen on,default 0.0.0.0
	Listen string `json:"listen,omitempty"`
	// The port for Fluent Bit to listen on.default 4318.
	// +kubebuilder:validation:Minimum:=1
	// +kubebuilder:validation:Maximum:=65535
	Port *int32 `json:"port,omitempty"`
	// Specify the key name to overwrite a tag. If set, the tag will be overwritten by a value of the key.
	Tagkey string `json:"tagKey,omitempty"`
	// Route trace data as a log message(default false).
	RawTraces *bool `json:"rawTraces,omitempty"`
	// Specify the maximum buffer size in KB to receive a JSON message(default 4M).
	// +kubebuilder:validation:Pattern:="^\\d+(k|K|KB|kb|m|M|MB|mb|g|G|GB|gb)?$"
	BufferMaxSize string `json:"bufferMaxSize,omitempty"`
	// This sets the chunk size for incoming incoming JSON messages. These chunks are then stored/managed in the space available by buffer_max_size(default 512K).
	// +kubebuilder:validation:Pattern:="^\\d+(k|K|KB|kb|m|M|MB|mb|g|G|GB|gb)?$"
	BufferChunkSize string `json:"bufferChunkSize,omitempty"`
	// It allows to set successful response code. 200, 201 and 204 are supported(default 201).
	SuccessfulResponseCode *int32 `json:"successfulResponseCode,omitempty"`
	// opentelemetry uses the tag value for incoming metrics.
	Tag string `json:"tag,omitempty"`
	// If true, tag will be created from uri. e.g. v1_metrics from /v1/metrics
	TagFromURI *bool `json:"tagFromURI,omitempty"`
}

func (*OpenTelemetry) Name() string {
	return "opentelemetry"
}

// implement Section() method
func (ot *OpenTelemetry) Params(_ plugins.SecretLoader) (*params.KVs, error) {
	kvs := params.NewKVs()

	plugins.InsertKVString(kvs, "listen", ot.Listen)
	plugins.InsertKVString(kvs, "tag_key", ot.Tagkey)
	plugins.InsertKVString(kvs, "buffer_max_size", ot.BufferMaxSize)
	plugins.InsertKVString(kvs, "buffer_chunk_size", ot.BufferChunkSize)
	plugins.InsertKVString(kvs, "tag", ot.Tag)

	plugins.InsertKVField(kvs, "port", ot.Port)
	plugins.InsertKVField(kvs, "raw_traces", ot.RawTraces)
	plugins.InsertKVField(kvs, "successful_response_code", ot.SuccessfulResponseCode)
	plugins.InsertKVField(kvs, "tag_from_uri", ot.TagFromURI)

	return kvs, nil
}
