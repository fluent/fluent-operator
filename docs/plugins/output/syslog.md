# HTTP

The Syslog output plugin allows you to deliver messages to Syslog servers, it supports RFC3164 and RFC5424 formats through different transports such as UDP, TCP or TLS.


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| host |  Host domain or IP address of the remote Syslog server. | string |
| port | TCP or UDP port of the remote Syslog server. | *int32 |
| mode | The desired transport type. | string |
| syslogFormat | The desired Syslog protocol format to use. | string |
| syslogMaxSize | Maximum size allowed per message. | *int32 |
| syslogSeverityKey | Name of the key from the original record that contains the Syslog severity number. | string |
| syslogFacilityKey | Key from the original record that contains the Syslog facility number. | string |
| syslogHostnameKey | Key name from the original record that contains the hostname that generated the message. | string |
| syslogAppnameKey | Key name from the original record that contains the Process ID that generated the message. | string |
| syslogProcessIDKey | Key name from the original record that contains the Message ID associated to the message. | string |
| syslogMessageIDKey | Key name from the original record that contains the Message ID associated to the message. | string |
| syslogSDKey | Key name from the original record that contains the Structured Data (SD) content. | string |
| syslogMessageKey | Key key name that contains the message to deliver. | string |
| tls | HTTP output plugin supports TTL/SSL, for more details about the properties available and general configuration, please refer to the TLS/SSL section. | *[plugins.TLS](../tls.md) |
