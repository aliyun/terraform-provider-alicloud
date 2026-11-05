package account

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage declares the resources and data sources owned by this package.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "account",
		SDKDataSources: []conns.SDKDataSource{
			{TypeName: "alicloud_account", Factory: dataSourceAlicloudAccount},
		},
	}
}
