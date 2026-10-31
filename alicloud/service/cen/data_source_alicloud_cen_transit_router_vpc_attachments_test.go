package cen_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudCenTransitRouterVpcAttachmentsDataSource(t *testing.T) {
	rand := tfacctest.RandInt()

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudCenTransitRouterVpcAttachmentsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cen_transit_router_vpc_attachment.default.id}"]`,
		}),
		FakeConfig: testAccCheckAliCloudCenTransitRouterVpcAttachmentsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cen_transit_router_vpc_attachment.default.id}_fake"]`,
		}),
	}
	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudCenTransitRouterVpcAttachmentsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_cen_transit_router_vpc_attachment.default.transit_router_attachment_name}"`,
		}),
		FakeConfig: testAccCheckAliCloudCenTransitRouterVpcAttachmentsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_cen_transit_router_vpc_attachment.default.transit_router_attachment_name}_fake"`,
		}),
	}
	vpcIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudCenTransitRouterVpcAttachmentsDataSourceName(rand, map[string]string{
			"vpc_id": `"${alicloud_cen_transit_router_vpc_attachment.default.vpc_id}"`,
		}),
		FakeConfig: testAccCheckAliCloudCenTransitRouterVpcAttachmentsDataSourceName(rand, map[string]string{
			"vpc_id": `"${alicloud_cen_transit_router_vpc_attachment.default.vpc_id}_fake"`,
		}),
	}
	transitRouterIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudCenTransitRouterVpcAttachmentsDataSourceName(rand, map[string]string{
			"transit_router_id": `"${alicloud_cen_transit_router_vpc_attachment.default.transit_router_id}"`,
		}),
		FakeConfig: testAccCheckAliCloudCenTransitRouterVpcAttachmentsDataSourceName(rand, map[string]string{
			"transit_router_id": `"${alicloud_cen_transit_router_vpc_attachment.default.transit_router_id}"`,
			"status":            `"Attaching"`,
		}),
	}
	transitRouterAttachmentIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudCenTransitRouterVpcAttachmentsDataSourceName(rand, map[string]string{
			"transit_router_attachment_id": `"${alicloud_cen_transit_router_vpc_attachment.default.transit_router_attachment_id}"`,
		}),
		FakeConfig: testAccCheckAliCloudCenTransitRouterVpcAttachmentsDataSourceName(rand, map[string]string{
			"transit_router_attachment_id": `"${alicloud_cen_transit_router_vpc_attachment.default.transit_router_attachment_id}_fake"`,
		}),
	}
	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudCenTransitRouterVpcAttachmentsDataSourceName(rand, map[string]string{
			"status": `"Attached"`,
		}),
		FakeConfig: testAccCheckAliCloudCenTransitRouterVpcAttachmentsDataSourceName(rand, map[string]string{
			"status": `"Attaching"`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudCenTransitRouterVpcAttachmentsDataSourceName(rand, map[string]string{
			"ids":                          `["${alicloud_cen_transit_router_vpc_attachment.default.id}"]`,
			"name_regex":                   `"${alicloud_cen_transit_router_vpc_attachment.default.transit_router_attachment_name}"`,
			"vpc_id":                       `"${alicloud_cen_transit_router_vpc_attachment.default.vpc_id}"`,
			"transit_router_id":            `"${alicloud_cen_transit_router_vpc_attachment.default.transit_router_id}"`,
			"transit_router_attachment_id": `"${alicloud_cen_transit_router_vpc_attachment.default.transit_router_attachment_id}"`,
			"status":                       `"Attached"`,
		}),
		FakeConfig: testAccCheckAliCloudCenTransitRouterVpcAttachmentsDataSourceName(rand, map[string]string{
			"ids":                          `["${alicloud_cen_transit_router_vpc_attachment.default.id}_fake"]`,
			"name_regex":                   `"${alicloud_cen_transit_router_vpc_attachment.default.transit_router_attachment_name}_fake"`,
			"vpc_id":                       `"${alicloud_cen_transit_router_vpc_attachment.default.vpc_id}_fake"`,
			"transit_router_id":            `"${alicloud_cen_transit_router_vpc_attachment.default.transit_router_id}"`,
			"transit_router_attachment_id": `"${alicloud_cen_transit_router_vpc_attachment.default.transit_router_attachment_id}_fake"`,
			"status":                       `"Attaching"`,
		}),
	}

	var existAliCloudCenTransitRouterVpcAttachmentsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                "1",
			"names.#":              "1",
			"attachments.#":        "1",
			"attachments.0.id":     acctest.CHECKSET,
			"attachments.0.cen_id": acctest.CHECKSET,
			"attachments.0.transit_router_attachment_id":          acctest.CHECKSET,
			"attachments.0.vpc_id":                                acctest.CHECKSET,
			"attachments.0.transit_router_id":                     acctest.CHECKSET,
			"attachments.0.resource_type":                         "VPC",
			"attachments.0.payment_type":                          "PayAsYouGo",
			"attachments.0.vpc_owner_id":                          acctest.CHECKSET,
			"attachments.0.auto_publish_route_enabled":            acctest.CHECKSET,
			"attachments.0.transit_router_attachment_name":        acctest.CHECKSET,
			"attachments.0.transit_router_attachment_description": acctest.CHECKSET,
			"attachments.0.status":                                acctest.CHECKSET,
			"attachments.0.zone_mappings.#":                       "2",
			"attachments.0.zone_mappings.0.vswitch_id":            acctest.CHECKSET,
			"attachments.0.zone_mappings.0.zone_id":               acctest.CHECKSET,
			"attachments.0.zone_mappings.1.vswitch_id":            acctest.CHECKSET,
			"attachments.0.zone_mappings.1.zone_id":               acctest.CHECKSET,
		}
	}

	var fakeAliCloudCenTransitRouterVpcAttachmentsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":         "0",
			"names.#":       "0",
			"attachments.#": "0",
		}
	}

	var alicloudCenTransitRouterVpcAttachmentsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_cen_transit_router_vpc_attachments.default",
		ExistMapFunc: existAliCloudCenTransitRouterVpcAttachmentsDataSourceNameMapFunc,
		FakeMapFunc:  fakeAliCloudCenTransitRouterVpcAttachmentsDataSourceNameMapFunc,
	}

	preCheck := func() {
		acctest.PreCheckWithRegions(t, true, []connectivity.Region{connectivity.Hangzhou})
		acctest.PreCheck(t)
	}

	alicloudCenTransitRouterVpcAttachmentsCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, vpcIdConf, transitRouterIdConf, transitRouterAttachmentIdConf, statusConf, allConf)
}

