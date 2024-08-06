package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Fcv3 ConcurrencyConfig. >>> Resource test cases, automatically generated.
// Case ConcurrencyConfig_online 7303
func TestAccAliCloudFcv3ConcurrencyConfig_basic7303(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_fcv3_concurrency_config.default"
	ra := resourceAttrInit(resourceId, AlicloudFcv3ConcurrencyConfigMap7303)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &Fcv3ServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeFcv3ConcurrencyConfig")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sfcv3concurrencyconfig%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudFcv3ConcurrencyConfigBasicDependence7303)
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
					"function_name":        "${alicloud_fcv3_function.function.function_name}",
					"reserved_concurrency": "100",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"function_name":        CHECKSET,
						"reserved_concurrency": "100",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"reserved_concurrency": "200",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"reserved_concurrency": "200",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"function_name":        "${alicloud_fcv3_function.function.function_name}",
					"reserved_concurrency": "210",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"function_name":        CHECKSET,
						"reserved_concurrency": "210",
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

var AlicloudFcv3ConcurrencyConfigMap7303 = map[string]string{}

func AlicloudFcv3ConcurrencyConfigBasicDependence7303(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_fcv3_function" "function" {
  memory_size = "512"
  cpu         = 0.5
  handler     = "index.Handler"
  code {
    zip_file = "UEsDBBQACAAIAAAAAAAAAAAAAAAAAAAAAAAIAAAAaW5kZXgucHmEkEFKxEAQRfd9ig9ZTCJOooIwDMwNXLqXnnQlaalUhU5lRj2KZ/FOXkESGR114bJ/P/7jV4b1xRq1hijtFpM1682cuNgPmgysbRulPT0fRxXnMtwrSPyeCdYRokSLnuMLJTTkbUqEvDMbxm1VdcRD6Tk+T1LW2ldB66knsYdA5iNX17ebm6tN2VnPhcswMPmREPuBacb+CiapLarAj9gT6/H97dVlCNScY3mtYvRkxdZlwDKDEnanPWVLdrdkeXEGlFEazVdfPVHaVeHc3N15CUwppwOJXeK7HshAB8NuOU7J6sP4SRXuH/EvbUfMiqMmDqv5M5FNSfAj/wgAAP//UEsHCPl//NYAAQAArwEAAFBLAQIUABQACAAIAAAAAAD5f/zWAAEAAK8BAAAIAAAAAAAAAAAAAAAAAAAAAABpbmRleC5weVBLBQYAAAAAAQABADYAAAA2AQAAAAA="
  }
  function_name = var.name
  runtime       = "python3.9"
  disk_size     = "512"
  log_config {
    log_begin_rule = "None"
  }
}
`, name)
}

// Case ConcurrencyConfig 7180
func TestAccAliCloudFcv3ConcurrencyConfig_basic7180(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_fcv3_concurrency_config.default"
	ra := resourceAttrInit(resourceId, AlicloudFcv3ConcurrencyConfigMap7180)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &Fcv3ServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeFcv3ConcurrencyConfig")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sfcv3concurrencyconfig%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudFcv3ConcurrencyConfigBasicDependence7180)
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
					"function_name":        "${alicloud_fcv3_function.function.function_name}",
					"reserved_concurrency": "200",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"function_name":        CHECKSET,
						"reserved_concurrency": "200",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"reserved_concurrency": "100",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"reserved_concurrency": "100",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"reserved_concurrency": "100",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"reserved_concurrency": "100",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"function_name":        "${alicloud_fcv3_function.function.function_name}",
					"reserved_concurrency": "100",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"function_name":        CHECKSET,
						"reserved_concurrency": "100",
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

var AlicloudFcv3ConcurrencyConfigMap7180 = map[string]string{}

func AlicloudFcv3ConcurrencyConfigBasicDependence7180(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_fcv3_function" "function" {
  memory_size = "512"
  cpu         = 0.5
  handler     = "index.Handler"
  code {
    zip_file = "UEsDBBQACAAIAAAAAAAAAAAAAAAAAAAAAAAIAAAAaW5kZXgucHmEkEFKxEAQRfd9ig9ZTCJOooIwDMwNXLqXnnQlaalUhU5lRj2KZ/FOXkESGR114bJ/P/7jV4b1xRq1hijtFpM1682cuNgPmgysbRulPT0fRxXnMtwrSPyeCdYRokSLnuMLJTTkbUqEvDMbxm1VdcRD6Tk+T1LW2ldB66knsYdA5iNX17ebm6tN2VnPhcswMPmREPuBacb+CiapLarAj9gT6/H97dVlCNScY3mtYvRkxdZlwDKDEnanPWVLdrdkeXEGlFEazVdfPVHaVeHc3N15CUwppwOJXeK7HshAB8NuOU7J6sP4SRXuH/EvbUfMiqMmDqv5M5FNSfAj/wgAAP//UEsHCPl//NYAAQAArwEAAFBLAQIUABQACAAIAAAAAAD5f/zWAAEAAK8BAAAIAAAAAAAAAAAAAAAAAAAAAABpbmRleC5weVBLBQYAAAAAAQABADYAAAA2AQAAAAA="
  }
  function_name = var.name
  runtime       = "python3.9"
  disk_size     = "512"
  log_config {
    log_begin_rule = "None"
  }
}
`, name)
}

