package function

import (
	"context"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/arn"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ function.Function = &arnParseFunction{}

// arnParseAttrTypes describes the object arn_parse returns. Its keys are the
// parameter names of arn_build, in the same order, so the output of one feeds
// straight into the other.
var arnParseAttrTypes = map[string]attr.Type{
	"service":    types.StringType,
	"region":     types.StringType,
	"account_id": types.StringType,
	"resource":   types.StringType,
}

// NewARNParseFunction returns the provider-defined function implementation of
// provider::alicloud::arn_parse.
func NewARNParseFunction() function.Function {
	return &arnParseFunction{}
}

type arnParseFunction struct{}

func (f *arnParseFunction) Metadata(ctx context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "arn_parse"
}

func (f *arnParseFunction) Definition(ctx context.Context, req function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:             "Parse an Alibaba Cloud Resource Name (ARN) into its constituent parts.",
		MarkdownDescription: "Splits an ARN of the form `acs:<service>:<region>:<account_id>:<resource>` into an object. The resource is the last section, so any colons it contains are kept as part of it. Returns an error if the ARN is not of that form; unlike `arn_build`, which formats whatever it is given, this function validates.",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:                "arn",
				MarkdownDescription: "ARN to parse, such as `acs:ecs:cn-hangzhou:123456789012****:instance/i-bp1234567890abcdef`.",
			},
		},
		Return: function.ObjectReturn{
			AttributeTypes: arnParseAttrTypes,
		},
	}
}

func (f *arnParseFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var s string

	resp.Error = function.ConcatFuncErrors(req.Arguments.Get(ctx, &s))
	if resp.Error != nil {
		return
	}

	parsed, err := arn.Parse(s)
	if err != nil {
		// Attributed to the argument, not to the call, so a configuration that parses
		// several ARNs points at the one that is malformed.
		resp.Error = function.NewArgumentFuncError(0, err.Error())
		return
	}

	result, diags := types.ObjectValue(arnParseAttrTypes, map[string]attr.Value{
		"service":    types.StringValue(parsed.Service),
		"region":     types.StringValue(parsed.Region),
		"account_id": types.StringValue(parsed.AccountID),
		"resource":   types.StringValue(parsed.Resource),
	})
	resp.Error = function.ConcatFuncErrors(function.FuncErrorFromDiags(ctx, diags))
	if resp.Error != nil {
		return
	}

	resp.Error = function.ConcatFuncErrors(resp.Result.Set(ctx, result))
}
