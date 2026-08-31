package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test CloudFirewall UserAlarmConfig. >>> Resource test cases, automatically generated.
// Case 告警测试_Notify 12512
// lintignore: AT001
func TestAccAliCloudCloudFirewallUserAlarmConfig_basic12512(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_user_alarm_config.default"
	ra := resourceAttrInit(resourceId, AlicloudCloudFirewallUserAlarmConfigMap12512)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudFirewallServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudFirewallUserAlarmConfig")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudfirewall%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudFirewallUserAlarmConfigBasicDependence12512)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
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
						"use_default_contact":        CHECKSET,
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
					"notify_config": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alarm_config.#":              "20",
						"alarm_config.0.alarm_value":  "off",
						"alarm_config.0.alarm_type":   "bandwidth",
						"alarm_config.0.alarm_period": "1",
						"use_default_contact":         CHECKSET,
						"notify_config.#":             REMOVEKEY,
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
	ra := resourceAttrInit(resourceId, AlicloudCloudFirewallUserAlarmConfigMap12506)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudFirewallServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudFirewallUserAlarmConfig")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudfirewall%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudFirewallUserAlarmConfigBasicDependence12506)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
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
						"use_default_contact": CHECKSET,
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

// Test CloudFirewall UserAlarmConfig. <<< Resource test cases, automatically generated.
