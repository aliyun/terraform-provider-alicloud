---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_waiting_room"
description: |-
  Provides a Alicloud ESA Waiting Room resource.
---

# alicloud_esa_waiting_room

Provides a ESA Waiting Room resource.



For information about ESA Waiting Room and how to use it, see [What is Waiting Room](https://next.api.alibabacloud.com/document/ESA/2024-09-10/CreateWaitingRoom).

-> **NOTE:** Available since v1.244.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
}

resource "alicloud_esa_site" "resource_Site_example" {
  site_name   = "terraform.site"
  instance_id = data.alicloud_esa_sites.default.sites.0.instance_id
  coverage    = "overseas"
  access_type = "NS"
}

resource "alicloud_esa_waiting_room" "default" {
  queuing_method     = "fifo"
  session_duration   = "5"
  total_active_users = "300"
  host_name_and_path {
    domain    = "sub_domain.com"
    path      = "/example"
    subdomain = "example_sub_domain.com."
  }
  host_name_and_path {
    domain    = "sub_domain.com"
    path      = "/example"
    subdomain = "example_sub_domain1.com."
  }
  host_name_and_path {
    path      = "/example"
    subdomain = "example_sub_domain2.com."
    domain    = "sub_domain.com"
  }

  waiting_room_type              = "default"
  new_users_per_minute           = "200"
  custom_page_html               = ""
  language                       = "zhcn"
  queuing_status_code            = "200"
  waiting_room_name              = "waitingroom_example"
  status                         = "off"
  site_id                        = alicloud_esa_site.resource_Site_example.id
  queue_all_enable               = "off"
  disable_session_renewal_enable = "off"
  description                    = "example"
  json_response_enable           = "off"
  cookie_name                    = "__aliwaitingroom_example"
}
```

## Argument Reference

The following arguments are supported:
* `cookie_name` - (Required) Custom Cookie name.
* `custom_page_html` - (Optional) User-defined waiting room page content, when the waiting room type is custom type, you need to enter. The incoming content needs to be base64 encoded.
* `description` - (Optional) Waiting room description.
* `disable_session_renewal_enable` - (Optional) Disable session renewal. Value:
  - `on`: open.
  - `off`: closed.
* `host_name_and_path` - (Required, List) Host name and path. See [`host_name_and_path`](#host_name_and_path) below.
* `json_response_enable` - (Optional) The JSON response. If the accept request header contains "application/json", JSON data is returned. Value:
  - `on`: open.
  - `off`: closed.
* `language` - (Optional) The language of the waiting room page. When the waiting room type is the default type, it needs to be passed in. The following types are supported:
  - `enus`: English.
  - `zhcn`: Simplified Chinese.
  - `zhhk`: Traditional Chinese.
* `new_users_per_minute` - (Required) Number of new users per minute.
* `queue_all_enable` - (Optional) All in line. Value:
  - `on`: open.
  - `off`: closed.
* `queuing_method` - (Required) Way of queuing. Value:
  - `random`: random.
  - `fifo`: first in, first out.
  - `Passthrough`: through.
  - `Reject-all`: reject all.
* `queuing_status_code` - (Required) Waiting room status code. Value:
  - `200`
  - `202`
  - `429`
* `session_duration` - (Required) Session duration in minutes.
* `site_id` - (Required, ForceNew) The site ID, which can be obtained by calling the [ListSites](https://help.aliyun.com/document_detail/2850189.html) interface.
* `status` - (Required) Waiting room enabled status. Value:
  - 'on': Enable waiting room
  - 'off': Disabled waiting room
* `total_active_users` - (Required) Total number of active users.
* `waiting_room_name` - (Required) The name of the waiting room.
* `waiting_room_type` - (Required) Waiting room type, support:
  - `default`: Indicates the default type.
  - `custom`: indicates a custom type.

### `host_name_and_path`

The host_name_and_path supports the following:
* `domain` - (Required) The domain name.
* `path` - (Required) The path.
* `subdomain` - (Required) The subdomain.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<site_id>:<waiting_room_id>`.
* `waiting_room_id` - The waiting room ID, which can be obtained by calling the [ListWaitingRooms](https://help.aliyun.com/document_detail/2850279.html) API.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Waiting Room.
* `delete` - (Defaults to 5 mins) Used when delete the Waiting Room.
* `update` - (Defaults to 5 mins) Used when update the Waiting Room.

## Import

ESA Waiting Room can be imported using the id, e.g.

```shell
$ terraform import alicloud_esa_waiting_room.example <site_id>:<waiting_room_id>
```