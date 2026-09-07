package alicloud

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudAlidnsAccessStrategiesDataSource(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	checkoutSupportedRegions(t, true, connectivity.AlidnsSupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAlidnsAccessStrategiesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_alidns_access_strategy.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudAlidnsAccessStrategiesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_alidns_access_strategy.default.id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAlidnsAccessStrategiesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_alidns_access_strategy.default.strategy_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudAlidnsAccessStrategiesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_alidns_access_strategy.default.strategy_name}_fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAlidnsAccessStrategiesDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_alidns_access_strategy.default.id}"]`,
			"name_regex": `"${alicloud_alidns_access_strategy.default.strategy_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudAlidnsAccessStrategiesDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_alidns_access_strategy.default.id}"]`,
			"name_regex": `"${alicloud_alidns_access_strategy.default.strategy_name}_fake"`,
		}),
	}
	var existAlicloudAlidnsAccessStrategiesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                           "1",
			"names.#":                                         "1",
			"strategies.#":                                    "1",
			"strategies.0.id":                                 CHECKSET,
			"strategies.0.access_strategy_id":                 CHECKSET,
			"strategies.0.create_time":                        CHECKSET,
			"strategies.0.create_timestamp":                   CHECKSET,
			"strategies.0.default_addr_pool_type":             "IPV4",
			"strategies.0.default_addr_pools.#":               "1",
			"strategies.0.default_addr_pools.0.addr_count":    CHECKSET,
			"strategies.0.default_addr_pools.0.addr_pool_id":  CHECKSET,
			"strategies.0.default_addr_pools.0.lba_weight":    "1",
			"strategies.0.default_addr_pools.0.name":          CHECKSET,
			"strategies.0.default_available_addr_num":         CHECKSET,
			"strategies.0.default_latency_optimization":       "",
			"strategies.0.default_lba_strategy":               "RATIO",
			"strategies.0.default_max_return_addr_num":        "0",
			"strategies.0.default_min_available_addr_num":     "1",
			"strategies.0.effective_addr_pool_group_type":     CHECKSET,
			"strategies.0.failover_addr_pool_type":            "IPV4",
			"strategies.0.failover_addr_pools.#":              "1",
			"strategies.0.failover_addr_pools.0.addr_count":   CHECKSET,
			"strategies.0.failover_addr_pools.0.addr_pool_id": CHECKSET,
			"strategies.0.failover_addr_pools.0.lba_weight":   "1",
			"strategies.0.failover_addr_pools.0.name":         CHECKSET,
			"strategies.0.failover_available_addr_num":        CHECKSET,
			"strategies.0.failover_latency_optimization":      "",
			"strategies.0.failover_lba_strategy":              "RATIO",
			"strategies.0.failover_max_return_addr_num":       "0",
			"strategies.0.failover_min_available_addr_num":    "1",
			"strategies.0.instance_id":                        CHECKSET,
			"strategies.0.access_mode":                        "AUTO",
			"strategies.0.lines.#":                            "1",
			"strategies.0.lines.0.group_code":                 CHECKSET,
			"strategies.0.lines.0.group_name":                 CHECKSET,
			"strategies.0.lines.0.line_code":                  "default",
			"strategies.0.lines.0.line_name":                  CHECKSET,
			"strategies.0.strategy_mode":                      "GEO",
			"strategies.0.strategy_name":                      fmt.Sprintf("tf-testAcc-%d", rand),
		}
	}
	var fakeAlicloudAlidnsAccessStrategiesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudAlidnsAccessStrategiesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_alidns_access_strategies.default",
		existMapFunc: existAlicloudAlidnsAccessStrategiesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudAlidnsAccessStrategiesDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
		testAccPreCheckWithTime(t, []int{1})
		testAccPreCheckWithEnvVariable(t, "ALICLOUD_ICP_DOMAIN_NAME")
	}
	alicloudAlidnsAccessStrategiesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, allConf)
}
func testAccCheckAlicloudAlidnsAccessStrategiesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAcc-%d"
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

resource "alicloud_alidns_address_pool" "ipv4" {
  count             = 2
  address_pool_name = var.name
  instance_id       = local.gtm_instance_id
  lba_strategy      = "RATIO"
  type              = "IPV4"
  address {
    attribute_info  = "{\"lineCodeRectifyType\":\"RECTIFIED\",\"lineCodes\":[\"os_namerica_us\"]}"
    remark          = "address_remark"
    address         = "1.1.1.1"
    mode            = "SMART"
    lba_weight      = 1
  }
}

resource "alicloud_alidns_access_strategy" "default" {
  strategy_name                   = var.name
  instance_id                     = local.gtm_instance_id
  default_addr_pool_type          = "IPV4"
  failover_lba_strategy           = "RATIO"
  failover_min_available_addr_num = 1
  strategy_mode                   = "GEO"
  default_lba_strategy            = "RATIO"
  failover_addr_pool_type         = "IPV4"
  failover_addr_pools {
    lba_weight   = 1
    addr_pool_id = alicloud_alidns_address_pool.ipv4.0.id
  }
  lines {
    line_code = "default"
  }
  default_addr_pools {
    lba_weight   = 1
    addr_pool_id = alicloud_alidns_address_pool.ipv4.1.id
  }
  default_min_available_addr_num  = 1
}

data "alicloud_alidns_access_strategies" "default" {	
	enable_details = true
	instance_id =  local.gtm_instance_id
	strategy_mode = "GEO"
	%s	
}
`, rand, os.Getenv("ALICLOUD_ICP_DOMAIN_NAME"), strings.Join(pairs, " \n "))
	return config
}
