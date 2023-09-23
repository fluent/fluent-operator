package input

import (
	"fmt"

	"github.com/fluent/fluent-operator/v2/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v2/apis/fluentbit/v1alpha2/plugins/params"
)

// +kubebuilder:object:generate:=true

// The tcp input plugin allows to retrieve structured JSON or raw messages over a TCP network interface (TCP port).
// **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/inputs/tcp**
type TCP struct {
	// Listener network interface,default 0.0.0.0
	Listen string `json:"listen,omitempty"`
	// TCP port where listening for connections,default 5170
	// +kubebuilder:validation:Minimum:=1
	// +kubebuilder:validation:Maximum:=65535
	Port *int32 `json:"port,omitempty"`
	// Specify the maximum buffer size in KB to receive a JSON message. If not set, the default size will be the value of Chunk_Size.
	// +kubebuilder:validation:Pattern:="^\\d+(k|K|KB|kb|m|M|MB|mb|g|G|GB|gb)?$"
	BufferSize string `json:"bufferSize,omitempty"`
	// By default the buffer to store the incoming JSON messages, do not allocate the maximum memory allowed, instead it allocate memory when is required.
	//The rounds of allocations are set by Chunk_Size in KB. If not set, Chunk_Size is equal to 32 (32KB).
	// +kubebuilder:validation:Pattern:="^\\d+(k|K|KB|kb|m|M|MB|mb|g|G|GB|gb)?$"
	ChunkSize string `json:"chunkSize,omitempty"`
	// Specify the expected payload format. It support the options json and none.
	// When using json, it expects JSON maps, when is set to none, it will split every record using the defined Separator (option below).
	Format string `json:"format,omitempty"`
	// When the expected Format is set to none, Fluent Bit needs a separator string to split the records. By default it uses the breakline character (LF or 0x10).
	Separator string `json:"separator,omitempty"`
}

func (_ *TCP) Name() string {
	return "tcp"
}

// Params implement Section() method
func (t *TCP) Params(sl plugins.SecretLoader) (*params.KVs, error) {
	kvs := params.NewKVs()
	if t.Listen != "" {
		kvs.Insert("listen", t.Listen)
	}
	if t.Port != nil {
		kvs.Insert("port", fmt.Sprint(*t.Port))
	}
	if t.BufferSize != "" {
		kvs.Insert("Buffer_Size", t.BufferSize)
	}
	if t.ChunkSize != "" {
		kvs.Insert("Chunk_Size", t.ChunkSize)
	}
	if t.Format != "" {
		kvs.Insert("Format", t.Format)
	}
	if t.Separator != "" {
		kvs.Insert("Separator", t.Separator)
	}
	return kvs, nil
}
