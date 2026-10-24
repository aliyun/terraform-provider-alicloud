package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	credentials "github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	resource.AddTestSweepers("alicloud_arms_notification_policy", &resource.Sweeper{
		Name: "alicloud_arms_notification_policy",
		F:    testSweepArmsNotificationPolicy,
	})
}

func testSweepArmsNotificationPolicy(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testacc",
	}

	request := make(map[string]interface{})
	query := map[string]interface{}{
		"RegionId": client.RegionId,
		"Page":     1,
		"Size":     PageSizeLarge,
	}
	var response map[string]interface{}
	for {
		action := "ListNotificationPolicies"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(2*time.Minute, func() *retry.RetryError {
			response, err = client.RpcGet("ARMS", "2019-08-08", action, query, request)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, query)
		if err != nil {
			log.Printf("[ERROR] %s got an error: %s", action, err)
			return nil
		}
		resp, err := jsonpath.Get("$.PageBean.NotificationPolicies", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.PageBean.NotificationPolicies", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item, ok := v.(map[string]interface{})
			if !ok {
				continue
			}
			name := fmt.Sprint(item["Name"])
			skip := true
			if !sweepAll() {
				for _, prefix := range prefixes {
					if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
						skip = false
						break
					}
				}
				if skip {
					log.Printf("[INFO] Skipping arms notification policy: %s ", name)
					continue
				}
			}
			log.Printf("[INFO] delete arms notification policy: %s ", name)
			deleteRequest := map[string]interface{}{
				"Id": fmt.Sprint(item["Id"]),
			}
			_, err = client.RpcPost("ARMS", "2019-08-08", "DeleteNotificationPolicy", nil, deleteRequest, false)
			if err != nil {
				log.Printf("[ERROR] Failed to delete arms notification policy (%s): %s", name, err)
			}
		}
		if len(result) < PageSizeLarge {
			break
		}
		query["Page"] = query["Page"].(int) + 1
	}

	return nil
}

// Test Arms NotificationPolicy. >>> Resource test cases, hand-written.
func TestAccAliCloudArmsNotificationPolicy_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_arms_notification_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudArmsNotificationPolicyMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ArmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeArmsNotificationPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sarmsnp%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudArmsNotificationPolicyBasicDependence)
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
					"name":                 name,
					"state":                "enable",
					"send_recover_message": true,
					"escalation_policy_id": 0,
					"repeat":               true,
					"repeat_interval":      600,
					"integration_id":       0,
					"directed_mode":        false,
					"notify_rule": []map[string]interface{}{
						{
							"notify_start_time": "00:00",
							"notify_end_time":   "23:59",
							"notify_channels":   []string{"dingTalk", "email", "sms", "tts", "webhook"},
							"notify_objects": []map[string]interface{}{
								{
									"notify_object_type": "ARMS_CONTACT",
									"notify_object_id":   "${alicloud_arms_alert_contact.default[0].id}",
									"notify_object_name": "${alicloud_arms_alert_contact.default[0].alert_contact_name}",
									"notify_channels":    []string{"email", "sms", "tts"},
								},
							},
						},
					},
					"matching_rules": []map[string]interface{}{
						{
							"matching_conditions": []map[string]interface{}{
								{
									"key":      "alertname",
									"value":    "tf-test-alert",
									"operator": "eq",
								},
							},
						},
					},
					"group_rule": []map[string]interface{}{
						{
							"grouping_fields": []string{"alertname"},
							"group_wait":      5,
							"group_interval":  30,
						},
					},
					"notify_template": []map[string]interface{}{
						{
							"email_title":           "Alert: {{ .commonLabels.alertname }}",
							"email_content":         "{{ .commonLabels.alertname }} fired",
							"email_recover_title":   "Recovered: {{ .commonLabels.alertname }}",
							"email_recover_content": "{{ .commonLabels.alertname }} recovered",
							"sms_content":           "SMS alert: {{ .commonLabels.alertname }}",
							"sms_recover_content":   "SMS recovered: {{ .commonLabels.alertname }}",
							"tts_content":           "TTS alert: {{ .commonLabels.alertname }}",
							"tts_recover_content":   "TTS recovered: {{ .commonLabels.alertname }}",
							"robot_content":         "Robot alert: {{ .commonLabels.alertname }}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                   name,
						"state":                  "enable",
						"notification_policy_id": CHECKSET,
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

// armsNotificationPolicyReorderChecks returns the state-check helper and destroy
// check shared by the TypeList Order Coverage tests below. TestCase field values
// (PreCheck/Providers/CheckDestroy) may be supplied by a helper; the coverage
// checker only parses the inline resource.TestCase{Steps:...} literal and the
// literal config maps passed to the resourceTestAccConfigFunc builder.
func armsNotificationPolicyReorderChecks(t *testing.T) (resourceAttrMapUpdate, resource.TestCheckFunc) {
	var v map[string]interface{}
	resourceId := "alicloud_arms_notification_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudArmsNotificationPolicyMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ArmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeArmsNotificationPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	return rac.resourceAttrMapUpdateSet(), rac.checkResourceDestroy()
}

func TestAccAliCloudArmsNotificationPolicy_matchingRulesOrder(t *testing.T) {
	testAccCheck, checkDestroy := armsNotificationPolicyReorderChecks(t)
	resourceId := "alicloud_arms_notification_policy.default"
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sarmsnpmr%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudArmsNotificationPolicyBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":            name,
					"repeat_interval": 600,
					"notify_rule": []map[string]interface{}{
						{
							"notify_start_time": "00:00",
							"notify_end_time":   "23:59",
							"notify_channels":   []string{"email"},
							"notify_objects": []map[string]interface{}{
								{"notify_object_type": "ARMS_CONTACT", "notify_object_id": "${alicloud_arms_alert_contact.default[0].id}", "notify_object_name": "${alicloud_arms_alert_contact.default[0].alert_contact_name}"},
							},
						},
					},
					"matching_rules": []map[string]interface{}{
						{"matching_conditions": []map[string]interface{}{
							{"key": "alertname", "value": "tf-test-alert-a", "operator": "eq"},
						}},
						{"matching_conditions": []map[string]interface{}{
							{"key": "clustername", "value": "tf-test-cluster-b", "operator": "neq"},
						}},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"name": name, "repeat_interval": "600", "escalation_policy_id": "0", "integration_id": "0", "directed_mode": "false"}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"matching_rules": []map[string]interface{}{
						{"matching_conditions": []map[string]interface{}{
							{"key": "clustername", "value": "tf-test-cluster-b", "operator": "neq"},
						}},
						{"matching_conditions": []map[string]interface{}{
							{"key": "alertname", "value": "tf-test-alert-a", "operator": "eq"},
						}},
					},
				}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"matching_rules": []map[string]interface{}{
						{"matching_conditions": []map[string]interface{}{
							{"key": "clustername", "value": "tf-test-cluster-b", "operator": "neq"},
						}},
						{"matching_conditions": []map[string]interface{}{
							{"key": "alertname", "value": "tf-test-alert-a", "operator": "eq"},
						}},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"name": name, "repeat_interval": "600", "escalation_policy_id": "0", "integration_id": "0", "directed_mode": "false"}),
				),
			},
		},
	})
}

