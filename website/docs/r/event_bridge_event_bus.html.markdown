---
subcategory: "Event Bridge"
layout: "alicloud"
page_title: "Alicloud: alicloud_event_bridge_event_bus"
sidebar_current: "docs-alicloud-resource-event-bridge-event-bus"
description: |-
  Provides a Alicloud Event Bridge Event Bus resource.
---

# alicloud_event_bridge_event_bus

Provides a Event Bridge Event Bus resource.

For information about Event Bridge Event Bus and how to use it, see [What is Event Bus](https://www.alibabacloud.com/help/en/eventbridge/latest/api-eventbridge-2020-04-01-createeventbus).

-> **NOTE:** Available since v1.129.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_event_bridge_event_bus&exampleId=810d6db0-9dd7-6d83-0721-70eb01a7beb2e3309d7c&activeTab=example&spm=docs.r.event_bridge_event_bus.0.810d6db09d&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}
resource "alicloud_event_bridge_event_bus" "example" {
  event_bus_name = var.name
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_event_bridge_event_bus&spm=docs.r.event_bridge_event_bus.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `description` - (Optional) The description of event bus.
* `event_bus_name` - (Required, ForceNew) The name of event bus. The length is limited to 2 ~ 127 characters, which can be composed of letters, numbers or hyphens (-)

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Event Bus. Its value is same as `event_bus_name`.

## Import

Event Bridge Event Bus can be imported using the id, e.g.

```shell
$ terraform import alicloud_event_bridge_event_bus.example <event_bus_name>
```
