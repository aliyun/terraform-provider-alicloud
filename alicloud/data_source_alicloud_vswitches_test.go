package alicloud

import (
	"fmt"
	"testing"

	"regexp"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudVSwitchesDataSource_all(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudVSwitchesDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_vswitches.foo"),
					resource.TestCheckResourceAttr("data.alicloud_vswitches.foo", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_vswitches.foo", "names.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_vswitches.foo", "names.0"),
					resource.TestCheckResourceAttr("data.alicloud_vswitches.foo", "vswitches.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_vswitches.foo", "vswitches.0.id"),
					resource.TestCheckResourceAttrSet("data.alicloud_vswitches.foo", "vswitches.0.vpc_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_vswitches.foo", "vswitches.0.zone_id"),
					resource.TestMatchResourceAttr("data.alicloud_vswitches.foo", "vswitches.0.name", regexp.MustCompile("^tf-testAcc-for-vswitch-datasourc")),
					resource.TestCheckResourceAttr("data.alicloud_vswitches.foo", "vswitches.0.instance_ids.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_vswitches.foo", "vswitches.0.cidr_block", "172.16.0.0/16"),
					resource.TestCheckResourceAttr("data.alicloud_vswitches.foo", "vswitches.0.description", ""),
					resource.TestCheckResourceAttr("data.alicloud_vswitches.foo", "vswitches.0.is_default", "false"),
					resource.TestCheckResourceAttrSet("data.alicloud_vswitches.foo", "vswitches.0.creation_time"),
				),
			},
		},
	})
}

func TestAccAlicloudVSwitchesDataSource_Name_Regex(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudVSwitchesDataSourceConfigNameRegex(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_vswitches.foo"),
					resource.TestCheckResourceAttr("data.alicloud_vswitches.foo", "vswitches.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_vswitches.foo", "vswitches.0.id"),
					resource.TestCheckResourceAttrSet("data.alicloud_vswitches.foo", "vswitches.0.vpc_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_vswitches.foo", "vswitches.0.zone_id"),
					resource.TestMatchResourceAttr("data.alicloud_vswitches.foo", "vswitches.0.name", regexp.MustCompile(fmt.Sprintf("tf-testAcc-for-vswitch-datasource-name-regex-%d", rand))),
					resource.TestCheckResourceAttr("data.alicloud_vswitches.foo", "vswitches.0.cidr_block", "172.16.1.0/24"),
					resource.TestCheckResourceAttr("data.alicloud_vswitches.foo", "vswitches.0.instance_ids.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_vswitches.foo", "vswitches.0.description", ""),
					resource.TestCheckResourceAttr("data.alicloud_vswitches.foo", "vswitches.0.is_default", "false"),
					resource.TestCheckResourceAttrSet("data.alicloud_vswitches.foo", "vswitches.0.creation_time"),
				),
			},
			{
				Config: testAccCheckAlicloudVSwitchesDataSourceConfigNameRegexEmpty(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_vswitches.foo"),
					resource.TestCheckResourceAttr("data.alicloud_vswitches.foo", "ids.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_vswitches.foo", "names.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_vswitches.foo", "vswitches.#", "0"),
				),
			},
		},
	})
}

