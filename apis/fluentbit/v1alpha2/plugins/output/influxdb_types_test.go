package output

import (
	"testing"

	"github.com/fluent/fluent-operator/v2/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v2/apis/fluentbit/v1alpha2/plugins/params"
	. "github.com/onsi/gomega"
)

func ptrAny[T any](obj T) *T {
	return &obj
}

func TestOutput_InfluxDB_Params(t *testing.T) {
	g := NewGomegaWithT(t)

	sl := plugins.NewSecretLoader(nil, "test namespace")

	dd := InfluxDB{
		Host:            "127.0.0.1",
		Port:            ptrAny(int32(8086)),
		Database:        "fluentbit",
		Bucket:          "buck",
		Org:             "orgnis",
		SequenceTag:     "_inc",
		TagKeys:         []string{"foo", "bar", "foo:bar"},
		AutoTags:        ptrAny(false),
		TagsListEnabled: ptrAny(true),
		TagsListKey:     "taglist_key",
	}

	expected := params.NewKVs()
	expected.Insert("Host", "127.0.0.1")
	expected.Insert("Port", "8086")
	expected.Insert("Database", "fluentbit")
	expected.Insert("Bucket", "buck")
	expected.Insert("Org", "orgnis")
	expected.Insert("Sequence_Tag", "_inc")
	expected.Insert("Tag_Keys", "foo bar foo:bar")
	expected.Insert("Auto_Tags", "false")
	expected.Insert("Tags_List_Enabled", "true")
	expected.Insert("Tags_List_Key", "taglist_key")

	kvs, err := dd.Params(sl)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(kvs).To(Equal(expected))

}
