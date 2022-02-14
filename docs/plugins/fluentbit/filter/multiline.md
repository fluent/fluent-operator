# Multiline

The Multiline Filter helps to concatenate messages that originally belong to one context but were split across multiple records or log lines. Common examples are stack traces or applications that print logs in multiple lines.


| Field | Description | Scheme |Default|
| ----- | ----------- | ------ | -----|
| multiline.parser | Specify one or multiple Multiline Parser definitions to apply to the content. You can specify multiple multiline parsers to detect different formats by separating them with a comma.  | string ||
| multiline.key_content| Key name that holds the content to process. Note that a Multiline Parser definition can already specify the key_content to use, but this option allows to overwrite that value for the purpose of the filter. | string ||
