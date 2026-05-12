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

// Package v1alpha2 contains API Schema definitions for the logging v1alpha2 API group
// +kubebuilder:object:generate=true
// +groupName=fluentbit.fluent.io
package v1alpha2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	// SchemeGroupVersion is group version used to register these objects
	SchemeGroupVersion = schema.GroupVersion{Group: "fluentbit.fluent.io", Version: "v1alpha2"}

	// SchemeBuilder is used to add go types to the GroupVersionKind scheme
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)

	// AddToScheme adds the types in this group-version to the given scheme.
	AddToScheme = SchemeBuilder.AddToScheme
)

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&ClusterFluentBitConfig{},
		&ClusterFluentBitConfigList{},
		&ClusterFilter{},
		&ClusterFilterList{},
		&ClusterInput{},
		&ClusterInputList{},
		&ClusterMultilineParser{},
		&ClusterMultilineParserList{},
		&ClusterOutput{},
		&ClusterOutputList{},
		&ClusterParser{},
		&ClusterParserList{},
		&Collector{},
		&CollectorList{},
		&Filter{},
		&FilterList{},
		&FluentBit{},
		&FluentBitList{},
		&FluentBitConfig{},
		&FluentBitConfigList{},
		&MultilineParser{},
		&MultilineParserList{},
		&Output{},
		&OutputList{},
		&Parser{},
		&ParserList{},
	)
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}
