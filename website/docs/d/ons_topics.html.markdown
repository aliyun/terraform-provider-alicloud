---
subcategory: "RocketMQ"
layout: "alicloud"
page_title: "Alicloud: alicloud_ons_topics"
sidebar_current: "docs-alicloud-datasource-ons-topics"
description: |-
    Provides a list of ons topics available to the user.
---

# alicloud\_ons\_topics

This data source provides a list of ONS Topics in an Alibaba Cloud account according to the specified filters.

-> **NOTE:** Available in 1.53.0+

## Example Usage

```terraform
variable "name" {
  default = "onsInstanceName"
}

variable "topic" {
  default = "onsTopicDatasourceName"
}

resource "alicloud_ons_instance" "default" {
  instance_name = var.name
  remark        = "default_ons_instance_remark"
}

resource "alicloud_ons_topic" "default" {
  topic_name   = var.topic
  instance_id  = alicloud_ons_instance.default.id
  message_type = 0
  remark       = "dafault_ons_topic_remark"
}

data "alicloud_ons_topics" "topics_ds" {
  instance_id = alicloud_ons_topic.default.instance_id
  name_regex  = var.topic
  output_file = "topics.txt"
}

output "first_topic_name" {
  value = data.alicloud_ons_topics.topics_ds.topics.0.topic_name
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) ID of the ONS Instance that owns the topics.
* `ids` - (Optional, Available in v1.97.0+)  A list of topic IDs to filter results.
* `name_regex` - (Optional) A regex string to filter results by the topic name. 
* `tags` - (Optional, Available in v1.97.0+) A map of tags assigned to the Ons instance.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of topic names.
* `topics` - A list of topics. Each element contains the following attributes:
  * `id` - The id of the topic.
  * `topic` - The name of the topic.
  * `topic_name` - The name of the topic.
  * `owner` - The ID of the topic owner, which is the Alibaba Cloud UID.
  * `relation` - The relation ID. Read [Fields in PublishInfoDo](https://www.alibabacloud.com/help/doc-detail/29590.html) for further details.
  * `relation_name` - The name of the relation, for example, owner, publishable, subscribable, and publishable and subscribable.
  * `message_type` - The type of the message. Read [Fields in PublishInfoDo](https://www.alibabacloud.com/help/doc-detail/29590.html) for further details.
  * `independent_naming` - Indicates whether namespaces are available. Read [Fields in PublishInfoDo](https://www.alibabacloud.com/help/doc-detail/29590.html) for further details.
  * `create_time` - Time of creation.
  * `remark` - Remark of the topic.
  * `tags` - A map of tags assigned to the Ons instance.
  * `perm` - This attribute is used to set the read-write mode for the topic.
