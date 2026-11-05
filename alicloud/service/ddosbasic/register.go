package ddosbasic

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "antiddos-public",
		SDKResources: []conns.SDKResource{
			{TypeName: "alicloud_ddos_basic_defense_threshold", Factory: ResourceDdosBasicDefenseThreshold},
			{TypeName: "alicloud_ddos_basic_threshold", Factory: ResourceDdosBasicThreshold},
		},
	}
}
