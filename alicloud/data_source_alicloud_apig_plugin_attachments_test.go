package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudApigPluginAttachmentDataSource(t *testing.T) {
	testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudApigPluginAttachmentSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_apig_plugin_attachment.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudApigPluginAttachmentSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_apig_plugin_attachment.default.id}_fake"]`,
		}),
	}

	EnvironmentIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudApigPluginAttachmentSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_apig_plugin_attachment.default.id}"]`,
			"environment_id": `"${alicloud_apig_gateway.default.environments.0.environment_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudApigPluginAttachmentSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_apig_plugin_attachment.default.id}_fake"]`,
			"environment_id": `"${alicloud_apig_gateway.default.environments.0.environment_id}_fake"`,
		}),
	}
	AttachResourceTypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudApigPluginAttachmentSourceConfig(rand, map[string]string{
			"ids":                  `["${alicloud_apig_plugin_attachment.default.id}"]`,
			"attach_resource_type": `"HttpApi"`,
		}),
		fakeConfig: testAccCheckAlicloudApigPluginAttachmentSourceConfig(rand, map[string]string{
			"ids":                  `["${alicloud_apig_plugin_attachment.default.id}_fake"]`,
			"attach_resource_type": `"HttpApi_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudApigPluginAttachmentSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_apig_plugin_attachment.default.id}"]`,
			"environment_id": `"${alicloud_apig_gateway.default.environments.0.environment_id}"`,

			"attach_resource_type": `"HttpApi"`,
		}),
		fakeConfig: testAccCheckAlicloudApigPluginAttachmentSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_apig_plugin_attachment.default.id}_fake"]`,
			"environment_id": `"${alicloud_apig_gateway.default.environments.0.environment_id}_fake"`,

			"attach_resource_type": `"HttpApi_fake"`,
		}),
	}

	ApigPluginAttachmentCheckInfo.dataSourceTestCheck(t, rand, idsConf, EnvironmentIdConf, AttachResourceTypeConf, allConf)
}

var existApigPluginAttachmentMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"attachments.#":                         "1",
		"attachments.0.attach_resource_ids.#":   CHECKSET,
		"attachments.0.plugin_class_info.#":     CHECKSET,
		"attachments.0.attach_resource_names.#": CHECKSET,
		"attachments.0.attach_resource_type":    CHECKSET,
		"attachments.0.environment_id":          CHECKSET,
		"attachments.0.plugin_attachment_id":    CHECKSET,
		"attachments.0.enable":                  CHECKSET,
		"attachments.0.plugin_info.#":           CHECKSET,
		"attachments.0.attach_resource_id":      CHECKSET,
	}
}

var fakeApigPluginAttachmentMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"attachments.#": "0",
	}
}

var ApigPluginAttachmentCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_apig_plugin_attachments.default",
	existMapFunc: existApigPluginAttachmentMapFunc,
	fakeMapFunc:  fakeApigPluginAttachmentMapFunc,
}

func testAccCheckAlicloudApigPluginAttachmentSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccApigPluginAttachment%d"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_apig_gateway" "default" {
  network_access_config {
    type = "Internet"
  }
  vswitch {
    vswitch_id = data.alicloud_vswitches.default.ids.0
  }
  zone_config {
    select_option = "Auto"
  }
  vpc {
    vpc_id = data.alicloud_vpcs.default.ids.0
  }
  gateway_type = "API"
  payment_type = "PayAsYouGo"
  gateway_name = var.name
  spec         = "apigw.small.x1"
  log_config {
    sls {
      enable = true
    }
  }
}

resource "alicloud_apig_plugin" "default" {
  plugin_class_id = "pls-crpqb35lhtgo800k2m86"
  gateway_id      = alicloud_apig_gateway.default.id
}


resource "alicloud_apig_http_api" "default" {
  http_api_name = var.name
  protocols     = ["HTTP"]
  type          = "Rest"
  base_path     = "/terraform-acc"
}

resource "alicloud_apig_plugin_attachment" "default" {
  attach_resource_ids  = ["${alicloud_apig_http_api.default.id}"]
  environment_id       = alicloud_apig_gateway.default.environments.0.environment_id
  enable               = true
  attach_resource_type = "HttpApi"
  plugin_info {
    plugin_config = "eyJ0ZXN0IjoiaGVsbG8ifQ=="
    plugin_id     = alicloud_apig_plugin.default.id
    gateway_id    = alicloud_apig_gateway.default.id
  }
}

data "alicloud_apig_plugin_attachments" "default" {
  %s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
