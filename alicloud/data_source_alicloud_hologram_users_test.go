package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudHologramUsersDataSource_basic(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tftestacchologramuser%d", rand)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.HologramSupportRegions)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAlicloudHologramUsersDataSourceConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_hologram_users.default", "users.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_hologram_users.default", "users.0.user_name", name),
					resource.TestCheckResourceAttrSet("data.alicloud_hologram_users.default", "users.0.id"),
					resource.TestCheckResourceAttrSet("data.alicloud_hologram_users.default", "users.0.super_user"),
				),
			},
		},
	})
}

func testAccAlicloudHologramUsersDataSourceConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "defaultVpc" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = var.name
}

resource "alicloud_vswitch" "defaultVswitch" {
  vpc_id       = alicloud_vpc.defaultVpc.id
  zone_id      = "cn-hangzhou-j"
  cidr_block   = "172.16.53.0/24"
  vswitch_name = var.name
}

resource "alicloud_hologram_instance" "default" {
  zone_id       = alicloud_vswitch.defaultVswitch.zone_id
  instance_name = var.name
  payment_type  = "PayAsYouGo"
  instance_type = "Warehouse"
  pricing_cycle = "Hour"
  cpu           = "32"
  gateway_count = "2"
  endpoints {
    type = "Intranet"
  }
  endpoints {
    type       = "VPCSingleTunnel"
    vswitch_id = alicloud_vswitch.defaultVswitch.id
    vpc_id     = alicloud_vswitch.defaultVswitch.vpc_id
  }
}

resource "alicloud_hologram_user" "default" {
  instance_id = alicloud_hologram_instance.default.id
  user_name   = var.name
  super_user  = false
}

data "alicloud_hologram_users" "default" {
  instance_id = alicloud_hologram_instance.default.id
  ids          = [alicloud_hologram_user.default.id]
}
`, name)
}
