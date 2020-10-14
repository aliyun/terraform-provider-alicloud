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
- create notmal nat gateway
```terraform
variable "name" {
  default = "natGatewayExampleName"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  name       = var.name
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
  vpc_id            = alicloud_vpc.default.id
  cidr_block        = "172.16.0.0/21"
  availability_zone = data.alicloud_zones.default.zones[0].id
  name              = var.name
}

resource "alicloud_nat_gateway" "default" {
  vpc_id = alicloud_vswitch.default.vpc_id
  name   = var.name
}
```

- create enhanced nat gateway
```terraform
variable "name" {
  default = "natGatewayExampleName"
}

resource "alicloud_vpc" "enhanced" {
  name       = var.name
  cidr_block = "10.0.0.0/8"
}

data "alicloud_enhanced_nat_available_zones" "enhanced"{
}

resource "alicloud_vswitch" "enhanced" {
  name              = var.name
  availability_zone = data.alicloud_enhanced_nat_available_zones.enhanced.zones.0.zone_id
  cidr_block        = "10.10.0.0/20"
  vpc_id            = alicloud_vpc.enhanced.id
}

resource "alicloud_nat_gateway" "enhanced" {
  depends_on           = [alicloud_vswitch.enhanced]
  vpc_id               = alicloud_vpc.enhanced.id
  specification        = "Small"
  name                 = var.name
  instance_charge_type = "PostPaid"
  vswitch_id           = alicloud_vswitch.enhanced.id
  nat_type             = "Enhanced"
}
```

## Argument Reference

The following arguments are supported:

* `vpc_id` - (Required, ForceNew) The VPC ID.
* `spec` - (Deprecated) It has been deprecated from provider version 1.7.1, and new field 'specification' can replace it.
* `specification` - (Optional) The specification of the nat gateway. Valid values are `Small`, `Middle` and `Large`. Default to `Small`. Details refer to [Nat Gateway Specification](https://www.alibabacloud.com/help/doc-detail/42757.htm).
* `name` - (Optional) Name of the nat gateway. The value can have a string of 2 to 128 characters, must contain only alphanumeric characters or hyphens, such as "-",".","_", and must not begin or end with a hyphen, and must not begin with http:// or https://. Defaults to null.
* `description` - (Optional) Description of the nat gateway, This description can have a string of 2 to 256 characters, It cannot begin with http:// or https://. Defaults to null.
* `bandwidth_packages` - (Optional) A list of bandwidth packages for the nat gatway. Only support nat gateway created before 00:00 on November 4, 2017. Available in v1.13.0+ and v1.7.1-.
* `instance_charge_type` - (Optional, ForceNew, Available in 1.45.0+) The billing method of the nat gateway. Valid values are "PrePaid" and "PostPaid". Default to "PostPaid".
* `period` - (Optional, ForceNew, Available in 1.45.0+) The duration that you will buy the resource, in month. It is valid when `instance_charge_type` is `PrePaid`. Default to 1. Valid values: [1-9, 12, 24, 36]. At present, the provider does not support modify "period" and you can do that via web console.
* `nat_type` - (Optional, ForceNew, Available in 1.102.0+) The type of nat gateway. Default to `Normal`. Valid values: [`Normal`, `Enhanced`].
* `vswitch_id` - (Optional, ForceNew, Available in 1.102.0+) The id of VSwitch.

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
* `bandwidth_package_ids` - A list ID of the bandwidth packages, and split them with commas.
* `snat_table_ids` - The nat gateway will auto create a snap and forward item, the `snat_table_ids` is the created one.
* `forward_table_ids` - The nat gateway will auto create a snap and forward item, the `forward_table_ids` is the created one.
* `nat_type` - The type of the nat gateway.
* `vswitch_id` - The ID of the VSwitch, if the `nat_type` is `Enhanced`, it will not be none. 

## Import

Nat gateway can be imported using the id, e.g.

```
$ terraform import alicloud_nat_gateway.example ngw-abc123456
```
