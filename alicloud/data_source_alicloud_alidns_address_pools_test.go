package alicloud

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudAlidnsAddressPoolsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(100, 999)
	checkoutSupportedRegions(t, true, connectivity.AlidnsSupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAlidnsAddressPoolsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_alidns_address_pool.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudAlidnsAddressPoolsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_alidns_address_pool.default.id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAlidnsAddressPoolsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_alidns_address_pool.default.id}"]`,
			"name_regex": `"${alicloud_alidns_address_pool.default.address_pool_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudAlidnsAddressPoolsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_alidns_address_pool.default.id}"]`,
			"name_regex": `"${alicloud_alidns_address_pool.default.address_pool_name}_fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAlidnsAddressPoolsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_alidns_address_pool.default.id}"]`,
			"name_regex": `"${alicloud_alidns_address_pool.default.address_pool_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudAlidnsAddressPoolsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_alidns_address_pool.default.id}_fake"]`,
			"name_regex": `"${alicloud_alidns_address_pool.default.address_pool_name}_fake"`,
		}),
	}
	var existAlicloudAlidnsAddressPoolsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                            "1",
			"names.#":                          "1",
			"pools.#":                          "1",
			"pools.0.address_pool_name":        fmt.Sprintf("tf-testacc%d", rand),
			"pools.0.id":                       CHECKSET,
			"pools.0.instance_id":              CHECKSET,
			"pools.0.address_pool_id":          CHECKSET,
			"pools.0.create_time":              CHECKSET,
			"pools.0.create_timestamp":         CHECKSET,
			"pools.0.update_time":              CHECKSET,
			"pools.0.update_timestamp":         CHECKSET,
			"pools.0.lba_strategy":             "RATIO",
			"pools.0.monitor_status":           "UNCONFIGURED",
			"pools.0.monitor_config_id":        "",
			"pools.0.type":                     "IPV4",
			"pools.0.address.#":                "1",
			"pools.0.address.0.address":        "1.1.1.1",
			"pools.0.address.0.remark":         "address_remark",
			"pools.0.address.0.mode":           "SMART",
			"pools.0.address.0.lba_weight":     "1",
			"pools.0.address.0.attribute_info": CHECKSET,
		}
	}
	var fakeAlicloudAlidnsAddressPoolsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudAlidnsAddressPoolsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_alidns_address_pools.default",
		existMapFunc: existAlicloudAlidnsAddressPoolsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudAlidnsAddressPoolsDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
		testAccPreCheckWithTime(t, []int{1})
		testAccPreCheckWithEnvVariable(t, "ALICLOUD_ICP_DOMAIN_NAME")
	}
	alicloudAlidnsAddressPoolsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, allConf)
}
func testAccCheckAlicloudAlidnsAddressPoolsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
  default = "tf-testacc%d"
}

variable "domain_name" {
  default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_cms_alarm_contact_group" "default" {
  alarm_contact_group_name = var.name
}

data "alicloud_alidns_gtm_instances" "default" {}

resource "alicloud_alidns_gtm_instance" "default" {
  count                   = length(data.alicloud_alidns_gtm_instances.default.ids) > 0 ? 0 : 1
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

locals {
  gtm_instance_id = length(data.alicloud_alidns_gtm_instances.default.ids) > 0 ? data.alicloud_alidns_gtm_instances.default.ids[0] : concat(alicloud_alidns_gtm_instance.default.*.id, [""])[0]
}

resource "alicloud_alidns_address_pool" "default" {
  address_pool_name = var.name
  instance_id       = local.gtm_instance_id
  lba_strategy      = "RATIO"
  type              = "IPV4"
  address {
    attribute_info = "{\"lineCodeRectifyType\":\"RECTIFIED\",\"lineCodes\":[\"os_namerica_us\"]}"
    remark     = "address_remark"
    address    = "1.1.1.1"
    mode       = "SMART"
    lba_weight = 1
  }
}

data "alicloud_alidns_address_pools" "default" {
  enable_details = true
  instance_id    = local.gtm_instance_id
	%s
}
`, rand, os.Getenv("ALICLOUD_ICP_DOMAIN_NAME"), strings.Join(pairs, " \n "))
	return config
}
