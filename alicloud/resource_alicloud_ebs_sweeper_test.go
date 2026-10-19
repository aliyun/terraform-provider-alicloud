package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// EBS (块存储/Elastic Block Storage, product "ebs" / "2021-07-30") sweeper
// gaps. All EBS OpenAPI actions are RPC POST. This file fills the following
// resources that lacked a sweeper:
//   - alicloud_ebs_solution_instance             (DescribeSolutionInstance / DeleteSolutionInstance)
//   - alicloud_ebs_enterprise_snapshot_policy    (DescribeEnterpriseSnapshotPolicy / DeleteEnterpriseSnapshotPolicy)
//   - alicloud_ebs_enterprise_snapshot_policy_attachment (DescribeEnterpriseSnapshotPolicy / UnbindEnterpriseSnapshotPolicy)
//   - alicloud_ebs_replica_group_drill           (DescribeReplicaGroupDrills / ClearReplicaGroupDrill)
//   - alicloud_ebs_replica_pair_drill            (DescribePairDrills / ClearPairDrill)
//   - alicloud_ebs_disk_replica_pair             (DescribeDiskReplicaPairs / DeleteDiskReplicaPair)
//
// DEFERRED (no sweeper written):
//   - alicloud_ebs_dedicated_block_storage_cluster: resource Delete is a no-op
//     ("Cannot destroy. Terraform will remove this resource from the state
//     file, however resources may remain.") — there is no delete API for a
//     dedicated block storage cluster, so sweeping has no effect.
//   - alicloud_ebs_disk_replica_group: already has an inline sweeper in
//     resource_alicloud_ebs_disk_replica_group_test.go (DescribeDiskReplicaGroups
//     / DeleteDiskReplicaGroup); not re-registered here to avoid init panic.
//
// Fleet-wipe pattern: SWEEP_ALL (ALICLOUD_SWEEP_ALL_RESOURCES=true) bypasses
// the tf-prefix filter so all non-NODELETING resources are swept; NORMAL mode
// keeps the prefix filter. The protectedSweepVpcName NODELETING guard (21
// markers) is always-active in both modes. All sweepers return nil — errors
// are logged and never propagated, so one failing item does not abort the
// sweep.
//
// Ordering (parent.Dependencies = [child], child runs first):
//   - alicloud_ebs_enterprise_snapshot_policy depends on
//     alicloud_ebs_enterprise_snapshot_policy_attachment (unbind disks before
//     deleting the policy).
//   - alicloud_ebs_disk_replica_pair depends on
//     alicloud_ebs_replica_pair_drill (clear pair drills before deleting the
//     pair).
//   - alicloud_ebs_replica_group_drill has no Dependencies because its parent
//     alicloud_ebs_disk_replica_group is an existing inline sweeper (in
//     resource_alicloud_ebs_disk_replica_group_test.go) that cannot be
//     modified here. It lists replica groups itself and is best-effort: if the
//     group sweeper already deleted a tf-test group, its orphaned drills
//     cannot be enumerated (DescribeReplicaGroupDrills requires GroupId).

func init() {
	resource.AddTestSweepers("alicloud_ebs_solution_instance", &resource.Sweeper{
		Name: "alicloud_ebs_solution_instance",
		F:    testSweepEbsSolutionInstance,
	})
	resource.AddTestSweepers("alicloud_ebs_enterprise_snapshot_policy_attachment", &resource.Sweeper{
		Name: "alicloud_ebs_enterprise_snapshot_policy_attachment",
		F:    testSweepEbsEnterpriseSnapshotPolicyAttachment,
	})
	resource.AddTestSweepers("alicloud_ebs_enterprise_snapshot_policy", &resource.Sweeper{
		Name:         "alicloud_ebs_enterprise_snapshot_policy",
		F:            testSweepEbsEnterpriseSnapshotPolicy,
		Dependencies: []string{"alicloud_ebs_enterprise_snapshot_policy_attachment"},
	})
	resource.AddTestSweepers("alicloud_ebs_replica_group_drill", &resource.Sweeper{
		Name: "alicloud_ebs_replica_group_drill",
		F:    testSweepEbsReplicaGroupDrill,
	})
	resource.AddTestSweepers("alicloud_ebs_replica_pair_drill", &resource.Sweeper{
		Name: "alicloud_ebs_replica_pair_drill",
		F:    testSweepEbsReplicaPairDrill,
	})
	resource.AddTestSweepers("alicloud_ebs_disk_replica_pair", &resource.Sweeper{
		Name:         "alicloud_ebs_disk_replica_pair",
		F:            testSweepEbsDiskReplicaPair,
		Dependencies: []string{"alicloud_ebs_replica_pair_drill"},
	})
}

