package pailangstudio

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the PAI LangStudio
// service. The product exposes ROA-style APIs on the 2024-07-10 version.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "PAILangStudio",
		Resources: []conns.Resource{
			{
				TypeName: "alicloud_pailangstudio_knowledge_base_job",
				Factory:  NewKnowledgeBaseJobResource,
			},
		},
	}
}
