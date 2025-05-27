---
subcategory: "Eflo"
layout: "alicloud"
page_title: "Alicloud: alicloud_eflo_vsc"
description: |-
  Provides a Alicloud Eflo Vsc resource.
---

# alicloud_eflo_vsc

Provides a Eflo Vsc resource.

Virtual Storage Channel.

For information about Eflo Vsc and how to use it, see [What is Vsc](https://www.alibabacloud.com/help/en/pai/developer-reference/api-eflo-controller-2022-12-15-createvsc).

-> **NOTE:** Available since v1.250.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_eflo_vsc" "default" {
  vsc_type = "primary"
  node_id  = "e01-cn-9me49omda01"
  vsc_name = var.name
}
```

## Argument Reference

The following arguments are supported:
* `node_id` - (Required, ForceNew) The ID of the Node.
* `resource_group_id` - (Optional) The ID of the resource group.
* `tags` - (Optional, Map) The tag of the resource.
* `vsc_name` - (Optional, ForceNew) The name of the Vsc.
* `vsc_type` - (Optional, ForceNew) The type of the Vsc. Default value: `primary`. Valid values: `primary`.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `status` - The status of the Vsc.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Vsc.
* `delete` - (Defaults to 5 mins) Used when delete the Vsc.
* `update` - (Defaults to 5 mins) Used when update the Vsc.

## Import

Eflo Vsc can be imported using the id, e.g.

```shell
$ terraform import alicloud_eflo_vsc.example <id>
```
