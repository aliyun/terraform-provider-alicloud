package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Fcv3 AsyncInvokeConfig. >>> Resource test cases, automatically generated.
// Case AsyncInvokeConfig_Base_Online 7336
func TestAccAliCloudFcv3AsyncInvokeConfig_basic7336(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_fcv3_async_invoke_config.default"
	ra := resourceAttrInit(resourceId, AliCloudFcv3AsyncInvokeConfigMap7336)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &Fcv3ServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeFcv3AsyncInvokeConfig")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sFcv3asyncinvokeconfig%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudFcv3AsyncInvokeConfigBasicDependence7336)
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
					"function_name": "${alicloud_fcv3_function.default.function_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"function_name": CHECKSET,
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
					"destination_config": []map[string]interface{}{
						{
							"on_success": []map[string]interface{}{
								{
									"destination": "acs:fc:${data.alicloud_regions.default.regions.0.id}:${data.alicloud_account.default.id}:functions/${alicloud_fcv3_function.success.function_name}",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"destination_config.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"destination_config": []map[string]interface{}{
						{
							"on_success": []map[string]interface{}{
								{
									"destination": "acs:fc:${data.alicloud_regions.default.regions.0.id}:${data.alicloud_account.default.id}:functions/${alicloud_fcv3_function.success.function_name}",
								},
							},
							"on_failure": []map[string]interface{}{
								{
									"destination": "acs:fc:${data.alicloud_regions.default.regions.0.id}:${data.alicloud_account.default.id}:functions/${alicloud_fcv3_function.failure.function_name}",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"destination_config.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"destination_config": []map[string]interface{}{
						{
							"on_success": []map[string]interface{}{
								{
									"destination": "acs:fc:${data.alicloud_regions.default.regions.0.id}:${data.alicloud_account.default.id}:functions/${alicloud_fcv3_function.failure.function_name}",
								},
							},
							"on_failure": []map[string]interface{}{
								{
									"destination": "acs:fc:${data.alicloud_regions.default.regions.0.id}:${data.alicloud_account.default.id}:functions/${alicloud_fcv3_function.success.function_name}",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"destination_config.#": "1",
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

func TestAccAliCloudFcv3AsyncInvokeConfig_basic7336_1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_fcv3_async_invoke_config.default"
	ra := resourceAttrInit(resourceId, AliCloudFcv3AsyncInvokeConfigMap7336)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &Fcv3ServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeFcv3AsyncInvokeConfig")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sFcv3asyncinvokeconfig%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudFcv3AsyncInvokeConfigBasicDependence7336)
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
					"function_name": "${alicloud_fcv3_function.default.function_name}",
					"destination_config": []map[string]interface{}{
						{
							"on_success": []map[string]interface{}{
								{
									"destination": "acs:fc:${data.alicloud_regions.default.regions.0.id}:${data.alicloud_account.default.id}:functions/${alicloud_fcv3_function.success.function_name}",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"function_name":        CHECKSET,
						"destination_config.#": "1",
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
					"destination_config": []map[string]interface{}{
						{
							"on_success": []map[string]interface{}{
								{
									"destination": "acs:fc:${data.alicloud_regions.default.regions.0.id}:${data.alicloud_account.default.id}:functions/${alicloud_fcv3_function.success.function_name}",
								},
							},
							"on_failure": []map[string]interface{}{
								{
									"destination": "acs:fc:${data.alicloud_regions.default.regions.0.id}:${data.alicloud_account.default.id}:functions/${alicloud_fcv3_function.failure.function_name}",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"destination_config.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"destination_config": []map[string]interface{}{
						{
							"on_success": []map[string]interface{}{
								{
									"destination": "acs:fc:${data.alicloud_regions.default.regions.0.id}:${data.alicloud_account.default.id}:functions/${alicloud_fcv3_function.failure.function_name}",
								},
							},
							"on_failure": []map[string]interface{}{
								{
									"destination": "acs:fc:${data.alicloud_regions.default.regions.0.id}:${data.alicloud_account.default.id}:functions/${alicloud_fcv3_function.success.function_name}",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"destination_config.#": "1",
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

func TestAccAliCloudFcv3AsyncInvokeConfig_basic7336_2(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_fcv3_async_invoke_config.default"
	ra := resourceAttrInit(resourceId, AliCloudFcv3AsyncInvokeConfigMap7336)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &Fcv3ServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeFcv3AsyncInvokeConfig")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sFcv3asyncinvokeconfig%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudFcv3AsyncInvokeConfigBasicDependence7336)
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
					"function_name": "${alicloud_fcv3_function.default.function_name}",
					"destination_config": []map[string]interface{}{
						{
							"on_failure": []map[string]interface{}{
								{
									"destination": "acs:fc:${data.alicloud_regions.default.regions.0.id}:${data.alicloud_account.default.id}:functions/${alicloud_fcv3_function.failure.function_name}",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"function_name":        CHECKSET,
						"destination_config.#": "1",
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
					"destination_config": []map[string]interface{}{
						{
							"on_success": []map[string]interface{}{
								{
									"destination": "acs:fc:${data.alicloud_regions.default.regions.0.id}:${data.alicloud_account.default.id}:functions/${alicloud_fcv3_function.success.function_name}",
								},
							},
							"on_failure": []map[string]interface{}{
								{
									"destination": "acs:fc:${data.alicloud_regions.default.regions.0.id}:${data.alicloud_account.default.id}:functions/${alicloud_fcv3_function.failure.function_name}",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"destination_config.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"destination_config": []map[string]interface{}{
						{
							"on_success": []map[string]interface{}{
								{
									"destination": "acs:fc:${data.alicloud_regions.default.regions.0.id}:${data.alicloud_account.default.id}:functions/${alicloud_fcv3_function.failure.function_name}",
								},
							},
							"on_failure": []map[string]interface{}{
								{
									"destination": "acs:fc:${data.alicloud_regions.default.regions.0.id}:${data.alicloud_account.default.id}:functions/${alicloud_fcv3_function.success.function_name}",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"destination_config.#": "1",
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

func TestAccAliCloudFcv3AsyncInvokeConfig_basic7336_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_fcv3_async_invoke_config.default"
	ra := resourceAttrInit(resourceId, AliCloudFcv3AsyncInvokeConfigMap7336)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &Fcv3ServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeFcv3AsyncInvokeConfig")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sFcv3asyncinvokeconfig%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudFcv3AsyncInvokeConfigBasicDependence7336)
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
					"function_name":                  "${alicloud_fcv3_function.default.function_name}",
					"async_task":                     "true",
					"max_async_event_age_in_seconds": "2",
					"max_async_retry_attempts":       "2",
					"qualifier":                      "LATEST",
					"destination_config": []map[string]interface{}{
						{
							"on_success": []map[string]interface{}{
								{
									"destination": "acs:fc:${data.alicloud_regions.default.regions.0.id}:${data.alicloud_account.default.id}:functions/${alicloud_fcv3_function.success.function_name}",
								},
							},
							"on_failure": []map[string]interface{}{
								{
									"destination": "acs:fc:${data.alicloud_regions.default.regions.0.id}:${data.alicloud_account.default.id}:functions/${alicloud_fcv3_function.failure.function_name}",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"function_name":                  CHECKSET,
						"async_task":                     "true",
						"max_async_event_age_in_seconds": "2",
						"max_async_retry_attempts":       "2",
						"qualifier":                      "LATEST",
						"destination_config.#":           "1",
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

var AliCloudFcv3AsyncInvokeConfigMap7336 = map[string]string{
	"create_time":        CHECKSET,
	"function_arn":       CHECKSET,
	"last_modified_time": CHECKSET,
}

func AliCloudFcv3AsyncInvokeConfigBasicDependence7336(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_regions" "default" {
  current = true
}

data "alicloud_account" "default" {
}

resource "alicloud_fcv3_function" "default" {
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

resource "alicloud_fcv3_function" "success" {
  memory_size = "512"
  cpu         = 0.5
  handler     = "index.Handler"
  code {
    zip_file = "UEsDBBQACAAIAAAAAAAAAAAAAAAAAAAAAAAIAAAAaW5kZXgucHmEkEFKxEAQRfd9ig9ZTCJOooIwDMwNXLqXnnQlaalUhU5lRj2KZ/FOXkESGR114bJ/P/7jV4b1xRq1hijtFpM1682cuNgPmgysbRulPT0fRxXnMtwrSPyeCdYRokSLnuMLJTTkbUqEvDMbxm1VdcRD6Tk+T1LW2ldB66knsYdA5iNX17ebm6tN2VnPhcswMPmREPuBacb+CiapLarAj9gT6/H97dVlCNScY3mtYvRkxdZlwDKDEnanPWVLdrdkeXEGlFEazVdfPVHaVeHc3N15CUwppwOJXeK7HshAB8NuOU7J6sP4SRXuH/EvbUfMiqMmDqv5M5FNSfAj/wgAAP//UEsHCPl//NYAAQAArwEAAFBLAQIUABQACAAIAAAAAAD5f/zWAAEAAK8BAAAIAAAAAAAAAAAAAAAAAAAAAABpbmRleC5weVBLBQYAAAAAAQABADYAAAA2AQAAAAA="
  }
  function_name = format("%%s_%%s", var.name, "success")
  runtime       = "python3.9"
  disk_size     = "512"
  log_config {
    log_begin_rule = "None"
  }
}

resource "alicloud_fcv3_function" "failure" {
  memory_size = "512"
  cpu         = 0.5
  handler     = "index.Handler"
  code {
    zip_file = "UEsDBBQACAAIAAAAAAAAAAAAAAAAAAAAAAAIAAAAaW5kZXgucHmEkEFKxEAQRfd9ig9ZTCJOooIwDMwNXLqXnnQlaalUhU5lRj2KZ/FOXkESGR114bJ/P/7jV4b1xRq1hijtFpM1682cuNgPmgysbRulPT0fRxXnMtwrSPyeCdYRokSLnuMLJTTkbUqEvDMbxm1VdcRD6Tk+T1LW2ldB66knsYdA5iNX17ebm6tN2VnPhcswMPmREPuBacb+CiapLarAj9gT6/H97dVlCNScY3mtYvRkxdZlwDKDEnanPWVLdrdkeXEGlFEazVdfPVHaVeHc3N15CUwppwOJXeK7HshAB8NuOU7J6sP4SRXuH/EvbUfMiqMmDqv5M5FNSfAj/wgAAP//UEsHCPl//NYAAQAArwEAAFBLAQIUABQACAAIAAAAAAD5f/zWAAEAAK8BAAAIAAAAAAAAAAAAAAAAAAAAAABpbmRleC5weVBLBQYAAAAAAQABADYAAAA2AQAAAAA="
  }
  function_name = format("%%s_%%s", var.name, "failure")
  runtime       = "python3.9"
  disk_size     = "512"
  log_config {
    log_begin_rule = "None"
  }
}
`, name)
}

// Test Fcv3 AsyncInvokeConfig. <<< Resource test cases, automatically generated.
