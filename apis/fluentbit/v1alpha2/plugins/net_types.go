package plugins

import (
	"fmt"

	"github.com/fluent/fluent-operator/v2/apis/fluentbit/v1alpha2/plugins/params"
)

// +kubebuilder:object:generate:=true

// Fluent Bit implements a unified networking interface that is exposed to components like plugins. These are the functions from https://docs.fluentbit.io/manual/administration/networking and can be used on various output plugins
type Net struct {
	// Set maximum time expressed in seconds to wait for a TCP connection to be established, this include the TLS handshake time.
	ConnectTimeout *int32 `json:"connectTimeout,omitempty"`
	// On connection timeout, specify if it should log an error. When disabled, the timeout is logged as a debug message.
	ConnectTimeoutLogError *bool `json:"connectTimeoutLogError,omitempty"`
	// Select the primary DNS connection type (TCP or UDP).
	// +kubebuilder:validation:Enum:="TCP";"UDP"
	DNSMode *string `json:"DNSMode,omitempty"`
	// Prioritize IPv4 DNS results when trying to establish a connection.
	DNSPreferIPv4 *bool `json:"DNSPreferIPv4,omitempty"`
	// Select the primary DNS resolver type (LEGACY or ASYNC).
	// +kubebuilder:validation:Enum:="LEGACY";"ASYNC"
	DNSResolver *string `json:"DNSResolver,omitempty"`
	// Enable or disable connection keepalive support. Accepts a boolean value: on / off.
	// +kubebuilder:validation:Enum:="on";"off"
	Keepalive *string `json:"keepalive,omitempty"`
	// Set maximum time expressed in seconds for an idle keepalive connection.
	KeepaliveIdleTimeout *int32 `json:"keepaliveIdleTimeout,omitempty"`
	// Set maximum number of times a keepalive connection can be used before it is retired.
	KeepaliveMaxRecycle *int32 `json:"keepaliveMaxRecycle,omitempty"`
	// Set maximum number of TCP connections that can be established per worker.
	MaxWorkerConnections *int32 `json:"maxWorkerConnections,omitempty"`
	// Specify network address to bind for data traffic.
	SourceAddress *string `json:"sourceAddress,omitempty"`
}

func (t *Net) Params(sl SecretLoader) (*params.KVs, error) {
	kvs := params.NewKVs()
	if t.ConnectTimeout != nil {
		kvs.Insert("net.connect_timeout", fmt.Sprint(*t.ConnectTimeout))
	}
	if t.ConnectTimeoutLogError != nil {
		kvs.Insert("net.connect_timeout_log_error", fmt.Sprint(*t.ConnectTimeoutLogError))
	}
	if t.DNSMode != nil {
		kvs.Insert("net.dns.mode", *t.DNSMode)
	}
	if t.DNSPreferIPv4 != nil {
		kvs.Insert("net.dns.prefer_ipv4", fmt.Sprint(*t.DNSPreferIPv4))
	}
	if t.DNSResolver != nil {
		kvs.Insert("net.dns.prefer_ipv4", *t.DNSResolver)
	}
	if t.Keepalive != nil {
		kvs.Insert("net.keepalive", *t.Keepalive)
	}
	if t.KeepaliveIdleTimeout != nil {
		kvs.Insert("net.keepalive_idle_timeout", fmt.Sprint(*t.KeepaliveIdleTimeout))
	}
	if t.KeepaliveMaxRecycle != nil {
		kvs.Insert("net.keepalive_max_recycle", fmt.Sprint(*t.KeepaliveMaxRecycle))
	}
	if t.MaxWorkerConnections != nil {
		kvs.Insert("net.max_worker_connections", fmt.Sprint(*t.MaxWorkerConnections))
	}
	if t.SourceAddress != nil {
		kvs.Insert("net.source_address", *t.SourceAddress)
	}
	return kvs, nil
}
