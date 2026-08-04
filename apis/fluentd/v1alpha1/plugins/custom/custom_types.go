package custom

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/fluent/fluent-operator/v3/apis/fluentd/v1alpha1/plugins"
	"github.com/fluent/fluent-operator/v3/apis/fluentd/v1alpha1/plugins/params"
)

type CustomPlugin struct {
	Config string `json:"config"`
}

func (c *CustomPlugin) Name() string {
	return ""
}

func (c *CustomPlugin) Params(_ plugins.SecretLoader) (*params.PluginStore, error) {
	if err := validateConfig(c.Config); err != nil {
		return nil, err
	}
	ps := params.NewPluginStore("")
	ps.Content = indentation(c.Config)
	return ps, nil
}

// validateConfig ensures a customPlugin config snippet cannot break out of the
// routing block it is rendered into. The text is written verbatim into the
// shared Fluentd configuration, so a snippet whose directive blocks are not
// self-contained could close its enclosing <label> and inject top-level
// directives (e.g. an `@type exec` source) that execute in the shared
// aggregator process (GHSA-9jg5-g9mc-7m7c). Merely indenting the text does not
// prevent this because the Fluentd config grammar is whitespace-insensitive.
//
// We require the snippet's directive nesting to be self-contained: the depth may
// never go negative (a closing tag with no matching opener would close the
// surrounding block) and must return to zero at the end (a dangling opener would
// otherwise swallow the sibling config that follows it).
//
// Only a tag that begins a line changes block nesting in the Fluentd config
// grammar; angle brackets appearing inside parameter values do not, so those are
// ignored to avoid rejecting legitimate configs.
func validateConfig(config string) error {
	depth := 0
	for i, raw := range strings.Split(config, "\n") {
		line := strings.TrimSpace(raw)
		// Only a line-leading tag changes block nesting; blank lines, comments
		// and parameter lines never begin with '<'.
		if !strings.HasPrefix(line, "<") {
			continue
		}
		if strings.HasPrefix(line, "</") {
			depth--
			if depth < 0 {
				return fmt.Errorf("invalid customPlugin config: closing directive %q on line %d has no matching opening directive and would close the enclosing block", line, i+1)
			}
		} else {
			depth++
		}
	}
	if depth != 0 {
		return fmt.Errorf("invalid customPlugin config: %d directive block(s) left unclosed; every opened block must be closed within the snippet", depth)
	}
	return nil
}

func indentation(config string) string {
	var buf bytes.Buffer
	for split := range strings.SplitSeq(config, "\n") {
		if split != "" {
			fmt.Fprintf(&buf, "  %s\n", split)
		}
	}
	return buf.String()
}
