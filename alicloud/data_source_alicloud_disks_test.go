package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"

	"fmt"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudDisksDataSource_ids(t *testing.T) {
	rand := acctest.RandIntRange(1000, 9999)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDisksDataSource_ids_exist(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_disks.disks"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_disks.disks", "disks.0.id"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.name", fmt.Sprintf("tf-testAccCheckAlicloudDisksDataSource_ids-%d", rand)),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.description", fmt.Sprintf("tf-testAccCheckAlicloudDisksDataSource_ids-%d_description", rand)),
					resource.TestCheckResourceAttrSet("data.alicloud_disks.disks", "disks.0.region_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_disks.disks", "disks.0.availability_zone"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.status", "Available"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.type", "data"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.category", "cloud_efficiency"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.encrypted", "off"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.size", "20"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.image_id", ""),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.snapshot_id", ""),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.instance_id", ""),
					resource.TestCheckResourceAttrSet("data.alicloud_disks.disks", "disks.0.creation_time"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.attached_time", ""),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.detached_time", ""),
					resource.TestCheckResourceAttrSet("data.alicloud_disks.disks", "disks.0.expiration_time"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.tags.%", "2"),
				),
			},
			{
				Config: testAccCheckAlicloudDisksDataSource_ids_fake(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_disks.disks"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.#", "0"),
				),
			},
		},
	})
}

func TestAccAlicloudDisksDataSource_name_regex(t *testing.T) {
	rand := acctest.RandIntRange(1000, 9999)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDisksDataSource_name_regex_exist(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_disks.disks"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_disks.disks", "disks.0.id"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.name", fmt.Sprintf("tf-testAccCheckAlicloudDisksDataSource_ids-%d", rand)),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.description", fmt.Sprintf("tf-testAccCheckAlicloudDisksDataSource_ids-%d_description", rand)),
					resource.TestCheckResourceAttrSet("data.alicloud_disks.disks", "disks.0.region_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_disks.disks", "disks.0.availability_zone"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.status", "Available"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.type", "data"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.category", "cloud_efficiency"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.encrypted", "off"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.size", "20"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.image_id", ""),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.snapshot_id", ""),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.instance_id", ""),
					resource.TestCheckResourceAttrSet("data.alicloud_disks.disks", "disks.0.creation_time"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.attached_time", ""),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.detached_time", ""),
					resource.TestCheckResourceAttrSet("data.alicloud_disks.disks", "disks.0.expiration_time"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.tags.%", "2"),
				),
			},
			{
				Config: testAccCheckAlicloudDisksDataSource_name_regex_fake(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_disks.disks"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.#", "0"),
				),
			},
		},
	})
}

func TestAccAlicloudDisksDataSource_type(t *testing.T) {
	rand := acctest.RandIntRange(1000, 9999)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDisksDataSource_type_exist(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_disks.disks"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_disks.disks", "disks.0.id"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.name", fmt.Sprintf("tf-testAccCheckAlicloudDisksDataSource_ids-%d", rand)),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.description", fmt.Sprintf("tf-testAccCheckAlicloudDisksDataSource_ids-%d_description", rand)),
					resource.TestCheckResourceAttrSet("data.alicloud_disks.disks", "disks.0.region_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_disks.disks", "disks.0.availability_zone"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.status", "Available"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.type", "data"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.category", "cloud_efficiency"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.encrypted", "off"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.size", "20"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.image_id", ""),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.snapshot_id", ""),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.instance_id", ""),
					resource.TestCheckResourceAttrSet("data.alicloud_disks.disks", "disks.0.creation_time"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.attached_time", ""),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.detached_time", ""),
					resource.TestCheckResourceAttrSet("data.alicloud_disks.disks", "disks.0.expiration_time"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.tags.%", "2"),
				),
			},
			{
				Config: testAccCheckAlicloudDisksDataSource_type_fake(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_disks.disks"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.#", "0"),
				),
			},
		},
	})
}

