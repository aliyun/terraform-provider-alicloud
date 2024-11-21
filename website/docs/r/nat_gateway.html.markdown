---
subcategory: "NAT Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_nat_gateway"
sidebar_current: "docs-alicloud-resource-nat-gateway"
description: |-
  Provides a resource to create a VPC NAT Gateway.
---

# alicloud_nat_gateway

Provides a resource to create a VPC NAT Gateway.

-> **NOTE:** Resource bandwidth packages will not be supported since 00:00 on November 4, 2017, and public IP can be replaced be elastic IPs.
If a Nat Gateway has already bought some bandwidth packages, it can not bind elastic IP and you have to submit the [work order](https://selfservice.console.aliyun.com/ticket/createIndex) to solve.
If you want to add public IP, you can use resource 'alicloud_eip_association' to bind several elastic IPs for one Nat Gateway.

-> **NOTE:** From version 1.7.1, this resource has deprecated bandwidth packages.
But, in order to manage stock bandwidth packages, version 1.13.0 re-support configuring 'bandwidth_packages'.

-> **NOTE:** Available since v1.37.0.

## Example Usage

Basic usage

- create enhanced nat gateway
<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_nat_gateway&exampleId=db0b3938-fe52-60c0-dbbc-5d69bbc282dbcb073981&activeTab=example&spm=docs.r.nat_gateway.0.db0b3938fe&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf_example"
}

data "alicloud_enhanced_nat_available_zones" "default" {
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.0.0.0/8"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  zone_id      = data.alicloud_enhanced_nat_available_zones.default.zones.0.zone_id
  cidr_block   = "10.10.0.0/20"
  vpc_id       = alicloud_vpc.default.id
}

resource "alicloud_nat_gateway" "default" {
  vpc_id           = alicloud_vpc.default.id
  nat_gateway_name = var.name
  payment_type     = "PayAsYouGo"
  vswitch_id       = alicloud_vswitch.default.id
  nat_type         = "Enhanced"
}
```

- transform nat from Normal to Enhanced
-> **NOTE:** You must set `nat_type` to `Enhanced` and set `vswitch_id`.

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_nat_gateway&exampleId=690738fa-c76a-212e-6a46-8df639ab287a9e7e8b40&activeTab=example&spm=docs.r.nat_gateway.1.690738fac7&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}

data "alicloud_enhanced_nat_available_zones" "default" {
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.0.0.0/8"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  zone_id      = data.alicloud_enhanced_nat_available_zones.default.zones.0.zone_id
  cidr_block   = "10.10.0.0/20"
  vpc_id       = alicloud_vpc.default.id
}

resource "alicloud_nat_gateway" "default" {
  vpc_id           = alicloud_vpc.default.id
  nat_gateway_name = var.name
  vswitch_id       = alicloud_vswitch.default.id
  nat_type         = "Enhanced"
}
```

### Deleting `alicloud_nat_gateway` or removing it from your configuration

The `alicloud_nat_gateway` resource allows you to manage `payment_type = "Subscription"` or `instance_charge_type = "Prepaid"` nat gateway, but Terraform cannot destroy it.
Deleting the subscription resource or removing it from your configuration
will remove it from your statefile and management, but will not destroy the Nat Gateway.
You can resume managing the subscription nat gateway via the AlibabaCloud Console.

## Argument Reference

The following arguments are supported:

