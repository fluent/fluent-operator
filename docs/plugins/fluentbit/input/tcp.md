# TCP

The tcp input plugin allows to retrieve structured JSON or raw messages over a TCP network interface (TCP port). **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/inputs/tcp**


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| listen | Listener network interface,default 0.0.0.0 | string |
| port | TCP port where listening for connections,default 5170 | *int32 |
| bufferSize | Specify the maximum buffer size in KB to receive a JSON message. If not set, the default size will be the value of Chunk_Size. | string |
| chunkSize | By default the buffer to store the incoming JSON messages, do not allocate the maximum memory allowed, instead it allocate memory when is required. The rounds of allocations are set by Chunk_Size in KB. If not set, Chunk_Size is equal to 32 (32KB). | string |
| format | Specify the expected payload format. It support the options json and none. When using json, it expects JSON maps, when is set to none, it will split every record using the defined Separator (option below). | string |
| separator | When the expected Format is set to none, Fluent Bit needs a separator string to split the records. By default it uses the breakline character (LF or 0x10). | string |
