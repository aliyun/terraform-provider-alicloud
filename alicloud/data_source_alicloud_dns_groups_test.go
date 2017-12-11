package alicloud

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
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
					resource.TestCheckResourceAttr("data.alicloud_dns_groups.group", "groups.0.group_id", "520fa32a-076b-4f80-854d-987046e223fe"),
					resource.TestCheckResourceAttr("data.alicloud_dns_groups.group", "groups.0.group_name", "yuy"),
				),
			},
		},
	})
}

const testAccCheckAlicloudDnsGroupsDataSourceNameRegexConfig = `
data "alicloud_dns_groups" "group" {
  name_regex = "^yu"
}`
