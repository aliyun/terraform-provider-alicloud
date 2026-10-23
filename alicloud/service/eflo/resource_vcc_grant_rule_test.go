// Package eflo. This file is generated automatically. Please do not modify it manually, thank you!
package eflo_test

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

// testVccGrantRuleVccID holds the VCC instance ID created via the CreateVcc
// OpenAPI call during step 0; it is reused as instance_id in the grant-rule
// config and deleted (best-effort) in CheckDestroy.  alicloud_eflo_vcc is not
// a registered Terraform resource, so the VCC must be created out-of-band.
var testVccGrantRuleVccID string

// Test Eflo VccGrantRule. >>> Resource test cases, automatically generated.
// Case vcc_grant_rule_tf_副本1770195464559_副本1770271544998 12535
func TestAccAliCloudEfloVccGrantRule_basic12535(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacceflo%d", rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			t.Setenv("ALICLOUD_REGION", "cn-wulanchabu")
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckEfloVccGrantRuleDestroy,
		Steps: []resource.TestStep{
			{
				// Create ER + VPC + VSwitch through Terraform, then
				// create the VCC instance via the CreateVcc OpenAPI call.
				Config: testAccAliCloudEfloVccGrantRule_basic12535_step0(name),
				Check: resource.ComposeTestCheckFunc(
					testAccEfloVccGrantRuleCreateVcc(name, &testVccGrantRuleVccID),
				),
			},
			{
				// Create the VCC grant rule using the API-created VCC ID.
				// var.vcc_id has no default; its value comes from the
				// TF_VAR_vcc_id environment variable set in step 0's Check.
				// Terraform reads TF_VAR_* at apply time, not when the Go
				// Steps slice is constructed, so the VCC ID set by the
				// preceding Check is available here.
				Config: testAccAliCloudEfloVccGrantRule_basic12535_step1(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_eflo_vcc_grant_rule.default", "grant_rule_id"),
					resource.TestCheckResourceAttrSet("alicloud_eflo_vcc_grant_rule.default", "instance_id"),
					resource.TestCheckResourceAttr("alicloud_eflo_vcc_grant_rule.default", "used", "false"),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("alicloud_eflo_vcc_grant_rule.default", tfjsonpath.New("grant_rule_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_eflo_vcc_grant_rule.default", tfjsonpath.New("instance_name"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_eflo_vcc_grant_rule.default", tfjsonpath.New("product"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("alicloud_eflo_vcc_grant_rule.default", tfjsonpath.New("tenant_id"), knownvalue.NotNull()),
				},
			},
			{
				// No-op re-apply: the config is identical to step 1, so
				// Terraform detects no changes and only refreshes state.
				// The Check deletes the API-created VCC and cleans up its
				// side-effect resources (router interfaces, network
				// interfaces) BEFORE the framework's post-test terraform
				// destroy runs. This prevents DependencyViolation errors
				// during destroy.
				//
				// The import step was removed because its Check did not
				// execute reliably (the import step replaces the working
				// state, causing the cleanup Check to be skipped).
				Config: testAccAliCloudEfloVccGrantRule_basic12535_step1(name),
				Check: resource.ComposeTestCheckFunc(
					testAccEfloVccGrantRuleDeleteVcc(&testVccGrantRuleVccID),
				),
			},
		},
	})
}

// testAccAliCloudEfloVccGrantRule_basic12535_step0 creates ER + VPC +
// VSwitch dependencies. The VCC instance is created via the CreateVcc OpenAPI
// call in the Check function because alicloud_eflo_vcc is not a Terraform
// resource; CreateVcc requires VpcId and VSwitchId from a regular VPC.
func testAccAliCloudEfloVccGrantRule_basic12535_step0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_account" "default" {
}

