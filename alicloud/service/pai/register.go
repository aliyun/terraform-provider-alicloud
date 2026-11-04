package pai

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage registers the PAI service resources relocated from the root
// alicloud package. The factories stay package-private; ServicePackage is in
// the same package and may call them directly.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "PAI",
		SDKResources: []conns.SDKResource{
			{TypeName: "alicloud_pai_service", Factory: resourceAliCloudPaiService},
		},
	}
}
