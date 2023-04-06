# Forward

Forward is the protocol used by Fluentd to route messages between peers. <br /> The forward output plugin allows to provide interoperability between Fluent Bit and Fluentd. <br /> **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/outputs/forward**


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| host | Target host where Fluent-Bit or Fluentd are listening for Forward messages. | string |
| port | TCP Port of the target service. | *int32 |
| timeAsInteger | Set timestamps in integer format, it enable compatibility mode for Fluentd v0.12 series. | *bool |
| sendOptions | Always send options (with \"size\"=count of messages) | *bool |
| requireAckResponse | Send \"chunk\"-option and wait for \"ack\" response from server. Enables at-least-once and receiving server can control rate of traffic. (Requires Fluentd v0.14.0+ server) | *bool |
| sharedKey | A key string known by the remote Fluentd used for authorization. | string |
| emptySharedKey | Use this option to connect to Fluentd with a zero-length secret. | *bool |
| username | Specify the username to present to a Fluentd server that enables user_auth. | *[plugins.Secret](../secret.md) |
| password | Specify the password corresponding to the username. | *[plugins.Secret](../secret.md) |
| selfHostname | Default value of the auto-generated certificate common name (CN). | string |
| tls |  | *[plugins.TLS](../tls.md) |
