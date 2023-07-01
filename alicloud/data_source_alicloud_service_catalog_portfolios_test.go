package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudServiceCatalogPortfolioDataSource(t *testing.T) {
	checkoutSupportedRegions(t, true, connectivity.ServiceCatalogProvisionedProductSupportRegions)
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudServiceCatalogPortfolioSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_service_catalog_portfolio.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudServiceCatalogPortfolioSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_service_catalog_portfolio.default.id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudServiceCatalogPortfolioSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_service_catalog_portfolio.default.portfolio_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudServiceCatalogPortfolioSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_service_catalog_portfolio.default.portfolio_name}_fake"`,
		}),
	}
	scopeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudServiceCatalogPortfolioSourceConfig(rand, map[string]string{
			"ids":   `["${alicloud_service_catalog_portfolio.default.id}"]`,
			"scope": `"Local"`,
		}),
		fakeConfig: testAccCheckAlicloudServiceCatalogPortfolioSourceConfig(rand, map[string]string{
			"ids":   `["${alicloud_service_catalog_portfolio.default.id}_fake"]`,
			"scope": `"Import"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudServiceCatalogPortfolioSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_service_catalog_portfolio.default.id}"]`,
			"name_regex": `"${alicloud_service_catalog_portfolio.default.portfolio_name}"`,
			"scope":      `"Local"`,
		}),
		fakeConfig: testAccCheckAlicloudServiceCatalogPortfolioSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_service_catalog_portfolio.default.id}_fake"]`,
			"name_regex": `"${alicloud_service_catalog_portfolio.default.portfolio_name}_fake"`,
			"scope":      `"Import"`,
		}),
	}

	ServiceCatalogPortfolioCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameRegexConf, scopeConf, allConf)
}

var existServiceCatalogPortfolioMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                       "1",
		"names.#":                     "1",
		"portfolios.#":                "1",
		"portfolios.0.id":             CHECKSET,
		"portfolios.0.create_time":    CHECKSET,
		"portfolios.0.description":    CHECKSET,
		"portfolios.0.portfolio_arn":  CHECKSET,
		"portfolios.0.portfolio_id":   CHECKSET,
		"portfolios.0.portfolio_name": CHECKSET,
		"portfolios.0.provider_name":  CHECKSET,
	}
}

var fakeServiceCatalogPortfolioMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":        "0",
		"names.#":      "0",
		"portfolios.#": "0",
	}
}

var ServiceCatalogPortfolioCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_service_catalog_portfolios.default",
	existMapFunc: existServiceCatalogPortfolioMapFunc,
	fakeMapFunc:  fakeServiceCatalogPortfolioMapFunc,
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