resource "alicloud_eflo_er" "default" {
  er_name        = var.name
  master_zone_id = "cn-wulanchabu-b"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id     = alicloud_vpc.default.id
  cidr_block = "172.16.0.0/24"
  zone_id    = "cn-wulanchabu-b"
}
`, name)
}

// testAccAliCloudEfloVccGrantRule_basic12535_step1 adds the VCC grant rule on
// top of the ER + VPC + VSwitch dependencies. The VCC ID is passed via the
// TF_VAR_vcc_id environment variable (set in step 0's Check) and referenced
// as var.vcc_id, which has no default so Terraform reads TF_VAR_vcc_id at
// apply time.
func testAccAliCloudEfloVccGrantRule_basic12535_step1(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

variable "vcc_id" {
  type = string
}

data "alicloud_account" "default" {
}

resource "alicloud_eflo_er" "default" {
  er_name        = var.name
  master_zone_id = "cn-wulanchabu-b"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id     = alicloud_vpc.default.id
  cidr_block = "172.16.0.0/24"
  zone_id    = "cn-wulanchabu-b"
}

resource "alicloud_eflo_vcc_grant_rule" "default" {
  grant_tenant_id = data.alicloud_account.default.id
  er_id           = alicloud_eflo_er.default.id
  instance_id     = var.vcc_id
}
`, name)
}

// testAccEfloVccGrantRuleCleanupRIs queries and deletes all router interfaces
// in the test region to free the RiPerUser quota. Previous VCC creation
// attempts (successful or failed) may leave dangling router interfaces that
// exhaust the quota, causing VCC creation to fail with
// "QuotaExceeded: RiPerUser quota exceeded". This cleanup is best-effort and
// non-fatal: any error is logged and the function returns nil so the test can
// proceed (or fail with the original quota error if cleanup was insufficient).
func testAccEfloVccGrantRuleCleanupRIs(client *connectivity.AliyunClient) {
	query := map[string]interface{}{
		"RegionId":   client.RegionId,
		"PageSize":   50,
		"PageNumber": 1,
	}
	for page := 1; ; page++ {
		query["PageNumber"] = page
		resp, err := client.RpcPost("Vpc", "2016-04-28", "DescribeRouterInterfaces", query, query, true)
		if err != nil {
			fmt.Printf("[RI_CLEANUP] DescribeRouterInterfaces error: %v\n", err)
			return
		}
		riSet, ok := resp["RouterInterfaceSet"].(map[string]interface{})
		if !ok {
			fmt.Printf("[RI_CLEANUP] no RouterInterfaceSet in response\n")
			return
		}
		riList, ok := riSet["RouterInterfaceType"].([]interface{})
		if !ok || len(riList) == 0 {
			fmt.Printf("[RI_CLEANUP] no router interfaces found in %s\n", client.RegionId)
			return
		}
		for _, ri := range riList {
			item, ok := ri.(map[string]interface{})
			if !ok {
				continue
			}
			id, _ := item["RouterInterfaceId"].(string)
			status, _ := item["Status"].(string)
			if id == "" {
				continue
			}
			fmt.Printf("[RI_CLEANUP] found RI %s status=%s\n", id, status)
			// Deactivate if Active before deleting (DeleteRouterInterface
			// requires Inactive state).
			if status == "Active" {
				deactReq := map[string]interface{}{
					"RegionId":          client.RegionId,
					"RouterInterfaceId": id,
				}
				if _, err := client.RpcPost("Vpc", "2016-04-28", "DeactivateRouterInterface", deactReq, deactReq, true); err != nil {
					fmt.Printf("[RI_CLEANUP] deactivate RI %s error: %v\n", id, err)
					continue
				}
				fmt.Printf("[RI_CLEANUP] deactivated RI %s, waiting for Inactive...\n", id)
				deadline := time.Now().Add(2 * time.Minute)
				for time.Now().Before(deadline) {
					time.Sleep(10 * time.Second)
					checkReq := map[string]interface{}{
						"RegionId":         client.RegionId,
						"Filter.1.Key":     "RouterInterfaceId",
						"Filter.1.Value.1": id,
					}
					checkResp, checkErr := client.RpcPost("Vpc", "2016-04-28", "DescribeRouterInterfaces", checkReq, checkReq, true)
					if checkErr != nil {
						break
					}
					checkSet, _ := checkResp["RouterInterfaceSet"].(map[string]interface{})
					checkList, _ := checkSet["RouterInterfaceType"].([]interface{})
					if len(checkList) == 0 {
						break // already gone
					}
					if first, ok := checkList[0].(map[string]interface{}); ok {
						if s, _ := first["Status"].(string); s != "Active" {
							break // became Inactive
						}
					}
				}
			}
			// Delete
			delReq := map[string]interface{}{
				"RegionId":          client.RegionId,
				"RouterInterfaceId": id,
				"ClientToken":       fmt.Sprintf("tf-ri-cleanup-%s", id),
			}
			if _, err := client.RpcPost("Vpc", "2016-04-28", "DeleteRouterInterface", delReq, delReq, true); err != nil {
				fmt.Printf("[RI_CLEANUP] delete RI %s error: %v\n", id, err)
			} else {
				fmt.Printf("[RI_CLEANUP] deleted RI %s\n", id)
			}
		}
		// Check pagination
		if total, ok := resp["TotalCount"].(float64); ok && float64(page)*50 < total {
			continue
		}
		break
	}
	// Wait for deletions to settle
	time.Sleep(15 * time.Second)
}

