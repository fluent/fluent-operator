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
	"github.com/fluent/fluent-operator/apis/fluentbit/v1alpha2/plugins"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sort"
	"strings"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type CustomPluginSpec struct {
	// plugin name
	PluginName string `json:"pluginName,omitempty"`
	// plugin type
	PluginType string `json:"pluginType,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:shortName=ccp,scope=Cluster
// +genclient
// +genclient:nonNamespaced

// ClusterCustomPlugin defines a custom plugin
type ClusterCustomPlugin struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec CustomPluginSpec `json:"spec,omitempty"`
}

// +kubebuilder:object:root=true

// ClusterCustomPluginList contains a list of ClusterCustomPlugin
type ClusterCustomPluginList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClusterCustomPlugin `json:"items"`
}

// +kubebuilder:object:generate:=false

// CustomPluginByName implements sort.Interface for []ClusterCustomPlugin based on the Name field.
type CustomPluginByName []ClusterCustomPlugin

func (a CustomPluginByName) Len() int           { return len(a) }
func (a CustomPluginByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a CustomPluginByName) Less(i, j int) bool { return a[i].Name < a[j].Name }

func (list ClusterCustomPluginList) Load(_ plugins.SecretLoader) (string, error) {
	var buf bytes.Buffer
	sort.Sort(CustomPluginByName(list.Items))
	for _, item := range list.Items {
		pluginConf, ok := item.Annotations["plugin.config"]
		if !ok {
			return "", fmt.Errorf("plugin.conf annotation not found for %s", item.Name)
		}
		pluginType := item.Spec.PluginType
		pluginName := item.Spec.PluginName
		if pluginType != "" {
			buf.WriteString(fmt.Sprintf("[%s]\n", firstUpper(pluginType)))
		}
		if pluginName != "" {
			buf.WriteString(fmt.Sprintf("    Name    %s\n", pluginName))
		}
		lines := strings.Split(pluginConf, "\n")
		for _, line := range lines {
			if len(line) == 0 {
				continue
			}
			buf.WriteString(fmt.Sprintf("    %s\n", strings.TrimSpace(line)))
		}
	}
	return buf.String(), nil
}

func (list ClusterCustomPluginList) hasOutput() bool {
	var hasOutput bool
	for _, item := range list.Items {
		if firstUpper(item.Spec.PluginType) == "Output" {
			hasOutput = true
			break
		}
	}
	return hasOutput
}

func firstUpper(str string) string {
	if len(str) == 0 {
		return str
	}
	return strings.ToUpper(str[:1]) + str[1:]
}

func init() {
	SchemeBuilder.Register(&ClusterCustomPlugin{}, &ClusterCustomPluginList{})
}
