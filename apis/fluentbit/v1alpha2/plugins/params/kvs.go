package params

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/fluent/fluent-operator/v2/pkg/utils"
)

type kvTransformFunc func(string, string) (string, string)

type SecretProvider interface {
	GetStringVal() string
}

type kvTransformFunc1 func(string, string) (string, SecretProvider)

type KVs struct {
	keys        []string
	values      []string
	Content     string
	YamlContent string
}

func NewKVs() *KVs {
	return &KVs{
		keys:   []string{},
		values: []string{},
	}
}

func (kvs *KVs) InsertMapValMap(m map[string]SecretProvider, f kvTransformFunc1) {
	if len(m) > 0 {
		keys := make([]string, 0, len(m))
		for k := range m {
			keys = append(keys, k)
		}

		sort.Strings(keys)

		for _, k := range keys {
			v := m[k]
			strval := v.GetStringVal()
			if f != nil {
				transformedKey, transformedVal := f(k, strval)

				if transformedVal != nil {
					strval = transformedVal.GetStringVal()
				} else {
					strval = "" // Default to an empty string if transformation returns nil
				}

				k = transformedKey
			}

			kvs.Insert(k, strval)

		}
	}
}

func (kvs *KVs) Insert(key, value string) {
	kvs.keys = append(kvs.keys, key)
	kvs.values = append(kvs.values, value)
}

func (kvs *KVs) InsertStringMap(m map[string]string, f kvTransformFunc) {
	if len(m) > 0 {
		keys := make([]string, 0, len(m))
		for k := range m {
			keys = append(keys, k)
		}

		sort.Strings(keys)

		for _, k := range keys {
			v := m[k]
			if f != nil {
				k, v = f(k, v)
			}
			kvs.Insert(k, v)
		}
	}
}

func (kvs *KVs) Merge(tail *KVs) {
	kvs.keys = append(kvs.keys, tail.keys...)
	kvs.values = append(kvs.values, tail.values...)
}

func (kvs *KVs) String() string {
	if kvs == nil {
		return ""
	}

	if kvs.Content != "" {
		return kvs.Content
	}

	var buf bytes.Buffer
	for i := 0; i < len(kvs.keys); i++ {
		buf.WriteString(fmt.Sprintf("    %s    %s\n", kvs.keys[i], kvs.values[i]))
	}
	return buf.String()
}

func (kvs *KVs) YamlString(depth int) string {
	if kvs == nil {
		return ""
	}
	if kvs.YamlContent != "" {
		return utils.AdjustYamlIndent(kvs.YamlContent, depth)
	}

	var buf bytes.Buffer

	// deduplicate to yaml format
	keyValuesMap := make(map[string][]string)
	for i, k := range kvs.keys {
		keyValuesMap[k] = append(keyValuesMap[k], kvs.values[i])
	}
	keyFinishedMap := make(map[string]bool)
	for _, key := range kvs.keys { // keep the order
		values := keyValuesMap[key]
		if len(values) == 1 {
			buf.WriteString(fmt.Sprintf("%s%s: %s\n", utils.YamlIndent(depth), strings.ToLower(key), values[0]))
		} else {
			if _, ok := keyFinishedMap[key]; ok { // avoid output multiple times
				continue
			}
			buf.WriteString(fmt.Sprintf("%s%s:\n", utils.YamlIndent(depth), strings.ToLower(key)))
			for _, value := range values {
				buf.WriteString(fmt.Sprintf("%s  - %s\n", utils.YamlIndent(depth), value))
			}
			keyFinishedMap[key] = true
		}
	}
	return buf.String()
}
