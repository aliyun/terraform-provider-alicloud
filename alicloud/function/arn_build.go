package function

import (
	"context"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/arn"
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
		MarkdownDescription: "Builds an ARN of the form `acs:<service>:<region>:<account_id>:<resource>`. The resource may itself contain colons.",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:                "service",
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
				Name:                "resource",
				MarkdownDescription: "Service-specific resource path, typically composed of a resource type and identifier, such as `instance/i-bp1234567890abcdef`.",
			},
		},
		Return: function.StringReturn{},
	}
}

func (f *arnBuildFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var service, region, accountID, resource string

	resp.Error = function.ConcatFuncErrors(req.Arguments.Get(ctx, &service, &region, &accountID, &resource))
	if resp.Error != nil {
		return
	}

	// String does not validate, so a caller passing an empty or colon-bearing component
	// still gets a malformed ARN back rather than an error — the documented behaviour of
	// this function, and the reason it can share the grammar with arn.Parse rather than
	// restating it.
	result := arn.ARN{
		Service:   service,
		Region:    region,
		AccountID: accountID,
		Resource:  resource,
	}

	resp.Error = function.ConcatFuncErrors(resp.Result.Set(ctx, result.String()))
}
