package filter

import (
	"testing"

	. "github.com/onsi/gomega"
	"kubesphere.io/fluentbit-operator/api/fluentbitoperator/v1alpha2/plugins"
	"kubesphere.io/fluentbit-operator/api/fluentbitoperator/v1alpha2/plugins/params"
)

func TestFilter_Modify_Params(t *testing.T) {
	g := NewGomegaWithT(t)

	sl := plugins.NewSecretLoader(nil, "test namespace", nil)

	mod := Modify{
		Conditions: []Condition{
			{
				KeyValueEquals: map[string]string{
					"kve1": "kvev1",
					"kve0": "kvev0",
					"kve2": "kvev2",
				},
			},
			{
				KeyDoesNotExist: map[string]string{
					"kdn1": "kdnv1",
					"kdn0": "kdnv0",
					"kdn2": "kdnv2",
				},
			},
		},
		Rules: []Rule{
			{
				Set: map[string]string{
					"sk0":      "skv0",
					"customer": "cus1",
					"app":      "foo",
				},
				Add: map[string]string{
					"add_k1": "k1v",
					"add_k2": "k2val",
					"add_k0": "k0value",
				},
				Rename: map[string]string{
					"rk1": "r1v",
					"rk0": "r0v",
					"rk3": "r3v",
					"rk2": "r2v",
				},
			},
		},
	}

	expected := params.NewKVs()
	expected.Insert("Condition", "Key_value_equals    kve0    kvev0")
	expected.Insert("Condition", "Key_value_equals    kve1    kvev1")
	expected.Insert("Condition", "Key_value_equals    kve2    kvev2")
	expected.Insert("Condition", "Key_does_not_exist    kdn0    kdnv0")
	expected.Insert("Condition", "Key_does_not_exist    kdn1    kdnv1")
	expected.Insert("Condition", "Key_does_not_exist    kdn2    kdnv2")
	expected.Insert("Set", "app    foo")
	expected.Insert("Set", "customer    cus1")
	expected.Insert("Set", "sk0    skv0")
	expected.Insert("Add", "add_k0    k0value")
	expected.Insert("Add", "add_k1    k1v")
	expected.Insert("Add", "add_k2    k2val")
	expected.Insert("Rename", "rk0    r0v")
	expected.Insert("Rename", "rk1    r1v")
	expected.Insert("Rename", "rk2    r2v")
	expected.Insert("Rename", "rk3    r3v")

	// we should not see any permutations in serialized kvs
	i := 0
	for i < 5 {
		kvs, err := mod.Params(sl)
		g.Expect(err).NotTo(HaveOccurred())
		g.Expect(kvs).To(Equal(expected))

		i++
	}
}
