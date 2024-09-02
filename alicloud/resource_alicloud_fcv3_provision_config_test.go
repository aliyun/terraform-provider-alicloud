package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Fcv3 ProvisionConfig. >>> Resource test cases, automatically generated.
// Case ProvisionConfig_Base 7182
func TestAccAliCloudFcv3ProvisionConfig_basic7182(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_fcv3_provision_config.default"
	ra := resourceAttrInit(resourceId, AlicloudFcv3ProvisionConfigMap7182)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &Fcv3ServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeFcv3ProvisionConfig")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sfcv3provisionconfig%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudFcv3ProvisionConfigBasicDependence7182)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
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
						"function_name":              CHECKSET,
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
  project_name = var.name
  description  = var.name
}

resource "alicloud_log_store" "default" {
  project_name          = alicloud_log_project.default.name
  logstore_name         = var.name
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
