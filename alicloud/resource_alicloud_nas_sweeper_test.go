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

// NAS (Apsara File Storage NAS, product "NAS" / API version "2017-06-26")
// gap sweepers. NAS is RPC-style; the 2017-06-26 version is the only one the
// provider uses (service_alicloud_nas_v2.go still calls 2017-06-26).
//
// Resources enumerated by the sweepers below (each has BOTH an enumerable
// list API and a destructive delete API):
//
//   region-scoped (no file-system parent):
//     alicloud_nas_snapshot        DescribeSnapshots            -> DeleteSnapshot
//                                  (extreme-type only; snapshots persist after
//                                  their source file system is removed, so they
//                                  are swept directly by SnapshotName)
//
//   file-system-scoped (walk DescribeFileSystems first; every child is gated
//   on the parent file system Description so children of NODELETING or non-test
//   file systems are preserved):
//     alicloud_nas_mount_target        DescribeMountTargets      -> DeleteMountTarget
//     alicloud_nas_fileset             DescribeFilesets          -> DeleteFileset
//     alicloud_nas_data_flow           DescribeDataFlows         -> DeleteDataFlow
//     alicloud_nas_lifecycle_policy    DescribeLifecyclePolicies -> DeleteLifecyclePolicy
//     alicloud_nas_access_point        DescribeAccessPoints      -> DeleteAccessPoint
//     alicloud_nas_protocol_service    DescribeProtocolService   -> DeleteProtocolService
//                                       (depends on alicloud_nas_protocol_mount_target so
//                                       the export dirs are removed before the service)
//     alicloud_nas_protocol_mount_target DescribeProtocolMountTarget -> DeleteProtocolMountTarget
//                                       (delete needs ExportId + FileSystemId +
//                                       ProtocolServiceId, all returned by the list)
//
//   access-group-scoped (walk DescribeAccessGroups first; every rule is gated
//   on the parent AccessGroupName so rules of NODELETING or non-test groups are
//   preserved):
//     alicloud_nas_access_rule         DescribeAccessRules       -> DeleteAccessRule
//                                       (delete needs AccessGroupName + AccessRuleId +
//                                       FileSystemType, all returned by the group list)
//
// DEFERRED (no sweeper written — documented here with the reason):
//   alicloud_nas_recycle_bin      only a per-file-system singleton feature toggle
//                                 (EnableRecycleBin / DisableAndCleanRecycleBin);
//                                 there is NO enumerable list API, only the
//                                 Get-by-id DescribeFileSystemGetRecycleBinAttribute.
//                                 The "delete" is a disable+clean toggle, not a
//                                 destroy, and the resource is removed with its
//                                 parent file system.
//   alicloud_nas_smb_acl_attachment  per-file-system singleton feature toggle
//                                 (EnableSmbAcl / DisableSmbAcl); NO enumerable
//                                 list API, only Get-by-id
//                                 DescribeFileSystemDescribeSmbAcl. Removed with
//                                 its parent file system.
//   alicloud_nas_log_analysis     although DescribeLogAnalysis is region-listable
//                                 and DeleteLogAnalysis exists, log analysis is a
//                                 per-file-system singleton feature config (one
//                                 entry per file system, id = FileSystemId) that
//                                 is removed when the parent file system is
//                                 deleted. A SWEEP_ALL pass would disable log
//                                 analysis on every non-NODELETING file system in
//                                 the account, which risks disrupting test
//                                 workloads that rely on it for no orphan-cleanup
//                                 value (the file_system sweeper already removes
//                                 these with the FS).
//
// PRE-EXISTING inline sweepers (registered in their own _test.go files — NOT
// re-registered here to avoid duplicate-name init panics):
//   alicloud_nas_file_system          alicloud/resource_alicloud_nas_file_system_test.go
//                                     (also deletes mount targets inline)
//   alicloud_nas_access_group         alicloud/resource_alicloud_nas_access_group_test.go
//                                     (depends on alicloud_nas_file_system)
//   alicloud_nas_auto_snapshot_policy alicloud/resource_alicloud_nas_auto_snapshot_policy_test.go

