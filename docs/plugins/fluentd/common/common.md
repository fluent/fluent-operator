# CommonFields

CommonFields defines the common parameters for all plugins


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| id | The @id parameter specifies a unique name for the configuration. | *string |
| type | The @type parameter specifies the type of the plugin. | *string |
| logLevel | The @log_level parameter specifies the plugin-specific logging level | *string |
# Time

Time defines the common parameters for the time plugin


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| timeType | parses/formats value according to this type, default is string | *string |
| timeFormat | Process value according to the specified format. This is available only when time_type is string | *string |
| localtime | If true, uses local time. | *bool |
| utc | If true, uses UTC. | *bool |
| timezone | Uses the specified timezone. | *string |
| timeFormatFallbacks | Uses the specified time format as a fallback in the specified order. You can parse undetermined time format by using time_format_fallbacks. This options is enabled when time_type is mixed. | *string |
# Inject

Inject defines the common parameters for the inject plugin The inject section can be under <match> or <filter> section.


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| hostnameKey | The field name to inject hostname | *string |
| hostname | Hostname value | *string |
| workerIdKey | The field name to inject worker_id | *string |
| tagKey | The field name to inject tag | *string |
| timeKey | The field name to inject time | *string |
# Security

Security defines the common parameters for the security plugin


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| selfHostname | The hostname. | *string |
| sharedKey | The shared key for authentication. | *string |
| userAuth | If true, user-based authentication is used. | *string |
| allowAnonymousSource | Allows the anonymous source. <client> sections are required, if disabled. | *string |
| user | Defines user section directly. | *User |
# User

User defines the common parameters for the user plugin


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| username |  | *[plugins.Secret](../secret.md) |
| password |  | *[plugins.Secret](../secret.md) |
# Transport

Transport defines the commont parameters for the transport plugin


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| protocol | The protocal name of this plugin, i.e: tls | *string |
| version |  | *string |
| ciphers |  | *string |
| insecure |  | *bool |
| caPath | for Cert signed by public CA | *string |
| certPath |  | *string |
| privateKeyPath |  | *string |
| privateKeyPassphrase |  | *string |
| clientCertAuth |  | *bool |
| caCertPath | for Cert generated | *string |
| caPrivateKeyPath |  | *string |
| caPrivateKeyPassphrase |  | *string |
| certVerifier | other parameters | *string |
# Client

Client defines the commont parameters for the client plugin


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| host | The IP address or hostname of the client. This is exclusive with Network. | *string |
| network | The network address specification. This is exclusive with Host. | *string |
| sharedKey | The shared key per client. | *string |
| users | The array of usernames. | *string |
# Auth

Auth defines the common parameters for the auth plugin


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| auth | The method for HTTP authentication. Now only basic. | *string |
| username | The username for basic authentication. | *[plugins.Secret](../secret.md) |
| password | The password for basic authentication. | *[plugins.Secret](../secret.md) |
# Server

Server defines the common parameters for the server plugin


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| host | Host defines the IP address or host name of the server. | *string |
| name | Name defines the name of the server. Used for logging and certificate verification in TLS transport (when the host is the address). | *string |
| port | Port defines the port number of the host. Note that both TCP packets (event stream) and UDP packets (heartbeat messages) are sent to this port. | *string |
| sharedKey | SharedKey defines the shared key per server. | *string |
| username | Username defines the username for authentication. | *[plugins.Secret](../secret.md) |
| password | Password defines the password for authentication. | *[plugins.Secret](../secret.md) |
| standby | Standby marks a node as the standby node for an Active-Standby model between Fluentd nodes. | *string |
| weight | Weight defines the load balancing weight | *string |
# SDCommon

SDCommon defines the common parameters for the ServiceDiscovery plugin


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| id | The @id parameter specifies a unique name for the configuration. | *string |
| type | The @type parameter specifies the type of the plugin. | *string |
| logLevel | The @log_level parameter specifies the plugin-specific logging level | *string |
# ServiceDiscovery

ServiceDiscovery defines various parameters for the ServiceDiscovery plugin. Fluentd has a pluggable system called Service Discovery that lets the user extend and reuse custom output service discovery.


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| server | The server section of this plugin | *Server |
# FileServiceDiscovery

FileServiceDiscovery defines the file type for the ServiceDiscovery plugin


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| path | The path of the target list. Default is '/etc/fluent/sd.yaml' | *string |
| confEncoding | The encoding of the configuration file. | *string |
# SrvServiceDiscovery

SrvServiceDiscovery defines the srv type for the ServiceDiscovery plugin


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| service | Service without the underscore in RFC2782. | *string |
| proto | Proto without the underscore in RFC2782. | *string |
| hostname | The name in RFC2782. | *string |
| dnsServerHost | DnsServerHost defines the hostname of the DNS server to request the SRV record. | *string |
| interval | Interval defines the interval of sending requests to DNS server. | *string |
| dnsLookup | DnsLookup resolves the hostname to IP address of the SRV's Target. | *string |
