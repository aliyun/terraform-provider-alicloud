// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Sls LogtailPipelineConfig. >>> Resource test cases, automatically generated.
// Case LogtailPipelineConfig_basic_test 10001
func TestAccAliCloudSlsLogtailPipelineConfig_basic10001(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_sls_logtail_pipeline_config.default"
	ra := resourceAttrInit(resourceId, AlicloudSlsLogtailPipelineConfigMap10001)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSlsLogtailPipelineConfig")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccsls%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSlsLogtailPipelineConfigBasicDependence10001)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"config_name":  name,
					"project_name": "${alicloud_log_project.defaultProject.project_name}",
					"log_sample":   "2024-03-11 10:00:00 INFO test log message",
					"global": []map[string]interface{}{
						{
							"topic_type":                  "filepath",
							"topic_format":                "/var/log/(.*).log",
							"enable_timestamp_nanosecond": false,
						},
					},
					"inputs": `[
						{
							"Type": "input_file",
							"FilePaths": ["/var/log/app/*.log"]
						}
					]`,
					"processors": `[
						{
							"Type": "processor_parse_json_native",
							"SourceKey": "content"
						}
					]`,
					"flushers": `[
						{
							"Type": "flusher_sls",
							"Logstore": "test-logstore"
						}
					]`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"config_name":  name,
						"project_name": CHECKSET,
						"log_sample":   "2024-03-11 10:00:00 INFO test log message",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"log_sample": "2024-03-11 11:00:00 DEBUG updated log message",
					"processors": `[
						{
							"Type": "processor_parse_json_native",
							"SourceKey": "content"
						},
						{
							"Type": "processor_filter_regex",
							"Include": {"content": "ERROR.*"}
						}
					]`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"log_sample": "2024-03-11 11:00:00 DEBUG updated log message",
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

var AlicloudSlsLogtailPipelineConfigMap10001 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudSlsLogtailPipelineConfigBasicDependence10001(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "project_name" {
  default = "project-for-pipeline-terraform"
}

resource "alicloud_log_project" "defaultProject" {
  description  = "project for pipeline config terraform test"
  project_name = var.project_name
}

resource "alicloud_log_store" "defaultLogstore" {
  project_name = alicloud_log_project.defaultProject.project_name
  logstore_name = "test-logstore"
  retention_period = 30
  shard_count = 2
}
`, name)
}

// Test Sls LogtailPipelineConfig. <<< Resource test cases, automatically generated.
