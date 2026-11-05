// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package config_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/config"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test Config ReportTemplate. >>> Resource test cases, automatically generated.
// Case 共亿测试ReportTemplate的第一个测试case 11967
func TestAccAliCloudConfigReportTemplate_basic11967(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_config_report_template.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudConfigReportTemplateMap11967)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return config.NewConfigServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeConfigReportTemplate")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccconfig%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudConfigReportTemplateBasicDependence11967)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"report_granularity": "AllInOne",
					"report_scope": []map[string]interface{}{
						{
							"key":        "RuleId",
							"value":      "cr-xxx",
							"match_type": "In",
						},
					},
					"report_file_formats":         "excel",
					"report_template_name":        name,
					"report_template_description": "test-desc",
					"subscription_frequency":      " ",
					"report_language":             "en-US",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"report_granularity":          "AllInOne",
						"report_scope.#":              "1",
						"report_file_formats":         "excel",
						"report_template_name":        name,
						"report_template_description": "test-desc",
						"subscription_frequency":      " ",
						"report_language":             "en-US",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"report_granularity": "GroupByAccount",
					"report_scope": []map[string]interface{}{
						{
							"key":        "CompliancePackId",
							"value":      "cp-xxxx",
							"match_type": "In",
						},
					},
					"report_template_name":        name + "_update",
					"report_template_description": "desc-updated",
					"subscription_frequency":      "0 0 0 * * ?",
					"report_language":             "zh-CN",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"report_granularity":          "GroupByAccount",
						"report_scope.#":              "1",
						"report_template_name":        name + "_update",
						"report_template_description": "desc-updated",
						"subscription_frequency":      "0 0 0 * * ?",
						"report_language":             "zh-CN",
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

var AlicloudConfigReportTemplateMap11967 = map[string]string{}

func AlicloudConfigReportTemplateBasicDependence11967(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Test Config ReportTemplate. <<< Resource test cases, automatically generated.
