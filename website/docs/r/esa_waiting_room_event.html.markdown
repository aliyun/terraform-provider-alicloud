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

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_esa_waiting_room_event&exampleId=20edd2d2-c7c1-982c-e81e-a6f67d2f191166378cb0&activeTab=example&spm=docs.r.esa_waiting_room_event.0.20edd2d2c7&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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
* `custom_page_html` - (Optional) User-defined waiting room page content, when the waiting room type is custom type, you need to enter. The incoming content needs to be base64 encoded.
* `description` - (Optional) Waiting room description.
* `disable_session_renewal_enable` - (Optional) Disable session renewal. Value:
  - `on`: open.
  - `off`: closed.
* `end_time` - (Required) The timestamp of the end time of the event.
* `json_response_enable` - (Optional) JSON response switch. Value:
  - `on`: open.
  - `off`: closed.
* `language` - (Optional) Default language setting. Values include:
  - `enus`: English.
  - `zhcn`: Simplified Chinese.
  - `zhhk`: Traditional Chinese.
* `new_users_per_minute` - (Required) Number of new users per minute.
* `pre_queue_enable` - (Optional) Pre-queue switch.
  - `on`: open.
  - `off`: closed.
* `pre_queue_start_time` - (Optional) Pre-queue start time.
* `queuing_method` - (Required) Way of queuing. Value:
  - `random`: random.
  - `fifo`: first in, first out.
  - `passthrough`: through.
  - `reject-all`: reject all.
* `queuing_status_code` - (Required) Waiting room status code. Value:
  - `200`
  - `202`
  - `429`
* `random_pre_queue_enable` - (Optional) Random queue switch.
  - `on`: open.
  - `off`: closed.
* `session_duration` - (Required) User session duration in minutes.
* `site_id` - (Required, ForceNew) The site ID, which can be obtained by calling the ListSites API.
* `start_time` - (Required) The timestamp of the event start time.
* `status` - (Required) Enabled status. Value:
  - `on`: Enable waiting room events
  - `off`: Disable waiting room events
* `total_active_users` - (Required) Total number of active users.
* `waiting_room_event_name` - (Required) Event name, custom event description.
* `waiting_room_id` - (Optional, ForceNew, Computed) Waiting room ID, used to identify a specific waiting room. It can be obtained by calling the [listwaitingroom](https://help.aliyun.com/document_detail/2850279.html) interface.
* `waiting_room_type` - (Required) Waiting room type. The following types are supported:
  - `default`: the default type.
  - `custom`: custom type.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<site_id>:<waiting_room_id>:<waiting_room_event_id>`.
* `waiting_room_event_id` - The waiting room event ID, which can be obtained by calling the [ListWaitingRoomEvents](https://help.aliyun.com/document_detail/2850279.html) operation.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Waiting Room Event.
* `delete` - (Defaults to 5 mins) Used when delete the Waiting Room Event.
* `update` - (Defaults to 5 mins) Used when update the Waiting Room Event.

## Import

ESA Waiting Room Event can be imported using the id, e.g.

```shell
$ terraform import alicloud_esa_waiting_room_event.example <site_id>:<waiting_room_id>:<waiting_room_event_id>
```