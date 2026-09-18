package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
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

	action := "ListNotificationPolicies"
	request := map[string]interface{}{
		"RegionId": client.RegionId,
		"Page":     1,
		"Size":     PageSizeXLarge,
	}
	var response map[string]interface{}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("ARMS", "2019-08-08", action, nil, request, true)
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
		log.Printf("[ERROR] %s failed: %v", action, err)
		return nil
	}
	resp, err := jsonpath.Get("$.PageBean.NotificationPolicies", response)
	if err != nil {
		log.Printf("[ERROR] %v", WrapError(err))
		return nil
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
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
				log.Printf("[INFO] Skipping notification policy: %s ", name)
				continue
			}
		}
		log.Printf("[INFO] delete notification policy: %s ", name)

		action = "DeleteNotificationPolicy"
		request = map[string]interface{}{
			"Id":       fmt.Sprint(item["Id"]),
			"RegionId": client.RegionId,
		}
		wait = incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(1*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("ARMS", "2019-08-08", action, nil, request, false)
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
			log.Printf("[ERROR] %s failed: %v", action, err)
		}
	}
	return nil
}

func TestAccAliCloudARMSNotificationPolicy_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_arms_notification_policy.default"
	ra := resourceAttrInit(resourceId, ArmsNotificationPolicyMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ArmsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "ListArmsNotificationPolicies")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccArmsNotificationPolicy%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, ArmsNotificationPolicyBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.ARMSSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name": "${var.name}",
					"notify_rule": []map[string]interface{}{
						{
							"notify_start_time": "00:00",
							"notify_end_time":   "23:59",
							"notify_channels":   []string{"dingTalk", "email"},
							"notify_objects": []map[string]interface{}{
								{
									"notify_object_type": "ARMS_CONTACT",
									"notify_object_id":   "${alicloud_arms_alert_contact.default.id}",
									"notify_object_name": "${var.name}",
								},
							},
						},
					},
					"group_rule": []map[string]interface{}{
						{
							"group_wait":      "5",
							"group_interval":  "30",
							"grouping_fields": []string{"alertname"},
						},
					},
					"matching_rules": []map[string]interface{}{
						{
							"matching_conditions": []map[string]interface{}{
								{
									"key":      "_aliyun_arms_involvedObject_kind",
									"value":    "app",
									"operator": "eq",
								},
							},
						},
					},
					"repeat":          "true",
					"repeat_interval": "600",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":              CHECKSET,
						"notify_rule.#":     "1",
						"group_rule.#":      "1",
						"matching_rules.#":  "1",
						"repeat":            "true",
						"repeat_interval":   "600",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"notify_rule": []map[string]interface{}{
						{
							"notify_start_time": "01:00",
							"notify_end_time":   "22:00",
							"notify_channels":   []string{"dingTalk"},
							"notify_objects": []map[string]interface{}{
								{
									"notify_object_type": "ARMS_CONTACT",
									"notify_object_id":   "${alicloud_arms_alert_contact.default.id}",
									"notify_object_name": "${var.name}",
								},
								{
									"notify_object_type": "ARMS_CONTACT_GROUP",
									"notify_object_id":   "${alicloud_arms_alert_contact_group.default.id}",
									"notify_object_name": "${var.name}",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"notify_rule.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"repeat_interval": "300",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"repeat_interval": "300",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"notify_template": []map[string]interface{}{
						{
							"email_title":           "CreateEmailTitle",
							"email_content":         "CreateEmailContent",
							"email_recover_title":   "CreateEmailRecoverTitle",
							"email_recover_content": "CreateEmailRecoverContent",
							"sms_content":           "CreateSmsContent",
							"sms_recover_content":   "CreateSmsRecoverContent",
							"tts_content":           "CreateTtsContent",
							"tts_recover_content":   "CreateTtsRecoverContent",
							"robot_content":         "CreateRobotContent",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"notify_template.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"send_recover_message": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"send_recover_message": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": "${var.name}",
					"notify_rule": []map[string]interface{}{
						{
							"notify_start_time": "00:00",
							"notify_end_time":   "23:59",
							"notify_channels":   []string{"dingTalk", "email"},
							"notify_objects": []map[string]interface{}{
								{
									"notify_object_type": "ARMS_CONTACT",
									"notify_object_id":   "${alicloud_arms_alert_contact.default.id}",
									"notify_object_name": "${var.name}",
								},
							},
						},
					},
					"group_rule": []map[string]interface{}{
						{
							"group_wait":      "5",
							"group_interval":  "30",
							"grouping_fields": []string{"alertname"},
						},
					},
					"repeat":          "true",
					"repeat_interval": "600",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":            CHECKSET,
						"notify_rule.#":   "1",
						"group_rule.#":    "1",
						"repeat":          "true",
						"repeat_interval": "600",
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

var ArmsNotificationPolicyMap = map[string]string{}

func ArmsNotificationPolicyBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
resource "alicloud_arms_alert_contact" "default" {
  alert_contact_name = "${var.name}"
  email = "${var.name}@aaa.com"
}
resource "alicloud_arms_alert_contact_group" "default" {
  alert_contact_group_name = "${var.name}"
  contact_ids = [alicloud_arms_alert_contact.default.id]
}
`, name)
}
