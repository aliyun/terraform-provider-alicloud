---
subcategory: "NAT Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_nat_ip"
sidebar_current: "docs-alicloud-resource-vpc-nat-ip"
description: |-
  Provides a Alicloud VPC Nat Ip resource.
---

# alicloud\_vpc\_nat\_ip

Provides a VPC Nat Ip resource.

For information about VPC Nat Ip and how to use it, see [What is Nat Ip](https://www.alibabacloud.com/help/doc-detail/281976.htm).

-> **NOTE:** Available in v1.136.0+.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_vpc_nat_ip&exampleId=4ada155c-b87d-bb11-dc6b-9edd006ac02a6b01dd02&activeTab=example&spm=docs.r.vpc_nat_ip.0.4ada155cb8&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
data "alicloud_zones" "example" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "example" {
  vpc_name   = "terraform-example"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "example" {
  vpc_id       = alicloud_vpc.example.id
  cidr_block   = "172.16.0.0/21"
  zone_id      = data.alicloud_zones.example.zones.0.id
  vswitch_name = "terraform-example"
}

resource "alicloud_nat_gateway" "example" {
  vpc_id               = alicloud_vpc.example.id
  internet_charge_type = "PayByLcu"
  nat_gateway_name     = "terraform-example"
  description          = "terraform-example"
  nat_type             = "Enhanced"
  vswitch_id           = alicloud_vswitch.example.id
  network_type         = "intranet"
}

resource "alicloud_vpc_nat_ip_cidr" "example" {
  nat_ip_cidr             = "192.168.0.0/16"
  nat_gateway_id          = alicloud_nat_gateway.example.id
  nat_ip_cidr_description = "terraform-example"
  nat_ip_cidr_name        = "terraform-example"
}

resource "alicloud_vpc_nat_ip" "example" {
  nat_ip             = "192.168.0.37"
  nat_gateway_id     = alicloud_nat_gateway.example.id
  nat_ip_description = "example_value"
  nat_ip_name        = "example_value"
  nat_ip_cidr        = alicloud_vpc_nat_ip_cidr.example.nat_ip_cidr
}

```

## Argument Reference

The following arguments are supported:

* `dry_run` - (Optional, Computed) Specifies whether to check the validity of the request without actually making the request.
* `nat_gateway_id` - (Required, ForceNew) The ID of the Virtual Private Cloud (VPC) NAT gateway for which you want to create the NAT IP address.
* `nat_ip` - (Optional, ForceNew, Computed) The NAT IP address that you want to create. If you do not specify an IP address, the system selects a random IP address from the specified CIDR block.
* `nat_ip_cidr` - (Optional, ForceNew) NAT IP ADDRESS of the address segment.
* `nat_ip_cidr_id` - (Optional) The ID of the CIDR block to which the NAT IP address belongs.
* `nat_ip_description` - (Optional) NAT IP ADDRESS description of information. Length is from `2` to `256` characters, must start with a letter or the Chinese at the beginning, but not at the` http://` Or `https://` at the beginning.
* `nat_ip_name` - (Optional) NAT IP ADDRESS the name of the root directory. Length is from `2` to `128` characters, must start with a letter or the Chinese at the beginning can contain numbers, half a period (.), underscore (_) and dash (-). But do not start with `http://` or `https://` at the beginning.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Nat Ip. The value formats as `<nat_gateway_id>:<nat_ip_id>`.
* `nat_ip_id` - Ihe ID of the Nat Ip.
* `status` - The status of the NAT IP address. Valid values: `Available`, `Deleting`, `Creating` and `Deleted`. 

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the Nat Ip.
* `delete` - (Defaults to 1 mins) Used when delete the Nat Ip.

## Import

VPC Nat Ip can be imported using the id, e.g.

```shell
$ terraform import alicloud_vpc_nat_ip.example <nat_gateway_id>:<nat_ip_id>
```
