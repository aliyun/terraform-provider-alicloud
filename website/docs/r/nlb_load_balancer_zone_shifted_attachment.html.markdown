---
subcategory: "Network Load Balancer (NLB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_nlb_load_balancer_zone_shifted_attachment"
description: |-
  Provides a Alicloud Network Load Balancer (NLB) Load Balancer Zone Shifted Attachment resource.
---

# alicloud_nlb_load_balancer_zone_shifted_attachment

Provides a Network Load Balancer (NLB) Load Balancer Zone Shifted Attachment resource.

Network load balancer shift zone.

For information about Network Load Balancer (NLB) Load Balancer Zone Shifted Attachment and how to use it, see [What is Load Balancer Zone Shifted Attachment](https://next.api.alibabacloud.com/document/Nlb/2022-04-30/StartShiftLoadBalancerZones).

-> **NOTE:** Available since v1.245.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_nlb_load_balancer_zone_shifted_attachment&exampleId=74c39d55-4bf3-9b2d-89c3-2306c1ae804ebb56b333&activeTab=example&spm=docs.r.nlb_load_balancer_zone_shifted_attachment.0.74c39d554b&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-beijing"
}

resource "alicloud_vpc" "vpc" {
  description = "example"
  cidr_block  = "10.0.0.0/8"
  enable_ipv6 = true
  vpc_name    = "tf-exampleacc-248"
}

resource "alicloud_vswitch" "vsw1" {
  vpc_id       = alicloud_vpc.vpc.id
  zone_id      = "cn-beijing-l"
  cidr_block   = "10.0.1.0/24"
  vswitch_name = "tf-exampleacc-41"
}

resource "alicloud_vswitch" "vsw2" {
  vpc_id               = alicloud_vpc.vpc.id
  zone_id              = "cn-beijing-k"
  cidr_block           = "10.0.2.0/24"
  vswitch_name         = "tf-exampleacc-301"
  ipv6_cidr_block_mask = "8"
}

resource "alicloud_nlb_load_balancer" "nlb" {
  zone_mappings {
    vswitch_id = alicloud_vswitch.vsw1.id
    zone_id    = alicloud_vswitch.vsw1.zone_id
  }
  zone_mappings {
    vswitch_id = alicloud_vswitch.vsw2.id
    zone_id    = alicloud_vswitch.vsw2.zone_id
  }
  vpc_id       = alicloud_vpc.vpc.id
  address_type = "Intranet"
}


resource "alicloud_nlb_load_balancer_zone_shifted_attachment" "default" {
  zone_id          = alicloud_vswitch.vsw1.zone_id
  vswitch_id       = alicloud_vswitch.vsw1.id
  load_balancer_id = alicloud_nlb_load_balancer.nlb.id
}
```

## Argument Reference

The following arguments are supported:
* `load_balancer_id` - (Required, ForceNew) Network load balancer id
* `vswitch_id` - (Required, ForceNew) The list of zones and vSwitch mappings
* `zone_id` - (Required, ForceNew) ZoneId

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<load_balancer_id>:<zone_id>:<vswitch_id>`.
* `status` - Status

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Load Balancer Zone Shifted Attachment.
* `delete` - (Defaults to 5 mins) Used when delete the Load Balancer Zone Shifted Attachment.

## Import

Network Load Balancer (NLB) Load Balancer Zone Shifted Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_nlb_load_balancer_zone_shifted_attachment.example <load_balancer_id>:<zone_id>:<vswitch_id>
```