package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudServiceCatalogProductVersionDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	checkoutSupportedRegions(t, true, connectivity.ServiceCatalogProvisionedProductSupportRegions)

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudServiceCatalogProductVersionSourceConfig(rand, map[string]string{
			"name_regex": `"1.0.0"`,
		}),
		fakeConfig: testAccCheckAlicloudServiceCatalogProductVersionSourceConfig(rand, map[string]string{
			"name_regex": `"1.0.0_fake"`,
		}),
	}

	ServiceCatalogProductVersionCheckInfo.dataSourceTestCheck(t, rand, allConf)
}

var existServiceCatalogProductVersionMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"versions.#":                              "1",
		"versions.0.id":                           CHECKSET,
		"versions.0.active":                       CHECKSET,
		"versions.0.create_time":                  CHECKSET,
		"versions.0.guidance":                     CHECKSET,
		"versions.0.product_id":                   CHECKSET,
		"versions.0.product_version_id":           CHECKSET,
		"versions.0.product_version_name":         CHECKSET,
		"versions.0.template_type":                CHECKSET,
		"versions.0.template_url":                 CHECKSET,
		"product_versions.#":                      "1",
		"product_versions.0.id":                   CHECKSET,
		"product_versions.0.active":               CHECKSET,
		"product_versions.0.create_time":          CHECKSET,
		"product_versions.0.guidance":             CHECKSET,
		"product_versions.0.product_id":           CHECKSET,
		"product_versions.0.product_version_id":   CHECKSET,
		"product_versions.0.product_version_name": CHECKSET,
		"product_versions.0.template_type":        CHECKSET,
		"product_versions.0.template_url":         CHECKSET,
	}
}

var fakeServiceCatalogProductVersionMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"versions.#":         "0",
		"product_versions.#": "0",
	}
}

var ServiceCatalogProductVersionCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_service_catalog_product_versions.default",
	existMapFunc: existServiceCatalogProductVersionMapFunc,
	fakeMapFunc:  fakeServiceCatalogProductVersionMapFunc,
}

func testAccCheckAlicloudServiceCatalogProductVersionSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
  default = "tf-testAccServiceCatalogProductVersion%d"
}

data "alicloud_service_catalog_end_user_products" "default" {
  name_regex = "ram模板创建"
}

data "alicloud_service_catalog_product_versions" "default" {
  product_id = data.alicloud_service_catalog_end_user_products.default.end_user_products.0.id
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
