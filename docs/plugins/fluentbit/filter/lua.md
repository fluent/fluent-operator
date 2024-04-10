# Lua

The Lua Filter allows you to modify the incoming records using custom Lua Scripts. <br /> **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/filters/lua**


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| script | Path to the Lua script that will be used. | v1.ConfigMapKeySelector |
| call | Lua function name that will be triggered to do filtering. It's assumed that the function is declared inside the Script defined above. | string |
| code | Inline LUA code instead of loading from a path via script. | string |
| typeIntKey | If these keys are matched, the fields are converted to integer. If more than one key, delimit by space. Note that starting from Fluent Bit v1.6 integer data types are preserved and not converted to double as in previous versions. | []string |
| protectedMode | If enabled, Lua script will be executed in protected mode. It prevents to crash when invalid Lua script is executed. Default is true. | *bool |
| timeAsTable | By default when the Lua script is invoked, the record timestamp is passed as a Floating number which might lead to loss precision when the data is converted back. If you desire timestamp precision enabling this option will pass the timestamp as a Lua table with keys sec for seconds since epoch and nsec for nanoseconds. | bool |
