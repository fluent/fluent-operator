# Grep

The Grep Filter plugin allows to match or exclude specific records based in regular expression patterns. <br /> **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/filters/grep**


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| regex | Keep records which field matches the regular expression. Value Format: FIELD REGEX | string |
| exclude | Exclude records which field matches the regular expression. Value Format: FIELD REGEX | string |
