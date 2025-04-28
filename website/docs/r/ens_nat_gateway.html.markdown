---
subcategory: "ENS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ens_nat_gateway"
description: |-
  Provides a Alicloud ENS Nat Gateway resource.
---

# alicloud_ens_nat_gateway

Provides a ENS Nat Gateway resource.

Nat gateway of ENS.

For information about ENS Nat Gateway and how to use it, see [What is Nat Gateway](https://next.api.alibabacloud.com/document/Ens/2017-11-10/CreateNatGateway).

-> **NOTE:** Available since v1.227.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ens_nat_gateway&exampleId=43c8bcc1-e0c4-65db-947f-404db031e6947b77c185&activeTab=example&spm=docs.r.ens_nat_gateway.0.43c8bcc1e0&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

variable "ens_region_id" {
  default = "cn-chenzhou-telecom_unicom_cmcc"
}

resource "alicloud_ens_network" "defaultObbrL7" {
  network_name  = var.name
  description   = var.name
  cidr_block    = "10.0.0.0/8"
  ens_region_id = var.ens_region_id
}

resource "alicloud_ens_vswitch" "defaulteFw783" {
  cidr_block    = "10.0.8.0/24"
  vswitch_name  = var.name
  ens_region_id = alicloud_ens_network.defaultObbrL7.ens_region_id
  network_id    = alicloud_ens_network.defaultObbrL7.id
}

resource "alicloud_ens_nat_gateway" "default" {
  vswitch_id    = alicloud_ens_vswitch.defaulteFw783.id
  ens_region_id = alicloud_ens_vswitch.defaulteFw783.ens_region_id
  network_id    = alicloud_ens_vswitch.defaulteFw783.network_id
  instance_type = "enat.default"
  nat_name      = var.name
}
```

## Argument Reference

The following arguments are supported:
* `ens_region_id` - (Required, ForceNew) The ID of the ENS node.
* `instance_type` - (Optional, ForceNew) NAT specifications. Value: `enat.default`.
* `nat_name` - (Optional) The name of the NAT gateway. The length is 1 to 128 characters, but it cannot start with 'http:// 'or 'https.
* `network_id` - (Required, ForceNew) The network ID.
* `vswitch_id` - (Required, ForceNew) The vSwitch ID.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - Creation time. UTC time, in the format of YYYY-MM-DDThh:mm:ssZ.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Nat Gateway.
* `delete` - (Defaults to 5 mins) Used when delete the Nat Gateway.
* `update` - (Defaults to 5 mins) Used when update the Nat Gateway.

## Import

ENS Nat Gateway can be imported using the id, e.g.

```shell
$ terraform import alicloud_ens_nat_gateway.example <id>
```