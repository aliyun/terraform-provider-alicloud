package alicloud_test

import (
	"fmt"
	"os"
	"testing"

	acc "github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// EAS virtual resource groups are exercised in cn-hangzhou, the region the
// existing alicloud_pai_service acceptance tests run in. ProtoV5ProviderFactories
// serves the muxed provider (SDK v2 half + framework half), so the framework-native
// alicloud_pai_virtual_resource resource is reachable from this test.
//
// The SDK v2 helper/resource harness is used (not terraform-plugin-testing/helper/resource)
// because non-test files in package alicloud still import the SDK helper/resource, whose
// init registers a -sweep flag; linking both harnesses panics at startup with
// "flag redefined: sweep". The SDK harness TestCase supports ProtoV5ProviderFactories.
func testAccPreCheckPaiVirtualResource(t *testing.T) {
	os.Setenv("ALICLOUD_REGION", "cn-hangzhou")
	acc.PreCheck(t)
}

// TestAccAliCloudPaiVirtualResource_basic covers the Create/Read/Update/Delete
// lifecycle and import of a virtual resource group backed by a public resource
// group instance type (the InstanceType path). This path needs no pre-existing
// resource, so SpotPriceLimit and DisableSpotProtectionPeriod — two of the
// governance-mapped parameters — are round-tripped end to end.
//
// QuotaId and ResourceId identify a Lingjun quota / EAS dedicated resource group
// respectively; neither is constructible in the test account, so they are wired
// in the schema and Create/Update/Read mappings but not exercised here. See the
// PR description for that coverage gap.
func TestAccAliCloudPaiVirtualResource_basic(t *testing.T) {
	r := tfacctest.RandIntRange(10000, 99999)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckPaiVirtualResource(t)
		},
		ProtoV5ProviderFactories: acc.ProtoV5ProviderFactories(),
		// CheckDestroy satisfies tfproviderlint AT001. Destroy itself is
		// exercised by the harness Destroy step: Delete calls
		// DeleteVirtualResource and suppresses the not-found error, so a
		// resource that no longer exists removes cleanly. A local SDK-typed
		// noop is used because acc.CheckDestroyNoop is plugin-testing-typed
		// (its *terraform.State is terraform-plugin-testing/terraform.State),
		// which does not assign to the SDKv2 helper/resource.TestCheckFunc
		// this TestCase uses; the local func matches the SDKv2 terraform.State.
		CheckDestroy: func(*terraform.State) error { return nil },
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudPaiVirtualResourceConfig(r, true, 6, 10.5, "ecs.s6-c1m2.xlarge"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_pai_virtual_resource.default", "id"),
					resource.TestCheckResourceAttrSet("alicloud_pai_virtual_resource.default", "virtual_resource_id"),
					resource.TestCheckResourceAttr("alicloud_pai_virtual_resource.default", "virtual_resource_name", fmt.Sprintf("tf-acc-%d", r)),
					resource.TestCheckResourceAttr("alicloud_pai_virtual_resource.default", "disable_spot_protection_period", "true"),
					resource.TestCheckResourceAttr("alicloud_pai_virtual_resource.default", "region_id", "cn-hangzhou"),
					resource.TestCheckResourceAttrSet("alicloud_pai_virtual_resource.default", "create_time"),
					resource.TestCheckResourceAttr("alicloud_pai_virtual_resource.default", "resources.#", "1"),
					resource.TestCheckResourceAttr("alicloud_pai_virtual_resource.default", "resources.0.region", "cn-hangzhou"),
					resource.TestCheckResourceAttr("alicloud_pai_virtual_resource.default", "resources.0.instance_type", "ecs.s6-c1m2.xlarge"),
					resource.TestCheckResourceAttr("alicloud_pai_virtual_resource.default", "resources.0.priority", "6"),
					resource.TestCheckResourceAttr("alicloud_pai_virtual_resource.default", "resources.0.spot_price_limit", "10.5"),
				),
			},
			{
				Config: testAccAliCloudPaiVirtualResourceConfig(r, false, 8, 5.5, "ecs.s6-c1m2.xlarge"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_pai_virtual_resource.default", "disable_spot_protection_period", "false"),
					resource.TestCheckResourceAttr("alicloud_pai_virtual_resource.default", "resources.0.priority", "8"),
					resource.TestCheckResourceAttr("alicloud_pai_virtual_resource.default", "resources.0.spot_price_limit", "5.5"),
				),
			},
			{
				// Priority and SpotPriceLimit are input-only for the public resource
				// group path: the EAS API accepts them on Create/Update but does not
				// echo them in DescribeVirtualResource, so a fresh import cannot
				// recover them. They are preserved across Create/Update/Read (see
				// virtualResourceState), but excluded from import verification.
				ResourceName:            "alicloud_pai_virtual_resource.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"resources.0.priority", "resources.0.spot_price_limit"},
			},
		},
	})
}

func testAccAliCloudPaiVirtualResourceConfig(r int, disableSpot bool, priority int, spotPriceLimit float64, instanceType string) string {
	return fmt.Sprintf(`
resource "alicloud_pai_virtual_resource" "default" {
  virtual_resource_name          = "tf-acc-%[1]d"
  disable_spot_protection_period = %[2]t
  resources {
    region           = "cn-hangzhou"
    instance_type    = "%[3]s"
    priority         = %[4]d
    spot_price_limit = %[5]f
  }
}
`, r, disableSpot, instanceType, priority, spotPriceLimit)
}
