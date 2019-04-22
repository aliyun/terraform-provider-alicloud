package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

var DBReadonlyMap = map[string]string{
	"instance_storage":      "30",
	"engine_version":        "5.6",
	"engine":                "MySQL",
	"port":                  "3306",
	"instance_name":         "tf-testAccDBInstance_vpc_ro",
	"instance_type":         "rds.mysql.t1.small",
	"parameters":            NOSET,
	"master_db_instance_id": CHECKSET,
	"zone_id":               CHECKSET,
	"vswitch_id":            CHECKSET,
	"connection_string":     CHECKSET,
}

func TestAccAlicloudDBReadonlyInstance_update(t *testing.T) {
	var instance *rds.DBInstanceAttribute
	resourceId := "alicloud_db_readonly_instance.default"
	ra := resourceAttrInit(resourceId, DBReadonlyMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBReadonlyInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDBReadonlyInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDBReadonlyInstance_vpc(testAccDBRInstance_vpc(RdsCommonTestCase)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			// upgrade storage
			{
				Config: testAccDBReadonlyInstance_upgrade(testAccDBRInstance_vpc(RdsCommonTestCase),
					"${alicloud_db_instance.default.instance_type}", "40"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"instance_storage": "40"}),
				),
			},
			// upgrade instanceType
			{
				Config: testAccDBReadonlyInstance_upgrade(testAccDBRInstance_vpc(RdsCommonTestCase),
					"rds.mysql.s1.small", "40"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"instance_type": "rds.mysql.s1.small"}),
				),
			},
			{
				Config: testAccDBReadonlyInstance_multiAZ_vpc(testAccDBInstance_vpc_multiAZ(DBMultiAZCommonTestCase)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":    "tf-testAccDBInstance_multiAZ_ro",
						"instance_storage": "30",
						"instance_type":    "rds.mysql.s2.large",
					}),
				),
			},
			{
				Config: testAccDBReadonlyInstance_vpc_instanceName(
					testAccDBRInstance_vpc(RdsCommonTestCase), "some_other_name"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"instance_name": "some_other_name", "instance_type": "rds.mysql.t1.small"}),
				),
			},
		},
	})

}

func TestAccAlicloudDBReadonlyInstance_multi(t *testing.T) {
	var instance *rds.DBInstanceAttribute
	resourceId := "alicloud_db_readonly_instance.default.4"
	ra := resourceAttrInit(resourceId, DBReadonlyMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDBReadonlyInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDBReadonlyInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDBReadonlyInstance_mulit(testAccDBRInstance_vpc(RdsCommonTestCase)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func testAccDBRInstance_vpc(common string) string {
	return fmt.Sprintf(`
	%s
	variable "creation" {
		default = "Rds"
	}
	variable "multi_az" {
		default = "false"
	}
	variable "name" {
		default = "tf-testAccDBInstance_vpc"
	}

	resource "alicloud_db_instance" "default" {
		engine = "MySQL"
		engine_version = "5.6"
		instance_type = "rds.mysql.t1.small"
		instance_storage = "20"
		instance_charge_type = "Postpaid"
		instance_name = "${var.name}"
		vswitch_id = "${alicloud_vswitch.default.id}"
		security_ips = ["10.168.1.12", "100.69.7.112"]
	}
	`, common)
}

func testAccDBReadonlyInstance_vpc(common string) string {
	return fmt.Sprintf(`
	%s
	resource "alicloud_db_readonly_instance" "default" {
		master_db_instance_id = "${alicloud_db_instance.default.id}"
		zone_id = "${alicloud_db_instance.default.zone_id}"
		engine_version = "${alicloud_db_instance.default.engine_version}"
		instance_type = "${alicloud_db_instance.default.instance_type}"
		instance_storage = "30"
		instance_name = "${var.name}_ro"
		vswitch_id = "${alicloud_vswitch.default.id}"
	}
	`, common)
}

func testAccDBReadonlyInstance_vpc_instanceName(common, instanceName string) string {
	return fmt.Sprintf(`
	%s
	resource "alicloud_db_readonly_instance" "default" {
		master_db_instance_id = "${alicloud_db_instance.default.id}"
		zone_id = "${alicloud_db_instance.default.zone_id}"
		engine_version = "${alicloud_db_instance.default.engine_version}"
		instance_type = "${alicloud_db_instance.default.instance_type}"
		instance_storage = "30"
		instance_name = "%s"
		vswitch_id = "${alicloud_vswitch.default.id}"
	}
	`, common, instanceName)
}

func testAccDBReadonlyInstance_upgrade(common, instanceType, storage string) string {
	return fmt.Sprintf(`
	%s
	resource "alicloud_db_readonly_instance" "default" {
		master_db_instance_id = "${alicloud_db_instance.default.id}"
		zone_id = "${alicloud_db_instance.default.zone_id}"
		engine_version = "${alicloud_db_instance.default.engine_version}"
		instance_type = "%s"
		instance_storage = "%s"
		instance_name = "${var.name}_ro"
		vswitch_id = "${alicloud_vswitch.default.id}"
	}
	`, common, instanceType, storage)
}

func testAccDBReadonlyInstance_multiAZ(common string) string {
	return fmt.Sprintf(`
	%s
	resource "alicloud_db_readonly_instance" "default" {
		master_db_instance_id = "${alicloud_db_instance.default.id}"
		zone_id = "${alicloud_db_instance.default.zone_id}"
		engine_version = "${alicloud_db_instance.default.engine_version}"
		instance_type = "${alicloud_db_instance.default.instance_type}"
		instance_storage = "30"
		instance_name = "${var.name}_ro"
	}
	`, common)
}

func testAccDBReadonlyInstance_multiAZ_vpc(common string) string {
	return fmt.Sprintf(`
	%s
	resource "alicloud_db_readonly_instance" "default" {
		master_db_instance_id = "${alicloud_db_instance.default.id}"
		zone_id = "${alicloud_db_instance.default.zone_id}"
		engine_version = "${alicloud_db_instance.default.engine_version}"
		instance_type = "${alicloud_db_instance.default.instance_type}"
		instance_storage = "30"
		instance_name = "${var.name}_ro"
		vswitch_id = "${alicloud_vswitch.default.id}"
	}
	`, common)
}

func testAccCheckDBReadonlyInstanceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	rdsService := RdsService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_db_readonly_instance" {
			continue
		}

		ins, err := rdsService.DescribeDBInstance(rs.Primary.ID)

		if ins != nil {
			return fmt.Errorf("Error DB Instance still exist")
		}

		// Verify the error is what we want
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}
	}

	return nil
}

func testAccDBReadonlyInstance_mulit(common string) string {
	return fmt.Sprintf(`
	%s
	resource "alicloud_db_readonly_instance" "default" {
		count = 5
		master_db_instance_id = "${alicloud_db_instance.default.id}"
		zone_id = "${alicloud_db_instance.default.zone_id}"
		engine_version = "${alicloud_db_instance.default.engine_version}"
		instance_type = "${alicloud_db_instance.default.instance_type}"
		instance_storage = "30"
		instance_name = "${var.name}_ro"
		vswitch_id = "${alicloud_vswitch.default.id}"
	}
	`, common)
}
