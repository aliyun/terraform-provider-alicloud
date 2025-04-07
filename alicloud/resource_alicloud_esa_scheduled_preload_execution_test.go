package alicloud

import (
	"fmt"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

// Test ESA ScheduledPreloadExecution. >>> Resource test cases, automatically generated.
// Case scheduledpreloadexecution_test
func TestAccAliCloudESAScheduledPreloadExecutionscheduledpreloadexecution_test(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_scheduled_preload_execution.default"
	ra := resourceAttrInit(resourceId, AliCloudESAScheduledPreloadExecutionscheduledpreloadexecution_testMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaScheduledPreloadExecution")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sESAScheduledPreloadExecution%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudESAScheduledPreloadExecutionscheduledpreloadexecution_testBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"slice_len":                "5",
					"end_time":                 "2024-06-04T10:02:09.000+08:00",
					"start_time":               "2024-06-04T00:00:00.000+08:00",
					"scheduled_preload_job_id": "${alicloud_esa_scheduled_preload_job.default.scheduled_preload_job_id}",
					"interval":                 "30",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"interval": "60",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"slice_len": "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"start_time": "2024-05-31T17:10:48.849+08:00",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"end_time": "2024-05-31T18:10:48.849+08:00",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
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

var AliCloudESAScheduledPreloadExecutionscheduledpreloadexecution_testMap = map[string]string{
	"id": CHECKSET,
}

func AliCloudESAScheduledPreloadExecutionscheduledpreloadexecution_testBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
}

resource "alicloud_esa_site" "default" {
  site_name   = "terraform.cn"
  instance_id = data.alicloud_esa_sites.default.sites.0.instance_id
  coverage    = "overseas"
  access_type = "NS"
}

resource "alicloud_esa_scheduled_preload_job" "default" {
  insert_way                 = "textBox"
  site_id                    = alicloud_esa_site.default.id
  scheduled_preload_job_name = "example_scheduledpreloadexecution_job"
  url_list                   = "http://example.gositecdn.cn/example/example.txt"
}

`, name)
}

// Test ESA ScheduledPreloadExecution. <<< Resource test cases, automatically generated.
