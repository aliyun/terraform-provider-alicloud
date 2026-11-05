package servicecatalog_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudServiceCatalogPortfolioDataSource(t *testing.T) {
	acctest.CheckoutSupportedRegions(t, true, connectivity.ServiceCatalogProvisionedProductSupportRegions)
	rand := tfacctest.RandIntRange(1000000, 9999999)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudServiceCatalogPortfolioSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_service_catalog_portfolio.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudServiceCatalogPortfolioSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_service_catalog_portfolio.default.id}_fake"]`,
		}),
	}
	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudServiceCatalogPortfolioSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_service_catalog_portfolio.default.portfolio_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudServiceCatalogPortfolioSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_service_catalog_portfolio.default.portfolio_name}_fake"`,
		}),
	}
	scopeConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudServiceCatalogPortfolioSourceConfig(rand, map[string]string{
			"ids":   `["${alicloud_service_catalog_portfolio.default.id}"]`,
			"scope": `"Local"`,
		}),
		FakeConfig: testAccCheckAlicloudServiceCatalogPortfolioSourceConfig(rand, map[string]string{
			"ids":   `["${alicloud_service_catalog_portfolio.default.id}_fake"]`,
			"scope": `"Import"`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudServiceCatalogPortfolioSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_service_catalog_portfolio.default.id}"]`,
			"name_regex": `"${alicloud_service_catalog_portfolio.default.portfolio_name}"`,
			"scope":      `"Local"`,
		}),
		FakeConfig: testAccCheckAlicloudServiceCatalogPortfolioSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_service_catalog_portfolio.default.id}_fake"]`,
			"name_regex": `"${alicloud_service_catalog_portfolio.default.portfolio_name}_fake"`,
			"scope":      `"Import"`,
		}),
	}

	ServiceCatalogPortfolioCheckInfo.DataSourceTestCheck(t, rand, idsConf, nameRegexConf, scopeConf, allConf)
}

var existServiceCatalogPortfolioMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                       "1",
		"names.#":                     "1",
		"portfolios.#":                "1",
		"portfolios.0.id":             acctest.CHECKSET,
		"portfolios.0.create_time":    acctest.CHECKSET,
		"portfolios.0.description":    acctest.CHECKSET,
		"portfolios.0.portfolio_arn":  acctest.CHECKSET,
		"portfolios.0.portfolio_id":   acctest.CHECKSET,
		"portfolios.0.portfolio_name": acctest.CHECKSET,
		"portfolios.0.provider_name":  acctest.CHECKSET,
	}
}

var fakeServiceCatalogPortfolioMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":        "0",
		"names.#":      "0",
		"portfolios.#": "0",
	}
}

var ServiceCatalogPortfolioCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_service_catalog_portfolios.default",
	ExistMapFunc: existServiceCatalogPortfolioMapFunc,
	FakeMapFunc:  fakeServiceCatalogPortfolioMapFunc,
}

func testAccCheckAlicloudServiceCatalogPortfolioSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccServiceCatalogPortfolio%d"
}

resource "alicloud_service_catalog_portfolio" "default" {
	provider_name = var.name
	portfolio_name = var.name
	description = var.name
}

data "alicloud_service_catalog_portfolios" "default" {
	sort_by = "CreateTime"
	sort_order = "Asc"
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
