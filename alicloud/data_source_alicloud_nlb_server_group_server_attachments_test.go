package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudNlbServerGroupServerAttachmentsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudNlbServerGroupServerAttachmentsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_nlb_server_group_server_attachment.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudNlbServerGroupServerAttachmentsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_nlb_server_group_server_attachment.default.id}_fake"]`,
		}),
	}
	serverGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudNlbServerGroupServerAttachmentsDataSourceName(rand, map[string]string{
			"ids":             `["${alicloud_nlb_server_group_server_attachment.default.id}"]`,
			"server_group_id": `"${alicloud_nlb_server_group_server_attachment.default.server_group_id}"`,
		}),
		fakeConfig: "",
	}
	serverIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudNlbServerGroupServerAttachmentsDataSourceName(rand, map[string]string{
			"server_ids": `["${alicloud_nlb_server_group_server_attachment.default.server_id}"]`,
			"ids":        `["${alicloud_nlb_server_group_server_attachment.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudNlbServerGroupServerAttachmentsDataSourceName(rand, map[string]string{
			"server_ids": `["${alicloud_nlb_server_group_server_attachment.default.server_id}_fake"]`,
			"ids":        `["${alicloud_nlb_server_group_server_attachment.default.id}"]`,
		}),
	}
	serverIpConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudNlbServerGroupServerAttachmentsDataSourceName(rand, map[string]string{
			"server_ips": `["${alicloud_nlb_server_group_server_attachment.default.server_ip}"]`,
			"ids":        `["${alicloud_nlb_server_group_server_attachment.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudNlbServerGroupServerAttachmentsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_nlb_server_group_server_attachment.default.id}"]`,
			"server_ips": `["${alicloud_nlb_server_group_server_attachment.default.server_ip}_fake"]`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudNlbServerGroupServerAttachmentsDataSourceName(rand, map[string]string{
			"server_ids":      `["${alicloud_nlb_server_group_server_attachment.default.server_id}"]`,
			"server_ips":      `["${alicloud_nlb_server_group_server_attachment.default.server_ip}"]`,
			"ids":             `["${alicloud_nlb_server_group_server_attachment.default.id}"]`,
			"server_group_id": `"${alicloud_nlb_server_group_server_attachment.default.server_group_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudNlbServerGroupServerAttachmentsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_nlb_server_group_server_attachment.default.id}_fake"]`,
			"server_ids": `["${alicloud_nlb_server_group_server_attachment.default.server_id}_fake"]`,
			"server_ips": `["${alicloud_nlb_server_group_server_attachment.default.server_ip}_fake"]`,
		}),
	}
	var existAlicloudNlbServerGroupServerAttachmentsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                         "1",
			"attachments.#":                 "1",
			"attachments.0.description":     fmt.Sprintf("tf-testAccServerGroupServerAttachment-%d", rand),
			"attachments.0.port":            "80",
			"attachments.0.server_group_id": CHECKSET,
			"attachments.0.id":              CHECKSET,
			"attachments.0.server_ip":       "10.0.0.0",
			"attachments.0.server_id":       "10.0.0.0",
			"attachments.0.server_type":     "Ip",
			"attachments.0.weight":          "100",
			"attachments.0.zone_id":         "",
			"attachments.0.status":          CHECKSET,
		}
	}
	var fakeAlicloudNlbServerGroupServerAttachmentsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#": "0",
		}
	}
	var AlicloudNlbServerGroupServerAttachmentsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_nlb_server_group_server_attachments.default",
		existMapFunc: existAlicloudNlbServerGroupServerAttachmentsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudNlbServerGroupServerAttachmentsDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	AlicloudNlbServerGroupServerAttachmentsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, serverGroupIdConf, serverIdConf, serverIpConf, allConf)
}
func testAccCheckAlicloudNlbServerGroupServerAttachmentsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {
	default = "tf-testAccServerGroupServerAttachment-%d"
}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_nlb_server_group" "default" {
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.0
  server_group_name = var.name
  server_group_type = "Ip"
  vpc_id            = data.alicloud_vpcs.default.ids.0
  scheduler         = "Wrr"
  protocol          = "TCP"
  health_check {
    health_check_enabled = false
  }
  address_ip_version = "Ipv4"
}

resource "alicloud_nlb_server_group_server_attachment" "default" {
  server_type     = "Ip"
  server_id       = "10.0.0.0"
  description     = var.name
  port            = 80
  server_group_id = alicloud_nlb_server_group.default.id
  weight          = 100
  server_ip       = "10.0.0.0"
}

data "alicloud_nlb_server_group_server_attachments" "default" {
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
