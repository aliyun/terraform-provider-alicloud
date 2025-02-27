---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_waiting_room_event"
description: |-
  Provides a Alicloud ESA Waiting Room Event resource.
---

# alicloud_esa_waiting_room_event

Provides a ESA Waiting Room Event resource.



For information about ESA Waiting Room Event and how to use it, see [What is Waiting Room Event](https://next.api.alibabacloud.com/document/ESA/2024-09-10/CreateWaitingRoomEvent).

-> **NOTE:** Available since v1.244.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
}

resource "alicloud_esa_site" "default" {
  site_name   = "chenxin0116.site"
  instance_id = data.alicloud_esa_sites.default.sites.0.instance_id
  coverage    = "overseas"
  access_type = "NS"
}

resource "alicloud_esa_waiting_room" "default" {
  status                         = "off"
  site_id                        = alicloud_esa_site.default.id
  json_response_enable           = "off"
  description                    = "example"
  waiting_room_type              = "default"
  disable_session_renewal_enable = "off"
  cookie_name                    = "__aliwaitingroom_example"
  waiting_room_name              = "waitingroom_example"
  queue_all_enable               = "off"
  queuing_status_code            = "200"
  custom_page_html               = ""
  new_users_per_minute           = "200"
  session_duration               = "5"
  language                       = "zhcn"
  total_active_users             = "300"
  queuing_method                 = "fifo"
  host_name_and_path {
    domain    = "sub_domain.com"
    path      = "/example"
    subdomain = "example_sub_domain.com."
  }
}

resource "alicloud_esa_waiting_room_event" "default" {
  waiting_room_id                = alicloud_esa_waiting_room.default.waiting_room_id
  end_time                       = "1719863200"
  waiting_room_event_name        = "WaitingRoomEvent_example"
  pre_queue_start_time           = ""
  random_pre_queue_enable        = "off"
  json_response_enable           = "off"
  site_id                        = alicloud_esa_site.default.id
  pre_queue_enable               = "off"
  description                    = "example"
  new_users_per_minute           = "200"
  queuing_status_code            = "200"
  custom_page_html               = ""
  language                       = "zhcn"
  total_active_users             = "300"
  waiting_room_type              = "default"
  start_time                     = "1719763200"
  status                         = "off"
  disable_session_renewal_enable = "off"
  queuing_method                 = "fifo"
  session_duration               = "5"
}
```

## Argument Reference

The following arguments are supported:
* `custom_page_html` - (Optional) The type of the waiting room. Valid values:

  - default
  - custom
* `description` - (Optional) Specifies whether to enable JSON response. Valid values:

  - on
  - off
* `disable_session_renewal_enable` - (Optional) The maximum duration for which a session remains valid after a user leaves the origin. Unit: minutes.
* `end_time` - (Required) The start time of the event. This value is a UNIX timestamp.
* `json_response_enable` - (Optional) The HTTP status code to return while a user is in the queue. Valid values:

  - 200
  - 202
  - 429
* `language` - (Optional) Specifies whether to enable random queuing.

  - on
  - off
* `new_users_per_minute` - (Required) The maximum number of active users.
* `pre_queue_enable` - (Optional) The end time of the event. This value is a UNIX timestamp.
* `pre_queue_start_time` - (Optional) Specifies whether to enable pre-queuing.

  - on
  - off
* `queuing_method` - (Required) Specifies whether to disable session renewal. Valid values:

  - on
  - off
* `queuing_status_code` - (Required) The queuing method. Valid values:

  - random: Users gain access to the origin randomly, regardless of the arrival time.
  - fifo: Users gain access to the origin in order of arrival.
  - passthrough: Users pass through the waiting room and go straight to the origin.
  - reject-all: All requests are blocked from accessing the origin.
* `random_pre_queue_enable` - (Optional) The start time for pre-queuing.
* `session_duration` - (Required) The maximum number of new users per minute.
* `site_id` - (Required, ForceNew, Int) The site ID, which can be obtained by calling the ListSites API.
* `start_time` - (Required) The content of the custom waiting room page. You must specify this parameter if you set WaitingRoomType to custom. The content must be Base64-encoded.
* `status` - (Required) The ID of the waiting room event, which can be obtained by calling the [ListWaitingRoomEvents](https://www.alibabacloud.com/help/en/doc-detail/2850279.html) operation.
* `total_active_users` - (Required) The name of the waiting room event.
* `waiting_room_event_name` - (Required) Specifies whether to enable the waiting room event. Valid values:

  -   `on`
  -   `off`
* `waiting_room_id` - (Optional, ForceNew, Computed) The website ID, which can be obtained by calling the [ListSites](https://www.alibabacloud.com/help/en/doc-detail/2850189.html) operation.
* `waiting_room_type` - (Required) The description of the waiting room.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<site_id>:<waiting_room_id>:<waiting_room_event_id>`.
* `waiting_room_event_id` - The unique ID of the waiting room, which can be obtained by calling the [ListWaitingRooms](https://www.alibabacloud.com/help/en/doc-detail/2850279.html) operation.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Waiting Room Event.
* `delete` - (Defaults to 5 mins) Used when delete the Waiting Room Event.
* `update` - (Defaults to 5 mins) Used when update the Waiting Room Event.

## Import

ESA Waiting Room Event can be imported using the id, e.g.

```shell
$ terraform import alicloud_esa_waiting_room_event.example <site_id>:<waiting_room_id>:<waiting_room_event_id>
```