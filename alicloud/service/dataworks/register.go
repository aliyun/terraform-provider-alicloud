package dataworks

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the DataWorks service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "DataWorks",
		Resources: []conns.Resource{
			{
				TypeName: "alicloud_dataworks_table",
				Factory:  NewTableResource,
			},
		},
	}
}
