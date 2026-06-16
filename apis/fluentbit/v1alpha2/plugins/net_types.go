package plugins

import (
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins/params"
)

// +kubebuilder:object:generate:=true

// Fluent Bit implements a unified networking interface that is exposed to components like plugins. These are the functions from https://docs.fluentbit.io/manual/administration/networking and can be used on various output plugins. These options are configured through each plugin's networking field (for example, the S3 output plugin).
type Networking struct {
	// Set maximum time expressed in seconds to wait for a TCP connection to be established, this include the TLS handshake time.
	ConnectTimeout *int32 `json:"connectTimeout,omitempty"`
	// On connection timeout, specify if it should log an error. When disabled, the timeout is logged as a debug message.
	ConnectTimeoutLogError *bool `json:"connectTimeoutLogError,omitempty"`
	// Set maximum time a connection can stay idle while assigned.
	IOTimeout *int32 `json:"ioTimeout,omitempty"`
	// Select the primary DNS connection type (TCP or UDP).
	// +kubebuilder:validation:Enum:="TCP";"UDP"
	DNSMode *string `json:"DNSMode,omitempty"`
	// Prioritize IPv4 DNS results when trying to establish a connection.
	DNSPreferIPv4 *bool `json:"DNSPreferIPv4,omitempty"`
	// Prioritize IPv6 DNS results when trying to establish a connection.
	DNSPreferIPv6 *bool `json:"DNSPreferIPv6,omitempty"`
	// Select the primary DNS resolver type (LEGACY or ASYNC).
	// +kubebuilder:validation:Enum:="LEGACY";"ASYNC"
	DNSResolver *string `json:"DNSResolver,omitempty"`
	// Enable or disable connection keepalive support. Accepts string enum values: on / off.
	// +kubebuilder:validation:Enum:="on";"off"
	Keepalive *string `json:"keepalive,omitempty"`
	// Set maximum time expressed in seconds for an idle keepalive connection.
	KeepaliveIdleTimeout *int32 `json:"keepaliveIdleTimeout,omitempty"`
	// Set maximum number of times a keepalive connection can be used before it is retired.
	KeepaliveMaxRecycle *int32 `json:"keepaliveMaxRecycle,omitempty"`
	// Set maximum number of TCP connections that can be established per worker.
	MaxWorkerConnections *int32 `json:"maxWorkerConnections,omitempty"`
	// Ignore the environment variables HTTP_PROXY, HTTPS_PROXY and NO_PROXY when set.
	ProxyEnvIgnore *bool `json:"proxyEnvIgnore,omitempty"`
	// Enable or disable Keepalive support. Accepts string enum values: on / off.
	// +kubebuilder:validation:Enum:="on";"off"
	TCPKeepalive *string `json:"tcpKeepalive,omitempty"`
	// Interval between the last data packet sent and the first TCP keepalive probe.
	TCPKeepaliveTime *int32 `json:"tcpKeepaliveTime,omitempty"`
	// Interval between TCP keepalive probes when no response is received on a keepidle probe.
	TCPKeepaliveInterval *int32 `json:"tcpKeepaliveInterval,omitempty"`
	// Number of unacknowledged probes to consider a connection dead.
	TCPKeepaliveProbes *int32 `json:"tcpKeepaliveProbes,omitempty"`
	// Specify network address to bind for data traffic.
	SourceAddress *string `json:"sourceAddress,omitempty"`
}

func (t *Networking) Params(sl SecretLoader) (*params.KVs, error) {
	kvs := params.NewKVs()
	InsertKVField(kvs, "net.connect_timeout", t.ConnectTimeout)
	InsertKVField(kvs, "net.connect_timeout_log_error", t.ConnectTimeoutLogError)
	InsertKVField(kvs, "net.io_timeout", t.IOTimeout)
	InsertKVField(kvs, "net.dns.mode", t.DNSMode)
	InsertKVField(kvs, "net.dns.prefer_ipv4", t.DNSPreferIPv4)
	InsertKVField(kvs, "net.dns.prefer_ipv6", t.DNSPreferIPv6)
	InsertKVField(kvs, "net.dns.resolver", t.DNSResolver)
	InsertKVField(kvs, "net.keepalive", t.Keepalive)
	InsertKVField(kvs, "net.keepalive_idle_timeout", t.KeepaliveIdleTimeout)
	InsertKVField(kvs, "net.keepalive_max_recycle", t.KeepaliveMaxRecycle)
	InsertKVField(kvs, "net.max_worker_connections", t.MaxWorkerConnections)
	InsertKVField(kvs, "net.proxy_env_ignore", t.ProxyEnvIgnore)
	InsertKVField(kvs, "net.tcp_keepalive", t.TCPKeepalive)
	InsertKVField(kvs, "net.tcp_keepalive_interval", t.TCPKeepaliveInterval)
	InsertKVField(kvs, "net.tcp_keepalive_probes", t.TCPKeepaliveProbes)
	InsertKVField(kvs, "net.tcp_keepalive_time", t.TCPKeepaliveTime)
	InsertKVField(kvs, "net.source_address", t.SourceAddress)
	return kvs, nil
}
