package plugins

// FluentBit Service name
const FluentBitService = "fluentbit_service"

// FluentBitServiceDefaultValues for FluentBit Service plugin
var FluentBitServiceDefaultValues = map[string]string{
	"flush":        "1",
	"daemon":       "Off",
	"log_level":    "info",
	"parsers_file": "parsers.conf",
}

// FluentBitServiceTemplate for FluentBit Service plugin
const FluentBitServiceTemplate = `
[SERVICE]
    Flush        {{ .flush }}
    Daemon       {{ .daemon }}
    Log_Level    {{ .log_level }}
    Parsers_File {{ .parsers_file }}
`
