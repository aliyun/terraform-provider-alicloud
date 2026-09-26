package pai

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the PAI service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "AIWorkSpace",
		DataSources: []conns.DataSource{
			{
				TypeName: "alicloud_pai_workspace_prompts",
				Factory:  NewPromptsDataSource,
			},
		},
		Resources: []conns.Resource{
			{
				TypeName: "alicloud_pai_workspace_prompt",
				Factory:  NewPromptResource,
			},
		},
	}
}
