# Syslog

Syslog output plugin allows you to deliver messages to Syslog servers. <br /> **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/outputs/syslog**


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| host | Host domain or IP address of the remote Syslog server. | string |
| port | TCP or UDP port of the remote Syslog server. | *int32 |
| mode | Mode of the desired transport type, the available options are tcp, tls and udp. | string |
| syslogFormat | Syslog protocol format to use, the available options are rfc3164 and rfc5424. | string |
| syslogMaxSize | Maximum size allowed per message, in bytes. | *int32 |
| syslogSeverityKey | Key from the original record that contains the Syslog severity number. | string |
| syslogFacilityKey | Key from the original record that contains the Syslog facility number. | string |
| syslogHostnameKey | Key name from the original record that contains the hostname that generated the message. | string |
| syslogAppnameKey | Key name from the original record that contains the application name that generated the message. | string |
| syslogProcessIDKey | Key name from the original record that contains the Process ID that generated the message. | string |
| syslogMessageIDKey | Key name from the original record that contains the Message ID associated to the message. | string |
| syslogSDKey | Key name from the original record that contains the Structured Data (SD) content. | string |
| syslogMessageKey | Key key name that contains the message to deliver. | string |
| tls | Syslog output plugin supports TTL/SSL, for more details about the properties available and general configuration, please refer to the TLS/SSL section. | *[plugins.TLS](../tls.md) |
