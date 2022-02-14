# Forward

Forward defines the out_forward Buffered Output plugin forwards events to other fluentd nodes.


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| servers | Servers defines the servers section, at least one is required | []*common.Server |
| serviceDiscovery | ServiceDiscovery defines the service_discovery section | *common.ServiceDiscovery |
| security | ServiceDiscovery defines the security section | *common.Security |
| requireAckResponse | Changes the protocol to at-least-once. The plugin waits the ack from destination's in_forward plugin. | *bool |
| ackResponseTimeout | This option is used when require_ack_response is true. This default value is based on popular tcp_syn_retries. | *string |
| sendTimeout | The timeout time when sending event logs. | *string |
| connectTimeout | The connection timeout for the socket. When the connection is timed out during the connection establishment, Errno::ETIMEDOUT error is raised. | *string |
| recoverWait | The wait time before accepting a server fault recovery. | *string |
| heartbeatType | Specifies the transport protocol for heartbeats. Set none to disable. | *string |
| heartbeatInterval | The interval of the heartbeat packer. | *string |
| phiFailureDetector | Use the \"Phi accrual failure detector\" to detect server failure. | *bool |
| phiThreshold | The threshold parameter used to detect server faults. | *uint16 |
| hardTimeout | The hard timeout used to detect server failure. The default value is equal to the send_timeout parameter. | *string |
| expireDnsCache | Sets TTL to expire DNS cache in seconds. Set 0 not to use DNS Cache. | *string |
| dnsRoundRobin | Enable client-side DNS round robin. Uniform randomly pick an IP address to send data when a hostname has several IP addresses. heartbeat_type udp is not available with dns_round_robintrue. Use heartbeat_type tcp or heartbeat_type none. | *bool |
| ignoreNetworkErrorsAtStartup | Ignores DNS resolution and errors at startup time. | *bool |
| tlsVersion | The default version of TLS transport. | *string |
| tlsCiphers | The cipher configuration of TLS transport. | *string |
| tlsInsecureMode | Skips all verification of certificates or not. | *bool |
| tlsAllowSelfSignedCert | Allows self-signed certificates or not. | *bool |
| tlsVerifyHostname | Verifies hostname of servers and certificates or not in TLS transport. | *bool |
| tlsCertPath | The additional CA certificate path for TLS. | *string |
| tlsClientCertPath | The client certificate path for TLS. | *string |
| tlsClientPrivateKeyPath | The client private key path for TLS. | *string |
| tlsClientPrivateKeyPassphrase | The TLS private key passphrase for the client. | *string |
| tlsCertThumbprint | The certificate thumbprint for searching from Windows system certstore. This parameter is for Windows only. | *string |
| tlsCertLogicalStoreName | The certificate logical store name on Windows system certstore. This parameter is for Windows only. | *string |
| tlsCertUseEnterpriseStore | Enables the certificate enterprise store on Windows system certstore. This parameter is for Windows only. | *bool |
| keepalive | Enables the keepalive connection. | *bool |
| keepaliveTimeout | Timeout for keepalive. Default value is nil which means to keep the connection alive as long as possible. | *string |
| verifyConnectionAtStartup | Verify that a connection can be made with one of out_forward nodes at the time of startup. | *bool |
