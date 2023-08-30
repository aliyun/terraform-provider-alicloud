---
subcategory: "Private Link"
layout: "alicloud"
page_title: "Alicloud: alicloud_privatelink_vpc_endpoint_connection"
sidebar_current: "docs-alicloud-resource-privatelink-vpc-endpoint-connection"
description: |-
  Provides a Alicloud Private Link Vpc Endpoint Connection resource.
---

# alicloud_privatelink_vpc_endpoint_connection

Provides a Private Link Vpc Endpoint Connection resource.

For information about Private Link Vpc Endpoint Connection and how to use it, see [What is Vpc Endpoint Connection](https://www.alibabacloud.com/help/en/privatelink/latest/api-privatelink-2020-04-15-enablevpcendpointzoneconnection).

-> **NOTE:** Available since v1.110.0.

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

resource "alicloud_privatelink_vpc_endpoint_connection" "example" {
  endpoint_id = alicloud_privatelink_vpc_endpoint.example.id
  service_id  = alicloud_privatelink_vpc_endpoint.example.service_id
  bandwidth   = "1024"
}
```

## Argument Reference

The following arguments are supported:

* `bandwidth` - (Optional) The Bandwidth.
* `dry_run` - (Optional) The dry run.
* `endpoint_id` - (Required, ForceNew) The ID of the Vpc Endpoint.
* `service_id` - (Required, ForceNew) The ID of the Vpc Endpoint Service.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Vpc Endpoint Connection. The value is formatted `<service_id>:<endpoint_id>`.
* `status` - The status of Vpc Endpoint Connection.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 11 mins) Used when create the Vpc Endpoint Connection.
* `delete` - (Defaults to 6 mins) Used when delete the Vpc Endpoint Connection.
* `update` - (Defaults to 4 mins) Used when update the Vpc Endpoint Connection.

## Import

Private Link Vpc Endpoint Connection can be imported using the id, e.g.

```shell
$ terraform import alicloud_privatelink_vpc_endpoint_connection.example <service_id>:<endpoint_id>
```
