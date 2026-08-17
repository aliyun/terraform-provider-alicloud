package function_test

import (
	"context"
	"testing"

	tffunction "github.com/aliyun/terraform-provider-alicloud/alicloud/function"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/framework"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	sdkschema "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestARNBuildFunction_Metadata(t *testing.T) {
	resp := &function.MetadataResponse{}

	tffunction.NewARNBuildFunction().Metadata(context.Background(), function.MetadataRequest{}, resp)

	if resp.Name != "arn_build" {
		t.Errorf("unexpected name: got %q, want %q", resp.Name, "arn_build")
	}
}

func TestARNBuildFunction_Definition(t *testing.T) {
	resp := &function.DefinitionResponse{}

	tffunction.NewARNBuildFunction().Definition(context.Background(), function.DefinitionRequest{}, resp)

	if len(resp.Definition.Parameters) != 4 {
		t.Fatalf("expected 4 parameters, got %d", len(resp.Definition.Parameters))
	}

	expected := []string{"service", "region", "account_id", "resource"}
	for i, param := range resp.Definition.Parameters {
		if got := param.GetName(); got != expected[i] {
			t.Errorf("parameter %d: got %q, want %q", i, got, expected[i])
		}
		// The Registry and the language server render MarkdownDescription, not
		// Description, so every parameter must carry the former.
		if got := param.GetMarkdownDescription(); got == "" {
			t.Errorf("parameter %d (%s): empty MarkdownDescription", i, expected[i])
		}
		// Every component is a plain string; anything else changes the documented
		// signature and silently breaks existing calls.
		if _, ok := param.(function.StringParameter); !ok {
			t.Errorf("parameter %d (%s): unexpected type %T", i, expected[i], param)
		}
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

	if _, ok := resp.Definition.Return.(function.StringReturn); !ok {
		t.Errorf("unexpected return type: %T", resp.Definition.Return)
	}
}

func TestARNBuildFunction_Run(t *testing.T) {
	cases := []struct {
		service, region, accountID, resource, expected string
	}{
		{"ram", "", "123456789012****", "role/MyRole",
			"acs:ram::123456789012****:role/MyRole"},
		{"ecs", "cn-hangzhou", "*", "instance/i-123456",
			"acs:ecs:cn-hangzhou:*:instance/i-123456"},
		{"oss", "*", "*", "my-bucket/*",
			"acs:oss:*:*:my-bucket/*"},
		{"fc", "cn-hangzhou", "123456789012****", "services/foo.LATEST/functions/bar",
			"acs:fc:cn-hangzhou:123456789012****:services/foo.LATEST/functions/bar"},
		{"kms", "*", "*", "*",
			"acs:kms:*:*:*"},
		{"mns", "cn-hangzhou", "123456789012****", "/queues/my-queue/messages",
			"acs:mns:cn-hangzhou:123456789012****:/queues/my-queue/messages"},
		// The resource is the last section, so colons inside it are part of the
		// resource path and must survive verbatim.
		{"log", "cn-hangzhou", "123456789012****", "project/my-project:logstore/my-logstore",
			"acs:log:cn-hangzhou:123456789012****:project/my-project:logstore/my-logstore"},
		// The two cases below pin down deliberate behaviour, not desirable output.
		// Like the AWS provider's arn_build, this function is a formatter and not a
		// validator: it never rejects an argument, so a caller that passes an empty
		// or colon-bearing component gets a malformed ARN back rather than an error.
		// Changing that is a breaking change and must break these cases first.
		{"", "", "", "", "acs::::"},
		{"ecs:extra", "cn-hangzhou", "123456789012****", "instance/i-123456",
			"acs:ecs:extra:cn-hangzhou:123456789012****:instance/i-123456"},
	}

	f := tffunction.NewARNBuildFunction()

	for _, c := range cases {
		req := function.RunRequest{
			Arguments: function.NewArgumentsData([]attr.Value{
				types.StringValue(c.service),
				types.StringValue(c.region),
				types.StringValue(c.accountID),
				types.StringValue(c.resource),
			}),
		}
		resp := &function.RunResponse{
			Result: function.NewResultData(types.StringNull()),
		}

		f.Run(context.Background(), req, resp)

		if resp.Error != nil {
			t.Errorf("arn_build(%q, %q, %q, %q) unexpected error: %v",
				c.service, c.region, c.accountID, c.resource, resp.Error)
			continue
		}

		if got := resp.Result.Value(); !got.Equal(types.StringValue(c.expected)) {
			t.Errorf("arn_build(%q, %q, %q, %q) = %v, want %q",
				c.service, c.region, c.accountID, c.resource, got, c.expected)
		}
	}
}

func TestARNBuildFunction_HCL(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		// Provider-defined functions do not exist before Terraform 1.8; without this
		// the test fails on an older CLI with a parse error on the provider:: prefix
		// rather than skipping.
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_8_0),
		},
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: arnBuildHCLConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckOutput("role_arn", "acs:ram::123456789012****:role/example"),
					resource.TestCheckOutput("instance_arn", "acs:ecs:cn-hangzhou:123456789012****:instance/i-bp1234567890abcdef"),
					resource.TestCheckOutput("bucket_arn", "acs:oss:*:*:my-bucket/*"),
					resource.TestCheckOutput("function_arn", "acs:fc:cn-hangzhou:123456789012****:services/foo.LATEST/functions/bar"),
					resource.TestCheckOutput("queue_arn", "acs:mns:cn-hangzhou:123456789012****:/queues/my-queue/messages"),
				),
			},
		},
	})
}

func protoV5ProviderFactories() map[string]func() (tfprotov5.ProviderServer, error) {
	return map[string]func() (tfprotov5.ProviderServer, error){
		"alicloud": providerserver.NewProtocol5WithError(framework.NewProvider(&sdkschema.Provider{})),
	}
}

const arnBuildHCLConfig = `
output "role_arn" {
  value = provider::alicloud::arn_build("ram", "", "123456789012****", "role/example")
}

output "instance_arn" {
  value = provider::alicloud::arn_build("ecs", "cn-hangzhou", "123456789012****", "instance/i-bp1234567890abcdef")
}

output "bucket_arn" {
  value = provider::alicloud::arn_build("oss", "*", "*", "my-bucket/*")
}

output "function_arn" {
  value = provider::alicloud::arn_build("fc", "cn-hangzhou", "123456789012****", "services/foo.LATEST/functions/bar")
}

output "queue_arn" {
  value = provider::alicloud::arn_build("mns", "cn-hangzhou", "123456789012****", "/queues/my-queue/messages")
}
`
