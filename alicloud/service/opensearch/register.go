package opensearch

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the OpenSearch service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "OpenSearch",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_open_search_app_group",
				Factory:  ResourceAppGroup,
			},
		},
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_open_search_app_groups",
				Factory:  DataSourceAppGroups,
			},
		},
	}
}
