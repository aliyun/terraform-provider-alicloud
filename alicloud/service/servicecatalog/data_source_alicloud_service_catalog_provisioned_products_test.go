package servicecatalog_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudServiceCatalogProvisionedProductDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(1000000, 9999999)
	acctest.CheckoutSupportedRegions(t, true, connectivity.ServiceCatalogProvisionedProductSupportRegions)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudServiceCatalogProvisionedProductSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_service_catalog_provisioned_product.default.id}"]`,
			"enable_details": `"true"`,
		}),
		FakeConfig: testAccCheckAlicloudServiceCatalogProvisionedProductSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_service_catalog_provisioned_product.default.id}_fake"]`,
		}),
	}

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudServiceCatalogProvisionedProductSourceConfig(rand, map[string]string{
			"name_regex":     `"${alicloud_service_catalog_provisioned_product.default.provisioned_product_name}"`,
			"enable_details": `"true"`,
		}),
		FakeConfig: testAccCheckAlicloudServiceCatalogProvisionedProductSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_service_catalog_provisioned_product.default.provisioned_product_name}_fake"`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudServiceCatalogProvisionedProductSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_service_catalog_provisioned_product.default.id}"]`,
			"name_regex":     `"${alicloud_service_catalog_provisioned_product.default.provisioned_product_name}"`,
			"enable_details": `"true"`,
		}),
		FakeConfig: testAccCheckAlicloudServiceCatalogProvisionedProductSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_service_catalog_provisioned_product.default.id}_fake"]`,
			"name_regex": `"${alicloud_service_catalog_provisioned_product.default.provisioned_product_name}_fake"`,
		}),
	}

	ServiceCatalogProvisionedProductCheckInfo.DataSourceTestCheck(t, rand, idsConf, nameRegexConf, allConf)
}

var existServiceCatalogProvisionedProductMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"products.#":                                                  "1",
		"products.0.id":                                               acctest.CHECKSET,
		"products.0.create_time":                                      acctest.CHECKSET,
		"products.0.last_provisioning_task_id":                        acctest.CHECKSET,
		"products.0.last_successful_provisioning_task_id":             acctest.CHECKSET,
		"products.0.last_task_id":                                     acctest.CHECKSET,
		"products.0.owner_principal_id":                               acctest.CHECKSET,
		"products.0.owner_principal_type":                             acctest.CHECKSET,
		"products.0.portfolio_id":                                     acctest.CHECKSET,
		"products.0.product_id":                                       acctest.CHECKSET,
		"products.0.product_name":                                     acctest.CHECKSET,
		"products.0.product_version_id":                               acctest.CHECKSET,
		"products.0.product_version_name":                             acctest.CHECKSET,
		"products.0.provisioned_product_arn":                          acctest.CHECKSET,
		"products.0.provisioned_product_id":                           acctest.CHECKSET,
		"products.0.provisioned_product_name":                         acctest.CHECKSET,
		"products.0.provisioned_product_type":                         acctest.CHECKSET,
		"products.0.stack_id":                                         acctest.CHECKSET,
		"products.0.stack_region_id":                                  acctest.CHECKSET,
		"products.0.status":                                           acctest.CHECKSET,
		"products.0.tags.%":                                           "1",
		"products.0.parameters.#":                                     "2",
		"provisioned_products.#":                                      "1",
		"provisioned_products.0.id":                                   acctest.CHECKSET,
		"provisioned_products.0.create_time":                          acctest.CHECKSET,
		"provisioned_products.0.last_provisioning_task_id":            acctest.CHECKSET,
		"provisioned_products.0.last_successful_provisioning_task_id": acctest.CHECKSET,
		"provisioned_products.0.last_task_id":                         acctest.CHECKSET,
		"provisioned_products.0.owner_principal_id":                   acctest.CHECKSET,
		"provisioned_products.0.owner_principal_type":                 acctest.CHECKSET,
		"provisioned_products.0.portfolio_id":                         acctest.CHECKSET,
		"provisioned_products.0.product_id":                           acctest.CHECKSET,
		"provisioned_products.0.product_name":                         acctest.CHECKSET,
		"provisioned_products.0.product_version_id":                   acctest.CHECKSET,
		"provisioned_products.0.product_version_name":                 acctest.CHECKSET,
		"provisioned_products.0.provisioned_product_arn":              acctest.CHECKSET,
		"provisioned_products.0.provisioned_product_id":               acctest.CHECKSET,
		"provisioned_products.0.provisioned_product_name":             acctest.CHECKSET,
		"provisioned_products.0.provisioned_product_type":             acctest.CHECKSET,
		"provisioned_products.0.stack_id":                             acctest.CHECKSET,
		"provisioned_products.0.stack_region_id":                      acctest.CHECKSET,
		"provisioned_products.0.status":                               acctest.CHECKSET,
		"provisioned_products.0.tags.%":                               "1",
		"provisioned_products.0.parameters.#":                         "2",
	}
}

var fakeServiceCatalogProvisionedProductMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"products.#":             "0",
		"provisioned_products.#": "0",
	}
}

var ServiceCatalogProvisionedProductCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_service_catalog_provisioned_products.default",
	ExistMapFunc: existServiceCatalogProvisionedProductMapFunc,
	FakeMapFunc:  fakeServiceCatalogProvisionedProductMapFunc,
}

func testAccCheckAlicloudServiceCatalogProvisionedProductSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	return testAccServiceCatalogDataSourceFixtureConfig(rand) + fmt.Sprintf(`
data "alicloud_service_catalog_launch_options" "default" {
  product_id = alicloud_service_catalog_product.default.id
  depends_on = [alicloud_service_catalog_principal_portfolio_association.default]
}

resource "alicloud_service_catalog_provisioned_product" "default" {
  provisioned_product_name = var.name
  stack_region_id          = "cn-hangzhou"
  product_version_id       = alicloud_service_catalog_product_version.default.id
  product_id               = alicloud_service_catalog_product.default.id
  portfolio_id             = data.alicloud_service_catalog_launch_options.default.launch_options.0.portfolio_id
  tags = {
    "v1" = "tf-test"
  }
  parameters {
    parameter_key   = "acl"
    parameter_value = "private"
  }
  parameters {
    parameter_key   = "bucket_name"
    parameter_value = lower(var.name)
  }
}

data "alicloud_service_catalog_provisioned_products" "default" {
%s
  depends_on = [alicloud_service_catalog_provisioned_product.default]
}
`, strings.Join(pairs, "\n   "))
}
