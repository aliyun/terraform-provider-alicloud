package fcv3_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/fcv3"
	schema "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdktf "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
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
		CheckDestroy:             rac.CheckResourceDestroy(),
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
							"metric_type":   "CPUUtilization",
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
						"always_allocate_cpu":        "false",
						"target":                     "2",
						"target_tracking_policies.#": "1",
						"scheduled_actions.#":        "1",
						"always_allocate_gpu":        "false",
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
// Step 1 is the minimal legal configuration required by AC3 (function_name and
// target only). Later steps exercise the Optional+Computed
// always_allocate_cpu/always_allocate_gpu fields: explicit true, explicit false
// while target stays unchanged, both false, back to true, and finally removal.
// Every successful configuration step is followed by ImportStateVerify.
// The full-coverage case basic7182 exercises every optional field but cannot
// pass on this account: its target_tracking_policies use
// metric_type=CPUUtilization and CPU autoScaling requires an account-side
// allowlist. This case deliberately omits CPU/GPU autoScaling policies.
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
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			// AC3: the minimal legal configuration sets only function_name and
			// target; the optional always-allocate fields stay unset.
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
			// AC4: explicit false must be stored and converge. Only the two
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
			// AC4: switching back to true must take effect and converge again.
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
// project quota. The full-coverage case basic7182 is held back by its
// CPUUtilization tracking policy, not by the SLS dependency: memberOrder and
// memberOrderPair reuse that SLS-based dependency and pass.
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
		CheckDestroy:             rac.CheckResourceDestroy(),
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
		CheckDestroy:             rac.CheckResourceDestroy(),
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

// TestAccAliCloudFcv3ProvisionConfig_alwaysAllocateCoupling records the
// service-side coupling observed on a function without GPU: a request that
// sets always_allocate_cpu=false together with always_allocate_gpu=true is
// normalized to false/false, so the asymmetric configuration cannot converge.
// The persistent diff is the expected result here, not a Provider defect: the
// argument reference documents that both values must be configured together.
// This case exists to keep that limitation executable and to give AC4 a
// first-class counterexample (QA-F02); it does not claim AC4 is satisfied.
func TestAccAliCloudFcv3ProvisionConfig_alwaysAllocateCoupling(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_fcv3_provision_config.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudFcv3ProvisionConfigMap7183)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return fcv3.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeProvisionConfig")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sfcv3provisionalloc%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudFcv3ProvisionConfigBasicDependence7183)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"function_name":       "${alicloud_fcv3_function.function.function_name}",
					"target":              "1",
					"always_allocate_cpu": "false",
					"always_allocate_gpu": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"always_allocate_cpu": "false",
						"always_allocate_gpu": "false",
					}),
				),
				// The requested gpu=true is normalized to false, so the
				// configuration keeps a diff after apply by design.
				ExpectNonEmptyPlan: true,
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

// TestUnitAliCloudFcv3ProvisionConfigRequest verifies the request payload built
// by the public CRUD entry points. Omitted Optional+Computed fields must be
// absent, an explicit false must be sent as false, and an update that only
// toggles one boolean must still carry the API-required target and the current
// values of the untouched fields.
func TestUnitAliCloudFcv3ProvisionConfigRequest(t *testing.T) {
	type capturedRequest struct {
		method string
		query  string
		body   map[string]interface{}
	}

	var mu sync.Mutex
	var requests []capturedRequest
	stored := map[string]interface{}{
		"functionArn":       "acs:fc:cn-hangzhou:1234567890123456:functions/unit-provision",
		"alwaysAllocateCPU": true,
		"alwaysAllocateGPU": true,
		"target":            float64(1),
	}

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

		mu.Lock()
		requests = append(requests, capturedRequest{method: req.Method, query: req.URL.RawQuery, body: body})
		if req.Method == http.MethodPut {
			// PutProvisionConfig replaces the whole configuration: keys missing
			// from the request are gone from the response.
			replaced := map[string]interface{}{
				"functionArn": "acs:fc:cn-hangzhou:1234567890123456:functions/unit-provision",
			}
			for key, value := range body {
				replaced[key] = value
			}
			stored = replaced
		}
		response := make(map[string]interface{}, len(stored))
		for key, value := range stored {
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
	assertListLength := func(body map[string]interface{}, key string, want int) {
		t.Helper()
		values, ok := body[key].([]interface{})
		if !ok {
			t.Fatalf("request body %q is not a list: %#v", key, body[key])
		}
		if len(values) != want {
			t.Fatalf("request body %q has %d entries, want %d: %#v", key, len(values), want, values)
		}
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
		diff, err := r.Diff(ctx, prior, resourceConfig, nil)
		if err != nil {
			t.Fatalf("diff: %v", err)
		}
		if diff == nil || diff.Empty() {
			t.Fatalf("expected a plan for %#v", config)
		}
		next, diags := r.Apply(ctx, prior, diff, client)
		if diags.HasError() {
			t.Fatalf("apply: %v", diags)
		}
		again, err := r.Diff(ctx, next, resourceConfig, nil)
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
	assertMissing(body, "alwaysAllocateCPU")
	assertMissing(body, "alwaysAllocateGPU")

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

	// CR-04: an update whose configuration omits target (for example right
	// after an import, where the state value came from the service) must send
	// the current state value instead of the zero value of an unset argument.
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

	// The qualifier selects the target object and must reach the request query
	// on every replacement PUT that the resource performs.
	body = applyUpdate(priorState(map[string]string{
		"qualifier":           "LATEST",
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
			"metric_target": float64(1),
			"metric_type":   "ProvisionedConcurrencyUtilization",
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
		"target_tracking_policies.0.metric_target": "1",
		"target_tracking_policies.0.metric_type":   "ProvisionedConcurrencyUtilization",
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
	assertListLength(body, "targetTrackingPolicies", 1)
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
	client, err := config.Client()
	if err != nil {
		t.Fatal(err)
	}
	endpoints.Store("fc", endpoint)
	return client
}
