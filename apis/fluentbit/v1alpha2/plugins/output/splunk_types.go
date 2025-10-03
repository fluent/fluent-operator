package output

import (
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins/params"
)

// +kubebuilder:object:generate:=true

// Splunk output plugin allows to ingest your records into a Splunk Enterprise service
// through the HTTP Event Collector (HEC) interface. <br />
// **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/outputs/splunk**
type Splunk struct {
	// IP address or hostname of the target OpenSearch instance, default `127.0.0.1`
	Host string `json:"host,omitempty"`
	// TCP port of the target Splunk instance, default `8088`
	// +kubebuilder:validation:Minimum:=1
	// +kubebuilder:validation:Maximum:=65535
	Port *int32 `json:"port,omitempty"`
	// Specify the Authentication Token for the HTTP Event Collector interface.
	SplunkToken *plugins.Secret `json:"splunkToken,omitempty"`
	// Buffer size used to receive Splunk HTTP responses: Default `2M`
	// +kubebuilder:validation:Pattern:="^\\d+(k|K|KB|kb|m|M|MB|mb|g|G|GB|gb)?$"
	HTTPBufferSize string `json:"httpBufferSize,omitempty"`
	// Set payload compression mechanism. The only available option is gzip.
	Compress string `json:"compress,omitempty"`
	// Specify X-Splunk-Request-Channel Header for the HTTP Event Collector interface.
	Channel string `json:"channel,omitempty"`
	// Optional username credential for access
	HTTPUser *plugins.Secret `json:"httpUser,omitempty"`
	// Password for user defined in HTTP_User
	HTTPPasswd *plugins.Secret `json:"httpPassword,omitempty"`
	// If the HTTP server response code is 400 (bad request) and this flag is enabled, it will print the full HTTP request
	// and response to the stdout interface. This feature is available for debugging purposes.
	HTTPDebugBadRequest *bool `json:"httpDebugBadRequest,omitempty"`
	// When enabled, the record keys and values are set in the top level of the map instead of under the event key. Refer to
	// the Sending Raw Events section from the docs more details to make this option work properly.
	SplunkSendRaw *bool `json:"splunkSendRaw,omitempty"`
	// Specify the key name that will be used to send a single value as part of the record.
	EventKey string `json:"eventKey,omitempty"`
	// Specify the key name that contains the host value. This option allows a record accessors pattern.
	EventHost string `json:"eventHost,omitempty"`
	// Set the source value to assign to the event data.
	EventSource string `json:"eventSource,omitempty"`
	// Set the sourcetype value to assign to the event data.
	EventSourcetype string `json:"eventSourcetype,omitempty"`
	// Set a record key that will populate 'sourcetype'. If the key is found, it will have precedence
	// over the value set in event_sourcetype.
	EventSourcetypeKey string `json:"eventSourcetypeKey,omitempty"`
	// The name of the index by which the event data is to be indexed.
	EventIndex string `json:"eventIndex,omitempty"`
	// Set a record key that will populate the index field. If the key is found, it will have precedence
	// over the value set in event_index.
	EventIndexKey string `json:"eventIndexKey,omitempty"`
	// Set event fields for the record. This option is an array and the format is "key_name
	// record_accessor_pattern".
	EventFields []string `json:"eventFields,omitempty"`

	// Enables dedicated thread(s) for this output. Default value `2` is set since version 1.8.13. For previous versions is 0.
	Workers      *int32 `json:"Workers,omitempty"`
	*plugins.TLS `json:"tls,omitempty"`
	// Include fluentbit networking options for this output-plugin
	*plugins.Networking `json:"networking,omitempty"`
}

// Name implement Section() method
func (*Splunk) Name() string {
	return "splunk"
}

// Params implement Section() method
func (o *Splunk) Params(sl plugins.SecretLoader) (*params.KVs, error) {
	kvs := params.NewKVs()

	if err := plugins.InsertKVSecret(kvs, "splunk_token", o.SplunkToken, sl); err != nil {
		return nil, err
	}
	if err := plugins.InsertKVSecret(kvs, "http_user", o.HTTPUser, sl); err != nil {
		return nil, err
	}
	if err := plugins.InsertKVSecret(kvs, "http_passwd", o.HTTPPasswd, sl); err != nil {
		return nil, err
	}
	if o.TLS != nil {
		tls, err := o.TLS.Params(sl)
		if err != nil {
			return nil, err
		}
		kvs.Merge(tls)
	}
	if o.Networking != nil {
		net, err := o.Networking.Params(sl)
		if err != nil {
			return nil, err
		}
		kvs.Merge(net)
	}

	plugins.InsertKVString(kvs, "host", o.Host)
	plugins.InsertKVString(kvs, "http_buffer_size", o.HTTPBufferSize)
	plugins.InsertKVString(kvs, "compress", o.Compress)
	plugins.InsertKVString(kvs, "channel", o.Channel)
	plugins.InsertKVString(kvs, "event_key", o.EventKey)
	plugins.InsertKVString(kvs, "event_host", o.EventHost)
	plugins.InsertKVString(kvs, "event_source", o.EventSource)
	plugins.InsertKVString(kvs, "event_sourcetype", o.EventSourcetype)
	plugins.InsertKVString(kvs, "event_sourcetype_key", o.EventSourcetypeKey)
	plugins.InsertKVString(kvs, "event_index", o.EventIndex)
	plugins.InsertKVString(kvs, "event_index_key", o.EventIndexKey)

	plugins.InsertKVField(kvs, "port", o.Port)
	plugins.InsertKVField(kvs, "http_debug_bad_request", o.HTTPDebugBadRequest)
	plugins.InsertKVField(kvs, "splunk_send_raw", o.SplunkSendRaw)
	plugins.InsertKVField(kvs, "workers", o.Workers)

	if len(o.EventFields) > 0 {
		for _, v := range o.EventFields {
			kvs.Insert("event_field", v)
		}
	}

	return kvs, nil
}
