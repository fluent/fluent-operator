# Http

Http defines the in_http Input plugin that listens to a TCP socket to receive the event stream.


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| transport | The transport section of http plugin | *common.Transport |
| parse | The parse section of http plugin | *common.Parse |
| port | The port to listen to, default is 9880. | *int32 |
| bind | The port to listen to, default is \"0.0.0.0\" | *string |
| bodySizeLimit | The size limit of the POSTed element. | *string |
| keepaliveTimeout | The timeout limit for keeping the connection alive. | *string |
| addHttpHeaders | Adds HTTP_ prefix headers to the record. | *bool |
| addRemoteAddr | Adds REMOTE_ADDR field to the record. The value of REMOTE_ADDR is the client's address. i.e: X-Forwarded-For: host1, host2 | *string |
| corsAllOrigins | Whitelist domains for CORS. | *string |
| corsAllowCredentials | Add Access-Control-Allow-Credentials header. It's needed when a request's credentials mode is include | *string |
| respondsWithEmptyImg | Responds with an empty GIF image of 1x1 pixel (rather than an empty string). | *bool |
