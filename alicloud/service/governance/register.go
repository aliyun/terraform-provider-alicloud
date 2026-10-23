package governance

import "github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"

// ServicePackage returns the registration declaration of the Governance service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "governance",
		Resources: []conns.Resource{
			{
				TypeName: "alicloud_governance_service_activation",
				Factory:  NewServiceActivationResource,
			},
		},
	}
}
