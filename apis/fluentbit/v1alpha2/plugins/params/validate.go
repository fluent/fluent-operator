package params

import (
	"fmt"
	"strings"
)

// ValidateNoNewline rejects a value containing CR/LF. Fluent Bit's classic config
// format is line-based with no way to escape a newline, so a newline in a
// tenant-supplied value could inject an arbitrary [SECTION] into the shared config
// (GHSA-2j8x-46rv-qmpq). field is used only for the error message.
func ValidateNoNewline(field, value string) error {
	if strings.ContainsAny(value, "\r\n") {
		return &InvalidValueError{Field: field}
	}
	return nil
}

// InvalidValueError is returned for a value that is not allowed in the config, as opposed to a failure to read
// something the config needs. The user has to change the value for the config to render.
type InvalidValueError struct {
	Field string
}

func (e *InvalidValueError) Error() string {
	return fmt.Sprintf("invalid value for %q: must not contain carriage-return or line-feed characters (would allow Fluent Bit config injection)", e.Field)
}
