package plugins

// FluentBit Output name
const FluentBitOutput = "fluentbit_output"

// FluentBitOutputDefaultValues for FluentBit Output plugin
var FluentBitOutputDefaultValues = map[string]string{
	"name":            "es",
	"match":           "kube.*",
	"host":            "elasticsearch-logging-data.kubesphere-logging-system.svc",
	"port":            "9200",
	"logstash_format": "On",
	"replace_dots":    "on",
	"retry_limit":     "False",
	"type":            "flb_type",
	"time_key":        "@timestamp",
	"logstash_prefix": "logstash",
}

// FluentBitOutputTemplate for FluentBit Output plugin
const FluentBitOutputTemplate = `
[OUTPUT]
    Name  {{ .name }}
    Match {{ .match }}
    Host  {{ .host }}
    Port  {{ .port }}
    Logstash_Format {{ .logstash_format }}
    Replace_Dots {{ .replace_dots }}
    Retry_Limit {{ .retry_limit }}
    Type  {{ .type }}
    Time_Key {{ .time_key }}
    Logstash_Prefix {{ .logstash_prefix }}
`
