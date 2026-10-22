package esa

import "github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"

func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{Name: "esa", SDKResources: []conns.SDKResource{{TypeName: "alicloud_esa_image_transform", Factory: ResourceImageTransform}}}
}
