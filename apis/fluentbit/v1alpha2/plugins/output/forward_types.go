package output

import (
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins/params"
)

// +kubebuilder:object:generate:=true

// Forward is the protocol used by Fluentd to route messages between peers. <br />
// The forward output plugin allows to provide interoperability between Fluent Bit and Fluentd. <br />
// **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/outputs/forward**
type Forward struct {
	// Target host where Fluent-Bit or Fluentd are listening for Forward messages.
	Host string `json:"host,omitempty"`
	// TCP Port of the target service.
	// +kubebuilder:validation:Minimum:=1
	// +kubebuilder:validation:Maximum:=65535
	Port *int32 `json:"port,omitempty"`
	// Overwrite the tag as we transmit. This allows the receiving pipeline start
	// fresh, or to attribute source.
	Tag string `json:"tag,omitempty"`
	// Set timestamps in integer format, it enable compatibility mode for Fluentd v0.12 series.
	TimeAsInteger *bool `json:"timeAsInteger,omitempty"`
	// Always send options (with "size"=count of messages)
	SendOptions *bool `json:"sendOptions,omitempty"`
	// Send "chunk"-option and wait for "ack" response from server.
	// Enables at-least-once and receiving server can control rate of traffic.
	// (Requires Fluentd v0.14.0+ server)
	RequireAckResponse *bool `json:"requireAckResponse,omitempty"`
	// A key string known by the remote Fluentd used for authorization.
	SharedKey string `json:"sharedKey,omitempty"`
	// Use this option to connect to Fluentd with a zero-length secret.
	EmptySharedKey *bool `json:"emptySharedKey,omitempty"`
	// Specify the username to present to a Fluentd server that enables user_auth.
	Username *plugins.Secret `json:"username,omitempty"`
	// Specify the password corresponding to the username.
	Password *plugins.Secret `json:"password,omitempty"`
	// Default value of the auto-generated certificate common name (CN).
	SelfHostname string `json:"selfHostname,omitempty"`
	*plugins.TLS `json:"tls,omitempty"`
	// Include fluentbit networking options for this output-plugin
	*plugins.Networking `json:"networking,omitempty"`
}

func (*Forward) Name() string {
	return "forward"
}

// implement Section() method
func (f *Forward) Params(sl plugins.SecretLoader) (*params.KVs, error) {
	kvs := params.NewKVs()

	if err := plugins.InsertKVSecret(kvs, "Username", f.Username, sl); err != nil {
		return nil, err
	}
	if err := plugins.InsertKVSecret(kvs, "Password", f.Password, sl); err != nil {
		return nil, err
	}
	if f.TLS != nil {
		tls, err := f.TLS.Params(sl)
		if err != nil {
			return nil, err
		}
		kvs.Merge(tls)
	}
	if f.Networking != nil {
		net, err := f.Networking.Params(sl)
		if err != nil {
			return nil, err
		}
		kvs.Merge(net)
	}

	plugins.InsertKVString(kvs, "Host", f.Host)
	plugins.InsertKVString(kvs, "Tag", f.Tag)
	plugins.InsertKVString(kvs, "Shared_Key", f.SharedKey)
	plugins.InsertKVString(kvs, "Self_Hostname", f.SelfHostname)

	plugins.InsertKVField(kvs, "Port", f.Port)
	plugins.InsertKVField(kvs, "Empty_Shared_Key", f.EmptySharedKey)
	plugins.InsertKVField(kvs, "Time_as_Integer", f.TimeAsInteger)
	plugins.InsertKVField(kvs, "Send_options", f.SendOptions)
	plugins.InsertKVField(kvs, "Require_ack_response", f.RequireAckResponse)

	return kvs, nil
}
