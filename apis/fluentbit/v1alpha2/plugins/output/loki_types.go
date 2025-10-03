package output

import (
	"fmt"
	"sort"
	"strings"

	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins/params"
)

// +kubebuilder:object:generate:=true

// The loki output plugin, allows to ingest your records into a Loki service. <br />
// **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/outputs/loki**
type Loki struct {
	// Loki hostname or IP address.
	Host string `json:"host"`
	// Loki TCP port
	// +kubebuilder:validation:Minimum:=1
	// +kubebuilder:validation:Maximum:=65535
	Port *int32 `json:"port,omitempty"`
	// Specify a custom HTTP URI. It must start with forward slash.
	Uri string `json:"uri,omitempty"`
	// Set HTTP basic authentication user name.
	HTTPUser *plugins.Secret `json:"httpUser,omitempty"`
	// Password for user defined in HTTP_User
	// Set HTTP basic authentication password
	HTTPPasswd *plugins.Secret `json:"httpPassword,omitempty"`
	// Set bearer token authentication token value.
	// Can be used as alterntative to HTTP basic authentication
	BearerToken *plugins.Secret `json:"bearerToken,omitempty"`
	// Tenant ID used by default to push logs to Loki.
	// If omitted or empty it assumes Loki is running in single-tenant mode and no X-Scope-OrgID header is sent.
	TenantID *plugins.Secret `json:"tenantID,omitempty"`
	// Stream labels for API request. It can be multiple comma separated of strings specifying  key=value pairs.
	// In addition to fixed parameters, it also allows to add custom record keys (similar to label_keys property).
	Labels []string `json:"labels,omitempty"`
	// Optional list of record keys that will be placed as stream labels.
	// This configuration property is for records key only.
	LabelKeys []string `json:"labelKeys,omitempty"`
	// Specify the label map file path. The file defines how to extract labels from each record.
	LabelMapPath string `json:"labelMapPath,omitempty"`
	// Optional list of keys to remove.
	RemoveKeys []string `json:"removeKeys,omitempty"`
	// If set to true and after extracting labels only a single key remains, the log line sent to Loki will be the value of that key in line_format.
	// +kubebuilder:validation:Enum:=on;off
	DropSingleKey string `json:"dropSingleKey,omitempty"`
	// Format to use when flattening the record to a log line. Valid values are json or key_value.
	// If set to json,  the log line sent to Loki will be the Fluent Bit record dumped as JSON.
	// If set to key_value, the log line will be each item in the record concatenated together (separated by a single space) in the format.
	// +kubebuilder:validation:Enum:=json;key_value
	LineFormat string `json:"lineFormat,omitempty"`
	// If set to true, it will add all Kubernetes labels to the Stream labels.
	// +kubebuilder:validation:Enum:=on;off
	AutoKubernetesLabels string `json:"autoKubernetesLabels,omitempty"`
	// Specify the name of the key from the original record that contains the Tenant ID.
	// The value of the key is set as X-Scope-OrgID of HTTP header. It is useful to set Tenant ID dynamically.
	TenantIDKey string `json:"tenantIDKey,omitempty"`
	// Stream structured metadata for API request. It can be multiple comma separated key=value pairs.
	// This is used for high cardinality data that isn't suited for using labels.
	// Only supported in Loki 3.0+ with schema v13 and TSDB storage.
	StructuredMetadata map[string]string `json:"structuredMetadata,omitempty"`
	// Optional list of record keys that will be placed as structured metadata.
	// This allows using record accessor patterns (e.g. $kubernetes['pod_name']) to reference record keys.
	StructuredMetadataKeys []string `json:"structuredMetadataKeys,omitempty"`
	*plugins.TLS           `json:"tls,omitempty"`
	// Include fluentbit networking options for this output-plugin
	*plugins.Networking `json:"networking,omitempty"`
	// Limit the maximum number of Chunks in the filesystem for the current output logical destination.
	TotalLimitSize string `json:"totalLimitSize,omitempty"`
	// Enables dedicated thread(s) for this output. Default value is set since version 1.8.13. For previous versions is 0.
	Workers *int32 `json:"workers,omitempty"`
}

