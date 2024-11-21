---
subcategory: "Private Zone"
layout: "alicloud"
page_title: "Alicloud: alicloud_pvtz_rule"
sidebar_current: "docs-alicloud-resource-pvtz-rule"
description: |-
  Provides a Alicloud PrivateZone Rule resource.
---

# alicloud_pvtz_rule

Provides a Private Zone Rule resource.

For information about Private Zone Rule and how to use it, see [What is Rule](https://www.alibabacloud.com/help/en/privatezone/latest/add-forwarding-rule).

-> **NOTE:** Available since v1.143.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_pvtz_rule&exampleId=df1ad7a0-8c69-8334-995d-09c45885f2051493b2ef&activeTab=example&spm=docs.r.pvtz_rule.0.df1ad7a08c&intl_lang=EN_US" target="_blank">
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
  vpc_name   = var.name
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
  count      = 2
  vpc_id     = alicloud_vpc.default.id
  cidr_block = cidrsubnet(alicloud_vpc.default.cidr_block, 8, count.index)
  zone_id    = data.alicloud_pvtz_resolver_zones.default.zones[count.index].zone_id
}

resource "alicloud_security_group" "default" {
  vpc_id = alicloud_vpc.default.id
  name   = var.name
}

resource "alicloud_pvtz_endpoint" "default" {
  endpoint_name     = "${var.name}-${random_integer.default.result}"
  security_group_id = alicloud_security_group.default.id
  vpc_id            = alicloud_vpc.default.id
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
```

## Argument Reference

The following arguments are supported:

* `endpoint_id` - (Required, ForceNew) The ID of the Endpoint.
* `forward_ips` - (Required) Forwarding target. See [`forward_ips`](#forward_ips) below.
* `rule_name` - (Required, ForceNew) The name of the resource.
* `type` - (Optional, ForceNew) The type of the rule. Valid values: `OUTBOUND`.
* `zone_name` - (Required, ForceNew) The name of the forwarding zone.

### `forward_ips`

The forward_ips supports the following:

* `ip` - (Required) The ip of the forwarding destination.
* `port` - (Required) The port of the forwarding destination.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Rule.

## Import

Private Zone Rule can be imported using the id, e.g.

```shell
$ terraform import alicloud_pvtz_rule.example <id>
```