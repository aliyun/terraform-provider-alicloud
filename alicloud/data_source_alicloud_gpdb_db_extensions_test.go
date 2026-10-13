// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccAlicloudGpdbDbExtensionDataSource(t *testing.T) {
	testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
	rand := acctest.RandIntRange(1000000, 9999999)

	// db_instance_id and database_name are required request parameters of
	// ListDatabaseExtensions rather than client-side filters, so they stay real in
	// every case; ids is the filter that is exercised against a non-matching value.
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudGpdbDbExtensionSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_gpdb_db_extension.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudGpdbDbExtensionSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_gpdb_db_extension.default.id}_fake"]`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudGpdbDbExtensionSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_gpdb_db_extension.default.id}"]`,
			"enable_details": `"true"`,
		}),
		fakeConfig: testAccCheckAlicloudGpdbDbExtensionSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_gpdb_db_extension.default.id}_fake"]`,
			"enable_details": `"true"`,
		}),
	}

	GpdbDbExtensionCheckInfo.dataSourceTestCheck(t, rand, idsConf, allConf)
}

var existGpdbDbExtensionMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                       "1",
		"extensions.#":                "1",
		"extensions.0.id":             CHECKSET,
		"extensions.0.status":         CHECKSET,
		"extensions.0.extension_name": CHECKSET,
		"extensions.0.description":    CHECKSET,
	}
}

var fakeGpdbDbExtensionMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":        "0",
		"extensions.#": "0",
	}
}

var GpdbDbExtensionCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_gpdb_db_extensions.default",
	existMapFunc: existGpdbDbExtensionMapFunc,
	fakeMapFunc:  fakeGpdbDbExtensionMapFunc,
}

func testAccCheckAlicloudGpdbDbExtensionSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccGpdbDbExtension%d"
}
resource "alicloud_vpc" "defaultiYyNGW" {
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "defaultPlruct" {
  vpc_id     = alicloud_vpc.defaultiYyNGW.id
  zone_id    = "cn-beijing-h"
  cidr_block = "192.168.1.0/24"
}

resource "alicloud_gpdb_instance" "defaultqsmpIy" {
  instance_spec         = "2C8G"
  seg_node_num          = "2"
  seg_storage_type      = "cloud_essd"
  instance_network_type = "VPC"
  db_instance_category  = "Basic"
  payment_type          = "PayAsYouGo"
  ssl_enabled           = "0"
  engine_version        = "6.0"
  engine                = "gpdb"
  zone_id               = "cn-beijing-h"
  vswitch_id            = alicloud_vswitch.defaultPlruct.id
  storage_size          = "50"
  master_cu             = "4"
  vpc_id                = alicloud_vpc.defaultiYyNGW.id
  db_instance_mode      = "StorageElastic"
}

resource "alicloud_gpdb_account" "defaultOwner" {
  account_name        = "tf_example"
  account_password    = "Example1234"
  account_description = "tf_example"
  db_instance_id      = alicloud_gpdb_instance.defaultqsmpIy.id
}

resource "alicloud_gpdb_database" "defaultPPmRVa" {
  owner              = alicloud_gpdb_account.defaultOwner.account_name
  database_name      = "seagull"
  db_instance_id     = alicloud_gpdb_instance.defaultqsmpIy.id
  character_set_name = "UTF8"
  collate            = "en_US.utf8"
  ctype              = "en_US.utf8"
}

resource "alicloud_gpdb_db_extension" "default" {
  extension_name    = "uuid-ossp"
  db_instance_id    = alicloud_gpdb_instance.defaultqsmpIy.id
  database_name     = alicloud_gpdb_database.defaultPPmRVa.database_name
  is_latest_version = true
}

data "alicloud_gpdb_db_extensions" "default" {
  db_instance_id = alicloud_gpdb_instance.defaultqsmpIy.id
  database_name  = alicloud_gpdb_db_extension.default.database_name
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
