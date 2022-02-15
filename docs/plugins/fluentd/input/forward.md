# Forward

Forward defines the in_forward Input plugin that listens to a TCP socket to receive the event stream.


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| transport | The transport section of forward plugin | *common.Transport |
| security | The security section of forward plugin | *common.Security |
| client | The security section of client plugin | *common.Client |
| user | The security section of user plugin | *common.User |
| port | The port to listen to, default is 24224. | *int32 |
| bind | The port to listen to, default is \"0.0.0.0\" | *string |
| tag | in_forward uses incoming event's tag by default (See Protocol Section). If the tag parameter is set, its value is used instead. | *string |
| addTagPrefix | Adds the prefix to the incoming event's tag. | *string |
| lingerTimeout | The timeout used to set the linger option. | *uint16 |
| resolveHostname | Tries to resolve hostname from IP addresses or not. | *bool |
| denyKeepalive | The connections will be disconnected right after receiving a message, if true. | *bool |
| sendKeepalivePacket | Enables the TCP keepalive for sockets. | *bool |
| chunkSizeLimit | The size limit of the received chunk. If the chunk size is larger than this value, the received chunk is dropped. | *string |
| chunkSizeWarnLimit | The warning size limit of the received chunk. If the chunk size is larger than this value, a warning message will be sent. | *string |
| skipInvalidEvent | Skips the invalid incoming event. | *bool |
| sourceAddressKey | The field name of the client's source address. If set, the client's address will be set to its key. | *string |
| sourceHostnameKey | The field name of the client's hostname. If set, the client's hostname will be set to its key. | *string |
