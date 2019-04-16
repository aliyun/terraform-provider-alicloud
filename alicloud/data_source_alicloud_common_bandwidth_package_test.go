package alicloud

import (
	"testing"

	"regexp"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudCommonBandwidthPackagesDataSourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudCommonBandwidthPackagesDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_common_bandwidth_packages.foo"),
					resource.TestCheckResourceAttr("data.alicloud_common_bandwidth_packages.foo", "packages.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_common_bandwidth_packages.foo", "packages.0.id"),
					resource.TestCheckResourceAttrSet("data.alicloud_common_bandwidth_packages.foo", "packages.0.isp"),
					resource.TestCheckResourceAttrSet("data.alicloud_common_bandwidth_packages.foo", "packages.0.creation_time"),
					resource.TestCheckResourceAttrSet("data.alicloud_common_bandwidth_packages.foo", "packages.0.status"),
					resource.TestCheckResourceAttrSet("data.alicloud_common_bandwidth_packages.foo", "packages.0.business_status"),
					resource.TestMatchResourceAttr("data.alicloud_common_bandwidth_packages.foo", "packages.0.name", regexp.MustCompile("^tf-testAcc-for-cbwp-datasourc")),
					resource.TestMatchResourceAttr("data.alicloud_common_bandwidth_packages.foo", "packages.0.description", regexp.MustCompile("^tf-testAcc-for-cbwp-datasourc")),
					resource.TestCheckResourceAttr("data.alicloud_common_bandwidth_packages.foo", "packages.0.bandwidth", "2"),
					resource.TestCheckResourceAttr("data.alicloud_common_bandwidth_packages.foo", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_common_bandwidth_packages.foo", "names.#", "1"),
				),
			},
			{
				Config: testAccCheckAlicloudCommonBandwidthPackagesDataSourceConfig_mismatch,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_common_bandwidth_packages.foo"),
					resource.TestCheckResourceAttr("data.alicloud_common_bandwidth_packages.foo", "packages.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_common_bandwidth_packages.foo", "ids.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_common_bandwidth_packages.foo", "names.#", "0"),
				),
			},
		},
	})
}

func TestAccAlicloudCommonBandwidthPackagesDataSourceNameRegex(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudCommonBandwidthPackagesDataSourceNameRegex,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_common_bandwidth_packages.foo"),
					resource.TestCheckResourceAttr("data.alicloud_common_bandwidth_packages.foo", "packages.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_common_bandwidth_packages.foo", "packages.0.id"),
					resource.TestCheckResourceAttrSet("data.alicloud_common_bandwidth_packages.foo", "packages.0.isp"),
					resource.TestCheckResourceAttrSet("data.alicloud_common_bandwidth_packages.foo", "packages.0.creation_time"),
					resource.TestCheckResourceAttrSet("data.alicloud_common_bandwidth_packages.foo", "packages.0.status"),
					resource.TestCheckResourceAttrSet("data.alicloud_common_bandwidth_packages.foo", "packages.0.business_status"),
					resource.TestMatchResourceAttr("data.alicloud_common_bandwidth_packages.foo", "packages.0.name", regexp.MustCompile("^tf-testAcc-for-cbwp-datasourc")),
					resource.TestMatchResourceAttr("data.alicloud_common_bandwidth_packages.foo", "packages.0.description", regexp.MustCompile("^tf-testAcc-for-cbwp-datasourc")),
					resource.TestCheckResourceAttr("data.alicloud_common_bandwidth_packages.foo", "packages.0.bandwidth", "2"),
					resource.TestCheckResourceAttr("data.alicloud_common_bandwidth_packages.foo", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_common_bandwidth_packages.foo", "names.#", "1"),
				),
			},
			{
				Config: testAccCheckAlicloudCommonBandwidthPackagesDataSourceNameRegex_mismatch,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_common_bandwidth_packages.foo"),
					resource.TestCheckResourceAttr("data.alicloud_common_bandwidth_packages.foo", "packages.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_common_bandwidth_packages.foo", "ids.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_common_bandwidth_packages.foo", "names.#", "0"),
				),
			},
		},
	})
}

