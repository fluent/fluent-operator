/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha2

import (
	"bytes"
	"kubesphere.io/fluentbit-operator/api/clusterfluentbitoperator/v1alpha2/plugins"
	"kubesphere.io/fluentbit-operator/api/clusterfluentbitoperator/v1alpha2/plugins/clusteroutput"

	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"reflect"
	"sort"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ClusterOutputSpec defines the desired state of ClusterOutput
type ClusterOutputSpec struct {
	// A pattern to match against the tags of incoming records.
	// It's case sensitive and support the star (*) character as a wildcard.
	Match string `json:"match,omitempty"`
	// A regular expression to match against the tags of incoming records.
	// Use this option if you want to use the full regex syntax.
	MatchRegex string `json:"matchRegex,omitempty"`
	// A user friendly alias name for this clusteroutput plugin.
	// Used in metrics for distinction of each configured clusteroutput.
	Alias string `json:"alias,omitempty"`
	// RetryLimit represents configuration for the scheduler which can be set independently on each clusteroutput section.
	// This option allows to disable retries or impose a limit to try N times and then discard the data after reaching that limit.
	RetryLimit string `json:"retry_limit,omitempty"`
	// Elasticsearch defines Elasticsearch Output configuration.
	Elasticsearch *clusteroutput.Elasticsearch `json:"es,omitempty"`
	// File defines File Output configuration.
	File *clusteroutput.File `json:"file,omitempty"`
	// Forward defines Forward Output configuration.
	Forward *clusteroutput.Forward `json:"forward,omitempty"`
	// HTTP defines HTTP Output configuration.
	HTTP *clusteroutput.HTTP `json:"http,omitempty"`
	// Kafka defines Kafka Output configuration.
	Kafka *clusteroutput.Kafka `json:"kafka,omitempty"`
	// Null defines Null Output configuration.
	Null *clusteroutput.Null `json:"null,omitempty"`
	// Stdout defines Stdout Output configuration.
	Stdout *clusteroutput.Stdout `json:"stdout,omitempty"`
	// TCP defines TCP Output configuration.
	TCP *clusteroutput.TCP `json:"tcp,omitempty"`
	// Loki defines Loki Output configuration.
	Loki *clusteroutput.Loki `json:"loki,omitempty"`
	// Syslog defines Syslog Output configuration.
	Syslog *clusteroutput.Syslog `json:"syslog,omitempty"`
	// DataDog defines DataDog Output configuration.
	DataDog *clusteroutput.DataDog `json:"datadog,omitempty"`
}

// ClusterOutputStatus defines the observed state of ClusterOutput
type ClusterOutputStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:resource:scope=Cluster

// ClusterOutput is the Schema for the clusteroutputs API
type ClusterOutput struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClusterOutputSpec   `json:"spec,omitempty"`
	Status ClusterOutputStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ClusterOutputList contains a list of ClusterOutput
type ClusterOutputList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClusterOutput `json:"items"`
}

// +kubebuilder:object:generate:=false
// OutputByName implements sort.Interface for []Output based on the Name field.
type ClusterOutputByName []ClusterOutput

func (a ClusterOutputByName) Len() int           { return len(a) }
func (a ClusterOutputByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ClusterOutputByName) Less(i, j int) bool { return a[i].Name < a[j].Name }

func (list ClusterOutputList) Load(sl plugins.SecretLoader) (string, error) {
	var buf bytes.Buffer

	sort.Sort(ClusterOutputByName(list.Items))

	for _, item := range list.Items {
		merge := func(p plugins.Plugin) error {
			if p == nil || reflect.ValueOf(p).IsNil() {
				return nil
			}

			buf.WriteString("[Output]\n")
			buf.WriteString(fmt.Sprintf("    Name    %s\n", p.Name()))
			if item.Spec.Match != "" {
				buf.WriteString(fmt.Sprintf("    Match    %s\n", item.Spec.Match))
			}
			if item.Spec.MatchRegex != "" {
				buf.WriteString(fmt.Sprintf("    Match_Regex    %s\n", item.Spec.MatchRegex))
			}
			if item.Spec.Alias != "" {
				buf.WriteString(fmt.Sprintf("    Alias    %s\n", item.Spec.Alias))
			}
			if item.Spec.RetryLimit != "" {
				buf.WriteString(fmt.Sprintf("    Retry_Limit    %s\n", item.Spec.RetryLimit))
			}
			kvs, err := p.Params(sl)
			if err != nil {
				return err
			}
			buf.WriteString(kvs.String())
			return nil
		}

		for i := 2; i < reflect.ValueOf(item.Spec).NumField(); i++ {
			p, _ := reflect.ValueOf(item.Spec).Field(i).Interface().(plugins.Plugin)
			if err := merge(p); err != nil {
				return "", err
			}
		}
	}

	return buf.String(), nil
}


func init() {
	SchemeBuilder.Register(&ClusterOutput{}, &ClusterOutputList{})
}
