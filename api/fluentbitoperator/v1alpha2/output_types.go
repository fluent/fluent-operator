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
	"kubesphere.io/fluentbit-operator/api/fluentbitoperator/v1alpha2/plugins"
	"kubesphere.io/fluentbit-operator/api/fluentbitoperator/v1alpha2/plugins/output"
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
	// A user friendly alias name for this output plugin.
	// Used in metrics for distinction of each configured output.
	Alias string `json:"alias,omitempty"`
	// RetryLimit represents configuration for the scheduler which can be set independently on each output section.
	// This option allows to disable retries or impose a limit to try N times and then discard the data after reaching that limit.
	RetryLimit string `json:"retry_limit,omitempty"`
	// Elasticsearch defines Elasticsearch Output configuration.
	Elasticsearch *output.Elasticsearch `json:"es,omitempty"`
	// File defines File Output configuration.
	File *output.File `json:"file,omitempty"`
	// Forward defines Forward Output configuration.
	Forward *output.Forward `json:"forward,omitempty"`
	// HTTP defines HTTP Output configuration.
	HTTP *output.HTTP `json:"http,omitempty"`
	// Kafka defines Kafka Output configuration.
	Kafka *output.Kafka `json:"kafka,omitempty"`
	// Null defines Null Output configuration.
	Null *output.Null `json:"null,omitempty"`
	// Stdout defines Stdout Output configuration.
	Stdout *output.Stdout `json:"stdout,omitempty"`
	// TCP defines TCP Output configuration.
	TCP *output.TCP `json:"tcp,omitempty"`
	// Loki defines Loki Output configuration.
	Loki *output.Loki `json:"loki,omitempty"`
	// Syslog defines Syslog Output configuration.
	Syslog *output.Syslog `json:"syslog,omitempty"`
	// DataDog defines DataDog Output configuration.
	DataDog *output.DataDog `json:"datadog,omitempty"`
	// Firehose defines Firehose Output configuration.
	Fireose *output.Firehose `json:"firehose,omitempty"`
}

// +kubebuilder:object:root=true
// +genclient

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

// +kubebuilder:object:generate:=false
// OutputByName implements sort.Interface for []Output based on the Name field.
type OutputByName []Output

func (a OutputByName) Len() int           { return len(a) }
func (a OutputByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a OutputByName) Less(i, j int) bool { return a[i].Name < a[j].Name }

func (list OutputList) Load(sl plugins.SecretLoader) (string, error) {
	var buf bytes.Buffer

	sort.Sort(OutputByName(list.Items))

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
	SchemeBuilder.Register(&Output{}, &OutputList{})
}
