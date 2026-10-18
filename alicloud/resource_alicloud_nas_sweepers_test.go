package alicloud

import (
	"fmt"
	"log"
	"strings"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// This file 补齐 missing NAS sweepers for the following registered resources
// that previously had no dedicated sweeper registration:
//   - alicloud_nas_snapshot
//   - alicloud_nas_mount_target
//   - alicloud_nas_fileset
//   - alicloud_nas_data_flow
//   - alicloud_nas_lifecycle_policy
//   - alicloud_nas_access_point
//   - alicloud_nas_protocol_service
//   - alicloud_nas_protocol_mount_target
//
// The existing sweepers for alicloud_nas_file_system, alicloud_nas_access_group
// and alicloud_nas_auto_snapshot_policy already work and are intentionally left
// untouched. Mount targets are deleted inline by the file system sweeper; the
// dedicated sweeper here covers orphans and explicit per-resource sweeps.
//
// FS-attached feature toggles (alicloud_nas_log_analysis,
// alicloud_nas_smb_acl_attachment, alicloud_nas_recycle_bin) and
// alicloud_nas_access_rule are not given dedicated sweepers: they are removed
// with their parent file system / access group and a SWEEP_ALL pass that
// disabled them on every file system would risk disrupting real workloads
// without adding cleanup value.

const nasSweepAPIVersion = "2017-06-26"

// nasSweepPrefixes mark a resource (by name/description) as test-created.
var nasSweepPrefixes = []string{"tf-testacc", "tf_testacc", "tf-test", "test-"}

// nasSweepNameSkip returns true when a resource must be preserved:
//   - its name carries a NODELETING / do-not-delete marker, or
//   - SWEEP_ALL is off and the name does not match a test prefix.
func nasSweepNameSkip(name string) bool {
	if protectedSweepVpcName(name) {
		return true
	}
	if !sweepAll() {
		for _, prefix := range nasSweepPrefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				return false
			}
		}
		return true
	}
	return false
}

// nasSweepChildSkip gates a child resource on its parent file system's
// description (the same field the file system sweeper uses) so that children
// of NODELETING or non-test file systems are preserved. The child's own
// name/description is also checked for an explicit do-not-delete marker.
func nasSweepChildSkip(fsDescription, childName string) bool {
	if protectedSweepVpcName(fsDescription) || protectedSweepVpcName(childName) {
		return true
	}
	if !sweepAll() {
		for _, prefix := range nasSweepPrefixes {
			if strings.HasPrefix(strings.ToLower(fsDescription), strings.ToLower(prefix)) {
				return false
			}
		}
		return true
	}
	return false
}

