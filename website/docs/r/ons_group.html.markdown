---
subcategory: "RocketMQ"
layout: "alicloud"
page_title: "Alicloud: alicloud_ons_group"
sidebar_current: "docs-alicloud-resource-ons-group"
description: |-
  Provides a Alicloud ONS Group resource.
---

# alicloud\_ons\_group

Provides an ONS group resource.

For more information about how to use it, see [RocketMQ Group Management API](https://www.alibabacloud.com/help/doc-detail/29616.html). 

-> **NOTE:** Available in 1.53.0+

## Example Usage

Basic Usage

```
variable "name" {
  default = "onsInstanceName"
}

variable "group_id" {
  default = "GID-onsGroupDatasourceName"
}

resource "alicloud_ons_instance" "default" {
  name = "${var.name}"
  remark = "default_ons_instance_remark"
}

resource "alicloud_ons_group" "default" {
  group_id = "${var.group_id}"
  instance_id = "${alicloud_ons_instance.default.id}"
  remark = "dafault_ons_group_remark"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) ID of the ONS Instance that owns the groups.
* `group_id` - (Required) Name of the group. Two groups on a single instance cannot have the same name. A `group_id` starts with "GID_" or "GID-", and contains letters, numbers, hyphens (-), and underscores (_).
* `remark` - (Optional) This attribute is a concise description of group. The length cannot exceed 256.
* `read_enable` - (Optional) This attribute is used to set the message reading enabled or disabled. It can only be set after the group is used by the client.

## Attributes Reference

The following attributes are exported:

* `id` - The `key` of the resource supplied above. The value is formulated as `<instance_id>:<group_id>`.

## Import

ONS GROUP can be imported using the id, e.g.

```
$ terraform import alicloud_ons_group.group MQ_INST_1234567890_Baso1234567:GID-onsGroupDemo
```
