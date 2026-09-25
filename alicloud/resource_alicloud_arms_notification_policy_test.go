package alicloud

import (
	"context"
	"encoding/json"
	"errors"
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
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	providerpkg "github.com/aliyun/terraform-provider-alicloud/alicloud/provider"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/arms"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/terraform"
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
	return sweepArmsNotificationPolicy(rawClient.(*connectivity.AliyunClient))
}

// sweepArmsNotificationPolicy lists every page before deleting anything: with
// numeric pagination, deleting while advancing Page shifts the remaining
// policies onto already visited pages and skips them. Deletion then confirms
// both the transport-level result and the DeleteNotificationPolicy IsSuccess
// flag, and failures are returned instead of being logged as a successful run.
func sweepArmsNotificationPolicy(client *connectivity.AliyunClient) error {
	prefixes := []string{
		"tf-testAcc",
		"tf_testacc",
		"terraform-example",
	}

	type policyRef struct {
		id   string
		name string
	}
	var targets []policyRef

	request := make(map[string]interface{})
	query := map[string]interface{}{
		"RegionId": client.RegionId,
		"Page":     1,
		"Size":     PageSizeLarge,
	}
	const listAction = "ListNotificationPolicies"
	for {
		var response map[string]interface{}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err := retry.RetryContext(context.Background(), 2*time.Minute, func() *retry.RetryError {
			var callErr error
			response, callErr = client.RpcGet("ARMS", "2019-08-08", listAction, query, request)
			if callErr != nil {
				if NeedRetry(callErr) {
					wait()
					return retry.RetryableError(callErr)
				}
				return retry.NonRetryableError(callErr)
			}
			return nil
		})
		addDebug(listAction, response, query)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "alicloud_arms_notification_policy", listAction, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.PageBean.NotificationPolicies", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, listAction, "$.PageBean.NotificationPolicies", response)
		}
		result, ok := resp.([]interface{})
		if !ok {
			return WrapErrorf(fmt.Errorf("unexpected %s response type %T for $.PageBean.NotificationPolicies", listAction, resp), FailedGetAttributeMsg, listAction, "$.PageBean.NotificationPolicies", response)
		}
		for _, v := range result {
			item, ok := v.(map[string]interface{})
			if !ok {
				return WrapErrorf(fmt.Errorf("unexpected %s element type %T in $.PageBean.NotificationPolicies", listAction, v), FailedGetAttributeMsg, listAction, "$.PageBean.NotificationPolicies", response)
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
			targets = append(targets, policyRef{id: fmt.Sprint(item["Id"]), name: name})
		}
		if len(result) < PageSizeLarge {
			break
		}
		query["Page"] = query["Page"].(int) + 1
	}

	const deleteAction = "DeleteNotificationPolicy"
	var cleanupErrs []error
	for _, target := range targets {
		log.Printf("[INFO] delete arms notification policy: %s ", target.name)
		deleteRequest := map[string]interface{}{
			"Id": target.id,
		}
		var response map[string]interface{}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err := retry.RetryContext(context.Background(), 2*time.Minute, func() *retry.RetryError {
			var callErr error
			response, callErr = client.RpcPost("ARMS", "2019-08-08", deleteAction, nil, deleteRequest, true)
			if callErr != nil {
				if NeedRetry(callErr) {
					wait()
					return retry.RetryableError(callErr)
				}
				return retry.NonRetryableError(callErr)
			}
			return nil
		})
		addDebug(deleteAction, response, deleteRequest)
		if err != nil {
			// A policy deleted between the listing and this call (or by a
			// concurrent cleanup) is already absent: typed NotFound is a
			// successful cleanup, not a failure.
			if errs.IsTypedNotFound(err) {
				log.Printf("[INFO] arms notification policy %s is already gone", target.id)
				continue
			}
			cleanupErrs = append(cleanupErrs, WrapErrorf(err, DefaultErrorMsg, target.id, deleteAction, AlibabaCloudSdkGoERROR))
			continue
		}
		if success, ok := response["IsSuccess"].(bool); !ok || !success {
			cleanupErrs = append(cleanupErrs, WrapErrorf(fmt.Errorf("delete response did not confirm IsSuccess=true"), DefaultErrorMsg, target.id, deleteAction, AlibabaCloudSdkGoERROR))
		}
	}

	return errors.Join(cleanupErrs...)
}

// testAccCheckArmsNotificationPolicyDestroy accepts only the NotFound that
// DescribeArmsNotificationPolicy constructs itself. A permission or other
// non-404 failure whose message contains "NotFound" must surface as a cleanup
// failure instead of pretending the policy was deleted.
func testAccCheckArmsNotificationPolicyDestroy(s *tfacctest.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	return checkArmsNotificationPolicyDestroyed(s, func(id string) error {
		_, err := (&ArmsServiceV2{client}).DescribeArmsNotificationPolicy(id)
		return err
	})
}

func checkArmsNotificationPolicyDestroyed(s *tfacctest.State, describe func(string) error) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_arms_notification_policy" {
			continue
		}
		err := describe(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("notification policy %s still exists", rs.Primary.ID)
		}
		if errs.IsProviderNotFound(err) {
			continue
		}
		return fmt.Errorf("describing notification policy %s failed: %s", rs.Primary.ID, err)
	}
	return nil
}