func TestAccAlicloudCommonBandwidthPackagesDataSourceIds(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudCommonBandwidthPackagesDataSourceIds,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_common_bandwidth_packages.foo"),
					resource.TestCheckResourceAttr("data.alicloud_common_bandwidth_packages.foo", "packages.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_common_bandwidth_packages.foo", "packages.0.id"),
					resource.TestCheckResourceAttrSet("data.alicloud_common_bandwidth_packages.foo", "packages.0.isp"),
					resource.TestCheckResourceAttrSet("data.alicloud_common_bandwidth_packages.foo", "packages.0.creation_time"),
					resource.TestCheckResourceAttrSet("data.alicloud_common_bandwidth_packages.foo", "packages.0.status"),
					resource.TestCheckResourceAttrSet("data.alicloud_common_bandwidth_packages.foo", "packages.0.business_status"),
					resource.TestMatchResourceAttr("data.alicloud_common_bandwidth_packages.foo", "packages.0.name", regexp.MustCompile("^tf-testAcc-for-cbwp-datasourc")),
					resource.TestMatchResourceAttr("data.alicloud_common_bandwidth_packages.foo", "packages.0.description", regexp.MustCompile("^tf-testAcc-for-cbwp-datasourc")),
					resource.TestCheckResourceAttr("data.alicloud_common_bandwidth_packages.foo", "packages.0.bandwidth", "2"),
					resource.TestCheckResourceAttr("data.alicloud_common_bandwidth_packages.foo", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_common_bandwidth_packages.foo", "names.#", "1"),
				),
			},
			{
				Config: testAccCheckAlicloudCommonBandwidthPackagesDataSourceIds_mismatch,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_common_bandwidth_packages.foo"),
					resource.TestCheckResourceAttr("data.alicloud_common_bandwidth_packages.foo", "packages.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_common_bandwidth_packages.foo", "ids.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_common_bandwidth_packages.foo", "names.#", "0"),
				),
			},
		},
	})
}

const testAccCheckAlicloudCommonBandwidthPackagesDataSourceConfig = `
variable "name" {
  default = "tf-testAcc-for-cbwp-datasource"
}

resource "alicloud_common_bandwidth_package" "foo" {
  bandwidth = "2"
  name = "${var.name}"
  description = "${var.name}"
}

data "alicloud_common_bandwidth_packages" "foo"  {
  name_regex = "${alicloud_common_bandwidth_package.foo.name}"
  ids = ["${alicloud_common_bandwidth_package.foo.id}"]
}
`
const testAccCheckAlicloudCommonBandwidthPackagesDataSourceConfig_mismatch = `
variable "name" {
  default = "tf-testAcc-for-cbwp-datasource"
}

resource "alicloud_common_bandwidth_package" "foo" {
  bandwidth = "2"
  name = "${var.name}"
  description = "${var.name}"
}

data "alicloud_common_bandwidth_packages" "foo"  {
  name_regex = "${alicloud_common_bandwidth_package.foo.name}-fake"
  ids = ["${alicloud_common_bandwidth_package.foo.id}-fake"]
}
`

const testAccCheckAlicloudCommonBandwidthPackagesDataSourceNameRegex = `
variable "name" {
  default = "tf-testAcc-for-cbwp-datasource"
}

resource "alicloud_common_bandwidth_package" "foo" {
  bandwidth = "2"
  name = "${var.name}"
  description = "${var.name}"
}

data "alicloud_common_bandwidth_packages" "foo"  {
  name_regex = "${alicloud_common_bandwidth_package.foo.name}"
}
`

const testAccCheckAlicloudCommonBandwidthPackagesDataSourceNameRegex_mismatch = `
variable "name" {
  default = "tf-testAcc-for-cbwp-datasource"
}

resource "alicloud_common_bandwidth_package" "foo" {
  bandwidth = "2"
  name = "${var.name}"
  description = "${var.name}"
}

data "alicloud_common_bandwidth_packages" "foo"  {
  name_regex = "${alicloud_common_bandwidth_package.foo.name}-fake"
}
`

const testAccCheckAlicloudCommonBandwidthPackagesDataSourceIds = `
variable "name" {
  default = "tf-testAcc-for-cbwp-datasource"
}

resource "alicloud_common_bandwidth_package" "foo" {
  bandwidth = "2"
  name = "${var.name}"
  description = "${var.name}"
}

data "alicloud_common_bandwidth_packages" "foo"  {
  ids = ["${alicloud_common_bandwidth_package.foo.id}"]
}
`

const testAccCheckAlicloudCommonBandwidthPackagesDataSourceIds_mismatch = `
variable "name" {
  default = "tf-testAcc-for-cbwp-datasource"
}

resource "alicloud_common_bandwidth_package" "foo" {
  bandwidth = "2"
  name = "${var.name}"
  description = "${var.name}"
}

data "alicloud_common_bandwidth_packages" "foo"  {
  ids = ["${alicloud_common_bandwidth_package.foo.id}-fake"]
}
`
