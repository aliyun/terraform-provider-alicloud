---
layout: "alicloud"
page_title: "Alicloud: alicloud_alikafka_topics"
sidebar_current: "docs-alicloud-datasource-alikafka-topics"
description: |-
    Provides a list of alikafka topics available to the user.
---

# alicloud\_alikafka\_topics

This data source provides a list of ALIKAFKA Topics in an Alibaba Cloud account according to the specified filters.

-> **NOTE:** Available in 1.56.0+

## Example Usage

```
data "alicloud_alikafka_topics" "topics_ds" {
  instance_id = "xxx"
  name_regex = "alikafkaTopicName"
  output_file = "topics.txt"
}

output "first_topic_name" {
  value = "${data.alicloud_alikafka_topics.topics_ds.topics.0.topic}"
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional) A regex string to filter results by the topic name. 
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of topic names.
* `topics` - A list of topics. Each element contains the following attributes:
  * `topic` - The name of the topic.
  * `create_time` - Time of creation.
  * `local_topic` - whether the current topic is kafka local topic or not.
  * `compact_topic` - whether the current topic is kafka compact topic or not.
  * `partition_num` - Partition number of the topic.
  * `remark` - Remark of the topic.
  * `status` - The current status code of the topic.
  * `status_name` - The current status description of the topic.