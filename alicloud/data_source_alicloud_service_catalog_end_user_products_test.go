package alicloud

import (
	"fmt"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudServiceCatalogEndUserProductDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	checkoutSupportedRegions(t, true, connectivity.ServiceCatalogProvisionedProductSupportRegions)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudServiceCatalogEndUserProductSourceConfig(rand, map[string]string{
			"name_regex": `"ram模板创建"`,
			"sort_order": `"Asc"`,
			"sort_by":    `"CreateTime"`,
		}),
		fakeConfig: testAccCheckAlicloudServiceCatalogEndUserProductSourceConfig(rand, map[string]string{
			"name_regex": `"ram模板创建_fake"`,
			"sort_order": `"Asc"`,
			"sort_by":    `"CreateTime"`,
		}),
	}

	ServiceCatalogEndUserProductCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf)
}

var existServiceCatalogEndUserProductMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                           "1",
		"end_user_products.#":             "1",
		"end_user_products.0.id":          CHECKSET,
		"end_user_products.0.create_time": CHECKSET,
		"end_user_products.0.has_default_launch_option": CHECKSET,
		"end_user_products.0.product_arn":               CHECKSET,
		"end_user_products.0.product_id":                CHECKSET,
		"end_user_products.0.product_name":              CHECKSET,
		"end_user_products.0.product_type":              CHECKSET,
		"end_user_products.0.provider_name":             CHECKSET,
	}
}

var fakeServiceCatalogEndUserProductMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":               "0",
		"end_user_products.#": "0",
	}
}

var ServiceCatalogEndUserProductCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_service_catalog_end_user_products.default",
	existMapFunc: existServiceCatalogEndUserProductMapFunc,
	fakeMapFunc:  fakeServiceCatalogEndUserProductMapFunc,
}

func testAccCheckAlicloudServiceCatalogEndUserProductSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccServiceCatalogEndUserProduct%d"
}

data "alicloud_service_catalog_end_user_products" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