func TestAccAlicloudDisksDataSource_category(t *testing.T) {
	rand := acctest.RandIntRange(1000, 9999)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDisksDataSource_category_exist(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_disks.disks"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_disks.disks", "disks.0.id"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.name", fmt.Sprintf("tf-testAccCheckAlicloudDisksDataSource_ids-%d", rand)),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.description", fmt.Sprintf("tf-testAccCheckAlicloudDisksDataSource_ids-%d_description", rand)),
					resource.TestCheckResourceAttrSet("data.alicloud_disks.disks", "disks.0.region_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_disks.disks", "disks.0.availability_zone"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.status", "Available"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.type", "data"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.category", "cloud_efficiency"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.encrypted", "off"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.size", "20"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.image_id", ""),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.snapshot_id", ""),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.instance_id", ""),
					resource.TestCheckResourceAttrSet("data.alicloud_disks.disks", "disks.0.creation_time"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.attached_time", ""),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.detached_time", ""),
					resource.TestCheckResourceAttrSet("data.alicloud_disks.disks", "disks.0.expiration_time"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.tags.%", "2"),
				),
			},
			{
				Config: testAccCheckAlicloudDisksDataSource_category_fake(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_disks.disks"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.#", "0"),
				),
			},
		},
	})
}

func TestAccAlicloudDisksDataSource_encrypted(t *testing.T) {
	rand := acctest.RandIntRange(1000, 9999)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDisksDataSource_encrypted_exist(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_disks.disks"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_disks.disks", "disks.0.id"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.name", fmt.Sprintf("tf-testAccCheckAlicloudDisksDataSource_ids-%d", rand)),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.description", fmt.Sprintf("tf-testAccCheckAlicloudDisksDataSource_ids-%d_description", rand)),
					resource.TestCheckResourceAttrSet("data.alicloud_disks.disks", "disks.0.region_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_disks.disks", "disks.0.availability_zone"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.status", "Available"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.type", "data"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.category", "cloud_efficiency"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.encrypted", "on"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.size", "20"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.image_id", ""),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.snapshot_id", ""),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.instance_id", ""),
					resource.TestCheckResourceAttrSet("data.alicloud_disks.disks", "disks.0.creation_time"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.attached_time", ""),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.detached_time", ""),
					resource.TestCheckResourceAttrSet("data.alicloud_disks.disks", "disks.0.expiration_time"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.tags.%", "2"),
				),
			},
			{
				Config: testAccCheckAlicloudDisksDataSource_encrypted_fake(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_disks.disks"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.#", "0"),
				),
			},
		},
	})
}

func TestAccAlicloudDisksDataSource_instance_id(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDisksDataSource_instance_id_exist(EcsInstanceCommonTestCase),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_disks.disks"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_disks.disks", "disks.0.id"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.name", "tf-testAccCheckAlicloudDisksDataSourceFilterByInstanceId"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.description", "tf-testAccCheckAlicloudDisksDataSourceFilterByInstanceId_description"),
					resource.TestCheckResourceAttrSet("data.alicloud_disks.disks", "disks.0.region_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_disks.disks", "disks.0.availability_zone"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.status", "In_use"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.type", "data"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.category", "cloud_efficiency"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.encrypted", "off"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.size", "20"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.image_id", ""),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.snapshot_id", ""),
					resource.TestCheckResourceAttrSet("data.alicloud_disks.disks", "disks.0.instance_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_disks.disks", "disks.0.creation_time"),
					resource.TestCheckResourceAttrSet("data.alicloud_disks.disks", "disks.0.attached_time"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.detached_time", ""),
					resource.TestCheckResourceAttrSet("data.alicloud_disks.disks", "disks.0.expiration_time"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.tags.%", "0"),
				),
			},
			{
				Config: testAccCheckAlicloudDisksDataSource_instance_id_fake(EcsInstanceCommonTestCase),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_disks.disks"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.#", "0"),
				),
			},
		},
	})
}

func TestAccAlicloudDisksDataSource_tags(t *testing.T) {
	rand := acctest.RandIntRange(1000, 9999)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDisksDataSource_tags_exist(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_disks.disks"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_disks.disks", "disks.0.id"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.name", fmt.Sprintf("tf-testAccCheckAlicloudDisksDataSource_ids-%d", rand)),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.description", fmt.Sprintf("tf-testAccCheckAlicloudDisksDataSource_ids-%d_description", rand)),
					resource.TestCheckResourceAttrSet("data.alicloud_disks.disks", "disks.0.region_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_disks.disks", "disks.0.availability_zone"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.status", "Available"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.type", "data"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.category", "cloud_efficiency"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.encrypted", "on"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.size", "20"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.image_id", ""),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.snapshot_id", ""),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.instance_id", ""),
					resource.TestCheckResourceAttrSet("data.alicloud_disks.disks", "disks.0.creation_time"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.attached_time", ""),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.detached_time", ""),
					resource.TestCheckResourceAttrSet("data.alicloud_disks.disks", "disks.0.expiration_time"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.tags.%", "2"),
				),
			},
			{
				Config: testAccCheckAlicloudDisksDataSource_tags_fake(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_disks.disks"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.#", "0"),
				),
			},
		},
	})
}

