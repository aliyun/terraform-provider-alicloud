package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudRdsDBBackupPolicyMySql(t *testing.T) {
	var v map[string]interface{}
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
		//CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":                       "${alicloud_db_instance.default.id}",
					"local_log_retention_hours":         "18",
					"high_space_usage_protection":       "Enable",
					"log_backup":                        "false",
					"enable_backup_log":                 "false",
					"preferred_backup_time":             "02:00Z-03:00Z",
					"backup_time":                       "02:00Z-03:00Z",
					"retention_period":                  "900",
					"backup_retention_period":           "900",
					"backup_period":                     []string{"Wednesday"},
					"log_retention_period":              "7",
					"log_backup_local_retention_number": "20",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":                       CHECKSET,
						"log_backup_local_retention_number": CHECKSET,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"log_backup":        "true",
					"enable_backup_log": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"log_backup":        "true",
						"enable_backup_log": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"backup_retention_period": "910",
					"retention_period":        "910",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"retention_period":        "910",
						"backup_retention_period": "910",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"backup_time":           "03:00Z-04:00Z",
					"preferred_backup_time": "03:00Z-04:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_time":           "03:00Z-04:00Z",
						"preferred_backup_time": "03:00Z-04:00Z",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"log_backup_retention_period": "8",
					"log_retention_period":        "8",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"log_backup_retention_period": "8",
						"log_retention_period":        "8",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"local_log_retention_hours": "24",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"local_log_retention_hours": "24",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"local_log_retention_space": "35",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"local_log_retention_space": "35",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"high_space_usage_protection": "Disable",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"high_space_usage_protection": "Disable",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"compress_type": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"compress_type": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"archive_backup_retention_period": "150",
					"archive_backup_keep_count":       "1",
					"archive_backup_keep_policy":      "ByMonth",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"archive_backup_retention_period": "150",
						"archive_backup_keep_count":       "1",
						"archive_backup_keep_policy":      "ByMonth",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"archive_backup_keep_policy": "KeepAll",
					"archive_backup_keep_count":  "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"archive_backup_keep_policy": "KeepAll",
						"archive_backup_keep_count":  "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"released_keep_policy": "Lastest",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"released_keep_policy": "Lastest",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"log_backup_local_retention_number": "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"log_backup_local_retention_number": "10",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":                     "${alicloud_db_instance.default.id}",
					"preferred_backup_period":         []string{"Wednesday", "Monday"},
					"backup_period":                   []string{"Wednesday", "Monday"},
					"backup_time":                     "03:00Z-04:00Z",
					"retention_period":                "910",
					"backup_retention_period":         "910",
					"log_backup":                      "true",
					"log_retention_period":            "8",
					"local_log_retention_hours":       "48",
					"high_space_usage_protection":     "Enable",
					"archive_backup_retention_period": "100",
					"archive_backup_keep_count":       "3",
					"archive_backup_keep_policy":      "ByWeek",
					"category":                        "Standard",
					"enable_backup_log":               "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"preferred_backup_period.#":       "2",
						"backup_period.#":                 "2",
						"backup_time":                     "03:00Z-04:00Z",
						"retention_period":                "910",
						"backup_retention_period":         "910",
						"log_backup":                      "true",
						"log_retention_period":            "8",
						"local_log_retention_hours":       "48",
						"high_space_usage_protection":     "Enable",
						"archive_backup_retention_period": "100",
						"archive_backup_keep_count":       "3",
						"archive_backup_keep_policy":      "ByWeek",
						"category":                        "Standard",
						"enable_backup_log":               "true",
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
data "alicloud_db_zones" "default"{
	engine = "MySQL"
	engine_version = "8.0"
	instance_charge_type = "PostPaid"
	category = "HighAvailability"
 	db_instance_storage_type = "local_ssd"
}

data "alicloud_db_instance_classes" "default" {
    zone_id = data.alicloud_db_zones.default.zones.0.id
	engine = "MySQL"
	engine_version = "8.0"
    category = "HighAvailability"
 	db_instance_storage_type = "local_ssd"
	instance_charge_type = "PostPaid"
}

data "alicloud_vpcs" "default" {
    name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_db_zones.default.zones.0.id
}

resource "alicloud_vswitch" "this" {
 count = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
 vswitch_name = var.name
 vpc_id = data.alicloud_vpcs.default.ids.0
 zone_id = data.alicloud_db_zones.default.ids.0
 cidr_block = cidrsubnet(data.alicloud_vpcs.default.vpcs.0.cidr_block, 8, 4)
}
locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids.0 : concat(alicloud_vswitch.this.*.id, [""])[0]
  zone_id = data.alicloud_db_zones.default.ids.0
}

