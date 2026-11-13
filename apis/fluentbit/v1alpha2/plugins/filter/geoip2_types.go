package filter

import (
	"fmt"

	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins/params"
)

// +kubebuilder:object:generate:=true

// The GeoIP2 Filter plugin looks up IP addresses of a record in a MaxMind GeoIP2 or GeoLite2 database <br />
// and appends the location information, such as the country or the city, to the record. <br />
// **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/filters/geoip2-filter**
type GeoIP2 struct {
	plugins.CommonParams `json:",inline"`
	// Path to the MaxMind database file (.mmdb), e.g. /fluent-bit/geoip/GeoLite2-City.mmdb.
	// The file can be mounted into the Fluent Bit pod with the volumes and volumeMounts of the FluentBit resource.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	Database string `json:"database"`
	// Keys of the record whose value is an IP address to look up.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinItems:=1
	// +kubebuilder:validation:items:Pattern:="^\\S+$"
	LookupKeys []string `json:"lookupKeys"`
	// Keys to append to the record, with the value looked up in the database.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinItems:=1
	Records []GeoIP2Record `json:"records"`
}

// +kubebuilder:object:generate:=true

// GeoIP2Record appends Key to the record, with the value found at Value in the database for the IP address in LookupKey.
type GeoIP2Record struct {
	// Key appended to the record.
	// +kubebuilder:validation:Pattern:="^\\S+$"
	Key string `json:"key"`
	// Lookup key holding the IP address to look up. It must be one of lookupKeys, otherwise the value is null.
	// +kubebuilder:validation:Pattern:="^\\S+$"
	LookupKey string `json:"lookupKey"`
	// Path of the value in the database, written as %{path}, e.g. %{country.names.en} or %{location.latitude}.
	// The path is at most 63 characters long.
	// +kubebuilder:validation:Pattern:="^%\\{[^\\s}]{1,63}\\}$"
	Value string `json:"value"`
}

func (*GeoIP2) Name() string {
	return "geoip2"
}

func (g *GeoIP2) Params(_ plugins.SecretLoader) (*params.KVs, error) {
	kvs := params.NewKVs()
	err := g.AddCommonParams(kvs)
	if err != nil {
		return kvs, err
	}
	plugins.InsertKVString(kvs, "Database", g.Database)
	for _, key := range g.LookupKeys {
		kvs.Insert("Lookup_Key", key)
	}
	for _, r := range g.Records {
		kvs.Insert("Record", fmt.Sprintf("%s %s %s", r.Key, r.LookupKey, r.Value))
	}
	return kvs, nil
}
