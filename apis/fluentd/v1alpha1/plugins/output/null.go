package output

// Null defines the parameters for out_null output plugin
type Null struct {
	// NeverFlush for testing to simulate the output plugin that never succeeds to flush.
	NeverFlush *bool `json:"neverFlush,omitempty"`
}
