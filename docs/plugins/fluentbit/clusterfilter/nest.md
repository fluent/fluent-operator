# Nest

The Nest Filter plugin allows you to operate on or with nested data. Its modes of operation are


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| operation | Select the operation nest or lift | string |
| wildcard | Nest records which field matches the wildcard | []string |
| nestUnder | Nest records matching the Wildcard under this key | string |
| nestedUnder | Lift records nested under the Nested_under key | string |
| addPrefix | Prefix affected keys with this string | string |
| removePrefix | Remove prefix from affected keys if it matches this string | string |