// testAccCheckArmsWebhookContactDestroyForCombinedTest mirrors the product-package
// webhook-contact cleanup check for the combined notification-policy test, which
// runs in the root package.
func testAccCheckArmsWebhookContactDestroyForCombinedTest(s *tfacctest.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_arms_webhook_contact" {
			continue
		}
		_, err := arms.NewService(client).DescribeWebhookContact(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("webhook contact %s still exists", rs.Primary.ID)
		}
		if errs.IsProviderNotFound(err) {
			continue
		}
		return fmt.Errorf("describing webhook contact %s failed: %s", rs.Primary.ID, err)
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
		CheckDestroy:  testAccCheckArmsNotificationPolicyDestroy,
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

// TestAccAliCloudArmsNotificationPolicy_directedMode covers the write-side scope
// of DirectedMode: a policy created with directed_mode=true must stay readable,
// importable and destroyable through the same policy id.
func TestAccAliCloudArmsNotificationPolicy_directedMode(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_arms_notification_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudArmsNotificationPolicyMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ArmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeArmsNotificationPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sarmsnpdm%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudArmsNotificationPolicyBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckArmsNotificationPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":                 name,
					"state":                "enable",
					"send_recover_message": true,
					"repeat":               true,
					"repeat_interval":      600,
					"directed_mode":        true,
					"notify_rule": []map[string]interface{}{
						{
							"notify_start_time": "00:00",
							"notify_end_time":   "23:59",
							"notify_channels":   []string{"email"},
							"notify_objects": []map[string]interface{}{
								{
									"notify_object_type": "ARMS_CONTACT",
									"notify_object_id":   "${alicloud_arms_alert_contact.default[0].id}",
									"notify_object_name": "${alicloud_arms_alert_contact.default[0].alert_contact_name}",
									"notify_channels":    []string{"email"},
								},
							},
						},
					},
					"matching_rules": []map[string]interface{}{
						{
							"matching_conditions": []map[string]interface{}{
								{
									"key":      "alertname",
									"value":    "tf-test-directed-alert",
									"operator": "eq",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                   name,
						"state":                  "enable",
						"directed_mode":          "true",
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
	return rac.resourceAttrMapUpdateSet(), testAccCheckArmsNotificationPolicyDestroy
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

// armsNotificationPolicyApplyCreate drives the public Create entry point with an
// explicit raw configuration, so request-level assertions cover the same path
// Terraform uses. schema.TestResourceDataRaw cannot carry a raw configuration,
// which is why the diff is built here and given one explicitly.
func armsNotificationPolicyApplyCreate(t *testing.T, client *connectivity.AliyunClient, raw map[string]interface{}, rawConfig cty.Value) *terraform.InstanceState {
	t.Helper()
	res := resourceAliCloudArmsNotificationPolicy()
	diff, err := res.Diff(context.Background(), nil, terraform.NewResourceConfigRaw(raw), nil)
	if err != nil {
		t.Fatalf("diff failed: %v", err)
	}
	if diff == nil {
		t.Fatal("expected a create diff")
	}
	diff.RawConfig = rawConfig
	state, diags := res.Apply(context.Background(), nil, diff, client)
	if diags.HasError() {
		t.Fatalf("apply diagnostics: %v", diags)
	}
	if state == nil {
		t.Fatal("expected state from apply")
	}
	return state
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

	t.Run("read_malformed_nested_elements", func(t *testing.T) {
		for _, fields := range []map[string]interface{}{
			{"MatchingRules": []interface{}{"not-an-object"}},
			{"MatchingRules": []interface{}{map[string]interface{}{"MatchingConditions": []interface{}{"not-an-object"}}}},
			{"NotifyRule": map[string]interface{}{"NotifyObjects": []interface{}{"not-an-object"}}},
			{"NotifyRule": map[string]interface{}{"NotifyObjects": []interface{}{map[string]interface{}{"NotifyChannels": "not-an-array"}}}},
			{"GroupRule": map[string]interface{}{"GroupingFields": "not-an-array"}},
		} {
			client := armsNotificationPolicyTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				policy := map[string]interface{}{"Id": 1234, "Name": "changed"}
				for key, value := range fields {
					policy[key] = value
				}
				response := map[string]interface{}{"PageBean": map[string]interface{}{
					"NotificationPolicies": []interface{}{policy},
				}}
				if err := json.NewEncoder(w).Encode(response); err != nil {
					t.Error(err)
				}
			})
			d := schema.TestResourceDataRaw(t, resourceAliCloudArmsNotificationPolicy().Schema, map[string]interface{}{"name": "unchanged"})
			d.SetId("1234")
			if err := resourceAliCloudArmsNotificationPolicyRead(d, client); err == nil {
				t.Fatalf("malformed nested response %v must fail the read", fields)
			}
			if d.Id() != "1234" {
				t.Fatalf("malformed nested response removed the id: %q", d.Id())
			}
			if got := d.Get("name"); got != "unchanged" {
				t.Fatalf("malformed nested response polluted state: %q", got)
			}
			if _, ok := d.GetOk("matching_rules"); ok {
				t.Fatal("malformed nested response wrote matching_rules state")
			}
			if _, ok := d.GetOk("notify_rule"); ok {
				t.Fatal("malformed nested response wrote notify_rule state")
			}
			if _, ok := d.GetOk("group_rule"); ok {
				t.Fatal("malformed nested response wrote group_rule state")
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
		groupRule, err := armsNpGroupRuleFromResponse(`{"GroupingFields":[],"GroupWait":5,"GroupInterval":30}`)
		if err != nil {
			t.Fatal(err)
		}
		if err := d.Set("group_rule", groupRule); err != nil {
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

	// directed_mode is Optional+Computed: an explicit false must be sent, while
	// an omitted argument must stay unset so the service default is preserved.
	// Both cases drive the public Create entry point so the assertions cover the
	// same path Terraform uses.
	t.Run("create_directed_mode_false_sent", func(t *testing.T) {
		client := armsNotificationPolicyTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			var response interface{}
			switch r.Form.Get("Action") {
			case "CreateOrUpdateNotificationPolicy":
				if got := r.Form.Get("DirectedMode"); got != "false" {
					t.Errorf("expected DirectedMode=false, got %q", got)
				}
				response = map[string]interface{}{
					"NotificationPolicy": map[string]interface{}{"Id": float64(3333)},
				}
			case "ListNotificationPolicies":
				response = map[string]interface{}{
					"PageBean": map[string]interface{}{
						"NotificationPolicies": []interface{}{
							map[string]interface{}{
								"Id":           float64(3333),
								"Name":         "test-np-directed-false",
								"DirectedMode": false,
								"Repeat":       true,
								"NotifyRule":   `{"NotifyChannels":["email"]}`,
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

		state := armsNotificationPolicyApplyCreate(t, client, map[string]interface{}{
			"name":            "test-np-directed-false",
			"repeat_interval": 600,
			"directed_mode":   false,
			"notify_rule": []interface{}{
				map[string]interface{}{
					"notify_channels": []interface{}{"email"},
				},
			},
		}, cty.ObjectVal(map[string]cty.Value{
			"name":            cty.StringVal("test-np-directed-false"),
			"repeat_interval": cty.NumberIntVal(600),
			"directed_mode":   cty.False,
		}))
		if state.ID != "3333" {
			t.Fatalf("expected id 3333, got %q", state.ID)
		}
		if got := state.Attributes["directed_mode"]; got != "false" {
			t.Fatalf("read back directed_mode=%q, want false", got)
		}
	})

	t.Run("create_directed_mode_omitted", func(t *testing.T) {
		client := armsNotificationPolicyTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			var response interface{}
			switch r.Form.Get("Action") {
			case "CreateOrUpdateNotificationPolicy":
				if got := r.Form.Get("DirectedMode"); got != "" {
					t.Errorf("omitted directed_mode must not be sent, got %q", got)
				}
				response = map[string]interface{}{
					"NotificationPolicy": map[string]interface{}{"Id": float64(4444)},
				}
			case "ListNotificationPolicies":
				response = map[string]interface{}{
					"PageBean": map[string]interface{}{
						"NotificationPolicies": []interface{}{
							map[string]interface{}{
								"Id":         float64(4444),
								"Name":       "test-np-directed-unset",
								"Repeat":     true,
								"NotifyRule": `{"NotifyChannels":["email"]}`,
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

		state := armsNotificationPolicyApplyCreate(t, client, map[string]interface{}{
			"name":            "test-np-directed-unset",
			"repeat_interval": 600,
			"notify_rule": []interface{}{
				map[string]interface{}{
					"notify_channels": []interface{}{"email"},
				},
			},
		}, cty.ObjectVal(map[string]cty.Value{
			"name":            cty.StringVal("test-np-directed-unset"),
			"repeat_interval": cty.NumberIntVal(600),
		}))
		if state.ID != "4444" {
			t.Fatalf("expected id 4444, got %q", state.ID)
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

	// The documented contract returns NotificationPolicies as an array. A null,
	// scalar or object shape must fail the read instead of being treated as an
	// empty page that clears the ID of a live policy.
	t.Run("read_malformed_list_container", func(t *testing.T) {
		for _, container := range []interface{}{nil, "not-an-array", map[string]interface{}{}, float64(1)} {
			client := armsNotificationPolicyTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				response := map[string]interface{}{
					"PageBean": map[string]interface{}{"NotificationPolicies": container},
				}
				if err := json.NewEncoder(w).Encode(response); err != nil {
					t.Error(err)
				}
			})
			d := schema.TestResourceDataRaw(t, resourceAliCloudArmsNotificationPolicy().Schema, map[string]interface{}{"name": "unchanged"})
			d.SetId("1234")
			if err := resourceAliCloudArmsNotificationPolicyRead(d, client); err == nil {
				t.Fatalf("malformed NotificationPolicies container %v must fail the read", container)
			}
			if d.Id() != "1234" {
				t.Fatalf("malformed container %v removed the id: %q", container, d.Id())
			}
			if got := d.Get("name"); got != "unchanged" {
				t.Fatalf("malformed container %v polluted state: %q", container, got)
			}
		}
	})

	t.Run("read_empty_array_clears_id", func(t *testing.T) {
		client := armsNotificationPolicyTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"PageBean": map[string]interface{}{"NotificationPolicies": []interface{}{}},
			})
		})
		d := schema.TestResourceDataRaw(t, resourceAliCloudArmsNotificationPolicy().Schema, map[string]interface{}{"name": "gone"})
		d.SetId("1234")
		if err := resourceAliCloudArmsNotificationPolicyRead(d, client); err != nil {
			t.Fatalf("an empty array is a valid not-found: %v", err)
		}
		if d.Id() != "" {
			t.Fatalf("empty array must clear the id, got %q", d.Id())
		}
	})

	// A successful response carrying both a response-level 404 code and a
	// malformed collection container must fail the read: the container type is
	// validated before the response code is interpreted.
	t.Run("read_code_404_with_malformed_container", func(t *testing.T) {
		for _, container := range []interface{}{nil, "not-an-array", map[string]interface{}{}, float64(1)} {
			client := armsNotificationPolicyTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				response := map[string]interface{}{
					"Code":     "404",
					"PageBean": map[string]interface{}{"NotificationPolicies": container},
				}
				if err := json.NewEncoder(w).Encode(response); err != nil {
					t.Error(err)
				}
			})
			d := schema.TestResourceDataRaw(t, resourceAliCloudArmsNotificationPolicy().Schema, map[string]interface{}{"name": "unchanged"})
			d.SetId("1234")
			if err := resourceAliCloudArmsNotificationPolicyRead(d, client); err == nil {
				t.Fatalf("Code 404 with malformed container %v must fail the read", container)
			}
			if d.Id() != "1234" {
				t.Fatalf("Code 404 with malformed container %v removed the id: %q", container, d.Id())
			}
			if got := d.Get("name"); got != "unchanged" {
				t.Fatalf("Code 404 with malformed container %v polluted state: %q", container, got)
			}
		}
	})

	// A confirmed array keeps the existing not-found semantics even when the
	// response body also carries a 404 code.
	t.Run("read_code_404_with_empty_array_clears_id", func(t *testing.T) {
		client := armsNotificationPolicyTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"Code":     "404",
				"PageBean": map[string]interface{}{"NotificationPolicies": []interface{}{}},
			})
		})
		d := schema.TestResourceDataRaw(t, resourceAliCloudArmsNotificationPolicy().Schema, map[string]interface{}{"name": "gone"})
		d.SetId("1234")
		if err := resourceAliCloudArmsNotificationPolicyRead(d, client); err != nil {
			t.Fatalf("an array with a 404 code is a valid not-found: %v", err)
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

// TestUnitArmsNotificationPolicyConfiguredBool pins the raw-configuration
// check that separates an omitted Optional+Computed boolean from an explicit
// value, including an explicit false.
func TestUnitArmsNotificationPolicyConfiguredBool(t *testing.T) {
	cases := []struct {
		name string
		raw  cty.Value
		want bool
	}{
		{"null config", cty.NullVal(cty.Object(map[string]cty.Type{"directed_mode": cty.Bool})), false},
		{"attribute absent", cty.ObjectVal(map[string]cty.Value{"name": cty.StringVal("x")}), false},
		{"attribute null", cty.ObjectVal(map[string]cty.Value{"directed_mode": cty.NullVal(cty.Bool)}), false},
		{"unknown attribute", cty.ObjectVal(map[string]cty.Value{"directed_mode": cty.DynamicVal}), false},
		{"explicit false", cty.ObjectVal(map[string]cty.Value{"directed_mode": cty.False}), true},
		{"explicit true", cty.ObjectVal(map[string]cty.Value{"directed_mode": cty.True}), true},
		{"non-object config", cty.StringVal("x"), false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := armsNpConfiguredBool(tc.raw, "directed_mode"); got != tc.want {
				t.Fatalf("armsNpConfiguredBool() = %v, want %v", got, tc.want)
			}
		})
	}
}

// TestUnitArmsNotificationPolicyNotFound pins the typed not-found
// classification: an error whose payload merely contains "404" must not be
// treated as NotFound, while a typed HTTP 404 keeps the idempotent
// delete/clear-state behaviour.
func TestUnitArmsNotificationPolicyNotFound(t *testing.T) {
	forbiddenWith404Metadata := func(w http.ResponseWriter) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"Code":    "APINotAuthorize",
			"Message": "The user is not authorized to perform this operation.",
			"Notes":   "unrelated metadata mentioning 404 and the path /api/v2/404",
		})
	}

	forbiddenWithNotFoundMessage := func(w http.ResponseWriter) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"Code":    "APINotAuthorize",
			"Message": "NotFound while resolving authorization for the caller",
		})
	}

	t.Run("read_preserves_id_on_403_with_404_metadata", func(t *testing.T) {
		client := armsNotificationPolicyTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			forbiddenWith404Metadata(w)
		})
		d := schema.TestResourceDataRaw(t, resourceAliCloudArmsNotificationPolicy().Schema, map[string]interface{}{"name": "kept"})
		d.SetId("1234")
		if err := resourceAliCloudArmsNotificationPolicyRead(d, client); err == nil {
			t.Fatal("403 must be returned as an error")
		}
		if d.Id() != "1234" {
			t.Fatalf("read removed a live policy from state: id=%q", d.Id())
		}
		if got := d.Get("name"); got != "kept" {
			t.Fatalf("read changed state after a failed describe: %q", got)
		}
	})

	t.Run("read_preserves_id_on_403_with_notfound_message", func(t *testing.T) {
		client := armsNotificationPolicyTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			forbiddenWithNotFoundMessage(w)
		})
		d := schema.TestResourceDataRaw(t, resourceAliCloudArmsNotificationPolicy().Schema, map[string]interface{}{"name": "kept"})
		d.SetId("1234")
		if err := resourceAliCloudArmsNotificationPolicyRead(d, client); err == nil {
			t.Fatal("403 must be returned as an error")
		}
		if d.Id() != "1234" {
			t.Fatalf("read removed a live policy from state: id=%q", d.Id())
		}
		if got := d.Get("name"); got != "kept" {
			t.Fatalf("read changed state after a failed describe: %q", got)
		}
	})

	t.Run("delete_fails_on_403_with_404_metadata", func(t *testing.T) {
		client := armsNotificationPolicyTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			forbiddenWith404Metadata(w)
		})
		d := schema.TestResourceDataRaw(t, resourceAliCloudArmsNotificationPolicy().Schema, nil)
		d.SetId("1234")
		if err := resourceAliCloudArmsNotificationPolicyDelete(d, client); err == nil {
			t.Fatal("403 must not be reported as a successful delete")
		}
	})

	t.Run("read_clears_id_on_typed_404", func(t *testing.T) {
		client := armsNotificationPolicyTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"Code": "404", "Message": "not found"})
		})
		d := schema.TestResourceDataRaw(t, resourceAliCloudArmsNotificationPolicy().Schema, map[string]interface{}{"name": "gone"})
		d.SetId("1234")
		if err := resourceAliCloudArmsNotificationPolicyRead(d, client); err != nil {
			t.Fatalf("typed 404 must be treated as NotFound: %s", err)
		}
		if d.Id() != "" {
			t.Fatalf("typed 404 must clear the id, got %q", d.Id())
		}
	})

	t.Run("delete_accepts_typed_404", func(t *testing.T) {
		client := armsNotificationPolicyTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"Code": "404", "Message": "not found"})
		})
		d := schema.TestResourceDataRaw(t, resourceAliCloudArmsNotificationPolicy().Schema, nil)
		d.SetId("1234")
		if err := resourceAliCloudArmsNotificationPolicyDelete(d, client); err != nil {
			t.Fatalf("typed 404 delete must be idempotent: %s", err)
		}
	})

	t.Run("check_destroy_fails_on_403_with_notfound_message", func(t *testing.T) {
		client := armsNotificationPolicyTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			forbiddenWithNotFoundMessage(w)
		})
		err := checkArmsNotificationPolicyDestroyed(armsNotificationPolicyDestroyState("1234"), func(id string) error {
			_, describeErr := (&ArmsServiceV2{client}).DescribeArmsNotificationPolicy(id)
			return describeErr
		})
		if err == nil {
			t.Fatal("a non-404 describe failure must fail the cleanup check")
		}
	})

	t.Run("check_destroy_accepts_typed_404", func(t *testing.T) {
		client := armsNotificationPolicyTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"Code": "404", "Message": "not found"})
		})
		err := checkArmsNotificationPolicyDestroyed(armsNotificationPolicyDestroyState("1234"), func(id string) error {
			_, describeErr := (&ArmsServiceV2{client}).DescribeArmsNotificationPolicy(id)
			return describeErr
		})
		if err != nil {
			t.Fatalf("a typed 404 must count as destroyed: %s", err)
		}
	})

	t.Run("check_destroy_fails_when_policy_still_exists", func(t *testing.T) {
		client := armsNotificationPolicyTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"PageBean": map[string]interface{}{
				"NotificationPolicies": []interface{}{map[string]interface{}{"Id": 1234, "Name": "still-here"}},
			}})
		})
		err := checkArmsNotificationPolicyDestroyed(armsNotificationPolicyDestroyState("1234"), func(id string) error {
			_, describeErr := (&ArmsServiceV2{client}).DescribeArmsNotificationPolicy(id)
			return describeErr
		})
		if err == nil {
			t.Fatal("an existing policy must fail the cleanup check")
		}
	})
}

