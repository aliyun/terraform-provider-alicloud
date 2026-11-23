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

func TestFloat64RequiresReplace(t *testing.T) {
	t.Parallel()

	// Only Raw.IsNull() is inspected by the modifier: a null Raw models "no
	// prior state" (create) or "no plan" (destroy); a known Raw models an update.
	rawNull := tftypes.NewValue(tftypes.String, nil)
	rawKnown := tftypes.NewValue(tftypes.String, "x")

	testCases := map[string]struct {
		request  planmodifier.Float64Request
		expected planmodifier.Float64Response
	}{
		"create": {
			// null prior state → creating → no replacement
			request: planmodifier.Float64Request{
				State:      tfsdk.State{Raw: rawNull},
				Plan:       tfsdk.Plan{Raw: rawKnown},
				StateValue: types.Float64Null(),
				PlanValue:  types.Float64Value(1),
			},
			expected: planmodifier.Float64Response{
				PlanValue:       types.Float64Value(1),
				RequiresReplace: false,
			},
		},
		"destroy": {
			// null plan → destroying → no replacement
			request: planmodifier.Float64Request{
				State:      tfsdk.State{Raw: rawKnown},
				Plan:       tfsdk.Plan{Raw: rawNull},
				StateValue: types.Float64Value(1),
				PlanValue:  types.Float64Null(),
			},
			expected: planmodifier.Float64Response{
				PlanValue:       types.Float64Null(),
				RequiresReplace: false,
			},
		},
		"no-change": {
			// update with equal plan/state → no replacement
			request: planmodifier.Float64Request{
				State:      tfsdk.State{Raw: rawKnown},
				Plan:       tfsdk.Plan{Raw: rawKnown},
				StateValue: types.Float64Value(1),
				PlanValue:  types.Float64Value(1),
			},
			expected: planmodifier.Float64Response{
				PlanValue:       types.Float64Value(1),
				RequiresReplace: false,
			},
		},
		"change": {
			// update with differing plan/state → replacement
			request: planmodifier.Float64Request{
				State:      tfsdk.State{Raw: rawKnown},
				Plan:       tfsdk.Plan{Raw: rawKnown},
				StateValue: types.Float64Value(1),
				PlanValue:  types.Float64Value(2),
			},
			expected: planmodifier.Float64Response{
				PlanValue:       types.Float64Value(2),
				RequiresReplace: true,
			},
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			resp := planmodifier.Float64Response{
				PlanValue: testCase.request.PlanValue,
			}

			float64planmodifier.RequiresReplace().PlanModifyFloat64(context.Background(), testCase.request, &resp)

			if diff := cmp.Diff(testCase.expected, resp); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}
