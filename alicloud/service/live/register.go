package live

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "live",
		// alicloud_live_ai_studio is a hand-written Framework resource; the
		// legacy alicloud_live_caster / alicloud_live_domain remain SDKv2.
		Resources: []conns.Resource{
			{
				TypeName: "alicloud_live_ai_studio",
				Factory:  NewLiveAIStudioResource,
			},
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
