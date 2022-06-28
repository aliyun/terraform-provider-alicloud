package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudHBRReplicationVaultRegionsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(100, 999)
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudHbrVaultReplicationRegionsDataSourceName(rand, map[string]string{}),
		fakeConfig:  "",
	}

	var existAlicloudHbrVaultReplicationRegionDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"regions.#":                       CHECKSET,
			"regions.0.replication_region_id": CHECKSET,
		}
	}
	var fakeHbrVaultReplicationRegionsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"regions.#": "0",
		}
	}
	var alicloudHbrVaultReplicationAccountBusesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_hbr_replication_vault_regions.default",
		existMapFunc: existAlicloudHbrVaultReplicationRegionDataSourceNameMapFunc,
		fakeMapFunc:  fakeHbrVaultReplicationRegionsMapFunc,
	}

	alicloudHbrVaultReplicationAccountBusesCheckInfo.dataSourceTestCheck(t, rand, allConf)
}
func testAccCheckAlicloudHbrVaultReplicationRegionsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
data "alicloud_hbr_replication_vault_regions" "default" {	
	%s
}
`, strings.Join(pairs, " \n "))
	return config
}
