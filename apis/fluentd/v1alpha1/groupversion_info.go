/*
Copyright 2022.

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

// Package v1alpha1 contains API Schema definitions for the  v1alpha1 API group
// +kubebuilder:object:generate=true
// +groupName=fluentd.fluent.io
package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// localSchemeBuilder replaces the deprecated sigs.k8s.io/controller-runtime/pkg/scheme.Builder,
// registering types with a fixed GroupVersion via init() calls in each type file.
type localSchemeBuilder struct {
	gv    schema.GroupVersion
	types []runtime.Object
}

func (b *localSchemeBuilder) Register(objects ...runtime.Object) *localSchemeBuilder {
	b.types = append(b.types, objects...)
	return b
}

func (b *localSchemeBuilder) AddToScheme(s *runtime.Scheme) error {
	s.AddKnownTypes(b.gv, b.types...)
	metav1.AddToGroupVersion(s, b.gv)
	return nil
}

var (
	// GroupVersion is group version used to register these objects
	SchemeGroupVersion = schema.GroupVersion{Group: "fluentd.fluent.io", Version: "v1alpha1"}

	// SchemeBuilder is used to add go types to the GroupVersionKind scheme
	SchemeBuilder = &localSchemeBuilder{gv: SchemeGroupVersion}

	// AddToScheme adds the types in this group-version to the given scheme.
	AddToScheme = SchemeBuilder.AddToScheme
)
