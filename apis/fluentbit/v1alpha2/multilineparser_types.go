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
	"github.com/fluent/fluent-operator/v2/apis/fluentbit/v1alpha2/plugins"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"reflect"
	"sort"
)

// +kubebuilder:object:root=true
// +kubebuilder:resource:shortName=fbmp
// +genclient

// MultilineParser is the Schema of namespace-level multiline parser API
type MultilineParser struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec MultilineParserSpec `json:"spec,omitempty"`
}

// +kubebuilder:object:root=true

// MultilineParserList contains a list of MultilineParser
type MultilineParserList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MultilineParser `json:"items"`
}

// +kubebuilder:object:generate:=false

type MultilineParserByName []MultilineParser

func (a MultilineParserByName) Len() int {
	return len(a)
}

func (a MultilineParserByName) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a MultilineParserByName) Less(i, j int) bool {
	return a[i].Name < a[j].Name
}

func (list MultilineParserList) Load(sl plugins.SecretLoader) (string, error) {
	var buf bytes.Buffer

	sort.Sort(MultilineParserByName(list.Items))

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
	SchemeBuilder.Register(&MultilineParser{}, &MultilineParserList{})
}
