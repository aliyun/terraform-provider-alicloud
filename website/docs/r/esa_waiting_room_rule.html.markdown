---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_waiting_room_rule"
description: |-
  Provides a Alicloud ESA Waiting Room Rule resource.
---

# alicloud_esa_waiting_room_rule

Provides a ESA Waiting Room Rule resource.



For information about ESA Waiting Room Rule and how to use it, see [What is Waiting Room Rule](https://next.api.alibabacloud.com/document/ESA/2024-09-10/CreateWaitingRoomRule).

-> **NOTE:** Available since v1.244.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_esa_waiting_room_rule&exampleId=4fcf141c-b3a7-abce-a519-49066e7d08b607966cc7&activeTab=example&spm=docs.r.esa_waiting_room_rule.0.4fcf141cb3&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
}

resource "alicloud_esa_site" "default" {
  site_name   = "terraform.site"
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

resource "alicloud_esa_waiting_room_rule" "default" {
  rule            = "(http.host eq \"video.example.com\")"
  waiting_room_id = alicloud_esa_waiting_room.default.waiting_room_id
  rule_name       = "WaitingRoomRule_example1"
  status          = "off"
  site_id         = alicloud_esa_site.default.id
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_esa_waiting_room_rule&spm=docs.r.esa_waiting_room_rule.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `rule` - (Required) The content of the rule, the implemented policy or conditional expression.
* `status` - (Required, Available since v1.265.0) Rule switch. When adding global configuration, this parameter does not need to be set. Value range:
  - on: open.
  - off: close.
* `rule_name` - (Required) Rule name, optional, used to query by waiting room bypass rule name.
* `site_id` - (Required, ForceNew) The site ID, which can be obtained by calling the ListSites API.
* `waiting_room_id` - (Required, ForceNew) Waiting room ID, used to identify a specific waiting room. It can be obtained by calling the [listwaitingroom](https://help.aliyun.com/document_detail/2850279.html) interface.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<site_id>:<waiting_room_id>:<waiting_room_rule_id>`.
* `waiting_room_rule_id` - The rule ID, which can be used to query a specific rule.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Waiting Room Rule.
* `delete` - (Defaults to 5 mins) Used when delete the Waiting Room Rule.
* `update` - (Defaults to 5 mins) Used when update the Waiting Room Rule.

## Import

ESA Waiting Room Rule can be imported using the id, e.g.

```shell
$ terraform import alicloud_esa_waiting_room_rule.example <site_id>:<waiting_room_id>:<waiting_room_rule_id>
```