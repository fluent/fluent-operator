package input

// The in_monitor_agent Input plugin exports Fluentd's internal metrics via REST API.
type MonitorAgent struct {
	// The port to listen to.
	Port *int64 `json:"port,omitempty"`
	// The bind address to listen to.
	Bind *string `json:"bind,omitempty"`
	// If you set this parameter, this plugin emits metrics as records.
	Tag *string `json:"tag,omitempty"`
	// The interval time between event emits. This will be used when "tag" is configured.
	EmitInterval *int64 `json:"emitInterval,omitempty"`
	// You can set this option to false to remove the config field from the response.
	IncludeConfig *bool `json:"includeConfig,omitempty"`
	// You can set this option to false to remove the retry field from the response.
	IncludeRetry *bool `json:"includeRetry,omitempty"`
}
