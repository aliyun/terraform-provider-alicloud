package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudTsdbInstancesDataSource(t *testing.T) {
	checkoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudTsdbInstancesDataSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_tsdb_instance.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudTsdbInstancesDataSourceConfig(rand, map[string]string{
			"ids": `["fake"]`,
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudTsdbInstancesDataSourceConfig(rand, map[string]string{
			"ids":    `["${alicloud_tsdb_instance.default.id}"]`,
			"status": `"ACTIVATION"`,
		}),
		fakeConfig: testAccCheckAlicloudTsdbInstancesDataSourceConfig(rand, map[string]string{
			"ids":    `["${alicloud_tsdb_instance.default.id}"]`,
			"status": `"CREATING"`,
		}),
	}

	engineTypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudTsdbInstancesDataSourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_tsdb_instance.default.id}"]`,
			"engine_type": `"tsdb_tsdb"`,
		}),
		fakeConfig: testAccCheckAlicloudTsdbInstancesDataSourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_tsdb_instance.default.id}"]`,
			"engine_type": `"tsdb_influxdb"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudTsdbInstancesDataSourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_tsdb_instance.default.id}"]`,
			"status":      `"ACTIVATION"`,
			"engine_type": `"tsdb_tsdb"`,
		}),
		fakeConfig: testAccCheckAlicloudTsdbInstancesDataSourceConfig(rand, map[string]string{
			"ids":         `["fake"]`,
			"status":      `"CREATING"`,
			"engine_type": `"tsdb_influxdb"`,
		}),
	}

	var existTsdbInstancesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                              "1",
			"instances.#":                        "1",
			"instances.0.auto_renew":             CHECKSET,
			"instances.0.cpu_number":             "",
			"instances.0.disk_category":          "",
			"instances.0.engine_type":            "tsdb_tsdb",
			"instances.0.expired_time":           CHECKSET,
			"instances.0.instance_alias":         fmt.Sprintf("tf-testaccTsdbInstance%d", rand),
			"instances.0.instance_class":         "tsdb.1x.basic",
			"instances.0.id":                     CHECKSET,
			"instances.0.instance_id":            CHECKSET,
			"instances.0.instance_storage":       "50",
			"instances.0.memory_size":            "",
			"instances.0.network_type":           CHECKSET,
			"instances.0.payment_type":           "PayAsYouGo",
			"instances.0.status":                 "ACTIVATION",
			"instances.0.vswitch_id":             CHECKSET,
			"instances.0.vpc_connection_address": CHECKSET,
			"instances.0.vpc_id":                 CHECKSET,
			"instances.0.zone_id":                CHECKSET,
		}
	}

	var fakeTsdbInstancesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"instances.#": "0",
		}
	}

	var tsdbInstancesRecordsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_tsdb_instances.default",
		existMapFunc: existTsdbInstancesMapFunc,
		fakeMapFunc:  fakeTsdbInstancesMapFunc,
	}

	var perCheck = func() {
		testAccPreCheck(t)
		testAccPreCheckWithRegions(t, true, connectivity.TsdbInstanceSupportRegions)
	}

	tsdbInstancesRecordsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, perCheck, idsConf, statusConf, engineTypeConf, allConf)

}

func testAccCheckAlicloudTsdbInstancesDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testaccTsdbInstance%d"
}

data "alicloud_tsdb_zones" "default" {}

resource "alicloud_vpc" "default" {
  cidr_block = "192.168.0.0/16"
  vpc_name = var.name
}

resource "alicloud_vswitch" "default" {
  availability_zone = data.alicloud_tsdb_zones.default.ids.0
  cidr_block = "192.168.1.0/24"
  vpc_id = alicloud_vpc.default.id
}

resource "alicloud_tsdb_instance" "default" {
  payment_type = "PayAsYouGo"
  vswitch_id = alicloud_vswitch.default.id
  instance_storage = "50"
  instance_class = "tsdb.1x.basic"
  engine_type = "tsdb_tsdb"
  instance_alias = var.name
}

data "alicloud_tsdb_instances" "default" {
	enable_details = "true"
  %s
}
`, rand, strings.Join(pairs, "\n  "))
	return config
}
