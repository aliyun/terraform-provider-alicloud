package function

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/function"
)

var _ function.Function = &arnBuildFunction{}

// NewARNBuildFunction returns the provider-defined function implementation of
// provider::alicloud::arn_build.
func NewARNBuildFunction() function.Function {
	return &arnBuildFunction{}
}

type arnBuildFunction struct{}

func (f *arnBuildFunction) Metadata(ctx context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "arn_build"
}

func (f *arnBuildFunction) Definition(ctx context.Context, req function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:             "Build an Alibaba Cloud Resource Name (ARN) from its constituent parts.",
		MarkdownDescription: "Builds an ARN of the form `acs:<ram_code>:<region>:<account_id>:<relative_id>`. The relative ID may itself contain colons.",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:                "ram_code",
				MarkdownDescription: "RAM code of the Alibaba Cloud service, such as `ecs`, `oss`, `ram`, `fc`, `mns`, `fnf`, or `kms`.",
			},
			function.StringParameter{
				Name:                "region",
				MarkdownDescription: "Region of the resource, such as `cn-hangzhou`. Use `*` to match every region, or an empty string for services that carry no region component, such as `ram`.",
			},
			function.StringParameter{
				Name:                "account_id",
				MarkdownDescription: "ID of the Alibaba Cloud account.",
			},
			function.StringParameter{
				Name:                "relative_id",
				MarkdownDescription: "Service-specific resource path, typically composed of a resource type and identifier, such as `instance/i-bp1234567890abcdef`.",
			},
		},
		Return: function.StringReturn{},
	}
}

func (f *arnBuildFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var ramCode, region, accountID, relativeID string

	resp.Error = function.ConcatFuncErrors(req.Arguments.Get(ctx, &ramCode, &region, &accountID, &relativeID))
	if resp.Error != nil {
		return
	}

	result := fmt.Sprintf("acs:%s:%s:%s:%s", ramCode, region, accountID, relativeID)

	resp.Error = function.ConcatFuncErrors(resp.Result.Set(ctx, result))
}
