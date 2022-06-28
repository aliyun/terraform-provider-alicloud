---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_command"
sidebar_current: "docs-alicloud-resource-ecs-command"
description: |-
  Provides a Alicloud ECS Command resource.
---

# alicloud\_ecs\_command

Provides a ECS Command resource.

For information about ECS Command and how to use it, see [What is Command](https://www.alibabacloud.com/help/en/doc-detail/64844.htm).

-> **NOTE:** Available in v1.116.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_ecs_command" "example" {
  name            = "tf-testAcc"
  command_content = "bHMK"
  description     = "For Terraform Test"
  type            = "RunShellScript"
  working_dir     = "/root"
}

```

## Argument Reference

The following arguments are supported:

* `command_content` - (Required, ForceNew) The Base64-encoded content of the command.
* `description` - (Optional, ForceNew) The description of command.
* `enable_parameter` - (Optional, ForceNew) Specifies whether to use custom parameters in the command to be created. Default to: false.                                                                                                                  
* `name` - (Required, ForceNew) The name of the command, which supports all character sets. It can be up to 128 characters in length.
* `timeout` - (Optional, ForceNew) The timeout period that is specified for the command to be run on ECS instances. Unit: seconds. Default to: `60`.
* `type` - (Required, ForceNew) The command type. Valid Values: `RunBatScript`, `RunPowerShellScript` and `RunShellScript`.
* `working_dir` - (Optional, ForceNew) The execution path of the command in the ECS instance.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Command.

## Import

ECS Command can be imported using the id, e.g.

```
$ terraform import alicloud_ecs_command.example <id>
```