const nasSweepAPIVersion = "2017-06-26"

// nasSweepPrefixes mark a resource (by name/description) as test-created.
var nasSweepPrefixes = []string{
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

// nasGetField returns the first non-empty string value of item[name] across the
// given name variants. NAS response field casing is consistent (PascalCase) but
// a few items echo back camelCase aliases, so callers pass both variants.
func nasGetField(item map[string]interface{}, names ...string) string {
	for _, name := range names {
		if v, ok := item[name]; ok {
			if s := fmt.Sprint(v); s != "" && s != "<nil>" {
				return s
			}
		}
	}
	return ""
}

// nasSweepNameOk is the fleet-wipe gate: never touch NODELETING names; in
// normal (non-SWEEP_ALL) mode only sweep tf-prefixed names. SWEEP_ALL bypasses
// the prefix filter so the whole test account (1511928242963727, cn-beijing)
// is fleet-wiped. Returns true when the resource should be swept.
func nasSweepNameOk(name string) bool {
	if name == "" {
		return false
	}
	if protectedSweepVpcName(name) {
		log.Printf("[INFO] Skipping nas resource (NODELETING): %s", name)
		return false
	}
	if !sweepAll() {
		skip := true
		for _, prefix := range nasSweepPrefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping nas resource (no prefix): %s", name)
			return false
		}
	}
	return true
}

// nasSweepChildOk gates a child resource on its parent file system Description
// (the same field the file_system sweeper uses) so children of NODELETING or
// non-test file systems are preserved. The child's own name is also checked for
// an explicit do-not-delete marker.
func nasSweepChildOk(parentDesc, childName string) bool {
	if protectedSweepVpcName(parentDesc) || protectedSweepVpcName(childName) {
		log.Printf("[INFO] Skipping nas child (NODELETING): parent=%q child=%q", parentDesc, childName)
		return false
	}
	if !sweepAll() {
		for _, prefix := range nasSweepPrefixes {
			if strings.HasPrefix(strings.ToLower(parentDesc), strings.ToLower(prefix)) {
				return true
			}
		}
		log.Printf("[INFO] Skipping nas child (parent no prefix): parent=%q child=%q", parentDesc, childName)
		return false
	}
	return true
}

