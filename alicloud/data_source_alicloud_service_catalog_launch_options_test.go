package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudServiceCatalogLaunchOptionDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	checkoutSupportedRegions(t, true, connectivity.ServiceCatalogProvisionedProductSupportRegions)

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudServiceCatalogLaunchOptionSourceConfig(rand, map[string]string{
			"name_regex": `"ram模板创建"`,
		}),
		fakeConfig: testAccCheckAlicloudServiceCatalogLaunchOptionSourceConfig(rand, map[string]string{
			"name_regex": `"ram模板创建_fake"`,
		}),
	}

	ServiceCatalogLaunchOptionCheckInfo.dataSourceTestCheck(t, rand, allConf)
}

var existServiceCatalogLaunchOptionMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"options.#":                       "1",
		"options.0.id":                    CHECKSET,
		"options.0.portfolio_id":          CHECKSET,
		"options.0.portfolio_name":        CHECKSET,
		"launch_options.#":                "1",
		"launch_options.0.id":             CHECKSET,
		"launch_options.0.portfolio_id":   CHECKSET,
		"launch_options.0.portfolio_name": CHECKSET,
	}
}

var fakeServiceCatalogLaunchOptionMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"options.#":        "0",
		"launch_options.#": "0",
	}
}

var ServiceCatalogLaunchOptionCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_service_catalog_launch_options.default",
	existMapFunc: existServiceCatalogLaunchOptionMapFunc,
	fakeMapFunc:  fakeServiceCatalogLaunchOptionMapFunc,
}

func testAccCheckAlicloudServiceCatalogLaunchOptionSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccServiceCatalogLaunchOption%d"
}

data "alicloud_service_catalog_end_user_products" "default" {
  name_regex = "ram模板创建"
}

data "alicloud_service_catalog_launch_options" "default" {
	product_id = data.alicloud_service_catalog_end_user_products.default.end_user_products.0.id
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
