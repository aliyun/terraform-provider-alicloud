---
subcategory: "Realtime Compute"
layout: "alicloud"
page_title: "Alicloud: alicloud_realtime_compute_vvp_instance"
description: |-
  Provides a Alicloud Realtime Compute Vvp Instance resource.
---

# alicloud_realtime_compute_vvp_instance

Provides a Realtime Compute Vvp Instance resource.

For information about Realtime Compute Vvp Instance and how to use it, see [What is Vvp Instance](https://next.api.alibabacloud.com/api/foasconsole/2019-06-01/CreateInstance).

-> **NOTE:** Available since v1.214.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = "cn-hangzhou-h"
}

resource "alicloud_oss_bucket" "defaultOSS" {
  bucket = "${var.name}-${random_integer.default.result}"
}

data "alicloud_resource_manager_resource_groups" "default" {
  status = "OK"
}

resource "alicloud_realtime_compute_vvp_instance" "default" {
  storage {
    oss {
      bucket = alicloud_oss_bucket.defaultOSS.bucket
    }
  }

  vvp_instance_name = "${var.name}-${random_integer.default.result}"
  vpc_id            = data.alicloud_vpcs.default.ids.0
  zone_id           = "cn-hangzhou-h"
  vswitch_ids = [
    "${data.alicloud_vswitches.default.ids.0}"
  ]
  payment_type = "PayAsYouGo"
}
```

### Deleting `alicloud_realtime_compute_vvp_instance` or removing it from your configuration

The `alicloud_realtime_compute_vvp_instance` resource allows you to manage  `payment_type = "Subscription"`  instance, but Terraform cannot destroy it.
Deleting the subscription resource or removing it from your configuration will remove it from your state file and management, but will not destroy the Instance.
You can resume managing the subscription instance via the AlibabaCloud Console.

## Argument Reference

The following arguments are supported:
* `duration` - (Optional) The number of subscription periods. If the payment type is PRE, this parameter is required.
* `payment_type` - (Required, ForceNew) The payment type of the resource.
* `pricing_cycle` - (Optional) The subscription period. If the payment type is PRE, this parameter is required.
* `resource_group_id` - (Optional, Computed) The resource group to which the newly purchased instance belongs.
* `resource_spec` - (Optional) Resource specifications. See [`resource_spec`](#resource_spec) below.
* `storage` - (Required, ForceNew) Store information. See [`storage`](#storage) below.
* `vswitch_ids` - (Required, ForceNew) Virtual Switch ID.
* `vpc_id` - (Required, ForceNew) The VPC ID of the user.
* `vvp_instance_name` - (Required, ForceNew) The name of the workspace.
* `zone_id` - (Required, ForceNew) The zone ID of the resource.
* `tags` - (Optional) The tags of the resource.

### `resource_spec`

The resource_spec supports the following:
* `cpu` - (Optional) CPU number.
* `memory_gb` - (Optional) Memory size.

### `storage`

The storage supports the following:
* `oss` - (Required, ForceNew) OSS stores information. See [`oss`](#storage-oss) below.

### `storage-oss`

The oss supports the following:
* `bucket` - (Required, ForceNew) OSS Bucket name.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the resource.
* `status` - The status of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Vvp Instance.
* `delete` - (Defaults to 5 mins) Used when delete the Vvp Instance.
* `update` - (Defaults to 5 mins) Used when update the Vvp Instance.

## Import

Realtime Compute Vvp Instance can be imported using the id, e.g.

```shell
$ terraform import alicloud_realtime_compute_vvp_instance.example <id>
```