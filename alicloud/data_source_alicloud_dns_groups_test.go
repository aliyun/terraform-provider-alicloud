package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudDnsGroupsDataSource_name_regex(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDnsGroupsDataSourceNameRegexConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_dns_groups.group"),
					resource.TestCheckResourceAttr("data.alicloud_dns_groups.group", "groups.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_dns_groups.group", "groups.0.group_name", "ALL"),
				),
			},
		},
	})
}

const testAccCheckAlicloudDnsGroupsDataSourceNameRegexConfig = `
data "alicloud_dns_groups" "group" {
  name_regex = "^ALL"
}`
