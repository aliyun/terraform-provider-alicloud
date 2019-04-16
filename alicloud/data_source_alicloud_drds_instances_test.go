package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudDRDSInstancesDataSourceNameRegex(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, connectivity.DrdsSupportedRegions)
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDRDSInstancesDataSourceNameRegex,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_drds_instances.dbs"),
					resource.TestCheckResourceAttr("data.alicloud_drds_instances.dbs", "instances.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_drds_instances.dbs", "instances.0.description", "tf-testAccCheckAlicloudDRDSInstancesDataSourceNameRegex"),
					resource.TestCheckResourceAttr("data.alicloud_drds_instances.dbs", "instances.0.type", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_drds_instances.dbs", "instances.0.zone_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_drds_instances.dbs", "instances.0.id"),
					resource.TestCheckResourceAttr("data.alicloud_drds_instances.dbs", "instances.0.network_type", "vpc"),
					resource.TestCheckResourceAttrSet("data.alicloud_drds_instances.dbs", "instances.0.create_time"),
					resource.TestCheckResourceAttr("data.alicloud_drds_instances.dbs", "ids.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_drds_instances.dbs", "ids.0"),
				),
			},
			{
				Config: testAccCheckAlicloudDRDSInstancesDataSourceNameRegexEmpty,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_drds_instances.dbs"),
					resource.TestCheckResourceAttr("data.alicloud_drds_instances.dbs", "instances.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_drds_instances.dbs", "ids.#", "0"),
				),
			},
		},
	})
}

func TestAccAlicloudDRDSInstancesDataSourceIds(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, connectivity.DrdsSupportedRegions)
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDRDSInstancesDataSourceIds,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_drds_instances.dbs"),
					resource.TestCheckResourceAttr("data.alicloud_drds_instances.dbs", "instances.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_drds_instances.dbs", "instances.0.description", "tf-testAccCheckAlicloudDRDSInstancesDataSourceIds"),
					resource.TestCheckResourceAttr("data.alicloud_drds_instances.dbs", "instances.0.type", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_drds_instances.dbs", "instances.0.zone_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_drds_instances.dbs", "instances.0.id"),
					resource.TestCheckResourceAttr("data.alicloud_drds_instances.dbs", "instances.0.network_type", "vpc"),
					resource.TestCheckResourceAttrSet("data.alicloud_drds_instances.dbs", "instances.0.create_time"),
					resource.TestCheckResourceAttr("data.alicloud_drds_instances.dbs", "ids.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_drds_instances.dbs", "ids.0"),
				),
			},
			{
				Config: testAccCheckAlicloudDRDSInstancesDataSourceIdsEmpty,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_drds_instances.dbs"),
					resource.TestCheckResourceAttr("data.alicloud_drds_instances.dbs", "instances.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_drds_instances.dbs", "ids.#", "0"),
				),
			},
		},
	})
}

func TestAccAlicloudDRDSInstancesDataSourceAll(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, connectivity.DrdsSupportedRegions)
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDRDSInstancesDataSourceAll,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_drds_instances.dbs"),
					resource.TestCheckResourceAttr("data.alicloud_drds_instances.dbs", "instances.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_drds_instances.dbs", "instances.0.description", "tf-testAccCheckAlicloudDRDSInstancesDataSourceAll"),
					resource.TestCheckResourceAttr("data.alicloud_drds_instances.dbs", "instances.0.type", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_drds_instances.dbs", "instances.0.zone_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_drds_instances.dbs", "instances.0.id"),
					resource.TestCheckResourceAttr("data.alicloud_drds_instances.dbs", "instances.0.network_type", "vpc"),
					resource.TestCheckResourceAttrSet("data.alicloud_drds_instances.dbs", "instances.0.create_time"),
					resource.TestCheckResourceAttr("data.alicloud_drds_instances.dbs", "ids.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_drds_instances.dbs", "ids.0"),
				),
			},
			{
				Config: testAccCheckAlicloudDRDSInstancesDataSourceAllEmpty,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_drds_instances.dbs"),
					resource.TestCheckResourceAttr("data.alicloud_drds_instances.dbs", "instances.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_drds_instances.dbs", "ids.#", "0"),
				),
			},
		},
	})
}

