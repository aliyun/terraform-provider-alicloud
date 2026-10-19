package ims_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudRamMfaDevices_basic(t *testing.T) {
	r := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-acc-ds-mfa-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudRamMfaDevicesConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_ram_mfa_devices.default", "mfa_devices.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_ram_mfa_devices.default", "mfa_devices.0.serial_number"),
					resource.TestCheckResourceAttrSet("data.alicloud_ram_mfa_devices.default", "mfa_devices.0.virtual_mfa_device_name"),
				),
			},
		},
	})
}

func testAccAliCloudRamMfaDevicesConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%[1]s"
}

resource "alicloud_ram_mfa_device" "default" {
  virtual_mfa_device_name = var.name
}

data "alicloud_ram_mfa_devices" "default" {
  name_regex = "^${alicloud_ram_mfa_device.default.virtual_mfa_device_name}$"
}
`, name)
}
