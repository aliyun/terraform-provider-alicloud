---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_nat_gateway"
sidebar_current: "docs-alicloud-resource-nat-gateway"
description: |-
  Provides a resource to create a VPC NAT Gateway.
---

# alicloud\_nat\_gateway

Provides a resource to create a VPC NAT Gateway.


-> **NOTE:** Resource bandwidth packages will not be supported since 00:00 on November 4, 2017, and public IP can be replaced be elastic IPs.
If a Nat Gateway has already bought some bandwidth packages, it can not bind elastic IP and you have to submit the [work order](https://selfservice.console.aliyun.com/ticket/createIndex) to solve.
If you want to add public IP, you can use resource 'alicloud_eip_association' to bind several elastic IPs for one Nat Gateway.

-> **NOTE:** From version 1.7.1, this resource has deprecated bandwidth packages.
But, in order to manage stock bandwidth packages, version 1.13.0 re-support configuring 'bandwidth_packages'.
    

## Example Usage

Basic usage

- create enhanced nat gateway
```terraform
variable "name" {
  default = "natGatewayExampleName"
}

resource "alicloud_vpc" "enhanced" {
  vpc_name   = var.name
  cidr_block = "10.0.0.0/8"
}

data "alicloud_enhanced_nat_available_zones" "enhanced"{
}

resource "alicloud_vswitch" "enhanced" {
  vswitch_name      = var.name
  zone_id = data.alicloud_enhanced_nat_available_zones.enhanced.zones.0.zone_id
  cidr_block        = "10.10.0.0/20"
  vpc_id            = alicloud_vpc.enhanced.id
}

resource "alicloud_nat_gateway" "enhanced" {
  depends_on           = [alicloud_vswitch.enhanced]
  vpc_id               = alicloud_vpc.enhanced.id
  specification        = "Small"
  nat_gateway_name     = var.name
  payment_type         = "PayAsYouGo"
  vswitch_id           = alicloud_vswitch.enhanced.id
  nat_type             = "Enhanced"
}
```

- transform nat from Normal to Enhanced
-> **NOTE:** You must set `nat_type` to `Enhanced` and set `vswitch_id`.
```
variable "name" {
  default = "nat-transform-to-enhanced"
}

data "alicloud_enhanced_nat_available_zones" "enhanced" {
}

resource "alicloud_vpc" "foo" {
  vpc_name   = var.name
  cidr_block = "10.0.0.0/8"
}

resource "alicloud_vswitch" "foo1" {
  depends_on        = [alicloud_vpc.foo]
  vswitch_name      = var.name
  zone_id           = data.alicloud_enhanced_nat_available_zones.enhanced.zones[1].zone_id
  cidr_block        = "10.10.0.0/20"
  vpc_id            = alicloud_vpc.foo.id
}

resource "alicloud_nat_gateway" "main" {
  depends_on           = [alicloud_vpc.foo,alicloud_vswitch.foo1]
  vpc_id               = alicloud_vpc.foo.id
  specification        = "Small"
  nat_gateway_name     = var.name
  nat_type             = "Enhanced"
  vswitch_id           = alicloud_vswitch.foo1.id
}
```

## Argument Reference

The following arguments are supported:

* `vpc_id` - (Required, ForceNew) The VPC ID.
* `spec` - (Deprecated, Remove from v1.121.0) It has been deprecated from provider version 1.7.1, and new field 'specification' can replace it.
* `specification` - (Optional) The specification of the nat gateway. Valid values are `Small`, `Middle` and `Large`. Default to `Small`. Effective when `internet_charge_type` is `PayBySpec`. Details refer to [Nat Gateway Specification](https://www.alibabacloud.com/help/doc-detail/42757.htm).
* `name` - (Optional,  Deprecated from v1.121.0+) Field `name` has been deprecated from provider version 1.121.0. New field `nat_gateway_name` instead.
* `nat_gateway_name` - (Optional, Available in 1.121.0+) Name of the nat gateway. The value can have a string of 2 to 128 characters, must contain only alphanumeric characters or hyphens, such as "-",".","_", and must not begin or end with a hyphen, and must not begin with http:// or https://. Defaults to null.
* `description` - (Optional) Description of the nat gateway, This description can have a string of 2 to 256 characters, It cannot begin with http:// or https://. Defaults to null.
* `bandwidth_packages` - (Optional, Remove from v1.121.0) A list of bandwidth packages for the nat gatway. Only support nat gateway created before 00:00 on November 4, 2017. Available in v1.13.0+ and v1.7.1-.
* `instance_charge_type` - (Optional, ForceNew,  Deprecated from v1.121.0+) Field `instance_charge_type` has been deprecated from provider version 1.121.0. New field `payment_type` instead.
* `payment_type` - (Optional, ForceNew, Available in 1.121.0+) The billing method of the NAT gateway. Valid values are `PayAsYouGo` and `Subscription`. Default to `PayAsYouGo`.
* `period` - (Optional, ForceNew, Available in 1.45.0+) The duration that you will buy the resource, in month. It is valid when `payment_type` is `Subscription`. Default to 1. Valid values: [1-9, 12, 24, 36]. At present, the provider does not support modify "period" and you can do that via web console.
* `nat_type` - (Optional, Available in 1.102.0+) The type of NAT gateway. Default to `Normal`. Valid values: [`Normal`, `Enhanced`].
* `vswitch_id` - (Optional, Available in 1.102.0+) The id of VSwitch.
* `internet_charge_type` - (Optional, ForceNew, Available in 1.121.0+) The internet charge type. Valid values `PayByLcu` and `PayBySpec`, default value is `PayBySpec`. The `PayByLcu` is only support enhanced NAT.
* `tags` - (Optional, Available in 1.121.0+) The tags of NAT gateway.

-> **NOTE:** The `Normal` Nat Gateway has been offline and please using `Enhanced` Nat Gateway to get the better performance. 

## Block bandwidth packages
The bandwidth package mapping supports the following:

* `ip_count` - (Required) The IP number of the current bandwidth package. Its value range from 1 to 50.
* `bandwidth` - (Required) The bandwidth value of the current bandwidth package. Its value range from 5 to 5000.
* `zone` - (Optional) The AZ for the current bandwidth. If this value is not specified, Terraform will set a random AZ.
* `public_ip_addresses` - (Computer) The public ip for bandwidth package. the public ip count equal `ip_count`, multi ip would complex with ",", such as "10.0.0.1,10.0.0.2".

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the nat gateway.
* `name` - The name of the nat gateway.
* `description` - The description of the nat gateway.
* `spec` - It has been deprecated from provider version 1.7.1.
* `specification` - The specification of the nat gateway.
* `vpc_id` - The VPC ID for the nat gateway.
* `bandwidth_package_ids` - (Remove from v1.121.0) A list ID of the bandwidth packages, and split them with commas.
* `snat_table_ids` - The nat gateway will auto create a snap and forward item, the `snat_table_ids` is the created one. It is format a list after v1.121.0+.
* `forward_table_ids` - The nat gateway will auto create a snap and forward item, the `forward_table_ids` is the created one. It is format a list after v1.121.0+.
* `nat_type` - The type of the nat gateway.
* `vswitch_id` - The ID of the VSwitch, if the `nat_type` is `Enhanced`, it will not be none. 
* `status` - (Available in 1.121.0+) The status of NAT gateway.

#### Timeouts

-> **NOTE:** Available in 1.121.0+.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 10 mins) Used when create the nat gateway.
* `update` - (Defaults to 10 mins) Used when update the nat gateway.
* `delete` - (Defaults to 10 mins) Used when delete the nat gateway.

## Import

Nat gateway can be imported using the id, e.g.

```
$ terraform import alicloud_nat_gateway.example ngw-abc123456
```