func TestAccAlicloudVSwitchesDataSource_vpcid(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudVSwitchesDataSourceConfigVPCID,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_vswitches.foo"),
					resource.TestCheckResourceAttr("data.alicloud_vswitches.foo", "vswitches.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_vswitches.foo", "vswitches.0.id"),
					resource.TestCheckResourceAttrSet("data.alicloud_vswitches.foo", "vswitches.0.vpc_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_vswitches.foo", "vswitches.0.zone_id"),
					resource.TestMatchResourceAttr("data.alicloud_vswitches.foo", "vswitches.0.name", regexp.MustCompile("^tf-testAcc-for-vswitch-datasource-VPC-ID")),
					resource.TestCheckResourceAttr("data.alicloud_vswitches.foo", "vswitches.0.cidr_block", "172.16.0.0/24"),
					resource.TestCheckResourceAttr("data.alicloud_vswitches.foo", "vswitches.0.instance_ids.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_vswitches.foo", "vswitches.0.description", ""),
					resource.TestCheckResourceAttr("data.alicloud_vswitches.foo", "vswitches.0.is_default", "false"),
					resource.TestCheckResourceAttrSet("data.alicloud_vswitches.foo", "vswitches.0.creation_time"),
				),
			},
			{
				Config: testAccCheckAlicloudVSwitchesDataSourceConfigVPCIDEmpty,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_vswitches.foo"),
					resource.TestCheckResourceAttr("data.alicloud_vswitches.foo", "ids.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_vswitches.foo", "names.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_vswitches.foo", "vswitches.#", "0"),
				),
			},
		},
	})
}
func TestAccAlicloudVSwitchesDataSource_Zone_ID(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudVSwitchesDataSourceConfigAvailabilityZoneID,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_vswitches.foo"),
					resource.TestCheckResourceAttr("data.alicloud_vswitches.foo", "vswitches.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_vswitches.foo", "vswitches.0.id"),
					resource.TestCheckResourceAttrSet("data.alicloud_vswitches.foo", "vswitches.0.vpc_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_vswitches.foo", "vswitches.0.zone_id"),
					resource.TestMatchResourceAttr("data.alicloud_vswitches.foo", "vswitches.0.name", regexp.MustCompile("^tf-testAcc-for-vswitch-datasource-availability_zone_ID")),
					resource.TestCheckResourceAttr("data.alicloud_vswitches.foo", "vswitches.0.instance_ids.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_vswitches.foo", "vswitches.0.cidr_block", "172.16.0.0/24"),
					resource.TestCheckResourceAttr("data.alicloud_vswitches.foo", "vswitches.0.description", ""),
					resource.TestCheckResourceAttr("data.alicloud_vswitches.foo", "vswitches.0.is_default", "false"),
					resource.TestCheckResourceAttrSet("data.alicloud_vswitches.foo", "vswitches.0.creation_time"),
				),
			},
			{
				Config: testAccCheckAlicloudVSwitchesDataSourceConfigAvailabilityZoneIDEmpty,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_vswitches.foo"),
					resource.TestCheckResourceAttr("data.alicloud_vswitches.foo", "ids.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_vswitches.foo", "names.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_vswitches.foo", "vswitches.#", "0"),
				),
			},
		},
	})
}

const testAccCheckAlicloudVSwitchesDataSourceConfig = `
variable "name" {
  default = "tf-testAcc-for-vswitch-datasource"
}
data "alicloud_zones" "default" {}

resource "alicloud_vpc" "vpc" {
  cidr_block = "172.16.0.0/16"
  name = "${var.name}"
}
resource "alicloud_vswitch" "vswitch" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/16"
  vpc_id = "${alicloud_vpc.vpc.id}"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}
data "alicloud_vswitches" "foo" {
  name_regex = "^tf-testAcc-.*-datasource"
  vpc_id = "${alicloud_vpc.vpc.id}"
  cidr_block = "${alicloud_vswitch.vswitch.cidr_block}"
  zone_id = "${data.alicloud_zones.default.zones.0.id}"
}
`

func testAccCheckAlicloudVSwitchesDataSourceConfigNameRegex(rand int) string {
	return fmt.Sprintf(
		`
variable "name" {
  default = "tf-testAcc-for-vswitch-datasource-name-regex-%d"
}
data "alicloud_zones" "default" {}

resource "alicloud_vpc" "vpc" {
  cidr_block = "172.16.0.0/16"
  name = "${var.name}"
}
resource "alicloud_vswitch" "vswitch" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/24"
  vpc_id = "${alicloud_vpc.vpc.id}"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}
resource "alicloud_vswitch" "vswitch1" {
  name = "${var.name}-A"
  cidr_block = "172.16.1.0/24"
  vpc_id = "${alicloud_vpc.vpc.id}"
  availability_zone = "${data.alicloud_zones.default.zones.1.id}"
}
data "alicloud_vswitches" "foo" {
  name_regex = "${alicloud_vswitch.vswitch1.name}"
}
`, rand)
}

func testAccCheckAlicloudVSwitchesDataSourceConfigNameRegexEmpty(rand int) string {
	return fmt.Sprintf(
		`
variable "name" {
  default = "tf-testAcc-for-vswitch-datasource-name-regex-%d"
}
data "alicloud_zones" "default" {}

resource "alicloud_vpc" "vpc" {
  cidr_block = "172.16.0.0/16"
  name = "${var.name}"
}
resource "alicloud_vswitch" "vswitch" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/24"
  vpc_id = "${alicloud_vpc.vpc.id}"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}
resource "alicloud_vswitch" "vswitch1" {
  name = "${var.name}-A"
  cidr_block = "172.16.1.0/24"
  vpc_id = "${alicloud_vpc.vpc.id}"
  availability_zone = "${data.alicloud_zones.default.zones.1.id}"
}
data "alicloud_vswitches" "foo" {
  name_regex = "${alicloud_vswitch.vswitch1.name}-fake"
}
`, rand)
}

