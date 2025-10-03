package filter

import (
	"crypto/md5"
	"fmt"

	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins/params"
)

// +kubebuilder:object:generate:=true

// Kubernetes filter allows to enrich your log files with Kubernetes metadata. <br />
// **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/filters/kubernetes**
type Kubernetes struct {
	plugins.CommonParams `json:",inline"`
	// Set the buffer size for HTTP client when reading responses from Kubernetes API server.
	// +kubebuilder:validation:Pattern:="^\\d+(k|K|KB|kb|m|M|MB|mb|g|G|GB|gb)?$"
	BufferSize string `json:"bufferSize,omitempty"`
	// API Server end-point
	KubeURL string `json:"kubeURL,omitempty"`
	// CA certificate file
	KubeCAFile string `json:"kubeCAFile,omitempty"`
	// Absolute path to scan for certificate files
	KubeCAPath string `json:"kubeCAPath,omitempty"`
	// Token file
	KubeTokenFile string `json:"kubeTokenFile,omitempty"`
	// When the source records comes from Tail input plugin,
	// this option allows to specify what's the prefix used in Tail configuration.
	KubeTagPrefix string `json:"kubeTagPrefix,omitempty"`
	// When enabled, it checks if the log field content is a JSON string map,
	// if so, it append the map fields as part of the log structure.
	MergeLog *bool `json:"mergeLog,omitempty"`
	// When Merge_Log is enabled, the filter tries to assume the log field from the incoming message is a JSON string message
	// and make a structured representation of it at the same level of the log field in the map.
	// Now if Merge_Log_Key is set (a string name), all the new structured fields taken from the original log content are inserted under the new key.
	MergeLogKey string `json:"mergeLogKey,omitempty"`
	// When Merge_Log is enabled, trim (remove possible \n or \r) field values.
	MergeLogTrim *bool `json:"mergeLogTrim,omitempty"`
	// Optional parser name to specify how to parse the data contained in the log key. Recommended use is for developers or testing only.
	MergeParser string `json:"mergeParser,omitempty"`
	// When Keep_Log is disabled, the log field is removed
	// from the incoming message once it has been successfully merged
	// (Merge_Log must be enabled as well).
	KeepLog *bool `json:"keepLog,omitempty"`
	// Debug level between 0 (nothing) and 4 (every detail).
	TLSDebug *int32 `json:"tlsDebug,omitempty"`
	// When enabled, turns on certificate validation when connecting to the Kubernetes API server.
	TLSVerify *bool `json:"tlsVerify,omitempty"`
	// When enabled, the filter reads logs coming in Journald format.
	UseJournal *bool `json:"useJournal,omitempty"`
	// When enabled, metadata will be fetched from K8s when docker_id is changed.
	CacheUseDockerId *bool `json:"cacheUseDockerId,omitempty"`
	// Set an alternative Parser to process record Tag and extract pod_name, namespace_name, container_name and docker_id.
	// The parser must be registered in a parsers file (refer to parser filter-kube-test as an example).
	RegexParser string `json:"regexParser,omitempty"`
	// Allow Kubernetes Pods to suggest a pre-defined Parser
	// (read more about it in Kubernetes Annotations section)
	K8SLoggingParser *bool `json:"k8sLoggingParser,omitempty"`
	// Allow Kubernetes Pods to exclude their logs from the log processor
	// (read more about it in Kubernetes Annotations section).
	K8SLoggingExclude *bool `json:"k8sLoggingExclude,omitempty"`
	// Include Kubernetes resource labels in the extra metadata.
	Labels *bool `json:"labels,omitempty"`
	// Include Kubernetes resource annotations in the extra metadata.
	Annotations *bool `json:"annotations,omitempty"`
	// If set, Kubernetes meta-data can be cached/pre-loaded from files in JSON format in this directory,
	// named as namespace-pod.meta
	KubeMetaPreloadCacheDir string `json:"kubeMetaPreloadCacheDir,omitempty"`
	// If set, use dummy-meta data (for test/dev purposes)
	DummyMeta *bool `json:"dummyMeta,omitempty"`
	// DNS lookup retries N times until the network start working
	DNSRetries *int32 `json:"dnsRetries,omitempty"`
	// DNS lookup interval between network status checks
	DNSWaitTime *int32 `json:"dnsWaitTime,omitempty"`
	// This is an optional feature flag to get metadata information from kubelet
	// instead of calling Kube Server API to enhance the log.
	// This could mitigate the Kube API heavy traffic issue for large cluster.
	UseKubelet *bool `json:"useKubelet,omitempty"`
	// kubelet port using for HTTP request, this only works when useKubelet is set to On.
	KubeletPort *int32 `json:"kubeletPort,omitempty"`
	// kubelet host using for HTTP request, this only works when Use_Kubelet set to On.
	KubeletHost string `json:"kubeletHost,omitempty"`
	// configurable TTL for K8s cached metadata. By default, it is set to 0
	// which means TTL for cache entries is disabled and cache entries are evicted at random
	// when capacity is reached. In order to enable this option, you should set the number to a time interval.
	// For example, set this value to 60 or 60s and cache entries which have been created more than 60s will be evicted.
	KubeMetaCacheTTL string `json:"kubeMetaCacheTTL,omitempty"`
	// configurable 'time to live' for the K8s token. By default, it is set to 600 seconds.
	// After this time, the token is reloaded from Kube_Token_File or the Kube_Token_Command.
	KubeTokenTTL string `json:"kubeTokenTTL,omitempty"`
	// Command to get Kubernetes authorization token.
	// By default, it will be NULL and we will use token file to get token.
	KubeTokenCommand string `json:"kubeTokenCommand,omitempty"`
	// Configurable TTL for K8s cached namespace metadata.
	// By default, it is set to 900 which means a 15min TTL for namespace cache entries.
	// Setting this to 0 will mean entries are evicted at random once the cache is full.
	KubeMetaNamespaceCacheTTL *int32 `json:"kubeMetaNamespaceCacheTTL,omitempty"`
	// Include Kubernetes namespace resource labels in the extra metadata.
	NamespaceLabels *bool `json:"namespaceLabels,omitempty"`
	// Include Kubernetes namespace resource annotations in the extra metadata.
	NamespaceAnnotations *bool `json:"namespaceAnnotations,omitempty"`
	// Include Kubernetes namespace metadata only and no pod metadata.
	// If this is set, the values of Labels and Annotations are ignored.
	NamespaceMetadataOnly *bool `json:"namespaceMetadataOnly,omitempty"`
	// Include Kubernetes owner references in the extra metadata.
	OwnerReferences *bool `json:"ownerReferences,omitempty"`
	// If true, Kubernetes metadata (e.g., pod_name, container_name, namespace_name etc) will be extracted from the tag itself.
	UseTagForMeta *bool `json:"useTagForMeta,omitempty"`
}

