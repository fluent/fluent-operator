package input

import (
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins/params"
)

// +kubebuilder:object:generate:=true

// Syslog input plugins allows to collect Syslog messages through a Unix socket server (UDP or TCP) or over the network using TCP or UDP. <br />
// **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/inputs/syslog**
type Syslog struct {
	// Defines transport protocol mode: unix_udp (UDP over Unix socket), unix_tcp (TCP over Unix socket), tcp or udp
	// +kubebuilder:validation:Enum:=unix_udp;unix_tcp;tcp;udp
	Mode string `json:"mode,omitempty"`
	// If Mode is set to tcp or udp, specify the network interface to bind, default: 0.0.0.0
	Listen string `json:"listen,omitempty"`
	// If Mode is set to tcp or udp, specify the TCP port to listen for incoming connections.
	// +kubebuilder:validation:Minimum:=1
	// +kubebuilder:validation:Maximum:=65535
	Port *int32 `json:"port,omitempty"`
	// If Mode is set to unix_tcp or unix_udp, set the absolute path to the Unix socket file.
	Path string `json:"path,omitempty"`
	// If Mode is set to unix_tcp or unix_udp, set the permission of the Unix socket file, default: 0644
	UnixPerm *int32 `json:"unixPerm,omitempty"`
	// Specify an alternative parser for the message. If Mode is set to tcp or udp then the default parser is syslog-rfc5424 otherwise syslog-rfc3164-local is used.
	// If your syslog messages have fractional seconds set this Parser value to syslog-rfc5424 instead.
	Parser string `json:"parser,omitempty"`
	// By default the buffer to store the incoming Syslog messages, do not allocate the maximum memory allowed, instead it allocate memory when is required.
	// The rounds of allocations are set by Buffer_Chunk_Size. If not set, Buffer_Chunk_Size is equal to 32000 bytes (32KB).
	// +kubebuilder:validation:Pattern:="^\\d+(k|K|KB|kb|m|M|MB|mb|g|G|GB|gb)?$"
	BufferChunkSize string `json:"bufferChunkSize,omitempty"`
	// Specify the maximum buffer size to receive a Syslog message. If not set, the default size will be the value of Buffer_Chunk_Size.
	// +kubebuilder:validation:Pattern:="^\\d+(k|K|KB|kb|m|M|MB|mb|g|G|GB|gb)?$"
	BufferMaxSize string `json:"bufferMaxSize,omitempty"`
	// Specify the maximum socket receive buffer size. If not set, the default value is OS-dependant,
	// but generally too low to accept thousands of syslog messages per second without loss on udp or unix_udp sockets. Note that on Linux the value is capped by sysctl net.core.rmem_max.
	// +kubebuilder:validation:Pattern:="^\\d+(k|K|KB|kb|m|M|MB|mb|g|G|GB|gb)?$"
	ReceiveBufferSize string `json:"receiveBufferSize,omitempty"`
	// Specify the key where the source address will be injected.
	SourceAddressKey string `json:"sourceAddressKey,omitempty"`
	// Specify TLS connector options.
	*plugins.TLS `json:"tls,omitempty"`
}

func (*Syslog) Name() string {
	return "syslog"
}

func (s *Syslog) Params(sl plugins.SecretLoader) (*params.KVs, error) {
	kvs := params.NewKVs()

	plugins.InsertKVString(kvs, "Mode", s.Mode)
	plugins.InsertKVString(kvs, "Listen", s.Listen)
	plugins.InsertKVString(kvs, "Path", s.Path)
	plugins.InsertKVString(kvs, "Parser", s.Parser)
	plugins.InsertKVString(kvs, "Buffer_Chunk_Size", s.BufferChunkSize)
	plugins.InsertKVString(kvs, "Buffer_Max_Size", s.BufferMaxSize)
	plugins.InsertKVString(kvs, "Receive_Buffer_Size", s.ReceiveBufferSize)
	plugins.InsertKVString(kvs, "Source_Address_Key", s.SourceAddressKey)

	plugins.InsertKVField(kvs, "Port", s.Port)
	plugins.InsertKVField(kvs, "Unix_Perm", s.UnixPerm)

	if s.TLS != nil {
		tls, err := s.TLS.Params(sl)
		if err != nil {
			return nil, err
		}
		kvs.Merge(tls)
	}

	return kvs, nil
}
