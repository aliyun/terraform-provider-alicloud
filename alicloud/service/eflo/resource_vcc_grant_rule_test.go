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
				ResourceName:            "alicloud_eflo_vcc_grant_rule.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
				// Delete the API-created VCC before Terraform's post-test
				// destroy runs. The VCC creates a route interface on the
				// VPC; without deleting the VCC first, DeleteVpc fails with
				// DependencyViolation.RouteInterface.
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

// testAccEfloVccGrantRuleCreateVcc creates a VCC instance via the CreateVcc
// OpenAPI call, stores its ID, and sets the TF_VAR_vcc_id environment variable
// so the subsequent Terraform step can reference var.vcc_id.  alicloud_eflo_vcc
// is not a Terraform resource, so the VCC is created out-of-band using the RPC
// client.
func testAccEfloVccGrantRuleCreateVcc(name string, vccID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)

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
			"ZoneId":         "cn-wulanchabu-a",
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
// that Terraform's post-test destroy can delete the VPC. The VCC creates a
// route interface on the VPC; without deleting the VCC first, DeleteVpc fails
// with DependencyViolation.RouteInterface. This Check runs in the import step,
// before the framework's automatic terraform destroy.
func testAccEfloVccGrantRuleDeleteVcc(vccID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if *vccID == "" {
			return nil
		}
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)
		req := map[string]interface{}{
			"RegionId": client.RegionId,
			"VccId":    *vccID,
		}
		_, _ = client.RpcPost("eflo", "2022-05-30", "DeleteVcc", req, req, true)
		*vccID = ""
		os.Unsetenv("TF_VAR_vcc_id")
		return nil
	}
}

// testAccEfloVccGrantRuleWaitVcc polls GetVcc until the VCC is Available.
// VCC Status values (per GetVcc API): Available, Not Available, Executing,
// Deleting. Only "Available" means the VCC is operable for grant rule
// creation; any other status causes CreateVccGrantRule to fail with
// Code 1009 (un-operable state). VCC provisioning latency is intermittent
// in cn-wulanchabu, so the timeout is 10 minutes.
func testAccEfloVccGrantRuleWaitVcc(client *connectivity.AliyunClient, vccID string) error {
	query := map[string]interface{}{
		"RegionId": client.RegionId,
		"VccId":    vccID,
	}
	timeout := 10 * time.Minute
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		resp, err := client.RpcPost("eflo", "2022-05-30", "GetVcc", query, query, true)
		if err != nil {
			// The VCC may not be queryable immediately after creation.
			time.Sleep(10 * time.Second)
			continue
		}
		content, _ := resp["Content"].(map[string]interface{})
		if content == nil {
			time.Sleep(10 * time.Second)
			continue
		}
		// GetVcc may return a single object or a list wrapper.
		if status, ok := content["Status"].(string); ok && status == "Available" {
			return nil
		}
		// Try list-style response: Content contains an array.
		for _, key := range []string{"Content", "VccList", "Data"} {
			if arr, ok := content[key].([]interface{}); ok && len(arr) > 0 {
				if item, ok := arr[0].(map[string]interface{}); ok {
					if status, ok := item["Status"].(string); ok && status == "Available" {
						return nil
					}
				}
			}
		}
		time.Sleep(10 * time.Second)
	}
	return fmt.Errorf("timed out waiting for VCC %s to become Available", vccID)
}

// testAccCheckEfloVccGrantRuleDestroy verifies the grant rule is gone (error
// code 1003 = not found) and best-effort deletes the API-created VCC.
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
	// import step's Check (e.g. test failed before reaching import), try to
	// delete it now to avoid dangling route interfaces on the VPC.
	if testVccGrantRuleVccID != "" {
		req := map[string]interface{}{
			"RegionId": client.RegionId,
			"VccId":    testVccGrantRuleVccID,
		}
		_, _ = client.RpcPost("eflo", "2022-05-30", "DeleteVcc", req, req, true)
		testVccGrantRuleVccID = ""
	}
	os.Unsetenv("TF_VAR_vcc_id")
	return nil
}

// Test Eflo VccGrantRule. <<< Resource test cases, automatically generated.
