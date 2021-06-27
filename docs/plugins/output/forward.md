# Forward

Forward is the protocol used by Fluentd to route messages between peers. The forward output plugin allows to provide interoperability between Fluent Bit and Fluentd.


| Field | Description | Scheme | Default |
| ----- | ----------- | ------ | ----- |
| host | Target host where Fluent-Bit or Fluentd are listening for Forward messages. | string | 127.0.0.1 |
| port | TCP Port of the target service. | *int32 | 24224 |
| timeAsInteger | Set timestamps in integer format, it enable compatibility mode for Fluentd v0.12 series. | *bool | False |
| sendOptions | Always send options (with \"size\"=count of messages) | *bool | False |
| requireAckResponse | Send \"chunk\"-option and wait for \"ack\" response from server. Enables at-least-once and receiving server can control rate of traffic. (Requires Fluentd v0.14.0+ server) | *bool | False |
| sharedKey | A key string known by the remote Fluentd used for authorization. | string |  |
| emptySharedKey | Use this option to connect to Fluentd with a zero-length secret. | *bool | False |
| username | Specify the username to present to a Fluentd server that enables user_auth. | *[plugins.Secret](../secret.md) |  |
| password | Specify the password corresponding to the username. | *[plugins.Secret](../secret.md) |  |
| selfHostname | Default value of the auto-generated certificate common name (CN). | string |  |
| tls | Enable or disable TLS support | *[plugins.TLS](../tls.md) | off |

