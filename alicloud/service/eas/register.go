package eas

import "github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"

// ServicePackage returns the registration declaration of the EAS service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Eas",
		Resources: []conns.Resource{
			{
				TypeName: "alicloud_pai_log_download_task",
				Factory:  NewLogDownloadTaskResource,
			},
		},
	}
}
