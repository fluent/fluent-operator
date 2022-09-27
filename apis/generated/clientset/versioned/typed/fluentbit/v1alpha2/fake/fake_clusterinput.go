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
// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1alpha2 "github.com/fluent/fluent-operator/apis/fluentbit/v1alpha2"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeClusterInputs implements ClusterInputInterface
type FakeClusterInputs struct {
	Fake *FakeFluentbitV1alpha2
}

var clusterinputsResource = schema.GroupVersionResource{Group: "fluentbit.fluent.io", Version: "v1alpha2", Resource: "clusterinputs"}

var clusterinputsKind = schema.GroupVersionKind{Group: "fluentbit.fluent.io", Version: "v1alpha2", Kind: "ClusterInput"}

// Get takes name of the clusterInput, and returns the corresponding clusterInput object, and an error if there is any.
func (c *FakeClusterInputs) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha2.ClusterInput, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(clusterinputsResource, name), &v1alpha2.ClusterInput{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha2.ClusterInput), err
}

// List takes label and field selectors, and returns the list of ClusterInputs that match those selectors.
func (c *FakeClusterInputs) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha2.ClusterInputList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(clusterinputsResource, clusterinputsKind, opts), &v1alpha2.ClusterInputList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha2.ClusterInputList{ListMeta: obj.(*v1alpha2.ClusterInputList).ListMeta}
	for _, item := range obj.(*v1alpha2.ClusterInputList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested clusterInputs.
func (c *FakeClusterInputs) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(clusterinputsResource, opts))
}

// Create takes the representation of a clusterInput and creates it.  Returns the server's representation of the clusterInput, and an error, if there is any.
func (c *FakeClusterInputs) Create(ctx context.Context, clusterInput *v1alpha2.ClusterInput, opts v1.CreateOptions) (result *v1alpha2.ClusterInput, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(clusterinputsResource, clusterInput), &v1alpha2.ClusterInput{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha2.ClusterInput), err
}

// Update takes the representation of a clusterInput and updates it. Returns the server's representation of the clusterInput, and an error, if there is any.
func (c *FakeClusterInputs) Update(ctx context.Context, clusterInput *v1alpha2.ClusterInput, opts v1.UpdateOptions) (result *v1alpha2.ClusterInput, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(clusterinputsResource, clusterInput), &v1alpha2.ClusterInput{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha2.ClusterInput), err
}

// Delete takes name of the clusterInput and deletes it. Returns an error if one occurs.
func (c *FakeClusterInputs) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteActionWithOptions(clusterinputsResource, name, opts), &v1alpha2.ClusterInput{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeClusterInputs) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(clusterinputsResource, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha2.ClusterInputList{})
	return err
}

// Patch applies the patch and returns the patched clusterInput.
func (c *FakeClusterInputs) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha2.ClusterInput, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(clusterinputsResource, name, pt, data, subresources...), &v1alpha2.ClusterInput{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha2.ClusterInput), err
}