var ebsSweeperPrefixes = []string{
	"tf-testAcc",
	"tf_testAcc",
	"tf-test",
	"tf_test",
	"testAcc",
	"tftestacc",
}

// ebsGetField returns the first non-empty string value of item[name] across
// the given name variants. EBS response field casing is consistent
// (PascalCase) but a few items echo back camelCase aliases, so callers pass
// both variants.
func ebsGetField(item map[string]interface{}, names ...string) string {
	for _, name := range names {
		if v, ok := item[name]; ok {
			if s := fmt.Sprint(v); s != "" && s != "<nil>" {
				return s
			}
		}
	}
	return ""
}

// ebsSweepNameOk applies the always-active NODELETING guard and the
// sweepAll()-aware tf-prefix filter. Returns true when the resource should
// be swept. SWEEP_ALL bypasses the prefix filter (fleet-wipe: clean all
// non-NODELETING); NORMAL mode keeps the prefix filter.
func ebsSweepNameOk(name string) bool {
	if name == "" {
		return false
	}
	if protectedSweepVpcName(name) {
		log.Printf("[INFO] Skipping ebs resource (NODELETING): %s", name)
		return false
	}
	if !sweepAll() {
		skip := true
		for _, prefix := range ebsSweeperPrefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping ebs resource (no prefix): %s", name)
			return false
		}
	}
	return true
}

// ebsRpcPost wraps an EBS RPC POST call with retry. NeedRetry errors are
// retried for 5 minutes with incremental backoff; other errors surface
// immediately.
func ebsRpcPost(client *connectivity.AliyunClient, action string, query, body map[string]interface{}, autoRetry bool) (map[string]interface{}, error) {
	var response map[string]interface{}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err := retry.Retry(5*time.Minute, func() *retry.RetryError {
		resp, err := client.RpcPost("ebs", "2021-07-30", action, query, body, autoRetry)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		response = resp
		return nil
	})
	addDebug(action, response, body)
	return response, err
}

// ebsDeleteOne performs a best-effort delete. NotFoundError and the given
// expected error codes are treated as success (resource already gone);
// other errors are logged and the sweep continues.
func ebsDeleteOne(client *connectivity.AliyunClient, action string, query, body map[string]interface{}, id string, okErrs []string) {
	_, err := ebsRpcPost(client, action, query, body, false)
	if err != nil {
		if NotFoundError(err) || (len(okErrs) > 0 && IsExpectedErrors(err, okErrs)) {
			log.Printf("[INFO] %s already gone (%s): %s", action, id, err)
			return
		}
		log.Printf("[ERROR] Failed to %s (%s): %s", action, id, err)
		return
	}
	log.Printf("[INFO] %s success: %s", action, id)
}

// testSweepEbsSolutionInstance lists solution instances via
// DescribeSolutionInstance (omitting SolutionInstanceId.1 to list all) and
// deletes tf-test-prefixed ones via DeleteSolutionInstance.
func testSweepEbsSolutionInstance(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	action := "DescribeSolutionInstance"
	body := map[string]interface{}{
		"RegionId":   client.RegionId,
		"MaxResults": PageSizeXLarge,
	}
	for {
		response, err := ebsRpcPost(client, action, nil, body, true)
		if err != nil {
			log.Printf("[ERROR] %s get an error: %#v", action, err)
			return nil
		}
		resp, jpErr := jsonpath.Get("$.Data", response)
		if jpErr != nil {
			log.Printf("[ERROR] Getting resource attribute by path %s failed!!! Body: %v.", "$.Data", response)
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item, ok := v.(map[string]interface{})
			if !ok {
				continue
			}
			name := ebsGetField(item, "Name")
			if !ebsSweepNameOk(name) {
				continue
			}
			instanceId := ebsGetField(item, "SolutionInstanceId", "InstanceId")
			if instanceId == "" {
				continue
			}
			ebsDeleteOne(client, "DeleteSolutionInstance",
				map[string]interface{}{"SolutionInstanceId": instanceId},
				map[string]interface{}{"RegionId": client.RegionId, "ClientToken": buildClientToken("DeleteSolutionInstance")},
				instanceId, nil)
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			body["NextToken"] = nextToken
		} else {
			break
		}
	}
	return nil
}

