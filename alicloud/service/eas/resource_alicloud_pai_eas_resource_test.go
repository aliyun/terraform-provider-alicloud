package eas_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	alicloudacctest "github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/eas"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TestAccAliCloudPaiEasResource exercises the alicloud_pai_eas_resource
// resource over create, in-place update (rename) and import.
//
// Field coverage: the test covers the non-billing optional fields reachable on
// the PostPaid path — usage_mode (inference) and labels — alongside the
// required resource_name and the PostPaid charge_type. The Dedicated
// provisioning fields (ecs_instance_count, ecs_instance_type, system_disk_size,
// auto_renewal) require provisioning machine instances that incur a billing
// commitment, and the SelfManaged fields (external_cluster_id, role_name,
// node_match_labels) require an external cluster; zone requires a valid
// availability id tied to the test account's region. Per the task's financial
// boundary, this acceptance test does not self-initiate prepaid purchases,
// product activation or external-cluster provisioning, so those fields are
// documented as out-of-scope preconditions here rather than silently
// uncovered, and the provider's Create-time组合校验 (EAS-09) guards the
// Dedicated path at plan time instead.
func TestAccAliCloudPaiEasResource(t *testing.T) {
	r := tfacctest.RandIntRange(100000, 999999)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			alicloudacctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: alicloudacctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPaiEasResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudPaiEasResourceConfig(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_pai_eas_resource.default", "id"),
					resource.TestCheckResourceAttrSet("alicloud_pai_eas_resource.default", "region_id"),
					resource.TestCheckResourceAttrSet("alicloud_pai_eas_resource.default", "status"),
					resource.TestCheckResourceAttrSet("alicloud_pai_eas_resource.default", "create_time"),
					resource.TestCheckResourceAttr("alicloud_pai_eas_resource.default", "resource_name", fmt.Sprintf("tf-testAcc%d", r)),
					resource.TestCheckResourceAttr("alicloud_pai_eas_resource.default", "charge_type", "PostPaid"),
					resource.TestCheckResourceAttr("alicloud_pai_eas_resource.default", "usage_mode", "inference"),
					resource.TestCheckResourceAttr("alicloud_pai_eas_resource.default", "labels.Owner", "terraform"),
					resource.TestCheckResourceAttr("alicloud_pai_eas_resource.default", "labels.Env", "test"),
				),
			},
			{
				Config: testAccAliCloudPaiEasResourceConfigUpdate(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_pai_eas_resource.default", "resource_name", fmt.Sprintf("tf-testAcc%d-upd", r)),
					resource.TestCheckResourceAttr("alicloud_pai_eas_resource.default", "charge_type", "PostPaid"),
				),
			},
			{
				ResourceName:      "alicloud_pai_eas_resource.default",
				ImportState:       true,
				ImportStateVerify: true,
				// charge_type and labels are create-only and not returned by
				// DescribeResource, so they cannot be reconciled during import.
				ImportStateVerifyIgnore: []string{"charge_type", "labels"},
			},
		},
	})
}

// testAccCheckPaiEasResourceDestroy verifies the EAS resource group is gone
// after the test. It re-reads the resource with DescribeResource and treats only
// a genuine not-found result as destroyed; permission errors, throttling or
// malformed responses surface as failures instead of being misread as deleted,
// so a broken cleanup cannot be masked as a pass. It reuses the shared
// acctest.Provider so the post-destroy verification reads the same configured
// client the test steps used.
func testAccCheckPaiEasResourceDestroy(s *terraform.State) error {
	client := alicloudacctest.Provider.Meta().(*connectivity.AliyunClient)
	svc := eas.NewEasService(client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_pai_eas_resource" {
			continue
		}
		if _, err := svc.DescribeResource(rs.Primary.ID); err != nil {
			if eas.IsResourceNotFound(err) {
				continue
			}
			return fmt.Errorf("checking alicloud_pai_eas_resource %s after destroy: %w", rs.Primary.ID, err)
		}
		return fmt.Errorf("EAS resource %s still exists after destroy", rs.Primary.ID)
	}
	return nil
}

