# Record

The parameters inside <record> directives are considered to be new key-value pairs


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| key | New field can be defined as key | *string |
| value | The value must from Record properties. See https://docs.fluentd.org/filter/record_transformer#less-than-record-greater-than-directive | *string |
# RecordTransformer

RecordTransformer defines the parameters for filter_record_transformer plugin


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| records |  | []*Record |
| enableRuby | When set to true, the full Ruby syntax is enabled in the ${...} expression. The default value is false. i.e: jsonized_record ${record.to_json} | *bool |
| autoTypecast | Automatically casts the field types. Default is false. This option is effective only for field values comprised of a single placeholder. | *bool |
| renewRecord | By default, the record transformer filter mutates the incoming data. However, if this parameter is set to true, it modifies a new empty hash instead. | *bool |
| renewTimeKey | renew_time_key foo overwrites the time of events with a value of the record field foo if exists. The value of foo must be a Unix timestamp. | *string |
| keepKeys | A list of keys to keep. Only relevant if renew_record is set to true. | *string |
| removeKeys | A list of keys to delete. Supports nested field via record_accessor syntax since v1.1.0. | *string |
