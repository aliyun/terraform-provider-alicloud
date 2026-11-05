package devopsrdc

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the devops-rdc service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "devops-rdc",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_rdc_organization",
				Factory:  ResourceRdcOrganization,
			},
		},
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_rdc_organizations",
				Factory:  DataSourceRdcOrganizations,
			},
		},
	}
}
