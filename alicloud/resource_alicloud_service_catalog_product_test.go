package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test ServiceCatalog Product. >>> Resource test cases, automatically generated.
// Case 恬裕测试 7497
func TestAccAliCloudServiceCatalogProduct_basic7497(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_service_catalog_product.default"
	checkoutSupportedRegions(t, true, connectivity.ServiceCatalogProvisionedProductSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudServiceCatalogProductMap7497)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ServiceCatalogServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeServiceCatalogProduct")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sservicecatalogproduct%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudServiceCatalogProductBasicDependence7497)
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
					"provider_name": "tianyu",
					"product_name":  name,
					"product_type":  "Ros",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"provider_name": "tianyu",
						"product_name":  name,
						"product_type":  "Ros",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "desc",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "desc",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"provider_name": "tianyu2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"provider_name": "tianyu2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "desc2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "desc2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"product_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"product_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"provider_name": "tianyu",
					"description":   "desc",
					"product_name":  name + "_update",
					"product_type":  "Ros",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"provider_name": "tianyu",
						"description":   "desc",
						"product_name":  name + "_update",
						"product_type":  "Ros",
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

var AlicloudServiceCatalogProductMap7497 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudServiceCatalogProductBasicDependence7497(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 测试创建/更新资源_副本1702974127376 5469
func TestAccAliCloudServiceCatalogProduct_basic5469(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_service_catalog_product.default"
	checkoutSupportedRegions(t, true, connectivity.ServiceCatalogProvisionedProductSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudServiceCatalogProductMap5469)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ServiceCatalogServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeServiceCatalogProduct")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sservicecatalogproduct%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudServiceCatalogProductBasicDependence5469)
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
					"provider_name": "TF测试",
					"product_name":  name,
					"product_type":  "Ros",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"provider_name": "TF测试",
						"product_name":  name,
						"product_type":  "Ros",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "TF测试描述",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "TF测试描述",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"provider_name": "TF测试-修改后",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"provider_name": "TF测试-修改后",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"product_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"product_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "TF测试描述-修改后",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "TF测试描述-修改后",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"provider_name": "TF测试",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"provider_name": "TF测试",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"product_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"product_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "TF测试描述",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "TF测试描述",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"provider_name": "TF测试",
					"product_name":  name + "_update",
					"product_type":  "Ros",
					"description":   "TF测试描述",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"provider_name": "TF测试",
						"product_name":  name + "_update",
						"product_type":  "Ros",
						"description":   "TF测试描述",
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

var AlicloudServiceCatalogProductMap5469 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudServiceCatalogProductBasicDependence5469(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 恬裕测试 7497  twin
func TestAccAliCloudServiceCatalogProduct_basic7497_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_service_catalog_product.default"
	checkoutSupportedRegions(t, true, connectivity.ServiceCatalogProvisionedProductSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudServiceCatalogProductMap7497)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ServiceCatalogServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeServiceCatalogProduct")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sservicecatalogproduct%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudServiceCatalogProductBasicDependence7497)
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
					"provider_name": "tianyu",
					"description":   "desc",
					"product_name":  name,
					"product_type":  "Ros",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"provider_name": "tianyu",
						"description":   "desc",
						"product_name":  name,
						"product_type":  "Ros",
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

// Case 测试创建/更新资源_副本1702974127376 5469  twin
func TestAccAliCloudServiceCatalogProduct_basic5469_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_service_catalog_product.default"
	checkoutSupportedRegions(t, true, connectivity.ServiceCatalogProvisionedProductSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudServiceCatalogProductMap5469)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ServiceCatalogServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeServiceCatalogProduct")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sservicecatalogproduct%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudServiceCatalogProductBasicDependence5469)
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
					"provider_name": "TF测试",
					"product_name":  name,
					"product_type":  "Ros",
					"description":   "TF测试描述",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"provider_name": "TF测试",
						"product_name":  name,
						"product_type":  "Ros",
						"description":   "TF测试描述",
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

// Case 恬裕测试 7497  raw
func TestAccAliCloudServiceCatalogProduct_basic7497_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_service_catalog_product.default"
	checkoutSupportedRegions(t, true, connectivity.ServiceCatalogProvisionedProductSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudServiceCatalogProductMap7497)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ServiceCatalogServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeServiceCatalogProduct")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sservicecatalogproduct%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudServiceCatalogProductBasicDependence7497)
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
					"provider_name": "tianyu",
					"description":   "desc",
					"product_name":  name,
					"product_type":  "Ros",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"provider_name": "tianyu",
						"description":   "desc",
						"product_name":  name,
						"product_type":  "Ros",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"provider_name": "tianyu2",
					"description":   "desc2",
					"product_name":  name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"provider_name": "tianyu2",
						"description":   "desc2",
						"product_name":  name + "_update",
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

// Case 测试创建/更新资源_副本1702974127376 5469  raw
func TestAccAliCloudServiceCatalogProduct_basic5469_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_service_catalog_product.default"
	checkoutSupportedRegions(t, true, connectivity.ServiceCatalogProvisionedProductSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudServiceCatalogProductMap5469)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ServiceCatalogServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeServiceCatalogProduct")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sservicecatalogproduct%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudServiceCatalogProductBasicDependence5469)
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
					"provider_name": "TF测试",
					"product_name":  name,
					"product_type":  "Ros",
					"description":   "TF测试描述",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"provider_name": "TF测试",
						"product_name":  name,
						"product_type":  "Ros",
						"description":   "TF测试描述",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"provider_name": "TF测试-修改后",
					"product_name":  name + "_update",
					"description":   "TF测试描述-修改后",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"provider_name": "TF测试-修改后",
						"product_name":  name + "_update",
						"description":   "TF测试描述-修改后",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"provider_name": "TF测试",
					"product_name":  name + "_update",
					"description":   "TF测试描述",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"provider_name": "TF测试",
						"product_name":  name + "_update",
						"description":   "TF测试描述",
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

// Test ServiceCatalog Product. <<< Resource test cases, automatically generated.