func TestAccAlicloudDisksDataSource_all(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDisksDataSource_all_exist(EcsInstanceCommonTestCase),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_disks.disks"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_disks.disks", "disks.0.id"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.name", "tf-testAccCheckAlicloudDisksDataSourceFilterByInstanceId"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.description", "tf-testAccCheckAlicloudDisksDataSourceFilterByInstanceId_description"),
					resource.TestCheckResourceAttrSet("data.alicloud_disks.disks", "disks.0.region_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_disks.disks", "disks.0.availability_zone"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.status", "In_use"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.type", "data"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.category", "cloud_efficiency"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.encrypted", "off"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.size", "20"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.image_id", ""),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.snapshot_id", ""),
					resource.TestCheckResourceAttrSet("data.alicloud_disks.disks", "disks.0.instance_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_disks.disks", "disks.0.creation_time"),
					resource.TestCheckResourceAttrSet("data.alicloud_disks.disks", "disks.0.attached_time"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.detached_time", ""),
					resource.TestCheckResourceAttrSet("data.alicloud_disks.disks", "disks.0.expiration_time"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.tags.%", "2"),
				),
			},
			{
				Config: testAccCheckAlicloudDisksDataSource_all_fake(EcsInstanceCommonTestCase),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_disks.disks"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.#", "0"),
				),
			},
		},
	})
}

func testAccCheckAlicloudDisksDataSource_ids_exist(rand int) string {
	return fmt.Sprintf(`
variable "name" {
	default = "tf-testAccCheckAlicloudDisksDataSource_ids-%d"
}

data "alicloud_zones" "az" {
	"available_resource_creation"= "VSwitch"
}

resource "alicloud_disk" "sample_disk" {
	availability_zone = "${data.alicloud_zones.az.zones.0.id}"
	category = "cloud_efficiency"
	name = "${var.name}"
    description = "${var.name}_description"
	size = "20"
	tags {
	    Name = "TerraformTest"
	    Name1 = "TerraformTest"
	}
}

data "alicloud_disks" "disks" {
    ids = [ "${alicloud_disk.sample_disk.id}" ]
}
`, rand)
}

func testAccCheckAlicloudDisksDataSource_ids_fake(rand int) string {
	return fmt.Sprintf(`
variable "name" {
	default = "tf-testAccCheckAlicloudDisksDataSource_ids-%d"
}

data "alicloud_zones" "az" {
	"available_resource_creation"= "VSwitch"
}

resource "alicloud_disk" "sample_disk" {
	availability_zone = "${data.alicloud_zones.az.zones.0.id}"
	category = "cloud_efficiency"
	name = "${var.name}"
    description = "${var.name}_description"
	size = "20"
	tags {
	    Name = "TerraformTest"
	    Name1 = "TerraformTest"
	}
}

data "alicloud_disks" "disks" {
    ids = [ "${alicloud_disk.sample_disk.id}_fake" ]
}
`, rand)
}

func testAccCheckAlicloudDisksDataSource_name_regex_exist(rand int) string {
	return fmt.Sprintf(`
variable "name" {
	default = "tf-testAccCheckAlicloudDisksDataSource_ids-%d"
}

data "alicloud_zones" "az" {
	"available_resource_creation"= "VSwitch"
}

resource "alicloud_disk" "sample_disk" {
	availability_zone = "${data.alicloud_zones.az.zones.0.id}"
	category = "cloud_efficiency"
	name = "${var.name}"
    description = "${var.name}_description"
	size = "20"
	tags {
	    Name = "TerraformTest"
	    Name1 = "TerraformTest"
	}
}

data "alicloud_disks" "disks" {
    name_regex = "${alicloud_disk.sample_disk.name}"
}
`, rand)
}

