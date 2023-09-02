# OpenTelemetry

The OpenTelemetry plugin allows you to ingest telemetry data as per the OTLP specification, <br /> from various OpenTelemetry exporters, the OpenTelemetry Collector, or Fluent Bit's OpenTelemetry output plugin. <br /> **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/inputs/opentelemetry**


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| listen | The address to listen on,default 0.0.0.0 | string |
| port | The port for Fluent Bit to listen on.default 4318. | *int32 |
| tagKey | Specify the key name to overwrite a tag. If set, the tag will be overwritten by a value of the key. | string |
| rawTraces | Route trace data as a log message(default false). | *bool |
| bufferMaxSize | Specify the maximum buffer size in KB to receive a JSON message(default 4M). | string |
| bufferChunkSize | This sets the chunk size for incoming incoming JSON messages. These chunks are then stored/managed in the space available by buffer_max_size(default 512K). | string |
| successfulResponseCode | It allows to set successful response code. 200, 201 and 204 are supported(default 201). | *int32 |
