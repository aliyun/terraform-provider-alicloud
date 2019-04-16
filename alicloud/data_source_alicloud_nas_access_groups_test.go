package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudNasAccessGroupDataSource_all(t *testing.T) {
	rand := acctest.RandIntRange(10000, 999999)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, connectivity.NasClassicSupportedRegions)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudAccessGroupsDataSourceAll(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_nas_access_groups.ag"),
					resource.TestCheckResourceAttr("data.alicloud_nas_access_groups.ag", "groups.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_nas_access_groups.ag", "groups.0.rule_count"),
					resource.TestCheckResourceAttr("data.alicloud_nas_access_groups.ag", "groups.0.type", "Classic"),
					resource.TestCheckResourceAttr("data.alicloud_nas_access_groups.ag", "groups.0.description", "tf-testAccAccessGroupsdatasource"),
					resource.TestCheckResourceAttr("data.alicloud_nas_access_groups.ag", "groups.0.id", fmt.Sprintf("tf-testAccAccessGroupsdatasource-%d", rand)),
					resource.TestCheckResourceAttrSet("data.alicloud_nas_access_groups.ag", "groups.0.mount_target_count"),
					resource.TestCheckResourceAttr("data.alicloud_nas_access_groups.ag", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_nas_access_groups.ag", "ids.0", fmt.Sprintf("tf-testAccAccessGroupsdatasource-%d", rand)),
				),
			},
			{
				Config: testAccCheckAlicloudAccessGroupsDataSourceAllEmpty(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_nas_access_groups.ag"),
					resource.TestCheckResourceAttr("data.alicloud_nas_access_groups.ag", "groups.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_nas_access_groups.ag", "ids.#", "0"),
				),
			},
		},
	})
}

func TestAccAlicloudNasAccessGroupDataSource_type(t *testing.T) {
	rand := acctest.RandIntRange(10000, 999999)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudAccessGroupsDataSourceType(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_nas_access_groups.ag"),
					resource.TestCheckResourceAttr("data.alicloud_nas_access_groups.ag", "groups.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_nas_access_groups.ag", "groups.0.rule_count"),
					resource.TestCheckResourceAttr("data.alicloud_nas_access_groups.ag", "groups.0.type", "Vpc"),
					resource.TestCheckResourceAttr("data.alicloud_nas_access_groups.ag", "groups.0.description", "tf-testAccAccessGroupsdatasource"),
					resource.TestCheckResourceAttr("data.alicloud_nas_access_groups.ag", "groups.0.id", fmt.Sprintf("tf-testAccAccessGroupsdatasource-%d", rand)),
					resource.TestCheckResourceAttrSet("data.alicloud_nas_access_groups.ag", "groups.0.mount_target_count"),
					resource.TestCheckResourceAttr("data.alicloud_nas_access_groups.ag", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_nas_access_groups.ag", "ids.0", fmt.Sprintf("tf-testAccAccessGroupsdatasource-%d", rand)),
				),
			},
			{
				Config: testAccCheckAlicloudAccessGroupsDataSourceTypeEmpty(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_nas_access_groups.ag"),
					resource.TestCheckResourceAttr("data.alicloud_nas_access_groups.ag", "groups.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_nas_access_groups.ag", "ids.#", "0"),
				),
			},
		},
	})
}

func TestAccAlicloudNasAccessGroupDataSource_description(t *testing.T) {
	rand := acctest.RandIntRange(10000, 999999)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudAccessGroupsDataSourceDescription(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_nas_access_groups.ag"),
					resource.TestCheckResourceAttr("data.alicloud_nas_access_groups.ag", "groups.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_nas_access_groups.ag", "groups.0.rule_count"),
					resource.TestCheckResourceAttr("data.alicloud_nas_access_groups.ag", "groups.0.type", "Vpc"),
					resource.TestCheckResourceAttr("data.alicloud_nas_access_groups.ag", "groups.0.description", "tf-testAccAccessGroupsdatasource"),
					resource.TestCheckResourceAttr("data.alicloud_nas_access_groups.ag", "groups.0.id", fmt.Sprintf("tf-testAccAccessGroupsdatasource-%d", rand)),
					resource.TestCheckResourceAttrSet("data.alicloud_nas_access_groups.ag", "groups.0.mount_target_count"),
					resource.TestCheckResourceAttr("data.alicloud_nas_access_groups.ag", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_nas_access_groups.ag", "ids.0", fmt.Sprintf("tf-testAccAccessGroupsdatasource-%d", rand)),
				),
			},
			{
				Config: testAccCheckAlicloudAccessGroupsDataSourceDescriptionEmpty(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_nas_access_groups.ag"),
					resource.TestCheckResourceAttr("data.alicloud_nas_access_groups.ag", "groups.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_nas_access_groups.ag", "ids.#", "0"),
				),
			},
		},
	})
}

