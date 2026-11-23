package float64planmodifier_test

import (
	"context"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/planmodifiers/float64planmodifier"
	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestFloat64UseStateForUnknown(t *testing.T) {
	t.Parallel()

	// Only State.Raw.IsNull() is inspected from the raw values: a null Raw
	// models "no prior state" (create); a known Raw models an update.
	rawNull := tftypes.NewValue(tftypes.String, nil)
	rawKnown := tftypes.NewValue(tftypes.String, "x")

	testCases := map[string]struct {
		request  planmodifier.Float64Request
		expected planmodifier.Float64Response
	}{
		"create": {
			// null prior state → creating → keep the unknown planned value
			request: planmodifier.Float64Request{
				State:       tfsdk.State{Raw: rawNull},
				Plan:        tfsdk.Plan{Raw: rawKnown},
				ConfigValue: types.Float64Null(),
				StateValue:  types.Float64Null(),
				PlanValue:   types.Float64Unknown(),
			},
			expected: planmodifier.Float64Response{
				PlanValue: types.Float64Unknown(),
			},
		},
		"known-plan": {
			// configured attribute → planned value already known → unchanged
			request: planmodifier.Float64Request{
				State:       tfsdk.State{Raw: rawKnown},
				Plan:        tfsdk.Plan{Raw: rawKnown},
				ConfigValue: types.Float64Value(1.5),
				StateValue:  types.Float64Value(1.5),
				PlanValue:   types.Float64Value(1.5),
			},
			expected: planmodifier.Float64Response{
				PlanValue: types.Float64Value(1.5),
			},
		},
		"unknown-config": {
			// unknown configuration value → do not touch, avoid interpolation issues
			request: planmodifier.Float64Request{
				State:       tfsdk.State{Raw: rawKnown},
				Plan:        tfsdk.Plan{Raw: rawKnown},
				ConfigValue: types.Float64Unknown(),
				StateValue:  types.Float64Value(5.5),
				PlanValue:   types.Float64Unknown(),
			},
			expected: planmodifier.Float64Response{
				PlanValue: types.Float64Unknown(),
			},
		},
		"null-config-unknown-plan": {
			// unconfigured computed attribute → copy the prior state value
			request: planmodifier.Float64Request{
				State:       tfsdk.State{Raw: rawKnown},
				Plan:        tfsdk.Plan{Raw: rawKnown},
				ConfigValue: types.Float64Null(),
				StateValue:  types.Float64Value(5.5),
				PlanValue:   types.Float64Unknown(),
			},
			expected: planmodifier.Float64Response{
				PlanValue: types.Float64Value(5.5),
			},
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			resp := planmodifier.Float64Response{
				PlanValue: testCase.request.PlanValue,
			}

			float64planmodifier.UseStateForUnknown().PlanModifyFloat64(context.Background(), testCase.request, &resp)

			if diff := cmp.Diff(testCase.expected, resp); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}
