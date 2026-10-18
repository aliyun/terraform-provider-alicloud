package hologram_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccAliCloudHologramUsersDataSource_basic creates a Hologram User and reads
// it back through the data source, filtered by id. ListUsers returns each user's
// name as "useName" and nothing else, so super_user is null; the assertions cover
// the id and user_name only.
func TestAccAliCloudHologramUsersDataSource_basic(t *testing.T) {
	r := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tftestaccuser%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			os.Setenv("ALICLOUD_REGION", "cn-hangzhou")
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccHologramUsersDataSourceConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_hologram_users.default", "users.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_hologram_users.default", "users.0.user_name", name),
					resource.TestCheckResourceAttrSet("data.alicloud_hologram_users.default", "users.0.id"),
					resource.TestCheckResourceAttr("data.alicloud_hologram_users.default", "names.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_hologram_users.default", "ids.#", "1"),
				),
			},
		},
	})
}

func testAccHologramUsersDataSourceConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_vpc" "defaultVpc" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = var.name
}

resource "alicloud_vswitch" "defaultVSwitch" {
  vpc_id       = alicloud_vpc.defaultVpc.id
  zone_id      = "cn-hangzhou-j"
  cidr_block   = "172.16.53.0/24"
  vswitch_name = var.name
}

resource "alicloud_hologram_instance" "default" {
  zone_id       = alicloud_vswitch.defaultVSwitch.zone_id
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
    vswitch_id = alicloud_vswitch.defaultVSwitch.id
    vpc_id     = alicloud_vswitch.defaultVSwitch.vpc_id
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
