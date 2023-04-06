# Parser

The Parser Filter plugin allows to parse field in event records. <br /> **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/filters/parser**


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| keyName | Specify field name in record to parse. | string |
| parser | Specify the parser name to interpret the field. Multiple Parser entries are allowed (split by comma). | string |
| preserveKey | Keep original Key_Name field in the parsed result. If false, the field will be removed. | *bool |
| reserveData | Keep all other original fields in the parsed result. If false, all other original fields will be removed. | *bool |
| unescapeKey | If the key is a escaped string (e.g: stringify JSON), unescape the string before to apply the parser. | *bool |
