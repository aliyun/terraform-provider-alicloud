---
subcategory: "Threat Detection"
layout: "alicloud"
page_title: "Alicloud: alicloud_threat_detection_client_file_protect"
description: |-
  Provides a Alicloud Threat Detection Client File Protect resource.
---

# alicloud_threat_detection_client_file_protect

Provides a Threat Detection Client File Protect resource. Client core file protection event monitoring, including file reading and writing, deletion, and permission change.

For information about Threat Detection Client File Protect and how to use it, see [What is Client File Protect](https://www.alibabacloud.com/help/zh/security-center/developer-reference/api-sas-2018-12-03-createfileprotectrule).

-> **NOTE:** Available since v1.212.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}


resource "alicloud_threat_detection_client_file_protect" "default" {
  status      = "0"
  file_paths  = ["/usr/local"]
  file_ops    = ["CREATE"]
  rule_action = "pass"
  proc_paths  = ["/usr/local"]
  alert_level = "0"
  switch_id   = "FILE_PROTECT_RULE_SWITCH_TYPE_1693474122929"
  rule_name   = "rule_example"
}
```

## Argument Reference

The following arguments are supported:
* `alert_level` - (Optional) 0 no alert 1 info 2 suspicious 3 critical.
* `file_ops` - (Required) file operation.
* `file_paths` - (Required) file path.
* `proc_paths` - (Required) process path.
* `rule_action` - (Required) rule action, pass or alert.
* `rule_name` - (Required) ruleName.
* `status` - (Optional, Computed) rule status 0 is disable 1 is enable.
* `switch_id` - (Optional, ForceNew) switch id.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Client File Protect.
* `delete` - (Defaults to 5 mins) Used when delete the Client File Protect.
* `update` - (Defaults to 5 mins) Used when update the Client File Protect.

## Import

Threat Detection Client File Protect can be imported using the id, e.g.

```shell
$ terraform import alicloud_threat_detection_client_file_protect.example <id>
```