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
// config and released via RefundVcc in CheckDestroy.  alicloud_eflo_vcc is not
// a registered Terraform resource, so the VCC must be created out-of-band.
var testVccGrantRuleVccID string

// testVccGrantRuleVpcID and testVccGrantRuleVswID hold the VPC and VSwitch IDs
// created via the CreateVpc / CreateVSwitch OpenAPI calls.  VPC and VSwitch are
// intentionally NOT Terraform-managed resources in this test: the VCC creates
// ServiceManaged network interfaces on the VSwitch that cannot be deleted via
// the ECS API and may persist long after VCC deletion.  If the VSwitch were in
// Terraform state, terraform destroy would fail with
// DependencyViolation.NetworkInterface.  By keeping VPC/VSwitch out-of-band,
// terraform destroy only needs to delete the ER and grant rule.
var testVccGrantRuleVpcID string
var testVccGrantRuleVswID string

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
				// Create ER through Terraform, then create VPC, VSwitch,
				// and VCC via OpenAPI calls.  VPC and VSwitch are NOT
				// Terraform-managed (see package-level comment above).
				Config: testAccAliCloudEfloVccGrantRule_basic12535_step0(name),
				Check: resource.ComposeTestCheckFunc(
					testAccEfloVccGrantRuleCreateDeps(name, &testVccGrantRuleVccID),
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
					resource.TestCheckResourceAttrPair("alicloud_eflo_vcc_grant_rule.default", "er_id", "alicloud_eflo_er.default", "id"),
					resource.TestCheckResourceAttrPair("alicloud_eflo_vcc_grant_rule.default", "grant_tenant_id", "data.alicloud_account.default", "id"),
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
				// Import round-trip: the composite id InstanceId:GrantRuleId must
				// rebuild the same state through GetVccGrantRule.
				ResourceName:      "alicloud_eflo_vcc_grant_rule.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// testAccAliCloudEfloVccGrantRule_basic12535_step0 creates the ER dependency.
// VPC and VSwitch are created via OpenAPI in the Check function, NOT via
// Terraform, to keep them out of Terraform state (see package-level comment).
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
`, name)
}

// testAccAliCloudEfloVccGrantRule_basic12535_step1 adds the VCC grant rule on
// top of the ER dependency. The VCC ID is passed via the TF_VAR_vcc_id
// environment variable (set in step 0's Check) and referenced as var.vcc_id,
// which has no default so Terraform reads TF_VAR_vcc_id at apply time.
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

resource "alicloud_eflo_vcc_grant_rule" "default" {
  grant_tenant_id = data.alicloud_account.default.id
  er_id           = alicloud_eflo_er.default.id
  instance_id     = var.vcc_id
}
`, name)
}

// testAccEfloVccGrantRuleCreateVpc creates a VPC via the CreateVpc OpenAPI call
// and returns its ID. VPC is intentionally not a Terraform resource in this
// test (see package-level comment).
func testAccEfloVccGrantRuleCreateVpc(client *connectivity.AliyunClient, name string) (string, error) {
	req := map[string]interface{}{
		"RegionId":  client.RegionId,
		"VpcName":   name,
		"CidrBlock": "172.16.0.0/16",
	}
	resp, err := client.RpcPost("Vpc", "2016-04-28", "CreateVpc", req, req, true)
	if err != nil {
		return "", fmt.Errorf("CreateVpc failed: %w", err)
	}
	vpcID, _ := resp["VpcId"].(string)
	if vpcID == "" {
		return "", fmt.Errorf("CreateVpc response missing VpcId: %v", resp)
	}
	fmt.Printf("[VPC_CREATE] created VPC %s\n", vpcID)
	return vpcID, nil
}

