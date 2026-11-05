package servicecatalog_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudServiceCatalogLaunchOptionDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(1000000, 9999999)
	acctest.CheckoutSupportedRegions(t, true, connectivity.ServiceCatalogProvisionedProductSupportRegions)

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudServiceCatalogLaunchOptionSourceConfig(rand, map[string]string{
			"name_regex": `var.name`,
		}),
		FakeConfig: testAccCheckAlicloudServiceCatalogLaunchOptionSourceConfig(rand, map[string]string{
			"name_regex": `format("%s_fake", var.name)`,
		}),
	}

	ServiceCatalogLaunchOptionCheckInfo.DataSourceTestCheck(t, rand, allConf)
}

var existServiceCatalogLaunchOptionMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"options.#":                       "1",
		"options.0.id":                    acctest.CHECKSET,
		"options.0.portfolio_id":          acctest.CHECKSET,
		"options.0.portfolio_name":        acctest.CHECKSET,
		"launch_options.#":                "1",
		"launch_options.0.id":             acctest.CHECKSET,
		"launch_options.0.portfolio_id":   acctest.CHECKSET,
		"launch_options.0.portfolio_name": acctest.CHECKSET,
	}
}

var fakeServiceCatalogLaunchOptionMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"options.#":        "0",
		"launch_options.#": "0",
	}
}

var ServiceCatalogLaunchOptionCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_service_catalog_launch_options.default",
	ExistMapFunc: existServiceCatalogLaunchOptionMapFunc,
	FakeMapFunc:  fakeServiceCatalogLaunchOptionMapFunc,
}

func testAccCheckAlicloudServiceCatalogLaunchOptionSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	return testAccServiceCatalogDataSourceFixtureConfig(rand) + fmt.Sprintf(`
data "alicloud_service_catalog_launch_options" "default" {
  product_id = alicloud_service_catalog_product.default.id
%s
  depends_on = [alicloud_service_catalog_principal_portfolio_association.default]
}
`, strings.Join(pairs, "\n   "))
}
