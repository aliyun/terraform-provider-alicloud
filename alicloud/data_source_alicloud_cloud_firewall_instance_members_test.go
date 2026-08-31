package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudCloudFirewallInstanceMemberDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudFirewallInstanceMemberSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_cloud_firewall_instance_member.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudCloudFirewallInstanceMemberSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_cloud_firewall_instance_member.default.id}_fake"]`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudFirewallInstanceMemberSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_cloud_firewall_instance_member.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudCloudFirewallInstanceMemberSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_cloud_firewall_instance_member.default.id}_fake"]`,
		}),
	}

	preCheck := func() {
		testAccPreCheck(t)
		// currently, international test account has not enabled RD
		testAccPreCheckWithAccountSiteType(t, DomesticSite)
	}

	CloudFirewallInstanceMemberCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, allConf)
}

var existCloudFirewallInstanceMemberMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"members.#":    "1",
		"members.0.id": CHECKSET,
	}
}

var fakeCloudFirewallInstanceMemberMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"members.#": "0",
	}
}

var CloudFirewallInstanceMemberCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_cloud_firewall_instance_members.default",
	existMapFunc: existCloudFirewallInstanceMemberMapFunc,
	fakeMapFunc:  fakeCloudFirewallInstanceMemberMapFunc,
}

func testAccCheckAlicloudCloudFirewallInstanceMemberSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccCFInstanceMember%d"
}

resource "alicloud_resource_manager_account" "default" {
  display_name = var.name
  abandon_able_check_id = ["SP_fc_fc"]
}

resource "alicloud_cloud_firewall_instance_member" "default" {
  member_desc = var.name
  member_uid  = alicloud_resource_manager_account.default.id
}

data "alicloud_cloud_firewall_instance_members" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
