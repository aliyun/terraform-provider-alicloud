---
subcategory: "Private Zone"
layout: "alicloud"
page_title: "Alicloud: alicloud_pvtz_zone_attachment"
sidebar_current: "docs-alicloud-resource-pvtz-zone-attachment"
description: |-
  Provides vpcs bound to Alicloud Private Zone resource.
---

# alicloud_pvtz_zone_attachment

Provides vpcs bound to Alicloud Private Zone resource.

-> **NOTE:** Available since v1.13.0+.

-> **NOTE:** Terraform will auto bind vpc to a Private Zone while it uses `alicloud_pvtz_zone_attachment` to build a Private Zone and VPC binding resource.

## Example Usage

Using `vpc_ids` to attach being in same region several vpc instances to a private zone

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_pvtz_zone_attachment&exampleId=02f04c34-15d7-20e1-1003-a4e6d3e74d3fc6ac4d63&activeTab=example&spm=docs.r.pvtz_zone_attachment.0.02f04c3415&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_pvtz_zone" "zone" {
  zone_name = "foo.example.com"
}

resource "alicloud_vpc" "first" {
  vpc_name   = "the-first-vpc"
  cidr_block = "172.16.0.0/12"
}
resource "alicloud_vpc" "second" {
  vpc_name   = "the-second-vpc"
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_pvtz_zone_attachment" "zone-attachment" {
  zone_id = alicloud_pvtz_zone.zone.id
  vpc_ids = [alicloud_vpc.first.id, alicloud_vpc.second.id]
}
```

Using `vpcs` to attach being in same region several vpc instances to a private zone

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_pvtz_zone_attachment&exampleId=3bd0e378-5295-5e66-2129-51bd7780f5511d3f6ebd&activeTab=example&spm=docs.r.pvtz_zone_attachment.1.3bd0e37852&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_pvtz_zone" "zone" {
  zone_name = "foo.example.com"
}

resource "alicloud_vpc" "first" {
  vpc_name   = "the-first-vpc"
  cidr_block = "172.16.0.0/12"
}
resource "alicloud_vpc" "second" {
  vpc_name   = "the-second-vpc"
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_pvtz_zone_attachment" "zone-attachment" {
  zone_id = alicloud_pvtz_zone.zone.id
  vpcs {
    vpc_id = alicloud_vpc.first.id
  }
  vpcs {
    vpc_id = alicloud_vpc.second.id
  }
}
```

Using `vpcs` to attach being in different regions several vpc instances to a private zone


<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_pvtz_zone_attachment&exampleId=fa7a7da7-a003-aa6f-1e19-5fd998e8b3727ed6e2ef&activeTab=example&spm=docs.r.pvtz_zone_attachment.2.fa7a7da7a0&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_pvtz_zone" "zone" {
  zone_name = "foo.example.com"
}

resource "alicloud_vpc" "first" {
  vpc_name   = "the-first-vpc"
  cidr_block = "172.16.0.0/12"
}
resource "alicloud_vpc" "second" {
  vpc_name   = "the-second-vpc"
  cidr_block = "172.16.0.0/16"
}

provider "alicloud" {
  alias  = "eu"
  region = "eu-central-1"
}

resource "alicloud_vpc" "third" {
  provider   = alicloud.eu
  vpc_name   = "the-third-vpc"
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_pvtz_zone_attachment" "zone-attachment" {
  zone_id = alicloud_pvtz_zone.zone.id
  vpcs {
    vpc_id = alicloud_vpc.first.id
  }
  vpcs {
    vpc_id = alicloud_vpc.second.id
  }
  vpcs {
    region_id = "eu-central-1"
    vpc_id    = alicloud_vpc.third.id
  }
}
```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Required, ForceNew) The name of the Private Zone Record.
* `vpc_ids` - (Optional, Conflict with `vpcs`) The id List of the VPC with the same region, for example:["vpc-1","vpc-2"]. 
* `vpcs` - (Optional, Conflict with `vpc_ids`, Available in 1.62.1+) See [`vpcs`](#vpcs) below.Recommend to use `vpcs`.
* `lang` - (Optional, ForceNew, Available in 1.62.1+) The language of code.
* `user_client_ip` - (Optional, ForceNew, Available in 1.62.1+) The user custom IP address.

### `vpcs`
The vpcs support the following:
* `vpc_id` - (Required) The Id of the vpc.
* `region_id` - (Optional) The region of the vpc. If not set, the current region will instead of.

## Timeouts

-> **NOTE:** Available in 1.110.0+.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when creating the Private Zone Attachment.
* `update` - (Defaults to 5 mins) Used when updating the Private Zone Attachment.
* `delete` - (Defaults to 5 mins) Used when terminating the Private Zone Attachment. 

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Private Zone VPC Attachment. It sames with `zone_id`.

## Import

Private Zone attachment can be imported using the id(same with `zone_id`), e.g.

```shell
$ terraform import alicloud_pvtz_zone_attachment.example abc123456
```
