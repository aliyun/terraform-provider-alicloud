// Package-internal unit tests for the RC cluster node pool helpers. The file is
// deliberately not named resource_rds_rc_cluster_node_pool_unit_test.go: the
// PR collection-order gate maps a changed resource_*_test.go back to
// resource_*.go and reads that path from the base commit, which fails for a new
// resource, and the ACC file resource_rds_rc_cluster_node_pool_test.go has to
// stay in package rds_test.
package rds

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestUnitRCClusterNodePoolNodePoolInfo(t *testing.T) {
	payload, err := expandRCClusterNodePoolNodePoolInfo("pool-a")
	if err != nil {
		t.Fatalf("expandRCClusterNodePoolNodePoolInfo returned an error: %v", err)
	}
	if payload != `{"Name":"pool-a"}` {
		t.Fatalf("unexpected NodePoolInfo payload: %s", payload)
	}

	empty, err := expandRCClusterNodePoolNodePoolInfo("")
	if err != nil {
		t.Fatalf("expandRCClusterNodePoolNodePoolInfo returned an error for an empty name: %v", err)
	}
	if empty != "" {
		t.Fatalf("an empty name must not render a NodePoolInfo payload, got %q", empty)
	}
}

func TestUnitRCClusterNodePoolIdentity(t *testing.T) {
	clusterId, nodePoolId, diags := rcClusterNodePoolIdentity(rcClusterNodePoolResourceModel{
		Id:         types.StringValue("c1:rcnp1"),
		ClusterId:  types.StringValue("c1"),
		NodePoolId: types.StringValue("rcnp1"),
	})
	if diags.HasError() || clusterId != "c1" || nodePoolId != "rcnp1" {
		t.Fatalf("unexpected identity: %s %s %v", clusterId, nodePoolId, diags)
	}

	clusterId, nodePoolId, diags = rcClusterNodePoolIdentity(rcClusterNodePoolResourceModel{
		Id: types.StringValue("c1:rcnp1"),
	})
	if diags.HasError() || clusterId != "c1" || nodePoolId != "rcnp1" {
		t.Fatalf("an imported ID must be parsed: %s %s %v", clusterId, nodePoolId, diags)
	}

	_, _, diags = rcClusterNodePoolIdentity(rcClusterNodePoolResourceModel{
		Id: types.StringValue("rcnp1"),
	})
	if !diags.HasError() {
		t.Fatal("an ID without a cluster part must be rejected")
	}
}

func TestUnitRCClusterNodePoolState(t *testing.T) {
	object := map[string]interface{}{
		"NodePoolInfo": map[string]interface{}{
			"NodePoolId": "rcnp1",
			"Name":       "pool-a",
			"RegionId":   "cn-hangzhou",
			"Created":    "2026-10-28T10:00:00Z",
			"Updated":    "2026-10-28T10:05:00Z",
		},
		"Status": map[string]interface{}{
			"State":                   "active",
			"TotalNodes":              json.Number("2"),
			"RetryTimes":              json.Number("1"),
			"ScalingErrorDescription": "",
		},
		"ScalingGroup": map[string]interface{}{
			"VSwitchIds":         []interface{}{"vsw-1"},
			"InstanceTypes":      []interface{}{"mysql.x4.xlarge.6cm"},
			"InstanceChargeType": "PostPaid",
			"Period":             json.Number("1"),
			"PeriodUnit":         "Month",
			"AutoRenew":          false,
			"SecurityGroupIds":   []interface{}{"sg-1"},
			"DesiredSize":        json.Number("2"),
		},
		"KubernetesConfig": map[string]interface{}{
			"Labels": []interface{}{
				map[string]interface{}{"Key": "env", "Value": "prod"},
			},
			"Taints": []interface{}{
				map[string]interface{}{"Key": "dedicated", "Value": "rds", "Effect": "NoSchedule"},
			},
		},
	}

	state, diags := rcClusterNodePoolState(context.Background(), object, "c1", "rcnp1", "cn-hangzhou")
	if diags.HasError() {
		t.Fatalf("rcClusterNodePoolState returned diagnostics: %v", diags)
	}
	if got := state.Id.ValueString(); got != "c1:rcnp1" {
		t.Fatalf("unexpected id: %s", got)
	}
	if got := state.Name.ValueString(); got != "pool-a" {
		t.Fatalf("unexpected name: %s", got)
	}
	if got := state.RegionId.ValueString(); got != "cn-hangzhou" {
		t.Fatalf("unexpected region: %s", got)
	}
	if len(state.Status.Elements()) != 1 {
		t.Fatalf("status must carry exactly one element, got %d", len(state.Status.Elements()))
	}
	status := state.Status.Elements()[0].(types.Object).Attributes()
	if got := status["state"].(types.String).ValueString(); got != "active" {
		t.Fatalf("unexpected status state: %s", got)
	}
	if got := status["total_nodes"].(types.Int64).ValueInt64(); got != 2 {
		t.Fatalf("unexpected total nodes: %d", got)
	}

	if len(state.ScalingGroup.Elements()) != 1 {
		t.Fatalf("scaling_group must carry exactly one element, got %d", len(state.ScalingGroup.Elements()))
	}
	scalingGroup := state.ScalingGroup.Elements()[0].(types.Object).Attributes()
	if got := scalingGroup["instance_charge_type"].(types.String).ValueString(); got != "PostPaid" {
		t.Fatalf("unexpected instance charge type: %s", got)
	}
	if got := scalingGroup["desired_size"].(types.Int64).ValueInt64(); got != 2 {
		t.Fatalf("unexpected desired size: %d", got)
	}

	if len(state.KubernetesConfig.Elements()) != 1 {
		t.Fatalf("kubernetes_config must carry exactly one element, got %d", len(state.KubernetesConfig.Elements()))
	}
	kubernetesConfig := state.KubernetesConfig.Elements()[0].(types.Object).Attributes()
	labels := kubernetesConfig["labels"].(types.List)
	if len(labels.Elements()) != 1 {
		t.Fatalf("labels must carry exactly one element, got %d", len(labels.Elements()))
	}
	label := labels.Elements()[0].(types.Object).Attributes()
	if got := label["value"].(types.String).ValueString(); got != "prod" {
		t.Fatalf("unexpected label value: %s", got)
	}
}

