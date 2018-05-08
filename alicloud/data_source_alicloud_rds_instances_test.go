package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudRdsInstancesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudRdsInstancesDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_rds_instances.rds"),
					resource.TestCheckResourceAttr("data.alicloud_rds_instances.rds", "instances.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_rds_instances.rds", "instances.0.db_instance_class", "rds.mysql.t1.small"),
					resource.TestCheckResourceAttr("data.alicloud_rds_instances.rds", "instances.0.engine", "MySQL"),
					resource.TestCheckResourceAttr("data.alicloud_rds_instances.rds", "instances.0.engine_version", "5.6"),
					resource.TestCheckResourceAttr("data.alicloud_rds_instances.rds", "instances.0.name", "data-server"),
					resource.TestCheckResourceAttr("data.alicloud_rds_instances.rds", "instances.0.pay_type", "Postpaid"),
				),
			},
		},
	})
}

const testAccCheckAlicloudRdsInstancesDataSourceConfig = `
data "alicloud_rds_instances" "rds" {
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