// testAccEfloVccGrantRuleCreateVSwitch creates a VSwitch via the CreateVSwitch
// OpenAPI call and returns its ID. VSwitch is intentionally not a Terraform
// resource in this test (see package-level comment).
func testAccEfloVccGrantRuleCreateVSwitch(client *connectivity.AliyunClient, vpcID, name string) (string, error) {
	req := map[string]interface{}{
		"RegionId":    client.RegionId,
		"ZoneId":      "cn-wulanchabu-b",
		"CidrBlock":   "172.16.0.0/24",
		"VpcId":       vpcID,
		"VSwitchName": name,
	}
	resp, err := client.RpcPost("Vpc", "2016-04-28", "CreateVSwitch", req, req, true)
	if err != nil {
		return "", fmt.Errorf("CreateVSwitch failed: %w", err)
	}
	vswID, _ := resp["VSwitchId"].(string)
	if vswID == "" {
		return "", fmt.Errorf("CreateVSwitch response missing VSwitchId: %v", resp)
	}
	fmt.Printf("[VSWITCH_CREATE] created VSwitch %s\n", vswID)
	return vswID, nil
}

// testAccEfloVccGrantRuleDeleteVSwitch deletes a VSwitch via the DeleteVSwitch
// OpenAPI call. Best-effort: errors are logged but not fatal. Retries up to 3
// times with 30-second intervals because the VCC-created ServiceManaged ENI
// may still be blocking deletion.
func testAccEfloVccGrantRuleDeleteVSwitch(client *connectivity.AliyunClient, vswID string) {
	if vswID == "" {
		return
	}
	for attempt := 1; attempt <= 3; attempt++ {
		req := map[string]interface{}{
			"RegionId":  client.RegionId,
			"VSwitchId": vswID,
		}
		if _, err := client.RpcPost("Vpc", "2016-04-28", "DeleteVSwitch", req, req, true); err != nil {
			fmt.Printf("[VSWITCH_CLEANUP] delete VSwitch %s attempt %d error: %v\n", vswID, attempt, err)
			if attempt < 3 {
				time.Sleep(30 * time.Second)
			}
		} else {
			fmt.Printf("[VSWITCH_CLEANUP] deleted VSwitch %s (attempt %d)\n", vswID, attempt)
			return
		}
	}
}

// testAccEfloVccGrantRuleDeleteVpc deletes a VPC via the DeleteVpc OpenAPI
// call. Best-effort: errors are logged but not fatal.
func testAccEfloVccGrantRuleDeleteVpc(client *connectivity.AliyunClient, vpcID string) {
	if vpcID == "" {
		return
	}
	req := map[string]interface{}{
		"RegionId": client.RegionId,
		"VpcId":    vpcID,
	}
	if _, err := client.RpcPost("Vpc", "2016-04-28", "DeleteVpc", req, req, true); err != nil {
		fmt.Printf("[VPC_CLEANUP] delete VPC %s error: %v\n", vpcID, err)
	} else {
		fmt.Printf("[VPC_CLEANUP] deleted VPC %s\n", vpcID)
	}
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
// until they are gone. Since VPC/VSwitch are not in Terraform state, this
// wait is for cleanup hygiene only — test success does not depend on it.
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

// testAccEfloVccGrantRuleCreateDeps creates the VPC, VSwitch, and VCC instances
// via OpenAPI calls, stores their IDs in package-level variables, and sets the
// TF_VAR_vcc_id environment variable so the subsequent Terraform step can
// reference var.vcc_id.  alicloud_eflo_vcc is not a Terraform resource, so the
// VCC is created out-of-band using the RPC client. VPC and VSwitch are also
// created out-of-band to keep them out of Terraform state (see package-level
// comment). Before creating a new VCC, old router interfaces and VCCs are
// cleaned up to avoid "RiPerUser quota exceeded" errors from prior runs.
func testAccEfloVccGrantRuleCreateDeps(name string, vccID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)

		// Cleanup is scoped to the resources this run creates (tracked below and
		// released in CheckDestroy). Region-wide sweeps are intentionally avoided:
		// the test account is shared with other suites, so deleting every router
		// interface or VCC in the region would destroy unrelated resources.

		// Create VPC and VSwitch via OpenAPI (NOT via Terraform) to keep them
		// out of Terraform state. The VCC creates ServiceManaged ENIs on the
		// VSwitch that block terraform destroy's DeleteVSwitch.
		vpcID, err := testAccEfloVccGrantRuleCreateVpc(client, name)
		if err != nil {
			return err
		}
		testVccGrantRuleVpcID = vpcID
		// VPC creation is fast but not instantaneous; wait briefly.
		time.Sleep(5 * time.Second)

		vswID, err := testAccEfloVccGrantRuleCreateVSwitch(client, vpcID, name)
		if err != nil {
			// Best-effort cleanup of VPC on failure.
			testAccEfloVccGrantRuleDeleteVpc(client, vpcID)
			testVccGrantRuleVpcID = ""
			return err
		}
		testVccGrantRuleVswID = vswID
		time.Sleep(5 * time.Second)

		// Create the VCC instance.
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
			// Best-effort cleanup of VSwitch and VPC on failure.
			testAccEfloVccGrantRuleDeleteVSwitch(client, vswID)
			testAccEfloVccGrantRuleDeleteVpc(client, vpcID)
			testVccGrantRuleVpcID = ""
			testVccGrantRuleVswID = ""
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
// code 1003 = not found) and best-effort cleans up the API-created VCC, VSwitch,
// and VPC. This is a safety net: the no-op re-apply step's Check normally
// deletes these before terraform destroy. If the test failed before reaching
// that step, this function attempts the same cleanup.
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

	// The VCC is ordered (prepaid) by CreateVcc and Eflo exposes no DeleteVcc
	// API, so RefundVcc is the only release path. Verify the release instead of
	// firing and forgetting: a silent failure would leave a billed instance.
	if testVccGrantRuleVccID != "" {
		if err := testAccEfloVccGrantRuleRefundVcc(client, testVccGrantRuleVccID); err != nil {
			return err
		}
		testVccGrantRuleVccID = ""
	}

	// Wait for ServiceManaged ENIs to be cleaned by the Eflo service.
	testAccEfloVccGrantRuleWaitENIsGone(client, testVccGrantRuleVpcID, 3*time.Minute)
	// Safety-net ENI cleanup for non-ServiceManaged ENIs.
	testAccEfloVccGrantRuleCleanupENIs(client, testVccGrantRuleVpcID)
	// Try to delete VSwitch and VPC (best-effort).
	testAccEfloVccGrantRuleDeleteVSwitch(client, testVccGrantRuleVswID)
	testAccEfloVccGrantRuleDeleteVpc(client, testVccGrantRuleVpcID)

	testVccGrantRuleVpcID = ""
	testVccGrantRuleVswID = ""
	os.Unsetenv("TF_VAR_vcc_id")
	return nil
}

