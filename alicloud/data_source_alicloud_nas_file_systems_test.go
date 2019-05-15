package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudNasFileSystem_DataSourceSourceStorageType(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudFileSystemsDataSourceStorageType,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_nas_file_systems.fs"),
					resource.TestCheckResourceAttrSet("data.alicloud_nas_file_systems.fs", "systems.0.id"),
					resource.TestCheckResourceAttrSet("data.alicloud_nas_file_systems.fs", "systems.0.region_id"),
					resource.TestCheckResourceAttr("data.alicloud_nas_file_systems.fs", "systems.0.description", "tf-testAccCheckAlicloudFileSystemsDataSource"),
					resource.TestCheckResourceAttrSet("data.alicloud_nas_file_systems.fs", "systems.0.protocol_type"),
					resource.TestCheckResourceAttr("data.alicloud_nas_file_systems.fs", "systems.0.storage_type", "Performance"),
					resource.TestCheckResourceAttrSet("data.alicloud_nas_file_systems.fs", "systems.0.metered_size"),
					resource.TestCheckResourceAttrSet("data.alicloud_nas_file_systems.fs", "systems.0.create_time"),
					resource.TestCheckResourceAttr("data.alicloud_nas_file_systems.fs", "ids.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_nas_file_systems.fs", "ids.0"),
				),
			},
			{
				Config: testAccCheckAlicloudFileSystemsDataSourceStorageTypeEmpty,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_nas_file_systems.fs"),
					resource.TestCheckResourceAttr("data.alicloud_nas_file_systems.fs", "systems.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_nas_file_systems.fs", "ids.#", "0"),
				),
			},
		},
	})
}

func TestAccAlicloudNasFileSystem_DataSourceSourceProtocolType(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudFileSystemsDataSourceProtocolType,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_nas_file_systems.fs"),
					resource.TestCheckResourceAttrSet("data.alicloud_nas_file_systems.fs", "systems.0.id"),
					resource.TestCheckResourceAttrSet("data.alicloud_nas_file_systems.fs", "systems.0.region_id"),
					resource.TestCheckResourceAttr("data.alicloud_nas_file_systems.fs", "systems.0.description", "tf-testAccCheckAlicloudFileSystemsDataSource"),
					resource.TestCheckResourceAttrSet("data.alicloud_nas_file_systems.fs", "systems.0.protocol_type"),
					resource.TestCheckResourceAttr("data.alicloud_nas_file_systems.fs", "systems.0.storage_type", "Performance"),
					resource.TestCheckResourceAttrSet("data.alicloud_nas_file_systems.fs", "systems.0.metered_size"),
					resource.TestCheckResourceAttrSet("data.alicloud_nas_file_systems.fs", "systems.0.create_time"),
					resource.TestCheckResourceAttr("data.alicloud_nas_file_systems.fs", "ids.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_nas_file_systems.fs", "ids.0"),
				),
			},
			{
				Config: testAccCheckAlicloudFileSystemsDataSourceProtocolTypeEmpty,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_nas_file_systems.fs", "ids.#", "0"),
					testAccCheckAlicloudDataSourceID("data.alicloud_nas_file_systems.fs"),
					resource.TestCheckResourceAttr("data.alicloud_nas_file_systems.fs", "systems.#", "0"),
				),
			},
		},
	})
}

func TestAccAlicloudNasFileSystem_DataSourceSourceDescription(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudFileSystemsDataSourceDescription,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_nas_file_systems.fs"),
					resource.TestCheckResourceAttrSet("data.alicloud_nas_file_systems.fs", "systems.0.id"),
					resource.TestCheckResourceAttrSet("data.alicloud_nas_file_systems.fs", "systems.0.region_id"),
					resource.TestCheckResourceAttr("data.alicloud_nas_file_systems.fs", "systems.0.description", "tf-testAccCheckAlicloudFileSystemsDataSource"),
					resource.TestCheckResourceAttrSet("data.alicloud_nas_file_systems.fs", "systems.0.protocol_type"),
					resource.TestCheckResourceAttr("data.alicloud_nas_file_systems.fs", "systems.0.storage_type", "Performance"),
					resource.TestCheckResourceAttrSet("data.alicloud_nas_file_systems.fs", "systems.0.metered_size"),
					resource.TestCheckResourceAttrSet("data.alicloud_nas_file_systems.fs", "systems.0.create_time"),
					resource.TestCheckResourceAttr("data.alicloud_nas_file_systems.fs", "ids.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_nas_file_systems.fs", "ids.0"),
				),
			},
			{
				Config: testAccCheckAlicloudFileSystemsDataSourceDescriptionEmpty,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_nas_file_systems.fs", "ids.#", "0"),
					testAccCheckAlicloudDataSourceID("data.alicloud_nas_file_systems.fs"),
					resource.TestCheckResourceAttr("data.alicloud_nas_file_systems.fs", "systems.#", "0"),
				),
			},
		},
	})
}

