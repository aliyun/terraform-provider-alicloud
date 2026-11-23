package fcv3_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/fcv3"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/sls"
	"github.com/hashicorp/go-cty/cty"
	schema "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdktf "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// Test Fcv3 ProvisionConfig. >>> Resource test cases, automatically generated.
// Case ProvisionConfig_Base 7182
func TestAccAliCloudFcv3ProvisionConfig_basic7182(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_fcv3_provision_config.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudFcv3ProvisionConfigMap7182)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return fcv3.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeProvisionConfig")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sfcv3provisionconfig%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudFcv3ProvisionConfigBasicDependence7182)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy: resource.ComposeTestCheckFunc(
			rac.CheckResourceDestroy(),
			checkProvisionConfigDependenciesDestroyed(),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"function_name":       "${alicloud_fcv3_function.function.function_name}",
					"always_allocate_cpu": "true",
					"target":              "1",
					"target_tracking_policies": []map[string]interface{}{
						{
							"name":          "t1",
							"start_time":    "2030-10-10T10:10:10Z",
							"end_time":      "2035-10-10T10:10:10Z",
							"min_capacity":  "0",
							"max_capacity":  "1",
							"metric_target": "1",
							"metric_type":   "ProvisionedConcurrencyUtilization",
						},
						{
							"name":          "t2",
							"start_time":    "2030-10-10T10:10:10Z",
							"end_time":      "2035-10-10T10:10:10Z",
							"min_capacity":  "0",
							"max_capacity":  "1",
							"metric_target": "1",
							"metric_type":   "ProvisionedConcurrencyUtilization",
						},
						{
							"name":          "t3",
							"start_time":    "2030-10-10T10:10:10",
							"end_time":      "2035-10-10T10:10:10",
							"min_capacity":  "0",
							"max_capacity":  "1",
							"metric_target": "1",
							"metric_type":   "ProvisionedConcurrencyUtilization",
							"time_zone":     "Asia/Shanghai",
						},
					},
					"scheduled_actions": []map[string]interface{}{
						{
							"name":                "s1",
							"start_time":          "2030-10-10T10:10:10Z",
							"end_time":            "2035-10-10T10:10:10Z",
							"schedule_expression": "cron(0 0 4 * * *)",
							"target":              "0",
						},
						{
							"name":                "s2",
							"start_time":          "2030-10-10T10:10:10Z",
							"end_time":            "2035-10-10T10:10:10Z",
							"schedule_expression": "cron(0 0 6 * * *)",
							"target":              "1",
						},
						{
							"name":                "s3",
							"start_time":          "2030-10-10T10:10:10",
							"end_time":            "2035-10-10T10:10:10",
							"schedule_expression": "cron(0 0 7 * * *)",
							"target":              "0",
							"time_zone":           "Asia/Shanghai",
						},
					},
					"qualifier":           "LATEST",
					"always_allocate_gpu": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"function_name":              acctest.CHECKSET,
						"always_allocate_cpu":        "true",
						"target":                     "1",
						"target_tracking_policies.#": "3",
						"scheduled_actions.#":        "3",
						"qualifier":                  "LATEST",
						"always_allocate_gpu":        "true",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"always_allocate_cpu": "false",
					"target":              "2",
					"target_tracking_policies": []map[string]interface{}{
						{
							"name":          "t1",
							"start_time":    "2031-10-10T10:10:10",
							"end_time":      "2036-10-10T10:10:10",
							"min_capacity":  "1",
							"max_capacity":  "2",
							"metric_target": "0.6",
							"metric_type":   "ProvisionedConcurrencyUtilization",
							"time_zone":     "Asia/Singapore",
						},
					},
					"scheduled_actions": []map[string]interface{}{
						{
							"name":                "s1",
							"start_time":          "2031-10-10T10:10:10",
							"end_time":            "2036-10-10T10:10:10",
							"schedule_expression": "cron(0 0 3 * * *)",
							"target":              "2",
							"time_zone":           "Asia/Singapore",
						},
					},
					"always_allocate_gpu": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"always_allocate_cpu":                      "false",
						"target":                                   "2",
						"target_tracking_policies.#":               "1",
						"target_tracking_policies.0.metric_type":   "ProvisionedConcurrencyUtilization",
						"target_tracking_policies.0.metric_target": "0.6",
						"scheduled_actions.#":                      "1",
						"always_allocate_gpu":                      "false",
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

var AlicloudFcv3ProvisionConfigMap7182 = map[string]string{}

func AlicloudFcv3ProvisionConfigBasicDependence7182(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_log_project" "default" {
  project_name = lower(var.name)
  description  = var.name
}

resource "alicloud_log_store" "default" {
  project_name          = alicloud_log_project.default.name
  logstore_name         = lower(var.name)
  shard_count           = 3
  auto_split            = true
  max_split_shard_count = 60
  append_meta           = true
}

resource "alicloud_fcv3_function" "function" {
  memory_size   = "512"
  cpu           = 0.5
  handler       = "index.handler"
  function_name = var.name
  runtime       = "python3.10"
  disk_size     = "512"
  code {
    zip_file = "UEsDBBQACAAIAAAAAAAAAAAAAAAAAAAAAAAIAAAAaW5kZXgucHmEkEFKxEAQRfd9ig9ZTCJOooIwDMwNXLqXnnQlaalUhU5lRj2KZ/FOXkESGR114bJ/P/7jV4b1xRq1hijtFpM1682cuNgPmgysbRulPT0fRxXnMtwrSPyeCdYRokSLnuMLJTTkbUqEvDMbxm1VdcRD6Tk+T1LW2ldB66knsYdA5iNX17ebm6tN2VnPhcswMPmREPuBacb+CiapLarAj9gT6/H97dVlCNScY3mtYvRkxdZlwDKDEnanPWVLdrdkeXEGlFEazVdfPVHaVeHc3N15CUwppwOJXeK7HshAB8NuOU7J6sP4SRXuH/EvbUfMiqMmDqv5M5FNSfAj/wgAAP//UEsHCPl//NYAAQAArwEAAFBLAQIUABQACAAIAAAAAAD5f/zWAAEAAK8BAAAIAAAAAAAAAAAAAAAAAAAAAABpbmRleC5weVBLBQYAAAAAAQABADYAAAA2AQAAAAA="
  }
  log_config {
    enable_instance_metrics = true
    enable_request_metrics  = true
    project                 = alicloud_log_project.default.project_name
    logstore                = alicloud_log_store.default.logstore_name
    log_begin_rule          = "None"
  }
}


