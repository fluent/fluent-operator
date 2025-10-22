package output

import (
	"strings"

	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v3/apis/fluentbit/v1alpha2/plugins/params"
)

// +kubebuilder:object:generate:=true

// Stackdriver is the Stackdriver output plugin, allows you to ingest your records into GCP Stackdriver. <br />
// **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/outputs/stackdriver**
type Stackdriver struct {
	// Path to GCP Credentials JSON file
	GoogleServiceCredentials string `json:"googleServiceCredentials,omitempty"`
	// Email associated with the service
	ServiceAccountEmail *plugins.Secret `json:"serviceAccountEmail,omitempty"`
	// Private Key associated with the service
	ServiceAccountSecret *plugins.Secret `json:"serviceAccountSecret,omitempty"`
	// Metadata Server Prefix
	MetadataServer string `json:"metadataServer,omitempty"`
	// GCP/AWS region to store data. Required if Resource is generic_node or generic_task
	Location string `json:"location,omitempty"`
	// Namespace identifier. Required if Resource is generic_node or generic_task
	Namespace string `json:"namespace,omitempty"`
	// Node identifier within the namespace. Required if Resource is generic_node or generic_task
	NodeID string `json:"nodeID,omitempty"`
	// Identifier for a grouping of tasks. Required if Resource is generic_task
	Job string `json:"job,omitempty"`
	// Identifier for a task within a namespace. Required if Resource is generic_task
	TaskID string `json:"taskID,omitempty"`
	// The GCP Project that should receive the logs
	ExportToProjectID string `json:"exportToProjectID,omitempty"`
	// Set resource types of data
	Resource string `json:"resource,omitempty"`
	// Name of the cluster that the pod is running in. Required if Resource is k8s_container, k8s_node, or k8s_pod
	K8sClusterName string `json:"k8sClusterName,omitempty"`
	// Location of the cluster that contains the pods/nodes. Required if Resource is k8s_container, k8s_node, or k8s_pod
	K8sClusterLocation string `json:"k8sClusterLocation,omitempty"`
	// Used by Stackdriver to find related labels and extract them to LogEntry Labels
	LabelsKey string `json:"labelsKey,omitempty"`
	// Optional list of comma separated of strings for key/value pairs
	Labels []string `json:"labels,omitempty"`
	// The value of this field is set as the logName field in Stackdriver
	LogNameKey string `json:"logNameKey,omitempty"`
	// Used to validate the tags of logs that when the Resource is k8s_container, k8s_node, or k8s_pod
	TagPrefix string `json:"tagPrefix,omitempty"`
	// Specify the key that contains the severity information for the logs
	SeverityKey string `json:"severityKey,omitempty"`
	// Rewrite the trace field to be formatted for use with GCP Cloud Trace
	AutoformatStackdriverTrace *bool `json:"autoformatStackdriverTrace,omitempty"`
	// Number of dedicated threads for the Stackdriver Output Plugin
	Workers *int32 `json:"workers,omitempty"`
	// A custom regex to extract fields from the local_resource_id of the logs
	CustomK8sRegex string `json:"customK8sRegex,omitempty"`
	// Optional list of comma separated strings. Setting these fields overrides the Stackdriver monitored resource API values
	ResourceLabels []string `json:"resourceLabels,omitempty"`
	// the key to used to select the text payload from the record
	TextPayloadKey string `json:"textPayloadKey,omitempty"`
}

// Name implement Section() method
func (*Stackdriver) Name() string {
	return "stackdriver"
}

// Params implement Section() method
func (o *Stackdriver) Params(sl plugins.SecretLoader) (*params.KVs, error) {
	kvs := params.NewKVs()

	if err := plugins.InsertKVSecret(kvs, "service_account_email", o.ServiceAccountEmail, sl); err != nil {
		return nil, err
	}
	if err := plugins.InsertKVSecret(kvs, "service_account_secret", o.ServiceAccountSecret, sl); err != nil {
		return nil, err
	}

	plugins.InsertKVString(kvs, "google_service_credentials", o.GoogleServiceCredentials)
	plugins.InsertKVString(kvs, "metadata_server", o.MetadataServer)
	plugins.InsertKVString(kvs, "location", o.Location)
	plugins.InsertKVString(kvs, "namespace", o.Namespace)
	plugins.InsertKVString(kvs, "node_id", o.NodeID)
	plugins.InsertKVString(kvs, "job", o.Job)
	plugins.InsertKVString(kvs, "task_id", o.TaskID)
	plugins.InsertKVString(kvs, "export_to_project_id", o.ExportToProjectID)
	plugins.InsertKVString(kvs, "resource", o.Resource)
	plugins.InsertKVString(kvs, "k8s_cluster_name", o.K8sClusterName)
	plugins.InsertKVString(kvs, "k8s_cluster_location", o.K8sClusterLocation)
	plugins.InsertKVString(kvs, "labels_key", o.LabelsKey)
	plugins.InsertKVString(kvs, "log_name_key", o.LogNameKey)
	plugins.InsertKVString(kvs, "tag_prefix", o.TagPrefix)
	plugins.InsertKVString(kvs, "severity_key", o.SeverityKey)
	plugins.InsertKVString(kvs, "custom_k8s_regex", o.CustomK8sRegex)

	plugins.InsertKVField(kvs, "autoformat_stackdriver_trace", o.AutoformatStackdriverTrace)
	plugins.InsertKVField(kvs, "workers", o.Workers)

	if len(o.Labels) > 0 {
		kvs.Insert("labels", strings.Join(o.Labels, ","))
	}
	if len(o.ResourceLabels) > 0 {
		kvs.Insert("resource_labels", strings.Join(o.ResourceLabels, ","))
	}
	plugins.InsertKVString(kvs, "text_payload_key", o.TextPayloadKey)

	return kvs, nil
}
