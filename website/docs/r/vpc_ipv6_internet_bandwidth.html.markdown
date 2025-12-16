---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_ipv6_internet_bandwidth"
description: |-
  Provides a Alicloud VPC Ipv6 Internet Bandwidth resource.
---

# alicloud_vpc_ipv6_internet_bandwidth

Provides a VPC Ipv6 Internet Bandwidth resource. Public network bandwidth of IPv6 address.

For information about VPC Ipv6 Internet Bandwidth and how to use it, see [What is Ipv6 Internet Bandwidth](https://www.alibabacloud.com/help/en/virtual-private-cloud/latest/allocateipv6internetbandwidth).

-> **NOTE:** Available since v1.143.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_vpc_ipv6_internet_bandwidth&exampleId=25df92c0-39fb-a6a9-b05a-8728773145e8732c3e14&activeTab=example&spm=docs.r.vpc_ipv6_internet_bandwidth.0.25df92c039&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_zones" "default" {
}

resource "alicloud_vpc" "default" {
  vpc_name    = var.name
  enable_ipv6 = "true"
  cidr_block  = "172.16.0.0/12"
}

resource "alicloud_vswitch" "vsw" {
  vpc_id               = alicloud_vpc.default.id
  cidr_block           = "172.16.0.0/21"
  zone_id              = data.alicloud_zones.default.zones.0.id
  vswitch_name         = var.name
  ipv6_cidr_block_mask = "22"
}

resource "alicloud_security_group" "group" {
  name        = var.name
  description = "foo"
  vpc_id      = alicloud_vpc.default.id
}

data "alicloud_instance_types" "default" {
  availability_zone                 = data.alicloud_zones.default.zones.0.id
  system_disk_category              = "cloud_efficiency"
  cpu_core_count                    = 4
  minimum_eni_ipv6_address_quantity = 1
}

data "alicloud_images" "default" {
  name_regex  = "^ubuntu_18.*64"
  most_recent = true
  owners      = "system"
}

resource "alicloud_instance" "vpc_instance" {
  availability_zone          = data.alicloud_zones.default.zones.0.id
  ipv6_address_count         = 1
  instance_type              = data.alicloud_instance_types.default.instance_types.0.id
  system_disk_category       = "cloud_efficiency"
  image_id                   = data.alicloud_images.default.images.0.id
  instance_name              = var.name
  vswitch_id                 = alicloud_vswitch.vsw.id
  internet_max_bandwidth_out = 10
  security_groups            = alicloud_security_group.group.*.id
}

resource "alicloud_vpc_ipv6_gateway" "example" {
  ipv6_gateway_name = "example_value"
  vpc_id            = alicloud_vpc.default.id
}

data "alicloud_vpc_ipv6_addresses" "default" {
  associated_instance_id = alicloud_instance.vpc_instance.id
  status                 = "Available"
}

resource "alicloud_vpc_ipv6_internet_bandwidth" "example" {
  ipv6_address_id      = data.alicloud_vpc_ipv6_addresses.default.addresses.0.id
  ipv6_gateway_id      = alicloud_vpc_ipv6_gateway.example.ipv6_gateway_id
  internet_charge_type = "PayByBandwidth"
  bandwidth            = "20"
}
```


📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_vpc_ipv6_internet_bandwidth&spm=docs.r.vpc_ipv6_internet_bandwidth.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `bandwidth` - (Required) The amount of Internet bandwidth resources of the IPv6 address, Unit: `Mbit/s`. Valid values: `1` to `5000`. **NOTE:** If `internet_charge_type` is set to `PayByTraffic`, the amount of Internet bandwidth resources of the IPv6 address is limited by the specification of the IPv6 gateway. `Small` (default): specifies the Free edition and the Internet bandwidth is from `1` to `500` Mbit/s. `Medium`: specifies the Medium edition and the Internet bandwidth is from `1` to `1000` Mbit/s. `Large`: specifies the Large edition and the Internet bandwidth is from `1` to `2000` Mbit/s.
* `internet_charge_type` - (Optional, ForceNew, Computed) The metering method of the Internet bandwidth resources of the IPv6 gateway. Valid values: `PayByBandwidth`, `PayByTraffic`.
* `ipv6_address_id` - (Required) The ID of the IPv6 address instance.
* `ipv6_gateway_id` - (Required, ForceNew) The ID of the IPv6 gateway to which the IPv6 address belongs.



## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `status` - The status of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Ipv6 Internet Bandwidth.
* `delete` - (Defaults to 5 mins) Used when delete the Ipv6 Internet Bandwidth.
* `update` - (Defaults to 5 mins) Used when update the Ipv6 Internet Bandwidth.

## Import

VPC Ipv6 Internet Bandwidth can be imported using the id, e.g.

```shell
$ terraform import alicloud_vpc_ipv6_internet_bandwidth.example <id>
```