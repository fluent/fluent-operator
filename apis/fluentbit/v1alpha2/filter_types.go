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
	"crypto/md5"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/fluent/fluent-operator/v2/apis/fluentbit/v1alpha2/plugins"
	"github.com/fluent/fluent-operator/v2/apis/fluentbit/v1alpha2/plugins/custom"
	"github.com/fluent/fluent-operator/v2/apis/fluentbit/v1alpha2/plugins/filter"
	"github.com/fluent/fluent-operator/v2/pkg/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// +kubebuilder:object:root=true
// +kubebuilder:resource:shortName=fbf
// +genclient

// Filter is the Schema for namespace level filter API
type Filter struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec FilterSpec `json:"spec,omitempty"`
}

// +kubebuilder:object:root=true

// FilterList contains a list of Filters
type FilterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Filter `json:"items"`
}

type NSFilterByName []Filter

func (a NSFilterByName) Len() int           { return len(a) }
func (a NSFilterByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a NSFilterByName) Less(i, j int) bool { return a[i].Name < a[j].Name }

func (list FilterList) Load(sl plugins.SecretLoader) (string, error) {
	var buf bytes.Buffer

	sort.Sort(NSFilterByName(list.Items))

	for _, item := range list.Items {
		merge := func(p plugins.Plugin) error {
			if p == nil || reflect.ValueOf(p).IsNil() {
				return nil
			}

			buf.WriteString("[Filter]\n")
			if p.Name() != "" {
				buf.WriteString(fmt.Sprintf("    Name    %s\n", p.Name()))
			}
			if item.Spec.Match != "" {
				buf.WriteString(fmt.Sprintf("    Match    %s\n", utils.GenerateNamespacedMatchExpr(item.Namespace, item.Spec.Match)))
			}
			if item.Spec.MatchRegex != "" {
				buf.WriteString(fmt.Sprintf("    Match_Regex    %s\n", utils.GenerateNamespacedMatchRegExpr(item.Namespace, item.Spec.MatchRegex)))
			}

			switch f := p.(type) {
			case *filter.Kubernetes:
				kubeTagPrefix := f.KubeTagPrefix
				if kubeTagPrefix == "" {
					kubeTagPrefix = "kube.var.log.containers."
				}
				f.KubeTagPrefix = fmt.Sprintf("%x.%s", md5.Sum([]byte(item.Namespace)), kubeTagPrefix)
				if f.RegexParser != "" {
					f.RegexParser = fmt.Sprintf("%s-%x", f.RegexParser, md5.Sum([]byte(item.Namespace)))
				}
			case *filter.Parser:
				parsers := strings.Split(f.Parser, ",")
				for i := range parsers {
					parsers[i] = strings.Trim(parsers[i], " ")
					parsers[i] = fmt.Sprintf("%s-%x", parsers[i], md5.Sum([]byte(item.Namespace)))
				}
				f.Parser = strings.Join(parsers, ",")
			case *custom.CustomPlugin:
				if f.Config != "" {
					f.Config = custom.MakeCustomConfigNamespaced(f.Config, item.Namespace)
				}
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
	SchemeBuilder.Register(&Filter{}, &FilterList{})
}
