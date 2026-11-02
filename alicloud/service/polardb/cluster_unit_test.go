package polardb

import (
	"testing"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestUnitPolarDBDistributedNodeSchema(t *testing.T) {
	clusterSchema := ResourceCluster().Schema
	for _, key := range []string{"cn_node_class", "dn_node_class", "cn_node_num", "dn_node_num"} {
		field := clusterSchema[key]
		if !field.Optional || !field.Computed || field.ForceNew {
			t.Fatalf("%s must be Optional and Computed without ForceNew", key)
		}
	}
	for _, key := range []string{"cn_node_num", "dn_node_num"} {
		if clusterSchema[key].Type != schema.TypeInt {
			t.Fatalf("%s must use TypeInt", key)
		}
	}
	if len(clusterSchema["db_node_count"].ConflictsWith) != 2 {
		t.Fatal("db_node_count must conflict with distributed node count fields")
	}
	for _, key := range []string{"modify_dist_nodes", "delete_sub_groups", "restart_cluster"} {
		if _, ok := clusterSchema[key]; ok {
			t.Fatalf("command-style field %s must not be part of desired state", key)
		}
	}
}
func TestUnitClassifyPolarDBDistributedNode(t *testing.T) {
	tests := []struct {
		name     string
		node     map[string]interface{}
		expected string
		wantErr  bool
	}{
		{name: "compute type", node: map[string]interface{}{"DBNodeId": "cn-1", "DBNodeType": "Compute"}, expected: "cn"},
		{name: "data type", node: map[string]interface{}{"DBNodeId": "dn-1", "DBNodeType": "Data"}, expected: "dn"},
		{name: "compute subgroup", node: map[string]interface{}{"DBNodeId": "cn-0", "SubGroupType": "Compute"}, expected: "cn"},
		{name: "data subgroup", node: map[string]interface{}{"DBNodeId": "dn-0", "SubGroupType": "Data"}, expected: "dn"},
		{name: "CN role", node: map[string]interface{}{"DBNodeId": "cn-2", "DBNodeRole": "CN"}, expected: "cn"},
		{name: "DN role", node: map[string]interface{}{"DBNodeId": "dn-2", "DBNodeRole": "DN"}, expected: "dn"},
		{name: "legacy CN", node: map[string]interface{}{"DBNodeId": "cn-3", "DBNodeClass": "polar.pg.x4.medium", "IsPrimaryCN": false}, expected: "cn"},
		{name: "legacy CN reader", node: map[string]interface{}{"DBNodeId": "cn-4", "DBNodeClass": "polar.pg.x4.medium", "DBNodeRole": "Reader", "IsPrimaryCN": false}, expected: "cn"},
		{name: "legacy DN", node: map[string]interface{}{"DBNodeId": "dn-3", "DBNodeClass": "polar.pg.x4.medium"}, expected: "dn"},
		{name: "unknown role", node: map[string]interface{}{"DBNodeId": "node-1", "DBNodeClass": "polar.pg.x4.medium", "DBNodeRole": "Reader"}, wantErr: true},
		{name: "missing role", node: map[string]interface{}{"DBNodeId": "node-2"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := classifyPolarDBDistributedNode(tt.node)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected classification error, got %q", actual)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected classification error: %s", err)
			}
			if actual != tt.expected {
				t.Fatalf("expected %q, got %q", tt.expected, actual)
			}
		})
	}
}
func TestUnitValidatePolarDBDistributedNodeCounts(t *testing.T) {
	clusterInfo := map[string]interface{}{"CnNodeCount": 1, "DnNodeCount": 2}
	if err := validatePolarDBDistributedNodeCounts(clusterInfo, 1, 2); err != nil {
		t.Fatalf("expected subgroup counts to match logical counts: %s", err)
	}
	if err := validatePolarDBDistributedNodeCounts(clusterInfo, 3, 0); err == nil {
		t.Fatal("expected silently misclassified nodes to be rejected")
	}
}
func TestUnitSelectPolarDBDistributedSubGroupsForDeletion(t *testing.T) {
	nodes := []map[string]interface{}{
		{"DBNodeId": "pi-primary-writer", "SubGroupName": "pc-test", "CreationTime": "2026-01-01T00:00:00Z"},
		{"DBNodeId": "pi-primary-reader", "SubGroupName": "pc-test", "CreationTime": "2026-01-01T00:00:00Z"},
		{"DBNodeId": "pi-extra-writer", "SubGroupName": "extra-cn", "CreationTime": "2026-01-02T00:00:00Z"},
		{"DBNodeId": "pi-extra-reader", "SubGroupName": "extra-cn", "CreationTime": "2026-01-02T00:00:00Z"},
		{"DBNodeId": "pi-newest", "SubGroupName": "newest-cn", "CreationTime": "2026-01-03T00:00:00Z"},
		{"DBNodeId": "pi-missing-subgroup"},
	}
	selected := selectPolarDBDistributedSubGroupsForDeletion(nodes, "cn", "pc-test", 1)
	if len(selected) != 1 || selected[0] != "newest-cn" {
		t.Fatalf("expected the newest non-primary logical subgroup, got %#v", selected)
	}
	selected = selectPolarDBDistributedSubGroupsForDeletion(nodes, "dn", "pc-test", 2)
	if len(selected) != 2 || selected[0] != "newest-cn" || selected[1] != "extra-cn" {
		t.Fatalf("expected deterministic DN selection without the primary or missing subgroup, got %#v", selected)
	}
}
func TestUnitPolarDBParameterModificationPendingStates(t *testing.T) {
	states := []string{"Rebooting", "Running"}
	nextState := 0
	stateConf := helper.BuildStateConf(
		polarDBParameterModificationPendingStates(),
		[]string{"Running"},
		time.Second,
		0,
		func() (interface{}, string, error) {
			state := states[nextState]
			nextState++
			return state, state, nil
		},
	)
	stateConf.PollInterval = time.Millisecond

	result, err := stateConf.WaitForState()
	if err != nil {
		t.Fatalf("expected Rebooting to be treated as pending: %s", err)
	}
	if result != "Running" {
		t.Fatalf("expected waiter to finish in Running, got %v", result)
	}
}
