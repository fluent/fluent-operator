package params

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/fluent/fluent-operator/v3/pkg/utils"
)

type kvTransformFunc func(string, string) (string, string)

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

func (kvs *KVs) String() (string, error) {
	if kvs == nil {
		return "", nil
	}

	// Content is a raw passthrough (CustomPlugin body), validated at its source and
	// legitimately multi-line, so it skips the per-value newline check below.
	if kvs.Content != "" {
		return kvs.Content, nil
	}

	var buf bytes.Buffer
	for i := 0; i < len(kvs.keys); i++ {
		if err := ValidateNoNewline(kvs.keys[i], kvs.values[i]); err != nil {
			return "", err
		}
		fmt.Fprintf(&buf, "    %s    %s\n", kvs.keys[i], kvs.values[i])
	}
	return buf.String(), nil
}

func (kvs *KVs) YamlString(depth int) (string, error) {
	if kvs == nil {
		return "", nil
	}
	// YamlContent is a raw passthrough (CustomPlugin body), validated at its source.
	if kvs.YamlContent != "" {
		return utils.AdjustYamlIndent(kvs.YamlContent, depth), nil
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
			if err := ValidateNoNewline(key, values[0]); err != nil {
				return "", err
			}
			fmt.Fprintf(&buf, "%s%s: %s\n", utils.YamlIndent(depth), strings.ToLower(key), values[0])
		} else {
			if _, ok := keyFinishedMap[key]; ok { // avoid output multiple times
				continue
			}
			fmt.Fprintf(&buf, "%s%s:\n", utils.YamlIndent(depth), strings.ToLower(key))
			for _, value := range values {
				if err := ValidateNoNewline(key, value); err != nil {
					return "", err
				}
				fmt.Fprintf(&buf, "%s  - %s\n", utils.YamlIndent(depth), value)
			}
			keyFinishedMap[key] = true
		}
	}
	return buf.String(), nil
}
