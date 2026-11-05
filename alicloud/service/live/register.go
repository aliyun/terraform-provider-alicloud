package live

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "live",
		Resources: []conns.Resource{
			{TypeName: "alicloud_live_private_line", Factory: NewLivePrivateLineResource},
		},
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