// testSweepEbsEnterpriseSnapshotPolicy lists enterprise snapshot policies via
// DescribeEnterpriseSnapshotPolicy (omitting PolicyIds to list all) and
// deletes tf-test-prefixed ones via DeleteEnterpriseSnapshotPolicy.
func testSweepEbsEnterpriseSnapshotPolicy(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	action := "DescribeEnterpriseSnapshotPolicy"
	body := map[string]interface{}{
		"RegionId":   client.RegionId,
		"MaxResults": PageSizeXLarge,
	}
	for {
		response, err := ebsRpcPost(client, action, nil, body, true)
		if err != nil {
			log.Printf("[ERROR] %s get an error: %#v", action, err)
			return nil
		}
		resp, jpErr := jsonpath.Get("$.Policies", response)
		if jpErr != nil {
			log.Printf("[ERROR] Getting resource attribute by path %s failed!!! Body: %v.", "$.Policies", response)
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item, ok := v.(map[string]interface{})
			if !ok {
				continue
			}
			name := ebsGetField(item, "Name")
			if !ebsSweepNameOk(name) {
				continue
			}
			policyId := ebsGetField(item, "PolicyId")
			if policyId == "" {
				continue
			}
			ebsDeleteOne(client, "DeleteEnterpriseSnapshotPolicy",
				map[string]interface{}{"PolicyId": policyId},
				map[string]interface{}{"RegionId": client.RegionId, "ClientToken": buildClientToken("DeleteEnterpriseSnapshotPolicy")},
				policyId, nil)
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			body["NextToken"] = nextToken
		} else {
			break
		}
	}
	return nil
}

// testSweepEbsEnterpriseSnapshotPolicyAttachment lists enterprise snapshot
// policies via DescribeEnterpriseSnapshotPolicy, and for each tf-test
// policy, unbinds each bound disk via UnbindEnterpriseSnapshotPolicy. The
// attachment resource id is policyId:diskId; there is no dedicated list API
// for attachments, so the parent policy's DiskIds array is the source of
// truth. Bindings are auto-removed when the parent policy is deleted, but
// unbinding first avoids DeleteEnterpriseSnapshotPolicy failing on bound
// disks.
func testSweepEbsEnterpriseSnapshotPolicyAttachment(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	action := "DescribeEnterpriseSnapshotPolicy"
	body := map[string]interface{}{
		"RegionId":   client.RegionId,
		"MaxResults": PageSizeXLarge,
	}
	for {
		response, err := ebsRpcPost(client, action, nil, body, true)
		if err != nil {
			log.Printf("[ERROR] %s get an error: %#v", action, err)
			return nil
		}
		resp, jpErr := jsonpath.Get("$.Policies", response)
		if jpErr != nil {
			log.Printf("[ERROR] Getting resource attribute by path %s failed!!! Body: %v.", "$.Policies", response)
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item, ok := v.(map[string]interface{})
			if !ok {
				continue
			}
			name := ebsGetField(item, "Name")
			if !ebsSweepNameOk(name) {
				continue
			}
			policyId := ebsGetField(item, "PolicyId")
			if policyId == "" {
				continue
			}
			diskIdsRaw, ok := item["DiskIds"]
			if !ok || diskIdsRaw == nil {
				continue
			}
			diskIds, _ := diskIdsRaw.([]interface{})
			for _, d := range diskIds {
				diskId := fmt.Sprint(d)
				if diskId == "" || diskId == "<nil>" {
					continue
				}
				ebsDeleteOne(client, "UnbindEnterpriseSnapshotPolicy",
					map[string]interface{}{"PolicyId": policyId, "DiskTargets.1": diskId},
					map[string]interface{}{"RegionId": client.RegionId, "ClientToken": buildClientToken("UnbindEnterpriseSnapshotPolicy")},
					policyId+":"+diskId, nil)
			}
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			body["NextToken"] = nextToken
		} else {
			break
		}
	}
	return nil
}