func TestAccAliCloudArmsNotificationPolicy_matchingConditionsOrder(t *testing.T) {
	testAccCheck, checkDestroy := armsNotificationPolicyReorderChecks(t)
	resourceId := "alicloud_arms_notification_policy.default"
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sarmsnpc%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudArmsNotificationPolicyBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":            name,
					"repeat_interval": 600,
					"notify_rule": []map[string]interface{}{
						{
							"notify_start_time": "00:00",
							"notify_end_time":   "23:59",
							"notify_channels":   []string{"email"},
							"notify_objects": []map[string]interface{}{
								{"notify_object_type": "ARMS_CONTACT", "notify_object_id": "${alicloud_arms_alert_contact.default[0].id}", "notify_object_name": "${alicloud_arms_alert_contact.default[0].alert_contact_name}"},
							},
						},
					},
					"matching_rules": []map[string]interface{}{
						{"matching_conditions": []map[string]interface{}{
							{"key": "alertname", "value": "tf-test-alert-a", "operator": "eq"},
							{"key": "clustername", "value": "tf-test-cluster-b", "operator": "neq"},
						}},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"name": name, "repeat_interval": "600", "escalation_policy_id": "0", "integration_id": "0", "directed_mode": "false"}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"matching_rules": []map[string]interface{}{
						{"matching_conditions": []map[string]interface{}{
							{"key": "clustername", "value": "tf-test-cluster-b", "operator": "neq"},
							{"key": "alertname", "value": "tf-test-alert-a", "operator": "eq"},
						}},
					},
				}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"matching_rules": []map[string]interface{}{
						{"matching_conditions": []map[string]interface{}{
							{"key": "clustername", "value": "tf-test-cluster-b", "operator": "neq"},
							{"key": "alertname", "value": "tf-test-alert-a", "operator": "eq"},
						}},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"name": name, "repeat_interval": "600", "escalation_policy_id": "0", "integration_id": "0", "directed_mode": "false"}),
				),
			},
		},
	})
}

func TestAccAliCloudArmsNotificationPolicy_groupingFieldsOrder(t *testing.T) {
	testAccCheck, checkDestroy := armsNotificationPolicyReorderChecks(t)
	resourceId := "alicloud_arms_notification_policy.default"
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sarmsnpgf%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudArmsNotificationPolicyBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":            name,
					"repeat_interval": 600,
					"group_rule":      []map[string]interface{}{{"grouping_fields": []string{"alertname", "clustername"}}},
					"notify_rule": []map[string]interface{}{
						{
							"notify_start_time": "00:00",
							"notify_end_time":   "23:59",
							"notify_channels":   []string{"email"},
							"notify_objects": []map[string]interface{}{
								{"notify_object_type": "ARMS_CONTACT", "notify_object_id": "${alicloud_arms_alert_contact.default[0].id}", "notify_object_name": "${alicloud_arms_alert_contact.default[0].alert_contact_name}"},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"name": name, "repeat_interval": "600", "escalation_policy_id": "0", "integration_id": "0", "directed_mode": "false"}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"group_rule": []map[string]interface{}{
						{"grouping_fields": []string{"clustername", "alertname"}},
					},
				}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"group_rule": []map[string]interface{}{
						{"grouping_fields": []string{"clustername", "alertname"}},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"name": name, "repeat_interval": "600", "escalation_policy_id": "0", "integration_id": "0", "directed_mode": "false"}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"group_rule": []map[string]interface{}{
						{"grouping_fields": []string{}},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"name": name, "repeat_interval": "600", "escalation_policy_id": "0", "integration_id": "0", "directed_mode": "false"}),
					resource.TestCheckResourceAttr(resourceId, "group_rule.0.grouping_fields.#", "0"),
				),
			},
		},
	})
}