`, name)
}

// Test Fcv3 ProvisionConfig. <<< Resource test cases, automatically generated.

// Case ProvisionConfig_Basic 7183
// Hand-written representative ACC covering the account-supported basic
// ProvisionConfig lifecycle (create/read/update/import/destroy).
// Step 1 is the minimal legal configuration (function_name and target only).
// Later steps exercise the Optional+Computed
// always_allocate_cpu/always_allocate_gpu fields: explicit true, explicit false
// while target stays unchanged, both false, back to true, and finally removal.
// Every successful configuration step is followed by ImportStateVerify.
func TestAccAliCloudFcv3ProvisionConfig_basic7183(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_fcv3_provision_config.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudFcv3ProvisionConfigMap7183)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return fcv3.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeProvisionConfig")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sfcv3provisionconfig%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudFcv3ProvisionConfigBasicDependence7183)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy: resource.ComposeTestCheckFunc(
			rac.CheckResourceDestroy(),
			checkProvisionConfigDependenciesDestroyed(),
		),
		Steps: []resource.TestStep{
			// The minimal legal configuration sets only function_name and target;
			// the optional always-allocate fields stay unset.
			{
				Config: testAccConfig(map[string]interface{}{
					"function_name": "${alicloud_fcv3_function.function.function_name}",
					"target":        "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"function_name":       acctest.CHECKSET,
						"always_allocate_cpu": acctest.CHECKSET,
						"always_allocate_gpu": acctest.CHECKSET,
						"target":              "1",
					}),
				),
			},
			// current/current_error change asynchronously while the reserved
			// instances become ready, so they are not compared on import.
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"current", "current_error"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"always_allocate_cpu": "true",
					"always_allocate_gpu": "true",
					"target":              "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"always_allocate_cpu": "true",
						"always_allocate_gpu": "true",
						"target":              "2",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"current", "current_error"},
			},
			// Explicit false must be stored and converge. Only the two
			// booleans change here, so target is unchanged and still must be
			// sent. The service keeps always_allocate_gpu consistent with
			// always_allocate_cpu, so both values move together (observed on a
			// function without GPU: true/false became true/true and false/true
			// became false/false).
			//
			// target is removed from the configuration in the same step: the
			// previous apply and refresh left 2 in the state, so the
			// replacement PUT must carry 2 and not the zero value of an unset
			// argument.
			{
				Config: testAccConfig(map[string]interface{}{
					"always_allocate_cpu": "false",
					"always_allocate_gpu": "false",
					"target":              acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"always_allocate_cpu": "false",
						"always_allocate_gpu": "false",
						"target":              "2",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"current", "current_error"},
			},
			// Switching back to true must take effect and converge again.
			{
				Config: testAccConfig(map[string]interface{}{
					"always_allocate_cpu": "true",
					"always_allocate_gpu": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"always_allocate_cpu": "true",
						"always_allocate_gpu": "true",
						"target":              "2",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"current", "current_error"},
			},
			// Removing all arguments keeps the service-side values because the
			// fields are Optional+Computed and must not produce a diff.
			{
				Config: testAccConfig(map[string]interface{}{
					"always_allocate_cpu": acctest.REMOVEKEY,
					"always_allocate_gpu": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"always_allocate_cpu": "true",
						"always_allocate_gpu": "true",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"current", "current_error"},
			},
		},
	})
}

var AlicloudFcv3ProvisionConfigMap7183 = map[string]string{}

// AlicloudFcv3ProvisionConfigBasicDependence7183 creates a minimal fcv3
// function without an SLS log_project/log_store (the function carries a
// log_config with only log_begin_rule, which does not require a separate SLS
// project). The provision-config lifecycle does not need function logging, so
// this keeps the fixture minimal and avoids depending on the account's SLS
// project quota.
func AlicloudFcv3ProvisionConfigBasicDependence7183(name string) string {
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

// TestAccAliCloudFcv3ProvisionConfig_createWithoutTarget records the write
// contract behind the Optional+Computed schema of target: PutProvisionConfig
// declares target as required, so creating the resource without it must be
// rejected by the service instead of silently writing the zero value. target
// may only be omitted for a resource that already exists (import or a removed
// argument), where the next update recovers the authoritative default from the
// same qualifier. The second step creates the same resource with target set,
// which proves the fixture is
// healthy and the first step failed only because target was missing.
func TestAccAliCloudFcv3ProvisionConfig_createWithoutTarget(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_fcv3_provision_config.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudFcv3ProvisionConfigMap7183)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return fcv3.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeProvisionConfig")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sfcv3provisionnotarget%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudFcv3ProvisionConfigBasicDependence7183)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy: resource.ComposeTestCheckFunc(
			rac.CheckResourceDestroy(),
			checkProvisionConfigDependenciesDestroyed(),
		),
		Steps: []resource.TestStep{
			// The service answers 400 InvalidArgument "Target is required but
			// not provided"; match that exact message so an unrelated create
			// failure cannot make this step pass.
			{
				Config: testAccConfig(map[string]interface{}{
					"function_name": "${alicloud_fcv3_function.function.function_name}",
				}),
				ExpectError: regexp.MustCompile("Target is required but not provided"),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"function_name": "${alicloud_fcv3_function.function.function_name}",
					"target":        "1",
					"qualifier":     "LATEST",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"target":    "1",
						"qualifier": "LATEST",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"current", "current_error"},
			},
		},
	})
}

// TestAccAliCloudFcv3ProvisionConfig_memberOrder covers convergent reordering of
// the configurable TypeList attributes target_tracking_policies and
// scheduled_actions: apply the members in one order, then apply a permutation.
func TestAccAliCloudFcv3ProvisionConfig_memberOrder(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_fcv3_provision_config.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudFcv3ProvisionConfigMap7183)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return fcv3.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeProvisionConfig")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sfcv3provisionorder%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudFcv3ProvisionConfigBasicDependence7182)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy: resource.ComposeTestCheckFunc(
			rac.CheckResourceDestroy(),
			checkProvisionConfigDependenciesDestroyed(),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"function_name": "${alicloud_fcv3_function.function.function_name}",
					"target":        "1",
					"target_tracking_policies": []map[string]interface{}{
						{
							"name":          "t1",
							"start_time":    "2030-10-10T10:10:10Z",
							"end_time":      "2035-10-10T10:10:10Z",
							"min_capacity":  "0",
							"max_capacity":  "1",
							"metric_target": "1",
							"metric_type":   "ProvisionedConcurrencyUtilization",
						},
						{
							"name":          "t2",
							"start_time":    "2030-10-10T10:10:10Z",
							"end_time":      "2035-10-10T10:10:10Z",
							"min_capacity":  "0",
							"max_capacity":  "1",
							"metric_target": "1",
							"metric_type":   "ProvisionedConcurrencyUtilization",
						},
					},
					"scheduled_actions": []map[string]interface{}{
						{
							"name":                "s1",
							"start_time":          "2030-10-10T10:10:10Z",
							"end_time":            "2035-10-10T10:10:10Z",
							"schedule_expression": "cron(0 0 4 * * *)",
							"target":              "0",
						},
						{
							"name":                "s2",
							"start_time":          "2030-10-10T10:10:10Z",
							"end_time":            "2035-10-10T10:10:10Z",
							"schedule_expression": "cron(0 0 6 * * *)",
							"target":              "1",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"target_tracking_policies.#": "2",
						"scheduled_actions.#":        "2",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"current", "current_error"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"target_tracking_policies": []map[string]interface{}{
						{
							"name":          "t2",
							"start_time":    "2030-10-10T10:10:10Z",
							"end_time":      "2035-10-10T10:10:10Z",
							"min_capacity":  "0",
							"max_capacity":  "1",
							"metric_target": "1",
							"metric_type":   "ProvisionedConcurrencyUtilization",
						},
						{
							"name":          "t1",
							"start_time":    "2030-10-10T10:10:10Z",
							"end_time":      "2035-10-10T10:10:10Z",
							"min_capacity":  "0",
							"max_capacity":  "1",
							"metric_target": "1",
							"metric_type":   "ProvisionedConcurrencyUtilization",
						},
					},
				}),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"current", "current_error"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scheduled_actions": []map[string]interface{}{
						{
							"name":                "s2",
							"start_time":          "2030-10-10T10:10:10Z",
							"end_time":            "2035-10-10T10:10:10Z",
							"schedule_expression": "cron(0 0 6 * * *)",
							"target":              "1",
						},
						{
							"name":                "s1",
							"start_time":          "2030-10-10T10:10:10Z",
							"end_time":            "2035-10-10T10:10:10Z",
							"schedule_expression": "cron(0 0 4 * * *)",
							"target":              "0",
						},
					},
				}),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"current", "current_error"},
			},
			// Removing a TypeList must clear it through the replacement
			// PUT and converge on an empty list. The configuration helper is
			// cumulative, so REMOVEKEY is required to drop the block from the
			// final HCL.
			{
				Config: testAccConfig(map[string]interface{}{
					"target_tracking_policies": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"target_tracking_policies.#": "0",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"current", "current_error"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scheduled_actions": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scheduled_actions.#": "0",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"current", "current_error"},
			},
		},
	})
}

// TestAccAliCloudFcv3ProvisionConfig_memberOrderPair keeps the two TypeList
// reorder applies adjacent, as required by scripts/collection-order/README.md
// "Required acceptance steps", which defines the reorder evidence as an apply
// of configuration A followed directly by an apply that only reorders the
// target list. The checker resets the pair on any intervening step that does
// not apply the configuration (ImportState included), so an ImportStateVerify
// between the two applies would drop the coverage. Per-step Import for the
// same lists is carried by TestAccAliCloudFcv3ProvisionConfig_memberOrder,
// which places an Import after every successful configuration step; this
// function supplies the gate-shaped convergence pair and verifies the final
// reordered state by ImportStateVerify.
func TestAccAliCloudFcv3ProvisionConfig_memberOrderPair(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_fcv3_provision_config.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudFcv3ProvisionConfigMap7183)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return fcv3.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeProvisionConfig")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sfcv3provisionpair%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudFcv3ProvisionConfigBasicDependence7182)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy: resource.ComposeTestCheckFunc(
			rac.CheckResourceDestroy(),
			checkProvisionConfigDependenciesDestroyed(),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"function_name": "${alicloud_fcv3_function.function.function_name}",
					"target":        "1",
					"target_tracking_policies": []map[string]interface{}{
						{
							"name":          "t1",
							"start_time":    "2030-10-10T10:10:10Z",
							"end_time":      "2035-10-10T10:10:10Z",
							"min_capacity":  "0",
							"max_capacity":  "1",
							"metric_target": "1",
							"metric_type":   "ProvisionedConcurrencyUtilization",
						},
						{
							"name":          "t2",
							"start_time":    "2030-10-10T10:10:10Z",
							"end_time":      "2035-10-10T10:10:10Z",
							"min_capacity":  "0",
							"max_capacity":  "1",
							"metric_target": "1",
							"metric_type":   "ProvisionedConcurrencyUtilization",
						},
					},
					"scheduled_actions": []map[string]interface{}{
						{
							"name":                "s1",
							"start_time":          "2030-10-10T10:10:10Z",
							"end_time":            "2035-10-10T10:10:10Z",
							"schedule_expression": "cron(0 0 4 * * *)",
							"target":              "0",
						},
						{
							"name":                "s2",
							"start_time":          "2030-10-10T10:10:10Z",
							"end_time":            "2035-10-10T10:10:10Z",
							"schedule_expression": "cron(0 0 6 * * *)",
							"target":              "1",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"target_tracking_policies.#": "2",
						"scheduled_actions.#":        "2",
					}),
				),
			},
			// Reorder one TypeList at a time so each reorder apply changes only
			// that list; the post-apply plan must converge.
			{
				Config: testAccConfig(map[string]interface{}{
					"target_tracking_policies": []map[string]interface{}{
						{
							"name":          "t2",
							"start_time":    "2030-10-10T10:10:10Z",
							"end_time":      "2035-10-10T10:10:10Z",
							"min_capacity":  "0",
							"max_capacity":  "1",
							"metric_target": "1",
							"metric_type":   "ProvisionedConcurrencyUtilization",
						},
						{
							"name":          "t1",
							"start_time":    "2030-10-10T10:10:10Z",
							"end_time":      "2035-10-10T10:10:10Z",
							"min_capacity":  "0",
							"max_capacity":  "1",
							"metric_target": "1",
							"metric_type":   "ProvisionedConcurrencyUtilization",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"target_tracking_policies.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scheduled_actions": []map[string]interface{}{
						{
							"name":                "s2",
							"start_time":          "2030-10-10T10:10:10Z",
							"end_time":            "2035-10-10T10:10:10Z",
							"schedule_expression": "cron(0 0 6 * * *)",
							"target":              "1",
						},
						{
							"name":                "s1",
							"start_time":          "2030-10-10T10:10:10Z",
							"end_time":            "2035-10-10T10:10:10Z",
							"schedule_expression": "cron(0 0 4 * * *)",
							"target":              "0",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scheduled_actions.#": "2",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"current", "current_error"},
			},
		},
	})
}

// waitForProvisionConfigTarget polls the service because a scheduled action
// takes effect asynchronously after the cron trigger, so apply alone does not
// guarantee that the policy-computed runtime target is already visible.
func waitForProvisionConfigTarget(name, wantTarget, wantDefaultTarget string) resource.TestCheckFunc {
	return func(*terraform.State) error {
		service := fcv3.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
		deadline := time.Now().Add(6 * time.Minute)
		var last string
		for {
			object, err := service.DescribeProvisionConfig(name)
			if err != nil {
				return err
			}
			last = fmt.Sprintf("target=%v defaultTarget=%v", object["target"], object["defaultTarget"])
			if fmt.Sprint(object["target"]) == wantTarget && fmt.Sprint(object["defaultTarget"]) == wantDefaultTarget {
				return nil
			}
			if time.Now().After(deadline) {
				return fmt.Errorf("provision config did not reach target=%s defaultTarget=%s: %s", wantTarget, wantDefaultTarget, last)
			}
			time.Sleep(10 * time.Second)
		}
	}
}

// TestAccAliCloudFcv3ProvisionConfig_policyTarget verifies that an active
// scheduled action computes a runtime target (2) that differs from the
// configured default (1). Read must
// keep the configured default in state and the replacement PUT performed by a
// boolean-only update must send that default, otherwise the temporary policy
// target becomes the persistent default. The first step also confirms the
// API-reported split (target=2, defaultTarget=1), which a future or inactive
// policy does not exercise.
func TestAccAliCloudFcv3ProvisionConfig_policyTarget(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_fcv3_provision_config.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudFcv3ProvisionConfigMap7183)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return fcv3.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeProvisionConfig")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sfcv3provisionpolicy%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudFcv3ProvisionConfigBasicDependence7183)

	activeScheduledAction := []map[string]interface{}{
		{
			"name":                "always-on",
			"start_time":          "2020-01-01T00:00:00Z",
			"end_time":            "2035-01-01T00:00:00Z",
			"schedule_expression": "cron(0 * * * * *)",
			"target":              "2",
		},
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy: resource.ComposeTestCheckFunc(
			rac.CheckResourceDestroy(),
			checkProvisionConfigDependenciesDestroyed(),
		),
		Steps: []resource.TestStep{
			// Configured default 1 with an active scheduled action targeting 2.
			// The state keeps the configured default even though the API
			// reports the policy-computed runtime target separately.
			{
				Config: testAccConfig(map[string]interface{}{
					"function_name":       "${alicloud_fcv3_function.function.function_name}",
					"target":              "1",
					"always_allocate_cpu": "true",
					"always_allocate_gpu": "true",
					"scheduled_actions":   activeScheduledAction,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"target":              "1",
						"scheduled_actions.#": "1",
					}),
					waitForProvisionConfigTarget(name, "2", "1"),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"current", "current_error"},
			},
			// Change only the allocation booleans and omit target from the
			// configuration: the replacement PUT must preserve defaultTarget=1
			// instead of writing back the runtime target 2.
			{
				Config: testAccConfig(map[string]interface{}{
					"always_allocate_cpu": "false",
					"always_allocate_gpu": "false",
					"target":              acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"always_allocate_cpu": "false",
						"always_allocate_gpu": "false",
						"target":              "1",
						"scheduled_actions.#": "1",
					}),
					waitForProvisionConfigTarget(name, "2", "1"),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"current", "current_error"},
			},
			// Removing the policy must clear the action list and converge on
			// the configured default instead of the previously computed
			// runtime target.
			{
				Config: testAccConfig(map[string]interface{}{
					"scheduled_actions": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scheduled_actions.#": "0",
						"target":              "1",
					}),
					waitForProvisionConfigTarget(name, "1", "1"),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"current", "current_error"},
			},
		},
	})
}

// TestAccAliCloudFcv3ProvisionConfig_targetZero verifies that an explicitly
// configured top-level target keeps working when it changes from a non-zero
// value to 0. The service must
// accept the update, report defaultTarget=0 and target=0, Read must keep 0 in
// state, and the post-apply plan must be empty before the state is imported.
func TestAccAliCloudFcv3ProvisionConfig_targetZero(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_fcv3_provision_config.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudFcv3ProvisionConfigMap7183)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return fcv3.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeProvisionConfig")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sfcv3provisionzero%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudFcv3ProvisionConfigBasicDependence7183)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy: resource.ComposeTestCheckFunc(
			rac.CheckResourceDestroy(),
			checkProvisionConfigDependenciesDestroyed(),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"function_name":       "${alicloud_fcv3_function.function.function_name}",
					"target":              "2",
					"always_allocate_cpu": "true",
					"always_allocate_gpu": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"target": "2",
					}),
					waitForProvisionConfigTarget(name, "2", "2"),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"current", "current_error"},
			},
			// Explicit non-zero to zero: the value must be sent, accepted and
			// read back as both the runtime target and the configured default.
			{
				Config: testAccConfig(map[string]interface{}{
					"target": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"target": "0",
					}),
					waitForProvisionConfigTarget(name, "0", "0"),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"current", "current_error"},
			},
		},
	})
}

// TestAccAliCloudFcv3ProvisionConfig_policyRemovalBaseline establishes a fresh
// baseline where the configured default is 1 and an active scheduled action
// computes the runtime target 2, then removes the policy without any earlier
// boolean-only update. The policy list must become empty, the configured
// default must stay 1, the post-apply plan must converge and the state must
// remain importable.
func TestAccAliCloudFcv3ProvisionConfig_policyRemovalBaseline(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_fcv3_provision_config.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudFcv3ProvisionConfigMap7183)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return fcv3.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeProvisionConfig")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sfcv3provisionbaseline%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudFcv3ProvisionConfigBasicDependence7183)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy: resource.ComposeTestCheckFunc(
			rac.CheckResourceDestroy(),
			checkProvisionConfigDependenciesDestroyed(),
		),
		Steps: []resource.TestStep{
			// The baseline is built by this step alone: configured default 1
			// with an active scheduled action that makes the runtime target 2.
			{
				Config: testAccConfig(map[string]interface{}{
					"function_name":       "${alicloud_fcv3_function.function.function_name}",
					"target":              "1",
					"always_allocate_cpu": "true",
					"always_allocate_gpu": "true",
					"scheduled_actions": []map[string]interface{}{
						{
							"name":                "always-on",
							"start_time":          "2020-01-01T00:00:00Z",
							"end_time":            "2035-01-01T00:00:00Z",
							"schedule_expression": "cron(0 * * * * *)",
							"target":              "2",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"target":              "1",
						"scheduled_actions.#": "1",
					}),
					waitForProvisionConfigTarget(name, "2", "1"),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"current", "current_error"},
			},
			// Removing the policy from the independently established baseline
			// must clear the list and converge on the configured default.
			{
				Config: testAccConfig(map[string]interface{}{
					"target":            acctest.REMOVEKEY,
					"scheduled_actions": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scheduled_actions.#": "0",
						"target":              "1",
					}),
					waitForProvisionConfigTarget(name, "1", "1"),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"current", "current_error"},
			},
		},
	})
}

func testAccProvisionConfigQualifierConfig(name, qualifier string, target int, independentTarget *int) string {
	config := AlicloudFcv3ProvisionConfigBasicDependence7183(name) + fmt.Sprintf(`
