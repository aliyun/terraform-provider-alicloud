package cs

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	builtinRepairPolicyRestartNode     = "r-repair-sys-comp-failure"
	builtinRepairPolicyWithoutReboots  = "r-repair-sys-comp-failure-without-reboots"
	autoRepairPolicyRestartNodeAddress = "management.auto_repair_policy.restart_node"
	autoRepairPolicyIDAddress          = "management.auto_repair_policy.auto_repair_policy_id"
)

// validateNodePoolAutoRepairPolicyConflict warns when the deprecated restart_node
// switch is configured together with auto_repair_policy_id. Both select the same
// builtin policy slot and the service lets auto_repair_policy_id win, so restart_node
// would be silently ignored. The check runs on the raw configuration because
// restart_node is Optional+Computed and the resolved value cannot tell an explicit
// setting apart from the value read back from the service.
func validateNodePoolAutoRepairPolicyConflict() schema.ValidateRawResourceConfigFunc {
	return func(_ context.Context, req schema.ValidateResourceConfigFuncRequest, resp *schema.ValidateResourceConfigFuncResponse) {
		management, ok := rawConfigBlock(req.RawConfig, "management")
		if !ok {
			return
		}
		policy, ok := rawConfigBlock(management, "auto_repair_policy")
		if !ok {
			return
		}
		policyID, _ := rawConfigString(policy, "auto_repair_policy_id")
		if policyID == "" {
			return
		}
		restartNode, ok := rawConfigAttribute(policy, "restart_node")
		if !ok || restartNode.IsNull() || !restartNode.IsKnown() {
			return
		}
		resp.Diagnostics = append(resp.Diagnostics, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "auto_repair_policy_id takes precedence over restart_node",
			Detail: fmt.Sprintf("%s is deprecated and ignored because %s is set to %q. Set auto_repair_policy_id to %q instead of restart_node = true, or to %q instead of restart_node = false, and remove restart_node from the configuration.",
				autoRepairPolicyRestartNodeAddress, autoRepairPolicyIDAddress, policyID,
				builtinRepairPolicyRestartNode, builtinRepairPolicyWithoutReboots),
			AttributePath: cty.GetAttrPath("management").IndexInt(0).GetAttr("auto_repair_policy").IndexInt(0).GetAttr("restart_node"),
		})
	}
}

// rawConfigBlock returns the single element of a MaxItems: 1 block from the raw cty
// configuration. It reports false when the block is absent, null or not yet known.
func rawConfigBlock(parent cty.Value, name string) (cty.Value, bool) {
	blocks, ok := rawConfigAttribute(parent, name)
	if !ok || blocks.IsNull() || !blocks.IsKnown() || !blocks.Type().IsListType() || blocks.LengthInt() == 0 {
		return cty.NilVal, false
	}
	return blocks.Index(cty.NumberIntVal(0)), true
}

// rawConfigAttribute returns an attribute of a raw cty object. It reports false when
// the attribute is missing or the value carries no object type, so callers never
// panic on unexpected input.
func rawConfigAttribute(parent cty.Value, name string) (cty.Value, bool) {
	if parent.IsNull() || !parent.IsKnown() || !parent.Type().IsObjectType() || !parent.Type().HasAttribute(name) {
		return cty.NilVal, false
	}
	return parent.GetAttr(name), true
}

func rawConfigString(parent cty.Value, name string) (string, bool) {
	value, ok := rawConfigAttribute(parent, name)
	if !ok || value.IsNull() || !value.IsKnown() || value.Type() != cty.String {
		return "", false
	}
	return value.AsString(), true
}

// autoRepairPolicyIDConfigured reports whether auto_repair_policy_id is present in the
// raw configuration. The service derives a builtin policy id from the deprecated
// restart_node switch for pools without an explicit binding and only injects it into
// Describe/List responses, but the value still lands in state. Echoing it back on the
// next update would persist the inferred binding, after which restart_node stops having
// any effect on the pool, so request builders must consult the raw configuration instead
// of the merged ResourceData value.
func autoRepairPolicyIDConfigured(rawConfig cty.Value) bool {
	management, ok := rawConfigBlock(rawConfig, "management")
	if !ok {
		return false
	}
	policy, ok := rawConfigBlock(management, "auto_repair_policy")
	if !ok {
		return false
	}
	value, ok := rawConfigAttribute(policy, "auto_repair_policy_id")
	if !ok || value.IsNull() {
		return false
	}
	if !value.IsKnown() {
		// Unknown at plan time, e.g. an id produced by another resource: keep
		// sending the resolved value rather than dropping the binding.
		return true
	}
	return value.Type() == cty.String && value.AsString() != ""
}
