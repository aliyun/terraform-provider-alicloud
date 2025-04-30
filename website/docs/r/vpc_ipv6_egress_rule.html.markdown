---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_ipv6_egress_rule"
description: |-
  Provides a Alicloud VPC Ipv6 Egress Rule resource.
---

# alicloud_vpc_ipv6_egress_rule

Provides a VPC Ipv6 Egress Rule resource. IPv6 address addition only active exit rule.

For information about VPC Ipv6 Egress Rule and how to use it, see [What is Ipv6 Egress Rule](https://www.alibabacloud.com/help/doc-detail/102200.htm).

-> **NOTE:** Available since v1.142.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_vpc_ipv6_egress_rule&exampleId=dbe71f0e-6c7b-32cf-29b0-91cc82868a1d9445417b&activeTab=example&spm=docs.r.vpc_ipv6_egress_rule.0.dbe71f0e6c&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
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

resource "alicloud_vpc" "default" {
  vpc_name    = var.name
  enable_ipv6 = "true"
  cidr_block  = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
  vpc_id               = alicloud_vpc.default.id
  cidr_block           = "172.16.0.0/21"
  zone_id              = data.alicloud_zones.default.zones.0.id
  vswitch_name         = var.name
  ipv6_cidr_block_mask = "64"
}

resource "alicloud_security_group" "default" {
  name        = var.name
  description = var.name
  vpc_id      = alicloud_vpc.default.id
}

resource "alicloud_instance" "default" {
  availability_zone          = data.alicloud_zones.default.zones.0.id
  ipv6_address_count         = 1
  instance_type              = data.alicloud_instance_types.default.instance_types.0.id
  system_disk_category       = "cloud_efficiency"
  image_id                   = data.alicloud_images.default.images.0.id
  instance_name              = var.name
  vswitch_id                 = alicloud_vswitch.default.id
  internet_max_bandwidth_out = 10
  security_groups            = [alicloud_security_group.default.id]
}

resource "alicloud_vpc_ipv6_gateway" "default" {
  ipv6_gateway_name = var.name
  vpc_id            = alicloud_vpc.default.id
}

data "alicloud_vpc_ipv6_addresses" "default" {
  associated_instance_id = alicloud_instance.default.id
  status                 = "Available"
}

resource "alicloud_vpc_ipv6_internet_bandwidth" "default" {
  ipv6_address_id      = data.alicloud_vpc_ipv6_addresses.default.addresses.0.id
  ipv6_gateway_id      = alicloud_vpc_ipv6_gateway.default.ipv6_gateway_id
  internet_charge_type = "PayByBandwidth"
  bandwidth            = "20"
}

resource "alicloud_vpc_ipv6_egress_rule" "default" {
  instance_id           = alicloud_vpc_ipv6_internet_bandwidth.default.ipv6_address_id
  ipv6_egress_rule_name = var.name
  description           = var.name
  ipv6_gateway_id       = alicloud_vpc_ipv6_internet_bandwidth.default.ipv6_gateway_id
  instance_type         = "Ipv6Address"
}
```

## Argument Reference

The following arguments are supported:
* `description` - (Optional, ForceNew) The description of the egress-only rule. The description must be `2` to `256` characters in length. It cannot start with `http://` or `https://`.
* `instance_id` - (Required, ForceNew) The ID of the IPv6 address to which you want to apply the egress-only rule.
* `instance_type` - (Optional, ForceNew, Computed) The type of instance to which you want to apply the egress-only rule. Valid values: `Ipv6Address`. `Ipv6Address` (default): an IPv6 address.
* `ipv6_egress_rule_name` - (Optional, ForceNew) The name of the egress-only rule. The name must be `2` to `128` characters in length, and can contain letters, digits, underscores (_), and hyphens (-). The name must start with a letter but cannot start with `http://` or `https://`.
* `ipv6_gateway_id` - (Required, ForceNew) The ID of the IPv6 gateway.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<ipv6_gateway_id>:<ipv6_egress_rule_id>`.
* `ipv6_egress_rule_id` - The ID of the IPv6 EgressRule.
* `status` - The status of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Ipv6 Egress Rule.
* `delete` - (Defaults to 5 mins) Used when delete the Ipv6 Egress Rule.

## Import

VPC Ipv6 Egress Rule can be imported using the id, e.g.

```shell
$ terraform import alicloud_vpc_ipv6_egress_rule.example <ipv6_gateway_id>:<ipv6_egress_rule_id>
```