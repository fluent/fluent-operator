package output

import (
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins/params"
)

// +kubebuilder:object:generate:=true

// The tcp output plugin allows to send records to a remote TCP server. <br />
// The payload can be formatted in different ways as required. <br />
// **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/outputs/tcp-and-tls**
type TCP struct {
	// Target host where Fluent-Bit or Fluentd are listening for Forward messages.
	Host string `json:"host,omitempty"`
	// TCP Port of the target service.
	// +kubebuilder:validation:Minimum:=1
	// +kubebuilder:validation:Maximum:=65535
	Port *int32 `json:"port,omitempty"`
	// Specify the data format to be printed. Supported formats are msgpack json, json_lines and json_stream.
	// +kubebuilder:validation:Enum:=msgpack;json;json_lines;json_stream
	Format string `json:"format,omitempty"`
	// TSpecify the name of the time key in the output record.
	// To disable the time key just set the value to false.
	JsonDateKey string `json:"jsonDateKey,omitempty"`
	// Specify the format of the date. Supported formats are double, epoch
	// and iso8601 (eg: 2018-05-30T09:39:52.000681Z)
	// +kubebuilder:validation:Enum:=double;epoch;iso8601
	JsonDateFormat string `json:"jsonDateFormat,omitempty"`
	*plugins.TLS   `json:"tls,omitempty"`
	// Include fluentbit networking options for this output-plugin
	*plugins.Networking `json:"networking,omitempty"`
}

func (*TCP) Name() string {
	return "tcp"
}

func (t *TCP) Params(sl plugins.SecretLoader) (*params.KVs, error) {
	kvs := params.NewKVs()

	if t.TLS != nil {
		tls, err := t.TLS.Params(sl)
		if err != nil {
			return nil, err
		}
		kvs.Merge(tls)
	}
	if t.Networking != nil {
		net, err := t.Networking.Params(sl)
		if err != nil {
			return nil, err
		}
		kvs.Merge(net)
	}
	if t.Networking != nil {
		net, err := t.Networking.Params(sl)
		if err != nil {
			return nil, err
		}
		kvs.Merge(net)
	}

	plugins.InsertKVString(kvs, "Host", t.Host)
	plugins.InsertKVString(kvs, "Format", t.Format)
	plugins.InsertKVString(kvs, "json_date_key", t.JsonDateKey)
	plugins.InsertKVString(kvs, "json_date_format", t.JsonDateFormat)

	plugins.InsertKVField(kvs, "Port", t.Port)

	return kvs, nil
}
