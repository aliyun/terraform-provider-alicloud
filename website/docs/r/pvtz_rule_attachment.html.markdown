---
subcategory: "Private Zone"
layout: "alicloud"
page_title: "Alicloud: alicloud_pvtz_rule_attachment"
sidebar_current: "docs-alicloud-resource-pvtz-rule-attachment"
description: |-
  Provides a Alicloud Private Zone Rule Attachment resource.
---

# alicloud_pvtz_rule_attachment

Provides a Private Zone Rule Attachment resource.

For information about Private Zone Rule Attachment and how to use it, see [What is Rule Attachment](https://www.alibabacloud.com/help/en/doc-detail/177601.htm).

-> **NOTE:** Available since v1.143.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_pvtz_rule_attachment&exampleId=0156c84d-c8e4-d634-bfcb-6d42b116d72e83fdaa56&activeTab=example&spm=docs.r.pvtz_rule_attachment.0.0156c84dc8&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "example_value"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

data "alicloud_pvtz_resolver_zones" "default" {
  status = "NORMAL"
}
data "alicloud_regions" "default" {
  current = true
}

resource "alicloud_vpc" "default" {
  count      = 3
  vpc_name   = var.name
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
  count      = 2
  vpc_id     = alicloud_vpc.default[2].id
  cidr_block = cidrsubnet(alicloud_vpc.default[2].cidr_block, 8, count.index)
  zone_id    = data.alicloud_pvtz_resolver_zones.default.zones[count.index].zone_id
}

resource "alicloud_security_group" "default" {
  vpc_id = alicloud_vpc.default[2].id
  name   = var.name
}

resource "alicloud_pvtz_endpoint" "default" {
  endpoint_name     = "${var.name}-${random_integer.default.result}"
  security_group_id = alicloud_security_group.default.id
  vpc_id            = alicloud_vpc.default[2].id
  vpc_region_id     = data.alicloud_regions.default.regions.0.id
  ip_configs {
    zone_id    = alicloud_vswitch.default[0].zone_id
    cidr_block = alicloud_vswitch.default[0].cidr_block
    vswitch_id = alicloud_vswitch.default[0].id
  }
  ip_configs {
    zone_id    = alicloud_vswitch.default[1].zone_id
    cidr_block = alicloud_vswitch.default[1].cidr_block
    vswitch_id = alicloud_vswitch.default[1].id
  }

}

resource "alicloud_pvtz_rule" "default" {
  endpoint_id = alicloud_pvtz_endpoint.default.id
  rule_name   = "${var.name}-${random_integer.default.result}"
  type        = "OUTBOUND"
  zone_name   = var.name
  forward_ips {
    ip   = "114.114.114.114"
    port = 8080
  }
}

resource "alicloud_pvtz_rule_attachment" "default" {
  rule_id = alicloud_pvtz_rule.default.id
  vpcs {
    region_id = data.alicloud_regions.default.regions.0.id
    vpc_id    = alicloud_vpc.default[0].id
  }
  vpcs {
    region_id = data.alicloud_regions.default.regions.0.id
    vpc_id    = alicloud_vpc.default[1].id
  }
}
```

## Argument Reference

The following arguments are supported:

* `rule_id` - (Required, ForceNew) The ID of the rule.
* `vpcs` - (Required) The List of the VPC. See [`vpcs`](#vpcs) below.

### `vpcs`

The vpcs supports the following:

* `vpc_id` - (Required) The ID of the VPC.  **NOTE:** The VPC that can be associated with the forwarding rule must belong to the same region as the Endpoint.
* `region_id` - (Required) The region of the vpc. If not set, the current region will instead of.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Rule Attachment. Its value is same as `rule_id`.

## Import

Private Zone Rule Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_pvtz_rule_attachment.example <rule_id>
```