data "alicloud_resource_manager_resource_groups" "default" {
	status = "OK"
}

resource "alicloud_security_group" "default" {
	security_group_name   = var.name
	vpc_id 				  = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_db_instance" "default" {
    engine = "MySQL"
	engine_version = "8.0"
 	db_instance_storage_type = "local_ssd"
	instance_type = data.alicloud_db_instance_classes.default.instance_classes.0.instance_class
	instance_storage = data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min
	vswitch_id = local.vswitch_id
	instance_name = var.name
	security_group_ids = alicloud_security_group.default.*.id
}
`, name)
}

func TestAccAliCloudRdsDBBackupPolicyPostgreSQL(t *testing.T) {
	var v map[string]interface{}
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
		//CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_id":              "${alicloud_db_instance.default.id}",
					"enable_backup_log":           "true",
					"local_log_retention_hours":   "1",
					"high_space_usage_protection": "Enable",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_id": CHECKSET,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"preferred_backup_period": []string{"Monday", "Wednesday"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"preferred_backup_period.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"preferred_backup_time": "10:00Z-11:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"preferred_backup_time": "10:00Z-11:00Z",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"backup_retention_period": "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_retention_period": "10",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"log_backup_retention_period": "7",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"log_backup_retention_period": "7",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"local_log_retention_hours": "24",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"local_log_retention_hours": "24",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"local_log_retention_space": "35",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"local_log_retention_space": "35",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"high_space_usage_protection": "Disable",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"high_space_usage_protection": "Disable",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"compress_type": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"compress_type": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"backup_interval": "60",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_interval": "60",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_backup_log": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_backup_log": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"preferred_backup_period":     []string{"Tuesday", "Wednesday", "Monday"},
					"preferred_backup_time":       "11:00Z-12:00Z",
					"backup_retention_period":     "20",
					"enable_backup_log":           "true",
					"log_backup_retention_period": "7",
					"local_log_retention_hours":   "48",
					"high_space_usage_protection": "Enable",
					"local_log_retention_space":   "40",
					"compress_type":               "1",
					"backup_interval":             "60",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"preferred_backup_period.#":   "3",
						"preferred_backup_time":       "11:00Z-12:00Z",
						"backup_retention_period":     "20",
						"enable_backup_log":           "true",
						"log_backup_retention_period": "7",
						"local_log_retention_hours":   "48",
						"high_space_usage_protection": "Enable",
						"local_log_retention_space":   "40",
						"compress_type":               "1",
						"backup_interval":             "60",
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
data "alicloud_db_zones" "default"{
	engine = "PostgreSQL"
	engine_version = "14.0"
	instance_charge_type = "PostPaid"
	category = "HighAvailability"
 	db_instance_storage_type = "cloud_essd"
}

data "alicloud_db_instance_classes" "default" {
    zone_id = data.alicloud_db_zones.default.zones.0.id
	engine = "PostgreSQL"
	engine_version = "14.0"
    category = "HighAvailability"
 	db_instance_storage_type = "cloud_essd"
	instance_charge_type = "PostPaid"
}

data "alicloud_vpcs" "default" {
    name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_db_zones.default.zones.0.id
}

resource "alicloud_vswitch" "this" {
 count = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
 vswitch_name = var.name
 vpc_id = data.alicloud_vpcs.default.ids.0
 zone_id = data.alicloud_db_zones.default.zones.0.id
 cidr_block = cidrsubnet(data.alicloud_vpcs.default.vpcs.0.cidr_block, 8, 4)
}
locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids.0 : concat(alicloud_vswitch.this.*.id, [""])[0]
  zone_id = data.alicloud_db_zones.default.ids[length(data.alicloud_db_zones.default.ids)-1]
}

data "alicloud_resource_manager_resource_groups" "default" {
	status = "OK"
}

resource "alicloud_security_group" "default" {
	security_group_name   = var.name
	vpc_id 				  = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_db_instance" "default" {
    engine = "PostgreSQL"
	engine_version = "14.0"
 	db_instance_storage_type = "cloud_essd"
	instance_type = data.alicloud_db_instance_classes.default.instance_classes.0.instance_class
	instance_storage = data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min
	vswitch_id = local.vswitch_id
	instance_name = var.name
	security_group_ids = alicloud_security_group.default.*.id

}
`, name)
}

