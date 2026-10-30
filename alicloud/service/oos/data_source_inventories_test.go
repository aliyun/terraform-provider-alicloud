package oos_test

import (
	"fmt"
	"os"
	"testing"

	acctest "github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudOosInventoriesDataSource_basic0(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacc%d", rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			os.Setenv("ALICLOUD_REGION", "cn-hangzhou")
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			// max_results filter, matching
			{
				Config: testAccInventoriesConfig(name, fmt.Sprintf(`max_results = 5`)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_oos_inventories.default", "inventories.#", "5"),
					resource.TestCheckResourceAttrSet("data.alicloud_oos_inventories.default", "next_token_result"),
				),
			},
			// aggregator filter, matching
			{
				Config: testAccInventoriesConfig(name, fmt.Sprintf(`aggregator = ["ACS:InstanceInformation.PlatformName"]`)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_oos_inventories.default", "inventories.#", "1"),
				),
			},
			// filter filter, matchingmak
			{
				Config: testAccInventoriesConfig(name, fmt.Sprintf(`
				filter {
				  name     = "ACS:InstanceInformation.Status"
				  operator = "Equal"
				  value    = ["Stopped"]
				}
				max_results = 1`)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_oos_inventories.default", "inventories.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_oos_inventories.default", "next_token_result"),
				),
			},
			// both filters, matching
			{
				Config: testAccInventoriesConfig(name, fmt.Sprintf(`aggregator = ["ACS:InstanceInformation.PlatformName"]
				filter {
				  name     = "ACS:InstanceInformation.Status"
				  operator = "Equal"
				  value    = ["Stopped"]
				}
				max_results = 1`)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_oos_inventories.default", "inventories.#", "1"),
				),
			},
		},
	})
}

func testAccInventoriesConfig(name string, filter string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_oos_inventories" "default" {
  %s
}
`, name, filter)
}
