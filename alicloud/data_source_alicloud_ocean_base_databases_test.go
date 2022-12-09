package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func SkipTestAccAlicloudOceanBaseDatabasesDataSource(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	checkoutSupportedRegions(t, true, connectivity.OceanBaseSupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudOceanBaseDatabaseDataSourceName(rand, map[string]string{
			"tenant_id":   `"t48xt2lc5vqrk"`,
			"ids":         `["${alicloud_ocean_base_database.default.id}"]`,
			"with_tables": "true",
		}),
		fakeConfig: testAccCheckAlicloudOceanBaseDatabaseDataSourceName(rand, map[string]string{
			"tenant_id":   `"t48xt2lc5vqrk"`,
			"ids":         `["${alicloud_ocean_base_database.default.id}_fake"]`,
			"with_tables": "true",
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudOceanBaseDatabaseDataSourceName(rand, map[string]string{
			"tenant_id":   `"t48xt2lc5vqrk"`,
			"name_regex":  `"${alicloud_ocean_base_database.default.database_name}"`,
			"with_tables": "true",
		}),
		fakeConfig: testAccCheckAlicloudOceanBaseDatabaseDataSourceName(rand, map[string]string{
			"tenant_id":   `"t48xt2lc5vqrk"`,
			"name_regex":  `"${alicloud_ocean_base_database.default.database_name}_fake"`,
			"with_tables": "true",
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudOceanBaseDatabaseDataSourceName(rand, map[string]string{
			"tenant_id":   `"t48xt2lc5vqrk"`,
			"ids":         `["${alicloud_ocean_base_database.default.id}"]`,
			"status":      `"${alicloud_ocean_base_database.default.status}"`,
			"with_tables": "true",
		}),
		fakeConfig: testAccCheckAlicloudOceanBaseDatabaseDataSourceName(rand, map[string]string{
			"tenant_id":   `"t48xt2lc5vqrk"`,
			"ids":         `["${alicloud_ocean_base_database.default.id}_fake"]`,
			"status":      `"DELETING"`,
			"with_tables": "true",
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudOceanBaseDatabaseDataSourceName(rand, map[string]string{
			"tenant_id":   `"t48xt2lc5vqrk"`,
			"ids":         `["${alicloud_ocean_base_database.default.id}"]`,
			"status":      `"${alicloud_ocean_base_database.default.status}"`,
			"name_regex":  `"${alicloud_ocean_base_database.default.database_name}"`,
			"with_tables": "true",
		}),
		fakeConfig: testAccCheckAlicloudOceanBaseDatabaseDataSourceName(rand, map[string]string{
			"tenant_id":   `"t48xt2lc5vqrk"`,
			"ids":         `["${alicloud_ocean_base_database.default.id}_fake"]`,
			"status":      `"DELETING"`,
			"name_regex":  `"${alicloud_ocean_base_database.default.database_name}_fake"`,
			"with_tables": "true",
		}),
	}
	var existAlicloudOceanBaseDatabasesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                     "1",
			"databases.#":               "1",
			"databases.0.collation":     "utf8mb4_general_ci",
			"databases.0.database_name": fmt.Sprintf("tfacccoceanbase%d", rand),
			"databases.0.id":            CHECKSET,
			"databases.0.description":   fmt.Sprintf("tfacccoceanbase%d", rand),
			"databases.0.encoding":      "utf8mb4",
			"databases.0.status":        "ONLINE",
			"databases.0.tenant_id":     CHECKSET,
			"databases.0.users.#":       NOSET,
			"databases.0.tables.#":      NOSET,
		}
	}
	var fakeAlicloudOceanBaseDatabasesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"databases.#": "0",
		}
	}
	var AlicloudOceanBaseDatabaseCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_ocean_base_databases.default",
		existMapFunc: existAlicloudOceanBaseDatabasesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudOceanBaseDatabasesDataSourceNameMapFunc,
	}
	preCheck := func() {
		checkoutSupportedRegions(t, true, connectivity.OceanBaseSupportRegions)
	}
	AlicloudOceanBaseDatabaseCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, statusConf, allConf)
}
func testAccCheckAlicloudOceanBaseDatabaseDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tfacccoceanbase%d"
}

resource "alicloud_ocean_base_database" "default" {
  collation     = "utf8mb4_general_ci"
  database_name = var.name
  description   = var.name
  encoding      = "utf8mb4"
  instance_id   = "ob48azzvzketpc"
  tenant_id     = "t48xt2lc5vqrk"
}

data "alicloud_ocean_base_databases" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
