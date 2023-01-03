package alicloud

import (
	"fmt"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudServiceCatalogProvisionedProductDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	checkoutSupportedRegions(t, true, connectivity.ServiceCatalogProvisionedProductSupportRegions)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudServiceCatalogProvisionedProductSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_service_catalog_provisioned_product.default.id}"]`,
			"enable_details": `"true"`,
		}),
		fakeConfig: testAccCheckAlicloudServiceCatalogProvisionedProductSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_service_catalog_provisioned_product.default.id}_fake"]`,
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudServiceCatalogProvisionedProductSourceConfig(rand, map[string]string{
			"name_regex":     `"${alicloud_service_catalog_provisioned_product.default.provisioned_product_name}"`,
			"enable_details": `"true"`,
		}),
		fakeConfig: testAccCheckAlicloudServiceCatalogProvisionedProductSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_service_catalog_provisioned_product.default.provisioned_product_name}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudServiceCatalogProvisionedProductSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_service_catalog_provisioned_product.default.id}"]`,
			"name_regex":     `"${alicloud_service_catalog_provisioned_product.default.provisioned_product_name}"`,
			"enable_details": `"true"`,
		}),
		fakeConfig: testAccCheckAlicloudServiceCatalogProvisionedProductSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_service_catalog_provisioned_product.default.id}_fake"]`,
			"name_regex": `"${alicloud_service_catalog_provisioned_product.default.provisioned_product_name}_fake"`,
		}),
	}

	ServiceCatalogProvisionedProductCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameRegexConf, allConf)
}

var existServiceCatalogProvisionedProductMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"products.#":                                      "1",
		"products.0.id":                                   CHECKSET,
		"products.0.create_time":                          CHECKSET,
		"products.0.last_provisioning_task_id":            CHECKSET,
		"products.0.last_successful_provisioning_task_id": CHECKSET,
		"products.0.last_task_id":                         CHECKSET,
		"products.0.owner_principal_id":                   CHECKSET,
		"products.0.owner_principal_type":                 CHECKSET,
		"products.0.portfolio_id":                         CHECKSET,
		"products.0.product_id":                           CHECKSET,
		"products.0.product_name":                         CHECKSET,
		"products.0.product_version_id":                   CHECKSET,
		"products.0.product_version_name":                 CHECKSET,
		"products.0.provisioned_product_arn":              CHECKSET,
		"products.0.provisioned_product_id":               CHECKSET,
		"products.0.provisioned_product_name":             CHECKSET,
		"products.0.provisioned_product_type":             CHECKSET,
		"products.0.stack_id":                             CHECKSET,
		"products.0.stack_region_id":                      CHECKSET,
		"products.0.status":                               CHECKSET,
		"products.0.tags.%":                               "1",
		"products.0.parameters.#":                         "1",
	}
}

var fakeServiceCatalogProvisionedProductMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"products.#": "0",
	}
}

var ServiceCatalogProvisionedProductCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_service_catalog_provisioned_products.default",
	existMapFunc: existServiceCatalogProvisionedProductMapFunc,
	fakeMapFunc:  fakeServiceCatalogProvisionedProductMapFunc,
}

func testAccCheckAlicloudServiceCatalogProvisionedProductSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccServiceCatalogProvisionedProduct%d"
}

data "alicloud_service_catalog_product_as_end_users" "default" {
  name_regex = "ram模板创建"
}

data "alicloud_service_catalog_product_versions" "default" {
  name_regex = "1.0.0"
  product_id = data.alicloud_service_catalog_product_as_end_users.default.users.0.id
}

data "alicloud_service_catalog_launch_options" "default" {
  product_id = data.alicloud_service_catalog_product_as_end_users.default.users.0.id
}

data "alicloud_ros_regions" "all" {}

resource "alicloud_service_catalog_provisioned_product" "default" {
  provisioned_product_name = var.name
  stack_region_id          = data.alicloud_ros_regions.all.regions.5.region_id
  product_version_id       = data.alicloud_service_catalog_product_versions.default.versions.0.id
  product_id               = data.alicloud_service_catalog_product_as_end_users.default.users.0.id
  portfolio_id             = data.alicloud_service_catalog_launch_options.default.options.0.portfolio_id
  tags = {
    "v1" = "tf-test"
  }
  parameters {
    parameter_key   = "role_name"
    parameter_value = var.name
  }
}

data "alicloud_service_catalog_provisioned_products" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
