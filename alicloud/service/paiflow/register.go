package paiflow

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the PAIFlow service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "PAIFlow",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_pai_flow_pipeline",
				Factory:  ResourcePaiFlowPipeline,
			},
		},
	}
}
