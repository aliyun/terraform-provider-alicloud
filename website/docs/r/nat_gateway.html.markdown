---
layout: "alicloud"
page_title: "Alicloud: alicloud_nat_gateway"
sidebar_current: "docs-alicloud-resource-nat-gateway"
description: |-
  Provides a resource to create a VPC NAT Gateway.
---

# alicloud\_nat\_gateway

Provides a resource to create a VPC NAT Gateway.

~> **NOTE:** From version 1.7.1, the resource deprecates bandwidth packages.
And if you want to add public IP, you can use resource 'alicloud_eip_association' to bind several elastic IPs for one Nat Gateway.

~> **NOTE:** Resource bandwidth packages will not be supported since 00:00 on November 4, 2017, and public IP can be replaced be elastic IPs.
If a Nat Gateway has already bought some bandwidth packages, it can not bind elastic IP and you have to submit the [work order](https://selfservice.console.aliyun.com/ticket/createIndex) to solve.


## Example Usage

Basic usage

```
resource "alicloud_vpc" "vpc" {
  name       = "tf_test_foo"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "vsw" {
  vpc_id            = "${alicloud_vpc.vpc.id}"
  cidr_block        = "172.16.0.0/21"
  availability_zone = "cn-beijing-b"
}

resource "alicloud_nat_gateway" "nat_gateway" {
  vpc_id = "${alicloud_vpc.vpc.id}"
  spec   = "Small"
  name   = "test_foo"
}
```

## Argument Reference

The following arguments are supported:

* `vpc_id` - (Required, Forces New Resorce) The VPC ID.
* `spec` - (Deprecated) It has been deprecated from provider version 1.7.1, and new field 'specification' can replace it.
* `specification` - (Optional) The specification of the nat gateway. Valid values are `Small`, `Middle` and `Large`. Default to `Small`. Details refer to [Nat Gateway Specification](https://www.alibabacloud.com/help/doc-detail/42757.htm).
* `name` - (Optional) Name of the nat gateway. The value can have a string of 2 to 128 characters, must contain only alphanumeric characters or hyphens, such as "-",".","_", and must not begin or end with a hyphen, and must not begin with http:// or https://. Defaults to null.
* `description` - (Optional) Description of the nat gateway, This description can have a string of 2 to 256 characters, It cannot begin with http:// or https://. Defaults to null.
* `bandwidth_packages` - (Deprecated) It has been deprecated from provider version 1.7.1. Resource 'alicloud_eip_association' can bind several elastic IPs for one Nat Gateway.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the nat gateway.
* `name` - The name of the nat gateway.
* `description` - The description of the nat gateway.
* `spec` - It has been deprecated from provider version 1.7.1.
* `specification` - The specification of the nat gateway.
* `vpc_id` - The VPC ID for the nat gateway.
* `bandwidth_package_ids` - A list ID of the bandwidth packages, and split them with commas
* `snat_table_ids` - The nat gateway will auto create a snap and forward item, the `snat_table_ids` is the created one.
* `forward_table_ids` - The nat gateway will auto create a snap and forward item, the `forward_table_ids` is the created one.

## Import

Nat gateway can be imported using the id, e.g.

```
$ terraform import alicloud_nat_gateway.example ngw-abc123456
```