package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func testAccCheckDBBackupPolicyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_db_backup_policy" {
			continue
		}
		request := rds.CreateDescribeBackupPolicyRequest()
		request.DBInstanceId = rs.Primary.ID
		_, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
			return rdsClient.DescribeBackupPolicy(request)
		})
		if err != nil {
			if IsExceptedError(err, InvalidDBInstanceIdNotFound) || IsExceptedError(err, InvalidDBInstanceNameNotFound) {
				continue
			}
			return WrapError(err)
		}
	}
	return nil
}

func TestAccAlicloudDBBackupPolicy_mysql(t *testing.T) {
	var v *rds.DescribeBackupPolicyResponse
	resourceId := "alicloud_db_backup_policy.default"
	serverFunc := func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serverFunc, "DescribeBackupPolicy")
	ra := resourceAttrInit(resourceId, nil)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := "tf-testAccDBbackuppolicy"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBBackupPolicyMysqlConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckDBBackupPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id": "${alicloud_db_instance.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"log_retention_period"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"backup_period": []string{"Wednesday"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_period.#":          "1",
						"backup_period.1970423419": "Wednesday",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"backup_time": "10:00Z-11:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_time": "10:00Z-11:00Z",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"retention_period": "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"retention_period": "10",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"log_backup": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"log_backup": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"log_backup": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"log_backup": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"log_retention_period": "7",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"log_retention_period": "7",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":          "${alicloud_db_instance.default.id}",
					"backup_period":        []string{"Tuesday", "Wednesday"},
					"backup_time":          "10:00Z-11:00Z",
					"retention_period":     "10",
					"log_backup":           "true",
					"log_retention_period": "7",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_period.#":          "2",
						"backup_period.1592931319": "Tuesday",
						"backup_period.1970423419": "Wednesday",
						"backup_time":              "10:00Z-11:00Z",
						"retention_period":         "10",
					}),
				),
			}},
	})
}

func resourceDBBackupPolicyMysqlConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
data "alicloud_zones" "default" {
  available_resource_creation = "Rds"
}

	data "alicloud_db_instance_engines" "default" {
  		instance_charge_type = "PostPaid"
  		engine               = "MySQL"
  		engine_version       = "5.6"
	}

	data "alicloud_db_instance_classes" "default" {
 	 	engine = "${data.alicloud_db_instance_engines.default.instance_engines.0.engine}"
		engine_version = "${data.alicloud_db_instance_engines.default.instance_engines.0.engine_version}"
	}

resource "alicloud_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vpc_id            = "${alicloud_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alicloud_db_instance_classes.default.instance_classes.0.zone_ids.0.sub_zone_ids.0}"
  name              = "${var.name}"
}
resource "alicloud_db_instance" "default" {
  	vswitch_id       = "${alicloud_vswitch.default.id}"
  	instance_name    = "${var.name}"
  	engine 			 = "${data.alicloud_db_instance_engines.default.instance_engines.0.engine}"
	engine_version   = "${data.alicloud_db_instance_engines.default.instance_engines.0.engine_version}"
	instance_type    = "${data.alicloud_db_instance_classes.default.instance_classes.0.instance_class}"
  	instance_storage = "${data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min}"
}`, name)
}

func TestAccAlicloudDBBackupPolicy_pgdb(t *testing.T) {
	var v *rds.DescribeBackupPolicyResponse
	resourceId := "alicloud_db_backup_policy.default"
	serverFunc := func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serverFunc, "DescribeBackupPolicy")
	ra := resourceAttrInit(resourceId, nil)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := "tf-testAccDBbackuppolicy"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBBackupPolicyPostgreSQLConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckDBBackupPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id": "${alicloud_db_instance.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"backup_period": []string{"Wednesday"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_period.#":          "1",
						"backup_period.1970423419": "Wednesday",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"backup_time": "10:00Z-11:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_time": "10:00Z-11:00Z",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"retention_period": "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"retention_period": "10",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"log_backup": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"log_backup": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"log_backup": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"log_backup": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"log_retention_period": "7",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"log_retention_period": "7",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"backup_period":    []string{"Tuesday", "Wednesday"},
					"backup_time":      "10:00Z-11:00Z",
					"retention_period": "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_period.#":          "2",
						"backup_period.1592931319": "Tuesday",
						"backup_period.1970423419": "Wednesday",
						"backup_time":              "10:00Z-11:00Z",
						"retention_period":         "10",
					}),
				),
			}},
	})
}

func resourceDBBackupPolicyPostgreSQLConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
data "alicloud_zones" "default" {
  available_resource_creation = "Rds"
}
data "alicloud_db_instance_engines" "default" {
	engine               = "PostgreSQL"
	engine_version       = "10.0"
	instance_charge_type = "PostPaid"
    multi_zone           = true
}

data "alicloud_db_instance_classes" "default" {
	engine               = "PostgreSQL"
	engine_version       = "10.0"
	category             = "HighAvailability"
	instance_charge_type = "PostPaid"
    multi_zone           = true
}

resource "alicloud_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vpc_id            = "${alicloud_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alicloud_db_instance_classes.default.instance_classes.0.zone_ids.0.sub_zone_ids.0}"
  name              = "${var.name}"
}
resource "alicloud_db_instance" "default" {
  	vswitch_id       = "${alicloud_vswitch.default.id}"
  	instance_name    = "${var.name}"
  	engine 			 = "${data.alicloud_db_instance_engines.default.instance_engines.0.engine}"
	engine_version   = "${data.alicloud_db_instance_engines.default.instance_engines.0.engine_version}"
	instance_type    = "${data.alicloud_db_instance_classes.default.instance_classes.0.instance_class}"
  	instance_storage = "${data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min}"
	zone_id          = "${data.alicloud_db_instance_classes.default.instance_classes.0.zone_ids.0.id}"
}`, name)
}

