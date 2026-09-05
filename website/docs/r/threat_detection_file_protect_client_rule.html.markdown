---
subcategory: "Threat Detection"
layout: "alicloud"
page_title: "Alicloud: alicloud_threat_detection_file_protect_client_rule"
description: |-
  Provides a Alicloud Threat Detection File Protect Client Rule resource.
---

# alicloud_threat_detection_file_protect_client_rule

Provides a Threat Detection File Protect Client Rule resource.

Client file protection event monitoring rules, including file read/write, delete, rename, and permission change.

For information about Threat Detection File Protect Client Rule and how to use it, see [What is File Protect Client Rule](https://next.api.alibabacloud.com/document/Sas/2018-12-03/CreateFileProtectClientRule).

-> **NOTE:** Available since v1.287.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_threat_detection_file_protect_client_rule" "default" {
  rule_name   = var.name
  rule_action = "pass"
  status      = 0
  alert_level = 1
  platform    = "linux"
  file_paths  = ["/opt/a", "/tmp/d"]
  file_ops    = ["READ", "WRITE"]
  proc_paths  = ["/usr/bin/java"]
}
```

## Argument Reference

The following arguments are supported:
* `alert_level` - (Optional, Int) 0 no alert 1 info 2 suspicious 3 critical
* `exclude_users` - (Optional, List) The list of users to exclude.
* `file_ops` - (Required, List) The operations that you want to perform on the files.
* `file_paths` - (Required, List) The paths to the files that you want to monitor. Wildcard characters are supported.
* `file_types` - (Optional, List) The list of file types to monitor.
* `platform` - (Optional, ForceNew) The type of the operating system. Valid values:

  - `windows`: Windows
  - `linux`: Linux
* `proc_paths` - (Required, List) The paths to the processes that you want to monitor. Wildcard characters are supported.
* `rule_action` - (Required) The action of the rule. Valid values: `pass`, `monitor`, `block`.
* `rule_name` - (Required) The name of the rule.
* `status` - (Required, Int) Specifies whether to enable the rule. Valid values:

  - `1`: yes
  - `0`: no

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `switch_id` - The switch ID of the rule, assigned by the server.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the File Protect Client Rule.
* `delete` - (Defaults to 5 mins) Used when delete the File Protect Client Rule.
* `update` - (Defaults to 5 mins) Used when update the File Protect Client Rule.

## Import

Threat Detection File Protect Client Rule can be imported using the id, e.g.

```shell
$ terraform import alicloud_threat_detection_file_protect_client_rule.example <id>
```