func TestAccAliCloudRdsDBBackupPolicySQLServer(t *testing.T) {
	var v map[string]interface{}
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
		//CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_id": "${alicloud_db_instance.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_id": CHECKSET,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"preferred_backup_period": []string{"Wednesday", "Monday"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"preferred_backup_period.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"preferred_backup_time": "10:00Z-11:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"preferred_backup_time": "10:00Z-11:00Z",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"backup_retention_period": "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_retention_period": "10",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"log_backup_frequency": "LogInterval",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"log_backup_frequency": "LogInterval",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"compress_type": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"compress_type": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"preferred_backup_period": []string{"Wednesday", "Tuesday"},
					"preferred_backup_time":   "11:00Z-12:00Z",
					"backup_retention_period": "13",
					"log_backup_frequency":    "LogInterval",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"preferred_backup_period.#": "2",
						"preferred_backup_time":     "11:00Z-12:00Z",
						"backup_retention_period":   "13",
						"log_backup_frequency":      "LogInterval",
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
data "alicloud_db_zones" "default"{
	engine               = "SQLServer"
	engine_version       = "2012"
	instance_charge_type = "PostPaid"
	category = "Basic"
 	db_instance_storage_type = "cloud_essd"
}

data "alicloud_db_instance_classes" "default" {
    zone_id = data.alicloud_db_zones.default.zones.0.id
	engine               = "SQLServer"
	engine_version       = "2012"
    category = "Basic"
 	db_instance_storage_type = "cloud_essd"
	instance_charge_type = "PostPaid"
}

data "alicloud_vpcs" "default" {
    name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_db_zones.default.zones.0.id
}

resource "alicloud_vswitch" "this" {
 count = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
 vswitch_name = var.name
 vpc_id = data.alicloud_vpcs.default.ids.0
 zone_id = data.alicloud_db_zones.default.ids.0
 cidr_block = cidrsubnet(data.alicloud_vpcs.default.vpcs.0.cidr_block, 8, 4)
}
locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids.0 : concat(alicloud_vswitch.this.*.id, [""])[0]
  zone_id = data.alicloud_db_zones.default.ids.0
}

data "alicloud_resource_manager_resource_groups" "default" {
	status = "OK"
}

resource "alicloud_security_group" "default" {
	security_group_name   = var.name
	vpc_id 				  = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_db_instance" "default" {
	engine               = "SQLServer"
	engine_version       = "2012"
 	db_instance_storage_type = "cloud_essd"
	instance_type = data.alicloud_db_instance_classes.default.instance_classes.0.instance_class
	instance_storage = data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min
	vswitch_id = local.vswitch_id
	instance_name = var.name
	security_group_ids = alicloud_security_group.default.*.id
}
`, name)
}

func TestAccAliCloudRdsDBBackupPolicySQLServerAlwaysOn(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_db_backup_policy.default"
	serverFunc := func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serverFunc, "DescribeBackupPolicy")
	ra := resourceAttrInit(resourceId, nil)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := "tf-testAccDBbackuppolicy"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBBackupPolicySQLServerAlwaysOnConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		//CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_id":               "${alicloud_db_instance.default.id}",
					"backup_priority":              "1",
					"enable_increment_data_backup": "false",
					"backup_method":                "Physical",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_id":               CHECKSET,
						"backup_priority":              "1",
						"enable_increment_data_backup": "false",
						"backup_method":                "Physical",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"backup_priority":              "2",
					"enable_increment_data_backup": "true",
					"backup_method":                "Snapshot",
					"log_backup_frequency":         "LogInterval",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_priority":              "2",
						"enable_increment_data_backup": "true",
						"backup_method":                "Snapshot",
						"log_backup_frequency":         "LogInterval",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
func resourceDBBackupPolicySQLServerAlwaysOnConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
data "alicloud_db_zones" "default"{
	engine               = "SQLServer"
	engine_version       = "2022_ent"
	instance_charge_type = "PostPaid"
	category             = "AlwaysOn"
 	db_instance_storage_type = "cloud_essd"
}

data "alicloud_db_instance_classes" "default" {
    zone_id = data.alicloud_db_zones.default.zones.0.id
	engine               = "SQLServer"
	engine_version       = "2022_ent"
    category             = "AlwaysOn"
 	db_instance_storage_type = "cloud_essd"
	instance_charge_type = "PostPaid"
}

data "alicloud_vpcs" "default" {
    name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_db_zones.default.zones.0.id
}

resource "alicloud_vswitch" "this" {
 count = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
 vswitch_name = var.name
 vpc_id = data.alicloud_vpcs.default.ids.0
 zone_id = data.alicloud_db_zones.default.ids.0
 cidr_block = cidrsubnet(data.alicloud_vpcs.default.vpcs.0.cidr_block, 8, 4)
}
locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids.0 : concat(alicloud_vswitch.this.*.id, [""])[0]
  zone_id = data.alicloud_db_zones.default.ids.0
}

data "alicloud_resource_manager_resource_groups" "default" {
	status = "OK"
}

resource "alicloud_security_group" "default" {
	security_group_name   = var.name
	vpc_id 				  = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_db_instance" "default" {
	engine               = "SQLServer"
	engine_version       = "2022_ent"
	category             = "AlwaysOn"
 	db_instance_storage_type = "cloud_essd"
	instance_type = data.alicloud_db_instance_classes.default.instance_classes.0.instance_class
	instance_storage = data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min
	vswitch_id = local.vswitch_id
	instance_name = var.name
	security_group_ids = alicloud_security_group.default.*.id
}
`, name)
}

