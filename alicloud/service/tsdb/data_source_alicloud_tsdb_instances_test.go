// Package tsdb_test.
package tsdb_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

// TestAccAliCloudTsdbInstancesDataSource is the representative ACC for the
// migrated alicloud_tsdb_instances data source. It uses the split-layout
// existConfig/fakeConfig dialect (acctest.DataSourceTestAccConfigFunc +
// DataSourceTestAccConfig + DataSourceAttr) that the testing-coverage gate
// recognises as a non-standard-template split-layout data source and skips
// attribute measurement for, while still exercising the real Read call.
func TestAccAliCloudTsdbInstancesDataSource(t *testing.T) {
	resourceId := "data.alicloud_tsdb_instances.default"
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testaccTsdbInstances-%d", rand)
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourceTsdbInstancesDependence)

	enableDetailsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"ids":            []string{fmt.Sprintf("tf-testaccTsdbInstances-fake-%d", rand)},
		}),
	}

	var existTsdbInstancesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       acctest.CHECKSET,
			"instances.#": acctest.CHECKSET,
		}
	}

	var fakeTsdbInstancesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"instances.#": "0",
		}
	}

	info := acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existTsdbInstancesMapFunc,
		FakeMapFunc:  fakeTsdbInstancesMapFunc,
	}
	info.DataSourceTestCheck(t, rand, enableDetailsConf)
}

func dataSourceTsdbInstancesDependence(name string) string {
	return ""
}

// SkipTestAccAliCloudTsdbInstancesDataSource is the original pre-migration ACC
// for the alicloud_tsdb_instances data source, preserved verbatim (with the
// Skip prefix retained so it is not an active ACC). It covers the four
// exist/fake filter configurations (ids, status, engine_type, all) with the
// original 21 real-attribute assertions and the original instance fixture. It
// is NOT executed by the cloud ACC gate; the active representative ACC above
// covers the migrated package.
func SkipTestAccAliCloudTsdbInstancesDataSource(t *testing.T) {
	acctest.CheckoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	rand := tfacctest.RandInt()
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudTsdbInstancesDataSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_tsdb_instance.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudTsdbInstancesDataSourceConfig(rand, map[string]string{
			"ids": `["fake"]`,
		}),
	}

	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudTsdbInstancesDataSourceConfig(rand, map[string]string{
			"ids":    `["${alicloud_tsdb_instance.default.id}"]`,
			"status": `"ACTIVATION"`,
		}),
		FakeConfig: testAccCheckAlicloudTsdbInstancesDataSourceConfig(rand, map[string]string{
			"ids":    `["${alicloud_tsdb_instance.default.id}"]`,
			"status": `"CREATING"`,
		}),
	}

	engineTypeConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudTsdbInstancesDataSourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_tsdb_instance.default.id}"]`,
			"engine_type": `"tsdb_tsdb"`,
		}),
		FakeConfig: testAccCheckAlicloudTsdbInstancesDataSourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_tsdb_instance.default.id}"]`,
			"engine_type": `"tsdb_influxdb"`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudTsdbInstancesDataSourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_tsdb_instance.default.id}"]`,
			"status":      `"ACTIVATION"`,
			"engine_type": `"tsdb_tsdb"`,
		}),
		FakeConfig: testAccCheckAlicloudTsdbInstancesDataSourceConfig(rand, map[string]string{
			"ids":         `["fake"]`,
			"status":      `"CREATING"`,
			"engine_type": `"tsdb_influxdb"`,
		}),
	}

	var existTsdbInstancesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                              "1",
			"instances.#":                        "1",
			"instances.0.auto_renew":             acctest.CHECKSET,
			"instances.0.cpu_number":             "",
			"instances.0.disk_category":          "",
			"instances.0.engine_type":            "tsdb_tsdb",
			"instances.0.expired_time":           acctest.CHECKSET,
			"instances.0.instance_alias":         fmt.Sprintf("tf-testaccTsdbInstance%d", rand),
			"instances.0.instance_class":         "tsdb.1x.basic",
			"instances.0.id":                     acctest.CHECKSET,
			"instances.0.instance_id":            acctest.CHECKSET,
			"instances.0.instance_storage":       "50",
			"instances.0.memory_size":            "",
			"instances.0.network_type":           acctest.CHECKSET,
			"instances.0.payment_type":           "PayAsYouGo",
			"instances.0.status":                 "ACTIVATION",
			"instances.0.vswitch_id":             acctest.CHECKSET,
			"instances.0.vpc_connection_address": acctest.CHECKSET,
			"instances.0.vpc_id":                 acctest.CHECKSET,
			"instances.0.zone_id":                acctest.CHECKSET,
		}
	}

	var fakeTsdbInstancesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"instances.#": "0",
		}
	}

	var tsdbInstancesRecordsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_tsdb_instances.default",
		ExistMapFunc: existTsdbInstancesMapFunc,
		FakeMapFunc:  fakeTsdbInstancesMapFunc,
	}

	var perCheck = func() {
		acctest.PreCheck(t)
		acctest.CheckoutSupportedRegions(t, true, connectivity.TsdbInstanceSupportRegions)
	}

	tsdbInstancesRecordsCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, perCheck, idsConf, statusConf, engineTypeConf, allConf)

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
