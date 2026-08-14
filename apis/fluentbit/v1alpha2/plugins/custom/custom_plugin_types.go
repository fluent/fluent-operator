package custom

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins/params"
	"github.com/fluent/fluent-operator/v3/pkg/utils"
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
		if err := validateClassicConfig(c.Config); err != nil {
			return nil, err
		}
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

// validateClassicConfig keeps a customPlugin config confined to the section the
// operator renders it into. The body is indented but the classic parser ignores
// leading whitespace, so a line opening a new [SECTION] or a Fluent Bit command
// (@INCLUDE/@SET) could inject an arbitrary plugin into the shared config
// (GHSA-2j8x-46rv-qmpq). A legitimate body holds only key-value properties, so
// reject any line starting with '[' or '@'.
func validateClassicConfig(config string) error {
	for i, raw := range strings.Split(config, "\n") {
		line := strings.TrimSpace(raw)
		if strings.HasPrefix(line, "[") {
			return fmt.Errorf("invalid customPlugin config: line %d %q opens a new config section; customPlugin config may only contain the key-value properties of the enclosing section", i+1, line)
		}
		if strings.HasPrefix(line, "@") {
			return fmt.Errorf("invalid customPlugin config: line %d %q is a Fluent Bit command, which is not permitted in customPlugin config", i+1, line)
		}
	}
	return nil
}

func indentation(str string) string {
	var buf bytes.Buffer
	for s := range strings.SplitSeq(str, "\n") {
		if s != "" {
			fmt.Fprintf(&buf, "    %s\n", strings.TrimSpace(s))
		}
	}
	return buf.String()
}

func MakeCustomConfigNamespaced(customConfig string, namespace string) string {
	var buf bytes.Buffer
	for section := range strings.SplitSeq(customConfig, "\n") {
		section = strings.TrimSpace(section)
		idx := strings.LastIndex(section, " ")
		if strings.HasPrefix(section, "Match_Regex") {
			fmt.Fprintf(&buf, "Match_Regex %s\n", utils.GenerateNamespacedMatchRegExpr(namespace, section[idx+1:]))
			continue
		}
		if strings.HasPrefix(section, "Match") {
			fmt.Fprintf(&buf, "Match %s\n", utils.GenerateNamespacedMatchExpr(namespace, section[idx+1:]))
			continue
		}
		fmt.Fprintf(&buf, "%s\n", section)
	}
	return buf.String()
}
