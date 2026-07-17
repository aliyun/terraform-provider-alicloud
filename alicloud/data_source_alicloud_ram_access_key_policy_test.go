package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudRamAccessKeyPolicyDataSource(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccramakpds%d", rand)
	resourceId := "data.alicloud_ram_access_key_policy.default"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudRamAccessKeyPolicyDataSourceConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceId, "access_key_policy"),
					resource.TestCheckResourceAttrSet(resourceId, "user_access_key_id"),
					resource.TestCheckResourceAttrSet(resourceId, "user_principal_name"),
				),
			},
		},
	})
}

func testAccAliCloudRamAccessKeyPolicyDataSourceConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_account" "default" {
}

resource "alicloud_ram_user" "default" {
  name = var.name
}

resource "alicloud_ram_access_key" "default" {
  user_name = alicloud_ram_user.default.name
}

resource "alicloud_ram_access_key_policy" "default" {
  user_access_key_id  = alicloud_ram_access_key.default.id
  user_principal_name = "${alicloud_ram_user.default.name}@${data.alicloud_account.default.id}.onaliyun.com"
  access_key_policy   = "{\"Version\":1,\"Status\":\"Active\",\"Statements\":[{\"Type\":\"ClassicWhiteList\",\"IPList\":[\"10.0.0.1/32\"]}]}"
}

data "alicloud_ram_access_key_policy" "default" {
  user_access_key_id  = alicloud_ram_access_key_policy.default.user_access_key_id
  user_principal_name = alicloud_ram_access_key_policy.default.user_principal_name
}
`, name)
}
