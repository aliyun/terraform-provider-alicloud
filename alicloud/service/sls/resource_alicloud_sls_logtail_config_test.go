// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package sls_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/sls"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test Sls LogtailConfig. >>> Resource test cases, automatically generated.
// Case Logtailconfig_terraform_覆盖度 10828
func TestAccAliCloudSlsLogtailConfig_basic10828(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_sls_logtail_config.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudSlsLogtailConfigMap10828)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return sls.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeSlsLogtailConfig")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccsls%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudSlsLogtailConfigBasicDependence10828)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"logtail_config_name": name,
					"input_type":          "file",
					"project_name":        "${alicloud_log_project.defaultuA28zS.project_name}",
					"output_detail": []map[string]interface{}{
						{
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
						"project_name":        acctest.CHECKSET,
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

func TestAccAliCloudSlsLogtailConfig_attributes(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_sls_logtail_config.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudSlsLogtailConfigMap10828)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return sls.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeSlsLogtailConfig")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccsls%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudSlsLogtailConfigBasicDependence10828)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"logtail_config_name": name,
					"input_type":          "file",
					"project_name":        "${alicloud_log_project.defaultuA28zS.project_name}",
					"create_time":         "1700000000",
					"last_modify_time":    1700000000,
					"log_sample":          "test-sample-log-content",
					"output_detail": []map[string]interface{}{
						{
							"logstore_name": "test",
							"endpoint":      "cn-hangzhou-intranet.log.aliyuncs.com",
							"region":        "cn-hangzhou",
						},
					},
					"output_type":  "LogService",
					"input_detail": "{\\\"adjustTimezone\\\":false,\\\"delayAlarmBytes\\\":0,\\\"delaySkipBytes\\\":0,\\\"discardNonUtf8\\\":false,\\\"discardUnmatch\\\":true,\\\"dockerFile\\\":false,\\\"enableRawLog\\\":false,\\\"enableTag\\\":false,\\\"fileEncoding\\\":\\\"utf8\\\",\\\"filePattern\\\":\\\"access*.log\\\",\\\"filterKey\\\":[\\\"key1\\\"],\\\"filterRegex\\\":[\\\"regex1\\\"],\\\"key\\\":[\\\"key1\\\",\\\"key2\\\"],\\\"localStorage\\\":true,\\\"logBeginRegex\\\":\\\".*\\\",\\\"logPath\\\":\\\"/var/log/httpd\\\",\\\"logTimezone\\\":\\\"\\\",\\\"logType\\\":\\\"common_reg_log\\\",\\\"maxDepth\\\":1000,\\\"maxSendRate\\\":-1,\\\"mergeType\\\":\\\"topic\\\",\\\"preserve\\\":true,\\\"preserveDepth\\\":0,\\\"priority\\\":0,\\\"regex\\\":\\\"(w+)(s+)\\\",\\\"sendRateExpire\\\":0,\\\"sensitive_keys\\\":[],\\\"tailExisted\\\":false,\\\"timeFormat\\\":\\\"%Y/%m/%d %H:%M:%S\\\",\\\"timeKey\\\":\\\"time\\\",\\\"topicFormat\\\":\\\"none\\\"}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"logtail_config_name": name,
						"input_type":          "file",
						"project_name":        acctest.CHECKSET,
						"output_type":         "LogService",
					}),
				),
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
  name        = var.name
}


`, name)
}

// Test Sls LogtailConfig. <<< Resource test cases, automatically generated.
