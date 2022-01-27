---
subcategory: "Elastic Desktop Service(EDS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecd_command"
sidebar_current: "docs-alicloud-resource-ecd-command"
description: |-
  Provides a Alicloud ECD Command resource.
---

# alicloud\_ecd\_command

Provides a ECD Command resource.

For information about ECD Command and how to use it, see [What is Command](https://help.aliyun.com/document_detail/188382.html).

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

```

## Argument Reference

The following arguments are supported:

* `command_content` - (Required, ForceNew) The Contents of the Script to Base64 Encoded Transmission.
* `command_type` - (Required, ForceNew) The Script Type. Valid values: `RunBatScript`, `RunPowerShellScript`.
* `content_encoding` - (Optional, Computed) That Returns the Data Encoding Method. Valid values: `Base64`, `PlainText`.
* `desktop_id` - (Required, ForceNew) The desktop id of the Desktop.
* `timeout` - (Optional) The timeout period for script execution the unit is seconds. Default to: `60`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Command.
* `status` - Script Is Executed in the Overall Implementation of the State. Valid values: `Pending`, `Failed`, `PartialFailed`, `Running`, `Stopped`, `Stopping`, `Finished`, `Success`.

## Import

ECD Command can be imported using the id, e.g.

```
$ terraform import alicloud_ecd_command.example <id>
```