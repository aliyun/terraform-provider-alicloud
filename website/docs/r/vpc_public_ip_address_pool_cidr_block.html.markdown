---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_public_ip_address_pool_cidr_block"
description: |-
  Provides a Alicloud VPC Public Ip Address Pool Cidr Block resource.
---

# alicloud_vpc_public_ip_address_pool_cidr_block

Provides a VPC Public Ip Address Pool Cidr Block resource. 
-> **NOTE:** Only users who have the required permissions can use the IP address pool feature of Elastic IP Address (EIP). To apply for the required permissions, [submit a ticket](https://smartservice.console.aliyun.com/service/create-ticket).

For information about VPC Public Ip Address Pool Cidr Block and how to use it, see [What is Public Ip Address Pool Cidr Block](https://www.alibabacloud.com/help/en/virtual-private-cloud/latest/429100).

-> **NOTE:** Available since v1.189.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-example"
}
data "alicloud_resource_manager_resource_groups" "default" {
  status = "OK"
}
resource "alicloud_vpc_public_ip_address_pool" "default" {
  description                 = var.name
  public_ip_address_pool_name = var.name
  isp                         = "BGP"
  resource_group_id           = data.alicloud_resource_manager_resource_groups.default.ids.0
}

resource "alicloud_vpc_public_ip_address_pool_cidr_block" "default" {
  public_ip_address_pool_id = alicloud_vpc_public_ip_address_pool.default.id
  cidr_block                = "47.118.126.0/25"
}
```


## Argument Reference

The following arguments are supported:
* `cidr_block` - (Optional, ForceNew) The CIDR block.
* `public_ip_address_pool_id` - (Required, ForceNew) The ID of the VPC Public IP address pool.



## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<public_ip_address_pool_id>:<cidr_block>`.
* `create_time` - The creation time of the resource.
* `status` - The status of the VPC Public Ip Address Pool Cidr Block.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Public Ip Address Pool Cidr Block.
* `delete` - (Defaults to 5 mins) Used when delete the Public Ip Address Pool Cidr Block.

## Import

VPC Public Ip Address Pool Cidr Block can be imported using the id, e.g.

```shell
$ terraform import alicloud_vpc_public_ip_address_pool_cidr_block.example <public_ip_address_pool_id>:<cidr_block>
```