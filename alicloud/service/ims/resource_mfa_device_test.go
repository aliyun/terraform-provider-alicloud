package ims_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudRamMFADevice_basic(t *testing.T) {
	r := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-acc-mfa-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudRamMFADeviceConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_ram_mfa_device.default", "id"),
					resource.TestCheckResourceAttr("alicloud_ram_mfa_device.default", "virtual_mfa_device_name", name),
				),
			},
			{
				ResourceName:      "alicloud_ram_mfa_device.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccAliCloudRamMFADeviceConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%[1]s"
}

resource "alicloud_ram_mfa_device" "default" {
  virtual_mfa_device_name = var.name
}
`, name)
}
