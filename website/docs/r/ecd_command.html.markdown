---
subcategory: "ECD"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecd_command"
sidebar_current: "docs-alicloud-resource-ecd-command"
description: |-
  Provides a Alicloud ECD Command resource.
---

# alicloud\_ecd\_command

Provides a ECD Command resource.

For information about ECD Command and how to use it, see [What is Command](https://help.aliyun.com/).

-> **NOTE:** Available in v1.144.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_ecd_command" "example" {
  command_content = "example_value"
  command_type    = "RunBatScript"
  desktop_id      = [example_value]
}

```

## Argument Reference

The following arguments are supported:

* `command_content` - (Required, ForceNew) The Contents of the Script to Base64 Encoded Transmission.
* `command_type` - (Required, ForceNew) The Script Type. Valid values: `RunBatScript`, `RunPowerShellScript`.
* `content_encoding` - (Optional, Computed, ForceNew) That Returns the Data Encoding Method. Valid values: `Base64`, `PlainText`.
* `desktop_id` - (Required, ForceNew) Cloud Desktop ID.
* `timeout` - (Optional) The timeout.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Command.
* `status` - Script Is Executed in the Overall Implementation of the State.

## Import

ECD Command can be imported using the id, e.g.

```
$ terraform import alicloud_ecd_command.example <id>
```