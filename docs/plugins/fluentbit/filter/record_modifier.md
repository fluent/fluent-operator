# RecordModifier

The Record Modifier Filter plugin allows to append fields or to exclude specific fields. <br /> RemoveKeys and WhitelistKeys are exclusive. <br /> **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/filters/record-modifier**


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| records | Append fields. This parameter needs key and value pair. | []string |
| removeKeys | If the key is matched, that field is removed. | []string |
| allowlistKeys | If the key is not matched, that field is removed. | []string |
| whitelistKeys | An alias of allowlistKeys for backwards compatibility. | []string |
| uuidKeys | If set, the plugin appends uuid to each record. The value assigned becomes the key in the map. | []string |
