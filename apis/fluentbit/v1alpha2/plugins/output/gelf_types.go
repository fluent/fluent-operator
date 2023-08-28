package output

import (
	"fmt"

	"github.com/fluent/fluent-operator/v2/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v2/apis/fluentbit/v1alpha2/plugins/params"
)

// +kubebuilder:object:generate:=true

// The Gelf output plugin allows to send logs in GELF format directly to a Graylog input using TLS, TCP or UDP protocols. <br />
// **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/outputs/gelf**
type Gelf struct {
	// IP address or hostname of the target Graylog server.
	Host string `json:"host,omitempty"`
	// The port that the target Graylog server is listening on.
	// +kubebuilder:validation:Minimum:=1
	// +kubebuilder:validation:Maximum:=65535
	Port *int32 `json:"port,omitempty"`
	// The protocol to use (tls, tcp or udp).
	// +kubebuilder:validation:Enum:=tls;tcp;udp
	Mode string `json:"mode,omitempty"`
	// ShortMessageKey is the key to use as the short message.
	ShortMessageKey string `json:"shortMessageKey,omitempty"`
	// TimestampKey is the key which its value is used as the timestamp of the message.
	TimestampKey string `json:"timestampKey,omitempty"`
	// HostKey is the key which its value is used as the name of the host, source or application that sent this message.
	HostKey string `json:"hostKey,omitempty"`
	// FullMessageKey is the key to use as the long message that can i.e. contain a backtrace.
	FullMessageKey string `json:"fullMessageKey,omitempty"`
	// LevelKey is the key to be used as the log level.
	LevelKey string `json:"levelKey,omitempty"`
	// If transport protocol is udp, it sets the size of packets to be sent.
	PacketSize *int32 `json:"packetSize,omitempty"`
	// If transport protocol is udp, it defines if UDP packets should be compressed.
	Compress     *bool `json:"compress,omitempty"`
	*plugins.TLS `json:"tls,omitempty"`
}

func (_ *Gelf) Name() string {
	return "gelf"
}

func (g *Gelf) Params(sl plugins.SecretLoader) (*params.KVs, error) {
	kvs := params.NewKVs()
	if g.Host != "" {
		kvs.Insert("Host", g.Host)
	}
	if g.Port != nil {
		kvs.Insert("Port", fmt.Sprint(*g.Port))
	}
	if g.Mode != "" {
		kvs.Insert("Mode", g.Mode)
	}
	if g.ShortMessageKey != "" {
		kvs.Insert("Gelf_Short_Message_Key", g.ShortMessageKey)
	}
	if g.TimestampKey != "" {
		kvs.Insert("Gelf_Timestamp_Key", g.TimestampKey)
	}
	if g.HostKey != "" {
		kvs.Insert("Gelf_Host_Key", g.HostKey)
	}
	if g.FullMessageKey != "" {
		kvs.Insert("Gelf_Full_Message_Key", g.FullMessageKey)
	}
	if g.LevelKey != "" {
		kvs.Insert("Gelf_Level_Key", g.LevelKey)
	}
	if g.PacketSize != nil {
		kvs.Insert("Packet_Size", fmt.Sprint(*g.PacketSize))
	}
	if g.Compress != nil {
		kvs.Insert("Compress", fmt.Sprint(*g.Compress))
	}
	if g.TLS != nil {
		tls, err := g.TLS.Params(sl)
		if err != nil {
			return nil, err
		}
		kvs.Merge(tls)
	}
	return kvs, nil
}