func (*Kubernetes) Name() string {
	return "kubernetes"
}

func (k *Kubernetes) Params(_ plugins.SecretLoader) (*params.KVs, error) {
	kvs := params.NewKVs()
	err := k.AddCommonParams(kvs)
	if err != nil {
		return kvs, err
	}

	plugins.InsertKVString(kvs, "Buffer_Size", k.BufferSize)
	plugins.InsertKVString(kvs, "Kube_URL", k.KubeURL)
	plugins.InsertKVString(kvs, "Kube_CA_File", k.KubeCAFile)
	plugins.InsertKVString(kvs, "Kube_CA_Path", k.KubeCAPath)
	plugins.InsertKVString(kvs, "Kube_Token_File", k.KubeTokenFile)
	plugins.InsertKVString(kvs, "Kube_Tag_Prefix", k.KubeTagPrefix)
	plugins.InsertKVString(kvs, "Merge_Log_Key", k.MergeLogKey)
	plugins.InsertKVString(kvs, "Merge_Parser", k.MergeParser)
	plugins.InsertKVString(kvs, "Regex_Parser", k.RegexParser)
	plugins.InsertKVString(kvs, "Kube_meta_preload_cache_dir", k.KubeMetaPreloadCacheDir)
	plugins.InsertKVString(kvs, "Kubelet_Host", k.KubeletHost)
	plugins.InsertKVString(kvs, "Kube_Token_TTL", k.KubeTokenTTL)
	plugins.InsertKVString(kvs, "Kube_Token_Command", k.KubeTokenCommand)

	plugins.InsertKVField(kvs, "Merge_Log", k.MergeLog)
	plugins.InsertKVField(kvs, "Merge_Log_Trim", k.MergeLogTrim)
	plugins.InsertKVField(kvs, "Keep_Log", k.KeepLog)
	plugins.InsertKVField(kvs, "tls.debug", k.TLSDebug)
	plugins.InsertKVField(kvs, "tls.verify", k.TLSVerify)
	plugins.InsertKVField(kvs, "Use_Journal", k.UseJournal)
	plugins.InsertKVField(kvs, "Cache_Use_Docker_Id", k.CacheUseDockerId)
	plugins.InsertKVField(kvs, "K8S-Logging.Parser", k.K8SLoggingParser)
	plugins.InsertKVField(kvs, "K8S-Logging.Exclude", k.K8SLoggingExclude)
	plugins.InsertKVField(kvs, "Labels", k.Labels)
	plugins.InsertKVField(kvs, "Annotations", k.Annotations)
	plugins.InsertKVField(kvs, "Dummy_Meta", k.DummyMeta)
	plugins.InsertKVField(kvs, "DNS_Retries", k.DNSRetries)
	plugins.InsertKVField(kvs, "DNS_Wait_Time", k.DNSWaitTime)
	plugins.InsertKVField(kvs, "Use_Kubelet", k.UseKubelet)
	plugins.InsertKVField(kvs, "Kubelet_Port", k.KubeletPort)
	plugins.InsertKVString(kvs, "Kube_Meta_Cache_TTL", k.KubeMetaCacheTTL)
	plugins.InsertKVField(kvs, "Kube_Meta_Namespace_Cache_TTL", k.KubeMetaNamespaceCacheTTL)
	plugins.InsertKVField(kvs, "Namespace_Labels", k.NamespaceLabels)
	plugins.InsertKVField(kvs, "Namespace_Annotations", k.NamespaceAnnotations)
	plugins.InsertKVField(kvs, "Namespace_Metadata_Only", k.NamespaceMetadataOnly)
	plugins.InsertKVField(kvs, "Owner_References", k.OwnerReferences)
	plugins.InsertKVField(kvs, "Use_Tag_For_Meta", k.UseTagForMeta)

	return kvs, nil
}

func (k *Kubernetes) MakeNamespaced(ns string) {
	if k.KubeTagPrefix == "" {
		k.KubeTagPrefix = "kube.var.log.containers."
	}
	k.KubeTagPrefix = fmt.Sprintf("%x.%s", md5.Sum([]byte(ns)), k.KubeTagPrefix)
	if k.RegexParser != "" {
		k.RegexParser = fmt.Sprintf("%s-%x", k.RegexParser, md5.Sum([]byte(ns)))
	}
}