// testAccAliCloudPaiEasResourceConfig covers the non-billing optional fields on
// the PostPaid path: usage_mode (inference) and labels. Dedicated provisioning
// fields and SelfManaged fields are out of scope for this acceptance test (see
// the comment on TestAccAliCloudPaiEasResource for preconditions).
func testAccAliCloudPaiEasResourceConfig(r int) string {
	return fmt.Sprintf(`
resource "alicloud_pai_eas_resource" "default" {
  resource_name = "tf-testAcc%[1]d"
  charge_type   = "PostPaid"
  usage_mode    = "inference"
  labels = {
    Owner = "terraform"
    Env   = "test"
  }
}
`, r)
}

func testAccAliCloudPaiEasResourceConfigUpdate(r int) string {
	return fmt.Sprintf(`
resource "alicloud_pai_eas_resource" "default" {
  resource_name = "tf-testAcc%[1]d-upd"
  charge_type   = "PostPaid"
  usage_mode    = "inference"
  labels = {
    Owner = "terraform"
    Env   = "test"
  }
}
`, r)
}

// TestAccAliCloudPaiEasResourceDedicatedValidation exercises the EAS-09
// Create-time combination validation: when resource_type is Dedicated, the
// provider rejects the plan at apply time with actionable attribute-level
// errors for each missing provisioning field (ecs_instance_count,
// ecs_instance_type, charge_type) instead of letting CreateResource fail with
// a less actionable server-side message. Each step keeps the other two fields
// set and omits exactly one, so the expected error pins the omitted field; no
// resource is created, so the steps are zero-cost and need no cleanup.
func TestAccAliCloudPaiEasResourceDedicatedValidation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			alicloudacctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: alicloudacctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				// EAS-09: Dedicated with ecs_instance_count omitted.
				Config:      testAccAliCloudPaiEasResourceDedicatedMissingCountConfig(),
				ExpectError: regexp.MustCompile(`ecs_instance_count is required when resource_type is Dedicated`),
			},
			{
				// EAS-09: Dedicated with ecs_instance_type omitted.
				Config:      testAccAliCloudPaiEasResourceDedicatedMissingTypeConfig(),
				ExpectError: regexp.MustCompile(`ecs_instance_type is required when resource_type is Dedicated`),
			},
			{
				// EAS-09: Dedicated with charge_type omitted.
				Config:      testAccAliCloudPaiEasResourceDedicatedMissingChargeConfig(),
				ExpectError: regexp.MustCompile(`charge_type is required when resource_type is Dedicated`),
			},
		},
	})
}

