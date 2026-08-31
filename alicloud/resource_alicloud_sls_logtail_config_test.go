// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Sls LogtailConfig. >>> Resource test cases, automatically generated.
// Case Logtailconfig_terraform_覆盖度 10828
func TestAccAliCloudSlsLogtailConfig_basic10828(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_sls_logtail_config.default"
	ra := resourceAttrInit(resourceId, AlicloudSlsLogtailConfigMap10828)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSlsLogtailConfig")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccsls%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSlsLogtailConfigBasicDependence10828)
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
					"logtail_config_name": name,
					"input_type":          "file",
					"project_name":        "${alicloud_log_project.defaultuA28zS.project_name}",
					"output_detail": []map[string]interface{}{
						{
							"endpoint":      "cn-hangzhou-intranet.log.aliyuncs.com",
							"region":        "cn-hangzhou",
							"logstore_name": "test",
						},
					},
					"output_type":  "LogService",
					"input_detail": "{\\\"adjustTimezone\\\":false,\\\"delayAlarmBytes\\\":0,\\\"delaySkipBytes\\\":0,\\\"discardNonUtf8\\\":false,\\\"discardUnmatch\\\":true,\\\"dockerFile\\\":false,\\\"enableRawLog\\\":false,\\\"enableTag\\\":false,\\\"fileEncoding\\\":\\\"utf8\\\",\\\"filePattern\\\":\\\"access*.log\\\",\\\"filterKey\\\":[\\\"key1\\\"],\\\"filterRegex\\\":[\\\"regex1\\\"],\\\"key\\\":[\\\"key1\\\",\\\"key2\\\"],\\\"localStorage\\\":true,\\\"logBeginRegex\\\":\\\".*\\\",\\\"logPath\\\":\\\"/var/log/httpd\\\",\\\"logTimezone\\\":\\\"\\\",\\\"logType\\\":\\\"common_reg_log\\\",\\\"maxDepth\\\":1000,\\\"maxSendRate\\\":-1,\\\"mergeType\\\":\\\"topic\\\",\\\"preserve\\\":true,\\\"preserveDepth\\\":0,\\\"priority\\\":0,\\\"regex\\\":\\\"(w+)(s+)\\\",\\\"sendRateExpire\\\":0,\\\"sensitive_keys\\\":[],\\\"tailExisted\\\":false,\\\"timeFormat\\\":\\\"%Y/%m/%d %H:%M:%S\\\",\\\"timeKey\\\":\\\"time\\\",\\\"topicFormat\\\":\\\"none\\\"}",
					//"input_detail": map[string]interface{}{
					//	"\"logType\"":       "common_reg_log",
					//	"\"logPath\"":       "/var/log/httpd",
					//	"\"filePattern\"":   "access*.log",
					//	"\"localStorage\"":  true,
					//	"\"timeFormat\"":    "%Y/%m/%d %H:%M:%S",
					//	"\"logBeginRegex\"": ".*",
					//	"\"regex\"":         "(w+)(s+)",
					//	"\"key\"":           "[\\\"key1\\\",\\\"key2\\\"]",
					//	"\"filterKey\"":     "[\\\"key1\\\"]",
					//	"\"filterRegex\"":   "[\\\"regex1\\\"]",
					//	"\"fileEncoding\"":  "utf8",
					//	"\"topicFormat\"":   "none",
					//},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"logtail_config_name": name,
						"input_type":          "file",
						"project_name":        CHECKSET,
						"output_type":         "LogService",
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

var AlicloudSlsLogtailConfigMap10828 = map[string]string{}

func AlicloudSlsLogtailConfigBasicDependence10828(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "project_name" {
  default = "project-for-logtail-terraform"
}

resource "alicloud_log_project" "defaultuA28zS" {
  description = "project for terrafrom test"
  name        = var.project_name
}


`, name)
}

// Test Sls LogtailConfig. <<< Resource test cases, automatically generated.
