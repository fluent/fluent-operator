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

// FakeClusterFilters implements ClusterFilterInterface
type FakeClusterFilters struct {
	Fake *FakeFluentbitV1alpha2
}

var clusterfiltersResource = schema.GroupVersionResource{Group: "fluentbit.fluent.io", Version: "v1alpha2", Resource: "clusterfilters"}

var clusterfiltersKind = schema.GroupVersionKind{Group: "fluentbit.fluent.io", Version: "v1alpha2", Kind: "ClusterFilter"}

// Get takes name of the clusterFilter, and returns the corresponding clusterFilter object, and an error if there is any.
func (c *FakeClusterFilters) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha2.ClusterFilter, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(clusterfiltersResource, name), &v1alpha2.ClusterFilter{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha2.ClusterFilter), err
}

// List takes label and field selectors, and returns the list of ClusterFilters that match those selectors.
func (c *FakeClusterFilters) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha2.ClusterFilterList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(clusterfiltersResource, clusterfiltersKind, opts), &v1alpha2.ClusterFilterList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha2.ClusterFilterList{ListMeta: obj.(*v1alpha2.ClusterFilterList).ListMeta}
	for _, item := range obj.(*v1alpha2.ClusterFilterList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested clusterFilters.
func (c *FakeClusterFilters) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(clusterfiltersResource, opts))
}

// Create takes the representation of a clusterFilter and creates it.  Returns the server's representation of the clusterFilter, and an error, if there is any.
func (c *FakeClusterFilters) Create(ctx context.Context, clusterFilter *v1alpha2.ClusterFilter, opts v1.CreateOptions) (result *v1alpha2.ClusterFilter, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(clusterfiltersResource, clusterFilter), &v1alpha2.ClusterFilter{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha2.ClusterFilter), err
}

// Update takes the representation of a clusterFilter and updates it. Returns the server's representation of the clusterFilter, and an error, if there is any.
func (c *FakeClusterFilters) Update(ctx context.Context, clusterFilter *v1alpha2.ClusterFilter, opts v1.UpdateOptions) (result *v1alpha2.ClusterFilter, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(clusterfiltersResource, clusterFilter), &v1alpha2.ClusterFilter{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha2.ClusterFilter), err
}

// Delete takes name of the clusterFilter and deletes it. Returns an error if one occurs.
func (c *FakeClusterFilters) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteActionWithOptions(clusterfiltersResource, name, opts), &v1alpha2.ClusterFilter{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeClusterFilters) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(clusterfiltersResource, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha2.ClusterFilterList{})
	return err
}

// Patch applies the patch and returns the patched clusterFilter.
func (c *FakeClusterFilters) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha2.ClusterFilter, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(clusterfiltersResource, name, pt, data, subresources...), &v1alpha2.ClusterFilter{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha2.ClusterFilter), err
}
