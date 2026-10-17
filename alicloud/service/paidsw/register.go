package paidsw

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the PAI DSW service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "PaiDswApi",
		Resources: []conns.Resource{
			{
				TypeName: "alicloud_pai_dsw_temp_file_task",
				Factory:  NewTempFileTaskResource,
			},
		},
	}
}
