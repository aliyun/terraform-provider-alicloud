package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudNasAccessRuleDataSource_ip(t *testing.T) {
	rand := acctest.RandIntRange(10000, 999999)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, connectivity.NasClassicSupportedRegions)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudAccessRulesDataSourceIp(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_nas_access_rules.rule"),
					resource.TestCheckResourceAttr("data.alicloud_nas_access_rules.rule", "rules.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_nas_access_rules.rule", "rules.0.source_cidr_ip", "168.1.1.0/16"),
					resource.TestCheckResourceAttr("data.alicloud_nas_access_rules.rule", "rules.0.priority", "2"),
					resource.TestCheckResourceAttrSet("data.alicloud_nas_access_rules.rule", "rules.0.access_rule_id"),
					resource.TestCheckResourceAttr("data.alicloud_nas_access_rules.rule", "rules.0.user_access", "no_squash"),
					resource.TestCheckResourceAttr("data.alicloud_nas_access_rules.rule", "rules.0.rw_access", "RDWR"),
					resource.TestCheckResourceAttr("data.alicloud_nas_access_rules.rule", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_nas_access_rules.rule", "ids.0", fmt.Sprintf("tf-testAccAccessGroupsdatasource-%d:1", rand)),
				),
			},
			{
				Config: testAccCheckAlicloudAccessRulesDataSourceIpEmpty(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_nas_access_rules.rule"),
					resource.TestCheckResourceAttr("data.alicloud_nas_access_rules.rule", "rules.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_nas_access_rules.rule", "ids.#", "0"),
				),
			},
		},
	})
}

func TestAccAlicloudNasAccessRuleDataSource_RWAccess(t *testing.T) {
	rand := acctest.RandIntRange(10000, 999999)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudAccessRulesDataSourceRWAccess(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_nas_access_rules.rule"),
					resource.TestCheckResourceAttr("data.alicloud_nas_access_rules.rule", "rules.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_nas_access_rules.rule", "rules.0.source_cidr_ip", "168.1.1.0/16"),
					resource.TestCheckResourceAttr("data.alicloud_nas_access_rules.rule", "rules.0.priority", "2"),
					resource.TestCheckResourceAttrSet("data.alicloud_nas_access_rules.rule", "rules.0.access_rule_id"),
					resource.TestCheckResourceAttr("data.alicloud_nas_access_rules.rule", "rules.0.user_access", "no_squash"),
					resource.TestCheckResourceAttr("data.alicloud_nas_access_rules.rule", "rules.0.rw_access", "RDWR"),
					resource.TestCheckResourceAttr("data.alicloud_nas_access_rules.rule", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_nas_access_rules.rule", "ids.0", fmt.Sprintf("tf-testAccAccessGroupsdatasource-%d:1", rand)),
				),
			},
			{
				Config: testAccCheckAlicloudAccessRulesDataSourceRWAccessEmpty(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_nas_access_rules.rule"),
					resource.TestCheckResourceAttr("data.alicloud_nas_access_rules.rule", "rules.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_nas_access_rules.rule", "ids.#", "0"),
				),
			},
		},
	})
}

func TestAccAlicloudNasAccessRuleDataSource_UserAccess(t *testing.T) {
	rand := acctest.RandIntRange(10000, 999999)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudAccessRulesDataSourceUserAccess(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_nas_access_rules.rule"),
					resource.TestCheckResourceAttr("data.alicloud_nas_access_rules.rule", "rules.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_nas_access_rules.rule", "rules.0.source_cidr_ip", "168.1.1.0/16"),
					resource.TestCheckResourceAttr("data.alicloud_nas_access_rules.rule", "rules.0.priority", "2"),
					resource.TestCheckResourceAttrSet("data.alicloud_nas_access_rules.rule", "rules.0.access_rule_id"),
					resource.TestCheckResourceAttr("data.alicloud_nas_access_rules.rule", "rules.0.user_access", "no_squash"),
					resource.TestCheckResourceAttr("data.alicloud_nas_access_rules.rule", "rules.0.rw_access", "RDWR"),
					resource.TestCheckResourceAttr("data.alicloud_nas_access_rules.rule", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_nas_access_rules.rule", "ids.0", fmt.Sprintf("tf-testAccAccessGroupsdatasource-%d:1", rand)),
				),
			},
			{
				Config: testAccCheckAlicloudAccessRulesDataSourceUserAccessEmpty(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_nas_access_rules.rule"),
					resource.TestCheckResourceAttr("data.alicloud_nas_access_rules.rule", "rules.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_nas_access_rules.rule", "ids.#", "0"),
				),
			},
		},
	})
}

