package input

import (
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins/params"
)

// +kubebuilder:object:generate:=true

// The KubernetesEvents input plugin allows you to collect kubernetes cluster events from kube-api server
// **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/inputs/kubernetes-events*
type KubernetesEvents struct {
	// Tag name associated to all records coming from this plugin.
	Tag string `json:"tag,omitempty"`
	// Set a database file to keep track of recorded Kubernetes events
	DB string `json:"db,omitempty"`
	// Set a database sync method. values: extra, full, normal and off
	DBSync string `json:"dbSync,omitempty"`
	// Set the polling interval for each channel.
	IntervalSec *int32 `json:"intervalSec,omitempty"`
	// Set the polling interval for each channel (sub seconds: nanoseconds).
	IntervalNsec *int64 `json:"intervalNsec,omitempty"`
	// API Server end-point
	KubeURL string `json:"kubeURL,omitempty"`
	// CA certificate file
	KubeCAFile string `json:"kubeCAFile,omitempty"`
	// Absolute path to scan for certificate files
	KubeCAPath string `json:"kubeCAPath,omitempty"`
	// Token file
	KubeTokenFile string `json:"kubeTokenFile,omitempty"`
	// configurable 'time to live' for the K8s token. By default, it is set to 600 seconds.
	// After this time, the token is reloaded from Kube_Token_File or the Kube_Token_Command.
	KubeTokenTTL string `json:"kubeTokenTTL,omitempty"`
	// kubernetes limit parameter for events query, no limit applied when set to 0.
	KubeRequestLimit *int32 `json:"kubeRequestLimit,omitempty"`
	// Kubernetes retention time for events.
	KubeRetentionTime string `json:"kubeRetentionTime,omitempty"`
	// Kubernetes namespace to query events from. Gets events from all namespaces by default
	KubeNamespace string `json:"kubeNamespace,omitempty"`
	// Debug level between 0 (nothing) and 4 (every detail).
	TLSDebug *int32 `json:"tlsDebug,omitempty"`
	// When enabled, turns on certificate validation when connecting to the Kubernetes API server.
	TLSVerify *bool `json:"tlsVerify,omitempty"`
	// Set optional TLS virtual host.
	TLSVhost string `json:"tlsVhost,omitempty"`
}

func (*KubernetesEvents) Name() string {
	return "kubernetes_events"
}

// implement Section() method
func (k *KubernetesEvents) Params(_ plugins.SecretLoader) (*params.KVs, error) {
	kvs := params.NewKVs()

	plugins.InsertKVString(kvs, "Tag", k.Tag)
	plugins.InsertKVString(kvs, "DB", k.DB)
	plugins.InsertKVString(kvs, "DB.Sync", k.DBSync)
	plugins.InsertKVString(kvs, "Kube_URL", k.KubeURL)
	plugins.InsertKVString(kvs, "Kube_CA_File", k.KubeCAFile)
	plugins.InsertKVString(kvs, "Kube_CA_Path", k.KubeCAPath)
	plugins.InsertKVString(kvs, "Kube_Token_File", k.KubeTokenFile)
	plugins.InsertKVString(kvs, "Kube_Token_TTL", k.KubeTokenTTL)
	plugins.InsertKVString(kvs, "Kube_Retention_Time", k.KubeRetentionTime)
	plugins.InsertKVString(kvs, "Kube_Namespace", k.KubeNamespace)
	plugins.InsertKVString(kvs, "tls.Vhost", k.TLSVhost)

	plugins.InsertKVField(kvs, "Interval_Sec", k.IntervalSec)
	plugins.InsertKVField(kvs, "Interval_Nsec", k.IntervalNsec)
	plugins.InsertKVField(kvs, "Kube_Request_Limit", k.KubeRequestLimit)
	plugins.InsertKVField(kvs, "tls.Debug", k.TLSDebug)
	plugins.InsertKVField(kvs, "tls.Verify", k.TLSVerify)

	return kvs, nil
}
