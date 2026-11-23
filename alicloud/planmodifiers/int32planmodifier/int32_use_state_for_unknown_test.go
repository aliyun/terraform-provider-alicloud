package int32planmodifier_test

import (
	"context"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/planmodifiers/int32planmodifier"
	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestInt32UseStateForUnknown(t *testing.T) {
	t.Parallel()

	// Only State.Raw.IsNull() is inspected from the raw values: a null Raw
	// models "no prior state" (create); a known Raw models an update.
	rawNull := tftypes.NewValue(tftypes.String, nil)
	rawKnown := tftypes.NewValue(tftypes.String, "x")

	testCases := map[string]struct {
		request  planmodifier.Int32Request
		expected planmodifier.Int32Response
	}{
		"create": {
			// null prior state → creating → keep the unknown planned value
			request: planmodifier.Int32Request{
				State:       tfsdk.State{Raw: rawNull},
				Plan:        tfsdk.Plan{Raw: rawKnown},
				ConfigValue: types.Int32Null(),
				StateValue:  types.Int32Null(),
				PlanValue:   types.Int32Unknown(),
			},
			expected: planmodifier.Int32Response{
				PlanValue: types.Int32Unknown(),
			},
		},
		"known-plan": {
			// configured attribute → planned value already known → unchanged
			request: planmodifier.Int32Request{
				State:       tfsdk.State{Raw: rawKnown},
				Plan:        tfsdk.Plan{Raw: rawKnown},
				ConfigValue: types.Int32Value(1),
				StateValue:  types.Int32Value(1),
				PlanValue:   types.Int32Value(1),
			},
			expected: planmodifier.Int32Response{
				PlanValue: types.Int32Value(1),
			},
		},
		"unknown-config": {
			// unknown configuration value → do not touch, avoid interpolation issues
			request: planmodifier.Int32Request{
				State:       tfsdk.State{Raw: rawKnown},
				Plan:        tfsdk.Plan{Raw: rawKnown},
				ConfigValue: types.Int32Unknown(),
				StateValue:  types.Int32Value(5),
				PlanValue:   types.Int32Unknown(),
			},
			expected: planmodifier.Int32Response{
				PlanValue: types.Int32Unknown(),
			},
		},
		"null-config-unknown-plan": {
			// unconfigured computed attribute → copy the prior state value
			request: planmodifier.Int32Request{
				State:       tfsdk.State{Raw: rawKnown},
				Plan:        tfsdk.Plan{Raw: rawKnown},
				ConfigValue: types.Int32Null(),
				StateValue:  types.Int32Value(5),
				PlanValue:   types.Int32Unknown(),
			},
			expected: planmodifier.Int32Response{
				PlanValue: types.Int32Value(5),
			},
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			resp := planmodifier.Int32Response{
				PlanValue: testCase.request.PlanValue,
			}

			int32planmodifier.UseStateForUnknown().PlanModifyInt32(context.Background(), testCase.request, &resp)

			if diff := cmp.Diff(testCase.expected, resp); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}
