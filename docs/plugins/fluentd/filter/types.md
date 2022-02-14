# FilterCommon

FilterCommon defines the common parameters for filter Plugins


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| logLevel | The @log_level parameter specifies the plugin-specific logging level | *string |
# Filter

Filter defines all types for filter Plugins


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| grep | The filter_grep filter plugin | *Grep |
| recordTransformer | The filter_record_transformer filter plugin | *RecordTransformer |
| parser | The filter_parser filter plugin | *Parser |
| stdout | The filter_stdout filter plugin | *Stdout |
