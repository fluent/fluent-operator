# MultilineParser

**For full documentation, refer to https://docs.fluentbit.io/manual/administration/configuring-fluent-bit/multiline-parsing**


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| type | Set the multiline mode, for now, we support the type regex. | string |
| parser | Name of a pre-defined parser that must be applied to the incoming content before applying the regex rule. If no parser is defined, it's assumed that's a raw text and not a structured message. | string |
| keyContent | For an incoming structured message, specify the key that contains the data that should be processed by the regular expression and possibly concatenated. | string |
| flushTimeout | Timeout in milliseconds to flush a non-terminated multiline buffer. Default is set to 5 seconds. | int |
| rules | Configure a rule to match a multiline pattern. The rule has a specific format described below. Multiple rules can be defined. | []Rule |
# Rule




| Field | Description | Scheme |
| ----- | ----------- | ------ |
| start |  | string |
| regex |  | string |
| next |  | string |
