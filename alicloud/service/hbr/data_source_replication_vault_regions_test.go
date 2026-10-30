package hbr_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudHbrReplicationVaultRegionsDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(100, 999)
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudHbrVaultReplicationRegionsDataSourceName(rand, map[string]string{}),
		FakeConfig:  "",
	}

	var existAlicloudHbrVaultReplicationRegionDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"regions.#":                       acctest.CHECKSET,
			"regions.0.replication_region_id": acctest.CHECKSET,
		}
	}
	var fakeHbrVaultReplicationRegionsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"regions.#": "0",
		}
	}
	var alicloudHbrVaultReplicationAccountBusesCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_hbr_replication_vault_regions.default",
		ExistMapFunc: existAlicloudHbrVaultReplicationRegionDataSourceNameMapFunc,
		FakeMapFunc:  fakeHbrVaultReplicationRegionsMapFunc,
	}

	alicloudHbrVaultReplicationAccountBusesCheckInfo.DataSourceTestCheck(t, rand, allConf)
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