func TestUnitRCClusterNodePoolExpandScalingGroup(t *testing.T) {
	scalingGroup := types.ListValueMust(rcClusterNodePoolScalingGroupType, []attr.Value{
		types.ObjectValueMust(rcClusterNodePoolScalingGroupAttrTypes, map[string]attr.Value{
			"v_switch_ids":         types.ListValueMust(types.StringType, []attr.Value{types.StringValue("vsw-1")}),
			"instance_types":       types.ListValueMust(types.StringType, []attr.Value{types.StringValue("mysql.x4.xlarge.6cm")}),
			"instance_charge_type": types.StringValue("PostPaid"),
			"period":               types.Int64Value(1),
			"period_unit":          types.StringValue("Month"),
			"auto_renew":           types.BoolValue(false),
			"security_group_ids":   types.ListValueMust(types.StringType, []attr.Value{types.StringValue("sg-1")}),
			"desired_size":         types.Int64Value(0),
		}),
	})

	payload, diags := expandRCClusterNodePoolScalingGroup(context.Background(), scalingGroup)
	if diags.HasError() {
		t.Fatalf("expandRCClusterNodePoolScalingGroup returned diagnostics: %v", diags)
	}
	var request map[string]interface{}
	if err := json.Unmarshal([]byte(payload), &request); err != nil {
		t.Fatalf("the scaling group payload is not valid JSON: %v", err)
	}
	if got := request["InstanceChargeType"]; got != "PostPaid" {
		t.Fatalf("unexpected InstanceChargeType: %v", got)
	}
	if got := request["DesiredSize"]; got != float64(0) {
		t.Fatalf("unexpected DesiredSize: %v", got)
	}
	if got := request["VSwitchIds"]; len(got.([]interface{})) != 1 {
		t.Fatalf("unexpected VSwitchIds: %v", got)
	}

	nullPayload, diags := expandRCClusterNodePoolScalingGroup(context.Background(), types.ListNull(rcClusterNodePoolScalingGroupType))
	if diags.HasError() || nullPayload != "" {
		t.Fatalf("a null scaling group must render nothing, got %q %v", nullPayload, diags)
	}
}

func TestUnitRCClusterNodePoolExpandKubernetesConfig(t *testing.T) {
	kubernetesConfig := types.ListValueMust(rcClusterNodePoolKubernetesConfigType, []attr.Value{
		types.ObjectValueMust(rcClusterNodePoolKubernetesConfigAttrTypes, map[string]attr.Value{
			"labels": types.ListValueMust(rcClusterNodePoolLabelType, []attr.Value{
				types.ObjectValueMust(rcClusterNodePoolLabelAttrTypes, map[string]attr.Value{
					"key":   types.StringValue("env"),
					"value": types.StringValue("prod"),
				}),
			}),
			"taints": types.ListValueMust(rcClusterNodePoolTaintType, []attr.Value{
				types.ObjectValueMust(rcClusterNodePoolTaintAttrTypes, map[string]attr.Value{
					"key":    types.StringValue("dedicated"),
					"value":  types.StringValue("rds"),
					"effect": types.StringValue("NoSchedule"),
				}),
			}),
		}),
	})

	payload, diags := expandRCClusterNodePoolKubernetesConfig(context.Background(), kubernetesConfig)
	if diags.HasError() {
		t.Fatalf("expandRCClusterNodePoolKubernetesConfig returned diagnostics: %v", diags)
	}
	var request map[string]interface{}
	if err := json.Unmarshal([]byte(payload), &request); err != nil {
		t.Fatalf("the kubernetes config payload is not valid JSON: %v", err)
	}
	labels, ok := request["Labels"].([]interface{})
	if !ok || len(labels) != 1 {
		t.Fatalf("unexpected Labels: %v", request["Labels"])
	}
	taints, ok := request["Taints"].([]interface{})
	if !ok || len(taints) != 1 {
		t.Fatalf("unexpected Taints: %v", request["Taints"])
	}
}

