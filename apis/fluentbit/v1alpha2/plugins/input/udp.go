package input

import (
	"fmt"

	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins/params"
)

// +kubebuilder:object:generate:=true

// The udp input plugin allows to retrieve structured JSON or raw messages over a UDP network interface (UDP port).
// **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/inputs/udp**

type UDP struct {
	// Listen Listener network interface, default: 0.0.0.0
	Listen *string `json:"listen,omitempty"`
	// Port Specify the UDP port where listening for connections, default: 5170
	// +kubebuilder:validation:Minimum:=1
	// +kubebuilder:validation:Maximum:=65535
	Port *int32 `json:"port,omitempty"`
	// BufferSize Specify the maximum buffer size in KB to receive a JSON message.
	// If not set, the default size will be the value of Chunk_Size.
	// +kubebuilder:validation:Pattern:="^\\d+(k|K|KB|kb|m|M|MB|mb|g|G|GB|gb)?$"
	BufferSize *string `json:"bufferSize,omitempty"`
	// By default the buffer to store the incoming JSON messages, do not allocate the maximum memory allowed,
	// instead it allocate memory when is required.
	// The rounds of allocations are set by Chunk_Size in KB. If not set, Chunk_Size is equal to 32 (32KB).
	// +kubebuilder:validation:Pattern:="^\\d+(k|K|KB|kb|m|M|MB|mb|g|G|GB|gb)?$"
	ChunkSize *string `json:"chunkSize,omitempty"`
	// Format Specify the expected payload format. It support the options json and none.
	// When using json, it expects JSON maps, when is set to none,
	// it will split every record using the defined Separator (option below).
	Format *string `json:"format,omitempty"`
	// Separator When the expected Format is set to none, Fluent Bit needs a separator string to split the records. By default it uses the breakline character (LF or 0x10).
	Separator *string `json:"separator,omitempty"`
	// SourceAddressKey Specify the key where the source address will be injected.
	SourceAddressKey *string `json:"sourceAddressKey,omitempty"`
	// Threaded mechanism allows input plugin to run in a separate thread which helps to desaturate the main pipeline.
	Threaded *string `json:"threaded,omitempty"`
}

func (_ *UDP) Name() string {
	return "udp"
}

func (u *UDP) Params(_ plugins.SecretLoader) (*params.KVs, error) {
	kvs := params.NewKVs()
	if u.Listen != nil {
		kvs.Insert("Listen", *u.Listen)
	}
	if u.Port != nil {
		kvs.Insert("Port", fmt.Sprint(*u.Port))
	}
	if u.BufferSize != nil {
		kvs.Insert("Buffer_Size", *u.BufferSize)
	}
	if u.ChunkSize != nil {
		kvs.Insert("Chunk_Size", *u.ChunkSize)
	}
	if u.Format != nil {
		kvs.Insert("Format", *u.Format)
	}
	if u.Separator != nil {
		kvs.Insert("Separator", *u.Separator)
	}
	if u.SourceAddressKey != nil {
		kvs.Insert("Source_Address_Key", *u.SourceAddressKey)
	}
	if u.Threaded != nil {
		kvs.Insert("Threaded", *u.Threaded)
	}
	return kvs, nil
}
