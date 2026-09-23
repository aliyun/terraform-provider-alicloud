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
	ra := resourceAttrInit(resourceId, AliCloudSchedulerxAppGroupMap9957)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SchedulerxServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSchedulerxAppGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sschedulerxappgroup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudSchedulerxAppGroupBasicDependence9957)
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
					"app_name":       name,
					"group_id":       name,
					"namespace":      "${alicloud_schedulerx_namespace.CreateNameSpace.namespace_uid}",
					"namespace_name": "default",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"app_name":  name,
						"group_id":  name,
						"namespace": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"app_version": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"app_version": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_log": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_log": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"monitor_config_json": "{\\\"sendChannel\\\":\\\"sms,ding\\\",\\\"alarmType\\\": \\\"Contacts\\\",\\\"webhookIsAtAll\\\": \\\"false\\\"}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"monitor_config_json": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"monitor_contacts_json": "[{\\\"name\\\":\\\"用户-手机\\\"},{\\\"name\\\":\\\"用户-钉钉\\\"}]",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"monitor_contacts_json": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"app_type", "delete_jobs", "max_concurrency", "namespace_name", "namespace_source", "schedule_busy_workers"},
			},
		},
	})
}

func TestAccAliCloudSchedulerxAppGroup_basic9957_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_schedulerx_app_group.default"
	ra := resourceAttrInit(resourceId, AliCloudSchedulerxAppGroupMap9957)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SchedulerxServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSchedulerxAppGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sschedulerxappgroup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudSchedulerxAppGroupBasicDependence9957)
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
					"group_id":              name,
					"description":           name,
					"monitor_contacts_json": "[{\\\"name\\\":\\\"用户-手机\\\"},{\\\"name\\\":\\\"用户-钉钉\\\"}]",
					"enable_log":            "true",
					"app_name":              name,
					"app_version":           "1",
					"namespace_name":        "default",
					"monitor_config_json":   "{\\\"sendChannel\\\":\\\"sms,ding\\\",\\\"alarmType\\\": \\\"Contacts\\\",\\\"webhookIsAtAll\\\": \\\"false\\\"}",
					"app_type":              "1",
					"max_concurrency":       "500",
					"max_jobs":              "100",
					"namespace_source":      "schedulerx",
					"schedule_busy_workers": "true",
					"delete_jobs":           "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"namespace":             CHECKSET,
						"group_id":              CHECKSET,
						"description":           name,
						"monitor_contacts_json": CHECKSET,
						"enable_log":            "true",
						"app_name":              name,
						"app_version":           "1",
						"max_jobs":              "100",
						"monitor_config_json":   CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"app_type", "delete_jobs", "max_concurrency", "namespace_name", "namespace_source", "schedule_busy_workers"},
			},
		},
	})
}

var AliCloudSchedulerxAppGroupMap9957 = map[string]string{
	"app_version": CHECKSET,
	"max_jobs":    CHECKSET,
}

func AliCloudSchedulerxAppGroupBasicDependence9957(name string) string {
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
