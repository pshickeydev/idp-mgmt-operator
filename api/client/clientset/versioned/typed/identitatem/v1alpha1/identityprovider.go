// Copyright Contributors to the Open Cluster Management project

// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"
	"time"

	scheme "github.com/identitatem/idp-mgmt-operator/api/client/clientset/versioned/scheme"
	v1alpha1 "github.com/identitatem/idp-mgmt-operator/api/identitatem/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// IdentityProvidersGetter has a method to return a IdentityProviderInterface.
// A group's client should implement this interface.
type IdentityProvidersGetter interface {
	IdentityProviders(namespace string) IdentityProviderInterface
}

// IdentityProviderInterface has methods to work with IdentityProvider resources.
type IdentityProviderInterface interface {
	Create(ctx context.Context, identityProvider *v1alpha1.IdentityProvider, opts v1.CreateOptions) (*v1alpha1.IdentityProvider, error)
	Update(ctx context.Context, identityProvider *v1alpha1.IdentityProvider, opts v1.UpdateOptions) (*v1alpha1.IdentityProvider, error)
	UpdateStatus(ctx context.Context, identityProvider *v1alpha1.IdentityProvider, opts v1.UpdateOptions) (*v1alpha1.IdentityProvider, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha1.IdentityProvider, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha1.IdentityProviderList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.IdentityProvider, err error)
	IdentityProviderExpansion
}

// identityProviders implements IdentityProviderInterface
type identityProviders struct {
	client rest.Interface
	ns     string
}

// newIdentityProviders returns a IdentityProviders
func newIdentityProviders(c *IdentitatemV1alpha1Client, namespace string) *identityProviders {
	return &identityProviders{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the identityProvider, and returns the corresponding identityProvider object, and an error if there is any.
func (c *identityProviders) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.IdentityProvider, err error) {
	result = &v1alpha1.IdentityProvider{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("identityproviders").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of IdentityProviders that match those selectors.
func (c *identityProviders) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.IdentityProviderList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.IdentityProviderList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("identityproviders").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested identityProviders.
func (c *identityProviders) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("identityproviders").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a identityProvider and creates it.  Returns the server's representation of the identityProvider, and an error, if there is any.
func (c *identityProviders) Create(ctx context.Context, identityProvider *v1alpha1.IdentityProvider, opts v1.CreateOptions) (result *v1alpha1.IdentityProvider, err error) {
	result = &v1alpha1.IdentityProvider{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("identityproviders").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(identityProvider).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a identityProvider and updates it. Returns the server's representation of the identityProvider, and an error, if there is any.
func (c *identityProviders) Update(ctx context.Context, identityProvider *v1alpha1.IdentityProvider, opts v1.UpdateOptions) (result *v1alpha1.IdentityProvider, err error) {
	result = &v1alpha1.IdentityProvider{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("identityproviders").
		Name(identityProvider.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(identityProvider).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *identityProviders) UpdateStatus(ctx context.Context, identityProvider *v1alpha1.IdentityProvider, opts v1.UpdateOptions) (result *v1alpha1.IdentityProvider, err error) {
	result = &v1alpha1.IdentityProvider{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("identityproviders").
		Name(identityProvider.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(identityProvider).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the identityProvider and deletes it. Returns an error if one occurs.
func (c *identityProviders) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("identityproviders").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *identityProviders) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("identityproviders").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched identityProvider.
func (c *identityProviders) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.IdentityProvider, err error) {
	result = &v1alpha1.IdentityProvider{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("identityproviders").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}