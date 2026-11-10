---
subcategory: "Private Zone"
layout: "alicloud"
page_title: "Alicloud: alicloud_pvtz_inbound_endpoint"
sidebar_current: "docs-alicloud-resource-pvtz-inbound-endpoint"
description: |-
  Provides an Alicloud Private Zone Inbound Endpoint resource.
---

# alicloud_pvtz_inbound_endpoint

Provides a Private Zone Inbound Endpoint resource.

For information about Private Zone Inbound Endpoint and how to use it, see [What is Inbound Endpoint](https://www.alibabacloud.com/help/en/privatezone/latest/add-inbound-endpoint).

-> **NOTE:** Available since v2.3.0.


## Example Usage

```terraform
terraform {
  required_providers {
    alicloud = {
      source = "aliyun/alicloud"
    }
  }
}

variable "name" {
  default = "terraform-example"
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
  vpc_id              = alicloud_vpc.default.id
  security_group_name = var.name
}

resource "alicloud_pvtz_inbound_endpoint" "default" {
  endpoint_name     = var.name
  security_group_id = alicloud_security_group.default.id
  vpc_id            = alicloud_vpc.default.id
  vpc_region_id     = data.alicloud_regions.default.regions.0.id
  spec_version      = "standard"
  ip_configs {
    az_id      = alicloud_vswitch.default[0].zone_id
    cidr_block = alicloud_vswitch.default[0].cidr_block
    vswitch_id = alicloud_vswitch.default[0].id
  }
  ip_configs {
    az_id      = alicloud_vswitch.default[1].zone_id
    cidr_block = alicloud_vswitch.default[1].cidr_block
    vswitch_id = alicloud_vswitch.default[1].id
  }
}
```

## Argument Reference

The following arguments are supported:

* `endpoint_name` - (Required, String) The name of the resource.
* `ip_configs` - (Required, Set) The Ip Configs. See [`ip_configs`](#ip_configs) below. **NOTE:** In order to ensure high availability, add at least 2 and up to 10.
* `security_group_id` - (Required, ForceNew, String) The ID of the Security Group.
* `vpc_id` - (Required, ForceNew, String) The VPC ID.
* `vpc_region_id` - (Required, ForceNew, String) The Region of the VPC.
* `spec_version` - (Required, ForceNew, String) The specification version. Valid values: `standard`.

### `ip_configs`

The ip_configs supports the following:

* `az_id` - (Required, String) The availability zone ID of the VSwitch.
* `cidr_block` - (Required, String) The subnet mask, in the CIDR format of the VSwitch.
* `ip` - (Optional, String) The IP address within the subnet mask range. If you omit it when the endpoint is created, the system assigns an IP address from the VSwitch. If you set it, the value is kept in the configuration and sent on update instead of being replaced by an automatically assigned address. If you only remove an explicitly configured `ip` from the configuration, the address that is currently in use is kept and no IP change is sent for that entry. If `ip_configs` changes while an entry has no explicit `ip`, that entry is sent without an address and the system assigns one.
* `vswitch_id` - (Required, String) The VSwitch ID.

## Attributes Reference

The following attributes are exported:

* `id` - (String) The resource ID in terraform of Inbound Endpoint.
* `status` - (String) The status of the resource. Valid values: `CHANGE_FAILED`, `CHANGE_INIT`, `EXCEPTION`, `FAILED`, `INIT`, `SUCCESS`.
* `endpoint_type` - (String) The type of the inbound endpoint reported by the service.

## Import

Private Zone Inbound Endpoint can be imported using the id, e.g.

```shell
$ terraform import alicloud_pvtz_inbound_endpoint.example <id>
```
