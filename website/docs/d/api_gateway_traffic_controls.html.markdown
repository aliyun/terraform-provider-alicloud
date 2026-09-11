---
subcategory: "Api Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_api_gateway_traffic_controls"
sidebar_current: "docs-alicloud-datasource-api-gateway-traffic-controls"
description: |-
  Provides a list of Api Gateway Traffic Controls to the user.
---

# alicloud_api_gateway_traffic_controls

This data source provides the Api Gateway Traffic Controls of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.236.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_api_gateway_traffic_controls" "ids" {
  ids = ["example_id"]
}

output "api_gateway_traffic_control_id_1" {
  value = data.alicloud_api_gateway_traffic_controls.ids.controls.0.id
}

data "alicloud_api_gateway_traffic_controls" "name" {
  traffic_control_name = "example_name"
}

output "api_gateway_traffic_control_id_2" {
  value = data.alicloud_api_gateway_traffic_controls.name.controls.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed) A list of Traffic Control IDs.
* `traffic_control_id` - (Optional, ForceNew) The ID of the Traffic Control.
* `traffic_control_name` - (Optional, ForceNew) The name of the Traffic Control.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `controls` - A list of Api Gateway Traffic Controls. Each element contains the following attributes: See [`controls`](#controls) below.

### `controls`

A list of Api Gateway Traffic Controls. Each element contains the following attributes:
* `id` - The ID of the Traffic Control.
* `traffic_control_id` - The ID of the Traffic Control.
* `traffic_control_name` - The name of the Traffic Control.
* `traffic_control_unit` - The unit of the Traffic Control.
* `api_default` - The default flow control value for each API.
* `user_default` - The default flow control value for each user.
* `app_default` - The default flow control value for each app.
* `description` - The description of the Traffic Control.
* `create_time` - The creation time of the Traffic Control.
* `modified_time` - The last modification time of the Traffic Control.
* `region_id` - The region ID of the Traffic Control.
