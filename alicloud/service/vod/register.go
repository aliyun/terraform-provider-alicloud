package vod

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the VOD service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "vod",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_vod_domain",
				Factory:  ResourceDomain,
			},
			{
				TypeName: "alicloud_vod_editing_project",
				Factory:  ResourceEditingProject,
			},
		},
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_vod_domains",
				Factory:  DataSourceDomains,
			},
		},
	}
}
