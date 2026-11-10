package cs

import (
	"context"
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	unitAutoRepairPolicyType     = cty.Object(map[string]cty.Type{"restart_node": cty.Bool, "auto_repair_policy_id": cty.String})
	unitAutoRepairPolicyListType = cty.List(unitAutoRepairPolicyType)
	unitManagementType           = cty.Object(map[string]cty.Type{"auto_repair_policy": unitAutoRepairPolicyListType})
	unitManagementListType       = cty.List(unitManagementType)
)

func unitNodePoolRawConfig(management cty.Value) cty.Value {
	return cty.ObjectVal(map[string]cty.Value{"management": management})
}

func unitManagementBlocks(policy cty.Value) cty.Value {
	return cty.ListVal([]cty.Value{cty.ObjectVal(map[string]cty.Value{"auto_repair_policy": policy})})
}

func unitAutoRepairPolicyBlocks(restartNode, policyID cty.Value) cty.Value {
	return cty.ListVal([]cty.Value{cty.ObjectVal(map[string]cty.Value{
		"restart_node":          restartNode,
		"auto_repair_policy_id": policyID,
	})})
}

// TestUnitValidateNodePoolAutoRepairPolicyConflict covers the raw-config warning emitted when
// the deprecated restart_node switch is configured together with auto_repair_policy_id.
// Only an explicitly configured restart_node can be confused with the value read back
// from the service, so absent, empty and unknown values must stay silent.
func TestUnitValidateNodePoolAutoRepairPolicyConflict(t *testing.T) {
	tests := []struct {
		name     string
		raw      cty.Value
		wantWarn int
	}{
		{
			name: "nil raw config",
			raw:  cty.NilVal,
		},
		{
			name: "management not configured",
			raw:  unitNodePoolRawConfig(cty.NullVal(unitManagementListType)),
		},
		{
			name: "management without auto_repair_policy",
			raw:  unitNodePoolRawConfig(unitManagementBlocks(cty.NullVal(unitAutoRepairPolicyListType))),
		},
		{
			name: "management is unknown",
			raw:  unitNodePoolRawConfig(cty.UnknownVal(unitManagementListType)),
		},
		{
			name: "only auto_repair_policy_id",
			raw:  unitNodePoolRawConfig(unitManagementBlocks(unitAutoRepairPolicyBlocks(cty.NullVal(cty.Bool), cty.StringVal(builtinRepairPolicyRestartNode)))),
		},
		{
			name: "only restart_node",
			raw:  unitNodePoolRawConfig(unitManagementBlocks(unitAutoRepairPolicyBlocks(cty.True, cty.NullVal(cty.String)))),
		},
		{
			name: "empty auto_repair_policy_id with restart_node",
			raw:  unitNodePoolRawConfig(unitManagementBlocks(unitAutoRepairPolicyBlocks(cty.False, cty.StringVal("")))),
		},
		{
			name: "unknown auto_repair_policy_id with restart_node",
			raw:  unitNodePoolRawConfig(unitManagementBlocks(unitAutoRepairPolicyBlocks(cty.True, cty.UnknownVal(cty.String)))),
		},
		{
			name: "unknown restart_node with auto_repair_policy_id",
			raw:  unitNodePoolRawConfig(unitManagementBlocks(unitAutoRepairPolicyBlocks(cty.UnknownVal(cty.Bool), cty.StringVal(builtinRepairPolicyRestartNode)))),
		},
		{
			name:     "restart_node true with auto_repair_policy_id",
			raw:      unitNodePoolRawConfig(unitManagementBlocks(unitAutoRepairPolicyBlocks(cty.True, cty.StringVal(builtinRepairPolicyRestartNode)))),
			wantWarn: 1,
		},
		{
			name:     "restart_node false with auto_repair_policy_id",
			raw:      unitNodePoolRawConfig(unitManagementBlocks(unitAutoRepairPolicyBlocks(cty.False, cty.StringVal(builtinRepairPolicyWithoutReboots)))),
			wantWarn: 1,
		},
		{
			name:     "restart_node with custom policy id",
			raw:      unitNodePoolRawConfig(unitManagementBlocks(unitAutoRepairPolicyBlocks(cty.True, cty.StringVal("rp-custom")))),
			wantWarn: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := &schema.ValidateResourceConfigFuncResponse{}
			validateNodePoolAutoRepairPolicyConflict()(context.Background(), schema.ValidateResourceConfigFuncRequest{RawConfig: tt.raw}, resp)

			require.Len(t, resp.Diagnostics, tt.wantWarn)
			for _, diagnostic := range resp.Diagnostics {
				assert.Equal(t, diag.Warning, diagnostic.Severity)
				assert.Equal(t, "auto_repair_policy_id takes precedence over restart_node", diagnostic.Summary)
				assert.Contains(t, diagnostic.Detail, builtinRepairPolicyRestartNode)
				assert.Contains(t, diagnostic.Detail, builtinRepairPolicyWithoutReboots)
			}
		})
	}
}

