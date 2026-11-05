package ecd_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudECDUsersDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudEcdUsersDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ecd_user.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudEcdUsersDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ecd_user.default.id}_fake"]`,
		}),
	}
	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudEcdUsersDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_ecd_user.default.id}"]`,
			"status": `"Unlocked"`,
		}),
		FakeConfig: testAccCheckAlicloudEcdUsersDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_ecd_user.default.id}"]`,
			"status": `"Locked"`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudEcdUsersDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_ecd_user.default.id}"]`,
			"status": `"Unlocked"`,
		}),
		FakeConfig: testAccCheckAlicloudEcdUsersDataSourceName(rand, map[string]string{
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
			"users.0.status":      acctest.CHECKSET,
		}
	}
	var fakeAlicloudEcdUsersDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"users.#": "0",
		}
	}
	var alicloudEcdUsersCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_ecd_users.default",
		ExistMapFunc: existAlicloudEcdUsersDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudEcdUsersDataSourceNameMapFunc,
	}

	alicloudEcdUsersCheckInfo.DataSourceTestCheck(t, rand, idsConf, statusConf, allConf)
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
