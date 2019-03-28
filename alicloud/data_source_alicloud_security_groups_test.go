package alicloud

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudSecurityGroupsDataSource_all(t *testing.T) {
	randnum := rand.Intn(1000)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckAlicloudSecurityGroupsDataSourceConfigWithName, randnum),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_security_groups.web"),
					resource.TestCheckResourceAttr("data.alicloud_security_groups.web", "groups.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_security_groups.web", "groups.0.vpc_id"),
					resource.TestCheckResourceAttr("data.alicloud_security_groups.web", "groups.0.name", fmt.Sprintf("tf-testAccCheckAlicloudSecurityGroupsDataSourceConfig%d", randnum)),
					resource.TestCheckResourceAttr("data.alicloud_security_groups.web", "groups.0.inner_access", "true"),
					resource.TestCheckResourceAttrSet("data.alicloud_security_groups.web", "groups.0.creation_time"),
					resource.TestCheckResourceAttr("data.alicloud_security_groups.web", "groups.0.description", "test security group"),
					resource.TestCheckResourceAttrSet("data.alicloud_security_groups.web", "groups.0.id"),
				),
			},
		},
	})
}

func TestAccAlicloudSecurityGroupsDataSource_vpc(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudSecurityGroupsDataSourceConfigWithVpc,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_security_groups.web"),
					resource.TestCheckResourceAttr("data.alicloud_security_groups.web", "groups.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_security_groups.web", "groups.0.vpc_id"),
					resource.TestCheckResourceAttr("data.alicloud_security_groups.web", "groups.0.name", "tf-testAccCheckAlicloudSecurityGroupsDataSourceConfig"),
					resource.TestCheckResourceAttr("data.alicloud_security_groups.web", "groups.0.description", "test security group"),
					resource.TestCheckResourceAttr("data.alicloud_security_groups.web", "groups.0.inner_access", "true"),
					resource.TestCheckResourceAttrSet("data.alicloud_security_groups.web", "groups.0.creation_time"),
					resource.TestCheckNoResourceAttr("data.alicloud_security_groups.web", "groups.0.tags"),
					resource.TestCheckResourceAttrSet("data.alicloud_security_groups.web", "groups.0.id"),
				),
			},
			{
				Config: testAccCheckAlicloudSecurityGroupsDataSourceConfigWithVpcEmpty,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_security_groups.web"),
					resource.TestCheckResourceAttr("data.alicloud_security_groups.web", "ids.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_security_groups.web", "names.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_security_groups.web", "groups.#", "0"),
				),
			},
		},
	})
}

func TestAccAlicloudSecurityGroupsDataSource_tags(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudSecurityGroupDataSourceTags,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_security_groups.web"),
					resource.TestCheckResourceAttr("data.alicloud_security_groups.web", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_security_groups.web", "names.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_security_groups.web", "groups.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_security_groups.web", "groups.0.vpc_id"),
					resource.TestCheckResourceAttr("data.alicloud_security_groups.web", "groups.0.name", "tf-testAccCheckAlicloudSecurityGroupsDataSourceConfig"),
					resource.TestCheckResourceAttr("data.alicloud_security_groups.web", "groups.0.tags.from", "datasource"),
					resource.TestCheckResourceAttr("data.alicloud_security_groups.web", "groups.0.tags.usage1", "test"),
					resource.TestCheckResourceAttr("data.alicloud_security_groups.web", "groups.0.tags.usage2", "test"),
					resource.TestCheckResourceAttr("data.alicloud_security_groups.web", "groups.0.inner_access", "true"),
					resource.TestCheckResourceAttrSet("data.alicloud_security_groups.web", "groups.0.creation_time"),
					resource.TestCheckResourceAttr("data.alicloud_security_groups.web", "groups.0.description", "test security group"),
					resource.TestCheckResourceAttrSet("data.alicloud_security_groups.web", "groups.0.id"),
				),
			},
			{
				Config: testAccCheckAlicloudSecurityGroupDataSourceTagsEmpty,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_security_groups.web"),
					resource.TestCheckResourceAttr("data.alicloud_security_groups.web", "ids.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_security_groups.web", "names.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_security_groups.web", "groups.#", "0"),
				),
			},
		},
	})
}