func testAccCheckAlicloudDisksDataSource_name_regex_fake(rand int) string {
	return fmt.Sprintf(`
variable "name" {
	default = "tf-testAccCheckAlicloudDisksDataSource_ids-%d"
}

data "alicloud_zones" "az" {
	"available_resource_creation"= "VSwitch"
}

resource "alicloud_disk" "sample_disk" {
	availability_zone = "${data.alicloud_zones.az.zones.0.id}"
	category = "cloud_efficiency"
	name = "${var.name}"
    description = "${var.name}_description"
	size = "20"
	tags {
	    Name = "TerraformTest"
	    Name1 = "TerraformTest"
	}
}

data "alicloud_disks" "disks" {
     name_regex = "${var.name}_fake"
}
`, rand)
}

func testAccCheckAlicloudDisksDataSource_type_exist(rand int) string {
	return fmt.Sprintf(`
variable "name" {
	default = "tf-testAccCheckAlicloudDisksDataSource_ids-%d"
}

data "alicloud_zones" "az" {
	"available_resource_creation"= "VSwitch"
}

resource "alicloud_disk" "sample_disk" {
	availability_zone = "${data.alicloud_zones.az.zones.0.id}"
	category = "cloud_efficiency"
	name = "${var.name}"
    description = "${var.name}_description"
	size = "20"
	tags {
	    Name = "TerraformTest"
	    Name1 = "TerraformTest"
	}
}

data "alicloud_disks" "disks" {
	 name_regex = "${alicloud_disk.sample_disk.name}"
     type = "data"
}
`, rand)
}

func testAccCheckAlicloudDisksDataSource_type_fake(rand int) string {
	return fmt.Sprintf(`
variable "name" {
	default = "tf-testAccCheckAlicloudDisksDataSource_ids-%d"
}

data "alicloud_zones" "az" {
	"available_resource_creation"= "VSwitch"
}

resource "alicloud_disk" "sample_disk" {
	availability_zone = "${data.alicloud_zones.az.zones.0.id}"
	category = "cloud_efficiency"
	name = "${var.name}"
    description = "${var.name}_description"
	size = "20"
	tags {
	    Name = "TerraformTest"
	    Name1 = "TerraformTest"
	}
}

data "alicloud_disks" "disks" {
     name_regex = "${alicloud_disk.sample_disk.name}"
     type = "system"
}
`, rand)
}

func testAccCheckAlicloudDisksDataSource_category_exist(rand int) string {
	return fmt.Sprintf(`
variable "name" {
	default = "tf-testAccCheckAlicloudDisksDataSource_ids-%d"
}

data "alicloud_zones" "az" {
	"available_resource_creation"= "VSwitch"
}

resource "alicloud_disk" "sample_disk" {
	availability_zone = "${data.alicloud_zones.az.zones.0.id}"
	category = "cloud_efficiency"
	name = "${var.name}"
    description = "${var.name}_description"
	size = "20"
	tags {
	    Name = "TerraformTest"
	    Name1 = "TerraformTest"
	}
}

data "alicloud_disks" "disks" {
	 name_regex = "${alicloud_disk.sample_disk.name}"
     category = "cloud_efficiency"
}
`, rand)
}

func testAccCheckAlicloudDisksDataSource_category_fake(rand int) string {
	return fmt.Sprintf(`
variable "name" {
	default = "tf-testAccCheckAlicloudDisksDataSource_ids-%d"
}

data "alicloud_zones" "az" {
	"available_resource_creation"= "VSwitch"
}

resource "alicloud_disk" "sample_disk" {
	availability_zone = "${data.alicloud_zones.az.zones.0.id}"
	category = "cloud_efficiency"
	name = "${var.name}"
    description = "${var.name}_description"
	size = "20"
	tags {
	    Name = "TerraformTest"
	    Name1 = "TerraformTest"
	}
}

data "alicloud_disks" "disks" {
	 name_regex = "${alicloud_disk.sample_disk.name}"
     category = "cloud"
}
`, rand)
}

func testAccCheckAlicloudDisksDataSource_encrypted_exist(rand int) string {
	return fmt.Sprintf(`
variable "name" {
	default = "tf-testAccCheckAlicloudDisksDataSource_ids-%d"
}

data "alicloud_zones" "az" {
	"available_resource_creation"= "VSwitch"
}

resource "alicloud_disk" "sample_disk" {
	availability_zone = "${data.alicloud_zones.az.zones.0.id}"
	category = "cloud_efficiency"
	name = "${var.name}"
    description = "${var.name}_description"
	size = "20"
	encrypted = "true"
	tags {
	    Name = "TerraformTest"
	    Name1 = "TerraformTest"
	}
}

data "alicloud_disks" "disks" {
	 name_regex = "${alicloud_disk.sample_disk.name}"
     encrypted = "on"
}
`, rand)
}