func TestAccAlicloudNasAccessRuleDataSource_All(t *testing.T) {
	rand := acctest.RandIntRange(10000, 999999)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudAccessRulesDataSourceAll(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_nas_access_rules.rule"),
					resource.TestCheckResourceAttr("data.alicloud_nas_access_rules.rule", "rules.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_nas_access_rules.rule", "rules.0.source_cidr_ip", "168.1.1.0/16"),
					resource.TestCheckResourceAttr("data.alicloud_nas_access_rules.rule", "rules.0.priority", "2"),
					resource.TestCheckResourceAttrSet("data.alicloud_nas_access_rules.rule", "rules.0.access_rule_id"),
					resource.TestCheckResourceAttr("data.alicloud_nas_access_rules.rule", "rules.0.user_access", "no_squash"),
					resource.TestCheckResourceAttr("data.alicloud_nas_access_rules.rule", "rules.0.rw_access", "RDWR"),
					resource.TestCheckResourceAttr("data.alicloud_nas_access_rules.rule", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_nas_access_rules.rule", "ids.0", fmt.Sprintf("tf-testAccAccessGroupsdatasource-%d:1", rand)),
				),
			},
			{
				Config: testAccCheckAlicloudAccessRulesDataSourceAllEmpty(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_nas_access_rules.rule"),
					resource.TestCheckResourceAttr("data.alicloud_nas_access_rules.rule", "rules.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_nas_access_rules.rule", "ids.#", "0"),
				),
			},
		},
	})
}

func testAccCheckAlicloudAccessRulesDataSourceIp(rand int) string {
	return fmt.Sprintf(`
	variable "name" {
        	default = "tf-testAccAccessGroupsdatasource-%d"
	}
	resource "alicloud_nas_access_group" "foo" {
        	name = "${var.name}"
	        type = "Classic"
	        description = "tf-testAccAccessGroupsdatasource"
	}
	resource "alicloud_nas_access_rule" "foo" {
        	access_group_name = "${alicloud_nas_access_group.foo.id}"
	        source_cidr_ip = "168.1.1.0/16"
        	rw_access_type = "RDWR"
	        user_access_type = "no_squash"
	        priority = 2
	}
	data "alicloud_nas_access_rules" "rule" {
		access_group_name = "${alicloud_nas_access_group.foo.id}"
		source_cidr_ip = "${alicloud_nas_access_rule.foo.source_cidr_ip}"
	}`, rand)
}

func testAccCheckAlicloudAccessRulesDataSourceRWAccess(rand int) string {
	return fmt.Sprintf(`
        variable "name" {
                default = "tf-testAccAccessGroupsdatasource-%d"
        }
        resource "alicloud_nas_access_group" "foo" {
                name = "${var.name}"
                type = "Vpc"
                description = "tf-testAccAccessGroupsdatasource"
        }
        resource "alicloud_nas_access_rule" "foo" {
                access_group_name = "${alicloud_nas_access_group.foo.id}"
                source_cidr_ip = "168.1.1.0/16"
                rw_access_type = "RDWR"
                user_access_type = "no_squash"
                priority = 2
        }
        data "alicloud_nas_access_rules" "rule" {
                access_group_name = "${alicloud_nas_access_group.foo.id}"
		rw_access = "${alicloud_nas_access_rule.foo.rw_access_type}"
        }`, rand)
}

func testAccCheckAlicloudAccessRulesDataSourceUserAccess(rand int) string {
	return fmt.Sprintf(`
        variable "name" {
                default = "tf-testAccAccessGroupsdatasource-%d"
        }
        resource "alicloud_nas_access_group" "foo" {
                name = "${var.name}"
                type = "Vpc"
                description = "tf-testAccAccessGroupsdatasource"
        }
        resource "alicloud_nas_access_rule" "foo" {
                access_group_name = "${alicloud_nas_access_group.foo.id}"
                source_cidr_ip = "168.1.1.0/16"
                rw_access_type = "RDWR"
                user_access_type = "no_squash"
                priority = 2
        }
        data "alicloud_nas_access_rules" "rule" {
                access_group_name = "${alicloud_nas_access_group.foo.id}"
                user_access = "${alicloud_nas_access_rule.foo.user_access_type}"
        }`, rand)
}

