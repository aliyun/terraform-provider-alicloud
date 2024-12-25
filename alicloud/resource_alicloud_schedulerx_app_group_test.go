package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Schedulerx AppGroup. >>> Resource test cases, automatically generated.
// Case 预发环境_20241220_杭州reigon 9654
func TestAccAliCloudSchedulerxAppGroup_basic9654(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_schedulerx_app_group.default"
	ra := resourceAttrInit(resourceId, AlicloudSchedulerxAppGroupMap9654)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SchedulerxServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSchedulerxAppGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sschedulerxappgroup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSchedulerxAppGroupBasicDependence9654)
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
					"namespace":             "${alicloud_schedulerx_namespace.CreateNameSpace.namespace_uid}",
					"group_id":              "test-appgroup-pop-autotest",
					"description":           "appgroup 资源用例自动生成",
					"monitor_contacts_json": "[{\\\"userName\\\":\\\"张三\\\",\\\"userPhone\\\":\\\"89756******\\\"},{\\\"userName\\\":\\\"李四\\\",\\\"ding\\\":\\\"http://www.example.com\\\"}]",
					"enable_log":            "false",
					"app_name":              "test-appgroup-pop-autotest",
					"app_version":           "1",
					"namespace_name":        "default",
					"monitor_config_json":   "{\\\"sendChannel\\\":\\\"sms,ding\\\"}",
					"app_type":              "1",
					"max_jobs":              "100",
					"namespace_source":      "schedulerx",
					"schedule_busy_workers": "false",
					"delete_jobs":           "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"namespace":             CHECKSET,
						"group_id":              "test-appgroup-pop-autotest",
						"description":           "appgroup 资源用例自动生成",
						"monitor_contacts_json": "[{\"userName\":\"张三\",\"userPhone\":\"89756******\"},{\"userName\":\"李四\",\"ding\":\"http://www.example.com\"}]",
						"enable_log":            "false",
						"app_name":              "test-appgroup-pop-autotest",
						"app_version":           "1",
						"namespace_name":        "default",
						"monitor_config_json":   "{\"sendChannel\":\"sms,ding\"}",
						"app_type":              "1",
						"max_jobs":              "100",
						"namespace_source":      "schedulerx",
						"schedule_busy_workers": "false",
						"delete_jobs":           "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":     "appgroup 资源用例自动生成_update",
					"app_version":     "2",
					"max_concurrency": "500",
					"delete_jobs":     "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":     "appgroup 资源用例自动生成_update",
						"app_version":     "2",
						"max_concurrency": "500",
						"delete_jobs":     "true",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"app_type", "enable_log", "max_concurrency", "monitor_config_json", "monitor_contacts_json", "namespace_name", "namespace_source", "schedule_busy_workers", "delete_jobs"},
			},
		},
	})
}

var AlicloudSchedulerxAppGroupMap9654 = map[string]string{}

func AlicloudSchedulerxAppGroupBasicDependence9654(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_schedulerx_namespace" "CreateNameSpace" {
  namespace_name = var.name
  description    = "由appgroup 资源测试用例前置步骤创建"
}


`, name)
}

// Test Schedulerx AppGroup. <<< Resource test cases, automatically generated.
