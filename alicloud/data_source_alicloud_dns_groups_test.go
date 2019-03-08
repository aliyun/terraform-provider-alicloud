package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudDnsGroupsDataSource_nameregexAll(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDnsGroupsDataSourceNameRegexAll,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_dns_groups.group"),
					resource.TestCheckResourceAttr("data.alicloud_dns_groups.group", "groups.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_dns_groups.group", "groups.0.group_id", ""),
					resource.TestCheckResourceAttr("data.alicloud_dns_groups.group", "groups.0.group_name", "ALL"),
				),
			},
		},
	})
}

func TestAccAlicloudDnsGroupsDataSource_name_regex(t *testing.T) {
	rand := acctest.RandIntRange(1000, 9999)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDnsGroupsDataSourceNameRegex_match(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_dns_groups.group"),
					resource.TestCheckResourceAttr("data.alicloud_dns_groups.group", "groups.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_dns_groups.group", "groups.0.group_id"),
					resource.TestCheckResourceAttr("data.alicloud_dns_groups.group", "groups.0.group_name", fmt.Sprintf("tf-testacc-%d", rand)),
				),
			},
			{
				Config: testAccCheckAlicloudDnsGroupsDataSourceNameRegex_Nonatch(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_dns_groups.group"),
					resource.TestCheckResourceAttr("data.alicloud_dns_groups.group", "groups.#", "0"),
				),
			},
		},
	})
}

func TestAccAlicloudDnsGroupsDataSource_empty(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDnsGroupsDataSourceNameRegexEmpty,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_dns_groups.group"),
					resource.TestCheckResourceAttr("data.alicloud_dns_groups.group", "groups.#", "0"),
					resource.TestCheckNoResourceAttr("data.alicloud_dns_groups.group", "groups.0.group_id"),
					resource.TestCheckNoResourceAttr("data.alicloud_dns_groups.group", "groups.0.group_name"),
				),
			},
		},
	})
}

const testAccCheckAlicloudDnsGroupsDataSourceNameRegexAll = `
data "alicloud_dns_groups" "group" {
  name_regex = "^ALL"
}`

func testAccCheckAlicloudDnsGroupsDataSourceNameRegex_match(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_dns_group" "foo" {
	  name = "tf-testacc%d-"
	}
	resource "alicloud_dns_group" "group" {
	  name = "tf-testacc-%d"
	}
	data "alicloud_dns_groups" "group" {
	  name_regex = "${alicloud_dns_group.group.name}"
	}`, rand, rand)
}

func testAccCheckAlicloudDnsGroupsDataSourceNameRegex_Nonatch(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_dns_group" "foo" {
	  name = "tf-testacc%d-"
	}
	resource "alicloud_dns_group" "group" {
	  name = "tf-testacc-%d"
	}
	data "alicloud_dns_groups" "group" {
	  name_regex = "anyother"
	}`, rand, rand)
}

const testAccCheckAlicloudDnsGroupsDataSourceNameRegexEmpty = `
data "alicloud_dns_groups" "group" {
  name_regex = "^tf-testacc-fake-name"
}`
