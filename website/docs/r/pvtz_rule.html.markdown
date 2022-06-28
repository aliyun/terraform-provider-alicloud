---
subcategory: "Private Zone"
layout: "alicloud"
page_title: "Alicloud: alicloud_pvtz_rule"
sidebar_current: "docs-alicloud-resource-pvtz-rule"
description: |-
  Provides a Alicloud PrivateZone Rule resource.
---

# alicloud\_pvtz\_rule

Provides a Private Zone Rule resource.

For information about Private Zone Rule and how to use it, see [What is Rule](https://www.alibabacloud.com/help/en/doc-detail/177601.htm).

-> **NOTE:** Available in v1.143.0+.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "example_value"
}

data "alicloud_pvtz_resolver_zones" "default" {
  status = "NORMAL"
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
  endpoint_name     = var.name
  security_group_id = alicloud_security_group.default.id
  vpc_id            = alicloud_vpc.default.id
  vpc_region_id     = "vpc_region_id"
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
  rule_name   = var.name
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
* `forward_ips` - (Required) Forwarding target. See the following `Block forward_ip`.
* `rule_name` - (Required, ForceNew) The name of the resource.
* `type` - (Optional, ForceNew) The type of the rule. Valid values: `OUTBOUND`.
* `zone_name` - (Required, ForceNew) The name of the forwarding zone.

#### Block forward_ips

The forward_ips supports the following:

* `ip` - (Required) The ip of the forwarding destination.
* `port` - (Required) The port of the forwarding destination.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Rule.

## Import

Private Zone Rule can be imported using the id, e.g.

```
$ terraform import alicloud_pvtz_rule.example <id>
```