package rds_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/rds"

	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudRdsDBInstancesDataSource(t *testing.T) {
	rand := tfacctest.RandInt()
	nameConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudDBInstanceDataSourceConfigMysql(rand, map[string]string{
			"name_regex": `"${alicloud_db_instance.default.instance_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudDBInstanceDataSourceConfigMysql(rand, map[string]string{
			"name_regex": `"^test1234"`,
		}),
	}
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudDBInstanceDataSourceConfigMysql(rand, map[string]string{
			"ids": `[ "${alicloud_db_instance.default.id}" ]`,
		}),
		FakeConfig: testAccCheckAlicloudDBInstanceDataSourceConfigMysql(rand, map[string]string{
			"ids": `[ "${alicloud_db_instance.default.id}-fake" ]`,
		}),
	}

	engineConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudDBInstanceDataSourceConfigMysql(rand, map[string]string{
			"name_regex": `"${alicloud_db_instance.default.instance_name}"`,
			"engine":     `"${alicloud_db_instance.default.engine}"`,
		}),
		FakeConfig: testAccCheckAlicloudDBInstanceDataSourceConfigMysql(rand, map[string]string{
			"name_regex": `"${alicloud_db_instance.default.instance_name}"`,
			"engine":     `"SQLServer"`,
		}),
	}

	vpcIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudDBInstanceDataSourceConfigMysql(rand, map[string]string{
			"name_regex": `"${alicloud_db_instance.default.instance_name}"`,
			"vpc_id":     `"${local.vpc_id}"`,
		}),
		FakeConfig: testAccCheckAlicloudDBInstanceDataSourceConfigMysql(rand, map[string]string{
			"name_regex": `"${alicloud_db_instance.default.instance_name}"`,
			"vpc_id":     `"unknow"`,
		}),
	}

	vswitchIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudDBInstanceDataSourceConfigMysql(rand, map[string]string{
			"name_regex": `"${alicloud_db_instance.default.instance_name}"`,
			"vswitch_id": `"${alicloud_db_instance.default.vswitch_id}"`,
		}),
		FakeConfig: testAccCheckAlicloudDBInstanceDataSourceConfigMysql(rand, map[string]string{
			"name_regex": `"${alicloud_db_instance.default.instance_name}"`,
			"vswitch_id": `"unknow"`,
		}),
	}

	tagsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudDBInstanceDataSourceConfigMysql(rand, map[string]string{
			"name_regex": `"${alicloud_db_instance.default.instance_name}"`,
			"tags": `{ 
						"key1" = "value1" 
						"key2" = "value2" 
					}`,
		}),
		FakeConfig: testAccCheckAlicloudDBInstanceDataSourceConfigMysql(rand, map[string]string{
			"name_regex": `"${alicloud_db_instance.default.instance_name}"`,
			"tags": `{ 
						"key1" = "value1_fake" 
						"key2" = "value2_fake" 
					}`,
		}),
	}

	//the parameter connection_mode has not stable default value. It's Standard at cn-hangzhou zone , but at ap-south-1 zone it is Safe.
	//connection_modeConf := acctest.DataSourceTestAccConfig{
	//	ExistConfig: testAccCheckAlicloudDBInstanceDataSourceConfigMysql(rand, map[string]string{
	//		"name_regex":      `"${alicloud_db_instance.default.instance_name}"`,
	//		"connection_mode": `"Standard"`,
	//	}),
	//	FakeConfig: testAccCheckAlicloudDBInstanceDataSourceConfigMysql(rand, map[string]string{
	//		"name_regex":      `"${alicloud_db_instance.default.instance_name}"`,
	//		"connection_mode": `"Safe"`,
	//	}),
	//}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudDBInstanceDataSourceConfigMysql(rand, map[string]string{
			"name_regex": `"${alicloud_db_instance.default.instance_name}"`,
			"vswitch_id": `"${alicloud_db_instance.default.vswitch_id}"`,
			"tags": `{ 
						"key1" = "value1" 
						"key2" = "value2" 
					}`,
			"engine": `"${alicloud_db_instance.default.engine}"`,
			"vpc_id": `"${local.vpc_id}"`,
			"ids":    `[ "${alicloud_db_instance.default.id}" ]`,
		}),
		FakeConfig: testAccCheckAlicloudDBInstanceDataSourceConfigMysql(rand, map[string]string{
			"name_regex": `"${alicloud_db_instance.default.instance_name}"`,
			"vswitch_id": `"${alicloud_db_instance.default.vswitch_id}"`,
			"tags": `{ 
						"key1" = "value1" 
						"key2" = "value2" 
					}`,
			"vpc_id": `"${local.vpc_id}"`,
			"engine": `"SQLServer"`,
			"ids":    `[ "${alicloud_db_instance.default.id}" ]`,
		}),
	}

	var existDBInstanceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                            "1",
			"names.#":                                          "1",
			"instances.#":                                      "1",
			"total_count":                                      acctest.CHECKSET,
			"instances.0.id":                                   acctest.CHECKSET,
			"instances.0.name":                                 fmt.Sprintf("tf-testAccDBInstanceConfig_%d", rand),
			"instances.0.db_type":                              acctest.CHECKSET,
			"instances.0.region_id":                            acctest.CHECKSET,
			"instances.0.create_time":                          acctest.CHECKSET,
			"instances.0.status":                               acctest.CHECKSET,
			"instances.0.engine":                               string(rds.MySQL),
			"instances.0.engine_version":                       "8.0",
			"instances.0.net_type":                             string(rds.Intranet),
			"instances.0.instance_type":                        acctest.CHECKSET,
			"instances.0.connection_mode":                      acctest.CHECKSET,
			"instances.0.availability_zone":                    acctest.CHECKSET,
			"instances.0.vpc_id":                               acctest.CHECKSET,
			"instances.0.vswitch_id":                           acctest.CHECKSET,
			"instances.0.charge_type":                          acctest.CHECKSET,
			"instances.0.connection_string":                    acctest.CHECKSET,
			"instances.0.port":                                 acctest.CHECKSET,
			"instances.0.db_instance_storage_type":             acctest.CHECKSET,
			"instances.0.instance_storage":                     acctest.CHECKSET,
			"instances.0.deletion_protection":                  acctest.CHECKSET,
			"instances.0.ha_mode":                              acctest.CHECKSET,
			"instances.0.sync_mode":                            acctest.CHECKSET,
			"instances.0.host_instance_infos.0.log_sync_time":  "",
			"instances.0.host_instance_infos.0.node_type":      acctest.CHECKSET,
			"instances.0.host_instance_infos.0.zone_id":        acctest.CHECKSET,
			"instances.0.host_instance_infos.0.sync_status":    acctest.CHECKSET,
			"instances.0.host_instance_infos.0.data_sync_time": "",
			"instances.0.host_instance_infos.0.node_id":        acctest.CHECKSET,
			"instances.0.host_instance_infos.0.region_id":      acctest.CHECKSET,
		}
	}

	var fakeDBInstanceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"instances.#": "0",
			"ids.#":       "0",
			"names.#":     "0",
		}
	}

	var DBInstanceCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_db_instances.dbs",
		ExistMapFunc: existDBInstanceMapFunc,
		FakeMapFunc:  fakeDBInstanceMapFunc,
	}

	DBInstanceCheckInfo.DataSourceTestCheck(t, rand, nameConf, idsConf, engineConf, vpcIdConf, vswitchIdConf, tagsConf, allConf)
}

func TestAccAliCloudRdsDBInstancesDataSourcePostgreSQLSSL(t *testing.T) {
	rand := tfacctest.RandInt()
	nameConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudDBInstanceDataSourceConfigPostgreSQL(rand, map[string]string{
			"name_regex": `"${alicloud_db_instance.default.instance_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudDBInstanceDataSourceConfigPostgreSQL(rand, map[string]string{
			"name_regex": `"^test1234"`,
		}),
	}
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudDBInstanceDataSourceConfigPostgreSQL(rand, map[string]string{
			"ids": `[ "${alicloud_db_instance.default.id}" ]`,
		}),
		FakeConfig: testAccCheckAlicloudDBInstanceDataSourceConfigPostgreSQL(rand, map[string]string{
			"ids": `[ "${alicloud_db_instance.default.id}-fake" ]`,
		}),
	}

	engineConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudDBInstanceDataSourceConfigPostgreSQL(rand, map[string]string{
			"name_regex": `"${alicloud_db_instance.default.instance_name}"`,
			"engine":     `"${alicloud_db_instance.default.engine}"`,
		}),
		FakeConfig: testAccCheckAlicloudDBInstanceDataSourceConfigPostgreSQL(rand, map[string]string{
			"name_regex": `"${alicloud_db_instance.default.instance_name}"`,
			"engine":     `"SQLServer"`,
		}),
	}

	vpcIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudDBInstanceDataSourceConfigPostgreSQL(rand, map[string]string{
			"name_regex": `"${alicloud_db_instance.default.instance_name}"`,
			"vpc_id":     `"${local.vpc_id}"`,
		}),
		FakeConfig: testAccCheckAlicloudDBInstanceDataSourceConfigPostgreSQL(rand, map[string]string{
			"name_regex": `"${alicloud_db_instance.default.instance_name}"`,
			"vpc_id":     `"unknow"`,
		}),
	}

	vswitchIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudDBInstanceDataSourceConfigPostgreSQL(rand, map[string]string{
			"name_regex": `"${alicloud_db_instance.default.instance_name}"`,
			"vswitch_id": `"${alicloud_db_instance.default.vswitch_id}"`,
		}),
		FakeConfig: testAccCheckAlicloudDBInstanceDataSourceConfigPostgreSQL(rand, map[string]string{
			"name_regex": `"${alicloud_db_instance.default.instance_name}"`,
			"vswitch_id": `"unknow"`,
		}),
	}

	tagsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudDBInstanceDataSourceConfigPostgreSQL(rand, map[string]string{
			"name_regex": `"${alicloud_db_instance.default.instance_name}"`,
			"tags": `{ 
						"key1" = "value1" 
						"key2" = "value2" 
					}`,
		}),
		FakeConfig: testAccCheckAlicloudDBInstanceDataSourceConfigPostgreSQL(rand, map[string]string{
			"name_regex": `"${alicloud_db_instance.default.instance_name}"`,
			"tags": `{ 
						"key1" = "value1_fake" 
						"key2" = "value2_fake" 
					}`,
		}),
	}

	//the parameter connection_mode has not stable default value. It's Standard at cn-hangzhou zone , but at ap-south-1 zone it is Safe.
	//connection_modeConf := acctest.DataSourceTestAccConfig{
	//	ExistConfig: testAccCheckAlicloudDBInstanceDataSourceConfigMysql(rand, map[string]string{
	//		"name_regex":      `"${alicloud_db_instance.default.instance_name}"`,
	//		"connection_mode": `"Standard"`,
	//	}),
	//	FakeConfig: testAccCheckAlicloudDBInstanceDataSourceConfigMysql(rand, map[string]string{
	//		"name_regex":      `"${alicloud_db_instance.default.instance_name}"`,
	//		"connection_mode": `"Safe"`,
	//	}),
	//}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudDBInstanceDataSourceConfigPostgreSQL(rand, map[string]string{
			"name_regex": `"${alicloud_db_instance.default.instance_name}"`,
			"vswitch_id": `"${alicloud_db_instance.default.vswitch_id}"`,
			"tags": `{ 
						"key1" = "value1" 
						"key2" = "value2" 
					}`,
			"engine": `"${alicloud_db_instance.default.engine}"`,
			"vpc_id": `"${local.vpc_id}"`,
			"ids":    `[ "${alicloud_db_instance.default.id}" ]`,
		}),
		FakeConfig: testAccCheckAlicloudDBInstanceDataSourceConfigPostgreSQL(rand, map[string]string{
			"name_regex": `"${alicloud_db_instance.default.instance_name}"`,
			"vswitch_id": `"${alicloud_db_instance.default.vswitch_id}"`,
			"tags": `{ 
						"key1" = "value1" 
						"key2" = "value2" 
					}`,
			"vpc_id": `"${local.vpc_id}"`,
			"engine": `"SQLServer"`,
			"ids":    `[ "${alicloud_db_instance.default.id}" ]`,
		}),
	}

	var existDBInstanceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                   "1",
			"names.#":                                 "1",
			"instances.#":                             "1",
			"total_count":                             acctest.CHECKSET,
			"instances.0.id":                          acctest.CHECKSET,
			"instances.0.name":                        fmt.Sprintf("tf-testAccDBInstanceConfig_%d", rand),
			"instances.0.db_type":                     acctest.CHECKSET,
			"instances.0.region_id":                   acctest.CHECKSET,
			"instances.0.create_time":                 acctest.CHECKSET,
			"instances.0.status":                      acctest.CHECKSET,
			"instances.0.engine":                      string(rds.PostgreSQL),
			"instances.0.engine_version":              "13.0",
			"instances.0.net_type":                    string(rds.Intranet),
			"instances.0.instance_type":               acctest.CHECKSET,
			"instances.0.connection_mode":             acctest.CHECKSET,
			"instances.0.availability_zone":           acctest.CHECKSET,
			"instances.0.vpc_id":                      acctest.CHECKSET,
			"instances.0.vswitch_id":                  acctest.CHECKSET,
			"instances.0.charge_type":                 acctest.CHECKSET,
			"instances.0.connection_string":           acctest.CHECKSET,
			"instances.0.port":                        acctest.CHECKSET,
			"instances.0.db_instance_storage_type":    acctest.CHECKSET,
			"instances.0.instance_storage":            acctest.CHECKSET,
			"instances.0.ssl_expire_time":             "",
			"instances.0.require_update":              "",
			"instances.0.acl":                         "",
			"instances.0.ca_type":                     "",
			"instances.0.client_ca_cert":              "",
			"instances.0.client_ca_cert_expire_time":  "",
			"instances.0.client_cert_revocation_list": "",
			"instances.0.last_modify_status":          "",
			"instances.0.modify_status_reason":        "",
			"instances.0.replication_acl":             "",
			"instances.0.require_update_item":         "",
			"instances.0.require_update_reason":       "",
			"instances.0.ssl_create_time":             "",
			"instances.0.ssl_enabled":                 "off",
			"instances.0.server_ca_url":               "",
			"instances.0.server_cert":                 "",
			"instances.0.server_key":                  "",
			"instances.0.deletion_protection":         acctest.CHECKSET,
		}
	}

	var fakeDBInstanceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"instances.#": "0",
			"ids.#":       "0",
			"names.#":     "0",
		}
	}

	var DBInstanceCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_db_instances.dbs",
		ExistMapFunc: existDBInstanceMapFunc,
		FakeMapFunc:  fakeDBInstanceMapFunc,
	}

	DBInstanceCheckInfo.DataSourceTestCheck(t, rand, nameConf, idsConf, engineConf, vpcIdConf, vswitchIdConf, tagsConf, allConf)
}

