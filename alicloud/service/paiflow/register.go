// Package paiflow registers the PAIFlow (2021-02-02) resources served by the
// Framework half of the provider. The legacy SDK v2 PAIFlow pipeline resource
// stays in package alicloud; new resources come in here as Framework-native.
package paiflow

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// ServicePackage declares the PAIFlow resources the Framework provider serves.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "PaiFlow",
		Resources: []conns.Resource{
			{
				TypeName: "alicloud_pai_run",
				Factory:  NewPaiRunResource,
			},
		},
	}
}

// NewPaiRunResource returns the PAIFlow pipeline run resource. The factory
// returns a pointer so the fwadapt client-injection guard accepts it.
func NewPaiRunResource() resource.Resource {
	return &paiRunResource{}
}
