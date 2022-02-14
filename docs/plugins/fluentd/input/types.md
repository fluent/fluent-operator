# InputCommon

InputCommon defines the common parameters for input plugins


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| id | The @id parameter specifies a unique name for the configuration. | *string |
| logLevel | The @log_level parameter specifies the plugin-specific logging level | *string |
| label | The @label parameter is to route the input events to <label> sections. | *string |
# Input

Input defines the supported types for input plugins


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| forward | in_forward plugin | *Forward |
| http | in_http plugin | *Http |
