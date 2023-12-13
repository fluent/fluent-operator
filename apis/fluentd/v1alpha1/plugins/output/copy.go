package output

// Copy defines the parameters for out_Copy plugin
type Copy struct {
	// CopyMode defines how to pass the events to <store> plugins.
	// +kubebuilder:validation:Enum:=no_copy;shallow;deep;marshal
	CopyMode *string `json:"copyMode"`
}
