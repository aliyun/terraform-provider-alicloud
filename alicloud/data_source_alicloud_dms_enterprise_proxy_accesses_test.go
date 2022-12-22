package alicloud

import (
	"fmt"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudDMSEnterpriseProxyAccessDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	checkoutSupportedRegions(t, true, connectivity.DMSEnterpriseProxyAccessSupportRegions)

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDMSEnterpriseProxyAccessSourceConfig(rand, map[string]string{
			"enable_details": `"true"`,
			"ids":            `["${alicloud_dms_enterprise_proxy_access.default.id}"]`,
			"proxy_id":       `"${data.alicloud_dms_enterprise_proxies.ids.proxies.0.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudDMSEnterpriseProxyAccessSourceConfig(rand, map[string]string{
			"ids":      `["${alicloud_dms_enterprise_proxy_access.default.id}_fake"]`,
			"proxy_id": `"${data.alicloud_dms_enterprise_proxies.ids.proxies.0.id}"`,
		}),
	}

	DMSEnterpriseProxyAccessCheckInfo.dataSourceTestCheck(t, rand, allConf)
}

var existDMSEnterpriseProxyAccessMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"accesses.#":                 "1",
		"accesses.0.id":              CHECKSET,
		"accesses.0.access_id":       CHECKSET,
		"accesses.0.access_secret":   CHECKSET,
		"accesses.0.instance_id":     CHECKSET,
		"accesses.0.proxy_access_id": CHECKSET,
		"accesses.0.proxy_id":        CHECKSET,
		"accesses.0.user_id":         CHECKSET,
		"accesses.0.user_name":       CHECKSET,
		"accesses.0.user_uid":        CHECKSET,
		"accesses.0.gmt_create":      CHECKSET,
	}
}

var fakeDMSEnterpriseProxyAccessMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"accesses.#": "0",
	}
}

var DMSEnterpriseProxyAccessCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_dms_enterprise_proxy_accesses.default",
	existMapFunc: existDMSEnterpriseProxyAccessMapFunc,
	fakeMapFunc:  fakeDMSEnterpriseProxyAccessMapFunc,
}

func testAccCheckAlicloudDMSEnterpriseProxyAccessSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccDMSEnterpriseProxyAccess%d"
}

data "alicloud_dms_enterprise_users" "dms_enterprise_users_ds" {
  role   = "USER"
  status = "NORMAL"
}
data "alicloud_dms_enterprise_proxies" "ids" {}

resource "alicloud_dms_enterprise_proxy_access" "default" {
  proxy_id       = data.alicloud_dms_enterprise_proxies.ids.proxies.0.id
  user_id        = data.alicloud_dms_enterprise_users.dms_enterprise_users_ds.users.0.user_id
}

data "alicloud_dms_enterprise_proxy_accesses" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