func testAccCheckAliCloudCenTransitRouterVpcAttachmentsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {
  default = "tf-testAcc-CenTransitRouterVpcAttachment-%d"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "192.168.1.0/24"
  zone_id      = "cn-hangzhou-h"
  vswitch_name = var.name
}

resource "alicloud_vswitch" "default_master" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "192.168.2.0/24"
  zone_id      = "cn-hangzhou-i"
  vswitch_name = "${var.name}-master"
}

resource "alicloud_cen_instance" "default" {
  cen_instance_name = var.name
  protection_level  = "REDUCED"
}

resource "alicloud_cen_transit_router" "default" {
  cen_id = alicloud_cen_instance.default.id
}

resource "alicloud_cen_transit_router_vpc_attachment" "default" {
  cen_id                                = alicloud_cen_instance.default.id
  vpc_id                                = alicloud_vpc.default.id
  transit_router_id                     = alicloud_cen_transit_router.default.transit_router_id
  transit_router_attachment_name        = var.name
  transit_router_attachment_description = var.name
  zone_mappings {
    vswitch_id = alicloud_vswitch.default_master.id
    zone_id    = alicloud_vswitch.default_master.zone_id
  }
  zone_mappings {
    vswitch_id = alicloud_vswitch.default.id
    zone_id    = alicloud_vswitch.default.zone_id
  }
}

data "alicloud_cen_transit_router_vpc_attachments" "default" {
  cen_id = alicloud_cen_instance.default.id
  %s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}

func TestAccAliCloudCenTransitRouterVpcAttachmentsDataSource_options(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)

	optionsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudCenTransitRouterVpcAttachmentsDataSourceOptions(rand),
	}

	existMapFunc := func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                "1",
			"names.#":                              "1",
			"attachments.#":                        "1",
			"attachments.0.options.#":              "1",
			"attachments.0.options.0.ipv6_support": "enable",
			"attachments.0.options.0.appliance_mode_support":      "enable",
			"attachments.0.transit_router_attachment_id":          acctest.CHECKSET,
			"attachments.0.transit_router_attachment_name":        acctest.CHECKSET,
			"attachments.0.transit_router_attachment_description": acctest.CHECKSET,
		}
	}

	fakeMapFunc := func(rand int) map[string]string {
		return map[string]string{
			"ids.#":         "0",
			"names.#":       "0",
			"attachments.#": "0",
		}
	}

	checkInfo := acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_cen_transit_router_vpc_attachments.default",
		ExistMapFunc: existMapFunc,
		FakeMapFunc:  fakeMapFunc,
	}

	preCheck := func() {
		acctest.PreCheckWithRegions(t, true, []connectivity.Region{connectivity.Hangzhou})
		acctest.PreCheck(t)
	}

	checkInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, optionsConf)
}

func testAccCheckAliCloudCenTransitRouterVpcAttachmentsDataSourceOptions(rand int) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testacc-cen-tr-vpc-attachment-options-%d"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_vpc" "default" {
  cidr_block  = "192.168.0.0/16"
  enable_ipv6 = true
  ipv6_isp    = "BGP"
  vpc_name    = var.name
}

resource "alicloud_vswitch" "default" {
  vpc_id               = alicloud_vpc.default.id
  cidr_block           = "192.168.3.0/24"
  zone_id              = "cn-hangzhou-h"
  vswitch_name         = var.name
  ipv6_cidr_block_mask = "3"
}

resource "alicloud_vswitch" "default_master" {
  vpc_id               = alicloud_vpc.default.id
  cidr_block           = "192.168.4.0/24"
  zone_id              = "cn-hangzhou-i"
  vswitch_name         = "${var.name}-master"
  ipv6_cidr_block_mask = "4"
}

resource "alicloud_cen_instance" "default" {
  cen_instance_name = var.name
}

resource "alicloud_cen_transit_router" "default" {
  cen_id = alicloud_cen_instance.default.id
}

resource "alicloud_cen_transit_router_vpc_attachment" "default" {
  cen_id                                = alicloud_cen_instance.default.id
  vpc_id                                = alicloud_vpc.default.id
  transit_router_id                     = alicloud_cen_transit_router.default.transit_router_id
  transit_router_attachment_name        = var.name
  transit_router_attachment_description = var.name
  payment_type                          = "PayAsYouGo"

  options {
    ipv6_support           = "enable"
    appliance_mode_support = "enable"
  }

  zone_mappings {
    vswitch_id = alicloud_vswitch.default.id
    zone_id    = alicloud_vswitch.default.zone_id
  }

  zone_mappings {
    vswitch_id = alicloud_vswitch.default_master.id
    zone_id    = alicloud_vswitch.default_master.zone_id
  }
}

data "alicloud_cen_transit_router_vpc_attachments" "default" {
  cen_id = alicloud_cen_instance.default.id
  ids    = [alicloud_cen_transit_router_vpc_attachment.default.id]
}
`, rand)
}
