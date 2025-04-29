---
subcategory: "Private Link"
layout: "alicloud"
page_title: "Alicloud: alicloud_privatelink_vpc_endpoint_zone"
description: |-
  Provides a Alicloud Private Link Vpc Endpoint Zone resource.
---

# alicloud_privatelink_vpc_endpoint_zone

Provides a Private Link Vpc Endpoint Zone resource.



For information about Private Link Vpc Endpoint Zone and how to use it, see [What is Vpc Endpoint Zone](https://www.alibabacloud.com/help/en/privatelink/latest/api-privatelink-2020-04-15-addzonetovpcendpoint).

-> **NOTE:** Available since v1.111.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_privatelink_vpc_endpoint_zone&exampleId=002be05b-cddb-26ee-e297-a2349d72c10955d7a11d&activeTab=example&spm=docs.r.privatelink_vpc_endpoint_zone.0.002be05bcd&intl_lang=EN_US" target="_blank">
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

resource "alicloud_privatelink_vpc_endpoint_zone" "example" {
  endpoint_id = alicloud_privatelink_vpc_endpoint.example.id
  vswitch_id  = alicloud_vswitch.example.id
  zone_id     = data.alicloud_zones.example.zones.0.id
}
```

## Argument Reference

The following arguments are supported:
* `dry_run` - (Optional) Specifies whether to perform only a dry run, without performing the actual request. Valid values:
  - `true`: performs only a dry run. The system checks the request for potential issues, including missing parameter values, incorrect request syntax, and service limits. If the request fails the dry run, an error message is returned. If the request passes the dry run, the DryRunOperation error code is returned.
  - **false (default)**: performs a dry run and performs the actual request. If the request passes the dry run, a 2xx HTTP status code is returned and the operation is performed.
* `endpoint_id` - (Required, ForceNew) The endpoint ID.
* `eni_ip` - (Optional, ForceNew, Computed, Available since v1.212.0) The IP address of the endpoint ENI.
* `vswitch_id` - (Required, ForceNew) The ID of the vSwitch in the zone.
* `zone_id` - (Optional, Computed, ForceNew) The zone ID.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<endpoint_id>:<zone_id>`.
* `status` - The state of the zone. 
* `region_id` - (Available since v1.235.0) The ID of the region to which the endpoint service belongs.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Vpc Endpoint Zone.
* `delete` - (Defaults to 5 mins) Used when delete the Vpc Endpoint Zone.

## Import

Private Link Vpc Endpoint Zone can be imported using the id, e.g.

```shell
$ terraform import alicloud_privatelink_vpc_endpoint_zone.example <endpoint_id>:<zone_id>
```