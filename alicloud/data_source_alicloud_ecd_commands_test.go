package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudECDCommandsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcdCommandsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ecd_command.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudEcdCommandsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ecd_command.default.id}_fake"]`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcdCommandsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_ecd_command.default.id}"]`,
			"status": `"Success"`,
		}),
		fakeConfig: testAccCheckAlicloudEcdCommandsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_ecd_command.default.id}"]`,
			"status": `"Running"`,
		}),
	}
	commandTypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcdCommandsDataSourceName(rand, map[string]string{
			"ids":          `["${alicloud_ecd_command.default.id}"]`,
			"command_type": `"RunPowerShellScript"`,
		}),
	}

	contentEncodingConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcdCommandsDataSourceName(rand, map[string]string{
			"ids":              `["${alicloud_ecd_command.default.id}"]`,
			"content_encoding": `"PlainText"`,
		}),
	}
	desktopIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcdCommandsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_ecd_command.default.id}"]`,
			"desktop_id": `"${alicloud_ecd_desktop.default.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudEcdCommandsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_ecd_command.default.id}"]`,
			"desktop_id": `"${alicloud_ecd_desktop.default.id}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcdCommandsDataSourceName(rand, map[string]string{
			"ids":              `["${alicloud_ecd_command.default.id}"]`,
			"status":           `"Success"`,
			"command_type":     `"RunPowerShellScript"`,
			"content_encoding": `"PlainText"`,
			"desktop_id":       `"${alicloud_ecd_desktop.default.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudEcdCommandsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_ecd_command.default.id}_fake"]`,
			"status":     `"Running"`,
			"desktop_id": `"${alicloud_ecd_desktop.default.id}_fake"`,
		}),
	}
	var existAlicloudEcdCommandsBusesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                          "1",
			"commands.#":                                     "1",
			"commands.0.id":                                  CHECKSET,
			"commands.0.command_content":                     "ipconfig",
			"commands.0.command_type":                        "RunPowerShellScript",
			"commands.0.invoke_id":                           CHECKSET,
			"commands.0.status":                              "Success",
			"commands.0.invoke_desktops.#":                   CHECKSET,
			"commands.0.invoke_desktops.0.desktop_id":        CHECKSET,
			"commands.0.invoke_desktops.0.dropped":           CHECKSET,
			"commands.0.invoke_desktops.0.error_code":        "",
			"commands.0.invoke_desktops.0.error_info":        "",
			"commands.0.invoke_desktops.0.exit_code":         CHECKSET,
			"commands.0.invoke_desktops.0.finish_time":       CHECKSET,
			"commands.0.invoke_desktops.0.invocation_status": "Success",
			"commands.0.invoke_desktops.0.output":            CHECKSET,
			"commands.0.invoke_desktops.0.repeats":           CHECKSET,
			"commands.0.invoke_desktops.0.start_time":        CHECKSET,
			"commands.0.invoke_desktops.0.stop_time":         "",
		}
	}
	var fakeAlicloudEcdCommandsBusesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"commands.#": "0",
		}
	}
	var alicloudEcdCommandsBusesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_ecd_commands.default",
		existMapFunc: existAlicloudEcdCommandsBusesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudEcdCommandsBusesDataSourceNameMapFunc,
	}
	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.EcdSupportRegions)
	}
	alicloudEcdCommandsBusesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, statusConf, commandTypeConf, contentEncodingConf, desktopIdConf, allConf)
}
func testAccCheckAlicloudEcdCommandsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {
  default = "tf-testacc-ecdcommands%d"
}

data "alicloud_ecd_bundles" "default"{
	bundle_type = "SYSTEM"
    name_regex  = "windows"
}
resource "alicloud_ecd_simple_office_site" "default" {
  cidr_block = "172.16.0.0/12"
  desktop_access_type = "Internet"
  office_site_name    = var.name
}
resource "alicloud_ecd_policy_group" "default" {
  policy_group_name = var.name
  clipboard = "readwrite"
  local_drive = "read"
  authorize_access_policy_rules{
    description= var.name
    cidr_ip=     "1.2.3.4/24"
  }
  authorize_security_policy_rules  {
    type=        "inflow"
    policy=      "accept"
    description=  var.name
    port_range= "80/80"
    ip_protocol= "TCP"
    priority=    "1"
    cidr_ip=     "0.0.0.0/0"
  }
}
resource "alicloud_ecd_desktop" "default" {
	office_site_id  = alicloud_ecd_simple_office_site.default.id
	policy_group_id = alicloud_ecd_policy_group.default.id
	bundle_id 		= data.alicloud_ecd_bundles.default.bundles.0.id
	desktop_name 	=  var.name
}

resource "alicloud_ecd_command" "default" {
	command_content = "ipconfig"
	command_type    = "RunPowerShellScript"
	desktop_id      = alicloud_ecd_desktop.default.id
}

data "alicloud_ecd_commands" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
