---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_instance"
description: |-
  Provides a Alicloud Cloud Enterprise Network (CEN) Cen Instance resource.
---

# alicloud_cen_instance

Provides a Cloud Enterprise Network (CEN) Cen Instance resource.



For information about Cloud Enterprise Network (CEN) Cen Instance and how to use it, see [What is Cen Instance](https://www.alibabacloud.com/help/en/cen/developer-reference/api-cbn-2017-09-12-createcen).

-> **NOTE:** Available since v1.15.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_cen_instance" "default" {
  cen_instance_name = var.name
  description       = var.name
}
```

## Argument Reference

The following arguments are supported:
* `cen_instance_name` - (Optional, Available since v1.98.0) The name of the CEN instance.
* `description` - (Optional) The description of the CEN instance.
* `protection_level` - (Optional, Available since v1.76.0) The level of CIDR block overlapping. Valid values:  REDUCED: Overlapped CIDR blocks are allowed. However, the overlapped CIDR blocks cannot be the same.
* `resource_group_id` - (Optional, Computed, Available since v1.232.0) The ID of the resource group
* `tags` - (Optional, Map, Available since v1.80.0) The tags of the CEN instance.

The following arguments will be discarded. Please use new fields as soon as possible:
* `name` - (Deprecated since v1.98.0). Field 'name' has been deprecated from provider version 1.246.0. New field 'cen_instance_name' instead.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The time when the CEN instance was created.
* `status` - The state of the CEN instance.   Creating: The CEN instance is being created. Active: The CEN instance is running. Deleting: The CEN instance is being deleted.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Cen Instance.
* `delete` - (Defaults to 5 mins) Used when delete the Cen Instance.
* `update` - (Defaults to 5 mins) Used when update the Cen Instance.

## Import

Cloud Enterprise Network (CEN) Cen Instance can be imported using the id, e.g.

```shell
$ terraform import alicloud_cen_instance.example <id>
```