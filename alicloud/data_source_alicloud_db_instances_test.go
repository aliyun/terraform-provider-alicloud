package alicloud

import (
	"testing"

	"fmt"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudDBInstancesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDBInstancesDataSourceConfig(RdsCommonTestCase),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_db_instances.dbs"),
					resource.TestCheckResourceAttr("data.alicloud_db_instances.dbs", "instances.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_db_instances.dbs", "instances.0.instance_type", "rds.mysql.t1.small"),
					resource.TestCheckResourceAttr("data.alicloud_db_instances.dbs", "instances.0.engine", string(MySQL)),
					resource.TestCheckResourceAttr("data.alicloud_db_instances.dbs", "instances.0.engine_version", "5.6"),
					resource.TestCheckResourceAttr("data.alicloud_db_instances.dbs", "instances.0.name", "tf-testAccCheckAlicloudDBInstancesDataSourceConfig"),
					resource.TestCheckResourceAttr("data.alicloud_db_instances.dbs", "instances.0.charge_type", string(Postpaid)),
				),
			},
		},
	})
}

func TestAccAlicloudDBInstancesDataSourceEmpty(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDBInstancesDataSourceEmpty,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_db_instances.dbs"),
					resource.TestCheckResourceAttr("data.alicloud_db_instances.dbs", "instances.#", "0"),
					resource.TestCheckNoResourceAttr("data.alicloud_db_instances.dbs", "instances.0.instance_type"),
					resource.TestCheckNoResourceAttr("data.alicloud_db_instances.dbs", "instances.0.engine"),
					resource.TestCheckNoResourceAttr("data.alicloud_db_instances.dbs", "instances.0.engine_version"),
					resource.TestCheckNoResourceAttr("data.alicloud_db_instances.dbs", "instances.0.name"),
					resource.TestCheckNoResourceAttr("data.alicloud_db_instances.dbs", "instances.0.charge_type"),
				),
			},
		},
	})
}

func testAccCheckAlicloudDBInstancesDataSourceConfig(common string) string {
	return fmt.Sprintf(`
	%s
	data "alicloud_db_instances" "dbs" {
	  name_regex = "${alicloud_db_instance.db.instance_name}"
	}
	variable "creation" {
		default = "Rds"
	}
	variable "name" {
		default = "tf-testAccCheckAlicloudDBInstancesDataSourceConfig"
	}
	resource "alicloud_db_instance" "db" {
	  engine               = "MySQL"
	  engine_version       = "5.6"
	  instance_type        = "rds.mysql.t1.small"
	  instance_storage     = "10"
	  instance_name        = "${var.name}"
	  instance_charge_type = "Postpaid"
	  vswitch_id = "${alicloud_vswitch.default.id}"
	}
	`, common)
}

const testAccCheckAlicloudDBInstancesDataSourceEmpty = `
data "alicloud_db_instances" "dbs" {
  name_regex = "^tf-testacc-fake-name*"
}
`
