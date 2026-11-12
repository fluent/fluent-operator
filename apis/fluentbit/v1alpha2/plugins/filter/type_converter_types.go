package filter

import (
	"fmt"

	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins/params"
)

// +kubebuilder:object:generate:=true
// +kubebuilder:validation:XValidation:rule="has(self.intKeys) || has(self.uintKeys) || has(self.floatKeys) || has(self.strKeys)",message="set at least one of intKeys, uintKeys, floatKeys and strKeys"

// The Type Converter Filter plugin converts the value of a key to another data type and appends it to the record under a new key. <br />
// At least one conversion is required. <br />
// **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/filters/type-converter**
type TypeConverter struct {
	plugins.CommonParams `json:",inline"`
	// Conversions of keys whose value is an integer.
	// +kubebuilder:validation:MinItems:=1
	IntKeys []TypeConversion `json:"intKeys,omitempty"`
	// Conversions of keys whose value is an unsigned integer.
	// +kubebuilder:validation:MinItems:=1
	UintKeys []TypeConversion `json:"uintKeys,omitempty"`
	// Conversions of keys whose value is a float.
	// +kubebuilder:validation:MinItems:=1
	FloatKeys []TypeConversion `json:"floatKeys,omitempty"`
	// Conversions of keys whose value is a string.
	// +kubebuilder:validation:MinItems:=1
	StrKeys []TypeConversion `json:"strKeys,omitempty"`
}

// +kubebuilder:object:generate:=true

// TypeConversion converts the value of FromKey to Type and appends it to the record as ToKey.
type TypeConversion struct {
	// Key whose value is converted.
	// +kubebuilder:validation:Pattern:="^\\S+$"
	FromKey string `json:"fromKey"`
	// Key the converted value is appended under.
	// +kubebuilder:validation:Pattern:="^\\S+$"
	ToKey string `json:"toKey"`
	// Data type the value is converted to.
	// +kubebuilder:validation:Enum:=int;uint;float;string
	Type string `json:"type"`
}

func (*TypeConverter) Name() string {
	return "type_converter"
}

func (t *TypeConverter) Params(_ plugins.SecretLoader) (*params.KVs, error) {
	kvs := params.NewKVs()
	err := t.AddCommonParams(kvs)
	if err != nil {
		return kvs, err
	}
	insert := func(key string, conversions []TypeConversion) {
		for _, c := range conversions {
			kvs.Insert(key, fmt.Sprintf("%s %s %s", c.FromKey, c.ToKey, c.Type))
		}
	}
	insert("Int_Key", t.IntKeys)
	insert("Uint_Key", t.UintKeys)
	insert("Float_Key", t.FloatKeys)
	insert("Str_Key", t.StrKeys)
	return kvs, nil
}
