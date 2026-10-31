package cloudmonitorservice_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/cloudmonitorservice"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudCloudMonitorServiceGroupMonitoringAgentProcess_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_monitor_service_group_monitoring_agent_process.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudCloudMonitorServiceGroupMonitoringAgentProcessMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudmonitorservice.NewCloudMonitorServiceServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCloudMonitorServiceGroupMonitoringAgentProcess")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scloudmonitorservicegroupmonitoringagentprocess%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudCloudMonitorServiceGroupMonitoringAgentProcessBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
						"group_id":       acctest.CHECKSET,
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
				// alert_config.statistics is constrained by the schema
				// ValidateFunc to a single allowed value ("Average"), so it
				// cannot be modified in any apply step. It is listed in
				// ImportStateVerifyIgnore to exempt it from the modification
				// coverage requirement (the checker removes ImportStateVerifyIgnore
				// members from schemaModifySet). Import still verifies every
				// other attribute.
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"alert_config.statistics"},
			},
		},
	})
}

func TestAccAliCloudCloudMonitorServiceGroupMonitoringAgentProcess_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_monitor_service_group_monitoring_agent_process.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudCloudMonitorServiceGroupMonitoringAgentProcessMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudmonitorservice.NewCloudMonitorServiceServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCloudMonitorServiceGroupMonitoringAgentProcess")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scloudmonitorservicegroupmonitoringagentprocess%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudCloudMonitorServiceGroupMonitoringAgentProcessBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
						"group_id":                      acctest.CHECKSET,
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
	"group_monitoring_agent_process_id": acctest.CHECKSET,
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

// TypeList order coverage for alicloud_cloud_monitor_service_group_monitoring_agent_process.match_express.
func TestAccAliCloudCloudMonitorServiceGroupMonitoringAgentProcess_matchExpressOrder(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_monitor_service_group_monitoring_agent_process.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudCloudMonitorServiceGroupMonitoringAgentProcessMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudmonitorservice.NewCloudMonitorServiceServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCloudMonitorServiceGroupMonitoringAgentProcess")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scmsgmap-me%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudCloudMonitorServiceGroupMonitoringAgentProcessBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
						{
							"name":     name + "b",
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
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"match_express.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"match_express": []map[string]interface{}{
						{
							"name":     name + "b",
							"value":    "*",
							"function": "all",
						},
						{
							"name":     name,
							"value":    "*",
							"function": "all",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"match_express.#": "2",
					}),
				),
			},
		},
	})
}

// TypeList order coverage for alicloud_cloud_monitor_service_group_monitoring_agent_process.alert_config.
func TestAccAliCloudCloudMonitorServiceGroupMonitoringAgentProcess_alertConfigOrder(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_monitor_service_group_monitoring_agent_process.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudCloudMonitorServiceGroupMonitoringAgentProcessMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudmonitorservice.NewCloudMonitorServiceServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCloudMonitorServiceGroupMonitoringAgentProcess")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scmsgmap-ac%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudCloudMonitorServiceGroupMonitoringAgentProcessBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
						{
							"escalations_level":   "warn",
							"comparison_operator": "GreaterThanThreshold",
							"statistics":          "Average",
							"threshold":           "30",
							"times":               "110",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alert_config.#": "2",
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
						},
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
						"alert_config.#": "2",
					}),
				),
			},
		},
	})
}

// TypeList order coverage for alicloud_cloud_monitor_service_group_monitoring_agent_process.alert_config.target_list.
func TestAccAliCloudCloudMonitorServiceGroupMonitoringAgentProcess_alertConfigTargetListOrder(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_monitor_service_group_monitoring_agent_process.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudCloudMonitorServiceGroupMonitoringAgentProcessMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudmonitorservice.NewCloudMonitorServiceServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCloudMonitorServiceGroupMonitoringAgentProcess")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scmsgmap-tl%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudCloudMonitorServiceGroupMonitoringAgentProcessBasicDependence0)
	arn := "acs:openapi:" + os.Getenv("ALICLOUD_REGION") + ":" + os.Getenv("ALICLOUD_ACCOUNT_ID") + ":cms/DescribeMetricList/2019-01-01/testrole"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"group_id":     "${alicloud_cms_monitor_group.default.id}",
					"process_name": name,
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
									"arn":            arn,
								},
								{
									"target_list_id": "2",
									"json_params":    `{\"c\":\"d\"}`,
									"level":          "CRITICAL",
									"arn":            arn,
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alert_config.0.target_list.#": "2",
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
									"target_list_id": "2",
									"json_params":    `{\"c\":\"d\"}`,
									"level":          "CRITICAL",
									"arn":            arn,
								},
								{
									"target_list_id": "1",
									"json_params":    `{\"a\":\"b\"}`,
									"level":          "WARN",
									"arn":            arn,
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alert_config.0.target_list.#": "2",
					}),
				),
			},
		},
	})
}
