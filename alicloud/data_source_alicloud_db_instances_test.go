package alicloud

import (
	"testing"

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
				Config: testAccCheckAlicloudDBInstancesDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_db_instances.dbs"),
					resource.TestCheckResourceAttr("data.alicloud_db_instances.dbs", "instances.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_db_instances.dbs", "instances.0.instance_type", "rds.mysql.t1.small"),
					resource.TestCheckResourceAttr("data.alicloud_db_instances.dbs", "instances.0.engine", string(MySQL)),
					resource.TestCheckResourceAttr("data.alicloud_db_instances.dbs", "instances.0.engine_version", "5.6"),
					resource.TestCheckResourceAttr("data.alicloud_db_instances.dbs", "instances.0.name", "data-server"),
					resource.TestCheckResourceAttr("data.alicloud_db_instances.dbs", "instances.0.charge_type", string(Postpaid)),
				),
			},
		},
	})
}

const testAccCheckAlicloudDBInstancesDataSourceConfig = `
data "alicloud_db_instances" "dbs" {
  name_regex = "${alicloud_db_instance.db.instance_name}"
}

resource "alicloud_db_instance" "db" {
  engine               = "MySQL"
  engine_version       = "5.6"
  instance_type        = "rds.mysql.t1.small"
  instance_storage     = "10"
  instance_name        = "data-server"
  instance_charge_type = "Postpaid"
}
`