// nasRpcPost wraps a NAS RPC POST call with retry. NeedRetry errors are retried
// for 5 minutes with incremental backoff; other errors surface immediately.
// Params go in the BODY slot (query=nil), matching the NAS data sources and the
// resource Create/Delete functions (RpcPost signature: product, version, action,
// query, body, autoRetry).
func nasRpcPost(client *connectivity.AliyunClient, action string, query, body map[string]interface{}, autoRetry bool) (map[string]interface{}, error) {
	var response map[string]interface{}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err := retry.Retry(5*time.Minute, func() *retry.RetryError {
		resp, err := client.RpcPost("NAS", nasSweepAPIVersion, action, query, body, autoRetry)
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

// nasDeleteOne performs a best-effort delete. NotFoundError and the given
// expected error codes are treated as success (resource already gone); other
// errors are logged and the sweep continues. NEVER returns an error that would
// abort the sweep.
func nasDeleteOne(client *connectivity.AliyunClient, action string, query, body map[string]interface{}, id string, okErrs []string) {
	_, err := nasRpcPost(client, action, query, body, false)
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

// nasListFileSystems enumerates every NAS file system in the sweep region.
func nasListFileSystems(client *connectivity.AliyunClient) ([]map[string]interface{}, error) {
	action := "DescribeFileSystems"
	body := map[string]interface{}{
		"RegionId":   client.Region,
		"PageSize":   PageSizeLarge,
		"PageNumber": 1,
	}
	var items []map[string]interface{}
	for {
		response, err := nasRpcPost(client, action, nil, body, true)
		if err != nil {
			return nil, fmt.Errorf("DescribeFileSystems: %s", err)
		}
		resp, jpErr := jsonpath.Get("$.FileSystems.FileSystem", response)
		if jpErr != nil {
			return nil, fmt.Errorf("DescribeFileSystems parse: %s", jpErr)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			if item, ok := v.(map[string]interface{}); ok {
				items = append(items, item)
			}
		}
		if len(result) < PageSizeLarge {
			return items, nil
		}
		body["PageNumber"] = body["PageNumber"].(int) + 1
	}
}

// nasListAccessGroups enumerates every NAS access group in the sweep region.
// Each item carries AccessGroupName + FileSystemType, both needed to list and
// delete its access rules.
func nasListAccessGroups(client *connectivity.AliyunClient) ([]map[string]interface{}, error) {
	action := "DescribeAccessGroups"
	body := map[string]interface{}{
		"PageSize":   PageSizeLarge,
		"PageNumber": 1,
	}
	var items []map[string]interface{}
	for {
		response, err := nasRpcPost(client, action, nil, body, true)
		if err != nil {
			return nil, fmt.Errorf("DescribeAccessGroups: %s", err)
		}
		resp, jpErr := jsonpath.Get("$.AccessGroups.AccessGroup", response)
		if jpErr != nil {
			return nil, fmt.Errorf("DescribeAccessGroups parse: %s", jpErr)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			if item, ok := v.(map[string]interface{}); ok {
				items = append(items, item)
			}
		}
		if len(result) < PageSizeLarge {
			return items, nil
		}
		body["PageNumber"] = body["PageNumber"].(int) + 1
	}
}

func init() {
	// Region-scoped.
	resource.AddTestSweepers("alicloud_nas_snapshot", &resource.Sweeper{
		Name: "alicloud_nas_snapshot",
		F:    testSweepNasSnapshots,
	})
	// File-system-scoped children.
	resource.AddTestSweepers("alicloud_nas_mount_target", &resource.Sweeper{
		Name: "alicloud_nas_mount_target",
		F:    testSweepNasMountTargets,
	})
	resource.AddTestSweepers("alicloud_nas_fileset", &resource.Sweeper{
		Name: "alicloud_nas_fileset",
		F:    testSweepNasFilesets,
	})
	resource.AddTestSweepers("alicloud_nas_data_flow", &resource.Sweeper{
		Name: "alicloud_nas_data_flow",
		F:    testSweepNasDataFlows,
	})
	resource.AddTestSweepers("alicloud_nas_lifecycle_policy", &resource.Sweeper{
		Name: "alicloud_nas_lifecycle_policy",
		F:    testSweepNasLifecyclePolicies,
	})
	resource.AddTestSweepers("alicloud_nas_access_point", &resource.Sweeper{
		Name: "alicloud_nas_access_point",
		F:    testSweepNasAccessPoints,
	})
	// protocol_service depends on protocol_mount_target so the export dirs are
	// removed BEFORE the protocol service (children-first DAG).
	resource.AddTestSweepers("alicloud_nas_protocol_service", &resource.Sweeper{
		Name: "alicloud_nas_protocol_service",
		F:    testSweepNasProtocolServices,
		Dependencies: []string{
			"alicloud_nas_protocol_mount_target",
		},
	})
	resource.AddTestSweepers("alicloud_nas_protocol_mount_target", &resource.Sweeper{
		Name: "alicloud_nas_protocol_mount_target",
		F:    testSweepNasProtocolMountTargets,
	})
	// Access-group-scoped child.
	resource.AddTestSweepers("alicloud_nas_access_rule", &resource.Sweeper{
		Name: "alicloud_nas_access_rule",
		F:    testSweepNasAccessRules,
	})
}

// testSweepNasSnapshots deletes extreme-type NAS snapshots. Snapshots are
// billed, independent resources that persist after their source file system is
// removed, so they get a dedicated region-scoped sweeper gated on SnapshotName.
func testSweepNasSnapshots(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	action := "DescribeSnapshots"
	body := map[string]interface{}{
		"FileSystemType": "extreme",
		"PageSize":       PageSizeLarge,
		"PageNumber":     1,
	}
	for {
		response, err := nasRpcPost(client, action, nil, body, true)
		if err != nil {
			log.Printf("[ERROR] %s get an error: %#v", action, err)
			return nil
		}
		resp, jpErr := jsonpath.Get("$.Snapshots.Snapshot", response)
		if jpErr != nil {
			log.Printf("[ERROR] Getting resource attribute by path %s failed!!! Body: %v.", "$.Snapshots.Snapshot", response)
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item, ok := v.(map[string]interface{})
			if !ok {
				continue
			}
			name := nasGetField(item, "SnapshotName", "Description")
			if !nasSweepNameOk(name) {
				continue
			}
			id := nasGetField(item, "SnapshotId")
			if id == "" {
				continue
			}
			nasDeleteOne(client, "DeleteSnapshot", nil,
				map[string]interface{}{"SnapshotId": id},
				id, []string{"InvalidFileSystem.NotFound"})
		}
		if len(result) < PageSizeLarge {
			return nil
		}
		body["PageNumber"] = body["PageNumber"].(int) + 1
	}
}

// testSweepNasMountTargets deletes mount targets, gated on the parent file
// system Description. Mount targets are also removed inline by the file_system
// sweeper; this sweeper covers orphaned targets and explicit per-resource
// sweeps. DescribeMountTargets uses PageSize/PageNumber pagination.
func testSweepNasMountTargets(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	fileSystems, err := nasListFileSystems(client)
	if err != nil {
		log.Printf("[ERROR] %s", err)
		return nil
	}
	for _, fs := range fileSystems {
		fsId := nasGetField(fs, "FileSystemId")
		fsDesc := nasGetField(fs, "Description")
		if !nasSweepChildOk(fsDesc, "") {
			continue
		}
		body := map[string]interface{}{
			"FileSystemId": fsId,
			"PageSize":     PageSizeLarge,
			"PageNumber":   1,
		}
		for {
			response, err := nasRpcPost(client, "DescribeMountTargets", nil, body, true)
			if err != nil {
				log.Printf("[ERROR] DescribeMountTargets of %s: %s", fsId, err)
				break
			}
			resp, jpErr := jsonpath.Get("$.MountTargets.MountTarget", response)
			if jpErr != nil {
				log.Printf("[ERROR] DescribeMountTargets parse of %s: %s", fsId, jpErr)
				break
			}
			result, _ := resp.([]interface{})
			for _, v := range result {
				item, ok := v.(map[string]interface{})
				if !ok {
					continue
				}
				domain := nasGetField(item, "MountTargetDomain")
				if domain == "" {
					continue
				}
				nasDeleteOne(client, "DeleteMountTarget", nil,
					map[string]interface{}{
						"FileSystemId":      fsId,
						"MountTargetDomain": domain,
					},
					domain, []string{"Forbidden.NasNotFound"})
			}
			if len(result) < PageSizeLarge {
				break
			}
			body["PageNumber"] = body["PageNumber"].(int) + 1
		}
	}
	return nil
}

// testSweepNasLifecyclePolicies deletes lifecycle policies, gated on the parent
// file system Description. DescribeLifecyclePolicies is GET-only in the data
// source and uses PageSize/PageNumber pagination; DeleteLifecyclePolicy is
// RpcPost and keyed by FileSystemId + LifecyclePolicyName.
func testSweepNasLifecyclePolicies(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	fileSystems, err := nasListFileSystems(client)
	if err != nil {
		log.Printf("[ERROR] %s", err)
		return nil
	}
	for _, fs := range fileSystems {
		fsId := nasGetField(fs, "FileSystemId")
		fsDesc := nasGetField(fs, "Description")
		if !nasSweepChildOk(fsDesc, "") {
			continue
		}
		query := map[string]interface{}{
			"FileSystemId": fsId,
			"PageSize":     PageSizeMedium,
			"PageNumber":   1,
		}
		for {
			// DescribeLifecyclePolicies is GET-only in the data source; call
			// RpcGet directly (nasRpcPost wraps RpcPost).
			response, err := client.RpcGet("NAS", nasSweepAPIVersion, "DescribeLifecyclePolicies", query, nil)
			if err != nil {
				log.Printf("[ERROR] DescribeLifecyclePolicies of %s: %s", fsId, err)
				break
			}
			resp, jpErr := jsonpath.Get("$.LifecyclePolicies", response)
			if jpErr != nil {
				log.Printf("[ERROR] DescribeLifecyclePolicies parse of %s: %s", fsId, jpErr)
				break
			}
			result, _ := resp.([]interface{})
			for _, v := range result {
				item, ok := v.(map[string]interface{})
				if !ok {
					continue
				}
				name := nasGetField(item, "LifecyclePolicyName")
				if !nasSweepChildOk(fsDesc, name) {
					continue
				}
				nasDeleteOne(client, "DeleteLifecyclePolicy", nil,
					map[string]interface{}{
						"FileSystemId":        fsId,
						"LifecyclePolicyName": name,
					},
					fsId+":"+name, []string{"InvalidFileSystem.NotFound", "InvalidLifecyclePolicy.NotFound"})
			}
			if len(result) < PageSizeMedium {
				break
			}
			query["PageNumber"] = query["PageNumber"].(int) + 1
		}
	}
	return nil
}

// testSweepNasAccessRules deletes access rules, gated on the parent access
// group AccessGroupName. Rules have no name field of their own, so the gate is
// purely on the parent group (tf-prefix in normal mode, all non-NODELETING
// under SWEEP_ALL). DeleteAccessRule needs AccessGroupName + AccessRuleId +
// FileSystemType, all returned by the group list and the rule list.
func testSweepNasAccessRules(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	groups, err := nasListAccessGroups(client)
	if err != nil {
		log.Printf("[ERROR] %s", err)
		return nil
	}
	for _, ag := range groups {
		agName := nasGetField(ag, "AccessGroupName")
		fsType := nasGetField(ag, "FileSystemType")
		if agName == "" {
			continue
		}
		if !nasSweepNameOk(agName) {
			continue
		}
		body := map[string]interface{}{
			"AccessGroupName": agName,
			"PageSize":        PageSizeLarge,
			"PageNumber":      1,
		}
		if fsType != "" {
			body["FileSystemType"] = fsType
		}
		for {
			response, err := nasRpcPost(client, "DescribeAccessRules", nil, body, true)
			if err != nil {
				log.Printf("[ERROR] DescribeAccessRules of %s: %s", agName, err)
				break
			}
			resp, jpErr := jsonpath.Get("$.AccessRules.AccessRule", response)
			if jpErr != nil {
				log.Printf("[ERROR] DescribeAccessRules parse of %s: %s", agName, jpErr)
				break
			}
			result, _ := resp.([]interface{})
			for _, v := range result {
				item, ok := v.(map[string]interface{})
				if !ok {
					continue
				}
				ruleId := nasGetField(item, "AccessRuleId")
				if ruleId == "" {
					continue
				}
				delBody := map[string]interface{}{
					"AccessGroupName": agName,
					"AccessRuleId":    ruleId,
				}
				if rfsType := nasGetField(item, "FileSystemType"); rfsType != "" {
					delBody["FileSystemType"] = rfsType
				} else if fsType != "" {
					delBody["FileSystemType"] = fsType
				}
				nasDeleteOne(client, "DeleteAccessRule", nil, delBody,
					agName+":"+ruleId, []string{"InvalidAccessRule.NotFound", "InvalidAccessGroup.NotFound"})
			}
			if len(result) < PageSizeLarge {
				break
			}
			body["PageNumber"] = body["PageNumber"].(int) + 1
		}
	}
	return nil
}

// nasSweepFsChildren is the shared driver for file-system-scoped child
// resources that share a NextToken/MaxResults list API and a single delete
// request. It walks every file system in the region and deletes child resources
// that belong to test file systems. Each child is gated on the parent file
// system Description and the child's own name field.
//
//	listAction   the Describe* list API (RpcPost, body slot).
//	delAction    the Delete* API (RpcPost, body slot).
//	respPath     the jsonpath to the child array in the list response.
//	nameOf       extracts the child's own name (for the do-not-delete marker).
//	listReq      builds the list body for a given file system id.
//	delOf        builds the (id, delete-body, okErrs) for a given child item.
func nasSweepFsChildren(region, listAction, delAction, respPath string,
	nameOf func(item map[string]interface{}) string,
	listReq func(fsId string) map[string]interface{},
	delOf func(item map[string]interface{}) (id string, body map[string]interface{}, okErrs []string)) error {

	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	fileSystems, err := nasListFileSystems(client)
	if err != nil {
		log.Printf("[ERROR] %s", err)
		return nil
	}
	for _, fs := range fileSystems {
		fsId := nasGetField(fs, "FileSystemId")
		fsDesc := nasGetField(fs, "Description")
		if !nasSweepChildOk(fsDesc, "") {
			continue
		}
		body := listReq(fsId)
		for {
			response, err := nasRpcPost(client, listAction, nil, body, true)
			if err != nil {
				// CPFS-only APIs (DescribeProtocolService /
				// DescribeProtocolMountTarget / DescribeFilesets /
				// DescribeDataFlows) error out for non-CPFS file systems;
				// log and move on to the next file system.
				log.Printf("[ERROR] %s of %s: %s", listAction, fsId, err)
				break
			}
			resp, jpErr := jsonpath.Get(respPath, response)
			if jpErr != nil {
				log.Printf("[ERROR] %s parse of %s: %s", listAction, fsId, jpErr)
				break
			}
			result, _ := resp.([]interface{})
			for _, v := range result {
				item, ok := v.(map[string]interface{})
				if !ok {
					continue
				}
				childName := nameOf(item)
				if !nasSweepChildOk(fsDesc, childName) {
					continue
				}
				id, delBody, okErrs := delOf(item)
				if id == "" {
					continue
				}
				nasDeleteOne(client, delAction, nil, delBody, id, okErrs)
			}
			nextToken, _ := response["NextToken"].(string)
			if nextToken == "" {
				break
			}
			body["NextToken"] = nextToken
		}
	}
	return nil
}

// testSweepNasFilesets deletes CPFS filesets, gated on the parent file system
// Description. Filesets can block CPFS file system deletion.
func testSweepNasFilesets(region string) error {
	return nasSweepFsChildren(region, "DescribeFilesets", "DeleteFileset", "$.Entries.Entrie",
		func(item map[string]interface{}) string {
			return nasGetField(item, "Description")
		},
		func(fsId string) map[string]interface{} {
			return map[string]interface{}{
				"FileSystemId": fsId,
				"MaxResults":   PageSizeLarge,
			}
		},
		func(item map[string]interface{}) (string, map[string]interface{}, []string) {
			id := nasGetField(item, "FsetId")
			return id, map[string]interface{}{
				"FileSystemId": nasGetField(item, "FileSystemId"),
				"FsetId":       id,
				"ClientToken":  buildClientToken("DeleteFileset"),
			}, []string{"InvalidFsetId.NotFound", "InvalidFilesystemType.NotSupport"}
		})
}

// testSweepNasDataFlows deletes CPFS data flows, gated on the parent file system
// Description. Data flows can block CPFS file system deletion.
func testSweepNasDataFlows(region string) error {
	return nasSweepFsChildren(region, "DescribeDataFlows", "DeleteDataFlow", "$.DataFlowInfo.DataFlow",
		func(item map[string]interface{}) string {
			return nasGetField(item, "Description")
		},
		func(fsId string) map[string]interface{} {
			return map[string]interface{}{
				"FileSystemId": fsId,
				"MaxResults":   PageSizeLarge,
			}
		},
		func(item map[string]interface{}) (string, map[string]interface{}, []string) {
			id := nasGetField(item, "DataFlowId")
			return id, map[string]interface{}{
				"FileSystemId": nasGetField(item, "FileSystemId"),
				"DataFlowId":   id,
				"ClientToken":  buildClientToken("DeleteDataFlow"),
			}, nil
		})
}

// testSweepNasAccessPoints deletes access points (general NAS NFS), gated on
// the parent file system Description and the access point's own
// AccessPointName.
func testSweepNasAccessPoints(region string) error {
	return nasSweepFsChildren(region, "DescribeAccessPoints", "DeleteAccessPoint", "$.AccessPoints",
		func(item map[string]interface{}) string {
			return nasGetField(item, "AccessPointName")
		},
		func(fsId string) map[string]interface{} {
			return map[string]interface{}{
				"FileSystemId": fsId,
				"MaxResults":   PageSizeLarge,
			}
		},
		func(item map[string]interface{}) (string, map[string]interface{}, []string) {
			id := nasGetField(item, "AccessPointId")
			return id, map[string]interface{}{
				"FileSystemId":  nasGetField(item, "FileSystemId"),
				"AccessPointId": id,
			}, []string{"InvalidFileSystem.NotFound"}
		})
}

// testSweepNasProtocolServices deletes CPFS protocol services, gated on the
// parent file system Description. Protocol services can block CPFS file system
// deletion. This sweeper runs AFTER the protocol_mount_target sweeper
// (Dependencies) so the export dirs are removed first.
func testSweepNasProtocolServices(region string) error {
	return nasSweepFsChildren(region, "DescribeProtocolService", "DeleteProtocolService", "$.ProtocolServices",
		func(item map[string]interface{}) string {
			return nasGetField(item, "Description")
		},
		func(fsId string) map[string]interface{} {
			return map[string]interface{}{
				"FileSystemId": fsId,
				"MaxResults":   PageSizeLarge,
			}
		},
		func(item map[string]interface{}) (string, map[string]interface{}, []string) {
			id := nasGetField(item, "ProtocolServiceId")
			return id, map[string]interface{}{
				"FileSystemId":      nasGetField(item, "FileSystemId"),
				"ProtocolServiceId": id,
				"ClientToken":       buildClientToken("DeleteProtocolService"),
			}, []string{"InvalidFileSystem.NotFound", "InvalidProtocolService.NotFound"}
		})
}

// testSweepNasProtocolMountTargets deletes CPFS protocol mount targets (export
// dirs), gated on the parent file system Description. They can block CPFS file
// system deletion and must be removed before their protocol service. The delete
// needs ExportId + FileSystemId + ProtocolServiceId, all returned by the list.
func testSweepNasProtocolMountTargets(region string) error {
	return nasSweepFsChildren(region, "DescribeProtocolMountTarget", "DeleteProtocolMountTarget", "$.ProtocolMountTargets",
		func(item map[string]interface{}) string {
			return nasGetField(item, "Description")
		},
		func(fsId string) map[string]interface{} {
			return map[string]interface{}{
				"FileSystemId": fsId,
				"MaxResults":   PageSizeLarge,
			}
		},
		func(item map[string]interface{}) (string, map[string]interface{}, []string) {
			id := nasGetField(item, "ExportId")
			return id, map[string]interface{}{
				"ExportId":          id,
				"FileSystemId":      nasGetField(item, "FileSystemId"),
				"ProtocolServiceId": nasGetField(item, "ProtocolServiceId"),
				"ClientToken":       buildClientToken("DeleteProtocolMountTarget"),
			}, []string{"InvalidFileSystem.NotFound", "InvalidProtocolService.NotFound"}
		})
}
