---
subcategory: "Global Accelerator (GA)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ga_custom_routing_endpoint"
sidebar_current: "docs-alicloud-resource-ga-custom-routing-endpoint"
description: |-
  Provides a Alicloud Global Accelerator (GA) Custom Routing Endpoint resource.
---

# alicloud_ga_custom_routing_endpoint

Provides a Global Accelerator (GA) Custom Routing Endpoint resource.

For information about Global Accelerator (GA) Custom Routing Endpoint and how to use it, see [What is Custom Routing Endpoint](https://www.alibabacloud.com/help/en/global-accelerator/latest/api-ga-2019-11-20-createcustomroutingendpoints).

-> **NOTE:** Available since v1.197.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ga_custom_routing_endpoint&exampleId=71cb934b-198f-e902-40cd-05389ea73052f78f7a14&activeTab=example&spm=docs.r.ga_custom_routing_endpoint.0.71cb934b19&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "region" {
  default = "cn-hangzhou"
}

provider "alicloud" {
  region = var.region
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  vpc_name   = "terraform-example"
  cidr_block = "172.17.3.0/24"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = "terraform-example"
  cidr_block   = "172.17.3.0/24"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_ga_accelerator" "default" {
  duration        = 1
  auto_use_coupon = true
  spec            = "1"
}

resource "alicloud_ga_bandwidth_package" "default" {
  bandwidth      = 100
  type           = "Basic"
  bandwidth_type = "Basic"
  payment_type   = "PayAsYouGo"
  billing_type   = "PayBy95"
  ratio          = 30
}

resource "alicloud_ga_bandwidth_package_attachment" "default" {
  accelerator_id       = alicloud_ga_accelerator.default.id
  bandwidth_package_id = alicloud_ga_bandwidth_package.default.id
}

resource "alicloud_ga_listener" "default" {
  accelerator_id = alicloud_ga_bandwidth_package_attachment.default.accelerator_id
  listener_type  = "CustomRouting"
  port_ranges {
    from_port = 10000
    to_port   = 16000
  }
}

resource "alicloud_ga_custom_routing_endpoint_group" "default" {
  accelerator_id                     = alicloud_ga_listener.default.accelerator_id
  listener_id                        = alicloud_ga_listener.default.id
  endpoint_group_region              = var.region
  custom_routing_endpoint_group_name = "terraform-example"
  description                        = "terraform-example"
}

resource "alicloud_ga_custom_routing_endpoint" "default" {
  endpoint_group_id          = alicloud_ga_custom_routing_endpoint_group.default.id
  endpoint                   = alicloud_vswitch.default.id
  type                       = "PrivateSubNet"
  traffic_to_endpoint_policy = "DenyAll"
}
```

## Argument Reference

The following arguments are supported:

* `endpoint_group_id` - (Required, ForceNew) The ID of the endpoint group in which to create endpoints.
* `endpoint` - (Required, ForceNew) The ID of the endpoint (vSwitch).
* `type` - (Required, ForceNew) The backend service type of the endpoint. Valid values: `PrivateSubNet`.
* `traffic_to_endpoint_policy` - (Optional) The access policy of traffic to the endpoint. Default value: `DenyAll`. Valid values:
  - `DenyAll`: denies all traffic to the endpoint.
  - `AllowAll`: allows all traffic to the endpoint.
  - `AllowCustom`: allows traffic only to specified destinations in the endpoint.
  
## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Custom Routing Endpoint. It formats as `<endpoint_group_id>:<custom_routing_endpoint_id>`.
* `accelerator_id` - The ID of the GA instance with which the endpoint is associated.
* `listener_id` - The ID of the listener with which the endpoint is associated.
* `custom_routing_endpoint_id` - The ID of the Custom Routing Endpoint.
* `status` - The status of the Custom Routing Endpoint.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Custom Routing Endpoint.
* `update` - (Defaults to 5 mins) Used when update the Custom Routing Endpoint.
* `delete` - (Defaults to 5 mins) Used when delete the Custom Routing Endpoint.

## Import

Global Accelerator (GA) Custom Routing Endpoint can be imported using the id, e.g.

```shell
$ terraform import alicloud_ga_custom_routing_endpoint.example <endpoint_group_id>:<custom_routing_endpoint_id>
```
