package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudDRDSInstancesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDRDSInstancesDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_drds_instances.dbs"),
					resource.TestCheckResourceAttr("data.alicloud_drds_instances.dbs", "instances.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_drds_instances.dbs", "instances.0.description", "tf-testAccCheckAlicloudDRDSInstancesDataSourceConfig"),
					resource.TestCheckResourceAttrSet("data.alicloud_db_instances.dbs", "instances.0.type"),
					resource.TestCheckResourceAttrSet("data.alicloud_db_instances.dbs", "instances.0.zone_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_db_instances.dbs", "instances.0.version"),
					resource.TestCheckResourceAttrSet("data.alicloud_db_instances.dbs", "instances.0.id"),
					resource.TestCheckResourceAttrSet("data.alicloud_db_instances.dbs", "instances.0.network_type"),
					resource.TestCheckResourceAttrSet("data.alicloud_db_instances.dbs", "instances.0.create_time"),
				),
			},
		},
	})
}

const testAccCheckAlicloudDRDSInstancesDataSourceConfig = `
 	data "alicloud_drds_instances" "dbs" {
  		name_regex = "${alicloud_drds_instance.db.description}"
	}
 	data "alicloud_zones" "default" {
		"available_resource_creation"= "VSwitch"
	}
 	variable "name" {
		default = "tf-testaccDrdsdatabase_vpc"
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
  		description = "tf-testAccCheckAlicloudDRDSInstancesDataSourceConfig"
  		zone_id = "${data.alicloud_zones.default.zones.0.id}"
  		instance_series = "drds.sn1.4c8g"
  		instance_charge_type = "PostPaid"
		vswitch_id = "${alicloud_vswitch.foo.id}"
  		specification = "drds.sn1.4c8g.8C16G"
}
 `
