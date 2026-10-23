package nas

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the NAS service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Nas",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_nas_file_system",
				Factory:  ResourceFileSystem,
			},
		},
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_nas_file_systems",
				Factory:  DataSourceFileSystems,
			},
		},
	}
}
