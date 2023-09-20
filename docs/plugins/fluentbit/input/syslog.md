# Syslog

Syslog input plugins allows to collect Syslog messages through a Unix socket server (UDP or TCP) or over the network using TCP or UDP. <br /> **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/inputs/syslog**


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| mode | Defines transport protocol mode: unix_udp (UDP over Unix socket), unix_tcp (TCP over Unix socket), tcp or udp | string |
| listen | If Mode is set to tcp or udp, specify the network interface to bind, default: 0.0.0.0 | string |
| port | If Mode is set to tcp or udp, specify the TCP port to listen for incoming connections. | *int32 |
| path | If Mode is set to unix_tcp or unix_udp, set the absolute path to the Unix socket file. | string |
| unixPerm | If Mode is set to unix_tcp or unix_udp, set the permission of the Unix socket file, default: 0644 | *int32 |
| parser | Specify an alternative parser for the message. If Mode is set to tcp or udp then the default parser is syslog-rfc5424 otherwise syslog-rfc3164-local is used. If your syslog messages have fractional seconds set this Parser value to syslog-rfc5424 instead. | string |
| bufferChunkSize | By default the buffer to store the incoming Syslog messages, do not allocate the maximum memory allowed, instead it allocate memory when is required. The rounds of allocations are set by Buffer_Chunk_Size. If not set, Buffer_Chunk_Size is equal to 32000 bytes (32KB). | string |
| bufferMaxSize | Specify the maximum buffer size to receive a Syslog message. If not set, the default size will be the value of Buffer_Chunk_Size. | string |
| receiveBufferSize | Specify the maximum socket receive buffer size. If not set, the default value is OS-dependant, but generally too low to accept thousands of syslog messages per second without loss on udp or unix_udp sockets. Note that on Linux the value is capped by sysctl net.core.rmem_max. | string |
| sourceAddressKey | Specify the key where the source address will be injected. | string |
