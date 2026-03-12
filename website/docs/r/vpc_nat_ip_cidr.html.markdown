---
subcategory: "NAT Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_nat_ip_cidr"
description: |-
  Provides a Alicloud Nat Gateway Nat Ip Cidr resource.
---

# alicloud_vpc_nat_ip_cidr

Provides a Nat Gateway Nat Ip Cidr resource.

NAT IP address segment.

For information about Nat Gateway Nat Ip Cidr and how to use it, see [What is Nat Ip Cidr](https://www.alibabacloud.com/help/doc-detail/281972.htm).

-> **NOTE:** Available since v1.136.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_vpc_nat_ip_cidr&exampleId=6049371f-2af9-3c61-f2a3-b369cb923f00ad4678d9&activeTab=example&spm=docs.r.vpc_nat_ip_cidr.0.6049371f2a&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  vpc_name   = "terraform-example"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/21"
  zone_id      = data.alicloud_zones.default.zones.0.id
  vswitch_name = "terraform-example"
}

resource "alicloud_nat_gateway" "default" {
  vpc_id               = alicloud_vpc.default.id
  internet_charge_type = "PayByLcu"
  nat_gateway_name     = "terraform-example"
  description          = "terraform-example"
  nat_type             = "Enhanced"
  vswitch_id           = alicloud_vswitch.default.id
  network_type         = "intranet"
}

resource "alicloud_vpc_nat_ip_cidr" "default" {
  nat_gateway_id   = alicloud_nat_gateway.default.id
  nat_ip_cidr_name = "terraform-example"
  nat_ip_cidr      = "192.168.0.0/16"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_vpc_nat_ip_cidr&spm=docs.r.vpc_nat_ip_cidr.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `dry_run` - (Optional) Specifies whether to only precheck this request. Valid values:
  - `true`: Sends a dry-run request without creating the NAT IP CIDR block. The system checks whether required parameters are specified, whether the request format is valid, and whether the request complies with service limits. If the check fails, an error is returned. If the check passes, the error code `DryRunOperation` is returned.
  - `false` (default): Sends a normal request. If the check passes, an HTTP 2xx status code is returned and the operation is performed.* `nat_gateway_id` - (Required, ForceNew) The ID of the Virtual Private Cloud (VPC) NAT gateway where you want to create the NAT CIDR block.

-> **NOTE:** This parameter only applies during resource creation, update or deletion. If modified in isolation without other property changes, Terraform will not trigger any action.

* `nat_gateway_id` - (Required, ForceNew) The ID of the VPC NAT gateway instance to which the NAT IP address block belongs.
* `nat_ip_cidr` - (Required, ForceNew) The NAT IP CIDR block to create.

The newly created CIDR block must meet the following requirements:
- It must belong to the 10.0.0.0/8, 172.16.0.0/12, or 192.168.0.0/16 CIDR blocks or their subnets.
- The subnet mask length must be between 16 and 32 bits.
- It must not overlap with the private CIDR block of the VPC to which the VPC NAT gateway belongs. If you need to translate a private IP address to another address within the VPC's private CIDR block, create a vSwitch in the corresponding VPC private CIDR block and then create a new VPC NAT gateway in that vSwitch to provide private address translation.
- If you want to use a public CIDR block as the NAT IP CIDR block, the CIDR block must belong to the customer CIDR block of the VPC to which the VPC NAT gateway belongs. For more information about customer CIDR blocks, see [What is a customer CIDR block?](https://help.aliyun.com/document_detail/185311.html).
* `nat_ip_cidr_description` - (Optional) The description of the NAT IP CIDR block to modify.
  The description must be 2 to 256 characters in length, start with a letter or Chinese character, and cannot start with `http://` or `https://`.
* `nat_ip_cidr_name` - (Required) The name of the NAT IP address block.
  The name must be 2 to 128 characters in length, and can contain letters, digits, periods (.), underscores (_), and hyphens (-). It must start with a letter or a Chinese character, and cannot start with `http://` or `https://`.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. The value is formulated as `<nat_gateway_id>:<nat_ip_cidr>`.
* `create_time` - (Available since v1.273.0) The time when the NAT IP CIDR block was created.
* `status` - The status of the NAT IP CIDR block to query.

## Timeouts

-> **NOTE:** Available since v1.273.0.

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Nat Ip Cidr.
* `delete` - (Defaults to 5 mins) Used when delete the Nat Ip Cidr.
* `update` - (Defaults to 5 mins) Used when update the Nat Ip Cidr.

## Import

Nat Gateway Nat Ip Cidr can be imported using the id, e.g.

```shell
$ terraform import alicloud_vpc_nat_ip_cidr.example <nat_gateway_id>:<nat_ip_cidr>
```
