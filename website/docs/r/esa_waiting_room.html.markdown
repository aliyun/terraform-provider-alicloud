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

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_esa_waiting_room&exampleId=a04eaebe-1739-0a22-5522-32ab1ae01dc72651f739&activeTab=example&spm=docs.r.esa_waiting_room.0.a04eaebe17&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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
* `cookie_name` - (Required) The name of the custom cookie.
* `custom_page_html` - (Optional) The HTML content or identifier of the custom queuing page. This parameter is valid only when `WaitingRoomType` is set to `custom`. The content must be URL-encoded.
* `description` - (Optional) Specifies whether to enable JSON response. If you set this parameter to on, a JSON body is returned for requests to the waiting room with the header Accept: application/json. Valid values:

  - on
  - off
* `disable_session_renewal_enable` - (Optional) The maximum duration for which a session remains valid after a user leaves the origin. Unit: minutes.
* `host_name_and_path` - (Required, List) The details of the hostname and path. See [`host_name_and_path`](#host_name_and_path) below.
* `json_response_enable` - (Optional) Indicates whether JSON response is enabled. If you set this parameter to on, a JSON body is returned for requests to the waiting room with the header Accept: application/json. Valid values:

  - on
  - off
* `language` - (Optional) The language of the waiting room page. This parameter is returned when the waiting room type is set to default. Valid values:

  - enus: English.
  - zhcn: Simplified Chinese.
  - zhhk: Traditional Chinese.
* `new_users_per_minute` - (Required) The maximum number of new users per minute.
* `queue_all_enable` - (Optional) Indicates whether all requests must be queued. Valid values:

  - on
  - off
* `queuing_method` - (Required) The queuing method. Valid values:

  - random: Users gain access to the origin randomly, regardless of the arrival time.
  - fifo: Users gain access to the origin in order of arrival.
  - passthrough: Users pass through the waiting room and go straight to the origin.
  - reject-all: Users are blocked from reaching the origin.
* `queuing_status_code` - (Required) The queuing method. Valid values:

  - random: Users gain access to the origin randomly, regardless of the arrival time.
  - fifo: Users gain access to the origin in order of arrival.
  - passthrough: Users pass through the waiting room and go straight to the origin.
  - reject-all: Users are blocked from reaching the origin.
* `session_duration` - (Required) The maximum duration for which a session remains valid after a user leaves the origin. Unit: minutes.
* `site_id` - (Required, ForceNew, Int) 
* `status` - (Required) The ID of the waiting room, which can be obtained by calling the [ListWaitingRooms](https://www.alibabacloud.com/help/en/doc-detail/2850279.html) operation.
* `total_active_users` - (Required) The maximum number of active users.
* `waiting_room_name` - (Required) Specifies whether to enable the waiting room. Valid values:

  - on
  - off
* `waiting_room_type` - (Required) The type of the waiting room. Valid values:

  - default
  - custom

### `host_name_and_path`

The host_name_and_path supports the following:
* `domain` - (Required) The domain name.
* `path` - (Required) The path.
* `subdomain` - (Required) The subdomain.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<site_id>:<waiting_room_id>`.
* `waiting_room_id` - The website ID, which can be obtained by calling the [ListSites](https://www.alibabacloud.com/help/en/doc-detail/2850189.html) operation.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Waiting Room.
* `delete` - (Defaults to 5 mins) Used when delete the Waiting Room.
* `update` - (Defaults to 5 mins) Used when update the Waiting Room.

## Import

ESA Waiting Room can be imported using the id, e.g.

```shell
$ terraform import alicloud_esa_waiting_room.example <site_id>:<waiting_room_id>
```