package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudServiceCatalogPortfolio_basic2593(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.ServiceCatalogProvisionedProductSupportRegions)
	resourceId := "alicloud_service_catalog_portfolio.default"
	ra := resourceAttrInit(resourceId, AlicloudServiceCatalogPortfolioMap2593)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ServicecatalogService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeServiceCatalogPortfolio")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sServiceCatalogPortfolio%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudServiceCatalogPortfolioBasicDependence2593)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"provider_name":  "${var.name}",
					"description":    "${var.name}",
					"portfolio_name": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"provider_name":  name,
						"description":    name,
						"portfolio_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"provider_name": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"provider_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"portfolio_name": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"portfolio_name": name + "_update",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudServiceCatalogPortfolioMap2593 = map[string]string{}

func AlicloudServiceCatalogPortfolioBasicDependence2593(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}
