package framework_test

import (
	"context"
	"maps"
	"slices"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/framework"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// TestUnitFrameworkProviderFunctions asserts that the provider advertises its
// provider-defined functions over the wire, through the same muxed server main.go
// builds. Each function's own tests exercise its implementation directly, or through a
// stand-in provider, so none of them notices if the registration in
// alicloudProvider.Functions is dropped — this test is the one that does.
func TestUnitFrameworkProviderFunctions(t *testing.T) {
	ctx := context.Background()

	serverFactory, err := framework.ProtoV5ProviderServerFactory(ctx, alicloud.Provider())
	if err != nil {
		t.Fatalf("muxing the SDK v2 and framework providers: %s", err)
	}

	resp, err := serverFactory().GetProviderSchema(ctx, &tfprotov5.GetProviderSchemaRequest{})
	if err != nil {
		t.Fatalf("GetProviderSchema on the muxed provider: %s", err)
	}

	expected := []struct {
		name   string
		params []string
	}{
		{"arn_build", []string{"ram_code", "region", "account_id", "relative_id"}},
	}

	for _, want := range expected {
		name, params := want.name, want.params

		fn, ok := resp.Functions[name]
		if !ok {
			t.Errorf("the muxed provider does not serve function %q; it serves %v",
				name, slices.Sorted(maps.Keys(resp.Functions)))
			continue
		}

		if fn.Summary == "" {
			t.Errorf("function %q: empty Summary", name)
		}

		if fn.Return == nil || !fn.Return.Type.Is(tftypes.String) {
			t.Errorf("function %q: unexpected return: %+v", name, fn.Return)
		}

		if len(fn.Parameters) != len(params) {
			t.Errorf("function %q: got %d parameters, want %d", name, len(fn.Parameters), len(params))
			continue
		}

		// Parameters are positional, so a reordering is a breaking change to every
		// existing call even though each individual name still resolves.
		for i, param := range fn.Parameters {
			if param.Name != params[i] {
				t.Errorf("function %q parameter %d: got %q, want %q", name, i, param.Name, params[i])
			}
			if !param.Type.Is(tftypes.String) {
				t.Errorf("function %q parameter %d (%s): got type %s, want string", name, i, param.Name, param.Type)
			}
		}
	}
}
