package int64planmodifier_test

import (
	"context"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/planmodifiers/int64planmodifier"
	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestInt64UseStateForUnknown(t *testing.T) {
	t.Parallel()

	// Only State.Raw.IsNull() is inspected from the raw values: a null Raw
	// models "no prior state" (create); a known Raw models an update.
	rawNull := tftypes.NewValue(tftypes.String, nil)
	rawKnown := tftypes.NewValue(tftypes.String, "x")

	testCases := map[string]struct {
		request  planmodifier.Int64Request
		expected planmodifier.Int64Response
	}{
		"create": {
			// null prior state → creating → keep the unknown planned value
			request: planmodifier.Int64Request{
				State:       tfsdk.State{Raw: rawNull},
				Plan:        tfsdk.Plan{Raw: rawKnown},
				ConfigValue: types.Int64Null(),
				StateValue:  types.Int64Null(),
				PlanValue:   types.Int64Unknown(),
			},
			expected: planmodifier.Int64Response{
				PlanValue: types.Int64Unknown(),
			},
		},
		"known-plan": {
			// configured attribute → planned value already known → unchanged
			request: planmodifier.Int64Request{
				State:       tfsdk.State{Raw: rawKnown},
				Plan:        tfsdk.Plan{Raw: rawKnown},
				ConfigValue: types.Int64Value(1),
				StateValue:  types.Int64Value(1),
				PlanValue:   types.Int64Value(1),
			},
			expected: planmodifier.Int64Response{
				PlanValue: types.Int64Value(1),
			},
		},
		"unknown-config": {
			// unknown configuration value → do not touch, avoid interpolation issues
			request: planmodifier.Int64Request{
				State:       tfsdk.State{Raw: rawKnown},
				Plan:        tfsdk.Plan{Raw: rawKnown},
				ConfigValue: types.Int64Unknown(),
				StateValue:  types.Int64Value(5),
				PlanValue:   types.Int64Unknown(),
			},
			expected: planmodifier.Int64Response{
				PlanValue: types.Int64Unknown(),
			},
		},
		"null-config-unknown-plan": {
			// unconfigured computed attribute → copy the prior state value
			request: planmodifier.Int64Request{
				State:       tfsdk.State{Raw: rawKnown},
				Plan:        tfsdk.Plan{Raw: rawKnown},
				ConfigValue: types.Int64Null(),
				StateValue:  types.Int64Value(5),
				PlanValue:   types.Int64Unknown(),
			},
			expected: planmodifier.Int64Response{
				PlanValue: types.Int64Value(5),
			},
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			resp := planmodifier.Int64Response{
				PlanValue: testCase.request.PlanValue,
			}

			int64planmodifier.UseStateForUnknown().PlanModifyInt64(context.Background(), testCase.request, &resp)

			if diff := cmp.Diff(testCase.expected, resp); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}
