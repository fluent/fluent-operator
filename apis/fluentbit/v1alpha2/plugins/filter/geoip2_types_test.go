package filter

import (
	"testing"

	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins/params"
	. "github.com/onsi/gomega"
)

func TestFilter_GeoIP2_Params(t *testing.T) {
	g := NewGomegaWithT(t)

	sl := plugins.NewSecretLoader(nil, "test namespace")

	geoip2 := GeoIP2{
		Database:   "/fluent-bit/geoip/GeoLite2-City.mmdb",
		LookupKeys: []string{"remote_addr", "x_forwarded_for"},
		Records: []GeoIP2Record{
			{Key: "country", LookupKey: "remote_addr", Value: "%{country.names.en}"},
			{Key: "latitude", LookupKey: "remote_addr", Value: "%{location.latitude}"},
			{Key: "client_country", LookupKey: "x_forwarded_for", Value: "%{country.iso_code}"},
		},
	}

	expected := params.NewKVs()
	expected.Insert("Database", "/fluent-bit/geoip/GeoLite2-City.mmdb")
	expected.Insert("Lookup_Key", "remote_addr")
	expected.Insert("Lookup_Key", "x_forwarded_for")
	expected.Insert("Record", "country remote_addr %{country.names.en}")
	expected.Insert("Record", "latitude remote_addr %{location.latitude}")
	expected.Insert("Record", "client_country x_forwarded_for %{country.iso_code}")

	kvs, err := geoip2.Params(sl)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(kvs).To(Equal(expected))
}

func TestFilter_GeoIP2_ParamsCommonParams(t *testing.T) {
	g := NewGomegaWithT(t)

	sl := plugins.NewSecretLoader(nil, "test namespace")

	geoip2 := GeoIP2{
		CommonParams: plugins.CommonParams{Alias: "geoip"},
		Database:     "/geoip/GeoLite2-Country.mmdb",
		LookupKeys:   []string{"ip"},
		Records: []GeoIP2Record{
			{Key: "country", LookupKey: "ip", Value: "%{country.iso_code}"},
		},
	}

	expected := params.NewKVs()
	expected.Insert("Alias", "geoip")
	expected.Insert("Database", "/geoip/GeoLite2-Country.mmdb")
	expected.Insert("Lookup_Key", "ip")
	expected.Insert("Record", "country ip %{country.iso_code}")

	kvs, err := geoip2.Params(sl)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(kvs).To(Equal(expected))
}
