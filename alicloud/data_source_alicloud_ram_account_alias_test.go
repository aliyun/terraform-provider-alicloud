package alicloud

import (
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudRAMAccountAliasDataSource(t *testing.T) {
	basicConf := dataSourceTestAccConfig{
		existConfig: testAccAlicloudRamAccountAliasDataSourceConfig(),
	}
	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.RamNoSkipRegions)
	}
	accountAliasCheckInfo.dataSourceTestCheckWithPreCheck(t, -1, preCheck, basicConf)
}

func testAccAlicloudRamAccountAliasDataSourceConfig() string {
	config := `
data "alicloud_ram_account_aliases" "default" {
}`
	return config
}

var existRamAccountAliasMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"account_alias": CHECKSET,
	}
}

var fakeRamAccountAliasMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"account_alias": "",
	}
}

var accountAliasCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_ram_account_aliases.default",
	existMapFunc: existRamAccountAliasMapFunc,
	fakeMapFunc:  fakeRamAccountAliasMapFunc,
}