// TestUnitAutoRepairPolicyIDConfigured covers the gate that decides whether the update
// request may carry auto_repair_policy_id. Only a value written in the configuration may
// be sent: the id read back from state can be one the service derived from restart_node,
// and persisting it would silently disable that switch.
func TestUnitAutoRepairPolicyIDConfigured(t *testing.T) {
	tests := []struct {
		name string
		raw  cty.Value
		want bool
	}{
		{
			name: "nil raw config",
			raw:  cty.NilVal,
		},
		{
			name: "management not configured",
			raw:  unitNodePoolRawConfig(cty.NullVal(unitManagementListType)),
		},
		{
			name: "management is unknown",
			raw:  unitNodePoolRawConfig(cty.UnknownVal(unitManagementListType)),
		},
		{
			name: "management without auto_repair_policy",
			raw:  unitNodePoolRawConfig(unitManagementBlocks(cty.NullVal(unitAutoRepairPolicyListType))),
		},
		{
			name: "only restart_node",
			raw:  unitNodePoolRawConfig(unitManagementBlocks(unitAutoRepairPolicyBlocks(cty.True, cty.NullVal(cty.String)))),
		},
		{
			name: "empty auto_repair_policy_id",
			raw:  unitNodePoolRawConfig(unitManagementBlocks(unitAutoRepairPolicyBlocks(cty.True, cty.StringVal("")))),
		},
		{
			name: "explicit builtin policy id",
			raw:  unitNodePoolRawConfig(unitManagementBlocks(unitAutoRepairPolicyBlocks(cty.NullVal(cty.Bool), cty.StringVal(builtinRepairPolicyRestartNode)))),
			want: true,
		},
		{
			name: "explicit custom policy id",
			raw:  unitNodePoolRawConfig(unitManagementBlocks(unitAutoRepairPolicyBlocks(cty.False, cty.StringVal("rp-custom")))),
			want: true,
		},
		{
			name: "unknown auto_repair_policy_id",
			raw:  unitNodePoolRawConfig(unitManagementBlocks(unitAutoRepairPolicyBlocks(cty.NullVal(cty.Bool), cty.UnknownVal(cty.String)))),
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, autoRepairPolicyIDConfigured(tt.raw))
		})
	}
}

// TestUnitAutoRepairPolicyUpdatePayload reproduces the update request built after a
// refresh stored the service-derived builtin policy id in state: the id must not be
// sent, otherwise the service persists it and restart_node loses effect on the pool.
// Other management changes, such as a new surge value, are enough to reach this code.
func TestUnitAutoRepairPolicyUpdatePayload(t *testing.T) {
	resourceSchema := ResourceAckNodepool().Schema
	nilConfig := cty.NullVal(cty.Object(map[string]cty.Type{"management": unitManagementListType}))

	state := func() *terraform.InstanceState {
		return &terraform.InstanceState{
			ID: "c-fake:np-test",
			Attributes: map[string]string{
				"id":                                "c-fake:np-test",
				"management.#":                      "1",
				"management.0.enable":               "true",
				"management.0.auto_repair":          "true",
				"management.0.surge":                "0",
				"management.0.auto_repair_policy.#": "1",
				"management.0.auto_repair_policy.0.restart_node":          "true",
				"management.0.auto_repair_policy.0.auto_repair_policy_id": builtinRepairPolicyRestartNode,
			},
		}
	}
	diff := func(raw cty.Value, policyDiff *terraform.ResourceAttrDiff) *terraform.InstanceDiff {
		attributes := map[string]*terraform.ResourceAttrDiff{
			"management.0.surge": {Old: "0", New: "1"},
		}
		if policyDiff != nil {
			attributes["management.0.auto_repair_policy.0.auto_repair_policy_id"] = policyDiff
		}
		return &terraform.InstanceDiff{Attributes: attributes, RawConfig: raw}
	}

	tests := []struct {
		name       string
		raw        cty.Value
		policyDiff *terraform.ResourceAttrDiff
		wantMerged string
		wantSent   string
	}{
		{
			name:       "id only injected into state",
			raw:        unitNodePoolRawConfig(unitManagementBlocks(unitAutoRepairPolicyBlocks(cty.True, cty.NullVal(cty.String)))),
			wantMerged: builtinRepairPolicyRestartNode,
		},
		{
			name:       "no raw config",
			raw:        nilConfig,
			wantMerged: builtinRepairPolicyRestartNode,
		},
		{
			name:       "empty id in configuration",
			raw:        unitNodePoolRawConfig(unitManagementBlocks(unitAutoRepairPolicyBlocks(cty.True, cty.StringVal("")))),
			wantMerged: builtinRepairPolicyRestartNode,
		},
		{
			name:       "explicit id in configuration",
			raw:        unitNodePoolRawConfig(unitManagementBlocks(unitAutoRepairPolicyBlocks(cty.True, cty.StringVal("rp-custom")))),
			policyDiff: &terraform.ResourceAttrDiff{Old: builtinRepairPolicyRestartNode, New: "rp-custom"},
			wantMerged: "rp-custom",
			wantSent:   "rp-custom",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := schema.InternalMap(resourceSchema).Data(state(), diff(tt.raw, tt.policyDiff))
			require.NoError(t, err)

			policy := data.Get("management").([]interface{})[0].(map[string]interface{})["auto_repair_policy"].([]interface{})[0].(map[string]interface{})
			require.Equal(t, tt.wantMerged, policy["auto_repair_policy_id"])

			payload := autoRepairPolicyUpdatePayload(data)
			if tt.wantSent == "" {
				assert.NotContains(t, payload, "auto_repair_policy_id")
			} else {
				assert.Equal(t, tt.wantSent, payload["auto_repair_policy_id"])
			}
			assert.Equal(t, true, payload["restart_node"])
		})
	}
}

func TestUnitNodePoolAutoRepairPolicySchema(t *testing.T) {
	resource := ResourceAckNodepool()
	require.Len(t, resource.ValidateRawResourceConfigFuncs, 1)

	management := resource.Schema["management"].Elem.(*schema.Resource).Schema
	policy := management["auto_repair_policy"].Elem.(*schema.Resource).Schema
	assert.Contains(t, policy["restart_node"].Deprecated, "2.3.0")
	assert.Contains(t, policy["restart_node"].Deprecated, builtinRepairPolicyRestartNode)
	assert.Contains(t, policy["restart_node"].Deprecated, builtinRepairPolicyWithoutReboots)
}
