package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudRdsDBInstancesDataSource(t *testing.T) {
	rand := acctest.RandInt()
	nameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDBInstanceDataSourceConfig_mysql(rand, map[string]string{
			"name_regex": `"${alicloud_db_instance.default.instance_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudDBInstanceDataSourceConfig_mysql(rand, map[string]string{
			"name_regex": `"^test1234"`,
		}),
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDBInstanceDataSourceConfig_mysql(rand, map[string]string{
			"ids": `[ "${alicloud_db_instance.default.id}" ]`,
		}),
		fakeConfig: testAccCheckAlicloudDBInstanceDataSourceConfig_mysql(rand, map[string]string{
			"ids": `[ "${alicloud_db_instance.default.id}-fake" ]`,
		}),
	}

	engineConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDBInstanceDataSourceConfig_mysql(rand, map[string]string{
			"name_regex": `"${alicloud_db_instance.default.instance_name}"`,
			"engine":     `"${alicloud_db_instance.default.engine}"`,
		}),
		fakeConfig: testAccCheckAlicloudDBInstanceDataSourceConfig_mysql(rand, map[string]string{
			"name_regex": `"${alicloud_db_instance.default.instance_name}"`,
			"engine":     `"SQLServer"`,
		}),
	}

	vpc_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDBInstanceDataSourceConfig_mysql(rand, map[string]string{
			"name_regex": `"${alicloud_db_instance.default.instance_name}"`,
			"vpc_id":     `"${alicloud_vswitch.default.vpc_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudDBInstanceDataSourceConfig_mysql(rand, map[string]string{
			"name_regex": `"${alicloud_db_instance.default.instance_name}"`,
			"vpc_id":     `"unknow"`,
		}),
	}

	vswitch_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDBInstanceDataSourceConfig_mysql(rand, map[string]string{
			"name_regex": `"${alicloud_db_instance.default.instance_name}"`,
			"vswitch_id": `"${alicloud_db_instance.default.vswitch_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudDBInstanceDataSourceConfig_mysql(rand, map[string]string{
			"name_regex": `"${alicloud_db_instance.default.instance_name}"`,
			"vswitch_id": `"unknow"`,
		}),
	}

	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDBInstanceDataSourceConfig_mysql(rand, map[string]string{
			"name_regex": `"${alicloud_db_instance.default.instance_name}"`,
			"tags": `{ 
						"key1" = "value1" 
						"key2" = "value2" 
					}`,
		}),
		fakeConfig: testAccCheckAlicloudDBInstanceDataSourceConfig_mysql(rand, map[string]string{
			"name_regex": `"${alicloud_db_instance.default.instance_name}"`,
			"tags": `{ 
						"key1" = "value1_fake" 
						"key2" = "value2_fake" 
					}`,
		}),
	}

	//the parameter connection_mode has not stable default value. It's Standard at cn-hangzhou zone , but at ap-south-1 zone it is Safe.
	//connection_modeConf := dataSourceTestAccConfig{
	//	existConfig: testAccCheckAlicloudDBInstanceDataSourceConfig_mysql(rand, map[string]string{
	//		"name_regex":      `"${alicloud_db_instance.default.instance_name}"`,
	//		"connection_mode": `"Standard"`,
	//	}),
	//	fakeConfig: testAccCheckAlicloudDBInstanceDataSourceConfig_mysql(rand, map[string]string{
	//		"name_regex":      `"${alicloud_db_instance.default.instance_name}"`,
	//		"connection_mode": `"Safe"`,
	//	}),
	//}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDBInstanceDataSourceConfig_mysql(rand, map[string]string{
			"name_regex": `"${alicloud_db_instance.default.instance_name}"`,
			"vswitch_id": `"${alicloud_db_instance.default.vswitch_id}"`,
			"tags": `{ 
						"key1" = "value1" 
						"key2" = "value2" 
					}`,
			"engine": `"${alicloud_db_instance.default.engine}"`,
			"vpc_id": `"${alicloud_vswitch.default.vpc_id}"`,
			"ids":    `[ "${alicloud_db_instance.default.id}" ]`,
		}),
		fakeConfig: testAccCheckAlicloudDBInstanceDataSourceConfig_mysql(rand, map[string]string{
			"name_regex": `"${alicloud_db_instance.default.instance_name}"`,
			"vswitch_id": `"${alicloud_db_instance.default.vswitch_id}"`,
			"tags": `{ 
						"key1" = "value1" 
						"key2" = "value2" 
					}`,
			"vpc_id": `"${alicloud_vswitch.default.vpc_id}"`,
			"engine": `"SQLServer"`,
			"ids":    `[ "${alicloud_db_instance.default.id}" ]`,
		}),
	}

	var existDBInstanceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                "1",
			"names.#":                              "1",
			"instances.#":                          "1",
			"instances.0.id":                       CHECKSET,
			"instances.0.name":                     fmt.Sprintf("tf-testAccDBInstanceConfig_%d", rand),
			"instances.0.db_type":                  CHECKSET,
			"instances.0.region_id":                CHECKSET,
			"instances.0.create_time":              CHECKSET,
			"instances.0.status":                   CHECKSET,
			"instances.0.engine":                   string(MySQL),
			"instances.0.engine_version":           "5.6",
			"instances.0.net_type":                 string(Intranet),
			"instances.0.instance_type":            "rds.mysql.s1.small",
			"instances.0.connection_mode":          CHECKSET,
			"instances.0.availability_zone":        CHECKSET,
			"instances.0.vpc_id":                   CHECKSET,
			"instances.0.vswitch_id":               CHECKSET,
			"instances.0.charge_type":              CHECKSET,
			"instances.0.connection_string":        CHECKSET,
			"instances.0.port":                     CHECKSET,
			"instances.0.db_instance_storage_type": CHECKSET,
			"instances.0.instance_storage":         CHECKSET,
		}
	}

	var fakeDBInstanceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"instances.#": "0",
			"ids.#":       "0",
			"names.#":     "0",
		}
	}

	var DBInstanceCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_db_instances.dbs",
		existMapFunc: existDBInstanceMapFunc,
		fakeMapFunc:  fakeDBInstanceMapFunc,
	}

	DBInstanceCheckInfo.dataSourceTestCheck(t, rand, nameConf, idsConf, engineConf, vpc_idConf, vswitch_idConf, tagsConf, allConf)
}

func TestAccAlicloudRdsDBInstancesDataSourcePostgreSQLSSL(t *testing.T) {
	rand := acctest.RandInt()
	nameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDBInstanceDataSourceConfig_PostgreSQL(rand, map[string]string{
			"name_regex": `"${alicloud_db_instance.default.instance_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudDBInstanceDataSourceConfig_PostgreSQL(rand, map[string]string{
			"name_regex": `"^test1234"`,
		}),
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDBInstanceDataSourceConfig_PostgreSQL(rand, map[string]string{
			"ids": `[ "${alicloud_db_instance.default.id}" ]`,
		}),
		fakeConfig: testAccCheckAlicloudDBInstanceDataSourceConfig_PostgreSQL(rand, map[string]string{
			"ids": `[ "${alicloud_db_instance.default.id}-fake" ]`,
		}),
	}

	engineConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDBInstanceDataSourceConfig_PostgreSQL(rand, map[string]string{
			"name_regex": `"${alicloud_db_instance.default.instance_name}"`,
			"engine":     `"${alicloud_db_instance.default.engine}"`,
		}),
		fakeConfig: testAccCheckAlicloudDBInstanceDataSourceConfig_PostgreSQL(rand, map[string]string{
			"name_regex": `"${alicloud_db_instance.default.instance_name}"`,
			"engine":     `"SQLServer"`,
		}),
	}

	vpc_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDBInstanceDataSourceConfig_PostgreSQL(rand, map[string]string{
			"name_regex": `"${alicloud_db_instance.default.instance_name}"`,
			"vpc_id":     `"${alicloud_vswitch.default.vpc_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudDBInstanceDataSourceConfig_PostgreSQL(rand, map[string]string{
			"name_regex": `"${alicloud_db_instance.default.instance_name}"`,
			"vpc_id":     `"unknow"`,
		}),
	}

	vswitch_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDBInstanceDataSourceConfig_PostgreSQL(rand, map[string]string{
			"name_regex": `"${alicloud_db_instance.default.instance_name}"`,
			"vswitch_id": `"${alicloud_db_instance.default.vswitch_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudDBInstanceDataSourceConfig_PostgreSQL(rand, map[string]string{
			"name_regex": `"${alicloud_db_instance.default.instance_name}"`,
			"vswitch_id": `"unknow"`,
		}),
	}

	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDBInstanceDataSourceConfig_PostgreSQL(rand, map[string]string{
			"name_regex": `"${alicloud_db_instance.default.instance_name}"`,
			"tags": `{ 
						"key1" = "value1" 
						"key2" = "value2" 
					}`,
		}),
		fakeConfig: testAccCheckAlicloudDBInstanceDataSourceConfig_PostgreSQL(rand, map[string]string{
			"name_regex": `"${alicloud_db_instance.default.instance_name}"`,
			"tags": `{ 
						"key1" = "value1_fake" 
						"key2" = "value2_fake" 
					}`,
		}),
	}

	//the parameter connection_mode has not stable default value. It's Standard at cn-hangzhou zone , but at ap-south-1 zone it is Safe.
	//connection_modeConf := dataSourceTestAccConfig{
	//	existConfig: testAccCheckAlicloudDBInstanceDataSourceConfig_mysql(rand, map[string]string{
	//		"name_regex":      `"${alicloud_db_instance.default.instance_name}"`,
	//		"connection_mode": `"Standard"`,
	//	}),
	//	fakeConfig: testAccCheckAlicloudDBInstanceDataSourceConfig_mysql(rand, map[string]string{
	//		"name_regex":      `"${alicloud_db_instance.default.instance_name}"`,
	//		"connection_mode": `"Safe"`,
	//	}),
	//}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDBInstanceDataSourceConfig_PostgreSQL(rand, map[string]string{
			"name_regex": `"${alicloud_db_instance.default.instance_name}"`,
			"vswitch_id": `"${alicloud_db_instance.default.vswitch_id}"`,
			"tags": `{ 
						"key1" = "value1" 
						"key2" = "value2" 
					}`,
			"engine": `"${alicloud_db_instance.default.engine}"`,
			"vpc_id": `"${alicloud_vswitch.default.vpc_id}"`,
			"ids":    `[ "${alicloud_db_instance.default.id}" ]`,
		}),
		fakeConfig: testAccCheckAlicloudDBInstanceDataSourceConfig_PostgreSQL(rand, map[string]string{
			"name_regex": `"${alicloud_db_instance.default.instance_name}"`,
			"vswitch_id": `"${alicloud_db_instance.default.vswitch_id}"`,
			"tags": `{ 
						"key1" = "value1" 
						"key2" = "value2" 
					}`,
			"vpc_id": `"${alicloud_vswitch.default.vpc_id}"`,
			"engine": `"SQLServer"`,
			"ids":    `[ "${alicloud_db_instance.default.id}" ]`,
		}),
	}

	var existDBInstanceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                   "1",
			"names.#":                                 "1",
			"instances.#":                             "1",
			"instances.0.id":                          CHECKSET,
			"instances.0.name":                        fmt.Sprintf("tf-testAccDBInstanceConfig_%d", rand),
			"instances.0.db_type":                     CHECKSET,
			"instances.0.region_id":                   CHECKSET,
			"instances.0.create_time":                 CHECKSET,
			"instances.0.status":                      CHECKSET,
			"instances.0.engine":                      string(PostgreSQL),
			"instances.0.engine_version":              "13.0",
			"instances.0.net_type":                    string(Intranet),
			"instances.0.instance_type":               "pg.n2.small.1",
			"instances.0.connection_mode":             CHECKSET,
			"instances.0.availability_zone":           CHECKSET,
			"instances.0.vpc_id":                      CHECKSET,
			"instances.0.vswitch_id":                  CHECKSET,
			"instances.0.charge_type":                 CHECKSET,
			"instances.0.connection_string":           CHECKSET,
			"instances.0.port":                        CHECKSET,
			"instances.0.db_instance_storage_type":    CHECKSET,
			"instances.0.instance_storage":            CHECKSET,
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
		}
	}

	var fakeDBInstanceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"instances.#": "0",
			"ids.#":       "0",
			"names.#":     "0",
		}
	}

	var DBInstanceCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_db_instances.dbs",
		existMapFunc: existDBInstanceMapFunc,
		fakeMapFunc:  fakeDBInstanceMapFunc,
	}

	DBInstanceCheckInfo.dataSourceTestCheck(t, rand, nameConf, idsConf, engineConf, vpc_idConf, vswitch_idConf, tagsConf, allConf)
}

