package filter

import (
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins/params"
)

// +kubebuilder:object:generate:=true

// Wasm Filter allows you to modify the incoming records using Wasm technology.
// **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/filters/wasm**
type Wasm struct {
	plugins.CommonParams `json:",inline"`
	// Path to the built Wasm program that will be used. This can be a relative path against the main configuration file.
	WasmPath string `json:"wasmPath,omitempty"`
	// Define event format to interact with Wasm programs: msgpack or json. Default: json
	EventFormat string `json:"eventFormat,omitempty"`
	// Wasm function name that will be triggered to do filtering. It's assumed that the function is built inside the Wasm program specified above.
	FunctionName string `json:"functionName,omitempty"`
	// Specify the whitelist of paths to be able to access paths from WASM programs.
	AccessiblePaths []string `json:"accessiblePaths,omitempty"`
	// Size of the heap size of Wasm execution. Review unit sizes for allowed values.
	// +kubebuilder:validation:Pattern:="^\\d+(k|K|KB|kb|m|M|MB|mb|g|G|GB|gb)?$"
	WasmHeapSize string `json:"wasmHeapSize,omitempty"`
	// Size of the stack size of Wasm execution. Review unit sizes for allowed values.
	// +kubebuilder:validation:Pattern:="^\\d+(k|K|KB|kb|m|M|MB|mb|g|G|GB|gb)?$"
	WasmStackSize string `json:"wasmStackSize,omitempty"`
}

// Name is the name of the filter plugin.
func (*Wasm) Name() string {
	return "wasm"
}

// Params represents the config options for the filter plugin.
func (w *Wasm) Params(_ plugins.SecretLoader) (*params.KVs, error) {
	kvs := params.NewKVs()
	err := w.AddCommonParams(kvs)
	if err != nil {
		return kvs, err
	}
	if w.WasmPath != "" {
		kvs.Insert("Wasm_Path", w.WasmPath)
	}
	if w.EventFormat != "" {
		kvs.Insert("Event_Format", w.EventFormat)
	}
	if w.FunctionName != "" {
		kvs.Insert("Function_Name", w.FunctionName)
	}
	for _, p := range w.AccessiblePaths {
		kvs.Insert("Accessible_Paths", p)
	}
	if w.WasmHeapSize != "" {
		kvs.Insert("Wasm_Heap_Size", w.WasmHeapSize)
	}
	if w.WasmStackSize != "" {
		kvs.Insert("Wasm_Stack_Size", w.WasmStackSize)
	}

	return kvs, nil
}
