package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudServiceCatalogPortfolio_basic2593(t *testing.T) {
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

// Test ServiceCatalog Portfolio. >>> Resource test cases, automatically generated.
// Case 恬裕测试 7495
func TestAccAliCloudServiceCatalogPortfolio_basic7495(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_service_catalog_portfolio.default"
	checkoutSupportedRegions(t, true, connectivity.ServiceCatalogProvisionedProductSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudServiceCatalogPortfolioMap7495)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ServiceCatalogServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeServiceCatalogPortfolio")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sservicecatalogportfolio%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudServiceCatalogPortfolioBasicDependence7495)
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
	"portfolio_arn": CHECKSET,
	"create_time":   CHECKSET,
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
	checkoutSupportedRegions(t, true, connectivity.ServiceCatalogProvisionedProductSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudServiceCatalogPortfolioMap5466)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ServiceCatalogServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeServiceCatalogPortfolio")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sservicecatalogportfolio%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudServiceCatalogPortfolioBasicDependence5466)
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
	"portfolio_arn": CHECKSET,
	"create_time":   CHECKSET,
}

func AlicloudServiceCatalogPortfolioBasicDependence5466(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Test ServiceCatalog Portfolio. <<< Resource test cases, automatically generated.
