package paiworkspace

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "PAIWorkspace",
		Resources: []conns.Resource{
			{TypeName: "alicloud_pai_workspace_dataset_job_config", Factory: NewDatasetJobConfigResource},
		},
	}
}