// implement Section() method
func (*Loki) Name() string {
	return "loki"
}

// implement Section() method
func (l *Loki) Params(sl plugins.SecretLoader) (*params.KVs, error) {
	kvs := params.NewKVs()

	if err := plugins.InsertKVSecret(kvs, "http_user", l.HTTPUser, sl); err != nil {
		return nil, err
	}
	if err := plugins.InsertKVSecret(kvs, "http_passwd", l.HTTPPasswd, sl); err != nil {
		return nil, err
	}
	if err := plugins.InsertKVSecret(kvs, "bearer_token", l.BearerToken, sl); err != nil {
		return nil, err
	}
	if err := plugins.InsertKVSecret(kvs, "tenant_id", l.TenantID, sl); err != nil {
		return nil, err
	}
	if l.TLS != nil {
		tls, err := l.TLS.Params(sl)
		if err != nil {
			return nil, err
		}
		kvs.Merge(tls)
	}
	if l.Networking != nil {
		net, err := l.Networking.Params(sl)
		if err != nil {
			return nil, err
		}
		kvs.Merge(net)
	}

	plugins.InsertKVString(kvs, "host", l.Host)
	plugins.InsertKVString(kvs, "uri", l.Uri)
	plugins.InsertKVString(kvs, "label_map_path", l.LabelMapPath)
	plugins.InsertKVString(kvs, "drop_single_key", l.DropSingleKey)
	plugins.InsertKVString(kvs, "line_format", l.LineFormat)
	plugins.InsertKVString(kvs, "auto_kubernetes_labels", l.AutoKubernetesLabels)
	plugins.InsertKVString(kvs, "tenant_id_key", l.TenantIDKey)
	plugins.InsertKVString(kvs, "storage.total_limit_size", l.TotalLimitSize)

	plugins.InsertKVField(kvs, "port", l.Port)
	plugins.InsertKVField(kvs, "workers", l.Workers)
	if len(l.Labels) > 0 {
		// Sort labels to ensure deterministic output
		sortedLabels := make([]string, len(l.Labels))
		copy(sortedLabels, l.Labels)

		// Sort labels alphabetically by the key part (before "=")
		sort.Slice(sortedLabels, func(i, j int) bool {
			iParts := strings.SplitN(sortedLabels[i], "=", 2)
			jParts := strings.SplitN(sortedLabels[j], "=", 2)

			// Special case: "environment" should come before "job"
			if iParts[0] == "environment" && jParts[0] == "job" {
				return true
			}
			if iParts[0] == "job" && jParts[0] == "environment" {
				return false
			}

			// Otherwise sort alphabetically
			return iParts[0] < jParts[0]
		})

		kvs.Insert("labels", strings.Join(sortedLabels, ","))
	}

	if len(l.LabelKeys) > 0 {
		kvs.Insert("label_keys", strings.Join(l.LabelKeys, ","))
	}
	if len(l.RemoveKeys) > 0 {
		kvs.Insert("remove_keys", strings.Join(l.RemoveKeys, ","))
	}

	if len(l.StructuredMetadata) > 0 {
		var metadataPairs []string
		for k, v := range l.StructuredMetadata {
			metadataPairs = append(metadataPairs, fmt.Sprintf("%s=%s", k, v))
		}
		if len(metadataPairs) > 0 {
			sort.Strings(metadataPairs)
			kvs.Insert("structured_metadata", strings.Join(metadataPairs, ","))
		}
	}

	if len(l.StructuredMetadataKeys) > 0 {
		kvs.Insert("structured_metadata_keys", strings.Join(l.StructuredMetadataKeys, ","))
	}

	return kvs, nil
}
