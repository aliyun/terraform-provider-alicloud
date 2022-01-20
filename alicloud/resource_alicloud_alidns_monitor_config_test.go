package alicloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudAlidnsMonitorConfig_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alidns_monitor_config.default"
	checkoutSupportedRegions(t, true, connectivity.AlidnsSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudAlidnsMonitorConfigMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AlidnsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAlidnsMonitorConfig")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAlidnsMonitorConfigBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithTime(t, []int{1})
			testAccPreCheckWithEnvVariable(t, "ALICLOUD_ICP_DOMAIN_NAME")
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"addr_pool_id":        "${alicloud_alidns_address_pool.default.id}",
					"evaluation_count":    "1",
					"interval":            "60",
					"timeout":             "5000",
					"protocol_type":       "TCP",
					"monitor_extend_info": `{\"failureRate\":50,\"port\":80}`,
					"isp_city_node": []map[string]interface{}{
						{
							"city_code": "503",
							"isp_code":  "465",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"addr_pool_id":        CHECKSET,
						"protocol_type":       "TCP",
						"evaluation_count":    "1",
						"interval":            "60",
						"timeout":             "5000",
						"monitor_extend_info": CHECKSET,
						"isp_city_node.#":     "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"evaluation_count": "3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"evaluation_count": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"timeout": "10000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"timeout": "10000",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"protocol_type":       "PING",
					"monitor_extend_info": `{\"packetNum\":20,\"packetLossRate\":10,\"failureRate\":50}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocol_type":       "PING",
						"monitor_extend_info": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"isp_city_node": []map[string]interface{}{
						{
							"city_code": "569",
							"isp_code":  "465",
						},
						{
							"city_code": "491",
							"isp_code":  "232",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"isp_city_node.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"protocol_type":       "TCP",
					"monitor_extend_info": `{\"failureRate\":50,\"port\":80}`,
					"evaluation_count":    "1",
					"timeout":             "5000",
					"isp_city_node": []map[string]interface{}{
						{
							"city_code": "503",
							"isp_code":  "465",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocol_type":       "TCP",
						"evaluation_count":    "1",
						"timeout":             "5000",
						"monitor_extend_info": CHECKSET,
						"isp_city_node.#":     "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"addr_pool_id"},
			},
		},
	})
}

var AlicloudAlidnsMonitorConfigMap0 = map[string]string{}

func AlicloudAlidnsMonitorConfigBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

variable "domain_name" {
  default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_cms_alarm_contact_group" "default" {
  alarm_contact_group_name = var.name
}

resource "alicloud_alidns_gtm_instance" "default" {
  instance_name           = var.name
  payment_type            = "Subscription"
  period                  = 1
  renewal_status          = "ManualRenewal"
  package_edition         = "ultimate"
  health_check_task_count = 100
  sms_notification_count  = 1000
  public_cname_mode       = "SYSTEM_ASSIGN"
  ttl                     = 60
  cname_type              = "PUBLIC"
  resource_group_id       = data.alicloud_resource_manager_resource_groups.default.groups.0.id
  alert_group             = [alicloud_cms_alarm_contact_group.default.alarm_contact_group_name]
  public_user_domain_name = var.domain_name
  alert_config {
    sms_notice      = true
    notice_type     = "ADDR_ALERT"
    email_notice    = true
    dingtalk_notice = true
  }
}

resource "alicloud_alidns_address_pool" "default" {
  address_pool_name = var.name
  instance_id       = alicloud_alidns_gtm_instance.default.id
  lba_strategy      = "RATIO"
  type              = "IPV4"
  address {
    attribute_info = "{\"lineCodeRectifyType\":\"RECTIFIED\",\"lineCodes\":[\"os_namerica_us\"]}"
    remark         = "address_remark"
    address        = "1.1.1.1"
    mode           = "SMART"
    lba_weight     = 1
  }
}
`, name, os.Getenv("ALICLOUD_ICP_DOMAIN_NAME"))
}