func TestAccAliCloudArmsNotificationPolicy_notifyChannelsOrder(t *testing.T) {
	testAccCheck, checkDestroy := armsNotificationPolicyReorderChecks(t)
	resourceId := "alicloud_arms_notification_policy.default"
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sarmsnpnc%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudArmsNotificationPolicyBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":            name,
					"repeat_interval": 600,
					"notify_rule": []map[string]interface{}{
						{
							"notify_start_time": "00:00",
							"notify_end_time":   "23:59",
							"notify_channels":   []string{"dingTalk", "email"},
							"notify_objects": []map[string]interface{}{
								{"notify_object_type": "ARMS_CONTACT", "notify_object_id": "${alicloud_arms_alert_contact.default[0].id}", "notify_object_name": "${alicloud_arms_alert_contact.default[0].alert_contact_name}"},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"name": name, "repeat_interval": "600", "escalation_policy_id": "0", "integration_id": "0", "directed_mode": "false"}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"notify_rule": []map[string]interface{}{
						{
							"notify_start_time": "00:00",
							"notify_end_time":   "23:59",
							"notify_channels":   []string{"email", "dingTalk"},
							"notify_objects": []map[string]interface{}{
								{"notify_object_type": "ARMS_CONTACT", "notify_object_id": "${alicloud_arms_alert_contact.default[0].id}", "notify_object_name": "${alicloud_arms_alert_contact.default[0].alert_contact_name}"},
							},
						},
					},
				}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"notify_rule": []map[string]interface{}{
						{
							"notify_start_time": "00:00",
							"notify_end_time":   "23:59",
							"notify_channels":   []string{"email", "dingTalk"},
							"notify_objects": []map[string]interface{}{
								{"notify_object_type": "ARMS_CONTACT", "notify_object_id": "${alicloud_arms_alert_contact.default[0].id}", "notify_object_name": "${alicloud_arms_alert_contact.default[0].alert_contact_name}"},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"name": name, "repeat_interval": "600", "escalation_policy_id": "0", "integration_id": "0", "directed_mode": "false"}),
				),
			},
		},
	})
}

func TestAccAliCloudArmsNotificationPolicy_notifyObjectsOrder(t *testing.T) {
	testAccCheck, checkDestroy := armsNotificationPolicyReorderChecks(t)
	resourceId := "alicloud_arms_notification_policy.default"
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sarmsnpno%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudArmsNotificationPolicyBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":            name,
					"repeat_interval": 600,
					"notify_rule": []map[string]interface{}{
						{"notify_start_time": "00:00", "notify_end_time": "23:59", "notify_channels": []string{"email"}, "notify_objects": []map[string]interface{}{
							{"notify_object_type": "ARMS_CONTACT", "notify_object_id": "${alicloud_arms_alert_contact.default[0].id}", "notify_object_name": "${alicloud_arms_alert_contact.default[0].alert_contact_name}"},
							{"notify_object_type": "ARMS_CONTACT", "notify_object_id": "${alicloud_arms_alert_contact.default[1].id}", "notify_object_name": "${alicloud_arms_alert_contact.default[1].alert_contact_name}"},
						}},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"name": name, "repeat_interval": "600", "escalation_policy_id": "0", "integration_id": "0", "directed_mode": "false"}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"notify_rule": []map[string]interface{}{
						{"notify_start_time": "00:00", "notify_end_time": "23:59", "notify_channels": []string{"email"}, "notify_objects": []map[string]interface{}{
							{"notify_object_type": "ARMS_CONTACT", "notify_object_id": "${alicloud_arms_alert_contact.default[1].id}", "notify_object_name": "${alicloud_arms_alert_contact.default[1].alert_contact_name}"},
							{"notify_object_type": "ARMS_CONTACT", "notify_object_id": "${alicloud_arms_alert_contact.default[0].id}", "notify_object_name": "${alicloud_arms_alert_contact.default[0].alert_contact_name}"},
						}},
					},
				}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"notify_rule": []map[string]interface{}{
						{"notify_start_time": "00:00", "notify_end_time": "23:59", "notify_channels": []string{"email"}, "notify_objects": []map[string]interface{}{
							{"notify_object_type": "ARMS_CONTACT", "notify_object_id": "${alicloud_arms_alert_contact.default[1].id}", "notify_object_name": "${alicloud_arms_alert_contact.default[1].alert_contact_name}"},
							{"notify_object_type": "ARMS_CONTACT", "notify_object_id": "${alicloud_arms_alert_contact.default[0].id}", "notify_object_name": "${alicloud_arms_alert_contact.default[0].alert_contact_name}"},
						}},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"name": name, "repeat_interval": "600", "escalation_policy_id": "0", "integration_id": "0", "directed_mode": "false"}),
				),
			},
		},
	})
}

