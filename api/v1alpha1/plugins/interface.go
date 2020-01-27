package plugins

import (
	"bytes"
	"fmt"
)

// +kubebuilder:object:generate=false

// The Plugin interface defines methods for transferring input, filter
// and output plugins to textual section content.
type Plugin interface {
	Name() string
	Params(SecretLoader) (*KVs, error)
}

type KVs struct {
	keys   []string
	values []string
}

func NewKVs() *KVs {
	return &KVs{
		keys:   []string{},
		values: []string{},
	}
}

func (kvs *KVs) Insert(key, value string) {
	kvs.keys = append(kvs.keys, key)
	kvs.values = append(kvs.values, value)
}

func (kvs *KVs) Merge(tail *KVs) {
	kvs.keys = append(kvs.keys, tail.keys...)
	kvs.values = append(kvs.values, tail.values...)
}

func (kvs *KVs) String() string {
	if kvs == nil {
		return ""
	}

	var buf bytes.Buffer
	for i := 0; i < len(kvs.keys); i++ {
		buf.WriteString(fmt.Sprintf("    %s    %s\n", kvs.keys[i], kvs.values[i]))
	}
	return buf.String()
}
