# Http

Http defines the parameters for out_http output plugin


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| auth | Auth section for this plugin | *common.Auth |
| endpoint | Endpoint defines the endpoint for HTTP request. If you want to use HTTPS, use https prefix. | *string |
| httpMethod | HttpMethod defines the method for HTTP request. | *string |
| proxy | Proxy defines the proxy for HTTP request. | *string |
| contentType | ContentType defines Content-Type for HTTP request. out_http automatically set Content-Type for built-in formatters when this parameter is not specified. | *string |
| jsonArray | JsonArray defines whether to use the array format of JSON or not | *bool |
| headers | Headers defines the additional headers for HTTP request. | *string |
| headersFromPlaceholders | Additional placeholder based headers for HTTP request. If you want to use tag or record field, use this parameter instead of headers. | *string |
| openTimeout | OpenTimeout defines the connection open timeout in seconds. | *uint16 |
| readTimeout | ReadTimeout defines the read timeout in seconds. | *uint16 |
| sslTimeout | SslTimeout defines the TLS timeout in seconds. | *uint16 |
| tlsCaCertPath | TlsCaCertPath defines the CA certificate path for TLS. | *string |
| tlsClientCertPath | TlsClientCertPath defines the client certificate path for TLS. | *string |
| tlsPrivateKeyPath | TlsPrivateKeyPath defines the client private key path for TLS. | *string |
| tlsPrivateKeyPassphrase | TlsPrivateKeyPassphrase defines the client private key passphrase for TLS. | *string |
| tlsVerifyMode | TlsVerifyMode defines the verify mode of TLS. | *string |
| tlsVersion | TlsVersion defines the default version of TLS transport. | *string |
| tlsCiphers | TlsCiphers defines the cipher suites configuration of TLS. | *string |
| errorResponseAsUnrecoverable | Raise UnrecoverableError when the response code is not SUCCESS. | *bool |
| retryableResponseCodes | The list of retryable response codes. If the response code is included in this list, out_http retries the buffer flush. | *string |
