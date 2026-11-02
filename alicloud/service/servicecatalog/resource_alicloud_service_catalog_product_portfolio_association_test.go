package servicecatalog_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/servicecatalog"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test ServiceCatalog ProductPortfolioAssociation. >>> Resource test cases, automatically generated.
// Case 贝熊测试 7473
func TestAccAliCloudServiceCatalogProductPortfolioAssociation_basic7473(t *testing.T) {
	var v map[string]interface{}
	acctest.CheckoutSupportedRegions(t, true, connectivity.ServiceCatalogProvisionedProductSupportRegions)
	resourceId := "alicloud_service_catalog_product_portfolio_association.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudServiceCatalogProductPortfolioAssociationMap7473)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return servicecatalog.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeServiceCatalogProductPortfolioAssociation")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sservicecatalogproductportfolioassociation%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudServiceCatalogProductPortfolioAssociationBasicDependence7473)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"portfolio_id": "${alicloud_service_catalog_portfolio.default0yAgJ8.id}",
					"product_id":   "${alicloud_service_catalog_product.defaultRetBJw.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"portfolio_id": acctest.CHECKSET,
						"product_id":   acctest.CHECKSET,
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
