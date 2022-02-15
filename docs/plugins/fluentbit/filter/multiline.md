# Multi




| Field | Description | Scheme |
| ----- | ----------- | ------ |
| parser | Specify one or multiple Multiline Parsing definitions to apply to the content. You can specify multiple multiline parsers to detect different formats by separating them with a comma. | string |
| keyContent | Key name that holds the content to process. Note that a Multiline Parser definition can already specify the key_content to use, but this option allows to overwrite that value for the purpose of the filter. | string |
