package alicloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudDMSEnterprisesDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	resourceId := "data.alicloud_dms_enterprise_instances.default"
	name := fmt.Sprintf("tf_testAccDmsEnterpriseInstancesDataSource_%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		name, dataSourceDmsEnterpriseInstancesConfigDependence)

	searchkeyConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"search_key": "${alicloud_dms_enterprise_instance.default.host}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"search_key": "${alicloud_dms_enterprise_instance.default.host}-fake",
		}),
	}
	instancealiasRegexConfConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"net_type":             "VPC",
			"instance_type":        "${alicloud_dms_enterprise_instance.default.instance_type}",
			"env_type":             "test",
			"instance_alias_regex": name,
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"net_type":             "VPC",
			"instance_type":        "${alicloud_dms_enterprise_instance.default.instance_type}",
			"env_type":             "test",
			"instance_alias_regex": name + "fake",
		}),
	}
	nameRegexConfConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"net_type":      "VPC",
			"instance_type": "${alicloud_dms_enterprise_instance.default.instance_type}",
			"env_type":      "test",
			"name_regex":    name,
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"net_type":      "VPC",
			"instance_type": "${alicloud_dms_enterprise_instance.default.instance_type}",
			"env_type":      "test",
			"name_regex":    name + "fake",
		}),
	}
	var existDmsEnterpriseInstancesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"instances.#":                   "1",
			"instances.0.data_link_name":    "",
			"instances.0.database_password": CHECKSET,
			"instances.0.database_user":     "tftestnormal",
			"instances.0.dba_id":            CHECKSET,
			"instances.0.dba_nick_name":     CHECKSET,
			"instances.0.ddl_online":        "0",
			"instances.0.ecs_instance_id":   "",
			"instances.0.ecs_region":        os.Getenv("ALICLOUD_REGION"),
			"instances.0.env_type":          "test",
			"instances.0.export_timeout":    CHECKSET,
			"instances.0.host":              CHECKSET,
			"instances.0.instance_alias":    CHECKSET,
			"instances.0.instance_id":       CHECKSET,
			"instances.0.instance_source":   "RDS",
			"instances.0.instance_type":     "mysql",
			"instances.0.port":              "3306",
			"instances.0.query_timeout":     CHECKSET,
			"instances.0.safe_rule_id":      CHECKSET,
			"instances.0.sid":               "",
			"instances.0.status":            CHECKSET,
			"instances.0.use_dsql":          "0",
			"instances.0.vpc_id":            "",
		}
	}

	var fakeDmsEnterpriseInstancesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"instances.#": "0",
		}
	}

	var DmsEnterpriseInstancesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existDmsEnterpriseInstancesMapFunc,
		fakeMapFunc:  fakeDmsEnterpriseInstancesMapFunc,
	}

	DmsEnterpriseInstancesCheckInfo.dataSourceTestCheck(t, rand, searchkeyConf, instancealiasRegexConfConf, nameRegexConfConf)
}

func dataSourceDmsEnterpriseInstancesConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}
	data "alicloud_account" "current" {
	}
	data "alicloud_db_zones" "default"{
		engine = "MySQL"
		engine_version = "8.0"
		instance_charge_type = "PostPaid"
		category = "HighAvailability"
		db_instance_storage_type = "cloud_essd"
	}
	
	data "alicloud_db_instance_classes" "default" {
		zone_id = data.alicloud_db_zones.default.zones.0.id
		engine = "MySQL"
		engine_version = "8.0"
		category = "HighAvailability"
		db_instance_storage_type = "cloud_essd"
		instance_charge_type = "PostPaid"
	}
	
	data "alicloud_vpcs" "default" {
	 name_regex = "^default-NODELETING"
	}
	data "alicloud_vswitches" "default" {
	  vpc_id = data.alicloud_vpcs.default.ids.0
	  zone_id = data.alicloud_db_zones.default.zones.0.id
	}
	
	resource "alicloud_security_group" "default" {
		name = var.name
		vpc_id = data.alicloud_vpcs.default.ids.0
	}
	
	resource "alicloud_db_instance" "instance" {
		engine = "MySQL"
		engine_version = "8.0"
		db_instance_storage_type = "cloud_essd"
		instance_type = data.alicloud_db_instance_classes.default.instance_classes.0.instance_class
		instance_storage = data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min
		vswitch_id       = data.alicloud_vswitches.default.ids.0
		instance_name    = var.name
		security_ips     = ["100.104.5.0/24","192.168.0.6"]
		tags = {
			"key1" = "value1"
			"key2" = "value2"
		}
	}
	
	resource "alicloud_db_account" "account" {
		instance_id = "${alicloud_db_instance.instance.id}"
		name        = "tftestnormal"
		password    = "Test12345"
		type        = "Normal"
	}

	data "alicloud_dms_user_tenants" "default" {
		status = "ACTIVE"
	}
	resource "alicloud_dms_enterprise_instance" "default" {
	  dba_uid           =  tonumber(data.alicloud_account.current.id)
	  host              =  "${alicloud_db_instance.instance.connection_string}"
	  port              =  "3306"
	  network_type      =	 "VPC"
	  safe_rule         =	"自由操作"
	  tid               =  data.alicloud_dms_user_tenants.default.ids.0
	  instance_type     =	 "mysql"
	  instance_source   =	 "RDS"
	  env_type          =	 "test"
	  database_user     =	 alicloud_db_account.account.name
	  database_password =	 alicloud_db_account.account.password
	  instance_alias    =	 var.name
	  query_timeout     =	 "70"
	  export_timeout    =	 "2000"
	  ecs_region        =	 "%s"
	  ddl_online        =	 "0"
	  use_dsql          =	 "0"
	  data_link_name    =	 ""
	}
`, name, defaultRegionToTest)
}
