package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Case AsyncInvokeConfig_Base 7133
func TestAccAliCloudFcv3AsyncInvokeConfig_basic7133(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_fcv3_async_invoke_config.default"
	ra := resourceAttrInit(resourceId, AlicloudFcv3AsyncInvokeConfigMap7133)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &Fcv3ServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeFcv3AsyncInvokeConfig")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sFcv3asyncinvokeconfig%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudFcv3AsyncInvokeConfigBasicDependence7133)
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
					"function_name": "${alicloud_fcv3_function.function.function_name}",
					"destination_config": []map[string]interface{}{
						{
							"on_failure": []map[string]interface{}{
								{
									"destination": "acs:fc:${data.alicloud_regions.current_regions.regions.0.id}:${data.alicloud_account.current.id}:functions/${alicloud_fcv3_function.function1.function_name}",
								},
							},
							"on_success": []map[string]interface{}{
								{
									"destination": "acs:fc:${data.alicloud_regions.current_regions.regions.0.id}:${data.alicloud_account.current.id}:functions/${alicloud_fcv3_function.function1.function_name}",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"function_name": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"max_async_retry_attempts": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"max_async_retry_attempts": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"max_async_event_age_in_seconds": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"max_async_event_age_in_seconds": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"async_task": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"async_task": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"destination_config": []map[string]interface{}{
						{
							"on_failure": []map[string]interface{}{
								{
									"destination": "acs:fc:${data.alicloud_regions.current_regions.regions.0.id}:${data.alicloud_account.current.id}:functions/${alicloud_fcv3_function.function2.function_name}",
								},
							},
							"on_success": []map[string]interface{}{
								{
									"destination": "acs:fc:${data.alicloud_regions.current_regions.regions.0.id}:${data.alicloud_account.current.id}:functions/${alicloud_fcv3_function.function2.function_name}",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"max_async_retry_attempts": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"max_async_retry_attempts": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"max_async_event_age_in_seconds": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"max_async_event_age_in_seconds": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"async_task": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"async_task": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"function_name": "${alicloud_fcv3_function.function.function_name}",
					"destination_config": []map[string]interface{}{
						{
							"on_failure": []map[string]interface{}{
								{
									"destination": "acs:fc:${data.alicloud_regions.current_regions.regions.0.id}:${data.alicloud_account.current.id}:functions/${alicloud_fcv3_function.function1.function_name}",
								},
							},
							"on_success": []map[string]interface{}{
								{
									"destination": "acs:fc:${data.alicloud_regions.current_regions.regions.0.id}:${data.alicloud_account.current.id}:functions/${alicloud_fcv3_function.function2.function_name}",
								},
							},
						},
					},
					"qualifier":                      "LATEST",
					"max_async_retry_attempts":       "1",
					"max_async_event_age_in_seconds": "1",
					"async_task":                     "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"function_name":                  CHECKSET,
						"qualifier":                      "LATEST",
						"max_async_retry_attempts":       "1",
						"max_async_event_age_in_seconds": "1",
						"async_task":                     "false",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"qualifier"},
			},
		},
	})
}

var AlicloudFcv3AsyncInvokeConfigMap7133 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudFcv3AsyncInvokeConfigBasicDependence7133(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_regions" "current_regions" {
  current = true
}

data "alicloud_account" "current" {
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

resource "alicloud_fcv3_function" "function1" {
  memory_size = "512"
  cpu         = 0.5
  handler     = "index.Handler"
  code {
    zip_file = "UEsDBBQACAAIAAAAAAAAAAAAAAAAAAAAAAAIAAAAaW5kZXgucHmEkEFKxEAQRfd9ig9ZTCJOooIwDMwNXLqXnnQlaalUhU5lRj2KZ/FOXkESGR114bJ/P/7jV4b1xRq1hijtFpM1682cuNgPmgysbRulPT0fRxXnMtwrSPyeCdYRokSLnuMLJTTkbUqEvDMbxm1VdcRD6Tk+T1LW2ldB66knsYdA5iNX17ebm6tN2VnPhcswMPmREPuBacb+CiapLarAj9gT6/H97dVlCNScY3mtYvRkxdZlwDKDEnanPWVLdrdkeXEGlFEazVdfPVHaVeHc3N15CUwppwOJXeK7HshAB8NuOU7J6sP4SRXuH/EvbUfMiqMmDqv5M5FNSfAj/wgAAP//UEsHCPl//NYAAQAArwEAAFBLAQIUABQACAAIAAAAAAD5f/zWAAEAAK8BAAAIAAAAAAAAAAAAAAAAAAAAAABpbmRleC5weVBLBQYAAAAAAQABADYAAAA2AQAAAAA="
  }
  function_name = format("%%s_%%s", var.name, "update1")
  runtime       = "python3.9"
  disk_size     = "512"
  log_config {
    log_begin_rule = "None"
  }
}

resource "alicloud_fcv3_function" "function2" {
  memory_size = "512"
  cpu         = 0.5
  handler     = "index.Handler"
  code {
    zip_file = "UEsDBBQACAAIAAAAAAAAAAAAAAAAAAAAAAAIAAAAaW5kZXgucHmEkEFKxEAQRfd9ig9ZTCJOooIwDMwNXLqXnnQlaalUhU5lRj2KZ/FOXkESGR114bJ/P/7jV4b1xRq1hijtFpM1682cuNgPmgysbRulPT0fRxXnMtwrSPyeCdYRokSLnuMLJTTkbUqEvDMbxm1VdcRD6Tk+T1LW2ldB66knsYdA5iNX17ebm6tN2VnPhcswMPmREPuBacb+CiapLarAj9gT6/H97dVlCNScY3mtYvRkxdZlwDKDEnanPWVLdrdkeXEGlFEazVdfPVHaVeHc3N15CUwppwOJXeK7HshAB8NuOU7J6sP4SRXuH/EvbUfMiqMmDqv5M5FNSfAj/wgAAP//UEsHCPl//NYAAQAArwEAAFBLAQIUABQACAAIAAAAAAD5f/zWAAEAAK8BAAAIAAAAAAAAAAAAAAAAAAAAAABpbmRleC5weVBLBQYAAAAAAQABADYAAAA2AQAAAAA="
  }
  function_name = format("%%s_%%s", var.name, "update2")
  runtime       = "python3.9"
  disk_size     = "512"
  log_config {
    log_begin_rule = "None"
  }
}

`, name)
}

// Case AsyncInvokeConfig_Base 7133  twin
func TestAccAliCloudFcv3AsyncInvokeConfig_basic7133_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_fcv3_async_invoke_config.default"
	ra := resourceAttrInit(resourceId, AlicloudFcv3AsyncInvokeConfigMap7133)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &Fcv3ServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeFcv3AsyncInvokeConfig")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sFcv3asyncinvokeconfig%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudFcv3AsyncInvokeConfigBasicDependence7133)
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
					"function_name": "${alicloud_fcv3_function.function.function_name}",
					"destination_config": []map[string]interface{}{
						{
							"on_failure": []map[string]interface{}{
								{
									"destination": "acs:fc:${data.alicloud_regions.current_regions.regions.0.id}:${data.alicloud_account.current.id}:functions/${alicloud_fcv3_function.function1.function_name}",
								},
							},
							"on_success": []map[string]interface{}{
								{
									"destination": "acs:fc:${data.alicloud_regions.current_regions.regions.0.id}:${data.alicloud_account.current.id}:functions/${alicloud_fcv3_function.function2.function_name}",
								},
							},
						},
					},
					"qualifier":                      "LATEST",
					"max_async_retry_attempts":       "1",
					"max_async_event_age_in_seconds": "1",
					"async_task":                     "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"function_name":                  CHECKSET,
						"qualifier":                      "LATEST",
						"max_async_retry_attempts":       "1",
						"max_async_event_age_in_seconds": "1",
						"async_task":                     "true",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"qualifier"},
			},
		},
	})
}

// Case AsyncInvokeConfig_Base 7133  raw
func TestAccAliCloudFcv3AsyncInvokeConfig_basic7133_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_fcv3_async_invoke_config.default"
	ra := resourceAttrInit(resourceId, AlicloudFcv3AsyncInvokeConfigMap7133)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &Fcv3ServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeFcv3AsyncInvokeConfig")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sfcv3asyncinvokeconfig%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudFcv3AsyncInvokeConfigBasicDependence7133)
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
					"function_name": "${alicloud_fcv3_function.function.function_name}",
					"destination_config": []map[string]interface{}{
						{
							"on_failure": []map[string]interface{}{
								{
									"destination": "acs:fc:${data.alicloud_regions.current_regions.regions.0.id}:${data.alicloud_account.current.id}:functions/${alicloud_fcv3_function.function1.function_name}",
								},
							},
							"on_success": []map[string]interface{}{
								{
									"destination": "acs:fc:${data.alicloud_regions.current_regions.regions.0.id}:${data.alicloud_account.current.id}:functions/${alicloud_fcv3_function.function1.function_name}",
								},
							},
						},
					},
					"qualifier":                      "LATEST",
					"max_async_retry_attempts":       "1",
					"max_async_event_age_in_seconds": "1",
					"async_task":                     "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"function_name":                  CHECKSET,
						"qualifier":                      "LATEST",
						"max_async_retry_attempts":       "1",
						"max_async_event_age_in_seconds": "1",
						"async_task":                     "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"destination_config": []map[string]interface{}{
						{
							"on_failure": []map[string]interface{}{
								{
									"destination": "acs:fc:${data.alicloud_regions.current_regions.regions.0.id}:${data.alicloud_account.current.id}:functions/${alicloud_fcv3_function.function2.function_name}",
								},
							},
							"on_success": []map[string]interface{}{
								{
									"destination": "acs:fc:${data.alicloud_regions.current_regions.regions.0.id}:${data.alicloud_account.current.id}:functions/${alicloud_fcv3_function.function2.function_name}",
								},
							},
						},
					},
					"max_async_retry_attempts":       "2",
					"max_async_event_age_in_seconds": "2",
					"async_task":                     "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"max_async_retry_attempts":       "2",
						"max_async_event_age_in_seconds": "2",
						"async_task":                     "true",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"qualifier"},
			},
		},
	})
}