func testAccCheckAlicloudDisksDataSource_encrypted_fake(rand int) string {
	return fmt.Sprintf(`
variable "name" {
	default = "tf-testAccCheckAlicloudDisksDataSource_ids-%d"
}

data "alicloud_zones" "az" {
	"available_resource_creation"= "VSwitch"
}

resource "alicloud_disk" "sample_disk" {
	availability_zone = "${data.alicloud_zones.az.zones.0.id}"
	category = "cloud_efficiency"
	name = "${var.name}"
    description = "${var.name}_description"
	size = "20"
	encrypted = "true"
	tags {
	    Name = "TerraformTest"
	    Name1 = "TerraformTest"
	}
}

data "alicloud_disks" "disks" {
	 name_regex = "${alicloud_disk.sample_disk.name}"
     encrypted = "off"
}
`, rand)
}

func testAccCheckAlicloudDisksDataSource_instance_id_exist(common string) string {
	return fmt.Sprintf(`	
	%s
	variable "name" {
		default = "tf-testAccCheckAlicloudDisksDataSourceFilterByInstanceId"
	}

	resource "alicloud_disk" "sample_disk" {
		availability_zone = "${data.alicloud_zones.default.zones.0.id}"
		category = "cloud_efficiency"
		name = "${var.name}"
	    description = "${var.name}_description"
		size = "20"
	}

	resource "alicloud_instance" "sample_instance" {
		vswitch_id = "${alicloud_vswitch.default.id}"
		private_ip = "172.16.0.10"
		image_id = "${data.alicloud_images.default.images.0.id}"
		instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
		instance_name = "${var.name}"
		system_disk_category = "cloud_efficiency"
		security_groups = ["${alicloud_security_group.default.id}"]
	}

	resource "alicloud_disk_attachment" "sample_disk_attachment" {
	  disk_id = "${alicloud_disk.sample_disk.id}"
	  instance_id = "${alicloud_instance.sample_instance.id}"
	}

	data "alicloud_disks" "disks" {
	    instance_id = "${alicloud_disk_attachment.sample_disk_attachment.instance_id}"
	    type = "data"
	}
	`, common)
}

func testAccCheckAlicloudDisksDataSource_instance_id_fake(common string) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccCheckAlicloudDisksDataSourceFilterByInstanceId"
	}

	resource "alicloud_disk" "sample_disk" {
		availability_zone = "${data.alicloud_zones.default.zones.0.id}"
		category = "cloud_efficiency"
		name = "${var.name}"
	    description = "${var.name}_description"
		size = "20"
	}

	resource "alicloud_instance" "sample_instance" {
		vswitch_id = "${alicloud_vswitch.default.id}"
		private_ip = "172.16.0.10"
		image_id = "${data.alicloud_images.default.images.0.id}"
		instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
		instance_name = "${var.name}"
		system_disk_category = "cloud_efficiency"
		security_groups = ["${alicloud_security_group.default.id}"]
	}

	resource "alicloud_disk_attachment" "sample_disk_attachment" {
	  disk_id = "${alicloud_disk.sample_disk.id}"
	  instance_id = "${alicloud_instance.sample_instance.id}"
	}

	data "alicloud_disks" "disks" {
	    instance_id = "${alicloud_disk_attachment.sample_disk_attachment.instance_id}_non"
	    type = "data"
	}
	`, common)
}

func testAccCheckAlicloudDisksDataSource_tags_exist(rand int) string {
	return fmt.Sprintf(`
variable "name" {
	default = "tf-testAccCheckAlicloudDisksDataSource_ids-%d"
}

data "alicloud_zones" "az" {
	"available_resource_creation"= "VSwitch"
}

resource "alicloud_disk" "sample_disk" {
	availability_zone = "${data.alicloud_zones.az.zones.0.id}"
	category = "cloud_efficiency"
	name = "${var.name}"
    description = "${var.name}_description"
	size = "20"
	encrypted = "true"
	tags {
	    Name = "TerraformTest"
	    Name1 = "TerraformTest"
	}
}

