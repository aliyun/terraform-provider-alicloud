package function_test

import (
	"context"
	"regexp"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tffunction "github.com/aliyun/terraform-provider-alicloud/alicloud/function"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestARNParseFunction_Metadata(t *testing.T) {
	resp := &function.MetadataResponse{}

	tffunction.NewARNParseFunction().Metadata(context.Background(), function.MetadataRequest{}, resp)

	if resp.Name != "arn_parse" {
		t.Errorf("unexpected name: got %q, want %q", resp.Name, "arn_parse")
	}
}

func TestARNParseFunction_Definition(t *testing.T) {
	ctx := context.Background()
	resp := &function.DefinitionResponse{}

	tffunction.NewARNParseFunction().Definition(ctx, function.DefinitionRequest{}, resp)

	if len(resp.Definition.Parameters) != 1 {
		t.Fatalf("expected 1 parameter, got %d", len(resp.Definition.Parameters))
	}

	param := resp.Definition.Parameters[0]
	if got := param.GetName(); got != "arn" {
		t.Errorf("parameter 0: got %q, want %q", got, "arn")
	}
	// The Registry and the language server render MarkdownDescription, not
	// Description, so the parameter must carry the former.
	if got := param.GetMarkdownDescription(); got == "" {
		t.Error("parameter 0 (arn): empty MarkdownDescription")
	}
	if _, ok := param.(function.StringParameter); !ok {
		t.Errorf("parameter 0 (arn): unexpected type %T", param)
	}

	if resp.Definition.Summary == "" {
		t.Error("empty Summary on the function definition")
	}

	if resp.Definition.MarkdownDescription == "" {
		t.Error("empty MarkdownDescription on the function definition")
	}

	if resp.Definition.VariadicParameter != nil {
		t.Errorf("unexpected variadic parameter: %v", resp.Definition.VariadicParameter)
	}

	// The returned object's attributes are arn_build's parameters, so the output of
	// one function feeds straight into the other. Deriving the expectation from
	// arn_build's own definition rather than restating it means a rename on either
	// side breaks here instead of quietly ending the composability.
	attributeTypes := arnParseReturnTypes(t)

	buildResp := &function.DefinitionResponse{}
	tffunction.NewARNBuildFunction().Definition(ctx, function.DefinitionRequest{}, buildResp)

	if len(attributeTypes) != len(buildResp.Definition.Parameters) {
		t.Fatalf("the return object has %d attributes, but arn_build takes %d parameters",
			len(attributeTypes), len(buildResp.Definition.Parameters))
	}

	for _, buildParam := range buildResp.Definition.Parameters {
		name := buildParam.GetName()

		attrType, ok := attributeTypes[name]
		if !ok {
			t.Errorf("the return object has no %q attribute, but arn_build takes one", name)
			continue
		}
		// Every section of an ARN is a plain string; anything else changes the
		// documented signature and breaks existing calls.
		if !attrType.Equal(types.StringType) {
			t.Errorf("return attribute %q: got type %s, want string", name, attrType)
		}
	}
}

func TestARNParseFunction_Run(t *testing.T) {
	cases := []struct {
		arn                                  string
		service, region, accountID, resource string
	}{
		{"acs:ram::123456789012****:role/MyRole",
			"ram", "", "123456789012****", "role/MyRole"},
		{"acs:ecs:cn-hangzhou:123456789012****:instance/i-bp1234567890abcdef",
			"ecs", "cn-hangzhou", "123456789012****", "instance/i-bp1234567890abcdef"},
		{"acs:oss:*:*:my-bucket/*",
			"oss", "*", "*", "my-bucket/*"},
		{"acs:mns:cn-hangzhou:123456789012****:/queues/my-queue/messages",
			"mns", "cn-hangzhou", "123456789012****", "/queues/my-queue/messages"},
		{"acs:kms:*:*:*",
			"kms", "*", "*", "*"},
		// A wildcard policy names no account, and both empty sections have to survive
		// as empty strings rather than collapsing the ARN into four sections.
		{"acs:ram:*::*",
			"ram", "*", "", "*"},
		// The resource is the last section, so colons inside it are part of the
		// resource path and stay with it.
		{"acs:log:cn-hangzhou:123456789012****:project/my-project:logstore/my-logstore",
			"log", "cn-hangzhou", "123456789012****", "project/my-project:logstore/my-logstore"},
	}

	attributeTypes := arnParseReturnTypes(t)
	f := tffunction.NewARNParseFunction()

	for _, c := range cases {
		resp := runARNParse(t, f, c.arn)

		if resp.Error != nil {
			t.Errorf("arn_parse(%q) unexpected error: %v", c.arn, resp.Error)
			continue
		}

		want := types.ObjectValueMust(attributeTypes, map[string]attr.Value{
			"service":    types.StringValue(c.service),
			"region":     types.StringValue(c.region),
			"account_id": types.StringValue(c.accountID),
			"resource":   types.StringValue(c.resource),
		})

		if got := resp.Result.Value(); !got.Equal(want) {
			t.Errorf("arn_parse(%q) = %v, want %v", c.arn, got, want)
		}
	}
}

func TestARNParseFunction_RunErrors(t *testing.T) {
	cases := []string{
		"",
		"acs:ecs:cn-hangzhou:123456789012****",
		// The likeliest paste-o.
		"arn:aws:s3:::my-bucket",
		// A truncated ARN names no resource. "*" is how a policy matches everything.
		"acs:ecs:cn-hangzhou:123456789012****:",
	}

	f := tffunction.NewARNParseFunction()

	for _, s := range cases {
		resp := runARNParse(t, f, s)

		if resp.Error == nil {
			t.Errorf("arn_parse(%q) = %v, want an error", s, resp.Result.Value())
			continue
		}

		if resp.Error.FunctionArgument == nil {
			t.Errorf("arn_parse(%q) error %q is not attributed to an argument", s, resp.Error.Text)
			continue
		}
		if got := *resp.Error.FunctionArgument; got != 0 {
			t.Errorf("arn_parse(%q) error is attributed to argument %d, want 0", s, got)
		}
	}
}

func arnParseReturnTypes(t *testing.T) map[string]attr.Type {
	t.Helper()

	resp := &function.DefinitionResponse{}
	tffunction.NewARNParseFunction().Definition(context.Background(), function.DefinitionRequest{}, resp)

	objectReturn, ok := resp.Definition.Return.(function.ObjectReturn)
	if !ok {
		t.Fatalf("arn_parse returns %T, want function.ObjectReturn", resp.Definition.Return)
	}

	return objectReturn.AttributeTypes
}

func runARNParse(t *testing.T, f function.Function, s string) *function.RunResponse {
	t.Helper()

	req := function.RunRequest{
		Arguments: function.NewArgumentsData([]attr.Value{types.StringValue(s)}),
	}
	resp := &function.RunResponse{
		Result: function.NewResultData(types.ObjectNull(arnParseReturnTypes(t))),
	}

	f.Run(context.Background(), req, resp)

	return resp
}

func TestARNParseFunction_HCL(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		// Provider-defined functions do not exist before Terraform 1.8; without this
		// the test fails on an older CLI with a parse error on the provider:: prefix
		// rather than skipping.
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_8_0),
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: arnParseHCLConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckOutput("service", "ecs"),
					resource.TestCheckOutput("region", "cn-hangzhou"),
					resource.TestCheckOutput("account_id", "123456789012****"),
					resource.TestCheckOutput("resource", "instance/i-bp1234567890abcdef"),
					resource.TestCheckOutput("role_region", ""),
					resource.TestCheckOutput("role_resource", "role/example"),
					resource.TestCheckOutput("logstore_resource", "project/my-project:logstore/my-logstore"),
					resource.TestCheckOutput("round_trip", "acs:log:cn-hangzhou:123456789012****:project/my-project:logstore/my-logstore"),
				),
			},
			{
				Config:      arnParseInvalidHCLConfig,
				ExpectError: regexp.MustCompile("is not an ARN"),
			},
		},
	})
}