const testAccCheckAlicloudVSwitchesDataSourceConfigAvailabilityZoneID = `
variable "name" {
  default = "tf-testAcc-for-vswitch-datasource-availability_zone_ID"
}
data "alicloud_zones" "default" {}

resource "alicloud_vpc" "vpc" {
  cidr_block = "172.16.0.0/16"
  name = "${var.name}"
}

resource "alicloud_vswitch" "vswitch" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/24"
  vpc_id = "${alicloud_vpc.vpc.id}"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}
resource "alicloud_vswitch" "vswitch1" {
  name = "${var.name}-A"
  cidr_block = "172.16.1.0/24"
  vpc_id = "${alicloud_vpc.vpc.id}"
  availability_zone = "${data.alicloud_zones.default.zones.1.id}"
}
data "alicloud_vswitches" "foo" {
	vpc_id = "${alicloud_vpc.vpc.id}"
  zone_id = "${alicloud_vswitch.vswitch.availability_zone}"
}
`

const testAccCheckAlicloudVSwitchesDataSourceConfigAvailabilityZoneIDEmpty = `
variable "name" {
  default = "tf-testAcc-for-vswitch-datasource-availability_zone_ID"
}
data "alicloud_zones" "default" {}

resource "alicloud_vpc" "vpc" {
  cidr_block = "172.16.0.0/16"
  name = "${var.name}"
}

resource "alicloud_vswitch" "vswitch" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/24"
  vpc_id = "${alicloud_vpc.vpc.id}"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}
resource "alicloud_vswitch" "vswitch1" {
  name = "${var.name}-A"
  cidr_block = "172.16.1.0/24"
  vpc_id = "${alicloud_vpc.vpc.id}"
  availability_zone = "${data.alicloud_zones.default.zones.1.id}"
}
data "alicloud_vswitches" "foo" {
	vpc_id = "${alicloud_vpc.vpc.id}"
  zone_id = "${alicloud_vswitch.vswitch.availability_zone}-fake"
}
`
const testAccCheckAlicloudVSwitchesDataSourceConfigVPCID = `
variable "name" {
  default = "tf-testAcc-for-vswitch-datasource-VPC-ID"
}
data "alicloud_zones" "default" {}

resource "alicloud_vpc" "vpc" {
  cidr_block = "172.16.0.0/16"
  name = "${var.name}"
}
resource "alicloud_vpc" "vpc1" {
  cidr_block = "192.168.0.0/16"
  name = "${var.name}-A"
}
resource "alicloud_vswitch" "vswitch" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/24"
  vpc_id = "${alicloud_vpc.vpc.id}"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}

resource "alicloud_vswitch" "vswitch1" {
  name = "${var.name}"
  cidr_block = "192.168.0.0/16"
  vpc_id = "${alicloud_vpc.vpc1.id}"
  availability_zone = "${data.alicloud_zones.default.zones.1.id}"
}
data "alicloud_vswitches" "foo" {
	vpc_id = "${alicloud_vswitch.vswitch.vpc_id}"
}
`

const testAccCheckAlicloudVSwitchesDataSourceConfigVPCIDEmpty = `
variable "name" {
  default = "tf-testAcc-for-vswitch-datasource-VPC-ID"
}
data "alicloud_zones" "default" {}

resource "alicloud_vpc" "vpc" {
  cidr_block = "172.16.0.0/16"
  name = "${var.name}"
}
resource "alicloud_vpc" "vpc1" {
  cidr_block = "192.168.0.0/16"
  name = "${var.name}-A"
}
resource "alicloud_vswitch" "vswitch" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/24"
  vpc_id = "${alicloud_vpc.vpc.id}"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}

resource "alicloud_vswitch" "vswitch1" {
  name = "${var.name}"
  cidr_block = "192.168.0.0/16"
  vpc_id = "${alicloud_vpc.vpc1.id}"
  availability_zone = "${data.alicloud_zones.default.zones.1.id}"
}
data "alicloud_vswitches" "foo" {
	vpc_id = "${alicloud_vswitch.vswitch.vpc_id}-fake"
}
`
