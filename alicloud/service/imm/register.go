package imm

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "imm",
		SDKResources: []conns.SDKResource{
			{
				TypeName: "alicloud_imm_project",
				Factory:  resourceAlicloudImmProject,
			},
		},
		SDKDataSources: []conns.SDKDataSource{
			{
				TypeName: "alicloud_imm_projects",
				Factory:  dataSourceAlicloudImmProjects,
			},
		},
	}
}
