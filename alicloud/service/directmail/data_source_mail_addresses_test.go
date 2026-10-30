package directmail_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func SkipTestAccAliCloudDirectMailMailAddressesDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudDirectMailMailAddressesDataSourceName(rand, map[string]string{
			"ids": `[alicloud_direct_mail_mail_address.default.id]`,
		}),
		FakeConfig: testAccCheckAlicloudDirectMailMailAddressesDataSourceName(rand, map[string]string{
			"ids": `["fake"]`,
		}),
	}
	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudDirectMailMailAddressesDataSourceName(rand, map[string]string{
			"key_word": `"${alicloud_direct_mail_mail_address.default.account_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudDirectMailMailAddressesDataSourceName(rand, map[string]string{
			"key_word": `"tf-fake@xxx.changes.com.cn"`,
		}),
	}
	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudDirectMailMailAddressesDataSourceName(rand, map[string]string{
			"ids":    `[alicloud_direct_mail_mail_address.default.id]`,
			"status": `"0"`,
		}),
		FakeConfig: testAccCheckAlicloudDirectMailMailAddressesDataSourceName(rand, map[string]string{
			"ids":    `[alicloud_direct_mail_mail_address.default.id]`,
			"status": `"1"`,
		}),
	}
	sendtypeConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudDirectMailMailAddressesDataSourceName(rand, map[string]string{
			"ids":      `[alicloud_direct_mail_mail_address.default.id]`,
			"sendtype": `"batch"`,
		}),
		FakeConfig: testAccCheckAlicloudDirectMailMailAddressesDataSourceName(rand, map[string]string{
			"ids":      `[alicloud_direct_mail_mail_address.default.id]`,
			"sendtype": `"trigger"`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudDirectMailMailAddressesDataSourceName(rand, map[string]string{
			"ids":      `[alicloud_direct_mail_mail_address.default.id]`,
			"key_word": `"${alicloud_direct_mail_mail_address.default.account_name}"`,
			"sendtype": `"batch"`,
			"status":   `0`,
		}),
		FakeConfig: testAccCheckAlicloudDirectMailMailAddressesDataSourceName(rand, map[string]string{
			"ids":      `["fake"]`,
			"key_word": `"tf-fake@xxx.changes.com.cn"`,
			"sendtype": `"trigger"`,
			"status":   `1`,
		}),
	}
	var existAlicloudDirectMailMailAddressesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"addresses.#":        "1",
			"addresses.0.id":     acctest.CHECKSET,
			"addresses.0.status": `0`,
		}
	}
	var fakeAlicloudDirectMailMailAddressesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"addresses.#": "0",
		}
	}
	var alicloudDirectMailMailAddressesCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_direct_mail_mail_addresses.default",
		ExistMapFunc: existAlicloudDirectMailMailAddressesDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudDirectMailMailAddressesDataSourceNameMapFunc,
	}
	preCheck := func() {
		acctest.PreCheck(t)
		acctest.PreCheckWithRegions(t, true, connectivity.DmSupportRegions)
	}
	alicloudDirectMailMailAddressesCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, statusConf, sendtypeConf, allConf)
}
func testAccCheckAlicloudDirectMailMailAddressesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAcc-%d"
}

resource "alicloud_direct_mail_mail_address" "default" {
	sendtype     = "batch"
	account_name = "tf-testacc%d@xxx.changes.com.cn"
}

data "alicloud_direct_mail_mail_addresses" "default" {
	%s
}
`, rand, rand, strings.Join(pairs, " \n "))
	return config
}
