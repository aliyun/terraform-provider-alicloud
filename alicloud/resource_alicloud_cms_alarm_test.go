package alicloud

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_cms_alarm", &resource.Sweeper{
		Name: "alicloud_cms_alarm",
		F:    testSweepCMSAlarms,
	})
}

func testSweepCMSAlarms(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	var alarms []cms.AlarmInDescribeMetricRuleList
	req := cms.CreateDescribeMetricRuleListRequest()
	req.RegionId = client.RegionId
	req.PageSize = requests.NewInteger(PageSizeLarge)
	req.Page = requests.NewInteger(1)
	for {
		raw, err := client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
			return cmsClient.DescribeMetricRuleList(req)
		})
		if err != nil {
			log.Printf("[ERROR] Error retrieving CMS Alarm: %s", err)
		}
		resp, _ := raw.(*cms.DescribeMetricRuleListResponse)
		if resp == nil || len(resp.Alarms.Alarm) < 1 {
			break
		}
		alarms = append(alarms, resp.Alarms.Alarm...)

		if len(resp.Alarms.Alarm) < PageSizeLarge {
			break
		}
		if page, err := getNextpageNumber(req.Page); err != nil {
			return WrapError(err)
		} else {
			req.Page = page
		}
	}

	for _, v := range alarms {
		name := v.RuleName
		id := v.RuleId
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip && v.AlertState == "INSUFFICIENT_DATA" {
			skip = false
		}
		if skip {
			log.Printf("[INFO] Skipping CMS Alarm: %s (%s)", name, id)
			continue
		}

		log.Printf("[INFO] Deleting CMS Alarm: %s (%s). Status: %s", name, id, v.AlertState)
		req := cms.CreateDeleteMetricRulesRequest()
		req.Id = &[]string{id}
		_, err := client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
			return cmsClient.DeleteMetricRules(req)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete CMS Alarm (%s (%s)): %s", name, id, err)
		}
	}
	return nil
}

