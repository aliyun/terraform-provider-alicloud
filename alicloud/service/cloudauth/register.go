package cloudauth

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the Cloudauth service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Cloudauth",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_cloudauth_face_config",
				Factory:  ResourceFaceConfig,
			},
		},
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_cloudauth_face_configs",
				Factory:  DataSourceFaceConfigs,
			},
		},
	}
}
