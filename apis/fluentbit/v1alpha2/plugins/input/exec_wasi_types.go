package input

import (
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins/params"
)

// +kubebuilder:object:generate:=true

// The exec_wasi input plugin, allows to execute WASM program that is WASI target like as external program and collects event logs from there.
// **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/inputs/exec-wasi**
type ExecWasi struct {
	// The place of a WASM program file.
	WASIPath string `json:"wasiPath,omitempty"`
	// Specify the name of a parser to interpret the entry as a structured message.
	Parser string `json:"parser,omitempty"`
	// Specify the whitelist of paths to be able to access paths from WASM programs.
	AccessiblePaths []string `json:"accessiblePaths,omitempty"`
	// Polling interval (seconds).
	IntervalSec *int32 `json:"intervalSec,omitempty"`
	// Polling interval (nanoseconds).
	IntervalNSec *int64 `json:"intervalNSec,omitempty"`
	// +kubebuilder:validation:Pattern:="^\\d+(k|K|KB|kb|m|M|MB|mb|g|G|GB|gb)?$"
	WasmHeapSize string `json:"wasmHeapSize,omitempty"`
	// Size of the stack size of Wasm execution. Review unit sizes for allowed values.
	// +kubebuilder:validation:Pattern:="^\\d+(k|K|KB|kb|m|M|MB|mb|g|G|GB|gb)?$"
	WasmStackSize string `json:"wasmStackSize,omitempty"`
	// Size of the buffer (check unit sizes for allowed values)
	// +kubebuilder:validation:Pattern:="^\\d+(k|K|KB|kb|m|M|MB|mb|g|G|GB|gb)?$"
	BufSize string `json:"bufSize,omitempty"`
	// Indicates whether to run this input in its own thread. Default: false.
	Threaded *bool `json:"threaded,omitempty"`
}

func (*ExecWasi) Name() string {
	return "exec_wasi"
}

// Params implement Section() method
func (w *ExecWasi) Params(_ plugins.SecretLoader) (*params.KVs, error) {
	kvs := params.NewKVs()

	plugins.InsertKVString(kvs, "WASI_Path", w.WASIPath)
	plugins.InsertKVString(kvs, "Parser", w.Parser)
	plugins.InsertKVString(kvs, "Wasm_Heap_Size", w.WasmHeapSize)
	plugins.InsertKVString(kvs, "Wasm_Stack_Size", w.WasmStackSize)
	plugins.InsertKVString(kvs, "Buf_Size", w.BufSize)

	plugins.InsertKVField(kvs, "Interval_Sec", w.IntervalSec)
	plugins.InsertKVField(kvs, "Interval_NSec", w.IntervalNSec)
	plugins.InsertKVField(kvs, "Threaded", w.Threaded)

	for _, p := range w.AccessiblePaths {
		kvs.Insert("Accessible_Paths", p)
	}

	return kvs, nil
}