func testAccCheckAlicloudDBInstanceDataSourceConfig_mysql(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`

data "alicloud_db_zones" "default"{
	engine = "MySQL"
	engine_version = "8.0"
	instance_charge_type = "PostPaid"
}

data "alicloud_db_instance_classes" "default" {
    zone_id = data.alicloud_db_zones.default.zones.0.id
	engine = "MySQL"
	engine_version = "8.0"
    category = "HighAvailability"
    storage_type = "local_ssd"
	instance_charge_type = "PostPaid"
}

data "alicloud_vpcs" "default" {
 name_regex = "^default-NODELETING"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_db_zones.default.zones.0.id
}

variable "name" {
	default = "tf-testAccDBInstanceConfig_%d"
}
resource "alicloud_db_instance" "default" {
	engine = "MySQL"
	engine_version = "8.0"
	instance_type = data.alicloud_db_instance_classes.default.ids.0
	instance_storage = "20"
	vswitch_id = data.alicloud_vswitches.default.ids.0
	instance_name = "${var.name}"
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

func testAccCheckAlicloudDBInstanceDataSourceConfig_PostgreSQL(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
data "alicloud_db_zones" "default"{
	engine = "PostgreSQL"
	engine_version = "13.0"
	instance_charge_type = "PostPaid"
}

data "alicloud_db_instance_classes" "default" {
    zone_id = data.alicloud_db_zones.default.zones.0.id
	engine = "PostgreSQL"
	engine_version = "13.0"
    category = "HighAvailability"
    storage_type = "local_ssd"
	instance_charge_type = "PostPaid"
}

data "alicloud_vpcs" "default" {
 name_regex = "^default-NODELETING"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_db_zones.default.zones.0.id
}

variable "name" {
	default = "tf-testAccDBInstanceConfig_%d"
}
resource "alicloud_db_instance" "default" {
	engine = "PostgreSQL"
	engine_version = "13.0"
	instance_type = data.alicloud_db_instance_classes.default.ids.0
	instance_storage = "20"
	vswitch_id = data.alicloud_vswitches.default.ids.0
	instance_name = "${var.name}"
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
