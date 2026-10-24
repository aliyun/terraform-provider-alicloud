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

// RDS family (RDS + R-kvstore/Redis + RdsAi) child-resource sweepers. Single
// file: alicloud/resource_alicloud_rds_sweeper_test.go.
//
// Prior commit f02a28496 added inline sweepers for alicloud_db_instance,
// alicloud_rds_ddr_instance, and alicloud_redis_tair_instance. The two
// "wrapper" resources alicloud_rds_clone_db_instance and
// alicloud_rds_upgrade_db_instance both point their F at testSweepDBInstances
// (the db_instance sweeper) and so are not real sweepers. alicloud_kvstore_instance
// has its own pre-existing inline sweeper. This file fills the remaining
// rds_*/kvstore_*/redis_* gaps.
//
// 13 IMPLEMENTED gap sweepers:
//   region-scoped (no instance walk):
//     rds_parameter_group, rds_whitelist_template, rds_custom_deployment_set,
//     rds_custom_disk (depends on rds_custom_disk_attachment),
//     rds_custom_disk_attachment, rds_ai_instance, redis_global_security_ip_group
//   instance-scoped (walk Running DB instances first):
//     rds_account, rds_backup, rds_db_instance_endpoint (depends on
//     rds_db_instance_endpoint_address), rds_db_instance_endpoint_address,
//     kvstore_account, redis_backup
//
// 8 DEFERRED (no real Delete API / cascade-cleaned / account-level singleton):
//   rds_db_node:                  no list-all API — DescribeRdsNode reuses
//                                 DescribeDBInstanceAttribute (by-instance) and
//                                 there is no per-node enumeration API.
//   rds_db_proxy:                 Delete = ModifyDBProxy(disable); singleton per
//                                 instance, cascade-cleaned by instance delete.
//   rds_db_proxy_public:          singleton per instance (public proxy address),
//                                 cascade-cleaned by instance delete.
//   rds_instance_cross_backup_policy: Delete = ModifyInstanceCrossBackupPolicy(disable);
//                                 singleton per instance.
//   rds_service_linked_role:      account-level singleton (one role per account).
//   kvstore_audit_log_config:     Delete = no-op (resource removes state only).
//   kvstore_backup_policy:        modify-only singleton, cascade-cleaned by instance delete.
//   kvstore_connection:           singleton public address, cascade-cleaned by instance delete.
//
// Products / versions:
//   RDS:        RpcPost/RpcGet product "Rds" version "2014-08-15" (RPC, params
//                in BODY slot — query=nil — matching existing RDS service code).
//   R-kvstore:   RpcPost product "R-kvstore" version "2015-01-01".
//   RdsAi:       RpcPost product "RdsAi" version "2025-05-07".
//
// red-line: only sweep test account 1511928242963727, cn-beijing. NEVER return
// error to abort the sweep — every sweeper returns nil after logging. NODELETING
// protection is always active via protectedSweepVpcName (in both prefix and
// SWEEP_ALL modes), so a NODELETING-named resource is never swept.

var rdsSweeperPrefixes = []string{
	"tf-testAcc",
	"tf_testAcc",
	"tf-example",
	"tf_example",
	"terraform-example",
	"tftestacc",
}

// rdsGetField returns the first non-empty value among the given keys on item.
// RDS / R-kvstore responses use consistent PascalCase, but a few actions
// return camelCase; this helper tolerates both by trying each key as-is.
func rdsGetField(item map[string]interface{}, names ...string) string {
	for _, name := range names {
		if v, ok := item[name]; ok {
			if s := fmt.Sprint(v); s != "" && s != "<nil>" {
				return s
			}
		}
	}
	return ""
}

// rdsSweepNameOk is the standard fleet-wipe gate: never touch NODELETING
// names, and in normal (non-SWEEP_ALL) mode only sweep resources whose name
// matches one of the tf-test prefixes. SWEEP_ALL bypasses the prefix filter so
// the whole test account (1511928242963727, cn-beijing) is fleet-wiped.
func rdsSweepNameOk(name string) bool {
	if name == "" {
		return false
	}
	if protectedSweepVpcName(name) {
		log.Printf("[INFO] Skipping rds resource (NODELETING): %s", name)
		return false
	}
	if !sweepAll() {
		skip := true
		for _, prefix := range rdsSweeperPrefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping rds resource (no prefix): %s", name)
			return false
		}
	}
	return true
}

// rdsRpcPost issues an RpcPost on the RDS product (2014-08-15) with
// retry.Retry + NeedRetry + incrementalWait. Params go in the BODY slot
// (query=nil), matching the pre-existing RDS service code and sweepers.
func rdsRpcPost(client *connectivity.AliyunClient, action string, params map[string]interface{}) (map[string]interface{}, error) {
	var response map[string]interface{}
	var err error
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Rds", "2014-08-15", action, nil, params, true)
		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, OperationDeniedDBStatus) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	return response, err
}

// rdsRpcGet issues an RpcGet on the RDS product (2014-08-15) with retry.
// DescribeRCDeploymentSets / DescribeRCDisks* are GET actions.
func rdsRpcGet(client *connectivity.AliyunClient, action string, query map[string]interface{}) (map[string]interface{}, error) {
	var response map[string]interface{}
	var err error
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcGet("Rds", "2014-08-15", action, query, nil)
		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, OperationDeniedDBStatus) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	return response, err
}