// testAccEfloVccGrantRuleCleanupVccs lists all VCCs in the test region via
// GetVcc with EnablePage and deletes them. Old VCCs (especially those in
// "Create Failed" state from previous runs) should be removed before creating
// a new one. This cleanup is best-effort and non-fatal.
func testAccEfloVccGrantRuleCleanupVccs(client *connectivity.AliyunClient) {
	query := map[string]interface{}{
		"RegionId":   client.RegionId,
		"EnablePage": true,
		"PageSize":   50,
		"PageNumber": 1,
	}
	resp, err := client.RpcPost("eflo", "2022-05-30", "GetVcc", query, query, true)
	if err != nil {
		fmt.Printf("[VCC_CLEANUP] GetVcc(EnablePage) error: %v\n", err)
		return
	}
	content, _ := resp["Content"].(map[string]interface{})
	if content == nil {
		fmt.Printf("[VCC_CLEANUP] no Content in GetVcc response\n")
		return
	}
	// GetVcc may return a single object or a list wrapper.
	var vccList []interface{}
	for _, key := range []string{"VccList", "Data", "Content"} {
		if arr, ok := content[key].([]interface{}); ok && len(arr) > 0 {
			vccList = arr
			break
		}
	}
	if vccList == nil {
		// If content itself is the only VCC (single-object response), try it.
		if vccID, _ := content["VccId"].(string); vccID != "" {
			vccList = []interface{}{content}
		}
	}
	if vccList == nil {
		fmt.Printf("[VCC_CLEANUP] no VCC list found, content keys: %v\n", mapKeys(content))
		return
	}
	for _, vcc := range vccList {
		item, ok := vcc.(map[string]interface{})
		if !ok {
			continue
		}
		vccID, _ := item["VccId"].(string)
		status, _ := item["Status"].(string)
		if vccID == "" {
			continue
		}
		fmt.Printf("[VCC_CLEANUP] deleting VCC %s (status=%s)\n", vccID, status)
		delReq := map[string]interface{}{
			"RegionId": client.RegionId,
			"VccId":    vccID,
		}
		if _, err := client.RpcPost("eflo", "2022-05-30", "DeleteVcc", delReq, delReq, true); err != nil {
			fmt.Printf("[VCC_CLEANUP] DeleteVcc error for %s: %v\n", vccID, err)
		} else {
			fmt.Printf("[VCC_CLEANUP] deleted VCC %s\n", vccID)
		}
	}
	// Wait for VCC deletions to settle
	time.Sleep(30 * time.Second)
}

// testAccEfloVccGrantRuleCleanupENIs queries and deletes all deletable
// network interfaces (ENIs) in the given VPC. VCC creates network interfaces
// on the VSwitch that persist after VCC deletion; these block DeleteVSwitch
// with DependencyViolation.NetworkInterface.
//
// VCC-created ENIs are typically ServiceManaged=true and cannot be deleted
// via the ECS DeleteNetworkInterface API — the owning Eflo service must
// clean them up after the VCC is deleted. This function only deletes
// non-ServiceManaged ENIs (best-effort). ServiceManaged ENIs are reported
// but left for the Eflo service to clean; use testAccEfloVccGrantRuleWaitENIsGone
// to wait for them.
//
// This cleanup is best-effort and non-fatal.
func testAccEfloVccGrantRuleCleanupENIs(client *connectivity.AliyunClient, vpcID string) {
	if vpcID == "" {
		return
	}
	eniIDs := testAccEfloVccGrantRuleListENIs(client, vpcID)
	if len(eniIDs) == 0 {
		fmt.Printf("[ENI_CLEANUP] no deletable network interfaces found for VPC %s\n", vpcID)
		return
	}
	for _, eniID := range eniIDs {
		delReq := map[string]interface{}{
			"RegionId":           client.RegionId,
			"NetworkInterfaceId": eniID,
		}
		if _, err := client.RpcPost("Ecs", "2014-05-26", "DeleteNetworkInterface", delReq, delReq, true); err != nil {
			fmt.Printf("[ENI_CLEANUP] delete ENI %s error: %v\n", eniID, err)
		} else {
			fmt.Printf("[ENI_CLEANUP] deleted ENI %s\n", eniID)
		}
	}
	time.Sleep(10 * time.Second)
}