func armsNotificationPolicyDestroyState(id string) *tfacctest.State {
	return &tfacctest.State{Modules: []*tfacctest.ModuleState{{
		Path: []string{"root"},
		Resources: map[string]*tfacctest.ResourceState{
			"alicloud_arms_notification_policy.default": {
				Type:    "alicloud_arms_notification_policy",
				Primary: &tfacctest.InstanceState{ID: id},
			},
		},
	}}}
}

// TestUnitArmsNotificationPolicySweeper pins the sweeper contracts: list every
// page before deleting, delete every matching policy exactly once across
// pages, and surface list/deletion failures including IsSuccess=false.
func TestUnitArmsNotificationPolicySweeper(t *testing.T) {
	t.Run("multi_page_deletes_every_match", func(t *testing.T) {
		pageOne := make([]map[string]interface{}, 0, PageSizeLarge)
		for i := 0; i < 40; i++ {
			pageOne = append(pageOne, map[string]interface{}{"Id": float64(1000 + i), "Name": fmt.Sprintf("tf-testAcc-np-%d", i)})
		}
		for i := 0; i < 10; i++ {
			pageOne = append(pageOne, map[string]interface{}{"Id": float64(2000 + i), "Name": fmt.Sprintf("prod-np-%d", i)})
		}
		pageTwo := []map[string]interface{}{
			{"Id": float64(3000), "Name": "tf-testAcc-np-last-0"},
			{"Id": float64(3001), "Name": "tf_testacc-np-last-1"},
			{"Id": float64(3002), "Name": "tf-testAcc-np-last-2"},
			{"Id": float64(4000), "Name": "prod-np-last-0"},
			{"Id": float64(4001), "Name": "prod-np-last-1"},
		}

		var mu sync.Mutex
		var events []string
		deleteCount := 0
		deleted := map[string]bool{}
		client := armsNotificationPolicyTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			switch r.Form.Get("Action") {
			case "ListNotificationPolicies":
				mu.Lock()
				events = append(events, "list:"+r.Form.Get("Page"))
				mu.Unlock()
				items := pageTwo
				if r.Form.Get("Page") == "1" {
					items = pageOne
				}
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"PageBean": map[string]interface{}{"NotificationPolicies": items},
				})
			case "DeleteNotificationPolicy":
				mu.Lock()
				events = append(events, "delete:"+r.Form.Get("Id"))
				deleteCount++
				deleted[r.Form.Get("Id")] = true
				mu.Unlock()
				_ = json.NewEncoder(w).Encode(map[string]interface{}{"IsSuccess": true})
			default:
				t.Errorf("unexpected action %q", r.Form.Get("Action"))
			}
		})

		if err := sweepArmsNotificationPolicy(client); err != nil {
			t.Fatalf("sweeper returned an error: %s", err)
		}
		firstDelete := -1
		lastList := -1
		var listedPages []string
		for i, event := range events {
			if strings.HasPrefix(event, "delete:") && firstDelete < 0 {
				firstDelete = i
			}
			if strings.HasPrefix(event, "list:") {
				lastList = i
				listedPages = append(listedPages, strings.TrimPrefix(event, "list:"))
			}
		}
		if want := []string{"1", "2"}; strings.Join(listedPages, ",") != strings.Join(want, ",") {
			t.Fatalf("expected a stable read-only pagination over pages %v, got %v", want, listedPages)
		}
		if firstDelete < 0 || firstDelete < lastList {
			t.Fatalf("sweeper deleted before finishing the listing: %v", events)
		}
		if deleteCount != 43 || len(deleted) != 43 {
			t.Fatalf("expected 43 matching policies deleted exactly once, got %d calls / %d ids", deleteCount, len(deleted))
		}
		for i := 0; i < 40; i++ {
			if !deleted[fmt.Sprint(1000+i)] {
				t.Fatalf("page one policy %d was skipped", 1000+i)
			}
		}
		for _, id := range []int{3000, 3001, 3002} {
			if !deleted[fmt.Sprint(id)] {
				t.Fatalf("page two policy %d was skipped", id)
			}
		}
		for _, id := range []int{2000, 2005, 4000, 4001} {
			if deleted[fmt.Sprint(id)] {
				t.Fatalf("retained policy %d was deleted", id)
			}
		}
	})

	t.Run("list_error_propagates", func(t *testing.T) {
		client := armsNotificationPolicyTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"Code": "APINotAuthorize", "Message": "denied"})
		})
		err := sweepArmsNotificationPolicy(client)
		if err == nil {
			t.Fatal("a list failure must not be reported as a successful cleanup")
		}
	})

	t.Run("delete_error_propagates", func(t *testing.T) {
		var mu sync.Mutex
		attempted := map[string]bool{}
		client := armsNotificationPolicyTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			switch r.Form.Get("Action") {
			case "ListNotificationPolicies":
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"PageBean": map[string]interface{}{"NotificationPolicies": []map[string]interface{}{
						{"Id": float64(5001), "Name": "tf-testAcc-np-fail"},
					}},
				})
			case "DeleteNotificationPolicy":
				mu.Lock()
				attempted[r.Form.Get("Id")] = true
				mu.Unlock()
				w.WriteHeader(http.StatusForbidden)
				_ = json.NewEncoder(w).Encode(map[string]interface{}{"Code": "NoPermission", "Message": "denied"})
			}
		})
		if err := sweepArmsNotificationPolicy(client); err == nil {
			t.Fatal("a deletion failure must not be reported as a successful cleanup")
		}
		if !attempted["5001"] {
			t.Fatal("sweeper did not attempt the matching policy")
		}
	})

	t.Run("delete_unconfirmed_is_error", func(t *testing.T) {
		for _, response := range []map[string]interface{}{{"IsSuccess": false}, {}} {
			client := armsNotificationPolicyTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				switch r.Form.Get("Action") {
				case "ListNotificationPolicies":
					_ = json.NewEncoder(w).Encode(map[string]interface{}{
						"PageBean": map[string]interface{}{"NotificationPolicies": []map[string]interface{}{
							{"Id": float64(6001), "Name": "tf-testAcc-np-unconfirmed"},
						}},
					})
				case "DeleteNotificationPolicy":
					_ = json.NewEncoder(w).Encode(response)
				}
			})
			if err := sweepArmsNotificationPolicy(client); err == nil {
				t.Fatalf("delete response %v must be reported as a failure", response)
			}
		}
	})

	t.Run("delete_not_found_is_idempotent", func(t *testing.T) {
		var mu sync.Mutex
		attempted := map[string]bool{}
		client := armsNotificationPolicyTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			switch r.Form.Get("Action") {
			case "ListNotificationPolicies":
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"PageBean": map[string]interface{}{"NotificationPolicies": []map[string]interface{}{
						{"Id": float64(7001), "Name": "tf-testAcc-np-already-gone"},
						{"Id": float64(7002), "Name": "tf-testAcc-np-live"},
					}},
				})
			case "DeleteNotificationPolicy":
				mu.Lock()
				attempted[r.Form.Get("Id")] = true
				mu.Unlock()
				if r.Form.Get("Id") == "7001" {
					w.WriteHeader(http.StatusNotFound)
					_ = json.NewEncoder(w).Encode(map[string]interface{}{"Code": "404", "Message": "not found"})
					return
				}
				_ = json.NewEncoder(w).Encode(map[string]interface{}{"IsSuccess": true})
			}
		})
		if err := sweepArmsNotificationPolicy(client); err != nil {
			t.Fatalf("an already-deleted policy must not fail the cleanup: %s", err)
		}
		if !attempted["7001"] || !attempted["7002"] {
			t.Fatalf("sweeper did not process every target: %v", attempted)
		}
	})

	t.Run("malformed_list_fails", func(t *testing.T) {
		for _, pageBean := range []map[string]interface{}{
			{"NotificationPolicies": "not-an-array"},
			{"NotificationPolicies": []interface{}{"not-an-object"}},
		} {
			client := armsNotificationPolicyTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode(map[string]interface{}{"PageBean": pageBean})
			})
			if err := sweepArmsNotificationPolicy(client); err == nil {
				t.Fatalf("malformed list response %v must fail the sweeper", pageBean)
			}
		}
	})

	t.Run("documented_example_prefix_is_swept", func(t *testing.T) {
		var mu sync.Mutex
		deleted := map[string]bool{}
		client := armsNotificationPolicyTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			switch r.Form.Get("Action") {
			case "ListNotificationPolicies":
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"PageBean": map[string]interface{}{"NotificationPolicies": []map[string]interface{}{
						{"Id": float64(8001), "Name": "tf-testAcc-np"},
						{"Id": float64(8002), "Name": "terraform-example"},
						{"Id": float64(8003), "Name": "prod-np"},
					}},
				})
			case "DeleteNotificationPolicy":
				mu.Lock()
				deleted[r.Form.Get("Id")] = true
				mu.Unlock()
				_ = json.NewEncoder(w).Encode(map[string]interface{}{"IsSuccess": true})
			}
		})
		if err := sweepArmsNotificationPolicy(client); err != nil {
			t.Fatalf("sweeper returned an error: %s", err)
		}
		if !deleted["8001"] || !deleted["8002"] {
			t.Fatalf("test and documented-example prefixes must be swept: %v", deleted)
		}
		if deleted["8003"] {
			t.Fatalf("unrelated policy must be retained: %v", deleted)
		}
	})
}

