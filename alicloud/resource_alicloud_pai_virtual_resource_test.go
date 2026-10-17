package alicloud_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccAliCloudPaiVirtualResource_basic exercises the create, read, update and
// import paths of the framework-native alicloud_pai_virtual_resource resource
// behind the muxed provider served by acctest.ProtoV5ProviderFactories.
//
// The test stays in package alicloud_test and uses the SDK v2 helper/resource
// TestCase (whose ProtoV5ProviderFactories field carries the muxed tfprotov5
// server factory map) so the binary links a single sweep init and drives the
// native framework resource through the mux. acctest.PreCheck supplies
// credentials; CheckDestroy is left nil (a deliberate no-op for this grouping
// resource, equivalent to acctest.CheckDestroyNoop but typed for this harness).
func TestAccAliCloudPaiVirtualResource_basic(t *testing.T) {
	rand := sdkacctest.RandInt()
	name := fmt.Sprintf("tf-testacc-paivr-%d", rand)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccPaiVirtualResourceConfig(name, false, 6, 10.05),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_pai_virtual_resource.example", "id"),
					resource.TestCheckResourceAttrSet("alicloud_pai_virtual_resource.example", "virtual_resource_id"),
					resource.TestCheckResourceAttr("alicloud_pai_virtual_resource.example", "virtual_resource_name", name),
					resource.TestCheckResourceAttr("alicloud_pai_virtual_resource.example", "cluster_id", "cn-hangzhou"),
					resource.TestCheckResourceAttr("alicloud_pai_virtual_resource.example", "disable_spot_protection_period", "false"),
					resource.TestCheckResourceAttr("alicloud_pai_virtual_resource.example", "resources.#", "1"),
					resource.TestCheckResourceAttr("alicloud_pai_virtual_resource.example", "resources.0.instance_type", "ecs.s6-c1m2.xlarge"),
					resource.TestCheckResourceAttr("alicloud_pai_virtual_resource.example", "resources.0.region", "cn-hangzhou"),
					resource.TestCheckResourceAttr("alicloud_pai_virtual_resource.example", "resources.0.priority", "6"),
					resource.TestCheckResourceAttr("alicloud_pai_virtual_resource.example", "resources.0.spot_price_limit", "10.05"),
				),
			},
			{
				Config: testAccPaiVirtualResourceConfig(name, true, 9, 5.5),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_pai_virtual_resource.example", "disable_spot_protection_period", "true"),
					resource.TestCheckResourceAttr("alicloud_pai_virtual_resource.example", "resources.0.priority", "9"),
					resource.TestCheckResourceAttr("alicloud_pai_virtual_resource.example", "resources.0.spot_price_limit", "5.5"),
				),
			},
			{
				ResourceName:      "alicloud_pai_virtual_resource.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccPaiVirtualResourceConfig(name string, disableSpotProtection bool, priority int, spotPriceLimit float64) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%[1]s"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_pai_virtual_resource" "example" {
  virtual_resource_name          = var.name
  disable_spot_protection_period = %[2]t

  resources {
    instance_type    = "ecs.s6-c1m2.xlarge"
    region           = "cn-hangzhou"
    priority         = %[3]d
    spot_price_limit = %[4]f
  }
}
`, name, disableSpotProtection, priority, spotPriceLimit)
}
