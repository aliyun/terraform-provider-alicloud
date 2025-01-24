---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_ipv4_cidr_block"
description: |-
  Provides a Alicloud VPC Ipv4 Cidr Block resource.
---

# alicloud_vpc_ipv4_cidr_block

Provides a VPC Ipv4 Cidr Block resource.

VPC IPv4 additional network segment.

For information about VPC Ipv4 Cidr Block and how to use it, see [What is Ipv4 Cidr Block](https://www.alibabacloud.com/help/en/virtual-private-cloud/latest/associatevpccidrblock).

-> **NOTE:** Available since v1.185.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_vpc_ipv4_cidr_block&exampleId=6745eb51-2032-9dde-29a3-d1bee81f4fc03b40e9fb&activeTab=example&spm=docs.r.vpc_ipv4_cidr_block.0.6745eb5120&intl_lang=EN_US" target="_blank">
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

resource "alicloud_vpc" "defaultvpc" {
  description = var.name
}

resource "alicloud_vpc_ipv4_cidr_block" "default" {
  secondary_cidr_block = "192.168.0.0/16"
  vpc_id               = alicloud_vpc.defaultvpc.id
}
```

## Argument Reference

The following arguments are supported:
* `ipv4_ipam_pool_id` - (Optional, Available since v1.242.0) The ID of the IP Address Manager (IPAM) pool that contains IPv4 addresses.
* `secondary_cidr_block` - (Optional, ForceNew, Computed) Additional network segment information.
* `secondary_cidr_mask` - (Optional, Int, Available since v1.242.0) Add an additional CIDR block from the IPAM address pool to the VPC by entering a mask.

-> **NOTE:**  Specify the IPAM address pool to add an additional CIDR block to the VPC. Enter at least one of the SecondaryCidrBlock or SecondaryCidrMask parameters.

* `vpc_id` - (Required, ForceNew) The ID of the VPC.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<vpc_id>:<secondary_cidr_block>`.
* `region_id` - The ID of the region where the VPC resides.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Ipv4 Cidr Block.
* `delete` - (Defaults to 5 mins) Used when delete the Ipv4 Cidr Block.

## Import

VPC Ipv4 Cidr Block can be imported using the id, e.g.

```shell
$ terraform import alicloud_vpc_ipv4_cidr_block.example <vpc_id>:<secondary_cidr_block>
```