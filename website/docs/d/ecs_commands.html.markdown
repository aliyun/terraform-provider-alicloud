---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_commands"
sidebar_current: "docs-alicloud-datasource-ecs-commands"
description: |-
  Provides a list of Ecs Commands to the user.
---

# alicloud\_ecs\_commands

This data source provides the Ecs Commands of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.116.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_ecs_commands" "example" {
  ids        = ["E2RY53-xxxx"]
  name_regex = "tf-testAcc"
}

output "first_ecs_command_id" {
  value = data.alicloud_ecs_commands.example.commands.0.id
}
```

## Argument Reference

The following arguments are supported:

* `content_encoding` - (Optional, ForceNew) The Base64-encoded content of the command.
* `description` - (Optional, ForceNew) The description of command.
* `ids` - (Optional, ForceNew, Computed)  A list of Command IDs.
* `name` - (Optional, ForceNew) The name of the command.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Command name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `command_provider` - (Optional, ForceNew) Public order provider.
* `type` - (Optional, ForceNew) The command type. Valid Values: `RunBatScript`, `RunPowerShellScript` and `RunShellScript`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Command names.
* `commands` - A list of Ecs Commands. Each element contains the following attributes:
	* `command_content` - The Base64-encoded content of the command.
	* `command_id` - The ID of the Command.
	* `description` - The description of command.
	* `enable_parameter` - Specifies whether to use custom parameters in the command to be created.
	* `id` - The ID of the Command.
	* `name` - The name of the command
	* `parameter_names` - A list of custom parameter names which are parsed from the command content specified when the command was being created.
	* `timeout` - The timeout period that is specified for the command to be run on ECS instances.
	* `type` - The command type.
	* `working_dir` - The execution path of the command in the ECS instance.
