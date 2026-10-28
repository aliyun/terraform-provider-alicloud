package cloudphone

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the Cloudphone service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "cloudphone",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_cloud_phone_image",
				Factory:  ResourceImage,
			},
			{
				TypeName: "alicloud_cloud_phone_instance",
				Factory:  ResourceInstance,
			},
			{
				TypeName: "alicloud_cloud_phone_instance_group",
				Factory:  ResourceInstanceGroup,
			},
			{
				TypeName: "alicloud_cloud_phone_key_pair",
				Factory:  ResourceKeyPair,
			},
			{
				TypeName: "alicloud_cloud_phone_policy",
				Factory:  ResourcePolicy,
			},
		},
	}
}