* `vpc_id` - (Required, ForceNew) The VPC ID.
* `specification` - (Optional) The specification of the nat gateway. Valid values are `Small`, `Middle` and `Large`. Effective when `internet_charge_type` is `PayBySpec` and `network_type` is `internet`. Details refer to [Nat Gateway Specification](https://help.aliyun.com/document_detail/203500.html).
* `nat_gateway_name` - (Optional, Available since v1.121.0) Name of the nat gateway. The value can have a string of 2 to 128 characters, must contain only alphanumeric characters or hyphens, such as "-",".","_", and must not begin or end with a hyphen, and must not begin with http:// or https://. Defaults to null.
* `description` - (Optional) Description of the nat gateway, This description can have a string of 2 to 256 characters, It cannot begin with http:// or https://. Defaults to null.
* `dry_run` - (Optional) Specifies whether to only precheck this request. Default value: `false`.
* `force` - (Optional) Specifies whether to forcefully delete the NAT gateway.
* `payment_type` - (Optional, ForceNew, Available since v1.121.0) The billing method of the NAT gateway. Valid values are `PayAsYouGo` and `Subscription`. Default to `PayAsYouGo`.
* `period` - (Optional, Available since v1.45.0) The duration that you will buy the resource, in month. It is valid when `payment_type` is `Subscription`. Valid values: [1-9, 12, 24, 36]. At present, the provider does not support modify "period" and you can do that via web console. **NOTE:** International station only supports `Subscription`.
-> **NOTE:** The attribute `period` is only used to create Subscription instance or modify the PayAsYouGo instance to Subscription. Once effect, it will not be modified that means running `terraform apply` will not effect the resource.
* `nat_type` - (Optional, Available since v1.102.0) The type of NAT gateway. Valid values: `Enhanced`. **NOTE:** From version 1.137.0, `nat_type` cannot be set to `Normal`.
* `vswitch_id` - (Optional, Available since v1.102.0) The id of VSwitch.
* `internet_charge_type` - (Optional, ForceNew, Available since v1.121.0) The internet charge type. Valid values `PayByLcu`. The `PayByLcu` is only support enhanced NAT. **NOTE:** From version 1.137.0, `internet_charge_type` cannot be set to `PayBySpec`.
* `tags` - (Optional, Available since v1.121.0) The tags of NAT gateway.
* `deletion_protection` - (Optional, Available since v1.124.4) Whether enable the deletion protection or not. Default value: `false`.
  - true: Enable deletion protection.
  - false: Disable deletion protection.
* `network_type` - (Optional, Available since v1.136.0) Indicates the type of the created NAT gateway. Valid values `internet` and `intranet`. `internet`: Internet NAT Gateway. `intranet`: VPC NAT Gateway.
* `eip_bind_mode` - (Optional, Available since v1.184.0) The EIP binding mode of the NAT gateway. Default value: `MULTI_BINDED`. Valid values:
  - `MULTI_BINDED`: Multi EIP network card visible mode.
  - `NAT`: EIP normal mode, compatible with IPv4 gateway.
* `icmp_reply_enabled` - (Optional, Bool, Available since v1.235.0) Specifies whether to enable ICMP retrieval. Default value: `true`. Valid values:
  - `true`: Enable.
  - `false`: Disable.
* `private_link_enabled` - (Optional, ForceNew, Bool, Available since v1.235.0) Specifies whether to enable PrivateLink. Default value: `false`. Valid values:
  - `true`: Enable.
  - `false`: Disable.
* `access_mode` - (Optional, ForceNew, Set, Available since v1.235.0) The access mode for reverse access to the VPC NAT gateway. See [`access_mode`](#access_mode) below.
* `name` - (Optional, ForceNew, Deprecated since v1.121.0) Field `name` has been deprecated from provider version 1.121.0. New field `nat_gateway_name` instead.
* `instance_charge_type` - (Optional, ForceNew, Deprecated since v1.121.0) Field `instance_charge_type` has been deprecated from provider version 1.121.0. New field `payment_type` instead.
* `spec` - (Removed since v1.121.0) The specification of the nat gateway. **NOTE:** Field `spec` has been deprecated from provider version 1.7.1, and it has been removed from provider version 1.121.0. New field `specification` instead.
* `bandwidth_package_ids` - (Removed since v1.121.0) The ID of the bandwidth package. **NOTE:** Field `bandwidth_package_ids` has been removed from provider version 1.121.0.
* `bandwidth_packages` - (Removed since v1.121.0) A list of bandwidth packages for the nat gatway. See [`bandwidth_packages`](#bandwidth_packages) below.

-> **NOTE:** Field `bandwidth_packages` has been removed from provider version 1.121.0.

-> **NOTE:** From version 1.194.0, `eip_bind_mode` can be modified. If the `eip_bind_mode` parameter is set to `MULTI_BINDED` when the NAT gateway is created, you can change the value of this parameter from `MULTI_BINDED` to `NAT`. If the `eip_bind_mode` parameter is set to `NAT` when the NAT gateway is created, you cannot change the value of this parameter from `NAT` to `MULTI_BINDED`.

-> **NOTE:** The `Normal` Nat Gateway has been offline and please using `Enhanced` Nat Gateway to get the better performance.

### `access_mode`

The access_mode supports the following:

* `mode_value` - (Optional, ForceNew) The mode of Access. Valid values:
  - `route`: Route mode.
  - `tunnel`: Tunnel mode.
**NOTE:** If `mode_value` is specified, `private_link_enabled` must be set to `true`.
* `tunnel_type` - (Optional, ForceNew) The type of Tunnel. Valid values: `geneve`. **NOTE:** `tunnel_type` takes effect only if `mode_value` is set to `tunnel`.

### `bandwidth_packages`

The bandwidth_packages mapping supports the following:

* `ip_count` - (Removed since v1.121.0) The IP number of the current bandwidth package. **NOTE:** Field `ip_count` has been removed from provider version 1.121.0.
* `bandwidth` - (Removed since v1.121.0) The bandwidth value of the current bandwidth package. **NOTE:** Field `bandwidth` has been removed from provider version 1.121.0.
* `zone` - (Removed since v1.121.0) The AZ for the current bandwidth. **NOTE:** Field `zone` has been removed from provider version 1.121.0.
* `public_ip_addresses` - (Removed since v1.121.0) The public ip for bandwidth package. **NOTE:** Field `public_ip_addresses` has been removed from provider version 1.121.0.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the nat gateway.
* `snat_table_ids` - The nat gateway will auto create a snat item.
* `forward_table_ids` - The nat gateway will auto create a forward item.
* `status` - (Available since v1.121.0) The status of NAT gateway.

## Timeouts

-> **NOTE:** Available since v1.121.0.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 10 mins) Used when create the nat gateway.
* `update` - (Defaults to 10 mins) Used when update the nat gateway.
* `delete` - (Defaults to 10 mins) Used when delete the nat gateway.

## Import

Nat gateway can be imported using the id, e.g.

```shell
$ terraform import alicloud_nat_gateway.example <id>
```
