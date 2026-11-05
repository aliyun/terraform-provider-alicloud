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

func TestAccAliCloudServiceCatalogPortfolio_basic2593(t *testing.T) {
	var v map[string]interface{}
	acctest.CheckoutSupportedRegions(t, true, connectivity.ServiceCatalogProvisionedProductSupportRegions)
	resourceId := "alicloud_service_catalog_portfolio.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudServiceCatalogPortfolioMap2593)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return servicecatalog.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeServiceCatalogPortfolio")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sServiceCatalogPortfolio%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudServiceCatalogPortfolioBasicDependence2593)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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

// Test ServiceCatalog Portfolio. >>> Resource test cases, automatically generated.
// Case 恬裕测试 7495
func TestAccAliCloudServiceCatalogPortfolio_basic7495(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_service_catalog_portfolio.default"
	acctest.CheckoutSupportedRegions(t, true, connectivity.ServiceCatalogProvisionedProductSupportRegions)
	ra := acctest.ResourceAttrInit(resourceId, AlicloudServiceCatalogPortfolioMap7495)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return servicecatalog.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeServiceCatalogPortfolio")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sservicecatalogportfolio%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudServiceCatalogPortfolioBasicDependence7495)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"provider_name":  "tianyu",
					"description":    "desc",
					"portfolio_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"provider_name":  "tianyu",
						"description":    "desc",
						"portfolio_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"provider_name":  "tianyu1",
					"description":    "desc2",
					"portfolio_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"provider_name":  "tianyu1",
						"description":    "desc2",
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

var AlicloudServiceCatalogPortfolioMap7495 = map[string]string{
	"portfolio_arn": acctest.CHECKSET,
	"create_time":   acctest.CHECKSET,
}

func AlicloudServiceCatalogPortfolioBasicDependence7495(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 产品组合接入TF录入用例 5466
func TestAccAliCloudServiceCatalogPortfolio_basic5466(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_service_catalog_portfolio.default"
	acctest.CheckoutSupportedRegions(t, true, connectivity.ServiceCatalogProvisionedProductSupportRegions)
	ra := acctest.ResourceAttrInit(resourceId, AlicloudServiceCatalogPortfolioMap5466)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return servicecatalog.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeServiceCatalogPortfolio")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sservicecatalogportfolio%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudServiceCatalogPortfolioBasicDependence5466)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"provider_name":  "TF测试",
					"description":    "TF测试产品组合描述",
					"portfolio_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"provider_name":  "TF测试",
						"description":    "TF测试产品组合描述",
						"portfolio_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"provider_name":  "TF测试-修改后",
					"description":    "TF测试描述-修改后",
					"portfolio_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"provider_name":  "TF测试-修改后",
						"description":    "TF测试描述-修改后",
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

var AlicloudServiceCatalogPortfolioMap5466 = map[string]string{
	"portfolio_arn": acctest.CHECKSET,
	"create_time":   acctest.CHECKSET,
}

func AlicloudServiceCatalogPortfolioBasicDependence5466(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Test ServiceCatalog Portfolio. <<< Resource test cases, automatically generated.
