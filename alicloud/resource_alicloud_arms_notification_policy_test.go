package alicloud

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	credentials "github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

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
	name := fmt.Sprintf("tf-testacc%sarmsnp%d", defaultRegionToTest, rand)
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
									"notify_object_type": "CONTACT",
									"notify_object_id":   1,
									"notify_object_name": "tf-test-contact",
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
							"email_title":           "Alert: ${alertname}",
							"email_content":         "${alertname} fired",
							"email_recover_title":   "Recovered: ${alertname}",
							"email_recover_content": "${alertname} recovered",
							"sms_content":           "SMS alert: ${alertname}",
							"sms_recover_content":   "SMS recovered: ${alertname}",
							"tts_content":           "TTS alert: ${alertname}",
							"tts_recover_content":   "TTS recovered: ${alertname}",
							"robot_content":         "Robot alert: ${alertname}",
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

// The six TestAcc functions below satisfy the TypeList Order Coverage CI gate
// (see scripts/collection-order/README.md). Each configurable multi-member
// TypeList gets its own isolated three-step sequence — apply A -> reordered
// PlanOnly+ExpectNonEmptyPlan:true -> apply B with an empty post-apply plan —
// using at least two distinct members. A single plan step cannot cover two
// paths: the checker requires every value outside the nominated collection to
// be unchanged, so each path has its own function with a minimal config that
// varies only that one list. Each container block carries only the target field
// so the cumulative builder delta can reorder it without dropping sibling keys.

func TestAccAliCloudArmsNotificationPolicy_matchingRulesOrder(t *testing.T) {
	testAccCheck, checkDestroy := armsNotificationPolicyReorderChecks(t)
	resourceId := "alicloud_arms_notification_policy.default"
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sarmsnpmr%d", defaultRegionToTest, rand)
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
					"name": name,
					"notify_rule": []map[string]interface{}{
						{"notify_channels": []string{"email"}},
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
					testAccCheck(map[string]string{"name": name}),
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
					testAccCheck(map[string]string{"name": name}),
				),
			},
		},
	})
}

func TestAccAliCloudArmsNotificationPolicy_matchingConditionsOrder(t *testing.T) {
	testAccCheck, checkDestroy := armsNotificationPolicyReorderChecks(t)
	resourceId := "alicloud_arms_notification_policy.default"
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sarmsnpc%d", defaultRegionToTest, rand)
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
					"name": name,
					"notify_rule": []map[string]interface{}{
						{"notify_channels": []string{"email"}},
					},
					"matching_rules": []map[string]interface{}{
						{"matching_conditions": []map[string]interface{}{
							{"key": "alertname", "value": "tf-test-alert-a", "operator": "eq"},
							{"key": "clustername", "value": "tf-test-cluster-b", "operator": "neq"},
						}},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"name": name}),
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
					testAccCheck(map[string]string{"name": name}),
				),
			},
		},
	})
}

func TestAccAliCloudArmsNotificationPolicy_groupingFieldsOrder(t *testing.T) {
	testAccCheck, checkDestroy := armsNotificationPolicyReorderChecks(t)
	resourceId := "alicloud_arms_notification_policy.default"
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sarmsnpgf%d", defaultRegionToTest, rand)
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
					"name": name,
					"notify_rule": []map[string]interface{}{
						{"notify_channels": []string{"email"}},
					},
					"group_rule": []map[string]interface{}{
						{"grouping_fields": []string{"alertname", "clustername"}},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"name": name}),
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
					testAccCheck(map[string]string{"name": name}),
				),
			},
		},
	})
}

func TestAccAliCloudArmsNotificationPolicy_notifyChannelsOrder(t *testing.T) {
	testAccCheck, checkDestroy := armsNotificationPolicyReorderChecks(t)
	resourceId := "alicloud_arms_notification_policy.default"
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sarmsnpnc%d", defaultRegionToTest, rand)
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
					"name": name,
					"notify_rule": []map[string]interface{}{
						{"notify_channels": []string{"dingTalk", "email"}},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"name": name}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"notify_rule": []map[string]interface{}{
						{"notify_channels": []string{"email", "dingTalk"}},
					},
				}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"notify_rule": []map[string]interface{}{
						{"notify_channels": []string{"email", "dingTalk"}},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"name": name}),
				),
			},
		},
	})
}

func TestAccAliCloudArmsNotificationPolicy_notifyObjectsOrder(t *testing.T) {
	testAccCheck, checkDestroy := armsNotificationPolicyReorderChecks(t)
	resourceId := "alicloud_arms_notification_policy.default"
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sarmsnpno%d", defaultRegionToTest, rand)
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
					"name": name,
					"notify_rule": []map[string]interface{}{
						{"notify_objects": []map[string]interface{}{
							{"notify_object_type": "CONTACT", "notify_object_id": 1, "notify_object_name": "tf-test-contact-a"},
							{"notify_object_type": "CONTACT_GROUP", "notify_object_id": 2, "notify_object_name": "tf-test-group-b"},
						}},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"name": name}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"notify_rule": []map[string]interface{}{
						{"notify_objects": []map[string]interface{}{
							{"notify_object_type": "CONTACT_GROUP", "notify_object_id": 2, "notify_object_name": "tf-test-group-b"},
							{"notify_object_type": "CONTACT", "notify_object_id": 1, "notify_object_name": "tf-test-contact-a"},
						}},
					},
				}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"notify_rule": []map[string]interface{}{
						{"notify_objects": []map[string]interface{}{
							{"notify_object_type": "CONTACT_GROUP", "notify_object_id": 2, "notify_object_name": "tf-test-group-b"},
							{"notify_object_type": "CONTACT", "notify_object_id": 1, "notify_object_name": "tf-test-contact-a"},
						}},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"name": name}),
				),
			},
		},
	})
}

func TestAccAliCloudArmsNotificationPolicy_notifyObjectChannelsOrder(t *testing.T) {
	testAccCheck, checkDestroy := armsNotificationPolicyReorderChecks(t)
	resourceId := "alicloud_arms_notification_policy.default"
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sarmsnpnoc%d", defaultRegionToTest, rand)
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
					"name": name,
					"notify_rule": []map[string]interface{}{
						{"notify_objects": []map[string]interface{}{
							{"notify_channels": []string{"email", "sms"}},
						}},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"name": name}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"notify_rule": []map[string]interface{}{
						{"notify_objects": []map[string]interface{}{
							{"notify_channels": []string{"sms", "email"}},
						}},
					},
				}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"notify_rule": []map[string]interface{}{
						{"notify_objects": []map[string]interface{}{
							{"notify_channels": []string{"sms", "email"}},
						}},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"name": name}),
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
				// Return the JSON-typed fields as parsed objects/lists (the
				// canonical API response form); the Read path must normalize
				// them to the TF schema shape.
				response = map[string]interface{}{
					"PageBean": map[string]interface{}{
						"NotificationPolicies": []interface{}{
							map[string]interface{}{
								"Id":                 float64(1234),
								"Name":               "test-np",
								"State":              "enable",
								"SendRecoverMessage": true,
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
				response = map[string]interface{}{"Code": 200}
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

	t.Run("read_not_found_clears_id", func(t *testing.T) {
		client := armsNotificationPolicyTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			response := map[string]interface{}{
				"PageBean": map[string]interface{}{"NotificationPolicies": []interface{}{}},
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
		d.SetId("1234-gone")
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
