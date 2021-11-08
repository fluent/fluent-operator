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
	"fmt"
	"reflect"
	"sort"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kubesphere.io/fluentbit-operator/apis/plugins"
	"kubesphere.io/fluentbit-operator/apis/plugins/filter"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ClusterFilterSpec defines the desired state of ClusterFilter
type ClusterFilterSpec struct {
	// A pattern to match against the tags of incoming records.
	// It's case-sensitive and support the star (*) character as a wildcard.
	Match string `json:"match,omitempty"`
	// A regular expression to match against the tags of incoming records.
	// Use this option if you want to use the full regex syntax.
	MatchRegex string `json:"matchRegex,omitempty"`
	// A set of clusterfilter plugins in order.
	FilterItems []FilterItem `json:"filters,omitempty"`
}

type FilterItem struct {
	// Grep defines Grep Filter configuration.
	Grep *filter.Grep `json:"grep,omitempty"`
	// RecordModifier defines Record Modifier Filter configuration.
	RecordModifier *filter.RecordModifier `json:"recordModifier,omitempty"`
	// Kubernetes defines Kubernetes Filter configuration.
	Kubernetes *filter.Kubernetes `json:"kubernetes,omitempty"`
	// Modify defines Modify Filter configuration.
	Modify *filter.Modify `json:"modify,omitempty"`
	// Nest defines Nest Filter configuration.
	Nest *filter.Nest `json:"nest,omitempty"`
	// Parser defines Parser Filter configuration.
	Parser *filter.Parser `json:"clusterparser,omitempty"`
	// Lua defines Lua Filter configuration.
	Lua *filter.Lua `json:"lua,omitempty"`
	// Throttle defines a Throttle configuration.
	Throttle *filter.Throttle `json:"throttle,omitempty"`
	// RewriteTag defines a RewriteTag configuration.
	RewriteTag *filter.RewriteTag `json:"rewriteTag,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:resource:scope=Cluster
//+genclient

// ClusterFilter is the Schema for the clusterfilters API
type ClusterFilter struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec ClusterFilterSpec `json:"spec,omitempty"`
}

//+kubebuilder:object:root=true

// ClusterFilterList contains a list of ClusterFilter
type ClusterFilterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClusterFilter `json:"items"`
}

// +kubebuilder:object:generate:=false

// FilterByName implements sort.Interface for []Filter based on the Name field.
type ClusterFilterByName []ClusterFilter

func (a ClusterFilterByName) Len() int           { return len(a) }
func (a ClusterFilterByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ClusterFilterByName) Less(i, j int) bool { return a[i].Name < a[j].Name }

func (list ClusterFilterList) Load(sl plugins.SecretLoader) (string, error) {
	var buf bytes.Buffer

	sort.Sort(ClusterFilterByName(list.Items))

	for _, item := range list.Items {
		merge := func(p plugins.Plugin) error {
			if p == nil || reflect.ValueOf(p).IsNil() {
				return nil
			}

			buf.WriteString("[Filter]\n")
			buf.WriteString(fmt.Sprintf("    Name    %s\n", p.Name()))
			if item.Spec.Match != "" {
				buf.WriteString(fmt.Sprintf("    Match    %s\n", item.Spec.Match))
			}
			if item.Spec.MatchRegex != "" {
				buf.WriteString(fmt.Sprintf("    Match_Regex    %s\n", item.Spec.MatchRegex))
			}
			kvs, err := p.Params(sl)
			if err != nil {
				return err
			}
			buf.WriteString(kvs.String())
			return nil
		}

		for _, elem := range item.Spec.FilterItems {
			for i := 0; i < reflect.ValueOf(elem).NumField(); i++ {
				p, _ := reflect.ValueOf(elem).Field(i).Interface().(plugins.Plugin)
				if err := merge(p); err != nil {
					return "", err
				}
			}
		}
	}

	return buf.String(), nil
}

func init() {
	SchemeBuilder.Register(&ClusterFilter{}, &ClusterFilterList{})
}