func TestAccAlicloudNasFileSystem_DataSourceSourceIds(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudFileSystemsDataSourceIds,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_nas_file_systems.fs"),
					resource.TestCheckResourceAttrSet("data.alicloud_nas_file_systems.fs", "systems.0.id"),
					resource.TestCheckResourceAttrSet("data.alicloud_nas_file_systems.fs", "systems.0.region_id"),
					resource.TestCheckResourceAttr("data.alicloud_nas_file_systems.fs", "systems.0.description", "tf-testAccCheckAlicloudFileSystemsDataSource"),
					resource.TestCheckResourceAttrSet("data.alicloud_nas_file_systems.fs", "systems.0.protocol_type"),
					resource.TestCheckResourceAttr("data.alicloud_nas_file_systems.fs", "systems.0.storage_type", "Capacity"),
					resource.TestCheckResourceAttrSet("data.alicloud_nas_file_systems.fs", "systems.0.metered_size"),
					resource.TestCheckResourceAttrSet("data.alicloud_nas_file_systems.fs", "systems.0.create_time"),
					resource.TestCheckResourceAttr("data.alicloud_nas_file_systems.fs", "ids.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_nas_file_systems.fs", "ids.0"),
				),
			},
			{
				Config: testAccCheckAlicloudFileSystemsDataSourceIdsEmpty,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_nas_file_systems.fs"),
					resource.TestCheckResourceAttr("data.alicloud_nas_file_systems.fs", "ids.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_nas_file_systems.fs", "systems.#", "0"),
				),
			},
		},
	})
}

func TestAccAlicloudNasFileSystem_DataSourceSourceAll(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudFileSystemsDataSourceAll,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_nas_file_systems.fs"),
					resource.TestCheckResourceAttrSet("data.alicloud_nas_file_systems.fs", "systems.0.id"),
					resource.TestCheckResourceAttrSet("data.alicloud_nas_file_systems.fs", "systems.0.region_id"),
					resource.TestCheckResourceAttr("data.alicloud_nas_file_systems.fs", "systems.0.description", "tf-testAccCheckAlicloudFileSystemsDataSource"),
					resource.TestCheckResourceAttrSet("data.alicloud_nas_file_systems.fs", "systems.0.protocol_type"),
					resource.TestCheckResourceAttr("data.alicloud_nas_file_systems.fs", "systems.0.storage_type", "Capacity"),
					resource.TestCheckResourceAttrSet("data.alicloud_nas_file_systems.fs", "systems.0.metered_size"),
					resource.TestCheckResourceAttrSet("data.alicloud_nas_file_systems.fs", "systems.0.create_time"),
					resource.TestCheckResourceAttr("data.alicloud_nas_file_systems.fs", "ids.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_nas_file_systems.fs", "ids.0"),
				),
			},
			{
				Config: testAccCheckAlicloudFileSystemsDataSourceAllEmpty,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_nas_file_systems.fs"),
					resource.TestCheckResourceAttr("data.alicloud_nas_file_systems.fs", "systems.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_nas_file_systems.fs", "ids.#", "0"),
				),
			},
		},
	})
}

