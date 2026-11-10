package cs

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestUnitFlattenContainerdConfig covers the API-response-to-schema conversion of the
// containerd_config block: only keys explicitly returned by the API are written into the
// element map (missing key = "not set"), numeric fields accept both json.Number and the
// float64 fallback and are rendered as decimal strings for the nullable TypeString fields,
// and an empty/nil response yields an empty list.
func TestUnitFlattenContainerdConfig(t *testing.T) {
	tests := []struct {
		name string
		raw  interface{}
		want []map[string]interface{}
	}{
		{
			name: "nil input returns empty list",
			raw:  nil,
			want: []map[string]interface{}{},
		},
		{
			name: "empty map returns empty list",
			raw:  map[string]interface{}{},
			want: []map[string]interface{}{},
		},
		{
			name: "non-map input returns empty list",
			raw:  "not-a-map",
			want: []map[string]interface{}{},
		},
		{
			name: "partial response only writes returned keys",
			raw: map[string]interface{}{
				"maxConcurrentDownloads": json.Number("3"),
			},
			want: []map[string]interface{}{{
				"max_concurrent_downloads": 3,
			}},
		},
		{
			name: "json.Number values render as decimal strings",
			raw: map[string]interface{}{
				"maxConcurrentDownloads": json.Number("5"),
				"limitCore":              json.Number("1"),
				"limitNoFile":            json.Number("65536"),
				"limitMemLock":           json.Number("128849018880"),
			},
			want: []map[string]interface{}{{
				"max_concurrent_downloads": 5,
				"limit_core":               "1",
				"limit_no_file":            "65536",
				"limit_mem_lock":           "128849018880",
			}},
		},
		{
			name: "float64 fallback values render as decimal strings",
			raw: map[string]interface{}{
				"maxConcurrentDownloads": float64(10),
				"limitCore":              float64(2),
				"limitNoFile":            float64(2048),
				"limitMemLock":           float64(65536),
			},
			want: []map[string]interface{}{{
				"max_concurrent_downloads": 10,
				"limit_core":               "2",
				"limit_no_file":            "2048",
				"limit_mem_lock":           "65536",
			}},
		},
		{
			name: "bool true renders as string true",
			raw: map[string]interface{}{
				"ignoreImageDefinedVolume": true,
			},
			want: []map[string]interface{}{{
				"ignore_image_defined_volume": "true",
			}},
		},
		{
			name: "bool false renders as string false",
			raw: map[string]interface{}{
				"ignoreImageDefinedVolume": false,
			},
			want: []map[string]interface{}{{
				"ignore_image_defined_volume": "false",
			}},
		},
		{
			name: "large json.Number is not truncated",
			raw: map[string]interface{}{
				"limitNoFile": json.Number("9007199254740991"),
			},
			want: []map[string]interface{}{{
				"limit_no_file": "9007199254740991",
			}},
		},
		{
			name: "non-empty registry lists are written",
			raw: map[string]interface{}{
				"registryMirrors":    []interface{}{"https://mirror.example.com"},
				"insecureRegistries": []interface{}{"registry.example.com:5000"},
			},
			want: []map[string]interface{}{{
				"registry_mirrors":    []interface{}{"https://mirror.example.com"},
				"insecure_registries": []interface{}{"registry.example.com:5000"},
			}},
		},
		{
			name: "empty or nil registry lists are skipped",
			raw: map[string]interface{}{
				"limitCore":          json.Number("1"),
				"registryMirrors":    []interface{}{},
				"insecureRegistries": nil,
			},
			want: []map[string]interface{}{{
				"limit_core": "1",
			}},
		},
		{
			name: "explicit zero values stay distinct from missing keys",
			raw: map[string]interface{}{
				"limitCore":                json.Number("0"),
				"ignoreImageDefinedVolume": false,
			},
			want: []map[string]interface{}{{
				"limit_core":                  "0",
				"ignore_image_defined_volume": "false",
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, FlattenContainerdConfig(tt.raw))
		})
	}
}

