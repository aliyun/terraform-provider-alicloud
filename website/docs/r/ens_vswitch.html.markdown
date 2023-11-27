---
subcategory: "ENS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ens_vswitch"
description: |-
  Provides a Alicloud ENS Vswitch resource.
---

# alicloud_ens_vswitch

Provides a ENS Vswitch resource. 

For information about ENS Vswitch and how to use it, see [What is Vswitch](https://www.alibabacloud.com/help/en/ens/developer-reference/api-createvswitch).

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

resource "alicloud_ens_network" "default23T2cD" {
  network_name = var.name

  description   = "VSwitchDescription_autotest"
  cidr_block    = "192.168.2.0/24"
  ens_region_id = "cn-hangzhou-44"
}


resource "alicloud_ens_vswitch" "default" {
  description  = "VSwitchDescription_autotest"
  cidr_block   = "192.168.2.0/24"
  vswitch_name = var.name

  ens_region_id = "cn-hangzhou-44"
  network_id    = alicloud_ens_network.default23T2cD.id
}
```

## Argument Reference

The following arguments are supported:
* `cidr_block` - (Required, ForceNew) IPv4 CIDR block of the VSwitch instanceExample value: 10.21.224.0/22.
* `description` - (Optional) Description of the VSwitch Instance.
* `ens_region_id` - (Required, ForceNew) ENS Region ID。.
* `network_id` - (Optional, ForceNew, Computed) Network ID of the vSwitch instanceExample value: vpc-25cdvfeq58pl * * * *.
* `vswitch_name` - (Optional) Name of the switch instance.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the VSwitch instance, in the UTC time format, yyyy-MM-ddTHH:mm:ssZ.
* `status` - Status of the switch instanceValues available:Pending: ConfiguringAvailable: Available.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Vswitch.
* `delete` - (Defaults to 5 mins) Used when delete the Vswitch.
* `update` - (Defaults to 5 mins) Used when update the Vswitch.

## Import

ENS Vswitch can be imported using the id, e.g.

```shell
$ terraform import alicloud_ens_vswitch.example <id>
```