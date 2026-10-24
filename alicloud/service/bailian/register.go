package bailian

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the Bailian (SFM) service.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Bailian",
		Resources: []conns.Resource{
			{
				TypeName: "alicloud_bailian_prompt_template",
				Factory:  NewPromptTemplateResource,
			},
			{
				TypeName: "alicloud_bailian_memory",
				Factory:  NewMemoryResource,
			},
			{
				TypeName: "alicloud_bailian_memory_node",
				Factory:  NewMemoryNodeResource,
			},
			{
				TypeName: "alicloud_bailian_connector",
				Factory:  NewConnectorResource,
			},
		},
	}
}
