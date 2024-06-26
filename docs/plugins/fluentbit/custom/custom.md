# CustomPlugin

CustomPlugin is used to support filter plugins that are not implemented yet. <br /> **For example usage, refer to https://github.com/fluent/fluent-operator/blob/master/docs/best-practice/custom-plugin.md**


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| config | Config holds any unsupported plugins classic configurations, if ConfigFileFormat is set to yaml, this filed will be ignored | string |
| yamlConfig | YamlConfig holds the unsupported plugins yaml configurations, it only works when the ConfigFileFormat is yaml | *plugins.Config |
