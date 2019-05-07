package alicloud

import (
	"testing"
)

func TestAccAlicloudRamAccountAliasDataSource(t *testing.T) {
	basicConf := dataSourceTestAccConfig{
		existConfig: testAccAlicloudRamAccountAliasDataSourceConfig(),
	}
	accountAliasCheckInfo.dataSourceTestCheck(t, -1, basicConf)
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