func testAccCheckAlicloudAccessRulesDataSourceAll(rand int) string {
	return fmt.Sprintf(`
        variable "name" {
                default = "tf-testAccAccessGroupsdatasource-%d"
        }
        resource "alicloud_nas_access_group" "foo" {
                name = "${var.name}"
                type = "Vpc"
                description = "tf-testAccAccessGroupsdatasource"
        }
        resource "alicloud_nas_access_rule" "foo" {
                access_group_name = "${alicloud_nas_access_group.foo.id}"
                source_cidr_ip = "168.1.1.0/16"
                rw_access_type = "RDWR"
                user_access_type = "no_squash"
                priority = 2
        }
        data "alicloud_nas_access_rules" "rule" {
                access_group_name = "${alicloud_nas_access_group.foo.id}"
                user_access = "${alicloud_nas_access_rule.foo.user_access_type}"
		rw_access = "${alicloud_nas_access_rule.foo.rw_access_type}"
		source_cidr_ip = "${alicloud_nas_access_rule.foo.source_cidr_ip}"
        }`, rand)
}

func testAccCheckAlicloudAccessRulesDataSourceAllEmpty(rand int) string {
	return fmt.Sprintf(`
        variable "name" {
                default = "tf-testAccAccessGroupsdatasource-%d"
        }
        resource "alicloud_nas_access_group" "foo" {
                name = "${var.name}"
                type = "Vpc"
                description = "tf-testAccAccessGroupsdatasource"
        }
        resource "alicloud_nas_access_rule" "foo" {
                access_group_name = "${alicloud_nas_access_group.foo.id}"
                source_cidr_ip = "168.1.1.0/16"
                rw_access_type = "RDWR"
                user_access_type = "no_squash"
                priority = 2
        }
        data "alicloud_nas_access_rules" "rule" {
                access_group_name = "${alicloud_nas_access_group.foo.id}"
                source_cidr_ip = "192.1.1.0/16"
		user_access = "root_squash"
		rw_access = "RDONLY"
        }`, rand)
}

func testAccCheckAlicloudAccessRulesDataSourceIpEmpty(rand int) string {
	return fmt.Sprintf(`
        variable "name" {
                default = "tf-testAccAccessGroupsdatasource-%d"
        }
        resource "alicloud_nas_access_group" "foo" {
                name = "${var.name}"
                type = "Classic"
                description = "tf-testAccAccessGroupsdatasource"
        }
        resource "alicloud_nas_access_rule" "foo" {
                access_group_name = "${alicloud_nas_access_group.foo.id}"
                source_cidr_ip = "168.1.1.0/16"
                rw_access_type = "RDWR"
                user_access_type = "no_squash"
                priority = 2
        }
        data "alicloud_nas_access_rules" "rule" {
                access_group_name = "${alicloud_nas_access_group.foo.id}"
                source_cidr_ip = "192.1.1.0/16"
        }`, rand)
}

func testAccCheckAlicloudAccessRulesDataSourceRWAccessEmpty(rand int) string {
	return fmt.Sprintf(`
        variable "name" {
                default = "tf-testAccAccessGroupsdatasource-%d"
        }
        resource "alicloud_nas_access_group" "foo" {
                name = "${var.name}"
                type = "Vpc"
                description = "tf-testAccAccessGroupsdatasource"
        }
        resource "alicloud_nas_access_rule" "foo" {
                access_group_name = "${alicloud_nas_access_group.foo.id}"
                source_cidr_ip = "168.1.1.0/16"
                rw_access_type = "RDWR"
                user_access_type = "no_squash"
                priority = 2
        }
        data "alicloud_nas_access_rules" "rule" {
                access_group_name = "${alicloud_nas_access_group.foo.id}"
                rw_access = "RDONLY"
        }`, rand)
}

func testAccCheckAlicloudAccessRulesDataSourceUserAccessEmpty(rand int) string {
	return fmt.Sprintf(`
        variable "name" {
                default = "tf-testAccAccessGroupsdatasource-%d"
        }
        resource "alicloud_nas_access_group" "foo" {
                name = "${var.name}"
                type = "Vpc"
                description = "tf-testAccAccessGroupsdatasource"
        }
        resource "alicloud_nas_access_rule" "foo" {
                access_group_name = "${alicloud_nas_access_group.foo.id}"
                source_cidr_ip = "168.1.1.0/16"
                rw_access_type = "RDWR"
                user_access_type = "no_squash"
                priority = 2
        }
        data "alicloud_nas_access_rules" "rule" {
                access_group_name = "${alicloud_nas_access_group.foo.id}"
                user_access = "root_squash"
        }`, rand)
}
