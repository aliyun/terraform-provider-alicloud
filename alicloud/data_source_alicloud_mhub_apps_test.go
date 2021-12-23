package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudMhubAppDataSource(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMhubAppDataSourceName(rand, map[string]string{
			"product_id": `"${alicloud_mhub_app.default.product_id}"`,
			"ids":        `["${alicloud_mhub_app.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudMhubAppDataSourceName(rand, map[string]string{
			"product_id": `"${alicloud_mhub_app.default.product_id}"`,
			"ids":        `["${alicloud_mhub_app.default.id}_fake"]`,
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMhubAppDataSourceName(rand, map[string]string{
			"product_id": `"${alicloud_mhub_app.default.product_id}"`,
			"name_regex": `"${alicloud_mhub_app.default.app_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudMhubAppDataSourceName(rand, map[string]string{
			"product_id": `"${alicloud_mhub_app.default.product_id}"`,
			"name_regex": `"${alicloud_mhub_app.default.app_name}_fake"`,
		}),
	}

	osTypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMhubAppDataSourceName(rand, map[string]string{
			"product_id": `"${alicloud_mhub_app.default.product_id}"`,
			"os_type":    `"Android"`,
		}),
		fakeConfig: testAccCheckAlicloudMhubAppDataSourceName(rand, map[string]string{
			"product_id": `"${alicloud_mhub_app.default.product_id}"`,
			"os_type":    `"iOS"`,
		}),
	}

	var existAlicloudMhubAppDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":               "1",
			"names.#":             "1",
			"apps.#":              "1",
			"apps.0.app_name":     fmt.Sprintf("tf_testaccmhub%d", rand),
			"apps.0.app_key":      CHECKSET,
			"apps.0.industry_id":  "0",
			"apps.0.package_name": "com.test.android",
			"apps.0.type":         "Android",
		}
	}
	var fakeAlicloudMhubAppDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
			"apps.#":  "0",
		}
	}

	var AlicloudMhubAppCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_mhub_apps.default",
		existMapFunc: existAlicloudMhubAppDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudMhubAppDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
		testAccPreCheckWithRegions(t, true, connectivity.MHUBSupportRegions)
	}

	AlicloudMhubAppCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, osTypeConf)
}
func testAccCheckAlicloudMhubAppDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf_testaccmhub%d"
}

resource "alicloud_mhub_product" "default"{
  product_name = var.name
}

resource "alicloud_mhub_app" "default"{
  app_name = var.name
  product_id = alicloud_mhub_product.default.id
  package_name = "com.test.android"
  type = "Android"
}

data "alicloud_mhub_apps" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