func TestAccAliCloudArmsNotificationPolicy_notifyObjectChannelsOrder(t *testing.T) {
	testAccCheck, checkDestroy := armsNotificationPolicyReorderChecks(t)
	resourceId := "alicloud_arms_notification_policy.default"
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sarmsnpnoc%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudArmsNotificationPolicyBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":            name,
					"repeat_interval": 600,
					"notify_rule": []map[string]interface{}{
						{"notify_start_time": "00:00", "notify_end_time": "23:59", "notify_channels": []string{"email"}, "notify_objects": []map[string]interface{}{
							{"notify_object_type": "ARMS_CONTACT", "notify_object_id": "${alicloud_arms_alert_contact.default[0].id}", "notify_object_name": "${alicloud_arms_alert_contact.default[0].alert_contact_name}", "notify_channels": []string{"email", "sms"}},
						}},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"name": name, "repeat_interval": "600", "escalation_policy_id": "0", "integration_id": "0", "directed_mode": "false"}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"notify_rule": []map[string]interface{}{
						{"notify_start_time": "00:00", "notify_end_time": "23:59", "notify_channels": []string{"email"}, "notify_objects": []map[string]interface{}{
							{"notify_object_type": "ARMS_CONTACT", "notify_object_id": "${alicloud_arms_alert_contact.default[0].id}", "notify_object_name": "${alicloud_arms_alert_contact.default[0].alert_contact_name}", "notify_channels": []string{"sms", "email"}},
						}},
					},
				}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"notify_rule": []map[string]interface{}{
						{"notify_start_time": "00:00", "notify_end_time": "23:59", "notify_channels": []string{"email"}, "notify_objects": []map[string]interface{}{
							{"notify_object_type": "ARMS_CONTACT", "notify_object_id": "${alicloud_arms_alert_contact.default[0].id}", "notify_object_name": "${alicloud_arms_alert_contact.default[0].alert_contact_name}", "notify_channels": []string{"sms", "email"}},
						}},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"name": name, "repeat_interval": "600", "escalation_policy_id": "0", "integration_id": "0", "directed_mode": "false"}),
				),
			},
		},
	})
}

var AlicloudArmsNotificationPolicyMap = map[string]string{}

func AlicloudArmsNotificationPolicyBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_arms_alert_contact" "default" {
  count              = 2
  alert_contact_name = "${var.name}-${count.index}"
  email              = "terraform-acc-${count.index}@example.invalid"
}
`, name)
}

func armsNotificationPolicyTestClient(t *testing.T, handler http.HandlerFunc) *connectivity.AliyunClient {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Error(err)
		}
		handler(w, r)
	}))
	t.Cleanup(server.Close)
	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").SetAccessKeyId("test-key").SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatal(err)
	}
	endpoints := new(sync.Map)
	endpoint := strings.TrimPrefix(server.URL, "http://")
	t.Setenv("NO_PROXY", endpoint)
	config := &connectivity.Config{
		AccessKey: "test-key", SecretKey: "test-secret", Credential: credential,
		RegionId: "cn-hangzhou", AccountType: "test", Protocol: "http",
		Endpoints: endpoints, SignVersion: new(sync.Map), SkipRegionValidation: true,
	}
	client, err := config.Client()
	if err != nil {
		t.Fatal(err)
	}
	endpoints.Store("arms", endpoint)
	return client
}

// TestUnitArmsNotificationPolicyCrud drives Create -> Read -> Delete against a
// mock ARMS API and asserts the 4 JSON-typed form params serialize to the
// published wire format (JSON string with camelCase nested keys), and that the
// Read path normalizes both parsed-object and JSON-string response variants.
func TestUnitArmsNotificationPolicyCrud(t *testing.T) {
	t.Run("delete_unconfirmed_keeps_id", func(t *testing.T) {
		for _, response := range []map[string]interface{}{{"IsSuccess": false}, {}, {"IsSuccess": "true"}} {
			client := armsNotificationPolicyTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Form.Get("Action") != "DeleteNotificationPolicy" || r.Form.Get("Id") != "1234" {
					t.Error("unexpected delete request")
				}
				w.Header().Set("Content-Type", "application/json")
				if err := json.NewEncoder(w).Encode(response); err != nil {
					t.Error(err)
				}
			})
			d := schema.TestResourceDataRaw(t, resourceAliCloudArmsNotificationPolicy().Schema, nil)
			d.SetId("1234")
			if err := resourceAliCloudArmsNotificationPolicyDelete(d, client); err == nil {
				t.Fatal("unconfirmed deletion must return an error")
			}
			if d.Id() != "1234" {
				t.Fatalf("failed deletion lost id: %q", d.Id())
			}
		}
	})

	t.Run("create_missing_id", func(t *testing.T) {
		client := armsNotificationPolicyTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Form.Get("Action") != "CreateOrUpdateNotificationPolicy" {
				t.Error("must not read a missing create id")
			}
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(map[string]interface{}{"NotificationPolicy": map[string]interface{}{}}); err != nil {
				t.Error(err)
			}
		})
		d := schema.TestResourceDataRaw(t, resourceAliCloudArmsNotificationPolicy().Schema, map[string]interface{}{
			"name": "test-np-missing-id", "repeat_interval": 600,
			"notify_rule": []interface{}{map[string]interface{}{"notify_channels": []interface{}{"email"}}},
		})
		if err := resourceAliCloudArmsNotificationPolicyCreate(d, client); err == nil {
			t.Fatal("missing create id must return an error")
		}
		if d.Id() != "" {
			t.Fatalf("missing create id polluted state: %q", d.Id())
		}
	})

	t.Run("read_invalid_group_rule", func(t *testing.T) {
		for _, value := range []interface{}{"{", "[]", true} {
			client := armsNotificationPolicyTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				response := map[string]interface{}{"PageBean": map[string]interface{}{
					"NotificationPolicies": []interface{}{map[string]interface{}{
						"Id": 1234, "Name": "changed", "GroupRule": value,
					}},
				}}
				if err := json.NewEncoder(w).Encode(response); err != nil {
					t.Error(err)
				}
			})
			d := schema.TestResourceDataRaw(t, resourceAliCloudArmsNotificationPolicy().Schema, map[string]interface{}{"name": "unchanged"})
			d.SetId("1234")
			if err := resourceAliCloudArmsNotificationPolicyRead(d, client); err == nil {
				t.Fatal("invalid group rule must return an error")
			}
			if d.Id() != "1234" || d.Get("name") != "unchanged" {
				t.Fatal("invalid response changed state")
			}
		}
	})

	t.Run("empty_grouping_fields", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, resourceAliCloudArmsNotificationPolicy().Schema, map[string]interface{}{
			"group_rule": []interface{}{map[string]interface{}{"grouping_fields": []interface{}{}}},
		})
		value, err := armsNpGroupRuleToJson(d.Get("group_rule").([]interface{}))
		if err != nil {
			t.Fatal(err)
		}
		assertJSONEqual(t, "GroupRule", `{"groupingFields":[],"groupWait":5,"groupInterval":30}`, value)
		if err := d.Set("group_rule", armsNpGroupRuleFromResponse(`{"GroupingFields":[],"GroupWait":5,"GroupInterval":30}`)); err != nil {
			t.Fatal(err)
		}
		if d.Get("group_rule.0.grouping_fields.#") != 0 {
			t.Fatal("empty grouping fields did not round trip")
		}
	})

	t.Run("create_missing_companion", func(t *testing.T) {
		for _, repeat := range []bool{true, false} {
			client := armsNotificationPolicyTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				t.Error("missing companion must be rejected before API request")
				w.WriteHeader(http.StatusBadRequest)
			})
			d := schema.TestResourceDataRaw(t, resourceAliCloudArmsNotificationPolicy().Schema, map[string]interface{}{
				"name": "test-np-missing-companion", "repeat": repeat,
				"notify_rule": []interface{}{map[string]interface{}{}},
			})
			if err := resourceAliCloudArmsNotificationPolicyCreate(d, client); err == nil {
				t.Fatal("missing companion must return an error")
			}
		}
	})

	t.Run("create_read_delete", func(t *testing.T) {
		var deleted string
		client := armsNotificationPolicyTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			var response interface{}
			switch r.Form.Get("Action") {
			case "CreateOrUpdateNotificationPolicy":
				if r.Form.Get("Name") != "test-np" {
					t.Errorf("unexpected Name: %q", r.Form.Get("Name"))
				}
				if r.Form.Get("State") != "enable" {
					t.Errorf("unexpected State: %q", r.Form.Get("State"))
				}
				if r.Form.Get("RepeatInterval") != "600" {
					t.Errorf("unexpected RepeatInterval: %q", r.Form.Get("RepeatInterval"))
				}
				if r.Form.Get("EscalationPolicyId") != "4321" {
					t.Errorf("unexpected EscalationPolicyId: %q", r.Form.Get("EscalationPolicyId"))
				}
				if r.Form.Get("IntegrationId") != "78" {
					t.Errorf("unexpected IntegrationId: %q", r.Form.Get("IntegrationId"))
				}
				if r.Form.Get("DirectedMode") != "true" {
					t.Errorf("unexpected DirectedMode: %q", r.Form.Get("DirectedMode"))
				}
				if r.Form.Get("RegionId") != "cn-hangzhou" {
					t.Errorf("unexpected RegionId: %q", r.Form.Get("RegionId"))
				}
				// NotifyRule: JSON object with camelCase nested keys.
				assertJSONEqual(t, "NotifyRule", `{"notifyStartTime":"00:00","notifyEndTime":"23:59","notifyChannels":["dingTalk","email"],"notifyObjects":[{"notifyObjectType":"CONTACT","notifyObjectId":123,"notifyObjectName":"ex","notifyChannels":["email"]}]}`, r.Form.Get("NotifyRule"))
				// MatchingRules: JSON array of {matchingConditions:[{key,value,operator}]}.
				assertJSONEqual(t, "MatchingRules", `[{"matchingConditions":[{"key":"alertname","value":"ex","operator":"eq"}]}]`, r.Form.Get("MatchingRules"))
				// GroupRule: JSON object with camelCase keys.
				assertJSONEqual(t, "GroupRule", `{"groupingFields":["alertname"],"groupWait":5,"groupInterval":30}`, r.Form.Get("GroupRule"))
				// NotifyTemplate: JSON object with camelCase keys.
				assertJSONEqual(t, "NotifyTemplate", `{"emailTitle":"Alert","emailContent":"fired","robotContent":"bot"}`, r.Form.Get("NotifyTemplate"))
				response = map[string]interface{}{
					"NotificationPolicy": map[string]interface{}{"Id": float64(1234)},
				}
			case "ListNotificationPolicies":
				if r.Form.Get("Page") == "" || r.Form.Get("Size") == "" {
					w.WriteHeader(http.StatusBadRequest)
					response = map[string]interface{}{
						"Code":    "MissingPage",
						"Message": "Page is mandatory for this action.",
					}
					break
				}
				if r.Form.Get("Ids") != "1234" {
					t.Errorf("unexpected Ids: %q", r.Form.Get("Ids"))
				}
				if r.Form.Get("Page") != "1" {
					t.Errorf("unexpected Page: %q", r.Form.Get("Page"))
				}
				if r.Form.Get("RegionId") != "cn-hangzhou" {
					t.Errorf("unexpected RegionId: %q", r.Form.Get("RegionId"))
				}
				// Return the JSON-typed fields as parsed objects/lists (the
				// canonical API response form); the Read path must normalize
				// them to the TF schema shape.
				response = map[string]interface{}{
					"PageBean": map[string]interface{}{
						"NotificationPolicies": []interface{}{
							map[string]interface{}{"Id": float64(9999), "Name": "decoy"},
							map[string]interface{}{
								"Id":                 float64(1234),
								"Name":               "test-np",
								"State":              "enable",
								"SendRecoverMessage": true,
								"Repeat":             true,
								"RepeatInterval":     float64(600),
								"EscalationPolicyId": float64(4321),
								"IntegrationId":      float64(78),
								"DirectedMode":       true,
								"MatchingRules": []interface{}{
									map[string]interface{}{
										"MatchingConditions": []interface{}{
											map[string]interface{}{
												"Key":      "alertname",
												"Value":    "ex",
												"Operator": "eq",
											},
										},
									},
								},
								"GroupRule": map[string]interface{}{
									"GroupingFields": []interface{}{"alertname"},
									"GroupWait":      float64(5),
									"GroupInterval":  float64(30),
								},
								"NotifyRule": map[string]interface{}{
									"NotifyStartTime": "00:00",
									"NotifyEndTime":   "23:59",
									"NotifyChannels":  []interface{}{"dingTalk", "email"},
									"NotifyObjects": []interface{}{
										map[string]interface{}{
											"NotifyObjectType": "CONTACT",
											"NotifyObjectId":   float64(123),
											"NotifyObjectName": "ex",
											"NotifyChannels":   []interface{}{"email"},
										},
									},
								},
								"NotifyTemplate": map[string]interface{}{
									"EmailTitle":   "Alert",
									"EmailContent": "fired",
									"RobotContent": "bot",
								},
							},
						},
					},
				}
			case "DeleteNotificationPolicy":
				deleted = r.Form.Get("Id")
				response = map[string]interface{}{"IsSuccess": true}
			default:
				t.Errorf("unexpected action: %q", r.Form.Get("Action"))
			}
			if err := json.NewEncoder(w).Encode(response); err != nil {
				t.Error(err)
			}
		})

		d := schema.TestResourceDataRaw(t, resourceAliCloudArmsNotificationPolicy().Schema, map[string]interface{}{
			"name":                 "test-np",
			"state":                "enable",
			"send_recover_message": true,
			"repeat_interval":      600,
			"escalation_policy_id": 4321,
			"integration_id":       78,
			"directed_mode":        true,
			"notify_rule": []interface{}{
				map[string]interface{}{
					"notify_start_time": "00:00",
					"notify_end_time":   "23:59",
					"notify_channels":   []interface{}{"dingTalk", "email"},
					"notify_objects": []interface{}{
						map[string]interface{}{
							"notify_object_type": "CONTACT",
							"notify_object_id":   123,
							"notify_object_name": "ex",
							"notify_channels":    []interface{}{"email"},
						},
					},
				},
			},
			"matching_rules": []interface{}{
				map[string]interface{}{
					"matching_conditions": []interface{}{
						map[string]interface{}{
							"key":      "alertname",
							"value":    "ex",
							"operator": "eq",
						},
					},
				},
			},
			"group_rule": []interface{}{
				map[string]interface{}{
					"grouping_fields": []interface{}{"alertname"},
					"group_wait":      5,
					"group_interval":  30,
				},
			},
			"notify_template": []interface{}{
				map[string]interface{}{
					"email_title":   "Alert",
					"email_content": "fired",
					"robot_content": "bot",
				},
			},
		})
		if err := resourceAliCloudArmsNotificationPolicyCreate(d, client); err != nil {
			t.Fatalf("create failed: %v", err)
		}
		if d.Id() != "1234" {
			t.Fatalf("expected id 1234, got %q", d.Id())
		}
		if got := d.Get("name"); got != "test-np" {
			t.Fatalf("read back name %q", got)
		}
		if got := d.Get("notification_policy_id"); got != "1234" {
			t.Fatalf("read back notification_policy_id %q", got)
		}
		if got := d.Get("repeat_interval"); got != 600 {
			t.Fatalf("read back repeat_interval %v", got)
		}
		if got := d.Get("escalation_policy_id"); got != 4321 {
			t.Fatalf("read back escalation_policy_id %v", got)
		}
		if got := d.Get("integration_id"); got != 78 {
			t.Fatalf("read back integration_id %v", got)
		}
		if got := d.Get("directed_mode"); got != true {
			t.Fatalf("read back directed_mode %v", got)
		}
		// matching_rules[*].matching_conditions[0].key == "alertname"
		mr := d.Get("matching_rules").([]interface{})
		if len(mr) != 1 {
			t.Fatalf("expected 1 matching rule, got %d", len(mr))
		}
		conds := mr[0].(map[string]interface{})["matching_conditions"].([]interface{})
		if len(conds) != 1 {
			t.Fatalf("expected 1 condition, got %d", len(conds))
		}
		if conds[0].(map[string]interface{})["key"] != "alertname" {
			t.Fatalf("unexpected condition key: %v", conds[0])
		}
		// group_rule[0].grouping_fields[0] == "alertname"
		gr := d.Get("group_rule").([]interface{})
		if len(gr) != 1 {
			t.Fatalf("expected 1 group rule, got %d", len(gr))
		}
		if fields := gr[0].(map[string]interface{})["grouping_fields"].([]interface{}); len(fields) != 1 || fields[0] != "alertname" {
			t.Fatalf("unexpected grouping_fields: %v", fields)
		}
		// notify_rule[0].notify_channels == ["dingTalk","email"]
		nr := d.Get("notify_rule").([]interface{})
		if len(nr) != 1 {
			t.Fatalf("expected 1 notify rule, got %d", len(nr))
		}
		if chans := nr[0].(map[string]interface{})["notify_channels"].([]interface{}); len(chans) != 2 || chans[0] != "dingTalk" || chans[1] != "email" {
			t.Fatalf("unexpected notify_channels: %v", chans)
		}
		// notify_template[0].email_title == "Alert"
		nt := d.Get("notify_template").([]interface{})
		if len(nt) != 1 {
			t.Fatalf("expected 1 notify template, got %d", len(nt))
		}
		if got := nt[0].(map[string]interface{})["email_title"]; got != "Alert" {
			t.Fatalf("unexpected email_title: %v", got)
		}
		if err := resourceAliCloudArmsNotificationPolicyDelete(d, client); err != nil {
			t.Fatalf("delete failed: %v", err)
		}
		if deleted != "1234" {
			t.Fatalf("expected DeleteNotificationPolicy Id=1234, got %q", deleted)
		}
	})

	// Explicit `false` on the Default:true booleans must reach the API
	// (GetOk would drop it as a zero value), and omitted optional blocks
	// must not be transmitted.
	t.Run("create_bool_false_sent", func(t *testing.T) {
		client := armsNotificationPolicyTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			var response interface{}
			switch r.Form.Get("Action") {
			case "CreateOrUpdateNotificationPolicy":
				if got := r.Form.Get("SendRecoverMessage"); got != "false" {
					t.Errorf("expected SendRecoverMessage=false, got %q", got)
				}
				if got := r.Form.Get("Repeat"); got != "false" {
					t.Errorf("expected Repeat=false, got %q", got)
				}
				if got := r.Form.Get("EscalationPolicyId"); got != "4321" {
					t.Errorf("expected EscalationPolicyId=4321, got %q", got)
				}
				for _, key := range []string{"GroupRule", "MatchingRules", "NotifyTemplate"} {
					if got := r.Form.Get(key); got != "" {
						t.Errorf("omitted block %s must not be sent, got %q", key, got)
					}
				}
				response = map[string]interface{}{
					"NotificationPolicy": map[string]interface{}{"Id": float64(2222)},
				}
			case "ListNotificationPolicies":
				response = map[string]interface{}{
					"PageBean": map[string]interface{}{
						"NotificationPolicies": []interface{}{
							map[string]interface{}{
								"Id":                 float64(2222),
								"Name":               "test-np-false",
								"SendRecoverMessage": false,
								"Repeat":             false,
								"EscalationPolicyId": float64(4321),
								"NotifyRule":         `{"NotifyChannels":["email"]}`,
							},
						},
					},
				}
			default:
				t.Errorf("unexpected action: %q", r.Form.Get("Action"))
			}
			if err := json.NewEncoder(w).Encode(response); err != nil {
				t.Error(err)
			}
		})

		d := schema.TestResourceDataRaw(t, resourceAliCloudArmsNotificationPolicy().Schema, map[string]interface{}{
			"name":                 "test-np-false",
			"send_recover_message": false,
			"repeat":               false,
			"escalation_policy_id": 4321,
			"notify_rule": []interface{}{
				map[string]interface{}{
					"notify_channels": []interface{}{"email"},
				},
			},
		})
		if err := resourceAliCloudArmsNotificationPolicyCreate(d, client); err != nil {
			t.Fatalf("create failed: %v", err)
		}
		if d.Id() != "2222" {
			t.Fatalf("expected id 2222, got %q", d.Id())
		}
		if v, _ := d.Get("send_recover_message").(bool); v {
			t.Fatalf("read back send_recover_message=%v, want false", d.Get("send_recover_message"))
		}
		if v, _ := d.Get("repeat").(bool); v {
			t.Fatalf("read back repeat=%v, want false", d.Get("repeat"))
		}
	})

	t.Run("read_not_found_clears_id", func(t *testing.T) {
		client := armsNotificationPolicyTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			response := map[string]interface{}{
				"PageBean": map[string]interface{}{"NotificationPolicies": []interface{}{
					map[string]interface{}{"Id": float64(9999), "Name": "decoy"},
				}},
			}
			if err := json.NewEncoder(w).Encode(response); err != nil {
				t.Error(err)
			}
		})
		d := schema.TestResourceDataRaw(t, resourceAliCloudArmsNotificationPolicy().Schema, map[string]interface{}{
			"name": "test-np-gone",
			"notify_rule": []interface{}{
				map[string]interface{}{},
			},
		})
		d.SetId("1234")
		if err := resourceAliCloudArmsNotificationPolicyRead(d, client); err != nil {
			t.Fatalf("read must clear the id on not found: %v", err)
		}
		if d.Id() != "" {
			t.Fatalf("expected id to be cleared, got %q", d.Id())
		}
	})

	// JSON-string response variant: the API may return any of the four
	// JSON-typed fields as a JSON string rather than a parsed object/array.
	// The Read path must handle both forms identically.
	t.Run("read_json_string_response", func(t *testing.T) {
		client := armsNotificationPolicyTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			response := map[string]interface{}{
				"PageBean": map[string]interface{}{
					"NotificationPolicies": []interface{}{
						map[string]interface{}{
							"Id":             float64(5678),
							"Name":           "test-np-str",
							"State":          "disable",
							"MatchingRules":  `[{"MatchingConditions":[{"Key":"alertname","Value":"ex","Operator":"eq"}]}]`,
							"GroupRule":      `{"GroupingFields":["alertname"],"GroupWait":5,"GroupInterval":30}`,
							"NotifyRule":     `{"NotifyStartTime":"00:00","NotifyEndTime":"23:59","NotifyChannels":["dingTalk"],"NotifyObjects":[]}`,
							"NotifyTemplate": `{"EmailTitle":"Alert","EmailContent":"fired","RobotContent":"bot"}`,
						},
					},
				},
			}
			if err := json.NewEncoder(w).Encode(response); err != nil {
				t.Error(err)
			}
		})
		d := schema.TestResourceDataRaw(t, resourceAliCloudArmsNotificationPolicy().Schema, map[string]interface{}{
			"name": "test-np-str",
			"notify_rule": []interface{}{
				map[string]interface{}{},
			},
		})
		d.SetId("5678")
		if err := resourceAliCloudArmsNotificationPolicyRead(d, client); err != nil {
			t.Fatalf("read failed: %v", err)
		}
		if got := d.Get("name"); got != "test-np-str" {
			t.Fatalf("read back name %q", got)
		}
		if got := d.Get("state"); got != "disable" {
			t.Fatalf("read back state %q", got)
		}
		mr := d.Get("matching_rules").([]interface{})
		if len(mr) != 1 {
			t.Fatalf("expected 1 matching rule from JSON string, got %d", len(mr))
		}
	})
}

// assertJSONEqual parses both expected and actual as JSON and compares deeply,
// avoiding string-comparison brittleness from map key ordering in json.Marshal.
func assertJSONEqual(t *testing.T, label, expected, actual string) {
	t.Helper()
	var ev, av interface{}
	if err := json.Unmarshal([]byte(expected), &ev); err != nil {
		t.Fatalf("%s: expected JSON not parseable: %v (expected=%q)", label, err, expected)
	}
	if err := json.Unmarshal([]byte(actual), &av); err != nil {
		t.Fatalf("%s: actual JSON not parseable: %v (actual=%q)", label, err, actual)
	}
	normalizeJSONNumbers(ev)
	normalizeJSONNumbers(av)
	if !deepJSONEqual(ev, av) {
		t.Errorf("%s: JSON mismatch\nexpected: %s\nactual:   %s", label, expected, actual)
	}
}

// normalizeJSONNumbers walks the parsed JSON tree and converts all float64
// values that are whole numbers back to int64, so 123 and 123.0 compare equal.
func normalizeJSONNumbers(v interface{}) {
	switch x := v.(type) {
	case map[string]interface{}:
		for k, val := range x {
			normalizeJSONNumbers(val)
			if f, ok := val.(float64); ok && f == float64(int64(f)) {
				x[k] = int64(f)
			}
		}
	case []interface{}:
		for i, val := range x {
			normalizeJSONNumbers(val)
			if f, ok := val.(float64); ok && f == float64(int64(f)) {
				x[i] = int64(f)
			}
		}
	case float64:
		// no-op at top level (handled in container cases)
	}
}

// deepJSONEqual compares two parsed JSON values structurally.
func deepJSONEqual(a, b interface{}) bool {
	switch av := a.(type) {
	case map[string]interface{}:
		bv, ok := b.(map[string]interface{})
		if !ok || len(av) != len(bv) {
			return false
		}
		for k, v := range av {
			if !deepJSONEqual(v, bv[k]) {
				return false
			}
		}
		return true
	case []interface{}:
		bv, ok := b.([]interface{})
		if !ok || len(av) != len(bv) {
			return false
		}
		for i, v := range av {
			if !deepJSONEqual(v, bv[i]) {
				return false
			}
		}
		return true
	default:
		return a == b
	}
}
