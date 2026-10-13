package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudRamMfaDevices_basic(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-acc-ds-mfa-%d", rand)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAlicloudRamMfaDevicesConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_ram_mfa_devices.default", "mfa_devices.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_ram_mfa_devices.default", "mfa_devices.0.serial_number"),
					resource.TestCheckResourceAttrSet("data.alicloud_ram_mfa_devices.default", "mfa_devices.0.virtual_mfa_device_name"),
				),
			},
		},
	})
}

func testAccAlicloudRamMfaDevicesConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_ram_mfa_device" "default" {
  virtual_mfa_device_name = var.name
}

data "alicloud_ram_mfa_devices" "default" {
  name_regex = "^${alicloud_ram_mfa_device.default.virtual_mfa_device_name}$"
}
`, name)
}
