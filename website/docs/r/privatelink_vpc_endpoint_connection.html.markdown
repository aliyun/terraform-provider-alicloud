---
subcategory: "Private Link"
layout: "alicloud"
page_title: "Alicloud: alicloud_privatelink_vpc_endpoint_connection"
description: |-
  Provides a Alicloud Private Link Vpc Endpoint Connection resource.
---

# alicloud_privatelink_vpc_endpoint_connection

Provides a Private Link Vpc Endpoint Connection resource.

vpc endpoint connection.

For information about Private Link Vpc Endpoint Connection and how to use it, see [What is Vpc Endpoint Connection](https://www.alibabacloud.com/help/en/privatelink/latest/api-privatelink-2020-04-15-enablevpcendpointzoneconnection).

-> **NOTE:** Available since v1.110.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_privatelink_vpc_endpoint_connection&exampleId=0847e1ba-7d2d-7ade-4392-292c3d8b7fdae9806382&activeTab=example&spm=docs.r.privatelink_vpc_endpoint_connection.0.0847e1ba7d&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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
* `bandwidth` - (Optional, Computed, Int) The bandwidth of the endpoint connection. Valid values: 1024 to 10240. Unit: Mbit/s.

  Note: The bandwidth of an endpoint connection is in the range of 100 to 10,240 Mbit/s. The default bandwidth is 1,024 Mbit/s. When the endpoint is connected to the endpoint service, the default bandwidth is the minimum bandwidth. In this case, the connection bandwidth range is 1,024 to 10,240 Mbit/s.
* `dry_run` - (Optional) Specifies whether to perform only a dry run, without performing the actual request. Valid values:
  - `true`: performs only a dry run. The system checks the request for potential issues, including missing parameter values, incorrect request syntax, and service limits. If the request fails the dry run, an error message is returned. If the request passes the dry run, the DryRunOperation error code is returned.
  - **false (default)**: performs a dry run and performs the actual request. If the request passes the dry run, a 2xx HTTP status code is returned and the operation is performed.
* `endpoint_id` - (Required, ForceNew) The endpoint ID.
* `service_id` - (Required, ForceNew) The endpoint service ID.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<service_id>:<endpoint_id>`.
* `status` - The state of the endpoint connection. 

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Vpc Endpoint Connection.
* `delete` - (Defaults to 5 mins) Used when delete the Vpc Endpoint Connection.
* `update` - (Defaults to 5 mins) Used when update the Vpc Endpoint Connection.

## Import

Private Link Vpc Endpoint Connection can be imported using the id, e.g.

```shell
$ terraform import alicloud_privatelink_vpc_endpoint_connection.example <service_id>:<endpoint_id>
```