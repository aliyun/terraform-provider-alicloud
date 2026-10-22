package oss

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the OSS service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Oss",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_oss_access_point",
				Factory:  ResourceAccessPoint,
			},
		},
		Resources: []conns.Resource{
			{
				TypeName: "alicloud_oss_access_point_policy",
				Factory:  NewAccessPointPolicyResource,
			},
		},
	}
}
