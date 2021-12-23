---
subcategory: "Elastic Desktop Service(EDS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecd_commands"
sidebar_current: "docs-alicloud-datasource-ecd-commands"
description: |-
  Provides a list of Ecd Commands to the user.
---

# alicloud\_ecd\_commands

This data source provides the Ecd Commands of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.146.0+.

## Example Usage

Basic Usage

```terraform

resource "alicloud_ecd_simple_office_site" "default" {
  cidr_block          = "172.16.0.0/12"
  desktop_access_type = "Internet"
  office_site_name    = "your_office_site_name"

}
data "alicloud_ecd_bundles" "default" {
  bundle_type = "SYSTEM"
  name_regex  = "windows"
}

resource "alicloud_ecd_policy_group" "default" {
  policy_group_name = "your_policy_group_name"
  clipboard         = "readwrite"
  local_drive       = "read"
  authorize_access_policy_rules {
    description = "example_value"
    cidr_ip     = "1.2.3.4/24"
  }
  authorize_security_policy_rules {
    type        = "inflow"
    policy      = "accept"
    description = "example_value"
    port_range  = "80/80"
    ip_protocol = "TCP"
    priority    = "1"
    cidr_ip     = "0.0.0.0/0"
  }
}
resource "alicloud_ecd_desktop" "default" {
  office_site_id  = alicloud_ecd_simple_office_site.default.id
  policy_group_id = alicloud_ecd_policy_group.default.id
  bundle_id       = data.alicloud_ecd_bundles.default.bundles.0.id
  desktop_name    = var.name
}

resource "alicloud_ecd_command" "default" {
  command_content = "ipconfig"
  command_type    = "RunPowerShellScript"
  desktop_id      = alicloud_ecd_desktop.default.id
}

data "alicloud_ecd_commands" "ids" {}
output "ecd_command_id_1" {
  value = data.alicloud_ecd_commands.ids.commands.0.id
}

```

## Argument Reference

The following arguments are supported:

* `command_type` - (Optional, ForceNew) The Script Type. Valid values: `RunBatScript`, `RunPowerShellScript`.
* `content_encoding` - (Optional, ForceNew) That Returns the Data Encoding Method. Valid values: `Base64`, `PlainText`.
* `ids` - (Optional, ForceNew, Computed)  A list of Command IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `desktop_id` - (Optional, ForceNew) The desktop id of the Desktop.
* `status` - (Optional, ForceNew) Script Is Executed in the Overall Implementation of the State. Valid values: `Pending`, `Failed`, `PartialFailed`, `Running`, `Stopped`, `Stopping`, `Finished`, `Success`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `commands` - A list of Ecd Commands. Each element contains the following attributes:
	* `command_content` - The Contents of the Script to Base64 Encoded Transmission.
	* `command_type` - The Script Type. Valid values: `RunBatScript`, `RunPowerShellScript`.
	* `create_time` - The Task of Creation Time.
	* `id` - The ID of the Command.
	* `invoke_desktops` - The Implementation of the Target Cloud Desktop Collection.
		* `error_code` - Command of the Failure Or Perform the Reason for the Failure of the Code.
		* `error_info` - Command of the Failure Or Perform the Reason for the Failure of the Details.
		* `exit_code` - Command of the Failure Or Perform the Reason for the Failure of the Details.
		* `start_time` - The Script Process on the Desktop, in the Start Timing of the Execution.
		* `stop_time` - If You Use the invocation Indicates That the Call of the Time.
		* `desktop_id` - The desktop id of the Desktop.
		* `dropped` - Output Field Text Length Exceeds 24 KB of Truncated Discarded Text Length.
		* `finish_time` - The Script Process until the End of Time.
		* `invocation_status` - A Single Cloud Desktop Script Progress Status.
		* `output` - Script the Output of the Process.
		* `repeats` - Command in the Desktop Implementation.
	* `invoke_id` - The invoke id of the Command.
	* `status` - Script Is Executed in the Overall Implementation of the State. Valid values: `Pending`, `Failed`, `PartialFailed`, `Running`, `Stopped`, `Stopping`, `Finished`, `Success`.