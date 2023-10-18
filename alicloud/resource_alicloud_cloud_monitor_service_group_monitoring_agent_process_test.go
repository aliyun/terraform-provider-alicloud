package alicloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudCloudMonitorServiceGroupMonitoringAgentProcess_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_monitor_service_group_monitoring_agent_process.default"
	ra := resourceAttrInit(resourceId, AliCloudCloudMonitorServiceGroupMonitoringAgentProcessMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudMonitorServiceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudMonitorServiceGroupMonitoringAgentProcess")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scloudmonitorservicegroupmonitoringagentprocess%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCloudMonitorServiceGroupMonitoringAgentProcessBasicDependence0)
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
					"group_id":     "${alicloud_cms_monitor_group.default.id}",
					"process_name": name,
					"alert_config": []map[string]interface{}{
						{
							"escalations_level":   "critical",
							"comparison_operator": "GreaterThanOrEqualToThreshold",
							"statistics":          "Average",
							"threshold":           "20",
							"times":               "100",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_id":       CHECKSET,
						"process_name":   name,
						"alert_config.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"alert_config": []map[string]interface{}{
						{
							"escalations_level":   "warn",
							"comparison_operator": "GreaterThanThreshold",
							"statistics":          "Average",
							"threshold":           "30",
							"times":               "110",
							"effective_interval":  "00:00-22:59",
							"silence_time":        "85800",
							"webhook":             "https://www.aliyun.com",
							"target_list": []map[string]interface{}{
								{
									"target_list_id": "1",
									"json_params":    `{\"a\":\"b\"}`,
									"level":          "WARN",
									"arn":            "acs:openapi:" + os.Getenv("ALICLOUD_REGION") + ":" + os.Getenv("ALICLOUD_ACCOUNT_ID") + ":cms/DescribeMetricList/2019-01-01/testrole",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alert_config.#": "1",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAliCloudCloudMonitorServiceGroupMonitoringAgentProcess_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_monitor_service_group_monitoring_agent_process.default"
	ra := resourceAttrInit(resourceId, AliCloudCloudMonitorServiceGroupMonitoringAgentProcessMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudMonitorServiceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudMonitorServiceGroupMonitoringAgentProcess")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scloudmonitorservicegroupmonitoringagentprocess%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCloudMonitorServiceGroupMonitoringAgentProcessBasicDependence0)
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
					"group_id":                      "${alicloud_cms_monitor_group.default.id}",
					"process_name":                  name,
					"match_express_filter_relation": "or",
					"match_express": []map[string]interface{}{
						{
							"name":     name,
							"value":    "*",
							"function": "all",
						},
					},
					"alert_config": []map[string]interface{}{
						{
							"escalations_level":   "critical",
							"comparison_operator": "GreaterThanOrEqualToThreshold",
							"statistics":          "Average",
							"threshold":           "20",
							"times":               "100",
							"effective_interval":  "00:00-22:59",
							"silence_time":        "85800",
							"webhook":             "https://www.aliyun.com",
							"target_list": []map[string]interface{}{
								{
									"target_list_id": "1",
									"json_params":    `{\"a\":\"b\"}`,
									"level":          "WARN",
									"arn":            "acs:openapi:" + os.Getenv("ALICLOUD_REGION") + ":" + os.Getenv("ALICLOUD_ACCOUNT_ID") + ":cms/DescribeMetricList/2019-01-01/testrole",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_id":                      CHECKSET,
						"process_name":                  name,
						"match_express_filter_relation": "or",
						"match_express.#":               "1",
						"alert_config.#":                "1",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var AliCloudCloudMonitorServiceGroupMonitoringAgentProcessMap0 = map[string]string{
	"group_monitoring_agent_process_id": CHECKSET,
}

func AliCloudCloudMonitorServiceGroupMonitoringAgentProcessBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	resource "alicloud_cms_alarm_contact_group" "default" {
  		alarm_contact_group_name = var.name
  		contacts                 = ["user", "user1", "user2"]
	}

	resource "alicloud_cms_monitor_group" "default" {
  		monitor_group_name = var.name
  		contact_groups     = [alicloud_cms_alarm_contact_group.default.id]
	}
`, name)
}
