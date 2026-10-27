package nlb_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudNlbServerGroupServerAttachmentsDataSource(t *testing.T) {
	rand := tfacctest.RandInt()
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudNlbServerGroupServerAttachmentsDataSourceName(rand, map[string]string{
			"ids": `["${local.attachment_id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudNlbServerGroupServerAttachmentsDataSourceName(rand, map[string]string{
			"ids": `["${local.attachment_id}_fake"]`,
		}),
	}
	serverGroupIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudNlbServerGroupServerAttachmentsDataSourceName(rand, map[string]string{
			"ids":             `["${local.attachment_id}"]`,
			"server_group_id": `"${alicloud_nlb_server_group_server_attachment.default.server_group_id}"`,
		}),
		FakeConfig: "",
	}
	serverIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudNlbServerGroupServerAttachmentsDataSourceName(rand, map[string]string{
			"server_ids": `["${alicloud_nlb_server_group_server_attachment.default.server_id}"]`,
			"ids":        `["${local.attachment_id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudNlbServerGroupServerAttachmentsDataSourceName(rand, map[string]string{
			"server_ids": `["${alicloud_nlb_server_group_server_attachment.default.server_id}_fake"]`,
			"ids":        `["${local.attachment_id}"]`,
		}),
	}
	serverIpConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudNlbServerGroupServerAttachmentsDataSourceName(rand, map[string]string{
			"server_ips": `["${alicloud_nlb_server_group_server_attachment.default.server_ip}"]`,
			"ids":        `["${local.attachment_id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudNlbServerGroupServerAttachmentsDataSourceName(rand, map[string]string{
			"ids":        `["${local.attachment_id}"]`,
			"server_ips": `["${alicloud_nlb_server_group_server_attachment.default.server_ip}_fake"]`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudNlbServerGroupServerAttachmentsDataSourceName(rand, map[string]string{
			"server_ids":      `["${alicloud_nlb_server_group_server_attachment.default.server_id}"]`,
			"server_ips":      `["${alicloud_nlb_server_group_server_attachment.default.server_ip}"]`,
			"ids":             `["${local.attachment_id}"]`,
			"server_group_id": `"${alicloud_nlb_server_group_server_attachment.default.server_group_id}"`,
		}),
		FakeConfig: testAccCheckAlicloudNlbServerGroupServerAttachmentsDataSourceName(rand, map[string]string{
			"ids":        `["${local.attachment_id}_fake"]`,
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
			"attachments.0.server_group_id": acctest.CHECKSET,
			"attachments.0.id":              acctest.CHECKSET,
			"attachments.0.server_ip":       "10.0.0.0",
			"attachments.0.server_id":       "10.0.0.0",
			"attachments.0.server_type":     "Ip",
			"attachments.0.weight":          "100",
			"attachments.0.zone_id":         "",
			"attachments.0.status":          acctest.CHECKSET,
		}
	}
	var fakeAlicloudNlbServerGroupServerAttachmentsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#": "0",
		}
	}
	var AlicloudNlbServerGroupServerAttachmentsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_nlb_server_group_server_attachments.default",
		ExistMapFunc: existAlicloudNlbServerGroupServerAttachmentsDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudNlbServerGroupServerAttachmentsDataSourceNameMapFunc,
	}

	preCheck := func() {
		acctest.PreCheck(t)
	}
	AlicloudNlbServerGroupServerAttachmentsCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, serverGroupIdConf, serverIdConf, serverIpConf, allConf)
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
    name_regex = "^default-NODELETING$"
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

locals {
  attachment_id = "${alicloud_nlb_server_group_server_attachment.default.server_group_id}:${alicloud_nlb_server_group_server_attachment.default.server_id}:${alicloud_nlb_server_group_server_attachment.default.server_type}:${alicloud_nlb_server_group_server_attachment.default.port}"
}

data "alicloud_nlb_server_group_server_attachments" "default" {
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
