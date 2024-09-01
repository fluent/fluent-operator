# Wasm

Wasm Filter allows you to modify the incoming records using Wasm technology. **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/filters/wasm**


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| wasmPath | Path to the built Wasm program that will be used. This can be a relative path against the main configuration file. | string |
| eventFormat | Define event format to interact with Wasm programs: msgpack or json. Default: json | string |
| functionName | Wasm function name that will be triggered to do filtering. It's assumed that the function is built inside the Wasm program specified above. | string |
| accessiblePaths | Specify the whitelist of paths to be able to access paths from WASM programs. | []string |
| wasmHeapSize | Size of the heap size of Wasm execution. Review unit sizes for allowed values. | string |
| wasmStackSize | Size of the stack size of Wasm execution. Review unit sizes for allowed values. | string |
