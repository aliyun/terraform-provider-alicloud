package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// ECS (Elastic Compute Service, product "Ecs" / "2014-05-26", RPC style)
// child-resource sweepers — fills the gaps left by the 12 pre-existing inline
// sweepers (activation, auto_snapshot_policy, deployment_set,
// dedicated_host_cluster, disk, image_component, image_pipeline,
// launch_template, network_interface, prefix_list, snapshot, snapshot_group in
// their own _test.go files). ECS is RPC-style:
// RpcPost(product, version, action, query=nil, body=params, autoRetry); the
// existing ECS sweepers put the params in the BODY slot (query=nil), and so do
// these.
//
// Seven gap sweepers are added here (all region-scoped, enumerable + deletable):
//   DescribeCapacityReservations -> ReleaseCapacityReservation     (PrivatePoolOptions.Id; list id field is PrivatePoolOptionsId)
//   DescribeCommands           -> DeleteCommand                    (CommandId)
//   DescribeDedicatedHosts     -> ReleaseDedicatedHost             (DedicatedHostId; PrePaid hosts error out, best-effort skip)
//   DescribeHpcClusters        -> DeleteHpcCluster                 (HpcClusterId + ClientToken)
//   DescribeInvocations        -> StopInvocation                   (InvokeId; stop, not delete — record persists)
//   DescribeKeyPairs           -> DeleteKeyPairs                   (KeyPairNames = JSON array string; id == KeyPairName)
//   DescribeNetworkInterfacePermissions -> DeleteNetworkInterfacePermission (NetworkInterfacePermissionId + Force; no Name field, gated on id -> only SWEEP_ALL)
//
// children-first DAG: command depends on invocation (running invocations
// stopped before the command is deleted).
//
// DEFERRED (12 gap resources — NOT registered, reasons below):
//  attachments (cascade-cleaned by their parent resource sweepers, no
//  standalone list-all API):
//   - ecs_auto_snapshot_policy_attachment: CancelAutoSnapshotPolicy (policy<->disk); cleaned by ecs_disk + ecs_auto_snapshot_policy sweepers
//   - ecs_disk_attachment: DetachDisk (disk<->instance); cleaned by ecs_disk + alicloud_instance sweepers
//   - ecs_key_pair_attachment: DetachKeyPair (keypair<->instance); no standalone list; cleaned by ecs_key_pair + alicloud_instance
//   - ecs_network_interface_attachment: DetachNetworkInterface (eni<->instance); cleaned by ecs_network_interface sweeper
//   - ecs_ram_role_attachment: DetachInstanceRamRole (ramrole<->instance); DescribeInstanceRamRole is per-instance only
//  account-level single settings (no region-scoped list-all):
//   - ecs_disk_default_kms_key: ResetDiskDefaultKMSKeyId (account-level default)
//   - ecs_disk_encryption_by_default: DisableDiskEncryptionByDefault (account-level toggle)
//   - ecs_session_manager_status: ModifyUserBusinessBehavior (account-level toggle)
//  refund/退订 (BssOpenApi RefundInstance — red-line excludes BSS refunds):
//   - ecs_elasticity_assurance: RefundInstance via BssOpenApi (prepaid reserved capacity)
//  no Delete API (TF Delete is a no-op "Cannot destroy ... resources may remain"):
//   - ecs_storage_capacity_unit: Prepaid SCU, Delete is a no-op
//   - ecs_image_pipeline_execution: transient execution record, Delete is a no-op (cascade-cleaned by pipeline)
//   - ecs_instance_set: creates plain ECS instances swept by the existing
//     alicloud_instance sweeper (tf-prefixed instances); id is base64-encoded
//     instance IDs with no list-by-set marker, so re-sweeping would duplicate alicloud_instance.

var ecsSweeperPrefixes = []string{
	"tf-testAcc",
	"tf_testAcc",
	"tf-test",
	"tf_test",
	"tf-example",
	"tf_example",
	"terraform-example",
	"tftestacc",
	"testAcc",
}

// ecsGetField returns the first non-empty value among the given keys on item,
// trying both as-is and case-insensitively (ECS responses use consistent
// PascalCase, but case-insensitive is safer across actions).
func ecsGetField(item map[string]interface{}, names ...string) string {
	if item == nil {
		return ""
	}
	for _, name := range names {
		if v, ok := item[name]; ok {
			if s := fmt.Sprint(v); s != "" && s != "<nil>" {
				return s
			}
		}
	}
	wanted := make(map[string]bool, len(names))
	for _, name := range names {
		wanted[strings.ToLower(name)] = true
	}
	for k, v := range item {
		if wanted[strings.ToLower(k)] {
			if s := fmt.Sprint(v); s != "" && s != "<nil>" {
				return s
			}
		}
	}
	return ""
}

