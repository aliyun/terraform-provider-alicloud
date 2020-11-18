package alicloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudCmsGroupMetricRule_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_group_metric_rule.default"
	ra := resourceAttrInit(resourceId, resourceAlicloudCmsGroupMetricRuleMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsGroupMetricRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sCmsGroupMetricRuletf-testacc-rule-name%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAlicloudCmsGroupMetricRuleBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithCmsGroupIdSetting(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"escalations": []map[string]interface{}{
						{
							"critical": []map[string]interface{}{
								{
									"comparison_operator": "GreaterThanOrEqualToThreshold",
									"statistics":          "Average",
									"threshold":           "90",
									"times":               "3",
								},
							},
						},
					},
					"group_id":               os.Getenv("CMS_GROUP_ID"),
					"group_metric_rule_name": "${var.name}",
					"category":               "ecs",
					"metric_name":            "cpu_total",
					"namespace":              "acs_ecs_dashboard",
					"rule_id":                "4a9a8978-a9cc-55ca-aa7c-530ccd91ae57",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"escalations.#":          "1",
						"group_id":               CHECKSET,
						"group_metric_rule_name": name,
						"category":               "ecs",
						"metric_name":            "cpu_total",
						"namespace":              "acs_ecs_dashboard",
						"rule_id":                "4a9a8978-a9cc-55ca-aa7c-530ccd91ae57",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"category", "interval"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"effective_interval": "00:00-22:59",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"effective_interval": "00:00-22:59",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"email_subject": "tf-testacc-rule-name-warning-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"email_subject": "tf-testacc-rule-name-warning-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"group_metric_rule_name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_metric_rule_name": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"metric_name": "cpu_idle",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"metric_name": "cpu_idle",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"no_effective_interval": "00:00-06:30",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"no_effective_interval": "00:00-06:30",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"period": "240",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"period": "240",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"silence_time": "85800",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"silence_time": "85800",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"webhook": "http://www.aliyun1.com",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"webhook": "http://www.aliyun1.com",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"effective_interval":     "00:00-23:59",
					"email_subject":          "tf-testacc-rule-name-warning",
					"group_id":               os.Getenv("CMS_GROUP_ID"),
					"group_metric_rule_name": "${var.name}",
					"category":               "ecs",
					"metric_name":            "cpu_total",
					"namespace":              "acs_ecs_dashboard",
					"no_effective_interval":  "00:00-05:30",
					"period":                 "60",
					"silence_time":           "86400",
					"webhook":                "http://www.aliyun.com",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"effective_interval":     "00:00-23:59",
						"email_subject":          "tf-testacc-rule-name-warning",
						"group_id":               CHECKSET,
						"group_metric_rule_name": name,
						"category":               "ecs",
						"metric_name":            "cpu_total",
						"namespace":              "acs_ecs_dashboard",
						"no_effective_interval":  "00:00-05:30",
						"period":                 "60",
						"silence_time":           "86400",
						"webhook":                "http://www.aliyun.com",
					}),
				),
			},
		},
	})
}

var resourceAlicloudCmsGroupMetricRuleMap = map[string]string{
	"contact_groups": CHECKSET,
	"dimensions":     CHECKSET,
	"email_subject":  CHECKSET,
	"period":         "300",
	"silence_time":   "86400",
	"status":         CHECKSET,
}

func resourceAlicloudCmsGroupMetricRuleBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
resource "alicloud_cms_alarm_contact_group" "this" {
  alarm_contact_group_name = "tf-testAccCms"    
  describe = "tf-testacc"   
  contacts = ["zhangsan","lisi","lll"] 
}
`, name)
}