// Case ConcurrencyConfig_online 7303  twin
func TestAccAliCloudFcv3ConcurrencyConfig_basic7303_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_fcv3_concurrency_config.default"
	ra := resourceAttrInit(resourceId, AlicloudFcv3ConcurrencyConfigMap7303)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &Fcv3ServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeFcv3ConcurrencyConfig")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sfcv3concurrencyconfig%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudFcv3ConcurrencyConfigBasicDependence7303)
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
					"function_name":        "${alicloud_fcv3_function.function.function_name}",
					"reserved_concurrency": "200",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"function_name":        CHECKSET,
						"reserved_concurrency": "200",
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

// Case ConcurrencyConfig 7180  twin
func TestAccAliCloudFcv3ConcurrencyConfig_basic7180_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_fcv3_concurrency_config.default"
	ra := resourceAttrInit(resourceId, AlicloudFcv3ConcurrencyConfigMap7180)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &Fcv3ServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeFcv3ConcurrencyConfig")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sfcv3concurrencyconfig%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudFcv3ConcurrencyConfigBasicDependence7180)
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
					"function_name":        "${alicloud_fcv3_function.function.function_name}",
					"reserved_concurrency": "143",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"function_name":        CHECKSET,
						"reserved_concurrency": "143",
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

// Case ConcurrencyConfig_online 7303  raw
func TestAccAliCloudFcv3ConcurrencyConfig_basic7303_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_fcv3_concurrency_config.default"
	ra := resourceAttrInit(resourceId, AlicloudFcv3ConcurrencyConfigMap7303)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &Fcv3ServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeFcv3ConcurrencyConfig")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sfcv3concurrencyconfig%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudFcv3ConcurrencyConfigBasicDependence7303)
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
					"function_name":        "${alicloud_fcv3_function.function.function_name}",
					"reserved_concurrency": "200",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"function_name":        CHECKSET,
						"reserved_concurrency": "200",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"reserved_concurrency": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"reserved_concurrency": "2",
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

// Case ConcurrencyConfig 7180  raw
func TestAccAliCloudFcv3ConcurrencyConfig_basic7180_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_fcv3_concurrency_config.default"
	ra := resourceAttrInit(resourceId, AlicloudFcv3ConcurrencyConfigMap7180)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &Fcv3ServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeFcv3ConcurrencyConfig")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sfcv3concurrencyconfig%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudFcv3ConcurrencyConfigBasicDependence7180)
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
					"function_name":        "${alicloud_fcv3_function.function.function_name}",
					"reserved_concurrency": "200",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"function_name":        CHECKSET,
						"reserved_concurrency": "200",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"reserved_concurrency": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"reserved_concurrency": "2",
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

// Test Fcv3 ConcurrencyConfig. <<< Resource test cases, automatically generated.