// TestUnitExpandContainerdConfig covers the schema-to-API-request conversion of the
// containerd_config block using the resource's real schema via schema.TestResourceDataRaw.
// The nullable string fields (limit_core/limit_no_file/limit_mem_lock/
// ignore_image_defined_volume) skip "" ("not set") and send explicit values, including
// "0"/"false"; removing the whole block yields an empty map ({} clears the cloud-side
// configuration under the API's full-replacement semantics). Invalid values bypassing the
// schema ValidateFunc layer surface as errors from strconv.ParseInt/ParseBool.
func TestUnitExpandContainerdConfig(t *testing.T) {
	resourceSchema := ResourceAckNodepool().Schema

	tests := []struct {
		name    string
		raw     map[string]interface{}
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name: "block absent returns empty map",
			raw:  map[string]interface{}{},
			want: map[string]interface{}{},
		},
		{
			name: "block removed (empty list) returns empty map",
			raw: map[string]interface{}{
				"containerd_config": []interface{}{},
			},
			want: map[string]interface{}{},
		},
		{
			name: "all fields unset sends no keys",
			raw: map[string]interface{}{
				"containerd_config": []interface{}{map[string]interface{}{
					"max_concurrent_downloads":    0,
					"limit_core":                  "",
					"limit_no_file":               "",
					"limit_mem_lock":              "",
					"ignore_image_defined_volume": "",
					"registry_mirrors":            []interface{}{},
					"insecure_registries":         []interface{}{},
				}},
			},
			want: map[string]interface{}{},
		},
		{
			name: "explicit zero and false are sent",
			raw: map[string]interface{}{
				"containerd_config": []interface{}{map[string]interface{}{
					"limit_core":                  "0",
					"ignore_image_defined_volume": "false",
				}},
			},
			want: map[string]interface{}{
				"limitCore":                int64(0),
				"ignoreImageDefinedVolume": false,
			},
		},
		{
			name: "normal values convert to int64, bool and lists",
			raw: map[string]interface{}{
				"containerd_config": []interface{}{map[string]interface{}{
					"max_concurrent_downloads":    5,
					"limit_core":                  "1",
					"limit_no_file":               "65536",
					"limit_mem_lock":              "128849018880",
					"ignore_image_defined_volume": "true",
					"registry_mirrors":            []interface{}{"https://mirror.example.com"},
					"insecure_registries":         []interface{}{"registry.example.com:5000"},
				}},
			},
			want: map[string]interface{}{
				"maxConcurrentDownloads":   int64(5),
				"limitCore":                int64(1),
				"limitNoFile":              int64(65536),
				"limitMemLock":             int64(128849018880),
				"ignoreImageDefinedVolume": true,
				"registryMirrors":          []interface{}{"https://mirror.example.com"},
				"insecureRegistries":       []interface{}{"registry.example.com:5000"},
			},
		},
		{
			name: "empty registry_mirrors list omits registryMirrors key",
			raw: map[string]interface{}{
				"containerd_config": []interface{}{map[string]interface{}{
					"limit_core":       "1",
					"registry_mirrors": []interface{}{},
				}},
			},
			want: map[string]interface{}{
				"limitCore": int64(1),
			},
		},
		{
			name: "invalid limit_core returns error",
			raw: map[string]interface{}{
				"containerd_config": []interface{}{map[string]interface{}{
					"limit_core": "abc",
				}},
			},
			wantErr: true,
		},
		{
			name: "invalid ignore_image_defined_volume returns error",
			raw: map[string]interface{}{
				"containerd_config": []interface{}{map[string]interface{}{
					"ignore_image_defined_volume": "maybe",
				}},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, resourceSchema, tt.raw)
			got, err := ExpandContainerdConfig(d)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

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

// TestUnitAutoRepairPolicyUpdatePayloadRateLimits covers the empty-value semantics of the
// two rate limit fields: an empty value is treated as unset and must not reach the request.
func TestUnitAutoRepairPolicyUpdatePayloadRateLimits(t *testing.T) {
	resourceSchema := ResourceAckNodepool().Schema
	nilConfig := cty.NullVal(cty.Object(map[string]cty.Type{"management": unitManagementListType}))

	tests := []struct {
		name          string
		parallel      string
		threshold     string
		wantParallel  string
		wantThreshold string
	}{
		{
			name:          "values from state are sent",
			parallel:      "5",
			threshold:     "50%",
			wantParallel:  "5",
			wantThreshold: "50%",
		},
		{
			name: "empty values are not sent",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state := &terraform.InstanceState{
				ID: "c-fake:np-test",
				Attributes: map[string]string{
					"id":                                "c-fake:np-test",
					"management.#":                      "1",
					"management.0.enable":               "true",
					"management.0.auto_repair":          "true",
					"management.0.surge":                "0",
					"management.0.auto_repair_policy.#": "1",
					"management.0.auto_repair_policy.0.max_parallel_repairing_nodes":  tt.parallel,
					"management.0.auto_repair_policy.0.max_unhealthy_nodes_threshold": tt.threshold,
				},
			}
			diff := &terraform.InstanceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"management.0.surge": {Old: "0", New: "1"},
				},
				RawConfig: nilConfig,
			}
			data, err := schema.InternalMap(resourceSchema).Data(state, diff)
			require.NoError(t, err)

			payload := autoRepairPolicyUpdatePayload(data)
			if tt.wantParallel == "" {
				assert.NotContains(t, payload, "max_parallel_repairing_nodes")
			} else {
				assert.Equal(t, tt.wantParallel, payload["max_parallel_repairing_nodes"])
			}
			if tt.wantThreshold == "" {
				assert.NotContains(t, payload, "max_unhealthy_nodes_threshold")
			} else {
				assert.Equal(t, tt.wantThreshold, payload["max_unhealthy_nodes_threshold"])
			}
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
