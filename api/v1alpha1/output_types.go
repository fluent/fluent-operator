/*

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

package v1alpha1

import (
	"bytes"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kubesphere.io/fluentbit-operator/api/v1alpha1/plugins"
	"kubesphere.io/fluentbit-operator/api/v1alpha1/plugins/output"
	"reflect"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// OutputSpec defines the desired state of Output
type OutputSpec struct {
	// A pattern to match against the tags of incoming records.
	// It's case sensitive and support the star (*) character as a wildcard.
	Match string `json:"match,omitempty"`
	// A regular expression to match against the tags of incoming records.
	// Use this option if you want to use the full regex syntax.
	MatchRegex string `json:"matchRegex,omitempty"`
	// Elasticsearch defines Elasticsearch Output configuration.
	Elasticsearch *output.Elasticsearch `json:"es,omitempty"`
	// Forward defines Forward Output configuration.
	Forward *output.Forward `json:"forward,omitempty"`
	// Kafka defines Kafka Output configuration.
	Kafka *output.Kafka `json:"kafka,omitempty"`
	// Null defines Null Output configuration.
	Null *output.Null `json:"null,omitempty"`
	// Stdout defines Stdout Output configuration.
	Stdout *output.Stdout `json:"stdout,omitempty"`
}

// +kubebuilder:object:root=true

// Output is the Schema for the outputs API
type Output struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec OutputSpec `json:"spec,omitempty"`
}

// +kubebuilder:object:root=true

// OutputList contains a list of Output
type OutputList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Output `json:"items"`
}

func (list OutputList) Load(sl plugins.SecretLoader) (string, error) {
	var buf bytes.Buffer

	for _, item := range list.Items {
		merge := func(p plugins.Plugin) error {
			if reflect.ValueOf(p).IsNil() {
				return nil
			}

			buf.WriteString("[Output]\n")
			buf.WriteString(fmt.Sprintf("    Name    %s\n", p.Name()))
			if item.Spec.Match != "" {
				buf.WriteString(fmt.Sprintf("    Match    %s\n", item.Spec.Match))
			}
			if item.Spec.MatchRegex != "" {
				buf.WriteString(fmt.Sprintf("    Match_Regexp    %s\n", item.Spec.MatchRegex))
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
	SchemeBuilder.Register(&Output{}, &OutputList{})
}
