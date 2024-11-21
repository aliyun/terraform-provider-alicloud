---
subcategory: "Event Bridge"
layout: "alicloud"
page_title: "Alicloud: alicloud_event_bridge_api_destination"
sidebar_current: "docs-alicloud-resource-event-bridge-api-destination"
description: |-
  Provides a Alicloud Event Bridge Api Destination resource.
---

# alicloud_event_bridge_api_destination

Provides a Event Bridge Api Destination resource. 

For information about Event Bridge Api Destination and how to use it, see [What is Api Destination](https://www.alibabacloud.com/help/en/eventbridge/latest/api-eventbridge-2020-04-01-createapidestination).

-> **NOTE:** Available since v1.211.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_event_bridge_api_destination&exampleId=2fb1b2d8-40ec-1a8a-7955-5cdab470a73af83cc4cd&activeTab=example&spm=docs.r.event_bridge_api_destination.0.2fb1b2d840&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = var.region
}

variable "region" {
  default = "cn-chengdu"
}

variable "name" {
  default = "terraform-example"
}

resource "alicloud_event_bridge_connection" "default" {
  connection_name = var.name
  network_parameters {
    network_type = "PublicNetwork"
  }
}

resource "alicloud_event_bridge_api_destination" "default" {
  connection_name      = alicloud_event_bridge_connection.default.connection_name
  api_destination_name = var.name
  description          = "test-api-destination-connection"
  http_api_parameters {
    endpoint = "http://127.0.0.1:8001"
    method   = "POST"
  }
}
```

## Argument Reference

The following arguments are supported:

* `connection_name` - (Required, ForceNew) The name of the connection.
* `api_destination_name` - (Required, ForceNew) The name of the API destination.
* `description` - (Optional) The description of the API destination.
* `http_api_parameters` - (Required, Set) The parameters that are configured for the API destination. See [`http_api_parameters`](#http_api_parameters) below.

### `http_api_parameters`

The http_api_parameters supports the following:

* `endpoint` - (Required) The endpoint of the API destination.
* `method` - (Required) The HTTP request method. Valid values: `GET`, `POST`, `HEAD`, `DELETE`, `PUT`, `PATCH`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Api Destination.
* `create_time` - The creation time of the Api Destination.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Api Destination.
* `update` - (Defaults to 5 mins) Used when update the Api Destination.
* `delete` - (Defaults to 5 mins) Used when delete the Api Destination.

## Import

Event Bridge Api Destination can be imported using the id, e.g.

```shell
$ terraform import alicloud_event_bridge_api_destination.example <id>
```
