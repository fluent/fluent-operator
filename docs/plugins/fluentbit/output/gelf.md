# Gelf

The Gelf output plugin allows to send logs in GELF format directly to a Graylog input using TLS, TCP or UDP protocols. <br /> **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/outputs/gelf**


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| host | IP address or hostname of the target Graylog server. | string |
| port | The port that the target Graylog server is listening on. | *int32 |
| mode | The protocol to use (tls, tcp or udp). | string |
| shortMessageKey | ShortMessageKey is the key to use as the short message. | string |
| timestampKey | TimestampKey is the key which its value is used as the timestamp of the message. | string |
| hostKey | HostKey is the key which its value is used as the name of the host, source or application that sent this message. | string |
| fullMessageKey | FullMessageKey is the key to use as the long message that can i.e. contain a backtrace. | string |
| levelKey | LevelKey is the key to be used as the log level. | string |
| packetSize | If transport protocol is udp, it sets the size of packets to be sent. | *int32 |
| compress | If transport protocol is udp, it defines if UDP packets should be compressed. | *bool |
| tls |  | *[plugins.TLS](../tls.md) |
