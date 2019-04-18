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

func TestAccAlicloudDBDatabase_update(t *testing.T) {
	var database *rds.Database
	resourceId := "alicloud_db_database.default"
	ra := resourceAttrInit(resourceId, dbDatabaseBasicMap)
	rc := resourceCheckInit(resourceId, &database, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDBDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDBDatabase_basic(RdsCommonTestCase),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				Config: testAccDBDatabase_description(RdsCommonTestCase),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"description": "from terraform"}),
				),
			},
		},
	})

}

func testAccCheckDBDatabaseDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	rdsService := RdsService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_db_database" {
			continue
		}

		parts := strings.Split(rs.Primary.ID, COLON_SEPARATED)

		if _, err := rdsService.DescribeDBDatabase(rs.Primary.ID); err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}

		return fmt.Errorf("Error database %s is still existing.", parts[1])
	}

	return nil
}

var dbDatabaseBasicMap = map[string]string{
	"instance_id":   CHECKSET,
	"name":          "tftestdatabase",
	"character_set": "utf8",
	"description":   "",
}

func testAccDBDatabase_basic(common string) string {
	return fmt.Sprintf(`
	%s
	variable "creation" {
		default = "Rds"
	}

	variable "name" {
		default = "tf-testAccDBdatabase_basic"
	}

	resource "alicloud_db_instance" "instance" {
		engine = "MySQL"
		engine_version = "5.6"
		instance_type = "rds.mysql.s1.small"
		instance_storage = "10"
		vswitch_id = "${alicloud_vswitch.default.id}"
		instance_name = "${var.name}"
	}

	resource "alicloud_db_database" "default" {
	  instance_id = "${alicloud_db_instance.instance.id}"
	  name = "tftestdatabase"
	}
	`, common)
}

func testAccDBDatabase_description(common string) string {
	return fmt.Sprintf(`
	%s
	variable "creation" {
		default = "Rds"
	}

	variable "name" {
		default = "tf-testAccDBdatabase_basic"
	}

	resource "alicloud_db_instance" "instance" {
		engine = "MySQL"
		engine_version = "5.6"
		instance_type = "rds.mysql.s1.small"
		instance_storage = "10"
		vswitch_id = "${alicloud_vswitch.default.id}"
		instance_name = "${var.name}"
	}

	resource "alicloud_db_database" "default" {
	  instance_id = "${alicloud_db_instance.instance.id}"
	  name = "tftestdatabase"
	  description = "from terraform"
	}
	`, common)
}
