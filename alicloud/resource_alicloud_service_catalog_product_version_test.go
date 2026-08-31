package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test ServiceCatalog ProductVersion. >>> Resource test cases, automatically generated.
// Case 产品版本测试 7448
func TestAccAliCloudServiceCatalogProductVersion_basic7448(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_service_catalog_product_version.default"
	checkoutSupportedRegions(t, true, connectivity.ServiceCatalogProvisionedProductSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudServiceCatalogProductVersionMap7448)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ServiceCatalogServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeServiceCatalogProductVersion")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sservicecatalogproductversion%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudServiceCatalogProductVersionBasicDependence7448)
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
					"guidance":             "Default",
					"template_url":         "oss://servicecatalog-cn-hangzhou/1466115886172051/terraform/template/tpl-bp1x4v3r44u7u7/template.json",
					"active":               "true",
					"description":          "产品版本测试",
					"product_version_name": name,
					"product_id":           "${alicloud_service_catalog_product.defaultmaeTcE.id}",
					"template_type":        "RosTerraformTemplate",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"guidance":             "Default",
						"template_url":         "oss://servicecatalog-cn-hangzhou/1466115886172051/terraform/template/tpl-bp1x4v3r44u7u7/template.json",
						"active":               "true",
						"description":          "产品版本测试",
						"product_version_name": name,
						"product_id":           CHECKSET,
						"template_type":        "RosTerraformTemplate",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"guidance":             "GuidanceAfter",
					"active":               "false",
					"description":          "产品版本测试-修改后",
					"product_version_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"guidance":             "GuidanceAfter",
						"active":               "false",
						"description":          "产品版本测试-修改后",
						"product_version_name": name + "_update",
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

var AlicloudServiceCatalogProductVersionMap7448 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudServiceCatalogProductVersionBasicDependence7448(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_service_catalog_product" "defaultmaeTcE" {
  provider_name = "贝熊"
  product_name  = var.name
  product_type  = "Ros"
}


`, name)
}

// Case 产品版本录入 5655
func TestAccAliCloudServiceCatalogProductVersion_basic5655(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_service_catalog_product_version.default"
	checkoutSupportedRegions(t, true, connectivity.ServiceCatalogProvisionedProductSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudServiceCatalogProductVersionMap5655)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ServiceCatalogServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeServiceCatalogProductVersion")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sservicecatalogproductversion%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudServiceCatalogProductVersionBasicDependence5655)
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
					"guidance":             "Default",
					"template_url":         "oss://servicecatalog-cn-hangzhou/1466115886172051/terraform/template/tpl-bp1c23hk4fv1bq/template.json",
					"description":          "产品版本测试",
					"product_version_name": name,
					"product_id":           "${alicloud_service_catalog_product.defaulthATkAc.id}",
					"template_type":        "RosTerraformTemplate",
					"active":               "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"guidance":             "Default",
						"template_url":         "oss://servicecatalog-cn-hangzhou/1466115886172051/terraform/template/tpl-bp1c23hk4fv1bq/template.json",
						"description":          "产品版本测试",
						"product_version_name": name,
						"product_id":           CHECKSET,
						"template_type":        "RosTerraformTemplate",
						"active":               "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"guidance":             "Recommended",
					"description":          "产品版本测试-修改后",
					"product_version_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"guidance":             "Recommended",
						"description":          "产品版本测试-修改后",
						"product_version_name": name + "_update",
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

var AlicloudServiceCatalogProductVersionMap5655 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudServiceCatalogProductVersionBasicDependence5655(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_service_catalog_product" "defaulthATkAc" {
  provider_name = "tf"
  description   = "tf-测试产品描述"
  product_name  = var.name
  product_type  = "Ros"
}


`, name)
}

// Test ServiceCatalog ProductVersion. <<< Resource test cases, automatically generated.
