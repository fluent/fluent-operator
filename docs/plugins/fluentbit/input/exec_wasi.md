# ExecWasi

The exec_wasi input plugin, allows to execute WASM program that is WASI target like as external program and collects event logs from there. **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/inputs/exec-wasi**


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| wasiPath | The place of a WASM program file. | string |
| parser | Specify the name of a parser to interpret the entry as a structured message. | string |
| accessiblePaths | Specify the whitelist of paths to be able to access paths from WASM programs. | []string |
| intervalSec | Polling interval (seconds). | *int32 |
| intervalNSec | Polling interval (nanoseconds). | *int64 |
| wasmHeapSize |  | string |
| wasmStackSize | Size of the stack size of Wasm execution. Review unit sizes for allowed values. | string |
| bufSize | Size of the buffer (check unit sizes for allowed values) | string |
| threaded | Indicates whether to run this input in its own thread. Default: false. | *bool |
