package plugins

import (
	"fmt"

	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins/params"
)

// InsertKVField is a generic helper to insert fields into KVs with type conversion
func InsertKVField[T any](kvs *params.KVs, key string, value *T) {
	if value != nil {
		converted := fmt.Sprint(*value)
		if converted != "" {
			kvs.Insert(key, converted)
		}
	}
}

// InsertKVString is a helper for string fields (replaces the old insertKVStore)
func InsertKVString(kvs *params.KVs, key string, value string) {
	if value != "" {
		kvs.Insert(key, value)
	}
}

// InsertKVSecret is a helper for secret fields with error handling
func InsertKVSecret(kvs *params.KVs, key string, secret *Secret, sl SecretLoader) error {
	if secret != nil {
		value, err := sl.LoadSecret(*secret)
		if err != nil {
			return err
		}
		kvs.Insert(key, value)
	}
	return nil
}