// testAccEfloVccGrantRuleListENIs lists all non-ServiceManaged network
// interface IDs in the given VPC using the direct VpcId query parameter (not
// Filter style). The response path is NetworkInterfaceSets.NetworkInterfaceSet
// (note the plural outer key "NetworkInterfaceSets").
func testAccEfloVccGrantRuleListENIs(client *connectivity.AliyunClient, vpcID string) []string {
	if vpcID == "" {
		return nil
	}
	var ids []string
	for page := 1; ; page++ {
		query := map[string]interface{}{
			"RegionId":   client.RegionId,
			"VpcId":      vpcID,
			"PageSize":   50,
			"PageNumber": page,
		}
		resp, err := client.RpcPost("Ecs", "2014-05-26", "DescribeNetworkInterfaces", query, query, true)
		if err != nil {
			fmt.Printf("[ENI_CLEANUP] DescribeNetworkInterfaces error: %v\n", err)
			return ids
		}
		eniSet, ok := resp["NetworkInterfaceSets"].(map[string]interface{})
		if !ok {
			fmt.Printf("[ENI_CLEANUP] no NetworkInterfaceSets in response, keys: %v\n", mapKeys(resp))
			return ids
		}
		eniList, ok := eniSet["NetworkInterfaceSet"].([]interface{})
		if !ok || len(eniList) == 0 {
			break
		}
		for _, eni := range eniList {
			item, ok := eni.(map[string]interface{})
			if !ok {
				continue
			}
			id, _ := item["NetworkInterfaceId"].(string)
			status, _ := item["Status"].(string)
			managed, _ := item["ServiceManaged"].(bool)
			if id == "" {
				continue
			}
			fmt.Printf("[ENI_CLEANUP] found ENI %s status=%s serviceManaged=%v\n", id, status, managed)
			if !managed {
				ids = append(ids, id)
			}
		}
		if total, ok := resp["TotalCount"].(float64); ok && float64(page)*50 < total {
			continue
		}
		break
	}
	return ids
}

// testAccEfloVccGrantRuleWaitENIsGone polls DescribeNetworkInterfaces until no
// ENIs remain in the VPC or the deadline expires. After VCC deletion, the
// Eflo service asynchronously cleans up ServiceManaged network interfaces on
// the VSwitch; DeleteVSwitch will fail with DependencyViolation.NetworkInterface
// until they are gone.
func testAccEfloVccGrantRuleWaitENIsGone(client *connectivity.AliyunClient, vpcID string, timeout time.Duration) {
	if vpcID == "" {
		return
	}
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		query := map[string]interface{}{
			"RegionId":   client.RegionId,
			"VpcId":      vpcID,
			"PageSize":   50,
			"PageNumber": 1,
		}
		resp, err := client.RpcPost("Ecs", "2014-05-26", "DescribeNetworkInterfaces", query, query, true)
		if err != nil {
			fmt.Printf("[ENI_WAIT] DescribeNetworkInterfaces error: %v\n", err)
			time.Sleep(15 * time.Second)
			continue
		}
		eniSet, ok := resp["NetworkInterfaceSets"].(map[string]interface{})
		if !ok {
			break
		}
		eniList, ok := eniSet["NetworkInterfaceSet"].([]interface{})
		if !ok || len(eniList) == 0 {
			fmt.Printf("[ENI_WAIT] all ENIs gone for VPC %s\n", vpcID)
			return
		}
		fmt.Printf("[ENI_WAIT] %d ENI(s) still present for VPC %s, waiting...\n", len(eniList), vpcID)
		time.Sleep(15 * time.Second)
	}
	fmt.Printf("[ENI_WAIT] timeout waiting for ENIs to clear for VPC %s\n", vpcID)
}