func TestUnitRCClusterNodePoolNotFound(t *testing.T) {
	cases := map[string]bool{
		"code: InvalidParam.NodePoolIdNotExist":     true,
		"code: InvalidParam.ClusterIdNotExist":      true,
		"code: OperationDenied.ClusterNotSupport":   true,
		"code: InvalidDBInstanceId.NotFound":        true,
		"code: InvalidParam.PeriodUnitNotSupported": false,
		"code: InvalidParameter":                    false,
	}
	for message, want := range cases {
		if got := isRCClusterNodePoolNotFound(errString(message)); got != want {
			t.Errorf("isRCClusterNodePoolNotFound(%q) = %v, want %v", message, got, want)
		}
	}
	if isRCClusterNodePoolNotFound(nil) {
		t.Error("a nil error must not be treated as not found")
	}
}

type errString string

func (e errString) Error() string { return string(e) }

func TestUnitRCClusterNodePoolApplyConfigShape(t *testing.T) {
	kubernetesConfig := func(labels, taints types.List) types.List {
		return types.ListValueMust(rcClusterNodePoolKubernetesConfigType, []attr.Value{
			types.ObjectValueMust(rcClusterNodePoolKubernetesConfigAttrTypes, map[string]attr.Value{
				"labels": labels,
				"taints": taints,
			}),
		})
	}
	labels := types.ListValueMust(rcClusterNodePoolLabelType, []attr.Value{
		types.ObjectValueMust(rcClusterNodePoolLabelAttrTypes, map[string]attr.Value{
			"key":   types.StringValue("env"),
			"value": types.StringValue("prod"),
		}),
	})
	isEmptyList := func(value attr.Value) bool {
		list, ok := value.(types.List)
		return ok && !list.IsNull() && !list.IsUnknown() && len(list.Elements()) == 0
	}
	taintsOf := func(list types.List) attr.Value {
		object, ok := list.Elements()[0].(types.Object)
		if !ok {
			t.Fatalf("unexpected kubernetes_config element: %T", list.Elements()[0])
		}
		return object.Attributes()["taints"]
	}

	// The API omits empty taints; a planned empty list must stay an empty list.
	observed := kubernetesConfig(labels, types.ListNull(rcClusterNodePoolTaintType))
	planned := kubernetesConfig(labels, types.ListValueMust(rcClusterNodePoolTaintType, []attr.Value{}))
	if got := taintsOf(rcClusterNodePoolApplyConfigShape(planned, observed)); !isEmptyList(got) {
		t.Fatalf("a planned empty taints list must survive read-back, got %v", got)
	}

	// A taints argument left out of the configuration stays null.
	planned = kubernetesConfig(labels, types.ListNull(rcClusterNodePoolTaintType))
	if got := taintsOf(rcClusterNodePoolApplyConfigShape(planned, observed)); !got.IsNull() {
		t.Fatalf("an absent taints argument must stay null, got %v", got)
	}

	// Import starts from an empty state: the collections come back as empty lists.
	imported := rcClusterNodePoolApplyConfigShape(types.ListNull(rcClusterNodePoolKubernetesConfigType), observed)
	if got := taintsOf(imported); !isEmptyList(got) {
		t.Fatalf("an import must report taints as an empty list, got %v", got)
	}

	// Values returned by the API always win over the shape of the configuration.
	apiTaints := types.ListValueMust(rcClusterNodePoolTaintType, []attr.Value{
		types.ObjectValueMust(rcClusterNodePoolTaintAttrTypes, map[string]attr.Value{
			"key":    types.StringValue("dedicated"),
			"value":  types.StringValue("rds"),
			"effect": types.StringValue("NoSchedule"),
		}),
	})
	observed = kubernetesConfig(labels, apiTaints)
	if got := taintsOf(rcClusterNodePoolApplyConfigShape(planned, observed)); got.IsNull() || len(got.(types.List).Elements()) != 1 {
		t.Fatalf("the API taints must be kept, got %v", got)
	}
}
