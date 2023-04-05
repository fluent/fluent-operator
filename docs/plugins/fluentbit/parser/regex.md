# Regex

The regex parser allows to define a custom Ruby Regular Expression that will use a named capture feature to define which content belongs to which key name. <br /> **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/parsers/regular-expression**


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| regex |  | string |
| timeKey | Time_Key | string |
| timeFormat | Time_Format, eg. %Y-%m-%dT%H:%M:%S %z | string |
| timeKeep | Time_Keep | *bool |
| timeOffset | Time_Offset, eg. +0200 | string |
| types |  | string |
