# Logfmt

The logfmt parser allows to parse the logfmt format described in https://brandur.org/logfmt . <br /> **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/parsers/logfmt**


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| timeKey | Time_Key | string |
| timeFormat | Time_Format, eg. %Y-%m-%dT%H:%M:%S %z | string |
| timeKeep | Time_Keep | *bool |
