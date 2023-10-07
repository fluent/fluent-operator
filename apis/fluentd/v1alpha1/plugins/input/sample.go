package input

// The in_sample input plugin generates sample events. It is useful for testing, debugging, benchmarking and getting started with Fluentd.
type Sample struct {
	// The tag of the event. The value is the tag assigned to the generated events.
	Tag *string `json:"tag,omitempty"`
	// The number of events in the event stream of each emit.
	Size *int64 `json:"size,omitempty"`
	// It configures how many events to generate per second.
	Rate *int64 `json:"rate,omitempty"`
	// If specified, each generated event has an auto-incremented key field.
	AutoIncrementKey *string `json:"autoIncrementKey,omitempty"`
	// The sample data to be generated. It should be either an array of JSON hashes or a single JSON hash. If it is an array of JSON hashes, the hashes in the array are cycled through in order.
	Sample *string `json:"sample,omitempty"`
}
