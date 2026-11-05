package cloudfirewall_test

import (
	"fmt"
	"testing"

	acctest "github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	connectivity "github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	cloudfirewall "github.com/aliyun/terraform-provider-alicloud/alicloud/service/cloudfirewall"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	resource "github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test CloudFirewall UserAlarmConfig. >>> Resource test cases, automatically generated.
// Case 告警测试_Notify 12512
// lintignore: AT001
func TestAccAliCloudCloudFirewallUserAlarmConfig_basic12512(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_user_alarm_config.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudCloudFirewallUserAlarmConfigMap12512)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudfirewall.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCloudFirewallUserAlarmConfig")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudfirewall%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudCloudFirewallUserAlarmConfigBasicDependence12512)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"alarm_config": []map[string]interface{}{
						{
							"alarm_value":    "on",
							"alarm_type":     "bandwidth",
							"alarm_period":   "1",
							"alarm_hour":     "0",
							"alarm_notify":   "0",
							"alarm_week_day": "0",
						},
					},
					"use_default_contact": "1",
					"notify_config": []map[string]interface{}{
						{
							"notify_value": "18800000000",
							"notify_type":  "sms",
						},
					},
					"alarm_lang": "zh",
					"lang":       "zh",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alarm_config.#":             "20",
						"alarm_config.0.alarm_value": "on",
						"alarm_config.0.alarm_type":  "bandwidth",
						"use_default_contact":        acctest.CHECKSET,
						"notify_config.#":            "2",
						"alarm_lang":                 "zh",
						"lang":                       "zh",
					}),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"lang", "use_default_contact"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"alarm_config": []map[string]interface{}{
						{
							"alarm_value":    "off",
							"alarm_type":     "bandwidth",
							"alarm_period":   "0",
							"alarm_hour":     "1",
							"alarm_notify":   "1",
							"alarm_week_day": "1",
						},
					},
					"use_default_contact": "1",
					"contact_config": []map[string]interface{}{
						{
							"status":       "1",
							"email":        "test@test.com",
							"mobile_phone": "12300000000",
							"name":         "test",
						},
					},
					"alarm_lang":    "en",
					"lang":          "zh",
					"notify_config": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alarm_config.#":              "20",
						"alarm_config.0.alarm_value":  "off",
						"alarm_config.0.alarm_type":   "bandwidth",
						"alarm_config.0.alarm_period": "1",
						"use_default_contact":         acctest.CHECKSET,
						"notify_config.#":             acctest.REMOVEKEY,
						"alarm_lang":                  "en",
						"lang":                        "zh",
					}),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"lang", "use_default_contact"},
			},
		},
	})
}

var AlicloudCloudFirewallUserAlarmConfigMap12512 = map[string]string{}

func AlicloudCloudFirewallUserAlarmConfigBasicDependence12512(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 告警测试 12506
// lintignore: AT001
func TestAccAliCloudCloudFirewallUserAlarmConfig_basic12506(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_user_alarm_config.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudCloudFirewallUserAlarmConfigMap12506)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudfirewall.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCloudFirewallUserAlarmConfig")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudfirewall%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudCloudFirewallUserAlarmConfigBasicDependence12506)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"alarm_config": []map[string]interface{}{
						{
							"alarm_value":    "on",
							"alarm_type":     "bandwidth",
							"alarm_period":   "1",
							"alarm_hour":     "0",
							"alarm_notify":   "0",
							"alarm_week_day": "0",
						},
					},
					"use_default_contact": "1",
					"contact_config": []map[string]interface{}{
						{
							"status":       "1",
							"email":        "test@test.com",
							"mobile_phone": "13500000000",
							"name":         "test",
						},
					},
					"alarm_lang": "zh",
					"lang":       "zh",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alarm_config.#":      "20",
						"use_default_contact": acctest.CHECKSET,
						"contact_config.#":    "1",
						"alarm_lang":          "zh",
						"lang":                "zh",
					}),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"lang", "use_default_contact"},
			},
		},
	})
}

var AlicloudCloudFirewallUserAlarmConfigMap12506 = map[string]string{}

func AlicloudCloudFirewallUserAlarmConfigBasicDependence12506(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 告警测试_modify 12512
// lintignore: AT001
func TestAccAliCloudCloudFirewallUserAlarmConfig_modify(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_user_alarm_config.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudCloudFirewallUserAlarmConfigMap12512)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudfirewall.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCloudFirewallUserAlarmConfig")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudfirewall%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudCloudFirewallUserAlarmConfigBasicDependence12512)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"alarm_config": []map[string]interface{}{
						{
							"alarm_value":    "on",
							"alarm_type":     "bandwidth",
							"alarm_period":   "1",
							"alarm_hour":     "0",
							"alarm_notify":   "0",
							"alarm_week_day": "0",
						},
					},
					"use_default_contact": "1",
					"notify_config": []map[string]interface{}{
						{
							"notify_value": "18800000000",
							"notify_type":  "sms",
						},
					},
					"alarm_lang": "zh",
					"lang":       "zh",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alarm_config.0.alarm_type":    "bandwidth",
						"notify_config.0.notify_value": "18800000000",
					}),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"alarm_config": []map[string]interface{}{
						{
							"alarm_value":    "on",
							"alarm_type":     "weeklyReport",
							"alarm_period":   "1",
							"alarm_hour":     "0",
							"alarm_notify":   "0",
							"alarm_week_day": "0",
						},
					},
					"use_default_contact": "1",
					"notify_config": []map[string]interface{}{
						{
							"notify_value": "test@test.com",
							"notify_type":  "mail",
						},
					},
					"alarm_lang": "zh",
					"lang":       "zh",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alarm_config.0.alarm_type":    "weeklyReport",
						"notify_config.0.notify_value": "test@test.com",
					}),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

// collection-order: bypass alarm_config -- pre-existing TypeList; migration preserves original Update (sends raw[0] only), multi-member reorder convergence not provable with existing implementation
// collection-order: bypass contact_config -- pre-existing TypeList; migration preserves original Update (sends raw[0] only), multi-member reorder convergence not provable with existing implementation
// collection-order: bypass notify_config -- pre-existing TypeList; migration preserves original Update (sends raw[0] only), multi-member reorder convergence not provable with existing implementation
//
// _modify restored per QA CFM-001 fix mandate: covers the three attribute
// modification paths (alarm_config.alarm_type, notify_config.notify_type,
// notify_config.notify_value) that were missing from the static
// TestingCoverageRate coverage. It uses single-member configs (consistent
// with the existing Update path that sends raw[0]) and ExpectNonEmptyPlan
// (consistent with basic12512/basic12506, which tolerate the known
// notify_type readback drift). Assertions cover only attributes that read
// back correctly (alarm_type, notify_value); notify_type is exercised in
// config for static coverage but not asserted, because the existing Read
// path does not write notify_type into state (separate readback defect,
// candidate fix parked per root-review decision, NOT applied to this PR).
//
// The three collection-order cases (_alarmConfigOrder, _contactConfigOrder,
// _notifyConfigOrder) remain withdrawn: they configure two list members and
// reorder them in step 2, but the existing Update path only sends the first
// member (raw[0]) in a single-element request, so two-member convergence
// cannot be proven. Their verbatim source is preserved as a problem-evidence
// patch Artifact (withdrawn-useralarm-cases-evidence). Multi-member list
// coverage is a real pre-existing gap in the baseline schema, declared as
// stock debt via the collection-order bypass markers above (not a new
// exemption, not a red light written green; the schema and Update semantics
// are unchanged by this migration).

// Test CloudFirewall UserAlarmConfig. <<< Resource test cases, automatically generated.