data "alicloud_disks" "disks" {
	 name_regex = "${alicloud_disk.sample_disk.name}"
     tags = "${alicloud_disk.sample_disk.tags}"
}
`, rand)
}

func testAccCheckAlicloudDisksDataSource_tags_fake(rand int) string {
	return fmt.Sprintf(`
variable "name" {
	default = "tf-testAccCheckAlicloudDisksDataSource_ids-%d"
}

data "alicloud_zones" "az" {
	"available_resource_creation"= "VSwitch"
}

resource "alicloud_disk" "sample_disk" {
	availability_zone = "${data.alicloud_zones.az.zones.0.id}"
	category = "cloud_efficiency"
	name = "${var.name}"
    description = "${var.name}_description"
	size = "20"
	encrypted = "true"
	tags {
	    Name = "TerraformTest"
	    Name1 = "TerraformTest"
	}
}

data "alicloud_disks" "disks" {
	 name_regex = "${alicloud_disk.sample_disk.name}"
     tags = {
	    Name = "TerraformTest_fake"
	    Name1 = "TerraformTest_fake"
	}
}
`, rand)
}

func testAccCheckAlicloudDisksDataSource_all_exist(common string) string {
	return fmt.Sprintf(`	
	%s
	variable "name" {
		default = "tf-testAccCheckAlicloudDisksDataSourceFilterByInstanceId"
	}

	resource "alicloud_disk" "sample_disk" {
		availability_zone = "${data.alicloud_zones.default.zones.0.id}"
		category = "cloud_efficiency"
		name = "${var.name}"
	    description = "${var.name}_description"
		size = "20"
		tags {
	    	Name = "TerraformTest"
			Name1 = "TerraformTest"
		}
	}

	resource "alicloud_instance" "sample_instance" {
		vswitch_id = "${alicloud_vswitch.default.id}"
		private_ip = "172.16.0.10"
		image_id = "${data.alicloud_images.default.images.0.id}"
		instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
		instance_name = "${var.name}"
		system_disk_category = "cloud_efficiency"
		security_groups = ["${alicloud_security_group.default.id}"]
	}

	resource "alicloud_disk_attachment" "sample_disk_attachment" {
	  disk_id = "${alicloud_disk.sample_disk.id}"
	  instance_id = "${alicloud_instance.sample_instance.id}"
	}

	data "alicloud_disks" "disks" {
		ids = [ "${alicloud_disk.sample_disk.id}" ]
		name_regex = "${alicloud_disk.sample_disk.name}"
		type = "data"
		category = "cloud_efficiency"
		encrypted = "off"
		tags = "${alicloud_disk.sample_disk.tags}"
	    instance_id = "${alicloud_disk_attachment.sample_disk_attachment.instance_id}"
	}
	`, common)
}

func testAccCheckAlicloudDisksDataSource_all_fake(common string) string {
	return fmt.Sprintf(`	
	%s
	variable "name" {
		default = "tf-testAccCheckAlicloudDisksDataSourceFilterByInstanceId"
	}

	resource "alicloud_disk" "sample_disk" {
		availability_zone = "${data.alicloud_zones.default.zones.0.id}"
		category = "cloud_efficiency"
		name = "${var.name}"
	    description = "${var.name}_description"
		size = "20"
		tags {
	    	Name = "TerraformTest"
			Name1 = "TerraformTest"
		}
	}

	resource "alicloud_instance" "sample_instance" {
		vswitch_id = "${alicloud_vswitch.default.id}"
		private_ip = "172.16.0.10"
		image_id = "${data.alicloud_images.default.images.0.id}"
		instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
		instance_name = "${var.name}"
		system_disk_category = "cloud_efficiency"
		security_groups = ["${alicloud_security_group.default.id}"]
	}

	resource "alicloud_disk_attachment" "sample_disk_attachment" {
	  disk_id = "${alicloud_disk.sample_disk.id}"
	  instance_id = "${alicloud_instance.sample_instance.id}"
	}

	data "alicloud_disks" "disks" {
		ids = [ "${alicloud_disk.sample_disk.id}" ]
		name_regex = "${alicloud_disk.sample_disk.name}"
		type = "data"
		category = "cloud_efficiency"
		encrypted = "on"
		tags = "${alicloud_disk.sample_disk.tags}"
	    instance_id = "${alicloud_disk_attachment.sample_disk_attachment.instance_id}"
	}
	`, common)
}
