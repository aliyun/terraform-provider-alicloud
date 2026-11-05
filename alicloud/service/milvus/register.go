package milvus

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the Milvus service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "milvus",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_milvus_instance",
				Factory:  ResourceInstance,
			},
		},
	}
}
