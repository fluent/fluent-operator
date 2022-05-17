# RecordModifier

The Record Modifier Filter plugin allows to append fields or to exclude specific fields. RemoveKeys and WhitelistKeys are exclusive.


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| records | Append fields. This parameter needs key and value pair. | []string |
| removeKeys | If the key is matched, that field is removed. | []string |
| whitelistKeys | If the key is not matched, that field is removed. | []string |
