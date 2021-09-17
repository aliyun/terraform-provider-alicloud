package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudARMSPrometheusAlertRule_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_arms_prometheus_alert_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudARMSPrometheusAlertRuleMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ArmsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeArmsPrometheusAlertRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sarmsprometheusalertrule%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudARMSPrometheusAlertRuleBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckPrePaidResources(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"prometheus_alert_rule_name": name,
					"cluster_id":                 "${data.alicloud_cs_managed_kubernetes_clusters.default.clusters.0.id}",
					"expression":                 "node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes * 100 < 10",
					"message":                    "node available memory is less than 10%",
					"duration":                   "1",
					"notify_type":                "ALERT_MANAGER",
					"type":                       name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"prometheus_alert_rule_name": name,
						"cluster_id":                 CHECKSET,
						"expression":                 "node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes * 100 < 10",
						"message":                    "node available memory is less than 10%",
						"duration":                   "1",
						"notify_type":                "ALERT_MANAGER",
						"type":                       name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"prometheus_alert_rule_name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"prometheus_alert_rule_name": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"duration": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"duration": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"labels": []map[string]interface{}{
						{
							"name":  "TF",
							"value": "test1",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"labels.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"annotations": []map[string]interface{}{
						{
							"name":  "TF",
							"value": "test1",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"annotations.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"message": "node available memory is less than 20%",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"message": "node available memory is less than 20%",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"expression": "node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes * 100 < 20",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"expression": "node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes * 100 < 20",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"prometheus_alert_rule_name": name,
					"type":                       name,
					"duration":                   "1",
					"labels": []map[string]interface{}{
						{
							"name":  "TF2",
							"value": "test2",
						},
					},
					"annotations": []map[string]interface{}{
						{
							"name":  "TF2",
							"value": "test2",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"prometheus_alert_rule_name": name,
						"duration":                   "1",
						"type":                       name,
						"labels.#":                   "1",
						"annotations.#":              "1",
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

func TestAccAlicloudARMSPrometheusAlertRule_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_arms_prometheus_alert_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudARMSPrometheusAlertRuleMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ArmsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeArmsPrometheusAlertRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sarmsprometheusalertrule%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudARMSPrometheusAlertRuleBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckPrePaidResources(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"prometheus_alert_rule_name": name,
					"cluster_id":                 "${data.alicloud_cs_managed_kubernetes_clusters.default.clusters.0.id}",
					"expression":                 "node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes * 100 < 10",
					"message":                    "node available memory is less than 10%",
					"duration":                   "1",
					"notify_type":                "DISPATCH_RULE",
					"dispatch_rule_id":           "${alicloud_arms_dispatch_rule.default.id}",
					"type":                       name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"prometheus_alert_rule_name": name,
						"cluster_id":                 CHECKSET,
						"expression":                 "node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes * 100 < 10",
						"message":                    "node available memory is less than 10%",
						"duration":                   "1",
						"notify_type":                "DISPATCH_RULE",
						"dispatch_rule_id":           CHECKSET,
						"type":                       name,
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

var AlicloudARMSPrometheusAlertRuleMap0 = map[string]string{}

func AlicloudARMSPrometheusAlertRuleBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
data "alicloud_cs_managed_kubernetes_clusters" "default" {
  name_regex = "Default"
}
`, name)
}

func AlicloudARMSPrometheusAlertRuleBasicDependence1(name string) string {
	return fmt.Sprintf(`
variable "name" {
 default = "%v"
}
data "alicloud_cs_managed_kubernetes_clusters" "default" {
  name_regex = "Default"
}
resource "alicloud_arms_alert_contact" "default" {
  alert_contact_name = var.name
  email              = "${var.name}@aaa.com"
}
resource "alicloud_arms_alert_contact_group" "default" {
  alert_contact_group_name = var.name
  contact_ids              = [alicloud_arms_alert_contact.default.id]
}

resource "alicloud_arms_dispatch_rule" "default" {
  dispatch_rule_name = var.name
  dispatch_type      = "CREATE_ALERT"
  group_rules {
    group_wait_time = 5
    group_interval  = 15
    repeat_interval = 100
    grouping_fields = [
      "alertname"]
  }
  label_match_expression_grid {
   label_match_expression_groups {
     label_match_expressions {
       key      = "_aliyun_arms_involvedObject_kind"
       value    = "app"
       operator = "eq"
     }
   }
  }

  notify_rules {
    notify_objects {
      notify_object_id = alicloud_arms_alert_contact.default.id
      notify_type      = "ARMS_CONTACT"
      name             = var.name
    }
    notify_objects {
      notify_object_id = alicloud_arms_alert_contact_group.default.id
      notify_type      = "ARMS_CONTACT_GROUP"
      name             = var.name
    }
    notify_channels = ["dingTalk", "wechat"]
  }
}
`, name)
}
