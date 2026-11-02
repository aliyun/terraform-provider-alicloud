package ram_test

import (
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAliCloudRamAccountAliasDataSource(t *testing.T) {
	basicConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccAlicloudRamAccountAliasDataSourceConfig(),
	}
	preCheck := func() {
		acctest.PreCheckWithRegions(t, true, connectivity.RamNoSkipRegions)
	}
	accountAliasCheckInfo.DataSourceTestCheckWithPreCheck(t, -1, preCheck, basicConf)
}

func testAccAlicloudRamAccountAliasDataSourceConfig() string {
	config := `
data "alicloud_ram_account_aliases" "default" {
}`
	return config
}

var existRamAccountAliasMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"account_alias": acctest.CHECKSET,
	}
}

var fakeRamAccountAliasMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"account_alias": "",
	}
}

var accountAliasCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_ram_account_aliases.default",
	ExistMapFunc: existRamAccountAliasMapFunc,
	FakeMapFunc:  fakeRamAccountAliasMapFunc,
}
