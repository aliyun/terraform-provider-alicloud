package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func SkipTestAccAlicloudDirectMailMailAddressesDataSource(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDirectMailMailAddressesDataSourceName(rand, map[string]string{
			"ids": `[alicloud_direct_mail_mail_address.default.id]`,
		}),
		fakeConfig: testAccCheckAlicloudDirectMailMailAddressesDataSourceName(rand, map[string]string{
			"ids": `["fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDirectMailMailAddressesDataSourceName(rand, map[string]string{
			"key_word": `"${alicloud_direct_mail_mail_address.default.account_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudDirectMailMailAddressesDataSourceName(rand, map[string]string{
			"key_word": `"tf-fake@xxx.changes.com.cn"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDirectMailMailAddressesDataSourceName(rand, map[string]string{
			"ids":    `[alicloud_direct_mail_mail_address.default.id]`,
			"status": `"0"`,
		}),
		fakeConfig: testAccCheckAlicloudDirectMailMailAddressesDataSourceName(rand, map[string]string{
			"ids":    `[alicloud_direct_mail_mail_address.default.id]`,
			"status": `"1"`,
		}),
	}
	sendtypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDirectMailMailAddressesDataSourceName(rand, map[string]string{
			"ids":      `[alicloud_direct_mail_mail_address.default.id]`,
			"sendtype": `"batch"`,
		}),
		fakeConfig: testAccCheckAlicloudDirectMailMailAddressesDataSourceName(rand, map[string]string{
			"ids":      `[alicloud_direct_mail_mail_address.default.id]`,
			"sendtype": `"trigger"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDirectMailMailAddressesDataSourceName(rand, map[string]string{
			"ids":      `[alicloud_direct_mail_mail_address.default.id]`,
			"key_word": `"${alicloud_direct_mail_mail_address.default.account_name}"`,
			"sendtype": `"batch"`,
			"status":   `0`,
		}),
		fakeConfig: testAccCheckAlicloudDirectMailMailAddressesDataSourceName(rand, map[string]string{
			"ids":      `["fake"]`,
			"key_word": `"tf-fake@xxx.changes.com.cn"`,
			"sendtype": `"trigger"`,
			"status":   `1`,
		}),
	}
	var existAlicloudDirectMailMailAddressesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"addresses.#":        "1",
			"addresses.0.id":     CHECKSET,
			"addresses.0.status": `0`,
		}
	}
	var fakeAlicloudDirectMailMailAddressesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"addresses.#": "0",
		}
	}
	var alicloudDirectMailMailAddressesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_direct_mail_mail_addresses.default",
		existMapFunc: existAlicloudDirectMailMailAddressesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudDirectMailMailAddressesDataSourceNameMapFunc,
	}
	preCheck := func() {
		testAccPreCheck(t)
		testAccPreCheckWithRegions(t, true, connectivity.DmSupportRegions)
	}
	alicloudDirectMailMailAddressesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, statusConf, sendtypeConf, allConf)
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