// ecsSweepNameOk is the fleet-wipe gate: never touch NODELETING names; in
// normal (non-SWEEP_ALL) mode only sweep tf-prefixed names. SWEEP_ALL bypasses
// the prefix filter so the whole test account (1511928242963727, cn-beijing) is
// fleet-wiped. protectedSweepVpcName is always active (both modes).
func ecsSweepNameOk(name string) bool {
	if name == "" {
		return false
	}
	if protectedSweepVpcName(name) {
		log.Printf("[INFO] Skipping ecs resource (NODELETING): %s", name)
		return false
	}
	if !sweepAll() {
		skip := true
		for _, prefix := range ecsSweeperPrefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping ecs resource (no prefix): %s", name)
			return false
		}
	}
	return true
}

// ecsRpcPost issues an RpcPost list/read with retry.Retry + NeedRetry +
// incrementalWait. Params go in the BODY slot (query=nil), matching the
// pre-existing ECS sweepers.
func ecsRpcPost(client *connectivity.AliyunClient, action string, params map[string]interface{}) (map[string]interface{}, error) {
	var response map[string]interface{}
	var err error
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Ecs", "2014-05-26", action, nil, params, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		addDebug(action, response, params)
		return nil
	})
	return response, err
}

// ecsDeletePost issues an RpcPost delete/release with retry. NotFound and
// expected-already-gone errors are surfaced to the caller, which treats them as
// success (best-effort: never returns an error that aborts the whole sweep).
func ecsDeletePost(client *connectivity.AliyunClient, action string, params map[string]interface{}) error {
	var response map[string]interface{}
	var err error
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Ecs", "2014-05-26", action, nil, params, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		addDebug(action, response, params)
		return nil
	})
	return err
}

