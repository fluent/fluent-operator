package input

import (
	"fmt"

	"github.com/fluent/fluent-operator/v2/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v2/apis/fluentbit/v1alpha2/plugins/params"
)

// +kubebuilder:object:generate:=true

// The HTTP input plugin allows you to send custom records to an HTTP endpoint.
// **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/inputs/http**
type HTTP struct {
	// The address to listen on,default 0.0.0.0
	Listen string `json:"listen,omitempty"`
	// The port for Fluent Bit to listen on,default 9880
	// +kubebuilder:validation:Minimum:=1
	// +kubebuilder:validation:Maximum:=65535
	Port *int32 `json:"port,omitempty"`
	// Specify the key name to overwrite a tag. If set, the tag will be overwritten by a value of the key.
	Tagkey string `json:"tagKey,omitempty"`
	// Specify the maximum buffer size in KB to receive a JSON message,default 4M.
	// +kubebuilder:validation:Pattern:="^\\d+(k|K|KB|kb|m|M|MB|mb|g|G|GB|gb)?$"
	BufferMaxSize string `json:"bufferMaxSize,omitempty"`
	// This sets the chunk size for incoming incoming JSON messages.
	//These chunks are then stored/managed in the space available by buffer_max_size,default 512K.
	// +kubebuilder:validation:Pattern:="^\\d+(k|K|KB|kb|m|M|MB|mb|g|G|GB|gb)?$"
	BufferChunkSize string `json:"bufferChunkSize,omitempty"`
	// It allows to set successful response code. 200, 201 and 204 are supported,default 201.
	SuccessfulResponseCode *int32 `json:"successfulResponseCode,omitempty"`
	// Add an HTTP header key/value pair on success. Multiple headers can be set. Example: X-Custom custom-answer.
	SuccessfulHeader string `json:"successfulHeader,omitempty"`
	*plugins.TLS     `json:"tls,omitempty"`
}

func (_ *HTTP) Name() string {
	return "http"
}

// Params implement Section() method
func (h *HTTP) Params(sl plugins.SecretLoader) (*params.KVs, error) {
	kvs := params.NewKVs()
	if h.Listen != "" {
		kvs.Insert("listen", h.Listen)
	}
	if h.Port != nil {
		kvs.Insert("port", fmt.Sprint(*h.Port))
	}
	if h.Tagkey != "" {
		kvs.Insert("tag_key", h.Tagkey)
	}
	if h.BufferMaxSize != "" {
		kvs.Insert("buffer_max_size", h.BufferMaxSize)
	}
	if h.BufferChunkSize != "" {
		kvs.Insert("buffer_chunk_size", h.BufferChunkSize)
	}
	if h.SuccessfulResponseCode != nil {
		kvs.Insert("successful_response_code", fmt.Sprint(*h.SuccessfulResponseCode))
	}
	if h.SuccessfulHeader != "" {
		kvs.Insert("success_header", h.SuccessfulHeader)
	}
	if h.TLS != nil {
		tls, err := h.TLS.Params(sl)
		if err != nil {
			return nil, err
		}
		kvs.Merge(tls)
	}
	return kvs, nil
}
