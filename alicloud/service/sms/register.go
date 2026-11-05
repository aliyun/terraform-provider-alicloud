package sms

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the SMS service.
// Name is the API product code ("Dysmsapi") used by RpcPost and the conns
// ServicePackage identifier; the package directory follows the canonical
// namespace (sms).
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Dysmsapi",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_sms_short_url",
				Factory:  ResourceSmsShortUrl,
			},
		},
	}
}
