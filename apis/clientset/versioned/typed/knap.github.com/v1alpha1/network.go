// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"
	"time"

	scheme "github.com/tliron/knap/apis/clientset/versioned/scheme"
	v1alpha1 "github.com/tliron/knap/resources/knap.github.com/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// NetworksGetter has a method to return a NetworkInterface.
// A group's client should implement this interface.
type NetworksGetter interface {
	Networks(namespace string) NetworkInterface
}

// NetworkInterface has methods to work with Network resources.
type NetworkInterface interface {
	Create(ctx context.Context, network *v1alpha1.Network, opts v1.CreateOptions) (*v1alpha1.Network, error)
	Update(ctx context.Context, network *v1alpha1.Network, opts v1.UpdateOptions) (*v1alpha1.Network, error)
	UpdateStatus(ctx context.Context, network *v1alpha1.Network, opts v1.UpdateOptions) (*v1alpha1.Network, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha1.Network, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha1.NetworkList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.Network, err error)
	NetworkExpansion
}

// networks implements NetworkInterface
type networks struct {
	client rest.Interface
	ns     string
}

// newNetworks returns a Networks
func newNetworks(c *KnapV1alpha1Client, namespace string) *networks {
	return &networks{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the network, and returns the corresponding network object, and an error if there is any.
func (c *networks) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.Network, err error) {
	result = &v1alpha1.Network{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("networks").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Networks that match those selectors.
func (c *networks) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.NetworkList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.NetworkList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("networks").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested networks.
func (c *networks) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("networks").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a network and creates it.  Returns the server's representation of the network, and an error, if there is any.
func (c *networks) Create(ctx context.Context, network *v1alpha1.Network, opts v1.CreateOptions) (result *v1alpha1.Network, err error) {
	result = &v1alpha1.Network{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("networks").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(network).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a network and updates it. Returns the server's representation of the network, and an error, if there is any.
func (c *networks) Update(ctx context.Context, network *v1alpha1.Network, opts v1.UpdateOptions) (result *v1alpha1.Network, err error) {
	result = &v1alpha1.Network{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("networks").
		Name(network.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(network).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *networks) UpdateStatus(ctx context.Context, network *v1alpha1.Network, opts v1.UpdateOptions) (result *v1alpha1.Network, err error) {
	result = &v1alpha1.Network{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("networks").
		Name(network.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(network).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the network and deletes it. Returns an error if one occurs.
func (c *networks) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("networks").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *networks) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("networks").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched network.
func (c *networks) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.Network, err error) {
	result = &v1alpha1.Network{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("networks").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
