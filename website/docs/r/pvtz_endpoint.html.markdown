---
subcategory: "Private Zone"
layout: "alicloud"
page_title: "Alicloud: alicloud_pvtz_endpoint"
sidebar_current: "docs-alicloud-resource-pvtz-endpoint"
description: |-
  Provides a Alicloud Private Zone Endpoint resource.
---

# alicloud\_pvtz\_endpoint

Provides a Private Zone Endpoint resource.

For information about Private Zone Endpoint and how to use it, see [What is Endpoint](https://www.alibabacloud.com/help/en/doc-detail/64611.htm).

-> **NOTE:** Available in v1.143.0+.


## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "example_name"
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

```

## Argument Reference

The following arguments are supported:

* `endpoint_name` - (Required) The name of the resource.
* `ip_configs` - (Required) The Ip Configs. See the following `Block ip_configs`. **NOTE:** In order to ensure high availability, add at least 2 and up to 6.
* `security_group_id` - (Required, ForceNew) The ID of the Security Group.
* `vpc_id` - (Required, ForceNew) The VPC ID.
* `vpc_region_id` - (Required, ForceNew) The Region of the VPC.

#### Block ip_configs

The ip_configs supports the following: 

* `cidr_block` - (Required) The Subnet mask.
* `ip` - (Optional, Computed) The IP address within the parameter range of the subnet mask.  It is recommended to use the IP address assigned by the system.
* `vswitch_id` - (Required) The Vswitch id.
* `zone_id` - (Required) The Zone ID.


## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Endpoint.
* `status` - The status of the resource. Valid values: `CHANGE_FAILED`, `CHANGE_INIT`, `EXCEPTION`, `FAILED`, `INIT`, `SUCCESS`.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 10 mins) Used when create the Endpoint.
* `update` - (Defaults to 10 mins) Used when update the Endpoint.

## Import

Private Zone Endpoint can be imported using the id, e.g.

```
$ terraform import alicloud_pvtz_endpoint.example <id>
```