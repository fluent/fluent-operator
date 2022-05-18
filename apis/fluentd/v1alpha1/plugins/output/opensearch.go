package output

import "github.com/fluent/fluent-operator/apis/fluentd/v1alpha1/plugins"

// OpenSearch defines the parameters for out_opensearch_out plugin
type Opensearch struct {
	// The hostname of your OpenSearch node (default: localhost).
	Host *string `json:"host,omitempty"`
	// The port number of your OpenSearch node (default: 9200).
	// +kubebuilder:validation:Minimum:=1
	// +kubebuilder:validation:Maximum:=65535
	Port *uint32 `json:"port,omitempty"`
	// Hosts defines a list of hosts if you want to connect to more than one OpenSearch nodes
	Hosts *string `json:"hosts,omitempty"`
	// Specify https if your OpenSearch endpoint supports SSL (default: http).
	Scheme *string `json:"scheme,omitempty"`
	// Path defines the REST API endpoint of OpenSearch to post write requests (default: nil).
	Path *string `json:"path,omitempty"`
	// IndexName defines the placeholder syntax of Fluentd plugin API. See https://docs.fluentd.org/configuration/buffer-section.
	IndexName *string `json:"indexName,omitempty"`
	// If true, Fluentd uses the conventional index name format logstash-%Y.%m.%d (default: false). This option supersedes the index_name option.
	LogstashFormat *bool `json:"logstashFormat,omitempty"`
	// LogstashPrefix defines the logstash prefix index name to write events when logstash_format is true (default: logstash).
	LogstashPrefix *string `json:"logstashPrefix,omitempty"`
	// Optional, The login credentials to connect to OpenSearch
	User *plugins.Secret `json:"user,omitempty"`
	// Optional, The login credentials to connect to OpenSearch
	Password *plugins.Secret `json:"password,omitempty"`
}
