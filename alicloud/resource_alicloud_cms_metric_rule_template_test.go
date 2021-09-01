package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers(
		"alicloud_cms_metric_rule_template",
		&resource.Sweeper{
			Name: "alicloud_cms_metric_rule_template",
			F:    testSweepCmsMetricRuleTemplate,
		})
}

func testSweepCmsMetricRuleTemplate(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "DescribeMetricRuleTemplateList"
	request := map[string]interface{}{}
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	var response map[string]interface{}
	conn, err := client.NewCmsClient()
	if err != nil {
		log.Printf("[ERROR] %s get an error: %#v", action, err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-01-01"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			log.Printf("[ERROR] %s get an error: %#v", action, err)
			return nil
		}

		resp, err := jsonpath.Get("$.Templates.Template", response)

		if formatInt(response["Total"]) != 0 && err != nil {
			log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.Templates.Template", action, err)
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})

			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["Name"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Cms Metric Rule Template: %s", item["Name"].(string))
				continue
			}

			action := "DeleteMetricRuleTemplate"
			request := map[string]interface{}{
				"TemplateId": item["TemplateId"],
			}
			_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				log.Printf("[ERROR] Failed to delete Cms Metric Rule Template (%s): %s", item["Name"].(string), err)
			}
			log.Printf("[INFO] Delete Cms Metric Rule Template success: %s ", item["Name"].(string))
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}

func TestAccAlicloudCloudMonitorServiceMetricRuleTemplate_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_metric_rule_template.default"
	ra := resourceAttrInit(resourceId, AlicloudCloudMonitorServiceMetricRuleTemplateMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsMetricRuleTemplate")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scloudmonitorservicemetricruletemplate%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudMonitorServiceMetricRuleTemplateBasicDependence0)
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
					"description":               "${var.name}",
					"metric_rule_template_name": "${var.name}",
					"alert_templates": []map[string]interface{}{
						{
							"category":    "ecs",
							"metric_name": "cpu_total",
							"namespace":   "acs_ecs_dashboard",
							"rule_name":   "tf_testAcc_new",
							"escalations": []map[string]interface{}{
								{
									"critical": []map[string]interface{}{
										{
											"comparison_operator": "GreaterThanThreshold",
											"statistics":          "Average",
											"threshold":           "90",
											"times":               "3",
										},
									},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":               name,
						"metric_rule_template_name": name,
						"alert_templates.#":         "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"alert_templates": []map[string]interface{}{
						{
							"category":    "ecs",
							"metric_name": "cpu_total",
							"namespace":   "acs_ecs_dashboard",
							"rule_name":   "tf_testAcc_update",
							"escalations": []map[string]interface{}{
								{
									"critical": []map[string]interface{}{
										{
											"comparison_operator": "GreaterThanThreshold",
											"statistics":          "Average",
											"threshold":           "80",
											"times":               "5",
										},
									},
									"info": []map[string]interface{}{
										{
											"comparison_operator": "GreaterThanThreshold",
											"statistics":          "Average",
											"threshold":           "80",
											"times":               "5",
										},
									},
									"warn": []map[string]interface{}{
										{
											"comparison_operator": "GreaterThanThreshold",
											"statistics":          "Average",
											"threshold":           "80",
											"times":               "5",
										},
									},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alert_templates.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"group_id":          "${local.group_id}",
					"silence_time":      "8640",
					"enable_start_time": "00",
					"enable_end_time":   "23",
					"notify_level":      "4",
					"apply_mode":        "GROUP_INSTANCE_FIRST",
					"webhook":           "https://www.aliyun.com",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"alert_templates": []map[string]interface{}{
						{
							"category":    "ecs",
							"metric_name": "cpu_total",
							"namespace":   "acs_ecs_dashboard",
							"rule_name":   "tf_testAcc_update",
							"escalations": []map[string]interface{}{
								{
									"critical": []map[string]interface{}{
										{
											"comparison_operator": "GreaterThanThreshold",
											"statistics":          "Average",
											"threshold":           "90",
											"times":               "3",
										},
									},
									"info": []map[string]interface{}{
										{
											"comparison_operator": "GreaterThanThreshold",
											"statistics":          "Average",
											"threshold":           "90",
											"times":               "3",
										},
									},
									"warn": []map[string]interface{}{
										{
											"comparison_operator": "GreaterThanThreshold",
											"statistics":          "Average",
											"threshold":           "90",
											"times":               "3",
										},
									},
								},
							},
						},
					},
					"description":               "${var.name}",
					"group_id":                  "${local.group_id}",
					"metric_rule_template_name": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alert_templates.#":         "1",
						"description":               name,
						"group_id":                  CHECKSET,
						"metric_rule_template_name": name,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"apply_mode", "notify_level", "enable_end_time", "silence_time", "enable_start_time", "group_id", "webhook"},
			},
		},
	})
}

var AlicloudCloudMonitorServiceMetricRuleTemplateMap0 = map[string]string{
	"alert_templates.#": CHECKSET,
}

func AlicloudCloudMonitorServiceMetricRuleTemplateBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
  ids = [data.alicloud_vpcs.default.vpcs.0.vswitch_ids.0]
}
data "alicloud_zones" "default" {}
resource "alicloud_vswitch" "default" {
  count = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vswitch_name = var.name
  cidr_block = "172.16.0.0/24"
  vpc_id = data.alicloud_vpcs.default.ids.0
  availability_zone = data.alicloud_zones.default.zones.0.id
  tags 		= {
		Created = "TF"
		For 	= "acceptance test"
  }
}

resource "alicloud_slb_load_balancer" "default" {
  load_balancer_name = var.name
  load_balancer_spec = "slb.s2.small"
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.default.*.id, [""])[0]
}
resource "alicloud_cms_monitor_group" "default" {
monitor_group_name = var.name
}
resource "alicloud_cms_monitor_group_instances" "default" {
  group_id = alicloud_cms_monitor_group.default.id
  instances {
    instance_id = alicloud_slb_load_balancer.default.id
    instance_name = alicloud_slb_load_balancer.default.name
    region_id = "%s"
    category = "slb"
  }
}

locals {
 group_id = alicloud_cms_monitor_group_instances.default.id
}

`, name, defaultRegionToTest)
}
