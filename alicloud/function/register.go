package function

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

// ServicePackage returns the registration declaration of the provider-level
// utility functions.
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{
		Name: "Core",
		Functions: []conns.Function{
			{
				TypeName: "arn_build",
				Factory:  NewARNBuildFunction,
			},
			{
				TypeName: "arn_parse",
				Factory:  NewARNParseFunction,
			},
		},
	}
}
