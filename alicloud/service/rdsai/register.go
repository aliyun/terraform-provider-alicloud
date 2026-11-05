package rdsai

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage is a placeholder for the rdsai service; declarations move here when the product adopts the split layout.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "RdsAi",
		SDKResources: []conns.SDKResource{
			{TypeName: "alicloud_rds_ai_instance", Factory: ResourceRdsAiInstance},
		},
	}
}