// testSweepEbsReplicaGroupDrill clears disaster-recovery drills started for
// tf-test replica groups. DescribeReplicaGroupDrills requires GroupId, so
// the sweeper first lists replica groups (DescribeDiskReplicaGroups),
// collects tf-test group ids, then for each lists and clears its drills via
// ClearReplicaGroupDrill. The drill itself has no name field; it is gated by
// the parent group name.
func testSweepEbsReplicaGroupDrill(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	// Step 1: list replica groups and collect tf-test group ids.
	groupAction := "DescribeDiskReplicaGroups"
	groupBody := map[string]interface{}{
		"RegionId":   client.RegionId,
		"MaxResults": PageSizeXLarge,
	}
	tfTestGroupIDs := make([]string, 0)
	for {
		response, err := ebsRpcPost(client, groupAction, nil, groupBody, true)
		if err != nil {
			log.Printf("[ERROR] %s get an error: %#v", groupAction, err)
			return nil
		}
		resp, jpErr := jsonpath.Get("$.ReplicaGroups", response)
		if jpErr != nil {
			log.Printf("[ERROR] Getting resource attribute by path %s failed!!! Body: %v.", "$.ReplicaGroups", response)
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item, ok := v.(map[string]interface{})
			if !ok {
				continue
			}
			name := ebsGetField(item, "GroupName")
			if !ebsSweepNameOk(name) {
				continue
			}
			if groupId := ebsGetField(item, "ReplicaGroupId"); groupId != "" {
				tfTestGroupIDs = append(tfTestGroupIDs, groupId)
			}
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			groupBody["NextToken"] = nextToken
		} else {
			break
		}
	}
	if len(tfTestGroupIDs) == 0 {
		return nil
	}

	// Step 2: for each tf-test group, list and clear its drills.
	drillAction := "DescribeReplicaGroupDrills"
	for _, groupId := range tfTestGroupIDs {
		drillQuery := map[string]interface{}{
			"GroupId":    groupId,
			"RegionId":   client.RegionId,
			"MaxResults": PageSizeXLarge,
		}
		for {
			response, err := ebsRpcPost(client, drillAction, drillQuery, nil, true)
			if err != nil {
				log.Printf("[ERROR] %s get an error (group %s): %#v", drillAction, groupId, err)
				break
			}
			resp, jpErr := jsonpath.Get("$.Drills", response)
			if jpErr != nil {
				log.Printf("[ERROR] Getting resource attribute by path %s failed (group %s)!!! Body: %v.", "$.Drills", groupId, response)
				break
			}
			result, _ := resp.([]interface{})
			for _, v := range result {
				item, ok := v.(map[string]interface{})
				if !ok {
					continue
				}
				drillId := ebsGetField(item, "DrillId")
				if drillId == "" {
					continue
				}
				ebsDeleteOne(client, "ClearReplicaGroupDrill",
					map[string]interface{}{"GroupId": groupId, "DrillId": drillId},
					map[string]interface{}{"RegionId": client.RegionId},
					groupId+":"+drillId, nil)
			}
			if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
				drillQuery["NextToken"] = nextToken
			} else {
				break
			}
		}
	}
	return nil
}

