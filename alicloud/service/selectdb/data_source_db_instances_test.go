package selectdb_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudSelectDBDbInstancesDataSource(t *testing.T) {
	rand := tfacctest.RandInt()
	dbInstanceIdsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudSelectDBDbInstancesDataSource(rand, map[string]string{
			"ids": `["${alicloud_selectdb_db_instance.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudSelectDBDbInstancesDataSource(rand, map[string]string{
			"ids": `["${alicloud_selectdb_db_instance.default.id}_fake"]`,
		}),
	}

	tagsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudSelectDBDbInstancesDataSource(rand, map[string]string{
			"tags": `{
							"Created" = "TF-update"
	    					"For"     = "acceptance-test-update"
						}`,
		}),
		FakeConfig: testAccCheckAlicloudSelectDBDbInstancesDataSource(rand, map[string]string{
			"tags": `{
							"Created" = "TF-update-fake"
	    					"For"     = "acceptance-test-update-fake"
						}`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudSelectDBDbInstancesDataSource(rand, map[string]string{
			"ids": `["${alicloud_selectdb_db_instance.default.id}"]`,
			"tags": `{
							"Created" = "TF-update"
	    					"For"     = "acceptance-test-update"
						}`,
		}),
		FakeConfig: testAccCheckAlicloudSelectDBDbInstancesDataSource(rand, map[string]string{
			"ids": `["${alicloud_selectdb_db_instance.default.id}_fake"]`,
			"tags": `{
							"Created" = "TF-update-fake"
	    					"For"     = "acceptance-test-update-fake"
						}`,
		}),
	}

	var existAlicloudSelectDBDbInstancesDataSourceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                               "1",
			"instances.#":                         "1",
			"instances.0.payment_type":            "PayAsYouGo",
			"instances.0.db_instance_id":          acctest.CHECKSET,
			"instances.0.db_instance_description": fmt.Sprintf("tf-testAccSelectDBDbInstance-%d", rand),
			"instances.0.engine":                  "selectdb",
			"instances.0.region_id":               acctest.CHECKSET,
			"instances.0.vswitch_id":              acctest.CHECKSET,
			"instances.0.vpc_id":                  acctest.CHECKSET,
			"instances.0.zone_id":                 acctest.CHECKSET,
		}
	}
	var fakeAlicloudSelectDBDbInstancesDataSourceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"instance_ids.#": "0",
			"names.#":        "0",
		}
	}
	var alicloudSelectDBDbInstancesCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_selectdb_db_instances.default",
		ExistMapFunc: existAlicloudSelectDBDbInstancesDataSourceMapFunc,
		FakeMapFunc:  fakeAlicloudSelectDBDbInstancesDataSourceMapFunc,
	}
	preCheck := func() {
		acctest.PreCheckWithRegions(t, true, []connectivity.Region{"eu-central-1"})
	}

	alicloudSelectDBDbInstancesCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, dbInstanceIdsConf, tagsConf, allConf)
}

func testAccCheckAlicloudSelectDBDbInstancesDataSource(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {
  default = "tf-testAccSelectDBDbInstance-%d"
}
data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}
data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_zones.default.ids.0
}
resource "alicloud_selectdb_db_instance" "default" {
  db_instance_class       = "selectdb.g.2xlarge"
  db_instance_description = var.name
  cache_size              = "400"
  engine_minor_version    = "4.0"
  payment_type            = "PayAsYouGo"
  vpc_id                  = "${data.alicloud_vpcs.default.ids.0}"
  zone_id                 = "${data.alicloud_zones.default.ids.0}"
  vswitch_id              = "${data.alicloud_vswitches.default.ids.0}"
  tags = {
    Created = "TF-update"
    For     = "acceptance-test-update"
  }
}

data "alicloud_selectdb_db_instances" "default" {
  %s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
