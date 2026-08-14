package v1alpha2

import "github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins/params"

// validateNoNewlines rejects any field value containing CR/LF. These section-header
// fields (Match, Alias, ...) are written directly to the config buffer, bypassing
// the KVs renderer, so they need the same newline check (GHSA-2j8x-46rv-qmpq).
func validateNoNewlines(fields map[string]string) error {
	for field, value := range fields {
		if err := params.ValidateNoNewline(field, value); err != nil {
			return err
		}
	}
	return nil
}
