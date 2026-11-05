package videosurveillancesystem

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the vs service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "videosurveillancesystem",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_video_surveillance_system_group",
				Factory:  ResourceVideoSurveillanceSystemGroup,
			},
		},
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_vs_service",
				Factory:  DataSourceVsService,
			},
			{
				TypeName: "alicloud_video_surveillance_system_groups",
				Factory:  DataSourceVideoSurveillanceSystemGroups,
			},
		},
	}
}
