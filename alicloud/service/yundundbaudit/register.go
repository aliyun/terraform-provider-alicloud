package yundundbaudit

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the Yundun-dbaudit service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Yundun-dbaudit",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_yundun_dbaudit_instance",
				Factory:  ResourceInstance,
			},
		},
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_yundun_dbaudit_instance",
				Factory:  DataSourceInstances,
			},
		},
	}
}
