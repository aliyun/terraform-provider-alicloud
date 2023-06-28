package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudRdsInstanceCrossBackupPolicyMySql(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rds_instance_cross_backup_policy.default"
	serverFunc := func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serverFunc, "DescribeInstanceCrossBackupPolicy")
	ra := resourceAttrInit(resourceId, nil)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := "tf-testAccRdsCrossBackupPolicy"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceRdsInstanceCrossBackupPolicyMysqlConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.RDSInstanceClassesSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":         "${alicloud_db_instance.default.id}",
					"cross_backup_region": "${data.alicloud_rds_cross_regions.regions.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":         CHECKSET,
						"cross_backup_region": CHECKSET,
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
					"cross_backup_region": "${data.alicloud_rds_cross_regions.regions.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cross_backup_region": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"log_backup_enabled": "Disabled",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"log_backup_enabled": "Disabled",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"retention": "15",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"retention": "15",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":         "${alicloud_db_instance.default.id}",
					"cross_backup_region": "${data.alicloud_rds_cross_regions.regions.ids.2}",
					"log_backup_enabled":  "Enable",
					"retention":           "30",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":         CHECKSET,
						"cross_backup_region": CHECKSET,
						"log_backup_enabled":  "Enable",
						"retention":           "30",
					}),
				),
			}},
	})
}

func resourceRdsInstanceCrossBackupPolicyMysqlConfigDependence(name string) string {
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

data "alicloud_rds_cross_regions" "regions" {
}

data "alicloud_vpcs" "default" {
 name_regex = "^default-NODELETING"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_db_zones.default.ids.0
}

data "alicloud_resource_manager_resource_groups" "default" {
	status = "OK"
}

resource "alicloud_db_instance" "default" {
    engine = "MySQL"
	engine_version = "8.0"
 	db_instance_storage_type = "local_ssd"
	instance_type = data.alicloud_db_instance_classes.default.instance_classes.0.instance_class
	instance_storage = data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min
	vswitch_id = data.alicloud_vswitches.default.ids.0
	instance_name = var.name
}
`, name)
}

func TestAccAlicloudRdsInstanceCrossBackupPolicyPostgreSQL(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rds_instance_cross_backup_policy.default"
	serverFunc := func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serverFunc, "DescribeInstanceCrossBackupPolicy")
	ra := resourceAttrInit(resourceId, nil)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := "tf-testAccRdsCrossBackupPolicy"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceRdsInstanceCrossBackupPolicyPostgreSQLConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.RDSInstanceClassesSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":         "${alicloud_db_instance.default.id}",
					"cross_backup_region": "${data.alicloud_rds_cross_regions.regions.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":         CHECKSET,
						"cross_backup_region": CHECKSET,
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
					"cross_backup_region": "${data.alicloud_rds_cross_regions.regions.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cross_backup_region": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"log_backup_enabled": "Disabled",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"log_backup_enabled": "Disabled",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"retention": "15",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"retention": "15",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":         "${alicloud_db_instance.default.id}",
					"cross_backup_region": "${data.alicloud_rds_cross_regions.regions.ids.2}",
					"log_backup_enabled":  "Enable",
					"retention":           "30",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":         CHECKSET,
						"cross_backup_region": CHECKSET,
						"log_backup_enabled":  "Enable",
						"retention":           "30",
					}),
				),
			}},
	})
}

func resourceRdsInstanceCrossBackupPolicyPostgreSQLConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
data "alicloud_db_zones" "default"{
	engine = "PostgreSQL"
	engine_version = "10.0"
	instance_charge_type = "PostPaid"
	category = "HighAvailability"
 	db_instance_storage_type = "local_ssd"
}

data "alicloud_db_instance_classes" "default" {
    zone_id = data.alicloud_db_zones.default.zones.0.id
	engine = "PostgreSQL"
	engine_version = "10.0"
    category = "HighAvailability"
 	db_instance_storage_type = "local_ssd"
	instance_charge_type = "PostPaid"
}

data "alicloud_rds_cross_regions" "regions" {
}

data "alicloud_vpcs" "default" {
 name_regex = "^default-NODELETING"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_db_zones.default.zones.0.id
}

data "alicloud_resource_manager_resource_groups" "default" {
	status = "OK"
}

resource "alicloud_db_instance" "default" {
	engine = "PostgreSQL"
	engine_version = "10.0"
 	db_instance_storage_type = "local_ssd"
	instance_type = data.alicloud_db_instance_classes.default.instance_classes.0.instance_class
	instance_storage = data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min
	vswitch_id = data.alicloud_vswitches.default.ids.0
	instance_name = var.name
}
`, name)
}