func TestAccAliCloudRdsDBBackupPolicyMariaDB(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_db_backup_policy.default"
	serverFunc := func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serverFunc, "DescribeBackupPolicy")
	name := "tf-testAccDBbackuppolicy"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBBackupPolicyMariaDBConfigDependence)
	ra := resourceAttrInit(resourceId, nil)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		//CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_id":              "${alicloud_db_instance.default.id}",
					"high_space_usage_protection": "Enable",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_id": CHECKSET,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"preferred_backup_period": []string{"Wednesday", "Monday"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"preferred_backup_period.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"preferred_backup_time": "10:00Z-11:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"preferred_backup_time": "10:00Z-11:00Z",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"backup_retention_period": "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_retention_period": "10",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"log_backup_retention_period": "7",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"log_backup_retention_period": "7",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"local_log_retention_hours": "24",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"local_log_retention_hours": "24",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"local_log_retention_space": "35",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"local_log_retention_space": "35",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"high_space_usage_protection": "Disable",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"high_space_usage_protection": "Disable",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"compress_type": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"compress_type": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_backup_log": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_backup_log": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"preferred_backup_period":     []string{"Wednesday", "Monday", "Tuesday"},
					"preferred_backup_time":       "11:00Z-12:00Z",
					"backup_retention_period":     "20",
					"enable_backup_log":           "true",
					"log_backup_retention_period": "7",
					"local_log_retention_hours":   "48",
					"high_space_usage_protection": "Enable",
					"local_log_retention_space":   "40",
					"compress_type":               "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"preferred_backup_period.#":   "3",
						"preferred_backup_time":       "11:00Z-12:00Z",
						"backup_retention_period":     "20",
						"enable_backup_log":           "true",
						"log_backup_retention_period": "7",
						"local_log_retention_hours":   "48",
						"high_space_usage_protection": "Enable",
						"local_log_retention_space":   "40",
						"compress_type":               "1",
					}),
				),
			}},
	})
}

func resourceDBBackupPolicyMariaDBConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
data "alicloud_db_zones" "default"{
  	engine               = "MariaDB"
  	engine_version       = "10.3"
	instance_charge_type = "PostPaid"
	category = "HighAvailability"
 	db_instance_storage_type = "cloud_essd"
}

data "alicloud_db_instance_classes" "default" {
    zone_id = data.alicloud_db_zones.default.zones.0.id
  	engine               = "MariaDB"
  	engine_version       = "10.3"
    category = "HighAvailability"
 	db_instance_storage_type = "cloud_essd"
	instance_charge_type = "PostPaid"
}

data "alicloud_vpcs" "default" {
    name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_db_zones.default.zones.0.id
}

resource "alicloud_vswitch" "this" {
 count = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
 vswitch_name = var.name
 vpc_id = data.alicloud_vpcs.default.ids.0
 zone_id = data.alicloud_db_zones.default.ids.0
 cidr_block = cidrsubnet(data.alicloud_vpcs.default.vpcs.0.cidr_block, 8, 4)
}
locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids.0 : concat(alicloud_vswitch.this.*.id, [""])[0]
  zone_id = data.alicloud_db_zones.default.ids.0
}

data "alicloud_resource_manager_resource_groups" "default" {
	status = "OK"
}

resource "alicloud_security_group" "default" {
	security_group_name   = var.name
	vpc_id 				  = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_db_instance" "default" {
  	engine               = "MariaDB"
  	engine_version       = "10.3"
 	db_instance_storage_type = "cloud_essd"
	instance_type = data.alicloud_db_instance_classes.default.instance_classes.0.instance_class
	instance_storage = data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min
	vswitch_id = local.vswitch_id
	instance_name = var.name
	security_group_ids = alicloud_security_group.default.*.id
}
`, name)
}
