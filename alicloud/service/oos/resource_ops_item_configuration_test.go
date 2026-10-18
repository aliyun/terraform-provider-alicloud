package oos_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudOosOpsItemConfiguration(t *testing.T) {
	r := tfacctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudOosOpsItemConfigurationConfig(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_oos_ops_item_configuration.default", "id"),
					resource.TestCheckResourceAttrSet("alicloud_oos_ops_item_configuration.default", "configuration_id"),
					resource.TestCheckResourceAttr("alicloud_oos_ops_item_configuration.default", "ops_item_configuration_name", fmt.Sprintf("tf-testacc%d", r)),
					resource.TestCheckResourceAttrSet("alicloud_oos_ops_item_configuration.default", "create_time"),
					resource.TestCheckResourceAttrSet("alicloud_oos_ops_item_configuration.default", "region_id"),
				),
			},
			{
				Config: testAccAliCloudOosOpsItemConfigurationConfigUpdate(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_oos_ops_item_configuration.default", "ops_item_configuration_name", fmt.Sprintf("tf-testacc%d-upd", r)),
					resource.TestCheckResourceAttrSet("alicloud_oos_ops_item_configuration.default", "resource_group_id"),
				),
			},
			{
				ResourceName:      "alicloud_oos_ops_item_configuration.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccAliCloudOosOpsItemConfigurationConfig(r int) string {
	return fmt.Sprintf(`
resource "alicloud_oos_ops_item_configuration" "default" {
  ops_item_configuration_name = "tf-testacc%[1]d"
  data_source                = jsonencode({ Type = "ALIYUN/CMS" })
}
`, r)
}

func testAccAliCloudOosOpsItemConfigurationConfigUpdate(r int) string {
	return fmt.Sprintf(`
data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_oos_ops_item_configuration" "default" {
  ops_item_configuration_name = "tf-testacc%[1]d-upd"
  data_source                = jsonencode({ Type = "ALIYUN/ARMS" })
  resource_group_id           = data.alicloud_resource_manager_resource_groups.default.groups.0.id
}
`, r)
}
