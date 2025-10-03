package input

import (
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins/params"
)

// +kubebuilder:object:generate:=true

// Forward defines the in_forward Input plugin that listens to TCP socket to receive the event stream.
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
	// Specify the path to unix socket to receive a forward message. If set, Listen and port are ignnored.
	UnixPath string `json:"unixPath,omitempty"`
	// Set the permission of unix socket file.
	UnixPerm string `json:"unixPerm,omitempty"`
	// Specify maximum buffer memory size used to receive a forward message.
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

func (*Forward) Name() string {
	return "forward"
}

// Params implement Section() method
func (f *Forward) Params(_ plugins.SecretLoader) (*params.KVs, error) {
	kvs := params.NewKVs()

	plugins.InsertKVField(kvs, "Port", f.Port)
	plugins.InsertKVString(kvs, "Listen", f.Listen)
	plugins.InsertKVString(kvs, "Tag", f.Tag)
	plugins.InsertKVString(kvs, "Tag_Prefix", f.TagPrefix)
	plugins.InsertKVString(kvs, "Unix_Path", f.UnixPath)
	plugins.InsertKVString(kvs, "Unix_Perm", f.UnixPerm)
	plugins.InsertKVString(kvs, "Buffer_Chunk_Size", f.BufferChunkSize)
	plugins.InsertKVString(kvs, "Buffer_Max_Size", f.BufferMaxSize)
	plugins.InsertKVString(kvs, "threaded", f.Threaded)

	return kvs, nil
}