resource "alicloud_fcv3_function_version" "default" {
  function_name = alicloud_fcv3_function.function.function_name
  description   = %q
}

resource "alicloud_fcv3_alias" "q1" {
  function_name = alicloud_fcv3_function.function.function_name
  version_id    = alicloud_fcv3_function_version.default.version_id
  alias_name    = "q1"
}

resource "alicloud_fcv3_alias" "q2" {
  function_name = alicloud_fcv3_function.function.function_name
  version_id    = alicloud_fcv3_function_version.default.version_id
  alias_name    = "q2"
}

resource "alicloud_fcv3_provision_config" "default" {
  function_name = alicloud_fcv3_function.function.function_name
  qualifier     = alicloud_fcv3_alias.%s.alias_name
  target        = %d
}
`, name, qualifier, target)
	if independentTarget != nil {
		config += fmt.Sprintf(`
resource "alicloud_fcv3_provision_config" "independent" {
  function_name = alicloud_fcv3_function.function.function_name
  qualifier     = alicloud_fcv3_alias.q2.alias_name
  target        = %d
}
`, *independentTarget)
	}
	return config
}

func testAccProvisionConfigLatestConfig(name, qualifier string, target int) string {
	qualifierConfig := ""
	if qualifier != "" {
		qualifierConfig = fmt.Sprintf("qualifier = %q", qualifier)
	}
	return AlicloudFcv3ProvisionConfigBasicDependence7183(name) + fmt.Sprintf(`
