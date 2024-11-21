---
subcategory: "Global Accelerator (GA)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ga_custom_routing_endpoint_traffic_policy"
sidebar_current: "docs-alicloud-resource-ga-custom-routing-endpoint-traffic-policy"
description: |-
  Provides a Alicloud Global Accelerator (GA) Custom Routing Endpoint Traffic Policy resource.
---

# alicloud_ga_custom_routing_endpoint_traffic_policy

Provides a Global Accelerator (GA) Custom Routing Endpoint Traffic Policy resource.

For information about Global Accelerator (GA) Custom Routing Endpoint Traffic Policy and how to use it, see [What is Custom Routing Endpoint Traffic Policy](https://www.alibabacloud.com/help/en/global-accelerator/latest/api-ga-2019-11-20-createcustomroutingendpointtrafficpolicies).

-> **NOTE:** Available since v1.197.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ga_custom_routing_endpoint_traffic_policy&exampleId=54c7150b-5960-6f53-4cc8-acb1eb2a46ca9f45d42a&activeTab=example&spm=docs.r.ga_custom_routing_endpoint_traffic_policy.0.54c7150b59&intl_lang=EN_US" target="_blank">
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

variable "name" {
  default = "tf-example"
}

data "alicloud_zones" "default" {
}

data "alicloud_ga_accelerators" "default" {
  status                 = "active"
  bandwidth_billing_type = "BandwidthPackage"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "192.168.192.0/24"
  zone_id      = data.alicloud_zones.default.ids.0
}

resource "alicloud_ga_bandwidth_package" "default" {
  bandwidth      = 100
  type           = "Basic"
  bandwidth_type = "Basic"
  payment_type   = "PayAsYouGo"
  billing_type   = "PayBy95"
  ratio          = 30
}

data "alicloud_regions" "default" {
  current = true
}

resource "alicloud_ga_bandwidth_package_attachment" "default" {
  accelerator_id       = data.alicloud_ga_accelerators.default.accelerators.1.id
  bandwidth_package_id = alicloud_ga_bandwidth_package.default.id
}

resource "alicloud_ga_listener" "default" {
  accelerator_id = alicloud_ga_bandwidth_package_attachment.default.accelerator_id
  listener_type  = "CustomRouting"
  port_ranges {
    from_port = 10000
    to_port   = 26000
  }
}

resource "alicloud_ga_custom_routing_endpoint_group" "default" {
  accelerator_id                     = alicloud_ga_listener.default.accelerator_id
  listener_id                        = alicloud_ga_listener.default.id
  endpoint_group_region              = data.alicloud_regions.default.regions.0.id
  custom_routing_endpoint_group_name = var.name
  description                        = var.name
}

resource "alicloud_ga_custom_routing_endpoint_group_destination" "default" {
  endpoint_group_id = alicloud_ga_custom_routing_endpoint_group.default.id
  protocols         = ["TCP"]
  from_port         = 1
  to_port           = 10
}

resource "alicloud_ga_custom_routing_endpoint" "default" {
  endpoint_group_id          = alicloud_ga_custom_routing_endpoint_group_destination.default.endpoint_group_id
  endpoint                   = alicloud_vswitch.default.id
  type                       = "PrivateSubNet"
  traffic_to_endpoint_policy = "AllowAll"
}

resource "alicloud_ga_custom_routing_endpoint_traffic_policy" "default" {
  endpoint_id = alicloud_ga_custom_routing_endpoint.default.custom_routing_endpoint_id
  address     = "192.168.192.2"
  port_ranges {
    from_port = 1
    to_port   = 2
  }
}
```

## Argument Reference

The following arguments are supported:

* `endpoint_id` - (Required, ForceNew) The ID of the Custom Routing Endpoint.
* `address` - (Required) The IP address of the destination to which traffic is allowed.
* `port_ranges` - (Optional, Set) Port rangeSee the following. See [`port_ranges`](#port_ranges) below.

### `port_ranges`

The port_ranges supports the following:

* `from_port` - (Optional, Int) The start port of the port range of the traffic destination. The specified port must fall within the port range of the specified endpoint group.
* `to_port` - (Optional, Int) The end port of the port range of the traffic destination. The specified port must fall within the port range of the specified endpoint group.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Custom Routing Endpoint Traffic Policy. It formats as `<endpoint_id>:<custom_routing_endpoint_traffic_policy_id>`.
* `accelerator_id` - The ID of the GA instance.
* `listener_id` - The ID of the listener.
* `endpoint_group_id` - The ID of the endpoint group.
* `custom_routing_endpoint_traffic_policy_id` - The ID of the Custom Routing Endpoint Traffic Policy.
* `status` - The status of the Custom Routing Endpoint Traffic Policy.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Custom Routing Endpoint Traffic Policy.
* `update` - (Defaults to 5 mins) Used when update the Custom Routing Endpoint Traffic Policy.
* `delete` - (Defaults to 5 mins) Used when delete the Custom Routing Endpoint Traffic Policy.

## Import

Global Accelerator (GA) Custom Routing Endpoint Traffic Policy can be imported using the id, e.g.

```shell
$ terraform import alicloud_ga_custom_routing_endpoint_traffic_policy.example <endpoint_id>:<custom_routing_endpoint_traffic_policy_id>
```
