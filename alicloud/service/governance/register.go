package governance

import "github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"

// ServicePackage returns the registration declaration of the Governance service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "governance",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_governance_account",
				Factory:  ResourceGovernanceAccount,
			},
			{
				TypeName: "alicloud_governance_baseline",
				Factory:  ResourceGovernanceBaseline,
			},
		},
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_governance_baselines",
				Factory:  DataSourceGovernanceBaselines,
			},
		},
		Resources: []conns.Resource{
			{
				TypeName: "alicloud_governance_service_activation",
				Factory:  NewServiceActivationResource,
			},
		},
	}
}
