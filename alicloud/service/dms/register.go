package dms

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the Dms service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Dms",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_dms_airflow",
				Factory:  ResourceAirflow,
			},
		},
	}
}
