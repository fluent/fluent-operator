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
| grokPattern | The pattern of grok. | *string |
| customPatternPath | Path to the file that includes custom grok patterns. | *string |
| grokFailureKey | The key has grok failure reason. | *string |
| multiLineStartRegexp | The regexp to match beginning of multiline. This is only for \"multiline_grok\". | *string |
| grokPatternSeries | Specify grok pattern series set. | *string |
| grok | Grok Sections | []Grok |
# Grok




| Field | Description | Scheme |
| ----- | ----------- | ------ |
| name | The name of this grok section. | *string |
| pattern | The pattern of grok. Required parameter. | *string |
| keepTimeKey | If true, keep time field in the record. | *bool |
| timeKey | Specify time field for event time. If the event doesn't have this field, current time is used. | *string |
| timeFormat | Process value using specified format. This is available only when time_type is string | *string |
| timeZone | Use specified timezone. one can parse/format the time value in the specified timezone. | *string |