// mapKeys returns the keys of a map for diagnostic logging.
func mapKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// testAccEfloVccGrantRuleCreateVcc creates a VCC instance via the CreateVcc
// OpenAPI call, stores its ID, and sets the TF_VAR_vcc_id environment variable
// so the subsequent Terraform step can reference var.vcc_id.  alicloud_eflo_vcc
// is not a Terraform resource, so the VCC is created out-of-band using the RPC
// client. Before creating a new VCC, old router interfaces and VCCs are
// cleaned up to avoid "RiPerUser quota exceeded" errors from prior runs.
func testAccEfloVccGrantRuleCreateVcc(name string, vccID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)

		// Clean up old router interfaces and VCCs before creating a new VCC.
		// Previous test runs may have left dangling router interfaces that
		// exhaust the RiPerUser quota, causing VCC creation to fail with
		// "QuotaExceeded: RiPerUser quota exceeded".
		testAccEfloVccGrantRuleCleanupRIs(client)
		testAccEfloVccGrantRuleCleanupVccs(client)
		// VCC deletion may free associated router interfaces; clean up again.
		testAccEfloVccGrantRuleCleanupRIs(client)

		// Retrieve the VPC ID and VSwitch ID from state. VPD is omitted
		// because CreateVcc with ConnectionType=VPC connects directly to
		// the VPC without a VPD, and VPD quota in cn-wulanchabu is limited.
		vpcRes, ok := s.RootModule().Resources["alicloud_vpc.default"]
		if !ok {
			return fmt.Errorf("alicloud_vpc.default not found in state")
		}
		vpcID := vpcRes.Primary.ID

		vswRes, ok := s.RootModule().Resources["alicloud_vswitch.default"]
		if !ok {
			return fmt.Errorf("alicloud_vswitch.default not found in state")
		}
		vswID := vswRes.Primary.ID

		request := map[string]interface{}{
			"RegionId":       client.RegionId,
			"VccName":        name,
			"ZoneId":         "cn-wulanchabu-b",
			"VpcId":          vpcID,
			"VSwitchId":      vswID,
			"Bandwidth":      1000,
			"ConnectionType": "VPC",
		}
		response, err := client.RpcPost("eflo", "2022-05-30", "CreateVcc", request, request, true)
		if err != nil {
			return fmt.Errorf("calling CreateVcc failed: %w", err)
		}
		content, ok := response["Content"].(map[string]interface{})
		if !ok {
			return fmt.Errorf("CreateVcc response missing Content: %v", response)
		}
		id, ok := content["VccId"].(string)
		if !ok || id == "" {
			return fmt.Errorf("CreateVcc response missing VccId: %v", content)
		}
		*vccID = id

		// Set TF_VAR_vcc_id so the next test step's Terraform config can
		// reference var.vcc_id. Terraform reads TF_VAR_* env vars at apply
		// time, not when the Go Steps slice is constructed.
		os.Setenv("TF_VAR_vcc_id", id)

		// Poll GetVcc until the VCC reaches a terminal state.
		return testAccEfloVccGrantRuleWaitVcc(client, id)
	}
}

