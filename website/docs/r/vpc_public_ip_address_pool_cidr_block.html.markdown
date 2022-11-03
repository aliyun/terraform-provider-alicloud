---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_public_ip_address_pool_cidr_block"
sidebar_current: "docs-alicloud-resource-vpc-public-ip-address-pool-cidr-block"
description: |-
  Provides a Alicloud VPC Public Ip Address Pool Cidr Block resource.
---

# alicloud\_vpc\_public\_ip\_address\_pool\_cidr\_block

Provides a VPC Public Ip Address Pool Cidr Block resource.

For information about VPC Public Ip Address Pool Cidr Block and how to use it, see [What is Public Ip Address Pool Cidr Block](https://www.alibabacloud.com/help/en/virtual-private-cloud/latest/429100).

-> **NOTE:** Available in v1.189.0+.

-> **NOTE:** Only users who have the required permissions can use the IP address pool feature of Elastic IP Address (EIP). To apply for the required permissions, [submit a ticket](https://smartservice.console.aliyun.com/service/create-ticket).

## Example Usage

Basic Usage

```terraform
resource "alicloud_vpc_public_ip_address_pool" "default" {
  public_ip_address_pool_name = "example_value"
  isp                         = "BGP"
  description                 = "example_value"
}

resource "alicloud_vpc_public_ip_address_pool_cidr_block" "default" {
  public_ip_address_pool_id = alicloud_vpc_public_ip_address_pool.default.id
  cidr_block                = "your_cidr_block"
}
```

## Argument Reference

The following arguments are supported:

* `public_ip_address_pool_id` - (Required, ForceNew) The ID of the VPC Public IP address pool.
* `cidr_block` - (Required, ForceNew) The CIDR block.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of VPC Public Ip Address Pool Cidr Block. The value formats as `<public_ip_address_pool_id>:<cidr_block>`.
* `status` - The status of the VPC Public Ip Address Pool Cidr Block.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 3 mins) Used when create the VPC Public Ip Address Pool Cidr Block.
* `delete` - (Defaults to 3 mins) Used when delete the VPC Public Ip Address Pool Cidr Block.

## Import

VPC Public Ip Address Pool Cidr Block can be imported using the id, e.g.

```shell
$ terraform import alicloud_vpc_public_ip_address_pool_cidr_block.example <public_ip_address_pool_id>:<cidr_block>
```
