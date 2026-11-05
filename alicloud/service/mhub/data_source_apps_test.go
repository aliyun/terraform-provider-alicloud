package mhub_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudMhubAppDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudMhubAppDataSourceName(rand, map[string]string{
			"product_id": `"${alicloud_mhub_app.default.product_id}"`,
			"ids":        `["${alicloud_mhub_app.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudMhubAppDataSourceName(rand, map[string]string{
			"product_id": `"${alicloud_mhub_app.default.product_id}"`,
			"ids":        `["${alicloud_mhub_app.default.id}_fake"]`,
		}),
	}

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudMhubAppDataSourceName(rand, map[string]string{
			"product_id": `"${alicloud_mhub_app.default.product_id}"`,
			"name_regex": `"${alicloud_mhub_app.default.app_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudMhubAppDataSourceName(rand, map[string]string{
			"product_id": `"${alicloud_mhub_app.default.product_id}"`,
			"name_regex": `"${alicloud_mhub_app.default.app_name}_fake"`,
		}),
	}

	osTypeConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudMhubAppDataSourceName(rand, map[string]string{
			"product_id": `"${alicloud_mhub_app.default.product_id}"`,
			"os_type":    `"Android"`,
		}),
		FakeConfig: testAccCheckAlicloudMhubAppDataSourceName(rand, map[string]string{
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
			"apps.0.app_key":      acctest.CHECKSET,
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

	var AlicloudMhubAppCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_mhub_apps.default",
		ExistMapFunc: existAlicloudMhubAppDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudMhubAppDataSourceNameMapFunc,
	}

	preCheck := func() {
		acctest.PreCheck(t)
		acctest.PreCheckWithRegions(t, true, connectivity.MHUBSupportRegions)
	}

	AlicloudMhubAppCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, osTypeConf)
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
