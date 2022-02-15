# Parser

Parser defines the parameters for filter_parser plugin


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| parse |  | *common.Parse |
| keyName | Specifies the field name in the record to parse. Required parameter. i.e: If set keyName to log, {\"key\":\"value\",\"log\":\"{\"time\":1622473200,\"user\":1}\"} => {\"user\":1} | *string |
| reserveTime | Keeps the original event time in the parsed result. Default is false. | *bool |
| reserveData | Keeps the original key-value pair in the parsed result. Default is false. i.e: If set keyName to log, reverseData to true, {\"key\":\"value\",\"log\":\"{\"user\":1,\"num\":2}\"} => {\"key\":\"value\",\"log\":\"{\"user\":1,\"num\":2}\",\"user\":1,\"num\":2} | *bool |
| removeKeyNameField | Removes key_name field when parsing is succeeded. | *bool |
| replaceInvalidSequence | If true, invalid string is replaced with safe characters and re-parse it. | *bool |
| injectKeyPrefix | Stores the parsed values with the specified key name prefix. | *string |
| hashValueField | Stores the parsed values as a hash value in a field. | *string |
| emitInvalidRecordToError | Emits invalid record to @ERROR label. Invalid cases are: key does not exist;the format is not matched;an unexpected error. If you want to ignore these errors, set false. | *bool |
