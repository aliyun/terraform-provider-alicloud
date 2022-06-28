package alicloud

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudAlidnsGtmInstancesDataSource(t *testing.T) {
	rand := acctest.RandInt()
	checkoutSupportedRegions(t, true, connectivity.AlidnsSupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAlidnsGtmInstancesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_alidns_gtm_instance.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudAlidnsGtmInstancesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_alidns_gtm_instance.default.id}_fake"]`,
		}),
	}
	resourceGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAlidnsGtmInstancesDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_alidns_gtm_instance.default.id}"]`,
			"resource_group_id": `"${alicloud_alidns_gtm_instance.default.resource_group_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudAlidnsGtmInstancesDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_alidns_gtm_instance.default.id}"]`,
			"resource_group_id": `"${alicloud_alidns_gtm_instance.default.resource_group_id}_fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAlidnsGtmInstancesDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_alidns_gtm_instance.default.id}"]`,
			"resource_group_id": `"${alicloud_alidns_gtm_instance.default.resource_group_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudAlidnsGtmInstancesDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_alidns_gtm_instance.default.id}_fake"]`,
			"resource_group_id": `"${alicloud_alidns_gtm_instance.default.resource_group_id}_fake"`,
		}),
	}
	var existAlicloudAlidnsGtmInstancesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                               "1",
			"instances.#":                         "1",
			"instances.0.alert_group.#":           "1",
			"instances.0.alert_config.#":          "1",
			"instances.0.instance_name":           fmt.Sprintf("tf-testAccGtmInstance-%d", rand),
			"instances.0.resource_group_id":       CHECKSET,
			"instances.0.ttl":                     "60",
			"instances.0.cname_type":              CHECKSET,
			"instances.0.create_time":             CHECKSET,
			"instances.0.expire_time":             CHECKSET,
			"instances.0.id":                      CHECKSET,
			"instances.0.instance_id":             CHECKSET,
			"instances.0.strategy_mode":           CHECKSET,
			"instances.0.payment_type":            "Subscription",
			"instances.0.public_cname_mode":       CHECKSET,
			"instances.0.public_rr":               CHECKSET,
			"instances.0.public_user_domain_name": CHECKSET,
			"instances.0.public_zone_name":        CHECKSET,
			"instances.0.package_edition":         "ultimate",
		}
	}
	var fakeAlicloudAlidnsGtmInstancesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"instances.#": "0",
		}
	}
	var alicloudAlidnsGtmInstancesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_alidns_gtm_instances.default",
		existMapFunc: existAlicloudAlidnsGtmInstancesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudAlidnsGtmInstancesDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
		testAccPreCheckWithTime(t, []int{1})
		testAccPreCheckWithEnvVariable(t, "ALICLOUD_ICP_DOMAIN_NAME")
	}
	alicloudAlidnsGtmInstancesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, resourceGroupIdConf, allConf)
}
func testAccCheckAlicloudAlidnsGtmInstancesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccGtmInstance-%d"
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

data "alicloud_alidns_gtm_instances" "default" {	
	%s
}
`, rand, os.Getenv("ALICLOUD_ICP_DOMAIN_NAME"), strings.Join(pairs, " \n "))
	return config
}