const testAccCheckAlicloudFileSystemsDataSourceStorageType = `
variable "description" {
  default = "tf-testAccCheckAlicloudFileSystemsDataSource"
}
data "alicloud_nas_protocols" "default" {
        type = "Performance"
}
resource "alicloud_nas_file_system" "foo" {
  description = "${var.description}"
  storage_type = "Performance"
  protocol_type = "${data.alicloud_nas_protocols.default.protocols.0}"
}
data "alicloud_nas_file_systems" "fs" {
  storage_type = "Performance"
  description_regex = "^${alicloud_nas_file_system.foo.description}"
}
`
const testAccCheckAlicloudFileSystemsDataSourceStorageTypeEmpty = `
variable "description" {
  default = "tf-testAccCheckAlicloudFileSystemsDataSource"
}
data "alicloud_nas_protocols" "default" {
        type = "Performance"
}
resource "alicloud_nas_file_system" "foo" {
  description = "${var.description}"
  storage_type = "Performance"
  protocol_type = "${data.alicloud_nas_protocols.default.protocols.0}"
}
data "alicloud_nas_file_systems" "fs" {
  storage_type = "Capacity"
  description_regex = "^${alicloud_nas_file_system.foo.description}"
}
`
const testAccCheckAlicloudFileSystemsDataSourceProtocolType = `
variable "description" {
  default = "tf-testAccCheckAlicloudFileSystemsDataSource"
}
data "alicloud_nas_protocols" "default" {
        type = "Performance"
}
resource "alicloud_nas_file_system" "foo" {
  description = "${var.description}"
  storage_type = "Performance"
  protocol_type = "${data.alicloud_nas_protocols.default.protocols.0}"
}
data "alicloud_nas_file_systems" "fs" {
  protocol_type = "${data.alicloud_nas_protocols.default.protocols.0}"
  description_regex = "^${alicloud_nas_file_system.foo.description}"
}
`
const testAccCheckAlicloudFileSystemsDataSourceProtocolTypeEmpty = `
variable "description" {
  default = "tf-testAccCheckAlicloudFileSystemsDataSource"
}
data "alicloud_nas_protocols" "default" {
        type = "Performance"
}
resource "alicloud_nas_file_system" "foo" {
  description = "${var.description}"
  storage_type = "Performance"
  protocol_type = "${data.alicloud_nas_protocols.default.protocols.0}"
}
data "alicloud_nas_file_systems" "fs" {
  protocol_type = "${data.alicloud_nas_protocols.default.protocols.2}"
  description_regex = "^${alicloud_nas_file_system.foo.description}"
}
`
const testAccCheckAlicloudFileSystemsDataSourceDescription = `
variable "description" {
  default = "tf-testAccCheckAlicloudFileSystemsDataSource"
}
data "alicloud_nas_protocols" "default" {
        type = "Performance"
}
resource "alicloud_nas_file_system" "foo" {
  description = "${var.description}"
  storage_type = "Performance"
  protocol_type = "${data.alicloud_nas_protocols.default.protocols.0}"
}
data "alicloud_nas_file_systems" "fs" {
  description_regex = "^${alicloud_nas_file_system.foo.description}"
}
`
const testAccCheckAlicloudFileSystemsDataSourceDescriptionEmpty = `
variable "description" {
  default = "tf-testAccCheckAlicloudFileSystemsDataSource"
}
data "alicloud_nas_protocols" "default" {
        type = "Performance"
}
resource "alicloud_nas_file_system" "foo" {
  description = "${var.description}"
  storage_type = "Performance"
  protocol_type = "${data.alicloud_nas_protocols.default.protocols.0}"
}
data "alicloud_nas_file_systems" "fs" {
  description_regex = "^${alicloud_nas_file_system.foo.description}-fake"
}
`
const testAccCheckAlicloudFileSystemsDataSourceIds = `
variable "description" {
  default = "tf-testAccCheckAlicloudFileSystemsDataSource"
}
data "alicloud_nas_protocols" "default" {
        type = "Capacity"
}
resource "alicloud_nas_file_system" "foo" {
  description = "${var.description}"
  storage_type = "Capacity"
  protocol_type = "${data.alicloud_nas_protocols.default.protocols.0}"
}
data "alicloud_nas_file_systems" "fs" {
  ids = ["${alicloud_nas_file_system.foo.id}"]
}
`
const testAccCheckAlicloudFileSystemsDataSourceIdsEmpty = `
variable "description" {
  default = "tf-testAccCheckAlicloudFileSystemsDataSource"
}
data "alicloud_nas_protocols" "default" {
        type = "Capacity"
}
resource "alicloud_nas_file_system" "foo" {
  description = "${var.description}"
  storage_type = "Capacity"
  protocol_type = "${data.alicloud_nas_protocols.default.protocols.0}"
}
data "alicloud_nas_file_systems" "fs" {
  ids = ["${alicloud_nas_file_system.foo.id}-fake"]
}
`

const testAccCheckAlicloudFileSystemsDataSourceAll = `
variable "description" {
  default = "tf-testAccCheckAlicloudFileSystemsDataSource"
}
data "alicloud_nas_protocols" "default" {
        type = "Capacity"
}
resource "alicloud_nas_file_system" "foo" {
  description = "${var.description}"
  storage_type = "Capacity"
  protocol_type = "${data.alicloud_nas_protocols.default.protocols.0}"
}
data "alicloud_nas_file_systems" "fs" {
  protocol_type = "${data.alicloud_nas_protocols.default.protocols.0}"
  storage_type = "Capacity"
  description_regex = "^${alicloud_nas_file_system.foo.description}"
  ids = ["${alicloud_nas_file_system.foo.id}"]
}
`
const testAccCheckAlicloudFileSystemsDataSourceAllEmpty = `
variable "description" {
  default = "tf-testAccCheckAlicloudFileSystemsDataSource"
}
data "alicloud_nas_protocols" "bar_1" {
        type = "Capacity"
}
resource "alicloud_nas_file_system" "foo" {
  description = "${var.description}"
  storage_type = "Capacity"
  protocol_type = "${data.alicloud_nas_protocols.bar_1.protocols.0}"
}
data "alicloud_nas_file_systems" "fs" {
  storage_type = "Performance"
  protocol_type = "${data.alicloud_nas_protocols.bar_1.protocols.0}"
  description_regex = "tf-testAccCheckAlicloudFile"
  ids = ["${alicloud_nas_file_system.foo.id}"]
}
`
