---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vswitch"
sidebar_current: "docs-alicloud-resource-vswitch"
description: |-
  Provides a Alicloud VPC Vswitch resource.
---

# alicloud_vswitch

Provides a VPC Vswitch resource. ## Module Support

You can use to the existing [vpc module](https://registry.terraform.io/modules/alibaba/vpc/alicloud)  to create a VPC and several VSwitches one-click.

For information about VPC Vswitch and how to use it, see [What is Vswitch](https://www.alibabacloud.com/help/en/virtual-private-cloud/latest/work-with-vswitches).

-> **NOTE:** Available since v1.0.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_vswitch&exampleId=079cdd19-b6f0-13c0-7a5d-4d84f66e1832cc04102f&activeTab=example&spm=docs.r.vswitch.0.079cdd19b6&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform

data "alicloud_zones" "foo" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "foo" {
  vpc_name   = "terraform-example"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "foo" {
  vswitch_name = "terraform-example"
  cidr_block   = "172.16.0.0/21"
  vpc_id       = alicloud_vpc.foo.id
  zone_id      = data.alicloud_zones.foo.zones.0.id
}
```

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_vswitch&exampleId=1a612892-3673-e700-0582-e5ae4685652a5322921c&activeTab=example&spm=docs.r.vswitch.1.1a61289236&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
data "alicloud_zones" "foo" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "vpc" {
  vpc_name   = "terraform-example"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vpc_ipv4_cidr_block" "cidr_blocks" {
  vpc_id               = alicloud_vpc.vpc.id
  secondary_cidr_block = "192.163.0.0/16"
}

resource "alicloud_vswitch" "island-nat" {
  vpc_id       = alicloud_vpc_ipv4_cidr_block.cidr_blocks.vpc_id
  cidr_block   = "172.16.0.0/21"
  zone_id      = data.alicloud_zones.foo.zones.0.id
  vswitch_name = "terraform-example"
  tags = {
    BuiltBy     = "example_value"
    cnm_version = "example_value"
    Environment = "example_value"
    ManagedBy   = "example_value"
  }
}
```

Create a switch associated with the additional network segment

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_vswitch&exampleId=30815915-965a-c6c1-bc30-cdf2de837f7b514f99f4&activeTab=example&spm=docs.r.vswitch.2.3081591596&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
data "alicloud_zones" "foo" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "foo" {
  vpc_name   = "terraform-example"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vpc_ipv4_cidr_block" "foo" {
  vpc_id               = alicloud_vpc.foo.id
  secondary_cidr_block = "192.163.0.0/16"
}

resource "alicloud_vswitch" "foo" {
  vpc_id     = alicloud_vpc_ipv4_cidr_block.foo.vpc_id
  cidr_block = "192.163.0.0/24"
  zone_id    = data.alicloud_zones.foo.zones.0.id
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_vswitch&spm=docs.r.vswitch.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `cidr_block` - (Optional, ForceNew) The IPv4 CIDR block of the VSwitch. **NOTE:** From version 1.233.0, if you do not set `is_default`, or set `is_default` to `false`, `cidr_block` is required.
* `description` - (Optional) The description of VSwitch.
* `zone_id` - (Optional, ForceNew, Available since v1.119.0) The AZ for the VSwitch. **Note:** Required for a VPC VSwitch.
* `enable_ipv6` - (Optional, Available since v1.201.0) Whether the IPv6 function is enabled in the switch. Value:
  - `true`: enables IPv6.
  - `false` (default): IPv6 is not enabled.
* `ipv6_cidr_block_mask` - (Optional, Available since v1.201.0) The IPv6 CIDR block of the VSwitch.
* `tags` - (Optional, Map, Available since v1.55.3) The tags of VSwitch.
* `vswitch_name` - (Optional, Available since v1.119.0) The name of the VSwitch.
* `vpc_id` - (Optional, ForceNew) The VPC ID. **NOTE:** From version 1.233.0, if you do not set `is_default`, or set `is_default` to `false`, `vpc_id` is required.
* `is_default` - (Optional, Bool, Available since v1.233.0) Specifies whether to create the default VSwitch. Default value: `false`. Valid values:
  - `true`: Creates a default vSwitch.
  - `false`: Creates a vSwitch.

The following arguments will be discarded. Please use new fields as soon as possible:
* `name` - (Deprecated since v1.119.0) Field `name` has been deprecated from provider version 1.119.0. New field `vswitch_name` instead.
* `availability_zone` - (Deprecated since v1.119.0) Field `availability_zone` has been deprecated from provider version 1.119.0. New field `zone_id` instead.


## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the VSwitch.
* `ipv6_cidr_block` - The IPv6 CIDR block of the VSwitch.
* `status` - The status of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Vswitch.
* `delete` - (Defaults to 5 mins) Used when delete the Vswitch.
* `update` - (Defaults to 5 mins) Used when update the Vswitch.

## Import

VPC Vswitch can be imported using the id, e.g.

```shell
$ terraform import alicloud_vpc_vswitch.example <id>
```
