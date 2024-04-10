# Multi




| Field | Description | Scheme |
| ----- | ----------- | ------ |
| parser | Specify one or multiple Multiline Parsing definitions to apply to the content. You can specify multiple multiline parsers to detect different formats by separating them with a comma. | string |
| keyContent | Key name that holds the content to process. Note that a Multiline Parser definition can already specify the key_content to use, but this option allows to overwrite that value for the purpose of the filter. | string |
| mode |  | string |
| buffer |  | bool |
| flushMs |  | int |
| emitterName | Name for the emitter input instance which re-emits the completed records at the beginning of the pipeline. | string |
| emitterType | The storage type for the emitter input instance. This option supports the values memory (default) and filesystem. | string |
| emitterMemBufLimit | Set a limit on the amount of memory in MB the emitter can consume if the outputs provide backpressure. The default for this limit is 10M. The pipeline will pause once the buffer exceeds the value of this setting. For example, if the value is set to 10MB then the pipeline will pause if the buffer exceeds 10M. The pipeline will remain paused until the output drains the buffer below the 10M limit. | int |