const arnParseInvalidHCLConfig = `
output "invalid" {
  value = provider::alicloud::arn_parse("arn:aws:s3:::my-bucket")
}
`

const arnParseHCLConfig = `
locals {
  instance = provider::alicloud::arn_parse("acs:ecs:cn-hangzhou:123456789012****:instance/i-bp1234567890abcdef")
  logstore = provider::alicloud::arn_parse("acs:log:cn-hangzhou:123456789012****:project/my-project:logstore/my-logstore")
  role     = provider::alicloud::arn_parse("acs:ram::123456789012****:role/example")
}

output "service" {
  value = local.instance.service
}

output "region" {
  value = local.instance.region
}

output "account_id" {
  value = local.instance.account_id
}

output "resource" {
  value = local.instance.resource
}

# A global service has no region, and the empty section must arrive as an empty
# string rather than shifting the sections along by one.
output "role_region" {
  value = local.role.region
}

output "role_resource" {
  value = local.role.resource
}

# The colons in a Log Service resource belong to the resource path.
output "logstore_resource" {
  value = local.logstore.resource
}

# arn_build is the inverse, so feeding the object straight back in has to reproduce
# the ARN the module started from.
output "round_trip" {
  value = provider::alicloud::arn_build(local.logstore.service, local.logstore.region, local.logstore.account_id, local.logstore.resource)
}
`
