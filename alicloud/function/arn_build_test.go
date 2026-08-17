package function

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	fwdatasource "github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	fwprovider "github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
	sdkresource "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	sdkschema "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestARNBuildFunction_Metadata(t *testing.T) {
	resp := &function.MetadataResponse{}

	NewARNBuildFunction().Metadata(context.Background(), function.MetadataRequest{}, resp)

	if resp.Name != "arn_build" {
		t.Errorf("unexpected name: got %q, want %q", resp.Name, "arn_build")
	}
}

func TestARNBuildFunction_Definition(t *testing.T) {
	resp := &function.DefinitionResponse{}

	NewARNBuildFunction().Definition(context.Background(), function.DefinitionRequest{}, resp)

	if len(resp.Definition.Parameters) != 4 {
		t.Fatalf("expected 4 parameters, got %d", len(resp.Definition.Parameters))
	}

	expected := []string{"ram_code", "region", "account_id", "relative_id"}
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
		ramCode, region, accountID, relativeID, expected string
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
		// The relative ID is the last component, so colons inside it are part of
		// the resource path and must survive verbatim.
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

	f := NewARNBuildFunction()

	for _, c := range cases {
		req := function.RunRequest{
			Arguments: function.NewArgumentsData([]attr.Value{
				types.StringValue(c.ramCode),
				types.StringValue(c.region),
				types.StringValue(c.accountID),
				types.StringValue(c.relativeID),
			}),
		}
		resp := &function.RunResponse{
			Result: function.NewResultData(types.StringNull()),
		}

		f.Run(context.Background(), req, resp)

		if resp.Error != nil {
			t.Errorf("arn_build(%q, %q, %q, %q) unexpected error: %v",
				c.ramCode, c.region, c.accountID, c.relativeID, resp.Error)
			continue
		}

		if got := resp.Result.Value(); !got.Equal(types.StringValue(c.expected)) {
			t.Errorf("arn_build(%q, %q, %q, %q) = %v, want %q",
				c.ramCode, c.region, c.accountID, c.relativeID, got, c.expected)
		}
	}
}

// harnessProvider is a minimal framework provider whose only purpose is to serve
// arn_build to a real terraform binary. It declares no schema and no credentials,
// so the HCL test below reaches no Alibaba Cloud API.
type harnessProvider struct{}

var _ fwprovider.ProviderWithFunctions = &harnessProvider{}

func (p *harnessProvider) Metadata(ctx context.Context, req fwprovider.MetadataRequest, resp *fwprovider.MetadataResponse) {
	resp.TypeName = "alicloud"
}

func (p *harnessProvider) Schema(ctx context.Context, req fwprovider.SchemaRequest, resp *fwprovider.SchemaResponse) {
}

func (p *harnessProvider) Configure(ctx context.Context, req fwprovider.ConfigureRequest, resp *fwprovider.ConfigureResponse) {
}

func (p *harnessProvider) Resources(ctx context.Context) []func() fwresource.Resource {
	return nil
}

func (p *harnessProvider) DataSources(ctx context.Context) []func() fwdatasource.DataSource {
	return nil
}

func (p *harnessProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{NewARNBuildFunction}
}

// The required_providers block is mandatory: Terraform resolves provider-defined
// functions only for providers explicitly declared in the module. The source must
// match the reattach registration the test harness performs, which defaults to
// registry.terraform.io/hashicorp/alicloud and is overridable through
// TF_ACC_PROVIDER_HOST / TF_ACC_PROVIDER_NAMESPACE. User-facing documentation
// shows the real aliyun/alicloud source instead.
const arnBuildHCLConfig = `
terraform {
  required_providers {
    alicloud = {
      source = "hashicorp/alicloud"
    }
  }
}

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

// TestARNBuildFunction_HCL drives arn_build from real HCL through a real terraform
// binary, with the SDK v2 provider muxed behind tf5muxserver exactly as the released
// binary serves it. This is the only test that proves the function is reachable as
// provider::alicloud::arn_build rather than merely correct in Go.
//
// It needs a terraform binary (1.8+) and TF_ACC=1, but no Alibaba Cloud credentials:
//
//	TF_ACC=1 go test ./alicloud/function/ -run TestARNBuildFunction_HCL -v
//
// AT001: the config creates no resources, so CheckDestroy has nothing to verify.
// AT005: this is a function, not a resource, and a TestAcc prefix would make the
// remote acceptance-test runner pick it up. Do not put a space after the comma.
// lintignore: AT001,AT005
func TestARNBuildFunction_HCL(t *testing.T) {
	sdkresource.Test(t, sdkresource.TestCase{
		ProtoV5ProviderFactories: map[string]func() (tfprotov5.ProviderServer, error){
			"alicloud": func() (tfprotov5.ProviderServer, error) {
				muxServer, err := tf5muxserver.NewMuxServer(context.Background(),
					(&sdkschema.Provider{}).GRPCProvider,
					providerserver.NewProtocol5(&harnessProvider{}),
				)
				if err != nil {
					return nil, err
				}
				return muxServer.ProviderServer(), nil
			},
		},
		Steps: []sdkresource.TestStep{
			{
				Config: arnBuildHCLConfig,
				Check: sdkresource.ComposeTestCheckFunc(
					sdkresource.TestCheckOutput("role_arn", "acs:ram::123456789012****:role/example"),
					sdkresource.TestCheckOutput("instance_arn", "acs:ecs:cn-hangzhou:123456789012****:instance/i-bp1234567890abcdef"),
					sdkresource.TestCheckOutput("bucket_arn", "acs:oss:*:*:my-bucket/*"),
					sdkresource.TestCheckOutput("function_arn", "acs:fc:cn-hangzhou:123456789012****:services/foo.LATEST/functions/bar"),
					sdkresource.TestCheckOutput("queue_arn", "acs:mns:cn-hangzhou:123456789012****:/queues/my-queue/messages"),
				),
			},
		},
	})
}
