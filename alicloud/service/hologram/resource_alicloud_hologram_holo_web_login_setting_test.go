package hologram_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/hologram"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// The HoloWeb login setting is a boolean toggle bound to one Hologres
// instance. The basic lifecycle covers enabling the setting (Create), toggling
// it off (Update), importing by instance id and reverting it to false on
// Destroy. The Hologres instance is created as a pay-as-you-go dependency and
// is removed by the same test sweep.
func TestAccAliCloudHologramHoloWebLoginSetting_basic(t *testing.T) {
	resourceId := "alicloud_hologram_holo_web_login_setting.default"
	rand := tfacctest.RandIntRange(10000, 99999)
	region := os.Getenv("ALICLOUD_REGION")
	if region == "" {
		region = "cn-hangzhou"
	}
	name := fmt.Sprintf("tf-testacc%shologramholowebloginsetting%d", region, rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.HologramSupportRegions)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckHologramHoloWebLoginSettingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccHologramHoloWebLoginSettingBasic(name, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceId, "instance_id"),
					resource.TestCheckResourceAttr(resourceId, "allow_external_accounts_login", "true"),
				),
			},
			{
				Config: testAccHologramHoloWebLoginSettingBasic(name, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "allow_external_accounts_login", "false"),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccHologramHoloWebLoginSettingBasic(name string, allowExternal bool) string {
	return fmt.Sprintf(`
%[1]s

resource "alicloud_hologram_instance" "default" {
  instance_type = "Standard"
  pricing_cycle = "Hour"
  cpu           = "32"
  instance_name = var.name
  payment_type  = "PayAsYouGo"
  zone_id       = alicloud_vswitch.defaultVSwitch.zone_id
  endpoints {
    type = "Intranet"
  }
  endpoints {
    type       = "VPCSingleTunnel"
    vswitch_id = alicloud_vswitch.defaultVSwitch.id
    vpc_id     = alicloud_vswitch.defaultVSwitch.vpc_id
  }
}

resource "alicloud_hologram_holo_web_login_setting" "default" {
  instance_id                   = alicloud_hologram_instance.default.id
  allow_external_accounts_login  = %[2]t
}
`, AlicloudHologramInstanceBasicDependence3920(name), allowExternal)
}

func testAccCheckHologramHoloWebLoginSettingDestroy(s *terraform.State) error {
	client := acctest.Provider.Meta().(*connectivity.AliyunClient)
	service := hologram.NewService(client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_hologram_holo_web_login_setting" {
			continue
		}
		instanceId := rs.Primary.Attributes["instance_id"]
		if instanceId == "" {
			instanceId = rs.Primary.ID
		}
		object, err := service.DescribeHoloWebLoginSetting(instanceId)
		if err != nil {
			// The instance is destroyed in the same sweep; once it is gone the
			// toggle is gone too.
			if errs.NotFoundError(err) {
				continue
			}
			return fmt.Errorf("describing hologram holo web login setting %s failed: %s", instanceId, err)
		}
		// Destroy reverts the toggle to false; a still-true setting means Delete
		// did not revert.
		if v, ok := object["AllowExternalAccountsLogin"].(bool); ok && !v {
			continue
		}
		return fmt.Errorf("hologram holo web login setting %s still allows external accounts login", instanceId)
	}
	return nil
}
