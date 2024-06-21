# Networking

Fluent Bit implements a unified networking interface that is exposed to components like plugins. These are the functions from https://docs.fluentbit.io/manual/administration/networking and can be used on various output plugins


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| connectTimeout | Set maximum time expressed in seconds to wait for a TCP connection to be established, this include the TLS handshake time. | *int32 |
| connectTimeoutLogError | On connection timeout, specify if it should log an error. When disabled, the timeout is logged as a debug message. | *bool |
| DNSMode | Select the primary DNS connection type (TCP or UDP). | *string |
| DNSPreferIPv4 | Prioritize IPv4 DNS results when trying to establish a connection. | *bool |
| DNSResolver | Select the primary DNS resolver type (LEGACY or ASYNC). | *string |
| keepalive | Enable or disable connection keepalive support. Accepts a boolean value: on / off. | *string |
| keepaliveIdleTimeout | Set maximum time expressed in seconds for an idle keepalive connection. | *int32 |
| keepaliveMaxRecycle | Set maximum number of times a keepalive connection can be used before it is retired. | *int32 |
| maxWorkerConnections | Set maximum number of TCP connections that can be established per worker. | *int32 |
| sourceAddress | Specify network address to bind for data traffic. | *string |
