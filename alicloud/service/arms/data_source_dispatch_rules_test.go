package arms_test

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudARMSDispatchRulesDataSource(t *testing.T) {

	rand := tfacctest.RandInt()
	resourceId := "data.alicloud_arms_dispatch_rules.default"
	name := fmt.Sprintf("tf-testacc-armsDispatchRules%v", rand)

	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourceArmsDispatchRulesConfigDependence)

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"name_regex":     "${alicloud_arms_dispatch_rule.default.dispatch_rule_name}",
			"enable_details": "true",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "fake_tf-testacc*",
		}),
	}

	nameConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"dispatch_rule_name": "${alicloud_arms_dispatch_rule.default.dispatch_rule_name}",
			"enable_details":     "true",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"dispatch_rule_name": "fake_tf-testacc*",
		}),
	}
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_arms_dispatch_rule.default.id}"},
			"enable_details": "true",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_arms_dispatch_rule.default.id}_fake"},
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"name_regex":         "${alicloud_arms_dispatch_rule.default.dispatch_rule_name}",
			"dispatch_rule_name": "${alicloud_arms_dispatch_rule.default.dispatch_rule_name}",
			"ids":                []string{"${alicloud_arms_dispatch_rule.default.id}"},
			"enable_details":     "true",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"name_regex":         "${alicloud_arms_dispatch_rule.default.dispatch_rule_name}",
			"dispatch_rule_name": "${alicloud_arms_dispatch_rule.default.dispatch_rule_name}_fake",
			"ids":                []string{"${alicloud_arms_dispatch_rule.default.id}"},
		}),
	}
	var existArmsDispatchRulesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                 "1",
			"names.#":                               "1",
			"rules.#":                               "1",
			"rules.0.id":                            acctest.CHECKSET,
			"rules.0.dispatch_rule_id":              acctest.CHECKSET,
			"rules.0.group_rules.#":                 "1",
			"rules.0.status":                        "enable",
			"rules.0.label_match_expression_grid.#": "1",
			"rules.0.notify_rules.#":                "1",
			"rules.0.notify_template.#":             "1",
			"rules.0.dispatch_rule_name":            name,
		}
	}

	var fakeArmsDispatchRulesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"rules.#": "0",
			"names.#": "0",
			"ids.#":   "0",
		}
	}

	var ArmsDispatchRulesCheckInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existArmsDispatchRulesMapFunc,
		FakeMapFunc:  fakeArmsDispatchRulesMapFunc,
	}
	preCheck := func() {
		acctest.PreCheck(t)
		acctest.PreCheckWithRegions(t, true, connectivity.ARMSSupportRegions)
	}
	ArmsDispatchRulesCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, nameRegexConf, nameConf, idsConf, allConf)
}

func dataSourceArmsDispatchRulesConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
 default = "%v"
}

resource "alicloud_arms_alert_robot" "default" {
  alert_robot_name = var.name
  robot_type       = "dingding"
  robot_addr       = "https://oapi.dingtalk.com/robot/send?access_token=1c704e23"
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
    grouping_fields = ["alertname"]
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
      notify_object_id = alicloud_arms_alert_robot.default.id
      notify_type      = "ARMS_ROBOT"
      name             = var.name
    }
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

  notify_template {
    email_title           = "example_email_title"
    email_content         = "example_email_content"
    email_recover_title   = "example_email_recover_title"
    email_recover_content = "example_email_recover_content"
    sms_content           = "example_sms_content"
    sms_recover_content   = "example_sms_recover_content"
    tts_content           = "example_tts_content"
    tts_recover_content   = "example_tts_recover_content"
    robot_content         = "example_robot_content"
  }
}
`, name)
}

// TestAccAliCloudARMSDispatchRulesDataSource_multiPage proves the data source reads
// past the first page: ARMS returns at most PageSizeXLarge policies per call, so the
// test seeds one rule more than a page and asserts the data source resolves every
// seeded rule. Seeding creates its own alert contact and rules, and cleans them up
// afterwards; an account that cannot seed is skipped instead of failing, because the
// seed is bookkeeping rather than the behaviour under test.
func TestAccAliCloudARMSDispatchRulesDataSource_multiPage(t *testing.T) {
	// Explicit ACC opt-in before any client or cloud access: acctest.PreCheck and
	// resource.Test only skip later, so without this gate a plain `go test` with
	// configured credentials would create the contact and 101 rules first.
	if os.Getenv("TF_ACC") == "" {
		t.Skip("TF_ACC is not set; skipping the multi-page acceptance test without any cloud access")
	}
	acctest.PreCheck(t)
	acctest.PreCheckWithRegions(t, true, connectivity.ARMSSupportRegions)

	region := os.Getenv("ALICLOUD_REGION")
	if region == "" {
		region = string(connectivity.Hangzhou)
	}
	rawClient, err := armsTestClientForRegion(region)
	if err != nil {
		t.Fatalf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	seededCount := helper.PageSizeXLarge + 1
	prefix := fmt.Sprintf("tf-testacc-multipage-%d", tfacctest.RandInt())
	contactID, contactName, err := createArmsAlertContact(client, prefix)
	if err != nil {
		t.Skipf("skipping the multi-page assertion: cannot create the seeding alert contact: %s", err)
	}
	defer deleteArmsAlertContact(client, contactID)
	seededIDs, err := seedArmsDispatchRules(client, prefix, armsDispatchRuleSeedTemplate(contactID, contactName), seededCount)
	defer deleteArmsDispatchRules(client, seededIDs)
	if err != nil {
		t.Skipf("skipping the multi-page assertion: cannot seed %d dispatch rules: %s", seededCount, err)
	}

	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
data "alicloud_arms_dispatch_rules" "default" {
  name_regex = "^%s-"
}
`, prefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_arms_dispatch_rules.default", "rules.#", strconv.Itoa(seededCount)),
					resource.TestCheckResourceAttr("data.alicloud_arms_dispatch_rules.default", "names.#", strconv.Itoa(seededCount)),
				),
			},
		},
	})
}