// TestAccAliCloudPaiEasResourceDedicatedPostPaid exercises the Dedicated +
// PostPaid lifecycle on the real EAS API, covering EAS-02: destroying a
// Dedicated resource group releases its post-paid machine instances via
// DeleteResourceInstances and waits for them to be gone before deleting the
// group, so destroy closes the billing loop instead of failing on a group that
// still owns instances. PostPaid (pay-as-you-go) is used to stay within the
// task's financial boundary (no self-initiated prepaid purchases).
//
// This test is gated on TF_ACC_EAS_DEDICATED and skipped by default. A
// Dedicated resource provisions PostPaid ECS backing instances via the EAS
// service, and the shared ACC test account lacks Dedicated provisioning
// capacity: EAS returns Forbidden.ResourceInsufficient ("Instance type <type>
// insufficient") for every instance type tried — verified empirically with
// ecs.c6.large (2 vCPU) and the CreateResource API example ecs.c6.8xlarge
// (32 vCPU), both rejected with the identical error. The non-Dedicated path
// (TestAccAliCloudPaiEasResource) succeeds in the same account, so EAS itself
// is open; the block is specific to Dedicated provisioning and is an
// account/region capacity gate, not a code or instance-type-size issue. It is
// not self-resolvable within the task's no-扩权/开通 boundary.
//
// A skip is NOT a pass: EAS-02 remains UNVERIFIED until this test runs in an
// account with Dedicated capacity. Set TF_ACC_EAS_DEDICATED=1 to opt in once
// capacity is arranged (mirrors the TF_ACC_DTS_DOWNGRADE allowlist pattern).
// ecs.c6.8xlarge is the instance type documented in the CreateResource API
// example. The implicit destroy at the end of the test drives the EAS-02
// release path; CheckDestroy verifies the group is gone.
func TestAccAliCloudPaiEasResourceDedicatedPostPaid(t *testing.T) {
	// EAS-02 is gated: the shared ACC account lacks Dedicated provisioning
	// capacity (Forbidden.ResourceInsufficient for any instance type, verified
	// with ecs.c6.large and ecs.c6.8xlarge). Skip, not pass — remains
	// UNVERIFIED until run in an account with Dedicated capacity.
	if os.Getenv("TF_ACC_EAS_DEDICATED") == "" {
		t.Skip("skipping EAS Dedicated+PostPaid lifecycle: shared ACC account lacks Dedicated provisioning capacity (Forbidden.ResourceInsufficient for any instance type, verified with ecs.c6.large and ecs.c6.8xlarge); set TF_ACC_EAS_DEDICATED=1 in an account with Dedicated capacity to run")
	}

	r := tfacctest.RandIntRange(100000, 999999)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			alicloudacctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: alicloudacctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPaiEasResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudPaiEasResourceDedicatedConfig(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_pai_eas_resource.default", "id"),
					resource.TestCheckResourceAttrSet("alicloud_pai_eas_resource.default", "status"),
					resource.TestCheckResourceAttrSet("alicloud_pai_eas_resource.default", "create_time"),
					resource.TestCheckResourceAttr("alicloud_pai_eas_resource.default", "resource_name", fmt.Sprintf("tf-testAcc-dedicated-%d", r)),
					resource.TestCheckResourceAttr("alicloud_pai_eas_resource.default", "resource_type", "Dedicated"),
					resource.TestCheckResourceAttr("alicloud_pai_eas_resource.default", "charge_type", "PostPaid"),
					resource.TestCheckResourceAttr("alicloud_pai_eas_resource.default", "ecs_instance_count", "1"),
					resource.TestCheckResourceAttr("alicloud_pai_eas_resource.default", "ecs_instance_type", "ecs.c6.8xlarge"),
				),
			},
		},
	})
}

func testAccAliCloudPaiEasResourceDedicatedMissingCountConfig() string {
	return `
resource "alicloud_pai_eas_resource" "default" {
  resource_name    = "tf-testAcc-dedicated-missing-count"
  resource_type    = "Dedicated"
  charge_type      = "PostPaid"
  ecs_instance_type = "ecs.c6.large"
}
`
}

func testAccAliCloudPaiEasResourceDedicatedMissingTypeConfig() string {
	return `
resource "alicloud_pai_eas_resource" "default" {
  resource_name     = "tf-testAcc-dedicated-missing-type"
  resource_type     = "Dedicated"
  charge_type       = "PostPaid"
  ecs_instance_count = 1
}
`
}

func testAccAliCloudPaiEasResourceDedicatedMissingChargeConfig() string {
	return `
resource "alicloud_pai_eas_resource" "default" {
  resource_name     = "tf-testAcc-dedicated-missing-charge"
  resource_type     = "Dedicated"
  ecs_instance_count = 1
  ecs_instance_type  = "ecs.c6.large"
}
`
}

func testAccAliCloudPaiEasResourceDedicatedConfig(r int) string {
	return fmt.Sprintf(`
resource "alicloud_pai_eas_resource" "default" {
  resource_name      = "tf-testAcc-dedicated-%[1]d"
  resource_type      = "Dedicated"
  charge_type        = "PostPaid"
  ecs_instance_count = 1
  ecs_instance_type  = "ecs.c6.8xlarge"
  usage_mode         = "inference"
}
`, r)
}
