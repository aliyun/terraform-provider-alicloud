---
subcategory: "ENS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ens_common_bandwidth_package"
description: |-
  Provides a Alicloud ENS Common Bandwidth Package resource.
---

# alicloud_ens_common_bandwidth_package

Provides a ENS Common Bandwidth Package resource.

ENS shared bandwidth package.

For information about ENS Common Bandwidth Package and how to use it, see [What is Common Bandwidth Package](https://next.api.alibabacloud.com/document/Ens/2017-11-10/CreateCommonBandwidthPackage).

-> **NOTE:** Available since v1.289.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_ens_common_bandwidth_package" "default" {
  bandwidth     = 10
  description   = var.name
  ens_region_id = "cn-xxx-7fj20x"
  name          = var.name
}
```

## Argument Reference

The following arguments are supported:
* `bandwidth` - (Required, Int) Speed limit bandwidth value, unit: Mbps.
* `description` - (Optional) Description information.
* `ens_region_id` - (Required, ForceNew) The ID of the ENS node.
* `name` - (Optional) The name of the resource.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `creation_time` - The creation time of the resource.
* `status` - The status of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Common Bandwidth Package.
* `delete` - (Defaults to 5 mins) Used when delete the Common Bandwidth Package.
* `update` - (Defaults to 5 mins) Used when update the Common Bandwidth Package.

## Import

ENS Common Bandwidth Package can be imported using the id, e.g.

```shell
$ terraform import alicloud_ens_common_bandwidth_package.example <bandwidth_package_id>
```