// ecsListPaged enumerates an action with PageNumber/PageSize pagination and
// returns the slice at path (a jsonpath like "$.Commands.Command"). A missing
// path is treated as an empty list (logged) so the sweeper continues with zero
// items rather than aborting.
func ecsListPaged(client *connectivity.AliyunClient, action string, base map[string]interface{}, path string) []map[string]interface{} {
	request := map[string]interface{}{}
	for k, v := range base {
		request[k] = v
	}
	request["RegionId"] = client.RegionId
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	out := make([]map[string]interface{}, 0)
	for {
		response, err := ecsRpcPost(client, action, request)
		if err != nil {
			log.Printf("[ERROR] ecs %s: %s", action, err)
			return out
		}
		raw, err := jsonpath.Get(path, response)
		if err != nil {
			log.Printf("[ERROR] ecs %s jsonpath %s: %s", action, path, err)
			return out
		}
		items, _ := raw.([]interface{})
		for _, v := range items {
			if m, ok := v.(map[string]interface{}); ok {
				out = append(out, m)
			}
		}
		if len(items) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return out
}

// ecsListToken enumerates an action with MaxResults/NextToken pagination and
// returns the slice at path (a jsonpath like "$.ImageComponent.ImageComponentSet").
func ecsListToken(client *connectivity.AliyunClient, action string, base map[string]interface{}, path string) []map[string]interface{} {
	request := map[string]interface{}{}
	for k, v := range base {
		request[k] = v
	}
	request["RegionId"] = client.RegionId
	request["MaxResults"] = PageSizeLarge
	out := make([]map[string]interface{}, 0)
	for {
		response, err := ecsRpcPost(client, action, request)
		if err != nil {
			log.Printf("[ERROR] ecs %s: %s", action, err)
			return out
		}
		raw, err := jsonpath.Get(path, response)
		if err != nil {
			log.Printf("[ERROR] ecs %s jsonpath %s: %s", action, path, err)
			return out
		}
		items, _ := raw.([]interface{})
		for _, v := range items {
			if m, ok := v.(map[string]interface{}); ok {
				out = append(out, m)
			}
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	return out
}

// ecsSweepList is the shared list+delete core: enumerate via listFn, gate each
// item by name (then id), and delete via deleteAction with params from
// buildParams. NotFound is treated as success; every other error is logged and
// skipped. NEVER returns an error that would abort the whole sweep.
func ecsSweepList(region string, listFn func(c *connectivity.AliyunClient) []map[string]interface{}, idKeys, nameKeys []string, deleteAction string, buildParams func(c *connectivity.AliyunClient, id string) map[string]interface{}) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	c := rawClient.(*connectivity.AliyunClient)
	items := listFn(c)
	for _, item := range items {
		id := ecsGetField(item, idKeys...)
		name := ecsGetField(item, nameKeys...)
		if id == "" {
			continue
		}
		if !ecsSweepNameOk(name) && !ecsSweepNameOk(id) {
			continue
		}
		if e := ecsDeletePost(c, deleteAction, buildParams(c, id)); e != nil {
			if NotFoundError(e) {
				log.Printf("[INFO] ecs %s (%s) already gone: %s", deleteAction, id, e)
			} else {
				log.Printf("[ERROR] Failed to delete ecs %s (%s): %s", deleteAction, id, e)
			}
			continue
		}
		log.Printf("[INFO] Delete ecs %s success: %s (%s)", deleteAction, id, name)
	}
	return nil
}

func init() {
	resource.AddTestSweepers("alicloud_ecs_capacity_reservation", &resource.Sweeper{
		Name: "alicloud_ecs_capacity_reservation", F: testSweepEcsCapacityReservations,
	})
	resource.AddTestSweepers("alicloud_ecs_command", &resource.Sweeper{
		Name: "alicloud_ecs_command", F: testSweepEcsCommands,
		// invocations reference a command; stop them before deleting the command
		Dependencies: []string{"alicloud_ecs_invocation"},
	})
	resource.AddTestSweepers("alicloud_ecs_dedicated_host", &resource.Sweeper{
		Name: "alicloud_ecs_dedicated_host", F: testSweepEcsDedicatedHosts,
	})
	resource.AddTestSweepers("alicloud_ecs_hpc_cluster", &resource.Sweeper{
		Name: "alicloud_ecs_hpc_cluster", F: testSweepEcsHpcClusters,
	})
	resource.AddTestSweepers("alicloud_ecs_invocation", &resource.Sweeper{
		Name: "alicloud_ecs_invocation", F: testSweepEcsInvocations,
	})
	resource.AddTestSweepers("alicloud_ecs_key_pair", &resource.Sweeper{
		Name: "alicloud_ecs_key_pair", F: testSweepEcsKeyPairs,
	})
	resource.AddTestSweepers("alicloud_ecs_network_interface_permission", &resource.Sweeper{
		Name: "alicloud_ecs_network_interface_permission", F: testSweepEcsNetworkInterfacePermissions,
	})
}

// ---- region-scoped sweepers ----

// testSweepEcsCapacityReservations: DescribeCapacityReservations
// (MaxResults/NextToken) -> $.CapacityReservationSet.CapacityReservationItem;
// delete ReleaseCapacityReservation (PrivatePoolOptions.Id). The list id field
// is PrivatePoolOptionsId but the delete param key is PrivatePoolOptions.Id.
func testSweepEcsCapacityReservations(region string) error {
	return ecsSweepList(region,
		func(c *connectivity.AliyunClient) []map[string]interface{} {
			return ecsListToken(c, "DescribeCapacityReservations", nil, "$.CapacityReservationSet.CapacityReservationItem")
		},
		[]string{"PrivatePoolOptionsId"},
		[]string{"PrivatePoolOptionsName", "Description"},
		"ReleaseCapacityReservation",
		func(c *connectivity.AliyunClient, id string) map[string]interface{} {
			return map[string]interface{}{"PrivatePoolOptions.Id": id, "RegionId": c.RegionId}
		},
	)
}

// testSweepEcsCommands: DescribeCommands (PageNumber/PageSize) ->
// $.Commands.Command; delete DeleteCommand (CommandId). Name field: Name.
func testSweepEcsCommands(region string) error {
	return ecsSweepList(region,
		func(c *connectivity.AliyunClient) []map[string]interface{} {
			return ecsListPaged(c, "DescribeCommands", nil, "$.Commands.Command")
		},
		[]string{"CommandId"},
		[]string{"Name", "Description"},
		"DeleteCommand",
		func(c *connectivity.AliyunClient, id string) map[string]interface{} {
			return map[string]interface{}{"CommandId": id, "RegionId": c.RegionId}
		},
	)
}

// testSweepEcsDedicatedHosts: DescribeDedicatedHosts (PageNumber/PageSize) ->
// $.DedicatedHosts.DedicatedHost; delete ReleaseDedicatedHost
// (DedicatedHostId). PrePaid hosts cannot be released (TF Delete errors out);
// the API call will fail for them and the error is logged + skipped
// (best-effort).
func testSweepEcsDedicatedHosts(region string) error {
	return ecsSweepList(region,
		func(c *connectivity.AliyunClient) []map[string]interface{} {
			return ecsListPaged(c, "DescribeDedicatedHosts", nil, "$.DedicatedHosts.DedicatedHost")
		},
		[]string{"DedicatedHostId"},
		[]string{"DedicatedHostName", "Description"},
		"ReleaseDedicatedHost",
		func(c *connectivity.AliyunClient, id string) map[string]interface{} {
			return map[string]interface{}{"DedicatedHostId": id, "RegionId": c.RegionId}
		},
	)
}

// testSweepEcsHpcClusters: DescribeHpcClusters (PageNumber/PageSize) ->
// $.HpcClusters.HpcCluster; delete DeleteHpcCluster (HpcClusterId + ClientToken).
func testSweepEcsHpcClusters(region string) error {
	return ecsSweepList(region,
		func(c *connectivity.AliyunClient) []map[string]interface{} {
			return ecsListPaged(c, "DescribeHpcClusters", nil, "$.HpcClusters.HpcCluster")
		},
		[]string{"HpcClusterId"},
		[]string{"Name", "Description"},
		"DeleteHpcCluster",
		func(c *connectivity.AliyunClient, id string) map[string]interface{} {
			return map[string]interface{}{
				"HpcClusterId": id,
				"RegionId":     c.RegionId,
				"ClientToken":  buildClientToken("DeleteHpcCluster"),
			}
		},
	)
}

// testSweepEcsInvocations: DescribeInvocations (PageNumber/PageSize) ->
// $.Invocations.Invocation; delete StopInvocation (InvokeId). StopInvocation is
// a stop, not a hard delete (the invocation record persists), so this only
// terminates running test invocations — best-effort.
func testSweepEcsInvocations(region string) error {
	return ecsSweepList(region,
		func(c *connectivity.AliyunClient) []map[string]interface{} {
			return ecsListPaged(c, "DescribeInvocations", nil, "$.Invocations.Invocation")
		},
		[]string{"InvokeId"},
		[]string{"CommandName", "CommandId"},
		"StopInvocation",
		func(c *connectivity.AliyunClient, id string) map[string]interface{} {
			return map[string]interface{}{"InvokeId": id, "RegionId": c.RegionId}
		},
	)
}

// testSweepEcsKeyPairs: DescribeKeyPairs (PageNumber/PageSize) ->
// $.KeyPairs.KeyPair; delete DeleteKeyPairs (KeyPairNames = JSON array string).
// The key pair id IS its name (KeyPairName), so name and id gate identically.
func testSweepEcsKeyPairs(region string) error {
	return ecsSweepList(region,
		func(c *connectivity.AliyunClient) []map[string]interface{} {
			return ecsListPaged(c, "DescribeKeyPairs", nil, "$.KeyPairs.KeyPair")
		},
		[]string{"KeyPairName"},
		[]string{"KeyPairName"},
		"DeleteKeyPairs",
		func(c *connectivity.AliyunClient, id string) map[string]interface{} {
			raw, _ := json.Marshal([]string{id})
			return map[string]interface{}{
				"KeyPairNames": string(raw),
				"RegionId":     c.RegionId,
			}
		},
	)
}

// testSweepEcsNetworkInterfacePermissions: DescribeNetworkInterfacePermissions
// (PageNumber/PageSize) -> $.NetworkInterfacePermissions.NetworkInterfacePermission;
// delete DeleteNetworkInterfacePermission (NetworkInterfacePermissionId + Force).
// These permissions have NO Name field (id is system-generated), so the prefix
// gate never matches in normal mode — they are only swept in SWEEP_ALL mode.
func testSweepEcsNetworkInterfacePermissions(region string) error {
	return ecsSweepList(region,
		func(c *connectivity.AliyunClient) []map[string]interface{} {
			return ecsListPaged(c, "DescribeNetworkInterfacePermissions", nil, "$.NetworkInterfacePermissions.NetworkInterfacePermission")
		},
		[]string{"NetworkInterfacePermissionId"},
		[]string{"NetworkInterfaceId"},
		"DeleteNetworkInterfacePermission",
		func(c *connectivity.AliyunClient, id string) map[string]interface{} {
			return map[string]interface{}{
				"NetworkInterfacePermissionId": id,
				"Force":                        true,
				"RegionId":                     c.RegionId,
			}
		},
	)
}