const testAccCheckAlicloudDRDSInstancesDataSourceNameRegex = `
 	data "alicloud_drds_instances" "dbs" {
  		name_regex = "${alicloud_drds_instance.db.description}"
	}
 	data "alicloud_zones" "default" {
		"available_resource_creation"= "VSwitch"
	}
 	variable "name" {
		default = "tf-testAccCheckAlicloudDRDSInstancesDataSourceNameRegex"
	}
 	resource "alicloud_vpc" "foo" {
		name = "${var.name}"
		cidr_block = "172.16.0.0/12"
	}
 	resource "alicloud_vswitch" "foo" {
 		vpc_id = "${alicloud_vpc.foo.id}"
 		cidr_block = "172.16.0.0/21"
 		availability_zone = "${data.alicloud_zones.default.zones.0.id}"
 		name = "${var.name}"
	}

 	resource "alicloud_drds_instance" "db" {
  		description = "${var.name}"
  		zone_id = "${data.alicloud_zones.default.zones.0.id}"
  		instance_series = "drds.sn1.4c8g"
  		instance_charge_type = "PostPaid"
		vswitch_id = "${alicloud_vswitch.foo.id}"
  		specification = "drds.sn1.4c8g.8C16G"
}
 `

const testAccCheckAlicloudDRDSInstancesDataSourceNameRegexEmpty = `
 	data "alicloud_drds_instances" "dbs" {
  		name_regex = "${alicloud_drds_instance.db.description}-fake"
	}
 	data "alicloud_zones" "default" {
		"available_resource_creation"= "VSwitch"
	}
 	variable "name" {
		default = "tf-testAccCheckAlicloudDRDSInstancesDataSourceNameRegex"
	}
 	resource "alicloud_vpc" "foo" {
		name = "${var.name}"
		cidr_block = "172.16.0.0/12"
	}
 	resource "alicloud_vswitch" "foo" {
 		vpc_id = "${alicloud_vpc.foo.id}"
 		cidr_block = "172.16.0.0/21"
 		availability_zone = "${data.alicloud_zones.default.zones.0.id}"
 		name = "${var.name}"
	}

 	resource "alicloud_drds_instance" "db" {
  		description = "${var.name}"
  		zone_id = "${data.alicloud_zones.default.zones.0.id}"
  		instance_series = "drds.sn1.4c8g"
  		instance_charge_type = "PostPaid"
		vswitch_id = "${alicloud_vswitch.foo.id}"
  		specification = "drds.sn1.4c8g.8C16G"
}
 `

const testAccCheckAlicloudDRDSInstancesDataSourceIds = `
 	data "alicloud_drds_instances" "dbs" {
  		ids = ["${alicloud_drds_instance.db.id}"]
	}
 	data "alicloud_zones" "default" {
		"available_resource_creation"= "VSwitch"
	}
 	variable "name" {
		default = "tf-testAccCheckAlicloudDRDSInstancesDataSourceIds"
	}
 	resource "alicloud_vpc" "foo" {
		name = "${var.name}"
		cidr_block = "172.16.0.0/12"
	}
 	resource "alicloud_vswitch" "foo" {
 		vpc_id = "${alicloud_vpc.foo.id}"
 		cidr_block = "172.16.0.0/21"
 		availability_zone = "${data.alicloud_zones.default.zones.0.id}"
 		name = "${var.name}"
	}

 	resource "alicloud_drds_instance" "db" {
  		description = "${var.name}"
  		zone_id = "${data.alicloud_zones.default.zones.0.id}"
  		instance_series = "drds.sn1.4c8g"
  		instance_charge_type = "PostPaid"
		vswitch_id = "${alicloud_vswitch.foo.id}"
  		specification = "drds.sn1.4c8g.8C16G"
}
 `

