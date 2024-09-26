# FilterCommon

FilterCommon defines the common parameters for the filter plugin.


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| logLevel | The @log_level parameter specifies the plugin-specific logging level | *string |
| tag | Which tag to be matched. | *string |
# Filter

Filter defines all available filter plugins and their parameters.


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| grep | The filter_grep filter plugin | *[Grep](#grep) |
| recordTransformer | The filter_record_transformer filter plugin | *[RecordTransformer](#recordtransformer) |
| parser | The filter_parser filter plugin | *[Parser](#parser) |
| stdout | The filter_stdout filter plugin | *[Stdout](#stdout) |
| customPlugin | Custom plugin type | *[custom.CustomPlugin](plugins/fluentd/custom/custom_plugin.md) |
