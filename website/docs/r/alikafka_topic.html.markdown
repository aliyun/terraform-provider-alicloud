---
subcategory: "AliKafka"
layout: "alicloud"
page_title: "Alicloud: alicloud_alikafka_topic"
description: |-
  Provides a Alicloud Alikafka Topic resource.
---

# alicloud_alikafka_topic

Provides a Alikafka Topic resource.

Topic in kafka.

For information about Alikafka Topic and how to use it, see [What is Topic](https://www.alibabacloud.com/help/en/message-queue-for-apache-kafka/latest/api-alikafka-2019-09-16-createtopic).

-> **NOTE:** Available since v1.56.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_alikafka_topic&exampleId=777edb69-3f55-dbe0-541e-4240ae1486ac9df26498&activeTab=example&spm=docs.r.alikafka_topic.0.777edb693f&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_security_group" "default" {
  vpc_id = alicloud_vpc.default.id
}

resource "alicloud_alikafka_instance" "default" {
  name            = var.name
  partition_num   = 50
  disk_type       = "1"
  disk_size       = "500"
  deploy_type     = "5"
  io_max          = "20"
  spec_type       = "professional"
  service_version = "2.2.0"
  vswitch_id      = alicloud_vswitch.default.id
  security_group  = alicloud_security_group.default.id
  config = jsonencode(
    {
      "enable.acl" : "true"
    }
  )
}

resource "alicloud_alikafka_topic" "default" {
  instance_id   = alicloud_alikafka_instance.default.id
  topic         = var.name
  remark        = var.name
  local_topic   = "true"
  compact_topic = "true"
  partition_num = "18"
  configs = jsonencode(
    {
      "retention.ms" : "3600000",
      "max.message.bytes" : "10485760",
      "message.timestamp.type" : "LogAppendTime"
    }
  )
  tags = {
    Created = "TF",
    For     = "example",
  }
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_alikafka_topic&spm=docs.r.alikafka_topic.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `compact_topic` - (Optional, ForceNew, Bool) The cleanup policy for the topic. This parameter is available only if you set the storage engine of the topic to Local storage. Valid values:
  - false: The delete cleanup policy is used.
  - true: The compact cleanup policy is used.
* `configs` - (Optional, Available since v1.262.1) The advanced configurations of the topic, as a JSON encoded object. The service carries every value as a string; a value declared as a number or a boolean is compared through its literal text, so it matches the string the service reports for it. The keys the service accepts and their value ranges are defined by [UpdateTopicConfig](https://www.alibabacloud.com/help/en/apsaramq-for-kafka/cloud-message-queue-for-kafka/developer-reference/api-alikafka-2019-09-16-updatetopicconfig) and depend on the instance:
  - Serverless instance: `retention.hours` (the retention period in hours, from `24` to `8760`), `max.message.bytes` (the maximum message size in bytes, from `1048576` to `10485760`), `message.timestamp.type` (`CreateTime` or `LogAppendTime`) and `message.timestamp.difference.max.ms` (the maximum difference allowed between the timestamp carried by a message and the time the server receives it - a message that exceeds it is rejected, and the key has no effect when `message.timestamp.type` is `LogAppendTime`).
  - Reserved instance: only a topic that uses the Local storage engine can be configured, and the supported keys are `retention.ms` (the retention period in milliseconds, from `3600000` to `31536000000`), `max.message.bytes`, `message.timestamp.type` and `message.timestamp.difference.max.ms`.
-> **NOTE:** The service reports the whole effective configuration of the topic, which can contain keys that were never declared in `configs` - a topic on a serverless instance reports `cloud.native.topic.type` and `replication-factor`, for example. Only the keys declared in `configs` are compared against the reported configuration, so those extra keys do not produce a diff.
-> **NOTE:** Removing a key from `configs` neither resets it on the topic nor produces a diff, because the update merges the submitted keys and offers no way to delete or reset a single key. Set the key to the value you want instead of removing it, or recreate the topic.
-> **NOTE:** A key that the service does not accept for the instance makes the update fail with `InvalidParameter.NotSupport`, so a configuration written for one instance type cannot be reused for the other as is - `retention.ms` and `retention.hours` in particular are not interchangeable.
-> **NOTE:** `replication-factor` is part of the configuration reported by the service, but it is not one of the keys accepted when the topic configuration is updated - the number of replicas is only taken when the topic is created. Declaring it in `configs` therefore cannot be used to make the reported configuration match the configuration.
* `instance_id` - (Required, ForceNew) The ID of the instance.
* `local_topic` - (Optional, ForceNew, Bool) The storage engine of the topic. Valid values:
  - false: Cloud storage.
  - true: Local storage.
* `partition_num` - (Optional, Int) The number of partitions in the topic.
* `remark` - (Required) The description of the topic.
* `tags` - (Optional, Map, Available since v1.63.0) A mapping of tags to assign to the resource.
* `topic` - (Required, ForceNew) The topic name.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<instance_id>:<topic>`.
* `create_time` - (Available since v1.262.1) The time when the topic was created.
* `region_id` - (Available since v1.262.1) The ID of the region where the instance resides.
* `status` - (Available since v1.262.1) The status of the service.

## Timeouts

-> **NOTE:** Available since v1.119.0.

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 15 mins) Used when create the Topic.
* `delete` - (Defaults to 16 mins) Used when delete the Topic.
* `update` - (Defaults to 5 mins) Used when update the Topic.

## Import

Alikafka Topic can be imported using the id, which consists of instance_id and topic, e.g.

```shell
$ terraform import alicloud_alikafka_topic.example <instance_id>:<topic>
```
