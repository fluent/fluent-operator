package filter

import (
	"testing"

	"github.com/fluent/fluent-operator/v2/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v2/apis/fluentbit/v1alpha2/plugins/params"
	. "github.com/onsi/gomega"
)

func TestFilter_RecordModifier_Params(t *testing.T) {
	g := NewGomegaWithT(t)

	sl := plugins.NewSecretLoader(nil, "test namespace")

	mod := RecordModifier{
		CommonParams: plugins.CommonParams{},
		Records: []string{
			"hostname ${HOSTNAME}",
			"product Awesome_Tool",
		},
		RemoveKeys: []string{
			"Swap.total",
			"Swap.free",
			"Swap.used",
		},
		AllowlistKeys: []string{
			"Mem.total",
			"Mem.free",
			"Mem.used",
		},
		WhitelistKeys: []string{
			"Disk.total",
			"Disk.free",
			"Disk.used",
		},
		UUIDKeys: []string{
			"ID1",
			"UiD2",
		},
	}

	expected := params.NewKVs()
	expected.Insert("Record", "hostname ${HOSTNAME}")
	expected.Insert("Record", "product Awesome_Tool")

	expected.Insert("Remove_key", "Swap.total")
	expected.Insert("Remove_key", "Swap.free")
	expected.Insert("Remove_key", "Swap.used")
	expected.Insert("Allowlist_key", "Mem.total")
	expected.Insert("Allowlist_key", "Mem.free")
	expected.Insert("Allowlist_key", "Mem.used")
	expected.Insert("Whitelist_key", "Disk.total")
	expected.Insert("Whitelist_key", "Disk.free")
	expected.Insert("Whitelist_key", "Disk.used")
	expected.Insert("Uuid_key", "ID1")
	expected.Insert("Uuid_key", "UiD2")

	kvs, err := mod.Params(sl)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(kvs).To(Equal(expected))
}
