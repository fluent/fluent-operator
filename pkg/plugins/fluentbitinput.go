package plugins

// FluentBit Input name
const FluentBitInput = "fluentbit_input"

// FluentBitInputDefaultValues for FluentBit Input plugin
var FluentBitInputDefaultValues = map[string]string{
	"name":             "tail",
	"path":             "/var/log/containers/*.log",
	"parser":           "docker",
	"tag":              "kube.*",
	"refresh_interval": "5",
	"mem_buf_limit":    "5MB",
	"skip_long_lines":  "On",
	"db":               "/tail-db/tail-containers-state.db",
	"db_sync":          "Normal",
}

// FluentBitInputTemplate for FluentBit Input plugin
const FluentBitInputTemplate = `
[INPUT]
    Name             {{ .name }}
    Path             {{ .path }}
    Parser           {{ .parser }}
    Tag              {{ .tag }}
    Refresh_Interval {{ .refresh_interval }}
    Mem_Buf_Limit    {{ .mem_buf_limit }}
    Skip_Long_Lines  {{ .skip_long_lines }}
    DB               {{ .db }}
    DB.Sync          {{ .db_sync }}
`
