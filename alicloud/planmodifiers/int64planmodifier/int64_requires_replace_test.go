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

func TestInt64RequiresReplace(t *testing.T) {
	t.Parallel()

	// Only Raw.IsNull() is inspected by the modifier: a null Raw models "no
	// prior state" (create) or "no plan" (destroy); a known Raw models an update.
	rawNull := tftypes.NewValue(tftypes.String, nil)
	rawKnown := tftypes.NewValue(tftypes.String, "x")

	testCases := map[string]struct {
		request  planmodifier.Int64Request
		expected planmodifier.Int64Response
	}{
		"create": {
			// null prior state → creating → no replacement
			request: planmodifier.Int64Request{
				State:      tfsdk.State{Raw: rawNull},
				Plan:       tfsdk.Plan{Raw: rawKnown},
				StateValue: types.Int64Null(),
				PlanValue:  types.Int64Value(1),
			},
			expected: planmodifier.Int64Response{
				PlanValue:       types.Int64Value(1),
				RequiresReplace: false,
			},
		},
		"destroy": {
			// null plan → destroying → no replacement
			request: planmodifier.Int64Request{
				State:      tfsdk.State{Raw: rawKnown},
				Plan:       tfsdk.Plan{Raw: rawNull},
				StateValue: types.Int64Value(1),
				PlanValue:  types.Int64Null(),
			},
			expected: planmodifier.Int64Response{
				PlanValue:       types.Int64Null(),
				RequiresReplace: false,
			},
		},
		"no-change": {
			// update with equal plan/state → no replacement
			request: planmodifier.Int64Request{
				State:      tfsdk.State{Raw: rawKnown},
				Plan:       tfsdk.Plan{Raw: rawKnown},
				StateValue: types.Int64Value(1),
				PlanValue:  types.Int64Value(1),
			},
			expected: planmodifier.Int64Response{
				PlanValue:       types.Int64Value(1),
				RequiresReplace: false,
			},
		},
		"change": {
			// update with differing plan/state → replacement
			request: planmodifier.Int64Request{
				State:      tfsdk.State{Raw: rawKnown},
				Plan:       tfsdk.Plan{Raw: rawKnown},
				StateValue: types.Int64Value(1),
				PlanValue:  types.Int64Value(2),
			},
			expected: planmodifier.Int64Response{
				PlanValue:       types.Int64Value(2),
				RequiresReplace: true,
			},
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			resp := planmodifier.Int64Response{
				PlanValue: testCase.request.PlanValue,
			}

			int64planmodifier.RequiresReplace().PlanModifyInt64(context.Background(), testCase.request, &resp)

			if diff := cmp.Diff(testCase.expected, resp); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}
