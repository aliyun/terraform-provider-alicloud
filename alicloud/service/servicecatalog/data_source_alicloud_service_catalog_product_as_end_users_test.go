package servicecatalog_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudServiceCatalogProductAsEndUserDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(1000000, 9999999)
	acctest.CheckoutSupportedRegions(t, true, connectivity.ServiceCatalogProvisionedProductSupportRegions)

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudServiceCatalogProductAsEndUserSourceConfig(rand, map[string]string{
			"name_regex": `var.name`,
		}),
		FakeConfig: testAccCheckAlicloudServiceCatalogProductAsEndUserSourceConfig(rand, map[string]string{
			"name_regex": `format("%s_fake", var.name)`,
		}),
	}

	ServiceCatalogProductAsEndUserCheckInfo.DataSourceTestCheck(t, rand, nameRegexConf)
}

var existServiceCatalogProductAsEndUserMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"users.#":                           "1",
		"users.0.id":                        acctest.CHECKSET,
		"users.0.create_time":               acctest.CHECKSET,
		"users.0.has_default_launch_option": acctest.CHECKSET,
		"users.0.product_arn":               acctest.CHECKSET,
		"users.0.product_id":                acctest.CHECKSET,
		"users.0.product_name":              acctest.CHECKSET,
		"users.0.product_type":              acctest.CHECKSET,
		"users.0.provider_name":             acctest.CHECKSET,
	}
}

var fakeServiceCatalogProductAsEndUserMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"users.#": "0",
	}
}

var ServiceCatalogProductAsEndUserCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_service_catalog_product_as_end_users.default",
	ExistMapFunc: existServiceCatalogProductAsEndUserMapFunc,
	FakeMapFunc:  fakeServiceCatalogProductAsEndUserMapFunc,
}

func testAccCheckAlicloudServiceCatalogProductAsEndUserSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	return testAccServiceCatalogDataSourceFixtureConfig(rand) + fmt.Sprintf(`
data "alicloud_service_catalog_product_as_end_users" "default" {
%s
  depends_on = [alicloud_service_catalog_principal_portfolio_association.default]
}
`, strings.Join(pairs, "\n   "))
}
