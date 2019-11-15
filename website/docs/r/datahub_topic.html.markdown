---
subcategory: "Datahub Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_datahub_topic"
sidebar_current: "docs-alicloud-resource-datahub-topic"
description: |-
  Provides a Alicloud datahub topic resource.
---

# alicloud\_datahub\_topic

The topic is the basic unit of Datahub data source and is used to define one kind of data or stream. It contains a set of subscriptions. You can manage the datahub source of an application by using topics. [Refer to details](https://help.aliyun.com/document_detail/47440.html).

## Example Usage

Basic Usage

- BLob Topic

  ```
resource "alicloud_datahub_topic" "example" {
  name         = "tf_datahub_topic"
  project_name = "tf_datahub_project"
  record_type  = "BLOB"
  shard_count  = 3
  life_cycle   = 7
  comment      = "created by terraform"
}
  ```
- Tuple Topic

  ```
resource "alicloud_datahub_topic" "example" {
  name         = "tf_datahub_topic"
  project_name = "tf_datahub_project"
  record_type  = "TUPLE"
  record_schema = {
    bigint_field    = "BIGINT"
    timestamp_field = "TIMESTAMP"
    string_field    = "STRING"
    double_field    = "DOUBLE"
    boolean_field   = "BOOLEAN"
  }
  shard_count = 3
  life_cycle  = 7
  comment     = "created by terraform"
}
  ```

## Argument Reference

The following arguments are supported:

* `name` - (Required, ForceNew) The name of the datahub topic. Its length is limited to 1-128 and only characters such as letters, digits and '_' are allowed. It is case-insensitive.
* `project_name` - (Required, ForceNew) The name of the datahub project that this topic belongs to. It is case-insensitive.
* `shard_count` - (Optional, ForceNew) The number of shards this topic contains. The permitted range of values is [1, 10]. The default value is 1.
* `life_cycle` - (Optional) How many days this topic lives. The permitted range of values is [1, 7]. The default value is 3.
* `record_type` - (Optional, ForceNew) The type of this topic. Its value must be one of {BLOB, TUPLE}. For BLOB topic, data will be organized as binary and encoded by BASE64. For TUPLE topic, data has fixed schema. The default value is "TUPLE" with a schema {STRING}.
* `record_schema` - (Optional, ForceNew) Schema of this topic, required only for TUPLE topic. Supported data types (case-insensitive) are:
  - BIGINT
  - STRING
  - BOOLEAN
  - DOUBLE
  - TIMESTAMP
* `comment` - (Optional) Comment of the datahub topic. It cannot be longer than 255 characters.

**Notes:** Currently `life_cycle` can not be modified and it will be supported in the next future.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the datahub topic. It was composed of project name and its name, and formats to `<project_name>:<name>`.
* `create_time` - Create time of the datahub topic. It is a human-readable string rather than 64-bits UTC.
* `last_modify_time` - Last modify time of the datahub topic. It is the same as *create_time* at the beginning. It is also a human-readable string rather than 64-bits UTC.

## Import

Datahub topic can be imported using the ID, e.g.

```
$ terraform import alicloud_datahub_topic.example tf_datahub_project:tf_datahub_topic
```
