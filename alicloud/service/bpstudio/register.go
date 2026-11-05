package bpstudio

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage wires the BPStudio SDKv2 resources and data sources into the
// shared provider registration.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "BPStudio",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_bp_studio_application",
				Factory:  ResourceApplication,
			},
		},
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_bp_studio_applications",
				Factory:  DataSourceApplications,
			},
		},
	}
}