// rdsDeletePost issues a delete RpcPost on the RDS product with retry. NotFound
// (cascade-cleaned) and expected error codes are treated as success so the
// sweep never aborts on stale resources. Never returns an error that aborts
// the sweep — callers log+continue on returned errors.
func rdsDeletePost(client *connectivity.AliyunClient, action string, params map[string]interface{}, expectedErrs []string) error {
	var err error
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		_, err = client.RpcPost("Rds", "2014-08-15", action, nil, params, true)
		if err != nil {
			if NotFoundError(err) || IsExpectedErrors(err, expectedErrs) {
				return nil
			}
			if NeedRetry(err) || IsExpectedErrors(err, OperationDeniedDBStatus) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	return err
}

// rdsKVPost issues an RpcPost on R-kvstore product (2015-01-01) with retry.
func rdsKVPost(client *connectivity.AliyunClient, action string, params map[string]interface{}) (map[string]interface{}, error) {
	var response map[string]interface{}
	var err error
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("R-kvstore", "2015-01-01", action, nil, params, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	return response, err
}

// rdsKVDelete issues a delete RpcPost on R-kvstore; returns nil on NotFound.
func rdsKVDelete(client *connectivity.AliyunClient, action string, params map[string]interface{}, expectedErrs []string) error {
	var err error
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		_, err = client.RpcPost("R-kvstore", "2015-01-01", action, nil, params, true)
		if err != nil {
			if NotFoundError(err) || IsExpectedErrors(err, expectedErrs) {
				return nil
			}
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	return err
}

// rdsAiPost issues an RpcPost on the RdsAi product (2025-05-07) with retry.
func rdsAiPost(client *connectivity.AliyunClient, action string, params map[string]interface{}) (map[string]interface{}, error) {
	var response map[string]interface{}
	var err error
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("RdsAi", "2025-05-07", action, nil, params, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	return response, err
}

// rdsListRunningInstanceIDs enumerates every Running RDS DBInstance in the
// region via DescribeDBInstances (PageNumber/PageSize pagination). Child
// resources live on Running instances; released instances have already
// cascade-deleted their children, so the walk intentionally skips non-Running.
func rdsListRunningInstanceIDs(client *connectivity.AliyunClient) ([]string, []map[string]interface{}, error) {
	action := "DescribeDBInstances"
	request := map[string]interface{}{
		"RegionId":   client.RegionId,
		"PageSize":   PageSizeLarge,
		"PageNumber": 1,
	}
	instanceIDs := make([]string, 0)
	instances := make([]map[string]interface{}, 0)
	for {
		response, err := rdsRpcPost(client, action, request)
		if err != nil {
			return nil, nil, fmt.Errorf("DescribeDBInstances: %s", err)
		}
		resp, err := jsonpath.Get("$.Items.DBInstance", response)
		if err != nil {
			return nil, nil, fmt.Errorf("DescribeDBInstances jsonpath $.Items.DBInstance: %s", err)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item, ok := v.(map[string]interface{})
			if !ok {
				continue
			}
			if rdsGetField(item, "DBInstanceStatus") != "Running" {
				continue
			}
			id := rdsGetField(item, "DBInstanceId")
			if id == "" {
				continue
			}
			instanceIDs = append(instanceIDs, id)
			instances = append(instances, item)
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return instanceIDs, instances, nil
}

// rdsListKVStoreInstanceIDs enumerates every KVStore / Redis / Tair instance in
// the region via DescribeInstances (R-kvstore 2015-01-01) across all instance
// types the provider creates (Redis, Memcache, Tair rdb/scm/essd). Tair
// accounts/backups are created via the same R-kvstore API surface as Redis, so
// they need to be walked too.
func rdsListKVStoreInstanceIDs(client *connectivity.AliyunClient) ([]string, error) {
	types := []string{string(KVStoreRedis), string(KVStoreMemcache), "tair_rdb", "tair_scm", "tair_essd"}
	instanceIDs := make([]string, 0)
	for _, instanceType := range types {
		request := map[string]interface{}{
			"RegionId":     client.RegionId,
			"InstanceType": instanceType,
			"PageSize":     PageSizeXLarge,
			"PageNumber":   1,
		}
		for {
			var response map[string]interface{}
			var err error
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = retry.Retry(5*time.Minute, func() *retry.RetryError {
				response, err = client.RpcPost("R-kvstore", "2015-01-01", "DescribeInstances", nil, request, true)
				if err != nil {
					if NeedRetry(err) {
						wait()
						return retry.RetryableError(err)
					}
					return retry.NonRetryableError(err)
				}
				return nil
			})
			if err != nil {
				return nil, fmt.Errorf("DescribeInstances (%s): %s", instanceType, err)
			}
			resp, err := jsonpath.Get("$.Instances.KVStoreInstance", response)
			if err != nil {
				return nil, fmt.Errorf("DescribeInstances jsonpath $.Instances.KVStoreInstance: %s", err)
			}
			result, _ := resp.([]interface{})
			for _, v := range result {
				if item, ok := v.(map[string]interface{}); ok {
					if id := rdsGetField(item, "InstanceId"); id != "" {
						instanceIDs = append(instanceIDs, id)
					}
				}
			}
			if len(result) < PageSizeXLarge {
				break
			}
			request["PageNumber"] = request["PageNumber"].(int) + 1
		}
	}
	return instanceIDs, nil
}

// rdsInstanceSweepable decides whether a tf-test instance's children should be
// swept: applies the NODELETING + prefix/sweepAll gate to the instance's
// DBInstanceDescription (RDS) / InstanceName (R-kvstore). Returns true if the
// instance's children should be swept.
func rdsInstanceSweepable(name string) bool {
	return rdsSweepNameOk(name)
}

func init() {
	// children-first DAG: disk_attachment sweeps before disk;
	// endpoint_address sweeps before endpoint.
	resource.AddTestSweepers("alicloud_rds_custom_disk", &resource.Sweeper{
		Name: "alicloud_rds_custom_disk",
		F:    testSweepRdsCustomDisks,
		Dependencies: []string{
			"alicloud_rds_custom_disk_attachment",
		},
	})
	resource.AddTestSweepers("alicloud_rds_custom_disk_attachment", &resource.Sweeper{
		Name: "alicloud_rds_custom_disk_attachment",
		F:    testSweepRdsCustomDiskAttachments,
	})
	resource.AddTestSweepers("alicloud_rds_custom_deployment_set", &resource.Sweeper{
		Name: "alicloud_rds_custom_deployment_set",
		F:    testSweepRdsCustomDeploymentSets,
	})
	resource.AddTestSweepers("alicloud_rds_parameter_group", &resource.Sweeper{
		Name: "alicloud_rds_parameter_group",
		F:    testSweepRdsParameterGroups,
	})
	resource.AddTestSweepers("alicloud_rds_whitelist_template", &resource.Sweeper{
		Name: "alicloud_rds_whitelist_template",
		F:    testSweepRdsWhitelistTemplates,
	})
	resource.AddTestSweepers("alicloud_rds_ai_instance", &resource.Sweeper{
		Name: "alicloud_rds_ai_instance",
		F:    testSweepRdsAiInstances,
	})
	resource.AddTestSweepers("alicloud_rds_account", &resource.Sweeper{
		Name: "alicloud_rds_account",
		F:    testSweepRdsAccounts,
	})
	resource.AddTestSweepers("alicloud_rds_backup", &resource.Sweeper{
		Name: "alicloud_rds_backup",
		F:    testSweepRdsBackups,
	})
	resource.AddTestSweepers("alicloud_rds_db_instance_endpoint", &resource.Sweeper{
		Name: "alicloud_rds_db_instance_endpoint",
		F:    testSweepRdsDBInstanceEndpoints,
		Dependencies: []string{
			"alicloud_rds_db_instance_endpoint_address",
		},
	})
	resource.AddTestSweepers("alicloud_rds_db_instance_endpoint_address", &resource.Sweeper{
		Name: "alicloud_rds_db_instance_endpoint_address",
		F:    testSweepRdsDBInstanceEndpointAddresses,
	})
	resource.AddTestSweepers("alicloud_kvstore_account", &resource.Sweeper{
		Name: "alicloud_kvstore_account",
		F:    testSweepKVStoreAccounts,
	})
	resource.AddTestSweepers("alicloud_redis_backup", &resource.Sweeper{
		Name: "alicloud_redis_backup",
		F:    testSweepRedisBackups,
	})
	resource.AddTestSweepers("alicloud_redis_global_security_ip_group", &resource.Sweeper{
		Name: "alicloud_redis_global_security_ip_group",
		F:    testSweepRedisGlobalSecurityIPGroups,
	})
}

// ---- region-scoped sweepers ----

// testSweepRdsCustomDeploymentSets: DescribeRCDeploymentSets (RpcGet, RegionId,
// PageNumber/PageSize max 50) -> $.DeploymentSets.DeploymentSet; delete
// DeleteRCDeploymentSet (DeploymentSetId + RegionId).
func testSweepRdsCustomDeploymentSets(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	c := rawClient.(*connectivity.AliyunClient)
	pageNum := 1
	for {
		query := map[string]interface{}{
			"RegionId":   c.RegionId,
			"PageNumber": pageNum,
			"PageSize":   50,
		}
		response, err := rdsRpcGet(c, "DescribeRCDeploymentSets", query)
		if err != nil {
			log.Printf("[ERROR] rds DescribeRCDeploymentSets: %s", err)
			return nil
		}
		resp, err := jsonpath.Get("$.DeploymentSets.DeploymentSet", response)
		if err != nil {
			log.Printf("[ERROR] rds DescribeRCDeploymentSets jsonpath: %s", err)
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item, _ := v.(map[string]interface{})
			id := rdsGetField(item, "DeploymentSetId")
			name := rdsGetField(item, "DeploymentSetName")
			if !rdsSweepNameOk(name) && !rdsSweepNameOk(id) {
				continue
			}
			if e := rdsDeletePost(c, "DeleteRCDeploymentSet", map[string]interface{}{
				"DeploymentSetId": id,
				"RegionId":        c.RegionId,
			}, []string{"DeploymentSet.NotFound", "InvalidParameter.InstanceId"}); e != nil {
				log.Printf("[ERROR] Failed to delete rds custom_deployment_set (%s): %s", id, e)
				continue
			}
			log.Printf("[INFO] Delete rds custom_deployment_set success: %s (%s)", id, name)
		}
		if len(result) < 50 {
			break
		}
		pageNum++
	}
	return nil
}

// testSweepRdsCustomDisks: DescribeRCDisks (RpcGet, RegionId, PageNumber/PageSize
// max 100) -> $.Disks[*]; delete DeleteRCDisk (DiskId + RegionId). Prepaid disks
// cannot be released (BSS refund — DEFER); skip DiskChargeType=="Prepaid".
// Runs AFTER rds_custom_disk_attachment so DetachRCDisk has already detached
// attached disks (DeleteRCDisk fails on In_use state).
func testSweepRdsCustomDisks(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	c := rawClient.(*connectivity.AliyunClient)
	pageNum := 1
	for {
		query := map[string]interface{}{
			"RegionId":   c.RegionId,
			"PageNumber": pageNum,
			"PageSize":   100,
		}
		response, err := rdsRpcGet(c, "DescribeRCDisks", query)
		if err != nil {
			log.Printf("[ERROR] rds DescribeRCDisks: %s", err)
			return nil
		}
		resp, err := jsonpath.Get("$.Disks", response)
		if err != nil {
			log.Printf("[ERROR] rds DescribeRCDisks jsonpath: %s", err)
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item, _ := v.(map[string]interface{})
			id := rdsGetField(item, "DiskId")
			name := rdsGetField(item, "DiskName")
			if !rdsSweepNameOk(name) && !rdsSweepNameOk(id) {
				continue
			}
			if rdsGetField(item, "DiskChargeType") == "Prepaid" {
				log.Printf("[INFO] Skipping rds custom_disk (Prepaid, BSS refund DEFER): %s", id)
				continue
			}
			if e := rdsDeletePost(c, "DeleteRCDisk", map[string]interface{}{
				"DiskId":   id,
				"RegionId": c.RegionId,
			}, []string{"InvalidDiskId.NotFound", "IncorrectDiskStatus"}); e != nil {
				log.Printf("[ERROR] Failed to delete rds custom_disk (%s): %s", id, e)
				continue
			}
			log.Printf("[INFO] Delete rds custom_disk success: %s (%s)", id, name)
		}
		if len(result) < 100 {
			break
		}
		pageNum++
	}
	return nil
}

// testSweepRdsCustomDiskAttachments: DescribeRCDisks (RpcGet, RegionId, no
// InstanceId filter) -> $.Disks[*]; for every disk in In_use state (attached),
// DetachRCDisk (DiskId + InstanceId + RegionId). Runs BEFORE
// rds_custom_disk so disks are detached before deletion attempt.
func testSweepRdsCustomDiskAttachments(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	c := rawClient.(*connectivity.AliyunClient)
	pageNum := 1
	for {
		query := map[string]interface{}{
			"RegionId":   c.RegionId,
			"PageNumber": pageNum,
			"PageSize":   100,
			"Status":     "In_use",
		}
		response, err := rdsRpcGet(c, "DescribeRCDisks", query)
		if err != nil {
			log.Printf("[ERROR] rds DescribeRCDisks (attachment): %s", err)
			return nil
		}
		resp, err := jsonpath.Get("$.Disks", response)
		if err != nil {
			log.Printf("[ERROR] rds DescribeRCDisks jsonpath (attachment): %s", err)
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item, _ := v.(map[string]interface{})
			diskID := rdsGetField(item, "DiskId")
			diskName := rdsGetField(item, "DiskName")
			instanceID := rdsGetField(item, "InstanceId")
			if diskID == "" || instanceID == "" {
				continue
			}
			if !rdsSweepNameOk(diskName) && !rdsSweepNameOk(diskID) {
				continue
			}
			if e := rdsDeletePost(c, "DetachRCDisk", map[string]interface{}{
				"DiskId":     diskID,
				"InstanceId": instanceID,
				"RegionId":   c.RegionId,
			}, []string{"InvalidDiskId.NotFound", "InvalidInstanceId.NotFound", "IncorrectDiskStatus"}); e != nil {
				log.Printf("[ERROR] Failed to detach rds custom_disk_attachment (%s/%s): %s", diskID, instanceID, e)
				continue
			}
			log.Printf("[INFO] Detach rds custom_disk_attachment success: %s/%s", diskID, instanceID)
		}
		if len(result) < 100 {
			break
		}
		pageNum++
	}
	return nil
}

// testSweepRdsParameterGroups: DescribeParameterGroups (RpcPost, RegionId) ->
// $.ParameterGroups.ParameterGroup; delete DeleteParameterGroup
// (ParameterGroupId + RegionId). Only ParameterGroupType==1 (user-created) is
// swept; system templates (type 0) and auto-backup templates (type 2) are
// preserved.
func testSweepRdsParameterGroups(region string) error {
	raw_client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	c := raw_client.(*connectivity.AliyunClient)
	response, err := rdsRpcPost(c, "DescribeParameterGroups", map[string]interface{}{
		"RegionId": c.RegionId,
	})
	if err != nil {
		log.Printf("[ERROR] rds DescribeParameterGroups: %s", err)
		return nil
	}
	resp, err := jsonpath.Get("$.ParameterGroups.ParameterGroup", response)
	if err != nil {
		log.Printf("[ERROR] rds DescribeParameterGroups jsonpath: %s", err)
		return nil
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item, _ := v.(map[string]interface{})
		id := rdsGetField(item, "ParameterGroupId")
		name := rdsGetField(item, "ParameterGroupName")
		// ParameterGroupType: 0 = system, 1 = user-created, 2 = auto-backup.
		// Only sweep user-created (type 1); others are server-managed.
		if rdsGetField(item, "ParameterGroupType") != "1" {
			log.Printf("[INFO] Skipping rds parameter_group (system/auto-backup): %s (%s)", id, name)
			continue
		}
		if !rdsSweepNameOk(name) && !rdsSweepNameOk(id) {
			continue
		}
		if e := rdsDeletePost(c, "DeleteParameterGroup", map[string]interface{}{
			"ParameterGroupId": id,
			"RegionId":         c.RegionId,
		}, []string{"ParamGroupsNotExist", "ParamGroupsNotExistError"}); e != nil {
			log.Printf("[ERROR] Failed to delete rds parameter_group (%s): %s", id, e)
			continue
		}
		log.Printf("[INFO] Delete rds parameter_group success: %s (%s)", id, name)
	}
	return nil
}

// testSweepRdsWhitelistTemplates: DescribeAllWhitelistTemplate (RpcPost,
// RegionId, MaxRecordsPerPage=1000, PageNumbers=1) -> $.Data.Templates; "delete"
// via ModifyWhitelistTemplate with IpWhitelist="" (the only operation available
// — ModifyWhitelistTemplate handles create/modify/delete per the API). The
// template is cleared (IpWhitelist emptied) rather than truly deleted — this
// matches the resource's Delete behavior.
func testSweepRdsWhitelistTemplates(region string) error {
	raw_client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	c := raw_client.(*connectivity.AliyunClient)
	response, err := rdsRpcPost(c, "DescribeAllWhitelistTemplate", map[string]interface{}{
		"RegionId":          c.RegionId,
		"MaxRecordsPerPage": 1000,
		"PageNumbers":       1,
	})
	if err != nil {
		log.Printf("[ERROR] rds DescribeAllWhitelistTemplate: %s", err)
		return nil
	}
	resp, err := jsonpath.Get("$.Data.Templates", response)
	if err != nil {
		log.Printf("[ERROR] rds DescribeAllWhitelistTemplate jsonpath: %s", err)
		return nil
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item, _ := v.(map[string]interface{})
		id := rdsGetField(item, "TemplateId")
		name := rdsGetField(item, "TemplateName")
		if !rdsSweepNameOk(name) && !rdsSweepNameOk(id) {
			continue
		}
		// Clear the IpWhitelist — there is no DeleteWhitelistTemplate API.
		if e := rdsDeletePost(c, "ModifyWhitelistTemplate", map[string]interface{}{
			"RegionId":     c.RegionId,
			"TemplateId":   id,
			"TemplateName": name,
			"IpWhitelist":  "",
		}, OperationDeniedDBStatus); e != nil {
			log.Printf("[ERROR] Failed to clear rds whitelist_template (%s): %s", id, e)
			continue
		}
		log.Printf("[INFO] Clear rds whitelist_template success: %s (%s)", id, name)
	}
	return nil
}

// testSweepRdsAiInstances: DescribeAppInstances (RdsAi 2025-05-07, RegionId,
// PageNumber/PageSize max 50) -> $.Instances[*] (top-level array — unusual);
// delete DeleteAppInstance (InstanceName + RegionId). RdsAi instances are
// pay-as-you-go (no BSS refund), so DeleteAppInstance is a plain delete.
func testSweepRdsAiInstances(region string) error {
	raw_client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	c := raw_client.(*connectivity.AliyunClient)
	pageNum := 1
	for {
		response, err := rdsAiPost(c, "DescribeAppInstances", map[string]interface{}{
			"RegionId":   c.RegionId,
			"PageNumber": pageNum,
			"PageSize":   50,
		})
		if err != nil {
			log.Printf("[ERROR] rdsai DescribeAppInstances: %s", err)
			return nil
		}
		resp, err := jsonpath.Get("$.Instances", response)
		if err != nil {
			log.Printf("[ERROR] rdsai DescribeAppInstances jsonpath: %s", err)
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item, _ := v.(map[string]interface{})
			id := rdsGetField(item, "InstanceName")
			name := rdsGetField(item, "AppName")
			if !rdsSweepNameOk(name) && !rdsSweepNameOk(id) {
				continue
			}
			// DeleteAppInstance is a plain PAYG delete — NOT a BSS refund/退订.
			var e error
			wait := incrementalWait(3*time.Second, 5*time.Second)
			e = retry.Retry(5*time.Minute, func() *retry.RetryError {
				_, e = c.RpcPost("RdsAi", "2025-05-07", "DeleteAppInstance", nil, map[string]interface{}{
					"InstanceName": id,
					"RegionId":     c.RegionId,
				}, true)
				if e != nil {
					if NotFoundError(e) || IsExpectedErrors(e, []string{"InvalidParameter", "MasterInstanceNotExist"}) {
						return nil
					}
					if NeedRetry(e) {
						wait()
						return retry.RetryableError(e)
					}
					return retry.NonRetryableError(e)
				}
				return nil
			})
			if e != nil {
				log.Printf("[ERROR] Failed to delete rds_ai_instance (%s): %s", id, e)
				continue
			}
			log.Printf("[INFO] Delete rds_ai_instance success: %s (%s)", id, name)
		}
		if len(result) < 50 {
			break
		}
		pageNum++
	}
	return nil
}

// testSweepRedisGlobalSecurityIPGroups: DescribeGlobalSecurityIPGroup (R-kvstore,
// RegionId) -> $.GlobalSecurityIPGroup[*]; delete DeleteGlobalSecurityIPGroup
// (GlobalSecurityGroupId + GlobalIgName + RegionId).
func testSweepRedisGlobalSecurityIPGroups(region string) error {
	raw_client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	c := raw_client.(*connectivity.AliyunClient)
	response, err := rdsKVPost(c, "DescribeGlobalSecurityIPGroup", map[string]interface{}{
		"RegionId": c.RegionId,
	})
	if err != nil {
		log.Printf("[ERROR] redis DescribeGlobalSecurityIPGroup: %s", err)
		return nil
	}
	resp, err := jsonpath.Get("$.GlobalSecurityIPGroup", response)
	if err != nil {
		log.Printf("[ERROR] redis DescribeGlobalSecurityIPGroup jsonpath: %s", err)
		return nil
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item, _ := v.(map[string]interface{})
		id := rdsGetField(item, "GlobalSecurityGroupId")
		name := rdsGetField(item, "GlobalIgName")
		if !rdsSweepNameOk(name) && !rdsSweepNameOk(id) {
			continue
		}
		if e := rdsKVDelete(c, "DeleteGlobalSecurityIPGroup", map[string]interface{}{
			"GlobalSecurityGroupId": id,
			"GlobalIgName":          name,
			"RegionId":              c.RegionId,
		}, []string{"RequiredParam.NotFound"}); e != nil {
			log.Printf("[ERROR] Failed to delete redis global_security_ip_group (%s): %s", id, e)
			continue
		}
		log.Printf("[INFO] Delete redis global_security_ip_group success: %s (%s)", id, name)
	}
	return nil
}

// ---- instance-scoped sweepers (walk Running DB instances first) ----

// testSweepRdsAccounts: walk Running DBInstances; per instance call
// DescribeAccounts (RpcPost, AccountName optional) -> $.Accounts.DBInstanceAccount;
// delete DeleteAccount (AccountName + DBInstanceId + SourceIp). Accounts are
// filtered by AccountName (tf-prefix gate OR sweepAll), so accounts on shared
// long-lived instances are still swept when tf-prefixed.
func testSweepRdsAccounts(region string) error {
	raw_client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	c := raw_client.(*connectivity.AliyunClient)
	_, instances, err := rdsListRunningInstanceIDs(c)
	if err != nil {
		log.Printf("[ERROR] rds list instances for account sweep: %s", err)
		return nil
	}
	for _, inst := range instances {
		instID := rdsGetField(inst, "DBInstanceId")
		instName := rdsGetField(inst, "DBInstanceDescription")
		// Instance-level gate: skip non-test instances entirely (their account
		// lists are not queried, so this also avoids sweeping accounts on
		// NODELETING instances). In SWEEP_ALL mode the gate is wide-open.
		if !rdsInstanceSweepable(instName) {
			continue
		}
		response, err := rdsRpcPost(c, "DescribeAccounts", map[string]interface{}{
			"RegionId":     c.RegionId,
			"DBInstanceId": instID,
			"SourceIp":     c.SourceIp,
		})
		if err != nil {
			log.Printf("[ERROR] DescribeAccounts (%s): %s", instID, err)
			continue
		}
		resp, err := jsonpath.Get("$.Accounts.DBInstanceAccount", response)
		if err != nil {
			log.Printf("[ERROR] DescribeAccounts (%s) jsonpath: %s", instID, err)
			continue
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item, _ := v.(map[string]interface{})
			accountName := rdsGetField(item, "AccountName")
			if accountName == "" {
				continue
			}
			if !rdsSweepNameOk(accountName) {
				continue
			}
			if e := rdsDeletePost(c, "DeleteAccount", map[string]interface{}{
				"AccountName":  accountName,
				"DBInstanceId": instID,
				"SourceIp":     c.SourceIp,
			}, []string{"InvalidDBInstanceId.NotFound", "InvalidDBInstanceName.NotFound", "InvalidAccountName.NotFound"}); e != nil {
				log.Printf("[ERROR] Failed to delete rds account (%s/%s): %s", instID, accountName, e)
				continue
			}
			log.Printf("[INFO] Delete rds account success: %s/%s", instID, accountName)
		}
	}
	return nil
}

// testSweepRdsBackups: walk Running DBInstances; per instance call
// DescribeBackups (RpcPost, BackupMode=Manual, StartTime/EndTime wide window)
// -> $.Items.Backup; delete DeleteBackup (BackupId + DBInstanceId). Only Manual
// backups are swept (system automated backups are server-managed and not Terraform
// resources). Backups have no name field, so the per-instance gate decides which
// instances' backups to delete (tf-prefixed instances OR sweepAll fleet-wipe).
func testSweepRdsBackups(region string) error {
	raw_client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	c := raw_client.(*connectivity.AliyunClient)
	_, instances, err := rdsListRunningInstanceIDs(c)
	if err != nil {
		log.Printf("[ERROR] rds list instances for backup sweep: %s", err)
		return nil
	}
	for _, inst := range instances {
		instID := rdsGetField(inst, "DBInstanceId")
		instName := rdsGetField(inst, "DBInstanceDescription")
		if !rdsInstanceSweepable(instName) {
			continue
		}
		pageNum := 1
		for {
			response, err := rdsRpcPost(c, "DescribeBackups", map[string]interface{}{
				"RegionId":     c.RegionId,
				"DBInstanceId": instID,
				"BackupMode":   "Manual",
				"StartTime":    "2010-01-01T00:00Z",
				"EndTime":      "2050-01-01T00:00Z",
				"PageSize":     100,
				"PageNumber":   pageNum,
			})
			if err != nil {
				log.Printf("[ERROR] DescribeBackups (%s): %s", instID, err)
				break
			}
			resp, err := jsonpath.Get("$.Items.Backup", response)
			if err != nil {
				log.Printf("[ERROR] DescribeBackups (%s) jsonpath: %s", instID, err)
				break
			}
			result, _ := resp.([]interface{})
			for _, v := range result {
				item, _ := v.(map[string]interface{})
				backupID := rdsGetField(item, "BackupId")
				if backupID == "" {
					continue
				}
				// Only delete backups whose StoreStatus is Enabled (per API doc).
				if storeStatus := rdsGetField(item, "StoreStatus"); storeStatus != "" && storeStatus != "Enabled" {
					log.Printf("[INFO] Skipping rds backup (StoreStatus=%s): %s/%s", storeStatus, instID, backupID)
					continue
				}
				if e := rdsDeletePost(c, "DeleteBackup", map[string]interface{}{
					"BackupId":     backupID,
					"DBInstanceId": instID,
				}, []string{"InvalidBackupId.NotFound", "BackupNotFound"}); e != nil {
					log.Printf("[ERROR] Failed to delete rds backup (%s/%s): %s", instID, backupID, e)
					continue
				}
				log.Printf("[INFO] Delete rds backup success: %s/%s", instID, backupID)
			}
			if len(result) < 100 {
				break
			}
			pageNum++
		}
	}
	return nil
}

// testSweepRdsDBInstanceEndpoints: walk Running DBInstances; per instance call
// DescribeDBInstanceEndpoints (RpcPost, DBInstanceEndpointId omitted -> all) ->
// $.Data.DBInstanceEndpoints.DBInstanceEndpoint; delete DeleteDBInstanceEndpoint
// (DBInstanceId + DBInstanceEndpointId). Primary endpoints are skipped (cannot
// be deleted); endpoints filtered by their ConnectionString prefix.
// Runs AFTER rds_db_instance_endpoint_address so public addresses are removed first.
func testSweepRdsDBInstanceEndpoints(region string) error {
	raw_client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	c := raw_client.(*connectivity.AliyunClient)
	_, instances, err := rdsListRunningInstanceIDs(c)
	if err != nil {
		log.Printf("[ERROR] rds list instances for endpoint sweep: %s", err)
		return nil
	}
	for _, inst := range instances {
		instID := rdsGetField(inst, "DBInstanceId")
		instName := rdsGetField(inst, "DBInstanceDescription")
		if !rdsInstanceSweepable(instName) {
			continue
		}
		response, err := rdsRpcPost(c, "DescribeDBInstanceEndpoints", map[string]interface{}{
			"RegionId":     c.RegionId,
			"DBInstanceId": instID,
			"SourceIp":     c.SourceIp,
		})
		if err != nil {
			log.Printf("[ERROR] DescribeDBInstanceEndpoints (%s): %s", instID, err)
			continue
		}
		// $.Data.DBInstanceEndpoints is a map containing the DBInstanceEndpoint array.
		nodeMap, err := jsonpath.Get("$.Data.DBInstanceEndpoints", response)
		if err != nil {
			log.Printf("[ERROR] DescribeDBInstanceEndpoints (%s) jsonpath: %s", instID, err)
			continue
		}
		m, ok := nodeMap.(map[string]interface{})
		if !ok {
			continue
		}
		result, _ := m["DBInstanceEndpoint"].([]interface{})
		for _, v := range result {
			item, _ := v.(map[string]interface{})
			endpointID := rdsGetField(item, "DBInstanceEndpointId", "EndpointId")
			if endpointID == "" {
				continue
			}
			// Primary endpoints cannot be deleted; skip them.
			if endpointType := rdsGetField(item, "EndpointType"); endpointType == "Primary" {
				continue
			}
			// Filter on ConnectionString (the endpoint's primary human-readable name).
			connStr := ""
			if addrs, ok := item["AddressItems"].(map[string]interface{}); ok {
				if items, ok := addrs["AddressItem"].([]interface{}); ok && len(items) > 0 {
					if first, ok := items[0].(map[string]interface{}); ok {
						connStr = rdsGetField(first, "ConnectionString")
					}
				}
			}
			if !rdsSweepNameOk(connStr) && !rdsSweepNameOk(endpointID) {
				continue
			}
			if e := rdsDeletePost(c, "DeleteDBInstanceEndpoint", map[string]interface{}{
				"RegionId":             c.RegionId,
				"DBInstanceId":         instID,
				"DBInstanceEndpointId": endpointID,
				"SourceIp":             c.SourceIp,
			}, []string{"InvalidDBInstanceEndpointId.NotFound", "InvalidDBInstance.NotFound"}); e != nil {
				log.Printf("[ERROR] Failed to delete rds db_instance_endpoint (%s/%s): %s", instID, endpointID, e)
				continue
			}
			log.Printf("[INFO] Delete rds db_instance_endpoint success: %s/%s", instID, endpointID)
		}
	}
	return nil
}

// testSweepRdsDBInstanceEndpointAddresses: walk Running DBInstances; per
// instance call DescribeDBInstanceEndpoints; for each endpoint iterate
// AddressItems.AddressItem looking for IpType=="Public" entries whose
// ConnectionString matches the tf-prefix gate; delete via
// DeleteDBInstanceEndpointAddress (DBInstanceId + DBInstanceEndpointId +
// ConnectionString + SourceIp). Runs BEFORE rds_db_instance_endpoint so the
// public address is removed before the endpoint itself is deleted.
func testSweepRdsDBInstanceEndpointAddresses(region string) error {
	raw_client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	c := raw_client.(*connectivity.AliyunClient)
	_, instances, err := rdsListRunningInstanceIDs(c)
	if err != nil {
		log.Printf("[ERROR] rds list instances for endpoint_address sweep: %s", err)
		return nil
	}
	for _, inst := range instances {
		instID := rdsGetField(inst, "DBInstanceId")
		instName := rdsGetField(inst, "DBInstanceDescription")
		if !rdsInstanceSweepable(instName) {
			continue
		}
		response, err := rdsRpcPost(c, "DescribeDBInstanceEndpoints", map[string]interface{}{
			"RegionId":     c.RegionId,
			"DBInstanceId": instID,
			"SourceIp":     c.SourceIp,
		})
		if err != nil {
			log.Printf("[ERROR] DescribeDBInstanceEndpoints (%s) for address sweep: %s", instID, err)
			continue
		}
		nodeMap, err := jsonpath.Get("$.Data.DBInstanceEndpoints", response)
		if err != nil {
			continue
		}
		m, ok := nodeMap.(map[string]interface{})
		if !ok {
			continue
		}
		result, _ := m["DBInstanceEndpoint"].([]interface{})
		for _, v := range result {
			item, _ := v.(map[string]interface{})
			endpointID := rdsGetField(item, "DBInstanceEndpointId", "EndpointId")
			if endpointID == "" {
				continue
			}
			addrs, ok := item["AddressItems"].(map[string]interface{})
			if !ok {
				continue
			}
			addrItems, _ := addrs["AddressItem"].([]interface{})
			for _, a := range addrItems {
				addr, _ := a.(map[string]interface{})
				if rdsGetField(addr, "IpType") != "Public" {
					continue
				}
				connStr := rdsGetField(addr, "ConnectionString")
				if connStr == "" {
					continue
				}
				if !rdsSweepNameOk(connStr) {
					continue
				}
				if e := rdsDeletePost(c, "DeleteDBInstanceEndpointAddress", map[string]interface{}{
					"RegionId":             c.RegionId,
					"DBInstanceId":         instID,
					"DBInstanceEndpointId": endpointID,
					"ConnectionString":     connStr,
					"SourceIp":             c.SourceIp,
				}, []string{"InvalidDBInstanceEndpointId.NotFound", "InvalidDBInstance.NotFound"}); e != nil {
					log.Printf("[ERROR] Failed to delete rds db_instance_endpoint_address (%s/%s/%s): %s", instID, endpointID, connStr, e)
					continue
				}
				log.Printf("[INFO] Delete rds db_instance_endpoint_address success: %s/%s/%s", instID, endpointID, connStr)
			}
		}
	}
	return nil
}

// testSweepKVStoreAccounts: walk every KVStore / Redis / Tair instance; per
// instance call DescribeAccounts (R-kvstore 2015-01-01) -> $.Accounts.Account[*];
// delete DeleteAccount (AccountName + InstanceId + RegionId). Accounts are
// filtered by AccountName (tf-prefix gate OR sweepAll), so accounts on shared
// instances are still swept when tf-prefixed.
func testSweepKVStoreAccounts(region string) error {
	raw_client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	c := raw_client.(*connectivity.AliyunClient)
	instanceIDs, err := rdsListKVStoreInstanceIDs(c)
	if err != nil {
		log.Printf("[ERROR] kvstore list instances for account sweep: %s", err)
		return nil
	}
	for _, instID := range instanceIDs {
		response, err := rdsKVPost(c, "DescribeAccounts", map[string]interface{}{
			"InstanceId": instID,
			"RegionId":   c.RegionId,
		})
		if err != nil {
			if NotFoundError(err) || IsExpectedErrors(err, []string{"InvalidInstanceId.NotFound"}) {
				continue
			}
			log.Printf("[ERROR] DescribeAccounts (%s): %s", instID, err)
			continue
		}
		resp, err := jsonpath.Get("$.Accounts.Account", response)
		if err != nil {
			continue
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item, _ := v.(map[string]interface{})
			accountName := rdsGetField(item, "AccountName")
			if accountName == "" {
				continue
			}
			if !rdsSweepNameOk(accountName) {
				continue
			}
			if e := rdsKVDelete(c, "DeleteAccount", map[string]interface{}{
				"AccountName": accountName,
				"InstanceId":  instID,
				"RegionId":    c.RegionId,
			}, []string{"InvalidInstanceId.NotFound", "InvalidAccountName.NotFound"}); e != nil {
				log.Printf("[ERROR] Failed to delete kvstore account (%s/%s): %s", instID, accountName, e)
				continue
			}
			log.Printf("[INFO] Delete kvstore account success: %s/%s", instID, accountName)
		}
	}
	return nil
}

// testSweepRedisBackups: walk every KVStore / Redis / Tair instance; per
// instance call DescribeBackups (R-kvstore 2015-01-01, StartTime/EndTime wide
// window) -> $.Backups.Backup[*]; delete DeleteBackup (BackupId + InstanceId).
// Backups have no name field, so the per-instance gate decides which instances'
// backups to delete. NOTE: Tair instances are also walked because the
// alicloud_redis_backup resource can be created against any R-kvstore instance.
func testSweepRedisBackups(region string) error {
	raw_client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	c := raw_client.(*connectivity.AliyunClient)
	instanceIDs, err := rdsListKVStoreInstanceIDs(c)
	if err != nil {
		log.Printf("[ERROR] redis list instances for backup sweep: %s", err)
		return nil
	}
	for _, instID := range instanceIDs {
		pageNum := 1
		for {
			response, err := rdsKVPost(c, "DescribeBackups", map[string]interface{}{
				"InstanceId": instID,
				"RegionId":   c.RegionId,
				"StartTime":  "2010-01-01T00:00Z",
				"EndTime":    "2050-01-01T00:00Z",
				"PageSize":   100,
				"PageNumber": pageNum,
			})
			if err != nil {
				if NotFoundError(err) || IsExpectedErrors(err, []string{"InvalidInstanceId.NotFound"}) {
					break
				}
				log.Printf("[ERROR] DescribeBackups (%s): %s", instID, err)
				break
			}
			resp, err := jsonpath.Get("$.Backups.Backup", response)
			if err != nil {
				break
			}
			result, _ := resp.([]interface{})
			for _, v := range result {
				item, _ := v.(map[string]interface{})
				backupID := rdsGetField(item, "BackupId")
				if backupID == "" {
					continue
				}
				if e := rdsKVDelete(c, "DeleteBackup", map[string]interface{}{
					"BackupId":   backupID,
					"InstanceId": instID,
				}, []string{"InvalidBackupId.NotFound", "BackupNotFound"}); e != nil {
					log.Printf("[ERROR] Failed to delete redis backup (%s/%s): %s", instID, backupID, e)
					continue
				}
				log.Printf("[INFO] Delete redis backup success: %s/%s", instID, backupID)
			}
			if len(result) < 100 {
				break
			}
			pageNum++
		}
	}
	return nil
}
