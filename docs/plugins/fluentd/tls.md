# TLS

Fluentd provides integrated support for Transport Layer Security (TLS) and it predecessor Secure Sockets Layer (SSL) respectively.


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| verify | Force certificate validation | *bool |
| debug | Set TLS debug verbosity level. It accept the following values: 0 (No debug), 1 (Error), 2 (State change), 3 (Informational) and 4 Verbose | *int32 |
| caFile | Absolute path to CA certificate file | string |
| caPath | Absolute path to scan for certificate files | string |
| crtFile | Absolute path to Certificate file | string |
| keyFile | Absolute path to private Key file | string |
| keyPassword | Optional password for tls.key_file file | *[Secret](secret.md) |
| vhost | Hostname to be used for TLS SNI extension | string |
