# RecordModifier

The Record Modifier Filter plugin allows to append fields or to exclude specific fields. <br /> RemoveKeys and WhitelistKeys are exclusive. <br /> **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/filters/record-modifier**


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| records | Append fields. This parameter needs key and value pair. | []string |
| removeKeys | If the key is matched, that field is removed. | []string |
| whitelistKeys | If the key is not matched, that field is removed. | []string |
