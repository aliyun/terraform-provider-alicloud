package computenest

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the ComputeNest service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "ComputeNest",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_compute_nest_service_instance",
				Factory:  ResourceServiceInstance,
			},
		},
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_compute_nest_service_instances",
				Factory:  DataSourceServiceInstances,
			},
		},
	}
}
