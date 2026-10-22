package dts

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the DTS service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Dts",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_dts_subscription_job",
				Factory:  ResourceSubscriptionJob,
			},
			{
				TypeName: "alicloud_dts_synchronization_job",
				Factory:  ResourceSynchronizationJob,
			},
		},
	}
}
