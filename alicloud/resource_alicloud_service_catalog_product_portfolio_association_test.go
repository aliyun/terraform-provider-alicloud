package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test ServiceCatalog ProductPortfolioAssociation. >>> Resource test cases, automatically generated.
// Case 贝熊测试 7473
func TestAccAliCloudServiceCatalogProductPortfolioAssociation_basic7473(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.ServiceCatalogProvisionedProductSupportRegions)
	resourceId := "alicloud_service_catalog_product_portfolio_association.default"
	ra := resourceAttrInit(resourceId, AlicloudServiceCatalogProductPortfolioAssociationMap7473)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ServiceCatalogServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeServiceCatalogProductPortfolioAssociation")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sservicecatalogproductportfolioassociation%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudServiceCatalogProductPortfolioAssociationBasicDependence7473)
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
					"portfolio_id": "${alicloud_service_catalog_portfolio.default0yAgJ8.id}",
					"product_id":   "${alicloud_service_catalog_product.defaultRetBJw.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"portfolio_id": CHECKSET,
						"product_id":   CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
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

var AlicloudServiceCatalogProductPortfolioAssociationMap7473 = map[string]string{}

func AlicloudServiceCatalogProductPortfolioAssociationBasicDependence7473(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_service_catalog_portfolio" "default0yAgJ8" {
  provider_name  = "tianyu"
  description    = "desc"
  portfolio_name = var.name
}

resource "alicloud_service_catalog_product" "defaultRetBJw" {
  provider_name = "tianyu"
  product_name  = format("%%s1", var.name)
  product_type  = "Ros"
}


`, name)
}

// Test ServiceCatalog ProductPortfolioAssociation. <<< Resource test cases, automatically generated.
