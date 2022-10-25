---
subcategory: "Message Notification Service (MNS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_message_service_topics"
sidebar_current: "docs-alicloud-datasource-message-service-topics"
description: |-
  Provides a list of Message Notification Service Topics to the user.
---

# alicloud\_message\_service\_topics

This data source provides the Message Notification Service Topics of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.188.0+.

## Example Usage

```terraform
data "alicloud_message_service_topics" "ids" {
  ids = ["example_id"]
}

output "topic_id_1" {
  value = data.alicloud_message_service_topics.ids.topics.0.id
}

data "alicloud_message_service_topics" "name" {
  topic_name = "tf-example"
}

output "topic_id_2" {
  value = data.alicloud_message_service_topics.name.topics.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed) A list of Topic IDs. Its element value is same as Topic Name.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Topic name.
* `topic_name` - (Optional, ForceNew) The name of the topic.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Topic names.
* `topics` - A list of Topics. Each element contains the following attributes:
  * `id` - The id of the Topic. Its value is same as Topic Name.
  * `topic_name` - The name of the topic.
  * `message_count` - The number of messages in the topic.
  * `max_message_size` - The maximum size of a message body that can be sent to the topic. Unit: bytes.
  * `message_retention_period` - The maximum period for which a message can be retained in the topic. A message that is sent to the topic can be retained for a specified period. After the specified period ends, the message is deleted no matter whether it is pushed to the specified endpoints. Unit: seconds.
  * `logging_enabled` - Indicates whether the log management feature is enabled.
  * `topic_url` - The url of the topic.
  * `topic_inner_url` - The inner url of the topic.
  * `last_modify_time` - The time when the topic was last modified. This value is a UNIX timestamp representing the number of milliseconds that have elapsed since the epoch time January 1, 1970, 00:00:00 UTC.
  * `create_time` - The time when the topic was created. This value is a UNIX timestamp representing the number of milliseconds that have elapsed since the epoch time January 1, 1970, 00:00:00 UTC.