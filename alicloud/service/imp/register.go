package imp

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the imp service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "imp",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_imp_app_template",
				Factory:  ResourceAppTemplate,
			},
		},
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_imp_app_templates",
				Factory:  DataSourceAppTemplates,
			},
		},
	}
}