func testAccCheckAlicloudAccessGroupsDataSourceAll(rand int) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "tf-testAccAccessGroupsdatasource-%d"
	}
	resource "alicloud_nas_access_group" "foo" {
		name = "${var.name}"
		type = "Classic"
		description = "tf-testAccAccessGroupsdatasource"
	}
	data "alicloud_nas_access_groups" "ag" {
		name_regex = "testAccAccessGroupsdatasource*"
		type = "${alicloud_nas_access_group.foo.type}"
		description = "tf-testAccAccessGroupsdatasource"
	}`, rand)
}

func testAccCheckAlicloudAccessGroupsDataSourceType(rand int) string {
	return fmt.Sprintf(`
	variable "name" {
        	default = "tf-testAccAccessGroupsdatasource-%d"
	}
	resource "alicloud_nas_access_group" "foo" {
        	name = "${var.name}"
	        type = "Vpc"
	        description = "tf-testAccAccessGroupsdatasource"
	}
	data "alicloud_nas_access_groups" "ag" {
        	name_regex = "testAccAccessGroupsdatasource*"
	        type = "${alicloud_nas_access_group.foo.type}"
	}`, rand)
}

func testAccCheckAlicloudAccessGroupsDataSourceDescription(rand int) string {
	return fmt.Sprintf(`
	variable "name" {
        	default = "tf-testAccAccessGroupsdatasource-%d"
	}
	resource "alicloud_nas_access_group" "foo" {
        	name = "${var.name}"
	        type = "Vpc"
	        description = "tf-testAccAccessGroupsdatasource"
	}
	data "alicloud_nas_access_groups" "ag" {
        	name_regex = "testAccAccessGroupsdatasource*"
	        description = "${alicloud_nas_access_group.foo.description}"
	}`, rand)
}

func testAccCheckAlicloudAccessGroupsDataSourceAllEmpty(rand int) string {
	return fmt.Sprintf(`
	variable "name" {
                 default = "tf-testAccAccessGroupsdatasource-%d"
        }
        resource "alicloud_nas_access_group" "foo" {
                name = "${var.name}"
                type = "Classic"
                description = "tf-testAccAccessGroupsdatasource"
        }
        data "alicloud_nas_access_groups" "ag" {
                name_regex = "testAccAccessGroupsdatasource*"
                type = "Vpc"
                description = "tf-testAccAccessGroups"
        }`, rand)
}

func testAccCheckAlicloudAccessGroupsDataSourceTypeEmpty(rand int) string {
	return fmt.Sprintf(`
        variable "name" {
                default = "tf-testAccAccessGroupsdatasource-%d"
        }
        resource "alicloud_nas_access_group" "foo" {
                name = "${var.name}"
                type = "Vpc"
                description = "tf-testAccAccessGroupsdatasource"
        }
        data "alicloud_nas_access_groups" "ag" {
                name_regex = "testAccAccessGroupsdatasource*"
                type = "Classis"
        }`, rand)
}

func testAccCheckAlicloudAccessGroupsDataSourceDescriptionEmpty(rand int) string {
	return fmt.Sprintf(`
        variable "name" {
                default = "tf-testAccAccessGroupsdatasource-%d"
        }
        resource "alicloud_nas_access_group" "foo" {
                name = "${var.name}"
                type = "Vpc"
                description = "tf-testAccAccessGroupsdatasource"
        }
        data "alicloud_nas_access_groups" "ag" {
                name_regex = "testAccAccessGroupsdatasource*"
                description = "tf-testAccAccessGroups"
        }`, rand)
}
