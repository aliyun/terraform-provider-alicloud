package quickbi

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the quickbi service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "quickbi-public",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_quick_bi_user",
				Factory:  ResourceQuickBiUser,
			},
		},
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_quick_bi_users",
				Factory:  DataSourceQuickBiUsers,
			},
		},
	}
}
