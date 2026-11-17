package provider_test

import (
	"context"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/sdkv2"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
)

// TestUnitFrameworkProviderSchemaMatchesSDK muxes the two provider servers exactly as
// main.go does and asks the result for its schema.
func TestUnitFrameworkProviderSchemaMatchesSDK(t *testing.T) {
	ctx := context.Background()

	serverFactory, err := provider.ProtoV5ProviderServerFactory(ctx, sdkv2.Provider())
	if err != nil {
		t.Fatalf("muxing the SDK v2 and framework providers: %s", err)
	}

	resp, err := serverFactory().GetProviderSchema(ctx, &tfprotov5.GetProviderSchemaRequest{})
	if err != nil {
		t.Fatalf("GetProviderSchema on the muxed provider: %s", err)
	}

	for _, d := range resp.Diagnostics {
		if d.Severity == tfprotov5.DiagnosticSeverityError {
			t.Errorf("GetProviderSchema on the muxed provider: %s: %s", d.Summary, d.Detail)
		}
	}
}
