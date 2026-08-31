---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_public_ip_address_pool_cidr_blocks"
sidebar_current: "docs-alicloud-datasource-vpc-public-ip-address-pool-cidr-blocks"
description: |-
  Provides a list of Vpc Public Ip Address Pool Cidr Blocks to the user.
---

# alicloud\_vpc\_public\_ip\_address\_pool\_cidr\_blocks

This data source provides the Vpc Public Ip Address Pool Cidr Blocks of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.189.0+.

-> **NOTE:** Only users who have the required permissions can use the IP address pool feature of Elastic IP Address (EIP). To apply for the required permissions, [submit a ticket](https://smartservice.console.aliyun.com/service/create-ticket).

## Example Usage

Basic Usage

```terraform
data "alicloud_vpc_public_ip_address_pool_cidr_blocks" "ids" {
  ids                       = ["example_id"]
  public_ip_address_pool_id = "example_value"
}

output "vpc_public_ip_address_pool_cidr_block_id_1" {
  value = data.alicloud_vpc_public_ip_address_pool_cidr_blocks.ids.blocks.0.id
}

data "alicloud_vpc_public_ip_address_pool_cidr_blocks" "cidrBlock" {
  public_ip_address_pool_id = "example_value"
  cidr_block                = "example_value"
}

output "vpc_public_ip_address_pool_cidr_block_id_2" {
  value = data.alicloud_vpc_public_ip_address_pool_cidr_blocks.cidrBlock.blocks.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Vpc Public Ip Address Pool Cidr Block IDs.
* `public_ip_address_pool_id` - (Required, ForceNew) The ID of the Vpc Public IP address pool.
* `cidr_block` - (Optional, ForceNew) The CIDR block.
* `status` - (Optional, ForceNew) The status of the CIDR block in the Vpc Public IP address pool. Valid values: `Created`, `Modifying`, `Deleting`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `blocks` - A list of Vpc Public Ip Address Pool Cidr Blocks. Each element contains the following attributes:
	* `id` - The ID of the Public Ip Address Pool Cidr Block.
	* `public_ip_address_pool_id` - The ID of the Vpc Public IP address pool.
	* `cidr_block` - The CIDR block.
	* `status` - The status of the CIDR block in the Vpc Public IP address pool.
	* `used_ip_num` - The total number of available IP addresses in the CIDR block.
	* `total_ip_num` - The number of occupied IP addresses in the CIDR block.
	* `create_time` - The time when the CIDR block was created. The time is displayed in YYYY-MM-DDThh:mm:ssZ format.
