---
subcategory: "Private Link"
layout: "alicloud"
page_title: "Alicloud: alicloud_privatelink_vpc_endpoint_zone"
sidebar_current: "docs-alicloud-resource-privatelink-vpc-endpoint-zone"
description: |-
  Provides a Alicloud Private Link Vpc Endpoint Zone resource.
---

# alicloud_privatelink_vpc_endpoint_zone

Provides a Private Link Vpc Endpoint Zone resource.

For information about Private Link Vpc Endpoint Zone and how to use it, see [What is Vpc Endpoint Zone](https://www.alibabacloud.com/help/en/privatelink/latest/api-privatelink-2020-04-15-addzonetovpcendpoint).

-> **NOTE:** Available since v1.111.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf_example"
}

data "alicloud_zones" "example" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_privatelink_vpc_endpoint_service" "example" {
  service_description    = var.name
  connect_bandwidth      = 103
  auto_accept_connection = false
}

resource "alicloud_vpc" "example" {
  vpc_name   = var.name
  cidr_block = "10.0.0.0/8"
}
resource "alicloud_vswitch" "example" {
  vswitch_name = var.name
  cidr_block   = "10.1.0.0/16"
  vpc_id       = alicloud_vpc.example.id
  zone_id      = data.alicloud_zones.example.zones.0.id
}

resource "alicloud_security_group" "example" {
  name   = var.name
  vpc_id = alicloud_vpc.example.id
}

resource "alicloud_slb_load_balancer" "example" {
  load_balancer_name = var.name
  vswitch_id         = alicloud_vswitch.example.id
  load_balancer_spec = "slb.s2.small"
  address_type       = "intranet"
}

resource "alicloud_privatelink_vpc_endpoint_service_resource" "example" {
  service_id    = alicloud_privatelink_vpc_endpoint_service.example.id
  resource_id   = alicloud_slb_load_balancer.example.id
  resource_type = "slb"
}

resource "alicloud_privatelink_vpc_endpoint" "example" {
  service_id         = alicloud_privatelink_vpc_endpoint_service_resource.example.service_id
  security_group_ids = [alicloud_security_group.example.id]
  vpc_id             = alicloud_vpc.example.id
  vpc_endpoint_name  = var.name
}

resource "alicloud_privatelink_vpc_endpoint_zone" "example" {
  endpoint_id = alicloud_privatelink_vpc_endpoint.example.id
  vswitch_id  = alicloud_vswitch.example.id
  zone_id     = data.alicloud_zones.example.zones.0.id
}
```

## Argument Reference

The following arguments are supported:

* `dry_run` - (Optional) The dry run.
* `endpoint_id` - (Required, ForceNew) The ID of the Vpc Endpoint.
* `vswitch_id` - (Required, ForceNew) The VSwitch id.
* `zone_id` - (Optional, Computed, ForceNew) The Zone Id.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Vpc Endpoint Zone. The value is formatted `<endpoint_id>:<zone_id>`.
* `status` - Status.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 4 mins) Used when create the Vpc Endpoint Zone.
* `delete` - (Defaults to 4 mins) Used when delete the Vpc Endpoint Zone.

## Import

Private Link Vpc Endpoint Zone can be imported using the id, e.g.

```shell
$ terraform import alicloud_privatelink_vpc_endpoint_zone.example <endpoint_id>:<zone_id>
```