var testAccCheckAlicloudSecurityGroupsDataSourceConfigWithName = `
variable "name" {
	default = "tf-testAccCheckAlicloudSecurityGroupsDataSourceConfig%d"
}
resource "alicloud_vpc" "foo" {
  cidr_block = "172.16.0.0/12"
  name = "${var.name}"
}

resource "alicloud_security_group" "test" {
  name        = "${var.name}"
  description = "test security group"
  vpc_id      = "${alicloud_vpc.foo.id}"
}

resource "alicloud_security_group" "unuse" {
  name        = "tf-testAccCheckAlicloudSecurityGroupsDataSourceConfig"
  description = "test data source"
  vpc_id      = "${alicloud_vpc.foo.id}"
}

data "alicloud_security_groups" "web" {
    name_regex = "${alicloud_security_group.test.name}"
}
`

const testAccCheckAlicloudSecurityGroupsDataSourceConfigWithVpc = `
variable "name" {
	default = "tf-testAccCheckAlicloudSecurityGroupsDataSourceConfig"
}

resource "alicloud_vpc" "foo" {
  cidr_block = "172.16.0.0/12"
  name = "${var.name}"
}

resource "alicloud_vpc" "test" {
  cidr_block = "192.168.0.0/16"
  name = "tf-testAccCheckAlicloudSecurityGroupsDataSourceConfig-unuse"
}

resource "alicloud_security_group" "test" {
  name        = "${var.name}"
  description = "test security group"
  vpc_id      = "${alicloud_vpc.foo.id}"
}

resource "alicloud_security_group" "unuse" {
  name        = "tf-testAccCheckAlicloudSecurityGroupsDataSourceConfig_unuse"
  description = "test data source"
  vpc_id      = "${alicloud_vpc.test.id}"
}

data "alicloud_security_groups" "web" {
  vpc_id    = "${alicloud_security_group.test.vpc_id}"
}
`

const testAccCheckAlicloudSecurityGroupsDataSourceConfigWithVpcEmpty = `
variable "name" {
	default = "tf-testAccCheckAlicloudSecurityGroupsDataSourceConfig"
}

resource "alicloud_vpc" "foo" {
  cidr_block = "172.16.0.0/12"
  name = "${var.name}"
}

resource "alicloud_vpc" "test" {
  cidr_block = "192.168.0.0/16"
  name = "tf-testAccCheckAlicloudSecurityGroupsDataSourceConfig-unuse"
}

resource "alicloud_security_group" "test" {
  name        = "${var.name}"
  description = "test security group"
  vpc_id      = "${alicloud_vpc.foo.id}"
}

resource "alicloud_security_group" "unuse" {
  name        = "tf-testAccCheckAlicloudSecurityGroupsDataSourceConfig_unuse"
  description = "test data source"
  vpc_id      = "${alicloud_vpc.test.id}"
}

data "alicloud_security_groups" "web" {
  vpc_id    = "${alicloud_security_group.test.vpc_id}-fake"
}
`

const testAccCheckAlicloudSecurityGroupDataSourceTags = `
variable "name" {
	default = "tf-testAccCheckAlicloudSecurityGroupsDataSourceConfig"
}
resource "alicloud_vpc" "tf_vpc_foo" {
  cidr_block = "172.16.0.0/12"
  name = "${var.name}"
}

resource "alicloud_security_group" "test" {
  name        = "${var.name}"
  description = "test security group"
  vpc_id      = "${alicloud_vpc.tf_vpc_foo.id}"
  tags {
		from = "datasource"
		usage1 = "test"
		usage2 = "test"
  }
}

data "alicloud_security_groups" "web" {
	tags = "${alicloud_security_group.test.tags}"
}
`

const testAccCheckAlicloudSecurityGroupDataSourceTagsEmpty = `
variable "name" {
	default = "tf-testAccCheckAlicloudSecurityGroupsDataSourceConfig"
}
resource "alicloud_vpc" "tf_vpc_foo" {
  cidr_block = "172.16.0.0/12"
  name = "${var.name}"
}

resource "alicloud_security_group" "test" {
  name        = "${var.name}"
  description = "test security group"
  vpc_id      = "${alicloud_vpc.tf_vpc_foo.id}"
  tags {
		from = "datasource"
		usage1 = "test"
		usage2 = "test"
  }
}

data "alicloud_security_groups" "web" {
	tags = {
		from = "${alicloud_security_group.test.tags.from}-fake"
	}
}
`
