package servicecatalog_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudServiceCatalogEndUserProductDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(1000000, 9999999)
	acctest.CheckoutSupportedRegions(t, true, connectivity.ServiceCatalogProvisionedProductSupportRegions)

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudServiceCatalogEndUserProductSourceConfig(rand, map[string]string{
			"name_regex": `var.name`,
			"sort_order": `"Asc"`,
			"sort_by":    `"CreateTime"`,
		}),
		FakeConfig: testAccCheckAlicloudServiceCatalogEndUserProductSourceConfig(rand, map[string]string{
			"name_regex": `format("%s_fake", var.name)`,
			"sort_order": `"Asc"`,
			"sort_by":    `"CreateTime"`,
		}),
	}

	ServiceCatalogEndUserProductCheckInfo.DataSourceTestCheck(t, rand, nameRegexConf)
}

var existServiceCatalogEndUserProductMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                           "1",
		"end_user_products.#":             "1",
		"end_user_products.0.id":          acctest.CHECKSET,
		"end_user_products.0.create_time": acctest.CHECKSET,
		"end_user_products.0.has_default_launch_option": acctest.CHECKSET,
		"end_user_products.0.product_arn":               acctest.CHECKSET,
		"end_user_products.0.product_id":                acctest.CHECKSET,
		"end_user_products.0.product_name":              acctest.CHECKSET,
		"end_user_products.0.product_type":              acctest.CHECKSET,
		"end_user_products.0.provider_name":             acctest.CHECKSET,
	}
}

var fakeServiceCatalogEndUserProductMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":               "0",
		"end_user_products.#": "0",
	}
}

var ServiceCatalogEndUserProductCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_service_catalog_end_user_products.default",
	ExistMapFunc: existServiceCatalogEndUserProductMapFunc,
	FakeMapFunc:  fakeServiceCatalogEndUserProductMapFunc,
}

// testAccServiceCatalogDataSourceFixtureConfig builds the shared fixture chain used by the
// servicecatalog end-user data source acceptance tests. It creates a portfolio, a product, a
// product version (using the service-owned ROS template), a product-portfolio association and a
// principal-portfolio association that grants the *caller* access to the portfolio, so that
// ListProductsAsEndUser/ListLaunchOptions return the created product for the test account.
func testAccServiceCatalogDataSourceFixtureConfig(rand int) string {
	return fmt.Sprintf(`
data "alicloud_caller_identity" "current" {}

locals {
  principal_type = data.alicloud_caller_identity.current.identity_type == "ram-role" ? "RamRole" : "RamUser"
}

variable "name" {
  default = "tf-testAccServiceCatalogDS%d"
}

resource "alicloud_service_catalog_portfolio" "default" {
  provider_name  = "tianyu"
  portfolio_name = var.name
}

resource "alicloud_service_catalog_product" "default" {
  provider_name = "tianyu"
  product_name  = var.name
  product_type  = "Ros"
}

resource "alicloud_service_catalog_product_version" "default" {
  product_id           = alicloud_service_catalog_product.default.id
  template_url         = "oss://servicecatalog-cn-hangzhou/1466115886172051/terraform/template/tpl-bp1x4v3r44u7u7/template.json"
  template_type        = "RosTerraformTemplate"
  product_version_name = "1.0.0"
  active               = true
  guidance             = "Default"
}

resource "alicloud_service_catalog_product_portfolio_association" "default" {
  portfolio_id = alicloud_service_catalog_portfolio.default.id
  product_id    = alicloud_service_catalog_product.default.id
}

resource "alicloud_service_catalog_principal_portfolio_association" "default" {
  portfolio_id   = alicloud_service_catalog_portfolio.default.id
  principal_id   = data.alicloud_caller_identity.current.id
  principal_type = local.principal_type
  depends_on = [
    alicloud_service_catalog_product_portfolio_association.default,
    alicloud_service_catalog_product_version.default,
  ]
}
`, rand)
}

func testAccCheckAlicloudServiceCatalogEndUserProductSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	return testAccServiceCatalogDataSourceFixtureConfig(rand) + fmt.Sprintf(`
data "alicloud_service_catalog_end_user_products" "default" {
%s
  depends_on = [alicloud_service_catalog_principal_portfolio_association.default]
}
`, strings.Join(pairs, "\n   "))
}
