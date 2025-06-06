---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_routine"
description: |-
  Provides a Alicloud ESA Routine resource.
---

# alicloud_esa_routine

Provides a ESA Routine resource.



For information about ESA Routine and how to use it, see [What is Routine](https://next.api.alibabacloud.com/document/ESA/2024-09-10/CreateRoutine).

-> **NOTE:** Available since v1.251.0.

## Example Usage

Basic Usage
```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_esa_routine" "default" {
  description = var.name
  name        = var.name
}
```

## Argument Reference

The following arguments are supported:
* `description` - (Optional, ForceNew) The routine name, which must be unique in the same account.
* `name` - (Required, ForceNew) Routine Name

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The time when the routine was created.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Routine.
* `delete` - (Defaults to 5 mins) Used when delete the Routine.

## Import

ESA Routine can be imported using the id, e.g.

```shell
$ terraform import alicloud_esa_routine.example <id>
```