---
subcategory: "Global Accelerator (GA)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ga_custom_routing_endpoint_group_destination"
sidebar_current: "docs-alicloud-resource-ga-custom-routing-endpoint-group-destination"
description: |-
  Provides a Alicloud Global Accelerator (GA) Custom Routing Endpoint Group Destination resource.
---

# alicloud_ga_custom_routing_endpoint_group_destination

Provides a Global Accelerator (GA) Custom Routing Endpoint Group Destination resource.

For information about Global Accelerator (GA) Custom Routing Endpoint Group Destination and how to use it, see [What is Custom Routing Endpoint Group Destination](https://www.alibabacloud.com/help/en/global-accelerator/latest/api-ga-2019-11-20-createcustomroutingendpointgroupdestinations).

-> **NOTE:** Available since v1.197.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ga_custom_routing_endpoint_group_destination&exampleId=70e988a4-74bc-2b5f-5435-02374a03eb4b6a3bc3ec&activeTab=example&spm=docs.r.ga_custom_routing_endpoint_group_destination.0.70e988a474&intl_lang=EN_US" target="_blank">
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

resource "alicloud_ga_custom_routing_endpoint_group_destination" "default" {
  endpoint_group_id = alicloud_ga_custom_routing_endpoint_group.default.id
  protocols         = ["TCP"]
  from_port         = 1
  to_port           = 2
}
```

## Argument Reference

The following arguments are supported:

* `endpoint_group_id` - (Required, ForceNew) The ID of the endpoint group.
* `protocols` - (Required, List) The backend service protocol of the endpoint group. Valid values: `TCP`, `UDP`, `TCP, UDP`.
* `from_port` - (Required, Int) The start port of the backend service port range of the endpoint group. The `from_port` value must be smaller than or equal to the `to_port` value. Valid values: `1` to `65499`.
* `to_port` - (Required, Int) The end port of the backend service port range of the endpoint group. The `from_port` value must be smaller than or equal to the `to_port` value. Valid values: `1` to `65499`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Custom Routing Endpoint Group Destination. It formats as `<endpoint_group_id>:<custom_routing_endpoint_group_destination_id>`.
* `accelerator_id` - The ID of the GA instance.
* `listener_id` - The ID of the listener.
* `custom_routing_endpoint_group_destination_id` - The ID of the Custom Routing Endpoint Group Destination.
* `status` - The status of the Custom Routing Endpoint Group Destination.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Custom Routing Endpoint Group Destination.
* `update` - (Defaults to 5 mins) Used when update the Custom Routing Endpoint Group Destination.
* `delete` - (Defaults to 5 mins) Used when delete the Custom Routing Endpoint Group Destination.

## Import

Global Accelerator (GA) Custom Routing Endpoint Group Destination can be imported using the id, e.g.

```shell
$ terraform import alicloud_ga_custom_routing_endpoint_group_destination.example <endpoint_group_id>:<custom_routing_endpoint_group_destination_id>
```
