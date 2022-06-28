package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudECDUsersDataSource(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	checkoutSupportedRegions(t, true, connectivity.EcdUserSupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcdUsersDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ecd_user.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudEcdUsersDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ecd_user.default.id}_fake"]`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcdUsersDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_ecd_user.default.id}"]`,
			"status": `"Unlocked"`,
		}),
		fakeConfig: testAccCheckAlicloudEcdUsersDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_ecd_user.default.id}"]`,
			"status": `"Locked"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcdUsersDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_ecd_user.default.id}"]`,
			"status": `"Unlocked"`,
		}),
		fakeConfig: testAccCheckAlicloudEcdUsersDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_ecd_user.default.id}_fake"]`,
			"status": `"Locked"`,
		}),
	}
	var existAlicloudEcdUsersDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":               "1",
			"users.#":             "1",
			"users.0.email":       "hello.uuu@aaa.com",
			"users.0.phone":       fmt.Sprintf("158016%d", rand),
			"users.0.end_user_id": fmt.Sprintf("tf_testaccecduser%d", rand),
			"users.0.status":      CHECKSET,
		}
	}
	var fakeAlicloudEcdUsersDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"users.#": "0",
		}
	}
	var alicloudEcdUsersCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_ecd_users.default",
		existMapFunc: existAlicloudEcdUsersDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudEcdUsersDataSourceNameMapFunc,
	}

	alicloudEcdUsersCheckInfo.dataSourceTestCheck(t, rand, idsConf, statusConf, allConf)
}
func testAccCheckAlicloudEcdUsersDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`

resource "alicloud_ecd_user" "default" {
	end_user_id = "tf_testaccecduser%d"
	email       = "hello.uuu@aaa.com"
	phone       = "158016%d"
	password    = "%d"
}

data "alicloud_ecd_users" "default" {	
	%s
}
`, rand, rand, rand, strings.Join(pairs, " \n "))
	return config
}
