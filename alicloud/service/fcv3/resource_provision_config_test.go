package fcv3_test

import (
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
// metric_type=CPUUtilization (CPU autoScaling requires an account-side
// allowlist), and its dependency creates an alicloud_log_project which hits
// ProjectQuotaExceed (account at the 300-project SLS quota). This case
// deliberately omits CPU/GPU autoScaling policies and uses a minimal function
// dependency without an SLS log_project.
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
			{
				Config: testAccConfig(map[string]interface{}{
					"always_allocate_cpu": "false",
					"always_allocate_gpu": "false",
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
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"current", "current_error"},
			},
			// Removing both arguments keeps the service-side values because the
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
// project). This avoids the account-side ProjectQuotaExceed (300-project SLS
// quota) that blocks the basic7182 dependency, while still providing a real
// function for the provision-config lifecycle target.
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

// TestUnitAliCloudFcv3ProvisionConfigRequest verifies the request payload built
// by the public CRUD entry points. Omitted Optional+Computed fields must be
// absent, an explicit false must be sent as false, and an update that only
// toggles one boolean must still carry the API-required target and the current
// values of the untouched fields.
func TestUnitAliCloudFcv3ProvisionConfigRequest(t *testing.T) {
	type capturedRequest struct {
		method string
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
		requests = append(requests, capturedRequest{method: req.Method, body: body})
		if req.Method == http.MethodPut {
			for _, key := range []string{"alwaysAllocateCPU", "alwaysAllocateGPU", "target"} {
				if value, ok := body[key]; ok {
					stored[key] = value
				}
			}
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
	update := func(current map[string]interface{}, changes map[string]interface{}) map[string]interface{} {
		t.Helper()
		data := schema.TestResourceDataRaw(t, r.Schema, current)
		data.SetId("unit-provision")
		for key, value := range changes {
			if err := data.Set(key, value); err != nil {
				t.Fatalf("set %s: %v", key, err)
			}
		}
		if err := r.Update(data, client); err != nil {
			t.Fatalf("update: %v", err)
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

	// Boolean-only update: target is unchanged but still present in the body,
	// and the untouched boolean keeps its configured value.
	body = update(map[string]interface{}{
		"function_name":       "unit-provision",
		"target":              1,
		"always_allocate_cpu": true,
		"always_allocate_gpu": true,
	}, map[string]interface{}{
		"always_allocate_cpu": false,
	})
	assertValue(body, "target", float64(1))
	assertValue(body, "alwaysAllocateCPU", false)
	assertValue(body, "alwaysAllocateGPU", true)

	// Explicit false for both fields must be sent verbatim.
	body = update(map[string]interface{}{
		"function_name":       "unit-provision",
		"target":              1,
		"always_allocate_cpu": false,
		"always_allocate_gpu": true,
	}, map[string]interface{}{
		"always_allocate_gpu": false,
	})
	assertValue(body, "alwaysAllocateCPU", false)
	assertValue(body, "alwaysAllocateGPU", false)
	assertValue(body, "target", float64(1))

	// Switching both back to true must also be sent verbatim.
	body = update(map[string]interface{}{
		"function_name":       "unit-provision",
		"target":              1,
		"always_allocate_cpu": false,
		"always_allocate_gpu": false,
	}, map[string]interface{}{
		"always_allocate_cpu": true,
		"always_allocate_gpu": true,
	})
	assertValue(body, "alwaysAllocateCPU", true)
	assertValue(body, "alwaysAllocateGPU", true)
	assertValue(body, "target", float64(1))
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
