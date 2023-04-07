package output

import (
	"fmt"
	"strings"

	"github.com/fluent/fluent-operator/v2/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v2/apis/fluentbit/v1alpha2/plugins/params"
)

// +kubebuilder:object:generate:=true

// The influxdb output plugin, allows to flush your records into a InfluxDB time series database. <br />
// **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/outputs/influxdb**
type InfluxDB struct {
	// IP address or hostname of the target InfluxDB service.
	// +kubebuilder:validation:Format="hostname"
	// +kubebuilder:validation:Format="ipv4"
	// +kubebuilder:validation:Format="ipv6"
	Host string `json:"host"`
	// TCP port of the target InfluxDB service.
	//  +kubebuilder:validation:Maximum=65536
	//  +kubebuilder:validation:Minimum=0
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
}

// Name implement Section() method
func (_ *InfluxDB) Name() string {
	return "influxdb"
}

func (o *InfluxDB) Params(sl plugins.SecretLoader) (*params.KVs, error) {
	kvs := params.NewKVs()
	// InfluxDB Validation
	if o.HTTPToken != nil {

	}
	if o.Host != "" {
		kvs.Insert("Host", o.Host)
	}
	if o.Port != nil {
		kvs.Insert("Port", fmt.Sprint(*o.Port))
	}
	if o.Database != "" {
		kvs.Insert("Database", o.Database)
	}
	if o.Bucket != "" {
		kvs.Insert("Bucket", o.Bucket)
	}
	if o.Org != "" {
		kvs.Insert("Org", o.Org)
	}
	if o.SequenceTag != "" {
		kvs.Insert("Sequence_Tag", o.SequenceTag)
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
	if o.HTTPToken != nil {
		pwd, err := sl.LoadSecret(*o.HTTPToken)
		if err != nil {
			return nil, err
		}
		kvs.Insert("HTTP_Token", pwd)
	}
	if o.TagKeys != nil {
		kvs.Insert("Tag_Keys", strings.Join(o.TagKeys, " "))
	}
	if o.AutoTags != nil {
		kvs.Insert("Auto_Tags", fmt.Sprint(*o.AutoTags))
	}
	if o.TagsListEnabled != nil {
		kvs.Insert("Tags_List_Enabled", fmt.Sprint(*o.TagsListEnabled))
	}
	if o.TagsListKey != "" {
		kvs.Insert("Tags_List_Key", o.TagsListKey)
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
