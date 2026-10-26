package servicecatalog_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudServiceCatalogProductVersionDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(1000000, 9999999)
	acctest.CheckoutSupportedRegions(t, true, connectivity.ServiceCatalogProvisionedProductSupportRegions)

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudServiceCatalogProductVersionSourceConfig(rand, map[string]string{
			"name_regex": `"1.0.0"`,
		}),
		FakeConfig: testAccCheckAlicloudServiceCatalogProductVersionSourceConfig(rand, map[string]string{
			"name_regex": `"1.0.0_fake"`,
		}),
	}

	ServiceCatalogProductVersionCheckInfo.DataSourceTestCheck(t, rand, allConf)
}

var existServiceCatalogProductVersionMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"versions.#":                              "1",
		"versions.0.id":                           acctest.CHECKSET,
		"versions.0.active":                       acctest.CHECKSET,
		"versions.0.create_time":                  acctest.CHECKSET,
		"versions.0.guidance":                     acctest.CHECKSET,
		"versions.0.product_id":                   acctest.CHECKSET,
		"versions.0.product_version_id":           acctest.CHECKSET,
		"versions.0.product_version_name":         acctest.CHECKSET,
		"versions.0.template_type":                acctest.CHECKSET,
		"versions.0.template_url":                 acctest.CHECKSET,
		"product_versions.#":                      "1",
		"product_versions.0.id":                   acctest.CHECKSET,
		"product_versions.0.active":               acctest.CHECKSET,
		"product_versions.0.create_time":          acctest.CHECKSET,
		"product_versions.0.guidance":             acctest.CHECKSET,
		"product_versions.0.product_id":           acctest.CHECKSET,
		"product_versions.0.product_version_id":   acctest.CHECKSET,
		"product_versions.0.product_version_name": acctest.CHECKSET,
		"product_versions.0.template_type":        acctest.CHECKSET,
		"product_versions.0.template_url":         acctest.CHECKSET,
	}
}

var fakeServiceCatalogProductVersionMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"versions.#":         "0",
		"product_versions.#": "0",
	}
}

var ServiceCatalogProductVersionCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_service_catalog_product_versions.default",
	ExistMapFunc: existServiceCatalogProductVersionMapFunc,
	FakeMapFunc:  fakeServiceCatalogProductVersionMapFunc,
}

func testAccCheckAlicloudServiceCatalogProductVersionSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	return testAccServiceCatalogDataSourceFixtureConfig(rand) + fmt.Sprintf(`
data "alicloud_service_catalog_product_versions" "default" {
  product_id = alicloud_service_catalog_product.default.id
%s
  depends_on = [alicloud_service_catalog_principal_portfolio_association.default]
}
`, strings.Join(pairs, "\n   "))
}