resource "alicloud_fcv3_provision_config" "default" {
  function_name = alicloud_fcv3_function.function.function_name
  target        = %d
  %s
}
`, target, qualifierConfig)
}

func provisionConfigFunctionName(state *terraform.State) (string, error) {
	rs, ok := state.RootModule().Resources["alicloud_fcv3_function.function"]
	if !ok || rs.Primary.ID == "" {
		return "", fmt.Errorf("cannot find the provision config test function")
	}
	return rs.Primary.ID, nil
}

func checkProvisionConfigQualifierTarget(qualifier string, want int) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		functionName, err := provisionConfigFunctionName(state)
		if err != nil {
			return err
		}
		service := fcv3.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
		object, err := service.DescribeProvisionConfigWithQualifier(functionName, qualifier)
		if err != nil {
			return fmt.Errorf("read provision config for qualifier %q: %w", qualifier, err)
		}
		target := object["defaultTarget"]
		if target == nil {
			target = object["target"]
		}
		if target == nil || fmt.Sprint(target) != fmt.Sprint(want) {
			return fmt.Errorf("provision config for qualifier %q has target %#v, want %d", qualifier, target, want)
		}
		return nil
	}
}

func checkProvisionConfigQualifierAbsent(qualifier string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		functionName, err := provisionConfigFunctionName(state)
		if err != nil {
			return err
		}
		service := fcv3.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
		_, err = service.DescribeProvisionConfigWithQualifier(functionName, qualifier)
		if errs.NotFoundError(err) {
			return nil
		}
		if err != nil {
			return fmt.Errorf("read old provision config for qualifier %q: %w", qualifier, err)
		}
		return fmt.Errorf("old provision config for qualifier %q still exists", qualifier)
	}
}

func checkProvisionConfigQualifiersDestroyed(qualifiers ...string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		functionName, err := provisionConfigFunctionName(state)
		if err != nil {
			return err
		}
		service := fcv3.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
		for _, qualifier := range qualifiers {
			_, err := service.DescribeProvisionConfigWithQualifier(functionName, qualifier)
			if errs.NotFoundError(err) {
				continue
			}
			if err != nil {
				return fmt.Errorf("check destroyed provision config for qualifier %q: %w", qualifier, err)
			}
			return fmt.Errorf("provision config for qualifier %q still exists after destroy", qualifier)
		}
		return nil
	}
}

// Named aliases cannot be imported by the historical function-name-only ID, so
// these create-only replacement cases use direct service assertions instead.
func TestAccAliCloudFcv3ProvisionConfig_qualifierReplacement(t *testing.T) {
	resourceId := "alicloud_fcv3_provision_config.default"
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sfcv3pcq%d", acctest.DefaultRegionToTest, rand)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy: resource.ComposeTestCheckFunc(
			checkProvisionConfigQualifiersDestroyed("q1", "q2"),
			checkProvisionConfigDependenciesDestroyed(),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccProvisionConfigQualifierConfig(name, "q1", 1, nil),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "qualifier", "q1"),
					checkProvisionConfigQualifierTarget("q1", 1),
				),
			},
			{
				Config: testAccProvisionConfigQualifierConfig(name, "q1", 2, nil),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceId, plancheck.ResourceActionUpdate),
					},
				},
				Check: checkProvisionConfigQualifierTarget("q1", 2),
			},
			{
				Config: testAccProvisionConfigQualifierConfig(name, "q2", 2, nil),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceId, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "qualifier", "q2"),
					checkProvisionConfigQualifierTarget("q2", 2),
					checkProvisionConfigQualifierAbsent("q1"),
				),
			},
		},
	})
}

func TestAccAliCloudFcv3ProvisionConfig_qualifierCollision(t *testing.T) {
	resourceId := "alicloud_fcv3_provision_config.default"
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sfcv3pcq%d", acctest.DefaultRegionToTest, rand)
	independentTarget := 0
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy: resource.ComposeTestCheckFunc(
			checkProvisionConfigQualifiersDestroyed("q1", "q2"),
			checkProvisionConfigDependenciesDestroyed(),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccProvisionConfigQualifierConfig(name, "q1", 1, &independentTarget),
				Check: resource.ComposeTestCheckFunc(
					checkProvisionConfigQualifierTarget("q1", 1),
					checkProvisionConfigQualifierTarget("q2", 0),
				),
			},
			{
				Config:      testAccProvisionConfigQualifierConfig(name, "q2", 2, &independentTarget),
				ExpectError: regexp.MustCompile(`(?s)qualifier.*q2.*refusing to overwrite an independent configuration`),
			},
			{
				Config: testAccProvisionConfigQualifierConfig(name, "q1", 1, &independentTarget),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "qualifier", "q1"),
					resource.TestCheckResourceAttr(resourceId, "target", "1"),
					checkProvisionConfigQualifierTarget("q1", 1),
					checkProvisionConfigQualifierTarget("q2", 0),
				),
			},
		},
	})
}

func TestAccAliCloudFcv3ProvisionConfig_qualifierEquivalence(t *testing.T) {
	resourceId := "alicloud_fcv3_provision_config.default"
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sfcv3pcq%d", acctest.DefaultRegionToTest, rand)
	importStep := func() resource.TestStep {
		return resource.TestStep{
			ResourceName:            resourceId,
			ImportState:             true,
			ImportStateVerify:       true,
			ImportStateVerifyIgnore: []string{"current", "current_error"},
		}
	}
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy: resource.ComposeTestCheckFunc(
			checkProvisionConfigQualifiersDestroyed("LATEST"),
			checkProvisionConfigDependenciesDestroyed(),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccProvisionConfigLatestConfig(name, "", 1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "qualifier", "LATEST"),
					checkProvisionConfigQualifierTarget("LATEST", 1),
				),
			},
			importStep(),
			{
				Config: testAccProvisionConfigLatestConfig(name, "LATEST", 1),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceId, plancheck.ResourceActionNoop),
					},
				},
			},
			importStep(),
			{
				Config: testAccProvisionConfigLatestConfig(name, "", 1),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceId, plancheck.ResourceActionNoop),
					},
				},
			},
			importStep(),
			{
				Config: testAccProvisionConfigLatestConfig(name, "", 2),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceId, plancheck.ResourceActionUpdate),
					},
				},
				Check: checkProvisionConfigQualifierTarget("LATEST", 2),
			},
			importStep(),
			{
				Config: testAccProvisionConfigLatestConfig(name, "LATEST", 2),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceId, plancheck.ResourceActionNoop),
					},
				},
			},
			importStep(),
		},
	})
}

// TestUnitProvisionConfigRequest verifies the request payload built by the
// public CRUD entry points. Omitted Optional+Computed fields must be absent, an
// explicit false must be sent as false, and an update that only toggles one
// boolean must still carry the API-required target and the current values of
// the untouched fields.
func TestUnitProvisionConfigRequest(t *testing.T) {
	type capturedRequest struct {
		method string
		query  string
		body   map[string]interface{}
	}

	var mu sync.Mutex
	var requests []capturedRequest
	functionArn := "acs:fc:cn-hangzhou:1234567890123456:functions/unit-provision"
	stored := map[string]interface{}{}
	// Empty and LATEST share stored; named qualifiers get independent objects.
	qualifiedStored := map[string]map[string]interface{}{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var body map[string]interface{}
		if req.Body != nil {
			raw, err := io.ReadAll(req.Body)
			if err != nil {
				t.Errorf("read request body: %v", err)
			}
			if len(raw) > 0 {
				if err := json.Unmarshal(raw, &body); err != nil {
					t.Errorf("decode request body: %v", err)
				}
			}
		}
		qualifier := req.URL.Query().Get("qualifier")

		mu.Lock()
		requests = append(requests, capturedRequest{method: req.Method, query: req.URL.RawQuery, body: body})
		object := stored
		if qualifier != "" && qualifier != "LATEST" {
			object = qualifiedStored[qualifier]
		}
		if req.Method == http.MethodPut {
			// PutProvisionConfig replaces the whole configuration: keys missing
			// from the request are gone from the response.
			replaced := map[string]interface{}{"functionArn": functionArn}
			for key, value := range body {
				replaced[key] = value
			}
			object = replaced
			if qualifier == "" || qualifier == "LATEST" {
				stored = replaced
			} else {
				qualifiedStored[qualifier] = replaced
			}
		}
		if req.Method == http.MethodDelete {
			object = map[string]interface{}{}
			if qualifier == "" || qualifier == "LATEST" {
				stored = object
			} else {
				qualifiedStored[qualifier] = object
			}
		}
		response := make(map[string]interface{}, len(object))
		for key, value := range object {
			response[key] = value
		}
		mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Errorf("encode response: %v", err)
		}
	}))
	defer server.Close()

	client := provisionConfigUnitTestClient(t, server)
	r := fcv3.ResourceProvisionConfig()
	ctx := context.Background()

	lastPut := func() map[string]interface{} {
		mu.Lock()
		defer mu.Unlock()
		for i := len(requests) - 1; i >= 0; i-- {
			if requests[i].method == http.MethodPut {
				return requests[i].body
			}
		}
		t.Fatal("no PUT request captured")
		return nil
	}
	lastPutQuery := func() string {
		mu.Lock()
		defer mu.Unlock()
		for i := len(requests) - 1; i >= 0; i-- {
			if requests[i].method == http.MethodPut {
				return requests[i].query
			}
		}
		t.Fatal("no PUT request captured")
		return ""
	}
	requestCount := func() int {
		mu.Lock()
		defer mu.Unlock()
		return len(requests)
	}
	requestsSince := func(start int) []capturedRequest {
		mu.Lock()
		defer mu.Unlock()
		return append([]capturedRequest(nil), requests[start:]...)
	}
	assertValue := func(body map[string]interface{}, key string, want interface{}) {
		t.Helper()
		got, ok := body[key]
		if !ok {
			t.Fatalf("request body is missing %q: %#v", key, body)
		}
		if got != want {
			t.Fatalf("request body %q = %#v, want %#v", key, got, want)
		}
	}
	assertMissing := func(body map[string]interface{}, key string) {
		t.Helper()
		if value, ok := body[key]; ok {
			t.Fatalf("request body must omit %q, got %#v", key, value)
		}
	}
	assertListLength := func(body map[string]interface{}, key string, want int) []interface{} {
		t.Helper()
		values, ok := body[key].([]interface{})
		if !ok {
			t.Fatalf("request body %q is not a list: %#v", key, body[key])
		}
		if len(values) != want {
			t.Fatalf("request body %q has %d entries, want %d: %#v", key, len(values), want, values)
		}
		return values
	}
	priorState := func(attributes map[string]string) *sdktf.InstanceState {
		merged := map[string]string{
			"id":            "unit-provision",
			"function_name": "unit-provision",
			"target":        "1",
		}
		for key, value := range attributes {
			merged[key] = value
		}
		return &sdktf.InstanceState{ID: "unit-provision", Attributes: merged}
	}
	applyUpdate := func(prior *sdktf.InstanceState, config map[string]interface{}) map[string]interface{} {
		t.Helper()
		resourceConfig := sdktf.NewResourceConfigRaw(config)
		requestsBeforeDiff := requestCount()
		diff, err := r.Diff(ctx, prior, resourceConfig, client)
		if err != nil {
			t.Fatalf("diff: %v", err)
		}
		if requestCount() != requestsBeforeDiff {
			t.Fatalf("in-place update planning made an occupancy read for %#v", config)
		}
		if diff == nil || diff.Empty() {
			t.Fatalf("expected a plan for %#v", config)
		}
		if diff.RequiresNew() {
			t.Fatalf("in-place update unexpectedly requires replacement for %#v: %#v", config, diff)
		}
		next, diags := r.Apply(ctx, prior, diff, client)
		if diags.HasError() {
			t.Fatalf("apply: %v", diags)
		}
		again, err := r.Diff(ctx, next, resourceConfig, client)
		if err != nil {
			t.Fatalf("diff after apply: %v", err)
		}
		if again != nil && !again.Empty() {
			t.Fatalf("configuration did not converge after apply: %#v", again)
		}
		return lastPut()
	}
	// applyUpdateWithRawConfig mirrors the protocol layer, which hands the
	// evaluated configuration to the provider as cty raw config before Apply.
	// The schema test helpers do not populate it, so set it explicitly here.
	applyUpdateWithRawConfig := func(prior *sdktf.InstanceState, config map[string]interface{}, rawConfig cty.Value) map[string]interface{} {
		t.Helper()
		resourceConfig := sdktf.NewResourceConfigRaw(config)
		requestsBeforeDiff := requestCount()
		diff, err := r.Diff(ctx, prior, resourceConfig, client)
		if err != nil {
			t.Fatalf("diff: %v", err)
		}
		if requestCount() != requestsBeforeDiff {
			t.Fatalf("in-place update planning made an occupancy read for %#v", config)
		}
		if diff == nil || diff.Empty() {
			t.Fatalf("expected a plan for %#v", config)
		}
		if diff.RequiresNew() {
			t.Fatalf("in-place update unexpectedly requires replacement for %#v: %#v", config, diff)
		}
		diff.RawConfig = rawConfig
		next, diags := r.Apply(ctx, prior, diff, client)
		if diags.HasError() {
			t.Fatalf("apply: %v", diags)
		}
		again, err := r.Diff(ctx, next, resourceConfig, client)
		if err != nil {
			t.Fatalf("diff after apply: %v", err)
		}
		if again != nil && !again.Empty() {
			t.Fatalf("configuration did not converge after apply: %#v", again)
		}
		return lastPut()
	}

	// Create with both fields omitted: the keys must not be sent at all.
	createData := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{
		"function_name": "unit-provision",
		"target":        1,
	})
	if err := r.Create(createData, client); err != nil {
		t.Fatalf("create: %v", err)
	}
	body := lastPut()
	assertValue(body, "target", float64(1))
	assertValue(body, "defaultTarget", float64(1))
	assertMissing(body, "alwaysAllocateCPU")
	assertMissing(body, "alwaysAllocateGPU")

	// Create with target omitted: the client must not substitute the zero
	// value. PutProvisionConfig declares target as required, so the real
	// service rejects this request; the request shape is what this unit test
	// can pin down.
	mu.Lock()
	stored = map[string]interface{}{}
	mu.Unlock()
	createData = schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{
		"function_name": "unit-provision",
	})
	if err := r.Create(createData, client); err != nil {
		t.Fatalf("create without target: %v", err)
	}
	body = lastPut()
	assertMissing(body, "target")
	assertMissing(body, "defaultTarget")

	mu.Lock()
	stored = map[string]interface{}{
		"functionArn":       functionArn,
		"alwaysAllocateCPU": true,
		"alwaysAllocateGPU": true,
		"target":            float64(1),
		"defaultTarget":     float64(1),
	}
	mu.Unlock()

	// Boolean-only update from a real prior state: target is unchanged, so the
	// old implementation omitted it; it must still be present.
	body = applyUpdate(priorState(map[string]string{
		"always_allocate_cpu": "true",
		"always_allocate_gpu": "true",
	}), map[string]interface{}{
		"function_name":       "unit-provision",
		"target":              1,
		"always_allocate_cpu": false,
		"always_allocate_gpu": true,
	})
	assertValue(body, "target", float64(1))
	assertValue(body, "defaultTarget", float64(1))
	assertValue(body, "alwaysAllocateCPU", false)
	assertValue(body, "alwaysAllocateGPU", true)

	// Explicit false for both fields must be sent verbatim.
	body = applyUpdate(priorState(map[string]string{
		"always_allocate_cpu": "false",
		"always_allocate_gpu": "true",
	}), map[string]interface{}{
		"function_name":       "unit-provision",
		"target":              1,
		"always_allocate_cpu": false,
		"always_allocate_gpu": false,
	})
	assertValue(body, "alwaysAllocateCPU", false)
	assertValue(body, "alwaysAllocateGPU", false)
	assertValue(body, "target", float64(1))

	// Switching both back to true must also be sent verbatim.
	body = applyUpdate(priorState(map[string]string{
		"always_allocate_cpu": "false",
		"always_allocate_gpu": "false",
	}), map[string]interface{}{
		"function_name":       "unit-provision",
		"target":              1,
		"always_allocate_cpu": true,
		"always_allocate_gpu": true,
	})
	assertValue(body, "alwaysAllocateCPU", true)
	assertValue(body, "alwaysAllocateGPU", true)
	assertValue(body, "target", float64(1))

	// An update whose configuration omits target must recover the current
	// service default instead of sending the zero value of an unset argument.
	mu.Lock()
	stored["target"] = float64(5)
	stored["defaultTarget"] = float64(5)
	mu.Unlock()
	body = applyUpdate(priorState(map[string]string{
		"target":              "5",
		"always_allocate_cpu": "true",
		"always_allocate_gpu": "true",
	}), map[string]interface{}{
		"function_name":       "unit-provision",
		"always_allocate_cpu": false,
		"always_allocate_gpu": true,
	})
	assertValue(body, "target", float64(5))

	// An explicitly configured zero is a real value and must still be sent.
	body = applyUpdate(priorState(map[string]string{
		"target":              "5",
		"always_allocate_cpu": "true",
		"always_allocate_gpu": "true",
	}), map[string]interface{}{
		"function_name":       "unit-provision",
		"target":              0,
		"always_allocate_cpu": true,
		"always_allocate_gpu": true,
	})
	assertValue(body, "target", float64(0))
	assertValue(body, "defaultTarget", float64(0))

	// The qualifier selects the target object for both the replacement PUT and
	// the subsequent refresh.
	mu.Lock()
	stored = map[string]interface{}{
		"functionArn":       functionArn,
		"alwaysAllocateCPU": false,
		"alwaysAllocateGPU": true,
		"target":            float64(1),
		"defaultTarget":     float64(1),
	}
	qualifiedStored["prod"] = map[string]interface{}{
		"functionArn":       functionArn,
		"alwaysAllocateCPU": false,
		"alwaysAllocateGPU": true,
		"target":            float64(1),
		"defaultTarget":     float64(1),
	}
	mu.Unlock()
	body = applyUpdate(priorState(map[string]string{
		"qualifier":           "prod",
		"always_allocate_cpu": "true",
		"always_allocate_gpu": "true",
	}), map[string]interface{}{
		"function_name":       "unit-provision",
		"target":              1,
		"qualifier":           "prod",
		"always_allocate_cpu": false,
		"always_allocate_gpu": true,
	})
	assertValue(body, "target", float64(1))
	if query := lastPutQuery(); !strings.Contains(query, "qualifier=prod") {
		t.Fatalf("PUT query %q does not carry the configured qualifier", query)
	}

	listConfig := map[string]interface{}{
		"target_tracking_policies": []interface{}{map[string]interface{}{
			"name":          "t1",
			"start_time":    "2030-10-10T10:10:10Z",
			"end_time":      "2035-10-10T10:10:10Z",
			"min_capacity":  0,
			"max_capacity":  1,
			"metric_target": 0.6,
			"metric_type":   "CPUUtilization",
		}},
		"scheduled_actions": []interface{}{map[string]interface{}{
			"name":                "s1",
			"start_time":          "2030-10-10T10:10:10Z",
			"end_time":            "2035-10-10T10:10:10Z",
			"schedule_expression": "cron(0 0 4 * * *)",
			"target":              0,
		}},
	}
	listState := map[string]string{
		"target_tracking_policies.#":               "1",
		"target_tracking_policies.0.name":          "t1",
		"target_tracking_policies.0.start_time":    "2030-10-10T10:10:10Z",
		"target_tracking_policies.0.end_time":      "2035-10-10T10:10:10Z",
		"target_tracking_policies.0.min_capacity":  "0",
		"target_tracking_policies.0.max_capacity":  "1",
		"target_tracking_policies.0.metric_target": "0.6",
		"target_tracking_policies.0.metric_type":   "CPUUtilization",
		"scheduled_actions.#":                      "1",
		"scheduled_actions.0.name":                 "s1",
		"scheduled_actions.0.start_time":           "2030-10-10T10:10:10Z",
		"scheduled_actions.0.end_time":             "2035-10-10T10:10:10Z",
		"scheduled_actions.0.schedule_expression":  "cron(0 0 4 * * *)",
		"scheduled_actions.0.target":               "0",
	}

	// A boolean-only update must keep configured lists in the replacement body.
	listPrior := priorState(map[string]string{
		"always_allocate_cpu": "true",
		"always_allocate_gpu": "true",
	})
	for key, value := range listState {
		listPrior.Attributes[key] = value
	}
	body = applyUpdate(listPrior, map[string]interface{}{
		"function_name":            "unit-provision",
		"target":                   1,
		"always_allocate_cpu":      false,
		"always_allocate_gpu":      true,
		"target_tracking_policies": listConfig["target_tracking_policies"],
		"scheduled_actions":        listConfig["scheduled_actions"],
	})
	assertValue(body, "target", float64(1))
	policies := assertListLength(body, "targetTrackingPolicies", 1)
	policy, ok := policies[0].(map[string]interface{})
	if !ok {
		t.Fatalf("expected a target tracking policy object, got %#v", policies[0])
	}
	assertValue(policy, "metricType", "CPUUtilization")
	assertValue(policy, "metricTarget", 0.6)
	assertListLength(body, "scheduledActions", 1)

	// Removing the lists must send empty arrays so the replacement clears them.
	removePrior := priorState(map[string]string{
		"always_allocate_cpu": "false",
		"always_allocate_gpu": "true",
	})
	for key, value := range listState {
		removePrior.Attributes[key] = value
	}
	body = applyUpdate(removePrior, map[string]interface{}{
		"function_name":       "unit-provision",
		"target":              1,
		"always_allocate_cpu": true,
		"always_allocate_gpu": true,
	})
	assertListLength(body, "targetTrackingPolicies", 0)
	assertListLength(body, "scheduledActions", 0)

	// A policy-computed runtime target must never be stored as the desired
	// configuration. Read maps target from the API defaultTarget, so the state
	// keeps the configured default while a scheduled action targets 2.
	mu.Lock()
	stored["defaultTarget"] = float64(1)
	stored["target"] = float64(2)
	mu.Unlock()
	readData := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{
		"function_name": "unit-provision",
		"target":        2,
	})
	readData.SetId("unit-provision")
	if err := r.Read(readData, client); err != nil {
		t.Fatalf("read: %v", err)
	}
	if got := readData.Get("target"); got != 1 {
		t.Fatalf("read mapped target = %#v, want the configured defaultTarget 1", got)
	}

	// Consequently a boolean-only update whose configuration omits target
	// sends the configured default instead of the runtime value reported by
	// the service. The prior state comes from the Read above, so this case
	// fails as soon as Read stores the policy-computed target again.
	priorPolicy := readData.State()
	body = applyUpdate(priorPolicy, map[string]interface{}{
		"function_name":       "unit-provision",
		"always_allocate_cpu": false,
		"always_allocate_gpu": true,
	})
	assertValue(body, "target", float64(1))
	assertValue(body, "defaultTarget", float64(1))

	// A response that omits defaultTarget while a scaling policy is active must
	// never let the zero value become the persistent default. Read keeps the
	// configured default already established in state.
	mu.Lock()
	stored = map[string]interface{}{
		"functionArn": "acs:fc:cn-hangzhou:1234567890123456:functions/unit-provision",
		"target":      float64(2),
		"scheduledActions": []interface{}{map[string]interface{}{
			"name":   "always-on",
			"target": float64(2),
		}},
	}
	mu.Unlock()
	keepData := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{
		"function_name": "unit-provision",
		"target":        1,
	})
	keepData.SetId("unit-provision")
	if err := r.Read(keepData, client); err != nil {
		t.Fatalf("read without defaultTarget: %v", err)
	}
	if got := keepData.Get("target"); got != 1 {
		t.Fatalf("read replaced the established default with %#v", got)
	}

	// A non-zero legacy state can contain the runtime target rather than the
	// persistent default. Update must recover the authoritative default instead
	// of trusting that inherited value.
	legacyData := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{
		"function_name": "unit-provision",
		"target":        2,
	})
	legacyData.SetId("unit-provision")
	if err := r.Read(legacyData, client); err != nil {
		t.Fatalf("read legacy state: %v", err)
	}
	if got := legacyData.Get("target"); got != 2 {
		t.Fatalf("read changed the unresolved legacy value to %#v", got)
	}
	mu.Lock()
	stored["defaultTarget"] = float64(1)
	mu.Unlock()
	body = applyUpdate(legacyData.State(), map[string]interface{}{
		"function_name":       "unit-provision",
		"always_allocate_cpu": false,
	})
	assertValue(body, "target", float64(1))
	assertValue(body, "defaultTarget", float64(1))

	// Without a scaling policy the effective target equals the configured
	// default, so a response that omits defaultTarget is still authoritative.
	// Read must replace a stale state value with the reported target, and the
	// following boolean-only update must PUT that value unchanged.
	mu.Lock()
	stored = map[string]interface{}{
		"functionArn":       functionArn,
		"alwaysAllocateCPU": true,
		"alwaysAllocateGPU": true,
		"target":            float64(3),
	}
	mu.Unlock()
	noPolicyData := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{
		"function_name": "unit-provision",
		"target":        1,
	})
	noPolicyData.SetId("unit-provision")
	if err := r.Read(noPolicyData, client); err != nil {
		t.Fatalf("read without defaultTarget and without a policy: %v", err)
	}
	if got := noPolicyData.Get("target"); got != 3 {
		t.Fatalf("read target = %#v, want the reported target 3 without a scaling policy", got)
	}
	body = applyUpdate(noPolicyData.State(), map[string]interface{}{
		"function_name":       "unit-provision",
		"always_allocate_cpu": false,
	})
	assertValue(body, "target", float64(3))
	assertValue(body, "defaultTarget", float64(3))

	// A stale zero in state must not replace the authoritative no-policy
	// fallback when Update reads the same response without defaultTarget.
	mu.Lock()
	stored = map[string]interface{}{
		"functionArn":       functionArn,
		"alwaysAllocateCPU": true,
		"alwaysAllocateGPU": true,
		"target":            float64(3),
	}
	mu.Unlock()
	body = applyUpdate(priorState(map[string]string{
		"target":              "0",
		"always_allocate_cpu": "true",
		"always_allocate_gpu": "true",
	}), map[string]interface{}{
		"function_name":       "unit-provision",
		"always_allocate_cpu": false,
	})
	assertValue(body, "target", float64(3))
	assertValue(body, "defaultTarget", float64(3))

	// An imported state has no established value because Read cannot map one
	// while the response omits defaultTarget and a policy is active. The
	// update path must ask the service for the configured default instead of
	// sending 0.
	mu.Lock()
	stored = map[string]interface{}{
		"functionArn": "acs:fc:cn-hangzhou:1234567890123456:functions/unit-provision",
		"target":      float64(2),
		"scheduledActions": []interface{}{map[string]interface{}{
			"name":   "always-on",
			"target": float64(2),
		}},
	}
	mu.Unlock()
	importData := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{
		"function_name": "unit-provision",
	})
	importData.SetId("unit-provision")
	if err := r.Read(importData, client); err != nil {
		t.Fatalf("import read: %v", err)
	}
	if got := importData.Get("target"); got != 0 {
		t.Fatalf("import read target = %#v, want the untouched zero value", got)
	}
	// The service reports the configured default to the update path even
	// though the earlier Read response omitted it.
	mu.Lock()
	stored["defaultTarget"] = float64(1)
	mu.Unlock()
	body = applyUpdate(importData.State(), map[string]interface{}{
		"function_name":       "unit-provision",
		"always_allocate_cpu": false,
	})
	assertValue(body, "target", float64(1))
	assertValue(body, "defaultTarget", float64(1))

	// When the service keeps omitting defaultTarget while a policy is active,
	// the update must fail instead of persisting the zero value.
	mu.Lock()
	stored = map[string]interface{}{
		"functionArn": "acs:fc:cn-hangzhou:1234567890123456:functions/unit-provision",
		"target":      float64(2),
		"scheduledActions": []interface{}{map[string]interface{}{
			"name":   "always-on",
			"target": float64(2),
		}},
	}
	mu.Unlock()
	missingData := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{
		"function_name": "unit-provision",
	})
	missingData.SetId("unit-provision")
	if err := r.Read(missingData, client); err != nil {
		t.Fatalf("read without defaultTarget and with a policy: %v", err)
	}
	missingPrior := missingData.State()
	missingConfig := map[string]interface{}{
		"function_name":       "unit-provision",
		"always_allocate_cpu": false,
	}
	missingDiff, err := r.Diff(ctx, missingPrior, sdktf.NewResourceConfigRaw(missingConfig), nil)
	if err != nil {
		t.Fatalf("diff for missing defaultTarget: %v", err)
	}
	if missingDiff == nil || missingDiff.Empty() {
		t.Fatal("expected a plan for the boolean-only update")
	}
	putStart := requestCount()
	_, diags := r.Apply(ctx, missingPrior, missingDiff, client)
	if !diags.HasError() {
		t.Fatal("update must fail when defaultTarget is missing and a policy is active")
	}
	// The failed update must stop before the replacement PUT: sending zero
	// would overwrite the configured default with an unintended value.
	for _, request := range requestsSince(putStart) {
		if request.method == http.MethodPut {
			t.Fatalf("a replacement PUT was sent although the configured default could not be determined: %#v", request.body)
		}
	}
	foundDefaultTargetReason := false
	for _, diagnostic := range diags {
		if strings.Contains(diagnostic.Summary, "defaultTarget") || strings.Contains(diagnostic.Detail, "defaultTarget") {
			foundDefaultTargetReason = true
		}
	}
	if !foundDefaultTargetReason {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}

	// Read and Update must address the configured qualifier. LATEST carries a
	// different default on purpose: using it would persist 7 on prod instead of
	// preserving the prod default 42, including when the inherited state is
	// non-zero.
	mu.Lock()
	stored = map[string]interface{}{
		"functionArn":       functionArn,
		"alwaysAllocateCPU": false,
		"alwaysAllocateGPU": true,
		"target":            float64(7),
		"defaultTarget":     float64(7),
	}
	qualifiedStored["prod"] = map[string]interface{}{
		"functionArn":       functionArn,
		"alwaysAllocateCPU": true,
		"alwaysAllocateGPU": true,
		"target":            float64(5),
		"defaultTarget":     float64(42),
	}
	mu.Unlock()

	qualifiedReadData := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{
		"function_name": "unit-provision",
		"qualifier":     "prod",
		"target":        7,
	})
	qualifiedReadData.SetId("unit-provision")
	start := requestCount()
	if err := r.Read(qualifiedReadData, client); err != nil {
		t.Fatalf("read configured qualifier: %v", err)
	}
	if got := qualifiedReadData.Get("target"); got != 42 {
		t.Fatalf("read target = %#v, want the prod default 42", got)
	}
	body = applyUpdate(qualifiedReadData.State(), map[string]interface{}{
		"function_name":       "unit-provision",
		"qualifier":           "prod",
		"always_allocate_cpu": false,
		"always_allocate_gpu": true,
	})
	assertValue(body, "target", float64(42))
	assertValue(body, "defaultTarget", float64(42))
	if query := lastPutQuery(); !strings.Contains(query, "qualifier=prod") {
		t.Fatalf("PUT query %q does not carry the configured qualifier", query)
	}
	qualifiedReads := 0
	for _, request := range requestsSince(start) {
		if request.method != http.MethodGet {
			continue
		}
		qualifiedReads++
		if !strings.Contains(request.query, "qualifier=prod") {
			t.Fatalf("GET query %q does not carry the configured qualifier", request.query)
		}
	}
	if qualifiedReads != 3 {
		t.Fatalf("qualified reads = %d, want refresh, default recovery and final refresh", qualifiedReads)
	}

	// Prior state 0, configuration 0 and no target change. The raw
	// configuration proves the zero is explicit, so the update must not read
	// the service default. The stored default is 7: falling back would send 7,
	// which fails the body assertion and the convergence check.
	mu.Lock()
	stored = map[string]interface{}{
		"functionArn":       functionArn,
		"alwaysAllocateCPU": true,
		"alwaysAllocateGPU": true,
		"target":            float64(7),
		"defaultTarget":     float64(7),
	}
	mu.Unlock()

	configType := r.CoreConfigSchema().ImpliedType()
	rawValues := make(map[string]cty.Value, len(configType.AttributeTypes()))
	for name, attrType := range configType.AttributeTypes() {
		rawValues[name] = cty.NullVal(attrType)
	}
	rawValues["function_name"] = cty.StringVal("unit-provision")
	rawValues["target"] = cty.NumberIntVal(0)
	rawValues["always_allocate_cpu"] = cty.False
	rawValues["always_allocate_gpu"] = cty.True

	zeroPrior := priorState(map[string]string{
		"target":              "0",
		"always_allocate_cpu": "true",
		"always_allocate_gpu": "true",
	})
	start = requestCount()
	body = applyUpdateWithRawConfig(zeroPrior, map[string]interface{}{
		"function_name":       "unit-provision",
		"target":              0,
		"always_allocate_cpu": false,
		"always_allocate_gpu": true,
	}, cty.ObjectVal(rawValues))
	assertValue(body, "target", float64(0))
	assertValue(body, "defaultTarget", float64(0))
	for _, request := range requestsSince(start) {
		if request.method == http.MethodPut {
			break
		}
		if request.method == http.MethodGet {
			t.Fatalf("update read the service default %q although the configuration set target explicitly", request.query)
		}
	}

	// Delete and its waiter must observe the same qualifier. The fixture keeps
	// LATEST present, so polling LATEST would never observe the prod deletion.
	mu.Lock()
	qualifiedStored["prod"] = map[string]interface{}{
		"functionArn":   functionArn,
		"defaultTarget": float64(42),
	}
	mu.Unlock()
	deleteData := r.Data(&sdktf.InstanceState{
		ID: "unit-provision",
		Attributes: map[string]string{
			"id":            "unit-provision",
			"function_name": "unit-provision",
			"qualifier":     "prod",
		},
	})
	start = requestCount()
	if err := r.Delete(deleteData, client); err != nil {
		t.Fatalf("delete configured qualifier: %v", err)
	}
	deleteRequests := requestsSince(start)
	if len(deleteRequests) != 2 {
		t.Fatalf("delete requests = %d, want DELETE and one qualifier-scoped read", len(deleteRequests))
	}
	if deleteRequests[0].method != http.MethodDelete || deleteRequests[1].method != http.MethodGet {
		t.Fatalf("unexpected delete sequence: %s then %s", deleteRequests[0].method, deleteRequests[1].method)
	}
	for _, request := range deleteRequests {
		if !strings.Contains(request.query, "qualifier=prod") {
			t.Fatalf("%s query %q does not carry the configured qualifier", request.method, request.query)
		}
	}
}

func TestUnitProvisionConfigQualifierReplacement(t *testing.T) {
	type capturedRequest struct {
		method    string
		qualifier string
	}
	functionArn := "acs:fc:cn-hangzhou:1234567890123456:functions/unit-provision"
	stored := map[string]map[string]interface{}{
		"q1": {
			"functionArn":   functionArn,
			"defaultTarget": float64(1),
		},
		"LATEST": {
			"functionArn":   functionArn,
			"defaultTarget": float64(9),
		},
	}
	var mu sync.Mutex
	var requests []capturedRequest
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		qualifier := req.URL.Query().Get("qualifier")
		if qualifier == "" {
			qualifier = "LATEST"
		}
		mu.Lock()
		requests = append(requests, capturedRequest{method: req.Method, qualifier: qualifier})
		object := stored[qualifier]
		response := map[string]interface{}{}
		for key, value := range object {
			response[key] = value
		}
		mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Errorf("encode response: %v", err)
		}
	}))
	defer server.Close()

	client := provisionConfigUnitTestClient(t, server)
	r := fcv3.ResourceProvisionConfig()
	ctx := context.Background()
	priorState := func(qualifier string) *sdktf.InstanceState {
		return &sdktf.InstanceState{
			ID: "unit-provision",
			Attributes: map[string]string{
				"id":            "unit-provision",
				"function_name": "unit-provision",
				"qualifier":     qualifier,
				"target":        "1",
			},
		}
	}
	config := func(qualifier *string, target int) map[string]interface{} {
		result := map[string]interface{}{
			"function_name": "unit-provision",
			"target":        target,
		}
		if qualifier != nil {
			result["qualifier"] = *qualifier
		}
		return result
	}
	latest := "LATEST"
	q1 := "q1"
	q2 := "q2"

	createDiff, err := r.Diff(ctx, nil, sdktf.NewResourceConfigRaw(config(nil, 1)), client)
	if err != nil {
		t.Fatalf("omitted-qualifier create diff: %v", err)
	}
	if createDiff == nil || createDiff.Attributes["qualifier"] == nil || createDiff.Attributes["qualifier"].New != "LATEST" {
		t.Fatalf("omitted qualifier must be planned as LATEST: %#v", createDiff)
	}

	for _, tc := range []struct {
		name          string
		oldQualifier  string
		newQualifier  *string
		wantEmptyDiff bool
	}{
		{name: "omitted to latest", oldQualifier: "", newQualifier: &latest, wantEmptyDiff: true},
		{name: "latest to omitted", oldQualifier: "LATEST", newQualifier: nil, wantEmptyDiff: true},
		{name: "same qualifier target update", oldQualifier: "q1", newQualifier: &q1, wantEmptyDiff: false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			target := 1
			if !tc.wantEmptyDiff {
				target = 2
			}
			mu.Lock()
			requests = nil
			mu.Unlock()
			diff, err := r.Diff(ctx, priorState(tc.oldQualifier), sdktf.NewResourceConfigRaw(config(tc.newQualifier, target)), client)
			if err != nil {
				t.Fatalf("diff: %v", err)
			}
			if tc.wantEmptyDiff {
				if diff != nil && !diff.Empty() {
					t.Fatalf("equivalent qualifiers produced a diff: %#v", diff)
				}
			} else {
				if diff == nil || diff.Empty() || diff.RequiresNew() {
					t.Fatalf("same-qualifier target change must be an in-place update: %#v", diff)
				}
			}
			mu.Lock()
			requestCount := len(requests)
			mu.Unlock()
			if requestCount != 0 {
				t.Fatalf("non-replacement diff made %d cloud reads, want 0", requestCount)
			}
		})
	}

	diff, err := r.Diff(ctx, priorState("q1"), sdktf.NewResourceConfigRaw(config(&q2, 2)), client)
	if err != nil {
		t.Fatalf("diff with absent target qualifier: %v", err)
	}
	if diff == nil || !diff.RequiresNew() {
		t.Fatalf("q1 to q2 must require replacement: %#v", diff)
	}

	mu.Lock()
	stored["q2"] = map[string]interface{}{
		"functionArn":   functionArn,
		"defaultTarget": float64(0),
	}
	requests = nil
	mu.Unlock()
	_, err = r.Diff(ctx, priorState("q1"), sdktf.NewResourceConfigRaw(config(&q2, 2)), client)
	if err == nil || !strings.Contains(err.Error(), "refusing to overwrite an independent configuration") {
		t.Fatalf("occupied target qualifier must reject planning, got %v", err)
	}
	mu.Lock()
	for _, request := range requests {
		if request.method != http.MethodGet {
			mu.Unlock()
			t.Fatalf("rejected plan made a write request: %#v", request)
		}
	}
	q1Target := stored["q1"]["defaultTarget"]
	q2Target := stored["q2"]["defaultTarget"]
	requests = nil
	mu.Unlock()
	if q1Target != float64(1) || q2Target != float64(0) {
		t.Fatalf("rejected plan changed source or target: q1=%#v q2=%#v", q1Target, q2Target)
	}

	protocolValue := func(values map[string]cty.Value) cty.Value {
		attributes := map[string]cty.Value{}
		for name, attributeType := range r.CoreConfigSchema().ImpliedType().AttributeTypes() {
			attributes[name] = cty.NullVal(attributeType)
		}
		for name, value := range values {
			attributes[name] = value
		}
		return cty.ObjectVal(attributes)
	}
	protocolConfig := func(values map[string]cty.Value) *sdktf.ResourceConfig {
		return sdktf.NewResourceConfigShimmed(protocolValue(values), r.CoreConfigSchema())
	}
	for _, tc := range []struct {
		name   string
		values map[string]cty.Value
	}{
		{
			name: "unknown qualifier",
			values: map[string]cty.Value{
				"function_name": cty.StringVal("unit-provision"),
				"qualifier":     cty.UnknownVal(cty.String),
				"target":        cty.NumberIntVal(2),
			},
		},
		{
			name: "unknown function name",
			values: map[string]cty.Value{
				"function_name": cty.UnknownVal(cty.String),
				"qualifier":     cty.StringVal("q2"),
				"target":        cty.NumberIntVal(2),
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			mu.Lock()
			requests = nil
			mu.Unlock()
			_, err := r.Diff(ctx, priorState("q1"), protocolConfig(tc.values), client)
			if err == nil || !strings.Contains(err.Error(), "must be known at plan time") {
				t.Fatalf("unknown replacement identity must reject planning, got %v", err)
			}
			mu.Lock()
			defer mu.Unlock()
			if len(requests) != 0 {
				t.Fatalf("unknown replacement identity made cloud requests: %#v", requests)
			}
			if stored["q1"]["defaultTarget"] != float64(1) || stored["q2"]["defaultTarget"] != float64(0) {
				t.Fatalf("unknown replacement identity changed source or target: %#v", stored)
			}
		})
	}

	for _, tc := range []struct {
		name           string
		priorQualifier string
		values         map[string]cty.Value
		wantError      string
		wantRequest    bool
	}{
		{
			name: "tainted replacement with occupied target",
			values: map[string]cty.Value{
				"function_name": cty.StringVal("unit-provision"),
				"qualifier":     cty.StringVal("q2"),
				"target":        cty.NumberIntVal(2),
			},
			wantError:   "refusing to overwrite an independent configuration",
			wantRequest: true,
		},
		{
			name: "tainted replacement with unknown qualifier",
			values: map[string]cty.Value{
				"function_name": cty.StringVal("unit-provision"),
				"qualifier":     cty.UnknownVal(cty.String),
				"target":        cty.NumberIntVal(2),
			},
			wantError: "must be known before deleting",
		},
		{
			name:           "tainted LATEST replacement with unknown qualifier",
			priorQualifier: "LATEST",
			values: map[string]cty.Value{
				"function_name": cty.StringVal("unit-provision"),
				"qualifier":     cty.UnknownVal(cty.String),
				"target":        cty.NumberIntVal(2),
			},
			wantError: "must be known before deleting",
		},
		{
			name: "tainted replacement with unknown function name",
			values: map[string]cty.Value{
				"function_name": cty.UnknownVal(cty.String),
				"qualifier":     cty.StringVal("q2"),
				"target":        cty.NumberIntVal(2),
			},
			wantError: "must be known before deleting",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			mu.Lock()
			requests = nil
			mu.Unlock()
			priorQualifier := tc.priorQualifier
			if priorQualifier == "" {
				priorQualifier = "q1"
			}
			state := priorState(priorQualifier)
			state.Tainted = true
			diff, err := r.Diff(ctx, state, protocolConfig(tc.values), client)
			if err != nil {
				t.Fatalf("tainted diff: %v", err)
			}
			diff.RawConfig = protocolValue(tc.values)
			_, diags := r.Apply(ctx, state, diff, client)
			if !diags.HasError() || !strings.Contains(fmt.Sprint(diags), tc.wantError) {
				t.Fatalf("tainted replacement must fail before source deletion, got %v", diags)
			}
			mu.Lock()
			defer mu.Unlock()
			if tc.wantRequest && len(requests) == 0 {
				t.Fatal("occupied tainted target did not perform its pre-delete read")
			}
			for _, request := range requests {
				if request.method != http.MethodGet {
					t.Fatalf("rejected tainted replacement made a write request: %#v", requests)
				}
			}
			if stored["q1"]["defaultTarget"] != float64(1) || stored["q2"]["defaultTarget"] != float64(0) || stored["LATEST"]["defaultTarget"] != float64(9) {
				t.Fatalf("rejected tainted replacement changed source or target: %#v", stored)
			}
		})
	}

	createData := schema.TestResourceDataRaw(t, r.Schema, config(&q2, 2))
	if err := r.Create(createData, client); err == nil || !strings.Contains(err.Error(), "refusing to overwrite an independent configuration") {
		t.Fatalf("create must reject an occupied qualifier before the upsert, got %v", err)
	}
	mu.Lock()
	defer mu.Unlock()
	for _, request := range requests {
		if request.method != http.MethodGet {
			t.Fatalf("rejected create made a write request: %#v", request)
		}
	}
	if stored["q1"]["defaultTarget"] != float64(1) || stored["q2"]["defaultTarget"] != float64(0) {
		t.Fatalf("rejected create changed source or target: %#v", stored)
	}
}

func TestUnitProvisionConfigFunctionNameImport(t *testing.T) {
	functionArn := "acs:fc:cn-hangzhou:1234567890123456:functions/unit-provision"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if qualifier := req.URL.Query().Get("qualifier"); qualifier != "" {
			t.Errorf("function-name import unexpectedly used qualifier %q", qualifier)
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"functionArn":   functionArn,
			"defaultTarget": float64(9),
		}); err != nil {
			t.Errorf("encode response: %v", err)
		}
	}))
	defer server.Close()

	r := fcv3.ResourceProvisionConfig()
	data := r.Data(&sdktf.InstanceState{
		ID: "unit-provision",
		Attributes: map[string]string{
			"id":            "unit-provision",
			"function_name": "unit-provision",
		},
	})
	if err := r.Read(data, provisionConfigUnitTestClient(t, server)); err != nil {
		t.Fatalf("function-name import read: %v", err)
	}
	if data.Id() != "unit-provision" || data.Get("function_name") != "unit-provision" {
		t.Fatalf("function-name identity changed: id=%q function_name=%#v", data.Id(), data.Get("function_name"))
	}
	if data.Get("qualifier") != "LATEST" || data.Get("target") != 9 {
		t.Fatalf("unexpected imported state: qualifier=%#v target=%#v", data.Get("qualifier"), data.Get("target"))
	}
}

// checkProvisionConfigDependenciesDestroyed verifies that the prerequisite
// resources created by the fixture are gone after the test teardown. The
// target CheckDestroy only inspects alicloud_fcv3_provision_config, so without
// this check a leaked function or SLS project/log store in the shared account
// would go unnoticed.
func checkProvisionConfigDependenciesDestroyed() resource.TestCheckFunc {
	return func(state *terraform.State) error {
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)
		fcv3Service := fcv3.NewService(client)
		slsService := sls.NewService(client)
		for _, rs := range state.RootModule().Resources {
			var describe func() error
			switch rs.Type {
			case "alicloud_fcv3_function":
				describe = func() error {
					_, err := fcv3Service.DescribeFunction(rs.Primary.ID)
					return err
				}
			case "alicloud_log_project":
				describe = func() error {
					_, err := slsService.DescribeLogProject(rs.Primary.ID)
					return err
				}
			case "alicloud_log_store":
				describe = func() error {
					_, err := slsService.DescribeLogStore(rs.Primary.ID)
					return err
				}
			case "alicloud_fcv3_alias":
				describe = func() error {
					_, err := fcv3Service.DescribeAlias(rs.Primary.ID)
					return err
				}
			case "alicloud_fcv3_function_version":
				describe = func() error {
					_, err := fcv3Service.DescribeFunctionVersion(rs.Primary.ID)
					return err
				}
			default:
				continue
			}
			if err := waitProvisionConfigDependencyDestroyed(rs.Type, rs.Primary.ID, describe); err != nil {
				return err
			}
		}
		return nil
	}
}

// waitProvisionConfigDependencyDestroyed tolerates the short window between a
// successful destroy call and the prerequisite disappearing from the service.
func waitProvisionConfigDependencyDestroyed(resourceType, id string, describe func() error) error {
	deadline := time.Now().Add(3 * time.Minute)
	for {
		err := describe()
		if err == nil {
			if time.Now().After(deadline) {
				return fmt.Errorf("%s %q still exists after destroy", resourceType, id)
			}
			time.Sleep(5 * time.Second)
			continue
		}
		if errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, "checking destroyed %s %q", resourceType, id)
	}
}

func provisionConfigUnitTestClient(t *testing.T, server *httptest.Server) *connectivity.AliyunClient {
	t.Helper()
	credential, err := credentials.NewCredential(new(credentials.Config).SetType("access_key").SetAccessKeyId("test-key").SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatal(err)
	}
	endpoint := strings.TrimPrefix(server.URL, "http://")
	t.Setenv("NO_PROXY", endpoint)
	endpoints := new(sync.Map)
	config := &connectivity.Config{
		AccessKey:            "test-key",
		SecretKey:            "test-secret",
		Credential:           credential,
		RegionId:             "cn-hangzhou",
		AccountId:            "1234567890123456",
		AccountType:          "test",
		Protocol:             "http",
		Endpoints:            endpoints,
		SignVersion:          new(sync.Map),
		SkipRegionValidation: true,
	}
	client, err := connectivity.NewClient(config)
	if err != nil {
		t.Fatal(err)
	}
	endpoints.Store("fc", endpoint)
	return client
}
