---
subcategory: "ECD"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecd_commands"
sidebar_current: "docs-alicloud-datasource-ecd-commands"
description: |-
  Provides a list of Ecd Commands to the user.
---

# alicloud\_ecd\_commands

This data source provides the Ecd Commands of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.144.0+.

## Example Usage

Basic Usage

```terraform
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
* `include_output` - (Optional, ForceNew) The include output.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew) Script Is Executed in the Overall Implementation of the State.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `commands` - A list of Ecd Commands. Each element contains the following attributes:
	* `command_content` - The Contents of the Script to Base64 Encoded Transmission.
	* `command_type` - The Script Type.
	* `create_time` - The Task of Creation Time.
	* `id` - The ID of the Command.
	* `invoke_desktops` - The Implementation of the Target Cloud Desktop Collection.
		* `error_code` - Command of the Failure Or Perform the Reason for the Failure of the Code.
		* `error_info` - Command of the Failure Or Perform the Reason for the Failure of the Details.
		* `exit_code` - Command of the Failure Or Perform the Reason for the Failure of the Details.
		* `start_time` - The Script Process on the Desktop, in the Start Timing of the Execution.
		* `stop_time` - If You Use the stopinvocation Indicates That the Call of the Time.
		* `desktop_id` - Cloud Desktop ID.
		* `dropped` - Output Field Text Length Exceeds 24 KB of Truncated Discarded Text Length.
		* `finish_time` - The Script Process until the End of Time.
		* `invocation_status` - A Single Cloud Desktop Script Progress Status.
		* `output` - Script the Output of the Process.
		* `repeats` - Command in the Desktop Implementation.
	* `invoke_id` - Execution ID.
	* `status` - Script Is Executed in the Overall Implementation of the State.