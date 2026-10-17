package paistudio

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the PaiStudio service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "PaiStudio",
		Resources: []conns.Resource{
			{
				TypeName: "alicloud_pai_algorithm_version",
				Factory:  NewAlgorithmVersionResource,
			},
		},
	}
}
