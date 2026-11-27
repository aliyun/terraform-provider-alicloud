package float32planmodifier_test

import (
	"context"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/planmodifiers/float32planmodifier"
	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestFloat32UseStateForUnknown(t *testing.T) {
	t.Parallel()

	// Only State.Raw.IsNull() is inspected from the raw values: a null Raw
	// models "no prior state" (create); a known Raw models an update.
	rawNull := tftypes.NewValue(tftypes.String, nil)
	rawKnown := tftypes.NewValue(tftypes.String, "x")

	testCases := map[string]struct {
		request  planmodifier.Float32Request
		expected planmodifier.Float32Response
	}{
		"create": {
			// null prior state → creating → keep the unknown planned value
			request: planmodifier.Float32Request{
				State:       tfsdk.State{Raw: rawNull},
				Plan:        tfsdk.Plan{Raw: rawKnown},
				ConfigValue: types.Float32Null(),
				StateValue:  types.Float32Null(),
				PlanValue:   types.Float32Unknown(),
			},
			expected: planmodifier.Float32Response{
				PlanValue: types.Float32Unknown(),
			},
		},
		"known-plan": {
			// configured attribute → planned value already known → unchanged
			request: planmodifier.Float32Request{
				State:       tfsdk.State{Raw: rawKnown},
				Plan:        tfsdk.Plan{Raw: rawKnown},
				ConfigValue: types.Float32Value(1.5),
				StateValue:  types.Float32Value(1.5),
				PlanValue:   types.Float32Value(1.5),
			},
			expected: planmodifier.Float32Response{
				PlanValue: types.Float32Value(1.5),
			},
		},
		"unknown-config": {
			// unknown configuration value → do not touch, avoid interpolation issues
			request: planmodifier.Float32Request{
				State:       tfsdk.State{Raw: rawKnown},
				Plan:        tfsdk.Plan{Raw: rawKnown},
				ConfigValue: types.Float32Unknown(),
				StateValue:  types.Float32Value(5.5),
				PlanValue:   types.Float32Unknown(),
			},
			expected: planmodifier.Float32Response{
				PlanValue: types.Float32Unknown(),
			},
		},
		"null-config-unknown-plan": {
			// unconfigured computed attribute → copy the prior state value
			request: planmodifier.Float32Request{
				State:       tfsdk.State{Raw: rawKnown},
				Plan:        tfsdk.Plan{Raw: rawKnown},
				ConfigValue: types.Float32Null(),
				StateValue:  types.Float32Value(5.5),
				PlanValue:   types.Float32Unknown(),
			},
			expected: planmodifier.Float32Response{
				PlanValue: types.Float32Value(5.5),
			},
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			resp := planmodifier.Float32Response{
				PlanValue: testCase.request.PlanValue,
			}

			float32planmodifier.UseStateForUnknown().PlanModifyFloat32(context.Background(), testCase.request, &resp)

			if diff := cmp.Diff(testCase.expected, resp); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}
