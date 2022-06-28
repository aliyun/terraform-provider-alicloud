---
subcategory: "Alikafka"
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

* `ids` - (Optional, ForceNew, Computed)  A list of ALIKAFKA Topics IDs, It is formatted to `<instance_id>:<topic>`.
* `instance_id` - (Required) ID of the instance.
* `name_regex` - (Optional) A regex string to filter results by the topic name.
* `topic` - (Optional) A topic to filter results by the topic name.  
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of topic names.
* `topics` - A list of topics. Each element contains the following attributes:
  * `id` - The ID of the topic, It is formatted to `<instance_id>:<topic>`.
  * `topic` - The name of the topic.
  * `create_time` - Time of creation.
  * `local_topic` - whether the current topic is kafka local topic or not.
  * `compact_topic` - whether the current topic is kafka compact topic or not.
  * `partition_num` - Partition number of the topic.
  * `remark` - Remark of the topic.
  * `status` - The current status code of the topic. There are three values to describe the topic status: 0 stands for the topic is in service, 1 stands for freezing and 2 stands for pause. 
  * `status_name` - The status_name of the topic.
  * `instance_id` - The instance_id of the instance.
  * `tags` - A mapping of tags to assign to the topic.