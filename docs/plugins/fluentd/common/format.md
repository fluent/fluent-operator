# FormatCommon

FormatCommon defines common parameters of format Plugins


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| id | The @id parameter specifies a unique name for the configuration. | *string |
| type | The @type parameter specifies the type of the plugin. | *string |
| logLevel | The @log_level parameter specifies the plugin-specific logging level | *string |
# Format

Format defines all types of Format Plugins


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| delimiter | Delimiter for each field. | *string |
| outputTag | Output tag field if true. | *bool |
| outputTime | Output time field if true. | *bool |
| timeType | Overwrites the default value in this plugin. | *string |
| timeFormat | Overwrites the default value in this plugin. | *string |
| newline | Specify newline characters. | *string |
