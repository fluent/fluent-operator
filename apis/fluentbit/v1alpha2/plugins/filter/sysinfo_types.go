package filter

import (
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins/params"
)

// +kubebuilder:object:generate:=true

// The Sysinfo Filter plugin appends system information, such as the Fluent Bit version or the hostname, to each record. <br />
// Only the keys that are set are appended. <br />
// **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/filters/sysinfo**
type Sysinfo struct {
	plugins.CommonParams `json:",inline"`
	// Key to append the Fluent Bit version under.
	FluentBitVersionKey string `json:"fluentBitVersionKey,omitempty"`
	// Key to append the operating system name under, e.g. linux.
	OSNameKey string `json:"osNameKey,omitempty"`
	// Key to append the hostname under. In Kubernetes this is the hostname of the Fluent Bit pod,
	// which is the pod name unless hostNetwork is enabled, in which case it is the node hostname.
	HostnameKey string `json:"hostnameKey,omitempty"`
	// Key to append the operating system version under. Only supported on Linux.
	OSVersionKey string `json:"osVersionKey,omitempty"`
	// Key to append the kernel version under. Only supported on Linux.
	KernelVersionKey string `json:"kernelVersionKey,omitempty"`
}

func (*Sysinfo) Name() string {
	return "sysinfo"
}

func (s *Sysinfo) Params(_ plugins.SecretLoader) (*params.KVs, error) {
	kvs := params.NewKVs()
	err := s.AddCommonParams(kvs)
	if err != nil {
		return kvs, err
	}
	plugins.InsertKVString(kvs, "Fluentbit_Version_Key", s.FluentBitVersionKey)
	plugins.InsertKVString(kvs, "Os_Name_Key", s.OSNameKey)
	plugins.InsertKVString(kvs, "Hostname_Key", s.HostnameKey)
	plugins.InsertKVString(kvs, "Os_Version_Key", s.OSVersionKey)
	plugins.InsertKVString(kvs, "Kernel_Version_Key", s.KernelVersionKey)
	return kvs, nil
}
