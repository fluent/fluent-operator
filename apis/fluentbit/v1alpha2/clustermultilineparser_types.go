/*
Copyright 2023.

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
	"fmt"
	"reflect"
	"sort"

	"github.com/fluent/fluent-operator/v2/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v2/apis/fluentbit/v1alpha2/plugins/multilineparser"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:object:generate:=true

type MultilineParserSpec struct {
	*multilineparser.MultilineParser `json:",inline"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:shortName=cfbmp,scope=Cluster
// +genclient
// +genclient:nonNamespaced

// ClusterMultilineParser is the Schema for the cluster-level multiline parser API
type ClusterMultilineParser struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec MultilineParserSpec `json:"spec,omitempty"`
}

// +kubebuilder:object:root=true

// ClusterMultilineParserList contains a list of ClusterMultilineParser
type ClusterMultilineParserList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClusterMultilineParser `json:"items"`
}

// +kubebuilder:object:generate:=false

// ClusterMultilineParserByName implements sort.Interface for []ClusterParser based on the Name field.
type ClusterMultilineParserByName []ClusterMultilineParser

func (a ClusterMultilineParserByName) Len() int {
	return len(a)
}

func (a ClusterMultilineParserByName) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ClusterMultilineParserByName) Less(i, j int) bool {
	return a[i].Name < a[j].Name
}

func (list ClusterMultilineParserList) Load(sl plugins.SecretLoader) (string, error) {
	var buf bytes.Buffer

	sort.Sort(ClusterMultilineParserByName(list.Items))

	for _, item := range list.Items {
		merge := func(p plugins.Plugin) error {
			if p == nil || reflect.ValueOf(p).IsNil() {
				return nil
			}

			buf.WriteString("[MULTILINE_PARSER]\n")
			buf.WriteString(fmt.Sprintf("    Name    %s\n", item.Name))

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
	SchemeBuilder.Register(&ClusterMultilineParser{}, &ClusterMultilineParserList{})
}