// nasListFileSystems enumerates every NAS file system in the region.
func nasListFileSystems(client *connectivity.AliyunClient) ([]map[string]interface{}, error) {
	action := "DescribeFileSystems"
	request := map[string]interface{}{
		"RegionId":   client.Region,
		"PageSize":   PageSizeLarge,
		"PageNumber": 1,
	}
	var items []map[string]interface{}
	for {
		response, err := client.RpcPost("NAS", nasSweepAPIVersion, action, nil, request, true)
		if err != nil {
			return nil, fmt.Errorf("error retrieving NAS file systems: %s", err)
		}
		resp, err := jsonpath.Get("$.FileSystems.FileSystem", response)
		if err != nil {
			return nil, fmt.Errorf("error parsing NAS file systems response: %s", err)
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
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
}

func init() {
	resource.AddTestSweepers("alicloud_nas_snapshot", &resource.Sweeper{
		Name: "alicloud_nas_snapshot",
		F:    testSweepNasSnapshots,
	})
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
	resource.AddTestSweepers("alicloud_nas_protocol_service", &resource.Sweeper{
		Name: "alicloud_nas_protocol_service",
		F:    testSweepNasProtocolServices,
	})
	resource.AddTestSweepers("alicloud_nas_protocol_mount_target", &resource.Sweeper{
		Name: "alicloud_nas_protocol_mount_target",
		F:    testSweepNasProtocolMountTargets,
	})
}

// testSweepNasSnapshots deletes extreme NAS snapshots. Snapshots are billed,
// independent resources that persist after the file system is removed, so they
// get a dedicated sweeper gated on SnapshotName.
func testSweepNasSnapshots(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	action := "DescribeSnapshots"
	request := map[string]interface{}{
		"FileSystemType": "extreme",
		"PageSize":       PageSizeLarge,
		"PageNumber":     1,
	}
	for {
		response, err := client.RpcPost("NAS", nasSweepAPIVersion, action, nil, request, true)
		if err != nil {
			log.Printf("[ERROR] Error listing NAS snapshots: %s", err)
			return nil
		}
		resp, err := jsonpath.Get("$.Snapshots.Snapshot", response)
		if err != nil {
			log.Printf("[ERROR] Error parsing NAS snapshots response: %s", err)
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item, ok := v.(map[string]interface{})
			if !ok {
				continue
			}
			name := fmt.Sprint(item["SnapshotName"])
			id := fmt.Sprint(item["SnapshotId"])
			if nasSweepNameSkip(name) {
				log.Printf("[INFO] Skipping NAS Snapshot: %s (%s)", name, id)
				continue
			}
			log.Printf("[INFO] Deleting NAS Snapshot: %s (%s)", name, id)
			delReq := map[string]interface{}{"SnapshotId": id}
			if _, err := client.RpcPost("NAS", nasSweepAPIVersion, "DeleteSnapshot", nil, delReq, false); err != nil {
				log.Printf("[ERROR] Failed to delete NAS Snapshot (%s): %s", id, err)
			}
		}
		if len(result) < PageSizeLarge {
			return nil
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
}

// testSweepNasMountTargets deletes mount targets, gated on the parent file
// system's description. Mount targets are also removed inline by the file
// system sweeper; this sweeper covers orphaned targets and explicit sweeps.
func testSweepNasMountTargets(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	fileSystems, err := nasListFileSystems(client)
	if err != nil {
		log.Printf("[ERROR] %s", err)
		return nil
	}
	for _, fs := range fileSystems {
		fsId := fmt.Sprint(fs["FileSystemId"])
		fsDesc := fmt.Sprint(fs["Description"])
		if nasSweepChildSkip(fsDesc, "") {
			log.Printf("[INFO] Skipping mount targets of file system %s (%s)", fsId, fsDesc)
			continue
		}
		request := map[string]interface{}{
			"FileSystemId": fsId,
			"RegionId":     client.Region,
			"PageSize":     PageSizeLarge,
			"PageNumber":   1,
		}
		for {
			response, err := client.RpcPost("NAS", nasSweepAPIVersion, "DescribeMountTargets", nil, request, true)
			if err != nil {
				log.Printf("[ERROR] Error listing mount targets of %s: %s", fsId, err)
				break
			}
			resp, err := jsonpath.Get("$.MountTargets.MountTarget", response)
			if err != nil {
				log.Printf("[ERROR] Error parsing mount targets of %s: %s", fsId, err)
				break
			}
			result, _ := resp.([]interface{})
			for _, v := range result {
				item, ok := v.(map[string]interface{})
				if !ok {
					continue
				}
				domain := fmt.Sprint(item["MountTargetDomain"])
				log.Printf("[INFO] Deleting NAS MountTarget: %s (fs %s)", domain, fsId)
				delReq := map[string]interface{}{
					"FileSystemId":      fsId,
					"MountTargetDomain": domain,
				}
				if _, err := client.RpcPost("NAS", nasSweepAPIVersion, "DeleteMountTarget", nil, delReq, true); err != nil {
					log.Printf("[ERROR] Failed to delete NAS MountTarget (%s): %s", domain, err)
				}
			}
			if len(result) < PageSizeLarge {
				break
			}
			request["PageNumber"] = request["PageNumber"].(int) + 1
		}
	}
	return nil
}

// testSweepNasFilesets deletes CPFS filesets, gated on the parent file
// system's description. Filesets can block CPFS file system deletion.
func testSweepNasFilesets(region string) error {
	return sweepNasFilesystemChildren(region, "DescribeFilesets", "$.Entries.Entrie",
		func(fs map[string]interface{}) map[string]interface{} {
			return map[string]interface{}{
				"FileSystemId": fmt.Sprint(fs["FileSystemId"]),
				"MaxResults":   PageSizeLarge,
			}
		}, func(item, fs map[string]interface{}) (string, map[string]interface{}, string) {
			id := fmt.Sprint(item["FsetId"])
			delReq := map[string]interface{}{
				"FileSystemId": fmt.Sprint(fs["FileSystemId"]),
				"FsetId":       id,
			}
			return id, delReq, "DeleteFileset"
		})
}

// testSweepNasDataFlows deletes CPFS data flows, gated on the parent file
// system's description. Data flows can block CPFS file system deletion.
func testSweepNasDataFlows(region string) error {
	return sweepNasFilesystemChildren(region, "DescribeDataFlows", "$.DataFlowInfo.DataFlow",
		func(fs map[string]interface{}) map[string]interface{} {
			return map[string]interface{}{
				"FileSystemId": fmt.Sprint(fs["FileSystemId"]),
				"MaxResults":   PageSizeLarge,
			}
		}, func(item, fs map[string]interface{}) (string, map[string]interface{}, string) {
			id := fmt.Sprint(item["DataFlowId"])
			delReq := map[string]interface{}{
				"FileSystemId": fmt.Sprint(fs["FileSystemId"]),
				"DataFlowId":   id,
			}
			return id, delReq, "DeleteDataFlow"
		})
}

// testSweepNasLifecyclePolicies deletes lifecycle policies, gated on the
// parent file system's description.
func testSweepNasLifecyclePolicies(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	fileSystems, err := nasListFileSystems(client)
	if err != nil {
		log.Printf("[ERROR] %s", err)
		return nil
	}
	for _, fs := range fileSystems {
		fsId := fmt.Sprint(fs["FileSystemId"])
		fsDesc := fmt.Sprint(fs["Description"])
		if nasSweepChildSkip(fsDesc, "") {
			log.Printf("[INFO] Skipping lifecycle policies of file system %s (%s)", fsId, fsDesc)
			continue
		}
		request := map[string]interface{}{
			"FileSystemId": fsId,
			"PageSize":     PageSizeMedium,
			"PageNumber":   1,
		}
		for {
			response, err := client.RpcGet("NAS", nasSweepAPIVersion, "DescribeLifecyclePolicies", request, nil)
			if err != nil {
				log.Printf("[ERROR] Error listing lifecycle policies of %s: %s", fsId, err)
				break
			}
			resp, err := jsonpath.Get("$.LifecyclePolicies", response)
			if err != nil {
				log.Printf("[ERROR] Error parsing lifecycle policies of %s: %s", fsId, err)
				break
			}
			result, _ := resp.([]interface{})
			for _, v := range result {
				item, ok := v.(map[string]interface{})
				if !ok {
					continue
				}
				name := fmt.Sprint(item["LifecyclePolicyName"])
				if nasSweepChildSkip(fsDesc, name) {
					log.Printf("[INFO] Skipping NAS LifecyclePolicy: %s (%s:%s)", name, fsId, name)
					continue
				}
				log.Printf("[INFO] Deleting NAS LifecyclePolicy: %s (%s:%s)", name, fsId, name)
				delReq := map[string]interface{}{
					"FileSystemId":        fsId,
					"LifecyclePolicyName": name,
				}
				if _, err := client.RpcPost("NAS", nasSweepAPIVersion, "DeleteLifecyclePolicy", nil, delReq, false); err != nil {
					log.Printf("[ERROR] Failed to delete NAS LifecyclePolicy (%s:%s): %s", fsId, name, err)
				}
			}
			if len(result) < PageSizeMedium {
				break
			}
			request["PageNumber"] = request["PageNumber"].(int) + 1
		}
	}
	return nil
}

// testSweepNasAccessPoints deletes access points (general NAS NFS), gated on
// the parent file system's description.
func testSweepNasAccessPoints(region string) error {
	return sweepNasFilesystemChildren(region, "DescribeAccessPoints", "$.AccessPoints",
		func(fs map[string]interface{}) map[string]interface{} {
			return map[string]interface{}{
				"FileSystemId": fmt.Sprint(fs["FileSystemId"]),
				"MaxResults":   PageSizeLarge,
			}
		}, func(item, fs map[string]interface{}) (string, map[string]interface{}, string) {
			id := fmt.Sprint(item["AccessPointId"])
			delReq := map[string]interface{}{
				"FileSystemId":  fmt.Sprint(fs["FileSystemId"]),
				"AccessPointId": id,
			}
			return id, delReq, "DeleteAccessPoint"
		})
}

// testSweepNasProtocolServices deletes CPFS protocol services, gated on the
// parent file system's description. Protocol services can block CPFS file
// system deletion.
func testSweepNasProtocolServices(region string) error {
	return sweepNasFilesystemChildren(region, "DescribeProtocolService", "$.ProtocolServices[*]",
		func(fs map[string]interface{}) map[string]interface{} {
			return map[string]interface{}{
				"FileSystemId": fmt.Sprint(fs["FileSystemId"]),
				"MaxResults":   PageSizeLarge,
			}
		}, func(item, fs map[string]interface{}) (string, map[string]interface{}, string) {
			id := fmt.Sprint(item["ProtocolServiceId"])
			delReq := map[string]interface{}{
				"FileSystemId":      fmt.Sprint(fs["FileSystemId"]),
				"ProtocolServiceId": id,
			}
			return id, delReq, "DeleteProtocolService"
		})
}

// testSweepNasProtocolMountTargets deletes CPFS protocol mount targets, gated
// on the parent file system's description. They can block CPFS file system
// deletion and must be removed before their protocol service.
func testSweepNasProtocolMountTargets(region string) error {
	return sweepNasFilesystemChildren(region, "DescribeProtocolMountTarget", "$.ProtocolMountTargets",
		func(fs map[string]interface{}) map[string]interface{} {
			return map[string]interface{}{
				"FileSystemId": fmt.Sprint(fs["FileSystemId"]),
				"MaxResults":   PageSizeLarge,
			}
		}, func(item, fs map[string]interface{}) (string, map[string]interface{}, string) {
			id := fmt.Sprint(item["ExportId"])
			delReq := map[string]interface{}{
				"ExportId":          id,
				"FileSystemId":      fmt.Sprint(fs["FileSystemId"]),
				"ProtocolServiceId": fmt.Sprint(item["ProtocolServiceId"]),
			}
			return id, delReq, "DeleteProtocolMountTarget"
		})
}

// sweepNasFilesystemChildren is the shared driver for per-file-system child
// resources that share a NextToken/MaxResults list API and a single-key delete
// request. It gates every child on the parent file system description.
func sweepNasFilesystemChildren(region, listAction, respPath string,
	listReq func(fs map[string]interface{}) map[string]interface{},
	delReq func(item, fs map[string]interface{}) (string, map[string]interface{}, string)) error {

	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	fileSystems, err := nasListFileSystems(client)
	if err != nil {
		log.Printf("[ERROR] %s", err)
		return nil
	}
	for _, fs := range fileSystems {
		fsId := fmt.Sprint(fs["FileSystemId"])
		fsDesc := fmt.Sprint(fs["Description"])
		if nasSweepChildSkip(fsDesc, "") {
			log.Printf("[INFO] Skipping %s of file system %s (%s)", listAction, fsId, fsDesc)
			continue
		}
		request := listReq(fs)
		for {
			response, err := client.RpcPost("NAS", nasSweepAPIVersion, listAction, nil, request, true)
			if err != nil {
				// CPFS-only APIs return errors for non-CPFS file systems; skip them.
				log.Printf("[ERROR] Error listing %s of %s: %s", listAction, fsId, err)
				break
			}
			resp, err := jsonpath.Get(respPath, response)
			if err != nil {
				log.Printf("[ERROR] Error parsing %s of %s: %s", listAction, fsId, err)
				break
			}
			result, _ := resp.([]interface{})
			for _, v := range result {
				item, ok := v.(map[string]interface{})
				if !ok {
					continue
				}
				id, delReqMap, delAction := delReq(item, fs)
				childName := fmt.Sprint(item["Description"])
				if nasSweepChildSkip(fsDesc, childName) {
					log.Printf("[INFO] Skipping %s: %s (%s)", delAction, childName, id)
					continue
				}
				log.Printf("[INFO] Deleting %s: %s (%s)", delAction, id, fsId)
				if _, err := client.RpcPost("NAS", nasSweepAPIVersion, delAction, nil, delReqMap, true); err != nil {
					log.Printf("[ERROR] Failed to delete %s (%s): %s", delAction, id, err)
				}
			}
			nextToken, _ := response["NextToken"].(string)
			if nextToken == "" {
				break
			}
			request["NextToken"] = nextToken
		}
	}
	return nil
}
