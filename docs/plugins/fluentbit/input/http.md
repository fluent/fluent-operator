# HTTP

The HTTP input plugin allows you to send custom records to an HTTP endpoint. **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/inputs/http**


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| listen | The address to listen on,default 0.0.0.0 | string |
| port | The port for Fluent Bit to listen on,default 9880 | *int32 |
| tagKey | Specify the key name to overwrite a tag. If set, the tag will be overwritten by a value of the key. | string |
| bufferMaxSize | Specify the maximum buffer size in KB to receive a JSON message,default 4M. | string |
| bufferChunkSize | This sets the chunk size for incoming incoming JSON messages. These chunks are then stored/managed in the space available by buffer_max_size,default 512K. | string |
| successfulResponseCode | It allows to set successful response code. 200, 201 and 204 are supported,default 201. | *int32 |
| successfulHeader | Add an HTTP header key/value pair on success. Multiple headers can be set. Example: X-Custom custom-answer. | string |
| tls |  | *[plugins.TLS](../tls.md) |
