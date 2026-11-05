package hologram

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the Hologram service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Hologram",
		Resources: []conns.Resource{
			{TypeName: "alicloud_hologram_holo_web_login_setting", Factory: NewHoloWebLoginSettingResource},
		},
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_hologram_instance",
				Factory:  ResourceInstance,
			},
		},
	}
}
