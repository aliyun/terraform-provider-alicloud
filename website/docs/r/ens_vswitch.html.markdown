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

-> **NOTE:** Available since v1.213.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ens_vswitch&exampleId=37d441b7-70c9-0b35-25fb-44d7ba2378f43d298368&activeTab=example&spm=docs.r.ens_vswitch.0.37d441b770&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_ens_network" "default" {
  network_name = var.name

  description   = var.name
  cidr_block    = "192.168.2.0/24"
  ens_region_id = "cn-chenzhou-telecom_unicom_cmcc"
}


resource "alicloud_ens_vswitch" "default" {
  description  = var.name
  cidr_block   = "192.168.2.0/24"
  vswitch_name = var.name

  ens_region_id = "cn-chenzhou-telecom_unicom_cmcc"
  network_id    = alicloud_ens_network.default.id
}
```

## Argument Reference

The following arguments are supported:
* `cidr_block` - (Required, ForceNew) The CIDR block of the vSwitch.
* `description` - (Optional) The description of the vSwitch.
* `ens_region_id` - (Required, ForceNew) ENS Region ID.
* `network_id` - (Optional, ForceNew) The ID of the network to which the vSwitch that you want to create belongs.
* `vswitch_name` - (Optional) The name of the vSwitch.

## Attributes Reference

The following attributes are exported:
* `id` - The resource ID in terraform of Vswitch.
* `create_time` - The time when the VPC was created.
* `status` - The status of the vSwitch.

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
