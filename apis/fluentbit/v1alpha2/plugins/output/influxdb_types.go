package output

import (
	"strings"

	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins/params"
)

// +kubebuilder:object:generate:=true

// The influxdb output plugin, allows to flush your records into a InfluxDB time series database. <br />
// **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/outputs/influxdb**
type InfluxDB struct {
	// IP address or hostname of the target InfluxDB service.
	Host string `json:"host"`
	// TCP port of the target InfluxDB service.
	//  +kubebuilder:validation:Maximum=65535
	//  +kubebuilder:validation:Minimum=1
	Port *int32 `json:"port,omitempty"`
	// InfluxDB database name where records will be inserted.
	Database string `json:"database,omitempty"`
	// InfluxDB bucket name where records will be inserted - if specified, database is ignored and v2 of API is used
	Bucket string `json:"bucket,omitempty"`
	// InfluxDB organization name where the bucket is (v2 only)
	Org string `json:"org,omitempty"`
	// The name of the tag whose value is incremented for the consecutive simultaneous events.
	SequenceTag string `json:"sequenceTag,omitempty"`
	// Optional username for HTTP Basic Authentication
	HTTPUser *plugins.Secret `json:"httpUser,omitempty"`
	// Password for user defined in HTTP_User
	HTTPPasswd *plugins.Secret `json:"httpPassword,omitempty"`
	// Authentication token used with InfluxDB v2 - if specified, both HTTPUser and HTTPPasswd are ignored
	HTTPToken *plugins.Secret `json:"httpToken,omitempty"`
	// List of keys that needs to be tagged
	TagKeys []string `json:"tagKeys,omitempty"`
	// Automatically tag keys where value is string.
	AutoTags *bool `json:"autoTags,omitempty"`
	// Dynamically tag keys which are in the string array at Tags_List_Key key.
	TagsListEnabled *bool `json:"tagsListEnabled,omitempty"`
	// Key of the string array optionally contained within each log record that contains tag keys for that record
	TagsListKey  string `json:"tagListKey,omitempty"`
	*plugins.TLS `json:"tls,omitempty"`
	// Include fluentbit networking options for this output-plugin
	*plugins.Networking `json:"networking,omitempty"`
}

// Name implement Section() method
func (*InfluxDB) Name() string {
	return "influxdb"
}

func (o *InfluxDB) Params(sl plugins.SecretLoader) (*params.KVs, error) {
	kvs := params.NewKVs()

	plugins.InsertKVString(kvs, "Host", o.Host)
	plugins.InsertKVField(kvs, "Port", o.Port)
	plugins.InsertKVString(kvs, "Database", o.Database)
	plugins.InsertKVString(kvs, "Bucket", o.Bucket)
	plugins.InsertKVString(kvs, "Org", o.Org)
	plugins.InsertKVString(kvs, "Sequence_Tag", o.SequenceTag)

	if len(o.TagKeys) > 0 {
		kvs.Insert("Tag_Keys", strings.Join(o.TagKeys, " "))
	}

	plugins.InsertKVField(kvs, "Auto_Tags", o.AutoTags)
	plugins.InsertKVField(kvs, "Tags_List_Enabled", o.TagsListEnabled)
	plugins.InsertKVString(kvs, "Tags_List_Key", o.TagsListKey)

	if err := plugins.InsertKVSecret(kvs, "HTTP_User", o.HTTPUser, sl); err != nil {
		return nil, err
	}
	if err := plugins.InsertKVSecret(kvs, "HTTP_Passwd", o.HTTPPasswd, sl); err != nil {
		return nil, err
	}
	if err := plugins.InsertKVSecret(kvs, "HTTP_Token", o.HTTPToken, sl); err != nil {
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

	return kvs, nil
}