const testAccCheckAlicloudDRDSInstancesDataSourceIdsEmpty = `
 	data "alicloud_drds_instances" "dbs" {
  		ids = ["${alicloud_drds_instance.db.id}-fake"]
	}
 	data "alicloud_zones" "default" {
		"available_resource_creation"= "VSwitch"
	}
 	variable "name" {
		default = "tf-testAccCheckAlicloudDRDSInstancesDataSourceIds"
	}
 	resource "alicloud_vpc" "foo" {
		name = "${var.name}"
		cidr_block = "172.16.0.0/12"
	}
 	resource "alicloud_vswitch" "foo" {
 		vpc_id = "${alicloud_vpc.foo.id}"
 		cidr_block = "172.16.0.0/21"
 		availability_zone = "${data.alicloud_zones.default.zones.0.id}"
 		name = "${var.name}"
	}

 	resource "alicloud_drds_instance" "db" {
  		description = "${var.name}"
  		zone_id = "${data.alicloud_zones.default.zones.0.id}"
  		instance_series = "drds.sn1.4c8g"
  		instance_charge_type = "PostPaid"
		vswitch_id = "${alicloud_vswitch.foo.id}"
  		specification = "drds.sn1.4c8g.8C16G"
}
 `

const testAccCheckAlicloudDRDSInstancesDataSourceAll = `
 	data "alicloud_drds_instances" "dbs" {
  		ids = ["${alicloud_drds_instance.db.id}"]
  		name_regex = "${alicloud_drds_instance.db.description}"
	}
 	data "alicloud_zones" "default" {
		"available_resource_creation"= "VSwitch"
	}
 	variable "name" {
		default = "tf-testAccCheckAlicloudDRDSInstancesDataSourceAll"
	}
 	resource "alicloud_vpc" "foo" {
		name = "${var.name}"
		cidr_block = "172.16.0.0/12"
	}
 	resource "alicloud_vswitch" "foo" {
 		vpc_id = "${alicloud_vpc.foo.id}"
 		cidr_block = "172.16.0.0/21"
 		availability_zone = "${data.alicloud_zones.default.zones.0.id}"
 		name = "${var.name}"
	}

 	resource "alicloud_drds_instance" "db" {
  		description = "${var.name}"
  		zone_id = "${data.alicloud_zones.default.zones.0.id}"
  		instance_series = "drds.sn1.4c8g"
  		instance_charge_type = "PostPaid"
		vswitch_id = "${alicloud_vswitch.foo.id}"
  		specification = "drds.sn1.4c8g.8C16G"
}
 `

const testAccCheckAlicloudDRDSInstancesDataSourceAllEmpty = `
 	data "alicloud_drds_instances" "dbs" {
  		ids = ["${alicloud_drds_instance.db.id}"]
  		name_regex = "${alicloud_drds_instance.db.description}-fake"
	}
 	data "alicloud_zones" "default" {
		"available_resource_creation"= "VSwitch"
	}
 	variable "name" {
		default = "tf-testAccCheckAlicloudDRDSInstancesDataSourceAll"
	}
 	resource "alicloud_vpc" "foo" {
		name = "${var.name}"
		cidr_block = "172.16.0.0/12"
	}
 	resource "alicloud_vswitch" "foo" {
 		vpc_id = "${alicloud_vpc.foo.id}"
 		cidr_block = "172.16.0.0/21"
 		availability_zone = "${data.alicloud_zones.default.zones.0.id}"
 		name = "${var.name}"
	}

 	resource "alicloud_drds_instance" "db" {
  		description = "${var.name}"
  		zone_id = "${data.alicloud_zones.default.zones.0.id}"
  		instance_series = "drds.sn1.4c8g"
  		instance_charge_type = "PostPaid"
		vswitch_id = "${alicloud_vswitch.foo.id}"
  		specification = "drds.sn1.4c8g.8C16G"
}
 `
