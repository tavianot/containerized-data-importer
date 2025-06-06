/*
Copyright 2018 The CDI Authors.

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

package v1beta1

import (
	"context"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	gentype "k8s.io/client-go/gentype"
	v1beta1 "kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"
	scheme "kubevirt.io/containerized-data-importer/pkg/client/clientset/versioned/scheme"
)

// CDIConfigsGetter has a method to return a CDIConfigInterface.
// A group's client should implement this interface.
type CDIConfigsGetter interface {
	CDIConfigs() CDIConfigInterface
}

// CDIConfigInterface has methods to work with CDIConfig resources.
type CDIConfigInterface interface {
	Create(ctx context.Context, cDIConfig *v1beta1.CDIConfig, opts v1.CreateOptions) (*v1beta1.CDIConfig, error)
	Update(ctx context.Context, cDIConfig *v1beta1.CDIConfig, opts v1.UpdateOptions) (*v1beta1.CDIConfig, error)
	// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
	UpdateStatus(ctx context.Context, cDIConfig *v1beta1.CDIConfig, opts v1.UpdateOptions) (*v1beta1.CDIConfig, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1beta1.CDIConfig, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1beta1.CDIConfigList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1beta1.CDIConfig, err error)
	CDIConfigExpansion
}

// cDIConfigs implements CDIConfigInterface
type cDIConfigs struct {
	*gentype.ClientWithList[*v1beta1.CDIConfig, *v1beta1.CDIConfigList]
}

// newCDIConfigs returns a CDIConfigs
func newCDIConfigs(c *CdiV1beta1Client) *cDIConfigs {
	return &cDIConfigs{
		gentype.NewClientWithList[*v1beta1.CDIConfig, *v1beta1.CDIConfigList](
			"cdiconfigs",
			c.RESTClient(),
			scheme.ParameterCodec,
			"",
			func() *v1beta1.CDIConfig { return &v1beta1.CDIConfig{} },
			func() *v1beta1.CDIConfigList { return &v1beta1.CDIConfigList{} }),
	}
}
