package plugins

// FluentBit Filter name
const FluentBitFilter = "fluentbit_filter"

// FluentBitFilterDefaultValues for FluentBit Filter plugin
var FluentBitFilterDefaultValues = map[string]string{
	"name":            "kubernetes",
	"match":           "kube.*",
	"kube_url":        "https://kubernetes.default.svc:443",
	"kube_ca_file":    "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt",
	"kube_token_file": "/var/run/secrets/kubernetes.io/serviceaccount/token",
}

// FluentBitFilterTemplate for FluentBit Filter plugin
const FluentBitFilterTemplate = `
[FILTER]
    Name                {{ .name }}
    Match               {{ .match }}
    Kube_URL            {{ .kube_url }}
    Kube_CA_File        {{ .kube_ca_file }}
    Kube_Token_File     {{ .kube_token_file }}
`
