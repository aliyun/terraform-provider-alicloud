package hologram_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccAliCloudHologramUser exercises the Hologram User resource end to end
// across the full CRUD surface: create with super_user=false, flip super_user to
// true (ForceNew recreate, exercising Delete+Create since Hologram 2022-06-01
// has no UpdateUser), and import. The import step ignores super_user because
// ListUsers never returns it.
func TestAccAliCloudHologramUser_basic(t *testing.T) {
	r := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tftestaccuser%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			// Hologram only supports cn-hangzhou; force it ahead of the credential
			// check so the SDK v2 default (cn-beijing) does not steer the run.
			os.Setenv("ALICLOUD_REGION", "cn-hangzhou")
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccHologramUserConfig(name, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_hologram_user.default", "id"),
					resource.TestCheckResourceAttr("alicloud_hologram_user.default", "user_name", name),
					resource.TestCheckResourceAttrSet("alicloud_hologram_user.default", "instance_id"),
					resource.TestCheckResourceAttr("alicloud_hologram_user.default", "super_user", "false"),
				),
			},
			{
				Config: testAccHologramUserConfig(name, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_hologram_user.default", "user_name", name),
					resource.TestCheckResourceAttr("alicloud_hologram_user.default", "super_user", "true"),
				),
			},
			{
				ResourceName:      "alicloud_hologram_user.default",
				ImportState:       true,
				ImportStateVerify: true,
				// super_user is not returned by ListUsers, so it cannot be
				// verified against state after a raw import.
				ImportStateVerifyIgnore: []string{"super_user"},
			},
		},
	})
}

func testAccHologramUserConfig(name string, superUser bool) string {
	flag := "false"
	if superUser {
		flag = "true"
	}
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
  super_user  = %s
}
`, name, flag)
}
