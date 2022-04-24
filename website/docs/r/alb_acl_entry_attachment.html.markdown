---
subcategory: "Application Load Balancer (ALB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_alb_acl_entry_attachment"
sidebar_current: "docs-alicloud-resource-alb-acl-entry-attachment"
description: |-
  Provides a Acl entry attachment resource.
---

# alicloud\_alb\_acl\_entry\_attachment

-> **NOTE:** Available in v1.166.0+.


For information about acl entry attachment and how to use it, see [Configure an acl entry](https://www.alibabacloud.com/help/en/server-load-balancer/latest/addentriestoacl).


## Example Usage

```
variable "name" {
  default = "terraformalbaclconfig"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_alb_acl" "default" {
  acl_name          = var.name
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.groups.0.id
}

resource "alicloud_alb_acl_entry_attachment" "default" {
  acl_id      = alicloud_alb_acl.default.id
  entry       = "168.10.10.0/24"
  description = var.name
}
```

## Argument Reference

The following arguments are supported:

* `acl_id` - (Required, ForceNew) The ID of the Acl.
* `entry` - (Required, ForceNew) The CIDR blocks.
* `description` - (Optional, ForceNew) The description of the entry.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource. The value formats as `<acl_id>:<entry>`.
* `status` - The Status of the resource.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the resource.
* `delete` - (Defaults to 5 mins) Used when delete the resource.

## Import

Acl entry attachment can be imported using the id, e.g.

```
$ terraform import alicloud_alb_acl_entry_attachment.example <acl_id>:<entry>
```
