# Modify

The Modify Filter plugin allows you to change records using rules and conditions. <br /> **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/filters/modify**


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| conditions | All conditions have to be true for the rules to be applied. | []Condition |
| rules | Rules are applied in the order they appear, with each rule operating on the result of the previous rule. | []Rule |
# Condition

The plugin supports the following conditions


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| keyExists | Is true if KEY exists | string |
| keyDoesNotExist | Is true if KEY does not exist | map[string]string |
| aKeyMatches | Is true if a key matches regex KEY | string |
| noKeyMatches | Is true if no key matches regex KEY | string |
| keyValueEquals | Is true if KEY exists and its value is VALUE | map[string]string |
| keyValueDoesNotEqual | Is true if KEY exists and its value is not VALUE | map[string]string |
| keyValueMatches | Is true if key KEY exists and its value matches VALUE | map[string]string |
| keyValueDoesNotMatch | Is true if key KEY exists and its value does not match VALUE | map[string]string |
| matchingKeysHaveMatchingValues | Is true if all keys matching KEY have values that match VALUE | map[string]string |
| matchingKeysDoNotHaveMatchingValues | Is true if all keys matching KEY have values that do not match VALUE | map[string]string |
# Rule

The plugin supports the following rules


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| set | Add a key/value pair with key KEY and value VALUE. If KEY already exists, this field is overwritten | map[string]string |
| add | Add a key/value pair with key KEY and value VALUE if KEY does not exist | map[string]string |
| remove | Remove a key/value pair with key KEY if it exists | string |
| removeWildcard | Remove all key/value pairs with key matching wildcard KEY | string |
| removeRegex | Remove all key/value pairs with key matching regexp KEY | string |
| rename | Rename a key/value pair with key KEY to RENAMED_KEY if KEY exists AND RENAMED_KEY does not exist | map[string]string |
| hardRename | Rename a key/value pair with key KEY to RENAMED_KEY if KEY exists. If RENAMED_KEY already exists, this field is overwritten | map[string]string |
| copy | Copy a key/value pair with key KEY to COPIED_KEY if KEY exists AND COPIED_KEY does not exist | map[string]string |
| hardCopy | Copy a key/value pair with key KEY to COPIED_KEY if KEY exists. If COPIED_KEY already exists, this field is overwritten | map[string]string |
