package cloudfirewall

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "CloudFirewall",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_cloud_firewall_instance_v2",
				Factory:  ResourceCloudFirewallInstanceV2,
			},
		},
	}
}
