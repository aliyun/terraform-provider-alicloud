package framework_test

import (
	"context"
	"maps"
	"slices"
	"testing"

	// Importing package alicloud pulls in terraform-plugin-sdk/v2/helper/resource, which
	// registers a -sweep flag in its init just as terraform-plugin-testing/helper/resource
	// does. A test binary linking both dies at startup with "flag redefined: sweep" before
	// any test runs, so no test in this package can use terraform-plugin-testing until
	// package alicloud moves its non-test files off helper/resource onto helper/retry.
	// Build a framework-only factory instead; see protoV5ProviderFactories in
	// alicloud/function/arn_build_test.go.
	"github.com/aliyun/terraform-provider-alicloud/alicloud"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/framework"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// TestUnitFrameworkProviderFunctions asserts that the provider advertises its
// provider-defined functions over the wire, through the same muxed server main.go
// builds. Each function's own unit tests exercise its implementation directly, so none
// of them notices if the registration in alicloudProvider.Functions is dropped. The HCL
// tests in alicloud/function do, but they need a terraform binary on PATH and they serve
// the framework provider on its own rather than the muxed pair. This test needs neither,
// and it is the only one that pins down the wire signature — parameter names, their
// order, and the return type — which a configuration calling the function positionally
// depends on.
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

	// arn_parse returns the four sections of an ARN as an object, keyed by arn_build's
	// parameter names so the two compose.
	arnSections := tftypes.Object{AttributeTypes: map[string]tftypes.Type{
		"service":    tftypes.String,
		"region":     tftypes.String,
		"account_id": tftypes.String,
		"resource":   tftypes.String,
	}}

	expected := []struct {
		name   string
		params []string
		ret    tftypes.Type
	}{
		{"arn_build", []string{"service", "region", "account_id", "resource"}, tftypes.String},
		{"arn_parse", []string{"arn"}, arnSections},
	}

	for _, want := range expected {
		name, params, ret := want.name, want.params, want.ret

		fn, ok := resp.Functions[name]
		if !ok {
			t.Errorf("the muxed provider does not serve function %q; it serves %v",
				name, slices.Sorted(maps.Keys(resp.Functions)))
			continue
		}

		if fn.Summary == "" {
			t.Errorf("function %q: empty Summary", name)
		}

		// Equal, not Is: Is compares only the kind of the type, so an object return
		// would satisfy it whatever its attributes are, and a renamed or dropped
		// attribute — a breaking change to every caller — would go unnoticed.
		switch {
		case fn.Return == nil:
			t.Errorf("function %q: no return, want %s", name, ret)
		case !fn.Return.Type.Equal(ret):
			t.Errorf("function %q: got return type %s, want %s", name, fn.Return.Type, ret)
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
