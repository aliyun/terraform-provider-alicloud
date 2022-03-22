---
subcategory: "Classic Load Balancer (CLB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_slb_acl_entry_attachment"
sidebar_current: "docs-alicloud-resource-slb-acl-entry-attachment"
description: |-
  Provides a Acl entry attachment resource.
---

# alicloud\_slb\_acl\_entry\_attachment

-> **NOTE:** Available in v1.162.0+.

-> **NOTE:** The maximum number of entries per acl is 300.

For information about acl entry attachment and how to use it, see [Configure an acl entry](https://www.alibabacloud.com/help/en/doc-detail/70023.html).


## Example Usage

```
variable "name" {
  default = "terraformslbaclconfig"
}

variable "ip_version" {
  default = "ipv4"
}

resource "alicloud_slb_acl" "default" {
  name       = var.name
  ip_version = var.ip_version
}

resource "alicloud_slb_acl_entry_attachment" "default" {
  acl_id  = alicloud_slb_acl.default.id
  entry   = "168.10.10.0/24"
  comment = "second"
}
```

## Argument Reference

The following arguments are supported:

* `acl_id` - (Required, ForceNew) The ID of the Acl.
* `entry` - (Required, ForceNew) The CIDR blocks.
* `comment` - (Optional, ForceNew) The comment of the entry.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource. The value formats as `<acl_id>:<entry>`.


### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the resource.
* `delete` - (Defaults to 5 mins) Used when delete the resource.

## Import

Acl entry attachment can be imported using the id, e.g.

```
$ terraform import alicloud_slb_acl_entry_attachment.example <acl_id>:<entry>
```