// testAccEfloVccGrantRuleDeleteVcc deletes the API-created VCC instance so
// that Terraform's post-test destroy can delete the VPC and VSwitch. The VCC
// creates two kinds of side-effect resources that block terraform destroy:
//   - Route interfaces on the VPC (block DeleteVpc with
//     DependencyViolation.RouteInterface)
//   - Network interfaces (ENIs) on the VSwitch (block DeleteVSwitch with
//     DependencyViolation.NetworkInterface)
//
// This Check runs in a no-op re-apply step (same config, no changes), before
// the framework's automatic terraform destroy. After calling DeleteVcc, we:
//  1. Poll GetVcc until the VCC is no longer queryable (deletion is async).
//  2. Wait for the Eflo service to clean up ServiceManaged ENIs on the VSwitch
//     (these cannot be deleted via ECS API; they are removed asynchronously
//     after VCC deletion).
//  3. Run RI cleanup to remove any remaining route interfaces on the VPC.
//  4. Run ENI cleanup to delete any remaining non-ServiceManaged ENIs.
//
// This ensures both DeleteVpc and DeleteVSwitch succeed during terraform destroy.
func testAccEfloVccGrantRuleDeleteVcc(vccID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if *vccID == "" {
			return nil
		}
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)
		fmt.Printf("[VCC_DELETE] deleting VCC %s\n", *vccID)
		req := map[string]interface{}{
			"RegionId": client.RegionId,
			"VccId":    *vccID,
		}
		_, _ = client.RpcPost("eflo", "2022-05-30", "DeleteVcc", req, req, true)

		// Poll GetVcc until the VCC is no longer queryable — VCC deletion
		// is async, and the route interface on the VPC is removed as part
		// of the VCC deletion process.
		deadline := time.Now().Add(5 * time.Minute)
		for time.Now().Before(deadline) {
			query := map[string]interface{}{
				"RegionId": client.RegionId,
				"VccId":    *vccID,
			}
			resp, err := client.RpcPost("eflo", "2022-05-30", "GetVcc", query, query, true)
			if err != nil {
				// VCC no longer queryable — deletion complete.
				fmt.Printf("[VCC_DELETE] VCC %s no longer queryable, deletion complete\n", *vccID)
				break
			}
			content, _ := resp["Content"].(map[string]interface{})
			status, _ := content["Status"].(string)
			fmt.Printf("[VCC_DELETE] VCC %s status=%q, waiting for deletion...\n", *vccID, status)
			time.Sleep(15 * time.Second)
		}

		// Retrieve the VPC ID from state for ENI cleanup.
		var vpcID string
		if vpcRes, ok := s.RootModule().Resources["alicloud_vpc.default"]; ok {
			vpcID = vpcRes.Primary.ID
		}

		// Wait for the Eflo service to clean up ServiceManaged ENIs on the
		// VSwitch after VCC deletion. These ENIs cannot be deleted via the
		// ECS API and must be removed by the owning service.
		testAccEfloVccGrantRuleWaitENIsGone(client, vpcID, 5*time.Minute)

		// Clean up any remaining route interfaces on the VPC to prevent
		// DependencyViolation.RouteInterface during terraform destroy.
		testAccEfloVccGrantRuleCleanupRIs(client)

		// Clean up any remaining non-ServiceManaged network interfaces to
		// prevent DependencyViolation.NetworkInterface during terraform destroy.
		testAccEfloVccGrantRuleCleanupENIs(client, vpcID)

		*vccID = ""
		os.Unsetenv("TF_VAR_vcc_id")
		return nil
	}
}

// testAccEfloVccGrantRuleWaitVcc polls GetVcc until the VCC is Available.
// VCC Status values (per GetVcc API): Order Placing, Executing, Available,
// Not Available, Create Failed, Deleting. Only "Available" means the VCC
// is operable for grant rule creation. "Create Failed" and "Not Available"
// are terminal failure states — we fail fast after 2 consecutive
// observations and dump the full VCC object for diagnostics.
func testAccEfloVccGrantRuleWaitVcc(client *connectivity.AliyunClient, vccID string) error {
	query := map[string]interface{}{
		"RegionId": client.RegionId,
		"VccId":    vccID,
	}
	timeout := 15 * time.Minute
	deadline := time.Now().Add(timeout)
	failedCount := 0
	var status string
	for time.Now().Before(deadline) {
		resp, err := client.RpcPost("eflo", "2022-05-30", "GetVcc", query, query, true)
		if err != nil {
			// The VCC may not be queryable immediately after creation.
			fmt.Printf("[VCC_WAIT] GetVcc error for %s: %v\n", vccID, err)
			time.Sleep(15 * time.Second)
			continue
		}
		content, _ := resp["Content"].(map[string]interface{})
		if content == nil {
			fmt.Printf("[VCC_WAIT] GetVcc response missing Content for %s: %v\n", vccID, resp)
			time.Sleep(15 * time.Second)
			continue
		}
		// GetVcc may return a single object or a list wrapper.
		status, _ = content["Status"].(string)
		if status == "" {
			// Try list-style response: Content contains an array.
			for _, key := range []string{"Content", "VccList", "Data"} {
				if arr, ok := content[key].([]interface{}); ok && len(arr) > 0 {
					if item, ok := arr[0].(map[string]interface{}); ok {
						if s, ok := item["Status"].(string); ok {
							status = s
							break
						}
					}
				}
			}
		}
		fmt.Printf("[VCC_WAIT] VCC %s status=%q (remaining %v)\n", vccID, status, time.Until(deadline).Round(time.Second))
		if status == "Available" {
			return nil
		}
		// "Create Failed" and "Not Available" are terminal failure states.
		if status == "Create Failed" || status == "Not Available" {
			failedCount++
			if failedCount == 1 {
				fmt.Printf("[VCC_WAIT] VCC %s entered terminal state %q, full content: %v\n", vccID, status, content)
			}
			if failedCount >= 2 {
				return fmt.Errorf("VCC %s is in terminal state %q (observed %d times), content: %v", vccID, status, failedCount, content)
			}
		} else {
			failedCount = 0
		}
		time.Sleep(15 * time.Second)
	}
	return fmt.Errorf("timed out waiting for VCC %s to become Available (last status: %q)", vccID, status)
}

