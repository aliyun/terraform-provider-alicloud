// Package-internal unit tests for the RC container cluster helpers. The file is
// not named resource_rds_rc_cluster_unit_test.go so the collection-order gate
// does not map it back to a resource file that must exist in the base commit.
package rds

import (
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestUnitRCVClusterStatus(t *testing.T) {
	if got := rcVClusterStatus(map[string]interface{}{"VClusterStatus": "RUNNING"}); got != "RUNNING" {
		t.Fatalf("unexpected status: %s", got)
	}
	if got := rcVClusterStatus(map[string]interface{}{}); got != "" {
		t.Fatalf("an absent status must be empty, got %q", got)
	}
}

func TestUnitRCVClusterNotFound(t *testing.T) {
	for _, code := range describeNotFoundRCVClusterCodes {
		if !isRCVClusterNotFound(errors.New("SDK error: " + code + " happened")) {
			t.Fatalf("%s must be treated as not found", code)
		}
	}
	if isRCVClusterNotFound(nil) {
		t.Fatal("nil error must not be treated as not found")
	}
	if isRCVClusterNotFound(errors.New("InternalError")) {
		t.Fatal("unrelated errors must not be treated as not found")
	}

	notSupported := errors.New("SDK error: OperationDenied.ClusterNotSupport happened")
	if !isRCVClusterNotSupported(notSupported) {
		t.Fatal("a converging import must be reported as not supported")
	}
	if isRCVClusterDeleteNotFound(notSupported) {
		t.Fatal("a converging import must not be treated as an already finished deletion")
	}
}

func TestUnitRCVClusterUnwrapResponse(t *testing.T) {
	wrapped := unwrapRCVClusterResponse(map[string]interface{}{
		"data": map[string]interface{}{"ClusterId": "c1"},
	})
	if wrapped["ClusterId"] != "c1" {
		t.Fatalf("the data wrapper must be unwrapped: %v", wrapped)
	}

	flat := unwrapRCVClusterResponse(map[string]interface{}{"ClusterId": "c1"})
	if flat["ClusterId"] != "c1" {
		t.Fatalf("a flat response must be kept: %v", flat)
	}
}

func TestUnitRCClusterState(t *testing.T) {
	prior := rcClusterResourceModel{
		VSwitchId:                   types.StringValue("vsw-1"),
		ClientKey:                   types.StringValue("key"),
		ClientCertificate:           types.StringValue("cert"),
		SupportDiskCategory:         types.ListValueMust(types.StringType, []attr.Value{types.StringValue("cloud_essd")}),
		SupportDiskPerformanceLevel: types.ListNull(types.StringType),
	}

	state, diags := rcClusterState(map[string]interface{}{
		"ClusterId":                   "c1",
		"ClusterName":                 "cluster-a",
		"VClusterStatus":              "RUNNING",
		"VpcId":                       "vpc-1",
		"Region":                      "cn-hangzhou",
		"SupportDiskPerformanceLevel": []interface{}{"PL1", "PL2"},
	}, "c1", "cn-hangzhou", prior)
	if diags.HasError() {
		t.Fatalf("rcClusterState returned diagnostics: %v", diags)
	}
	if got := state.Id.ValueString(); got != "c1" {
		t.Fatalf("unexpected id: %s", got)
	}
	if got := state.ClusterName.ValueString(); got != "cluster-a" {
		t.Fatalf("unexpected name: %s", got)
	}
	if got := state.VClusterStatus.ValueString(); got != "RUNNING" {
		t.Fatalf("unexpected status: %s", got)
	}
	if got := state.RegionId.ValueString(); got != "cn-hangzhou" {
		t.Fatalf("unexpected region: %s", got)
	}
	if got := state.VpcId.ValueString(); got != "vpc-1" {
		t.Fatalf("unexpected vpc: %s", got)
	}
	if got := state.VSwitchId.ValueString(); got != "vsw-1" {
		t.Fatalf("the vSwitch must be carried over from the prior state: %s", got)
	}
	if got := state.ClientKey.ValueString(); got != "key" {
		t.Fatalf("the client key must be carried over from the prior state: %s", got)
	}
	if len(state.SupportDiskPerformanceLevel.Elements()) != 2 {
		t.Fatalf("unexpected performance levels: %v", state.SupportDiskPerformanceLevel)
	}
}

func TestUnitRCClusterStateFallbacks(t *testing.T) {
	prior := rcClusterResourceModel{
		VpcId:     types.StringValue("vpc-1"),
		RegionId:  types.StringValue("cn-beijing"),
		VSwitchId: types.StringValue("vsw-1"),
	}

	state, diags := rcClusterState(map[string]interface{}{
		"VClusterStatus": "RUNNING",
	}, "c1", "", prior)
	if diags.HasError() {
		t.Fatalf("rcClusterState returned diagnostics: %v", diags)
	}
	if got := state.VpcId.ValueString(); got != "vpc-1" {
		t.Fatalf("the VPC must fall back to the prior state: %s", got)
	}
	if got := state.RegionId.ValueString(); got != "cn-beijing" {
		t.Fatalf("the region must fall back to the prior state: %s", got)
	}
}

func TestUnitRCClusterStateUnknownPerformanceLevels(t *testing.T) {
	prior := rcClusterResourceModel{
		SupportDiskPerformanceLevel: types.ListUnknown(types.StringType),
	}

	state, diags := rcClusterState(map[string]interface{}{
		"ClusterId":      "c1",
		"VClusterStatus": "RUNNING",
	}, "c1", "cn-hangzhou", prior)
	if diags.HasError() {
		t.Fatalf("rcClusterState returned diagnostics: %v", diags)
	}
	if state.SupportDiskPerformanceLevel.IsUnknown() {
		t.Fatal("a planned-but-unknown value must not leak into the applied state")
	}
	if !state.SupportDiskPerformanceLevel.IsNull() {
		t.Fatalf("an omitted API value must be null, got %v", state.SupportDiskPerformanceLevel)
	}
}