func testAccCheckAlicloudDBInstanceDataSourceConfigMysql(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccDBInstanceConfig_%d"
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
  vpc_id = data.alicloud_vpcs.default.ids.0
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids.0 : concat(alicloud_vswitch.this.*.id, [""])[0]
  zone_id = data.alicloud_db_zones.default.ids[length(data.alicloud_db_zones.default.ids)-1]
}

data "alicloud_resource_manager_resource_groups" "default" {
	status = "OK"
}

resource "alicloud_security_group" "default" {
	name   = var.name
	vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_db_instance" "default" {
    engine = "MySQL"
	engine_version = "8.0"
 	db_instance_storage_type = "cloud_essd"
	instance_type = data.alicloud_db_instance_classes.default.instance_classes.0.instance_class
	instance_storage = data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.0.min
	vswitch_id = local.vswitch_id
	instance_name = var.name
	tags = {
		"key1" = "value1"
		"key2" = "value2"
	}
}
data "alicloud_db_instances" "dbs" {
  %s
}
`, rand, strings.Join(pairs, "\n  "))
	return config
}

func testAccCheckAlicloudDBInstanceDataSourceConfigPostgreSQL(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccDBInstanceConfig_%d"
}
data "alicloud_db_zones" "default"{
	engine = "PostgreSQL"
	engine_version = "13.0"
	instance_charge_type = "PostPaid"
	category = "HighAvailability"
 	db_instance_storage_type = "cloud_essd"
}

data "alicloud_db_instance_classes" "default" {
    zone_id = data.alicloud_db_zones.default.zones.0.id
	engine = "PostgreSQL"
	engine_version = "13.0"
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
  vpc_id = data.alicloud_vpcs.default.ids.0
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids.0 : concat(alicloud_vswitch.this.*.id, [""])[0]
  zone_id = data.alicloud_db_zones.default.ids[length(data.alicloud_db_zones.default.ids)-1]
}

data "alicloud_resource_manager_resource_groups" "default" {
	status = "OK"
}

resource "alicloud_security_group" "default" {
	name   = var.name
	vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_db_instance" "default" {
	engine = "PostgreSQL"
	engine_version = "13.0"
 	db_instance_storage_type = "cloud_essd"
	instance_type = data.alicloud_db_instance_classes.default.instance_classes.0.instance_class
	instance_storage = data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.0.min
	vswitch_id = local.vswitch_id
	instance_name = var.name
	tags = {
		"key1" = "value1"
		"key2" = "value2"
	}
}
data "alicloud_db_instances" "dbs" {
  %s
}
`, rand, strings.Join(pairs, "\n  "))
	return config
}
