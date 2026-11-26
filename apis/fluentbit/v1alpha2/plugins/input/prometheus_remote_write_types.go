package input

import (
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins/params"
)

// +kubebuilder:object:generate:=true

// The Prometheus Remote Write input plugin receives metrics sent with the Prometheus remote write protocol, <br />
// for example by Prometheus, Grafana Alloy or Fluent Bit's Prometheus Remote Write output. <br />
// **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/inputs/prometheus-remote-write**
type PrometheusRemoteWrite struct {
	// Tag name associated to all metrics coming from this plugin.
	Tag string `json:"tag,omitempty"`
	// The address to listen on, defaults to 0.0.0.0.
	Listen string `json:"listen,omitempty"`
	// The port for Fluent Bit to listen on, defaults to 8080.
	// Also add it to the ports of the FluentBit resource so that the Fluent Bit Service exposes it.
	// +kubebuilder:validation:Minimum:=1
	// +kubebuilder:validation:Maximum:=65535
	Port *int32 `json:"port,omitempty"`
	// The HTTP path that accepts remote write requests, e.g. /api/prom/push.
	// If not set, requests are accepted on any path.
	// +kubebuilder:validation:Pattern:="^/\\S*$"
	URI string `json:"uri,omitempty"`
	// Maximum size of a remote write request, defaults to 4M.
	// +kubebuilder:validation:Pattern:="^\\d+(k|K|KB|kb|m|M|MB|mb|g|G|GB|gb)?$"
	BufferMaxSize string `json:"bufferMaxSize,omitempty"`
	// Size of the chunks allocated to read a request, up to bufferMaxSize, defaults to 512K.
	// +kubebuilder:validation:Pattern:="^\\d+(k|K|KB|kb|m|M|MB|mb|g|G|GB|gb)?$"
	BufferChunkSize string `json:"bufferChunkSize,omitempty"`
	// HTTP status code returned for a successfully ingested request, defaults to 201.
	// +kubebuilder:validation:Enum:=200;201;204
	SuccessfulResponseCode *int32 `json:"successfulResponseCode,omitempty"`
	// Run the input in its own thread instead of the main event loop.
	Threaded *bool `json:"threaded,omitempty"`
	// TLS settings of the listener. Set crtFile and keyFile to accept remote write requests over HTTPS.
	*plugins.TLS `json:"tls,omitempty"`
}

func (*PrometheusRemoteWrite) Name() string {
	return "prometheus_remote_write"
}

// Params implement Section() method
func (p *PrometheusRemoteWrite) Params(sl plugins.SecretLoader) (*params.KVs, error) {
	kvs := params.NewKVs()

	plugins.InsertKVString(kvs, "tag", p.Tag)
	plugins.InsertKVString(kvs, "listen", p.Listen)
	plugins.InsertKVField(kvs, "port", p.Port)
	plugins.InsertKVString(kvs, "uri", p.URI)
	plugins.InsertKVString(kvs, "buffer_max_size", p.BufferMaxSize)
	plugins.InsertKVString(kvs, "buffer_chunk_size", p.BufferChunkSize)
	plugins.InsertKVField(kvs, "successful_response_code", p.SuccessfulResponseCode)
	plugins.InsertKVField(kvs, "threaded", p.Threaded)

	if p.TLS != nil {
		tls, err := p.TLS.Params(sl)
		if err != nil {
			return nil, err
		}
		kvs.Merge(tls)
	}

	return kvs, nil
}
