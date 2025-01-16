package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Schedulerx AppGroup. >>> Resource test cases, automatically generated.
// Case 预发环境_20250110_乌兰察布(代码只部署到乌兰察布，用这个用例测试) 9957
func TestAccAliCloudSchedulerxAppGroup_basic9957(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_schedulerx_app_group.default"
	ra := resourceAttrInit(resourceId, AlicloudSchedulerxAppGroupMap9957)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SchedulerxServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSchedulerxAppGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sschedulerxappgroup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSchedulerxAppGroupBasicDependence9957)
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
					"monitor_contacts_json": "[{\\\"name\\\":\\\"用户-手机\\\"},{\\\"name\\\":\\\"用户-钉钉\\\"}]",
					"enable_log":            "false",
					"app_name":              "test-appgroup-pop-autotest",
					"app_version":           "1",
					"namespace_name":        "default",
					"monitor_config_json":   "{\\\"sendChannel\\\":\\\"sms,ding\\\",\\\"alarmType\\\": \\\"Contacts\\\",\\\"webhookIsAtAll\\\": \\\"false\\\"}",
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
						"monitor_contacts_json": CHECKSET,
						"enable_log":            "false",
						"app_name":              "test-appgroup-pop-autotest",
						"app_version":           CHECKSET,
						"namespace_name":        "default",
						"monitor_config_json":   CHECKSET,
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
					"description":           "appgroup 资源用例自动生成_update",
					"monitor_contacts_json": "[{\\\"name\\\":\\\"用户-飞书\\\"},{\\\"name\\\":\\\"用户-钉钉\\\"}]",
					"app_version":           "2",
					"monitor_config_json":   "{\\\"sendChannel\\\":\\\"sms,ding\\\",\\\"alarmType\\\": \\\"Contacts\\\",\\\"webhookIsAtAll\\\": \\\"true\\\"}",
					"max_concurrency":       "500",
					"delete_jobs":           "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":           "appgroup 资源用例自动生成_update",
						"monitor_contacts_json": CHECKSET,
						"app_version":           CHECKSET,
						"monitor_config_json":   CHECKSET,
						"max_concurrency":       "500",
						"delete_jobs":           "true",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"app_type", "enable_log", "max_concurrency", "namespace_name", "namespace_source", "schedule_busy_workers", "delete_jobs"},
			},
		},
	})
}

var AlicloudSchedulerxAppGroupMap9957 = map[string]string{}

func AlicloudSchedulerxAppGroupBasicDependence9957(name string) string {
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
