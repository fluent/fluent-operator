package input

import (
	"strconv"

	"github.com/fluent/fluent-operator/v2/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v2/apis/fluentbit/v1alpha2/plugins/params"
)

// +kubebuilder:objct:generate:=true
// Forward defines the in_forward Input plugin that listens to TCP socket to recieve the event stream
type Forward struct {
	Port            int32  `json:"port,omitempty"`
	Listen          string `json:"listen,omitempty"`
	// in_forward uses the tag value for incoming logs. If not set it uses tag from incoming log.
	Tag             string `json:"tag,omitempty"`
	// Adds the prefix to incoming event's tag
	TagPrefix       string `json:"tagPrefix,omitempty"`
	// Specify the path to unix socket to recieve a forward message. If set, Listen and port are ignnored.
	UnixPath        string `json:"unixPath,omitempty"`
	// Set the permission of unix socket file.
	UnixPerm        string `json:"unixPerm,omitempty"`
	// Specify maximum buffer memory size used to recieve a forward message. Provide value in bytes.
	BufferMaxSize   int    `json:"bufferMaxSize,omitempty"`
	BufferChunkSize int    `json:"bufferchunkSize,omitempty"`
	// Threaded mechanism allows input plugin to run in a separate thread which helps to desaturate the main pipeline.
	Threaded        string `json:"threaded,omitempty"`
}

func (_ *Forward) Name() string {
	return "forward"
}

// Params implement Section() method
func (f *Forward) Params(_ plugins.SecretLoader) (*params.KVs, error) {
	kvs := params.NewKVs()
	if f.Port >= 0 && f.Port < 65535 {
		kvs.Insert("Port", strconv.Itoa(int(f.Port)))
	}
	if f.Listen != "" {
		kvs.Insert("Listen", f.Listen)
	}
	if f.Tag != "" {
		kvs.Insert("Tag", f.Tag)
	}
	if f.TagPrefix != "" {
		kvs.Insert("TagPrefix", f.TagPrefix)
	}
	if f.UnixPath != "" {
		kvs.Insert("UnixPath", f.UnixPath)
	}
	if f.UnixPerm != "" {
		kvs.Insert("UnixPerm", f.UnixPerm)
	}
	if f.BufferChunkSize > 0 {
		kvs.Insert("BufferChunkSize", strconv.Itoa(f.BufferChunkSize))
	}
	if f.BufferMaxSize > 0 {
		kvs.Insert("BufferMaxSize", strconv.Itoa(f.BufferMaxSize))
	}
	if f.Threaded != "" {
		kvs.Insert("threaded", f.Threaded)
	}
	return kvs, nil
}
