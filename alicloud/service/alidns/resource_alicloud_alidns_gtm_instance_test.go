package alidns_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/alidns"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudAlidnsGtmInstance_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alidns_gtm_instance.default"
	acctest.CheckoutSupportedRegions(t, true, connectivity.AlidnsSupportRegions)
	ra := acctest.ResourceAttrInit(resourceId, AlicloudAlidnsGtmInstanceMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return alidns.NewAlidnsService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeAlidnsGtmInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%salidnsgtminstance%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudAlidnsGtmInstanceBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithEnvVariable(t, "ALICLOUD_ICP_DOMAIN_NAME")
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"force_update":            "true",
					"lang":                    "Subscription",
					"public_rr":               "60",
					"public_zone_name":        name,
					"renew_period":            "1",
					"instance_name":           name,
					"payment_type":            "Subscription",
					"period":                  "1",
					"renewal_status":          "ManualRenewal",
					"package_edition":         "ultimate",
					"health_check_task_count": "100",
					"sms_notification_count":  "1000",
					"public_cname_mode":       "SYSTEM_ASSIGN",
					"ttl":                     "60",
					"cname_type":              "PUBLIC",
					"resource_group_id":       "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
					"alert_group":             []string{"${alicloud_cms_alarm_contact_group.default.0.alarm_contact_group_name}"},
					"public_user_domain_name": "${var.domain_name}",
					"strategy_mode":           "GEO",
					"alert_config": []map[string]interface{}{
						{
							"sms_notice":      "true",
							"notice_type":     "ADDR_ALERT",
							"email_notice":    "true",
							"dingtalk_notice": "true",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":           name,
						"cname_type":              "PUBLIC",
						"ttl":                     "60",
						"alert_group.#":           "1",
						"alert_config.#":          "1",
						"resource_group_id":       acctest.CHECKSET,
						"public_cname_mode":       "SYSTEM_ASSIGN",
						"strategy_mode":           "GEO",
						"public_user_domain_name": acctest.CHECKSET,
						"public_rr":               acctest.CHECKSET,
						"public_zone_name":        acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cname_type":              "PUBLIC-upd",
					"public_rr":               "61",
					"public_user_domain_name": "60",
					"public_zone_name":        "60",
					"alert_config": []map[string]interface{}{
						{
							"dingtalk_notice": "false",
							"email_notice":    "false",
							"sms_notice":      "false",
						},
					},
					"instance_name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ttl": "300",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ttl": "300",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"strategy_mode": "LATENCY",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"strategy_mode": "LATENCY",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"alert_group": []string{"${alicloud_cms_alarm_contact_group.default.0.alarm_contact_group_name}", "${alicloud_cms_alarm_contact_group.default.1.alarm_contact_group_name}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alert_group.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"alert_config": []map[string]interface{}{
						{
							"sms_notice":      "true",
							"notice_type":     "ADDR_RESUME",
							"email_notice":    "true",
							"dingtalk_notice": "true",
						},
						{
							"sms_notice":      "true",
							"notice_type":     "ADDR_POOL_GROUP_UNAVAILABLE",
							"email_notice":    "true",
							"dingtalk_notice": "true",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alert_config.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name":     name,
					"public_cname_mode": "CUSTOM",
					"ttl":               "60",
					"alert_group":       []string{"${alicloud_cms_alarm_contact_group.default.0.alarm_contact_group_name}"},
					"strategy_mode":     "GEO",
					"alert_config": []map[string]interface{}{
						{
							"sms_notice":      "true",
							"notice_type":     "ADDR_ALERT",
							"email_notice":    "true",
							"dingtalk_notice": "true",
						},
					},
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":     name,
						"ttl":               "60",
						"alert_config.#":    "1",
						"alert_group.#":     "1",
						"strategy_mode":     "GEO",
						"resource_group_id": acctest.CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"lang", "force_update", "health_check_task_count", "sms_notification_count", "period", "renewal_status"},
			},
		},
	})
}

var AlicloudAlidnsGtmInstanceMap0 = map[string]string{}

func AlicloudAlidnsGtmInstanceBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
variable "domain_name" {
  default = "%s"
}
data "alicloud_resource_manager_resource_groups" "default" {}
resource "alicloud_cms_alarm_contact_group" "default" {
  count                    = 2
  alarm_contact_group_name = join("-", [var.name, count.index])
}
`, name, os.Getenv("ALICLOUD_ICP_DOMAIN_NAME"))
}

// lintignore: R001
