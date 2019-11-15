---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_instance_attachment"
sidebar_current: "docs-alicloud-resource-cen-instance-attachment"
description: |-
  Provides a Alicloud CEN child instance attachment resource.
---

# alicloud\_cen_instance_attachment

Provides a CEN child instance attachment resource.

## Example Usage

Basic Usage

```
# Create a new instance-attachment and use it to attach one child instance to a new CEN
variable "name" {
  default = "tf-testAccCenInstanceAttachmentBasic"
}

resource "alicloud_cen_instance" "cen" {
  name        = "${var.name}"
  description = "terraform01"
}

resource "alicloud_vpc" "vpc" {
  name       = "${var.name}"
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_cen_instance_attachment" "foo" {
  instance_id              = "${alicloud_cen_instance.cen.id}"
  child_instance_id        = "${alicloud_vpc.vpc.id}"
  child_instance_region_id = "cn-beijing"
}
```
## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The ID of the CEN.
* `child_instance_id` - (Required, ForceNew) The ID of the child instance to attach.
* `child_instance_region_id` - (Required, ForceNew) The region ID of the child instance to attach.
* `child_instance_owner_id` - (Optional, Available in 1.42.0+) The uid of the child instance. Only used when attach a child instance of other account.

->**NOTE:** Ensure that the child instance is not used in Express Connect.

## Attributes Reference

The following attributes are exported:

- `id` - ID of the resource, formatted as `<instance_id>:<child_instance_id>`.

## Import

CEN instance can be imported using the id, e.g.

```
$ terraform import alicloud_cen_instance_attachment.example cen-abc123456:vpc-abc123456
```
