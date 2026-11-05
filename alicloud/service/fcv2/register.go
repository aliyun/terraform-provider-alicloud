package fcv2

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the FCV2 service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "FCV2",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_fcv2_function",
				Factory:  ResourceFcv2Function,
			},
		},
		Resources: []conns.Resource{
			{
				TypeName: "alicloud_fcv2_function_on_demand_config",
				Factory:  NewFunctionOnDemandConfigResource,
			},
		},
	}
}