func TestAccAlicloudDBBackupPolicy_SQLServer(t *testing.T) {
	var v *rds.DescribeBackupPolicyResponse
	resourceId := "alicloud_db_backup_policy.default"
	serverFunc := func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serverFunc, "DescribeBackupPolicy")
	ra := resourceAttrInit(resourceId, nil)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := "tf-testAccDBbackuppolicy"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBBackupPolicySQLServerConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckDBBackupPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id": "${alicloud_db_instance.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"backup_period": []string{"Wednesday"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_period.#":          "1",
						"backup_period.1970423419": "Wednesday",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"backup_time": "10:00Z-11:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_time": "10:00Z-11:00Z",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"retention_period": "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"retention_period": "10",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"backup_period":    []string{"Wednesday", "Tuesday"},
					"backup_time":      "10:00Z-11:00Z",
					"retention_period": "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_period.#":          "2",
						"backup_period.1592931319": "Tuesday",
						"backup_period.1970423419": "Wednesday",
						"backup_time":              "10:00Z-11:00Z",
						"retention_period":         "10",
					}),
				),
			}},
	})
}

func resourceDBBackupPolicySQLServerConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
data "alicloud_zones" "default" {
  available_resource_creation = "Rds"
}
data "alicloud_db_instance_engines" "default" {
	engine               = "SQLServer"
	engine_version       = "2008r2"
}

data "alicloud_db_instance_classes" "default" {
	engine = "${data.alicloud_db_instance_engines.default.instance_engines.0.engine}"
	engine_version = "${data.alicloud_db_instance_engines.default.instance_engines.0.engine_version}"
}
resource "alicloud_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vpc_id            = "${alicloud_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alicloud_db_instance_classes.default.instance_classes.0.zone_ids.0.sub_zone_ids.0}"
  name              = "${var.name}"
}
resource "alicloud_db_instance" "default" {
  	vswitch_id       = "${alicloud_vswitch.default.id}"
  	instance_name    = "${var.name}"
  	engine 			 = "${data.alicloud_db_instance_engines.default.instance_engines.0.engine}"
	engine_version   = "${data.alicloud_db_instance_engines.default.instance_engines.0.engine_version}"
	instance_type    = "${data.alicloud_db_instance_classes.default.instance_classes.0.instance_class}"
  	instance_storage = "${data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min}"
}`, name)
}

// Unknown current resource exists
func TestAccAlicloudDBBackupPolicy_PPAS(t *testing.T) {
	var v *rds.DescribeBackupPolicyResponse
	resourceId := "alicloud_db_backup_policy.default"
	serverFunc := func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serverFunc, "DescribeBackupPolicy")
	name := "tf-testAccDBbackuppolicy"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBBackupPolicyPPASConfigDependence)
	ra := resourceAttrInit(resourceId, nil)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckDBBackupPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id": "${alicloud_db_instance.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"backup_period": []string{"Wednesday"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_period.#":          "1",
						"backup_period.1970423419": "Wednesday",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"backup_time": "10:00Z-11:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_time": "10:00Z-11:00Z",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"retention_period": "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"retention_period": "10",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"log_backup": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"log_backup": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"log_backup": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"log_backup": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"log_retention_period": "7",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"log_retention_period": "7",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"backup_period":    []string{"Wednesday", "Tuesday"},
					"backup_time":      "10:00Z-11:00Z",
					"retention_period": "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_period.#":          "2",
						"backup_period.1592931319": "Tuesday",
						"backup_period.1970423419": "Wednesday",
						"backup_time":              "10:00Z-11:00Z",
						"retention_period":         "10",
					}),
				),
			}},
	})
}

func resourceDBBackupPolicyPPASConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
data "alicloud_zones" "default" {
  available_resource_creation = "Rds"
}
data "alicloud_db_instance_engines" "default" {
	engine               = "PPAS"
	engine_version       = "10.0"
    instance_charge_type = "PostPaid"
    multi_zone           = true
}

data "alicloud_db_instance_classes" "default" {
	engine               = "PPAS"
	engine_version       = "10.0"
    instance_charge_type = "PostPaid"
    multi_zone           = true
}
resource "alicloud_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vpc_id            = "${alicloud_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alicloud_db_instance_classes.default.instance_classes.0.zone_ids.0.sub_zone_ids.0}"
  name              = "${var.name}"
}
resource "alicloud_db_instance" "default" {
  	vswitch_id       = "${alicloud_vswitch.default.id}"
  	instance_name    = "${var.name}"
  	engine 			 = "${data.alicloud_db_instance_engines.default.instance_engines.0.engine}"
	engine_version   = "${data.alicloud_db_instance_engines.default.instance_engines.0.engine_version}"
	instance_type    = "${data.alicloud_db_instance_classes.default.instance_classes.0.instance_class}"
  	instance_storage = "${data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min}"
	zone_id          = "${data.alicloud_db_instance_classes.default.instance_classes.0.zone_ids.0.id}"
}`, name)
}
