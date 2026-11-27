package arms_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/arms"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdktf "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"
)

func TestAccAliCloudARMSPrometheusAlertRule_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_arms_prometheus_alert_rule.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudARMSPrometheusAlertRuleMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return arms.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeArmsPrometheusAlertRule")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sarmsprometheusalertrule%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudARMSPrometheusAlertRuleBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckPrePaidResources(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             armsCheckPrometheusAlertRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"prometheus_alert_rule_name": name,
					"cluster_id":                 "${alicloud_arms_prometheus.default.cluster_id}",
					"expression":                 "node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes * 100 < 10",
					"message":                    "node available memory is less than 10%",
					"duration":                   "1",
					"notify_type":                "ALERT_MANAGER",
					"type":                       name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"prometheus_alert_rule_name": name,
						"cluster_id":                 acctest.CHECKSET,
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
					"dispatch_rule_id":           "${alicloud_arms_dispatch_rule.default.id}",
					"notify_type":                "DISPATCH_RULE",
					"prometheus_alert_rule_name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					// the check map accumulates across steps, so the values this
					// step changes must be re-asserted with their new expectations
					testAccCheck(map[string]string{
						"prometheus_alert_rule_name": name + "update",
						"notify_type":                "DISPATCH_RULE",
						"dispatch_rule_id":           acctest.CHECKSET,
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
						// UpdatePrometheusAlertRule overwrites the whole rule on the server side,
						// so the unchanged notify_type and dispatch_rule_id must still be carried
						// with the labels update instead of being cleared.
						"notify_type":      "DISPATCH_RULE",
						"dispatch_rule_id": acctest.CHECKSET,
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
					// This step rewrites the annotations in place without the
					// server-derived "message" entry; the raw response must
					// still carry it.
					testAccCheckArmsPrometheusAlertRuleMessageAnnotation("node available memory is less than 10%"),
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
						"labels.#":                   "1",
						"annotations.#":              "1",
					}),
				),
			},
			{
				// The configuration is unchanged, so planning it again must produce an empty plan.
				Config:             testAccConfig(map[string]interface{}{}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAliCloudARMSPrometheusAlertRule_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_arms_prometheus_alert_rule.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudARMSPrometheusAlertRuleMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return arms.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeArmsPrometheusAlertRule")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sarmsprometheusalertrule%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudARMSPrometheusAlertRuleBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckPrePaidResources(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             armsCheckPrometheusAlertRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"prometheus_alert_rule_name": name,
					"cluster_id":                 "${alicloud_arms_prometheus.default.cluster_id}",
					"expression":                 "node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes * 100 < 10",
					"message":                    "node available memory is less than 10%",
					"duration":                   "1",
					"notify_type":                "DISPATCH_RULE",
					"dispatch_rule_id":           "${alicloud_arms_dispatch_rule.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"prometheus_alert_rule_name": name,
						"cluster_id":                 acctest.CHECKSET,
						"expression":                 "node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes * 100 < 10",
						"message":                    "node available memory is less than 10%",
						"duration":                   "1",
						"notify_type":                "DISPATCH_RULE",
						"dispatch_rule_id":           acctest.CHECKSET,
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

// notify_type is Optional+Computed and the server always materializes a value
// (default ALERT_MANAGER). This case covers the configuration that never sets
// notify_type or dispatch_rule_id: with the conditional carry alone, an update
// triggered by another field omitted NotifyType from the request, the server
// reset it to ALERT_MANAGER, and the plan then reported a permanent removal
// diff (observed as `- notify_type = "ALERT_MANAGER" -> null` in ACC task
// 34812). Marking the attribute Computed makes the state value authoritative
// so Update carries it in full; the assertions here are the empty refresh
// plan after each step and the import round-trip.
func TestAccAliCloudARMSPrometheusAlertRule_unconfiguredNotifyType(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_arms_prometheus_alert_rule.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudARMSPrometheusAlertRuleMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return arms.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeArmsPrometheusAlertRule")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sarmsprometheusalertrule%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudARMSPrometheusAlertRuleBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckPrePaidResources(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             armsCheckPrometheusAlertRuleDestroy,
		Steps: []resource.TestStep{
			{
				// notify_type and dispatch_rule_id are deliberately absent: the
				// server materializes its default and Read fills the state.
				Config: testAccConfig(map[string]interface{}{
					"prometheus_alert_rule_name": name,
					"cluster_id":                 "${alicloud_arms_prometheus.default.cluster_id}",
					"expression":                 "node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes * 100 < 10",
					"message":                    "node available memory is less than 10%",
					"duration":                   "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"prometheus_alert_rule_name": name,
						"cluster_id":                 acctest.CHECKSET,
						"expression":                 "node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes * 100 < 10",
						"message":                    "node available memory is less than 10%",
						"duration":                   "1",
					}),
				),
			},
			{
				Config:             testAccConfig(map[string]interface{}{}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Update triggered by an unrelated field (labels; the rule name
				// and other ForceNew attributes cannot trigger it) while
				// notify_type and dispatch_rule_id remain unconfigured: the
				// request must carry the state's NotifyType instead of omitting
				// it, and the plan after the update must stay empty.
				Config: testAccConfig(map[string]interface{}{
					"labels": []map[string]interface{}{
						{
							"name":  "severity",
							"value": "P3",
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
				Config:             testAccConfig(map[string]interface{}{}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// UpdatePrometheusAlertRule is a full-overwrite API: fields omitted from the
// request are cleared on the server side. This case guards the regression
// where an update triggered by one field wiped the other updatable fields
// (labels, annotations, notify_type, dispatch_rule_id) and made every
// following plan report a new diff.
func TestAccAliCloudARMSPrometheusAlertRule_updateFullOverwrite(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_arms_prometheus_alert_rule.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudARMSPrometheusAlertRuleMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return arms.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeArmsPrometheusAlertRule")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sarmsprometheusalertrule%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudARMSPrometheusAlertRuleBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckPrePaidResources(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             armsCheckPrometheusAlertRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"prometheus_alert_rule_name": name,
					"cluster_id":                 "${alicloud_arms_prometheus.default.cluster_id}",
					"expression":                 "node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes * 100 < 10",
					"message":                    "node available memory is less than 10%",
					"duration":                   "1",
					"notify_type":                "ALERT_MANAGER",
					"labels": []map[string]interface{}{
						{
							"name":  "severity",
							"value": "P3",
						},
						{
							"name":  "alert_type",
							"value": "kubernetes pod memory",
						},
					},
					"annotations": []map[string]interface{}{
						{
							"name":  "description",
							"value": "pod memory usage more than 85%",
						},
						{
							"name":  "summary",
							"value": "pod memory alert",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"labels.#":      "2",
						"annotations.#": "2",
						"notify_type":   "ALERT_MANAGER",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// An update triggered by a labels change must carry the unchanged
				// annotations and notify_type as well.
				Config: testAccConfig(map[string]interface{}{
					"labels": []map[string]interface{}{
						{
							"name":  "severity",
							"value": "P2",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"labels.#":      "1",
						"annotations.#": "2",
						"notify_type":   "ALERT_MANAGER",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:             testAccConfig(map[string]interface{}{}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			{
				// Removing elements must shrink the server-side labels set. The
				// config generator renders a set of objects as HCL blocks and has
				// no zero-block form, so an empty set is expressed by removing the
				// argument entirely.
				Config: testAccConfig(map[string]interface{}{
					"labels": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"labels.#":      "0",
						"annotations.#": "2",
						"notify_type":   "ALERT_MANAGER",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:             testAccConfig(map[string]interface{}{}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			{
				// The same shrink transitions must hold for annotations: dropping
				// one element shrinks the server-side set, and the unchanged
				// (already empty) labels and notify_type must not come back.
				Config: testAccConfig(map[string]interface{}{
					"annotations": []map[string]interface{}{
						{
							"name":  "description",
							"value": "pod memory usage more than 85%",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"labels.#":      "0",
						"annotations.#": "1",
						"notify_type":   "ALERT_MANAGER",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:             testAccConfig(map[string]interface{}{}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			{
				// Removing the annotations argument entirely must clear the
				// server-side set via the "[]" overwrite payload.
				Config: testAccConfig(map[string]interface{}{
					"annotations": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"labels.#":      "0",
						"annotations.#": "0",
						"notify_type":   "ALERT_MANAGER",
					}),
					// The "[]" overwrite must still leave the server-derived
					// "message" annotation in the raw response: the update
					// carries Message, and the annotation is documented to
					// mirror it.
					testAccCheckArmsPrometheusAlertRuleMessageAnnotation("node available memory is less than 10%"),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:             testAccConfig(map[string]interface{}{}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			{
				// Switching notify_type to DISPATCH_RULE must persist the dispatch rule id.
				Config: testAccConfig(map[string]interface{}{
					"notify_type":      "DISPATCH_RULE",
					"dispatch_rule_id": "${alicloud_arms_dispatch_rule.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"notify_type":      "DISPATCH_RULE",
						"dispatch_rule_id": acctest.CHECKSET,
						"annotations.#":    "0",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:             testAccConfig(map[string]interface{}{}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			{
				// Switching back to ALERT_MANAGER removes dispatch_rule_id from the
				// configuration; its diff is suppressed for non-DISPATCH_RULE notify types.
				Config: testAccConfig(map[string]interface{}{
					"notify_type":      "ALERT_MANAGER",
					"dispatch_rule_id": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"notify_type":   "ALERT_MANAGER",
						"annotations.#": "0",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:             testAccConfig(map[string]interface{}{}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

var AlicloudARMSPrometheusAlertRuleMap0 = map[string]string{}

// testAccCheckArmsPrometheusAlertRuleMessageAnnotation verifies against the
// raw Describe response that the server-side "message" annotation (which
// Read deliberately filters out of the state) survives the in-place updates:
// the update request always carries Message and rewrites the annotations
// without a "message" entry, and the ListPrometheusAlertRules response
// documents the message annotation as mirroring the Message field, so it
// must stay present with the current message instead of being silently
// dropped by the full overwrite.
func testAccCheckArmsPrometheusAlertRuleMessageAnnotation(message string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources["alicloud_arms_prometheus_alert_rule.default"]
		if !ok {
			return fmt.Errorf("resource alicloud_arms_prometheus_alert_rule.default is not found in the state")
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource id of alicloud_arms_prometheus_alert_rule.default is empty")
		}
		armsService := arms.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
		object, err := armsService.DescribeArmsPrometheusAlertRule(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("describe alicloud_arms_prometheus_alert_rule %s got an error: %v", rs.Primary.ID, err)
		}
		annotations, ok := object["Annotations"].([]interface{})
		if !ok {
			return fmt.Errorf("the raw response of %s carries no Annotations list: %v", rs.Primary.ID, object["Annotations"])
		}
		for _, val := range annotations {
			item, ok := val.(map[string]interface{})
			if !ok {
				continue
			}
			if item["Name"] == "message" {
				if fmt.Sprint(item["Value"]) != message {
					return fmt.Errorf("the message annotation value %q does not match the rule message %q", fmt.Sprint(item["Value"]), message)
				}
				return nil
			}
		}
		return fmt.Errorf("the server annotations %v lost the message annotation derived from the rule message %q", annotations, message)
	}
}

// armsCheckPrometheusAlertRuleDestroy adapts the destroy check to the
// absence contract pinned by TestUnitArmsPrometheusAlertRuleParameterIllegal.
// By the time it runs, the post-test destroy has already removed the alert
// rule (terraform deletes it before its parent prometheus cluster) and the
// cluster itself, so ListPrometheusAlertRules fails with ParameterIllegal
// for the deleted cluster while the rule-level rescue inside
// DescribeArmsPrometheusAlertRule fails to classify the deleted rule — the
// error the contract now propagates instead of a synthetic not-found. A live
// rule still fails the check: the rescue returns it, so only the ambiguous
// ParameterIllegal wrap is tolerated here, and a rule that the list reports
// missing from a live cluster keeps failing with the not-found error.
func armsCheckPrometheusAlertRuleDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_arms_prometheus_alert_rule" {
			continue
		}
		armsService := arms.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
		_, err := armsService.DescribeArmsPrometheusAlertRule(rs.Primary.ID)
		if err != nil {
			if errs.NotFoundError(err) {
				continue
			}
			if !errs.IsExpectedErrors(err, []string{"ParameterIllegal"}) {
				return errs.WrapError(err)
			}
			continue
		}
		return fmt.Errorf("the alert rule %s %s was not destroyed", rs.Type, rs.Primary.ID)
	}
	return nil
}

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
resource "alicloud_arms_prometheus" "default" {
  cluster_type        = "remote-write"
  cluster_name        = var.name
  grafana_instance_id = "free"
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
  # ARMS adds a _aliyun_arms_alert_rule_id label matcher to the rule whenever
  # an alert rule binds it, which the dispatch rule schema does not model;
  # the fixture ignores that server-side mutation so the refresh plan of the
  # alert rule steps stays empty.
  lifecycle {
    ignore_changes = [label_match_expression_grid]
  }
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
    notify_channels   = ["dingTalk", "wechat"]
    notify_start_time = "00:00"
    notify_end_time   = "23:59"
  }
}
`, name)
}

// lintignore: R001
func TestUnitAlicloudARMSPrometheusAlertRule(t *testing.T) {
	p := unitTestProviderResources()
	rand := tfacctest.RandIntRange(10000, 99999)
	d, _ := schema.InternalMap(p["alicloud_arms_prometheus_alert_rule"].Schema).Data(nil, nil)
	dCreate, _ := schema.InternalMap(p["alicloud_arms_prometheus_alert_rule"].Schema).Data(nil, nil)
	dCreate.MarkNewResource()
	for key, value := range map[string]interface{}{
		"prometheus_alert_rule_name": "prometheus_alert_rule_name",
		"cluster_id":                 "ClusterId",
		"expression":                 "node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes * 100 < 10",
		"message":                    "node available memory is less than 10%",
		"duration":                   "1",
		"notify_type":                "DISPATCH_RULE",
		"type":                       "type",
		"dispatch_rule_id":           "1",
		"annotations": []map[string]interface{}{
			{
				"name":  "TF",
				"value": "test1",
			},
		},
		"labels": []map[string]interface{}{
			{
				"name":  "TF2",
				"value": "test2",
			},
		},
	} {
		err := dCreate.Set(key, value)
		assert.Nil(t, err)
		err = d.Set(key, value)
		assert.Nil(t, err)
	}
	region := os.Getenv("ALICLOUD_REGION")
	rawClient, err := armsTestClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}
	rawClient = rawClient.(*connectivity.AliyunClient)
	ReadMockResponse := map[string]interface{}{
		"PrometheusAlertRules": []interface{}{
			map[string]interface{}{
				// AlertId must be a numeric string: Read converts the id part
				// of the resource id with convert.FormatInt, and a successful
				// Create/Update ends with Read, so a non-numeric mock id would
				// fail those subcases at the added error propagation.
				"AlertId":               "123456",
				"ClusterId":             "ClusterId",
				"PrometheusAlertRuleId": "123456",
				"Duration":              "duration",
				"Expression":            "expression",
				"Annotations": []interface{}{
					map[string]interface{}{
						"Name":  "TF",
						"Value": "test1",
					},
				},
				"Labels": []interface{}{
					map[string]interface{}{
						"Name":  "TF",
						"Value": "test1",
					},
				},
				"DispatchRuleId": "1",
				"Message":        "message",
				"NotifyType":     "notify_type",
				"AlertName":      "prometheus_alert_rule_name",
				"Status":         "1",
				"Type":           "type",
			},
		},
	}

	responseMock := map[string]func(errorCode string) (map[string]interface{}, error){
		"RetryError": func(errorCode string) (map[string]interface{}, error) {
			return nil, &tea.SDKError{
				Code:       helper.String(errorCode),
				Data:       helper.String(errorCode),
				Message:    helper.String(errorCode),
				StatusCode: tea.Int(400),
			}
		},
		"NotFoundError": func(errorCode string) (map[string]interface{}, error) {
			return nil, errs.GetNotFoundErrorFromString(errs.GetNotFoundMessage("alicloud_arms_prometheus_alert_rule", "123456"))
		},
		"NoRetryError": func(errorCode string) (map[string]interface{}, error) {
			return nil, &tea.SDKError{
				Code:       helper.String(errorCode),
				Data:       helper.String(errorCode),
				Message:    helper.String(errorCode),
				StatusCode: tea.Int(400),
			}
		},
		"CreateNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			result["PrometheusAlertRule"] = map[string]interface{}{
				"AlertId": "123456",
			}
			return result, nil
		},
		"UpdateNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			return result, nil
		},
		"DeleteNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			return result, nil
		},
		"ReadNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			return result, nil
		},
		"AnnotationsName": func(errorCode string) (map[string]interface{}, error) {
			result := map[string]interface{}{
				"PrometheusAlertRules": []interface{}{
					map[string]interface{}{
						"AlertId":               "123456",
						"ClusterId":             "ClusterId",
						"PrometheusAlertRuleId": "123456",
						"Annotations": []interface{}{
							map[string]interface{}{
								"Name":  "message",
								"Value": "test1",
							},
						},
					},
				},
			}
			result["PrometheusAlertRule"] = map[string]interface{}{
				"AlertId": "123456",
			}
			return result, nil
		},
	}
	// Create
	t.Run("CreateClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewArmsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       helper.String("loadEndpoint error"),
				Data:       helper.String("loadEndpoint error"),
				Message:    helper.String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})
		err := arms.ResourcePrometheusAlertRuleCreate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("CreateAbnormal", func(t *testing.T) {
		retryFlag := true
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["CreateNormal"]("")
		})
		err := arms.ResourcePrometheusAlertRuleCreate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("CreateNormal", func(t *testing.T) {
		retryFlag := false
		noRetryFlag := false
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["CreateNormal"]("")
		})
		err := arms.ResourcePrometheusAlertRuleCreate(dCreate, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	t.Run("CreateAnnotationsNameMock", func(t *testing.T) {
		retryFlag := false
		noRetryFlag := false
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["AnnotationsName"]("")
		})
		err := arms.ResourcePrometheusAlertRuleCreate(dCreate, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	// Set ID for Update and Delete Method
	d.SetId(fmt.Sprint("ClusterId", ":", "123456"))
	// Update
	t.Run("UpdateClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewArmsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       helper.String("loadEndpoint error"),
				Data:       helper.String("loadEndpoint error"),
				Message:    helper.String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})

		err := arms.ResourcePrometheusAlertRuleUpdate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	t.Run("UpdatePrometheusAlertRuleAbnormal", func(t *testing.T) {
		diff := sdktf.NewInstanceDiff()
		for _, key := range []string{"duration", "expression", "message", "prometheus_alert_rule_name", "annotations", "dispatch_rule_id", "labels", "notify_type"} {
			switch p["alicloud_arms_prometheus_alert_rule"].Schema[key].Type {
			case schema.TypeString:
				diff.Attributes[key] = &sdktf.ResourceAttrDiff{Old: d.Get(key).(string), New: d.Get(key).(string) + "_update"}
			case schema.TypeBool:
				diff.Attributes[key] = &sdktf.ResourceAttrDiff{Old: strconv.FormatBool(d.Get(key).(bool)), New: strconv.FormatBool(true)}
			case schema.TypeInt:
				diff.Attributes[key] = &sdktf.ResourceAttrDiff{Old: strconv.Itoa(d.Get(key).(int)), New: strconv.Itoa(2)}
			case schema.TypeMap:
				diff.Attributes["tags.%"] = &sdktf.ResourceAttrDiff{Old: "0", New: "2"}
				diff.Attributes["tags.For"] = &sdktf.ResourceAttrDiff{Old: "", New: "Test"}
				diff.Attributes["tags.Created"] = &sdktf.ResourceAttrDiff{Old: "", New: "TF"}
			case schema.TypeSet:
				diff.Attributes[fmt.Sprintf("%s.#", key)] = &sdktf.ResourceAttrDiff{Old: "1", New: "1"}
				for _, ipConfig := range d.Get(key).(*schema.Set).List() {
					ipConfigArg := ipConfig.(map[string]interface{})
					for field := range p["alicloud_arms_prometheus_alert_rule"].Schema[key].Elem.(*schema.Resource).Schema {
						diff.Attributes[fmt.Sprintf("%s.%d.%s", key, rand, field)] = &sdktf.ResourceAttrDiff{Old: ipConfigArg[field].(string), New: ipConfigArg[field].(string) + "_update"}
					}
				}
			}
		}
		resourceData1, _ := schema.InternalMap(p["alicloud_arms_prometheus_alert_rule"].Schema).Data(nil, diff)
		resourceData1.SetId(d.Id())
		retryFlag := true
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["UpdateNormal"]("")
		})
		err := arms.ResourcePrometheusAlertRuleUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	t.Run("UpdatePrometheusAlertRuleNormal", func(t *testing.T) {
		diff := sdktf.NewInstanceDiff()
		for _, key := range []string{"duration", "expression", "message", "prometheus_alert_rule_name", "annotations", "dispatch_rule_id", "labels", "notify_type"} {
			switch p["alicloud_arms_prometheus_alert_rule"].Schema[key].Type {
			case schema.TypeString:
				diff.Attributes[key] = &sdktf.ResourceAttrDiff{Old: d.Get(key).(string), New: d.Get(key).(string) + "_update"}
			case schema.TypeBool:
				diff.Attributes[key] = &sdktf.ResourceAttrDiff{Old: strconv.FormatBool(d.Get(key).(bool)), New: strconv.FormatBool(true)}
			case schema.TypeInt:
				diff.Attributes[key] = &sdktf.ResourceAttrDiff{Old: strconv.Itoa(d.Get(key).(int)), New: strconv.Itoa(3)}
			case schema.TypeMap:
				diff.Attributes["tags.%"] = &sdktf.ResourceAttrDiff{Old: "0", New: "2"}
				diff.Attributes["tags.For"] = &sdktf.ResourceAttrDiff{Old: "", New: "Test"}
				diff.Attributes["tags.Created"] = &sdktf.ResourceAttrDiff{Old: "", New: "TF"}
			case schema.TypeSet:
				diff.Attributes[fmt.Sprintf("%s.#", key)] = &sdktf.ResourceAttrDiff{Old: "1", New: "1"}
				for _, ipConfig := range d.Get(key).(*schema.Set).List() {
					ipConfigArg := ipConfig.(map[string]interface{})
					for field := range p["alicloud_arms_prometheus_alert_rule"].Schema[key].Elem.(*schema.Resource).Schema {
						diff.Attributes[fmt.Sprintf("%s.%d.%s", key, rand, field)] = &sdktf.ResourceAttrDiff{Old: ipConfigArg[field].(string), New: ipConfigArg[field].(string) + "_update"}
					}
				}
			}
		}
		resourceData1, _ := schema.InternalMap(p["alicloud_arms_prometheus_alert_rule"].Schema).Data(nil, diff)
		resourceData1.SetId(d.Id())
		retryFlag := false
		noRetryFlag := false
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["UpdateNormal"]("")
		})
		err := arms.ResourcePrometheusAlertRuleUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	t.Run("UpdateMockParseResourceId", func(t *testing.T) {
		resourceData1, _ := schema.InternalMap(p["alicloud_arms_prometheus_alert_rule"].Schema).Data(nil, nil)
		resourceData1.SetId("MockId")
		retryFlag := false
		noRetryFlag := false
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("RetryError")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["UpdateNormal"]("")
		})
		err := arms.ResourcePrometheusAlertRuleUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	// Delete
	t.Run("DeleteClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewArmsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       helper.String("loadEndpoint error"),
				Data:       helper.String("loadEndpoint error"),
				Message:    helper.String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})
		err := arms.ResourcePrometheusAlertRuleDelete(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("DeleteMockAbnormal", func(t *testing.T) {
		retryFlag := true
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["DeleteNormal"]("")
		})
		err := arms.ResourcePrometheusAlertRuleDelete(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("DeleteMockNormal", func(t *testing.T) {
		retryFlag := false
		noRetryFlag := false
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["DeleteNormal"]("")
		})
		err := arms.ResourcePrometheusAlertRuleDelete(d, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	t.Run("DeleteMockParseResourceId", func(t *testing.T) {
		resourceData1, _ := schema.InternalMap(p["alicloud_arms_prometheus_alert_rule"].Schema).Data(nil, nil)
		resourceData1.SetId("MockId")
		retryFlag := false
		noRetryFlag := false
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("RetryError")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["DeleteNormal"]("")
		})
		err := arms.ResourcePrometheusAlertRuleDelete(resourceData1, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	//Read
	t.Run("ReadDescribeArmsPrometheusAlertRuleNotFound", func(t *testing.T) {
		patcheDorequest := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			NotFoundFlag := true
			noRetryFlag := false
			if NotFoundFlag {
				return responseMock["NotFoundError"]("ResourceNotfound")
			} else if noRetryFlag {
				return responseMock["NoRetryError"]("NoRetryError")
			}
			return responseMock["ReadNormal"]("")
		})
		err := arms.ResourcePrometheusAlertRuleRead(d, rawClient)
		patcheDorequest.Reset()
		assert.Nil(t, err)
	})

	t.Run("ReadDescribeArmsPrometheusAlertRuleAbnormal", func(t *testing.T) {
		patcheDorequest := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			retryFlag := false
			noRetryFlag := true
			if retryFlag {
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["ReadNormal"]("")
		})
		err := arms.ResourcePrometheusAlertRuleRead(d, rawClient)
		patcheDorequest.Reset()
		assert.NotNil(t, err)
	})

	t.Run("ReadMockParseResourceId", func(t *testing.T) {
		resourceData1, _ := schema.InternalMap(p["alicloud_arms_prometheus_alert_rule"].Schema).Data(nil, nil)
		resourceData1.SetId("MockId")
		patcheDorequest := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			retryFlag := false
			noRetryFlag := false
			if retryFlag {
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["ReadNormal"]("")
		})
		patcheDescribeArmsPrometheusAlertRule := gomonkey.ApplyMethod(reflect.TypeOf(&arms.Service{}), "DescribeArmsPrometheusAlertRule", func(*arms.Service, string) (map[string]interface{}, error) {
			return responseMock["ReadNormal"]("")
		})
		err := arms.ResourcePrometheusAlertRuleRead(resourceData1, rawClient)
		patcheDorequest.Reset()
		patcheDescribeArmsPrometheusAlertRule.Reset()
		assert.NotNil(t, err)
	})
}

// TestUnitArmsPrometheusAlertRuleParameterIllegal pins the absence
// classification behind the ParameterIllegal handling: live reproduction
// (https://jarivs.alibaba.net/Jarvis/tasks/11159f8c-9f21-4a86-a5d6-74297743fab8/)
// proved alert rules survive the deletion of their parent prometheus
// cluster, and a second round
// (https://jarivs.alibaba.net/Jarvis/tasks/dc99c0f7-00d7-438f-88fb-45c45fddcc4b/)
// proved the rule-level lookup cannot classify the rule either:
// DescribePrometheusAlertRule also returns ParameterIllegal for a live
// orphaned rule (while DeletePrometheusAlertRule on the same id returns
// 200), so neither error ever proves the rule is gone. Only a SUCCESSFUL
// rule-level lookup positively confirms the rule and completes the refresh;
// every other outcome keeps the original list error and the state, and the
// id is cleared only when the list itself succeeds without the rule. It
// drives the production Read through a local mock ARMS endpoint, so unlike
// the gomonkey subtests above it also runs in CI without cloud credentials.
func TestUnitArmsPrometheusAlertRuleParameterIllegal(t *testing.T) {
	parameterIllegal := func(w http.ResponseWriter) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"Code":    "ParameterIllegal",
			"Message": "the request parameters are illegal or the cluster does not exist",
		})
	}
	ruleAlive := func(w http.ResponseWriter) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"PrometheusAlertRule": map[string]interface{}{
				"AlertId":        "123456",
				"ClusterId":      "ClusterId",
				"Duration":       "1",
				"Expression":     "expression",
				"Message":        "message",
				"NotifyType":     "notify_type",
				"AlertName":      "prometheus_alert_rule_name",
				"Status":         "1",
				"Type":           "type",
				"DispatchRuleId": "0",
			},
		})
	}
	ruleLookupDenied := func(w http.ResponseWriter) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"Code": "APINotAuthorize", "Message": "denied"})
	}

	// listFailsWithParameterIllegal mocks ListPrometheusAlertRules failing
	// with ParameterIllegal while the DescribePrometheusAlertRule lookup
	// behaves as ruleLookup describes.
	listFailsWithParameterIllegal := func(t *testing.T, ruleLookup func(w http.ResponseWriter)) *connectivity.AliyunClient {
		return armsNotificationPolicyTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Form.Get("Action") == "DescribePrometheusAlertRule" {
				ruleLookup(w)
				return
			}
			parameterIllegal(w)
		})
	}

	// listSucceedsWithoutRule mocks a live, responding cluster whose
	// ListPrometheusAlertRules succeeds but does not contain the rule — the
	// only classification that legitimately clears the id.
	listSucceedsWithoutRule := func(t *testing.T) *connectivity.AliyunClient {
		return armsNotificationPolicyTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Form.Get("Action") == "ListPrometheusAlertRules" {
				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"PrometheusAlertRules": []interface{}{},
				})
				return
			}
			parameterIllegal(w)
		})
	}

	ruleId := "ClusterId:123456"

	t.Run("cluster_gone_rule_alive_keeps_id", func(t *testing.T) {
		// The reported defect: the parent prometheus cluster is deleted,
		// the list fails with ParameterIllegal, but the rule-level lookup
		// succeeds — the rule is alive and must stay in the state.
		client := listFailsWithParameterIllegal(t, ruleAlive)
		d := schema.TestResourceDataRaw(t, arms.ResourcePrometheusAlertRule().Schema, nil)
		d.SetId(ruleId)
		if err := arms.ResourcePrometheusAlertRuleRead(d, client); err != nil {
			t.Fatalf("a live rule must survive a deleted parent cluster: %s", err)
		}
		if d.Id() != ruleId {
			t.Fatalf("read dropped a live rule whose parent cluster is gone: %q", d.Id())
		}
	})

	t.Run("orphaned_rule_lookup_ambiguous_keeps_id", func(t *testing.T) {
		// The second reported defect: the list fails with ParameterIllegal
		// AND the rule-level lookup also returns ParameterIllegal — the
		// exact response a live orphaned rule produces (round-2/3
		// reproduction records). The lookup is inconclusive, so the read
		// must fail with the original list error and keep the id instead
		// of clearing it.
		client := listFailsWithParameterIllegal(t, parameterIllegal)
		d := schema.TestResourceDataRaw(t, arms.ResourcePrometheusAlertRule().Schema, nil)
		d.SetId(ruleId)
		err := arms.ResourcePrometheusAlertRuleRead(d, client)
		if err == nil {
			t.Fatal("an inconclusive rule lookup must fail the read")
		}
		if !strings.Contains(err.Error(), "ParameterIllegal") {
			t.Fatalf("the original ParameterIllegal must be reported, got: %s", err)
		}
		if d.Id() != ruleId {
			t.Fatalf("read dropped a rule whose existence is unknown: %q", d.Id())
		}
	})

	t.Run("rule_lookup_denied_keeps_id", func(t *testing.T) {
		// The rule-level lookup fails with an unrelated error (403): the
		// original list error is reported and the id is kept.
		client := listFailsWithParameterIllegal(t, ruleLookupDenied)
		d := schema.TestResourceDataRaw(t, arms.ResourcePrometheusAlertRule().Schema, nil)
		d.SetId(ruleId)
		err := arms.ResourcePrometheusAlertRuleRead(d, client)
		if err == nil {
			t.Fatal("an inconclusive rule lookup must fail the read")
		}
		if !strings.Contains(err.Error(), "ParameterIllegal") {
			t.Fatalf("the original ParameterIllegal must be reported, got: %s", err)
		}
		if d.Id() != ruleId {
			t.Fatalf("read dropped a rule whose existence is unknown: %q", d.Id())
		}
	})

	t.Run("rule_absent_in_live_cluster_clears_id", func(t *testing.T) {
		// The only legitimate absence classification: the cluster is
		// alive, the list succeeds, and the rule is not in it.
		client := listSucceedsWithoutRule(t)
		d := schema.TestResourceDataRaw(t, arms.ResourcePrometheusAlertRule().Schema, nil)
		d.SetId(ruleId)
		if err := arms.ResourcePrometheusAlertRuleRead(d, client); err != nil {
			t.Fatalf("a rule missing from a live cluster's list must be treated as not found: %s", err)
		}
		if d.Id() != "" {
			t.Fatalf("read kept the id of a rule absent from a live cluster list: %q", d.Id())
		}
	})
}