// TestUnitArmsNotificationPolicySchemaEnums pins that every value the
// documentation lists as valid is enforced by the schema, and that the
// documented enum sets match the published ARMS 2019-08-08 contract.
func TestUnitArmsNotificationPolicySchemaEnums(t *testing.T) {
	res := resourceAliCloudArmsNotificationPolicy()

	notifyRule := res.Schema["notify_rule"].Elem.(*schema.Resource)
	notifyObjects := notifyRule.Schema["notify_objects"].Elem.(*schema.Resource)
	matchingRules := res.Schema["matching_rules"].Elem.(*schema.Resource)
	matchingConditions := matchingRules.Schema["matching_conditions"].Elem.(*schema.Resource)

	cases := []struct {
		name    string
		field   *schema.Schema
		valid   []string
		invalid []string
	}{
		{
			name:    "state",
			field:   res.Schema["state"],
			valid:   []string{"enable", "disable"},
			invalid: []string{"enabled", "Enable"},
		},
		{
			name:    "notify_rule.notify_channels",
			field:   notifyRule.Schema["notify_channels"].Elem.(*schema.Schema),
			valid:   []string{"dingTalk", "email", "sms", "tts", "webhook"},
			invalid: []string{"dingding", "webHook"},
		},
		{
			name:    "notify_rule.notify_objects.notify_object_type",
			field:   notifyObjects.Schema["notify_object_type"],
			valid:   []string{"CONTACT", "CONTACT_GROUP", "ARMS_CONTACT", "ARMS_CONTACT_GROUP", "DING_ROBOT_GROUP", "CONTACT_SCHEDULE"},
			invalid: []string{"ROBOT", "contact"},
		},
		{
			name:    "notify_rule.notify_objects.notify_channels",
			field:   notifyObjects.Schema["notify_channels"].Elem.(*schema.Schema),
			valid:   []string{"email", "sms", "tts"},
			invalid: []string{"dingTalk", "webhook"},
		},
		{
			name:    "matching_rules.matching_conditions.operator",
			field:   matchingConditions.Schema["operator"],
			valid:   []string{"eq", "neq", "in", "nin", "re", "nre"},
			invalid: []string{"equals", "EQ"},
		},
	}

	for _, c := range cases {
		if c.field.ValidateFunc == nil {
			t.Fatalf("%s must declare a ValidateFunc", c.name)
		}
		for _, value := range c.valid {
			if _, errs := c.field.ValidateFunc(value, c.name); len(errs) != 0 {
				t.Fatalf("%s rejected documented value %q: %v", c.name, value, errs)
			}
		}
		for _, value := range c.invalid {
			if _, errs := c.field.ValidateFunc(value, c.name); len(errs) == 0 {
				t.Fatalf("%s accepted undocumented value %q", c.name, value)
			}
		}
	}
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

// TestAccAliCloudArmsNotificationPolicy_webhookContact covers the combined
// lifecycle with alicloud_arms_webhook_contact: the notification policy is
// created from the webhook contact's real output, a second plan is empty, and
// Terraform's dependency ordering destroys the policy before the webhook
// contact it references.
func TestAccAliCloudArmsNotificationPolicy_webhookContact(t *testing.T) {
	var np map[string]interface{}
	var wc map[string]interface{}
	npId := "alicloud_arms_notification_policy.default"
	wcId := "alicloud_arms_webhook_contact.default"

	npRa := resourceAttrInit(npId, AlicloudArmsNotificationPolicyMap)
	npRc := resourceCheckInitWithDescribeMethod(npId, &np, func() interface{} {
		return &ArmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeArmsNotificationPolicy")
	npRac := resourceAttrCheckInit(npRc, npRa)
	npCheck := npRac.resourceAttrMapUpdateSet()

	wcRa := resourceAttrInit(wcId, nil)
	wcRc := resourceCheckInitWithDescribeMethod(wcId, &wc, func() interface{} {
		return arms.NewService(testAccProvider.Meta().(*connectivity.AliyunClient))
	}, "DescribeWebhookContact")
	wcRac := resourceAttrCheckInit(wcRc, wcRa)
	wcCheck := wcRac.resourceAttrMapUpdateSet()

	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sarmswebhooknp%d", defaultRegionToTest, rand)
	config := armsNotificationPolicyWebhookContactDependence(name)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName:            npId,
		ProtoV5ProviderFactories: testAccMuxedProviderFactories(),
		CheckDestroy:             resource.ComposeTestCheckFunc(testAccCheckArmsNotificationPolicyDestroy, testAccCheckArmsWebhookContactDestroyForCombinedTest),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					npCheck(map[string]string{
						"name":                            name,
						"state":                           "enable",
						"notify_rule.0.notify_channels.#": "1",
						"notify_rule.0.notify_channels.0": "webhook",
						"notify_rule.0.notify_objects.#":  "1",
						"notify_rule.0.notify_objects.0.notify_object_type": "CONTACT",
						"notify_rule.0.notify_objects.0.notify_object_name": name,
					}),
					wcCheck(map[string]string{
						"webhook_contact_name": name,
						"method":               "Post",
						"url":                  "https://example.com/webhook",
						"body":                 "alert fired",
						"recover_body":         "alert recovered",
					}),
					resource.TestCheckResourceAttrPair(npId, "notify_rule.0.notify_objects.0.notify_object_id", wcId, "id"),
				),
			},
			{
				Config:   config,
				PlanOnly: true,
			},
		},
	})
}

