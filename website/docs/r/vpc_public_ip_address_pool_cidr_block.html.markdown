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

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_vpc_public_ip_address_pool_cidr_block&exampleId=5e9a0216-17a4-abd0-46b3-65c44987b54835c0346b&activeTab=example&spm=docs.r.vpc_public_ip_address_pool_cidr_block.0.5e9a021617&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}

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

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_vpc_public_ip_address_pool_cidr_block&spm=docs.r.vpc_public_ip_address_pool_cidr_block.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `cidr_block` - (Optional, ForceNew, Computed) The CIDR block.
* `cidr_mask` - (Optional, Available since v1.219.0) IP address and network segment mask. After you enter the mask, the system automatically allocates the IP address network segment. Value range: **24** to **28**.
-> **NOTE:**  **CidrBlock** and **CidrMask** cannot be configured at the same time. Select one of them to configure.
* `public_ip_address_pool_id` - (Required, ForceNew) The ID of the VPC Public IP address pool.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<public_ip_address_pool_id>:<cidr_block>`.
* `create_time` - The creation time of the resource.
* `status` - The status of the VPC Public Ip Address Pool Cidr Block.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Public Ip Address Pool Cidr Block.
* `delete` - (Defaults to 5 mins) Used when delete the Public Ip Address Pool Cidr Block.

## Import

VPC Public Ip Address Pool Cidr Block can be imported using the id, e.g.

```shell
$ terraform import alicloud_vpc_public_ip_address_pool_cidr_block.example <public_ip_address_pool_id>:<cidr_block>
```