func TestAccAliCloudCmsAlarm_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_alarm.default"
	ra := resourceAttrInit(resourceId, AliCloudCmsAlarmMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAlarm")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%s-rule-name%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCmsAlarmBasicDependence0)
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
					"name":           name,
					"project":        "acs_ecs_dashboard",
					"metric":         "disk_writebytes",
					"contact_groups": []string{"${alicloud_cms_monitor_group.default.id}"},
					"escalations_critical": []map[string]interface{}{
						{
							"statistics": "Average",
							"threshold":  "90",
							"times":      "1",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                   name,
						"project":                "acs_ecs_dashboard",
						"metric":                 "disk_writebytes",
						"contact_groups.#":       "1",
						"escalations_critical.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"contact_groups": []string{"${alicloud_cms_monitor_group.update.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"contact_groups.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"metric_dimensions": `[{\"instanceId\":\"` + "${alicloud_instance.default.id}" + `\"},{\"device\":\"/dev/vda1\"}]`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"metric_dimensions": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"effective_interval": "06:00-20:00",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"effective_interval": "06:00-20:00",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"period": "900",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"period": "900",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"silence_time": "300",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"silence_time": "300",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"webhook": "https://www.aliyun.com",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"webhook": "https://www.aliyun.com",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enabled": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enabled": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Alarm",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Alarm",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"escalations_critical": []map[string]interface{}{
						{
							"comparison_operator": ">",
							"statistics":          "Maximum",
							"threshold":           "35",
							"times":               "2",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"escalations_critical.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"escalations_info": []map[string]interface{}{
						{
							"comparison_operator": ">=",
							"statistics":          "Minimum",
							"threshold":           "20",
							"times":               "3",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"escalations_info.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"escalations_warn": []map[string]interface{}{
						{
							"comparison_operator": "<",
							"statistics":          "Average",
							"threshold":           "30",
							"times":               "5",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"escalations_warn.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"targets": []map[string]interface{}{
						{
							"target_id":   "1",
							"json_params": `{\"a\":\"b\"}`,
							"level":       "Warn",
							"arn":         "acs:openapi:" + os.Getenv("ALICLOUD_REGION") + ":" + os.Getenv("ALICLOUD_ACCOUNT_ID") + ":cms/DescribeMetricList/2019-01-01/testrole",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"targets.#": "1",
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

func TestAccAliCloudCmsAlarm_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_alarm.default"
	ra := resourceAttrInit(resourceId, AliCloudCmsAlarmMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAlarm")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%s-rule-name%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCmsAlarmBasicDependence0)
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
					"name":               name,
					"project":            "acs_ecs_dashboard",
					"metric":             "disk_writebytes",
					"contact_groups":     []string{"${alicloud_cms_monitor_group.default.id}"},
					"metric_dimensions":  `[{\"instanceId\":\"` + "${alicloud_instance.default.id}" + `\"},{\"device\":\"/dev/vda1\"}]`,
					"effective_interval": "06:00-20:00",
					"period":             "900",
					"silence_time":       "300",
					"webhook":            "https://www.aliyun.com",
					"enabled":            "true",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Alarm",
					},
					"escalations_critical": []map[string]interface{}{
						{
							"comparison_operator": "<=",
							"statistics":          "Average",
							"threshold":           "90",
							"times":               "1",
						},
					},
					"escalations_info": []map[string]interface{}{
						{
							"comparison_operator": "!=",
							"statistics":          "Minimum",
							"threshold":           "20",
							"times":               "3",
						},
					},
					"escalations_warn": []map[string]interface{}{
						{
							"comparison_operator": "GreaterThanYesterday",
							"statistics":          "Average",
							"threshold":           "30",
							"times":               "5",
						},
					},
					"targets": []map[string]interface{}{
						{
							"target_id":   "1",
							"json_params": `{\"a\":\"b\"}`,
							"level":       "Warn",
							"arn":         "acs:openapi:" + os.Getenv("ALICLOUD_REGION") + ":" + os.Getenv("ALICLOUD_ACCOUNT_ID") + ":cms/DescribeMetricList/2019-01-01/testrole",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                   name,
						"project":                "acs_ecs_dashboard",
						"metric":                 "disk_writebytes",
						"contact_groups.#":       "1",
						"metric_dimensions":      CHECKSET,
						"effective_interval":     "06:00-20:00",
						"period":                 "900",
						"silence_time":           "300",
						"webhook":                "https://www.aliyun.com",
						"enabled":                "true",
						"targets.#":              "1",
						"tags.%":                 "2",
						"tags.Created":           "TF",
						"tags.For":               "Alarm",
						"escalations_critical.#": "1",
						"escalations_info.#":     "1",
						"escalations_warn.#":     "1",
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

func TestAccAliCloudCmsAlarm_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_alarm.default"
	ra := resourceAttrInit(resourceId, AliCloudCmsAlarmMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAlarm")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%s-rule-name%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCmsAlarmBasicDependence0)
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
					"name":           name,
					"project":        "acs_prometheus",
					"metric":         "AliyunEcs_cpu_total",
					"contact_groups": []string{"${alicloud_cms_monitor_group.default.id}"},
					"prometheus": []map[string]interface{}{
						{
							"prom_ql": name,
							"level":   "Critical",
							"times":   "1",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":             name,
						"project":          "acs_prometheus",
						"metric":           "AliyunEcs_cpu_total",
						"contact_groups.#": "1",
						"prometheus.#":     "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"contact_groups": []string{"${alicloud_cms_monitor_group.update.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"contact_groups.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"effective_interval": "06:00-20:00",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"effective_interval": "06:00-20:00",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"period": "900",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"period": "900",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"silence_time": "300",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"silence_time": "300",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"webhook": "https://www.aliyun.com",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"webhook": "https://www.aliyun.com",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enabled": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enabled": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Alarm",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Alarm",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"prometheus": []map[string]interface{}{
						{
							"prom_ql": name + "update",
							"level":   "Warn",
							"times":   "2",
							"annotations": map[string]string{
								"Created": "TF",
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"prometheus.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"targets": []map[string]interface{}{
						{
							"target_id":   "1",
							"json_params": `{\"a\":\"b\"}`,
							"level":       "Warn",
							"arn":         "acs:openapi:" + os.Getenv("ALICLOUD_REGION") + ":" + os.Getenv("ALICLOUD_ACCOUNT_ID") + ":cms/DescribeMetricList/2019-01-01/testrole",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"targets.#": "1",
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

func TestAccAliCloudCmsAlarm_basic1_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_alarm.default"
	ra := resourceAttrInit(resourceId, AliCloudCmsAlarmMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAlarm")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%s-rule-name%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCmsAlarmBasicDependence0)
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
					"name":               name,
					"project":            "acs_prometheus",
					"metric":             "AliyunEcs_cpu_total",
					"contact_groups":     []string{"${alicloud_cms_monitor_group.default.id}"},
					"effective_interval": "06:00-20:00",
					"period":             "900",
					"silence_time":       "300",
					"webhook":            "https://www.aliyun.com",
					"enabled":            "true",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Alarm",
					},
					"prometheus": []map[string]interface{}{
						{
							"prom_ql": name,
							"level":   "Critical",
							"times":   "1",
							"annotations": map[string]string{
								"Created": "TF",
							},
						},
					},
					"targets": []map[string]interface{}{
						{
							"target_id":   "1",
							"json_params": `{\"a\":\"b\"}`,
							"level":       "Warn",
							"arn":         "acs:openapi:" + os.Getenv("ALICLOUD_REGION") + ":" + os.Getenv("ALICLOUD_ACCOUNT_ID") + ":cms/DescribeMetricList/2019-01-01/testrole",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":               name,
						"project":            "acs_prometheus",
						"metric":             "AliyunEcs_cpu_total",
						"contact_groups.#":   "1",
						"effective_interval": "06:00-20:00",
						"period":             "900",
						"silence_time":       "300",
						"webhook":            "https://www.aliyun.com",
						"enabled":            "true",
						"prometheus.#":       "1",
						"targets.#":          "1",
						"tags.%":             "2",
						"tags.Created":       "TF",
						"tags.For":           "Alarm",
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

func TestAccAliCloudCmsAlarm_basic2(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_alarm.default"
	ra := resourceAttrInit(resourceId, AliCloudCmsAlarmMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAlarm")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%s-rule-name%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCmsAlarmBasicDependence0)
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
					"name":           name,
					"project":        "acs_ecs_dashboard",
					"metric":         "disk_writebytes",
					"contact_groups": []string{"${alicloud_cms_monitor_group.default.id}"},
					"escalations_critical": []map[string]interface{}{
						{
							"comparison_operator": "LessThanYesterday",
							"statistics":          "Average",
							"threshold":           "90",
							"times":               "1",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                   name,
						"project":                "acs_ecs_dashboard",
						"metric":                 "disk_writebytes",
						"contact_groups.#":       "1",
						"escalations_critical.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"contact_groups": []string{"${alicloud_cms_monitor_group.update.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"contact_groups.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"dimensions": map[string]string{
						"instanceId": "${alicloud_instance.default.id}",
						"device":     "/dev/vda1",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dimensions.%":          "2",
						"dimensions.instanceId": CHECKSET,
						"dimensions.device":     "/dev/vda1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"start_time": "6",
					"end_time":   "20",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"start_time": "6",
						"end_time":   "20",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"period": "900",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"period": "900",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"silence_time": "300",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"silence_time": "300",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"webhook": "https://www.aliyun.com",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"webhook": "https://www.aliyun.com",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enabled": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enabled": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Alarm",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Alarm",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"escalations_critical": []map[string]interface{}{
						{
							"comparison_operator": "GreaterThanLastWeek",
							"statistics":          "Maximum",
							"threshold":           "35",
							"times":               "2",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"escalations_critical.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"escalations_info": []map[string]interface{}{
						{
							"comparison_operator": "LessThanLastWeek",
							"statistics":          "Minimum",
							"threshold":           "20",
							"times":               "3",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"escalations_info.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"escalations_warn": []map[string]interface{}{
						{
							"comparison_operator": "GreaterThanLastPeriod",
							"statistics":          "Average",
							"threshold":           "30",
							"times":               "5",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"escalations_warn.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"targets": []map[string]interface{}{
						{
							"target_id":   "1",
							"json_params": `{\"a\":\"b\"}`,
							"level":       "Warn",
							"arn":         "acs:openapi:" + os.Getenv("ALICLOUD_REGION") + ":" + os.Getenv("ALICLOUD_ACCOUNT_ID") + ":cms/DescribeMetricList/2019-01-01/testrole",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"targets.#": "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"start_time", "end_time"},
			},
		},
	})
}

func TestAccAliCloudCmsAlarm_basic2_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_alarm.default"
	ra := resourceAttrInit(resourceId, AliCloudCmsAlarmMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAlarm")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%s-rule-name%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCmsAlarmBasicDependence0)
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
					"name":           name,
					"project":        "acs_ecs_dashboard",
					"metric":         "disk_writebytes",
					"contact_groups": []string{"${alicloud_cms_monitor_group.default.id}"},
					"dimensions": map[string]string{
						"instanceId": "${alicloud_instance.default.id}",
						"device":     "/dev/vda1",
					},
					"start_time":   "6",
					"end_time":     "20",
					"period":       "900",
					"silence_time": "300",
					"webhook":      "https://www.aliyun.com",
					"enabled":      "true",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Alarm",
					},
					"escalations_critical": []map[string]interface{}{
						{
							"comparison_operator": "LessThanLastPeriod",
							"statistics":          "Average",
							"threshold":           "90",
							"times":               "1",
						},
					},
					"escalations_info": []map[string]interface{}{
						{
							"comparison_operator": ">",
							"statistics":          "Minimum",
							"threshold":           "20",
							"times":               "3",
						},
					},
					"escalations_warn": []map[string]interface{}{
						{
							"comparison_operator": ">=",
							"statistics":          "Average",
							"threshold":           "30",
							"times":               "5",
						},
					},
					"targets": []map[string]interface{}{
						{
							"target_id":   "1",
							"json_params": `{\"a\":\"b\"}`,
							"level":       "Warn",
							"arn":         "acs:openapi:" + os.Getenv("ALICLOUD_REGION") + ":" + os.Getenv("ALICLOUD_ACCOUNT_ID") + ":cms/DescribeMetricList/2019-01-01/testrole",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                   name,
						"project":                "acs_ecs_dashboard",
						"metric":                 "disk_writebytes",
						"contact_groups.#":       "1",
						"dimensions.%":           "2",
						"dimensions.instanceId":  CHECKSET,
						"dimensions.device":      "/dev/vda1",
						"start_time":             "6",
						"end_time":               "20",
						"period":                 "900",
						"silence_time":           "300",
						"webhook":                "https://www.aliyun.com",
						"enabled":                "true",
						"targets.#":              "1",
						"tags.%":                 "2",
						"tags.Created":           "TF",
						"tags.For":               "Alarm",
						"escalations_critical.#": "1",
						"escalations_info.#":     "1",
						"escalations_warn.#":     "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"start_time", "end_time"},
			},
		},
	})
}

var AliCloudCmsAlarmMap0 = map[string]string{
	"status": CHECKSET,
}

func AliCloudCmsAlarmBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_zones" "default" {
  		available_disk_category     = "cloud_efficiency"
  		available_resource_creation = "VSwitch"
	}

	data "alicloud_images" "default" {
  		most_recent = true
  		owners      = "system"
	}

	data "alicloud_instance_types" "default" {
  		availability_zone = data.alicloud_zones.default.zones.0.id
  		image_id          = data.alicloud_images.default.images.0.id
	}

	resource "alicloud_cms_alarm_contact_group" "default" {
  		alarm_contact_group_name = var.name
  		describe                 = "tf-testacc"
  		contacts                 = ["test1", "test2", "test3"]
	}

	resource "alicloud_cms_monitor_group" "default" {
  		monitor_group_name = var.name
  		contact_groups     = [alicloud_cms_alarm_contact_group.default.id]
	}

	resource "alicloud_cms_alarm_contact_group" "update" {
  		alarm_contact_group_name = "${var.name}-update"
  		describe                 = "tf-testacc"
  		contacts                 = ["test5", "test6", "test7"]
	}

	resource "alicloud_cms_monitor_group" "update" {
  		monitor_group_name = "${var.name}-update"
  		contact_groups     = [alicloud_cms_alarm_contact_group.update.id]
	}

	resource "alicloud_vpc" "default" {
  		vpc_name   = var.name
  		cidr_block = "192.168.0.0/16"
	}

	resource "alicloud_vswitch" "default" {
  		vswitch_name = var.name
  		vpc_id       = alicloud_vpc.default.id
  		cidr_block   = "192.168.192.0/24"
  		zone_id      = data.alicloud_zones.default.zones.0.id
	}

	resource "alicloud_security_group" "default" {
  		name   = var.name
  		vpc_id = alicloud_vpc.default.id
	}

	resource "alicloud_instance" "default" {
  		image_id                   = data.alicloud_images.default.images.0.id
  		instance_type              = data.alicloud_instance_types.default.instance_types.0.id
  		security_groups            = alicloud_security_group.default.*.id
  		internet_charge_type       = "PayByTraffic"
  		internet_max_bandwidth_out = "10"
  		availability_zone          = data.alicloud_instance_types.default.instance_types.0.availability_zones.0
  		instance_charge_type       = "PostPaid"
  		system_disk_category       = "cloud_efficiency"
  		vswitch_id                 = alicloud_vswitch.default.id
  		instance_name              = var.name
	}
`, name)
}