// testSweepEbsReplicaPairDrill clears disaster-recovery drills started for
// tf-test replica pairs. DescribePairDrills requires PairId, so the sweeper
// first lists replica pairs (DescribeDiskReplicaPairs), collects tf-test
// pair ids, then for each lists and clears its drills via ClearPairDrill.
// The drill itself has no name field; it is gated by the parent pair name.
func testSweepEbsReplicaPairDrill(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	// Step 1: list replica pairs and collect tf-test pair ids.
	pairAction := "DescribeDiskReplicaPairs"
	pairBody := map[string]interface{}{
		"RegionId":   client.RegionId,
		"MaxResults": PageSizeXLarge,
	}
	tfTestPairIDs := make([]string, 0)
	for {
		response, err := ebsRpcPost(client, pairAction, nil, pairBody, true)
		if err != nil {
			log.Printf("[ERROR] %s get an error: %#v", pairAction, err)
			return nil
		}
		resp, jpErr := jsonpath.Get("$.ReplicaPairs", response)
		if jpErr != nil {
			log.Printf("[ERROR] Getting resource attribute by path %s failed!!! Body: %v.", "$.ReplicaPairs", response)
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item, ok := v.(map[string]interface{})
			if !ok {
				continue
			}
			name := ebsGetField(item, "PairName")
			if !ebsSweepNameOk(name) {
				continue
			}
			if pairId := ebsGetField(item, "ReplicaPairId"); pairId != "" {
				tfTestPairIDs = append(tfTestPairIDs, pairId)
			}
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			pairBody["NextToken"] = nextToken
		} else {
			break
		}
	}
	if len(tfTestPairIDs) == 0 {
		return nil
	}

	// Step 2: for each tf-test pair, list and clear its drills.
	drillAction := "DescribePairDrills"
	for _, pairId := range tfTestPairIDs {
		drillQuery := map[string]interface{}{
			"PairId":     pairId,
			"RegionId":   client.RegionId,
			"MaxResults": PageSizeXLarge,
		}
		for {
			response, err := ebsRpcPost(client, drillAction, drillQuery, nil, true)
			if err != nil {
				log.Printf("[ERROR] %s get an error (pair %s): %#v", drillAction, pairId, err)
				break
			}
			resp, jpErr := jsonpath.Get("$.Drills", response)
			if jpErr != nil {
				log.Printf("[ERROR] Getting resource attribute by path %s failed (pair %s)!!! Body: %v.", "$.Drills", pairId, response)
				break
			}
			result, _ := resp.([]interface{})
			for _, v := range result {
				item, ok := v.(map[string]interface{})
				if !ok {
					continue
				}
				drillId := ebsGetField(item, "DrillId")
				if drillId == "" {
					continue
				}
				ebsDeleteOne(client, "ClearPairDrill",
					map[string]interface{}{"PairId": pairId, "DrillId": drillId},
					map[string]interface{}{"RegionId": client.RegionId},
					pairId+":"+drillId, nil)
			}
			if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
				drillQuery["NextToken"] = nextToken
			} else {
				break
			}
		}
	}
	return nil
}

// testSweepEbsDiskReplicaPair lists disk async-replication pairs via
// DescribeDiskReplicaPairs and deletes tf-test-prefixed ones via
// DeleteDiskReplicaPair. Depends on alicloud_ebs_replica_pair_drill so pair
// drills are cleared before the pair is deleted.
func testSweepEbsDiskReplicaPair(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	action := "DescribeDiskReplicaPairs"
	body := map[string]interface{}{
		"RegionId":   client.RegionId,
		"MaxResults": PageSizeXLarge,
	}
	for {
		response, err := ebsRpcPost(client, action, nil, body, true)
		if err != nil {
			log.Printf("[ERROR] %s get an error: %#v", action, err)
			return nil
		}
		resp, jpErr := jsonpath.Get("$.ReplicaPairs", response)
		if jpErr != nil {
			log.Printf("[ERROR] Getting resource attribute by path %s failed!!! Body: %v.", "$.ReplicaPairs", response)
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item, ok := v.(map[string]interface{})
			if !ok {
				continue
			}
			name := ebsGetField(item, "PairName")
			if !ebsSweepNameOk(name) {
				continue
			}
			pairId := ebsGetField(item, "ReplicaPairId")
			if pairId == "" {
				continue
			}
			ebsDeleteOne(client, "DeleteDiskReplicaPair",
				nil,
				map[string]interface{}{"ReplicaPairId": pairId, "RegionId": client.RegionId, "ClientToken": buildClientToken("DeleteDiskReplicaPair")},
				pairId, nil)
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			body["NextToken"] = nextToken
		} else {
			break
		}
	}
	return nil
}
