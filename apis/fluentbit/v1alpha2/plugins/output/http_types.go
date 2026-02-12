package output

import (
	"fmt"

	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins/params"
)

// +kubebuilder:object:generate:=true

// The http output plugin allows to flush your records into a HTTP endpoint. <br />
// For now the functionality is pretty basic and it issues a POST request
// with the data records in MessagePack (or JSON) format. <br />
// **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/outputs/http**
type HTTP struct {
	// IP address or hostname of the target HTTP Server
	Host string `json:"host,omitempty"`
	// Basic Auth Username
	HTTPUser *plugins.Secret `json:"httpUser,omitempty"`
	// Basic Auth Password. Requires HTTP_User to be set
	HTTPPasswd *plugins.Secret `json:"httpPassword,omitempty"`
	// TCP port of the target HTTP Server
	// +kubebuilder:validation:Minimum:=1
	// +kubebuilder:validation:Maximum:=65535
	Port *int32 `json:"port,omitempty"`
	// Specify an HTTP Proxy. The expected format of this value is http://host:port.
	// Note that https is not supported yet.
	Proxy string `json:"proxy,omitempty"`
	// Specify an optional HTTP URI for the target web server, e.g: /something
	Uri string `json:"uri,omitempty"`
	// Set payload compression mechanism. Option available is 'gzip'
	Compress string `json:"compress,omitempty"`
	// Specify the data format to be used in the HTTP request body, by default it uses msgpack.
	// Other supported formats are json, json_stream and json_lines and gelf.
	// +kubebuilder:validation:Enum:=msgpack;json;json_stream;json_lines;gelf
	Format string `json:"format,omitempty"`
	// Specify if duplicated headers are allowed.
	// If a duplicated header is found, the latest key/value set is preserved.
	AllowDuplicatedHeaders *bool `json:"allowDuplicatedHeaders,omitempty"`
	// Specify an optional HTTP header field for the original message tag.
	HeaderTag string `json:"headerTag,omitempty"`
	// Add a HTTP header key/value pair. Multiple headers can be set.
	Headers map[string]string `json:"headers,omitempty"`
	// Specify the name of the time key in the output record.
	// To disable the time key just set the value to false.
	JsonDateKey string `json:"jsonDateKey,omitempty"`
	// Specify the format of the date. Supported formats are double, epoch
	// and iso8601 (eg: 2018-05-30T09:39:52.000681Z)
	JsonDateFormat string `json:"jsonDateFormat,omitempty"`
	// Specify the key to use for timestamp in gelf format
	GelfTimestampKey string `json:"gelfTimestampKey,omitempty"`
	// Specify the key to use for the host in gelf format
	GelfHostKey string `json:"gelfHostKey,omitempty"`
	// Specify the key to use as the short message in gelf format
	GelfShortMessageKey string `json:"gelfShortMessageKey,omitempty"`
	// Specify the key to use for the full message in gelf format
	GelfFullMessageKey string `json:"gelfFullMessageKey,omitempty"`
	// Specify the key to use for the level in gelf format
	GelfLevelKey string `json:"gelfLevelKey,omitempty"`
	// HTTP output plugin supports TTL/SSL, for more details about the properties available
	// and general configuration, please refer to the TLS/SSL section.
	*plugins.TLS `json:"tls,omitempty"`
	// Include fluentbit networking options for this output-plugin
	*plugins.Networking `json:"networking,omitempty"`
	// Limit the maximum number of Chunks in the filesystem for the current output logical destination.
	TotalLimitSize string `json:"totalLimitSize,omitempty"`
}

// implement Name method
func (*HTTP) Name() string {
	return "http"
}

// implement Params method
func (h *HTTP) Params(sl plugins.SecretLoader) (*params.KVs, error) {
	const header = "header"

	kvs := params.NewKVs()

	if err := plugins.InsertKVSecret(kvs, "http_User", h.HTTPUser, sl); err != nil {
		return nil, err
	}
	if err := plugins.InsertKVSecret(kvs, "http_Passwd", h.HTTPPasswd, sl); err != nil {
		return nil, err
	}

	plugins.InsertKVString(kvs, "host", h.Host)
	plugins.InsertKVField(kvs, "port", h.Port)
	plugins.InsertKVString(kvs, "uri", h.Uri)
	plugins.InsertKVString(kvs, "format", h.Format)

	kvs.InsertStringMap(h.Headers, func(k, v string) (string, string) {
		return header, fmt.Sprintf(" %s    %s", k, v)
	})

	plugins.InsertKVString(kvs, "json_date_key", h.JsonDateKey)
	plugins.InsertKVString(kvs, "json_date_format", h.JsonDateFormat)
	plugins.InsertKVString(kvs, "Proxy", h.Proxy)
	plugins.InsertKVString(kvs, "compress", h.Compress)
	plugins.InsertKVString(kvs, "header_tag", h.HeaderTag)
	plugins.InsertKVString(kvs, "gelf_timestamp_key", h.GelfTimestampKey)
	plugins.InsertKVString(kvs, "gelf_host_key", h.GelfHostKey)
	plugins.InsertKVString(kvs, "gelf_short_message_key", h.GelfShortMessageKey)
	plugins.InsertKVString(kvs, "gelf_full_message_key", h.GelfFullMessageKey)
	plugins.InsertKVString(kvs, "gelf_level_key", h.GelfLevelKey)
	plugins.InsertKVField(kvs, "allow_duplicated_headers", h.AllowDuplicatedHeaders)

	if h.TLS != nil {
		tls, err := h.TLS.Params(sl)
		if err != nil {
			return nil, err
		}
		kvs.Merge(tls)
	}
	if h.Networking != nil {
		net, err := h.Networking.Params(sl)
		if err != nil {
			return nil, err
		}
		kvs.Merge(net)
	}

	plugins.InsertKVString(kvs, "storage.total_limit_size", h.TotalLimitSize)

	return kvs, nil
}