// seedArmsDispatchRules creates count dispatch rules sharing prefix and returns the
// ids it managed to create, so a partial failure still cleans up after the test.
func seedArmsDispatchRules(client *connectivity.AliyunClient, prefix string, template map[string]interface{}, count int) ([]string, error) {
	ids := make([]string, 0, count)
	for index := 0; index < count; index++ {
		template["name"] = fmt.Sprintf("%s-%03d", prefix, index)
		response, err := createArmsDispatchRule(client, template)
		if err != nil {
			return ids, err
		}
		ids = append(ids, fmt.Sprint(response["DispatchRuleId"]))
	}
	return ids, nil
}

// armsDispatchRuleSeedTemplate is the smallest dispatch rule the account accepts,
// wired to the seeding alert contact.
func armsDispatchRuleSeedTemplate(contactID, contactName string) map[string]interface{} {
	return map[string]interface{}{
		"name":         "",
		"dispatchType": "CREATE_ALERT",
		"groupRules": []interface{}{
			map[string]interface{}{
				"groupWait":      5,
				"groupInterval":  15,
				"groupingFields": []interface{}{"alertname"},
				"repeatInterval": 61,
			},
		},
		"labelMatchExpressionGrid": map[string]interface{}{
			"labelMatchExpressionGroups": []interface{}{
				map[string]interface{}{
					"labelMatchExpressions": []interface{}{
						map[string]interface{}{"key": "alertname", "operator": "eq", "value": "multipage"},
					},
				},
			},
		},
		"notifyRules": []interface{}{
			map[string]interface{}{
				"notifyChannels":  []interface{}{"dingTalk"},
				"notifyObjects":   []interface{}{map[string]interface{}{"notifyType": "ARMS_CONTACT", "name": contactName, "notifyObjectId": contactID}},
				"notifyStartTime": "00:00",
				"notifyEndTime":   "23:59",
			},
		},
	}
}

func createArmsAlertContact(client *connectivity.AliyunClient, prefix string) (string, string, error) {
	contactName := prefix + "-contact"
	response, err := client.RpcPost("ARMS", "2019-08-08", "CreateAlertContact", nil, map[string]interface{}{
		"RegionId":    client.RegionId,
		"ContactName": contactName,
		"Email":       contactName + "@example.com",
	}, false)
	if err != nil {
		return "", "", err
	}
	return fmt.Sprint(response["ContactId"]), contactName, nil
}

func deleteArmsAlertContact(client *connectivity.AliyunClient, id string) {
	if id == "" {
		return
	}
	if _, err := client.RpcPost("ARMS", "2019-08-08", "DeleteAlertContact", nil, map[string]interface{}{
		"ContactId": id,
		"RegionId":  client.RegionId,
	}, false); err != nil {
		log.Printf("[WARN] failed to delete the seeding alert contact %s: %s", id, err)
	}
}

func createArmsDispatchRule(client *connectivity.AliyunClient, rule map[string]interface{}) (map[string]interface{}, error) {
	ruleJSON, err := convert.ConvertMaptoJsonString(rule)
	if err != nil {
		return nil, err
	}
	return client.RpcPost("ARMS", "2019-08-08", "CreateDispatchRule", nil, map[string]interface{}{
		"RegionId":     client.RegionId,
		"DispatchRule": ruleJSON,
	}, false)
}

func deleteArmsDispatchRules(client *connectivity.AliyunClient, ids []string) {
	for _, id := range ids {
		if _, err := client.RpcPost("ARMS", "2019-08-08", "DeleteDispatchRule", nil, map[string]interface{}{
			"Id":       id,
			"RegionId": client.RegionId,
		}, false); err != nil {
			log.Printf("[WARN] failed to delete seeded dispatch rule %s: %s", id, err)
		}
	}
}

// TestAccAliCloudARMSDispatchRulesDataSource_nameSubstring distinguishes the two
// possible dispatch_rule_name semantics: the query is a proper substring of the
// created rule's name, so server-side substring matching returns the rule while a
// client-side exact comparison would return nothing.
func TestAccAliCloudARMSDispatchRulesDataSource_nameSubstring(t *testing.T) {
	rand := tfacctest.RandInt()
	name := fmt.Sprintf("tf-testacc-armsNameSub%v", rand)
	resourceId := "data.alicloud_arms_dispatch_rules.default"
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourceArmsDispatchRulesNameSubstringDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.ARMSSupportRegions)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				// The interpolation both computes a proper substring and ties the
				// data source to the rule so it is read after the rule exists.
				Config: testAccConfig(map[string]interface{}{
					"dispatch_rule_name": "${substr(alicloud_arms_dispatch_rule.default.dispatch_rule_name, 0, length(alicloud_arms_dispatch_rule.default.dispatch_rule_name) - 2)}",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "rules.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "rules.0.dispatch_rule_name", name),
				),
			},
		},
	})
}

func dataSourceArmsDispatchRulesNameSubstringDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
resource "alicloud_arms_alert_contact" "default" {
  alert_contact_name = var.name
  email              = "${var.name}@example.com"
}
resource "alicloud_arms_dispatch_rule" "default" {
  dispatch_rule_name = var.name
  dispatch_type      = "CREATE_ALERT"
  group_rules {
    group_wait_time = 5
    group_interval  = 15
    repeat_interval = 61
    grouping_fields = ["alertname"]
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
    notify_channels   = ["dingTalk"]
    notify_start_time = "00:00"
    notify_end_time   = "23:59"
  }
}
`, name)
}
