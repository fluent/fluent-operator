package output

import (
	"fmt"

	"github.com/fluent/fluent-operator/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/apis/fluentbit/v1alpha2/plugins/params"
)

// +kubebuilder:object:generate:=true

// OpenSearch is the opensearch output plugin, allows to ingest your records into an OpenSearch database.
type Splunk struct {
	// IP address or hostname of the target OpenSearch instance, default `127.0.0.1`
	Host string `json:"host,omitempty"`
	// TCP port of the target OpenSearch instance, default `9200`
	// +kubebuilder:validation:Minimum:=1
	// +kubebuilder:validation:Maximum:=65535
	Port *int32 `json:"port,omitempty"`
	// *@TODO* OpenSearch accepts new data on HTTP query path "/_bulk".
	// But it is also possible to serve OpenSearch behind a reverse proxy on a subpath.
	// This option defines such path on the fluent-bit side.
	// It simply adds a path prefix in the indexing HTTP POST URI.
	SplunkToken string `json:"splunk_token,omitempty"`
	// Specify the buffer size used to read the response from the OpenSearch HTTP service.
	// This option is useful for debugging purposes where is required to read full responses,
	// note that response size grows depending of the number of records inserted.
	// To set an unlimited amount of memory set this value to False,
	// otherwise the value must be according to the Unit Size specification.
	// +kubebuilder:validation:Pattern:="^\\d+(k|K|KB|kb|m|M|MB|mb|g|G|GB|gb)?$"
	BufferSize string `json:"bufferSize,omitempty"`
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
	HTTPDebugBadRequest bool `json:http_debug_bad_request,omitempty`
	//When enabled, the record keys and values are set in the top level of the map instead of under the event key. Refer to
	// the Sending Raw Events section from the docs for more details to make this option work properly.
	SplunkSendRaw bool `json:splunk_send_raw,omitempty`
	//Specify the key name that will be used to send a single value as part of the record.
	EventKey string `json:event_key,omitempty`
	//Specify the key name that contains the host value. This option allows a record accessors pattern.
	EventHost string `json:event_host,omitempty`
	//Set the source value to assign to the event data.
	EventSource string `json:event_source,omitempty`
	//Set the sourcetype value to assign to the event data.
	EventSourcetype string `json:event_sourcetype,omitempty`
	// Set a record key that will populate 'sourcetype'. If the key is found, it will have precedence
	// over the value set in event_sourcetype.
	EventSourcetypeKey string `json:event_sourcetype_key,omitempty`
	// The name of the index by which the event data is to be indexed.
	EventIndex string `json:event_index,omitempty`
	// Set a record key that will populate the index field. If the key is found, it will have precedence
	// over the value set in event_index.
	EventIndexKey string `json:event_index_key,omitempty`
	//Set event fields for the record. This option can be set multiple times and the format is key_name
	// record_accessor_pattern.
	EventField string `json:event_field,omitempty`

	// Enables dedicated thread(s) for this output. Default value is set since version 1.8.13. For previous versions is 0.
	Workers *int32 `json:"Workers,omitempty"`
	//@todo check if below code is neccessary
	*plugins.TLS `json:"tls,omitempty"`
}

// Name implement Section() method
func (_ *Splunk) Name() string {
	return "splunk"
}

// Params implement Section() method
func (o *Splunk) Params(sl plugins.SecretLoader) (*params.KVs, error) {
	kvs := params.NewKVs()
	if o.Host != "" {
		kvs.Insert("Host", o.Host)
	}
	if o.Port != nil {
		kvs.Insert("Port", fmt.Sprint(*o.Port))
	}
	if o.SplunkToken != "" {
		kvs.Insert("Splunk_Token", o.SplunkToken)
	}
	if o.BufferSize != "" {
		kvs.Insert("Buffer_Size", o.BufferSize)
	}
	if o.HTTPUser != nil {
		u, err := sl.LoadSecret(*o.HTTPUser)
		if err != nil {
			return nil, err
		}
		kvs.Insert("HTTP_User", u)
	}
	if o.HTTPPasswd != nil {
		pwd, err := sl.LoadSecret(*o.HTTPPasswd)
		if err != nil {
			return nil, err
		}
		kvs.Insert("HTTP_Passwd", pwd)
	}
	if o.Compress != "" {
		kvs.Insert("Index", o.Compress)
	}
	if o.Channel != "" {
		kvs.Insert("Type", o.Channel)
	}
	if o.HTTPDebugBadRequest {
		kvs.Insert("HTTP_Debug_Bad_Request", "On")
	} else {
		kvs.Insert("HTTP_Debug_Bad_Request", "Off")
	}
	if o.SplunkSendRaw {
		kvs.Insert("Splunk_Send_Raw", "On")
	} else {
		kvs.Insert("Splunk_Send_Raw", "Off")
	}

	if o.EventKey != "" {
		kvs.Insert("Event_Key", o.EventKey)
	}
	if o.EventHost != "" {
		kvs.Insert("Event_Host", o.EventHost)
	}
	if o.EventSource != "" {
		kvs.Insert("Event_Source", o.EventSource)
	}

	if o.Workers != nil {
		kvs.Insert("Workers", fmt.Sprint(*o.Workers))
	}
	if o.TLS != nil {
		tls, err := o.TLS.Params(sl)
		if err != nil {
			return nil, err
		}
		kvs.Merge(tls)
	}
	return kvs, nil
}