// testAccMuxedProviderFactories serves the SDK v2 and Framework halves behind one
// protocol v5 server, so one configuration may combine an SDK v2 resource with a
// Framework resource (alicloud_arms_webhook_contact is implemented in
// alicloud/service/arms). acctest.ProtoV5ProviderFactories is unusable from this
// package: importing alicloud/acctest would close an import cycle through the
// root provider.
func testAccMuxedProviderFactories() map[string]func() (tfprotov5.ProviderServer, error) {
	return map[string]func() (tfprotov5.ProviderServer, error){
		"alicloud": func() (tfprotov5.ProviderServer, error) {
			serverFactory, err := providerpkg.ProtoV5ProviderServerFactory(context.Background(), testAccProvider)
			if err != nil {
				return nil, err
			}
			return serverFactory(), nil
		},
	}
}

func armsNotificationPolicyWebhookContactDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_arms_webhook_contact" "default" {
  webhook_contact_name = var.name
  method               = "Post"
  url                  = "https://example.com/webhook"
  body                 = "alert fired"
  recover_body         = "alert recovered"
}

resource "alicloud_arms_notification_policy" "default" {
  name            = var.name
  repeat_interval = 600
  state           = "enable"

  notify_rule {
    notify_start_time = "00:00"
    notify_end_time   = "23:59"
    notify_channels   = ["webhook"]

    notify_objects {
      notify_object_type = "CONTACT"
      notify_object_id   = alicloud_arms_webhook_contact.default.id
      notify_object_name = alicloud_arms_webhook_contact.default.webhook_contact_name
    }
  }

  matching_rules {
    matching_conditions {
      key      = "alertname"
      value    = "example-alert"
      operator = "eq"
    }
  }

  group_rule {
    grouping_fields = ["alertname"]
    group_wait      = 5
    group_interval  = 30
  }
}
`, name)
}
