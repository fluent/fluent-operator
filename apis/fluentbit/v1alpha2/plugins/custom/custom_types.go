package custom

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/fluent/fluent-operator/v2/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v2/apis/fluentbit/v1alpha2/plugins/params"
	"github.com/fluent/fluent-operator/v2/pkg/utils"
	"sigs.k8s.io/yaml"
)

// +kubebuilder:object:generate:=true

// CustomPlugin is used to support filter plugins that are not implemented yet. <br />
// **For example usage, refer to https://github.com/fluent/fluent-operator/blob/master/docs/best-practice/custom-plugin.md**
type CustomPlugin struct {
	// Config holds any unsupported plugins classic configurations,
	// if ConfigFileFormat is set to yaml, this filed will be ignored
	Config string `json:"config,omitempty"`
	// YamlConfig holds the unsupported plugins yaml configurations, it only works when the ConfigFileFormat is yaml
	// +kubebuilder:pruning:PreserveUnknownFields
	YamlConfig *plugins.Config `json:"yamlConfig,omitempty"`
}

func (c *CustomPlugin) Name() string {
	return ""
}

func (c *CustomPlugin) Params(_ plugins.SecretLoader) (*params.KVs, error) {
	kvs := params.NewKVs()
	if c.Config != "" {
		kvs.Content = indentation(c.Config)
	} else if c.YamlConfig != nil {
		yamlConfig, err := yaml.Marshal(c.YamlConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize yaml config: %s", err)
		}
		kvs.YamlContent = string(yamlConfig)
	}

	return kvs, nil
}

func (c *CustomPlugin) MakeNamespaced(ns string) {
	if c.Config != "" {
		c.Config = MakeCustomConfigNamespaced(c.Config, ns)
	}
	if c.YamlConfig != nil {
		if match, ok := c.YamlConfig.Data["match"]; ok {
			c.YamlConfig.Data["match"] = utils.GenerateNamespacedMatchExpr(ns, match.(string))
		}
		if matchRegex, ok := c.YamlConfig.Data["match_regex"]; ok {
			c.YamlConfig.Data["match_regex"] = utils.GenerateNamespacedMatchExpr(ns, matchRegex.(string))
		}
	}
}

func indentation(str string) string {
	splits := strings.Split(str, "\n")
	var buf bytes.Buffer
	for _, i := range splits {
		if i != "" {
			buf.WriteString(fmt.Sprintf("    %s\n", strings.TrimSpace(i)))
		}
	}
	return buf.String()
}

func MakeCustomConfigNamespaced(customConfig string, namespace string) string {
	var buf bytes.Buffer
	sections := strings.Split(customConfig, "\n")
	for _, section := range sections {
		section = strings.TrimSpace(section)
		idx := strings.LastIndex(section, " ")
		if strings.HasPrefix(section, "Match_Regex") {
			buf.WriteString(fmt.Sprintf("Match_Regex %s\n", utils.GenerateNamespacedMatchRegExpr(namespace, section[idx+1:])))
			continue
		}
		if strings.HasPrefix(section, "Match") {
			buf.WriteString(fmt.Sprintf("Match %s\n", utils.GenerateNamespacedMatchExpr(namespace, section[idx+1:])))
			continue
		}
		buf.WriteString(fmt.Sprintf("%s\n", section))
	}
	return buf.String()
}
