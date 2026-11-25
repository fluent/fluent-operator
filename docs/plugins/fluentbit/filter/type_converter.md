# TypeConverter

The Type Converter Filter plugin converts the value of a key to another data type and appends it to the record under a new key. <br /> At least one conversion is required. <br /> **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/filters/type-converter**


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| intKeys | Conversions of keys whose value is an integer. | [][TypeConversion](#typeconversion) |
| uintKeys | Conversions of keys whose value is an unsigned integer. | [][TypeConversion](#typeconversion) |
| floatKeys | Conversions of keys whose value is a float. | [][TypeConversion](#typeconversion) |
| strKeys | Conversions of keys whose value is a string. | [][TypeConversion](#typeconversion) |
# TypeConversion

TypeConversion converts the value of FromKey to Type and appends it to the record as ToKey.


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| fromKey | Key whose value is converted. | string |
| toKey | Key the converted value is appended under. | string |
| type | Data type the value is converted to. | string |
