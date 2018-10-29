package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudDBDatabase_basic(t *testing.T) {
	var database rds.Database

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_db_database.db",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDBDatabaseDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDBDatabase_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDBDatabaseExists(
						"alicloud_db_database.db", &database),
					resource.TestCheckResourceAttr("alicloud_db_database.db", "character_set", "utf8"),
				),
			},
		},
	})

}

func testAccCheckDBDatabaseExists(n string, d *rds.Database) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No DB ID is set")
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		rdsService := RdsService{client}
		parts := strings.Split(rs.Primary.ID, COLON_SEPARATED)
		db, err := rdsService.DescribeDatabaseByName(parts[0], parts[1])

		if err != nil {
			return err
		}

		*d = *db
		return nil
	}
}

func testAccCheckDBDatabaseDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	rdsService := RdsService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_db_database" {
			continue
		}

		parts := strings.Split(rs.Primary.ID, COLON_SEPARATED)

		if _, err := rdsService.DescribeDatabaseByName(parts[0], parts[1]); err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}

		return fmt.Errorf("Error database %s is still existing.", parts[1])
	}

	return nil
}

const testAccDBDatabase_basic = `
variable "name" {
	default = "tf-testaccdbdatabase_basic"
}
data "alicloud_zones" "default" {
	"available_resource_creation"= "Rds"
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

resource "alicloud_db_instance" "instance" {
	engine = "MySQL"
	engine_version = "5.6"
	instance_type = "rds.mysql.t1.small"
	instance_storage = "10"
  	vswitch_id = "${alicloud_vswitch.foo.id}"
  	instance_name = "${var.name}"
}

resource "alicloud_db_database" "db" {
  instance_id = "${alicloud_db_instance.instance.id}"
  name = "${var.name}"
  description = "from terraform"
}
`
