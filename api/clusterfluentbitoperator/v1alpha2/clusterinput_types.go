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
	"kubesphere.io/fluentbit-operator/api/clusterfluentbitoperator/v1alpha2/plugins/clusterinput"

	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"reflect"
	"sort"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ClusterInputSpec defines the desired state of ClusterInput
type ClusterInputSpec struct {
	// A user friendly alias name for this clusterinput plugin.
	// Used in metrics for distinction of each configured clusterinput.
	Alias string `json:"alias,omitempty"`
	// Dummy defines Dummy Input configuration.
	Dummy *clusterinput.Dummy `json:"dummy,omitempty"`
	// Tail defines Tail Input configuration.
	Tail *clusterinput.Tail `json:"tail,omitempty"`
	// Systemd defines Systemd Input configuration.
	Systemd *clusterinput.Systemd `json:"systemd,omitempty"`
}

// ClusterInputStatus defines the observed state of ClusterInput
type ClusterInputStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:resource:scope=Cluster

// ClusterInput is the Schema for the clusterinputs API
type ClusterInput struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClusterInputSpec   `json:"spec,omitempty"`
	Status ClusterInputStatus `json:"status,omitempty"`
}

// +kubebuilder:object:generate:=false
// InputByName implements sort.Interface for []Input based on the Name field.
type ClusterInputByName []ClusterInput

func (a ClusterInputByName) Len() int           { return len(a) }
func (a ClusterInputByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ClusterInputByName) Less(i, j int) bool { return a[i].Name < a[j].Name }

//+kubebuilder:object:root=true

// ClusterInputList contains a list of ClusterInput
type ClusterInputList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClusterInput `json:"items"`
}

func (list ClusterInputList) Load(sl plugins.SecretLoader) (string, error) {
	var buf bytes.Buffer

	sort.Sort(ClusterInputByName(list.Items))

	for _, item := range list.Items {
		merge := func(p plugins.Plugin) error {
			if p == nil || reflect.ValueOf(p).IsNil() {
				return nil
			}

			buf.WriteString("[Input]\n")
			buf.WriteString(fmt.Sprintf("    Name    %s\n", p.Name()))
			if item.Spec.Alias != "" {
				buf.WriteString(fmt.Sprintf("    Alias    %s\n", item.Spec.Alias))
			}
			kvs, err := p.Params(sl)
			if err != nil {
				return err
			}
			buf.WriteString(kvs.String())
			return nil
		}

		for i := 0; i < reflect.ValueOf(item.Spec).NumField(); i++ {
			p, _ := reflect.ValueOf(item.Spec).Field(i).Interface().(plugins.Plugin)
			if err := merge(p); err != nil {
				return "", err
			}
		}
	}

	return buf.String(), nil
}

func init() {
	SchemeBuilder.Register(&ClusterInput{}, &ClusterInputList{})
}
