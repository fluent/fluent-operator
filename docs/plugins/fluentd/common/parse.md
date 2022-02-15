# ParseCommon

ParseCommon defines the common parameters for the parse plugin


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| id | The @id parameter specifies a unique name for the configuration. | *string |
| type | The @type parameter specifies the type of the plugin. | *string |
| logLevel | The @log_level parameter specifies the plugin-specific logging level | *string |
# Parse

Parse defines various parameters for the parse plugin


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| expression | Specifies the regular expression for matching logs. Regular expression also supports i and m suffix. | *string |
| types | Specify types for converting field into another, i.e: types user_id:integer,paid:bool,paid_usd_amount:float | *string |
| timeKey | Specify time field for event time. If the event doesn't have this field, current time is used. | *string |
| estimateCurrentEvent | If true, use Fluent::Eventnow(current time) as a timestamp when time_key is specified. | *bool |
| keepTimeKey | If true, keep time field in th record. | *bool |
| timeout | Specify timeout for parse processing. | *string |
