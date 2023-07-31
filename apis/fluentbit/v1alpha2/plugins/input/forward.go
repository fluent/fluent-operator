package input

import (
	"fmt"

	"github.com/fluent/fluent-operator/v2/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v2/apis/fluentbit/v1alpha2/plugins/params"
)

// +kubebuilder:object:generate:=true

// Forward defines the in_forward Input plugin that listens to TCP socket to recieve the event stream.
// **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/inputs/forward**
type Forward struct {
	// Port for forward plugin instance.
	// +kubebuilder:validation:Minimum:=1
	// +kubebuilder:validation:Maximum:=65535
	Port *int32 `json:"port,omitempty"`
	// Listener network interface.
	Listen string `json:"listen,omitempty"`
	// in_forward uses the tag value for incoming logs. If not set it uses tag from incoming log.
	Tag string `json:"tag,omitempty"`
	// Adds the prefix to incoming event's tag
	TagPrefix string `json:"tagPrefix,omitempty"`
	// Specify the path to unix socket to recieve a forward message. If set, Listen and port are ignnored.
	UnixPath string `json:"unixPath,omitempty"`
	// Set the permission of unix socket file.
	UnixPerm string `json:"unixPerm,omitempty"`
	// Specify maximum buffer memory size used to recieve a forward message.
	// The value must be according to the Unit Size specification.
	// +kubebuilder:validation:Pattern:="^\\d+(k|K|KB|kb|m|M|MB|mb|g|G|GB|gb)?$"
	BufferMaxSize string `json:"bufferMaxSize,omitempty"`
	// Set the initial buffer size to store incoming data.
	// This value is used too to increase buffer size as required.
	// The value must be according to the Unit Size specification.
	// +kubebuilder:validation:Pattern:="^\\d+(k|K|KB|kb|m|M|MB|mb|g|G|GB|gb)?$"
	BufferChunkSize string `json:"bufferchunkSize,omitempty"`
	// Threaded mechanism allows input plugin to run in a separate thread which helps to desaturate the main pipeline.
	Threaded string `json:"threaded,omitempty"`
}

func (_ *Forward) Name() string {
	return "forward"
}

// Params implement Section() method
func (f *Forward) Params(_ plugins.SecretLoader) (*params.KVs, error) {
	kvs := params.NewKVs()
	if f.Port != nil {
		kvs.Insert("Port", fmt.Sprint(*f.Port))
	}
	if f.Listen != "" {
		kvs.Insert("Listen", f.Listen)
	}
	if f.Tag != "" {
		kvs.Insert("Tag", f.Tag)
	}
	if f.TagPrefix != "" {
		kvs.Insert("Tag_Prefix", f.TagPrefix)
	}
	if f.UnixPath != "" {
		kvs.Insert("Unix_Path", f.UnixPath)
	}
	if f.UnixPerm != "" {
		kvs.Insert("Unix_Perm", f.UnixPerm)
	}
	if f.BufferChunkSize != "" {
		kvs.Insert("Buffer_Chunk_Size", f.BufferChunkSize)
	}
	if f.BufferMaxSize != "" {
		kvs.Insert("Buffer_Max_Size", f.BufferMaxSize)
	}
	if f.Threaded != "" {
		kvs.Insert("threaded", f.Threaded)
	}
	return kvs, nil
}
