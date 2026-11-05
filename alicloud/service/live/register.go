package live

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "live",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_live_caster",
				Factory:  ResourceLiveCaster,
			},
			{
				TypeName: "alicloud_live_domain",
				Factory:  ResourceLiveDomain,
			},
		},
	}
}