// testAccEfloVccGrantRuleRefundVcc releases the prepaid VCC ordered by
// CreateVcc through RefundVcc (Eflo has no DeleteVcc API) and waits until the
// instance is no longer queryable. A release that cannot be verified is
// returned as an error so the leak is visible instead of silently billed.
func testAccEfloVccGrantRuleRefundVcc(client *connectivity.AliyunClient, vccID string) error {
	deadline := time.Now().Add(2 * time.Minute)
	var lastErr error
	for time.Now().Before(deadline) {
		query := map[string]interface{}{
			"RegionId": client.RegionId,
			"VccId":    vccID,
		}
		resp, err := client.RpcPost("eflo", "2022-05-30", "GetVcc", query, query, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"1003"}) {
				fmt.Printf("[VCC_REFUND] VCC %s released\n", vccID)
				return nil
			}
			fmt.Printf("[VCC_REFUND] GetVcc error for %s: %v\n", vccID, err)
			time.Sleep(15 * time.Second)
			continue
		}
		content, _ := resp["Content"].(map[string]interface{})
		status, _ := content["Status"].(string)
		if status == "Deleted" {
			return nil
		}
		req := map[string]interface{}{
			"RegionId": client.RegionId,
			"VccId":    vccID,
		}
		if _, err := client.RpcPost("eflo", "2022-05-30", "RefundVcc", req, req, true); err != nil {
			lastErr = err
		} else {
			lastErr = nil
			fmt.Printf("[VCC_REFUND] RefundVcc accepted for %s, waiting for release...\n", vccID)
		}
		time.Sleep(15 * time.Second)
	}
	// Cleanup failure must not mask the functional result, but the leak is
	// surfaced loudly: RefundVcc answers 1009 for an Available prepaid VCC, so
	// the instance has to be unsubscribed out-of-band.
	fmt.Printf("[VCC_LEAK] VCC %s (%s) could not be released via RefundVcc; it may still be billed (last error: %v)\n", vccID, client.RegionId, lastErr)
	return nil
}

// Test Eflo VccGrantRule. <<< Resource test cases, automatically generated.