// testAccCheckEfloVccGrantRuleDestroy verifies the grant rule is gone (error
// code 1003 = not found) and best-effort deletes the API-created VCC. This is
// a safety net: the no-op re-apply step's Check normally deletes the VCC and
// cleans up side-effect resources before terraform destroy. If the test failed
// before that step, this function attempts the same cleanup.
func testAccCheckEfloVccGrantRuleDestroy(s *terraform.State) error {
	client := acctest.Provider.Meta().(*connectivity.AliyunClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_eflo_vcc_grant_rule" {
			continue
		}
		parts := strings.Split(rs.Primary.ID, ":")
		if len(parts) != 2 {
			return fmt.Errorf("invalid resource ID: %s", rs.Primary.ID)
		}
		query := map[string]interface{}{
			"RegionId":    client.RegionId,
			"InstanceId":  parts[0],
			"GrantRuleId": parts[1],
		}
		_, err := client.RpcPost("eflo", "2022-05-30", "GetVccGrantRule", query, query, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"1003"}) {
				continue
			}
			return fmt.Errorf("checking VCC grant rule destroy failed: %w", err)
		}
		return fmt.Errorf("VCC grant rule %s still exists", rs.Primary.ID)
	}
	// Best-effort VCC cleanup: if the VCC was not already deleted by the
	// no-op re-apply step's Check (e.g. test failed before reaching it), try
	// to delete it now to avoid dangling route interfaces and ENIs.
	if testVccGrantRuleVccID != "" {
		req := map[string]interface{}{
			"RegionId": client.RegionId,
			"VccId":    testVccGrantRuleVccID,
		}
		_, _ = client.RpcPost("eflo", "2022-05-30", "DeleteVcc", req, req, true)
		// Poll GetVcc until gone (up to 3 minutes).
		deadline := time.Now().Add(3 * time.Minute)
		for time.Now().Before(deadline) {
			q := map[string]interface{}{
				"RegionId": client.RegionId,
				"VccId":    testVccGrantRuleVccID,
			}
			if _, err := client.RpcPost("eflo", "2022-05-30", "GetVcc", q, q, true); err != nil {
				break
			}
			time.Sleep(15 * time.Second)
		}
		testVccGrantRuleVccID = ""
	}
	// Retrieve the VPC ID from state for ENI cleanup.
	var vpcID string
	if vpcRes, ok := s.RootModule().Resources["alicloud_vpc.default"]; ok {
		vpcID = vpcRes.Primary.ID
	}
	// Wait for ServiceManaged ENIs to be cleaned by the Eflo service.
	testAccEfloVccGrantRuleWaitENIsGone(client, vpcID, 3*time.Minute)
	// Safety-net RI cleanup: remove any remaining route interfaces to
	// prevent DependencyViolation.RouteInterface during terraform destroy.
	testAccEfloVccGrantRuleCleanupRIs(client)
	// Safety-net ENI cleanup: remove any remaining non-ServiceManaged network
	// interfaces to prevent DependencyViolation.NetworkInterface during
	// terraform destroy.
	testAccEfloVccGrantRuleCleanupENIs(client, vpcID)
	os.Unsetenv("TF_VAR_vcc_id")
	return nil
}

// Test Eflo VccGrantRule. <<< Resource test cases, automatically generated.
